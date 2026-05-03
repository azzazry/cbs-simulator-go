package services

import (
	"database/sql"
	"fmt"
	"time"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

type LoginRequest struct {
	CIF string `json:"cif"`
	PIN string `json:"pin"`
}

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

type RegisterResponse struct {
	CIF      string `json:"cif"`
	FullName string `json:"full_name"`
	Message  string `json:"message"`
}

func Authenticate(req LoginRequest, clientIP string) (*LoginResponse, error) {
	locked, err := IsAccountLocked(req.CIF)
	if err == nil && locked {
		failedCount, _ := GetFailedAttemptCount(req.CIF)
		return nil, fmt.Errorf("account is locked due to %d failed login attempts. Please unlock via e-KYC verification", failedCount)
	}

	var customer models.Customer
	query := `SELECT id, cif, full_name, pin, status FROM customers WHERE cif = $1`

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

	if customer.Status != "active" {
		return nil, fmt.Errorf("account is %s", customer.Status)
	}

	if !utils.VerifyPIN(req.PIN, customer.PIN) {
		RecordLoginAttempt(req.CIF, clientIP, false)

		failedCount, _ := GetFailedAttemptCount(req.CIF)
		remaining := 3 - failedCount
		if remaining <= 0 {
			return nil, fmt.Errorf("invalid PIN. Account is now locked. Please unlock via e-KYC verification")
		}
		return nil, fmt.Errorf("invalid CIF or PIN. %d attempts remaining", remaining)
	}

	RecordLoginAttempt(req.CIF, clientIP, true)

	role, err := GetPrimaryRole(req.CIF)
	if err != nil {
		role = "customer"
	}

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

func LogoutUser(jti, cif string, expiresAt time.Time) error {
	return BlacklistToken(jti, cif, expiresAt)
}

func GetCustomerByCIF(cif string) (*models.Customer, error) {
	var customer models.Customer

	query := `SELECT id, cif, full_name, id_card_number, phone_number, email,
	          address, date_of_birth, status, created_at, updated_at
	          FROM customers WHERE cif = $1`

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

func ChangePIN(cif, oldPIN, newPIN string) error {
	if err := ValidatePINPolicy(newPIN); err != nil {
		return fmt.Errorf("PIN policy violation: %v", err)
	}

	var currentPIN string
	err := database.DB.QueryRow("SELECT pin FROM customers WHERE cif = $1", cif).Scan(&currentPIN)
	if err != nil {
		return fmt.Errorf("customer not found")
	}

	if !utils.VerifyPIN(oldPIN, currentPIN) {
		return fmt.Errorf("incorrect old PIN")
	}

	hashedPIN, err := utils.HashPIN(newPIN)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`UPDATE customers SET pin = $1, updated_at = NOW() WHERE cif = $2`, hashedPIN, cif)
	return err
}

func ResetPIN(cif, newPIN, verificationID string) error {
	if err := ValidateEKYCForUnlock(verificationID, cif); err != nil {
		return fmt.Errorf("e-KYC validation failed: %v", err)
	}

	if err := ValidatePINPolicy(newPIN); err != nil {
		return fmt.Errorf("PIN policy violation: %v", err)
	}

	hashedPIN, err := utils.HashPIN(newPIN)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`UPDATE customers SET pin = $1, updated_at = NOW() WHERE cif = $2`, hashedPIN, cif)
	return err
}

func SelfServiceUnlockAccount(cif, otp, verificationID string) error {
	if err := ValidateEKYCForUnlock(verificationID, cif); err != nil {
		return fmt.Errorf("e-KYC validation failed: %v", err)
	}

	if err := VerifyOTP(cif, otp, "unlock_account"); err != nil {
		return fmt.Errorf("OTP verification failed: %v", err)
	}

	return UnlockAccount(cif)
}

func RegisterCustomer(req RegisterRequest) (*RegisterResponse, error) {
	if err := ValidatePINPolicy(req.PIN); err != nil {
		return nil, fmt.Errorf("PIN policy violation: %v", err)
	}

	var existingCIF string
	err := database.DB.QueryRow("SELECT cif FROM customers WHERE cif = $1", req.CIF).Scan(&existingCIF)
	if err == nil {
		return nil, fmt.Errorf("CIF already exists")
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	hashedPIN, err := utils.HashPIN(req.PIN)
	if err != nil {
		return nil, fmt.Errorf("failed to hash PIN: %v", err)
	}

	query := `INSERT INTO customers (cif, full_name, id_card_number, phone_number, email, address, date_of_birth, pin, status)
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active')`

	_, err = database.DB.Exec(query, req.CIF, req.FullName, req.IDCardNumber, req.PhoneNumber,
		req.Email, req.Address, req.DateOfBirth, hashedPIN)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %v", err)
	}

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
