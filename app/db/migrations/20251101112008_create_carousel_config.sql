-- +goose Up
CREATE TABLE IF NOT EXISTS carousel_config (
    id SERIAL PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert default row so we always have one config record
INSERT INTO carousel_config (enabled) VALUES (TRUE);

-- +goose Down
DROP TABLE IF EXISTS carousel_config;
