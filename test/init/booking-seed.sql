-- Create necessary tables for testing if they don't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Coworking table
CREATE TABLE IF NOT EXISTS coworking (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            TEXT NOT NULL,
    address         TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Place table
CREATE TABLE IF NOT EXISTS place (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coworking_id    UUID NOT NULL REFERENCES coworking(id) ON DELETE CASCADE,
    label           TEXT NOT NULL,
    place_type      TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed data for booking-service tests
-- Inserts test coworking and place for integration tests
INSERT INTO coworking (
        id,
        name,
        address,
        is_active,
        created_at,
        updated_at
    )
VALUES (
        '550e8400-e29b-41d4-a716-446655440000'::uuid,
        'Test Coworking Space',
        '123 Test Street, Test City, State 12345',
        true,
        NOW(),
        NOW()
    ) ON CONFLICT DO NOTHING;
-- Insert test places
INSERT INTO place (
        id,
        coworking_id,
        label,
        place_type,
        is_active,
        created_at,
        updated_at
    )
VALUES (
        '550e8400-e29b-41d4-a716-446655441001'::uuid,
        '550e8400-e29b-41d4-a716-446655440000'::uuid,
        'Open Desk A',
        'open_desk',
        true,
        NOW(),
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655441002'::uuid,
        '550e8400-e29b-41d4-a716-446655440000'::uuid,
        'Meeting Room B',
        'meeting_room',
        true,
        NOW(),
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655441003'::uuid,
        '550e8400-e29b-41d4-a716-446655440000'::uuid,
        'Private Office C',
        'private_office',
        true,
        NOW(),
        NOW()
    ) ON CONFLICT DO NOTHING;