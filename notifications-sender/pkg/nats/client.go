package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Service struct {
	streamName string
	sourceName string
	conn       *nats.Conn
	js         nats.JetStreamContext
	logger     *zap.Logger
}

func NewService(streamName, sourceName string, logger *zap.Logger) *Service {
	return &Service{
		streamName: streamName,
		logger:     logger,
	}
}

func (s *Service) Connect(natsAddress string) error {
	if s.streamName == "" {
		return fmt.Errorf("stream is not defined")
	}

	opts := []nats.Option{
		nats.Name(s.sourceName),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Minute),
		nats.Timeout(5 * time.Second),
		nats.Compression(true),
		nats.PingInterval(time.Minute),
		nats.MaxPingsOutstanding(nats.DefaultMaxPingOut),
	}

	var err error
	s.conn, err = nats.Connect(natsAddress, opts...)
	if err != nil {
		return fmt.Errorf("error connecting to NATS at %s: %w", natsAddress, err)
	}

	s.js, err = s.conn.JetStream()
	if err != nil {
		return fmt.Errorf("error creating stream context for nats server: %w", err)
	}
	s.logger.Info("Connected to NATS: ", zap.String("nats_address", natsAddress))

	//TODO: change that before delivering to production
	//this is needed because there is no easy way to initially create nats stream and subjects.
	//For purposes, this should be handled in different way, for example with Terraform
	streamConfig := &nats.StreamConfig{
		Name: "notifications-stream",
		Subjects: []string{
			"EMAIL_MESSAGE",
			"SMS_MESSAGE",
			"SLACK_MESSAGE",
		},
		MaxBytes: 1024,
	}

	_, err = s.js.StreamInfo(streamConfig.Name)
	if err == nil {
		return nil
	}
	_, err = s.js.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("error adding stream %w", err)
	}
	s.logger.Info("Adding stream to nats: ", zap.String("stream_name", streamConfig.Name))
	return nil
}
