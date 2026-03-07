package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"database/sql"
	"fmt"
	"time"
)

// CalculateTransferFee calculates the fee for a transfer OUT to another bank (interbank transfer)
// This bank is ONE specific bank. When customers transfer to other banks, they pay a fee.
// Intrabank (same bank) transfers are always FREE
func CalculateTransferFee(destinationBankCode string, amount float64) (float64, error) {
	query := `SELECT fee_type, fee_amount, fee_percentage, minimum_amount, maximum_amount 
	          FROM transfer_fees 
	          WHERE destination_bank_code = ? 
	          AND is_active = 1 
	          AND effective_from <= datetime('now')
	          AND (effective_to IS NULL OR effective_to > datetime('now'))
	          LIMIT 1`

	var feeType string
	var feeAmount, feePercentage, minAmount, maxAmount float64

	err := database.DB.QueryRow(query, destinationBankCode).Scan(
		&feeType, &feeAmount, &feePercentage, &minAmount, &maxAmount,
	)

	if err != nil {
		// Return default fee if not found (Rp 5000)
		return 5000, nil
	}

	// Check minimum and maximum amount constraints
	if amount < minAmount || amount > maxAmount {
		// Use default fee if amount is outside bounds
		return 5000, nil
	}

	// Calculate fee based on type
	var fee float64
	if feeType == "percentage" {
		fee = (amount * feePercentage) / 100
	} else {
		fee = feeAmount
	}

	return fee, nil
}

// CalculateServiceFee calculates the fee for a service (e-wallet top-up, VA payment, QRIS, etc)
func CalculateServiceFee(serviceCode string, amount float64) (float64, error) {
	query := `SELECT fee_type, fee_amount, fee_percentage, minimum_amount, maximum_amount 
	          FROM service_fees 
	          WHERE service_code = ? 
	          AND is_active = 1 
	          AND effective_from <= datetime('now')
	          AND (effective_to IS NULL OR effective_to > datetime('now'))
	          LIMIT 1`

	var feeType string
	var feeAmount, feePercentage, minAmount, maxAmount float64

	err := database.DB.QueryRow(query, serviceCode).Scan(
		&feeType, &feeAmount, &feePercentage, &minAmount, &maxAmount,
	)

	if err != nil {
		// Return default fee if not found
		return 0, fmt.Errorf("service fee not found: %w", err)
	}

	// Check minimum and maximum amount constraints
	if amount < minAmount || amount > maxAmount {
		return 0, fmt.Errorf("amount %f is outside allowed range (%.0f - %.0f)", amount, minAmount, maxAmount)
	}

	// Calculate fee based on type
	var fee float64
	if feeType == "percentage" {
		fee = (amount * feePercentage) / 100
	} else {
		fee = feeAmount
	}

	return fee, nil
}

// GetTransferFee retrieves a specific transfer fee configuration
func GetTransferFee(destinationBankCode string) (*models.TransferFee, error) {
	query := `SELECT id, destination_bank_code, destination_bank_name, fee_type, fee_amount, 
	          fee_percentage, minimum_amount, maximum_amount, description, is_active, 
	          effective_from, effective_to, notes, created_at, updated_at
	          FROM transfer_fees 
	          WHERE destination_bank_code = ? AND is_active = 1
	          LIMIT 1`

	var fee models.TransferFee
	var effectiveTo *time.Time

	err := database.DB.QueryRow(query, destinationBankCode).Scan(
		&fee.ID, &fee.DestinationBankCode, &fee.DestinationBankName, &fee.FeeType,
		&fee.FeeAmount, &fee.FeePercentage, &fee.MinimumAmount, &fee.MaximumAmount,
		&fee.Description, &fee.IsActive, &fee.EffectiveFrom, &effectiveTo, &fee.Notes,
		&fee.CreatedAt, &fee.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("transfer fee not found: %w", err)
	}

	fee.EffectiveTo = effectiveTo
	return &fee, nil
}

// GetAllTransferFees retrieves all transfer fees
func GetAllTransferFees() ([]models.TransferFee, error) {
	query := `SELECT id, destination_bank_code, destination_bank_name, fee_type, fee_amount, 
	          fee_percentage, minimum_amount, maximum_amount, description, is_active, 
	          effective_from, effective_to, notes, created_at, updated_at
	          FROM transfer_fees 
	          WHERE is_active = 1
	          ORDER BY destination_bank_name`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query transfer fees: %w", err)
	}
	defer rows.Close()

	var fees []models.TransferFee
	for rows.Next() {
		var fee models.TransferFee
		var effectiveTo *time.Time

		err := rows.Scan(
			&fee.ID, &fee.DestinationBankCode, &fee.DestinationBankName, &fee.FeeType,
			&fee.FeeAmount, &fee.FeePercentage, &fee.MinimumAmount, &fee.MaximumAmount,
			&fee.Description, &fee.IsActive, &fee.EffectiveFrom, &effectiveTo, &fee.Notes,
			&fee.CreatedAt, &fee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fee: %w", err)
		}

		fee.EffectiveTo = effectiveTo
		fees = append(fees, fee)
	}

	return fees, nil
}

// UpdateTransferFee updates transfer fee configuration
func UpdateTransferFee(destinationBankCode string, feeAmount float64, feeType string) error {
	query := `UPDATE transfer_fees 
	          SET fee_amount = ?, fee_type = ?, updated_at = datetime('now')
	          WHERE destination_bank_code = ?`

	_, err := database.DB.Exec(query, feeAmount, feeType, destinationBankCode)
	if err != nil {
		return fmt.Errorf("failed to update transfer fee: %w", err)
	}

	return nil
}

// GetServiceFee retrieves a specific service fee configuration
func GetServiceFee(serviceCode string) (*models.ServiceFee, error) {
	query := `SELECT id, service_code, service_name, service_type, fee_type, fee_amount, 
	          fee_percentage, minimum_amount, maximum_amount, is_active, effective_from, 
	          effective_to, notes, created_at, updated_at
	          FROM service_fees 
	          WHERE service_code = ? AND is_active = 1
	          LIMIT 1`

	var fee models.ServiceFee
	var effectiveTo *time.Time

	err := database.DB.QueryRow(query, serviceCode).Scan(
		&fee.ID, &fee.ServiceCode, &fee.ServiceName, &fee.ServiceType, &fee.FeeType,
		&fee.FeeAmount, &fee.FeePercentage, &fee.MinimumAmount, &fee.MaximumAmount,
		&fee.IsActive, &fee.EffectiveFrom, &effectiveTo, &fee.Notes,
		&fee.CreatedAt, &fee.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("service fee not found: %w", err)
	}

	fee.EffectiveTo = effectiveTo
	return &fee, nil
}

// GetAllServiceFees retrieves all service fees
func GetAllServiceFees(serviceType string) ([]models.ServiceFee, error) {
	query := `SELECT id, service_code, service_name, service_type, fee_type, fee_amount, 
	          fee_percentage, minimum_amount, maximum_amount, is_active, effective_from, 
	          effective_to, notes, created_at, updated_at
	          FROM service_fees 
	          WHERE is_active = 1`

	if serviceType != "" {
		query += ` AND service_type = ?`
	}

	query += ` ORDER BY service_type, service_code`

	var rows interface{}
	var err error

	if serviceType != "" {
		rows, err = database.DB.Query(query, serviceType)
	} else {
		rows, err = database.DB.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query service fees: %w", err)
	}

	// Cast to *sql.Rows type
	sqlRows := rows.(*sql.Rows)
	defer sqlRows.Close()

	var fees []models.ServiceFee
	for sqlRows.Next() {
		var fee models.ServiceFee
		var effectiveTo *time.Time

		err := sqlRows.Scan(
			&fee.ID, &fee.ServiceCode, &fee.ServiceName, &fee.ServiceType, &fee.FeeType,
			&fee.FeeAmount, &fee.FeePercentage, &fee.MinimumAmount, &fee.MaximumAmount,
			&fee.IsActive, &fee.EffectiveFrom, &effectiveTo, &fee.Notes,
			&fee.CreatedAt, &fee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fee: %w", err)
		}

		fee.EffectiveTo = effectiveTo
		fees = append(fees, fee)
	}

	return fees, nil
}

// UpdateServiceFee updates service fee configuration
func UpdateServiceFee(serviceCode string, feeAmount float64, feePercentage float64, feeType string) error {
	query := `UPDATE service_fees 
	          SET fee_amount = ?, fee_percentage = ?, fee_type = ?, updated_at = datetime('now')
	          WHERE service_code = ?`

	_, err := database.DB.Exec(query, feeAmount, feePercentage, feeType, serviceCode)
	if err != nil {
		return fmt.Errorf("failed to update service fee: %w", err)
	}

	return nil
}

// GetServiceFeesByType retrieves all service fees by type
func GetServiceFeesByType(serviceType string) ([]models.ServiceFee, error) {
	query := `SELECT id, service_code, service_name, service_type, fee_type, fee_amount, 
	          fee_percentage, minimum_amount, maximum_amount, is_active, effective_from, 
	          effective_to, notes, created_at, updated_at
	          FROM service_fees 
	          WHERE service_type = ? AND is_active = 1
	          ORDER BY service_code`

	rows, err := database.DB.Query(query, serviceType)
	if err != nil {
		return nil, fmt.Errorf("failed to query service fees: %w", err)
	}
	defer rows.Close()

	var fees []models.ServiceFee
	for rows.Next() {
		var fee models.ServiceFee
		var effectiveTo *time.Time

		err := rows.Scan(
			&fee.ID, &fee.ServiceCode, &fee.ServiceName, &fee.ServiceType, &fee.FeeType,
			&fee.FeeAmount, &fee.FeePercentage, &fee.MinimumAmount, &fee.MaximumAmount,
			&fee.IsActive, &fee.EffectiveFrom, &effectiveTo, &fee.Notes,
			&fee.CreatedAt, &fee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fee: %w", err)
		}

		fee.EffectiveTo = effectiveTo
		fees = append(fees, fee)
	}

	return fees, nil
}
