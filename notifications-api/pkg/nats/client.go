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
	return nil
}
