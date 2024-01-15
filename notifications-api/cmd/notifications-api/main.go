package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mkermilska/notifications-api/internal/web"
	"github.com/mkermilska/notifications-api/pkg/nats"
	"github.com/mkermilska/notifications-api/pkg/service"
)

const (
	serviceID = "notifications-api"
)

var cli struct {
	Debug      int    `kong:"short='d',env='DEBUG',default=0,help='Run in debug mode'"`
	HTTPPort   int    `kong:"short='t',env='HTTP_PORT',default='59191',help='HTTP server port'"`
	NATSHost   string `kong:"env='NATS_HOST',default='localhost',help='NATS host'"`
	NATSPort   int    `kong:"env='NATS_PORT',default='4222',help='NATS port'"`
	NATSStream string `kong:"env='NATS_STREAM',default='notifications-stream',help='NATS stream'"`
}

func main() {
	kong.Parse(&cli, kong.Name(serviceID), kong.Description("API receiving notifications"), kong.UsageOnError())
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
	logger.Info("starting notifications API")

	natsSvc := nats.NewService(cli.NATSStream, serviceID, logger)
	err = natsSvc.Connect(fmt.Sprintf("%s:%d", cli.NATSHost, cli.NATSPort))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	natsPublisher := nats.NewPublisher(*natsSvc, logger)

	notificationSvc := service.NewNotificationService(natsPublisher, logger)

	server := web.New(
		cli.HTTPPort,
		notificationSvc,
		logger)
	if err != nil {
		logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}

	server.Start()
}
