CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    coins INTEGER NOT NULL DEFAULT 1000,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
