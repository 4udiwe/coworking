-- +goose Up
------------------------------------------------
-- RAW EVENTS
------------------------------------------------
CREATE TABLE booking_events (
    event_id UUID,
    event_type String,
    booking_id UUID,
    coworking_id UUID,
    place_id UUID,
    user_id UUID,
    start_time DateTime,
    end_time DateTime,
    status String,
    occurred_at DateTime
) ENGINE = MergeTree PARTITION BY toYYYYMM(start_time)
ORDER BY (coworking_id, place_id, start_time);
------------------------------------------------
-- BOOKING STATE (Lifecycle)
------------------------------------------------
CREATE TABLE booking_state (
    booking_id UUID,
    coworking_id UUID,
    place_id UUID,
    user_id UUID,
    start_time DateTime,
    end_time DateTime,
    status Enum8(
        'created' = 1,
        'cancelled' = 2,
        'completed' = 3
    ),
    updated_at DateTime
) ENGINE = ReplacingMergeTree(updated_at)
ORDER BY (coworking_id, place_id, start_time, booking_id);
------------------------------------------------
-- AGGREGATED TABLES (COWORKING LEVEL)
------------------------------------------------
CREATE TABLE coworking_hourly (
    coworking_id UUID,
    hour UInt8,
    bookings UInt64
) ENGINE = SummingMergeTree
ORDER BY (coworking_id, hour);
CREATE TABLE coworking_weekday (
    coworking_id UUID,
    weekday UInt8,
    bookings UInt64
) ENGINE = SummingMergeTree
ORDER BY (coworking_id, weekday);
CREATE TABLE coworking_heatmap (
    coworking_id UUID,
    weekday UInt8,
    hour UInt8,
    bookings UInt64
) ENGINE = SummingMergeTree
ORDER BY (coworking_id, weekday, hour);
------------------------------------------------
-- OPTIONAL PLACE ANALYTICS
------------------------------------------------
CREATE TABLE place_heatmap (
    place_id UUID,
    weekday UInt8,
    hour UInt8,
    bookings UInt64
) ENGINE = SummingMergeTree
ORDER BY (place_id, weekday, hour);
------------------------------------------------
-- MATERIALIZED VIEWS
------------------------------------------------
-- Hourly load for coworking
CREATE MATERIALIZED VIEW coworking_hourly_mv TO coworking_hourly AS
SELECT coworking_id,
    arrayJoin(
        range(
            toHour(start_time),
            toHour(end_time)
        )
    ) AS hour,
    1 AS bookings
FROM booking_state
WHERE status != 'cancelled';
------------------------------------------------
-- Weekday popularity
CREATE MATERIALIZED VIEW coworking_weekday_mv TO coworking_weekday AS
SELECT coworking_id,
    toDayOfWeek(start_time) AS weekday,
    count() AS bookings
FROM booking_state
WHERE status != 'cancelled'
GROUP BY coworking_id,
    weekday;
------------------------------------------------
-- Global heatmap
CREATE MATERIALIZED VIEW coworking_heatmap_mv TO coworking_heatmap AS
SELECT coworking_id,
    toDayOfWeek(start_time) AS weekday,
    arrayJoin(
        range(
            toHour(start_time),
            toHour(end_time)
        )
    ) AS hour,
    1 AS bookings
FROM booking_state
WHERE status != 'cancelled';
------------------------------------------------
-- Place heatmap (optional)
CREATE MATERIALIZED VIEW place_heatmap_mv TO place_heatmap AS
SELECT place_id,
    toDayOfWeek(start_time) AS weekday,
    arrayJoin(
        range(
            toHour(start_time),
            toHour(end_time)
        )
    ) AS hour,
    1 AS bookings
FROM booking_state
WHERE status != 'cancelled';
-- +goose Down
DROP VIEW place_heatmap_mv;
DROP VIEW coworking_heatmap_mv;
DROP VIEW coworking_weekday_mv;
DROP VIEW coworking_hourly_mv;
DROP TABLE place_heatmap;
DROP TABLE coworking_heatmap;
DROP TABLE coworking_weekday;
DROP TABLE coworking_hourly;
DROP TABLE booking_state;
DROP TABLE booking_events;