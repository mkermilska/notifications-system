package nats

import "encoding/json"

type Event struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Time     int64             `json:"time"`
	Payload  json.RawMessage   `json:"payload"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type EmailNotification struct {
	To      string `json:"to"`
	Cc      string `json:"cc"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type SMSNotification struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type SlackNotification struct {
	ChannelID string `json:"channel_id"`
	Message   string `json:"message"`
}
