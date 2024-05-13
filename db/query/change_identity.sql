-- name: CreateChangeIdRequest :exec
INSERT INTO change_identifier_requests (user_id, identifier, token, type, expires_at)
VALUES ($1, $2, $3, $4, $5);


-- name: GetChangeIdRequestByID :one
SELECT * FROM change_identifier_requests WHERE id = $1 LIMIT 1;

-- name: GetChangeIdRequestByToken :one
SELECT * FROM change_identifier_requests WHERE token = $1 LIMIT 1;

-- name: UpdateChangeIdRequest :exec
UPDATE change_identifier_requests SET used = $1 WHERE id = $2;

-- name: UpdateChangeIdByToken :exec
UPDATE change_identifier_requests SET used = $1 WHERE token = $2;

-- name: DeleteChangeIdRequestByID :exec
DELETE FROM change_identifier_requests WHERE id = $1;
