-- +goose Up
CREATE TABLE campsites (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    image_url TEXT,
    location TEXT,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    available_from DATE,
    available_to DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS campsites;
