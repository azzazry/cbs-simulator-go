package utils

import (
	"io"
	"log"
	"os"
	"time"
)

// Logger represents server logger with separate info and error logs
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

var AppLogger *Logger

// InitLogger initializes the logger with file and console output
func InitLogger() {
	// Create logs directory if not exists
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Println("Warning: Could not create logs directory")
	}

	// Open log files
	infoFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Warning: Could not open app.log")
	}

	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Warning: Could not open error.log")
	}

	// Multi-writer: ke file dan stdout
	AppLogger = &Logger{
		infoLog:  log.New(io.MultiWriter(infoFile, os.Stdout), "[INFO] ", log.LstdFlags|log.Lshortfile),
		errorLog: log.New(io.MultiWriter(errorFile, os.Stderr), "[ERROR] ", log.LstdFlags|log.Lshortfile),
	}
}

// LogAPI logs API request/response
func (l *Logger) LogAPI(method, path, clientIP string, status int, latency time.Duration) {
	l.infoLog.Printf("API | Method: %s | Path: %s | IP: %s | Status: %d | Latency: %v",
		method, path, clientIP, status, latency)
}

// LogAPIWithBody logs API request/response with request/response body
func (l *Logger) LogAPIWithBody(method, path, clientIP string, status int, latency time.Duration, reqBody, resBody string) {
	if len(reqBody) > 200 {
		reqBody = reqBody[:200] + "..."
	}
	if len(resBody) > 200 {
		resBody = resBody[:200] + "..."
	}
	l.infoLog.Printf("API | Method: %s | Path: %s | IP: %s | Status: %d | Latency: %v | Req: %s | Res: %s",
		method, path, clientIP, status, latency, reqBody, resBody)
}

// LogTransaction logs transaction with CIF, transaction ID, type, amount and status
func (l *Logger) LogTransaction(cif, transactionID, transactionType string, amount float64, status string) {
	l.infoLog.Printf("TRANSACTION | CIF: %s | TxID: %s | Type: %s | Amount: %.0f | Status: %s",
		cif, transactionID, transactionType, amount, status)
}

// LogNotification logs notification sent to customer
func (l *Logger) LogNotification(cif, notificationType, message string) {
	l.infoLog.Printf("NOTIFICATION | CIF: %s | Type: %s | Message: %s",
		cif, notificationType, message)
}

// LogError logs operation error with context
func (l *Logger) LogError(operation, errorMsg string) {
	l.errorLog.Printf("OPERATION: %s | ERROR: %s", operation, errorMsg)
}

// LogErrorWithContext logs error with additional context
func (l *Logger) LogErrorWithContext(operation, cif, errorMsg string) {
	l.errorLog.Printf("OPERATION: %s | CIF: %s | ERROR: %s", operation, cif, errorMsg)
}

// LogInfoWithContext logs info with CIF context
func (l *Logger) LogInfoWithContext(operation, cif, message string) {
	l.infoLog.Printf("%s | CIF: %s | %s", operation, cif, message)
}

// LogAuthEvent logs authentication related events
func (l *Logger) LogAuthEvent(cif, event string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	l.infoLog.Printf("AUTH | CIF: %s | Event: %s | Status: %s", cif, event, status)
}

// LogBillPayment logs bill payment events
func (l *Logger) LogBillPayment(cif, billID string, amount float64, status string) {
	l.infoLog.Printf("BILL_PAYMENT | CIF: %s | BillID: %s | Amount: %.0f | Status: %s",
		cif, billID, amount, status)
}

// LogLoanEvent logs loan-related events
func (l *Logger) LogLoanEvent(cif, loanID, event string) {
	l.infoLog.Printf("LOAN | CIF: %s | LoanID: %s | Event: %s", cif, loanID, event)
}

// LogDepositEvent logs deposit-related events
func (l *Logger) LogDepositEvent(cif, depositID, event string) {
	l.infoLog.Printf("DEPOSIT | CIF: %s | DepositID: %s | Event: %s", cif, depositID, event)
}
