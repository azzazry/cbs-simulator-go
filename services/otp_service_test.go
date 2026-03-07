package services_test

import (
	"testing"

	"cbs-simulator/services"
)

func TestGenerateOTP(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	otp, err := services.GenerateOTP("CIF001", "unlock_account", "sms")
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	if len(otp) != 6 {
		t.Errorf("OTP should be 6 digits, got %d characters: '%s'", len(otp), otp)
	}

	// OTP should be all digits
	for _, c := range otp {
		if c < '0' || c > '9' {
			t.Errorf("OTP should contain only digits, got '%s'", otp)
			break
		}
	}
}

func TestVerifyOTP_Valid(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// Generate OTP
	otp, _ := services.GenerateOTP("CIF001", "unlock_account", "sms")

	// Verify it
	err := services.VerifyOTP("CIF001", otp, "unlock_account")
	if err != nil {
		t.Errorf("OTP verification should succeed: %v", err)
	}
}

func TestVerifyOTP_WrongCode(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// Generate OTP
	services.GenerateOTP("CIF001", "unlock_account", "sms")

	// Verify with wrong code
	err := services.VerifyOTP("CIF001", "000000", "unlock_account")
	if err == nil {
		t.Error("Should fail with wrong OTP code")
	}
}

func TestVerifyOTP_AlreadyUsed(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// Generate and verify OTP
	otp, _ := services.GenerateOTP("CIF001", "unlock_account", "sms")
	services.VerifyOTP("CIF001", otp, "unlock_account")

	// Try to use again
	err := services.VerifyOTP("CIF001", otp, "unlock_account")
	if err == nil {
		t.Error("Should fail: OTP already used")
	}
}

func TestVerifyOTP_WrongType(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// Generate OTP for unlock
	otp, _ := services.GenerateOTP("CIF001", "unlock_account", "sms")

	// Try to verify as reset_pin
	err := services.VerifyOTP("CIF001", otp, "reset_pin")
	if err == nil {
		t.Error("Should fail: wrong OTP type")
	}
}

func TestVerifyOTP_WrongCIF(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// Generate OTP for CIF001
	otp, _ := services.GenerateOTP("CIF001", "unlock_account", "sms")

	// Try to verify as CIF002
	err := services.VerifyOTP("CIF002", otp, "unlock_account")
	if err == nil {
		t.Error("Should fail: wrong CIF")
	}
}
