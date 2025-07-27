-- +goose Up
CREATE TABLE guests (
	id SERIAL PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	cin TEXT,
	created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS guests;

