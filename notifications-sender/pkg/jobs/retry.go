package jobs

import (
	"sync"
	"time"

	"github.com/mkermilska/notifications-sender/pkg/client"
	"github.com/mkermilska/notifications-sender/pkg/service"
	"go.uber.org/zap"
)

const (
	statusFailed      = "failed"
	checkIntervalMins = 15
)

type NotificationsRetry struct {
	notificationSvc *service.NotificationService
	emailSvc        *client.EmailService
	smsSvc          *client.SMSService
	slackSvc        *client.SlackService
	checkInterval   time.Duration
	logger          *zap.Logger
	doneChan        chan struct{}
	wg              sync.WaitGroup
}

func NewNotificationsRetry(notificationSvc *service.NotificationService, emailSvc *client.EmailService,
	smsSvc *client.SMSService, slackSvc *client.SlackService, logger *zap.Logger) *NotificationsRetry {
	return &NotificationsRetry{
		notificationSvc: notificationSvc,
		smsSvc:          smsSvc,
		emailSvc:        emailSvc,
		slackSvc:        slackSvc,
		checkInterval:   checkIntervalMins * time.Minute,
		logger:          logger,
		doneChan:        make(chan struct{}),
	}
}

func (n *NotificationsRetry) Retry() {
	n.logger.Info("Starting retry over failed notifications")
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		ticker := time.NewTicker(n.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-n.doneChan:
				return
			case <-ticker.C:
				n.logger.Info("Tenants synchronizer checking for failed tenants configuration")
				n.RetryNotifications()
			}
		}
	}()
}

func (n *NotificationsRetry) RetryNotifications() {
	params := map[string]interface{}{
		"status": "failed",
	}
	_, err := n.notificationSvc.GetNotifications(params)
	if err != nil {
		n.logger.Error("Background cron job: Error getting failed notifications")
	}

	//TODO: loop through notifications and retry sending each one
	// for _, notification := range failedNotifications {

	// }
}
