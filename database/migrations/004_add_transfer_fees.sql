-- Migration: Add transfer fees table
CREATE TABLE IF NOT EXISTS transfer_fees (
    id                    BIGSERIAL PRIMARY KEY,
    destination_bank_code VARCHAR(10)    NOT NULL UNIQUE,
    destination_bank_name VARCHAR(100)   NOT NULL,
    fee_type              VARCHAR(20)    NOT NULL DEFAULT 'flat',
    fee_amount            DECIMAL(15, 2) NOT NULL,
    fee_percentage        DECIMAL(5, 2),
    minimum_amount        DECIMAL(15, 2) DEFAULT 0,
    maximum_amount        DECIMAL(15, 2) DEFAULT 999999999,
    description           VARCHAR(255),
    is_active             BOOLEAN        DEFAULT TRUE,
    effective_from        TIMESTAMPTZ    DEFAULT NOW(),
    effective_to          TIMESTAMPTZ,
    notes                 VARCHAR(255),
    created_at            TIMESTAMPTZ    DEFAULT NOW(),
    updated_at            TIMESTAMPTZ    DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transfer_fees_code   ON transfer_fees(destination_bank_code);
CREATE INDEX IF NOT EXISTS idx_transfer_fees_active  ON transfer_fees(is_active);

CREATE TABLE IF NOT EXISTS service_fees (
    id              BIGSERIAL PRIMARY KEY,
    service_code    VARCHAR(50)    NOT NULL UNIQUE,
    service_name    VARCHAR(100)   NOT NULL,
    service_type    VARCHAR(50)    NOT NULL,
    fee_type        VARCHAR(20)    NOT NULL DEFAULT 'flat',
    fee_amount      DECIMAL(15, 2),
    fee_percentage  DECIMAL(5, 2),
    minimum_amount  DECIMAL(15, 2) DEFAULT 0,
    maximum_amount  DECIMAL(15, 2) DEFAULT 999999999,
    is_active       BOOLEAN        DEFAULT TRUE,
    effective_from  TIMESTAMPTZ    DEFAULT NOW(),
    effective_to    TIMESTAMPTZ,
    notes           VARCHAR(255),
    created_at      TIMESTAMPTZ    DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_service_fees_code   ON service_fees(service_code);
CREATE INDEX IF NOT EXISTS idx_service_fees_type   ON service_fees(service_type);
CREATE INDEX IF NOT EXISTS idx_service_fees_active ON service_fees(is_active);
