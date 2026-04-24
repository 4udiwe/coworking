package get_weekday_loaded

import (
	"context"

	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetWeekdayLoad(ctx context.Context, coworkingID uuid.UUID) (map[int]int, error)
}
