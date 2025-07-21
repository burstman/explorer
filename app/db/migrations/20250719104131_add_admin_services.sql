-- +goose Up
CREATE TABLE services (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	price DECIMAL(10, 2) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS services;
