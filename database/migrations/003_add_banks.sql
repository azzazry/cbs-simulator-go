-- Migration: Add banks master data table
CREATE TABLE IF NOT EXISTS banks (
    id         BIGSERIAL PRIMARY KEY,
    bank_code  VARCHAR(10)  NOT NULL UNIQUE,
    bank_name  VARCHAR(100) NOT NULL,
    swift_code VARCHAR(11),
    is_active  BOOLEAN      DEFAULT TRUE,
    created_at TIMESTAMPTZ  DEFAULT NOW(),
    updated_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_banks_code ON banks(bank_code);
