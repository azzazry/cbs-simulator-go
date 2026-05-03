-- Security Tables - PostgreSQL

-- Token Blacklist
CREATE TABLE IF NOT EXISTS token_blacklist (
    id              BIGSERIAL PRIMARY KEY,
    token_jti       VARCHAR(50)  UNIQUE NOT NULL,
    cif             VARCHAR(20)  NOT NULL,
    expires_at      TIMESTAMPTZ  NOT NULL,
    blacklisted_at  TIMESTAMPTZ  DEFAULT NOW()
);

-- Roles Table
CREATE TABLE IF NOT EXISTS roles (
    id          BIGSERIAL PRIMARY KEY,
    role_name   VARCHAR(50)  UNIQUE NOT NULL,
    description TEXT,
    is_active   BOOLEAN      DEFAULT TRUE,
    created_at  TIMESTAMPTZ  DEFAULT NOW()
);

-- User Roles
CREATE TABLE IF NOT EXISTS user_roles (
    id          BIGSERIAL PRIMARY KEY,
    cif         VARCHAR(20)  NOT NULL,
    role_id     INTEGER      NOT NULL,
    assigned_by VARCHAR(20),
    assigned_at TIMESTAMPTZ  DEFAULT NOW(),
    FOREIGN KEY (cif) REFERENCES customers(cif),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    UNIQUE(cif, role_id)
);

-- Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    cif             VARCHAR(20),
    action          VARCHAR(100) NOT NULL,
    resource        VARCHAR(100),
    resource_id     VARCHAR(100),
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    request_method  VARCHAR(10),
    request_path    VARCHAR(255),
    request_body    TEXT,
    response_status INTEGER,
    details         TEXT,
    created_at      TIMESTAMPTZ  DEFAULT NOW()
);

-- Transaction Limits
CREATE TABLE IF NOT EXISTS transaction_limits (
    id                   BIGSERIAL PRIMARY KEY,
    role_name            VARCHAR(50)    NOT NULL,
    transaction_type     VARCHAR(50)    NOT NULL,
    daily_limit          DECIMAL(18, 2) DEFAULT 0,
    per_transaction_limit DECIMAL(18, 2) DEFAULT 0,
    monthly_limit        DECIMAL(18, 2) DEFAULT 0,
    is_active            BOOLEAN        DEFAULT TRUE,
    created_at           TIMESTAMPTZ    DEFAULT NOW(),
    updated_at           TIMESTAMPTZ    DEFAULT NOW(),
    UNIQUE(role_name, transaction_type)
);

-- Login Attempts
CREATE TABLE IF NOT EXISTS login_attempts (
    id           BIGSERIAL PRIMARY KEY,
    cif          VARCHAR(20)  NOT NULL,
    ip_address   VARCHAR(45),
    attempt_type VARCHAR(30)  DEFAULT 'pin',
    is_success   BOOLEAN      DEFAULT FALSE,
    attempted_at TIMESTAMPTZ  DEFAULT NOW()
);

-- OTP Codes
CREATE TABLE IF NOT EXISTS otp_codes (
    id         BIGSERIAL PRIMARY KEY,
    cif        VARCHAR(20) NOT NULL,
    otp_code   VARCHAR(10) NOT NULL,
    otp_type   VARCHAR(30) NOT NULL,
    channel    VARCHAR(20) DEFAULT 'sms',
    is_used    BOOLEAN     DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- e-KYC Verifications
CREATE TABLE IF NOT EXISTS ekyc_verifications (
    id                  BIGSERIAL PRIMARY KEY,
    verification_id     VARCHAR(50)  UNIQUE NOT NULL,
    cif                 VARCHAR(20)  NOT NULL,
    id_card_number      VARCHAR(20)  NOT NULL,
    verification_type   VARCHAR(30)  NOT NULL,
    verification_status VARCHAR(20)  DEFAULT 'pending',
    verified_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_token_blacklist_jti    ON token_blacklist(token_jti);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires ON token_blacklist(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_roles_cif         ON user_roles(cif);
CREATE INDEX IF NOT EXISTS idx_audit_logs_cif         ON audit_logs(cif);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created     ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action      ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_login_attempts_cif     ON login_attempts(cif);
CREATE INDEX IF NOT EXISTS idx_login_attempts_time    ON login_attempts(attempted_at);
CREATE INDEX IF NOT EXISTS idx_otp_codes_cif          ON otp_codes(cif);
CREATE INDEX IF NOT EXISTS idx_otp_codes_type         ON otp_codes(otp_type);
CREATE INDEX IF NOT EXISTS idx_ekyc_verification_id   ON ekyc_verifications(verification_id);
CREATE INDEX IF NOT EXISTS idx_ekyc_cif               ON ekyc_verifications(cif);
