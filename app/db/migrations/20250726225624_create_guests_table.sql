-- +goose Up
CREATE TABLE guest (
	id SERIAL PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	cin TEXT,
	created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS guest;

