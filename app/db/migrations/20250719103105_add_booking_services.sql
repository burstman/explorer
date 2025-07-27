-- +goose Up
CREATE TABLE booking_services (
	id SERIAL PRIMARY KEY,
	booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
	service_id INTEGER NOT NULL REFERENCES service(id) ON DELETE CASCADE,
	quantity INTEGER NOT NULL CHECK (quantity >= 0), --  support multiple of a service
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS booking_services;
