package booking_service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	layout_model "github.com/4udiwe/cowoking/booking-service/internal/layout/model"
	"github.com/4udiwe/cowoking/booking-service/internal/repository"
	"github.com/4udiwe/cowoking/booking-service/pkg/json_schema_validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const ExpectedLayoutFormatVersion = 1

type BookingService struct {
	bookingRepo     BookingRepository
	placeRepo       PlaceRepository
	coworkingRepo   CoworkingRepository
	outboxRepo      OutboxRepo
	layoutValidator json_schema_validator.Validator
	txManager       transactor.Transactor
}

func New(
	bookingRepo BookingRepository,
	placeRepo PlaceRepository,
	coworkingRepo CoworkingRepository,
	outboxRepo OutboxRepo,
	layoutValidator json_schema_validator.Validator,
	txManager transactor.Transactor,
) *BookingService {
	return &BookingService{
		bookingRepo:     bookingRepo,
		placeRepo:       placeRepo,
		coworkingRepo:   coworkingRepo,
		outboxRepo:      outboxRepo,
		layoutValidator: layoutValidator,
		txManager:       txManager,
	}
}

func (s *BookingService) CreateCoworking(ctx context.Context, coworking entity.Coworking) error {
	logrus.Infof("Creating coworking: %s", coworking.Name)

	err := s.coworkingRepo.Create(ctx, coworking)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrCoworkingAlreadyExists
		}
		logrus.Errorf("Failed to create coworking: %v", err)
		return ErrCannotCreateCoworking
	}

	return nil
}

func (s *BookingService) UpdateCoworking(ctx context.Context, coworking entity.Coworking) error {
	logrus.Infof("Updating coworking: %s", coworking.Name)

	err := s.coworkingRepo.Update(ctx, coworking)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to update coworking: %v", err)
		return ErrCannotUpdateCoworking
	}

	return nil
}

func (s *BookingService) GetCoworking(ctx context.Context, coworkingID uuid.UUID) (entity.Coworking, error) {
	logrus.Infof("Getting coworking with ID: %s", coworkingID.String())

	coworking, err := s.coworkingRepo.GetByID(ctx, coworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return entity.Coworking{}, ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get coworking: %v", err)
		return entity.Coworking{}, ErrCannotFetchCoworking
	}

	return coworking, nil
}

func (s *BookingService) ListCoworkings(ctx context.Context) ([]entity.Coworking, error) {
	logrus.Info("Listing coworkings")

	coworkings, err := s.coworkingRepo.List(ctx)
	if err != nil {
		logrus.Errorf("Failed to list coworkings: %v", err)
		return nil, ErrCannotFetchCoworking
	}

	return coworkings, nil
}

func (s *BookingService) CreateLayoutVersion(ctx context.Context, layout entity.CoworkingLayout) error {
	logrus.Infof("Creating layout version for coworking ID: %s", layout.CoworkingID)

	// Check if coworking exists
	_, err := s.coworkingRepo.GetByID(ctx, layout.CoworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get coworking by ID: %v", err)
		return ErrCannotCreateLayout
	}

	// Validate layout JSON schema
	if err := s.layoutValidator.Validate(layout.Layout); err != nil {
		logrus.Errorf("Failed to validate layout JSON schema: %v", err)
		return ErrInvalidLayoutSchema
	}

	var parsed layout_model.Layout

	if err := json.Unmarshal(layout.Layout, &parsed); err != nil {
		logrus.Errorf("Failed to unmarshal and validate layout JSON: %v", err)
		return ErrInvalidLayoutSchema
	}

	if parsed.FormatVersion != ExpectedLayoutFormatVersion {
		logrus.Errorf("Invalid layout format version: expected %d, got %d", ExpectedLayoutFormatVersion, parsed.FormatVersion)
		return ErrInvalidLayoutSchemaVersion
	}

	// Validate layout matches current coworking places
	// 1. Получаем все места coworking из БД
	places, err := s.placeRepo.GetByCoworking(ctx, layout.CoworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get places by coworking: %v", err)
		return ErrCannotCreateLayout
	}

	// 2. Формируем set мест из БД
	dbSet := make(map[string]struct{}, len(places))
	for _, p := range places {
		dbSet[p.ID.String()] = struct{}{}
	}

	// 3. Проверяем layout
	layoutSet := make(map[string]struct{}, len(parsed.Places))

	for _, p := range parsed.Places {

		// 3.1 Проверка что место существует в БД
		if _, ok := dbSet[p.ID]; !ok {
			logrus.Errorf("Layout contains unknown place ID: %s", p.ID)
			return ErrInvalidLayoutSchema
		}

		// 3.2 Проверка на дубликаты в layout
		if _, exists := layoutSet[p.ID]; exists {
			logrus.Errorf("Duplicate place ID in layout: %s", p.ID)
			return ErrInvalidLayoutSchema
		}

		layoutSet[p.ID] = struct{}{}

		// 3.3 Удаляем из dbSet проверенные места, чтобы в конце проверить что layout содержит все места из БД
		delete(dbSet, p.ID)
	}

	// 4. Проверка что layout содержит ВСЕ места
	if len(dbSet) > 0 {
		logrus.Errorf("Layout missing %d places from DB", len(dbSet))
		return ErrInvalidLayoutSchema
	}

	// Create new layout version
	err = s.coworkingRepo.CreateLayoutVersion(ctx, layout)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to create layout version: %v", err)
		return ErrCannotCreateLayout
	}

	return nil
}

func (s *BookingService) GetLatestLayout(ctx context.Context, coworkingID uuid.UUID) (entity.CoworkingLayout, error) {
	logrus.Infof("Getting latest layout for coworking ID: %s", coworkingID.String())

	layout, err := s.coworkingRepo.GetLatestLayout(ctx, coworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return entity.CoworkingLayout{}, ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get latest layout: %v", err)
		return entity.CoworkingLayout{}, ErrCannotFetchLayout
	}

	return layout, nil
}

func (s *BookingService) GetLayoutByVersion(ctx context.Context, coworkingID uuid.UUID, version int) (entity.CoworkingLayout, error) {
	logrus.Infof("Getting layout version %d for coworking ID: %s", version, coworkingID.String())

	layout, err := s.coworkingRepo.GetLayoutByVersion(ctx, coworkingID, version)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return entity.CoworkingLayout{}, ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get layout by version: %v", err)
		return entity.CoworkingLayout{}, ErrCannotFetchLayout
	}

	return layout, nil
}

func (s *BookingService) ListLayoutVersions(ctx context.Context, coworkingID uuid.UUID) ([]entity.CoworkingLayoutVersionTime, error) {
	logrus.Infof("Listing layout versions for coworking ID: %s", coworkingID.String())

	layoutVersions, err := s.coworkingRepo.ListLayoutVersions(ctx, coworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return nil, ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to list layout versions: %v", err)
		return nil, ErrCannotFetchLayout
	}

	return layoutVersions, nil
}

func (s *BookingService) CreatePlacesBatch(ctx context.Context, places []entity.Place) error {
	logrus.Infof("Creating batch of %d places for coworking ID: %s", len(places), places[0].Coworking.ID)

	err := s.placeRepo.CreateBatch(ctx, places)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrPlaceAlreadyExists
		}
		logrus.Errorf("Failed to create batch of places: %v", err)
		return ErrCannotCreatePlace
	}

	return nil
}

func (s *BookingService) SetPlaceActive(ctx context.Context, placeID uuid.UUID, active bool) error {
	logrus.Infof("Setting place ID %s active status to %t", placeID.String(), active)

	return s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		hasActiveBookings, err := s.placeRepo.CheckHasActiveBookings(ctx, placeID)
		if err != nil {
			return err
		}
		if !active && hasActiveBookings {
			return ErrPlaceHasActiveBookings
		}

		err = s.placeRepo.SetActive(ctx, placeID, active)
		if err != nil {
			if errors.Is(err, repository.ErrPlaceNotFound) {
				return ErrPlaceNotFound
			}
			logrus.Errorf("Failed to set place active status: %v", err)
			return ErrCannotUpdatePlace
		}

		return nil
	})
}

func (s *BookingService) GetPlacesByCoworking(ctx context.Context, coworkingID uuid.UUID) ([]entity.Place, error) {
	logrus.Infof("Getting places for coworking ID: %s", coworkingID.String())

	places, err := s.placeRepo.GetByCoworking(ctx, coworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return nil, ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get places by coworking: %v", err)
		return nil, ErrCannotFetchPlace
	}

	return places, nil
}

func (s *BookingService) GetPlaceByID(ctx context.Context, placeID uuid.UUID) (entity.Place, error) {
	logrus.Infof("Getting place with ID: %s", placeID.String())

	place, err := s.placeRepo.GetByID(ctx, placeID)
	if err != nil {
		if errors.Is(err, repository.ErrPlaceNotFound) {
			return entity.Place{}, ErrPlaceNotFound
		}
		logrus.Errorf("Failed to get place by ID: %v", err)
		return entity.Place{}, ErrCannotFetchPlace
	}

	return place, nil
}

func (s *BookingService) GetAvailablePlacesByCoworking(ctx context.Context, coworkingID uuid.UUID, start, end time.Time) ([]entity.Place, error) {
	logrus.Infof("Getting available places for coworking ID: %s between %s and %s", coworkingID.String(), start.Format(time.RFC3339), end.Format(time.RFC3339))

	places, err := s.placeRepo.GetAvailableByCoworking(ctx, coworkingID, start, end)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return nil, ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to get available places by coworking: %v", err)
		return nil, ErrCannotFetchPlace
	}

	return places, nil
}

func (s *BookingService) CreateBooking(ctx context.Context, booking entity.Booking) error {
	logrus.Infof("Creating booking for user ID: %s and place ID: %s", booking.UserID.String(), booking.Place.ID.String())

	if booking.StartTime.After(booking.EndTime) {
		return ErrBookingStartTimeAfterEndTime
	} else if booking.StartTime.Equal(booking.EndTime) {
		return ErrBookingStartTimeEqualEndTime
	} else if booking.StartTime.Before(time.Now().UTC()) {
		return ErrBookingStartTimeInPast
	} else if booking.StartTime.Minute() != 0 || booking.StartTime.Second() != 0 || booking.StartTime.Nanosecond() != 0 ||
		booking.EndTime.Minute() != 0 || booking.EndTime.Second() != 0 || booking.EndTime.Nanosecond() != 0 {
		return ErrBookingTimeNotMultipleOfHour
	} else if booking.EndTime.Sub(booking.StartTime) < time.Hour {
		return ErrBookingDurationLessThanOneHour
	} else if booking.EndTime.Sub(booking.StartTime) > 3*time.Hour {
		return ErrBookingDurationMoreThanThreeHours
	}

	booking.Status = entity.BookingStatusActive

	return s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check if place is active
		place, err := s.placeRepo.GetByID(ctx, booking.Place.ID)
		if err != nil {
			if errors.Is(err, repository.ErrPlaceNotFound) {
				return ErrPlaceNotFound
			}
			logrus.Errorf("Failed to get place by ID: %v", err)
			return ErrCannotCreateBooking
		}

		if !place.IsActive {
			return ErrPlaceInactive
		}

		// Check if coworking is active
		coworking, err := s.coworkingRepo.GetByID(ctx, place.Coworking.ID)
		if err != nil {
			if errors.Is(err, repository.ErrCoworkingNotFound) {
				return ErrCoworkingNotFound
			}
			logrus.Errorf("Failed to get coworking by ID: %v", err)
			return ErrCannotCreateBooking
		}

		if !coworking.IsActive {
			return ErrCoworkingInactive
		}

		// Create booking
		bookingID, err := s.bookingRepo.Create(ctx, booking)
		if err != nil {
			if errors.Is(err, repository.ErrBookingTimeConflict) {
				return ErrBookingTimeConflict
			}
			logrus.Errorf("Failed to create booking: %v", err)
			return ErrCannotCreateBooking
		}

		// Fetch booking, place, cowoking info
		booking, err := s.bookingRepo.GetByID(ctx, bookingID)
		if err != nil {
			return ErrCannotCreateBooking
		}

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "booking",
			AggregateID:   booking.ID,
			EventType:     "created",
			Payload: map[string]any{
				"bookingId": booking.ID,
				"userId":    booking.UserID,
				"placeId":   booking.Place.ID,
				"startTime": booking.StartTime,
				"endTime":   booking.EndTime,
			},
			Status:    entity.OutboxStatus{ID: 1, Name: "pending"},
			CreatedAt: time.Now(),
		}
		if err := s.outboxRepo.Create(ctx, ev); err != nil {
			logrus.Errorf("Failed to create outbox event: %v", err)
			return ErrCannotCreateBooking
		}

		return nil
	})
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID uuid.UUID, reason *string) error {
	logrus.Infof("Canceling booking with ID: %s", bookingID.String())

	return s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check if booking exists and is active
		booking, err := s.bookingRepo.GetByID(ctx, bookingID)
		if err != nil {
			if errors.Is(err, repository.ErrBookingNotFound) {
				return ErrBookingNotFound
			}
			logrus.Errorf("Failed to get booking by ID: %v", err)
			return ErrCannotCancelBooking
		}

		switch booking.Status {
		case entity.BookingStatusCancelled:
			return ErrBookingAlreadyCancelled
		case entity.BookingStatusCompleted:
			return ErrBookingAlreadyCompleted
		}

		// Cancel the booking
		err = s.bookingRepo.Cancel(ctx, bookingID, reason)

		if err != nil {
			logrus.Errorf("Failed to cancel booking: %v", err)
			return ErrCannotCancelBooking
		}

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "booking",
			AggregateID:   booking.ID,
			EventType:     "cancelled",
			Payload: map[string]any{
				"bookingId": booking.ID,
				"reason":    reason,
			},
			Status:    entity.OutboxStatus{ID: 1, Name: "pending"},
			CreatedAt: time.Now(),
		}
		if err := s.outboxRepo.Create(ctx, ev); err != nil {
			logrus.Errorf("Failed to create outbox event: %v", err)
			return ErrCannotCancelBooking
		}

		return nil
	})
}

func (s *BookingService) CompleteBooking(ctx context.Context, bookingID uuid.UUID) error {
	logrus.Infof("Completing booking with ID: %s", bookingID.String())

	return s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check if booking exists and is active
		booking, err := s.bookingRepo.GetByID(ctx, bookingID)

		if err != nil {
			if errors.Is(err, repository.ErrBookingNotFound) {
				return ErrBookingNotFound
			}
			logrus.Errorf("Failed to get booking by ID: %v", err)
			return ErrCannotCompleteBooking
		}

		switch booking.Status {
		case entity.BookingStatusCancelled:
			return ErrBookingAlreadyCancelled
		case entity.BookingStatusCompleted:
			return ErrBookingAlreadyCompleted
		}

		// Mark the booking as completed
		err = s.bookingRepo.MarkCompleted(ctx, bookingID)
		if err != nil {
			logrus.Errorf("Failed to complete booking: %v", err)
			return ErrCannotCompleteBooking
		}

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "booking",
			AggregateID:   booking.ID,
			EventType:     "completed",
			Payload: map[string]any{
				"bookingId": booking.ID,
			},
			Status:    entity.OutboxStatus{ID: 1, Name: "pending"},
			CreatedAt: time.Now(),
		}
		if err := s.outboxRepo.Create(ctx, ev); err != nil {
			logrus.Errorf("Failed to create outbox event: %v", err)
			return ErrCannotCancelBooking
		}

		return nil
	})
}

func (s *BookingService) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (entity.Booking, error) {
	logrus.Infof("Getting booking with ID: %s", bookingID.String())

	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			return entity.Booking{}, ErrBookingNotFound
		}
		logrus.Errorf("Failed to get booking by ID: %v", err)
		return entity.Booking{}, ErrCannotFetchBooking
	}

	return booking, nil
}

func (s *BookingService) ListBookingsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {
	logrus.Infof("Listing bookings for user ID: %s", userID.String())

	bookings, err := s.bookingRepo.ListByUser(ctx, userID)
	if err != nil {
		logrus.Errorf("Failed to list bookings by user: %v", err)
		return nil, ErrCannotFetchBooking
	}

	return bookings, nil
}

func (s *BookingService) SetCoworkingActive(ctx context.Context, coworkingID uuid.UUID) error {
	logrus.Infof("Setting coworking ID %s active status to %t", coworkingID.String(), true)

	err := s.coworkingRepo.SetActive(ctx, coworkingID, true)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to set coworking active status: %v", err)
		return ErrCannotUpdateCoworking
	}

	return nil
}

func (s *BookingService) SetCoworkingInactive(ctx context.Context, coworkingID uuid.UUID) error {
	logrus.Infof("Setting coworking ID %s active status to %t", coworkingID.String(), false)

	return s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		hasActiveBookings, err := s.coworkingRepo.CheckHasActiveBookings(ctx, coworkingID)
		if err != nil {
			return err
		}
		if hasActiveBookings {
			return ErrCoworkingHasActiveBookings
		}

		err = s.coworkingRepo.SetActive(ctx, coworkingID, false)
		if err != nil {
			return ErrCannotUpdateCoworking
		}

		return nil
	})
}

func (s *BookingService) RollbackLatestLayoutVersion(ctx context.Context, coworkingID uuid.UUID) error {
	logrus.Infof("Rolling back latest layout version for coworking ID: %s", coworkingID.String())

	err := s.coworkingRepo.RollbackLatestLayoutVersion(ctx, coworkingID)
	if err != nil {
		if errors.Is(err, repository.ErrCoworkingNotFound) {
			return ErrCoworkingNotFound
		}
		logrus.Errorf("Failed to rollback latest layout version: %v", err)
		return ErrCannotUpdateCoworking
	}

	return nil
}
