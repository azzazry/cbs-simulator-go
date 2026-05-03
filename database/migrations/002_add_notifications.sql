-- Notifications table untuk menyimpan transaction notifications
CREATE TABLE IF NOT EXISTS notifications (
    id                BIGSERIAL PRIMARY KEY,
    cif               TEXT        NOT NULL,
    notification_type TEXT        NOT NULL, -- 'transfer', 'payment', 'deposit', 'loan', etc
    title             TEXT        NOT NULL,
    message           TEXT        NOT NULL,
    transaction_id    TEXT,
    is_read           BOOLEAN     DEFAULT FALSE,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (cif) REFERENCES customers(cif)
);

-- FCM Device Tokens untuk push notifications
CREATE TABLE IF NOT EXISTS fcm_tokens (
    id           BIGSERIAL PRIMARY KEY,
    cif          TEXT        NOT NULL,
    device_token TEXT        NOT NULL,
    device_type  TEXT,                -- 'android', 'ios', 'web'
    device_name  TEXT,
    is_active    BOOLEAN     DEFAULT TRUE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (cif) REFERENCES customers(cif),
    UNIQUE(device_token)
);

-- Notification preferences untuk user setting
CREATE TABLE IF NOT EXISTS notification_preferences (
    id                    BIGSERIAL PRIMARY KEY,
    cif                   TEXT        NOT NULL,
    transfer_notification BOOLEAN     DEFAULT TRUE,
    payment_notification  BOOLEAN     DEFAULT TRUE,
    deposit_notification  BOOLEAN     DEFAULT TRUE,
    loan_notification     BOOLEAN     DEFAULT TRUE,
    promotion_notification BOOLEAN    DEFAULT TRUE,
    updated_at            TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (cif) REFERENCES customers(cif),
    UNIQUE(cif)
);

CREATE INDEX IF NOT EXISTS idx_notifications_cif     ON notifications(cif);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_fcm_tokens_cif        ON fcm_tokens(cif);
CREATE INDEX IF NOT EXISTS idx_fcm_tokens_active     ON fcm_tokens(is_active);
