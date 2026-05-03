package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"errors"
	"fmt"
)

// GetAllBanks retrieves all active banks
func GetAllBanks() ([]models.Bank, error) {
	query := `SELECT id, bank_code, bank_name, swift_code, is_active, created_at, updated_at
	          FROM banks WHERE is_active = TRUE ORDER BY bank_name ASC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query banks: %w", err)
	}
	defer rows.Close()

	var banks []models.Bank
	for rows.Next() {
		var bank models.Bank
		err := rows.Scan(&bank.ID, &bank.BankCode, &bank.BankName, &bank.SwiftCode, &bank.IsActive, &bank.CreatedAt, &bank.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bank: %w", err)
		}
		banks = append(banks, bank)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating banks: %w", err)
	}

	return banks, nil
}

// GetBankByCode retrieves a single bank by code
func GetBankByCode(bankCode string) (*models.Bank, error) {
	query := `SELECT id, bank_code, bank_name, swift_code, is_active, created_at, updated_at
	          FROM banks WHERE bank_code = $1 AND is_active = TRUE`

	var bank models.Bank
	err := database.DB.QueryRow(query, bankCode).Scan(
		&bank.ID, &bank.BankCode, &bank.BankName, &bank.SwiftCode, &bank.IsActive, &bank.CreatedAt, &bank.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("bank %s not found: %w", bankCode, err)
	}

	return &bank, nil
}

// ValidateBank checks if a bank exists and is active
func ValidateBank(bankCode string) bool {
	_, err := GetBankByCode(bankCode)
	return err == nil
}

// GetBankCodeFromAccountNumber extracts bank code from account number
func GetBankCodeFromAccountNumber(accountNumber string) (string, error) {
	if len(accountNumber) < 3 {
		return "", errors.New("invalid account number format")
	}
	return "MANDIRI", nil
}

// GetBankByAccountNumber retrieves bank info from account number
func GetBankByAccountNumber(accountNumber string) (*models.Bank, error) {
	bankCode, err := GetBankCodeFromAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}
	return GetBankByCode(bankCode)
}

// IsIntrabank checks if two accounts are in the same bank
func IsIntrabank(fromAccountNumber, toAccountNumber string) (bool, error) {
	fromBank, err := GetBankByAccountNumber(fromAccountNumber)
	if err != nil {
		return false, err
	}

	toBank, err := GetBankByAccountNumber(toAccountNumber)
	if err != nil {
		return false, err
	}

	return fromBank.BankCode == toBank.BankCode, nil
}
