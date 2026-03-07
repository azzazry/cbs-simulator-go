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

**Headers:** `Authorization: Bearer <access_token>`

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
- **Application:** Real-time calculation before deduction
- **Management:** Use `/api/v1/admin/fees/transfer/calculate` to check fee for any amount

**Example with Fee Calculation:**

Request:
```bash
POST /api/v1/admin/fees/transfer/calculate
Content-Type: application/json

{
  "destination_bank_code": "BCA",
  "amount": 500000
}
```

Response:
```json
{
  "status": "success",
  "data": {
    "destination_bank_code": "BCA",
    "destination_bank_name": "Bank BCA",
    "transfer_amount": 500000,
    "fee": 5000,
    "total_amount": 505000
  }
}
```

**Transfer Response (with fee deducted):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "transaction_id": "TRX20260307XYZ",
    "transaction_type": "transfer_inter",
    "from_account_number": "1001234567",
    "to_account_number": "1234567890",
    "amount": 500000.00,
    "currency": "IDR",
    "description": "Transfer ke Bank Lain",
    "reference_number": "REF20260307001",
    "status": "success",
    "transaction_date": "2026-03-07T10:30:00Z",
    "settlement_date": "2026-03-07",
    "fee": 5000.00,
    "created_at": "2026-03-07T10:30:00Z"
  }
}
```

**Supported Banks:**
- MANDIRI, BCA, BRI, CIMB, DANAMON, OCBC, MEGA, PERMATA
- UOB, COMMONWEALTH, PANIN, MAYBANK, BTN, SUMITOMO, DBS
- CITIBANK, HSBC, MIZUHO, JPMORGAN

**Possible Errors:**
- 400: Insufficient balance (including fee)
- 404: Account not found
- 400: Account not active
- 400: Invalid bank code

### Get Transaction Details
Get specific transaction information.

**Endpoint:** `GET /transfers/:transaction_id`

**Example:** `GET /transfers/TRX20260301001`

---

## Bill Payment (PPOB)

### Get Biller List
Get list of available billers.

**Endpoint:** `GET /bills/billers`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "code": "PLN",
      "name": "PT PLN (Persero)",
      "category": "Electricity"
    },
    {
      "code": "PDAM",
      "name": "PDAM",
      "category": "Water"
    }
  ]
}
```

### Bill Inquiry
Check bill amount before payment.

**Endpoint:** `GET /bills/inquiry?biller_code=PLN&customer_number=123456789012`

**Query Parameters:**
- `biller_code`: Biller code (PLN, PDAM, TELKOM, etc)
- `customer_number`: Customer ID at biller

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "biller_code": "PLN",
    "biller_name": "PT PLN (Persero)",
    "customer_number": "123456789012",
    "bill_number": "BILL202603001",
    "bill_amount": 450000.00,
    "admin_fee": 2500.00,
    "total_amount": 452500.00,
    "bill_period": "2026-02",
    "due_date": "2026-03-20",
    "status": "unpaid"
  }
}
```

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

---

## Cards

### Get All Cards by CIF
Get all cards for a customer.

**Endpoint:** `GET /customers/:cif/cards`

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

### Get Card Details
Get specific card information.

**Endpoint:** `GET /cards/:card_number`

### Block Card
Block a card for security.

**Endpoint:** `POST /cards/block`

**Request Body:**
```json
{
  "card_number": "4234567890123456"
}
```

### Unblock Card
Unblock a previously blocked card.

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
Get all loans for a customer.

**Endpoint:** `GET /customers/:cif/loans`

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
Get specific loan information.

**Endpoint:** `GET /loans/:loan_number`

---

## Deposits

### Get All Deposits by CIF
Get all time deposits for a customer.

**Endpoint:** `GET /customers/:cif/deposits`

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
Get specific deposit information.

**Endpoint:** `GET /deposits/:deposit_number`

---

## Health Check

### Server Health
Check if server is running.

**Endpoint:** `GET /health`

**Success Response:**
```json
{
  "status": "healthy",
  "service": "CBS Simulator",
  "version": "1.0.0"
}
```

---

## Notifications

### Get Notification History
Get notification history with pagination.

**Endpoint:** `GET /notifications/:cif?limit=20&offset=0`

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

### Get Unread Notification Count
Get count of unread notifications.

**Endpoint:** `GET /notifications/:cif/count`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "unread_count": 5
  }
}
```

### Mark Notification as Read
Mark a notification as read.

**Endpoint:** `POST /notifications/read`

**Request Body:**
```json
{
  "notification_id": 1
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Notification marked as read"
}
```

### Register FCM Token
Register device FCM token for push notifications.

**Endpoint:** `POST /notifications/fcm-token`

**Request Body:**
```json
{
  "cif": "CIF001",
  "device_token": "exxxxxxxxxxxxxxxxxxxxxxxxx",
  "device_type": "android",
  "device_name": "Samsung Galaxy A12"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "FCM token registered successfully"
}
```

### Get Notification Preferences
Get user notification settings.

**Endpoint:** `GET /notifications/:cif/preferences`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "cif": "CIF001",
    "transfer_notification": 1,
    "payment_notification": 1,
    "deposit_notification": 1,
    "loan_notification": 1,
    "promotion_notification": 1,
    "updated_at": "2026-03-07T10:00:00Z"
  }
}
```

### Update Notification Preferences
Update user notification settings.

**Endpoint:** `PUT /notifications/:cif/preferences`

**Request Body:**
```json
{
  "transfer_notification": 1,
  "payment_notification": 1,
  "deposit_notification": 0,
  "loan_notification": 1,
  "promotion_notification": 0
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Notification preferences updated successfully"
}
```

---

## Error Codes

| HTTP Status | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input or business logic error |
| 401 | Unauthorized - Missing/invalid/expired JWT token |
| 403 | Forbidden - Insufficient role permissions (RBAC) |
| 404 | Not Found - Resource not found |
| 429 | Too Many Requests - Rate limit exceeded (60 req/min) |
| 500 | Internal Server Error |

---

## Testing with cURL

### Complete Flow Example

```bash
# 1. Health Check
curl http://localhost:8080/health

# 2. Login (simpan access_token dari response)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"cif": "CIF001", "pin": "123456"}'

# Set token (ganti dengan access_token dari response login)
export TOKEN="eyJhbGci..."

# 3. Get Profile
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"

# 4. Get Accounts
curl http://localhost:8080/api/v1/customers/CIF001/accounts \
  -H "Authorization: Bearer $TOKEN"

# 5. Check Balance
curl http://localhost:8080/api/v1/accounts/1001234567 \
  -H "Authorization: Bearer $TOKEN"

# 6. Transfer Money
curl -X POST http://localhost:8080/api/v1/transfers/intra \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_number": "1001234567",
    "to_account_number": "1001234568",
    "amount": 100000,
    "description": "Test transfer"
  }'

# 7. Check Statement
curl "http://localhost:8080/api/v1/accounts/1001234567/statement?limit=5" \
  -H "Authorization: Bearer $TOKEN"

# 8. Pay Bill
curl -X POST http://localhost:8080/api/v1/bills/pay \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "account_number": "1001234567",
    "bill_number": "BILL202603001"
  }'

# 9. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

---

## Admin Management APIs

> **🔒 Admin endpoints membutuhkan role `admin` atau `supervisor`.**
> Lihat [API_SECURITY.md](API_SECURITY.md) untuk endpoint keamanan admin (audit logs, roles, transaction limits).

### Get All Supported Banks
Get list of all supported banks for transfers.

**Endpoint:** `GET /admin/banks`

**Headers:** `Authorization: Bearer <admin_token>`

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "bank_code": "MANDIRI",
      "bank_name": "Bank Mandiri",
      "swift_code": "BMRIIDJA",
      "is_active": 1,
      "created_at": "2026-03-07T00:00:00Z",
      "updated_at": "2026-03-07T00:00:00Z"
    },
    {
      "id": 2,
      "bank_code": "BCA",
      "bank_name": "Bank BCA",
      "swift_code": "BCAIDJA",
      "is_active": 1,
      "created_at": "2026-03-07T00:00:00Z",
      "updated_at": "2026-03-07T00:00:00Z"
    }
  ]
}
```

### Transfer Fee Management

#### Get All Transfer Fees
**Endpoint:** `GET /admin/fees/transfer`

Returns all configured transfer fees between bank pairs with current fee amounts.

#### Update Transfer Fee
Change fee for transfers between specific banks.

**Endpoint:** `PUT /admin/fees/transfer`

**Request:**
```json
{
  "from_bank_code": "MANDIRI",
  "to_bank_code": "BCA",
  "fee_amount": 7500,
  "fee_type": "flat"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Transfer fee updated successfully",
  "data": {
    "from_bank_code": "MANDIRI",
    "to_bank_code": "BCA",
    "fee_amount": 7500,
    "fee_type": "flat"
  }
}
```

#### Calculate Transfer Fee
Calculate the fee for a specific transfer amount between two banks.

**Endpoint:** `POST /admin/fees/transfer/calculate`

**Request:**
```json
{
  "from_bank_code": "MANDIRI",
  "to_bank_code": "BCA",
  "amount": 1000000
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "from_bank_code": "MANDIRI",
    "to_bank_code": "BCA",
    "transfer_amount": 1000000,
    "fee": 5000,
    "total_amount": 1005000
  }
}
```

### Service Fee Management

#### Get All Service Fees
Get fees for all services (e-wallet, e-money, VA, QRIS, etc).

**Endpoint:** `GET /admin/fees/services?type=topup_ewallet`

Optional query parameter `type` to filter by service type:
- `topup_ewallet` - E-wallet top-ups (OVO, DANA, GoPay)
- `topup_emoney` - E-money top-ups (LinkAja, Mandiri e-Money)
- `payment_va` - Virtual Account payments
- `qris_payment` - QRIS payments

#### Update Service Fee
Change fee configuration for a service.

**Endpoint:** `PUT /admin/fees/services`

**Request (Flat Fee):**
```json
{
  "service_code": "TOPUP_OVO",
  "fee_amount": 3000,
  "fee_percentage": 0,
  "fee_type": "flat"
}
```

**Request (Percentage Fee):**
```json
{
  "service_code": "QRIS_PAYMENT",
  "fee_amount": 0,
  "fee_percentage": 1.5,
  "fee_type": "percentage"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Service fee updated successfully",
  "data": {
    "service_code": "TOPUP_OVO",
    "fee_amount": 3000,
    "fee_percentage": 0,
    "fee_type": "flat"
  }
}
```

#### Calculate Service Fee
Calculate fee for a service transaction.

**Endpoint:** `POST /admin/fees/services/calculate`

**Request:**
```json
{
  "service_code": "TOPUP_OVO",
  "amount": 100000
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "service_code": "TOPUP_OVO",
    "service_amount": 100000,
    "fee": 2500,
    "total_amount": 102500
  }
}
```

#### Get Fee Statistics
View fee configuration statistics.

**Endpoint:** `GET /admin/fees/statistics?type=transfer`

Query parameters:
- `type` - `transfer` or `service` (default: transfer)
- `service_type` - Optional filter for service fees by type

**Response:**
```json
{
  "status": "success",
  "data": {
    "total_transfer_routes": 21,
    "transfers": [...]
  }
}
```

---

## Service Codes Reference

### E-Wallet Services
| Code | Name | Fee |
|------|------|-----|
| TOPUP_OVO | OVO Top-Up | Rp 2,500 |
| TOPUP_DANA | DANA Top-Up | Rp 2,500 |
| TOPUP_GOPAY | GoPay Top-Up | Rp 2,500 |

### E-Money Services
| Code | Name | Fee |
|------|------|-----|
| TOPUP_LINKAJA | LinkAja Top-Up | Rp 2,500 |
| TOPUP_MANDIRIEMONEY | Mandiri e-Money Top-Up | Rp 2,500 |

### Virtual Account Services
| Code | Name | Fee |
|------|------|-----|
| PAYMENT_VA_MANDIRI | Mandiri VA Payment | Free |
| PAYMENT_VA_BCA | BCA VA Payment | Free |
| PAYMENT_VA_BRI | BRI VA Payment | Free |

### Digital Payment Services
| Code | Name | Fee |
|------|------|-----|
| QRIS_PAYMENT | QRIS Payment | 1% |

---

## Postman Collection

Import this JSON to Postman for quick testing:

[Link to Postman collection would go here]

---

## Related Documentation

- [API_SECURITY.md](API_SECURITY.md) — Security endpoints (JWT, OTP, e-KYC, RBAC, audit)

---

## Phase 2: Core Banking Endpoints

### General Ledger

#### Daftar Chart of Accounts
**Endpoint:** `GET /gl/chart-of-accounts?type=asset`

```bash
curl http://localhost:8080/api/v1/gl/chart-of-accounts -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {"account_code": "111", "account_name": "Kas", "account_type": "asset", "normal_balance": "debit"}
  ]
}
```

#### Journal Entries
**Endpoint:** `GET /gl/journal-entries?date_from=2026-01-01&date_to=2026-12-31&page=1`

#### Detail Jurnal
**Endpoint:** `GET /gl/journal-entries/:id`

#### Trial Balance (Neraca Saldo)
**Endpoint:** `GET /gl/trial-balance?date=2026-03-07`

```bash
curl http://localhost:8080/api/v1/gl/trial-balance -H "Authorization: Bearer $TOKEN"
```

**Response:**
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

```bash
curl http://localhost:8080/api/v1/customers/CIF001/overview -H "Authorization: Bearer $TOKEN"
```

**Response:** Data lengkap nasabah termasuk accounts, loans, deposits, cards, roles, dan data tambahan.

#### Update Data Tambahan Nasabah
**Endpoint:** `PUT /customers/:cif/extended`

```bash
curl -X PUT http://localhost:8080/api/v1/customers/CIF001/extended \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"mother_maiden_name":"Sari Dewi","nationality":"WNI","occupation":"Engineer","monthly_income":25000000,"npwp":"12.345.678.9-012.000"}'
```

#### Cari Nasabah
**Endpoint:** `GET /customers/search?q=Budi`

---

### Bunga & Simulasi

#### Daftar Suku Bunga
**Endpoint:** `GET /interest/rates?product_type=savings`

**Response:**
```json
{
  "status": "success",
  "data": [
    {"product_type": "savings", "product_name": "Tabungan Reguler", "base_rate": 1.00, "min_balance": 0, "max_balance": 100000000},
    {"product_type": "savings", "product_name": "Tabungan Reguler", "base_rate": 2.00, "min_balance": 100000000}
  ]
}
```

#### Simulasi Bunga
**Endpoint:** `POST /interest/calculate`

```bash
curl -X POST http://localhost:8080/api/v1/interest/calculate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_type":"deposit","principal":100000000,"tenor_months":12}'
```

**Response:**
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

```bash
curl -X POST http://localhost:8080/api/v1/standing-instructions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"cif":"CIF001","from_account":"1001234567","instruction_type":"transfer","to_account":"1001234568","amount":500000,"frequency":"monthly","execution_day":1,"start_date":"2026-04-01"}'
```

#### Daftar SI Nasabah
**Endpoint:** `GET /standing-instructions/:cif`

#### Pause SI
**Endpoint:** `PUT /standing-instructions/:si_number/pause`

#### Batalkan SI
**Endpoint:** `DELETE /standing-instructions/:si_number`

#### Riwayat Eksekusi SI
**Endpoint:** `GET /standing-instructions/:si_number/history`

---

### Account Management

#### Buka Rekening
**Endpoint:** `POST /accounts/open`

```bash
curl -X POST http://localhost:8080/api/v1/accounts/open \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"cif":"CIF001","account_type":"savings","currency":"IDR","initial_deposit":1000000,"branch":"JKT001"}'
```

#### Tutup Rekening
**Endpoint:** `POST /accounts/:account_number/close`
> Rekening harus bersaldo 0 untuk bisa ditutup.

#### Daftar Rekening Dormant
**Endpoint:** `GET /accounts/dormant`

#### Aktifkan Kembali
**Endpoint:** `POST /accounts/:account_number/reactivate`

---

### EOD Processing (Admin Only)

#### Jalankan EOD
**Endpoint:** `POST /admin/eod/run`

```bash
curl -X POST http://localhost:8080/api/v1/admin/eod/run \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"process_date":"2026-03-07"}'
```

**Response:**
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
**Endpoint:** `GET /admin/eod/history?date_from=2026-01-01&date_to=2026-12-31`

---

**Last Updated:** Maret 2026 (Phase 2 Core Banking Update)
