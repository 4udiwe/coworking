-- +goose Up
-- +goose StatementBegin
-- ==============================
-- EXTENSIONS
-- ==============================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- ==============================
-- COWORKINGS
-- ==============================

CREATE TABLE IF NOT EXISTS coworking (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            TEXT NOT NULL,
    address         TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_coworking_active
    ON coworking(is_active);

-- ==============================
-- COWORKING LAYOUTS
-- ==============================

CREATE TABLE IF NOT EXISTS coworking_layout (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coworking_id    UUID NOT NULL REFERENCES coworking(id) ON DELETE CASCADE,
    version         INT NOT NULL,
    layout          JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_layout_version UNIQUE (coworking_id, version)
);

CREATE INDEX idx_layouts_coworking
    ON coworking_layout(coworking_id);

-- ==============================
-- PLACES
-- ==============================

CREATE TABLE IF NOT EXISTS place (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coworking_id    UUID NOT NULL REFERENCES coworking(id) ON DELETE CASCADE,
    label           TEXT NOT NULL,
    place_type      TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_places_coworking
    ON place(coworking_id);

CREATE INDEX idx_places_active
    ON place(is_active);

-- ==============================
-- BOOKINGS
-- ==============================

CREATE TABLE IF NOT EXISTS booking_status (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE
);

INSERT INTO booking_status (id, name) VALUES
    (1, 'active'),
    (2, 'cancelled'),
    (3, 'completed')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS booking (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL,
    place_id        UUID NOT NULL REFERENCES place(id) ON DELETE CASCADE,

    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,

    status_id       INT NOT NULL REFERENCES booking_status(id) DEFAULT 1,
    cancel_reason   TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    cancelled_at    TIMESTAMPTZ,

    -- ==============================
    -- INVARIANTS
    -- ==============================

    CONSTRAINT chk_time_order
        CHECK (end_time > start_time),

    CONSTRAINT chk_duration_hours
        CHECK (
            EXTRACT(EPOCH FROM (end_time - start_time)) IN (3600, 7200, 10800)
        )
);

-- ==============================
-- EXCLUSION CONSTRAINT
-- (NO OVERLAPPING ACTIVE BOOKINGS)
-- ==============================

ALTER TABLE booking
ADD CONSTRAINT no_overlapping_active_bookings
EXCLUDE USING gist (
    place_id WITH =,
    tstzrange(start_time, end_time) WITH &&
)
WHERE (status_id = 1);


-- ==============================
-- INDEXES
-- ==============================

CREATE INDEX idx_booking_place_time
    ON booking(place_id, start_time, end_time);

CREATE INDEX idx_booking_user
    ON booking(user_id);

CREATE INDEX idx_booking_status
    ON booking(status_id);

CREATE INDEX idx_booking_time_range
    ON booking(start_time, end_time);

-- ================================
--  OUTBOX
-- ================================
CREATE TABLE outbox_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO outbox_status (name) VALUES
('pending'), ('processed'), ('failed');

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
DROP TABLE IF EXISTS booking_status CASCADE;
DROP TABLE IF EXISTS booking CASCADE;
DROP TABLE IF EXISTS place CASCADE;
DROP TABLE IF EXISTS coworking_layout CASCADE;
DROP TABLE IF EXISTS coworking CASCADE;

DROP TYPE IF EXISTS booking_status;

-- +goose StatementEnd