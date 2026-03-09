package scheduler_service

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/cowoking/scheduler-service/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrCannotCreateTimer = errors.New("cannot create timer")
)

type SchedulerService struct {
	timerRepo TimerRepository
	txManager transactor.Transactor

	remindBefore time.Duration
}

func New(
	timerRepo TimerRepository,
	txManager transactor.Transactor,
	remindBefore time.Duration,
) *SchedulerService {
	return &SchedulerService{
		timerRepo:    timerRepo,
		txManager:    txManager,
		remindBefore: remindBefore,
	}
}

func (s *SchedulerService) HandleCreatedBooking(
	ctx context.Context,
	bookingID, userID, placeID uuid.UUID,
	startTime, endTime time.Time,
) error {

	logrus.Infof("Handling booking created: %s", bookingID)

	reminderTimer := entity.Timer{
		BookingID: bookingID,
		UserID:    &userID,

		Type: entity.TimerType{
			ID:   entity.TimerTypeBookingReminderID,
			Name: entity.TimerTypeBookingReminderName,
		},
		TriggerAt: startTime.Add(-s.remindBefore),
	}

	expireTimer := entity.Timer{
		BookingID: bookingID,
		UserID:    &userID,

		Type: entity.TimerType{
			ID:   entity.TimerTypeBookingExpireID,
			Name: entity.TimerTypeBoookingExpireName,
		},
		TriggerAt: endTime,
	}

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {

		_, err := s.timerRepo.Create(ctx, reminderTimer)
		if err != nil {
			logrus.Errorf("failed to create reminder timer: %v", err)
			return err
		}

		_, err = s.timerRepo.Create(ctx, expireTimer)
		if err != nil {
			logrus.Errorf("failed to create expire timer: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return ErrCannotCreateTimer
	}

	logrus.Infof("Timers created for booking %s", bookingID)

	return nil
}

func (s *SchedulerService) HandleCancelledBooking(
	ctx context.Context,
	bookingID uuid.UUID,
) error {

	logrus.Infof("Handling booking cancelled: %s", bookingID)

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {

		err := s.timerRepo.CancelByBooking(ctx, bookingID)
		if err != nil {
			logrus.Errorf("failed to cancel timers for booking %s: %v", bookingID, err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	logrus.Infof("Timers cancelled for booking %s", bookingID)

	return nil
}
