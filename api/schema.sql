DROP DATABASE IF EXISTS cryptstax_db CASCADE;
CREATE DATABASE IF NOT EXISTS cryptstax_db;
SET DATABASE = cryptstax_db;

CREATE TABLE IF NOT EXISTS channels (
    id INT PRIMARY KEY DEFAULT unique_rowid(),
    address VARCHAR NOT NULL,
    created VARCHAR NOT NULL,
    creator VARCHAR NOT NULL,
    is_archived BOOLEAN NOT NULL,
    is_channel BOOLEAN NOT NULL,
    is_general BOOLEAN NOT NULL,
    is_member BOOLEAN NOT NULL,
    is_mpim BOOLEAN NOT NULL,
    is_org_shared BOOLEAN NOT NULL,
    is_private BOOLEAN NOT NULL,
    is_shared BOOLEAN NOT NULL
);