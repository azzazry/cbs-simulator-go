# CBS Simulator - Security API Documentation

## Overview

Phase 1 Security menambahkan fitur autentikasi JWT, RBAC, audit trail, rate limiting, PIN policy, OTP, dan e-KYC ke CBS Simulator. Dokumen ini menjelaskan semua endpoint keamanan yang tersedia.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Semua **protected endpoint** membutuhkan JWT token di header:

```
Authorization: Bearer <access_token>
```

Token didapat dari response login. Access token berlaku 15 menit, refresh token 7 hari.

---

## Response Format

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

## Test Credentials

| CIF | Nama | PIN | Role |
|-----|------|-----|------|
| `CIF001` | Budi Santoso | `123456` | admin, customer |
| `CIF002` | Siti Nurhaliza | `123456` | customer |
| `CIF003` | Ahmad Wijaya | `123456` | supervisor, customer |
| `CIF004` | Dewi Lestari | `123456` | customer |
| `CIF005` | Rizki Pratama | `123456` | customer |

**ID Card (KTP) Numbers:**
- CIF001: `3201011234567890`
- CIF002: `3201021234567891`
- CIF003: `3201031234567892`
- CIF004: `3201041234567893`
- CIF005: `3201051234567894`

---

## 🔐 Auth Endpoints (Public)

### Login
Authenticate customer dan dapatkan JWT token pair.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "cif": "CIF001",
  "pin": "123456"
}
```

**Success Response (HTTP 200):**
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
- 401: Account is locked (3x failed login)
- 401: Account is suspended/inactive

**Lockout Behavior:**
- 3x login gagal → account locked
- Response akan menunjukkan sisa attempt: `"invalid CIF or PIN. 2 attempts remaining"`
- Setelah locked: `"account is locked due to 3 failed login attempts. Please unlock via e-KYC verification"`

---

### Register
Buat akun nasabah baru dengan PIN policy enforcement.

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

**PIN Policy Rules:**
- Harus 6 digit angka
- Tidak boleh berurutan (contoh: `123456`, `654321`)
- Tidak boleh angka sama semua (contoh: `111111`, `222222`)

**Error Responses:**
- 400: CIF already exists
- 400: PIN policy violation

---

## 🔓 Self-Service Unlock Flow (Public)

Flow untuk nasabah yang akunnya terkunci setelah 3x gagal login. Tidak membutuhkan JWT.

### Step 1: Verifikasi e-KYC (KTP)
Verifikasi identitas nasabah menggunakan nomor KTP.

**Endpoint:** `POST /auth/ekyc/verify`

**Request Body:**
```json
{
  "cif": "CIF002",
  "id_card_number": "3201021234567891"
}
```

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "verification_id": "abc123-def456-...",
    "status": "verified",
    "customer_name": "Siti Nurhaliza",
    "message": "e-KYC verification successful"
  }
}
```

> **Penting:** Simpan `verification_id` untuk Step 3.

**Error Responses:**
- 400: ID card number does not match
- 400: Customer not found

---

### Step 2: Request OTP
Generate OTP dan kirim ke nasabah (di simulator, OTP dikembalikan langsung di response).

**Endpoint:** `POST /auth/otp/request`

**Request Body:**
```json
{
  "cif": "CIF002",
  "otp_type": "unlock_account",
  "channel": "sms"
}
```

**Parameters:**
- `otp_type`: `unlock_account` atau `reset_pin`
- `channel`: `sms` atau `email` (opsional, default: `sms`)

**Success Response:**
```json
{
  "status": "success",
  "message": "OTP sent via sms",
  "data": {
    "otp_type": "unlock_account",
    "channel": "sms",
    "otp": "482931"
  }
}
```

> **Note:** Field `otp` hanya muncul di mode simulator. Di production, OTP dikirim via SMS/Email saja.

---

### Step 3a: Unlock Account
Buka kunci akun menggunakan OTP dan verification_id dari e-KYC.

**Endpoint:** `POST /auth/unlock`

**Request Body:**
```json
{
  "cif": "CIF002",
  "otp": "482931",
  "verification_id": "abc123-def456-..."
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Account unlocked successfully. You can now login with your PIN."
}
```

**Error Responses:**
- 400: OTP verification failed (expired, incorrect, already used)
- 400: e-KYC validation failed

---

### Step 3b: Reset PIN
Reset PIN tanpa perlu PIN lama (menggunakan e-KYC + OTP).

**Endpoint:** `POST /auth/reset-pin`

**Request Body:**
```json
{
  "cif": "CIF002",
  "new_pin": "975312",
  "verification_id": "abc123-def456-..."
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "PIN reset successfully. You can now login with your new PIN."
}
```

**Error Responses:**
- 400: PIN policy violation
- 400: e-KYC validation failed

---

### Verify OTP (Standalone)
Verifikasi OTP secara terpisah.

**Endpoint:** `POST /auth/otp/verify`

**Request Body:**
```json
{
  "cif": "CIF002",
  "otp": "482931",
  "otp_type": "unlock_account"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "OTP verified successfully"
}
```

---

## 🔑 Auth Endpoints (Protected)

Semua endpoint di bawah membutuhkan `Authorization: Bearer <access_token>`.

### Logout
Invalidasi token saat ini (blacklisting).

**Endpoint:** `POST /auth/logout`

**Headers:** `Authorization: Bearer <access_token>`

**Success Response:**
```json
{
  "status": "success",
  "message": "Logged out successfully"
}
```

> Setelah logout, token yang sama tidak bisa digunakan lagi (di-blacklist).

---

### Refresh Token
Generate token pair baru dari refresh token yang valid.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

**Error Responses:**
- 401: Invalid or expired refresh token
- 401: Token has been revoked

---

### Get Profile
Mendapatkan profil nasabah yang sedang login beserta roles.

**Endpoint:** `GET /auth/profile`

**Headers:** `Authorization: Bearer <access_token>`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "customer": {
      "id": 1,
      "cif": "CIF001",
      "full_name": "Budi Santoso",
      "id_card_number": "3201011234567890",
      "phone_number": "081234567890",
      "email": "budi.santoso@email.com",
      "address": "Jl. Sudirman No. 123, Jakarta Selatan",
      "date_of_birth": "1985-03-15",
      "status": "active"
    },
    "roles": [
      {
        "role_name": "admin",
        "description": "System administrator with full access"
      },
      {
        "role_name": "customer",
        "description": "Regular bank customer with standard access"
      }
    ]
  }
}
```

---

### Change PIN
Ganti PIN dengan validasi PIN lama dan PIN policy.

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

**Success Response:**
```json
{
  "status": "success",
  "message": "PIN changed successfully"
}
```

**Error Responses:**
- 400: Incorrect old PIN
- 400: PIN policy violation

---

## 🛡️ Admin Endpoints

Membutuhkan JWT token dari user dengan role **admin** atau **supervisor**.

### Get Audit Logs
Lihat log aktivitas sistem dengan filter dan paginasi.

**Endpoint:** `GET /admin/audit-logs`

**Headers:** `Authorization: Bearer <admin_token>`

**Query Parameters:**
- `cif` (optional): Filter berdasarkan CIF nasabah
- `action` (optional): Filter berdasarkan jenis aksi
- `page` (optional): Halaman (default: 1)
- `page_size` (optional): Ukuran per halaman (default: 20)

**Example:** `GET /admin/audit-logs?cif=CIF002&page=1&page_size=10`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "audit_logs": [
      {
        "id": 1,
        "cif": "CIF002",
        "action": "POST /api/v1/auth/login",
        "resource": "/api/v1/auth/login",
        "ip_address": "127.0.0.1",
        "user_agent": "curl/8.0",
        "request_method": "POST",
        "request_path": "/api/v1/auth/login",
        "response_status": 200,
        "created_at": "2026-03-07T10:30:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "page_size": 10
  }
}
```

**Error Response:**
- 403: Insufficient permissions (role bukan admin/supervisor)

---

### Get Transaction Limits
Lihat semua transaction limit atau filter berdasarkan role.

**Endpoint:** `GET /admin/transaction-limits`

**Headers:** `Authorization: Bearer <admin_token>`

**Query Parameters:**
- `role` (optional): Filter berdasarkan role name (`customer`, `teller`)

**Example:** `GET /admin/transaction-limits?role=customer`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "role_name": "customer",
      "transaction_type": "intra_transfer",
      "daily_limit": 50000000,
      "per_transaction_limit": 25000000,
      "monthly_limit": 500000000
    },
    {
      "id": 2,
      "role_name": "customer",
      "transaction_type": "inter_transfer",
      "daily_limit": 25000000,
      "per_transaction_limit": 10000000,
      "monthly_limit": 250000000
    }
  ]
}
```

**Default Limits (Customer):**

| Jenis Transaksi | Per Transaksi | Harian | Bulanan |
|-----------------|--------------|--------|---------|
| Intra Transfer | Rp 25 juta | Rp 50 juta | Rp 500 juta |
| Inter Transfer | Rp 10 juta | Rp 25 juta | Rp 250 juta |
| Bill Payment | Rp 5 juta | Rp 10 juta | Rp 100 juta |
| QRIS Payment | Rp 2 juta | Rp 5 juta | Rp 50 juta |
| VA Payment | Rp 25 juta | Rp 50 juta | Rp 500 juta |
| E-Wallet Topup | Rp 1 juta | Rp 2 juta | Rp 20 juta |
| E-Money Topup | Rp 1 juta | Rp 2 juta | Rp 20 juta |

> Admin dan Supervisor tidak memiliki limit (bypass).

---

### Update Transaction Limit
Update limit transaksi tertentu.

**Endpoint:** `PUT /admin/transaction-limits`

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "id": 1,
  "daily_limit": 100000000,
  "per_transaction_limit": 50000000,
  "monthly_limit": 1000000000
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Transaction limit updated"
}
```

---

### Get Roles
Lihat semua role yang tersedia, atau role user tertentu.

**Endpoint:** `GET /admin/roles`

**Headers:** `Authorization: Bearer <admin_token>`

**Query Parameters:**
- `cif` (optional): Lihat role user tertentu

**Example 1 (semua role):** `GET /admin/roles`

**Success Response:**
```json
{
  "status": "success",
  "data": [
    { "id": 1, "role_name": "customer", "description": "Regular bank customer with standard access" },
    { "id": 2, "role_name": "teller", "description": "Bank teller with counter operation access" },
    { "id": 3, "role_name": "admin", "description": "System administrator with full access" },
    { "id": 4, "role_name": "supervisor", "description": "Branch supervisor with approval authority" }
  ]
}
```

**Example 2 (role user):** `GET /admin/roles?cif=CIF001`

**Success Response:**
```json
{
  "status": "success",
  "data": {
    "cif": "CIF001",
    "roles": [
      { "role_name": "admin", "description": "System administrator with full access" },
      { "role_name": "customer", "description": "Regular bank customer with standard access" }
    ]
  }
}
```

---

### Assign Role
Tambahkan role ke user tertentu.

**Endpoint:** `POST /admin/roles/assign`

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "cif": "CIF002",
  "role_name": "teller"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Role teller assigned to CIF002"
}
```

**Available Roles:** `customer`, `teller`, `admin`, `supervisor`

---

### Admin Force-Unlock Account
Buka kunci akun secara paksa (bypass e-KYC + OTP).

**Endpoint:** `POST /admin/unlock-account`

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "cif": "CIF002"
}
```

**Success Response:**
```json
{
  "status": "success",
  "message": "Account CIF002 unlocked by admin"
}
```

---

## 🔒 Middleware & Security

### Rate Limiting

Semua endpoint dilindungi rate limiter: **60 request / menit per IP** (configurable via `RATE_LIMIT_PER_MINUTE`).

**Error Response (HTTP 429):**
```json
{
  "status": "error",
  "message": "Rate limit exceeded. Too many requests.",
  "retry_after": 45
}
```

### JWT Token Lifecycle

```
Login → Access Token (15 min) + Refresh Token (7 hari)
  ↓
Access Token expired → POST /auth/refresh → New Token Pair
  ↓
Logout → Token di-blacklist (tidak bisa dipakai lagi)
```

### Protected Endpoint Errors

| HTTP Status | Condition |
|-------------|-----------|
| 401 | Token missing, invalid, expired, atau revoked |
| 403 | Role tidak memenuhi (RBAC) |
| 429 | Rate limit exceeded |

**401 Response:**
```json
{
  "status": "error",
  "message": "Authorization header is required"
}
```

**403 Response:**
```json
{
  "status": "error",
  "message": "Insufficient permissions. Required role: admin or supervisor"
}
```

### Audit Trail

Semua request ke protected endpoint dicatat otomatis ke audit log:
- CIF user yang melakukan aksi
- HTTP method + path
- IP address + User-Agent
- Request body (sensitive fields di-mask)
- Response status code
- Timestamp

---

## 🔧 Environment Configuration

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `JWT_SECRET` | - | Secret key untuk sign JWT (wajib diubah di production) |
| `JWT_ACCESS_EXPIRY_MINUTES` | `15` | Masa berlaku access token (menit) |
| `JWT_REFRESH_EXPIRY_HOURS` | `168` | Masa berlaku refresh token (jam, 7 hari) |
| `RATE_LIMIT_PER_MINUTE` | `60` | Maks request per menit per IP |
| `MAX_LOGIN_ATTEMPTS` | `3` | Maks login gagal sebelum locked |
| `LOCKOUT_DURATION_MINUTES` | `30` | Durasi lockout (menit) |
| `PIN_MIN_LENGTH` | `6` | Panjang minimum PIN |
| `PIN_MAX_LENGTH` | `6` | Panjang maksimum PIN |
| `OTP_EXPIRY_MINUTES` | `5` | Masa berlaku OTP (menit) |
| `OTP_LENGTH` | `6` | Jumlah digit OTP |

---

## Route Map

### Public Routes (No Auth)

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/auth/login` | Login | Autentikasi + JWT |
| POST | `/auth/register` | Register | Registrasi nasabah |
| POST | `/auth/otp/request` | RequestOTP | Generate OTP |
| POST | `/auth/otp/verify` | VerifyOTP | Verifikasi OTP |
| POST | `/auth/ekyc/verify` | VerifyEKYC | Verifikasi KTP |
| POST | `/auth/unlock` | UnlockAccount | Self-service unlock |
| POST | `/auth/reset-pin` | ResetPIN | Reset PIN via e-KYC |

### Protected Routes (JWT Required)

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/auth/logout` | Logout | Blacklist token |
| POST | `/auth/refresh` | RefreshToken | Refresh JWT pair |
| POST | `/auth/change-pin` | ChangePIN | Ganti PIN |
| GET | `/auth/profile` | GetProfile | Profil + roles |

### Admin Routes (JWT + RBAC: admin/supervisor)

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/admin/audit-logs` | GetAuditLogs | Audit trail |
| GET | `/admin/transaction-limits` | GetTransactionLimits | Lihat limits |
| PUT | `/admin/transaction-limits` | UpdateTransactionLimit | Update limit |
| GET | `/admin/roles` | GetUserRoles | Lihat roles |
| POST | `/admin/roles/assign` | AssignRole | Tambah role |
| POST | `/admin/unlock-account` | AdminUnlockAccount | Force unlock |

---

## Database Security Tables

| Tabel | Deskripsi |
|-------|-----------|
| `token_blacklist` | JWT token yang sudah di-revoke |
| `roles` | Daftar role (customer, teller, admin, supervisor) |
| `user_roles` | Mapping CIF ↔ Role (many-to-many) |
| `audit_logs` | Log aktivitas immutable |
| `transaction_limits` | Limit transaksi per role per jenis |
| `login_attempts` | Tracking login gagal untuk lockout |
| `otp_codes` | One-Time Password records |
| `ekyc_verifications` | Catatan verifikasi e-KYC |
