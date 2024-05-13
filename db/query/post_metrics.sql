

-- name: CreateMetric :exec
INSERT INTO posts_metrics (post_id, views, likes, comments, reposts) VALUES ($1, $2, $3, $4, $5);

-- name: IncrementViews :exec
UPDATE posts_metrics SET views = views + 1 WHERE post_id = $1;

-- name: IncrementReposts :one
UPDATE posts_metrics SET reposts = reposts + 1 WHERE post_id = $1
RETURNING reposts AS total_reposts;

-- name: IncrementComments :one
UPDATE posts_metrics SET comments = comments + 1 WHERE post_id = $1
RETURNING *;

-- name: IncrementLikes :exec
UPDATE posts_metrics SET likes = likes + 1 WHERE post_id = $1;