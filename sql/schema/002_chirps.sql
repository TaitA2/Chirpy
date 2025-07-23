-- +goose Up
create table chirps(
    id uuid unique primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id uuid not null references users on delete cascade,
    constraint fk_user_id foreign key(user_id) references users (id)
);

-- +goose Down
drop table chirps;

