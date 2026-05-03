package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"database/sql"
	"fmt"
	"time"
)

// CalculateTransferFee calculates outbound interbank transfer fee
func CalculateTransferFee(destinationBankCode string, amount float64) (float64, error) {
	query := `SELECT fee_type, fee_amount, fee_percentage, minimum_amount, maximum_amount
	          FROM transfer_fees
	          WHERE destination_bank_code = $1
	          AND is_active = TRUE
	          AND effective_from <= NOW()
	          AND (effective_to IS NULL OR effective_to > NOW())
	          LIMIT 1`

	var feeType string
	var feeAmount, feePercentage, minAmount, maxAmount float64

	err := database.DB.QueryRow(query, destinationBankCode).Scan(
		&feeType, &feeAmount, &feePercentage, &minAmount, &maxAmount,
	)

	if err != nil {
		return 5000, nil
	}

	if amount < minAmount || amount > maxAmount {
		return 5000, nil
	}

	var fee float64
	if feeType == "percentage" {
		fee = (amount * feePercentage) / 100
	} else {
		fee = feeAmount
	}

	return fee, nil
}

// CalculateServiceFee calculates service fee (e-wallet, VA, QRIS, etc)
func CalculateServiceFee(serviceCode string, amount float64) (float64, error) {
	query := `SELECT fee_type, fee_amount, fee_percentage, minimum_amount, maximum_amount
	          FROM service_fees
	          WHERE service_code = $1
	          AND is_active = TRUE
	          AND effective_from <= NOW()
	          AND (effective_to IS NULL OR effective_to > NOW())
	          LIMIT 1`

	var feeType string
	var feeAmount, feePercentage, minAmount, maxAmount float64

	err := database.DB.QueryRow(query, serviceCode).Scan(
		&feeType, &feeAmount, &feePercentage, &minAmount, &maxAmount,
	)

	if err != nil {
		return 0, fmt.Errorf("service fee not found: %w", err)
	}

	if amount < minAmount || amount > maxAmount {
		return 0, fmt.Errorf("amount %f is outside allowed range (%.0f - %.0f)", amount, minAmount, maxAmount)
	}

	var fee float64
	if feeType == "percentage" {
		fee = (amount * feePercentage) / 100
	} else {
		fee = feeAmount
	}

	return fee, nil
}

func GetTransferFee(destinationBankCode string) (*models.TransferFee, error) {
	query := `SELECT id, destination_bank_code, destination_bank_name, fee_type, fee_amount,
	          fee_percentage, minimum_amount, maximum_amount, description, is_active,
	          effective_from, effective_to, notes, created_at, updated_at
	          FROM transfer_fees
	          WHERE destination_bank_code = $1 AND is_active = TRUE
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

func GetAllTransferFees() ([]models.TransferFee, error) {
	query := `SELECT id, destination_bank_code, destination_bank_name, fee_type, fee_amount,
	          fee_percentage, minimum_amount, maximum_amount, description, is_active,
	          effective_from, effective_to, notes, created_at, updated_at
	          FROM transfer_fees
	          WHERE is_active = TRUE
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

func UpdateTransferFee(destinationBankCode string, feeAmount float64, feeType string) error {
	query := `UPDATE transfer_fees
	          SET fee_amount = $1, fee_type = $2, updated_at = NOW()
	          WHERE destination_bank_code = $3`

	_, err := database.DB.Exec(query, feeAmount, feeType, destinationBankCode)
	if err != nil {
		return fmt.Errorf("failed to update transfer fee: %w", err)
	}

	return nil
}

func GetServiceFee(serviceCode string) (*models.ServiceFee, error) {
	query := `SELECT id, service_code, service_name, service_type, fee_type, fee_amount,
	          fee_percentage, minimum_amount, maximum_amount, is_active, effective_from,
	          effective_to, notes, created_at, updated_at
	          FROM service_fees
	          WHERE service_code = $1 AND is_active = TRUE
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

func GetAllServiceFees(serviceType string) ([]models.ServiceFee, error) {
	args := []interface{}{}
	query := `SELECT id, service_code, service_name, service_type, fee_type, fee_amount,
	          fee_percentage, minimum_amount, maximum_amount, is_active, effective_from,
	          effective_to, notes, created_at, updated_at
	          FROM service_fees
	          WHERE is_active = TRUE`

	if serviceType != "" {
		query += ` AND service_type = $1`
		args = append(args, serviceType)
	}

	query += ` ORDER BY service_type, service_code`

	rows, err := database.DB.Query(query, args...)
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

func UpdateServiceFee(serviceCode string, feeAmount float64, feePercentage float64, feeType string) error {
	query := `UPDATE service_fees
	          SET fee_amount = $1, fee_percentage = $2, fee_type = $3, updated_at = NOW()
	          WHERE service_code = $4`

	_, err := database.DB.Exec(query, feeAmount, feePercentage, feeType, serviceCode)
	if err != nil {
		return fmt.Errorf("failed to update service fee: %w", err)
	}

	return nil
}

func GetServiceFeesByType(serviceType string) ([]models.ServiceFee, error) {
	query := `SELECT id, service_code, service_name, service_type, fee_type, fee_amount,
	          fee_percentage, minimum_amount, maximum_amount, is_active, effective_from,
	          effective_to, notes, created_at, updated_at
	          FROM service_fees
	          WHERE service_type = $1 AND is_active = TRUE
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

// GetServiceFeeRows helper used in GetAllServiceFees (replaces interface{} cast workaround)
func GetServiceFeeRows(serviceType string) (*sql.Rows, error) {
	if serviceType != "" {
		return database.DB.Query(
			`SELECT id FROM service_fees WHERE service_type = $1 AND is_active = TRUE`, serviceType)
	}
	return database.DB.Query(`SELECT id FROM service_fees WHERE is_active = TRUE`)
}
