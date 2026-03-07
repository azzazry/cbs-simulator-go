package services

import (
	"cbs-simulator/database"
	"fmt"
	"time"
)

type OpenAccountRequest struct {
	CIF            string  `json:"cif" binding:"required"`
	AccountType    string  `json:"account_type" binding:"required"`
	Currency       string  `json:"currency"`
	InitialDeposit float64 `json:"initial_deposit"`
	Branch         string  `json:"branch"`
}

type OpenAccountResponse struct {
	AccountNumber  string  `json:"account_number"`
	CIF            string  `json:"cif"`
	AccountType    string  `json:"account_type"`
	InitialDeposit float64 `json:"initial_deposit"`
	Message        string  `json:"message"`
}

func OpenAccount(req OpenAccountRequest) (*OpenAccountResponse, error) {
	// Validate customer
	var count int
	database.DB.QueryRow(`SELECT COUNT(*) FROM customers WHERE cif = ? AND status = 'active'`, req.CIF).Scan(&count)
	if count == 0 {
		return nil, fmt.Errorf("customer not found or inactive")
	}

	if req.Currency == "" {
		req.Currency = "IDR"
	}

	// Generate account number
	var maxID int
	database.DB.QueryRow(`SELECT COALESCE(MAX(id), 0) FROM accounts`).Scan(&maxID)
	accountNumber := fmt.Sprintf("100%07d", maxID+1)

	today := time.Now().UTC().Format("2006-01-02")

	_, err := database.DB.Exec(`INSERT INTO accounts (account_number, cif, account_type, currency, balance, avail_balance, status, opened_date, branch) 
	                            VALUES (?, ?, ?, ?, ?, ?, 'active', ?, ?)`,
		accountNumber, req.CIF, req.AccountType, req.Currency, req.InitialDeposit, req.InitialDeposit, today, req.Branch)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	// GL entry if initial deposit > 0
	if req.InitialDeposit > 0 {
		CreateJournalEntry(today, fmt.Sprintf("Account opening deposit %s", accountNumber),
			"opening", accountNumber, "SYSTEM", []JournalLineInput{
				{AccountCode: "111", DebitAmount: req.InitialDeposit, Description: "Kas masuk setoran awal"},
				{AccountCode: "211", CreditAmount: req.InitialDeposit, Description: fmt.Sprintf("Tabungan %s", accountNumber)},
			})
	}

	return &OpenAccountResponse{
		AccountNumber:  accountNumber,
		CIF:            req.CIF,
		AccountType:    req.AccountType,
		InitialDeposit: req.InitialDeposit,
		Message:        "Account opened successfully",
	}, nil
}

func CloseAccount(accountNumber, reason string) error {
	var balance float64
	var status string
	err := database.DB.QueryRow(`SELECT balance, status FROM accounts WHERE account_number = ?`, accountNumber).Scan(&balance, &status)
	if err != nil {
		return fmt.Errorf("account not found")
	}
	if status == "closed" {
		return fmt.Errorf("account already closed")
	}
	if balance > 0 {
		return fmt.Errorf("account has remaining balance of %.2f. Please withdraw or transfer first", balance)
	}

	today := time.Now().UTC().Format("2006-01-02")
	_, err = database.DB.Exec(`UPDATE accounts SET status = 'closed', updated_at = CURRENT_TIMESTAMP WHERE account_number = ?`, accountNumber)
	if err != nil {
		return fmt.Errorf("failed to close account: %v", err)
	}

	CreateJournalEntry(today, fmt.Sprintf("Account closure %s: %s", accountNumber, reason),
		"closing", accountNumber, "SYSTEM", []JournalLineInput{
			{AccountCode: "211", DebitAmount: 0, Description: "Account closure"},
			{AccountCode: "211", CreditAmount: 0, Description: "Account closure"},
		})

	return nil
}

func GetDormantAccounts() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`SELECT a.account_number, a.cif, c.full_name, a.account_type, a.balance, a.status, a.updated_at
	                                 FROM accounts a JOIN customers c ON a.cif = c.cif WHERE a.status = 'dormant' ORDER BY a.updated_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var acctNum, cif, name, acctType, acctStatus, updatedAt string
		var balance float64
		rows.Scan(&acctNum, &cif, &name, &acctType, &balance, &acctStatus, &updatedAt)
		results = append(results, map[string]interface{}{
			"account_number": acctNum, "cif": cif, "full_name": name,
			"account_type": acctType, "balance": balance, "status": acctStatus, "updated_at": updatedAt,
		})
	}
	return results, nil
}

func ReactivateAccount(accountNumber string) error {
	result, err := database.DB.Exec(`UPDATE accounts SET status = 'active', updated_at = CURRENT_TIMESTAMP WHERE account_number = ? AND status = 'dormant'`, accountNumber)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account not found or not dormant")
	}
	return nil
}
