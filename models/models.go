package models

import (
	"time"
)

// Customer represents a bank customer
type Customer struct {
	ID           int       `json:"id" db:"id"`
	CIF          string    `json:"cif" db:"cif"` // Customer Information File number
	FullName     string    `json:"full_name" db:"full_name"`
	IDCardNumber string    `json:"id_card_number" db:"id_card_number"` // KTP
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	Email        string    `json:"email" db:"email"`
	Address      string    `json:"address" db:"address"`
	DateOfBirth  string    `json:"date_of_birth" db:"date_of_birth"`
	PIN          string    `json:"-" db:"pin"`         // Hashed PIN for mobile banking
	Status       string    `json:"status" db:"status"` // active, blocked, closed
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Account represents customer's bank account
type Account struct {
	ID            int       `json:"id" db:"id"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	CIF           string    `json:"cif" db:"cif"`
	AccountType   string    `json:"account_type" db:"account_type"` // savings, checking, loan, deposit
	Currency      string    `json:"currency" db:"currency"`         // IDR, USD, etc
	Balance       float64   `json:"balance" db:"balance"`
	AvailBalance  float64   `json:"avail_balance" db:"avail_balance"` // Available balance (after holds)
	Status        string    `json:"status" db:"status"`               // active, blocked, closed
	OpenedDate    string    `json:"opened_date" db:"opened_date"`
	Branch        string    `json:"branch" db:"branch"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents all financial transactions
type Transaction struct {
	ID                int       `json:"id" db:"id"`
	TransactionID     string    `json:"transaction_id" db:"transaction_id"`     // Unique transaction ID
	TransactionType   string    `json:"transaction_type" db:"transaction_type"` // transfer, payment, withdrawal, etc
	FromAccountNumber string    `json:"from_account_number" db:"from_account_number"`
	ToAccountNumber   string    `json:"to_account_number" db:"to_account_number"`
	Amount            float64   `json:"amount" db:"amount"`
	Currency          string    `json:"currency" db:"currency"`
	Description       string    `json:"description" db:"description"`
	ReferenceNumber   string    `json:"reference_number" db:"reference_number"`
	Status            string    `json:"status" db:"status"` // pending, success, failed, reversed
	TransactionDate   time.Time `json:"transaction_date" db:"transaction_date"`
	SettlementDate    string    `json:"settlement_date" db:"settlement_date"`
	Fee               float64   `json:"fee" db:"fee"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// Card represents debit/credit cards
type Card struct {
	ID            int       `json:"id" db:"id"`
	CardNumber    string    `json:"card_number" db:"card_number"`
	CIF           string    `json:"cif" db:"cif"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	CardType      string    `json:"card_type" db:"card_type"`   // debit, credit
	CardBrand     string    `json:"card_brand" db:"card_brand"` // visa, mastercard
	CardLimit     float64   `json:"card_limit" db:"card_limit"`
	AvailLimit    float64   `json:"avail_limit" db:"avail_limit"`
	ExpiryDate    string    `json:"expiry_date" db:"expiry_date"`
	Status        string    `json:"status" db:"status"` // active, blocked, expired
	CVV           string    `json:"-" db:"cvv"`
	PIN           string    `json:"-" db:"pin"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Loan represents customer loans
type Loan struct {
	ID                int       `json:"id" db:"id"`
	LoanNumber        string    `json:"loan_number" db:"loan_number"`
	CIF               string    `json:"cif" db:"cif"`
	AccountNumber     string    `json:"account_number" db:"account_number"`
	LoanType          string    `json:"loan_type" db:"loan_type"` // personal, mortgage, business, etc
	PrincipalAmount   float64   `json:"principal_amount" db:"principal_amount"`
	OutstandingAmount float64   `json:"outstanding_amount" db:"outstanding_amount"`
	InterestRate      float64   `json:"interest_rate" db:"interest_rate"`
	MonthlyPayment    float64   `json:"monthly_payment" db:"monthly_payment"`
	TenorMonths       int       `json:"tenor_months" db:"tenor_months"`
	RemainingMonths   int       `json:"remaining_months" db:"remaining_months"`
	DisbursementDate  string    `json:"disbursement_date" db:"disbursement_date"`
	MaturityDate      string    `json:"maturity_date" db:"maturity_date"`
	NextPaymentDate   string    `json:"next_payment_date" db:"next_payment_date"`
	Status            string    `json:"status" db:"status"` // active, paid_off, overdue
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Deposit represents time deposits
type Deposit struct {
	ID              int       `json:"id" db:"id"`
	DepositNumber   string    `json:"deposit_number" db:"deposit_number"`
	CIF             string    `json:"cif" db:"cif"`
	PrincipalAmount float64   `json:"principal_amount" db:"principal_amount"`
	InterestRate    float64   `json:"interest_rate" db:"interest_rate"`
	TenorMonths     int       `json:"tenor_months" db:"tenor_months"`
	OpenDate        string    `json:"open_date" db:"open_date"`
	MaturityDate    string    `json:"maturity_date" db:"maturity_date"`
	MaturityAmount  float64   `json:"maturity_amount" db:"maturity_amount"`
	AutoRenew       bool      `json:"auto_renew" db:"auto_renew"`
	Status          string    `json:"status" db:"status"`                 // active, matured, closed
	LinkedAccount   string    `json:"linked_account" db:"linked_account"` // For interest payment
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// BillPayment represents PPOB bills
type BillPayment struct {
	ID             int       `json:"id" db:"id"`
	BillerCode     string    `json:"biller_code" db:"biller_code"` // PLN, PDAM, Telkom, etc
	BillerName     string    `json:"biller_name" db:"biller_name"`
	CustomerNumber string    `json:"customer_number" db:"customer_number"` // Customer ID at biller
	BillNumber     string    `json:"bill_number" db:"bill_number"`
	BillAmount     float64   `json:"bill_amount" db:"bill_amount"`
	AdminFee       float64   `json:"admin_fee" db:"admin_fee"`
	TotalAmount    float64   `json:"total_amount" db:"total_amount"`
	BillPeriod     string    `json:"bill_period" db:"bill_period"`
	DueDate        string    `json:"due_date" db:"due_date"`
	Status         string    `json:"status" db:"status"` // unpaid, paid
	TransactionID  string    `json:"transaction_id" db:"transaction_id"`
	PaymentDate    string    `json:"payment_date" db:"payment_date"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Notification represents transaction notifications
type Notification struct {
	ID               int       `json:"id" db:"id"`
	CIF              string    `json:"cif" db:"cif"`
	NotificationType string    `json:"notification_type" db:"notification_type"` // transfer, payment, deposit, loan, etc
	Title            string    `json:"title" db:"title"`
	Message          string    `json:"message" db:"message"`
	TransactionID    string    `json:"transaction_id" db:"transaction_id"`
	IsRead           int       `json:"is_read" db:"is_read"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// FCMToken represents device FCM registration token for push notifications
type FCMToken struct {
	ID          int       `json:"id" db:"id"`
	CIF         string    `json:"cif" db:"cif"`
	DeviceToken string    `json:"device_token" db:"device_token"`
	DeviceType  string    `json:"device_type" db:"device_type"` // android, ios, web
	DeviceName  string    `json:"device_name" db:"device_name"`
	IsActive    int       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NotificationPreference represents user notification settings
type NotificationPreference struct {
	ID                    int       `json:"id" db:"id"`
	CIF                   string    `json:"cif" db:"cif"`
	TransferNotification  int       `json:"transfer_notification" db:"transfer_notification"`
	PaymentNotification   int       `json:"payment_notification" db:"payment_notification"`
	DepositNotification   int       `json:"deposit_notification" db:"deposit_notification"`
	LoanNotification      int       `json:"loan_notification" db:"loan_notification"`
	PromotionNotification int       `json:"promotion_notification" db:"promotion_notification"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}
