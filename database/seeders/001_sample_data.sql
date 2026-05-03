-- CBS Simulator Sample Data
-- Indonesian Banking Sample Data

-- Insert Sample Customers
INSERT INTO customers (cif, full_name, id_card_number, phone_number, email, address, date_of_birth, pin, status) VALUES
('CIF001', 'Budi Santoso', '3201011234567890', '081234567890', 'budi.santoso@email.com', 'Jl. Sudirman No. 123, Jakarta Selatan', '1985-03-15', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa', 'active'),
('CIF002', 'Siti Nurhaliza', '3201021234567891', '081234567891', 'siti.nur@email.com', 'Jl. Gatot Subroto No. 45, Jakarta Pusat', '1990-07-20', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa', 'active'),
('CIF003', 'Ahmad Wijaya', '3201031234567892', '081234567892', 'ahmad.wijaya@email.com', 'Jl. Thamrin No. 78, Jakarta Pusat', '1988-11-10', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa', 'active'),
('CIF004', 'Dewi Lestari', '3201041234567893', '081234567893', 'dewi.lestari@email.com', 'Jl. Kuningan No. 56, Jakarta Selatan', '1992-05-25', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa', 'active'),
('CIF005', 'Rizki Pratama', '3201051234567894', '081234567894', 'rizki.p@email.com', 'Jl. Casablanca No. 88, Jakarta Selatan', '1987-09-08', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa', 'active');

-- Insert Sample Accounts (Savings, Checking)
INSERT INTO accounts (account_number, cif, account_type, currency, balance, avail_balance, status, opened_date, branch) VALUES
-- Budi's accounts
('1001234567', 'CIF001', 'savings', 'IDR', 25000000.00, 25000000.00, 'active', '2020-01-15', 'Jakarta Sudirman'),
('2001234567', 'CIF001', 'checking', 'IDR', 10000000.00, 10000000.00, 'active', '2020-01-15', 'Jakarta Sudirman'),
-- Siti's accounts
('1001234568', 'CIF002', 'savings', 'IDR', 50000000.00, 50000000.00, 'active', '2019-06-20', 'Jakarta Gatsu'),
('3001234568', 'CIF002', 'loan', 'IDR', -150000000.00, -150000000.00, 'active', '2021-03-10', 'Jakarta Gatsu'),
-- Ahmad's accounts
('1001234569', 'CIF003', 'savings', 'IDR', 75000000.00, 75000000.00, 'active', '2018-11-05', 'Jakarta Thamrin'),
('2001234569', 'CIF003', 'checking', 'IDR', 30000000.00, 30000000.00, 'active', '2018-11-05', 'Jakarta Thamrin'),
-- Dewi's accounts
('1001234570', 'CIF004', 'savings', 'IDR', 15000000.00, 15000000.00, 'active', '2021-02-14', 'Jakarta Kuningan'),
('4001234570', 'CIF004', 'deposit', 'IDR', 100000000.00, 100000000.00, 'active', '2022-05-20', 'Jakarta Kuningan'),
-- Rizki's accounts
('1001234571', 'CIF005', 'savings', 'IDR', 40000000.00, 40000000.00, 'active', '2019-08-30', 'Jakarta Casablanca'),
('2001234571', 'CIF005', 'checking', 'IDR', 20000000.00, 20000000.00, 'active', '2019-08-30', 'Jakarta Casablanca');

-- Insert Sample Cards
INSERT INTO cards (card_number, cif, account_number, card_type, card_brand, card_limit, avail_limit, expiry_date, status, cvv, pin) VALUES
('4234567890123456', 'CIF001', '1001234567', 'debit', 'visa', 10000000.00, 10000000.00, '12/2026', 'active', '123', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa'),
('5234567890123456', 'CIF002', '1001234568', 'credit', 'mastercard', 20000000.00, 15000000.00, '06/2027', 'active', '456', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa'),
('4234567890123457', 'CIF003', '1001234569', 'debit', 'visa', 15000000.00, 15000000.00, '03/2026', 'active', '789', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa'),
('5234567890123457', 'CIF004', '1001234570', 'debit', 'mastercard', 5000000.00, 5000000.00, '09/2025', 'active', '321', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa'),
('4234567890123458', 'CIF005', '1001234571', 'credit', 'visa', 30000000.00, 25000000.00, '11/2027', 'active', '654', '$2a$10$b/m9xktETKD2Km5HTXMSKeVRbrMgB8kO3fWpWD9kmGwOyuK5mWFpa');

-- Insert Sample Loans
INSERT INTO loans (loan_number, cif, account_number, loan_type, principal_amount, outstanding_amount, interest_rate, monthly_payment, tenor_months, remaining_months, disbursement_date, maturity_date, next_payment_date, status) VALUES
('LOAN001', 'CIF002', '3001234568', 'mortgage', 200000000.00, 150000000.00, 8.50, 2500000.00, 120, 72, '2021-03-10', '2031-03-10', '2026-04-10', 'active'),
('LOAN002', 'CIF003', '1001234569', 'personal', 50000000.00, 30000000.00, 12.00, 1500000.00, 36, 20, '2023-06-15', '2026-06-15', '2026-04-15', 'active'),
('LOAN003', 'CIF005', '1001234571', 'business', 100000000.00, 75000000.00, 10.00, 2000000.00, 60, 36, '2022-11-20', '2027-11-20', '2026-04-20', 'active');

-- Insert Sample Deposits
INSERT INTO deposits (deposit_number, cif, principal_amount, interest_rate, tenor_months, open_date, maturity_date, maturity_amount, auto_renew, status, linked_account) VALUES
('DEP001', 'CIF001', 50000000.00, 5.50, 12, '2025-03-01', '2026-03-01', 52750000.00, true, 'active', '1001234567'),
('DEP002', 'CIF004', 100000000.00, 6.00, 24, '2022-05-20', '2024-05-20', 112360000.00, false, 'matured', '1001234570'),
('DEP003', 'CIF005', 75000000.00, 5.75, 18, '2024-09-10', '2026-03-10', 81468750.00, true, 'active', '1001234571');

-- Insert Sample Bill Payments (PPOB)
INSERT INTO bill_payments (biller_code, biller_name, customer_number, bill_number, bill_amount, admin_fee, total_amount, bill_period, due_date, status) VALUES
('PLN', 'PT PLN (Persero)', '123456789012', 'BILL202603001', 450000.00, 2500.00, 452500.00, '2026-02', '2026-03-20', 'unpaid'),
('PDAM', 'PDAM DKI Jakarta', '987654321098', 'BILL202603002', 125000.00, 1500.00, 126500.00, '2026-02', '2026-03-15', 'unpaid'),
('TELKOM', 'PT Telkom Indonesia', '021-12345678', 'BILL202603003', 300000.00, 2000.00, 302000.00, '2026-02', '2026-03-25', 'unpaid'),
('INDIHOME', 'IndiHome', '1234567890', 'BILL202603004', 500000.00, 2500.00, 502500.00, '2026-02', '2026-03-18', 'unpaid'),
('BPJS', 'BPJS Kesehatan', '0001234567890', 'BILL202603005', 150000.00, 0.00, 150000.00, '2026-03', '2026-03-31', 'unpaid');

-- Insert Sample Transactions (Transaction History)
INSERT INTO transactions (transaction_id, transaction_type, from_account_number, to_account_number, amount, currency, description, reference_number, status, transaction_date, settlement_date, fee) VALUES
('TRX20260301001', 'transfer_intra', '1001234567', '1001234568', 1000000.00, 'IDR', 'Transfer ke Siti', 'REF001', 'success', '2026-03-01 10:30:00', '2026-03-01', 0.00),
('TRX20260301002', 'transfer_inter', '1001234569', '1234567890', 2000000.00, 'IDR', 'Transfer ke Bank Lain', 'REF002', 'success', '2026-03-01 14:15:00', '2026-03-02', 6500.00),
('TRX20260302001', 'payment_bill', '1001234567', NULL, 452500.00, 'IDR', 'Bayar PLN', 'BILL202603001', 'success', '2026-03-02 09:00:00', '2026-03-02', 0.00),
('TRX20260302002', 'withdrawal_atm', '1001234570', NULL, 500000.00, 'IDR', 'Tarik Tunai ATM', 'ATM001', 'success', '2026-03-02 16:45:00', '2026-03-02', 0.00),
('TRX20260303001', 'transfer_intra', '1001234571', '1001234567', 5000000.00, 'IDR', 'Transfer Gaji', 'REF003', 'success', '2026-03-03 08:00:00', '2026-03-03', 0.00);
