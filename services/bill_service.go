package services

import (
	"database/sql"
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
	"cbs-simulator/utils"
)

// GetBillByCustomerNumber retrieves unpaid bill by customer number
func GetBillByCustomerNumber(billerCode, customerNumber string) (*models.BillPayment, error) {
	var bill models.BillPayment
	
	query := `SELECT id, biller_code, biller_name, customer_number, bill_number, 
	          bill_amount, admin_fee, total_amount, bill_period, due_date, status, 
	          transaction_id, payment_date, created_at, updated_at 
	          FROM bill_payments 
	          WHERE biller_code = ? AND customer_number = ? AND status = 'unpaid' 
	          ORDER BY due_date ASC LIMIT 1`
	
	err := database.DB.QueryRow(query, billerCode, customerNumber).Scan(
		&bill.ID, &bill.BillerCode, &bill.BillerName, &bill.CustomerNumber,
		&bill.BillNumber, &bill.BillAmount, &bill.AdminFee, &bill.TotalAmount,
		&bill.BillPeriod, &bill.DueDate, &bill.Status, &bill.TransactionID,
		&bill.PaymentDate, &bill.CreatedAt, &bill.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no unpaid bill found for customer")
		}
		return nil, err
	}
	
	return &bill, nil
}

// GetAllUnpaidBills retrieves all unpaid bills
func GetAllUnpaidBills() ([]models.BillPayment, error) {
	query := `SELECT id, biller_code, biller_name, customer_number, bill_number, 
	          bill_amount, admin_fee, total_amount, bill_period, due_date, status, 
	          transaction_id, payment_date, created_at, updated_at 
	          FROM bill_payments 
	          WHERE status = 'unpaid' 
	          ORDER BY due_date ASC`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var bills []models.BillPayment
	for rows.Next() {
		var bill models.BillPayment
		err := rows.Scan(
			&bill.ID, &bill.BillerCode, &bill.BillerName, &bill.CustomerNumber,
			&bill.BillNumber, &bill.BillAmount, &bill.AdminFee, &bill.TotalAmount,
			&bill.BillPeriod, &bill.DueDate, &bill.Status, &bill.TransactionID,
			&bill.PaymentDate, &bill.CreatedAt, &bill.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}
	
	return bills, nil
}

// PayBill processes bill payment
func PayBill(accountNumber, billNumber string) (*models.Transaction, error) {
	// Get bill details
	var bill models.BillPayment
	query := `SELECT id, biller_code, biller_name, customer_number, bill_number, 
	          bill_amount, admin_fee, total_amount, bill_period, due_date, status 
	          FROM bill_payments WHERE bill_number = ? AND status = 'unpaid'`
	
	err := database.DB.QueryRow(query, billNumber).Scan(
		&bill.ID, &bill.BillerCode, &bill.BillerName, &bill.CustomerNumber,
		&bill.BillNumber, &bill.BillAmount, &bill.AdminFee, &bill.TotalAmount,
		&bill.BillPeriod, &bill.DueDate, &bill.Status,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bill not found or already paid")
		}
		return nil, err
	}
	
	// Validate account
	account, err := GetAccountBalance(accountNumber)
	if err != nil {
		return nil, fmt.Errorf("account not found: %v", err)
	}
	
	if account.Status != "active" {
		return nil, fmt.Errorf("account is not active")
	}
	
	// Check balance
	if account.AvailBalance < bill.TotalAmount {
		return nil, fmt.Errorf("insufficient balance")
	}
	
	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	
	// Debit from account
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - ?, avail_balance = avail_balance - ?, 
	                  updated_at = CURRENT_TIMESTAMP WHERE account_number = ?`,
		bill.TotalAmount, bill.TotalAmount, accountNumber)
	if err != nil {
		return nil, err
	}
	
	// Create transaction record
	transactionID := utils.GenerateTransactionID()
	referenceNumber := utils.GenerateReferenceNumber()
	settlementDate := utils.GetCurrentDate()
	description := fmt.Sprintf("Payment for %s - %s", bill.BillerName, bill.BillPeriod)
	
	result, err := tx.Exec(`INSERT INTO transactions (transaction_id, transaction_type, 
	                        from_account_number, to_account_number, amount, currency, 
	                        description, reference_number, status, settlement_date, fee) 
	                        VALUES (?, 'payment_bill', ?, NULL, ?, 'IDR', ?, ?, 'success', ?, ?)`,
		transactionID, accountNumber, bill.TotalAmount, description,
		referenceNumber, settlementDate, bill.AdminFee)
	
	if err != nil {
		return nil, err
	}
	
	// Update bill status
	paymentDate := utils.GetCurrentDate()
	_, err = tx.Exec(`UPDATE bill_payments SET status = 'paid', transaction_id = ?, 
	                  payment_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		transactionID, paymentDate, bill.ID)
	
	if err != nil {
		return nil, err
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	
	// Simulate processing delay
	utils.SimulateDelay(800)
	
	// Retrieve and return transaction
	trxID, _ := result.LastInsertId()
	return GetTransactionByID(int(trxID))
}

// GetBillerList returns list of available billers
func GetBillerList() []map[string]string {
	return []map[string]string{
		{"code": "PLN", "name": "PT PLN (Persero)", "category": "Electricity"},
		{"code": "PDAM", "name": "PDAM", "category": "Water"},
		{"code": "TELKOM", "name": "PT Telkom Indonesia", "category": "Phone"},
		{"code": "INDIHOME", "name": "IndiHome", "category": "Internet"},
		{"code": "BPJS", "name": "BPJS Kesehatan", "category": "Insurance"},
		{"code": "PGN", "name": "PT PGN", "category": "Gas"},
		{"code": "SPEEDY", "name": "Speedy", "category": "Internet"},
		{"code": "FIRSTMEDIA", "name": "First Media", "category": "TV Cable"},
	}
}
