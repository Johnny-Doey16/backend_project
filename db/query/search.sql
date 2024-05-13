-- name: GetSearchResult :many
WITH search_data AS (
  SELECT
    'users' AS source,
    au.id AS id,
    au.id AS post_user_id,
    NULL AS uid,
    ep.header_image_url AS header_image_url,
    ep.about AS about,
    au.is_verified AS is_verified,
    au.username AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url AS image_url,
    0 AS members_count,
    -- u.first_name || ' ' || u.last_name AS name,
    u.first_name AS name,
    au.user_type AS user_type,
    au.created_at AS created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
    EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = au.id) AS is_following,
    NULL AS city,
    NULL AS state,
    FALSE AS post_liked
  FROM authentications au
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $2 AND bu.blocked_user_id = au.id)
    OR (bu.blocking_user_id = au.id AND bu.blocked_user_id = $2)
  JOIN users u ON au.id = u.user_id
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE to_tsvector('english', au.username || ' ' || u.first_name || ' ' || ' ' || ep.about) @@ plainto_tsquery('english', $1)
  AND au.deleted_at IS NULL
  AND au.suspended_at IS NULL
  AND bu.id IS NULL
  
  UNION ALL
  
  SELECT
    'churches' AS source,
    c.auth_id AS id,
    c.auth_id AS post_user_id,
    NULL AS uid,
    ep.header_image_url,
    ep.about,
    a.is_verified,
    NULL AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url,
    c.members_count,
    c.name,
    a.user_type,
    a.created_at,
    ucm.active AS is_member,
    c.id AS church_id,
    a.username AS church_username,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
    EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = a.id) AS is_following,
    l.city,
    l.state,
    FALSE AS post_liked
  FROM churches c
  JOIN authentications a ON c.auth_id = a.id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $2 AND bu.blocked_user_id = c.auth_id)
    OR (bu.blocking_user_id = c.auth_id AND bu.blocked_user_id = $2)
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN church_locations l ON a.id = l.auth_id
  LEFT JOIN user_church_membership ucm ON ucm.church_id = c.id AND ucm.user_id = $2
  WHERE to_tsvector('english', c.name || ' ' || a.username || ' ' || ep.about) @@ plainto_tsquery('english', $1)
  AND a.deleted_at IS NULL
  AND a.suspended_at IS NULL
  AND bu.id IS NULL
  
  UNION ALL
  
  SELECT
    'posts' AS source,
    p.id,
    a.id AS uid,
    p.user_id AS post_user_id,
    NULL AS header_image_url,
    ep.about,
    a.is_verified,
    a.username,
    p.content,
    pi.caption AS captions,
    ep.image_url,
    0 AS members_count,
    CASE
        WHEN a.user_type = 'user' THEN u.first_name-- || ' ' || u.last_name
        WHEN a.user_type = 'churchAdmin' THEN c.name
        ELSE NULL
    END AS name,
    a.user_type,
    p.created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    (SELECT json_agg(DISTINCT image_url) FROM posts_images WHERE post_id = p.id) AS post_images,
    pm.comments,
    pm.views,
    pm.likes,
    pm.reposts,
    NULL AS is_following,
    NULL AS city,
    NULL AS state,
    EXISTS (
        SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $2
    ) AS post_liked
  FROM posts p
  JOIN authentications a ON p.user_id = a.id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $2 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $2)
  JOIN users u ON a.id = u.user_id
  LEFT JOIN churches c ON a.id = c.auth_id
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN posts_metrics pm ON p.id = pm.post_id
  LEFT JOIN posts_images pi ON p.id = pi.post_id
  WHERE to_tsvector('english', p.content || ' ' || a.username || ' ' || u.first_name ) @@ plainto_tsquery('english', $1)
  AND p.deleted_at IS NULL
  AND p.suspended_at IS NULL
  AND bu.id IS NULL
)
SELECT
  source,
  json_agg(row) AS data
FROM (
  SELECT DISTINCT ON (source, id)
    source,
    id,
    post_user_id,
    header_image_url,
    about,
    is_verified,
    username,
    content,
    captions,
    image_url,
    members_count,
    user_type,
    created_at,
    is_member,
    church_id,
    church_username,
    name,
    post_images,
    comments,
    views,
    likes,
    reposts,
    is_following,
    city,
    state,
    post_liked
  FROM search_data
  ORDER BY source, id
  LIMIT $3
  OFFSET $4
) AS row
GROUP BY source;

-- TODO: Limit and offset for users, church, posts individually rather than for all.


-- name: GetSearchResultOldNoFTS :many
WITH search_data AS (
  SELECT
    'users' AS source,
    au.id AS id,
    au.id AS post_user_id,
    ep.header_image_url AS header_image_url,
    ep.about AS about,
    au.is_verified AS is_verified,
    au.username AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url AS image_url,
    0 AS members_count,
    --  u.first_name || ' ' || u.last_name AS name,
     u.first_name AS name,
    au.user_type AS user_type,
    au.created_at AS created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
    EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = au.id) AS is_following,
    NULL AS city,
    NULL AS state,
    FALSE AS post_liked
  FROM authentications au
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $2 AND bu.blocked_user_id = au.id)
    OR (bu.blocking_user_id = au.id AND bu.blocked_user_id = $2)
  JOIN users u ON au.id = u.user_id
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE (
    au.username ILIKE '%'|| $1 ||'%' OR
    u.first_name ILIKE '%'|| $1 ||'%' OR
    ep.about ILIKE '%'|| $1 ||'%'
  )
  AND au.deleted_at IS NULL
  AND au.suspended_at IS NULL
  AND bu.id IS NULL
  
  UNION ALL
  
  SELECT
    'churches' AS source,
    c.auth_id AS id,
    c.auth_id AS post_user_id,
    ep.header_image_url,
    ep.about,
    a.is_verified,
    NULL AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url,
    c.members_count,
    c.name,
    a.user_type,
    a.created_at,
    ucm.active AS is_member,
    c.id AS church_id,
    a.username AS church_username,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
    EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = a.id) AS is_following,
    l.city,
    l.state,
    FALSE AS post_liked
  FROM churches c
  JOIN authentications a ON c.auth_id = a.id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $2 AND bu.blocked_user_id = c.auth_id)
    OR (bu.blocking_user_id = c.auth_id AND bu.blocked_user_id = $2)
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN church_locations l ON a.id = l.auth_id
  LEFT JOIN user_church_membership ucm ON ucm.church_id = c.id AND ucm.user_id = $2
  WHERE (
    c.name ILIKE '%'|| $1 ||'%' OR
    ep.about ILIKE '%'|| $1 ||'%' OR
    a.username ILIKE '%'|| $1 ||'%'
  )
  AND a.deleted_at IS NULL
  AND a.suspended_at IS NULL
  AND bu.id IS NULL
  
  UNION ALL
  
  SELECT
    'posts' AS source,
    p.id,
    p.user_id AS post_user_id,
    NULL AS header_image_url,
    ep.about,
    a.is_verified,
    a.username,
    p.content,
    pi.caption AS captions,
    ep.image_url,
    0 AS members_count,
    CASE
        WHEN a.user_type = 'user' THEN u.first_name-- || ' ' || u.last_name
        WHEN a.user_type = 'churchAdmin' THEN c.name
        ELSE NULL
    END AS name,
    a.user_type,
    p.created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    (SELECT json_agg(DISTINCT image_url) FROM posts_images WHERE post_id = p.id) AS post_images,
    pm.comments,
    pm.views,
    pm.likes,
    pm.reposts,
    NULL AS is_following,
    NULL AS city,
    NULL AS state,
    EXISTS (
        SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $2
    ) AS post_liked
  FROM posts p
  JOIN authentications a ON p.user_id = a.id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $2 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $2)
  LEFT JOIN users u ON a.id = u.user_id
  LEFT JOIN churches c ON a.id = c.auth_id
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN posts_metrics pm ON p.id = pm.post_id
  LEFT JOIN posts_images pi ON p.id = pi.post_id
  WHERE (
    p.content ILIKE '%'|| $1 ||'%'
    -- OR
    -- a.username ILIKE '%'|| $1 ||'%' OR
    -- u.first_name ILIKE '%'|| $1 ||'%' OR
    -- u.last_name ILIKE '%'|| $1 ||'%'
  )
  AND p.deleted_at IS NULL
  AND p.suspended_at IS NULL
  AND bu.id IS NULL
  GROUP BY p.id, a.id, u.id, ep.id, pm.id, pi.caption, p.content, c.name
)
SELECT
  source,
  json_agg(row) AS data
FROM (
  SELECT DISTINCT ON (source, id)
    source,
    id,
    post_user_id,
    header_image_url,
    about,
    is_verified,
    username,
    content,
    captions,
    image_url,
    members_count,
    user_type,
    created_at,
    is_member,
    church_id,
    church_username,
    name,
    post_images,
    comments,
    views,
    likes,
    reposts,
    is_following,
    city,
    state,
    post_liked
  FROM search_data
  ORDER BY source, id
  LIMIT $3
  OFFSET $4
) AS row
GROUP BY source;