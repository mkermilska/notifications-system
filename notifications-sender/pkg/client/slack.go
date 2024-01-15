package client

import (
	"fmt"

	"github.com/mkermilska/notifications-sender/pkg/nats"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type SlackAuthOpts struct {
	Token string
}

type SlackService struct {
	slackClient *slack.Client
	logger      *zap.Logger
}

func NewSlackService(opts SlackAuthOpts, logger *zap.Logger) *SlackService {
	client := startSlackClient(opts)
	return &SlackService{
		slackClient: client,
		logger:      logger,
	}
}

func startSlackClient(opts SlackAuthOpts) *slack.Client {
	api := slack.New(opts.Token)
	return api
}

func (s *SlackService) SendSlackMessage(slackMsg *nats.SlackNotification) error {
	_, _, err := s.slackClient.PostMessage(slackMsg.ChannelID, slack.MsgOptionText(slackMsg.Message, false))
	if err != nil {
		return fmt.Errorf("Error sending slack message: %w", err)
	}
	return nil
}
