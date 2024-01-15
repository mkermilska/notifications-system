package service

import (
	"encoding/json"

	"github.com/gofrs/uuid"
	nats2 "github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/mkermilska/notifications-sender/pkg/client"
	"github.com/mkermilska/notifications-sender/pkg/db"
	"github.com/mkermilska/notifications-sender/pkg/nats"
)

type EventHandler struct {
	notificationSvc *NotificationService
	smsSvc          *client.SMSService
	emailSvc        *client.EmailService
	slackSvc        *client.SlackService
	logger          *zap.Logger
}

func NewEventHandler(notificationSvc *NotificationService, smsSvc *client.SMSService, emailSvc *client.EmailService,
	slackSvc *client.SlackService, logger *zap.Logger) *EventHandler {
	eventHandler := &EventHandler{
		notificationSvc: notificationSvc,
		smsSvc:          smsSvc,
		emailSvc:        emailSvc,
		slackSvc:        slackSvc,
		logger:          logger,
	}
	return eventHandler
}

func (h *EventHandler) EmailMessageEventHandler(m *nats2.Msg) {
	h.logger.Debug("Nats event received: EMAIL_MESSAGE")
	event := &nats.Event{}
	if err := json.Unmarshal(m.Data, event); err != nil {
		h.logger.Error("Error unmarshalling incoming event", zap.Error(err))
		return
	}
	email := nats.EmailNotification{}
	if err := json.Unmarshal(event.Payload, &email); err != nil {
		h.logger.Error("Error unmarshalling email payload", zap.Error(err))
		return
	}
	if err := m.Ack(); err != nil {
		h.logger.Error("Not able to acknowledge message was delivered", zap.Error(err))
		return
	}

	if err := h.emailSvc.SendEmail(&email); err != nil {
		if err := h.notificationSvc.UpsertNotification(&db.Notification{
			ID:       uuid.Must(uuid.NewV4()).String(),
			Type:     "email",
			Message:  email.Message,
			Receiver: email.To,
			Status:   "failed",
		}); err != nil {
			h.logger.Error("Not able to store failed email notification", zap.Error(err))
		}
	}
}

func (h *EventHandler) SMSMessageEventHandler(m *nats2.Msg) {
	h.logger.Debug("Nats event received: SMS_MESSAGE")
	event := &nats.Event{}
	if err := json.Unmarshal(m.Data, event); err != nil {
		h.logger.Error("Error unmarshalling incoming event", zap.Error(err))
		return
	}
	sms := nats.SMSNotification{}
	if err := json.Unmarshal(event.Payload, &sms); err != nil {
		h.logger.Error("Error unmarshalling sms payload", zap.Error(err))
		return
	}

	if err := m.Ack(); err != nil {
		h.logger.Error("Not able to acknowledge message was delivered", zap.Error(err))
		return
	}

	if err := h.smsSvc.SendSMS(&sms); err != nil {
		if err := h.notificationSvc.UpsertNotification(&db.Notification{
			ID:       uuid.Must(uuid.NewV4()).String(),
			Type:     "sms",
			Message:  sms.Message,
			Receiver: sms.Phone,
			Status:   "failed",
		}); err != nil {
			h.logger.Error("Not able to store failed sms notification", zap.Error(err))
		}
	}
}

func (h *EventHandler) SlackMessageEventHandler(m *nats2.Msg) {
	h.logger.Debug("Nats event received: SLACK_MESSAGE")
	event := &nats.Event{}
	if err := json.Unmarshal(m.Data, event); err != nil {
		h.logger.Error("Error unmarshalling incoming event", zap.Error(err))
		return
	}
	slack := nats.SlackNotification{}
	if err := json.Unmarshal(event.Payload, &slack); err != nil {
		h.logger.Error("Error unmarshalling slack payload", zap.Error(err))
		return
	}

	if err := m.Ack(); err != nil {
		h.logger.Error("Not able to acknowledge message was delivered", zap.Error(err))
		return
	}

	if err := h.slackSvc.SendSlackMessage(&slack); err != nil {
		if err := h.notificationSvc.UpsertNotification(&db.Notification{
			ID:       uuid.Must(uuid.NewV4()).String(),
			Type:     "slack",
			Message:  slack.Message,
			Receiver: slack.ChannelID,
			Status:   "failed",
		}); err != nil {
			h.logger.Error("Not able to store failed notification", zap.Error(err))
		}
	}
}
