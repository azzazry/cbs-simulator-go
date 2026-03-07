package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
	"time"
)

// CheckTransactionLimit validates a transaction against configured limits
func CheckTransactionLimit(cif, transactionType string, amount float64) error {
	// Get user's primary role
	role, err := GetPrimaryRole(cif)
	if err != nil {
		role = "customer"
	}

	// Admin and supervisor bypass limits
	if role == "admin" || role == "supervisor" {
		return nil
	}

	// Get limit for this role and transaction type
	var limit models.TransactionLimit
	query := `SELECT id, role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit 
	          FROM transaction_limits WHERE role_name = ? AND transaction_type = ? AND is_active = 1`
	err = database.DB.QueryRow(query, role, transactionType).Scan(
		&limit.ID, &limit.RoleName, &limit.TransactionType,
		&limit.DailyLimit, &limit.PerTransactionLimit, &limit.MonthlyLimit)
	if err != nil {
		// No limit configured, allow
		return nil
	}

	// Check per-transaction limit
	if limit.PerTransactionLimit > 0 && amount > limit.PerTransactionLimit {
		return fmt.Errorf("amount exceeds per-transaction limit of Rp %.0f", limit.PerTransactionLimit)
	}

	// Check daily limit
	if limit.DailyLimit > 0 {
		dailyUsage, err := GetDailyUsage(cif, transactionType)
		if err != nil {
			return fmt.Errorf("failed to check daily usage: %v", err)
		}
		if dailyUsage+amount > limit.DailyLimit {
			return fmt.Errorf("transaction would exceed daily limit of Rp %.0f (used: Rp %.0f)", limit.DailyLimit, dailyUsage)
		}
	}

	// Check monthly limit
	if limit.MonthlyLimit > 0 {
		monthlyUsage, err := GetMonthlyUsage(cif, transactionType)
		if err != nil {
			return fmt.Errorf("failed to check monthly usage: %v", err)
		}
		if monthlyUsage+amount > limit.MonthlyLimit {
			return fmt.Errorf("transaction would exceed monthly limit of Rp %.0f (used: Rp %.0f)", limit.MonthlyLimit, monthlyUsage)
		}
	}

	return nil
}

// GetDailyUsage returns total transaction amount for today
func GetDailyUsage(cif, transactionType string) (float64, error) {
	today := time.Now().Format("2006-01-02")

	// Map transaction type to actual transaction_type in transactions table
	txTypeMap := map[string]string{
		"intra_transfer": "intra_transfer",
		"inter_transfer": "inter_transfer",
		"bill_payment":   "bill_payment",
		"qris_payment":   "qris_payment",
		"va_payment":     "va_payment",
		"ewallet_topup":  "ewallet_topup",
		"emoney_topup":   "emoney_topup",
	}

	dbTxType := transactionType
	if mapped, ok := txTypeMap[transactionType]; ok {
		dbTxType = mapped
	}

	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM transactions 
	          WHERE from_account_number IN (SELECT account_number FROM accounts WHERE cif = ?)
	          AND transaction_type = ? AND DATE(transaction_date) = ? AND status = 'success'`
	err := database.DB.QueryRow(query, cif, dbTxType, today).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetMonthlyUsage returns total transaction amount for current month
func GetMonthlyUsage(cif, transactionType string) (float64, error) {
	monthStart := time.Now().Format("2006-01") + "-01"

	dbTxType := transactionType
	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM transactions 
	          WHERE from_account_number IN (SELECT account_number FROM accounts WHERE cif = ?)
	          AND transaction_type = ? AND DATE(transaction_date) >= ? AND status = 'success'`
	err := database.DB.QueryRow(query, cif, dbTxType, monthStart).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetTransactionLimits returns all limits for a role
func GetTransactionLimits(roleName string) ([]models.TransactionLimit, error) {
	query := `SELECT id, role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit, 
	          is_active, created_at, updated_at FROM transaction_limits WHERE role_name = ?`
	rows, err := database.DB.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []models.TransactionLimit
	for rows.Next() {
		var l models.TransactionLimit
		if err := rows.Scan(&l.ID, &l.RoleName, &l.TransactionType, &l.DailyLimit,
			&l.PerTransactionLimit, &l.MonthlyLimit, &l.IsActive, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		limits = append(limits, l)
	}
	return limits, nil
}

// GetAllTransactionLimits returns all limits
func GetAllTransactionLimits() ([]models.TransactionLimit, error) {
	query := `SELECT id, role_name, transaction_type, daily_limit, per_transaction_limit, monthly_limit, 
	          is_active, created_at, updated_at FROM transaction_limits ORDER BY role_name, transaction_type`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []models.TransactionLimit
	for rows.Next() {
		var l models.TransactionLimit
		if err := rows.Scan(&l.ID, &l.RoleName, &l.TransactionType, &l.DailyLimit,
			&l.PerTransactionLimit, &l.MonthlyLimit, &l.IsActive, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		limits = append(limits, l)
	}
	return limits, nil
}

// UpdateTransactionLimit updates a specific limit
func UpdateTransactionLimit(id int, dailyLimit, perTxLimit, monthlyLimit float64) error {
	query := `UPDATE transaction_limits SET daily_limit = ?, per_transaction_limit = ?, monthly_limit = ?, 
	          updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := database.DB.Exec(query, dailyLimit, perTxLimit, monthlyLimit, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("transaction limit not found")
	}
	return nil
}
