-- Migration: Add transfer fees table (Revised for Single Bank)
-- Purpose: Store dynamic fees for transfers OUT to other banks
-- Architecture: This CBS system represents ONE bank (e.g., Bank Daerah Jaya)
-- We only store fees for OUR bank → OTHER banks (OUTBOUND transfers)
-- Intrabank transfers are FREE
-- Incoming transfers from other banks are FREE (received by our customers)
-- Created: 2026-03-07

CREATE TABLE IF NOT EXISTS transfer_fees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    destination_bank_code VARCHAR(10) NOT NULL UNIQUE,
    destination_bank_name VARCHAR(100) NOT NULL,
    fee_type VARCHAR(20) NOT NULL DEFAULT 'flat', -- 'flat' for fixed amount, 'percentage' for %
    fee_amount DECIMAL(15, 2) NOT NULL,
    fee_percentage DECIMAL(5, 2), -- For percentage-based fees
    minimum_amount DECIMAL(15, 2) DEFAULT 0,
    maximum_amount DECIMAL(15, 2) DEFAULT 999999999,
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT 1,
    effective_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    effective_to TIMESTAMP,
    notes VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookup
CREATE INDEX IF NOT EXISTS idx_transfer_fees_code ON transfer_fees(destination_bank_code);
CREATE INDEX IF NOT EXISTS idx_transfer_fees_active ON transfer_fees(is_active);

-- Create table for other service fees (e-wallet, e-money, VA, QRIS)
CREATE TABLE IF NOT EXISTS service_fees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_code VARCHAR(50) NOT NULL UNIQUE,
    service_name VARCHAR(100) NOT NULL,
    service_type VARCHAR(50) NOT NULL, -- 'topup_ewallet', 'topup_emoney', 'payment_va', 'qris_payment', etc
    fee_type VARCHAR(20) NOT NULL DEFAULT 'flat', -- 'flat' or 'percentage'
    fee_amount DECIMAL(15, 2), -- For flat fee
    fee_percentage DECIMAL(5, 2), -- For percentage fee
    minimum_amount DECIMAL(15, 2) DEFAULT 0,
    maximum_amount DECIMAL(15, 2) DEFAULT 999999999,
    is_active BOOLEAN DEFAULT 1,
    effective_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    effective_to TIMESTAMP,
    notes VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_service_fees_code ON service_fees(service_code);
CREATE INDEX IF NOT EXISTS idx_service_fees_type ON service_fees(service_type);
CREATE INDEX IF NOT EXISTS idx_service_fees_active ON service_fees(is_active);

-- Data inserted via seeder: 003_sample_transfer_fees.sql
