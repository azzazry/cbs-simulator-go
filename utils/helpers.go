package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateTransactionID generates a unique transaction ID with Indonesian banking format
func GenerateTransactionID() string {
	now := time.Now()
	return fmt.Sprintf("TRX%s%s", now.Format("20060102"), uuid.New().String()[:8])
}

// GenerateAccountNumber generates a bank account number
func GenerateAccountNumber(accountType string) string {
	prefix := "1" // savings
	switch accountType {
	case "checking":
		prefix = "2"
	case "loan":
		prefix = "3"
	case "deposit":
		prefix = "4"
	}
	
	now := time.Now()
	return fmt.Sprintf("%s%s%04d", prefix, now.Format("060102"), time.Now().Nanosecond()%10000)
}

// GenerateCIF generates a Customer Information File number
func GenerateCIF() string {
	now := time.Now()
	return fmt.Sprintf("CIF%s%04d", now.Format("060102"), time.Now().Nanosecond()%10000)
}

// GenerateLoanNumber generates a loan number
func GenerateLoanNumber() string {
	now := time.Now()
	return fmt.Sprintf("LOAN%s%04d", now.Format("060102"), time.Now().Nanosecond()%10000)
}

// GenerateDepositNumber generates a deposit number
func GenerateDepositNumber() string {
	now := time.Now()
	return fmt.Sprintf("DEP%s%04d", now.Format("060102"), time.Now().Nanosecond()%10000)
}

// GenerateReferenceNumber generates a reference number for transactions
func GenerateReferenceNumber() string {
	return fmt.Sprintf("REF%s", uuid.New().String()[:12])
}

// HashPIN hashes a PIN using bcrypt
func HashPIN(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPIN verifies a PIN against a hash
func VerifyPIN(pin, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin))
	return err == nil
}

// FormatCurrency formats amount to Indonesian Rupiah format
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("Rp %.2f", amount)
}

// GetCurrentDate returns current date in YYYY-MM-DD format
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurrentDateTime returns current datetime
func GetCurrentDateTime() time.Time {
	return time.Now()
}

// AddMonths adds months to a date string
func AddMonths(dateStr string, months int) (string, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", err
	}
	
	newDate := date.AddDate(0, months, 0)
	return newDate.Format(layout), nil
}

// CalculateMaturityAmount calculates deposit maturity amount
func CalculateMaturityAmount(principal float64, interestRate float64, tenorMonths int) float64 {
	monthlyRate := interestRate / 100 / 12
	return principal * (1 + (monthlyRate * float64(tenorMonths)))
}

// ValidateAccountNumber validates account number format
func ValidateAccountNumber(accountNumber string) bool {
	return len(accountNumber) >= 10 && len(accountNumber) <= 20
}

// SimulateDelay simulates processing delay for realistic behavior
func SimulateDelay(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
