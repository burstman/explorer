-- +goose Up
-- +goose StatementBegin
CREATE TABLE payment_responses (
    id SERIAL PRIMARY KEY,
    payment_ref VARCHAR(100) UNIQUE,
    status VARCHAR(50),
    amount INT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
-- +goose StatementEnd
