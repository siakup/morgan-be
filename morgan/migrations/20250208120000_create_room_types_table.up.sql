CREATE SCHEMA IF NOT EXISTS master;

CREATE TABLE IF NOT EXISTS master.room_types
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    description TEXT,
    is_active  BOOLEAN      NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_by UUID,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE UNIQUE INDEX idx_room_types_name_active ON master.room_types (name) WHERE deleted_at IS NULL;

COMMENT ON TABLE master.room_types IS 'Master data: room types for scheduling/booking';
