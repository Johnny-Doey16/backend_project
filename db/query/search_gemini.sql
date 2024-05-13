WITH search_data AS (
  SELECT
    'users' AS source,
    au.id AS id,
    ep.header_image_url AS header_image_url,
    ep.about AS about,
    au.is_verified AS is_verified,
    au.username AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url AS image_url,
    0 AS members_count,
    u.first_name AS first_name,
    u.last_name AS last_name,
    au.user_type AS user_type,
    au.created_at AS created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    NULL AS name,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
     EXISTS (
    SELECT 1
    FROM follow fu
    WHERE fu.follower_user_id = $2 -- current user
      AND fu.following_user_id = au.id
  ) AS is_following,
    -- EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = au.id) AS is_following
    NULL AS city,
    NULL AS state


  FROM authentications au
  JOIN users u ON au.id = u.user_id
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE (
    au.email ILIKE '%'|| $1 ||'%' OR
    au.phone ILIKE '%'|| $1 ||'%' OR
    au.username ILIKE '%'|| $1 ||'%' OR
    u.first_name ILIKE '%'|| $1 ||'%' OR
    u.last_name ILIKE '%'|| $1 ||'%' OR
    ep.about ILIKE '%'|| $1 ||'%'
  )
  UNION ALL
  SELECT
    'churches' AS source,
    c.auth_id AS auth_id,
    ep.header_image_url AS header_image_url,
    ep.about AS about,
    a.is_verified AS is_verified,
    NULL AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url AS image_url,
    c.members_count AS members_count,
    NULL AS first_name,
    NULL AS last_name,
    a.user_type AS user_type,
    a.created_at AS created_at, -- TODO: Not showing data in user, church and post
    ucm.active AS is_member,
    c.id AS church_id,
    a.username AS church_username,
    c.name AS name,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
     (
    SELECT COUNT(*)
    FROM follow fu
    WHERE fu .follower_user_id = $2 -- current user
      AND fu .following_user_id = a.id
  ) > 0 AS is_following,
    -- EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = au.id) AS is_following



    l.city AS city,
    l.state AS state
    -- l.address AS address,
    -- l.postalcode AS postalcode,
    -- l.lga AS lga,
    -- l.country AS country,
  FROM churches c
  JOIN authentications a ON c.auth_id = a.id
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN church_locations l ON a.id = l.auth_id
  LEFT JOIN user_church_membership ucm ON ucm.church_id = c.id AND ucm.user_id = $2
  WHERE (
    c.name ILIKE '%'|| $1 ||'%' OR
    ep.about ILIKE '%'|| $1 ||'%'
  )

  UNION ALL
  SELECT
    'posts' AS source,
    p.id AS id,
    NULL AS header_image_url,
    ep.about AS about,
    a.is_verified AS is_verified,
    a.username AS username,
    p.content AS content,
    pi.caption AS captions,
    ep.image_url AS image_url,
    0 AS members_count,
    u.first_name AS first_name,
    u.last_name AS last_name,
    a.user_type AS user_type,
    p.created_at AS created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    NULL AS name,
    -- pi.image_url AS post_images,
    -- array_agg(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL) AS post_images,
    (SELECT json_agg(DISTINCT image_url) FROM posts_images WHERE post_id = p.id) AS post_images,


    pm.comments AS comments,
    pm.views AS views,
    pm.likes AS likes,
    pm.reposts AS reposts,
    NULL AS is_following,
    NULL AS city,
    NULL AS state

    
  FROM posts p
  JOIN authentications a ON p.user_id = a.id
  JOIN users u ON a.id = u.user_id --! check
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN posts_metrics pm ON p.id = pm.post_id
  LEFT JOIN posts_images pi ON p.id = pi.post_id
-- INNER JOIN posts_images pi ON p.id = pi.post_id
  WHERE p.content ILIKE '%'|| $1 ||'%'
  GROUP BY p.id, a.id, u.id, ep.id, pm.id, pi.caption, p.content -- Make sure all selected columns are in the GROUP BY clause

)
SELECT
  source,
  json_agg(row) AS data
FROM (
  SELECT DISTINCT ON (source, id)
    source,
    id,
    header_image_url,
    about,
    is_verified,
    username,
    content,
    captions,
    image_url,
    members_count,
    first_name,
    last_name,
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
    state
  FROM search_data
) AS row
GROUP BY source;

-- Claude
-- WITH search_term AS (
--   SELECT to_tsvector('english', ?) AS query
-- ),
-- SELECT
--   json_build_object(
--     'Users', (
--       SELECT json_agg(
--         json_build_object(
--           'id', u.id,
--           'user_id', a.id,
--           'first_name', u.first_name,
--           'last_name', u.last_name,
--           'email', a.email,
--           'phone', a.phone,
--           'username', a.username,
--           'user_type', a.user_type,
--           'is_suspended', a.is_suspended,
--           'is_deleted', a.is_deleted,
--           'is_verified', a.is_verified,
--           'is_email_verified', a.is_email_verified,
--           'profile', (
--             SELECT json_build_object(
--               'image_url', ep.image_url,
--               'following_count', ep.following_count,
--               'followers_count', ep.followers_count,
--               'entity_type', ep.entity_type,
--               'posts_count', ep.posts_count,
--               'header_image_url', ep.header_image_url,
--               'about', ep.about,
--               'website', ep.website
--             )
--             FROM entity_profiles ep
--             WHERE ep.user_id = a.id
--           ),
--           'is_following', (
--             SELECT EXISTS (
--               SELECT 1
--               FROM follow f
--               WHERE f.following_user_id = a.id AND f.follower_user_id = (SELECT id FROM authentications WHERE id = ?)
--             )
--           ),
--           'is_followed_by', (
--             SELECT EXISTS (
--               SELECT 1
--               FROM follow f
--               WHERE f.follower_user_id = a.id AND f.following_user_id = (SELECT id FROM authentications WHERE id = ?)
--             )
--           )
--         )
--       )
--       FROM users u
--       JOIN authentications a ON u.user_id = a.id
--       WHERE a.email @@ (SELECT query FROM search_term)
--         OR a.phone @@ (SELECT query FROM search_term)
--         OR a.username @@ (SELECT query FROM search_term)
--         OR u.first_name @@ (SELECT query FROM search_term)
--         OR u.last_name @@ (SELECT query FROM search_term)
--     ),
--     'Churches', (
--       SELECT json_agg(
--         json_build_object(
--           'id', c.id,
--           'name', c.name,
--           'members_count', c.members_count,
--           'denomination_id', c.denomination_id,
--           'auth', (
--             SELECT json_build_object(
--               'id', a.id,
--               'email', a.email,
--               'phone', a.phone,
--               'username', a.username,
--               'user_type', a.user_type,
--               'is_suspended', a.is_suspended,
--               'is_deleted', a.is_deleted,
--               'is_verified', a.is_verified,
--               'is_email_verified', a.is_email_verified
--             )
--             FROM authentications a
--             WHERE a.id = c.auth_id
--           ),
--           'profile', (
--             SELECT json_build_object(
--               'image_url', ep.image_url,
--               'following_count', ep.following_count,
--               'followers_count', ep.followers_count,
--               'entity_type', ep.entity_type,
--               'posts_count', ep.posts_count,
--               'header_image_url', ep.header_image_url,
--               'about', ep.about,
--               'website', ep.website
--             )
--             FROM entity_profiles ep
--             WHERE ep.user_id = c.auth_id
--           ),
--           'location', (
--             SELECT json_build_object(
--               'address', cl.address,
--               'city', cl.city,
--               'postalCode', cl.postalCode,
--               'lga', cl.lga,
--               'state', cl.state,
--               'country', cl.country,
--               'location', ST_AsGeoJSON(cl.location)::jsonb
--             )
--             FROM church_locations cl
--             WHERE cl.auth_id = c.auth_id
--           ),
--           'is_member', (
--             SELECT EXISTS (
--               SELECT 1
--               FROM user_church_membership ucm
--               WHERE ucm.church_id = c.id AND ucm.user_id = (SELECT id FROM authentications WHERE id = ?)
--             )
--           ),
--           'is_following', (
--             SELECT EXISTS (
--               SELECT 1
--               FROM follow f
--               WHERE f.following_user_id = c.auth_id AND f.follower_user_id = (SELECT id FROM authentications WHERE id = ?)
--             )
--           ),
--           'is_followed_by', (
--             SELECT EXISTS (
--               SELECT 1
--               FROM follow f
--               WHERE f.follower_user_id = c.auth_id AND f.following_user_id = (SELECT id FROM authentications WHERE id = ?)
--             )
--           )
--         )
--       )
--       FROM churches c
--       WHERE c.name @@ (SELECT query FROM search_term)
--     ),
--     'Posts', (
--       SELECT json_agg(
--         json_build_object(
--           'id', p.id,
--           'user_id', p.user_id,
--           'content', p.content,
--           'total_images', p.total_images,
--           'created_at', p.created_at,
--           'updated_at', p.updated_at,
--           'suspended_at', p.suspended_at,
--           'deleted_at', p.deleted_at,
--           'images', (
--             SELECT json_agg(
--               json_build_object(
--                 'id', pi.id,
--                 'image_url', pi.image_url,
--                 'caption', pi.caption
--               )
--             )
--             FROM posts_images pi
--             WHERE pi.post_id = p.id
--           ),
--           'metrics', (
--             SELECT json_build_object(
--               'views', pm.views,
--               'likes', pm.likes,
--               'comments', pm.comments,
--               'reposts', pm.reposts
--             )
--             FROM posts_metrics pm
--             WHERE pm.post_id = p.id
--           )
--         )
--       )
--       FROM posts p
--       WHERE p.content @@ (SELECT query FROM search_term)
--     )
--   ) AS result;
-- au.updated_at AS updated_at,
--     au.created_at AS created_at,
--     au.is_suspended AS is_suspended,
--     au.is_deleted AS is_deleted,
--     au.is_verified AS is_verified,
--     au.is_email_verified AS is_email_verified,
    -- au.login_attempts AS login_attempts,

-- name: GetSearchResultOld :many
WITH search_data AS (
  SELECT
    'users' AS source,
    au.id AS id,
    ep.header_image_url AS header_image_url,
    ep.about AS about,
    au.is_verified AS is_verified,
    au.username AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url AS image_url,
    0 AS members_count,
    u.first_name AS first_name,
    u.last_name AS last_name,
    au.user_type AS user_type,
    au.created_at AS created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    NULL AS name,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
    EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = au.id) AS is_following,
    NULL AS city,
    NULL AS state


  FROM authentications au
  JOIN users u ON au.id = u.user_id
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE (
    au.username ILIKE '%'|| $1 ||'%' OR
    u.first_name ILIKE '%'|| $1 ||'%' OR
    u.last_name ILIKE '%'|| $1 ||'%' OR
    ep.about ILIKE '%'|| $1 ||'%'
  )
  UNION ALL
  SELECT
    'churches' AS source,
    c.auth_id AS auth_id,
    ep.header_image_url AS header_image_url,
    ep.about AS about,
    a.is_verified AS is_verified,
    NULL AS username,
    NULL AS content,
    NULL AS captions,
    ep.image_url AS image_url,
    c.members_count AS members_count,
    NULL AS first_name,
    NULL AS last_name,
    a.user_type AS user_type,
    a.created_at AS created_at,
    ucm.active AS is_member,
    c.id AS church_id,
    a.username AS church_username,
    c.name AS name,
    NULL::json AS post_images,
    0 AS comments,
    0 AS views,
    0 AS likes,
    0 AS reposts,
    EXISTS (SELECT 1 FROM follow fu WHERE fu.follower_user_id = $2 AND fu.following_user_id = a.id) AS is_following,
    l.city AS city,
    l.state AS state

  FROM churches c
  JOIN authentications a ON c.auth_id = a.id
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN church_locations l ON a.id = l.auth_id
  LEFT JOIN user_church_membership ucm ON ucm.church_id = c.id AND ucm.user_id = $2
  WHERE (
    c.name ILIKE '%'|| $1 ||'%' OR
    ep.about ILIKE '%'|| $1 ||'%' OR
    a.username ILIKE '%'|| $1 ||'%'
  )

  UNION ALL
  SELECT
    'posts' AS source,
    p.id AS id,
    NULL AS header_image_url,
    ep.about AS about,
    a.is_verified AS is_verified,
    a.username AS username,
    p.content AS content,
    pi.caption AS captions,
    ep.image_url AS image_url,
    0 AS members_count,
    u.first_name AS first_name,
    u.last_name AS last_name,
    a.user_type AS user_type,
    p.created_at AS created_at,
    NULL AS is_member,
    NULL AS church_id,
    NULL AS church_username,
    NULL AS name,
    (SELECT json_agg(DISTINCT image_url) FROM posts_images WHERE post_id = p.id) AS post_images,
    pm.comments AS comments,
    pm.views AS views,
    pm.likes AS likes,
    pm.reposts AS reposts,
    NULL AS is_following,
    NULL AS city,
    NULL AS state
    
  FROM posts p
  JOIN authentications a ON p.user_id = a.id
  JOIN users u ON a.id = u.user_id --! check
  JOIN entity_profiles ep ON a.id = ep.user_id
  JOIN posts_metrics pm ON p.id = pm.post_id
  LEFT JOIN posts_images pi ON p.id = pi.post_id
-- INNER JOIN posts_images pi ON p.id = pi.post_id
  WHERE
  (
    p.content ILIKE '%'|| $1 ||'%' OR
    a.username ILIKE '%'|| $1 ||'%' OR
    u.first_name ILIKE '%'|| $1 ||'%' OR
    u.last_name ILIKE '%'|| $1 ||'%'
  )
  GROUP BY p.id, a.id, u.id, ep.id, pm.id, pi.caption, p.content -- Make sure all selected columns are in the GROUP BY clause

)
SELECT
  source,
  json_agg(row) AS data
FROM (
  SELECT DISTINCT ON (source, id)
    source,
    id,
    header_image_url,
    about,
    is_verified,
    username,
    content,
    captions,
    image_url,
    members_count,
    first_name,
    last_name,
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
    state
  FROM search_data
  -- ORDER BY source, id -- Add an ORDER BY clause for DISTINCT ON
  LIMIT $3 -- Limit the number of rows per page
  OFFSET $4 -- Offset for pagination
) AS row
GROUP BY source;
-- LIMIT $3 OFFSET $4;

-- ! GPT3
-- name: GetSearchResult2 :many
WITH SearchResults AS (
  SELECT
    'Users' AS category,
    A.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    U.first_name,
    U.last_name,
    U.user_id,
    F.follower_user_id AS is_follower,
    F.following_user_id AS is_following
  FROM authentications A
  LEFT JOIN users U ON A.id = U.user_id
  LEFT JOIN entity_profiles E ON U.user_id = E.user_id
  LEFT JOIN follow F ON A.id = F.follower_user_id
  WHERE
    A.email ILIKE '%your_keyword%' OR
    A.username ILIKE '%your_keyword%' OR
    U.first_name ILIKE '%your_keyword%' OR
    U.last_name ILIKE '%your_keyword%'
  
  UNION ALL

  SELECT
    'Churches' AS category,
    C.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    NULL AS first_name,
    NULL AS last_name,
    NULL AS user_id,
    F.follower_user_id AS is_follower,
    F.following_user_id AS is_following
  FROM churches C
  INNER JOIN authentications A ON C.auth_id = A.id
  LEFT JOIN entity_profiles E ON A.id = E.user_id
  LEFT JOIN follow F ON A.id = F.follower_user_id
  WHERE
    C.name ILIKE '%your_keyword%'

  UNION ALL

  SELECT
    'Posts' AS category,
    P.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    U.first_name,
    U.last_name,
    U.user_id,
    F.follower_user_id AS is_follower,
    F.following_user_id AS is_following
  FROM posts P
  INNER JOIN authentications A ON P.user_id = A.id
  LEFT JOIN users U ON A.id = U.user_id
  LEFT JOIN entity_profiles E ON U.user_id = E.user_id
  LEFT JOIN follow F ON A.id = F.follower_user_id
  WHERE
    P.content ILIKE '%your_keyword%'
)
SELECT
  category,
  json_agg(
    json_build_object(
      'result_id', result_id,
      'email', email,
      'phone', phone,
      'username', username,
      'profile_image', profile_image,
      'following_count', following_count,
      'followers_count', followers_count,
      'entity_type', entity_type,
      'first_name', first_name,
      'last_name', last_name,
      'user_id', user_id,
      'is_follower', is_follower,
      'is_following', is_following
    )
  ) AS results
FROM SearchResults
GROUP BY category;


-- name: GetSearchResult2New :many
WITH SearchResults AS (
  SELECT
    'Users' AS category,
    A.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    U.first_name,
    U.last_name,
    U.user_id,
    F.follower_user_id AS is_follower,
    F.following_user_id AS is_following
  FROM authentications A
  LEFT JOIN users U ON A.id = U.user_id
  LEFT JOIN entity_profiles E ON U.user_id = E.user_id
  LEFT JOIN follow F ON A.id = F.follower_user_id
  WHERE
    A.email ILIKE '%' || $1 || '%' OR
    A.username ILIKE '%' || $1 || '%' OR
    U.first_name ILIKE '%' || $1 || '%' OR
    U.last_name ILIKE '%' || $1 || '%'
  
  UNION ALL

  SELECT
    'Churches' AS category,
    C.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    NULL AS first_name,
    NULL AS last_name,
    NULL AS user_id,
    F.follower_user_id AS is_follower,
    F.following_user_id AS is_following
  FROM churches C
  INNER JOIN authentications A ON C.auth_id = A.id
  LEFT JOIN entity_profiles E ON A.id = E.user_id
  LEFT JOIN follow F ON A.id = F.follower_user_id
  WHERE
    C.name ILIKE '%' || $1 || '%'

  UNION ALL

  SELECT
    'Posts' AS category,
    P.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    U.first_name,
    U.last_name,
    U.user_id,
    F.follower_user_id AS is_follower,
    F.following_user_id AS is_following
  FROM posts P
  INNER JOIN authentications A ON P.user_id = A.id
  LEFT JOIN users U ON A.id = U.user_id
  LEFT JOIN entity_profiles E ON U.user_id = E.user_id
  LEFT JOIN follow F ON A.id = F.follower_user_id
  WHERE
    P.content ILIKE '%' || $1 || '%'
)
SELECT
  category,
  json_agg(
    json_build_object(
      'result_id', result_id,
      'email', email,
      'phone', phone,
      'username', username,
      'profile_image', profile_image,
      'following_count', following_count,
      'followers_count', followers_count,
      'entity_type', entity_type,
      'first_name', first_name,
      'last_name', last_name,
      'user_id', user_id,
      'is_follower', is_follower,
      'is_following', is_following
    )
  ) AS results
FROM SearchResults
GROUP BY category;

-- name: GetSearchResult2New1 :many
WITH SearchResults AS (
  SELECT
    'Users' AS category,
    A.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    U.first_name,
    U.last_name,
    U.user_id,
    CASE WHEN F1.follower_user_id IS NOT NULL THEN true ELSE false END AS is_follower,
    CASE WHEN F2.following_user_id IS NOT NULL THEN true ELSE false END AS is_following
  FROM authentications A
  LEFT JOIN users U ON A.id = U.user_id
  LEFT JOIN entity_profiles E ON U.user_id = E.user_id
  LEFT JOIN follow F1 ON A.id = F1.follower_user_id AND F1.following_user_id = $2
  LEFT JOIN follow F2 ON A.id = F2.following_user_id AND F2.follower_user_id = $2
  WHERE
    (A.email ILIKE '%' || $1 || '%' OR
    A.username ILIKE '%' || $1 || '%' OR
    U.first_name ILIKE '%' || $1 || '%' OR
    U.last_name ILIKE '%' || $1 || '%')
  
  UNION ALL

  SELECT
    'Churches' AS category,
    C.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    NULL AS first_name,
    NULL AS last_name,
    NULL AS user_id,
    CASE WHEN F1.follower_user_id IS NOT NULL THEN true ELSE false END AS is_follower,
    CASE WHEN F2.following_user_id IS NOT NULL THEN true ELSE false END AS is_following
  FROM churches C
  INNER JOIN authentications A ON C.auth_id = A.id
  LEFT JOIN entity_profiles E ON A.id = E.user_id
  LEFT JOIN follow F1 ON A.id = F1.follower_user_id AND F1.following_user_id = $2
  LEFT JOIN follow F2 ON A.id = F2.following_user_id AND F2.follower_user_id = $2
  WHERE
    C.name ILIKE '%' || $1 || '%'

  UNION ALL

  SELECT
    'Posts' AS category,
    P.id AS result_id,
    A.email,
    A.phone,
    A.username,
    E.image_url AS profile_image,
    E.following_count,
    E.followers_count,
    E.entity_type,
    U.first_name,
    U.last_name,
    U.user_id,
    CASE WHEN F1.follower_user_id IS NOT NULL THEN true ELSE false END AS is_follower,
    CASE WHEN F2.following_user_id IS NOT NULL THEN true ELSE false END AS is_following
  FROM posts P
  INNER JOIN authentications A ON P.user_id = A.id
  LEFT JOIN users U ON A.id = U.user_id
  LEFT JOIN entity_profiles E ON U.user_id = E.user_id
  LEFT JOIN follow F1 ON A.id = F1.follower_user_id AND F1.following_user_id = $2
  LEFT JOIN follow F2 ON A.id = F2.following_user_id AND F2.follower_user_id = $2
  WHERE
    P.content ILIKE '%' || $1 || '%'
)
SELECT
  category,
  json_agg(
    json_build_object(
      'result_id', result_id,
      'email', email,
      'phone', phone,
      'username', username,
      'profile_image', profile_image,
      'following_count', following_count,
      'followers_count', followers_count,
      'entity_type', entity_type,
      'first_name', first_name,
      'last_name', last_name,
      'user_id', user_id,
      'is_follower', is_follower,
      'is_following', is_following
    )
  ) AS results
FROM SearchResults
GROUP BY category;
