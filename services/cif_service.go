package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"database/sql"
	"fmt"
)

// GetSingleCustomerView returns a complete view of a customer
func GetSingleCustomerView(cif string) (*models.SingleCustomerView, error) {
	// Get base customer
	customer, err := GetCustomerByCIF(cif)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %v", err)
	}

	view := &models.SingleCustomerView{
		Customer: *customer,
	}

	// Get extended data
	ext, _ := GetCustomerExtended(cif)
	view.Extended = ext

	// Get accounts
	rows, err := database.DB.Query(`SELECT id, account_number, cif, account_type, currency, balance, avail_balance, status, opened_date, branch 
	                                 FROM accounts WHERE cif = ?`, cif)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a models.Account
			rows.Scan(&a.ID, &a.AccountNumber, &a.CIF, &a.AccountType, &a.Currency,
				&a.Balance, &a.AvailBalance, &a.Status, &a.OpenedDate, &a.Branch)
			view.Accounts = append(view.Accounts, a)
		}
	}

	// Get loans
	loanRows, err := database.DB.Query(`SELECT id, loan_number, cif, account_number, loan_type, principal_amount, outstanding_amount, 
	                                     interest_rate, monthly_payment, tenor_months, remaining_months, disbursement_date, 
	                                     maturity_date, next_payment_date, status FROM loans WHERE cif = ?`, cif)
	if err == nil {
		defer loanRows.Close()
		for loanRows.Next() {
			var l models.Loan
			loanRows.Scan(&l.ID, &l.LoanNumber, &l.CIF, &l.AccountNumber, &l.LoanType,
				&l.PrincipalAmount, &l.OutstandingAmount, &l.InterestRate, &l.MonthlyPayment,
				&l.TenorMonths, &l.RemainingMonths, &l.DisbursementDate, &l.MaturityDate,
				&l.NextPaymentDate, &l.Status)
			view.Loans = append(view.Loans, l)
		}
	}

	// Get deposits
	depRows, err := database.DB.Query(`SELECT id, deposit_number, cif, principal_amount, interest_rate, tenor_months,
	                                    open_date, maturity_date, maturity_amount, auto_renew, status, linked_account
	                                    FROM deposits WHERE cif = ?`, cif)
	if err == nil {
		defer depRows.Close()
		for depRows.Next() {
			var d models.Deposit
			depRows.Scan(&d.ID, &d.DepositNumber, &d.CIF, &d.PrincipalAmount, &d.InterestRate,
				&d.TenorMonths, &d.OpenDate, &d.MaturityDate, &d.MaturityAmount,
				&d.AutoRenew, &d.Status, &d.LinkedAccount)
			view.Deposits = append(view.Deposits, d)
		}
	}

	// Get cards
	cardRows, err := database.DB.Query(`SELECT id, card_number, cif, account_number, card_type, card_brand, 
	                                     card_limit, avail_limit, expiry_date, status FROM cards WHERE cif = ?`, cif)
	if err == nil {
		defer cardRows.Close()
		for cardRows.Next() {
			var c models.Card
			cardRows.Scan(&c.ID, &c.CardNumber, &c.CIF, &c.AccountNumber, &c.CardType,
				&c.CardBrand, &c.CardLimit, &c.AvailLimit, &c.ExpiryDate, &c.Status)
			view.Cards = append(view.Cards, c)
		}
	}

	// Get roles
	roles, _ := GetUserRoles(cif)
	view.Roles = roles

	return view, nil
}

// GetCustomerExtended returns extended CIF data
func GetCustomerExtended(cif string) (*models.CustomerExtended, error) {
	var ext models.CustomerExtended
	query := `SELECT id, cif, mother_maiden_name, nationality, occupation, employer_name, monthly_income,
	           source_of_funds, risk_profile, segment, branch_code, rm_code, npwp, last_kyc_date, next_kyc_date,
	           created_at, updated_at FROM customer_extended WHERE cif = ?`

	err := database.DB.QueryRow(query, cif).Scan(
		&ext.ID, &ext.CIF, &ext.MotherMaidenName, &ext.Nationality, &ext.Occupation,
		&ext.EmployerName, &ext.MonthlyIncome, &ext.SourceOfFunds, &ext.RiskProfile,
		&ext.Segment, &ext.BranchCode, &ext.RMCode, &ext.NPWP, &ext.LastKYCDate,
		&ext.NextKYCDate, &ext.CreatedAt, &ext.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ext, nil
}

// UpdateCustomerExtended updates or inserts extended CIF data
func UpdateCustomerExtended(cif string, ext models.CustomerExtended) error {
	// Check if exists
	existing, _ := GetCustomerExtended(cif)

	if existing != nil {
		// Update
		query := `UPDATE customer_extended SET 
		          mother_maiden_name=?, nationality=?, occupation=?, employer_name=?, monthly_income=?,
		          source_of_funds=?, risk_profile=?, segment=?, branch_code=?, rm_code=?, npwp=?,
		          updated_at=CURRENT_TIMESTAMP WHERE cif=?`
		_, err := database.DB.Exec(query, ext.MotherMaidenName, ext.Nationality, ext.Occupation,
			ext.EmployerName, ext.MonthlyIncome, ext.SourceOfFunds, ext.RiskProfile,
			ext.Segment, ext.BranchCode, ext.RMCode, ext.NPWP, cif)
		return err
	}

	// Insert
	query := `INSERT INTO customer_extended (cif, mother_maiden_name, nationality, occupation, employer_name, monthly_income,
	           source_of_funds, risk_profile, segment, branch_code, rm_code, npwp)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, cif, ext.MotherMaidenName, ext.Nationality, ext.Occupation,
		ext.EmployerName, ext.MonthlyIncome, ext.SourceOfFunds, ext.RiskProfile,
		ext.Segment, ext.BranchCode, ext.RMCode, ext.NPWP)
	return err
}

// SearchCustomers searches by name, CIF, KTP, or phone
func SearchCustomers(query string) ([]models.Customer, error) {
	searchQuery := `SELECT id, cif, full_name, id_card_number, phone_number, email, address, date_of_birth, status 
	                FROM customers 
	                WHERE cif LIKE ? OR full_name LIKE ? OR id_card_number LIKE ? OR phone_number LIKE ?
	                ORDER BY full_name LIMIT 20`

	pattern := "%" + query + "%"
	rows, err := database.DB.Query(searchQuery, pattern, pattern, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		rows.Scan(&c.ID, &c.CIF, &c.FullName, &c.IDCardNumber, &c.PhoneNumber,
			&c.Email, &c.Address, &c.DateOfBirth, &c.Status)
		customers = append(customers, c)
	}

	return customers, nil
}
