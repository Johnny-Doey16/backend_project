-- name: CreateComment :one
INSERT INTO "post_comments" (
  "user_id",
  "post_id",
  "comment_text",
  "created_at"
) VALUES (
  $1, $2, $3, NOW()
)
RETURNING *;

-- name: GetCommentByID :one
SELECT * FROM "post_comments"
WHERE "id" = $1;

-- name: ListCommentsByPostID :many
SELECT * FROM "post_comments"
WHERE "post_id" = $1
ORDER BY "created_at" DESC
LIMIT $2 OFFSET $3;

-- name: UpdateComment :exec
UPDATE "post_comments"
SET
  "comment_text" = $2,
  "updated_at" = NOW()
WHERE "id" = $1;

-- name: DeleteComment :exec
DELETE FROM "post_comments"
WHERE "id" = $1;

-- name: RecommendComments :many
SELECT
  pc.id AS comment_id,
  pc.comment_text,
  pc.created_at AS comment_created_at,
  pc.updated_at AS comment_updated_at,
  u.id AS user_id,
  u.username AS commenter_username,
  ep.image_url,
  us.first_name,
  CASE
    WHEN f.following_user_id IS NOT NULL THEN 1
    ELSE 0
  END AS is_following_commenter
FROM
  post_comments pc
JOIN authentications u ON
  pc.user_id = u.id

JOIN entity_profiles ep ON u.id = ep.user_id
JOIN users us ON u.id = us.user_id

LEFT JOIN follow f ON
  pc.user_id = f.following_user_id AND f.follower_user_id = $1
WHERE
  pc.post_id = $2
ORDER BY
  is_following_commenter DESC,
  pc.created_at DESC
  LIMIT $3 OFFSET $4;


-- SELECT
--   pc.id,
--   pc.comment_text,
--   pc.created_at,
--   u.username AS commenter_username,
--   CASE
--     WHEN f.following_user_id IS NOT NULL THEN 1
--     ELSE 0
--   END AS is_following_commenter,
--   COALESCE(cl.likes_count, 0) AS likes_count,
--   COALESCE(cl.likes_count, 0) +
--   CASE
--     WHEN f.following_user_id IS NOT NULL THEN 100
--     ELSE 0
--   END AS relevance_score
-- FROM
--   post_comments pc
-- LEFT JOIN authentications u ON
--   pc.user_id = u.id
-- LEFT JOIN (
--   SELECT
--     comment_id,
--     COUNT(*) AS likes_count
--   FROM
--     comment_likes
--   GROUP BY
--     comment_id
-- ) cl ON
--   pc.id = cl.comment_id
-- LEFT JOIN follow f ON
--   pc.user_id = f.following_user_id
--   AND f.follower_user_id = $1
-- WHERE
--   pc.post_id = $2
-- ORDER BY
--   relevance_score DESC,
--   pc.created_at DESC
--   LIMIT $3 OFFSET $4;