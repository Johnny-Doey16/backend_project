-- name: CreatePost :one
INSERT INTO posts (
    id, user_id, content, total_images, created_at
    )
VALUES ($1, $2, $3, $4, now())
RETURNING *;


-- name: IsUserFollowing :one
SELECT *
FROM follow
WHERE follower_user_id = $1 AND following_user_id = $2
LIMIT 1;


-- name: GetPostsByFollowing :many
SELECT p.*
FROM posts p
LEFT JOIN follow f ON p.user_id = f.following_user_id
WHERE f.follower_user_id = $1 OR p.user_id = $1
ORDER BY p.created_at DESC;

-- name: GetPostByUserID :many
SELECT
  p.*,
  CASE
    WHEN u.user_type = 'user' THEN us.first_name
    WHEN u.user_type = 'churchAdmin' THEN c.name
  END AS name,
  u.username,
  u.is_verified,
  u.user_type,
  us.first_name AS user_first_name,  -- Rename to avoid ambiguity
  us.last_name,
  ep.entity_type,
  ep.following_count,
  ep.followers_count,
  ep.image_url AS user_image_url,
  json_agg(pi.image_url) AS post_image_urls,
  pm.views,
  pm.likes,
  pm.comments,
  pm.reposts,
  EXISTS (
        SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1
    ) AS post_liked
FROM (
  SELECT DISTINCT ON (ps.id)
    ps.*
  FROM posts ps
  WHERE ps.user_id = $1
    AND ps.deleted_at IS NULL
    AND ps.suspended_at IS NULL
  ORDER BY ps.id, ps.created_at DESC
  LIMIT $2 OFFSET $3
) p
JOIN authentications u ON p.user_id = u.id
LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $1 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $1)
JOIN entity_profiles ep ON p.user_id = ep.user_id
LEFT JOIN users us ON u.id = us.user_id
LEFT JOIN churches c ON u.id = c.auth_id
LEFT JOIN posts_images pi ON p.id = pi.post_id
LEFT JOIN posts_metrics pm ON p.id = pm.post_id
GROUP BY
    p.id,
  p.user_id,
  p.content,
  p.created_at,
  p.total_images,
  p.updated_at,
  p.suspended_at,
  p.deleted_at,
  u.username,
  u.is_verified,
  u.user_type,
  us.first_name,
  us.last_name,
  c.name,
  ep.entity_type,
  ep.following_count,
  ep.followers_count,
  ep.image_url,
  pm.views,
  pm.likes,
  pm.comments,
  pm.reposts,
  post_liked
ORDER BY p.created_at DESC;


-- name: CheckPostStatus :one
SELECT
    id AS post_id,
    CASE
        WHEN deleted_at IS NOT NULL THEN 'deleted'
        WHEN suspended_at IS NOT NULL THEN 'suspended'
        ELSE 'active'
    END AS post_status
FROM
    posts
WHERE
    id = $1;

-- name: CreatePostMention :exec
INSERT INTO post_mentions (mentioned_user_id, post_id)
SELECT mentioned_user_id, $2
FROM UNNEST($1::uuid[]) AS mentioned_user_id;



-- name: GetPostById :one
WITH Post_With_Images AS (
    SELECT
        p.id AS post_id,
        p.content,
        p.total_images,
        p.created_at AS post_created_at,
        p.updated_at AS post_updated_at,
        p.suspended_at AS post_suspended_at,
        p.deleted_at AS post_deleted_at,
        pi.image_url AS post_image_url,
        pi.caption AS post_image_caption,
        pm.views,
        pm.likes,
        pm.comments,
        pm.reposts,
        a.username,
        e.image_url,
        e.user_id,
        a.user_type,
        a.is_verified,
        json_agg(pi.image_url) AS image_urls, -- Aggregate image URLs into a JSON array
        e.followers_count,
        e.following_count

    FROM
        posts p
    LEFT JOIN
        posts_images pi ON p.id = pi.post_id
    LEFT JOIN
        posts_metrics pm ON p.id = pm.post_id
    INNER JOIN
        authentications a ON p.user_id = a.id
    INNER JOIN 
      entity_profiles e ON p.user_id = e.user_id
    WHERE
        p.id = $1
    GROUP BY
        p.id, p.content, p.total_images, p.created_at, p.updated_at, p.suspended_at, p.deleted_at, e.user_id,
        pm.views, pm.likes, pm.comments, pm.reposts, a.username, e.image_url, a.user_type, pi.image_url, pi.caption,
        a.is_verified, e.followers_count, e.following_count
)
SELECT * FROM Post_With_Images;

-- SELECT 
--     p.id AS post_id,
--     p.content,
--     p.total_images,
--     p.created_at AS post_created_at,
--     p.updated_at AS post_updated_at,
--     p.suspended_at AS post_suspended_at,
--     p.deleted_at AS post_deleted_at,
--     pi.image_url AS post_image_url,
--     pi.caption AS post_image_caption,
--     pm.views,
--     pm.likes,
--     pm.comments,
--     pm.reposts,
--     a.username,
--     e.image_url,
--     e.user_id,
--     a.user_type,
--     a.is_verified,
--     e.followers_count,
--     e.following_count
-- FROM 
--     posts p
-- LEFT JOIN 
--     posts_images pi ON p.id = pi.post_id
-- LEFT JOIN 
--     posts_metrics pm ON p.id = pm.post_id
-- INNER JOIN 
--     authentications a ON p.user_id = a.id
-- INNER JOIN 
--     entity_profiles e ON p.user_id = e.user_id
-- WHERE 
--     p.id = $1;
