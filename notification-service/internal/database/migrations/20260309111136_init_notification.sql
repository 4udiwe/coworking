-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==============================
-- NOTIFICATION TYPES
-- ==============================

CREATE TABLE notification_type (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE
);

INSERT INTO notification_type (id, name) VALUES
(1, 'booking_created'),
(2, 'booking_cancelled'),
(3, 'booking_reminder'),
(4, 'booking_expired');

-- ==============================
-- NOTIFICATION STATUS
-- ==============================

CREATE TABLE notification_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO notification_status (name) VALUES
('unread'),
('read');

-- ==============================
-- USER DEVICES
-- ==============================

CREATE TABLE user_device (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    device_token TEXT NOT NULL,
    platform VARCHAR(16) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(device_token)
);

CREATE INDEX idx_user_device_user
    ON user_device (user_id);

-- ==============================
-- NOTIFICATIONS
-- ==============================

CREATE TABLE notification (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    notification_type_id SMALLINT NOT NULL
        REFERENCES notification_type(id),

    title TEXT NOT NULL,
    body TEXT NOT NULL,

    payload JSONB,

    status_id SMALLINT NOT NULL
        REFERENCES notification_status(id)
        DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ NULL
);

-- ==============================
-- INDEXES
-- ==============================

CREATE INDEX idx_notification_user
    ON notification (user_id);

CREATE INDEX idx_notification_status
    ON notification (status_id);

CREATE INDEX idx_notification_user_status
    ON notification (user_id, status_id);

CREATE INDEX idx_notification_created_at
    ON notification (created_at);

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

DROP TABLE IF EXISTS notification CASCADE;
DROP TABLE IF EXISTS user_device CASCADE;

DROP TABLE IF EXISTS notification_status CASCADE;
DROP TABLE IF EXISTS notification_type CASCADE;

-- +goose StatementEnd