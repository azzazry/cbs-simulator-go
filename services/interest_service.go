package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
	"math"
	"time"
)

// CalculateDailyInterest calculates and records daily interest for an account
func CalculateDailyInterest(accountNumber, accrualDate string) (*models.InterestAccrual, error) {
	var balance float64
	var accountType string
	err := database.DB.QueryRow(`SELECT balance, account_type FROM accounts WHERE account_number = $1 AND status = 'active'`,
		accountNumber).Scan(&balance, &accountType)
	if err != nil {
		return nil, fmt.Errorf("account not found or inactive")
	}

	if balance <= 0 {
		return nil, fmt.Errorf("no positive balance for interest calculation")
	}

	productType := "savings"
	if accountType == "giro" || accountType == "checking" {
		productType = "checking"
	}

	rate, err := GetApplicableRate(productType, balance, 0)
	if err != nil {
		return nil, fmt.Errorf("no applicable interest rate: %v", err)
	}

	dailyInterest := math.Round((balance*rate/100/365)*10000) / 10000

	// Hitung bulan berjalan menggunakan date range (PostgreSQL tidak support LIKE pada DATE)
	monthStart := accrualDate[:8] + "01"
	var existingAccrued float64
	database.DB.QueryRow(`SELECT COALESCE(SUM(daily_interest), 0) FROM interest_accruals
	                      WHERE account_number = $1 AND accrual_date >= $2 AND accrual_date < $3`,
		accountNumber, monthStart, accrualDate).Scan(&existingAccrued)

	accruedInterest := existingAccrued + dailyInterest

	// INSERT ON CONFLICT menggantikan INSERT OR REPLACE dari SQLite
	_, err = database.DB.Exec(`INSERT INTO interest_accruals
	                           (account_number, accrual_date, product_type, balance, rate, daily_interest, accrued_interest)
	                           VALUES ($1, $2, $3, $4, $5, $6, $7)
	                           ON CONFLICT (account_number, accrual_date) DO UPDATE
	                           SET product_type = EXCLUDED.product_type,
	                               balance = EXCLUDED.balance,
	                               rate = EXCLUDED.rate,
	                               daily_interest = EXCLUDED.daily_interest,
	                               accrued_interest = EXCLUDED.accrued_interest`,
		accountNumber, accrualDate, productType, balance, rate, dailyInterest, accruedInterest)
	if err != nil {
		return nil, fmt.Errorf("failed to record accrual: %v", err)
	}

	return &models.InterestAccrual{
		AccountNumber:   accountNumber,
		AccrualDate:     accrualDate,
		ProductType:     productType,
		Balance:         balance,
		Rate:            rate,
		DailyInterest:   dailyInterest,
		AccruedInterest: accruedInterest,
	}, nil
}

// GetApplicableRate finds the interest rate for a product type, balance, and tenor
func GetApplicableRate(productType string, balance float64, tenorMonths int) (float64, error) {
	var query string
	var args []interface{}

	if productType == "deposit" {
		query = `SELECT base_rate FROM interest_rates
		         WHERE product_type = $1 AND is_active = TRUE AND tenor_months = $2
		         ORDER BY effective_date DESC LIMIT 1`
		args = []interface{}{productType, tenorMonths}
	} else {
		query = `SELECT base_rate FROM interest_rates
		         WHERE product_type = $1 AND is_active = TRUE
		         AND min_balance <= $2 AND (max_balance IS NULL OR max_balance >= $3)
		         ORDER BY effective_date DESC, min_balance DESC LIMIT 1`
		args = []interface{}{productType, balance, balance}
	}

	var rate float64
	err := database.DB.QueryRow(query, args...).Scan(&rate)
	if err != nil {
		return 0, fmt.Errorf("no applicable rate found for %s", productType)
	}

	return rate, nil
}

// AccrueInterestForAllAccounts runs daily interest accrual for all active savings accounts
func AccrueInterestForAllAccounts(date string) (int, int, error) {
	rows, err := database.DB.Query(`SELECT account_number FROM accounts WHERE status = 'active' AND balance > 0
	                                 AND account_type IN ('savings', 'checking', 'giro', 'tabungan')`)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var processed, failed int
	for rows.Next() {
		var accountNumber string
		rows.Scan(&accountNumber)

		_, err := CalculateDailyInterest(accountNumber, date)
		if err != nil {
			failed++
		} else {
			processed++
		}
	}

	return processed, failed, nil
}

// PostMonthlyInterest posts accumulated interest to account balances
func PostMonthlyInterest(yearMonth string) (int, error) {
	// Gunakan date range untuk menggantikan LIKE pada PostgreSQL DATE column
	startDate := yearMonth + "-01"
	endDate := yearMonth + "-31"

	query := `SELECT account_number, SUM(daily_interest) as total_interest
	          FROM interest_accruals
	          WHERE accrual_date >= $1 AND accrual_date <= $2 AND is_posted = FALSE
	          GROUP BY account_number`

	rows, err := database.DB.Query(query, startDate, endDate)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	today := time.Now().UTC().Format("2006-01-02")
	var posted int

	for rows.Next() {
		var accountNumber string
		var totalInterest float64
		rows.Scan(&accountNumber, &totalInterest)

		if totalInterest <= 0 {
			continue
		}

		totalInterest = math.Round(totalInterest*100) / 100

		tx, err := database.DB.Begin()
		if err != nil {
			continue
		}

		_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1, avail_balance = avail_balance + $2, updated_at = NOW()
		                  WHERE account_number = $3`, totalInterest, totalInterest, accountNumber)
		if err != nil {
			tx.Rollback()
			continue
		}

		_, err = tx.Exec(`UPDATE interest_accruals SET is_posted = TRUE
		                  WHERE account_number = $1 AND accrual_date >= $2 AND accrual_date <= $3 AND is_posted = FALSE`,
			accountNumber, startDate, endDate)
		if err != nil {
			tx.Rollback()
			continue
		}

		tx.Commit()

		CreateJournalEntry(today, fmt.Sprintf("Monthly interest posting for %s", accountNumber),
			"interest", accountNumber, "SYSTEM", []JournalLineInput{
				{AccountCode: "511", DebitAmount: totalInterest, Description: "Beban bunga tabungan"},
				{AccountCode: "211", CreditAmount: totalInterest, Description: fmt.Sprintf("Bunga %s untuk %s", yearMonth, accountNumber)},
			})

		posted++
	}

	return posted, nil
}

// CalculateDepositInterest calculates interest for a time deposit
func CalculateDepositInterest(depositNumber string) (float64, error) {
	var principal, interestRate float64
	var tenorMonths int
	err := database.DB.QueryRow(`SELECT principal_amount, interest_rate, tenor_months FROM deposits WHERE deposit_number = $1`,
		depositNumber).Scan(&principal, &interestRate, &tenorMonths)
	if err != nil {
		return 0, fmt.Errorf("deposit not found")
	}

	interest := principal * (interestRate / 100) * float64(tenorMonths) / 12
	return math.Round(interest*100) / 100, nil
}

// CalculateLoanInterest calculates loan interest using annuity method
func CalculateLoanInterest(loanNumber string) (float64, error) {
	var outstanding, interestRate float64
	err := database.DB.QueryRow(`SELECT outstanding_amount, interest_rate FROM loans WHERE loan_number = $1`,
		loanNumber).Scan(&outstanding, &interestRate)
	if err != nil {
		return 0, fmt.Errorf("loan not found")
	}

	monthlyInterest := outstanding * interestRate / 12 / 100
	return math.Round(monthlyInterest*100) / 100, nil
}

// GetInterestRates returns active interest rates with optional product type filter
func GetInterestRates(productType string) ([]models.InterestRate, error) {
	args := []interface{}{}
	query := `SELECT id, product_type, product_name, rate_type, base_rate, min_balance, max_balance,
	           tenor_months, effective_date, is_active FROM interest_rates WHERE is_active = TRUE`

	if productType != "" {
		query += " AND product_type = $1"
		args = append(args, productType)
	}

	query += " ORDER BY product_type, min_balance, tenor_months"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []models.InterestRate
	for rows.Next() {
		var r models.InterestRate
		rows.Scan(&r.ID, &r.ProductType, &r.ProductName, &r.RateType, &r.BaseRate,
			&r.MinBalance, &r.MaxBalance, &r.TenorMonths, &r.EffectiveDate, &r.IsActive)
		rates = append(rates, r)
	}

	return rates, nil
}

// InterestSimulation represents a simulation result
type InterestSimulation struct {
	ProductType     string  `json:"product_type"`
	Principal       float64 `json:"principal"`
	Rate            float64 `json:"rate"`
	TenorMonths     int     `json:"tenor_months"`
	TotalInterest   float64 `json:"total_interest"`
	MaturityAmount  float64 `json:"maturity_amount"`
	MonthlyInterest float64 `json:"monthly_interest,omitempty"`
}

// SimulateInterest calculates interest for a given scenario
func SimulateInterest(productType string, principal float64, tenorMonths int) (*InterestSimulation, error) {
	rate, err := GetApplicableRate(productType, principal, tenorMonths)
	if err != nil {
		return nil, err
	}

	totalInterest := principal * (rate / 100) * float64(tenorMonths) / 12
	totalInterest = math.Round(totalInterest*100) / 100

	sim := &InterestSimulation{
		ProductType:    productType,
		Principal:      principal,
		Rate:           rate,
		TenorMonths:    tenorMonths,
		TotalInterest:  totalInterest,
		MaturityAmount: principal + totalInterest,
	}

	if tenorMonths > 0 {
		sim.MonthlyInterest = math.Round(totalInterest/float64(tenorMonths)*100) / 100
	}

	return sim, nil
}
