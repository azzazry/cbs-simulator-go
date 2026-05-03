# CBS Simulator - API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

> **Security:** Semua endpoint kecuali login, register, dan health check membutuhkan JWT token.
> Lihat [API_SECURITY.md](API_SECURITY.md) untuk dokumentasi lengkap fitur keamanan.

## Response Format

All endpoints return JSON with consistent format:

### Success Response
```json
{
  "status": "success",
  "data": { ... }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description"
}
```

---

## Authentication

Semua protected endpoint membutuhkan header:
```
Authorization: Bearer <access_token>
```

### Login
Authenticate customer with CIF and PIN. Returns JWT token pair.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "cif": "CIF001",
  "pin": "123456"
}
```

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "cif": "CIF001",
    "full_name": "Budi Santoso",
    "role": "admin",
    "message": "Login successful",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

**Error Responses:**
- 401: Invalid CIF or PIN
- 401: Account locked (3x gagal login). Lihat [Self-Service Unlock](API_SECURITY.md#step-1-verifikasi-e-kyc-ktp)

### Register
Create a new customer account with PIN policy enforcement.

**Endpoint:** `POST /auth/register`

**PIN Policy:** 6 digit, tidak boleh berurutan (123456), tidak boleh angka sama semua (111111).

**Request Body:**
```json
{
  "cif": "CIF006",
  "full_name": "John Doe",
  "id_card_number": "3201061234567895",
  "phone_number": "081234567895",
  "email": "john.doe@email.com",
  "address": "Jl. Merdeka No. 100, Jakarta",
  "date_of_birth": "1993-06-15",
  "pin": "246813"
}
```

**Success Response (HTTP 201):**
```json
{
  "status": "success",
  "data": {
    "cif": "CIF006",
    "full_name": "John Doe",
    "message": "Customer registered successfully"
  }
}
```

**Error Responses:**
- 400: CIF already exists
- 400: PIN policy violation

### Change PIN
Change customer PIN (requires authentication).

**Endpoint:** `POST /auth/change-pin`

**Headers:** `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "cif": "CIF001",
  "old_pin": "123456",
  "new_pin": "654321"
}
```

**Error Responses:**
- 400: Incorrect old PIN
- 400: PIN policy violation (new PIN)

### Logout
Invalidate current JWT token.

**Endpoint:** `POST /auth/logout`

**Headers:** `Authorization: Bearer <access_token>`

**Success Response:**
```json
{
  "status": "success",
  "message": "Logged out successfully"
}
```

### Refresh Token
Generate new token pair from refresh token.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Get Profile
Get authenticated user profile and roles.

**Endpoint:** `GET /auth/profile`

**Headers:** `Authorization: Bearer <access_token>`

> **More security endpoints:** Lihat [API_SECURITY.md](API_SECURITY.md) untuk e-KYC, OTP, unlock, dan reset PIN.

---

## Customer Profile

> **🔒 Semua endpoint di bawah membutuhkan:** `Authorization: Bearer <access_token>`

### Get Customer Profile
Get customer information.

**Endpoint:** `GET /customers/:cif`

**Example:** `GET /customers/CIF001`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "cif": "CIF001",
    "full_name": "Budi Santoso",
    "id_card_number": "3201011234567890",
    "phone_number": "081234567890",
    "email": "budi.santoso@email.com",
    "address": "Jl. Sudirman No. 123, Jakarta Selatan",
    "date_of_birth": "1985-03-15",
    "status": "active"
  }
}
```

---

## Accounts

### Get All Accounts by CIF
Get all accounts for a customer.

**Endpoint:** `GET /customers/:cif/accounts`

**Example:** `GET /customers/CIF001/accounts`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "account_number": "1001234567",
      "cif": "CIF001",
      "account_type": "savings",
      "currency": "IDR",
      "balance": 25000000.00,
      "avail_balance": 25000000.00,
      "status": "active",
      "opened_date": "2020-01-15",
      "branch": "Jakarta Sudirman"
    }
  ]
}
```

### Get Account Balance
Get specific account balance and details.

**Endpoint:** `GET /accounts/:account_number`

**Example:** `GET /accounts/1001234567`

### Get Account Statement
Get transaction history with pagination.

**Endpoint:** `GET /accounts/:account_number/statement?limit=20&offset=0`

**Query Parameters:**
- `limit` (optional): Number of records (default: 20)
- `offset` (optional): Pagination offset (default: 0)

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "transaction_id": "TRX20260301001",
      "transaction_type": "transfer_intra",
      "from_account_number": "1001234567",
      "to_account_number": "1001234568",
      "amount": 1000000.00,
      "currency": "IDR",
      "description": "Transfer ke Siti",
      "reference_number": "REF001",
      "status": "success",
      "transaction_date": "2026-03-01T10:30:00Z",
      "settlement_date": "2026-03-01",
      "fee": 0.00
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "count": 5
  }
}
```

---

## Transfers

### Intrabank Transfer
Transfer within the same bank (no fee).

**Endpoint:** `POST /transfers/intra`

**Request Body:**
```json
{
  "from_account_number": "1001234567",
  "to_account_number": "1001234568",
  "amount": 100000.00,
  "description": "Transfer untuk Siti"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Transfer successful",
  "data": {
    "id": 6,
    "transaction_id": "TRX20260306001",
    "transaction_type": "transfer_intra",
    "from_account_number": "1001234567",
    "to_account_number": "1001234568",
    "amount": 100000.00,
    "currency": "IDR",
    "description": "Transfer untuk Siti",
    "reference_number": "REF12345678",
    "status": "success",
    "transaction_date": "2026-03-06T10:30:00Z",
    "settlement_date": "2026-03-06",
    "fee": 0.00
  }
}
```

**Error Responses:**
- 400: Insufficient balance
- 404: Account not found
- 400: Account not active

### Interbank Transfer
Transfer to other bank (with dynamic fee based on destination bank configuration).

**Endpoint:** `POST /transfers/inter`

**Request Body:**
```json
{
  "from_account_number": "1001234567",
  "to_account_number": "1234567890",
  "destination_bank_code": "BCA",
  "amount": 500000.00,
  "description": "Transfer ke Bank Lain"
}
```

**Fee Structure (Dynamic - Configurable via Admin Panel):**
- **Default:** Rp 5,000 per transfer (domestic)
- **International:** Rp 10,000 per transfer
- **Fee Type:** Flat amount (can be changed to percentage)

**Supported Banks:**
- MANDIRI, BCA, BRI, CIMB, DANAMON, OCBC, MEGA, PERMATA
- UOB, COMMONWEALTH, PANIN, MAYBANK, BTN, SUMITOMO, DBS
- CITIBANK, HSBC, MIZUHO, JPMORGAN

**Error Responses:**
- 400: Insufficient balance (including fee)
- 404: Account not found
- 400: Account not active
- 400: Invalid bank code

---

## Bill Payment (PPOB)

### Pay Bill
Process bill payment.

**Endpoint:** `POST /bills/pay`

**Request Body:**
```json
{
  "account_number": "1001234567",
  "bill_number": "BILL202603001"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Bill payment successful",
  "data": {
    "transaction_id": "TRX20260306002",
    "amount": 452500.00,
    "description": "Payment for PT PLN (Persero) - 2026-02",
    "status": "success"
  }
}
```

### Get Bill History
Get bill payment history.

**Endpoint:** `GET /bills/history`

---

## Cards

### Get All Cards by CIF

**Endpoint:** `GET /cards/:cif`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "card_number": "4234********3456",
      "cif": "CIF001",
      "account_number": "1001234567",
      "card_type": "debit",
      "card_brand": "visa",
      "card_limit": 10000000.00,
      "avail_limit": 10000000.00,
      "expiry_date": "12/2026",
      "status": "active"
    }
  ]
}
```

### Block Card

**Endpoint:** `POST /cards/block`

**Request Body:**
```json
{
  "card_number": "4234567890123456"
}
```

### Unblock Card

**Endpoint:** `POST /cards/unblock`

**Request Body:**
```json
{
  "card_number": "4234567890123456"
}
```

---

## Loans

### Get All Loans by CIF

**Endpoint:** `GET /loans/:cif`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "loan_number": "LOAN001",
      "cif": "CIF002",
      "account_number": "3001234568",
      "loan_type": "mortgage",
      "principal_amount": 200000000.00,
      "outstanding_amount": 150000000.00,
      "interest_rate": 8.50,
      "monthly_payment": 2500000.00,
      "tenor_months": 120,
      "remaining_months": 72,
      "disbursement_date": "2021-03-10",
      "maturity_date": "2031-03-10",
      "next_payment_date": "2026-04-10",
      "status": "active"
    }
  ]
}
```

### Get Loan Details

**Endpoint:** `GET /loans/detail/:loan_number`

---

## Deposits

### Get All Deposits by CIF

**Endpoint:** `GET /deposits/:cif`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "deposit_number": "DEP001",
      "cif": "CIF001",
      "principal_amount": 50000000.00,
      "interest_rate": 5.50,
      "tenor_months": 12,
      "open_date": "2025-03-01",
      "maturity_date": "2026-03-01",
      "maturity_amount": 52750000.00,
      "auto_renew": true,
      "status": "active",
      "linked_account": "1001234567"
    }
  ]
}
```

### Get Deposit Details

**Endpoint:** `GET /deposits/detail/:deposit_number`

---

## Notifications

### Get Notification History

**Endpoint:** `GET /notifications/:cif?limit=20&offset=0`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "cif": "CIF001",
      "notification_type": "transfer",
      "title": "Transfer Berhasil",
      "message": "Anda mengirim Rp 100.000 ke 1001234568",
      "transaction_id": "TRX20260307001",
      "is_read": 0,
      "created_at": "2026-03-07T10:30:00Z"
    }
  ]
}
```

### Mark Notification as Read

**Endpoint:** `POST /notifications/read`

**Request Body:**
```json
{
  "notification_id": 1
}
```

---

## FCM & Device Management

> **🔒 Semua endpoint membutuhkan:** `Authorization: Bearer <access_token>`
> CIF diambil otomatis dari JWT token, tidak perlu dikirim di request body.

### Register Device Token
Daftarkan FCM token device untuk menerima push notification.

**Endpoint:** `POST /fcm/register`

**Request Body:**
```json
{
  "device_token": "fcm_token_dari_firebase_sdk",
  "device_type": "android",
  "device_name": "Samsung Galaxy S24"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Device token registered successfully"
}
```

> **Note:** Jika `device_token` sudah ada, data akan diupdate (upsert).

### Unregister Device Token
Nonaktifkan device token (misal saat logout dari device).

**Endpoint:** `DELETE /fcm/unregister`

**Request Body:**
```json
{
  "device_token": "fcm_token_yang_mau_dinonaktifkan"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Device token unregistered successfully"
}
```

### Get Registered Devices
Lihat semua device yang terdaftar untuk akun ini.

**Endpoint:** `GET /fcm/devices`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "cif": "CIF001",
      "device_token": "fcm_token_dari_firebase_sdk",
      "device_type": "android",
      "device_name": "Samsung Galaxy S24",
      "is_active": true,
      "created_at": "2026-05-04T01:00:00Z",
      "updated_at": "2026-05-04T01:00:00Z"
    }
  ]
}
```

**Testing tanpa mobile app (Postman):**
```bash
# Register device token dummy
curl -X POST http://localhost:8080/api/v1/fcm/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "device_token": "test-device-token-001",
    "device_type": "android",
    "device_name": "Test Device"
  }'
```

---

## Payments

### QRIS Payment

**Endpoint:** `POST /payments/qris`

### Virtual Account Payment

**Endpoint:** `POST /payments/va`

### E-Wallet Top Up

**Endpoint:** `POST /payments/ewallet/topup`

### E-Money Top Up

**Endpoint:** `POST /payments/emoney/topup`

---

## Health Check

**Endpoint:** `GET /health`

**Success Response:**
```json
{
  "status": "ok",
  "service": "CBS Simulator"
}
```

---

## Admin Management APIs

> **🔒 Admin endpoints membutuhkan role `admin` atau `supervisor`.**

### Audit Logs

**Endpoint:** `GET /admin/audit-logs`

### Transaction Limits

**Endpoint:** `GET /admin/transaction-limits`

**Endpoint:** `PUT /admin/transaction-limits`

### Role Management

**Endpoint:** `GET /admin/roles`

**Endpoint:** `POST /admin/roles/assign`

### Unlock Account

**Endpoint:** `POST /admin/unlock-account`

### EOD Processing

#### Jalankan EOD
**Endpoint:** `POST /admin/eod/run`

**Request Body:**
```json
{
  "process_date": "2026-03-07"
}
```

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "process_date": "2026-03-07",
    "overall_status": "completed",
    "processes": [
      {"process_type": "interest_accrual", "status": "completed", "records_processed": 5},
      {"process_type": "si_execution", "status": "completed", "records_processed": 2},
      {"process_type": "dormant_check", "status": "completed", "records_processed": 0}
    ]
  }
}
```

#### Status EOD
**Endpoint:** `GET /admin/eod/status/:date`

#### Riwayat EOD
**Endpoint:** `GET /admin/eod/history`

---

## Phase 2: Core Banking Endpoints

### General Ledger

#### Chart of Accounts
**Endpoint:** `GET /gl/chart-of-accounts`

#### Journal Entries
**Endpoint:** `GET /gl/journal-entries`

#### Detail Jurnal
**Endpoint:** `GET /gl/journal-entries/:id`

#### Trial Balance
**Endpoint:** `GET /gl/trial-balance`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "accounts": [
      {"account_code": "111", "account_name": "Kas", "debit_balance": 1000000, "credit_balance": 0}
    ],
    "total_debit": 1000000,
    "total_credit": 1000000,
    "is_balanced": true
  }
}
```

#### Saldo Akun GL
**Endpoint:** `GET /gl/account-balance/:code`

---

### CIF Enhancement

#### Single Customer View
**Endpoint:** `GET /customers/:cif/overview`

#### Update Data Tambahan Nasabah
**Endpoint:** `PUT /customers/:cif/extended`

**Request Body:**
```json
{
  "mother_maiden_name": "Sari Dewi",
  "nationality": "WNI",
  "occupation": "Engineer",
  "monthly_income": 25000000,
  "npwp": "12.345.678.9-012.000"
}
```

#### Cari Nasabah
**Endpoint:** `GET /customers/search?q=Budi`

---

### Bunga & Simulasi

#### Daftar Suku Bunga
**Endpoint:** `GET /interest/rates`

#### Simulasi Bunga
**Endpoint:** `POST /interest/calculate`

**Request Body:**
```json
{
  "product_type": "deposit",
  "principal": 100000000,
  "tenor_months": 12
}
```

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "product_type": "deposit",
    "principal": 100000000,
    "rate": 4.50,
    "tenor_months": 12,
    "total_interest": 4500000,
    "maturity_amount": 104500000,
    "monthly_interest": 375000
  }
}
```

---

### Standing Instructions

#### Buat SI Baru
**Endpoint:** `POST /standing-instructions`

**Request Body:**
```json
{
  "cif": "CIF001",
  "from_account": "1001234567",
  "instruction_type": "transfer",
  "to_account": "1001234568",
  "amount": 500000,
  "frequency": "monthly",
  "execution_day": 1,
  "start_date": "2026-04-01"
}
```

#### Daftar SI Nasabah
**Endpoint:** `GET /standing-instructions/by-cif/:cif`

#### Pause SI
**Endpoint:** `PUT /standing-instructions/:id/pause`

#### Batalkan SI
**Endpoint:** `DELETE /standing-instructions/:id`

#### Riwayat Eksekusi SI
**Endpoint:** `GET /standing-instructions/:id/history`

---

### Account Management

#### Buka Rekening
**Endpoint:** `POST /accounts/open`

**Request Body:**
```json
{
  "cif": "CIF001",
  "account_type": "savings",
  "currency": "IDR",
  "initial_deposit": 1000000,
  "branch": "JKT001"
}
```

#### Tutup Rekening
**Endpoint:** `POST /accounts/:account_number/close`

> Rekening harus bersaldo 0 untuk bisa ditutup.

#### Daftar Rekening Dormant
**Endpoint:** `GET /accounts/dormant`

#### Aktifkan Kembali
**Endpoint:** `POST /accounts/:account_number/reactivate`

---

## Error Codes

| HTTP Status | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input atau business logic error |
| 401 | Unauthorized - Missing/invalid/expired JWT token |
| 403 | Forbidden - Insufficient role permissions |
| 404 | Not Found |
| 429 | Too Many Requests - Rate limit exceeded (60 req/min) |
| 500 | Internal Server Error |

---

## Testing dengan cURL

```bash
# 1. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"cif": "CIF001", "pin": "123456"}'

export TOKEN="eyJhbGci..."

# 2. Register device FCM
curl -X POST http://localhost:8080/api/v1/fcm/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_token": "test-device-001", "device_type": "android", "device_name": "Test Device"}'

# 3. Lihat devices
curl http://localhost:8080/api/v1/fcm/devices \
  -H "Authorization: Bearer $TOKEN"

# 4. Transfer (push notification akan terkirim ke device terdaftar)
curl -X POST http://localhost:8080/api/v1/transfers/intra \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_account_number": "1001234567", "to_account_number": "1001234568", "amount": 100000, "description": "Test"}'

# 5. Unregister device
curl -X DELETE http://localhost:8080/api/v1/fcm/unregister \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_token": "test-device-001"}'
```

---

## Related Documentation

- [API_SECURITY.md](API_SECURITY.md) — Security endpoints (JWT, OTP, e-KYC, RBAC, audit)

---

**Last Updated:** Mei 2026 (FCM Device Management Update)