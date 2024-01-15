package service

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	database "github.com/mkermilska/notifications-sender/pkg/db"
)

type NotificationService struct {
	notificationsRepository *database.NotificationsRepository
	logger                  *zap.Logger
}

func NewNotificationService(db *sqlx.DB, logger *zap.Logger) *NotificationService {
	notificationsRepository := database.NewNotificationsRepository(db, logger)
	return &NotificationService{
		notificationsRepository: notificationsRepository,
		logger:                  logger,
	}
}

func (n NotificationService) UpsertNotification(notification *database.Notification) error {
	err := n.notificationsRepository.UpsertNotification(notification)
	if err != nil {
		return err
	}
	return nil
}

func (n NotificationService) GetNotifications(params map[string]interface{}) ([]database.Notification, error) {
	notifications, err := n.notificationsRepository.GetNotifications(params)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
