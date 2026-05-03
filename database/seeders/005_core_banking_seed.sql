-- Phase 2: Core Banking Seed Data
-- Chart of Accounts + Interest Rates + CIF Extended

-- =============================================
-- CHART OF ACCOUNTS (Kode 3-digit standar perbankan Indonesia)
-- Mengacu pada PAPI/PSAK klasifikasi
-- Level: 1=Header, 2=Sub-header, 3=Detail
-- =============================================

-- ASET (1xx)
INSERT INTO chart_of_accounts (account_code, account_name, account_type, parent_code, level, normal_balance) VALUES
('100', 'ASET', 'asset', NULL, 1, 'debit'),
('110', 'Kas dan Setara Kas', 'asset', '100', 2, 'debit'),
('111', 'Kas', 'asset', '110', 3, 'debit'),
('112', 'Giro Pada Bank Indonesia', 'asset', '110', 3, 'debit'),
('113', 'Giro Pada Bank Lain', 'asset', '110', 3, 'debit'),
('114', 'Penempatan Pada Bank Lain', 'asset', '110', 3, 'debit'),
('120', 'Kredit Yang Diberikan', 'asset', '100', 2, 'debit'),
('121', 'Kredit Modal Kerja', 'asset', '120', 3, 'debit'),
('122', 'Kredit Investasi', 'asset', '120', 3, 'debit'),
('123', 'Kredit Konsumsi (KPR)', 'asset', '120', 3, 'debit'),
('124', 'Kredit Mikro/UMKM', 'asset', '120', 3, 'debit'),
('125', 'Cadangan Kerugian Penurunan Nilai (CKPN)', 'asset', '120', 3, 'debit'),
('130', 'Surat Berharga', 'asset', '100', 2, 'debit'),
('131', 'Surat Berharga Negara (SBN)', 'asset', '130', 3, 'debit'),
('132', 'Obligasi Korporasi', 'asset', '130', 3, 'debit'),
('140', 'Aset Tetap', 'asset', '100', 2, 'debit'),
('141', 'Tanah dan Bangunan', 'asset', '140', 3, 'debit'),
('142', 'Peralatan dan Inventaris', 'asset', '140', 3, 'debit'),
('143', 'Akumulasi Penyusutan', 'asset', '140', 3, 'debit'),
('150', 'Aset Lain-lain', 'asset', '100', 2, 'debit'),
('151', 'Bunga Yang Akan Diterima', 'asset', '150', 3, 'debit'),
('152', 'Biaya Dibayar Dimuka', 'asset', '150', 3, 'debit');

-- KEWAJIBAN / LIABILITAS (2xx)
INSERT INTO chart_of_accounts (account_code, account_name, account_type, parent_code, level, normal_balance) VALUES
('200', 'KEWAJIBAN', 'liability', NULL, 1, 'credit'),
('210', 'Dana Pihak Ketiga (DPK)', 'liability', '200', 2, 'credit'),
('211', 'Tabungan', 'liability', '210', 3, 'credit'),
('212', 'Giro Nasabah', 'liability', '210', 3, 'credit'),
('213', 'Deposito Berjangka', 'liability', '210', 3, 'credit'),
('220', 'Kewajiban Segera', 'liability', '200', 2, 'credit'),
('221', 'Kewajiban Kepada Bank Lain', 'liability', '220', 3, 'credit'),
('222', 'Hutang Pajak', 'liability', '220', 3, 'credit'),
('223', 'Hutang Bunga', 'liability', '220', 3, 'credit'),
('230', 'Kewajiban Lain-lain', 'liability', '200', 2, 'credit'),
('231', 'Pendapatan Diterima Dimuka', 'liability', '230', 3, 'credit');

-- EKUITAS (3xx)
INSERT INTO chart_of_accounts (account_code, account_name, account_type, parent_code, level, normal_balance) VALUES
('300', 'EKUITAS', 'equity', NULL, 1, 'credit'),
('310', 'Modal', 'equity', '300', 2, 'credit'),
('311', 'Modal Disetor', 'equity', '310', 3, 'credit'),
('312', 'Tambahan Modal Disetor', 'equity', '310', 3, 'credit'),
('320', 'Laba', 'equity', '300', 2, 'credit'),
('321', 'Laba Ditahan', 'equity', '320', 3, 'credit'),
('322', 'Laba Tahun Berjalan', 'equity', '320', 3, 'credit');

-- PENDAPATAN (4xx)
INSERT INTO chart_of_accounts (account_code, account_name, account_type, parent_code, level, normal_balance) VALUES
('400', 'PENDAPATAN', 'revenue', NULL, 1, 'credit'),
('410', 'Pendapatan Bunga', 'revenue', '400', 2, 'credit'),
('411', 'Pendapatan Bunga Kredit', 'revenue', '410', 3, 'credit'),
('412', 'Pendapatan Bunga Penempatan', 'revenue', '410', 3, 'credit'),
('413', 'Pendapatan Bunga Surat Berharga', 'revenue', '410', 3, 'credit'),
('420', 'Pendapatan Operasional Lainnya', 'revenue', '400', 2, 'credit'),
('421', 'Pendapatan Provisi dan Komisi', 'revenue', '420', 3, 'credit'),
('422', 'Pendapatan Fee Transaksi', 'revenue', '420', 3, 'credit'),
('423', 'Pendapatan Fee Transfer', 'revenue', '420', 3, 'credit'),
('424', 'Pendapatan Fee Administrasi', 'revenue', '420', 3, 'credit'),
('430', 'Pendapatan Non-Operasional', 'revenue', '400', 2, 'credit'),
('431', 'Keuntungan Penjualan Aset', 'revenue', '430', 3, 'credit'),
('432', 'Pendapatan Lain-lain', 'revenue', '430', 3, 'credit');

-- BEBAN (5xx)
INSERT INTO chart_of_accounts (account_code, account_name, account_type, parent_code, level, normal_balance) VALUES
('500', 'BEBAN', 'expense', NULL, 1, 'debit'),
('510', 'Beban Bunga', 'expense', '500', 2, 'debit'),
('511', 'Beban Bunga Tabungan', 'expense', '510', 3, 'debit'),
('512', 'Beban Bunga Giro', 'expense', '510', 3, 'debit'),
('513', 'Beban Bunga Deposito', 'expense', '510', 3, 'debit'),
('514', 'Beban Bunga Pinjaman', 'expense', '510', 3, 'debit'),
('520', 'Beban Operasional', 'expense', '500', 2, 'debit'),
('521', 'Beban Gaji dan Tunjangan', 'expense', '520', 3, 'debit'),
('522', 'Beban Sewa', 'expense', '520', 3, 'debit'),
('523', 'Beban Penyusutan', 'expense', '520', 3, 'debit'),
('524', 'Beban Umum dan Administrasi', 'expense', '520', 3, 'debit'),
('525', 'Beban CKPN (Provisi Kerugian)', 'expense', '520', 3, 'debit'),
('530', 'Beban Non-Operasional', 'expense', '500', 2, 'debit'),
('531', 'Beban Pajak', 'expense', '530', 3, 'debit'),
('532', 'Beban Lain-lain', 'expense', '530', 3, 'debit');

-- =============================================
-- INTEREST RATES (Suku Bunga)
-- =============================================

-- Tabungan (Savings)
INSERT INTO interest_rates (product_type, product_name, rate_type, base_rate, min_balance, max_balance, effective_date) VALUES
('savings', 'Tabungan Reguler', 'tiered', 0.50, 0, 10000000, '2026-01-01'),
('savings', 'Tabungan Reguler', 'tiered', 1.00, 10000000, 100000000, '2026-01-01'),
('savings', 'Tabungan Reguler', 'tiered', 1.50, 100000000, 1000000000, '2026-01-01'),
('savings', 'Tabungan Reguler', 'tiered', 2.00, 1000000000, NULL, '2026-01-01');

-- Giro (Checking)
INSERT INTO interest_rates (product_type, product_name, rate_type, base_rate, min_balance, effective_date) VALUES
('checking', 'Giro Rupiah', 'fixed', 0.50, 0, '2026-01-01')
ON CONFLICT DO NOTHING;

-- Deposito (Time Deposit)
INSERT INTO interest_rates (product_type, product_name, rate_type, base_rate, tenor_months, effective_date) VALUES
('deposit', 'Deposito 1 Bulan', 'fixed', 3.00, 1, '2026-01-01'),
('deposit', 'Deposito 3 Bulan', 'fixed', 3.50, 3, '2026-01-01'),
('deposit', 'Deposito 6 Bulan', 'fixed', 4.00, 6, '2026-01-01'),
('deposit', 'Deposito 12 Bulan', 'fixed', 4.50, 12, '2026-01-01'),
('deposit', 'Deposito 24 Bulan', 'fixed', 5.00, 24, '2026-01-01')
ON CONFLICT DO NOTHING;

-- Kredit (Loan)
INSERT INTO interest_rates (product_type, product_name, rate_type, base_rate, effective_date) VALUES
('loan', 'KPR Fixed 5 Tahun', 'fixed', 7.50, '2026-01-01'),
('loan', 'Kredit Modal Kerja', 'floating', 9.00, '2026-01-01'),
('loan', 'Kredit Konsumsi', 'fixed', 10.50, '2026-01-01'),
('loan', 'Kredit Mikro/UMKM', 'fixed', 6.00, '2026-01-01'),
('loan', 'Kredit Investasi', 'floating', 8.50, '2026-01-01')
ON CONFLICT DO NOTHING;

-- =============================================
-- CIF EXTENDED DATA FOR EXISTING CUSTOMERS
-- =============================================
INSERT INTO customer_extended (cif, mother_maiden_name, nationality, occupation, employer_name, monthly_income, source_of_funds, risk_profile, segment, branch_code, npwp) VALUES
('CIF001', 'Sari Dewi', 'WNI', 'Software Engineer', 'PT Telkom Indonesia', 25000000, 'Gaji', 'low', 'affluent', 'JKT001', '12.345.678.9-012.000'),
('CIF002', 'Aminah', 'WNI', 'Dokter', 'RS Pondok Indah', 35000000, 'Gaji', 'low', 'priority', 'JKT002', '23.456.789.0-123.000'),
('CIF003', 'Fatimah', 'WNI', 'Pengusaha', 'CV Wijaya Mandiri', 50000000, 'Usaha', 'medium', 'priority', 'JKT003', '34.567.890.1-234.000'),
('CIF004', 'Kartini', 'WNI', 'Guru', 'SDN 01 Jakarta', 8000000, 'Gaji', 'low', 'mass', 'JKT004', '45.678.901.2-345.000'),
('CIF005', 'Ratna', 'WNI', 'Wiraswasta', 'Toko Pratama', 15000000, 'Usaha', 'medium', 'mass', 'JKT005', '56.789.012.3-456.000')
ON CONFLICT DO NOTHING;
