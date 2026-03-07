# CBS Simulator - Core Banking System

CBS Simulator adalah simulasi core banking yang lengkap dibangun dengan Go, Gin, dan SQLite. Cocok untuk development dan testing aplikasi mobile banking.

## Features

Core Banking:
- Account management dan inquiry saldo
- Intrabank dan interbank transfer dengan fee
- PPOB/bill payment (PLN, PDAM, Telkom, BPJS, dll)
- Card management (inquiry, block/unblock)
- Loan dan deposit inquiry
- Push notifications dengan FCM
- QRIS Payment (QR Code Indonesian Standard)
- Virtual Account (VA) Payment (Mandiri, BCA, BRI)
- E-Wallet Top-up (OVO, DANA, GoPay)
- E-Money Top-up (LinkAja, Mandiri e-Money)

Technical:
- RESTful API dengan response format konsisten
- Authentication dengan PIN
- SQLite database dengan auto-migration
- Transaction history dengan pagination
- Docker support
- Indonesian banking format
- Dynamic fee management system

## Prerequisites

- Go 1.21 or higher
- SQLite3
- Docker & Docker Compose (optional)

## Quick Start

### Option 1: Pre-Built Executable (Windows)

**cbs-simulator.exe** adalah executable yang sudah dikompilasi dan siap pakai:

```cmd
# Langsung jalankan executable
.\cbs-simulator.exe

# Server akan berjalan di http://localhost:8080
```

**Apa itu cbs-simulator.exe?**
- File binary hasil kompilasi dari Go source code
- Tidak perlu install Go, tinggal jalankan langsung
- Sempurna untuk testing & demo
- Database otomatis dibuat di `./database/cbs.db`
- Logs disimpan di `./logs/`

### Option 2: Run dengan Go

```bash
go run main.go
```

Server running at http://localhost:8080

### Option 3: Docker

```bash
docker-compose up -d
docker-compose logs -f
docker-compose down
```

## Database

Database SQLite otomatis dibuat di ./database/cbs.db dengan:
- Complete schema
- 5 sample customers
- Sample accounts, cards, loans, deposits, bills

Sample Customers (PIN: 123456):
- CIF001: Budi Santoso (accounts 1001234567, 2001234567)
- CIF002: Siti Nurhaliza (accounts 1001234568, 3001234568)
- CIF003: Ahmad Wijaya (accounts 1001234569, 2001234569)
- CIF004: Dewi Lestari (accounts 1001234570, 4001234570)
- CIF005: Rizki Pratama (accounts 1001234571, 2001234571)

## API Endpoints

Auth:
- POST /api/v1/auth/login
- POST /api/v1/auth/register
- POST /api/v1/auth/change-pin

Customer & Accounts:
- GET /api/v1/customers/:cif
- GET /api/v1/customers/:cif/accounts
- GET /api/v1/accounts/:account_number
- GET /api/v1/accounts/:account_number/statement

Transfers:
- POST /api/v1/transfers/intra
- POST /api/v1/transfers/inter
- GET /api/v1/transfers/:transaction_id

Bill Payments:
- GET /api/v1/bills/billers
- GET /api/v1/bills/inquiry
- POST /api/v1/bills/pay

Cards:
- GET /api/v1/customers/:cif/cards
- GET /api/v1/cards/:card_number
- POST /api/v1/cards/block
- POST /api/v1/cards/unblock

Loans & Deposits:
- GET /api/v1/customers/:cif/loans
- GET /api/v1/customers/:cif/deposits
- POST /api/v1/payments/qris (QRIS Payment)
- POST /api/v1/payments/va (Virtual Account Payment)
- POST /api/v1/payments/ewallet/topup (E-Wallet Top-up)
- GET /api/v1/payments/ewallet/providers
- POST /api/v1/payments/emoney/topup (E-Money Top-up)
- GET /api/v1/payments/emoney/providers
- GET /api/v1/payments/va/providers

Health:
- GET /health

See docs/API.md for complete documentation with request/response examples.

## Testing Endpoints

### Automated Test Scripts

**test_payment_features.bat** (Windows) dan **test_payment_features.sh** (Linux/Mac)

Apa itu test scripts?
- Automated testing untuk semua endpoint
- Menggunakan curl untuk send HTTP requests
- Menampilkan response JSON yang rapi dengan jq
- Mencakup 4 fitur payment utama:
  - **Provider endpoints** - List semua payment providers
  - **QRIS payment** - QR code payment testing
  - **Virtual Account** - Bank VA payment testing
  - **E-Wallet** - OVO, DANA, GoPay top-up testing
  - **E-Money** - LinkAja, Mandiri e-Money testing

**Cara menjalankan:**

Windows:
```cmd
.\test_payment_features.bat
```

Linux/Mac:
```bash
chmod +x test_payment_features.sh
./test_payment_features.sh
```

**Output:**
- ✅ Sukses: Response JSON ditampilkan
- ❌ Error: Pesan error ditampilkan dengan jelas
- ⏱️ Processing time ditampilkan di setiap test

**API Response Format

Success:
```json
{
  "status": "success",
  "data": {}
}
```

Error:
```json
{
  "status": "error",
  "message": "Error description"
}
```

## Project Structure

```
cbs-simulator/
├── api/
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Middleware
│   └── routes/         # Routes
├── config/             # Configuration
├── database/           # Database setup
│   ├── migrations/     # Schema
│   └── seeders/        # Sample data
├── models/             # Data models
├── services/           # Business logic
├── utils/              # Helpers
├── docs/               # Documentation
├── main.go
├── go.mod
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Configuration

Edit .env file:

```env
SERVER_PORT=8080
DATABASE_PATH=./database/cbs.db
JWT_SECRET=your-secret-key
ENVIRONMENT=development
```

## Database Schema

Tables:
- customers
- accounts
- transactions
- cards
- loans
- deposits
- bill_payments

See database/migrations/001_init_schema.sql for complete schema.
