-- CBS Simulator Database Schema
-- SQLite Database for Development

-- Customers Table
CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cif VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    id_card_number VARCHAR(20) UNIQUE NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    date_of_birth DATE,
    pin VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Accounts Table
CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    cif VARCHAR(20) NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    currency VARCHAR(3) DEFAULT 'IDR',
    balance DECIMAL(18,2) DEFAULT 0.00,
    avail_balance DECIMAL(18,2) DEFAULT 0.00,
    status VARCHAR(20) DEFAULT 'active',
    opened_date DATE NOT NULL,
    branch VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif)
);

-- Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    transaction_type VARCHAR(30) NOT NULL,
    from_account_number VARCHAR(20),
    to_account_number VARCHAR(20),
    amount DECIMAL(18,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'IDR',
    description TEXT,
    reference_number VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    transaction_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    settlement_date DATE,
    fee DECIMAL(10,2) DEFAULT 0.00,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Cards Table
CREATE TABLE IF NOT EXISTS cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    card_number VARCHAR(16) UNIQUE NOT NULL,
    cif VARCHAR(20) NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    card_type VARCHAR(20) NOT NULL,
    card_brand VARCHAR(20) NOT NULL,
    card_limit DECIMAL(18,2) DEFAULT 0.00,
    avail_limit DECIMAL(18,2) DEFAULT 0.00,
    expiry_date VARCHAR(7) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    cvv VARCHAR(3) NOT NULL,
    pin VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif),
    FOREIGN KEY (account_number) REFERENCES accounts(account_number)
);

-- Loans Table
CREATE TABLE IF NOT EXISTS loans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    loan_number VARCHAR(20) UNIQUE NOT NULL,
    cif VARCHAR(20) NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    loan_type VARCHAR(30) NOT NULL,
    principal_amount DECIMAL(18,2) NOT NULL,
    outstanding_amount DECIMAL(18,2) NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,
    monthly_payment DECIMAL(18,2) NOT NULL,
    tenor_months INTEGER NOT NULL,
    remaining_months INTEGER NOT NULL,
    disbursement_date DATE NOT NULL,
    maturity_date DATE NOT NULL,
    next_payment_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif),
    FOREIGN KEY (account_number) REFERENCES accounts(account_number)
);

-- Deposits Table
CREATE TABLE IF NOT EXISTS deposits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deposit_number VARCHAR(20) UNIQUE NOT NULL,
    cif VARCHAR(20) NOT NULL,
    principal_amount DECIMAL(18,2) NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,
    tenor_months INTEGER NOT NULL,
    open_date DATE NOT NULL,
    maturity_date DATE NOT NULL,
    maturity_amount DECIMAL(18,2) NOT NULL,
    auto_renew BOOLEAN DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    linked_account VARCHAR(20),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif),
    FOREIGN KEY (linked_account) REFERENCES accounts(account_number)
);

-- Bill Payments Table
CREATE TABLE IF NOT EXISTS bill_payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    biller_code VARCHAR(20) NOT NULL,
    biller_name VARCHAR(100) NOT NULL,
    customer_number VARCHAR(50) NOT NULL,
    bill_number VARCHAR(50) UNIQUE NOT NULL,
    bill_amount DECIMAL(18,2) NOT NULL,
    admin_fee DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(18,2) NOT NULL,
    bill_period VARCHAR(20),
    due_date DATE,
    status VARCHAR(20) DEFAULT 'unpaid',
    transaction_id VARCHAR(50),
    payment_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create Indexes for Performance
CREATE INDEX IF NOT EXISTS idx_customers_cif ON customers(cif);
CREATE INDEX IF NOT EXISTS idx_accounts_cif ON accounts(cif);
CREATE INDEX IF NOT EXISTS idx_accounts_number ON accounts(account_number);
CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions(from_account_number);
CREATE INDEX IF NOT EXISTS idx_transactions_to ON transactions(to_account_number);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_cards_cif ON cards(cif);
CREATE INDEX IF NOT EXISTS idx_loans_cif ON loans(cif);
CREATE INDEX IF NOT EXISTS idx_deposits_cif ON deposits(cif);
