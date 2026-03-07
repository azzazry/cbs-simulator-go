# CBS Simulator - API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

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

### Login
Authenticate customer with CIF and PIN.

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
    "message": "Login successful",
    "token": "token_CIF001"
  }
}
```

### Register
Create a new customer account.

**Endpoint:** `POST /auth/register`

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
  "pin": "123456"
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

**Error Response (CIF already exists):**
```json
{
  "status": "error",
  "message": "CIF already exists"
}
```

### Change PIN
Change customer PIN.

**Endpoint:** `POST /auth/change-pin`

**Request Body:**
```json
{
  "cif": "CIF001",
  "old_pin": "123456",
  "new_pin": "654321"
}
```

---

## Customer Profile

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
Transfer to other bank (with fee).

**Endpoint:** `POST /transfers/inter`

**Request Body:**
```json
{
  "from_account_number": "1001234567",
  "to_account_number": "1234567890",
  "amount": 500000.00,
  "description": "Transfer ke Bank Lain",
  "transfer_type": "skn"
}
```

**Transfer Types:**
- `inter`: Regular interbank (fee: Rp 6,500)
- `skn`: SKN - Sistem Kliring Nasional (fee: Rp 3,500)
- `rtgs`: RTGS - Real Time Gross Settlement (fee: Rp 25,000)

**Success Response:** Same format as intrabank transfer, with fee applied.

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
| 400 | Bad Request - Invalid input or business logic error |
| 401 | Unauthorized - Invalid credentials |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

---

## Testing with cURL

### Complete Flow Example

```bash
# 1. Health Check
curl http://localhost:8080/health

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"cif": "CIF001", "pin": "123456"}'

# 3. Get Accounts
curl http://localhost:8080/api/v1/customers/CIF001/accounts

# 4. Check Balance
curl http://localhost:8080/api/v1/accounts/1001234567

# 5. Transfer Money
curl -X POST http://localhost:8080/api/v1/transfers/intra \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_number": "1001234567",
    "to_account_number": "1001234568",
    "amount": 100000,
    "description": "Test transfer"
  }'

# 6. Check Statement
curl "http://localhost:8080/api/v1/accounts/1001234567/statement?limit=5"

# 7. Get Bills
curl "http://localhost:8080/api/v1/bills/inquiry?biller_code=PLN&customer_number=123456789012"

# 8. Pay Bill
curl -X POST http://localhost:8080/api/v1/bills/pay \
  -H "Content-Type: application/json" \
  -d '{
    "account_number": "1001234567",
    "bill_number": "BILL202603001"
  }'
```

---

## Postman Collection

Import this JSON to Postman for quick testing:

[Link to Postman collection would go here]

---

**Last Updated:** March 2026
