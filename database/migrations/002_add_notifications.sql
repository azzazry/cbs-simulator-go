-- Notifications table untuk menyimpan transaction notifications
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cif TEXT NOT NULL,
    notification_type TEXT NOT NULL, -- 'transfer', 'payment', 'deposit', 'loan', etc
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    transaction_id TEXT,
    is_read INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif)
);

-- FCM Device Tokens untuk push notifications
CREATE TABLE IF NOT EXISTS fcm_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cif TEXT NOT NULL,
    device_token TEXT NOT NULL,
    device_type TEXT, -- 'android', 'ios', 'web'
    device_name TEXT,
    is_active INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif),
    UNIQUE(device_token)
);

-- Notification preferences untuk user setting
CREATE TABLE IF NOT EXISTS notification_preferences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cif TEXT NOT NULL,
    transfer_notification INTEGER DEFAULT 1,
    payment_notification INTEGER DEFAULT 1,
    deposit_notification INTEGER DEFAULT 1,
    loan_notification INTEGER DEFAULT 1,
    promotion_notification INTEGER DEFAULT 1,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cif) REFERENCES customers(cif),
    UNIQUE(cif)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_notifications_cif ON notifications(cif);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_fcm_tokens_cif ON fcm_tokens(cif);
CREATE INDEX IF NOT EXISTS idx_fcm_tokens_active ON fcm_tokens(is_active);
