CREATE SCHEMA IF NOT EXISTS organization;

CREATE TABLE IF NOT EXISTS organization.domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    status BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,

    created_by UUID NULL,
    updated_by UUID NULL,
    deleted_by UUID NULL
);

CREATE INDEX idx_shift_sessions_status 
ON organization.domains(status);
