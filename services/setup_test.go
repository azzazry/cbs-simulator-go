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
