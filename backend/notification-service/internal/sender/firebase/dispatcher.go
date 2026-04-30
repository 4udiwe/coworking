package firebase_sender

import (
	"context"
	"encoding/json"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/4udiwe/coworking/notification-service/internal/sender"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// DefaultDispatcher implements the sender.Dispatcher interface using Firebase Cloud Messaging
type DefaultDispatcher struct {
	pushSender PushSender
	deviceRepo DeviceRepository
}

type PushSender interface {
	Send(ctx context.Context, msg sender.PushMessage) error
}

type DeviceRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserDevice, error)
	DeleteByToken(ctx context.Context, token string) error
}

func NewDefaultDispatcher(pushSender PushSender, deviceRepo DeviceRepository) *DefaultDispatcher {
	return &DefaultDispatcher{
		pushSender: pushSender,
		deviceRepo: deviceRepo,
	}
}

func (d *DefaultDispatcher) Dispatch(ctx context.Context, notification entity.Notification) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"user_id":         notification.UserID,
	}).Debug("dispatching notification to user devices")

	devices, err := d.deviceRepo.FindByUserID(ctx, notification.UserID)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch user devices for dispatch")
		return err
	}

	if len(devices) == 0 {
		logrus.WithField("user_id", notification.UserID).Debug("no devices found for user")
		return nil
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(notification.Payload, &rawMap); err != nil {
		logrus.WithField("payload", notification.Payload).WithError(err).Warn("failed to unmarshal notification payload")
		rawMap = make(map[string]interface{})
	}

	payloadMap := make(map[string]string)
	for k, v := range rawMap {
		switch val := v.(type) {
		case string:
			payloadMap[k] = val
		default:
			b, _ := json.Marshal(val)
			payloadMap[k] = string(b)
		}
	}

	for _, device := range devices {
		logrus.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"device_token":    device.DeviceToken,
		}).Debug("sending push to device")

		err := d.pushSender.Send(ctx, sender.PushMessage{
			Token:          device.DeviceToken,
			Title:          notification.Title,
			Body:           notification.Body,
			NotificationID: notification.ID.String(),
			ActionURL:      *notification.ActionURL,
			Data:           payloadMap,
		})

		if err != nil {
			if isInvalidToken(err) {
				logrus.WithField("device_token", device.DeviceToken).Warn("invalid token, removing device")
				deleteErr := d.deviceRepo.DeleteByToken(ctx, device.DeviceToken)
				if deleteErr != nil {
					logrus.WithError(deleteErr).Warn("failed to delete invalid token")
				}
				continue
			}
			logrus.WithError(err).WithField("device_token", device.DeviceToken).Error("failed to dispatch notification")
		}
	}

	return nil
}

func isInvalidToken(err error) bool {
	return err == sender.ErrInvalidToken
}
