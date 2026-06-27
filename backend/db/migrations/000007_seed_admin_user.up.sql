-- Seed default admin user
-- Password: admin123
-- BCrypt hash generated with cost 10
INSERT INTO users (name, phone_number, role, password_hash, created_at, updated_at)
VALUES (
    'Administrator',
    '08123456789',
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (phone_number) DO NOTHING;
