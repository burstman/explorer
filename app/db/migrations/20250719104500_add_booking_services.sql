-- +goose Up
CREATE TABLE booking_services (
	id SERIAL PRIMARY KEY,
	booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
	service_id INTEGER NOT NULL REFERENCES services(id),
	quantity INTEGER DEFAULT 1, --  support multiple of a service
	UNIQUE (booking_id, service_id)
);

-- +goose Down
DROP TABLE IF EXISTS booking_services;
