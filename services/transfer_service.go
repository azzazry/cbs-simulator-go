package services

import (
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

// TransferRequest represents a fund transfer request
type TransferRequest struct {
	FromAccountNumber   string  `json:"from_account_number"`
	ToAccountNumber     string  `json:"to_account_number"`
	Amount              float64 `json:"amount"`
	Description         string  `json:"description"`
	TransferType        string  `json:"transfer_type"` // intra, inter, rtgs, skn
	DestinationBankCode string  `json:"destination_bank_code"` // For interbank transfers
}

// ProcessIntraBankTransfer processes intrabank transfer (within same bank)
func ProcessIntraBankTransfer(req TransferRequest) (*models.Transaction, error) {
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

	if fromAccount.AvailBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1, avail_balance = avail_balance - $2,
	                  updated_at = NOW() WHERE account_number = $3`,
		req.Amount, req.Amount, req.FromAccountNumber)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1, avail_balance = avail_balance + $2,
	                  updated_at = NOW() WHERE account_number = $3`,
		req.Amount, req.Amount, req.ToAccountNumber)
	if err != nil {
		return nil, err
	}

	transactionID := utils.GenerateTransactionID()
	referenceNumber := utils.GenerateReferenceNumber()
	settlementDate := utils.GetCurrentDate()

	var trxID int64
	err = tx.QueryRow(`INSERT INTO transactions (transaction_id, transaction_type,
	                   from_account_number, to_account_number, amount, currency,
	                   description, reference_number, status, settlement_date, fee)
	                   VALUES ($1, 'transfer_intra', $2, $3, $4, 'IDR', $5, $6, 'success', $7, 0.00)
	                   RETURNING id`,
		transactionID, req.FromAccountNumber, req.ToAccountNumber, req.Amount,
		req.Description, referenceNumber, settlementDate).Scan(&trxID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	utils.SimulateDelay(500)

	transaction, _ := GetTransactionByID(int(trxID))

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
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	if req.DestinationBankCode == "" {
		return nil, fmt.Errorf("destination_bank_code is required for interbank transfers")
	}

	fee, err := CalculateTransferFee(req.DestinationBankCode, req.Amount)
	if err != nil {
		fee = 5000.00
	}

	totalAmount := req.Amount + fee

	if fromAccount.AvailBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance (transfer: %.0f + fee: %.0f = %.0f)", req.Amount, fee, totalAmount)
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1, avail_balance = avail_balance - $2,
	                  updated_at = NOW() WHERE account_number = $3`,
		totalAmount, totalAmount, req.FromAccountNumber)
	if err != nil {
		return nil, err
	}

	transactionID := utils.GenerateTransactionID()
	referenceNumber := utils.GenerateReferenceNumber()
	settlementDate := utils.GetCurrentDate()

	var trxID int64
	err = tx.QueryRow(`INSERT INTO transactions (transaction_id, transaction_type,
	                   from_account_number, to_account_number, amount, currency,
	                   description, reference_number, status, settlement_date, fee)
	                   VALUES ($1, $2, $3, $4, $5, 'IDR', $6, $7, 'success', $8, $9)
	                   RETURNING id`,
		transactionID, "transfer_inter", req.FromAccountNumber, req.ToAccountNumber,
		req.Amount, req.Description, referenceNumber, settlementDate, fee).Scan(&trxID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	utils.SimulateDelay(1000)

	transaction, _ := GetTransactionByID(int(trxID))

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
	          FROM transactions WHERE id = $1`

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

// GetTransactionByTransactionID retrieves transaction by transaction ID string
func GetTransactionByTransactionID(transactionID string) (*models.Transaction, error) {
	var trx models.Transaction

	query := `SELECT id, transaction_id, transaction_type, from_account_number,
	          to_account_number, amount, currency, description, reference_number,
	          status, transaction_date, settlement_date, fee, created_at
	          FROM transactions WHERE transaction_id = $1`

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
	err := database.DB.QueryRow(`SELECT cif FROM accounts WHERE account_number = $1`, accountNumber).Scan(&cif)
	if err != nil {
		return "", err
	}
	return cif, nil
}
