-- name: CreateChirp :one
insert into chirps(created_at, updated_at, body, user_id)
values (NOW(), NOW(), $1, $2)
returning *;

-- name: GetChirpsByUser :many
select * from chirps where user_id = $1 order by created_at desc;

