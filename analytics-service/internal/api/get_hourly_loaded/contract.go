package get_hourly_loaded

import (
	"context"

	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetHourlyLoad(ctx context.Context, coworkingID uuid.UUID) (map[int]int, error)
}
