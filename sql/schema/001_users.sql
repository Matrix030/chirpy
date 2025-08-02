-- +goose Up
create extension if not exists "pgcrypto";
create table "users" (
	id uuid primary key default gen_random_uuid(),
	created_at timestamp not null,
	updated_at timestamp not null,
	email text not null unique
);

-- +goose Down
drop table "users";
