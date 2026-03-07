# CBS Simulator - Core Banking System

CBS Simulator adalah simulasi core banking yang lengkap dibangun dengan Go, Gin, dan SQLite. Cocok untuk development dan testing aplikasi mobile banking.

## Features

Core Banking:
- Account management dan inquiry saldo
- Intrabank dan interbank transfer dengan fee
- PPOB/bill payment (PLN, PDAM, Telkom, BPJS, dll)
- Card management (inquiry, block/unblock)
- Loan dan deposit inquiry

Technical:
- RESTful API dengan response format konsisten
- Authentication dengan PIN
- SQLite database dengan auto-migration
- Transaction history dengan pagination
- Docker support
- Indonesian banking format

## Prerequisites

- Go 1.21 or higher
- SQLite3
- Docker & Docker Compose (optional)

## Quick Start

Run Locally:
```bash
cp .env.example .env
go mod download
go run main.go
```

Server running at http://localhost:8080

Run with Docker:
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

Health:
- GET /health

See docs/API.md for complete documentation with request/response examples.

## API Response Format

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
