package db

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNotificationRepository_CreateNotification(t *testing.T) {
	mockNotification := &Notification{
		ID:       "842e8b6cef6005042216b0fddeff3cd3286c717e4ae66e60c6287eda18f05c95",
		Type:     "email",
		Message:  "test message",
		Receiver: "test-email@gmail.com",
		Status:   "failed",
	}

	nr := NewNotificationsRepository(db, zap.NewNop())
	err := nr.UpsertNotification(mockNotification)
	require.Nil(t, err, "Error creating notification")
}
