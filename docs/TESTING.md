# CBS Simulator - Testing Guide

## Quick Start

```bash
# Jalankan semua test
go test ./...

# Jalankan test security services saja
go test ./services/

# Verbose mode (detail per-test)
go test ./services/ -v

# Jalankan test spesifik
go test ./services/ -v -run TestAuthenticate_Lockout
```

---

## Test Coverage

### Security Services (`services/*_test.go`)

| Test File | Tests | Coverage |
|-----------|-------|----------|
| `pin_policy_test.go` | 8 | PIN validation: valid, too short/long, non-digit, sequential, repeating, empty |
| `jwt_service_test.go` | 5 | Token generate, validate, blacklist, refresh, wrong token type |
| `rbac_service_test.go` | 9 | Get/has/primary role, assign role, get all roles, role by name |
| `otp_service_test.go` | 5 | OTP generate, verify valid/wrong/used/wrong type/wrong CIF |
| `auth_service_test.go` | 12 | Login, lockout 3x, register, change PIN, logout, get customer |
| `setup_test.go` | — | In-memory SQLite database + seed data helper |
| **Total** | **39** | |

---

## Test Infrastructure

Tests menggunakan **in-memory SQLite** (`:memory:`) sehingga:
- ✅ Tidak butuh file database
- ✅ Setiap test bersih (fresh database)
- ✅ Cepat (~5 detik untuk 39 test)

### Test Data

| CIF | Nama | PIN | Role |
|-----|------|-----|------|
| `CIF001` | Budi Santoso | `123456` | admin, customer |
| `CIF002` | Siti Nurhaliza | `123456` | customer |
| `CIF003` | Ahmad Wijaya | `123456` | supervisor, customer |

---

## Menjalankan Test Tertentu

```bash
# PIN Policy saja
go test ./services/ -v -run TestValidatePINPolicy

# JWT saja
go test ./services/ -v -run TestGenerateTokenPair
go test ./services/ -v -run TestValidateToken
go test ./services/ -v -run TestRefreshAccessToken

# RBAC saja
go test ./services/ -v -run TestHasRole
go test ./services/ -v -run TestGetPrimaryRole
go test ./services/ -v -run TestAssignRole

# OTP saja
go test ./services/ -v -run TestGenerateOTP
go test ./services/ -v -run TestVerifyOTP

# Auth (login, lockout, register, dll)
go test ./services/ -v -run TestAuthenticate
go test ./services/ -v -run TestRegisterCustomer
go test ./services/ -v -run TestChangePIN
go test ./services/ -v -run TestLogoutUser
```

---

## Expected Output

```
=== RUN   TestAuthenticate_Success
--- PASS: TestAuthenticate_Success (0.11s)
=== RUN   TestAuthenticate_Lockout3Attempts
--- PASS: TestAuthenticate_Lockout3Attempts (0.21s)
...
PASS
ok   cbs-simulator/services    5.722s
```

Semua test harus **PASS**. Jika ada yang **FAIL**, cek error message untuk detail.
