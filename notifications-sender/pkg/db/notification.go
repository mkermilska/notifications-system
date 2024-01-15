package db

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Notification struct {
	ID        string    `db:"id"`
	Type      string    `db:"type"`
	Message   string    `db:"message"`
	Receiver  string    `db:"receiver"`
	Status    string    `db:"status"`
	Retries   int       `db:"retries"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type NotificationsRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewNotificationsRepository(db *sqlx.DB, logger *zap.Logger) *NotificationsRepository {
	return &NotificationsRepository{db: db, logger: logger}
}

func (nr *NotificationsRepository) UpsertNotification(notification *Notification) error {
	nr.logger.Debug("Create notification: ", zap.String("notification_id", notification.ID))
	_, err := nr.db.NamedExec(
		`INSERT INTO notifications (id, type, message, receiver, status, retries)  
		VALUES (:id, :type, :message, :receiver, :status, 0)
		ON DUPLICATE KEY UPDATE status = :status, retries = retries + 1`, &notification)
	if err != nil {
		return fmt.Errorf("error persisting notification: %w", err)
	}
	return nil
}

func (nr *NotificationsRepository) GetNotifications(params map[string]interface{}) ([]Notification, error) {
	var query bytes.Buffer
	query.WriteString("SELECT * FROM notifications WHERE 1=1")
	args := make([]interface{}, 0)

	if status, exists := params["status"]; exists {
		query.WriteString(` AND status = ? `)
		args = append(args, status)
	}

	notifications := make([]Notification, 0)
	notificationsQuery, notificationsArgs, err := sqlx.In(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("error parsing notification query parameters: %w", err)
	}
	err = nr.db.Select(&notifications, nr.db.Rebind(notificationsQuery), notificationsArgs...)
	if err != nil {
		return nil, fmt.Errorf("error getting configured tenants: %w", err)
	}
	return notifications, nil
}
