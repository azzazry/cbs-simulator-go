package services

import (
	"database/sql"
	"fmt"

	"cbs-simulator/database"
	"cbs-simulator/models"
)

// ===== CARD SERVICES =====

// GetCardsByCIF retrieves all cards for a customer
func GetCardsByCIF(cif string) ([]models.Card, error) {
	query := `SELECT id, card_number, cif, account_number, card_type, card_brand,
	          card_limit, avail_limit, expiry_date, status, created_at, updated_at
	          FROM cards WHERE cif = $1`

	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.CIF, &card.AccountNumber,
			&card.CardType, &card.CardBrand, &card.CardLimit, &card.AvailLimit,
			&card.ExpiryDate, &card.Status, &card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		card.CardNumber = maskCardNumber(card.CardNumber)
		cards = append(cards, card)
	}

	return cards, nil
}

// GetCardByNumber retrieves card details
func GetCardByNumber(cardNumber string) (*models.Card, error) {
	var card models.Card

	query := `SELECT id, card_number, cif, account_number, card_type, card_brand,
	          card_limit, avail_limit, expiry_date, status, created_at, updated_at
	          FROM cards WHERE card_number = $1`

	err := database.DB.QueryRow(query, cardNumber).Scan(
		&card.ID, &card.CardNumber, &card.CIF, &card.AccountNumber,
		&card.CardType, &card.CardBrand, &card.CardLimit, &card.AvailLimit,
		&card.ExpiryDate, &card.Status, &card.CreatedAt, &card.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, err
	}

	card.CardNumber = maskCardNumber(card.CardNumber)
	return &card, nil
}

// BlockCard blocks a card
func BlockCard(cardNumber string) error {
	query := `UPDATE cards SET status = 'blocked', updated_at = NOW()
	          WHERE card_number = $1`

	result, err := database.DB.Exec(query, cardNumber)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("card not found")
	}

	return nil
}

// UnblockCard unblocks a card
func UnblockCard(cardNumber string) error {
	query := `UPDATE cards SET status = 'active', updated_at = NOW()
	          WHERE card_number = $1`

	result, err := database.DB.Exec(query, cardNumber)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("card not found")
	}

	return nil
}

func maskCardNumber(cardNumber string) string {
	if len(cardNumber) != 16 {
		return cardNumber
	}
	return cardNumber[:4] + "********" + cardNumber[12:]
}

// ===== LOAN SERVICES =====

// GetLoansByCIF retrieves all loans for a customer
func GetLoansByCIF(cif string) ([]models.Loan, error) {
	query := `SELECT id, loan_number, cif, account_number, loan_type, principal_amount,
	          outstanding_amount, interest_rate, monthly_payment, tenor_months,
	          remaining_months, disbursement_date, maturity_date, next_payment_date,
	          status, created_at, updated_at
	          FROM loans WHERE cif = $1`

	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var loan models.Loan
		err := rows.Scan(
			&loan.ID, &loan.LoanNumber, &loan.CIF, &loan.AccountNumber,
			&loan.LoanType, &loan.PrincipalAmount, &loan.OutstandingAmount,
			&loan.InterestRate, &loan.MonthlyPayment, &loan.TenorMonths,
			&loan.RemainingMonths, &loan.DisbursementDate, &loan.MaturityDate,
			&loan.NextPaymentDate, &loan.Status, &loan.CreatedAt, &loan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

// GetLoanByNumber retrieves loan details
func GetLoanByNumber(loanNumber string) (*models.Loan, error) {
	var loan models.Loan

	query := `SELECT id, loan_number, cif, account_number, loan_type, principal_amount,
	          outstanding_amount, interest_rate, monthly_payment, tenor_months,
	          remaining_months, disbursement_date, maturity_date, next_payment_date,
	          status, created_at, updated_at
	          FROM loans WHERE loan_number = $1`

	err := database.DB.QueryRow(query, loanNumber).Scan(
		&loan.ID, &loan.LoanNumber, &loan.CIF, &loan.AccountNumber,
		&loan.LoanType, &loan.PrincipalAmount, &loan.OutstandingAmount,
		&loan.InterestRate, &loan.MonthlyPayment, &loan.TenorMonths,
		&loan.RemainingMonths, &loan.DisbursementDate, &loan.MaturityDate,
		&loan.NextPaymentDate, &loan.Status, &loan.CreatedAt, &loan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan not found")
		}
		return nil, err
	}

	return &loan, nil
}

// ===== DEPOSIT SERVICES =====

// GetDepositsByCIF retrieves all deposits for a customer
func GetDepositsByCIF(cif string) ([]models.Deposit, error) {
	query := `SELECT id, deposit_number, cif, principal_amount, interest_rate,
	          tenor_months, open_date, maturity_date, maturity_amount, auto_renew,
	          status, linked_account, created_at, updated_at
	          FROM deposits WHERE cif = $1`

	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposits []models.Deposit
	for rows.Next() {
		var deposit models.Deposit
		err := rows.Scan(
			&deposit.ID, &deposit.DepositNumber, &deposit.CIF, &deposit.PrincipalAmount,
			&deposit.InterestRate, &deposit.TenorMonths, &deposit.OpenDate,
			&deposit.MaturityDate, &deposit.MaturityAmount, &deposit.AutoRenew,
			&deposit.Status, &deposit.LinkedAccount, &deposit.CreatedAt, &deposit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

// GetDepositByNumber retrieves deposit details
func GetDepositByNumber(depositNumber string) (*models.Deposit, error) {
	var deposit models.Deposit

	query := `SELECT id, deposit_number, cif, principal_amount, interest_rate,
	          tenor_months, open_date, maturity_date, maturity_amount, auto_renew,
	          status, linked_account, created_at, updated_at
	          FROM deposits WHERE deposit_number = $1`

	err := database.DB.QueryRow(query, depositNumber).Scan(
		&deposit.ID, &deposit.DepositNumber, &deposit.CIF, &deposit.PrincipalAmount,
		&deposit.InterestRate, &deposit.TenorMonths, &deposit.OpenDate,
		&deposit.MaturityDate, &deposit.MaturityAmount, &deposit.AutoRenew,
		&deposit.Status, &deposit.LinkedAccount, &deposit.CreatedAt, &deposit.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deposit not found")
		}
		return nil, err
	}

	return &deposit, nil
}

// CreateDeposit creates a new time deposit
func CreateDeposit(cif, linkedAccount string, principalAmount, interestRate float64, tenorMonths int, autoRenew bool) (*models.Deposit, error) {
	return nil, fmt.Errorf("create deposit not yet implemented in simulator")
}
