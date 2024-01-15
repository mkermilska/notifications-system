package nats

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Publisher struct {
	natsSvc Service
	logger  *zap.Logger
}

func NewPublisher(natsSvc Service, logger *zap.Logger) *Publisher {
	publisher := &Publisher{
		natsSvc: natsSvc,
		logger:  logger,
	}
	return publisher
}

func (p *Publisher) PublishEvent(eventType string, payload interface{}) error {
	msgID := uuid.Must(uuid.NewV4()).String()
	event := Event{
		ID:      msgID,
		Type:    eventType,
		Time:    time.Now().Unix(),
		Payload: payload,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshalling nats event: %w", err)
	}

	natsMsg := nats.NewMsg(event.Type)
	natsMsg.Data = eventData
	pubOpts := []nats.PubOpt{nats.MsgId(msgID), nats.ExpectStream(p.natsSvc.streamName)}
	_, err = p.natsSvc.js.PublishMsg(natsMsg, pubOpts...)
	if err != nil {
		return fmt.Errorf("error publishing event to nats: %w", err)
	}
	p.logger.Debug("Event published to NATS: ", zap.ByteString("event_data", eventData))
	return nil
}
