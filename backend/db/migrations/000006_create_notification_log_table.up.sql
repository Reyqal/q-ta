CREATE TABLE notification_log (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    channel VARCHAR(50) NOT NULL DEFAULT 'whatsapp',
    message TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'simulated_sent' CHECK (status IN ('simulated_sent', 'sent', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
