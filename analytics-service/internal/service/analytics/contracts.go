package analytics_service

import (
	"context"

	"github.com/4udiwe/coworking/analytics-service/internal/entity"
	"github.com/google/uuid"
)

// AnalyticsRepository определяет интерфейс для работы с репозиторием аналитики бронирований.
type AnalyticsRepository interface {
	GetCoworkingHourlyLoad(ctx context.Context, coworkingID uuid.UUID) (map[int]int, error)
	GetCoworkingWeekdayLoad(ctx context.Context, coworkingID uuid.UUID) (map[int]int, error)
	GetCoworkingHeatmap(ctx context.Context, coworkingID uuid.UUID) ([]entity.HeatmapCell, error)
	GetPlaceHeatmap(ctx context.Context, placeID uuid.UUID) ([]entity.HeatmapCell, error)
	InsertEvents(ctx context.Context, events []entity.BookingEvent) error
	InsertBookingState(ctx context.Context, events []entity.BookingEvent) error
}
