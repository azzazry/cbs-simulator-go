-- Phase 2: Core Banking Enhancement - PostgreSQL

-- =============================================
-- 1. GENERAL LEDGER (GL)
-- =============================================

CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id             BIGSERIAL PRIMARY KEY,
    account_code   VARCHAR(10)  UNIQUE NOT NULL,
    account_name   VARCHAR(100) NOT NULL,
    account_type   VARCHAR(20)  NOT NULL,
    parent_code    VARCHAR(10),
    level          INTEGER      DEFAULT 1,
    normal_balance VARCHAR(10)  NOT NULL,
    is_active      BOOLEAN      DEFAULT TRUE,
    created_at     TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS journal_entries (
    id             BIGSERIAL PRIMARY KEY,
    journal_number VARCHAR(30)  UNIQUE NOT NULL,
    entry_date     DATE         NOT NULL,
    description    TEXT,
    reference_type VARCHAR(30),
    reference_id   VARCHAR(50),
    posted_by      VARCHAR(20),
    status         VARCHAR(20)  DEFAULT 'posted',
    reversed_by    INTEGER,
    created_at     TIMESTAMPTZ  DEFAULT NOW(),
    FOREIGN KEY (reversed_by) REFERENCES journal_entries(id)
);

CREATE TABLE IF NOT EXISTS journal_lines (
    id            BIGSERIAL PRIMARY KEY,
    journal_id    INTEGER        NOT NULL,
    account_code  VARCHAR(10)    NOT NULL,
    debit_amount  DECIMAL(18, 2) DEFAULT 0,
    credit_amount DECIMAL(18, 2) DEFAULT 0,
    description   TEXT,
    FOREIGN KEY (journal_id) REFERENCES journal_entries(id),
    FOREIGN KEY (account_code) REFERENCES chart_of_accounts(account_code)
);

-- =============================================
-- 2. CIF ENHANCEMENT
-- =============================================

CREATE TABLE IF NOT EXISTS customer_extended (
    id                BIGSERIAL PRIMARY KEY,
    cif               VARCHAR(20)    UNIQUE NOT NULL,
    mother_maiden_name VARCHAR(100),
    nationality       VARCHAR(50)    DEFAULT 'WNI',
    occupation        VARCHAR(100),
    employer_name     VARCHAR(100),
    monthly_income    DECIMAL(18, 2),
    source_of_funds   VARCHAR(100),
    risk_profile      VARCHAR(20)    DEFAULT 'low',
    segment           VARCHAR(30)    DEFAULT 'mass',
    branch_code       VARCHAR(20),
    rm_code           VARCHAR(20),
    npwp              VARCHAR(20),
    last_kyc_date     DATE,
    next_kyc_date     DATE,
    created_at        TIMESTAMPTZ    DEFAULT NOW(),
    updated_at        TIMESTAMPTZ    DEFAULT NOW(),
    FOREIGN KEY (cif) REFERENCES customers(cif)
);

-- =============================================
-- 3. INTEREST CALCULATION ENGINE
-- =============================================

CREATE TABLE IF NOT EXISTS interest_rates (
    id           BIGSERIAL PRIMARY KEY,
    product_type VARCHAR(30)    NOT NULL,
    product_name VARCHAR(50),
    rate_type    VARCHAR(20)    NOT NULL,
    base_rate    DECIMAL(8, 4)  NOT NULL,
    min_balance  DECIMAL(18, 2) DEFAULT 0,
    max_balance  DECIMAL(18, 2),
    tenor_months INTEGER,
    effective_date DATE         NOT NULL,
    is_active    BOOLEAN        DEFAULT TRUE,
    created_at   TIMESTAMPTZ    DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS interest_accruals (
    id               BIGSERIAL PRIMARY KEY,
    account_number   VARCHAR(20)    NOT NULL,
    accrual_date     DATE           NOT NULL,
    product_type     VARCHAR(30)    NOT NULL,
    balance          DECIMAL(18, 2) NOT NULL,
    rate             DECIMAL(8, 4)  NOT NULL,
    daily_interest   DECIMAL(18, 4) NOT NULL,
    accrued_interest DECIMAL(18, 4) NOT NULL,
    is_posted        BOOLEAN        DEFAULT FALSE,
    created_at       TIMESTAMPTZ    DEFAULT NOW(),
    UNIQUE(account_number, accrual_date)
);

-- =============================================
-- 4. STANDING INSTRUCTIONS (SI)
-- =============================================

CREATE TABLE IF NOT EXISTS standing_instructions (
    id                   BIGSERIAL PRIMARY KEY,
    si_number            VARCHAR(30)    UNIQUE NOT NULL,
    cif                  VARCHAR(20)    NOT NULL,
    from_account         VARCHAR(20)    NOT NULL,
    instruction_type     VARCHAR(30)    NOT NULL,
    to_account           VARCHAR(20),
    to_bank_code         VARCHAR(20),
    amount               DECIMAL(18, 2) NOT NULL,
    description          TEXT,
    frequency            VARCHAR(20)    NOT NULL,
    execution_day        INTEGER,
    start_date           DATE           NOT NULL,
    end_date             DATE,
    next_execution_date  DATE           NOT NULL,
    total_executed       INTEGER        DEFAULT 0,
    total_failed         INTEGER        DEFAULT 0,
    last_execution_date  DATE,
    last_status          VARCHAR(20),
    status               VARCHAR(20)    DEFAULT 'active',
    created_at           TIMESTAMPTZ    DEFAULT NOW(),
    FOREIGN KEY (cif) REFERENCES customers(cif)
);

CREATE TABLE IF NOT EXISTS si_executions (
    id             BIGSERIAL PRIMARY KEY,
    si_number      VARCHAR(30)    NOT NULL,
    execution_date DATE           NOT NULL,
    amount         DECIMAL(18, 2) NOT NULL,
    transaction_id VARCHAR(50),
    status         VARCHAR(20)    NOT NULL,
    error_message  TEXT,
    created_at     TIMESTAMPTZ    DEFAULT NOW(),
    FOREIGN KEY (si_number) REFERENCES standing_instructions(si_number)
);

-- =============================================
-- 5. END OF DAY (EOD) PROCESSING
-- =============================================

CREATE TABLE IF NOT EXISTS eod_logs (
    id                BIGSERIAL PRIMARY KEY,
    process_date      DATE        NOT NULL,
    process_type      VARCHAR(30) NOT NULL,
    status            VARCHAR(20) NOT NULL,
    records_processed INTEGER     DEFAULT 0,
    records_failed    INTEGER     DEFAULT 0,
    started_at        TIMESTAMPTZ,
    completed_at      TIMESTAMPTZ,
    error_message     TEXT,
    created_at        TIMESTAMPTZ DEFAULT NOW()
);

-- =============================================
-- INDEXES
-- =============================================
CREATE INDEX IF NOT EXISTS idx_coa_type               ON chart_of_accounts(account_type);
CREATE INDEX IF NOT EXISTS idx_coa_parent             ON chart_of_accounts(parent_code);
CREATE INDEX IF NOT EXISTS idx_journal_date           ON journal_entries(entry_date);
CREATE INDEX IF NOT EXISTS idx_journal_ref            ON journal_entries(reference_type, reference_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_journal  ON journal_lines(journal_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_account  ON journal_lines(account_code);
CREATE INDEX IF NOT EXISTS idx_customer_ext_cif       ON customer_extended(cif);
CREATE INDEX IF NOT EXISTS idx_interest_rates_type    ON interest_rates(product_type);
CREATE INDEX IF NOT EXISTS idx_interest_accruals_account ON interest_accruals(account_number);
CREATE INDEX IF NOT EXISTS idx_interest_accruals_date ON interest_accruals(accrual_date);
CREATE INDEX IF NOT EXISTS idx_si_cif                 ON standing_instructions(cif);
CREATE INDEX IF NOT EXISTS idx_si_next_exec           ON standing_instructions(next_execution_date);
CREATE INDEX IF NOT EXISTS idx_si_exec_number         ON si_executions(si_number);
CREATE INDEX IF NOT EXISTS idx_eod_date               ON eod_logs(process_date);
