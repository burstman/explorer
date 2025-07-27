-- +goose Up
CREATE TABLE service (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	price DECIMAL(10, 2) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP

);

-- +goose Down
DROP TABLE IF EXISTS service;
