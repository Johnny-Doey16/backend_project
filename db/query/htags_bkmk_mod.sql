-- name: CreateHashtag :one
INSERT INTO "hashtag" ("hash_tag") VALUES ($1) RETURNING *;

-- name: GetHashtag :one
SELECT * FROM "hashtag" WHERE "id" = $1 LIMIT 1;

-- name: UpdateHashtag :exec
UPDATE "hashtag" SET "hash_tag" = $2 WHERE "id" = $1;

-- name: DeleteHashtag :exec
DELETE FROM "hashtag" WHERE "id" = $1;



-- name: CreatePostHashtag :one
INSERT INTO "post_hashtag" ("post_id", "hashtag_id") VALUES ($1, $2) RETURNING *;

-- name: GetPostHashtags :many
SELECT * FROM "post_hashtag" WHERE "post_id" = $1;

-- name: DeletePostHashtag :exec
DELETE FROM "post_hashtag" WHERE "id" = $1;


-- name: CreateBookmark :exec
INSERT INTO "bookmarks" ("user_id", "post_id") VALUES ($1, $2);-- ON CONFLICT DO NOTHING;

-- name: GetUserBookmarks :many
SELECT
    p.id AS post_id,
    p.content,
    p.created_at AS post_created_at,
    p.deleted_at AS post_deleted_at,
    p.suspended_at AS post_suspended_at,
    b.created_at AS bookmarked_at,
    p.user_id AS author_user_id,
    u.username AS author_username,
    us.first_name AS author_first_name,
    ep.image_url AS author_image_url,
    json_agg(pi.image_url) AS post_image_urls,
    pm.likes,
    pm.views,
    pm.comments,
    pm.reposts,
    EXISTS (
        SELECT 1
        FROM "likes" l
        WHERE l.post_id = p.id AND l.user_id = $1
    ) AS liked
FROM
    "bookmarks" b
JOIN "posts" p ON b.post_id = p.id
JOIN "authentications" u ON p.user_id = u.id
JOIN "users" us ON u.id = us.user_id
JOIN "entity_profiles" ep ON u.id = ep.user_id
LEFT JOIN "posts_images" pi ON p.id = pi.post_id
LEFT JOIN "posts_metrics" pm ON p.id = pm.post_id
WHERE
    b.user_id = $1 AND
    p.suspended_at IS NULL AND
    p.deleted_at IS NULL
GROUP BY
    p.id, b.created_at, p.content, p.created_at, p.user_id, u.username, us.first_name, ep.image_url, pm.likes, pm.views, pm.comments, pm.reposts
ORDER BY bookmarked_at DESC
LIMIT $2 OFFSET $3;


-- name: DeleteBookmark :one
DELETE FROM "bookmarks" WHERE "id" = $1 RETURNING *;



-- name: BlockUserSM :one
INSERT INTO "blocked_users" ("blocking_user_id", "blocked_user_id", "reason") VALUES ($1, $2, $3) RETURNING *;

-- name: GetBlockedUsersByBlocker :many
SELECT * FROM "blocked_users" WHERE "blocking_user_id" = $1;

-- name: UnblockUser :exec
DELETE FROM "blocked_users" WHERE "id" = $1;



-- name: ReportPost :one
INSERT INTO "reported_posts" ("post_id", "user_id", "reason") VALUES ($1, $2, $3) RETURNING *;

-- TODO: Admin functions
-- name: GetReportedPosts :many
SELECT * FROM "reported_posts";

-- name: DeleteReportedPost :exec
DELETE FROM "reported_posts" WHERE "id" = $1;


-- TODO: Admin functions
-- name: AddToModerationQueue :one
INSERT INTO "moderation_queue" ("user_id", "post_id", "report_reason", "status") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetModerationQueue :many
SELECT * FROM "moderation_queue" WHERE "status" = $1;

-- name: UpdateModerationQueueStatus :exec
UPDATE "moderation_queue" SET "status" = $2 WHERE "id" = $1;

-- name: RemoveFromModerationQueue :exec
DELETE FROM "moderation_queue" WHERE "id" = $1;