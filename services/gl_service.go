package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
	"math"
	"time"
)

// JournalLineInput represents input for creating a journal line
type JournalLineInput struct {
	AccountCode  string  `json:"account_code"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	Description  string  `json:"description"`
}

// CreateJournalEntry creates a double-entry journal posting
func CreateJournalEntry(entryDate, description, refType, refID, postedBy string, lines []JournalLineInput) (*models.JournalEntry, error) {
	// Validate: total debit must equal total credit
	var totalDebit, totalCredit float64
	for _, line := range lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}

	if math.Abs(totalDebit-totalCredit) > 0.01 {
		return nil, fmt.Errorf("journal not balanced: debit=%.2f, credit=%.2f", totalDebit, totalCredit)
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("journal must have at least 2 lines")
	}

	// Generate journal number
	journalNumber := fmt.Sprintf("JRN-%s-%s", time.Now().UTC().Format("20060102"), fmt.Sprintf("%06d", time.Now().UnixNano()%1000000))

	// Begin transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert journal entry
	result, err := tx.Exec(`INSERT INTO journal_entries (journal_number, entry_date, description, reference_type, reference_id, posted_by, status) 
	                         VALUES (?, ?, ?, ?, ?, ?, 'posted')`,
		journalNumber, entryDate, description, refType, refID, postedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create journal entry: %v", err)
	}

	journalID, _ := result.LastInsertId()

	// Insert journal lines
	for _, line := range lines {
		_, err := tx.Exec(`INSERT INTO journal_lines (journal_id, account_code, debit_amount, credit_amount, description) 
		                    VALUES (?, ?, ?, ?, ?)`,
			journalID, line.AccountCode, line.DebitAmount, line.CreditAmount, line.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to create journal line: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit journal: %v", err)
	}

	entry := &models.JournalEntry{
		ID:            int(journalID),
		JournalNumber: journalNumber,
		EntryDate:     entryDate,
		Description:   description,
		ReferenceType: refType,
		ReferenceID:   refID,
		PostedBy:      postedBy,
		Status:        "posted",
	}

	return entry, nil
}

// GetJournalEntries retrieves journal entries with optional filters
func GetJournalEntries(dateFrom, dateTo, refType string, page, pageSize int) ([]models.JournalEntry, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	query := `SELECT id, journal_number, entry_date, description, reference_type, reference_id, posted_by, status, created_at 
	          FROM journal_entries WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM journal_entries WHERE 1=1`
	args := []interface{}{}

	if dateFrom != "" {
		query += " AND entry_date >= ?"
		countQuery += " AND entry_date >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND entry_date <= ?"
		countQuery += " AND entry_date <= ?"
		args = append(args, dateTo)
	}
	if refType != "" {
		query += " AND reference_type = ?"
		countQuery += " AND reference_type = ?"
		args = append(args, refType)
	}

	var total int
	database.DB.QueryRow(countQuery, args...).Scan(&total)

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []models.JournalEntry
	for rows.Next() {
		var e models.JournalEntry
		err := rows.Scan(&e.ID, &e.JournalNumber, &e.EntryDate, &e.Description,
			&e.ReferenceType, &e.ReferenceID, &e.PostedBy, &e.Status, &e.CreatedAt)
		if err != nil {
			continue
		}
		entries = append(entries, e)
	}

	return entries, total, nil
}

// GetJournalLines retrieves the lines for a journal entry
func GetJournalLines(journalID int) ([]models.JournalLine, error) {
	query := `SELECT jl.id, jl.journal_id, jl.account_code, coa.account_name, jl.debit_amount, jl.credit_amount, jl.description
	          FROM journal_lines jl
	          JOIN chart_of_accounts coa ON jl.account_code = coa.account_code
	          WHERE jl.journal_id = ?`

	rows, err := database.DB.Query(query, journalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []models.JournalLine
	for rows.Next() {
		var l models.JournalLine
		err := rows.Scan(&l.ID, &l.JournalID, &l.AccountCode, &l.AccountName, &l.DebitAmount, &l.CreditAmount, &l.Description)
		if err != nil {
			continue
		}
		lines = append(lines, l)
	}

	return lines, nil
}

// GetAccountBalance returns the balance for a GL account
func GetGLAccountBalance(accountCode string) (float64, error) {
	query := `SELECT COALESCE(SUM(debit_amount), 0) - COALESCE(SUM(credit_amount), 0)
	          FROM journal_lines jl
	          JOIN journal_entries je ON jl.journal_id = je.id
	          WHERE jl.account_code = ? AND je.status = 'posted'`

	var balance float64
	err := database.DB.QueryRow(query, accountCode).Scan(&balance)
	if err != nil {
		return 0, err
	}

	// For liability/equity/revenue accounts, flip the sign
	var normalBalance string
	database.DB.QueryRow(`SELECT normal_balance FROM chart_of_accounts WHERE account_code = ?`, accountCode).Scan(&normalBalance)
	if normalBalance == "credit" {
		balance = -balance
	}

	return balance, nil
}

// GetTrialBalance generates a trial balance report
func GetTrialBalance(asOfDate string) ([]models.TrialBalance, error) {
	query := `SELECT coa.account_code, coa.account_name, coa.account_type,
	              COALESCE(SUM(jl.debit_amount), 0) as total_debit,
	              COALESCE(SUM(jl.credit_amount), 0) as total_credit
	          FROM chart_of_accounts coa
	          LEFT JOIN journal_lines jl ON coa.account_code = jl.account_code
	          LEFT JOIN journal_entries je ON jl.journal_id = je.id AND je.status = 'posted'`

	args := []interface{}{}
	if asOfDate != "" {
		query += ` AND je.entry_date <= ?`
		args = append(args, asOfDate)
	}

	query += ` WHERE coa.level = 3 AND coa.is_active = 1
	           GROUP BY coa.account_code, coa.account_name, coa.account_type
	           HAVING total_debit > 0 OR total_credit > 0
	           ORDER BY coa.account_code`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []models.TrialBalance
	for rows.Next() {
		var tb models.TrialBalance
		err := rows.Scan(&tb.AccountCode, &tb.AccountName, &tb.AccountType, &tb.DebitBalance, &tb.CreditBalance)
		if err != nil {
			continue
		}
		report = append(report, tb)
	}

	return report, nil
}

// GetChartOfAccounts returns the chart of accounts with optional type filter
func GetChartOfAccounts(accountType string) ([]models.ChartOfAccount, error) {
	query := `SELECT id, account_code, account_name, account_type, parent_code, level, normal_balance, is_active, created_at
	          FROM chart_of_accounts WHERE is_active = 1`
	args := []interface{}{}

	if accountType != "" {
		query += " AND account_type = ?"
		args = append(args, accountType)
	}

	query += " ORDER BY account_code"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.ChartOfAccount
	for rows.Next() {
		var a models.ChartOfAccount
		err := rows.Scan(&a.ID, &a.AccountCode, &a.AccountName, &a.AccountType,
			&a.ParentCode, &a.Level, &a.NormalBalance, &a.IsActive, &a.CreatedAt)
		if err != nil {
			continue
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}

// ReverseJournalEntry creates a reversal entry for an existing journal
func ReverseJournalEntry(journalID int, reversedBy string) (*models.JournalEntry, error) {
	// Get original lines
	lines, err := GetJournalLines(journalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original journal lines: %v", err)
	}

	// Reverse: swap debit and credit
	var reversedLines []JournalLineInput
	for _, l := range lines {
		reversedLines = append(reversedLines, JournalLineInput{
			AccountCode:  l.AccountCode,
			DebitAmount:  l.CreditAmount,
			CreditAmount: l.DebitAmount,
			Description:  fmt.Sprintf("Reversal: %s", l.Description),
		})
	}

	today := time.Now().UTC().Format("2006-01-02")
	entry, err := CreateJournalEntry(today, fmt.Sprintf("Reversal of journal ID %d", journalID), "reversal", fmt.Sprintf("%d", journalID), reversedBy, reversedLines)
	if err != nil {
		return nil, err
	}

	// Mark original as reversed
	database.DB.Exec(`UPDATE journal_entries SET status = 'reversed', reversed_by = ? WHERE id = ?`, entry.ID, journalID)

	return entry, nil
}
