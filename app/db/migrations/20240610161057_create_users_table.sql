-- +goose Up
create table if not exists users(
	id SERIAL PRIMARY key,
	role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
	email text unique,
	password_hash text,
	first_name text not null,
	last_name text not null,
	phone_number VARCHAR(15),
	social_id VARCHAR(64) ,
	provider TEXT,
	CIN VARCHAR(8),
	last_verification_sent_at TIMESTAMP NULL,
	email_verified_at timestamp,
	created_at timestamp not null,
	updated_at timestamp not null,
	deleted_at timestamp
);

-- +goose Down
drop table if exists users;
