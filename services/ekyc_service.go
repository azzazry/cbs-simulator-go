package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// EKYCVerifyRequest represents an e-KYC verification request
type EKYCVerifyRequest struct {
	CIF              string `json:"cif" binding:"required"`
	IDCardNumber     string `json:"id_card_number" binding:"required"`
	VerificationType string `json:"verification_type"` // unlock_account, reset_pin
}

// EKYCVerifyResponse represents the e-KYC verification result
type EKYCVerifyResponse struct {
	VerificationID string `json:"verification_id"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}

// VerifyEKYC performs e-KYC verification by matching KTP with customer data
func VerifyEKYC(req EKYCVerifyRequest) (*EKYCVerifyResponse, error) {
	if req.VerificationType == "" {
		req.VerificationType = "unlock_account"
	}

	var storedIDCard, fullName string
	query := `SELECT id_card_number, full_name FROM customers WHERE cif = $1`
	err := database.DB.QueryRow(query, req.CIF).Scan(&storedIDCard, &fullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to verify customer: %v", err)
	}

	verificationID := fmt.Sprintf("EKYC-%s", uuid.New().String()[:12])
	status := "rejected"
	message := "ID card number does not match our records"

	if storedIDCard == req.IDCardNumber {
		status = "verified"
		message = fmt.Sprintf("e-KYC verification successful for %s", fullName)
	}

	// PostgreSQL: gunakan CASE WHEN untuk verified_at
	insertQuery := `INSERT INTO ekyc_verifications (verification_id, cif, id_card_number, verification_type, verification_status, verified_at)
	                VALUES ($1, $2, $3, $4, $5, CASE WHEN $6 = 'verified' THEN NOW() ELSE NULL END)`
	_, err = database.DB.Exec(insertQuery, verificationID, req.CIF, req.IDCardNumber, req.VerificationType, status, status)
	if err != nil {
		return nil, fmt.Errorf("failed to record e-KYC verification: %v", err)
	}

	if status == "rejected" {
		return nil, fmt.Errorf("%s", message)
	}

	return &EKYCVerifyResponse{
		VerificationID: verificationID,
		Status:         status,
		Message:        message,
	}, nil
}

// GetEKYCVerification retrieves a verification by ID
func GetEKYCVerification(verificationID string) (*models.EKYCVerification, error) {
	var v models.EKYCVerification
	query := `SELECT id, verification_id, cif, id_card_number, verification_type, verification_status, verified_at, created_at
	          FROM ekyc_verifications WHERE verification_id = $1`
	err := database.DB.QueryRow(query, verificationID).Scan(
		&v.ID, &v.VerificationID, &v.CIF, &v.IDCardNumber,
		&v.VerificationType, &v.VerificationStatus, &v.VerifiedAt, &v.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("verification not found")
		}
		return nil, err
	}
	return &v, nil
}

// ValidateEKYCForUnlock checks if an e-KYC verification is valid for account unlock
func ValidateEKYCForUnlock(verificationID, cif string) error {
	v, err := GetEKYCVerification(verificationID)
	if err != nil {
		return fmt.Errorf("invalid verification ID")
	}

	if v.CIF != cif {
		return fmt.Errorf("verification does not belong to this customer")
	}

	if v.VerificationStatus != "verified" {
		return fmt.Errorf("e-KYC verification was not successful")
	}

	return nil
}
