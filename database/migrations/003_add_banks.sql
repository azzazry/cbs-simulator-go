-- Migration: Add banks master data table
-- Purpose: Store list of supported banks with SWIFT codes
-- Created: 2026-03-07

CREATE TABLE IF NOT EXISTS banks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bank_code VARCHAR(10) NOT NULL UNIQUE,
    bank_name VARCHAR(100) NOT NULL,
    swift_code VARCHAR(11),
    is_active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on bank_code for faster lookup
CREATE INDEX IF NOT EXISTS idx_banks_code ON banks(bank_code);

-- Data inserted via seeder: 002_sample_banks.sql