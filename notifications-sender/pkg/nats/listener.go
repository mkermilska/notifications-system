package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Listener struct {
	natsSvc Service
	logger  *zap.Logger
}

func NewListener(natsSvc Service, logger *zap.Logger) *Listener {
	listener := &Listener{
		natsSvc: natsSvc,
		logger:  logger,
	}
	return listener
}

func (l *Listener) Subscribe(eventName string, handler func(m *nats.Msg)) {
	inbox := nats.NewInbox()
	l.logger.Info("Listening for NATS event", zap.String("event", eventName))
	_, err := l.natsSvc.js.QueueSubscribe(
		eventName,
		fmt.Sprintf("%s-%s-queue", "als", eventName),
		handler,
		nats.Durable(fmt.Sprintf("%s-%s", "als", eventName)),
		nats.DeliverSubject(inbox),
		nats.DeliverNew(),
		nats.AckWait(5*time.Second),
	)
	if err != nil {
		l.logger.Error(fmt.Sprintf("Error with subscription for event %s, error: %s\n", eventName, err.Error()))
	}
}
