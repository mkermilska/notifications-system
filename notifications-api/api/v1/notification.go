package v1

type Notification struct {
	Message string `json:"message"`
	SMS     SMS    `json:"sms,omitempty"`
	Email   Email  `json:"email,omitempty"`
	Slack   Slack  `json:"slack,omitempty"`
}

type SMS struct {
	Phone string `json:"phone"`
}

type Email struct {
	To      string `json:"to"`
	Cc      string `json:"cc"`
	Subject string `json:"subject"`
}

type Slack struct {
	ChannelID string `json:"channel_id"`
}
