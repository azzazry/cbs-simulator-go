package services

import (
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

// TransferRequest represents a fund transfer request
type TransferRequest struct {
	FromAccountNumber string  `json:"from_account_number"`
	ToAccountNumber   string  `json:"to_account_number"`
	Amount            float64 `json:"amount"`
	Description       string  `json:"description"`
	TransferType      string  `json:"transfer_type"` // intra, inter, rtgs, skn
}

// ProcessIntraBankTransfer processes intrabank transfer (within same bank)
func ProcessIntraBankTransfer(req TransferRequest) (*models.Transaction, error) {
	// Validate accounts
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	toAccount, err := GetAccountBalance(req.ToAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("destination account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	if toAccount.Status != "active" {
		return nil, fmt.Errorf("destination account is not active")
	}

	// Check balance
	if fromAccount.AvailBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Debit from source account
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - ?, avail_balance = avail_balance - ?, 
	                  updated_at = CURRENT_TIMESTAMP WHERE account_number = ?`,
		req.Amount, req.Amount, req.FromAccountNumber)
	if err != nil {
		return nil, err
	}

	// Credit to destination account
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + ?, avail_balance = avail_balance + ?, 
	                  updated_at = CURRENT_TIMESTAMP WHERE account_number = ?`,
		req.Amount, req.Amount, req.ToAccountNumber)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transactionID := utils.GenerateTransactionID()
	referenceNumber := utils.GenerateReferenceNumber()
	settlementDate := utils.GetCurrentDate()

	result, err := tx.Exec(`INSERT INTO transactions (transaction_id, transaction_type, 
	                        from_account_number, to_account_number, amount, currency, 
	                        description, reference_number, status, settlement_date, fee) 
	                        VALUES (?, 'transfer_intra', ?, ?, ?, 'IDR', ?, ?, 'success', ?, 0.00)`,
		transactionID, req.FromAccountNumber, req.ToAccountNumber, req.Amount,
		req.Description, referenceNumber, settlementDate)

	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Simulate processing delay
	utils.SimulateDelay(500)

	// Retrieve transaction
	trxID, _ := result.LastInsertId()
	transaction, _ := GetTransactionByID(int(trxID))

	// Trigger notifications for sender and receiver
	fromCIF, _ := GetCIFByAccountNumber(req.FromAccountNumber)
	toCIF, _ := GetCIFByAccountNumber(req.ToAccountNumber)

	if fromCIF != "" {
		SaveNotification(fromCIF, "transfer", "Transfer Berhasil",
			fmt.Sprintf("Anda mengirim Rp %.0f ke %s", req.Amount, req.ToAccountNumber),
			transactionID)
		SendPushNotification(fromCIF, "Transfer Berhasil",
			fmt.Sprintf("Rp %.0f telah ditransfer ke %s", req.Amount, req.ToAccountNumber))
	}

	if toCIF != "" {
		SaveNotification(toCIF, "transfer", "Uang Masuk",
			fmt.Sprintf("Anda menerima Rp %.0f dari %s", req.Amount, req.FromAccountNumber),
			transactionID)
		SendPushNotification(toCIF, "Uang Masuk",
			fmt.Sprintf("Rp %.0f diterima dari %s", req.Amount, req.FromAccountNumber))
	}

	return transaction, nil
}

// ProcessInterBankTransfer processes interbank transfer (to other banks)
func ProcessInterBankTransfer(req TransferRequest) (*models.Transaction, error) {
	// Validate source account
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	// Fee for interbank transfer
	var fee float64
	switch req.TransferType {
	case "rtgs":
		fee = 25000.00 // RTGS fee
	case "skn":
		fee = 3500.00 // SKN fee
	default:
		fee = 6500.00 // Default interbank fee
	}

	totalAmount := req.Amount + fee

	// Check balance
	if fromAccount.AvailBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance (including fee)")
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Debit from source account (amount + fee)
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - ?, avail_balance = avail_balance - ?, 
	                  updated_at = CURRENT_TIMESTAMP WHERE account_number = ?`,
		totalAmount, totalAmount, req.FromAccountNumber)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transactionID := utils.GenerateTransactionID()
	referenceNumber := utils.GenerateReferenceNumber()
	settlementDate := utils.GetCurrentDate()

	// For interbank, settlement may be next day
	if req.TransferType == "skn" {
		settlementDate, _ = utils.AddMonths(settlementDate, 0)
	}

	result, err := tx.Exec(`INSERT INTO transactions (transaction_id, transaction_type, 
	                        from_account_number, to_account_number, amount, currency, 
	                        description, reference_number, status, settlement_date, fee) 
	                        VALUES (?, ?, ?, ?, ?, 'IDR', ?, ?, 'success', ?, ?)`,
		transactionID, "transfer_inter", req.FromAccountNumber, req.ToAccountNumber,
		req.Amount, req.Description, referenceNumber, settlementDate, fee)

	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Simulate processing delay (longer for interbank)
	utils.SimulateDelay(1000)

	// Retrieve transaction
	trxID, _ := result.LastInsertId()
	transaction, _ := GetTransactionByID(int(trxID))

	// Trigger notification for sender
	fromCIF, _ := GetCIFByAccountNumber(req.FromAccountNumber)
	if fromCIF != "" {
		SaveNotification(fromCIF, "transfer", "Transfer Antar Bank Berhasil",
			fmt.Sprintf("Rp %.0f + biaya Rp %.0f ditransfer. Settlement: %s", req.Amount, fee, "1-2 hari kerja"),
			transactionID)
		SendPushNotification(fromCIF, "Transfer Antar Bank Berhasil",
			fmt.Sprintf("Rp %.0f telah diproses", req.Amount))
	}

	return transaction, nil
}

// GetTransactionByID retrieves transaction by ID
func GetTransactionByID(id int) (*models.Transaction, error) {
	var trx models.Transaction

	query := `SELECT id, transaction_id, transaction_type, from_account_number, 
	          to_account_number, amount, currency, description, reference_number, 
	          status, transaction_date, settlement_date, fee, created_at 
	          FROM transactions WHERE id = ?`

	err := database.DB.QueryRow(query, id).Scan(
		&trx.ID, &trx.TransactionID, &trx.TransactionType,
		&trx.FromAccountNumber, &trx.ToAccountNumber, &trx.Amount,
		&trx.Currency, &trx.Description, &trx.ReferenceNumber,
		&trx.Status, &trx.TransactionDate, &trx.SettlementDate,
		&trx.Fee, &trx.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &trx, nil
}

// GetTransactionByTransactionID retrieves transaction by transaction ID
func GetTransactionByTransactionID(transactionID string) (*models.Transaction, error) {
	var trx models.Transaction

	query := `SELECT id, transaction_id, transaction_type, from_account_number, 
	          to_account_number, amount, currency, description, reference_number, 
	          status, transaction_date, settlement_date, fee, created_at 
	          FROM transactions WHERE transaction_id = ?`

	err := database.DB.QueryRow(query, transactionID).Scan(
		&trx.ID, &trx.TransactionID, &trx.TransactionType,
		&trx.FromAccountNumber, &trx.ToAccountNumber, &trx.Amount,
		&trx.Currency, &trx.Description, &trx.ReferenceNumber,
		&trx.Status, &trx.TransactionDate, &trx.SettlementDate,
		&trx.Fee, &trx.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &trx, nil
}

// GetCIFByAccountNumber retrieves CIF from account number
func GetCIFByAccountNumber(accountNumber string) (string, error) {
	var cif string
	err := database.DB.QueryRow(`SELECT cif FROM accounts WHERE account_number = ?`, accountNumber).Scan(&cif)
	if err != nil {
		return "", err
	}
	return cif, nil
}
