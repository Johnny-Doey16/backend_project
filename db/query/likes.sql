-- name: LikePost :many
WITH upsert_like AS (
    INSERT INTO likes (user_id, post_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, post_id) DO NOTHING
    RETURNING post_id
),
deleted_like AS (
    DELETE FROM likes
    WHERE (user_id, post_id) = ($1, $2)
    RETURNING post_id
),
updated_likes AS (
    UPDATE posts_metrics
    SET likes = CASE
        WHEN EXISTS (SELECT 1 FROM upsert_like) THEN likes + 1
        WHEN EXISTS (SELECT 1 FROM deleted_like) THEN GREATEST(likes - 1, 0)
        ELSE likes
    END
    WHERE post_id = $2
    RETURNING *, EXISTS (SELECT 1 FROM upsert_like) AS liked
)
SELECT *, CASE WHEN liked THEN true ELSE false END AS liked FROM updated_likes;


-- name: AddLike :exec
INSERT INTO likes (user_id, post_id)
VALUES ($1, $2)
ON CONFLICT (user_id, post_id) DO NOTHING;

-- name: GetUserLikedPost :one
SELECT * FROM likes
WHERE user_id = $1 AND
post_id = $2;

-- name: IncrementLikeCount :one
UPDATE posts_metrics
SET likes = likes + 1
WHERE post_id = $1
RETURNING *;

-- name: RemoveLike :exec
DELETE FROM likes
WHERE user_id = $1 AND post_id = $2;

-- name: DecrementLikeCount :one
UPDATE posts_metrics
SET likes = GREATEST(likes - 1, 0) -- Prevents the likes count from going negative
WHERE post_id = $1
RETURNING *;