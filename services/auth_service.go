package services

import (
	"database/sql"
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	CIF string `json:"cif"`
	PIN string `json:"pin"`
}

// LoginResponse represents login response with customer info
type LoginResponse struct {
	CIF      string `json:"cif"`
	FullName string `json:"full_name"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"` // For JWT implementation
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

// Authenticate verifies customer credentials
func Authenticate(req LoginRequest) (*LoginResponse, error) {
	var customer models.Customer

	query := `SELECT id, cif, full_name, pin, status 
	          FROM customers WHERE cif = ?`

	err := database.DB.QueryRow(query, req.CIF).Scan(
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
		return nil, fmt.Errorf("invalid CIF or PIN")
	}

	// In production, generate JWT token here
	// For simulator, we'll use simple response
	response := &LoginResponse{
		CIF:      customer.CIF,
		FullName: customer.FullName,
		Message:  "Login successful",
		Token:    fmt.Sprintf("token_%s", customer.CIF), // Mock token
	}

	return response, nil
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

// ChangePIN changes customer PIN
func ChangePIN(cif, oldPIN, newPIN string) error {
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
// RegisterCustomer creates a new customer account
func RegisterCustomer(req RegisterRequest) (*RegisterResponse, error) {
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
	
	response := &RegisterResponse{
		CIF:      req.CIF,
		FullName: req.FullName,
		Message:  "Customer registered successfully",
	}
	
	return response, nil
}