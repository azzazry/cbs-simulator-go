package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"time"
)

type EODResult struct {
	ProcessDate   string             `json:"process_date"`
	Processes     []EODProcessResult `json:"processes"`
	OverallStatus string             `json:"overall_status"`
}

type EODProcessResult struct {
	ProcessType      string `json:"process_type"`
	Status           string `json:"status"`
	RecordsProcessed int    `json:"records_processed"`
	RecordsFailed    int    `json:"records_failed"`
	ErrorMessage     string `json:"error_message,omitempty"`
}

func RunEOD(processDate string) (*EODResult, error) {
	result := &EODResult{ProcessDate: processDate, OverallStatus: "completed"}

	result.Processes = append(result.Processes, runEODProcess(processDate, "interest_accrual", func() (int, int, error) {
		return AccrueInterestForAllAccounts(processDate)
	}))

	result.Processes = append(result.Processes, runEODProcess(processDate, "si_execution", func() (int, int, error) {
		return ExecutePendingSI(processDate)
	}))

	parsedDate, _ := time.Parse("2006-01-02", processDate)
	nextDay := parsedDate.AddDate(0, 0, 1)
	if nextDay.Day() == 1 {
		yearMonth := parsedDate.Format("2006-01")
		result.Processes = append(result.Processes, runEODProcess(processDate, "monthly_interest_posting", func() (int, int, error) {
			posted, err := PostMonthlyInterest(yearMonth)
			return posted, 0, err
		}))
	}

	result.Processes = append(result.Processes, runEODProcess(processDate, "dormant_check", func() (int, int, error) {
		return CheckDormantAccounts()
	}))

	for _, p := range result.Processes {
		if p.Status == "failed" {
			result.OverallStatus = "completed_with_errors"
		}
	}
	return result, nil
}

func runEODProcess(date, processType string, fn func() (int, int, error)) EODProcessResult {
	startedAt := time.Now().UTC()
	database.DB.Exec(`INSERT INTO eod_logs (process_date, process_type, status, started_at) VALUES (?, ?, 'running', ?)`,
		date, processType, startedAt)

	processed, failed, err := fn()
	completedAt := time.Now().UTC()
	status, errMsg := "completed", ""
	if err != nil {
		status, errMsg = "failed", err.Error()
	}

	database.DB.Exec(`UPDATE eod_logs SET status=?, records_processed=?, records_failed=?, completed_at=?, error_message=? WHERE process_date=? AND process_type=? AND status='running'`,
		status, processed, failed, completedAt, errMsg, date, processType)

	return EODProcessResult{ProcessType: processType, Status: status, RecordsProcessed: processed, RecordsFailed: failed, ErrorMessage: errMsg}
}

func CheckDormantAccounts() (int, int, error) {
	cutoff := time.Now().UTC().AddDate(-1, 0, 0).Format("2006-01-02")
	rows, err := database.DB.Query(`SELECT a.account_number FROM accounts a WHERE a.status='active' AND NOT EXISTS (
		SELECT 1 FROM transactions t WHERE (t.from_account_number=a.account_number OR t.to_account_number=a.account_number) AND t.transaction_date > ?) AND a.balance > 0`, cutoff)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	var marked int
	for rows.Next() {
		var acct string
		rows.Scan(&acct)
		database.DB.Exec(`UPDATE accounts SET status='dormant', updated_at=CURRENT_TIMESTAMP WHERE account_number=?`, acct)
		marked++
	}
	return marked, 0, nil
}

func GetEODStatus(date string) ([]models.EODLog, error) {
	rows, err := database.DB.Query(`SELECT id, process_date, process_type, status, records_processed, records_failed, started_at, completed_at, error_message FROM eod_logs WHERE process_date=? ORDER BY id`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []models.EODLog
	for rows.Next() {
		var l models.EODLog
		rows.Scan(&l.ID, &l.ProcessDate, &l.ProcessType, &l.Status, &l.RecordsProcessed, &l.RecordsFailed, &l.StartedAt, &l.CompletedAt, &l.ErrorMessage)
		logs = append(logs, l)
	}
	return logs, nil
}

func GetEODHistory(dateFrom, dateTo string, page, pageSize int) ([]models.EODLog, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	query := `SELECT id, process_date, process_type, status, records_processed, records_failed, started_at, completed_at, error_message FROM eod_logs WHERE 1=1`
	countQ := `SELECT COUNT(*) FROM eod_logs WHERE 1=1`
	args := []interface{}{}
	if dateFrom != "" {
		query += " AND process_date >= ?"
		countQ += " AND process_date >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND process_date <= ?"
		countQ += " AND process_date <= ?"
		args = append(args, dateTo)
	}
	var total int
	database.DB.QueryRow(countQ, args...).Scan(&total)
	query += " ORDER BY process_date DESC, id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var logs []models.EODLog
	for rows.Next() {
		var l models.EODLog
		rows.Scan(&l.ID, &l.ProcessDate, &l.ProcessType, &l.Status, &l.RecordsProcessed, &l.RecordsFailed, &l.StartedAt, &l.CompletedAt, &l.ErrorMessage)
		logs = append(logs, l)
	}
	return logs, total, nil
}
