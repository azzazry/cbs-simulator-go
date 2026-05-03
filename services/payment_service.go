package services

import (
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

type QRISPaymentRequest struct {
	FromAccountNumber string  `json:"from_account_number" binding:"required"`
	MerchantName      string  `json:"merchant_name" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,min=1"`
	Description       string  `json:"description"`
	QRISCode          string  `json:"qris_code" binding:"required"`
}

type VAPaymentRequest struct {
	FromAccountNumber string  `json:"from_account_number" binding:"required"`
	DestinationVACode string  `json:"destination_va_code" binding:"required"`
	VABankCode        string  `json:"va_bank_code" binding:"required"`
	BeneficiaryName   string  `json:"beneficiary_name" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,min=1"`
	Description       string  `json:"description"`
}

type EWalletTopupRequest struct {
	FromAccountNumber string  `json:"from_account_number" binding:"required"`
	EWalletProvider   string  `json:"ewallet_provider" binding:"required"`
	PhoneNumber       string  `json:"phone_number" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,min=1"`
	Description       string  `json:"description"`
}

type EMoneyTopupRequest struct {
	FromAccountNumber string  `json:"from_account_number" binding:"required"`
	EMoneyProvider    string  `json:"emoney_provider" binding:"required"`
	CardNumber        string  `json:"card_number" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,min=1"`
	Description       string  `json:"description"`
}

// ProcessQRISPayment processes QRIS payment transaction
func ProcessQRISPayment(req QRISPaymentRequest) (*models.Transaction, error) {
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	fee, err := CalculateServiceFee("QRIS_PAYMENT", req.Amount)
	if err != nil {
		fee = (req.Amount * 1) / 100
	}

	totalAmount := req.Amount + fee

	if fromAccount.AvailBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance (amount: %.0f + fee: %.0f = %.0f)", req.Amount, fee, totalAmount)
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
		transactionID, "payment_qris", req.FromAccountNumber, req.QRISCode,
		req.Amount, req.Description, referenceNumber, settlementDate, fee).Scan(&trxID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	utils.SimulateDelay(800)

	transaction, _ := GetTransactionByID(int(trxID))

	fromCIF, _ := GetCIFByAccountNumber(req.FromAccountNumber)
	if fromCIF != "" {
		SaveNotification(fromCIF, "payment", "QRIS Payment Berhasil",
			fmt.Sprintf("Pembayaran QRIS ke %s sebesar Rp %.0f telah berhasil", req.MerchantName, req.Amount),
			transactionID)
		SendPushNotification(fromCIF, "QRIS Payment Berhasil",
			fmt.Sprintf("Rp %.0f dibayarkan ke %s", req.Amount, req.MerchantName))
	}

	return transaction, nil
}

// ProcessVAPayment processes Virtual Account payment transaction
func ProcessVAPayment(req VAPaymentRequest) (*models.Transaction, error) {
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	fee := 0.0
	totalAmount := req.Amount + fee

	if fromAccount.AvailBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance")
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
		transactionID, "payment_va", req.FromAccountNumber, req.DestinationVACode,
		req.Amount, req.Description, referenceNumber, settlementDate, fee).Scan(&trxID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	utils.SimulateDelay(800)

	transaction, _ := GetTransactionByID(int(trxID))

	fromCIF, _ := GetCIFByAccountNumber(req.FromAccountNumber)
	if fromCIF != "" {
		SaveNotification(fromCIF, "payment", "Virtual Account Payment Berhasil",
			fmt.Sprintf("Pembayaran ke VA %s atas nama %s sebesar Rp %.0f telah berhasil",
				req.DestinationVACode, req.BeneficiaryName, req.Amount),
			transactionID)
		SendPushNotification(fromCIF, "VA Payment Berhasil",
			fmt.Sprintf("Rp %.0f dibayarkan ke %s", req.Amount, req.BeneficiaryName))
	}

	return transaction, nil
}

// ProcessEWalletTopup processes e-wallet top-up transaction
func ProcessEWalletTopup(req EWalletTopupRequest) (*models.Transaction, error) {
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	serviceCode := fmt.Sprintf("TOPUP_%s", req.EWalletProvider)
	fee, err := CalculateServiceFee(serviceCode, req.Amount)
	if err != nil {
		fee = 2500.0
	}

	totalAmount := req.Amount + fee

	if fromAccount.AvailBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance (amount: %.0f + fee: %.0f = %.0f)", req.Amount, fee, totalAmount)
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
		transactionID, "topup_ewallet", req.FromAccountNumber, req.PhoneNumber,
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
		SaveNotification(fromCIF, "topup", "Top-up E-Wallet Berhasil",
			fmt.Sprintf("Top-up %s sebesar Rp %.0f ke %s telah berhasil", req.EWalletProvider, req.Amount, req.PhoneNumber),
			transactionID)
		SendPushNotification(fromCIF, "Top-up E-Wallet Berhasil",
			fmt.Sprintf("Rp %.0f top-up %s", req.Amount, req.EWalletProvider))
	}

	return transaction, nil
}

// ProcessEMoneyTopup processes e-money top-up transaction
func ProcessEMoneyTopup(req EMoneyTopupRequest) (*models.Transaction, error) {
	fromAccount, err := GetAccountBalance(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found: %v", err)
	}

	if fromAccount.Status != "active" {
		return nil, fmt.Errorf("source account is not active")
	}

	serviceCode := fmt.Sprintf("TOPUP_%s", req.EMoneyProvider)
	fee, err := CalculateServiceFee(serviceCode, req.Amount)
	if err != nil {
		fee = 2500.0
	}

	totalAmount := req.Amount + fee

	if fromAccount.AvailBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance (amount: %.0f + fee: %.0f = %.0f)", req.Amount, fee, totalAmount)
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
		transactionID, "topup_emoney", req.FromAccountNumber, req.CardNumber,
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
		SaveNotification(fromCIF, "topup", "Top-up E-Money Berhasil",
			fmt.Sprintf("Top-up %s sebesar Rp %.0f ke %s telah berhasil", req.EMoneyProvider, req.Amount, req.CardNumber),
			transactionID)
		SendPushNotification(fromCIF, "Top-up E-Money Berhasil",
			fmt.Sprintf("Rp %.0f top-up %s", req.Amount, req.EMoneyProvider))
	}

	return transaction, nil
}
