-- name: CreateAnnouncements :exec
INSERT INTO announcements (
    id, user_id, title, content, total_images, created_at
    )
VALUES ($1, $2, $3, $4, 0, now());


-- name: GetAnnouncementsForUser :many
SELECT 
    a.*,
    c.name AS church_name,
    ep.image_url AS church_image_url,
    u.is_verified AS church_verified,
    u.username
FROM 
    announcements a
JOIN
  entity_profiles ep ON a.user_id = ep.user_id
JOIN 
    authentications u ON a.user_id = u.id
JOIN 
    churches c ON u.id = c.auth_id
JOIN 
    user_church_membership m ON m.church_id = c.id
WHERE 
    m.user_id = $1
    AND u.is_suspended = false
    AND u.is_deleted = false
    AND a.suspended_at IS NULL
    AND a.deleted_at IS NULL
    ORDER BY a.created_at DESC
      LIMIT $2 OFFSET $3;



-- name: GetAnnouncementById :one
SELECT * FROM announcements
WHERE id = $1
LIMIT 1;

-- name: GetAnnouncementsByUserID :many
SELECT
  p.*,
  c.name,
  u.username,
  u.is_verified,
  ep.entity_type,
  ep.following_count,
  ep.followers_count,
  ep.image_url AS user_image_url
FROM (
  SELECT DISTINCT ON (ps.id)
    ps.*
  FROM announcements ps
  WHERE ps.user_id = $1
  ORDER BY ps.id, ps.created_at DESC
  LIMIT $2 OFFSET $3
) p
JOIN authentications u ON p.user_id = u.id
JOIN entity_profiles ep ON p.user_id = ep.user_id
LEFT JOIN churches c ON u.id = c.auth_id
GROUP BY
    p.id,
  p.user_id,
  p.total_images,
  p.title,
  p.content,
  p.created_at,
  p.updated_at,
  p.suspended_at,
  p.deleted_at,
  u.username,
  u.is_verified,
  u.user_type,
  c.name,
  ep.entity_type,
  ep.following_count,
  ep.followers_count,
  ep.image_url
ORDER BY p.created_at DESC;


-- name: CheckAnnouncementStatus :one
SELECT
    id AS post_id,
    CASE
        WHEN deleted_at IS NOT NULL THEN 'deleted'
        WHEN suspended_at IS NOT NULL THEN 'suspended'
        ELSE 'active'
    END AS post_status
FROM
    announcements
WHERE
    id = $1;
