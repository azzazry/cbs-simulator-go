package services

import (
	"database/sql"
	"fmt"
	"time"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	CIF string `json:"cif"`
	PIN string `json:"pin"`
}

// LoginResponse represents login response with JWT tokens
type LoginResponse struct {
	CIF          string `json:"cif"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RegisterRequest represents customer registration details
type RegisterRequest struct {
	CIF          string `json:"cif" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	IDCardNumber string `json:"id_card_number" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Address      string `json:"address" binding:"required"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"`
	PIN          string `json:"pin" binding:"required"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	CIF      string `json:"cif"`
	FullName string `json:"full_name"`
	Message  string `json:"message"`
}

// Authenticate verifies customer credentials with lockout and JWT
func Authenticate(req LoginRequest, clientIP string) (*LoginResponse, error) {
	// Check if account is locked (3x failed attempts)
	locked, err := IsAccountLocked(req.CIF)
	if err == nil && locked {
		failedCount, _ := GetFailedAttemptCount(req.CIF)
		return nil, fmt.Errorf("account is locked due to %d failed login attempts. Please unlock via e-KYC verification", failedCount)
	}

	var customer models.Customer
	query := `SELECT id, cif, full_name, pin, status 
	          FROM customers WHERE cif = ?`

	err = database.DB.QueryRow(query, req.CIF).Scan(
		&customer.ID, &customer.CIF, &customer.FullName,
		&customer.PIN, &customer.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid CIF or PIN")
		}
		return nil, err
	}

	// Check account status
	if customer.Status != "active" {
		return nil, fmt.Errorf("account is %s", customer.Status)
	}

	// Verify PIN
	if !utils.VerifyPIN(req.PIN, customer.PIN) {
		// Record failed attempt
		RecordLoginAttempt(req.CIF, clientIP, false)

		// Check if now locked after this failure
		failedCount, _ := GetFailedAttemptCount(req.CIF)
		remaining := 3 - failedCount
		if remaining <= 0 {
			return nil, fmt.Errorf("invalid PIN. Account is now locked. Please unlock via e-KYC verification")
		}
		return nil, fmt.Errorf("invalid CIF or PIN. %d attempts remaining", remaining)
	}

	// Record successful login
	RecordLoginAttempt(req.CIF, clientIP, true)

	// Get user role
	role, err := GetPrimaryRole(req.CIF)
	if err != nil {
		role = "customer"
	}

	// Generate JWT token pair
	tokenPair, err := GenerateTokenPair(req.CIF, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	response := &LoginResponse{
		CIF:          customer.CIF,
		FullName:     customer.FullName,
		Role:         role,
		Message:      "Login successful",
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}

	return response, nil
}

// LogoutUser blacklists the current token
func LogoutUser(jti, cif string, expiresAt time.Time) error {
	return BlacklistToken(jti, cif, expiresAt)
}

// GetCustomerByCIF retrieves customer details
func GetCustomerByCIF(cif string) (*models.Customer, error) {
	var customer models.Customer

	query := `SELECT id, cif, full_name, id_card_number, phone_number, email, 
	          address, date_of_birth, status, created_at, updated_at 
	          FROM customers WHERE cif = ?`

	err := database.DB.QueryRow(query, cif).Scan(
		&customer.ID, &customer.CIF, &customer.FullName, &customer.IDCardNumber,
		&customer.PhoneNumber, &customer.Email, &customer.Address,
		&customer.DateOfBirth, &customer.Status, &customer.CreatedAt, &customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, err
	}

	return &customer, nil
}

// ChangePIN changes customer PIN with policy validation
func ChangePIN(cif, oldPIN, newPIN string) error {
	// Validate new PIN against policy
	if err := ValidatePINPolicy(newPIN); err != nil {
		return fmt.Errorf("PIN policy violation: %v", err)
	}

	// Get current PIN
	var currentPIN string
	err := database.DB.QueryRow("SELECT pin FROM customers WHERE cif = ?", cif).Scan(&currentPIN)
	if err != nil {
		return fmt.Errorf("customer not found")
	}

	// Verify old PIN
	if !utils.VerifyPIN(oldPIN, currentPIN) {
		return fmt.Errorf("incorrect old PIN")
	}

	// Hash new PIN
	hashedPIN, err := utils.HashPIN(newPIN)
	if err != nil {
		return err
	}

	// Update PIN
	_, err = database.DB.Exec(`UPDATE customers SET pin = ?, updated_at = CURRENT_TIMESTAMP 
	                           WHERE cif = ?`, hashedPIN, cif)
	return err
}

// ResetPIN resets PIN after e-KYC + OTP verification (no old PIN needed)
func ResetPIN(cif, newPIN, verificationID string) error {
	// Validate e-KYC verification
	if err := ValidateEKYCForUnlock(verificationID, cif); err != nil {
		return fmt.Errorf("e-KYC validation failed: %v", err)
	}

	// Validate new PIN policy
	if err := ValidatePINPolicy(newPIN); err != nil {
		return fmt.Errorf("PIN policy violation: %v", err)
	}

	// Hash new PIN
	hashedPIN, err := utils.HashPIN(newPIN)
	if err != nil {
		return err
	}

	// Update PIN
	_, err = database.DB.Exec(`UPDATE customers SET pin = ?, updated_at = CURRENT_TIMESTAMP WHERE cif = ?`, hashedPIN, cif)
	return err
}

// SelfServiceUnlockAccount unlocks account after e-KYC + OTP verification
func SelfServiceUnlockAccount(cif, otp, verificationID string) error {
	// Validate e-KYC verification
	if err := ValidateEKYCForUnlock(verificationID, cif); err != nil {
		return fmt.Errorf("e-KYC validation failed: %v", err)
	}

	// Verify OTP
	if err := VerifyOTP(cif, otp, "unlock_account"); err != nil {
		return fmt.Errorf("OTP verification failed: %v", err)
	}

	// Unlock account
	return UnlockAccount(cif)
}

// RegisterCustomer creates a new customer account
func RegisterCustomer(req RegisterRequest) (*RegisterResponse, error) {
	// Validate PIN policy
	if err := ValidatePINPolicy(req.PIN); err != nil {
		return nil, fmt.Errorf("PIN policy violation: %v", err)
	}

	// Check if CIF already exists
	var existingCIF string
	err := database.DB.QueryRow("SELECT cif FROM customers WHERE cif = ?", req.CIF).Scan(&existingCIF)
	if err == nil {
		return nil, fmt.Errorf("CIF already exists")
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Hash PIN
	hashedPIN, err := utils.HashPIN(req.PIN)
	if err != nil {
		return nil, fmt.Errorf("failed to hash PIN: %v", err)
	}

	// Insert customer
	query := `INSERT INTO customers (cif, full_name, id_card_number, phone_number, email, address, date_of_birth, pin, status) 
	         VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'active')`

	_, err = database.DB.Exec(query, req.CIF, req.FullName, req.IDCardNumber, req.PhoneNumber,
		req.Email, req.Address, req.DateOfBirth, hashedPIN)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %v", err)
	}

	// Assign default customer role
	role, _ := GetRoleByName("customer")
	if role != nil {
		AssignRole(req.CIF, role.ID, "SYSTEM")
	}

	response := &RegisterResponse{
		CIF:      req.CIF,
		FullName: req.FullName,
		Message:  "Customer registered successfully",
	}

	return response, nil
}
