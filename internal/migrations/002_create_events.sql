CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    page VARCHAR(255) NOT NULL,
    amount NUMERIC,
    created_at TIMESTAMP DEFAULT NOW()
    );