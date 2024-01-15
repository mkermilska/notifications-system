package client

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mkermilska/notifications-sender/pkg/nats"
	"go.uber.org/zap"
)

type EmailAuthOpts struct {
	ApiKey string
	Domain string
	Sender string
}

type EmailService struct {
	mgClient *mailgun.MailgunImpl
	Sender   string
	logger   *zap.Logger
}

func NewEmailService(opts EmailAuthOpts, logger *zap.Logger) *EmailService {
	client := startEmailClient(opts)
	return &EmailService{
		mgClient: client,
		Sender:   opts.Sender,
		logger:   logger,
	}
}

func startEmailClient(opts EmailAuthOpts) *mailgun.MailgunImpl {
	mg := mailgun.NewMailgun(opts.Domain, opts.ApiKey)
	return mg
}

func (e *EmailService) SendEmail(email *nats.EmailNotification) error {
	message := e.mgClient.NewMessage(
		e.Sender,
		email.Subject,
		email.Message,
		email.To,
	)
	resp, id, err := e.mgClient.Send(context.TODO(), message)
	if err != nil {
		return fmt.Errorf("Error sending email: %w", err)
	}

	e.logger.Info("Email sent successfully.", zap.String("Response", resp), zap.String("ID", id))
	return nil
}
