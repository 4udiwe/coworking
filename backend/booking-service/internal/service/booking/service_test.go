package booking_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/4udiwe/cowoking/booking-service/internal/repository"
	"github.com/4udiwe/cowoking/booking-service/internal/service/booking/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// ============================================================================
// HELPERS
// ============================================================================

type transactionTracker struct {
	withInTransactionCalled bool
}

func (tt *transactionTracker) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tt.withInTransactionCalled = true
	return fn(ctx)
}

type dummyTransactor struct{}

func (d dummyTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func stringPtr(s string) *string {
	return &s
}

// ============================================================================
// TESTS: CreateBooking Validations
// ============================================================================

func TestCreateBooking_Validations(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Hour)
	validStart := now.Add(2 * time.Hour)
	validEnd := now.Add(3 * time.Hour)

	tests := []struct {
		name      string
		start     time.Time
		end       time.Time
		wantError error
		desc      string
	}{
		{
			name:      "start_after_end",
			start:     validEnd,
			end:       validStart,
			wantError: ErrBookingStartTimeAfterEndTime,
			desc:      "Стартовое время позже времени окончания",
		},
		{
			name:      "start_equal_end",
			start:     validStart,
			end:       validStart,
			wantError: ErrBookingStartTimeEqualEndTime,
			desc:      "Стартовое время равно времени окончания",
		},
		{
			name:      "start_in_past",
			start:     now.Add(-1 * time.Hour),
			end:       now.Add(1 * time.Hour),
			wantError: ErrBookingStartTimeInPast,
			desc:      "Стартовое время в прошлом",
		},
		{
			name:      "start_with_minutes",
			start:     validStart.Add(15 * time.Minute),
			end:       validEnd,
			wantError: ErrBookingTimeNotMultipleOfHour,
			desc:      "Стартовое время не кратно часу (содержит минуты)",
		},
		{
			name:      "end_with_seconds",
			start:     validStart,
			end:       validEnd.Add(30 * time.Second),
			wantError: ErrBookingTimeNotMultipleOfHour,
			desc:      "Время окончания не кратно часу (содержит секунды)",
		},
		{
			name:      "duration_less_than_one_hour",
			start:     validStart,
			end:       validStart.Add(30 * time.Minute),
			wantError: ErrBookingTimeNotMultipleOfHour,
			desc:      "Длительность бронирования и время не кратны часу",
		},
		{
			name:      "duration_more_than_three_hours",
			start:     validStart,
			end:       validStart.Add(4 * time.Hour),
			wantError: ErrBookingDurationMoreThanThreeHours,
			desc:      "Длительность бронирования более трех часов",
		},
	}

	svc := &BookingService{
		txManager: dummyTransactor{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateBooking(context.Background(), entity.Booking{
				StartTime: tt.start,
				EndTime:   tt.end,
			})
			if !errors.Is(err, tt.wantError) {
				t.Errorf("CreateBooking() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}
		})
	}
}

// ============================================================================
// TESTS: CreateBooking Business Logic
// ============================================================================

func TestCreateBooking_BusinessLogic(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Hour)
	coworkingID := uuid.New()
	placeID := uuid.New()
	userID := uuid.New()
	bookingID := uuid.New()

	activeCoworking := entity.Coworking{ID: coworkingID, IsActive: true}
	inactiveCoworking := entity.Coworking{ID: coworkingID, IsActive: false}
	activePlace := entity.Place{ID: placeID, IsActive: true, Coworking: activeCoworking}
	inactivePlace := entity.Place{ID: placeID, IsActive: false, Coworking: activeCoworking}

	tests := []struct {
		name      string
		setup     func(*mocks.MockBookingRepository, *mocks.MockPlaceRepository, *mocks.MockOutboxRepo)
		place     entity.Place
		wantError error
		wantInTx  bool
		desc      string
	}{
		{
			name: "successful_booking_creation",
			setup: func(br *mocks.MockBookingRepository, pr *mocks.MockPlaceRepository, or *mocks.MockOutboxRepo) {
				pr.EXPECT().GetByID(gomock.Any(), placeID).Return(activePlace, nil)
				br.EXPECT().Create(gomock.Any(), gomock.Any()).Return(bookingID, nil)
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(entity.Booking{
					ID:        bookingID,
					UserID:    userID,
					Place:     activePlace,
					StartTime: now.Add(2 * time.Hour),
					EndTime:   now.Add(3 * time.Hour),
					Status:    entity.BookingStatusActive,
				}, nil)
				or.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			place:     activePlace,
			wantError: nil,
			wantInTx:  true,
			desc:      "Успешное создание бронирования с активным местом и коворкингом",
		},
		{
			name: "place_not_found",
			setup: func(br *mocks.MockBookingRepository, pr *mocks.MockPlaceRepository, or *mocks.MockOutboxRepo) {
				pr.EXPECT().GetByID(gomock.Any(), placeID).Return(entity.Place{}, repository.ErrPlaceNotFound)
			},
			place:     activePlace,
			wantError: ErrPlaceNotFound,
			wantInTx:  true,
			desc:      "Место не найдено в базе данных",
		},
		{
			name: "inactive_place",
			setup: func(br *mocks.MockBookingRepository, pr *mocks.MockPlaceRepository, or *mocks.MockOutboxRepo) {
				pr.EXPECT().GetByID(gomock.Any(), placeID).Return(inactivePlace, nil)
			},
			place:     inactivePlace,
			wantError: ErrPlaceInactive,
			wantInTx:  true,
			desc:      "Место неактивно",
		},
		{
			name: "inactive_coworking",
			setup: func(br *mocks.MockBookingRepository, pr *mocks.MockPlaceRepository, or *mocks.MockOutboxRepo) {
				pr.EXPECT().GetByID(gomock.Any(), placeID).Return(entity.Place{
					ID:        placeID,
					IsActive:  true,
					Coworking: inactiveCoworking,
				}, nil)
			},
			place:     activePlace,
			wantError: ErrCoworkingInactive,
			wantInTx:  true,
			desc:      "Коворкинг неактивен",
		},
		{
			name: "booking_time_conflict",
			setup: func(br *mocks.MockBookingRepository, pr *mocks.MockPlaceRepository, or *mocks.MockOutboxRepo) {
				pr.EXPECT().GetByID(gomock.Any(), placeID).Return(activePlace, nil)
				br.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.Nil, repository.ErrBookingTimeConflict)
			},
			place:     activePlace,
			wantError: ErrBookingTimeConflict,
			wantInTx:  true,
			desc:      "Конфликт времени бронирования",
		},
		{
			name: "place_repository_error",
			setup: func(br *mocks.MockBookingRepository, pr *mocks.MockPlaceRepository, or *mocks.MockOutboxRepo) {
				pr.EXPECT().GetByID(gomock.Any(), placeID).Return(entity.Place{}, errors.New("database error"))
			},
			place:     activePlace,
			wantError: ErrCannotCreateBooking,
			wantInTx:  true,
			desc:      "Ошибка репозитория при получении места",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBooking := mocks.NewMockBookingRepository(ctrl)
			mockPlace := mocks.NewMockPlaceRepository(ctrl)
			mockOutbox := mocks.NewMockOutboxRepo(ctrl)
			txTracker := &transactionTracker{}

			tt.setup(mockBooking, mockPlace, mockOutbox)

			svc := &BookingService{
				bookingRepo: mockBooking,
				placeRepo:   mockPlace,
				outboxRepo:  mockOutbox,
				txManager:   txTracker,
			}

			booking := entity.Booking{
				UserID:    userID,
				Place:     tt.place,
				StartTime: now.Add(2 * time.Hour),
				EndTime:   now.Add(3 * time.Hour),
			}

			err := svc.CreateBooking(context.Background(), booking)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("CreateBooking() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}

			if tt.wantInTx && !txTracker.withInTransactionCalled {
				t.Errorf("Expected transaction to be called | %s", tt.desc)
			}
		})
	}
}

// ============================================================================
// TESTS: CancelBooking
// ============================================================================

func TestCancelBooking(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Hour)
	bookingID := uuid.New()
	coworkingID := uuid.New()
	placeID := uuid.New()
	userID := uuid.New()

	activeBooking := entity.Booking{
		ID:        bookingID,
		UserID:    userID,
		Status:    entity.BookingStatusActive,
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(3 * time.Hour),
		Place: entity.Place{
			ID: placeID,
			Coworking: entity.Coworking{
				ID: coworkingID,
			},
		},
	}

	tests := []struct {
		name      string
		reason    *string
		setup     func(*mocks.MockBookingRepository, *mocks.MockOutboxRepo)
		wantError error
		wantInTx  bool
		desc      string
	}{
		{
			name:   "successful_cancellation",
			reason: nil,
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(activeBooking, nil)
				br.EXPECT().Cancel(gomock.Any(), bookingID, nil).Return(nil)
				or.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Успешная отмена активного бронирования",
		},
		{
			name:   "booking_not_found",
			reason: nil,
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(entity.Booking{}, repository.ErrBookingNotFound)
			},
			wantError: ErrBookingNotFound,
			wantInTx:  true,
			desc:      "Бронирование не найдено",
		},
		{
			name:   "booking_already_cancelled",
			reason: nil,
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				cancelledBooking := activeBooking
				cancelledBooking.Status = entity.BookingStatusCancelled
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(cancelledBooking, nil)
			},
			wantError: ErrBookingAlreadyCancelled,
			wantInTx:  true,
			desc:      "Бронирование уже отменено",
		},
		{
			name:   "booking_already_completed",
			reason: nil,
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				completedBooking := activeBooking
				completedBooking.Status = entity.BookingStatusCompleted
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(completedBooking, nil)
			},
			wantError: ErrBookingAlreadyCompleted,
			wantInTx:  true,
			desc:      "Бронирование уже завершено",
		},
		{
			name:   "cancel_with_reason",
			reason: stringPtr("User requested cancellation"),
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(activeBooking, nil)
				br.EXPECT().Cancel(gomock.Any(), bookingID, stringPtr("User requested cancellation")).Return(nil)
				or.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Отмена бронирования с указанной причиной",
		},
		{
			name:   "repository_error",
			reason: nil,
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(activeBooking, nil)
				br.EXPECT().Cancel(gomock.Any(), bookingID, nil).Return(errors.New("database error"))
			},
			wantError: ErrCannotCancelBooking,
			wantInTx:  true,
			desc:      "Ошибка репозитория при отмене",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBooking := mocks.NewMockBookingRepository(ctrl)
			mockOutbox := mocks.NewMockOutboxRepo(ctrl)
			txTracker := &transactionTracker{}

			tt.setup(mockBooking, mockOutbox)

			svc := &BookingService{
				bookingRepo: mockBooking,
				outboxRepo:  mockOutbox,
				txManager:   txTracker,
			}

			err := svc.CancelBooking(context.Background(), bookingID, tt.reason)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("CancelBooking() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}

			if tt.wantInTx && !txTracker.withInTransactionCalled {
				t.Errorf("Expected transaction to be called | %s", tt.desc)
			}
		})
	}
}

// ============================================================================
// TESTS: CompleteBooking
// ============================================================================

func TestCompleteBooking(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Hour)
	bookingID := uuid.New()
	coworkingID := uuid.New()
	placeID := uuid.New()
	userID := uuid.New()

	activeBooking := entity.Booking{
		ID:        bookingID,
		UserID:    userID,
		Status:    entity.BookingStatusActive,
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(3 * time.Hour),
		Place: entity.Place{
			ID: placeID,
			Coworking: entity.Coworking{
				ID: coworkingID,
			},
		},
	}

	tests := []struct {
		name      string
		setup     func(*mocks.MockBookingRepository, *mocks.MockOutboxRepo)
		wantError error
		wantInTx  bool
		desc      string
	}{
		{
			name: "successful_completion",
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(activeBooking, nil)
				br.EXPECT().MarkCompleted(gomock.Any(), bookingID).Return(nil)
				or.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Успешное завершение активного бронирования",
		},
		{
			name: "booking_not_found",
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(entity.Booking{}, repository.ErrBookingNotFound)
			},
			wantError: ErrBookingNotFound,
			wantInTx:  true,
			desc:      "Бронирование не найдено",
		},
		{
			name: "booking_already_cancelled",
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				cancelledBooking := activeBooking
				cancelledBooking.Status = entity.BookingStatusCancelled
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(cancelledBooking, nil)
			},
			wantError: ErrBookingAlreadyCancelled,
			wantInTx:  true,
			desc:      "Бронирование уже отменено",
		},
		{
			name: "booking_already_completed",
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				completedBooking := activeBooking
				completedBooking.Status = entity.BookingStatusCompleted
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(completedBooking, nil)
			},
			wantError: ErrBookingAlreadyCompleted,
			wantInTx:  true,
			desc:      "Бронирование уже завершено",
		},
		{
			name: "mark_completed_error",
			setup: func(br *mocks.MockBookingRepository, or *mocks.MockOutboxRepo) {
				br.EXPECT().GetByID(gomock.Any(), bookingID).Return(activeBooking, nil)
				br.EXPECT().MarkCompleted(gomock.Any(), bookingID).Return(errors.New("database error"))
			},
			wantError: ErrCannotCompleteBooking,
			wantInTx:  true,
			desc:      "Ошибка при отметке бронирования как завершенного",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBooking := mocks.NewMockBookingRepository(ctrl)
			mockOutbox := mocks.NewMockOutboxRepo(ctrl)
			txTracker := &transactionTracker{}

			tt.setup(mockBooking, mockOutbox)

			svc := &BookingService{
				bookingRepo: mockBooking,
				outboxRepo:  mockOutbox,
				txManager:   txTracker,
			}

			err := svc.CompleteBooking(context.Background(), bookingID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("CompleteBooking() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}

			if tt.wantInTx && !txTracker.withInTransactionCalled {
				t.Errorf("Expected transaction to be called | %s", tt.desc)
			}
		})
	}
}

// ============================================================================
// TESTS: SetPlaceActive
// ============================================================================

func TestSetPlaceActive(t *testing.T) {
	placeID := uuid.New()

	tests := []struct {
		name      string
		active    bool
		setup     func(*mocks.MockPlaceRepository)
		wantError error
		wantInTx  bool
		desc      string
	}{
		{
			name:   "activate_place",
			active: true,
			setup: func(pr *mocks.MockPlaceRepository) {
				pr.EXPECT().CheckHasActiveBookings(gomock.Any(), placeID).Return(false, nil)
				pr.EXPECT().SetActive(gomock.Any(), placeID, true).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Активация места без проверок",
		},
		{
			name:   "deactivate_place_without_active_bookings",
			active: false,
			setup: func(pr *mocks.MockPlaceRepository) {
				pr.EXPECT().CheckHasActiveBookings(gomock.Any(), placeID).Return(false, nil)
				pr.EXPECT().SetActive(gomock.Any(), placeID, false).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Деактивация места без активных бронирований",
		},
		{
			name:   "deactivate_place_with_active_bookings",
			active: false,
			setup: func(pr *mocks.MockPlaceRepository) {
				pr.EXPECT().CheckHasActiveBookings(gomock.Any(), placeID).Return(true, nil)
			},
			wantError: ErrPlaceHasActiveBookings,
			wantInTx:  true,
			desc:      "Попытка деактивации места с активными бронированиями",
		},
		{
			name:   "place_not_found",
			active: true,
			setup: func(pr *mocks.MockPlaceRepository) {
				pr.EXPECT().CheckHasActiveBookings(gomock.Any(), placeID).Return(false, nil)
				pr.EXPECT().SetActive(gomock.Any(), placeID, true).Return(repository.ErrPlaceNotFound)
			},
			wantError: ErrPlaceNotFound,
			wantInTx:  true,
			desc:      "Место не найдено",
		},
		{
			name:   "check_bookings_error",
			active: false,
			setup: func(pr *mocks.MockPlaceRepository) {
				pr.EXPECT().CheckHasActiveBookings(gomock.Any(), placeID).Return(false, errors.New("database error"))
			},
			wantError: errors.New("database error"),
			wantInTx:  true,
			desc:      "Ошибка при проверке активных бронирований",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPlace := mocks.NewMockPlaceRepository(ctrl)
			txTracker := &transactionTracker{}

			tt.setup(mockPlace)

			svc := &BookingService{
				placeRepo: mockPlace,
				txManager: txTracker,
			}

			err := svc.SetPlaceActive(context.Background(), placeID, tt.active)
			if tt.wantError != nil && (err == nil || err.Error() != tt.wantError.Error()) {
				t.Errorf("SetPlaceActive() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			} else if (err == nil) != (tt.wantError == nil) {
				t.Errorf("SetPlaceActive() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}

			if tt.wantInTx && !txTracker.withInTransactionCalled {
				t.Errorf("Expected transaction to be called | %s", tt.desc)
			}
		})
	}
}

// ============================================================================
// TESTS: SetCoworkingActive
// ============================================================================

func TestSetCoworkingActive(t *testing.T) {
	coworkingID := uuid.New()

	tests := []struct {
		name      string
		isActive  bool
		setup     func(*mocks.MockCoworkingRepository)
		wantError error
		wantInTx  bool
		desc      string
	}{
		{
			name:     "activate_coworking",
			isActive: true,
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().SetActive(gomock.Any(), coworkingID, true).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Активация коворкинга",
		},
		{
			name:     "deactivate_coworking_without_active_bookings",
			isActive: false,
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().CheckHasActiveBookings(gomock.Any(), coworkingID).Return(false, nil)
				cr.EXPECT().SetActive(gomock.Any(), coworkingID, false).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Деактивация коворкинга без активных бронирований",
		},
		{
			name:     "deactivate_coworking_with_active_bookings",
			isActive: false,
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().CheckHasActiveBookings(gomock.Any(), coworkingID).Return(true, nil)
			},
			wantError: ErrCoworkingHasActiveBookings,
			wantInTx:  true,
			desc:      "Попытка деактивации коворкинга с активными бронированиями",
		},
		{
			name:     "coworking_not_found",
			isActive: true,
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().SetActive(gomock.Any(), coworkingID, true).Return(repository.ErrCoworkingNotFound)
			},
			wantError: ErrCoworkingNotFound,
			wantInTx:  true,
			desc:      "Коворкинг не найден",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCoworking := mocks.NewMockCoworkingRepository(ctrl)
			txTracker := &transactionTracker{}

			tt.setup(mockCoworking)

			svc := &BookingService{
				coworkingRepo: mockCoworking,
				txManager:     txTracker,
			}

			err := svc.SetCoworkingActive(context.Background(), coworkingID, tt.isActive)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("SetCoworkingActive() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}

			if tt.wantInTx && !txTracker.withInTransactionCalled {
				t.Errorf("Expected transaction to be called | %s", tt.desc)
			}
		})
	}
}

// ============================================================================
// TESTS: SetLayoutVersionToActive
// ============================================================================

func TestSetLayoutVersionToActive(t *testing.T) {
	coworkingID := uuid.New()
	layoutVersion := 2

	tests := []struct {
		name      string
		setup     func(*mocks.MockCoworkingRepository)
		wantError error
		wantInTx  bool
		desc      string
	}{
		{
			name: "successful_activation",
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().DisableAllLayoutsByCoworking(gomock.Any(), coworkingID).Return(nil)
				cr.EXPECT().SetLayoutActiveByVersion(gomock.Any(), coworkingID, layoutVersion).Return(nil)
			},
			wantError: nil,
			wantInTx:  true,
			desc:      "Успешная активация версии layout",
		},
		{
			name: "disable_layouts_error",
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().DisableAllLayoutsByCoworking(gomock.Any(), coworkingID).Return(errors.New("database error"))
			},
			wantError: ErrCannotSetActiveLayout,
			wantInTx:  true,
			desc:      "Ошибка при отключении других layout версий",
		},
		{
			name: "layout_not_found",
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().DisableAllLayoutsByCoworking(gomock.Any(), coworkingID).Return(nil)
				cr.EXPECT().SetLayoutActiveByVersion(gomock.Any(), coworkingID, layoutVersion).Return(repository.ErrLayoutNotFound)
			},
			wantError: ErrLayoutNotFound,
			wantInTx:  true,
			desc:      "Layout версия не найдена",
		},
		{
			name: "set_active_error",
			setup: func(cr *mocks.MockCoworkingRepository) {
				cr.EXPECT().DisableAllLayoutsByCoworking(gomock.Any(), coworkingID).Return(nil)
				cr.EXPECT().SetLayoutActiveByVersion(gomock.Any(), coworkingID, layoutVersion).Return(errors.New("database error"))
			},
			wantError: ErrCannotSetActiveLayout,
			wantInTx:  true,
			desc:      "Ошибка при установки активности layout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCoworking := mocks.NewMockCoworkingRepository(ctrl)
			txTracker := &transactionTracker{}

			tt.setup(mockCoworking)

			svc := &BookingService{
				coworkingRepo: mockCoworking,
				txManager:     txTracker,
			}

			err := svc.SetLayoutVersionToActive(context.Background(), coworkingID, layoutVersion)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("SetLayoutVersionToActive() error = %v, wantErr %v | %s", err, tt.wantError, tt.desc)
			}

			if tt.wantInTx && !txTracker.withInTransactionCalled {
				t.Errorf("Expected transaction to be called | %s", tt.desc)
			}
		})
	}
}
