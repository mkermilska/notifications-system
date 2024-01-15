package service

import (
	"testing"

	"go.uber.org/zap"

	apiv1 "github.com/mkermilska/notifications-api/api/v1"
)

func TestNotificationService_ValidateNotification(t *testing.T) {
	notificationSvc := &NotificationService{
		natsPublisher: nil,
		logger:        zap.NewNop(),
	}

	tests := map[string]struct {
		input   *apiv1.Notification
		wantErr bool
	}{
		"Valid email": {
			input: &apiv1.Notification{
				Email: apiv1.Email{
					To: "kermilska@gmail.com",
				},
			},
			wantErr: false,
		},
		"Invalid Email - missing @": {
			input: &apiv1.Notification{
				Email: apiv1.Email{
					To: "kermilskagmail.com",
				},
			},
			wantErr: true,
		},
		"Valid SMS starting with +": {
			input: &apiv1.Notification{
				SMS: apiv1.SMS{
					Phone: "+35985525667",
				},
			},
			wantErr: false,
		},
		"Valid SMS without +": {
			input: &apiv1.Notification{
				SMS: apiv1.SMS{
					Phone: "085525667",
				},
			},
			wantErr: false,
		},
		"Invalid SMS": {
			input: &apiv1.Notification{
				SMS: apiv1.SMS{
					Phone: "08552k5667",
				},
			},
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := notificationSvc.ValidateNotification(test.input)
			if err != nil && !test.wantErr {
				t.Fatalf("Test case %s failed: %s", name, err)
			}
			if err == nil && test.wantErr {
				t.Fatalf("Test case %s failed: %s", name, err)
			}
		})
	}
}
