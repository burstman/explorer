-- +goose Up
CREATE TABLE campsite_buses (
    id SERIAL PRIMARY KEY,
    campsite_id INT NOT NULL REFERENCES campsites(id) ON DELETE CASCADE,
    bus_type_id INT NOT NULL REFERENCES bus_types(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS campsite_buses;
