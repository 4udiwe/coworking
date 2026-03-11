package notification_service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	notification_repository "github.com/4udiwe/coworking/notification-service/internal/repository/notification"
	"github.com/4udiwe/coworking/notification-service/internal/sender"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type NotificationService struct {
	notificationRepo NotificationRepository
	deviceRepo       DeviceRepository
	outboxRepo       OutboxRepository
	sender           PushSender

	txManager transactor.Transactor
}

func New(
	notificationRepo NotificationRepository,
	deviceRepo DeviceRepository,
	outboxRepo OutboxRepository,
	sender PushSender,
	txManager transactor.Transactor,
) *NotificationService {

	return &NotificationService{
		notificationRepo: notificationRepo,
		deviceRepo:       deviceRepo,
		outboxRepo:       outboxRepo,
		sender:           sender,
		txManager:        txManager,
	}
}

func (s *NotificationService) RegisterDevice(
	ctx context.Context,
	userID uuid.UUID,
	deviceToken string,
	platform string,
) error {

	logrus.WithFields(logrus.Fields{
		"user_id": userID.String(),
	}).Info("registering device")

	device := entity.UserDevice{
		UserID:      userID,
		DeviceToken: deviceToken,
		Platform:    platform,
	}

	_, err := s.deviceRepo.Create(ctx, device)
	if err != nil {
		logrus.WithError(err).Error("failed to create device")
		return ErrCannotRegisterDevice
	}

	logrus.WithField("user_id", userID.String()).Info("device registered")

	return nil
}

func (s NotificationService) NotifyUser(ctx context.Context, notificationID uuid.UUID) error {
	logrus.WithField("notification_id", notificationID).Info("notifying user")

	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, notification_repository.ErrNotificationNotFound) {
			return ErrNotificationNotFound
		}
		logrus.WithError(err).Error("failed to fetch notification")
		return ErrCannotFetchNotification
	}

	devices, err := s.deviceRepo.FindByUserID(ctx, notification.UserID)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch user devices")
	}

	for _, d := range devices {
		logrus.WithFields(logrus.Fields{
			"device":  d.Platform,
			"user_id": d.UserID,
			"title":   notification.Title,
		}).Debug("notifying user device")

		var payloadMap map[string]string
		json.Unmarshal(notification.Payload, &payloadMap)
		err = s.sender.Send(ctx, sender.PushMessage{
			Token: d.DeviceToken,
			Title: notification.Title,
			Body:  notification.Body,
			Data:  payloadMap,
		})

		if err != nil {
			if errors.Is(err, sender.ErrInvalidToken) {
				logrus.WithField("device", d.Platform).Warn("invalid token found")
				deleteErr := s.deviceRepo.DeleteByToken(ctx, d.DeviceToken)
				if deleteErr != nil {
					logrus.WithField("device_token", d.DeviceToken).Warn("failed to delete invalid token")
				}
				continue
			}
			logrus.WithError(err).WithField("device", d.UserID).Error("failed to notify user device")
		}
	}
	logrus.WithField("notification_id", notificationID).Info("notifying user")

	return nil
}

// func (s *NotificationService) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]entity.UserDevice, error) {
// 	logrus.WithField("user_id", userID.String()).Info("fetching user devices")

// 	devices, err := s.deviceRepo.FindByUserID(ctx, userID)
// 	if err != nil {
// 		logrus.WithError(err).Error("failed to fetch user devices")
// 	}

// 	logrus.WithField("devices_number", len(devices)).Info("devices fetched")
// 	return devices, nil
// }

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	notification entity.Notification,
) error {

	logrus.WithFields(logrus.Fields{
		"user_id": notification.UserID.String(),
		"type":    notification.Type,
	}).Info("creating notification")

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {

		// Create notification
		id, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {

			logrus.WithError(err).
				Error("failed to create notification")

			return err
		}

		// Create outbox event
		event := entity.OutboxEvent{
			AggregateType: "notification",
			AggregateID:   id,
			EventType:     "created",
			Payload: map[string]any{
				"notificataionId": id,
			},
		}

		err = s.outboxRepo.Create(ctx, event)
		if err != nil {

			logrus.WithError(err).
				Error("failed to create outbox event")

			return err
		}

		return nil
	})

	if err != nil {
		return ErrCannotCreateNotification
	}

	logrus.WithField("user_id", notification.UserID.String()).
		Info("notification created")

	return nil
}

func (s *NotificationService) MarkRead(
	ctx context.Context,
	notificationID uuid.UUID,
) error {

	logrus.WithField("notification_id", notificationID.String()).Info("marking notification as read")

	err := s.notificationRepo.MarkRead(ctx, notificationID)
	if err != nil {

		logrus.WithError(err).
			Error("failed to mark notification read")

		return ErrCannotMarkRead

	}

	logrus.WithField("notification_id", notificationID.String()).
		Info("notification marked read")

	return nil
}

func (s *NotificationService) FetchUnreadNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error) {
	logrus.WithField("user_id", userID.String()).Info("fetching unread notifications for user")

	notifications, err := s.notificationRepo.FetchUnreadByUser(ctx, userID, limit)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch unread notifications")
		return nil, ErrCannotFetchNotification
	}

	logrus.WithField("number", len(notifications)).Info("unread notifications fetched")

	return notifications, nil
}

// func (s *NotificationService) GetNotificationByID(
// 	ctx context.Context,
// 	notificationID uuid.UUID,
// ) (entity.Notification, error) {
// 	logrus.WithField("notification_id", notificationID.String()).Info("fetching notification")

// 	logrus.WithField("notification_id", notificationID.String()).Info("notification fetched")
// 	return notification, nil
// }
