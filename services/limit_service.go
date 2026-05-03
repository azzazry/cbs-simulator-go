package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"database/sql"
	"fmt"
)

// GetTransactionLimit returns the limit for a role and transaction type
func GetTransactionLimit(roleName, transactionType string) (*models.TransactionLimit, error) {
	var limit models.TransactionLimit
	query := `SELECT id, role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit, is_active, created_at, updated_at
	          FROM transaction_limits WHERE role_name = $1 AND transaction_type = $2 AND is_active = TRUE`

	err := database.DB.QueryRow(query, roleName, transactionType).Scan(
		&limit.ID, &limit.RoleName, &limit.TransactionType,
		&limit.DailyLimit, &limit.PerTransactionLimit,
		&limit.MonthlyLimit, &limit.IsActive, &limit.CreatedAt, &limit.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no limit found for role=%s, type=%s", roleName, transactionType)
		}
		return nil, err
	}
	return &limit, nil
}

// GetAllTransactionLimits returns all active transaction limits
func GetAllTransactionLimits() ([]models.TransactionLimit, error) {
	query := `SELECT id, role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit, is_active, created_at, updated_at
	          FROM transaction_limits WHERE is_active = TRUE ORDER BY role_name, transaction_type`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []models.TransactionLimit
	for rows.Next() {
		var l models.TransactionLimit
		err := rows.Scan(&l.ID, &l.RoleName, &l.TransactionType,
			&l.DailyLimit, &l.PerTransactionLimit,
			&l.MonthlyLimit, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		limits = append(limits, l)
	}
	return limits, nil
}

// CheckTransactionLimit checks if an amount is within per-transaction limit for a CIF
func CheckTransactionLimit(cif, transactionType string, amount float64) error {
	role, err := GetPrimaryRole(cif)
	if err != nil {
		role = "customer"
	}

	limit, err := GetTransactionLimit(role, transactionType)
	if err != nil {
		return nil // No limit configured = allow
	}

	if limit.PerTransactionLimit > 0 && amount > limit.PerTransactionLimit {
		return fmt.Errorf("amount %.0f exceeds per-transaction limit of %.0f", amount, limit.PerTransactionLimit)
	}

	// Daily limit check
	if limit.DailyLimit > 0 {
		var dailyTotal float64
		database.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions
		                      WHERE from_account_number IN (SELECT account_number FROM accounts WHERE cif = $1)
		                      AND transaction_type ILIKE $2
		                      AND transaction_date >= CURRENT_DATE
		                      AND status = 'success'`, cif, "%"+transactionType+"%").Scan(&dailyTotal)

		if dailyTotal+amount > limit.DailyLimit {
			return fmt.Errorf("transaction would exceed daily limit of %.0f (used: %.0f)", limit.DailyLimit, dailyTotal)
		}
	}

	return nil
}

// UpdateTransactionLimit updates or inserts a transaction limit (upsert)
func UpdateTransactionLimit(roleName, transactionType string, daily, perTrx, monthly float64) error {
	query := `INSERT INTO transaction_limits (role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit)
	          VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (role_name, transaction_type) DO UPDATE
	          SET daily_limit = EXCLUDED.daily_limit,
	              per_transaction_limit = EXCLUDED.per_transaction_limit,
	              monthly_limit = EXCLUDED.monthly_limit,
	              updated_at = NOW()`

	_, err := database.DB.Exec(query, roleName, transactionType, daily, perTrx, monthly)
	return err
}
