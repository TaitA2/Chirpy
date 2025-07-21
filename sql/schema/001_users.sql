-- +goose Up
create table users(
    id uuid unique primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text unique not null
);

-- +goose Down
drop table users;
