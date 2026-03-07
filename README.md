# CBS Simulator - Core Banking System

CBS Simulator adalah simulasi core banking system lengkap yang dibangun dengan Go, Gin, dan SQLite. Dirancang sebagai backend untuk development dan testing aplikasi mobile banking, dengan fitur keamanan dan core banking yang menyerupai CBS bank sesungguhnya.

## Fitur

### Keamanan (Phase 1)
- **JWT Authentication** — Access + refresh token, blacklisting
- **RBAC** — 4 role: admin, supervisor, teller, customer
- **Account Lockout** — 3x gagal login = terkunci
- **Self-service Unlock** — Via e-KYC + OTP
- **PIN Policy** — 6 digit, tanpa sequential/repeating
- **Audit Trail** — Logging otomatis setiap aktivitas
- **Rate Limiting** — 60 request/menit per IP
- **Transaction Limits** — Batas per role (harian, per-transaksi, bulanan)

### Core Banking (Phase 2)
- **General Ledger** — Double-entry bookkeeping, 53 Chart of Accounts (kode 3-digit PAPI)
- **CIF Enhancement** — Single customer view, data tambahan (NPWP, pendapatan, profil risiko)
- **Mesin Bunga** — Accrual harian, tiered rates, simulasi bunga deposito/kredit
- **Standing Instructions** — Auto-debit, transfer terjadwal (daily/weekly/monthly/quarterly)
- **EOD Processing** — Batch: interest accrual → SI execution → monthly posting → dormant check
- **Account Management** — Pembukaan/penutupan rekening dengan jurnal GL

### Operasi Perbankan
- Transfer intra-bank dan inter-bank dengan fee
- PPOB/bill payment (PLN, PDAM, Telkom, BPJS)
- QRIS Payment, Virtual Account, E-Wallet, E-Money topup
- Card management (inquiry, block/unblock)
- Loan dan deposit inquiry
- Push notifications (FCM)

## Prerequisites

- Go 1.21+
- SQLite3
- Docker & Docker Compose (opsional)

## Quick Start

### Opsi 1: Pre-Built Executable (Windows)

```cmd
# Set JWT secret (wajib)
set JWT_SECRET=secret-key-kamu

# Jalankan
.\cbs-simulator.exe
```

Server berjalan di `http://localhost:8080`

### Opsi 2: Build dari Source

```bash
go build -o cbs-simulator.exe .
./cbs-simulator.exe
```

### Opsi 3: Docker

```bash
docker-compose up -d
```

## Environment Variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| `SERVER_PORT` | 8080 | Port server |
| `DATABASE_PATH` | ./database/cbs.db | Path database |
| `JWT_SECRET` | **(wajib)** | Secret key JWT |
| `ENVIRONMENT` | development | development/production |
| `MAX_LOGIN_ATTEMPTS` | 3 | Batas gagal login |
| `LOCKOUT_DURATION_MINUTES` | 30 | Durasi lockout |
| `RATE_LIMIT_PER_MINUTE` | 60 | Rate limit per IP |

## Database

Database SQLite otomatis dibuat dengan 25 tabel dan sample data.

**Sample Customers** (PIN: `123456`):

| CIF | Nama | Akun |
|-----|------|------|
| CIF001 | Budi Santoso | 1001234567, 2001234567 |
| CIF002 | Siti Nurhaliza | 1001234568, 3001234568 |
| CIF003 | Ahmad Wijaya | 1001234569, 2001234569 |
| CIF004 | Dewi Lestari | 1001234570, 4001234570 |
| CIF005 | Rizki Pratama | 1001234571, 2001234571 |

## API Endpoints

### Auth (Public)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/auth/login` | Login, return JWT |
| POST | `/api/v1/auth/register` | Registrasi nasabah baru |
| POST | `/api/v1/auth/otp/request` | Request OTP |
| POST | `/api/v1/auth/otp/verify` | Verifikasi OTP |
| POST | `/api/v1/auth/ekyc/verify` | Verifikasi e-KYC |
| POST | `/api/v1/auth/unlock` | Unlock akun (e-KYC + OTP) |
| POST | `/api/v1/auth/reset-pin` | Reset PIN via e-KYC |

### Auth (Protected — JWT Required)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/auth/logout` | Logout, revoke token |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/change-pin` | Ganti PIN |
| GET | `/api/v1/auth/profile` | Profil user |

### Banking
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/customers/:cif` | Data nasabah |
| GET | `/api/v1/accounts/:account_number` | Saldo rekening |
| GET | `/api/v1/accounts/:account_number/transactions` | Mutasi |
| POST | `/api/v1/transfers/intra` | Transfer intra-bank |
| POST | `/api/v1/transfers/inter` | Transfer inter-bank |
| POST | `/api/v1/bills/pay` | Bayar tagihan |

### General Ledger
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/gl/chart-of-accounts` | Daftar akun GL |
| GET | `/api/v1/gl/journal-entries` | Jurnal entries |
| GET | `/api/v1/gl/journal-entries/:id` | Detail jurnal |
| GET | `/api/v1/gl/trial-balance` | Neraca saldo |
| GET | `/api/v1/gl/account-balance/:code` | Saldo akun GL |

### CIF Enhancement
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/customers/:cif/overview` | Single customer view |
| PUT | `/api/v1/customers/:cif/extended` | Update data tambahan |
| GET | `/api/v1/customers/search?q=` | Cari nasabah |

### Bunga & Simulasi
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/interest/rates` | Daftar suku bunga |
| POST | `/api/v1/interest/calculate` | Simulasi bunga |

### Standing Instructions
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/standing-instructions` | Buat SI baru |
| GET | `/api/v1/standing-instructions/:cif` | Daftar SI nasabah |
| PUT | `/api/v1/standing-instructions/:si/pause` | Pause SI |
| DELETE | `/api/v1/standing-instructions/:si` | Batalkan SI |
| GET | `/api/v1/standing-instructions/:si/history` | Riwayat eksekusi |

### Account Management
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/accounts/open` | Buka rekening baru |
| POST | `/api/v1/accounts/:account_number/close` | Tutup rekening |
| GET | `/api/v1/accounts/dormant` | Daftar rekening dormant |
| POST | `/api/v1/accounts/:account_number/reactivate` | Aktifkan kembali |

### Admin (Require Role: admin/supervisor)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/admin/audit-logs` | Audit trail |
| GET | `/api/v1/admin/transaction-limits` | Batas transaksi |
| PUT | `/api/v1/admin/transaction-limits` | Update batas |
| GET | `/api/v1/admin/roles` | Daftar role |
| POST | `/api/v1/admin/roles/assign` | Assign role |
| POST | `/api/v1/admin/unlock-account` | Force unlock |
| POST | `/api/v1/admin/eod/run` | Jalankan EOD |
| GET | `/api/v1/admin/eod/status/:date` | Status EOD |
| GET | `/api/v1/admin/eod/history` | Riwayat EOD |

### Payment Channels
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/payments/qris` | QRIS Payment |
| POST | `/api/v1/payments/va` | Virtual Account |
| POST | `/api/v1/payments/ewallet/topup` | E-Wallet Top-up |
| POST | `/api/v1/payments/emoney/topup` | E-Money Top-up |

**Lihat `docs/API.md` untuk dokumentasi lengkap dengan contoh request/response.**

## Testing

```bash
# Jalankan semua 60 test
go test ./services/ -v

# Test spesifik
go test ./services/ -v -run TestCreateJournalEntry
```

Lihat `docs/TESTING.md` untuk detail lengkap.

## Project Structure

```
cbs-simulator/
├── api/
│   ├── handlers/          # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── banking_handler.go
│   │   └── core_banking_handler.go
│   ├── middleware/         # JWT, RBAC, Audit, Rate Limiter
│   └── routes/
├── config/                 # Environment config
├── database/
│   ├── migrations/         # 6 migration files (25 tabel)
│   └── seeders/            # Sample data
├── models/                 # 29 data models
├── services/               # 20 service files + tests
├── utils/                  # Hash, helpers
├── docs/                   # API.md, API_SECURITY.md, TESTING.md
├── main.go
├── Dockerfile
└── docker-compose.yml
```

## Dokumentasi

| File | Keterangan |
|------|------------|
| `docs/API.md` | Dokumentasi API lengkap |
| `docs/API_SECURITY.md` | Dokumentasi endpoint keamanan |
| `docs/TESTING.md` | Panduan testing |
| `docs/BANKING_FEES.md` | Daftar biaya transaksi |
| `docs/FLUTTER_INTEGRATION.md` | Panduan integrasi Flutter |

## Response Format

```json
// Sukses
{"status": "success", "data": {}}

// Error
{"status": "error", "message": "Deskripsi error"}

// Login sukses
{"status": "success", "data": {"access_token": "...", "refresh_token": "...", "role": "customer"}}
```

## License

MIT License
