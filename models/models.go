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

// Bank represents supported banks for transfers
type Bank struct {
	ID        int       `json:"id" db:"id"`
	BankCode  string    `json:"bank_code" db:"bank_code"`
	BankName  string    `json:"bank_name" db:"bank_name"`
	SwiftCode string    `json:"swift_code" db:"swift_code"`
	IsActive  int       `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TransferFee represents dynamic fees for interbank transfers (OUTBOUND to other banks)
// This system represents ONE specific bank. Customers transfer money OUT to destination banks.
// Intrabank (within our bank) transfers are FREE.
// Incoming transfers from other banks are FREE (received by our customers).
type TransferFee struct {
	ID                  int        `json:"id" db:"id"`
	DestinationBankCode string     `json:"destination_bank_code" db:"destination_bank_code"` // Target bank code
	DestinationBankName string     `json:"destination_bank_name" db:"destination_bank_name"` // Target bank name
	FeeType             string     `json:"fee_type" db:"fee_type"`                           // flat or percentage
	FeeAmount           float64    `json:"fee_amount" db:"fee_amount"`
	FeePercentage       float64    `json:"fee_percentage" db:"fee_percentage"`
	MinimumAmount       float64    `json:"minimum_amount" db:"minimum_amount"`
	MaximumAmount       float64    `json:"maximum_amount" db:"maximum_amount"`
	Description         string     `json:"description" db:"description"`
	IsActive            int        `json:"is_active" db:"is_active"`
	EffectiveFrom       time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveTo         *time.Time `json:"effective_to" db:"effective_to"`
	Notes               string     `json:"notes" db:"notes"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// ServiceFee represents fees for other services (e-wallet, e-money, VA, QRIS, etc)
type ServiceFee struct {
	ID            int        `json:"id" db:"id"`
	ServiceCode   string     `json:"service_code" db:"service_code"` // TOPUP_OVO, TOPUP_DANA, etc
	ServiceName   string     `json:"service_name" db:"service_name"`
	ServiceType   string     `json:"service_type" db:"service_type"` // topup_ewallet, topup_emoney, payment_va, qris_payment
	FeeType       string     `json:"fee_type" db:"fee_type"`         // flat, percentage
	FeeAmount     float64    `json:"fee_amount" db:"fee_amount"`
	FeePercentage float64    `json:"fee_percentage" db:"fee_percentage"`
	MinimumAmount float64    `json:"minimum_amount" db:"minimum_amount"`
	MaximumAmount float64    `json:"maximum_amount" db:"maximum_amount"`
	IsActive      int        `json:"is_active" db:"is_active"`
	EffectiveFrom time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to" db:"effective_to"`
	Notes         string     `json:"notes" db:"notes"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ==================== SECURITY MODELS ====================

// TokenBlacklist represents revoked JWT tokens
type TokenBlacklist struct {
	ID            int       `json:"id" db:"id"`
	TokenJTI      string    `json:"token_jti" db:"token_jti"`
	CIF           string    `json:"cif" db:"cif"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	BlacklistedAt time.Time `json:"blacklisted_at" db:"blacklisted_at"`
}

// Role represents system roles (customer, teller, admin, supervisor)
type Role struct {
	ID          int       `json:"id" db:"id"`
	RoleName    string    `json:"role_name" db:"role_name"`
	Description string    `json:"description" db:"description"`
	IsActive    int       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserRole represents the assignment of roles to users
type UserRole struct {
	ID         int       `json:"id" db:"id"`
	CIF        string    `json:"cif" db:"cif"`
	RoleID     int       `json:"role_id" db:"role_id"`
	AssignedBy string    `json:"assigned_by" db:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}

// AuditLog represents immutable activity logs
type AuditLog struct {
	ID             int       `json:"id" db:"id"`
	CIF            string    `json:"cif" db:"cif"`
	Action         string    `json:"action" db:"action"`
	Resource       string    `json:"resource" db:"resource"`
	ResourceID     string    `json:"resource_id" db:"resource_id"`
	IPAddress      string    `json:"ip_address" db:"ip_address"`
	UserAgent      string    `json:"user_agent" db:"user_agent"`
	RequestMethod  string    `json:"request_method" db:"request_method"`
	RequestPath    string    `json:"request_path" db:"request_path"`
	RequestBody    string    `json:"request_body,omitempty" db:"request_body"`
	ResponseStatus int       `json:"response_status" db:"response_status"`
	Details        string    `json:"details" db:"details"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// TransactionLimit represents per-role transaction limits
type TransactionLimit struct {
	ID                  int       `json:"id" db:"id"`
	RoleName            string    `json:"role_name" db:"role_name"`
	TransactionType     string    `json:"transaction_type" db:"transaction_type"`
	DailyLimit          float64   `json:"daily_limit" db:"daily_limit"`
	PerTransactionLimit float64   `json:"per_transaction_limit" db:"per_transaction_limit"`
	MonthlyLimit        float64   `json:"monthly_limit" db:"monthly_limit"`
	IsActive            int       `json:"is_active" db:"is_active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// LoginAttempt represents login attempt tracking for lockout
type LoginAttempt struct {
	ID          int       `json:"id" db:"id"`
	CIF         string    `json:"cif" db:"cif"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	AttemptType string    `json:"attempt_type" db:"attempt_type"`
	IsSuccess   int       `json:"is_success" db:"is_success"`
	AttemptedAt time.Time `json:"attempted_at" db:"attempted_at"`
}

// OTPCode represents OTP codes for MFA and account unlock
type OTPCode struct {
	ID        int       `json:"id" db:"id"`
	CIF       string    `json:"cif" db:"cif"`
	OTPCode   string    `json:"otp_code,omitempty" db:"otp_code"`
	OTPType   string    `json:"otp_type" db:"otp_type"` // login_mfa, unlock_account, reset_pin
	Channel   string    `json:"channel" db:"channel"`   // sms, email
	IsUsed    int       `json:"is_used" db:"is_used"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EKYCVerification represents e-KYC verification for self-service unlock
type EKYCVerification struct {
	ID                 int        `json:"id" db:"id"`
	VerificationID     string     `json:"verification_id" db:"verification_id"`
	CIF                string     `json:"cif" db:"cif"`
	IDCardNumber       string     `json:"id_card_number" db:"id_card_number"`
	VerificationType   string     `json:"verification_type" db:"verification_type"`     // unlock_account, reset_pin
	VerificationStatus string     `json:"verification_status" db:"verification_status"` // pending, verified, rejected
	VerifiedAt         *time.Time `json:"verified_at" db:"verified_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

// =============================================
// Phase 2: Core Banking Models
// =============================================

// ChartOfAccount represents a general ledger account in the CoA
type ChartOfAccount struct {
	ID            int       `json:"id" db:"id"`
	AccountCode   string    `json:"account_code" db:"account_code"`
	AccountName   string    `json:"account_name" db:"account_name"`
	AccountType   string    `json:"account_type" db:"account_type"`
	ParentCode    *string   `json:"parent_code" db:"parent_code"`
	Level         int       `json:"level" db:"level"`
	NormalBalance string    `json:"normal_balance" db:"normal_balance"`
	IsActive      int       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// JournalEntry represents a double-entry journal posting
type JournalEntry struct {
	ID            int       `json:"id" db:"id"`
	JournalNumber string    `json:"journal_number" db:"journal_number"`
	EntryDate     string    `json:"entry_date" db:"entry_date"`
	Description   string    `json:"description" db:"description"`
	ReferenceType string    `json:"reference_type" db:"reference_type"`
	ReferenceID   string    `json:"reference_id" db:"reference_id"`
	PostedBy      string    `json:"posted_by" db:"posted_by"`
	Status        string    `json:"status" db:"status"`
	ReversedBy    *int      `json:"reversed_by" db:"reversed_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// JournalLine represents a single debit or credit line in a journal entry
type JournalLine struct {
	ID           int     `json:"id" db:"id"`
	JournalID    int     `json:"journal_id" db:"journal_id"`
	AccountCode  string  `json:"account_code" db:"account_code"`
	AccountName  string  `json:"account_name,omitempty"`
	DebitAmount  float64 `json:"debit_amount" db:"debit_amount"`
	CreditAmount float64 `json:"credit_amount" db:"credit_amount"`
	Description  string  `json:"description" db:"description"`
}

// CustomerExtended represents additional CIF data for banking
type CustomerExtended struct {
	ID               int       `json:"id" db:"id"`
	CIF              string    `json:"cif" db:"cif"`
	MotherMaidenName string    `json:"mother_maiden_name" db:"mother_maiden_name"`
	Nationality      string    `json:"nationality" db:"nationality"`
	Occupation       string    `json:"occupation" db:"occupation"`
	EmployerName     string    `json:"employer_name" db:"employer_name"`
	MonthlyIncome    float64   `json:"monthly_income" db:"monthly_income"`
	SourceOfFunds    string    `json:"source_of_funds" db:"source_of_funds"`
	RiskProfile      string    `json:"risk_profile" db:"risk_profile"`
	Segment          string    `json:"segment" db:"segment"`
	BranchCode       string    `json:"branch_code" db:"branch_code"`
	RMCode           string    `json:"rm_code" db:"rm_code"`
	NPWP             string    `json:"npwp" db:"npwp"`
	LastKYCDate      *string   `json:"last_kyc_date" db:"last_kyc_date"`
	NextKYCDate      *string   `json:"next_kyc_date" db:"next_kyc_date"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// InterestRate represents interest rate configuration
type InterestRate struct {
	ID            int      `json:"id" db:"id"`
	ProductType   string   `json:"product_type" db:"product_type"`
	ProductName   string   `json:"product_name" db:"product_name"`
	RateType      string   `json:"rate_type" db:"rate_type"`
	BaseRate      float64  `json:"base_rate" db:"base_rate"`
	MinBalance    float64  `json:"min_balance" db:"min_balance"`
	MaxBalance    *float64 `json:"max_balance" db:"max_balance"`
	TenorMonths   *int     `json:"tenor_months" db:"tenor_months"`
	EffectiveDate string   `json:"effective_date" db:"effective_date"`
	IsActive      int      `json:"is_active" db:"is_active"`
}

// InterestAccrual represents daily interest calculation record
type InterestAccrual struct {
	ID              int     `json:"id" db:"id"`
	AccountNumber   string  `json:"account_number" db:"account_number"`
	AccrualDate     string  `json:"accrual_date" db:"accrual_date"`
	ProductType     string  `json:"product_type" db:"product_type"`
	Balance         float64 `json:"balance" db:"balance"`
	Rate            float64 `json:"rate" db:"rate"`
	DailyInterest   float64 `json:"daily_interest" db:"daily_interest"`
	AccruedInterest float64 `json:"accrued_interest" db:"accrued_interest"`
	IsPosted        int     `json:"is_posted" db:"is_posted"`
}

// StandingInstruction represents a recurring payment instruction
type StandingInstruction struct {
	ID                int       `json:"id" db:"id"`
	SINumber          string    `json:"si_number" db:"si_number"`
	CIF               string    `json:"cif" db:"cif"`
	FromAccount       string    `json:"from_account" db:"from_account"`
	InstructionType   string    `json:"instruction_type" db:"instruction_type"`
	ToAccount         string    `json:"to_account" db:"to_account"`
	ToBankCode        string    `json:"to_bank_code" db:"to_bank_code"`
	Amount            float64   `json:"amount" db:"amount"`
	Description       string    `json:"description" db:"description"`
	Frequency         string    `json:"frequency" db:"frequency"`
	ExecutionDay      int       `json:"execution_day" db:"execution_day"`
	StartDate         string    `json:"start_date" db:"start_date"`
	EndDate           *string   `json:"end_date" db:"end_date"`
	NextExecutionDate string    `json:"next_execution_date" db:"next_execution_date"`
	TotalExecuted     int       `json:"total_executed" db:"total_executed"`
	TotalFailed       int       `json:"total_failed" db:"total_failed"`
	LastExecutionDate *string   `json:"last_execution_date" db:"last_execution_date"`
	LastStatus        string    `json:"last_status" db:"last_status"`
	Status            string    `json:"status" db:"status"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// SIExecution represents an execution log for a standing instruction
type SIExecution struct {
	ID            int       `json:"id" db:"id"`
	SINumber      string    `json:"si_number" db:"si_number"`
	ExecutionDate string    `json:"execution_date" db:"execution_date"`
	Amount        float64   `json:"amount" db:"amount"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"`
	Status        string    `json:"status" db:"status"`
	ErrorMessage  string    `json:"error_message" db:"error_message"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// EODLog represents end-of-day processing log
type EODLog struct {
	ID               int        `json:"id" db:"id"`
	ProcessDate      string     `json:"process_date" db:"process_date"`
	ProcessType      string     `json:"process_type" db:"process_type"`
	Status           string     `json:"status" db:"status"`
	RecordsProcessed int        `json:"records_processed" db:"records_processed"`
	RecordsFailed    int        `json:"records_failed" db:"records_failed"`
	StartedAt        *time.Time `json:"started_at" db:"started_at"`
	CompletedAt      *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage     string     `json:"error_message" db:"error_message"`
}

// SingleCustomerView represents the full customer overview
type SingleCustomerView struct {
	Customer Customer          `json:"customer"`
	Extended *CustomerExtended `json:"extended"`
	Accounts []Account         `json:"accounts"`
	Loans    []Loan            `json:"loans"`
	Deposits []Deposit         `json:"deposits"`
	Cards    []Card            `json:"cards"`
	Roles    []string          `json:"roles"`
}

// TrialBalance represents a trial balance report line
type TrialBalance struct {
	AccountCode   string  `json:"account_code"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"account_type"`
	DebitBalance  float64 `json:"debit_balance"`
	CreditBalance float64 `json:"credit_balance"`
}
