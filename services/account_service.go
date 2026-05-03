package services

import (
	"database/sql"
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

func GetAccountBalance(accountNumber string) (*models.Account, error) {
	var account models.Account

	query := `SELECT id, account_number, cif, account_type, currency, balance,
	          avail_balance, status, opened_date, branch, created_at, updated_at
	          FROM accounts WHERE account_number = $1 AND status = 'active'`

	err := database.DB.QueryRow(query, accountNumber).Scan(
		&account.ID, &account.AccountNumber, &account.CIF, &account.AccountType,
		&account.Currency, &account.Balance, &account.AvailBalance, &account.Status,
		&account.OpenedDate, &account.Branch, &account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}

	return &account, nil
}

func GetAccountsByCIF(cif string) ([]models.Account, error) {
	query := `SELECT id, account_number, cif, account_type, currency, balance,
	          avail_balance, status, opened_date, branch, created_at, updated_at
	          FROM accounts WHERE cif = $1`

	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID, &account.AccountNumber, &account.CIF, &account.AccountType,
			&account.Currency, &account.Balance, &account.AvailBalance, &account.Status,
			&account.OpenedDate, &account.Branch, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetAccountStatement(accountNumber string, limit int, offset int) ([]models.Transaction, error) {
	query := `SELECT id, transaction_id, transaction_type, from_account_number,
	          to_account_number, amount, currency, description, reference_number,
	          status, transaction_date, settlement_date, fee, created_at
	          FROM transactions
	          WHERE (from_account_number = $1 OR to_account_number = $2)
	          AND status = 'success'
	          ORDER BY transaction_date DESC
	          LIMIT $3 OFFSET $4`

	rows, err := database.DB.Query(query, accountNumber, accountNumber, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var trx models.Transaction
		err := rows.Scan(
			&trx.ID, &trx.TransactionID, &trx.TransactionType,
			&trx.FromAccountNumber, &trx.ToAccountNumber, &trx.Amount,
			&trx.Currency, &trx.Description, &trx.ReferenceNumber,
			&trx.Status, &trx.TransactionDate, &trx.SettlementDate,
			&trx.Fee, &trx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, trx)
	}

	return transactions, nil
}

func UpdateAccountBalance(accountNumber string, amount float64) error {
	query := `UPDATE accounts
	          SET balance = balance + $1,
	              avail_balance = avail_balance + $2,
	              updated_at = NOW()
	          WHERE account_number = $3`

	_, err := database.DB.Exec(query, amount, amount, accountNumber)
	return err
}

func CheckSufficientBalance(accountNumber string, amount float64) (bool, error) {
	account, err := GetAccountBalance(accountNumber)
	if err != nil {
		return false, err
	}

	return account.AvailBalance >= amount, nil
}

func CreateAccount(cif, accountType, currency, branch string, initialBalance float64) (*models.Account, error) {
	accountNumber := utils.GenerateAccountNumber(accountType)
	openDate := utils.GetCurrentDate()

	query := `INSERT INTO accounts (account_number, cif, account_type, currency,
	          balance, avail_balance, status, opened_date, branch)
	          VALUES ($1, $2, $3, $4, $5, $6, 'active', $7, $8)`

	_, err := database.DB.Exec(query, accountNumber, cif, accountType, currency,
		initialBalance, initialBalance, openDate, branch)

	if err != nil {
		return nil, err
	}

	return GetAccountBalance(accountNumber)
}
