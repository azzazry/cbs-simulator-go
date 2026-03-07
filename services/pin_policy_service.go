package services

import (
	"cbs-simulator/config"
	"cbs-simulator/database"
	"fmt"
	"time"
)

// ValidatePINPolicy validates a PIN against security policy
func ValidatePINPolicy(pin string) error {
	minLen := config.AppConfig.PINMinLength
	maxLen := config.AppConfig.PINMaxLength

	// Check length
	if len(pin) < minLen || len(pin) > maxLen {
		return fmt.Errorf("PIN must be %d digits", minLen)
	}

	// Check all digits
	for _, c := range pin {
		if c < '0' || c > '9' {
			return fmt.Errorf("PIN must contain only digits")
		}
	}

	// Check sequential (ascending: 123456, descending: 654321)
	isAscending := true
	isDescending := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			isAscending = false
		}
		if pin[i] != pin[i-1]-1 {
			isDescending = false
		}
	}
	if isAscending || isDescending {
		return fmt.Errorf("PIN cannot be sequential numbers")
	}

	// Check all same digits (111111, 222222)
	allSame := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return fmt.Errorf("PIN cannot be all the same digit")
	}

	return nil
}

// RecordLoginAttempt records a login attempt
func RecordLoginAttempt(cif, ip string, success bool) error {
	isSuccess := 0
	if success {
		isSuccess = 1
	}
	query := `INSERT INTO login_attempts (cif, ip_address, attempt_type, is_success) VALUES (?, ?, 'pin', ?)`
	_, err := database.DB.Exec(query, cif, ip, isSuccess)
	return err
}

// IsAccountLocked checks if account is locked due to failed login attempts
func IsAccountLocked(cif string) (bool, error) {
	maxAttempts := config.AppConfig.MaxLoginAttempts
	lockoutMinutes := config.AppConfig.LockoutDurationMinutes

	// Use UTC to match SQLite's CURRENT_TIMESTAMP format
	windowStart := time.Now().UTC().Add(-time.Duration(lockoutMinutes) * time.Minute).Format("2006-01-02 15:04:05")

	// Count failed attempts in the lockout window since last success
	query := `SELECT COUNT(*) FROM login_attempts 
	          WHERE cif = ? AND is_success = 0 AND attempted_at > ?
	          AND attempted_at > COALESCE(
	              (SELECT MAX(attempted_at) FROM login_attempts WHERE cif = ? AND is_success = 1),
	              '2000-01-01'
	          )`
	var failedCount int
	err := database.DB.QueryRow(query, cif, windowStart, cif).Scan(&failedCount)
	if err != nil {
		return false, err
	}

	return failedCount >= maxAttempts, nil
}

// GetFailedAttemptCount returns the number of recent failed login attempts
func GetFailedAttemptCount(cif string) (int, error) {
	lockoutMinutes := config.AppConfig.LockoutDurationMinutes
	// Use UTC to match SQLite's CURRENT_TIMESTAMP format
	windowStart := time.Now().UTC().Add(-time.Duration(lockoutMinutes) * time.Minute).Format("2006-01-02 15:04:05")

	query := `SELECT COUNT(*) FROM login_attempts 
	          WHERE cif = ? AND is_success = 0 AND attempted_at > ?
	          AND attempted_at > COALESCE(
	              (SELECT MAX(attempted_at) FROM login_attempts WHERE cif = ? AND is_success = 1),
	              '2000-01-01'
	          )`
	var count int
	err := database.DB.QueryRow(query, cif, windowStart, cif).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UnlockAccount resets the lockout by recording a successful "unlock" attempt
func UnlockAccount(cif string) error {
	// Insert a successful attempt to reset the counter
	query := `INSERT INTO login_attempts (cif, ip_address, attempt_type, is_success) VALUES (?, 'system', 'unlock', 1)`
	_, err := database.DB.Exec(query, cif)
	return err
}
