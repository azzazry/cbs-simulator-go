package services_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"cbs-simulator/config"
	"cbs-simulator/database"
	"cbs-simulator/utils"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database with security tables for testing
func setupTestDB(t *testing.T) {
	t.Helper()

	var err error
	database.DB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Create customers table
	_, err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cif VARCHAR(20) UNIQUE NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			id_card_number VARCHAR(20) UNIQUE NOT NULL,
			phone_number VARCHAR(20) NOT NULL,
			email VARCHAR(100),
			address TEXT,
			date_of_birth DATE,
			pin VARCHAR(255) NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create customers table: %v", err)
	}

	// Create security tables
	securitySQL := `
		CREATE TABLE IF NOT EXISTS token_blacklist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token_jti VARCHAR(50) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			expires_at DATETIME NOT NULL,
			blacklisted_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_name VARCHAR(50) UNIQUE NOT NULL,
			description TEXT,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cif VARCHAR(20) NOT NULL,
			role_id INTEGER NOT NULL,
			assigned_by VARCHAR(20),
			assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(cif, role_id)
		);

		CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cif VARCHAR(20),
			action VARCHAR(100) NOT NULL,
			resource VARCHAR(100),
			resource_id VARCHAR(100),
			ip_address VARCHAR(45),
			user_agent TEXT,
			request_method VARCHAR(10),
			request_path VARCHAR(255),
			request_body TEXT,
			response_status INTEGER,
			details TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS transaction_limits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_name VARCHAR(50) NOT NULL,
			transaction_type VARCHAR(50) NOT NULL,
			daily_limit DECIMAL(18,2) DEFAULT 0,
			per_transaction_limit DECIMAL(18,2) DEFAULT 0,
			monthly_limit DECIMAL(18,2) DEFAULT 0,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(role_name, transaction_type)
		);

		CREATE TABLE IF NOT EXISTS login_attempts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cif VARCHAR(20) NOT NULL,
			ip_address VARCHAR(45),
			attempt_type VARCHAR(30) DEFAULT 'pin',
			is_success INTEGER DEFAULT 0,
			attempted_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS otp_codes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cif VARCHAR(20) NOT NULL,
			otp_code VARCHAR(10) NOT NULL,
			otp_type VARCHAR(30) NOT NULL,
			channel VARCHAR(20) DEFAULT 'sms',
			is_used INTEGER DEFAULT 0,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS ekyc_verifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verification_id VARCHAR(50) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			id_card_number VARCHAR(20) NOT NULL,
			verification_type VARCHAR(30) NOT NULL,
			verification_status VARCHAR(20) DEFAULT 'pending',
			verified_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			transaction_id VARCHAR(50) UNIQUE NOT NULL,
			transaction_type VARCHAR(30) NOT NULL,
			from_account_number VARCHAR(20),
			to_account_number VARCHAR(20),
			amount DECIMAL(18,2) NOT NULL,
			currency VARCHAR(3) DEFAULT 'IDR',
			description TEXT,
			reference_number VARCHAR(50),
			status VARCHAR(20) DEFAULT 'pending',
			transaction_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			settlement_date DATE,
			fee DECIMAL(18,2) DEFAULT 0
		);
	`
	_, err = database.DB.Exec(securitySQL)
	if err != nil {
		t.Fatalf("Failed to create security tables: %v", err)
	}

	// Create Phase 2: Core Banking tables
	coreBankingSQL := `
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_number VARCHAR(20) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			account_type VARCHAR(20) NOT NULL,
			currency VARCHAR(3) DEFAULT 'IDR',
			balance DECIMAL(18,2) DEFAULT 0.00,
			avail_balance DECIMAL(18,2) DEFAULT 0.00,
			status VARCHAR(20) DEFAULT 'active',
			opened_date DATE NOT NULL,
			branch VARCHAR(50),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS chart_of_accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_code VARCHAR(10) UNIQUE NOT NULL,
			account_name VARCHAR(100) NOT NULL,
			account_type VARCHAR(20) NOT NULL,
			parent_code VARCHAR(10),
			level INTEGER DEFAULT 1,
			normal_balance VARCHAR(10) NOT NULL,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS journal_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			journal_number VARCHAR(30) UNIQUE NOT NULL,
			entry_date DATE NOT NULL,
			description TEXT,
			reference_type VARCHAR(30),
			reference_id VARCHAR(50),
			posted_by VARCHAR(20),
			status VARCHAR(20) DEFAULT 'posted',
			reversed_by INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS journal_lines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			journal_id INTEGER NOT NULL,
			account_code VARCHAR(10) NOT NULL,
			debit_amount DECIMAL(18,2) DEFAULT 0,
			credit_amount DECIMAL(18,2) DEFAULT 0,
			description TEXT
		);

		CREATE TABLE IF NOT EXISTS customer_extended (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cif VARCHAR(20) UNIQUE NOT NULL,
			mother_maiden_name VARCHAR(100),
			nationality VARCHAR(50) DEFAULT 'WNI',
			occupation VARCHAR(100),
			employer_name VARCHAR(100),
			monthly_income DECIMAL(18,2),
			source_of_funds VARCHAR(100),
			risk_profile VARCHAR(20) DEFAULT 'low',
			segment VARCHAR(30) DEFAULT 'mass',
			branch_code VARCHAR(20),
			rm_code VARCHAR(20),
			npwp VARCHAR(20),
			last_kyc_date DATE,
			next_kyc_date DATE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS interest_rates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_type VARCHAR(30) NOT NULL,
			product_name VARCHAR(50),
			rate_type VARCHAR(20) NOT NULL,
			base_rate DECIMAL(8,4) NOT NULL,
			min_balance DECIMAL(18,2) DEFAULT 0,
			max_balance DECIMAL(18,2),
			tenor_months INTEGER,
			effective_date DATE NOT NULL,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS interest_accruals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_number VARCHAR(20) NOT NULL,
			accrual_date DATE NOT NULL,
			product_type VARCHAR(30) NOT NULL,
			balance DECIMAL(18,2) NOT NULL,
			rate DECIMAL(8,4) NOT NULL,
			daily_interest DECIMAL(18,4) NOT NULL,
			accrued_interest DECIMAL(18,4) NOT NULL,
			is_posted INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(account_number, accrual_date)
		);

		CREATE TABLE IF NOT EXISTS standing_instructions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			si_number VARCHAR(30) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			from_account VARCHAR(20) NOT NULL,
			instruction_type VARCHAR(30) NOT NULL,
			to_account VARCHAR(20),
			to_bank_code VARCHAR(20),
			amount DECIMAL(18,2) NOT NULL,
			description TEXT,
			frequency VARCHAR(20) NOT NULL,
			execution_day INTEGER,
			start_date DATE NOT NULL,
			end_date DATE,
			next_execution_date DATE NOT NULL,
			total_executed INTEGER DEFAULT 0,
			total_failed INTEGER DEFAULT 0,
			last_execution_date DATE,
			last_status VARCHAR(20),
			status VARCHAR(20) DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS si_executions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			si_number VARCHAR(30) NOT NULL,
			execution_date DATE NOT NULL,
			amount DECIMAL(18,2) NOT NULL,
			transaction_id VARCHAR(50),
			status VARCHAR(20) NOT NULL,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS eod_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			process_date DATE NOT NULL,
			process_type VARCHAR(30) NOT NULL,
			status VARCHAR(20) NOT NULL,
			records_processed INTEGER DEFAULT 0,
			records_failed INTEGER DEFAULT 0,
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			card_number VARCHAR(16) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			account_number VARCHAR(20) NOT NULL,
			card_type VARCHAR(20) NOT NULL,
			card_brand VARCHAR(20) NOT NULL,
			card_limit DECIMAL(18,2) DEFAULT 0.00,
			avail_limit DECIMAL(18,2) DEFAULT 0.00,
			expiry_date VARCHAR(7) NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			cvv VARCHAR(3) NOT NULL,
			pin VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS loans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			loan_number VARCHAR(20) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			account_number VARCHAR(20) NOT NULL,
			loan_type VARCHAR(30) NOT NULL,
			principal_amount DECIMAL(18,2) NOT NULL,
			outstanding_amount DECIMAL(18,2) NOT NULL,
			interest_rate DECIMAL(5,2) NOT NULL,
			monthly_payment DECIMAL(18,2) NOT NULL,
			tenor_months INTEGER NOT NULL,
			remaining_months INTEGER NOT NULL,
			disbursement_date DATE NOT NULL,
			maturity_date DATE NOT NULL,
			next_payment_date DATE NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS deposits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			deposit_number VARCHAR(20) UNIQUE NOT NULL,
			cif VARCHAR(20) NOT NULL,
			principal_amount DECIMAL(18,2) NOT NULL,
			interest_rate DECIMAL(5,2) NOT NULL,
			tenor_months INTEGER NOT NULL,
			open_date DATE NOT NULL,
			maturity_date DATE NOT NULL,
			maturity_amount DECIMAL(18,2) NOT NULL,
			auto_renew BOOLEAN DEFAULT 0,
			status VARCHAR(20) DEFAULT 'active',
			linked_account VARCHAR(20),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = database.DB.Exec(coreBankingSQL)
	if err != nil {
		t.Fatalf("Failed to create core banking tables: %v", err)
	}

	// Seed test data
	seedTestData(t)
}

func seedTestData(t *testing.T) {
	t.Helper()

	// Hash PIN "123456"
	hashedPIN, err := utils.HashPIN("123456")
	if err != nil {
		t.Fatalf("Failed to hash PIN: %v", err)
	}

	// Insert test customers
	_, err = database.DB.Exec(`
		INSERT INTO customers (cif, full_name, id_card_number, phone_number, email, address, date_of_birth, pin, status) VALUES
		('CIF001', 'Budi Santoso', '3201011234567890', '081234567890', 'budi@email.com', 'Jakarta', '1985-03-15', ?, 'active'),
		('CIF002', 'Siti Nurhaliza', '3201021234567891', '081234567891', 'siti@email.com', 'Jakarta', '1990-07-20', ?, 'active'),
		('CIF003', 'Ahmad Wijaya', '3201031234567892', '081234567892', 'ahmad@email.com', 'Jakarta', '1988-11-10', ?, 'active')
	`, hashedPIN, hashedPIN, hashedPIN)
	if err != nil {
		t.Fatalf("Failed to seed customers: %v", err)
	}

	// Seed roles
	_, err = database.DB.Exec(`
		INSERT INTO roles (role_name, description) VALUES 
		('customer', 'Regular customer'),
		('teller', 'Bank teller'),
		('admin', 'System administrator'),
		('supervisor', 'Branch supervisor')
	`)
	if err != nil {
		t.Fatalf("Failed to seed roles: %v", err)
	}

	// Assign roles
	_, err = database.DB.Exec(`
		INSERT INTO user_roles (cif, role_id, assigned_by) VALUES
		('CIF001', 1, 'SYSTEM'),
		('CIF001', 3, 'SYSTEM'),
		('CIF002', 1, 'SYSTEM'),
		('CIF003', 1, 'SYSTEM'),
		('CIF003', 4, 'SYSTEM')
	`)
	if err != nil {
		t.Fatalf("Failed to seed user roles: %v", err)
	}

	// Seed transaction limits
	_, err = database.DB.Exec(`
		INSERT INTO transaction_limits (role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit) VALUES
		('customer', 'intra_transfer', 50000000, 25000000, 500000000),
		('customer', 'inter_transfer', 25000000, 10000000, 250000000)
	`)
	if err != nil {
		t.Fatalf("Failed to seed transaction limits: %v", err)
	}

	// Phase 2: Seed accounts
	_, err = database.DB.Exec(`
		INSERT INTO accounts (account_number, cif, account_type, currency, balance, avail_balance, status, opened_date, branch) VALUES
		('1000000001', 'CIF001', 'savings', 'IDR', 50000000.00, 50000000.00, 'active', '2025-01-01', 'JKT001'),
		('1000000002', 'CIF002', 'savings', 'IDR', 25000000.00, 25000000.00, 'active', '2025-01-01', 'JKT002'),
		('1000000003', 'CIF003', 'savings', 'IDR', 100000000.00, 100000000.00, 'active', '2025-03-01', 'JKT003')
	`)
	if err != nil {
		t.Fatalf("Failed to seed accounts: %v", err)
	}

	// Phase 2: Seed Chart of Accounts
	_, err = database.DB.Exec(`
		INSERT INTO chart_of_accounts (account_code, account_name, account_type, parent_code, level, normal_balance) VALUES
		('100', 'ASET', 'asset', NULL, 1, 'debit'),
		('110', 'Kas dan Setara Kas', 'asset', '100', 2, 'debit'),
		('111', 'Kas', 'asset', '110', 3, 'debit'),
		('200', 'KEWAJIBAN', 'liability', NULL, 1, 'credit'),
		('210', 'Dana Pihak Ketiga', 'liability', '200', 2, 'credit'),
		('211', 'Tabungan', 'liability', '210', 3, 'credit'),
		('400', 'PENDAPATAN', 'revenue', NULL, 1, 'credit'),
		('422', 'Pendapatan Fee Transaksi', 'revenue', '400', 3, 'credit'),
		('500', 'BEBAN', 'expense', NULL, 1, 'debit'),
		('511', 'Beban Bunga Tabungan', 'expense', '500', 3, 'debit')
	`)
	if err != nil {
		t.Fatalf("Failed to seed CoA: %v", err)
	}

	// Phase 2: Seed interest rates
	_, err = database.DB.Exec(`
		INSERT INTO interest_rates (product_type, product_name, rate_type, base_rate, min_balance, max_balance, effective_date) VALUES
		('savings', 'Tabungan Reguler', 'tiered', 1.00, 0, 100000000, '2026-01-01'),
		('savings', 'Tabungan Reguler', 'tiered', 2.00, 100000000, NULL, '2026-01-01');
		INSERT INTO interest_rates (product_type, product_name, rate_type, base_rate, tenor_months, effective_date) VALUES
		('deposit', 'Deposito 12 Bulan', 'fixed', 4.50, 12, '2026-01-01')
	`)
	if err != nil {
		t.Fatalf("Failed to seed interest rates: %v", err)
	}

	// Phase 2: Seed deposits & loans for testing
	_, err = database.DB.Exec(`
		INSERT INTO deposits (deposit_number, cif, principal_amount, interest_rate, tenor_months, open_date, maturity_date, maturity_amount, status, linked_account) VALUES
		('DEP001', 'CIF001', 100000000, 4.50, 12, '2025-01-01', '2026-01-01', 104500000, 'active', '1000000001')
	`)
	if err != nil {
		t.Fatalf("Failed to seed deposits: %v", err)
	}

	_, err = database.DB.Exec(`
		INSERT INTO loans (loan_number, cif, account_number, loan_type, principal_amount, outstanding_amount, interest_rate, monthly_payment, tenor_months, remaining_months, disbursement_date, maturity_date, next_payment_date, status) VALUES
		('LOAN001', 'CIF001', '1000000001', 'kpr', 500000000, 450000000, 7.50, 5800000, 120, 108, '2025-01-01', '2035-01-01', '2026-04-01', 'active')
	`)
	if err != nil {
		t.Fatalf("Failed to seed loans: %v", err)
	}
}

func teardownTestDB(t *testing.T) {
	t.Helper()
	if database.DB != nil {
		database.DB.Close()
	}
}

// ensureConfig ensures config is loaded for tests
func ensureConfig() {
	if config.AppConfig == nil {
		os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")
		os.Setenv("MAX_LOGIN_ATTEMPTS", "3")
		os.Setenv("PIN_MIN_LENGTH", "6")
		os.Setenv("PIN_MAX_LENGTH", "6")
		os.Setenv("OTP_EXPIRY_MINUTES", "5")
		os.Setenv("OTP_LENGTH", "6")
		os.Setenv("LOCKOUT_DURATION_MINUTES", "30")
		config.AppConfig = &config.Config{
			JWTSecret:              "test-secret-key-for-unit-tests",
			JWTAccessExpiry:        15,
			JWTRefreshExpiry:       168,
			RateLimitPerMinute:     60,
			MaxLoginAttempts:       3,
			LockoutDurationMinutes: 30,
			PINMinLength:           6,
			PINMaxLength:           6,
			OTPExpiry:              5,
			OTPLength:              6,
		}
		log.Println("Test config initialized")
	}
}
