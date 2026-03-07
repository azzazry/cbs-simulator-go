# CBS Simulator - Testing Guide

## Quick Start

```bash
# Jalankan semua 60 test
go test ./services/ -v

# Jalankan test spesifik
go test ./services/ -v -run TestCreateJournalEntry

# Build check
go build ./...
go vet ./...
```

---

## Test Coverage

### Phase 1: Security Services

| Test File | Tests | Coverage |
|-----------|-------|----------|
| `pin_policy_test.go` | 8 | PIN validation: valid, too short/long, non-digit, sequential, repeating, empty |
| `jwt_service_test.go` | 5 | Token generate, validate, blacklist, refresh, wrong token type |
| `rbac_service_test.go` | 9 | Get/has/primary role, assign role, get all roles, role by name |
| `otp_service_test.go` | 5 | OTP generate, verify valid/wrong/used/wrong type/wrong CIF |
| `auth_service_test.go` | 12 | Login, lockout 3x, register, change PIN, logout, get customer |
| **Subtotal** | **39** | |

### Phase 2: Core Banking Services

| Test File | Tests | Coverage |
|-----------|-------|----------|
| `core_banking_test.go` | 21 | GL journal (balanced/unbalanced), trial balance, CoA, interest (daily, tiered, deposit, loan, simulate), CIF (single view, search, extended), SI (create, pause/cancel), account (open, close), EOD (full run) |
| **Subtotal** | **21** | |

### Total: 60 tests

---

## Test Infrastructure

Tests menggunakan **in-memory SQLite** (`:memory:`) sehingga:
- Tidak butuh file database
- Setiap test bersih (fresh database)
- Cepat (~4 detik untuk 60 test)

### Test Data

| CIF | Nama | PIN | Role | Akun |
|-----|------|-----|------|------|
| `CIF001` | Budi Santoso | `123456` | admin, customer | 1000000001 (50jt) |
| `CIF002` | Siti Nurhaliza | `123456` | customer | 1000000002 (25jt) |
| `CIF003` | Ahmad Wijaya | `123456` | supervisor, customer | 1000000003 (100jt) |

### Seed Data Phase 2
- 10 Chart of Accounts (kode 3-digit PAPI)
- 3 Interest rates (tiered savings + deposito)
- 1 Deposito (DEP001, 100jt, 4.5%, 12 bulan)
- 1 Kredit (LOAN001, KPR 500jt, 7.5%)

---

## Menjalankan Test Tertentu

```bash
# Phase 1: Security
go test ./services/ -v -run TestAuthenticate
go test ./services/ -v -run TestValidatePINPolicy
go test ./services/ -v -run TestGenerateTokenPair
go test ./services/ -v -run TestHasRole
go test ./services/ -v -run TestVerifyOTP

# Phase 2: General Ledger
go test ./services/ -v -run TestCreateJournalEntry
go test ./services/ -v -run TestGetTrialBalance
go test ./services/ -v -run TestGetChartOfAccounts

# Phase 2: Interest
go test ./services/ -v -run TestCalculateDailyInterest
go test ./services/ -v -run TestGetApplicableRate
go test ./services/ -v -run TestSimulateInterest

# Phase 2: CIF, SI, Account, EOD
go test ./services/ -v -run TestGetSingleCustomerView
go test ./services/ -v -run TestCreateStandingInstruction
go test ./services/ -v -run TestOpenAccount
go test ./services/ -v -run TestRunEOD
```

---

## Expected Output

```
=== RUN   TestAuthenticate_Success
--- PASS: TestAuthenticate_Success (0.11s)
=== RUN   TestCreateJournalEntry_Balanced
--- PASS: TestCreateJournalEntry_Balanced (0.06s)
=== RUN   TestCalculateDailyInterest
    core_banking_test.go:127: Account balance=50000000, rate=1.00%, daily interest=1369.8630
--- PASS: TestCalculateDailyInterest (0.06s)
=== RUN   TestRunEOD
    core_banking_test.go:394: EOD result: completed, processes: 3
--- PASS: TestRunEOD (0.06s)
...
PASS
ok   cbs-simulator/services    4.246s
```

Semua test harus **PASS**. Jika ada yang **FAIL**, cek error message untuk detail.
