-- Security Seed Data
-- CBS Simulator Phase 1: Security & Compliance

-- Seed Roles
INSERT OR IGNORE INTO roles (role_name, description) VALUES ('customer', 'Regular bank customer with standard access');
INSERT OR IGNORE INTO roles (role_name, description) VALUES ('teller', 'Bank teller with counter operation access');
INSERT OR IGNORE INTO roles (role_name, description) VALUES ('admin', 'System administrator with full access');
INSERT OR IGNORE INTO roles (role_name, description) VALUES ('supervisor', 'Branch supervisor with approval authority');

-- Assign customer role to all existing customers
INSERT OR IGNORE INTO user_roles (cif, role_id, assigned_by) 
SELECT c.cif, r.id, 'SYSTEM' FROM customers c, roles r WHERE r.role_name = 'customer';

-- Assign admin role to CIF001 (Budi Santoso) for testing
INSERT OR IGNORE INTO user_roles (cif, role_id, assigned_by) 
SELECT 'CIF001', r.id, 'SYSTEM' FROM roles r WHERE r.role_name = 'admin';

-- Assign supervisor role to CIF003 (Ahmad Wijaya) for testing  
INSERT OR IGNORE INTO user_roles (cif, role_id, assigned_by) 
SELECT 'CIF003', r.id, 'SYSTEM' FROM roles r WHERE r.role_name = 'supervisor';

-- Default Transaction Limits for Customer role
INSERT OR IGNORE INTO transaction_limits (role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit) VALUES
('customer', 'intra_transfer', 50000000, 25000000, 500000000),
('customer', 'inter_transfer', 25000000, 10000000, 250000000),
('customer', 'bill_payment', 10000000, 5000000, 100000000),
('customer', 'qris_payment', 5000000, 2000000, 50000000),
('customer', 'va_payment', 50000000, 25000000, 500000000),
('customer', 'ewallet_topup', 2000000, 1000000, 20000000),
('customer', 'emoney_topup', 2000000, 1000000, 20000000);

-- Default Transaction Limits for Teller role (higher limits)
INSERT OR IGNORE INTO transaction_limits (role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit) VALUES
('teller', 'intra_transfer', 500000000, 100000000, 5000000000),
('teller', 'inter_transfer', 250000000, 50000000, 2500000000),
('teller', 'bill_payment', 100000000, 50000000, 1000000000),
('teller', 'qris_payment', 50000000, 20000000, 500000000),
('teller', 'va_payment', 500000000, 100000000, 5000000000),
('teller', 'ewallet_topup', 20000000, 10000000, 200000000),
('teller', 'emoney_topup', 20000000, 10000000, 200000000);

-- Supervisor and Admin have no enforced limits (controlled by role check)
