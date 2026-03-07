package services_test

import (
	"testing"

	"cbs-simulator/services"
)

func TestAuthenticate_Success(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	resp, err := services.Authenticate(services.LoginRequest{
		CIF: "CIF001",
		PIN: "123456",
	}, "127.0.0.1")

	if err != nil {
		t.Fatalf("Login should succeed: %v", err)
	}
	if resp.CIF != "CIF001" {
		t.Errorf("CIF should be 'CIF001', got '%s'", resp.CIF)
	}
	if resp.AccessToken == "" {
		t.Error("Access token should not be empty")
	}
	if resp.RefreshToken == "" {
		t.Error("Refresh token should not be empty")
	}
	if resp.Role != "admin" {
		t.Errorf("Role should be 'admin', got '%s'", resp.Role)
	}
}

func TestAuthenticate_WrongPIN(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := services.Authenticate(services.LoginRequest{
		CIF: "CIF001",
		PIN: "999999",
	}, "127.0.0.1")

	if err == nil {
		t.Error("Login with wrong PIN should fail")
	}
}

func TestAuthenticate_NonExistentCIF(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := services.Authenticate(services.LoginRequest{
		CIF: "CIFXXX",
		PIN: "123456",
	}, "127.0.0.1")

	if err == nil {
		t.Error("Login with nonexistent CIF should fail")
	}
}

func TestAuthenticate_Lockout3Attempts(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// 3 failed attempts — each should return error
	var lastErr error
	for i := 0; i < 3; i++ {
		_, lastErr = services.Authenticate(services.LoginRequest{
			CIF: "CIF002",
			PIN: "wrong",
		}, "127.0.0.1")
		if lastErr == nil {
			t.Fatalf("Attempt %d with wrong PIN should fail", i+1)
		}
	}

	// The 3rd failure should trigger lockout
	// Now verify: even with correct PIN, should be locked
	_, err := services.Authenticate(services.LoginRequest{
		CIF: "CIF002",
		PIN: "123456",
	}, "127.0.0.1")

	if err == nil {
		t.Error("Account should be locked after 3 failed attempts, even with correct PIN")
	}
	t.Logf("Lockout error: %v", err)
}

func TestAuthenticate_LockoutMessage(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// 2 failed attempts
	for i := 0; i < 2; i++ {
		services.Authenticate(services.LoginRequest{
			CIF: "CIF002",
			PIN: "wrong",
		}, "127.0.0.1")
	}

	// 3rd attempt should mention remaining attempts or lock
	_, err := services.Authenticate(services.LoginRequest{
		CIF: "CIF002",
		PIN: "wrong",
	}, "127.0.0.1")

	if err == nil {
		t.Error("Should fail on 3rd bad attempt")
	}
}

func TestRegisterCustomer_Success(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	resp, err := services.RegisterCustomer(services.RegisterRequest{
		CIF:          "CIF099",
		FullName:     "Test User",
		IDCardNumber: "3201991234567899",
		PhoneNumber:  "081299999999",
		Email:        "test@email.com",
		Address:      "Jl. Test",
		DateOfBirth:  "2000-01-01",
		PIN:          "975312",
	})

	if err != nil {
		t.Fatalf("Registration should succeed: %v", err)
	}
	if resp.CIF != "CIF099" {
		t.Errorf("CIF should be 'CIF099', got '%s'", resp.CIF)
	}
}

func TestRegisterCustomer_DuplicateCIF(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := services.RegisterCustomer(services.RegisterRequest{
		CIF:          "CIF001", // already exists
		FullName:     "Duplicate",
		IDCardNumber: "9999999999999999",
		PhoneNumber:  "081200000000",
		Email:        "dup@email.com",
		Address:      "Jl. Dup",
		DateOfBirth:  "2000-01-01",
		PIN:          "975312",
	})

	if err == nil {
		t.Error("Should fail: CIF already exists")
	}
}

func TestRegisterCustomer_WeakPIN(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := services.RegisterCustomer(services.RegisterRequest{
		CIF:          "CIF098",
		FullName:     "Weak PIN User",
		IDCardNumber: "3201981234567898",
		PhoneNumber:  "081298888888",
		Email:        "weak@email.com",
		Address:      "Jl. Weak",
		DateOfBirth:  "2000-01-01",
		PIN:          "111111", // repeating digits
	})

	if err == nil {
		t.Error("Should fail: weak PIN (repeating)")
	}
}

func TestChangePIN_Success(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	err := services.ChangePIN("CIF001", "123456", "975312")
	if err != nil {
		t.Fatalf("PIN change should succeed: %v", err)
	}

	// Login with new PIN should work
	_, err = services.Authenticate(services.LoginRequest{
		CIF: "CIF001",
		PIN: "975312",
	}, "127.0.0.1")

	if err != nil {
		t.Errorf("Login with new PIN should succeed: %v", err)
	}
}

func TestChangePIN_WrongOldPIN(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	err := services.ChangePIN("CIF001", "wrong_pin", "975312")
	if err == nil {
		t.Error("Should fail: wrong old PIN")
	}
}

func TestChangePIN_WeakNewPIN(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	err := services.ChangePIN("CIF001", "123456", "111111")
	if err == nil {
		t.Error("Should fail: weak new PIN")
	}
}

func TestLogoutUser(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// Login
	resp, _ := services.Authenticate(services.LoginRequest{
		CIF: "CIF001",
		PIN: "123456",
	}, "127.0.0.1")

	// Validate to get claims
	claims, _ := services.ValidateToken(resp.AccessToken)

	// Logout
	err := services.LogoutUser(claims.ID, "CIF001", claims.ExpiresAt.Time)
	if err != nil {
		t.Fatalf("Logout should succeed: %v", err)
	}

	// Token should be rejected now
	_, err = services.ValidateToken(resp.AccessToken)
	if err == nil {
		t.Error("Token should be rejected after logout")
	}
}

func TestGetCustomerByCIF(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	customer, err := services.GetCustomerByCIF("CIF001")
	if err != nil {
		t.Fatalf("Should find customer: %v", err)
	}
	if customer.FullName != "Budi Santoso" {
		t.Errorf("Name should be 'Budi Santoso', got '%s'", customer.FullName)
	}
}

func TestGetCustomerByCIF_NotFound(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := services.GetCustomerByCIF("CIFXXX")
	if err == nil {
		t.Error("Should fail: customer not found")
	}
}
