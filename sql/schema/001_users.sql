-- +goose Up
create table users(
    id uuid unique primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text unique not null,
    hashed_password text not null,
    is_chirpy_red boolean default false
);

-- +goose Down
drop table users;
