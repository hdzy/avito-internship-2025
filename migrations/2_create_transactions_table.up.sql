-- Таблица транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id),
    type VARCHAR(20) NOT NULL,  -- transfer / purchase
    amount INTEGER NOT NULL,
    counterparty VARCHAR(50),   -- контрагент (другой сотрудник для перевода)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
