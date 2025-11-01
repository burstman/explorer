-- +goose Up
ALTER TABLE users
    RENAME COLUMN social_link TO social_id;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS provider TEXT;

ALTER TABLE users
    ALTER COLUMN password_hash DROP NOT NULL;


