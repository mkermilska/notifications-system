package client

import (
	"fmt"
	"net/http"

	"github.com/mkermilska/notifications-sender/pkg/nats"
	"github.com/nexmo-community/nexmo-go"
	"go.uber.org/zap"
)

type SMSAuthOpts struct {
	ApiKey    string
	ApiSecret string
}

type SMSService struct {
	client *nexmo.Client
	logger *zap.Logger
}

func NewSMSService(opts SMSAuthOpts, logger *zap.Logger) *SMSService {
	client := startSMSClient(opts)
	return &SMSService{
		client: client,
		logger: logger,
	}
}

func startSMSClient(opts SMSAuthOpts) *nexmo.Client {
	auth := nexmo.NewAuthSet()
	auth.SetAPISecret(opts.ApiKey, opts.ApiSecret)

	client := nexmo.NewClient(http.DefaultClient, auth)
	return client
}

func (s *SMSService) SendSMS(sms *nats.SMSNotification) error {
	smsReq := nexmo.SendSMSRequest{
		To:   sms.Phone,
		Text: sms.Message,
	}

	_, _, err := s.client.SMS.SendSMS(smsReq)
	if err != nil {
		return fmt.Errorf("Error sending sms: %w", err)
	}
	return nil
}
