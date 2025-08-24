-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
returning *;

-- name: ResetUsers :exec
delete from users *;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: UpdateUser :exec
update users
set email = $2, hashed_password = $3, updated_at = NOW()
where id = $1;

-- name: MakeUserRed :exec
update users
set is_chirpy_red = true
where id = $1;
