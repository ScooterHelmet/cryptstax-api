DROP DATABASE IF EXISTS cryptstax_db CASCADE;
CREATE DATABASE IF NOT EXISTS cryptstax_db;
SET DATABASE = cryptstax_db;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email STRING UNIQUE,
    username STRING UNIQUE
);

CREATE TABLE IF NOT EXISTS verification_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (email, username) VALUES
    ('john@cryptstax.local', 'john_doe');