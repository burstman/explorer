-- +goose Up
CREATE TABLE guests (
	id SERIAL PRIMARY KEY,
	booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	cin TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS guests;
