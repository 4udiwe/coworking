package notification_service

import (
	"context"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/4udiwe/coworking/notification-service/internal/sender"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// DefaultPushService is the default implementation of PushService
type DefaultPushService struct {
	dispatcher sender.Dispatcher
}

func NewDefaultPushService(dispatcher sender.Dispatcher) *DefaultPushService {
	return &DefaultPushService{
		dispatcher: dispatcher,
	}
}

func (s *DefaultPushService) SendToUser(
	ctx context.Context,
	userID uuid.UUID,
	notification entity.Notification,
) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"user_id":         userID,
	}).Info("sending push notification to user")

	err := s.dispatcher.Dispatch(ctx, notification)
	if err != nil {
		logrus.WithError(err).Error("failed to dispatch notification")
		return err
	}

	return nil
}
