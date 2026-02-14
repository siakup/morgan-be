CREATE TABLE IF NOT EXISTS hr.shift_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    start TIME NOT NULL,
    "end" TIME NOT NULL,
    status BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,

    created_by UUID NULL,
    updated_by UUID NULL,
    deleted_by UUID NULL
);

CREATE INDEX idx_shift_sessions_status 
ON hr.shift_sessions(status);
