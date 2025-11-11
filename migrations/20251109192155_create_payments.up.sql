CREATE TABLE IF NOT EXISTS payment_module.payments (
    id VARCHAR(255) PRIMARY KEY,
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_status ON payment_module.payments(status);
CREATE INDEX idx_payments_currency ON payment_module.payments(currency);
CREATE INDEX idx_payments_created_at ON payment_module.payments(created_at);
