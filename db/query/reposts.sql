-- name: CreateRepost :one
INSERT INTO repost (user_id, original_post_id)
VALUES ($1, $2)
RETURNING id, user_id, original_post_id, created_at;

-- name: GetRepostByID :one
SELECT id, user_id, original_post_id, created_at
FROM repost
WHERE id = $1;

-- name: GetRepostsByUserID :many
SELECT id, user_id, original_post_id, created_at
FROM repost
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteRepost :exec
DELETE FROM repost
WHERE id = $1;