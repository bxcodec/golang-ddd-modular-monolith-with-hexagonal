CREATE TABLE IF NOT EXISTS payment_settings (
    id VARCHAR(255) PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL,
    setting_value VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(setting_key, currency)
);

CREATE INDEX idx_payment_settings_key ON payment_settings(setting_key);
CREATE INDEX idx_payment_settings_currency ON payment_settings(currency);
CREATE INDEX idx_payment_settings_status ON payment_settings(status);
CREATE INDEX idx_payment_settings_created_at ON payment_settings(created_at);

