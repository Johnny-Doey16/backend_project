-- name: IncrementViewCount :exec
UPDATE posts_metrics
SET views = views + 1
WHERE post_id = $1;

-- name: TrackPostView :one
-- WITH inserted_view AS (
--     INSERT INTO views (user_id, post_id)
--     VALUES ($1, $2)
--     ON CONFLICT (user_id, post_id) DO UPDATE SET viewed_at = CURRENT_TIMESTAMP
--     RETURNING post_id, viewed_at
-- ),
-- updated_views AS (
--     UPDATE posts_metrics
--     SET views = views + 1
--     WHERE post_id = $2
--     RETURNING *, true AS viewed
-- )
-- SELECT * FROM updated_views;