-- +goose Up
-- +goose StatementBegin
CREATE TABLE konnect_payment_responses (
    id SERIAL PRIMARY KEY,
    payment_ref VARCHAR(100) UNIQUE,
    status VARCHAR(50),
    amount INT,
    payment_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE konnect_payment_responses;
-- +goose StatementEnd
