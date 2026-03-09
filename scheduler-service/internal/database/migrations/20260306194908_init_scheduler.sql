-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==============================
-- TIMER TYPES
-- ==============================

CREATE TABLE timer_type (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO timer_type (id, name) VALUES
(1, 'booking_reminder'),
(2, 'booking_expire');

-- ==============================
-- TIMER STATUS
-- ==============================

CREATE TABLE timer_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO timer_status (name) VALUES
('pending'),
('triggered'),
('cancelled');

-- ==============================
-- TIMERS
-- ==============================

CREATE TABLE timer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    timer_type_id SMALLINT NOT NULL REFERENCES timer_type(id),

    booking_id UUID NOT NULL,
    user_id UUID NULL,

    trigger_at TIMESTAMPTZ NOT NULL,

    payload JSONB,

    status_id SMALLINT NOT NULL REFERENCES timer_status(id) DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    triggered_at TIMESTAMPTZ NULL,
    cancelled_at TIMESTAMPTZ NULL
);

-- ==============================
-- INDEXES
-- ==============================

CREATE INDEX idx_timers_trigger_at
    ON timer (trigger_at);

CREATE INDEX idx_timers_status
    ON timer (status_id);

CREATE INDEX idx_timers_booking
    ON timer (booking_id);

CREATE INDEX idx_timers_pending_trigger
    ON timer (status_id, trigger_at);

-- ==============================
-- OUTBOX
-- ==============================

CREATE TABLE outbox_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO outbox_status (name) VALUES
('pending'),
('processed'),
('failed');

CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id UUID NOT NULL,

    event_type VARCHAR(128) NOT NULL,

    payload JSONB NOT NULL,

    status_id SMALLINT NOT NULL REFERENCES outbox_status(id) DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_outbox_status_id ON outbox (status_id);
CREATE INDEX idx_outbox_created_at ON outbox (created_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS outbox_status CASCADE;

DROP TABLE IF EXISTS timer CASCADE;
DROP TABLE IF EXISTS timer_status CASCADE;
DROP TABLE IF EXISTS timer_type CASCADE;

-- +goose StatementEnd