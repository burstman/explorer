-- +goose Up
create table if not exists users(
	id SERIAL PRIMARY key,
	role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
	email text unique not null,
	password_hash text not null,
	first_name text not null,
	last_name text not null,
	phone_number VARCHAR(15),
	social_Link VARCHAR(64),
	CIN VARCHAR(8),
	email_verified_at timestamp,
	created_at timestamp not null,
	updated_at timestamp not null,
	deleted_at timestamp
);

-- +goose Down
drop table if exists users;
