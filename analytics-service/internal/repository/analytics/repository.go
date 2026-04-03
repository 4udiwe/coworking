package analytics_repository

import (
	"context"

	"github.com/4udiwe/coworking/analytics-service/internal/entity"
	"github.com/4udiwe/coworking/analytics-service/pkg/clickhouse"
	"github.com/google/uuid"
)

type AnalyticsRepository struct {
	ch *clickhouse.ClickHouse
}

func New(ch *clickhouse.ClickHouse) *AnalyticsRepository {
	return &AnalyticsRepository{ch: ch}
}

func (r *AnalyticsRepository) InsertEvents(
	ctx context.Context,
	events []entity.BookingEvent,
) error {

	batch, err := r.ch.PrepareBatch(ctx, "INSERT INTO booking_events")
	if err != nil {
		return err
	}

	for _, e := range events {

		if err := batch.Append(
			e.EventID,
			e.EventType,
			e.BookingID,
			e.CoworkingID,
			e.PlaceID,
			e.UserID,
			e.StartTime,
			e.EndTime,
			e.BookingStatus,
			e.Occurred,
		); err != nil {
			return err
		}
	}

	return batch.Send()
}

func (r *AnalyticsRepository) InsertBookingState(
	ctx context.Context,
	events []entity.BookingEvent,
) error {

	batch, err := r.ch.PrepareBatch(ctx, "INSERT INTO booking_state")
	if err != nil {
		return err
	}

	for _, e := range events {

		if err := batch.Append(
			e.BookingID,
			e.CoworkingID,
			e.PlaceID,
			e.UserID,
			e.StartTime,
			e.EndTime,
			e.BookingStatus,
			e.Occurred,
		); err != nil {
			return err
		}
	}

	return batch.Send()
}

func (r *AnalyticsRepository) GetCoworkingHourlyLoad(
	ctx context.Context,
	coworkingID uuid.UUID,
) (map[int]int, error) {

	rows, err := r.ch.Conn().Query(ctx,
		`
        SELECT hour, sum(bookings)
        FROM coworking_hourly
        WHERE coworking_id = ?
        GROUP BY hour
        ORDER BY hour
        `,
		coworkingID,
	)
	if err != nil {
		return nil, err
	}

	result := make(map[int]int)

	for rows.Next() {

		var hour uint8
		var count uint64

		if err := rows.Scan(&hour, &count); err != nil {
			return nil, err
		}

		result[int(hour)] = int(count)
	}

	return result, nil
}

func (r *AnalyticsRepository) GetCoworkingHourlyLoadByWeekday(
	ctx context.Context,
	coworkingID uuid.UUID,
	weekday int,
) (map[int]int, error) {

	rows, err := r.ch.Conn().Query(ctx,
		`
        SELECT hour, sum(bookings)
        FROM coworking_heatmap
        WHERE coworking_id = ? AND weekday = ?
        GROUP BY hour
        ORDER BY hour
        `,
		coworkingID,
		weekday,
	)
	if err != nil {
		return nil, err
	}

	result := make(map[int]int)

	for rows.Next() {

		var hour uint8
		var count uint64

		if err := rows.Scan(&hour, &count); err != nil {
			return nil, err
		}

		result[int(hour)] = int(count)
	}

	return result, nil
}

func (r *AnalyticsRepository) GetCoworkingWeekdayLoad(
	ctx context.Context,
	coworkingID uuid.UUID,
) (map[int]int, error) {

	rows, err := r.ch.Conn().Query(ctx,
		`
        SELECT weekday, sum(bookings)
        FROM coworking_weekday
        WHERE coworking_id = ?
        GROUP BY weekday
        `,
		coworkingID,
	)

	if err != nil {
		return nil, err
	}

	result := make(map[int]int)

	for rows.Next() {

		var weekday uint8
		var count uint64

		rows.Scan(&weekday, &count)

		result[int(weekday)] = int(count)
	}

	return result, nil
}

func (r *AnalyticsRepository) GetCoworkingHeatmap(
	ctx context.Context,
	coworkingID uuid.UUID,
) ([]entity.HeatmapCell, error) {

	rows, err := r.ch.Conn().Query(ctx,
		`
        SELECT weekday, hour, sum(bookings)
        FROM coworking_heatmap
        WHERE coworking_id = ?
        GROUP BY weekday, hour
        `,
		coworkingID,
	)
	if err != nil {
		return nil, err
	}

	var result []entity.HeatmapCell

	for rows.Next() {

		var cell entity.HeatmapCell

		if err := rows.Scan(&cell.Weekday, &cell.Hour, &cell.Count); err != nil {
			return nil, err
		}

		result = append(result, cell)
	}

	return result, nil
}

func (r *AnalyticsRepository) GetPlaceHeatmap(
	ctx context.Context,
	placeID uuid.UUID,
) ([]entity.HeatmapCell, error) {

	rows, err := r.ch.Conn().Query(ctx,
		`
        SELECT weekday, hour, sum(bookings)
        FROM place_heatmap
        WHERE place_id = ?
        GROUP BY weekday, hour
        `,
		placeID,
	)
	if err != nil {
		return nil, err
	}

	var result []entity.HeatmapCell

	for rows.Next() {

		var cell entity.HeatmapCell

		if err := rows.Scan(&cell.Weekday, &cell.Hour, &cell.Count); err != nil {
			return nil, err
		}

		result = append(result, cell)
	}

	return result, nil
}
