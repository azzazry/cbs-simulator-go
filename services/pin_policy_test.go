package services_test

import (
	"testing"

	"cbs-simulator/services"
)

func TestValidatePINPolicy_ValidPIN(t *testing.T) {
	ensureConfig()

	validPINs := []string{"246813", "975312", "482931", "371592"}
	for _, pin := range validPINs {
		err := services.ValidatePINPolicy(pin)
		if err != nil {
			t.Errorf("PIN %s should be valid, got error: %v", pin, err)
		}
	}
}

func TestValidatePINPolicy_TooShort(t *testing.T) {
	ensureConfig()

	err := services.ValidatePINPolicy("123")
	if err == nil {
		t.Error("PIN '123' should fail: too short")
	}
}

func TestValidatePINPolicy_TooLong(t *testing.T) {
	ensureConfig()

	err := services.ValidatePINPolicy("12345678")
	if err == nil {
		t.Error("PIN '12345678' should fail: too long")
	}
}

func TestValidatePINPolicy_NonDigit(t *testing.T) {
	ensureConfig()

	err := services.ValidatePINPolicy("12abc6")
	if err == nil {
		t.Error("PIN '12abc6' should fail: contains non-digits")
	}
}

func TestValidatePINPolicy_SequentialAscending(t *testing.T) {
	ensureConfig()

	err := services.ValidatePINPolicy("123456")
	if err == nil {
		t.Error("PIN '123456' should fail: sequential ascending")
	}
}

func TestValidatePINPolicy_SequentialDescending(t *testing.T) {
	ensureConfig()

	err := services.ValidatePINPolicy("654321")
	if err == nil {
		t.Error("PIN '654321' should fail: sequential descending")
	}
}

func TestValidatePINPolicy_RepeatingDigits(t *testing.T) {
	ensureConfig()

	repeatingPINs := []string{"111111", "222222", "999999", "000000"}
	for _, pin := range repeatingPINs {
		err := services.ValidatePINPolicy(pin)
		if err == nil {
			t.Errorf("PIN '%s' should fail: all repeating digits", pin)
		}
	}
}

func TestValidatePINPolicy_EmptyPIN(t *testing.T) {
	ensureConfig()

	err := services.ValidatePINPolicy("")
	if err == nil {
		t.Error("Empty PIN should fail")
	}
}
