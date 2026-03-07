-- CBS Simulator Sample Transfer Fees and Service Fees
-- Insert dynamic fees for interbank transfers and services

-- Transfer Fees: Fees for OUR bank → OTHER banks (OUTBOUND transfers)
-- Intrabank transfers are FREE
-- Incoming transfers from other banks are FREE (received by our customers)
INSERT OR IGNORE INTO transfer_fees 
    (destination_bank_code, destination_bank_name, fee_type, fee_amount, description, notes) 
VALUES
    -- Major Banks - Standard Domestic Fee (Rp 5,000)
    ('BCA', 'Bank BCA', 'flat', 5000, 'Transfer to BCA', 'Domestic SKNT'),
    ('MANDIRI', 'Bank Mandiri', 'flat', 5000, 'Transfer to Mandiri', 'Domestic SKNT'),
    ('BRI', 'Bank BRI', 'flat', 5000, 'Transfer to BRI', 'Domestic SKNT'),
    ('CIMB', 'Bank CIMB Niaga', 'flat', 5000, 'Transfer to CIMB', 'Domestic SKNT'),
    ('DANAMON', 'Bank Danamon', 'flat', 5000, 'Transfer to Danamon', 'Domestic SKNT'),
    ('PERMATA', 'Bank Permata', 'flat', 5000, 'Transfer to Permata', 'Domestic SKNT'),
    ('OCBC', 'Bank OCBC NISP', 'flat', 5000, 'Transfer to OCBC', 'Domestic SKNT'),
    ('MEGA', 'Bank Mega', 'flat', 5000, 'Transfer to Mega', 'Domestic SKNT'),
    ('UOB', 'Bank UOB Indonesia', 'flat', 5000, 'Transfer to UOB', 'Domestic SKNT'),
    ('PANIN', 'Bank Panin', 'flat', 5000, 'Transfer to Panin', 'Domestic SKNT'),
    ('COMMONWEALTH', 'Bank Commonwealth', 'flat', 5000, 'Transfer to Commonwealth', 'Domestic SKNT'),
    ('MAYBANK', 'Maybank Indonesia', 'flat', 5000, 'Transfer to Maybank', 'Domestic SKNT'),
    ('BTN', 'Bank BTN', 'flat', 5000, 'Transfer to BTN', 'Domestic SKNT'),
    ('SUMITOMO', 'Bank BTMU', 'flat', 5000, 'Transfer to BTMU', 'Domestic SKNT'),
    ('DBS', 'DBS Bank', 'flat', 5000, 'Transfer to DBS', 'Domestic SKNT'),
    
    -- International Banks - Higher Fee (Rp 10,000 - SWIFT)
    ('CITIBANK', 'Citibank Indonesia', 'flat', 10000, 'Transfer to Citibank', 'SWIFT International'),
    ('HSBC', 'HSBC Bank Indonesia', 'flat', 10000, 'Transfer to HSBC', 'SWIFT International'),
    ('MIZUHO', 'Bank Mizuho Indonesia', 'flat', 10000, 'Transfer to Mizuho', 'SWIFT International'),
    ('JPMORGAN', 'JP Morgan Chase', 'flat', 10000, 'Transfer to JP Morgan', 'SWIFT International');

-- Service Fees: Fees for other services (e-wallet, e-money, VA, QRIS)
INSERT OR IGNORE INTO service_fees 
    (service_code, service_name, service_type, fee_type, fee_amount, notes) 
VALUES
    ('TOPUP_OVO', 'Top Up OVO', 'topup_ewallet', 'flat', 2500, 'OVO e-wallet top up'),
    ('TOPUP_DANA', 'Top Up DANA', 'topup_ewallet', 'flat', 2500, 'DANA e-wallet top up'),
    ('TOPUP_GOPAY', 'Top Up GoPay', 'topup_ewallet', 'flat', 2500, 'GoPay e-wallet top up'),
    ('TOPUP_LINKAJA', 'Top Up LinkAja', 'topup_emoney', 'flat', 2500, 'LinkAja e-money top up'),
    ('TOPUP_MANDIRIEMONEY', 'Top Up Mandiri e-Money', 'topup_emoney', 'flat', 2500, 'Mandiri e-money top up'),
    ('PAYMENT_VA_MANDIRI', 'Payment Virtual Account Mandiri', 'payment_va', 'flat', 0, 'Virtual Account payment - Mandiri'),
    ('PAYMENT_VA_BCA', 'Payment Virtual Account BCA', 'payment_va', 'flat', 0, 'Virtual Account payment - BCA'),
    ('PAYMENT_VA_BRI', 'Payment Virtual Account BRI', 'payment_va', 'flat', 0, 'Virtual Account payment - BRI'),
    ('QRIS_PAYMENT', 'QRIS Payment', 'qris_payment', 'percentage', 1, 'QRIS Payment - 1% charge');
