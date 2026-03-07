# Banking Fees & Interbank Transfer Guide (Single Bank Model)

Panduan lengkap tentang sistem fee dinamis untuk transfer antar bank dan berbagai layanan di CBS Simulator.

**IMPORTANT ARCHITECTURE NOTE:**  
Sistem ini merepresentasikan **SATU BANK SPESIFIK** (e.g., Bank Daerah Jaya, Bank Pribadi, dll).  
- ✅ **Intrabank Transfer:** Rp 0 (GRATIS)  
- ✅ **Outbound to Other Banks:** Dynamic fee (Rp 5,000-10,000)  
- ✅ **Inbound from Other Banks:** Rp 0 (GRATIS - pengirim bayar)  

Fee hanya disimpan untuk **DESTINATION BANKS** tempat customer kami mengirim uang.

---

## 📋 Daftar Bank Tujuan Transfer (Destination Banks)

Customer kami dapat transfer OUT ke 19 bank utama di Indonesia:

| Bank Code | Bank Name | Fee | Category |
|-----------|-----------|-----|----------|
| MANDIRI | Bank Mandiri | Rp 5,000 | Domestic (SKNT) |
| BCA | Bank BCA | Rp 5,000 | Domestic (SKNT) |
| BRI | Bank BRI | Rp 5,000 | Domestic (SKNT) |
| CIMB | Bank CIMB Niaga | Rp 5,000 | Domestic (SKNT) |
| DANAMON | Bank Danamon | Rp 5,000 | Domestic (SKNT) |
| OCBC | Bank OCBC NISP | Rp 5,000 | Domestic (SKNT) |
| MEGA | Bank Mega | Rp 5,000 | Domestic (SKNT) |
| PERMATA | Bank Permata | Rp 5,000 | Domestic (SKNT) |
| UOB | Bank UOB Indonesia | Rp 5,000 | Domestic (SKNT) |
| COMMONWEALTH | Bank Commonwealth | Rp 5,000 | Domestic (SKNT) |
| PANIN | Bank Panin | Rp 5,000 | Domestic (SKNT) |
| MAYBANK | Maybank Indonesia | Rp 5,000 | Domestic (SKNT) |
| BTN | Bank BTN | Rp 5,000 | Domestic (SKNT) |
| SUMITOMO | Bank BTMU | Rp 5,000 | Domestic (SKNT) |
| DBS | DBS Bank | Rp 5,000 | Domestic (SKNT) |
| CITIBANK | Citibank Indonesia | Rp 10,000 | International (SWIFT) |
| HSBC | HSBC Bank Indonesia | Rp 10,000 | International (SWIFT) |
| MIZUHO | Bank Mizuho Indonesia | Rp 10,000 | International (SWIFT) |
| JPMORGAN | JP Morgan Chase | Rp 10,000 | International (SWIFT) |

---

## 💰 Fee Structure (Struktur Biaya)

### A. Transfer Fees (Biaya Transfer OUTBOUND)

#### Core Logic:
```
✓ Intrabank Transfer (Rek 1001 → 1002):    Rp 0 (GRATIS)
✓ Outbound Transfer (Our Bank → BCA):      Rp 5,000 (Dynamic per destination bank)
✓ Inbound Transfer (Mandiri → Our Bank):   Rp 0 (GRATIS - pengirim yang bayar)
```

#### Simplified Destination-Based Model:
Untuk setiap **destination bank**, kami simpan **1 fee entry**:
```
Example 1: BCA as destination
- Destination Bank Code: BCA
- Destination Bank Name: Bank BCA
- Fee Type: Flat
- Fee Amount: Rp 5,000
- Min/Max: Rp 0 - Rp 999,999,999
- Effective: Now - NULL (no expiration)

Example 2: CITIBANK as destination (International)
- Destination Bank Code: CITIBANK
- Destination Bank Name: Citibank Indonesia
- Fee Type: Flat
- Fee Amount: Rp 10,000 (International SWIFT)
- Min/Max: Rp 0 - Rp 999,999,999
- Effective: Now - NULL
```

#### Key Features:
- ✅ **Simplified:** Hanya satu entry per destination bank
- ✅ **Dynamic:** Can be updated via admin API
- ✅ **Real-time:** Calculated before transaction
- ✅ **Single perspective:** OUR bank sending to OTHER banks
- ✅ **No reverse:** Inbound transfers tidak ada fee (inbound is free)

---

### B. Service Fees (Biaya Layanan Lainnya)

#### 1. E-Wallet Top-Up
| Service | Code | Fee | Type | Min | Max |
|---------|------|-----|------|-----|-----|
| OVO Top-Up | TOPUP_OVO | Rp 2,500 | Flat | 0 | 10M |
| DANA Top-Up | TOPUP_DANA | Rp 2,500 | Flat | 0 | 10M |
| GoPay Top-Up | TOPUP_GOPAY | Rp 2,500 | Flat | 0 | 10M |

#### 2. E-Money Top-Up
| Service | Code | Fee | Type | Min | Max |
|---------|------|-----|------|-----|-----|
| LinkAja | TOPUP_LINKAJA | Rp 2,500 | Flat | 0 | 10M |
| Mandiri e-Money | TOPUP_MANDIRIEMONEY | Rp 2,500 | Flat | 0 | 10M |

#### 3. Virtual Account (VA) Payment
| Service | Code | Fee | Type | Min | Max |
|---------|------|-----|------|-----|-----|
| Mandiri VA | PAYMENT_VA_MANDIRI | Rp 0 | Flat | 0 | 999M |
| BCA VA | PAYMENT_VA_BCA | Rp 0 | Flat | 0 | 999M |
| BRI VA | PAYMENT_VA_BRI | Rp 0 | Flat | 0 | 999M |

#### 4. QRIS Payment
| Service | Code | Fee | Type | Min | Max |
|---------|------|-----|------|-----|-----|
| QRIS Payment | QRIS_PAYMENT | 1% | Percentage | 0 | 100M |

---

## 🔧 Admin API Endpoints

### 1. Get All Banks

**Endpoint:** `GET /api/v1/admin/banks`

**Response:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/banks
```

**Response Body:**
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

---

### 2. Get All Transfer Fees (Destination Banks)

**Endpoint:** `GET /api/v1/admin/fees/transfer`

Shows all destination banks we support for OUTBOUND transfers and their fees.

**Response:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/fees/transfer
```

**Response Body:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "destination_bank_code": "BCA",
      "destination_bank_name": "Bank BCA",
      "fee_type": "flat",
      "fee_amount": 5000,
      "fee_percentage": 0,
      "minimum_amount": 0,
      "maximum_amount": 999999999,
      "is_active": true,
      "effective_from": "2026-03-07T00:00:00Z",
      "effective_to": null,
      "description": "BCA - SKNT Domestic Transfer",
      "notes": "Standard fee for BCA destination",
      "created_at": "2026-03-07T00:00:00Z",
      "updated_at": "2026-03-07T00:00:00Z"
    },
    {
      "id": 2,
      "destination_bank_code": "MANDIRI",
      "destination_bank_name": "Bank Mandiri",
      "fee_type": "flat",
      "fee_amount": 5000,
      "fee_percentage": 0,
      "minimum_amount": 0,
      "maximum_amount": 999999999,
      "is_active": true,
      "effective_from": "2026-03-07T00:00:00Z",
      "effective_to": null,
      "description": "MANDIRI - SKNT Domestic Transfer",
      "notes": "Standard fee for Mandiri destination",
      "created_at": "2026-03-07T00:00:00Z",
      "updated_at": "2026-03-07T00:00:00Z"
    }
  ]
}
```

---

### 3. Update Transfer Fee

**Endpoint:** `PUT /api/v1/admin/fees/transfer`

Update fee for a **specific destination bank**:

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/admin/fees/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "destination_bank_code": "BCA",
    "fee_amount": 7500,
    "fee_type": "flat"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Transfer fee for BCA updated successfully",
  "data": {
    "destination_bank_code": "BCA",
    "destination_bank_name": "Bank BCA",
    "fee_amount": 7500,
    "fee_type": "flat"
  }
}
```

---

### 4. Calculate Transfer Fee

**Endpoint:** `POST /api/v1/admin/fees/transfer/calculate`

Calculate fee for transferring to a **specific destination bank**:

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/fees/transfer/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "destination_bank_code": "BCA",
    "amount": 1000000
  }'
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "destination_bank_code": "BCA",
    "destination_bank_name": "Bank BCA",
    "transfer_amount": 1000000,
    "fee": 5000,
    "total_amount": 1005000
  }
}
```

---

### 5. Get All Service Fees

**Endpoint:** `GET /api/v1/admin/fees/services`

**Query Parameters (Optional):**
- `type`: Filter by service type (topup_ewallet, topup_emoney, payment_va, qris_payment)

**Request:**
```bash
# Get all service fees
curl -X GET http://localhost:8080/api/v1/admin/fees/services

# Get only e-wallet fees
curl -X GET "http://localhost:8080/api/v1/admin/fees/services?type=topup_ewallet"

# Get only e-money fees
curl -X GET "http://localhost:8080/api/v1/admin/fees/services?type=topup_emoney"

# Get only VA payment fees
curl -X GET "http://localhost:8080/api/v1/admin/fees/services?type=payment_va"
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "service_code": "TOPUP_OVO",
      "service_name": "Top Up OVO",
      "service_type": "topup_ewallet",
      "fee_type": "flat",
      "fee_amount": 2500,
      "fee_percentage": 0,
      "minimum_amount": 0,
      "maximum_amount": 10000000,
      "is_active": 1,
      "effective_from": "2026-03-07T00:00:00Z",
      "effective_to": null,
      "notes": "OVO e-wallet top up",
      "created_at": "2026-03-07T00:00:00Z",
      "updated_at": "2026-03-07T00:00:00Z"
    }
  ]
}
```

---

### 6. Update Service Fee

**Endpoint:** `PUT /api/v1/admin/fees/services`

**Request (Flat Fee):**
```bash
curl -X PUT http://localhost:8080/api/v1/admin/fees/services \
  -H "Content-Type: application/json" \
  -d '{
    "service_code": "TOPUP_OVO",
    "fee_amount": 3000,
    "fee_percentage": 0,
    "fee_type": "flat"
  }'
```

**Request (Percentage Fee):**
```bash
curl -X PUT http://localhost:8080/api/v1/admin/fees/services \
  -H "Content-Type: application/json" \
  -d '{
    "service_code": "QRIS_PAYMENT",
    "fee_amount": 0,
    "fee_percentage": 1.5,
    "fee_type": "percentage"
  }'
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

---

### 7. Calculate Service Fee

**Endpoint:** `POST /api/v1/admin/fees/services/calculate`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/fees/services/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "service_code": "TOPUP_OVO",
    "amount": 100000
  }'
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

**Request (QRIS - Percentage):**
```bash
curl -X POST http://localhost:8080/api/v1/admin/fees/services/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "service_code": "QRIS_PAYMENT",
    "amount": 5000000
  }'
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "service_code": "QRIS_PAYMENT",
    "service_amount": 5000000,
    "fee": 50000,
    "total_amount": 5050000
  }
}
```

---

### 8. Get Fee Statistics

**Endpoint:** `GET /api/v1/admin/fees/statistics`

**Query Parameters (Optional):**
- `type`: `transfer` (default) or `service`
- `service_type`: Filter by service type (for service type only)

**Request:**
```bash
# Transfer fee statistics
curl -X GET http://localhost:8080/api/v1/admin/fees/statistics?type=transfer

# Service fee statistics
curl -X GET http://localhost:8080/api/v1/admin/fees/statistics?type=service

# E-wallet service statistics
curl -X GET "http://localhost:8080/api/v1/admin/fees/statistics?type=service&service_type=topup_ewallet"
```

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

## 🔄 Fee Integration with User Transactions

### User Transfer Flow with Dynamic Fee

```
1. User initiates transfer:
   POST /api/v1/transfers/inter
   {
     "from_account_number": "1001234567",
     "to_account_number": "2001234567",  (Different bank)
     "amount": 1000000,
     "description": "Transfer to friend"
   }

2. System processes:
   a. Validates both accounts
   b. Determines destination bank from to_account_number (IMPORTANT: SEE NOTES BELOW)
   c. Calls CalculateTransferFee(destination_bank_code, amount)
   d. Retrieves fee from database: Rp 5,000
   e. Checks balance: 1,000,000 + 5,000 = 1,005,000 ✓
   f. Debits account: 1,005,000 (transfer + fee)
   g. Creates transaction with fee recorded

3. Response to user:
   {
     "status": "success",
     "data": {
       "transaction_id": "TRX20260307ABC",
       "amount": 1000000,
       "fee": 5000,
       "total_debit": 1005000,
       "status": "success"
     }
   }

4. Statement shows:
   - Transaction Amount: Rp 1,000,000
   - Fee (Admin): Rp 5,000
   - Total: Rp 1,005,000
```

---

## ⚠️ IMPLEMENTATION NOTES (IMPORTANT)

### PENDING: Destination Bank Detection

**Current Status:**  
In `services/transfer_service.go`, the `ProcessInterBankTransfer()` function currently has:
```go
destinationBankCode := "MANDIRI" // TODO: Should be determined from to_account_number
```

**Problem:**  
All transfers will calculate fee for MANDIRI regardless of actual destination bank.

**Solution Options:**

**Option 1: Parse Account Number**
Map account number prefixes to bank codes:
```go
func DetermineBankFromAccountNumber(accountNumber string) string {
    prefix := strings.ToUpper(accountNumber[:4])
    bankMap := map[string]string{
        "1001": "BCA",
        "1009": "BNI",
        "0023": "MANDIRI",
        "0062": "PERMATA",
        // ... etc
    }
    return bankMap[prefix]
}
```

**Option 2: Accept Bank Code in Request (Recommended)**
Modify transfer request to include bank code:
```go
type InterBankTransferRequest struct {
    FromAccountNumber    string  `json:"from_account_number"`
    ToAccountNumber      string  `json:"to_account_number"`
    DestinationBankCode  string  `json:"destination_bank_code"` // ← ADD THIS
    Amount               float64 `json:"amount"`
    Description          string  `json:"description"`
}
```

**Option 3: Support Both**
Accept bank code if provided, else infer from account number.

---

## 💡 Use Cases & Examples

### Scenario 1: Change Outbound Fee Policy
Company decides to charge more for transfers to OCBC bank:

```bash
# Current fee: Rp 5,000
# New fee: Rp 10,000 (because OCBC is considered premium)

curl -X PUT http://localhost:8080/api/v1/admin/fees/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "destination_bank_code": "OCBC",
    "fee_amount": 10000,
    "fee_type": "flat"
  }'
```

**Effective immediately** - next transfer to OCBC from our bank will charge Rp 10,000

---

### Scenario 2: Promotional E-Wallet Fee
Company wants to promote OVO top-up with temporary fee reduction:

```bash
# Current: Rp 2,500
# Promo: Rp 1,000 (50% off)

curl -X PUT http://localhost:8080/api/v1/admin/fees/services \
  -H "Content-Type: application/json" \
  -d '{
    "service_code": "TOPUP_OVO",
    "fee_amount": 1000,
    "fee_percentage": 0,
    "fee_type": "flat"
  }'
```

---

### Scenario 3: Percentage-Based Fee for Large Transfers
System supports percentage-based fees for QRIS:

```bash
# QRIS: 1% fee on amount
# User transfers: Rp 10,000,000
# Fee calculated: 10,000,000 × 1% = Rp 100,000

curl -X POST http://localhost:8080/api/v1/admin/fees/services/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "service_code": "QRIS_PAYMENT",
    "amount": 10000000
  }'

Response: { "fee": 100000, "total_amount": 10100000 }
```

---

## 🎯 Recommended Initial Fee Configuration

### For Production:

**Interbank Transfer Fees:**
- **Premium Banks** (BCA, Mandiri): Rp 5,000
- **Standard Banks** (BRI, CIMB, Danamon): Rp 5,000
- **Regional Banks**: Rp 7,500
- **International Banks** (CITI, HSBC, JP Morgan): Rp 10,000

**E-Wallet Top-Up:**
- **OVO, DANA, GoPay**: Rp 2,500 each

**E-Money Top-Up:**
- **LinkAja, Mandiri e-Money**: Rp 2,500 each

**VA Payment:**
- **All**: Free (Rp 0)

**QRIS Payment:**
- **All**: 0.5% to 1% commission

---

## 📊 Database Schema

### Table: banks
```sql
CREATE TABLE banks (
  id INTEGER PRIMARY KEY,
  bank_code VARCHAR(10) UNIQUE NOT NULL,
  bank_name VARCHAR(100) NOT NULL,
  swift_code VARCHAR(11),
  is_active BOOLEAN DEFAULT 1,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Table: transfer_fees
```sql
CREATE TABLE transfer_fees (
  id INTEGER PRIMARY KEY,
  from_bank_code VARCHAR(10) NOT NULL,
  to_bank_code VARCHAR(10) NOT NULL,
  transaction_type VARCHAR(50),
  fee_type VARCHAR(20), -- 'flat' or 'percentage'
  fee_amount DECIMAL(15,2),
  fee_percentage DECIMAL(5,2),
  minimum_amount DECIMAL(15,2),
  maximum_amount DECIMAL(15,2),
  is_active BOOLEAN DEFAULT 1,
  effective_from TIMESTAMP,
  effective_to TIMESTAMP,
  notes VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Table: service_fees
```sql
CREATE TABLE service_fees (
  id INTEGER PRIMARY KEY,
  service_code VARCHAR(50) UNIQUE NOT NULL,
  service_name VARCHAR(100),
  service_type VARCHAR(50), -- 'topup_ewallet', 'topup_emoney', etc
  fee_type VARCHAR(20), -- 'flat' or 'percentage'
  fee_amount DECIMAL(15,2),
  fee_percentage DECIMAL(5,2),
  minimum_amount DECIMAL(15,2),
  maximum_amount DECIMAL(15,2),
  is_active BOOLEAN DEFAULT 1,
  effective_from TIMESTAMP,
  effective_to TIMESTAMP,
  notes VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

---

## 🔐 Security Considerations

> ⚠️ **IMPORTANT**: Admin endpoints are currently unprotected for testing.

**Before production deployment:**
1. Add authentication middleware to `/api/v1/admin` routes
2. Implement role-based access control (RBAC)
3. Log all fee changes with admin user audit trail
4. Require approval for fee changes above threshold
5. Enable rate limiting on admin endpoints

---

## 📈 Future Enhancements

- [ ] Time-based fee scheduling (fees change at specific dates)
- [ ] Regional/zone-based fees (different for Java, Sumatra, etc)
- [ ] Volume-based discounts (lower fees for high transaction volumes)
- [ ] Customer segment-based fees (Premium/Standard/Basic tiers)
- [ ] Fee approval workflow (manager approval required)
- [ ] Real-time fee change notifications
- [ ] Fee analytics dashboard
- [ ] A/B testing for different fee structures

---

**Last Updated:** 2026-03-07  
**Version:** 1.0  
**Status:** Complete & Ready
