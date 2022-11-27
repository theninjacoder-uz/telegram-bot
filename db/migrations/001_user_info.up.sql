CREATE TYPE lang AS ENUM('UZB', 'ENG', 'RUS', 'KIR');

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT UNIQUE NOT NULL,
    phone_number VARCHAR(24),
    language lang,
    tin VARCHAR(24),
    state INT DEFAULT 0,
    is_verified BOOLEAN DEFAULT false
);