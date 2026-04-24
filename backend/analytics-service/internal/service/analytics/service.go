package analytics_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/4udiwe/coworking/analytics-service/internal/entity"
)

// AnalyticsService предоставляет бизнес-логику для получения аналитики по бронированиям.
type AnalyticsService struct {
	repo AnalyticsRepository
}

// New создает новый экземпляр AnalyticsService.
func New(bookingRepo AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: bookingRepo}
}

// InsertEvents вставляет пакет событий бронирования в аналитическую базу данных.
// Принимает контекст и срез сущностей BookingEvent.
// В случае ошибки возвращает ErrCannotInsertEvents.
func (s *AnalyticsService) InsertEvents(ctx context.Context, events []entity.BookingEvent) error {
	logrus.Infof("Inserting events amount: %v", len(events))
	err := s.repo.InsertEvents(ctx, events)
	if err != nil {
		logrus.WithError(err).Error("Insert events failed")
		return ErrCannotInsertEvents
	}

	err = s.repo.InsertBookingState(ctx, events)
	if err != nil {
		logrus.WithError(err).Error("Insert state failed")
		return ErrCannotInsertEvents
	}

	logrus.Info("Insert performed")

	return nil
}

// GetHourlyLoad возвращает распределение количества бронирований по часам дня для коворкинга.
// Если weekday не указан (nil), возвращает усреднённые часы по всем дням недели.
// Если weekday указан (1-7), возвращает часы только для конкретного дня недели.
// Возвращает карту, где ключ - час суток (0-23), значение - количество бронирований.
// В случае ошибки возвращает ErrCannotFetchInfo.
func (s *AnalyticsService) GetHourlyLoad(ctx context.Context, coworkingID uuid.UUID, weekday *int) (map[int]int, error) {
	if weekday != nil {
		logrus.Infof("Getting hourly load for coworking ID: %s, weekday: %d", coworkingID, *weekday)
		result, err := s.repo.GetCoworkingHourlyLoadByWeekday(ctx, coworkingID, *weekday)
		if err != nil {
			logrus.Errorf("Failed to get hourly load by weekday: %v", err)
			return nil, ErrCannotFetchInfo
		}
		return result, nil
	}

	logrus.Infof("Getting hourly load for coworking ID: %s", coworkingID)
	result, err := s.repo.GetCoworkingHourlyLoad(ctx, coworkingID)
	if err != nil {
		logrus.Errorf("Failed to get hourly load: %v", err)
		return nil, ErrCannotFetchInfo
	}

	return result, nil
}

// GetWeekdayLoad возвращает распределение количества бронирований по дням недели для коворкинга.
// Возвращает карту, где ключ - день недели (1-7, понедельник=1), значение - количество бронирований.
// В случае ошибки возвращает ErrCannotFetchInfo.
func (s *AnalyticsService) GetWeekdayLoad(ctx context.Context, coworkingID uuid.UUID) (map[int]int, error) {
	logrus.Infof("Getting weekday load for coworking ID: %s", coworkingID)

	result, err := s.repo.GetCoworkingWeekdayLoad(ctx, coworkingID)
	if err != nil {
		logrus.Errorf("Failed to get weekday load: %v", err)
		return nil, ErrCannotFetchInfo
	}

	return result, nil
}

// GetHeatmap возвращает данные для построения тепловой карты загруженности (день недели × час) коворкинга.
// Возвращает срез HeatmapCell, где каждый элемент содержит день недели, час и количество бронирований.
// В случае ошибки возвращает ErrCannotFetchInfo.
func (s *AnalyticsService) GetCoworkingHeatmap(ctx context.Context, coworkingID uuid.UUID) ([]entity.HeatmapCell, error) {
	logrus.Infof("Getting heatmap for coworking ID: %s", coworkingID)

	result, err := s.repo.GetCoworkingHeatmap(ctx, coworkingID)
	if err != nil {
		logrus.Errorf("Failed to get heatmap: %v", err)
		return nil, ErrCannotFetchInfo
	}

	return result, nil
}

// GetHeatmap возвращает данные для построения тепловой карты загруженности (день недели × час) места.
// Возвращает срез HeatmapCell, где каждый элемент содержит день недели, час и количество бронирований.
// В случае ошибки возвращает ErrCannotFetchInfo.
func (s *AnalyticsService) GetPlaceHeatmap(ctx context.Context, placeID uuid.UUID) ([]entity.HeatmapCell, error) {
	logrus.Infof("Getting heatmap for place ID: %s", placeID)

	result, err := s.repo.GetPlaceHeatmap(ctx, placeID)
	if err != nil {
		logrus.Errorf("Failed to get heatmap: %v", err)
		return nil, ErrCannotFetchInfo
	}

	return result, nil
}
