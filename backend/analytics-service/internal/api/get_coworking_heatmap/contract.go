package get_coworking_heatmap

import (
	"context"

	"github.com/4udiwe/coworking/analytics-service/internal/entity"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetCoworkingHeatmap(ctx context.Context, coworkingID uuid.UUID) ([]entity.HeatmapCell, error)
}
