package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
	"time"
)

// CreateSIRequest represents a request to create a standing instruction
type CreateSIRequest struct {
	CIF             string  `json:"cif" binding:"required"`
	FromAccount     string  `json:"from_account" binding:"required"`
	InstructionType string  `json:"instruction_type" binding:"required"`
	ToAccount       string  `json:"to_account"`
	ToBankCode      string  `json:"to_bank_code"`
	Amount          float64 `json:"amount" binding:"required"`
	Description     string  `json:"description"`
	Frequency       string  `json:"frequency" binding:"required"`
	ExecutionDay    int     `json:"execution_day"`
	StartDate       string  `json:"start_date" binding:"required"`
	EndDate         string  `json:"end_date"`
}

// CreateStandingInstruction creates a new recurring payment instruction
func CreateStandingInstruction(req CreateSIRequest) (*models.StandingInstruction, error) {
	// Validate frequency
	validFreqs := map[string]bool{"daily": true, "weekly": true, "monthly": true, "quarterly": true}
	if !validFreqs[req.Frequency] {
		return nil, fmt.Errorf("invalid frequency: %s (must be daily, weekly, monthly, quarterly)", req.Frequency)
	}

	// Validate account exists
	var accountCount int
	database.DB.QueryRow(`SELECT COUNT(*) FROM accounts WHERE account_number = ? AND cif = ? AND status = 'active'`,
		req.FromAccount, req.CIF).Scan(&accountCount)
	if accountCount == 0 {
		return nil, fmt.Errorf("source account not found or inactive")
	}

	// Generate SI number
	siNumber := fmt.Sprintf("SI-%s-%06d", time.Now().UTC().Format("20060102"), time.Now().UnixNano()%1000000)

	// Calculate next execution date
	nextExec := req.StartDate
	if req.ExecutionDay == 0 {
		req.ExecutionDay = 1
	}

	query := `INSERT INTO standing_instructions 
	          (si_number, cif, from_account, instruction_type, to_account, to_bank_code, amount, description,
	           frequency, execution_day, start_date, end_date, next_execution_date, status)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'active')`

	var endDate interface{} = nil
	if req.EndDate != "" {
		endDate = req.EndDate
	}

	_, err := database.DB.Exec(query, siNumber, req.CIF, req.FromAccount, req.InstructionType,
		req.ToAccount, req.ToBankCode, req.Amount, req.Description,
		req.Frequency, req.ExecutionDay, req.StartDate, endDate, nextExec)
	if err != nil {
		return nil, fmt.Errorf("failed to create standing instruction: %v", err)
	}

	return &models.StandingInstruction{
		SINumber:          siNumber,
		CIF:               req.CIF,
		FromAccount:       req.FromAccount,
		InstructionType:   req.InstructionType,
		ToAccount:         req.ToAccount,
		Amount:            req.Amount,
		Frequency:         req.Frequency,
		NextExecutionDate: nextExec,
		Status:            "active",
	}, nil
}

// ExecutePendingSI executes all standing instructions due on the given date
func ExecutePendingSI(date string) (int, int, error) {
	rows, err := database.DB.Query(`SELECT si_number, cif, from_account, instruction_type, to_account, to_bank_code, amount, description, frequency, execution_day
	                                 FROM standing_instructions 
	                                 WHERE status = 'active' AND next_execution_date <= ?`, date)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var success, failed int

	for rows.Next() {
		var si models.StandingInstruction
		rows.Scan(&si.SINumber, &si.CIF, &si.FromAccount, &si.InstructionType,
			&si.ToAccount, &si.ToBankCode, &si.Amount, &si.Description, &si.Frequency, &si.ExecutionDay)

		// Try to execute the instruction
		txID := fmt.Sprintf("SI-%s-%s", si.SINumber, date)
		execStatus := "success"
		var errMsg string

		// Check balance
		var balance float64
		database.DB.QueryRow(`SELECT avail_balance FROM accounts WHERE account_number = ?`, si.FromAccount).Scan(&balance)

		if balance < si.Amount {
			execStatus = "failed"
			errMsg = "insufficient balance"
			failed++
		} else {
			// Execute transfer
			_, err := database.DB.Exec(`UPDATE accounts SET balance = balance - ?, avail_balance = avail_balance - ?, updated_at = CURRENT_TIMESTAMP 
			                            WHERE account_number = ?`, si.Amount, si.Amount, si.FromAccount)
			if err != nil {
				execStatus = "failed"
				errMsg = err.Error()
				failed++
			} else {
				// Credit to target if intra-bank
				if si.ToAccount != "" && si.ToBankCode == "" {
					database.DB.Exec(`UPDATE accounts SET balance = balance + ?, avail_balance = avail_balance + ?, updated_at = CURRENT_TIMESTAMP 
					                  WHERE account_number = ?`, si.Amount, si.Amount, si.ToAccount)
				}

				// GL entry
				today := time.Now().UTC().Format("2006-01-02")
				CreateJournalEntry(today, fmt.Sprintf("SI execution %s: %s", si.SINumber, si.Description),
					"si_transfer", txID, "SYSTEM", []JournalLineInput{
						{AccountCode: "211", DebitAmount: si.Amount, Description: fmt.Sprintf("Debit SI %s", si.FromAccount)},
						{AccountCode: "211", CreditAmount: si.Amount, Description: fmt.Sprintf("Credit SI %s", si.ToAccount)},
					})

				success++
			}
		}

		// Record execution
		database.DB.Exec(`INSERT INTO si_executions (si_number, execution_date, amount, transaction_id, status, error_message)
		                  VALUES (?, ?, ?, ?, ?, ?)`, si.SINumber, date, si.Amount, txID, execStatus, errMsg)

		// Update SI record
		nextExec := calculateNextExecution(date, si.Frequency, si.ExecutionDay)
		if execStatus == "success" {
			database.DB.Exec(`UPDATE standing_instructions SET total_executed = total_executed + 1, last_execution_date = ?, last_status = ?, next_execution_date = ? WHERE si_number = ?`,
				date, execStatus, nextExec, si.SINumber)
		} else {
			database.DB.Exec(`UPDATE standing_instructions SET total_failed = total_failed + 1, last_execution_date = ?, last_status = ?, next_execution_date = ? WHERE si_number = ?`,
				date, execStatus, nextExec, si.SINumber)
		}
	}

	return success, failed, nil
}

// calculateNextExecution computes the next execution date based on frequency
func calculateNextExecution(currentDate, frequency string, executionDay int) string {
	t, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return currentDate
	}

	switch frequency {
	case "daily":
		t = t.AddDate(0, 0, 1)
	case "weekly":
		t = t.AddDate(0, 0, 7)
	case "monthly":
		t = t.AddDate(0, 1, 0)
		if executionDay > 0 {
			t = time.Date(t.Year(), t.Month(), executionDay, 0, 0, 0, 0, time.UTC)
		}
	case "quarterly":
		t = t.AddDate(0, 3, 0)
		if executionDay > 0 {
			t = time.Date(t.Year(), t.Month(), executionDay, 0, 0, 0, 0, time.UTC)
		}
	}

	return t.Format("2006-01-02")
}

// GetSIByCIF returns all standing instructions for a customer
func GetSIByCIF(cif string) ([]models.StandingInstruction, error) {
	rows, err := database.DB.Query(`SELECT id, si_number, cif, from_account, instruction_type, to_account, to_bank_code,
	                                 amount, description, frequency, execution_day, start_date, end_date,
	                                 next_execution_date, total_executed, total_failed, last_execution_date, last_status, status, created_at
	                                 FROM standing_instructions WHERE cif = ? ORDER BY created_at DESC`, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructions []models.StandingInstruction
	for rows.Next() {
		var si models.StandingInstruction
		rows.Scan(&si.ID, &si.SINumber, &si.CIF, &si.FromAccount, &si.InstructionType,
			&si.ToAccount, &si.ToBankCode, &si.Amount, &si.Description, &si.Frequency,
			&si.ExecutionDay, &si.StartDate, &si.EndDate, &si.NextExecutionDate,
			&si.TotalExecuted, &si.TotalFailed, &si.LastExecutionDate, &si.LastStatus,
			&si.Status, &si.CreatedAt)
		instructions = append(instructions, si)
	}

	return instructions, nil
}

// PauseSI pauses a standing instruction
func PauseSI(siNumber string) error {
	result, err := database.DB.Exec(`UPDATE standing_instructions SET status = 'paused' WHERE si_number = ? AND status = 'active'`, siNumber)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("standing instruction not found or not active")
	}
	return nil
}

// ResumeSI resumes a paused standing instruction
func ResumeSI(siNumber string) error {
	result, err := database.DB.Exec(`UPDATE standing_instructions SET status = 'active' WHERE si_number = ? AND status = 'paused'`, siNumber)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("standing instruction not found or not paused")
	}
	return nil
}

// CancelSI cancels a standing instruction
func CancelSI(siNumber string) error {
	result, err := database.DB.Exec(`UPDATE standing_instructions SET status = 'cancelled' WHERE si_number = ? AND status IN ('active', 'paused')`, siNumber)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("standing instruction not found or already cancelled")
	}
	return nil
}

// GetSIExecutionHistory returns execution history for a standing instruction
func GetSIExecutionHistory(siNumber string) ([]models.SIExecution, error) {
	rows, err := database.DB.Query(`SELECT id, si_number, execution_date, amount, transaction_id, status, error_message, created_at
	                                 FROM si_executions WHERE si_number = ? ORDER BY execution_date DESC`, siNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var execs []models.SIExecution
	for rows.Next() {
		var e models.SIExecution
		rows.Scan(&e.ID, &e.SINumber, &e.ExecutionDate, &e.Amount, &e.TransactionID, &e.Status, &e.ErrorMessage, &e.CreatedAt)
		execs = append(execs, e)
	}

	return execs, nil
}
