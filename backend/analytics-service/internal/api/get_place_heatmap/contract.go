package get_place_heatmap

import (
	"context"

	"github.com/4udiwe/coworking/analytics-service/internal/entity"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetPlaceHeatmap(ctx context.Context, placeID uuid.UUID) ([]entity.HeatmapCell, error)
}
