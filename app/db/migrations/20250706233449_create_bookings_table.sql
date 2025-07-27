-- +goose Up
CREATE TABLE bookings (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	camp_id INTEGER NOT NULL REFERENCES campsites(id) ON DELETE CASCADE,
	special_request TEXT,
	total_price DECIMAL(10, 2),
	status TEXT DEFAULT 'pending',
	payment_status TEXT DEFAULT 'pending', 
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS bookings;

