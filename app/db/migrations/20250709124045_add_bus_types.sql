-- +goose Up
CREATE TABLE bus_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP

);

-- +goose Down
DROP TABLE IF EXISTS bus_types;
