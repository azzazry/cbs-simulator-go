package services

import (
	"cbs-simulator/config"
	"cbs-simulator/database"
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP creates and stores a new OTP code
// In production, this would send via SMS/email. For simulator, OTP is returned in response.
func GenerateOTP(cif, otpType, channel string) (string, error) {
	// Invalidate any existing active OTPs for this CIF and type
	InvalidateOTP(cif, otpType)

	// Generate random OTP
	otpLength := config.AppConfig.OTPLength
	otp := generateRandomOTP(otpLength)

	// Calculate expiry
	expiryMinutes := config.AppConfig.OTPExpiry
	expiresAt := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	// Store in database
	query := `INSERT INTO otp_codes (cif, otp_code, otp_type, channel, expires_at) VALUES (?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, cif, otp, otpType, channel, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to store OTP: %v", err)
	}

	// In production: send via SMS gateway or email service
	// For simulator: return the OTP directly (shown in response)
	return otp, nil
}

// VerifyOTP validates an OTP code
func VerifyOTP(cif, otp, otpType string) error {
	var id int
	var expiresAt time.Time
	var isUsed int

	query := `SELECT id, expires_at, is_used FROM otp_codes 
	          WHERE cif = ? AND otp_code = ? AND otp_type = ? 
	          ORDER BY created_at DESC LIMIT 1`
	err := database.DB.QueryRow(query, cif, otp, otpType).Scan(&id, &expiresAt, &isUsed)
	if err != nil {
		return fmt.Errorf("invalid OTP code")
	}

	if isUsed == 1 {
		return fmt.Errorf("OTP has already been used")
	}

	if time.Now().After(expiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	// Mark as used
	_, err = database.DB.Exec(`UPDATE otp_codes SET is_used = 1 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %v", err)
	}

	return nil
}

// InvalidateOTP invalidates all active OTPs for a CIF and type
func InvalidateOTP(cif, otpType string) {
	database.DB.Exec(`UPDATE otp_codes SET is_used = 1 WHERE cif = ? AND otp_type = ? AND is_used = 0`, cif, otpType)
}

// generateRandomOTP generates a random numeric OTP of specified length
func generateRandomOTP(length int) string {
	digits := "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}
