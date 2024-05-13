-- name: CreateUserNames :exec
INSERT INTO users (
    user_id, first_name, last_name
    )
VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT * FROM users WHERE user_id = $1;

-- name: UpdateUserNames :exec
UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'),first_name),
    last_name = COALESCE(sqlc.narg('last_name'),last_name)
WHERE
    user_id = sqlc.arg('user_id');