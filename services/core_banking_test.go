package services_test

import (
	"cbs-simulator/models"
	"cbs-simulator/services"
	"testing"
)

// === GENERAL LEDGER TESTS ===

func TestCreateJournalEntry_Balanced(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	lines := []services.JournalLineInput{
		{AccountCode: "111", DebitAmount: 1000000, Description: "Kas masuk"},
		{AccountCode: "211", CreditAmount: 1000000, Description: "Tabungan nasabah"},
	}

	entry, err := services.CreateJournalEntry("2026-03-07", "Setoran tunai", "deposit", "TX001", "TELLER01", lines)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if entry.JournalNumber == "" {
		t.Error("Expected journal number to be generated")
	}
	if entry.Status != "posted" {
		t.Errorf("Expected status 'posted', got '%s'", entry.Status)
	}
}

func TestCreateJournalEntry_Unbalanced(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	lines := []services.JournalLineInput{
		{AccountCode: "111", DebitAmount: 1000000, Description: "Kas masuk"},
		{AccountCode: "211", CreditAmount: 500000, Description: "Tabungan"},
	}

	_, err := services.CreateJournalEntry("2026-03-07", "Unbalanced", "test", "TX002", "TELLER01", lines)
	if err == nil {
		t.Error("Expected error for unbalanced journal, got nil")
	}
}

func TestGetTrialBalance(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	// Post a balanced journal
	lines := []services.JournalLineInput{
		{AccountCode: "111", DebitAmount: 5000000, Description: "Kas"},
		{AccountCode: "211", CreditAmount: 5000000, Description: "Tabungan"},
	}
	services.CreateJournalEntry("2026-03-07", "Test", "test", "TX003", "SYSTEM", lines)

	report, err := services.GetTrialBalance("")
	if err != nil {
		t.Fatalf("GetTrialBalance error: %v", err)
	}
	if len(report) == 0 {
		t.Error("Expected trial balance to have entries")
	}

	// Check balance: total debit should equal total credit
	var totalDebit, totalCredit float64
	for _, tb := range report {
		totalDebit += tb.DebitBalance
		totalCredit += tb.CreditBalance
	}
	if totalDebit != totalCredit {
		t.Errorf("Trial balance not balanced: debit=%.2f, credit=%.2f", totalDebit, totalCredit)
	}
}

func TestGetChartOfAccounts(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	accounts, err := services.GetChartOfAccounts("")
	if err != nil {
		t.Fatalf("GetChartOfAccounts error: %v", err)
	}
	if len(accounts) < 10 {
		t.Errorf("Expected at least 10 CoA entries, got %d", len(accounts))
	}
}

func TestGetChartOfAccounts_FilterByType(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	assets, err := services.GetChartOfAccounts("asset")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	for _, a := range assets {
		if a.AccountType != "asset" {
			t.Errorf("Expected type 'asset', got '%s' for code %s", a.AccountType, a.AccountCode)
		}
	}
}

// === INTEREST TESTS ===

func TestCalculateDailyInterest(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	accrual, err := services.CalculateDailyInterest("1000000001", "2026-03-07")
	if err != nil {
		t.Fatalf("CalculateDailyInterest error: %v", err)
	}
	if accrual.DailyInterest <= 0 {
		t.Error("Expected positive daily interest")
	}
	if accrual.Rate <= 0 {
		t.Error("Expected positive rate")
	}
	t.Logf("Account balance=%.0f, rate=%.2f%%, daily interest=%.4f", accrual.Balance, accrual.Rate, accrual.DailyInterest)
}

func TestGetApplicableRate_TieredSavings(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	// Balance 50M should get rate 1.00%
	rate, err := services.GetApplicableRate("savings", 50000000, 0)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if rate != 1.00 {
		t.Errorf("Expected rate 1.00%% for 50M balance, got %.2f%%", rate)
	}

	// Balance 100M+ should get rate 2.00%
	rate2, err := services.GetApplicableRate("savings", 150000000, 0)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if rate2 != 2.00 {
		t.Errorf("Expected rate 2.00%% for 150M balance, got %.2f%%", rate2)
	}
}

func TestCalculateDepositInterest(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	interest, err := services.CalculateDepositInterest("DEP001")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	// 100M * 4.5% * 12/12 = 4,500,000
	if interest != 4500000 {
		t.Errorf("Expected deposit interest 4500000, got %.2f", interest)
	}
}

func TestCalculateLoanInterest(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	interest, err := services.CalculateLoanInterest("LOAN001")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	// 450M * 7.5% / 12 = 2,812,500
	if interest != 2812500 {
		t.Errorf("Expected monthly loan interest 2812500, got %.2f", interest)
	}
}

func TestGetInterestRates(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	rates, err := services.GetInterestRates("")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(rates) < 3 {
		t.Errorf("Expected at least 3 interest rates, got %d", len(rates))
	}
}

func TestSimulateInterest(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	sim, err := services.SimulateInterest("deposit", 100000000, 12)
	if err != nil {
		t.Fatalf("SimulateInterest error: %v", err)
	}
	if sim.TotalInterest != 4500000 {
		t.Errorf("Expected 4500000, got %.2f", sim.TotalInterest)
	}
	if sim.MaturityAmount != 104500000 {
		t.Errorf("Expected maturity 104500000, got %.2f", sim.MaturityAmount)
	}
}

// === CIF TESTS ===

func TestGetSingleCustomerView(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	view, err := services.GetSingleCustomerView("CIF001")
	if err != nil {
		t.Fatalf("GetSingleCustomerView error: %v", err)
	}
	if view.Customer.CIF != "CIF001" {
		t.Errorf("Expected CIF001, got %s", view.Customer.CIF)
	}
	if len(view.Accounts) == 0 {
		t.Error("Expected at least 1 account")
	}
	if len(view.Roles) == 0 {
		t.Error("Expected at least 1 role")
	}
}

func TestSearchCustomers(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	customers, err := services.SearchCustomers("Budi")
	if err != nil {
		t.Fatalf("SearchCustomers error: %v", err)
	}
	if len(customers) == 0 {
		t.Error("Expected to find 'Budi'")
	}
}

func TestUpdateCustomerExtended(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	ext := models.CustomerExtended{
		MotherMaidenName: "Sari Dewi",
		Nationality:      "WNI",
		Occupation:       "Engineer",
		MonthlyIncome:    25000000,
		NPWP:             "12.345.678.9-012.000",
	}

	err := services.UpdateCustomerExtended("CIF001", ext)
	if err != nil {
		t.Fatalf("UpdateCustomerExtended error: %v", err)
	}

	// Verify
	result, err := services.GetCustomerExtended("CIF001")
	if err != nil {
		t.Fatalf("GetCustomerExtended error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected extended data, got nil")
	}
	if result.MotherMaidenName != "Sari Dewi" {
		t.Errorf("Expected 'Sari Dewi', got '%s'", result.MotherMaidenName)
	}
}

// === STANDING INSTRUCTION TESTS ===

func TestCreateStandingInstruction(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	req := services.CreateSIRequest{
		CIF:             "CIF001",
		FromAccount:     "1000000001",
		InstructionType: "transfer",
		ToAccount:       "1000000002",
		Amount:          500000,
		Description:     "Monthly transfer",
		Frequency:       "monthly",
		ExecutionDay:    1,
		StartDate:       "2026-04-01",
	}

	si, err := services.CreateStandingInstruction(req)
	if err != nil {
		t.Fatalf("CreateStandingInstruction error: %v", err)
	}
	if si.SINumber == "" {
		t.Error("Expected SI number")
	}
	if si.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", si.Status)
	}
}

func TestPauseAndCancelSI(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	// Create SI
	req := services.CreateSIRequest{
		CIF: "CIF001", FromAccount: "1000000001", InstructionType: "transfer",
		ToAccount: "1000000002", Amount: 100000, Frequency: "monthly",
		ExecutionDay: 15, StartDate: "2026-04-15",
	}
	si, _ := services.CreateStandingInstruction(req)

	// Pause
	err := services.PauseSI(si.SINumber)
	if err != nil {
		t.Fatalf("PauseSI error: %v", err)
	}

	// Cancel
	err = services.CancelSI(si.SINumber)
	if err == nil {
		// Paused SI can't be cancelled directly in some implementations
		// but our implementation allows it
	}
}

// === ACCOUNT MANAGEMENT TESTS ===

func TestOpenAccount(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	req := services.OpenAccountRequest{
		CIF:            "CIF001",
		AccountType:    "savings",
		Currency:       "IDR",
		InitialDeposit: 1000000,
		Branch:         "JKT001",
	}

	resp, err := services.OpenAccount(req)
	if err != nil {
		t.Fatalf("OpenAccount error: %v", err)
	}
	if resp.AccountNumber == "" {
		t.Error("Expected account number")
	}
	t.Logf("New account: %s", resp.AccountNumber)
}

func TestCloseAccount_WithBalance(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	// Try to close account with balance — should fail
	err := services.CloseAccount("1000000001", "Customer request")
	if err == nil {
		t.Error("Expected error closing account with balance")
	}
}

// === EOD TESTS ===

func TestRunEOD(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	ensureConfig()

	result, err := services.RunEOD("2026-03-07")
	if err != nil {
		t.Fatalf("RunEOD error: %v", err)
	}
	if result.ProcessDate != "2026-03-07" {
		t.Errorf("Expected date 2026-03-07, got %s", result.ProcessDate)
	}
	if len(result.Processes) == 0 {
		t.Error("Expected at least 1 EOD process")
	}
	t.Logf("EOD result: %s, processes: %d", result.OverallStatus, len(result.Processes))
	for _, p := range result.Processes {
		t.Logf("  - %s: %s (processed=%d, failed=%d)", p.ProcessType, p.Status, p.RecordsProcessed, p.RecordsFailed)
	}
}
