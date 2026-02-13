CREATE TABLE IF NOT EXISTS master.domains (
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

CREATE INDEX idx_domains_status 
ON master.domains(status);
