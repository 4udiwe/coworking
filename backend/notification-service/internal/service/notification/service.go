package notification_service

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	notification_repository "github.com/4udiwe/coworking/notification-service/internal/repository/notification"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type NotificationService struct {
	notificationRepo NotificationRepository
	deviceRepo       DeviceRepository
	outboxRepo       OutboxRepository
	pushService      PushService

	txManager transactor.Transactor
}

func New(
	notificationRepo NotificationRepository,
	deviceRepo DeviceRepository,
	outboxRepo OutboxRepository,
	pushService PushService,
	txManager transactor.Transactor,
) *NotificationService {

	return &NotificationService{
		notificationRepo: notificationRepo,
		deviceRepo:       deviceRepo,
		outboxRepo:       outboxRepo,
		pushService:      pushService,
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
		return ErrCannotRegisterDevice
	}

	logrus.WithField("user_id", userID.String()).Info("device registered")

	return nil
}

func (s *NotificationService) NotifyUser(ctx context.Context, notificationID uuid.UUID) error {
	logrus.WithField("notification_id", notificationID).Info("notifying user")

	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, notification_repository.ErrNotificationNotFound) {
			return ErrNotificationNotFound
		}
		logrus.WithError(err).Error("failed to fetch notification")
		return ErrCannotFetchNotification
	}

	err = s.pushService.SendToUser(ctx, notification.UserID, notification)
	if err != nil {
		logrus.WithError(err).Error("failed to send push notification")
		return err
	}

	logrus.WithField("notification_id", notificationID).Info("user notified successfully")
	return nil
}

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
				"notificataionId":  id,
				"notificationType": notification.Type,
				"userId":           notification.UserID,
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

func (s *NotificationService) FetchNotifications(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int,
	isRead *bool,
) ([]entity.Notification, error) {
	logrus.WithFields(logrus.Fields{
		"user_id": userID.String(),
		"limit":   limit,
		"offset":  offset,
		"isRead":  isRead,
	}).Info("fetching notifications for user")

	notifications, err := s.notificationRepo.FetchByUser(ctx, userID, limit, offset, isRead)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch notifications")
		return nil, ErrCannotFetchNotification
	}

	logrus.WithField("number", len(notifications)).Info("notifications fetched")
	return notifications, nil
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	logrus.WithField("user_id", userID.String()).Info("marking all notifications as read")

	err := s.notificationRepo.MarkAllRead(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to mark all notifications as read")
		return ErrCannotMarkRead
	}

	logrus.WithField("user_id", userID.String()).Info("all notifications marked as read")
	return nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	logrus.WithField("user_id", userID.String()).Info("fetching unread notification count")

	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch unread count")
		return 0, ErrCannotFetchNotification
	}

	logrus.WithFields(logrus.Fields{
		"user_id": userID.String(),
		"count":   count,
	}).Info("unread count fetched")

	return count, nil
}

func (s *NotificationService) FetchNotificationsAfterDate(
	ctx context.Context,
	userID uuid.UUID,
	since time.Time,
	limit, offset int,
) ([]entity.Notification, error) {
	logrus.WithFields(logrus.Fields{
		"user_id": userID.String(),
		"since":   since,
		"limit":   limit,
		"offset":  offset,
	}).Info("fetching notifications after date")

	notifications, err := s.notificationRepo.FetchAfterDate(ctx, userID, since, limit, offset)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch notifications after date")
		return nil, ErrCannotFetchNotification
	}

	logrus.WithField("number", len(notifications)).Info("notifications after date fetched")
	return notifications, nil
}
