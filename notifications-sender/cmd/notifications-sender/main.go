package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mkermilska/notifications-sender/pkg/client"
	"github.com/mkermilska/notifications-sender/pkg/db"
	"github.com/mkermilska/notifications-sender/pkg/jobs"
	"github.com/mkermilska/notifications-sender/pkg/nats"
	"github.com/mkermilska/notifications-sender/pkg/service"
)

const (
	serviceID                    = "notifications-sender"
	natsSubjectEmailNotification = "EMAIL_MESSAGE"
	natsSubjectSMSNotification   = "SMS_MESSAGE"
	natsSubjectSlackNotification = "SLACK_MESSAGE"
)

var cli struct {
	Debug        int    `kong:"short='d',env='DEBUG',default=0,help='Run in debug mode'"`
	EmailEnabled bool   `kong:"short='e',env='EMAIL_ENABLED',default=true,help='Enable notifications via email channel'"`
	SMSEnabled   bool   `kong:"short='s',env='SMS_ENABLED',default=true,help='Enable notifications via sms channel'"`
	SlackEnabled bool   `kong:"short='m',env='SLACK_ENABLED',default=false,help='Enable notifications via slack channel'"`
	HTTPPort     int    `kong:"short='p',env='HTTP_PORT',default='57171',help='HTTP server port'"`
	NATSHost     string `kong:"short='h',env='NATS_HOST',default='localhost',help='NATS host'"`
	NATSPort     int    `kong:"short='t',env='NATS_PORT',default='4222',help='NATS port'"`
	NATSStream   string `kong:"short='r',env='NATS_STREAM',default='notifications-stream',help='NATS stream'"`
	EmailAPIKey  string `kong:"short='a',env='EMAIL_API_KEY',default='',help='Email client API key'"`
	EmailDomain  string `kong:"short='o',env='EMAIL_DOMAIN',default='',help='Email client domain'"`
	EmailSender  string `kong:"short='i',env='EMAIL_SENDER',default='postmaster@sandbox526f83b07a964bea9f8af7f5da7f3721.mailgun.org>', help='Email sender'"`
	SMSAPIKey    string `kong:"short='k',env='SMS_API_KEY',default='',help='SMS client API key'"`
	SMSAPISecret string `kong:"short='c',env='SMS_API_SECRET',default='',help='SMS client API secret'"`
	SlackToken   string `kong:"short='l',env='SLACK_TOKEN',default='',help='Slack client access token'"`
	DBHost       string `kong:"short='h',env='DB_HOST',default='127.0.0.1',help='DB server host'"`
	DBPort       int    `kong:"short='r',env='DB_PORT',default='5434',help='DB server port'"`
	DBName       string `kong:"short='n',env='DB_NAME',default='testingwithrentals',help='DB name'"`
	DBUsername   string `kong:"short='u',env='DB_USERNAME',default='root',help='DB username'"`
	DBPassword   string `kong:"short='p',env='DB_PASSWORD',default='root',help='DB password'"`
	//Disabled by default, functionality is not finalized
	RetryEnabled bool `kong:"short='y',env='RETRY_ENABLED',default=false,help='Enable sending messages retry'"`
}

func main() {
	kong.Parse(&cli, kong.Name(serviceID), kong.Description("Service sending notifications"), kong.UsageOnError())
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.TimeKey = "time"
	logCfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	if cli.Debug != 0 {
		logCfg.Level.SetLevel(zap.DebugLevel)
	}

	logger, err := logCfg.Build()
	if err != nil {
		logger.Fatal("Failed to initialize logger. Exiting " + err.Error())
	}
	logger.Info("starting notifications sender")

	db, err := db.StartDBStore(db.StartUpOptions{
		DBHost:         cli.DBHost,
		DBPort:         cli.DBPort,
		DBName:         cli.DBName,
		DBUsername:     cli.DBUsername,
		DBPassword:     cli.DBPassword,
		SkipMigrations: false,
	})
	if err != nil {
		logger.Fatal("Failed to start database", zap.Error(err))
	}

	notificationSvc := service.NewNotificationService(db, logger)

	var smsSvc *client.SMSService
	if cli.SMSEnabled {
		smsSvc = client.NewSMSService(client.SMSAuthOpts{
			ApiKey:    cli.SMSAPIKey,
			ApiSecret: cli.SMSAPISecret,
		}, logger)
	}

	var emailSvc *client.EmailService
	if cli.EmailEnabled {
		emailSvc = client.NewEmailService(client.EmailAuthOpts{
			ApiKey: cli.EmailAPIKey,
			Domain: cli.EmailDomain,
			Sender: cli.EmailSender,
		}, logger)
	}

	var slackSvc *client.SlackService
	if cli.SlackEnabled {
		slackSvc = client.NewSlackService(client.SlackAuthOpts{
			Token: cli.SlackToken,
		}, logger)
	}
	eventHandler := service.NewEventHandler(notificationSvc, smsSvc, emailSvc, slackSvc, logger)

	natsSvc := nats.NewService(cli.NATSStream, serviceID, logger)
	err = natsSvc.Connect(fmt.Sprintf("%s:%d", cli.NATSHost, cli.NATSPort))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	natsListener := nats.NewListener(*natsSvc, logger)
	go listenForEvents(natsListener, eventHandler)

	if cli.RetryEnabled {
		notificationsRetry := jobs.NewNotificationsRetry(notificationSvc, emailSvc, smsSvc, slackSvc, logger)
		notificationsRetry.Retry()
	}

	select {}
}

func listenForEvents(natsListener *nats.Listener, eventHandler *service.EventHandler) {
	if cli.EmailEnabled {
		natsListener.Subscribe(natsSubjectEmailNotification, eventHandler.EmailMessageEventHandler)
	}

	if cli.SMSEnabled {
		natsListener.Subscribe(natsSubjectSMSNotification, eventHandler.SMSMessageEventHandler)
	}

	if cli.SlackEnabled {
		natsListener.Subscribe(natsSubjectSlackNotification, eventHandler.SlackMessageEventHandler)
	}
}
