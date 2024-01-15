package service

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"

	"go.uber.org/zap"

	apiv1 "github.com/mkermilska/notifications-api/api/v1"
	"github.com/mkermilska/notifications-api/pkg/nats"
)

const (
	natsEmailNotification = "EMAIL_MESSAGE"
	natsSMSNotification   = "SMS_MESSAGE"
	natsSlackNotification = "SLACK_MESSAGE"
	phoneNumberRegex      = `^\+?[0-9]+$`
)

type NotificationService struct {
	natsPublisher *nats.Publisher
	logger        *zap.Logger
}

func NewNotificationService(natsPublisher *nats.Publisher, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		natsPublisher: natsPublisher,
		logger:        logger,
	}
}

func (n *NotificationService) ValidateNotification(notification *apiv1.Notification) error {
	errMsg := ""
	if !reflect.ValueOf(notification.Email).IsZero() {
		emailErr := n.validateEmailBody(&notification.Email)
		if emailErr != nil {
			errMsg += fmt.Sprintln(emailErr.Error())
		}
	}
	if !reflect.ValueOf(notification.SMS).IsZero() {
		smsErr := n.validateSMSBody(&notification.SMS)
		if smsErr != nil {
			errMsg += fmt.Sprintln(smsErr.Error())
		}
	}
	if errMsg != "" {
		return fmt.Errorf("Body validation errors: %s", errMsg)
	}
	return nil
}

// basic validation of To email receiver, Cc and other fields needs to be added
func (n *NotificationService) validateEmailBody(email *apiv1.Email) error {
	_, err := mail.ParseAddress(email.To)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

// basic validation of the phone number - if it contains only digits, + symbol is allowed at the beginning
func (n *NotificationService) validateSMSBody(sms *apiv1.SMS) error {
	fmt.Println(sms.Phone)
	valid, err := regexp.MatchString(phoneNumberRegex, sms.Phone)
	if err != nil {
		return fmt.Errorf("regexp error while validating sms body: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid sms phone number")
	}
	return nil
}

func (n *NotificationService) PublishNotification(notification *apiv1.Notification) error {
	if !reflect.ValueOf(notification.Email).IsZero() {
		if err := n.natsPublisher.PublishEvent(natsEmailNotification, &nats.EmailNotification{
			To:      notification.Email.To,
			Cc:      notification.Email.Cc,
			Subject: notification.Email.Subject,
			Message: notification.Message,
		}); err != nil {
			return fmt.Errorf("Error publishing message to nats subject %s: %w", natsEmailNotification, err)
		}
	}
	if !reflect.ValueOf(notification.SMS).IsZero() {
		if err := n.natsPublisher.PublishEvent(natsSMSNotification, &nats.SMSNotification{
			Phone:   notification.SMS.Phone,
			Message: notification.Message,
		}); err != nil {
			return fmt.Errorf("Error publishing message to nats subject %s: %w", natsSMSNotification, err)
		}
	}
	if !reflect.ValueOf(notification.Slack).IsZero() {
		if err := n.natsPublisher.PublishEvent(natsSlackNotification, &nats.SlackNotification{
			Message: notification.Message,
		}); err != nil {
			return fmt.Errorf("Error publishing message to nats subject %s: %w", natsSlackNotification, err)
		}
	}
	return nil
}
