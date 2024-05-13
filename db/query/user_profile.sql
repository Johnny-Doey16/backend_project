-- name: CreateEntityProfile :exec
INSERT INTO entity_profiles (
    user_id, image_url, entity_type, header_image_url, about, website
    )
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserProfileByUID :one
SELECT * FROM entity_profiles WHERE user_id = $1 LIMIT 1;

-- name: UpdateImgEntityProfile :exec
UPDATE entity_profiles SET image_url = $2, about = $3, website = $4, header_image_url = $5 WHERE user_id = $1;

-- name: DeleteUserProfileByID :exec
DELETE FROM entity_profiles WHERE user_id = $1;

-- name: UpdateIncreaseFollowers :exec
UPDATE entity_profiles SET following_count = following_count + 1 WHERE user_id = $1;

-- name: UpdateIncreaseFollowing :exec
UPDATE entity_profiles SET followers_count = followers_count + 1 WHERE user_id = $1;

-- name: DecreaseFollowers :exec
UPDATE entity_profiles SET following_count = following_count - 1 WHERE user_id = $1;

-- name: DecreaseFollowing :exec
UPDATE entity_profiles SET followers_count = followers_count - 1 WHERE user_id = $1;


-- name: GetUserProfile :one
SELECT
  u.id,
  u.username,
  u.email,
  u.phone,
  u.created_at,
  u.user_type,
  CASE
    WHEN u.user_type = 'user' THEN us.first_name
    WHEN u.user_type = 'churchAdmin' THEN c.name
    ELSE NULL
  END AS first_name,
  CASE
    WHEN u.user_type = 'user' THEN us.last_name
    ELSE NULL
  END AS last_name,
  ep.image_url,
  u.is_verified,
  ep.following_count,
  ep.followers_count,
  ep.posts_count,
  ep.header_image_url,
  ep.website,
  ep.about,
  EXISTS (
    SELECT 1
    FROM follow f
    WHERE f.follower_user_id = $1 AND f.following_user_id = u.id
  ) AS is_following,
  EXISTS (
    SELECT 1
    FROM follow f
    WHERE f.follower_user_id = u.id AND f.following_user_id = $1
  ) AS is_followed,
  EXISTS (
        SELECT 1
        FROM user_church_membership ucm
        WHERE ucm.user_id = $1 AND ucm.church_id = c.id
    ) AS is_member,
  -- Include church location details for churchAdmin type
  cl.address AS church_address,
  cl.city AS church_city,
  cl.postalCode AS church_postalCode,
  cl.lga AS church_lga,
  cl.state AS church_state,
  cl.country AS church_country,
  cl.location AS church_location,
  c.id AS church_id,
  c.denomination_id AS church_denomination_id,
  c.members_count As church_members_count,
  -- ! Added
  accounts.account_name,
  accounts.account_number,
  accounts.bank_name
FROM
  authentications u
JOIN
  entity_profiles ep ON u.id = ep.user_id
  -- ! Added
LEFT JOIN
  accounts ON u.id = accounts.user_id
LEFT JOIN
  users us ON u.id = us.user_id AND u.user_type = 'user'
LEFT JOIN
  churches c ON u.id = c.auth_id AND u.user_type = 'churchAdmin'
-- Joining with church_locations to get church location details
LEFT JOIN
  church_locations cl ON c.auth_id = cl.auth_id
WHERE
--   u.id = $2 AND
  (u.username = $2 OR u.id::text = $2) AND
  NOT EXISTS (
    SELECT 1
    FROM blocked_users bu
    WHERE bu.blocking_user_id = u.id AND bu.blocked_user_id = $1
  );







-- TODO: Remove too
-- name: GetUserProfileByUsername :one
SELECT
  u.id,
  u.username,
  u.email,
  u.phone,
  u.created_at,
  u.user_type,
  CASE
    WHEN u.user_type = 'user' THEN us.first_name
    WHEN u.user_type = 'churchAdmin' THEN c.name
    ELSE NULL
  END AS first_name,
  CASE
    WHEN u.user_type = 'user' THEN us.last_name
    ELSE NULL
  END AS last_name,
  ep.image_url,
  u.is_verified,
  ep.following_count,
  ep.followers_count,
  ep.posts_count,
  ep.header_image_url,
  ep.website,
  ep.about,
  EXISTS (
    SELECT 1
    FROM follow f
    WHERE f.follower_user_id = $1 AND f.following_user_id = u.id
  ) AS is_following,
  EXISTS (
    SELECT 1
    FROM follow f
    WHERE f.follower_user_id = u.id AND f.following_user_id = $1
  ) AS is_followed,
  EXISTS (
        SELECT 1
        FROM user_church_membership ucm
        WHERE ucm.user_id = $1 AND ucm.church_id = c.id
    ) AS is_member,
  cl.address AS church_address,
  cl.city AS church_city,
  cl.postalCode AS church_postalCode,
  cl.lga AS church_lga,
  cl.state AS church_state,
  cl.country AS church_country,
  cl.location AS church_location,
  c.id AS church_id,
  c.denomination_id AS church_denomination_id,
  c.members_count As church_members_count
FROM
  authentications u
JOIN
  entity_profiles ep ON u.id = ep.user_id
LEFT JOIN
  users us ON u.id = us.user_id AND u.user_type = 'user'
LEFT JOIN
  churches c ON u.id = c.auth_id AND u.user_type = 'churchAdmin'
-- Joining with church_locations to get church location details
LEFT JOIN
  church_locations cl ON c.auth_id = cl.auth_id
WHERE
  u.username = $2 AND
  NOT EXISTS (
    SELECT 1
    FROM blocked_users bu
    WHERE bu.blocking_user_id = u.id AND bu.blocked_user_id = $1
  );



















-- TODO: Old Remove down
-- name: GetUserProfileOld :one
SELECT
    u.id,
    u.username,
    u.email,
    u.phone,
    u.created_at,
    CASE
        WHEN u.user_type = 'user' THEN us.first_name
        WHEN u.user_type = 'churchAdmin' THEN c.name
        ELSE NULL
    END AS first_name,
    CASE
        WHEN u.user_type = 'user' THEN us.last_name
        ELSE NULL
    END AS last_name,
    ep.image_url,
    u.is_verified,
    ep.following_count,
    ep.followers_count,

    ep.posts_count,
    ep.header_image_url,
    ep.website,
    ep.about,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = $1 AND f.following_user_id = u.id
    ) AS is_following,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = u.id AND f.following_user_id = $1
    ) AS is_followed
FROM
    authentications u
JOIN
    entity_profiles ep ON u.id = ep.user_id
LEFT JOIN
    users us ON u.id = us.user_id AND u.user_type = 'user'
LEFT JOIN
    churches c ON u.id = c.auth_id AND u.user_type = 'churchAdmin'
WHERE
    u.id = $2 AND
    NOT EXISTS (
        SELECT 1
        FROM blocked_users bu
        WHERE bu.blocking_user_id = u.id AND bu.blocked_user_id = $1
    );

-- name: GetUserProfileByUsernameOld :one
SELECT
    u.id,
    u.username,
    u.email,
    u.phone,
    u.created_at,
    CASE
        WHEN u.user_type = 'user' THEN us.first_name
        WHEN u.user_type = 'churchAdmin' THEN c.name
        ELSE NULL
    END AS first_name,
    CASE
        WHEN u.user_type = 'user' THEN us.last_name
        ELSE NULL
    END AS last_name,
    ep.image_url,
    u.is_verified,
    ep.following_count,
    ep.followers_count,

    ep.posts_count,
    ep.header_image_url,
    ep.website,
    ep.about,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = $1 AND f.following_user_id = u.id
    ) AS is_following,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = u.id AND f.following_user_id = $1
    ) AS is_followed
FROM
    authentications u
JOIN
    entity_profiles ep ON u.id = ep.user_id
LEFT JOIN
    users us ON u.id = us.user_id AND u.user_type = 'user'
LEFT JOIN
    churches c ON u.id = c.auth_id AND u.user_type = 'churchAdmin'
WHERE
    u.username = $2 AND
    NOT EXISTS (
        SELECT 1
        FROM blocked_users bu
        WHERE bu.blocking_user_id = u.id AND bu.blocked_user_id = $1
    );
