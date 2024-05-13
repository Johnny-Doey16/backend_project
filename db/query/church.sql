-- name: CreateNewChurch :exec
INSERT INTO churches (auth_id, denomination_id, name)
VALUES ($1, $2, $3);

-- name: CreateNewChurchLocation :exec
INSERT INTO church_locations (auth_id, address, city, postalCode, state, country, location, lga)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);


-- *** *** *** *** *** *** *** ***


-- name: GetUserAndChurchMembership :one
SELECT u.*, m.*
FROM users u
LEFT JOIN user_church_membership m ON u.user_id = m.user_id AND m.active
WHERE u.user_id = $1
ORDER BY m.join_date DESC
LIMIT 1;


-- name: UpdateChurchForUser :one
WITH valid_user AS (
  SELECT user_id
  FROM users u
  JOIN churches c ON u.denomination_id = c.denomination_id
  WHERE u.user_id = $2
    AND u.denomination_id IS NOT NULL
    AND c.id = $1
),
deactivated AS (
  UPDATE user_church_membership
  SET active = FALSE, leave_date = NOW()
  WHERE user_id IN (SELECT user_id FROM valid_user) AND active
  RETURNING user_id
),
updated_users AS (
  UPDATE users
  SET church_id = $1, last_church_change = NOW()
  FROM deactivated d
  WHERE users.user_id IN (SELECT user_id FROM valid_user)
    AND users.church_id IS DISTINCT FROM $1
    AND (users.last_church_change IS NULL OR users.last_church_change < (NOW() - INTERVAL '6 months'))
  RETURNING id
)
INSERT INTO user_church_membership (user_id, church_id, join_date, active)
SELECT $2, $1, NOW(), TRUE
FROM updated_users
RETURNING *;


-- name: UpdateChurchForUserOld :one
WITH valid_user AS (
  SELECT user_id
  FROM users
  JOIN churches ON users.denomination_id = churches.denomination_id
  WHERE users.user_id = $2
    AND users.denomination_id IS NOT NULL
    AND churches.id = $1
),
deactivated AS (
  UPDATE user_church_membership
  SET active = FALSE, leave_date = NOW()
  WHERE user_id IN (SELECT user_id FROM valid_user) AND active
  RETURNING user_id
),
updated_users AS (
  UPDATE users
  SET church_id = $1, last_church_change = NOW()
  FROM deactivated
  WHERE users.user_id IN (SELECT user_id FROM valid_user)
    AND users.church_id IS DISTINCT FROM $1
    AND (users.last_church_change IS NULL OR users.last_church_change < NOW() - INTERVAL '6 months')
  RETURNING id
)
INSERT INTO user_church_membership (user_id, church_id, join_date, active)
SELECT $2, $1, NOW(), TRUE
FROM updated_users
RETURNING *;



-- name: GetUserChurchMembership :one
SELECT * FROM user_church_membership
WHERE user_id = $1 AND church_id = $2 AND active
ORDER BY join_date DESC
LIMIT 1;



-- name: LeaveChurch :exec
UPDATE user_church_membership
SET active = FALSE, leave_date = NOW()
WHERE user_id = $1 AND active
  AND (join_date < NOW() - INTERVAL '6 month');



-- name: CreateChurchForUser :one
WITH valid_user AS (
  SELECT user_id
  FROM users
  JOIN churches ON users.denomination_id = churches.denomination_id
  WHERE user_id = $1
    AND users.denomination_id IS NOT NULL
    AND churches.id = $2
),
updated_users AS (
  UPDATE users
  SET church_id = $2, last_church_change = NOW()
  WHERE user_id IN (SELECT user_id FROM valid_user)
    AND users.church_id IS DISTINCT FROM $2
    AND (users.last_church_change IS NULL OR users.last_church_change < NOW() - INTERVAL '6 month')
  RETURNING id
)
INSERT INTO user_church_membership (user_id, church_id, join_date, active)
SELECT $1, $2, NOW(), TRUE
FROM updated_users
WHERE NOT EXISTS (
  SELECT 1 FROM user_church_membership
  WHERE user_id = $1 AND church_id = $2 AND active
)
AND (
  SELECT COUNT(*)
  FROM user_church_membership
  WHERE user_id = $1 AND active AND join_date >= NOW() - INTERVAL '6 month'
) = 0
RETURNING *;

-- name: SearchChurches :many
SELECT
  c.id,
  c.name,
  c.denomination_id,
  cl.address,
  cl.city,
  cl.postalCode,
  cl.state,
  cl.country,
  cl.location,
  cl.lga
FROM
  churches c
JOIN
  church_locations cl ON c.auth_id = cl.auth_id
WHERE
  c.name ILIKE '%' || $1 || '%'
  OR CAST(c.denomination_id AS TEXT) ILIKE '%' || $1 || '%'
  -- ORDER BY bookmarked_at DESC
  LIMIT $2 OFFSET $3;
-- Above can also search by denomination id


-- name: GetNearbyChurches :many
SELECT
  c.id,
  c.auth_id,
  c.name,
  c.denomination_id,
  cl.address,
  cl.city,
  cl.postalCode,
  cl.state,
  cl.country,
  cl.location,
  cl.lga,
  ST_Distance(cl.location, ST_MakePoint($1, $2)::geography) AS distance,
  a.username,
  ep.image_url
FROM
  churches c
JOIN
  church_locations cl ON c.auth_id = cl.auth_id
JOIN
  authentications a ON c.auth_id = a.id
JOIN
  entity_profiles ep ON a.id = ep.user_id
WHERE
  ST_DWithin(
    cl.location,
    ST_MakePoint($1, $2)::geography,
    $3
)
ORDER BY
  ST_Distance(cl.location, ST_MakePoint($1, $2)::geography);

-- name: GetChurchProfile :one
SELECT
  c.id,
  c.auth_id,
  c.name,
  c.members_count,
  c.denomination_id,
  cl.address,
  cl.city,
  cl.postalCode,
  cl.state,
  cl.country,
  cl.location,
  cl.lga,
  a.username,
  a.phone,
  a.is_verified,
  ep.image_url,
  ep.header_image_url,
  ep.posts_count,
  ep.website,
  ep.about,
  ep.following_count,
  ep.followers_count,
  EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = $2 AND f.following_user_id = c.auth_id
    ) AS is_following,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = c.auth_id AND f.following_user_id = $2
    ) AS is_followed,
    EXISTS (
        SELECT 1
        FROM user_church_membership ucm
        WHERE ucm.user_id = $2 AND ucm.church_id = c.id
    ) AS is_member
  
FROM
  churches c
JOIN
  church_locations cl ON c.auth_id = cl.auth_id
JOIN
  authentications a ON c.auth_id = a.id
JOIN
  entity_profiles ep ON a.id = ep.user_id
WHERE
  c.auth_id = $1;

-- name: GetUserChurch :one
WITH church_membership AS (
  SELECT
    ucm.church_id,
    c.denomination_id,
    c.members_count,
    c.name AS church_name,
    c.auth_id AS church_auth_id,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = $1 AND f.following_user_id = c.auth_id
    ) AS is_following,
    EXISTS (
        SELECT 1
        FROM follow f
        WHERE f.follower_user_id = c.auth_id AND f.following_user_id = $1
    ) AS is_followed
  FROM
    user_church_membership ucm
  JOIN
    churches c ON ucm.church_id = c.id
  WHERE
    ucm.user_id = $1 -- Replace $4 with the actual user_id
    AND ucm.active = true
)

SELECT
  cm.church_id,
  cm.denomination_id,
  cm.church_name,
  cm.members_count,
  cm.church_auth_id,
  cm.is_followed,
  cm.is_following,
  a.email,
  a.username,
  a.is_verified,
  a.phone,
  ep.image_url,
  ep.following_count,
  ep.followers_count,
  ep.header_image_url,
  ep.posts_count,
  ep.website,
  ep.about,

  cl.address,
  cl.city,
  cl.postalcode,
  cl.state,
  cl.country,
  cl.lga,
  
  acc.account_name,
  acc.account_number,
  acc.bank_name,
  acc.total_coin
FROM
  church_membership cm
JOIN
  authentications a ON cm.church_auth_id = a.id
JOIN
  entity_profiles ep ON cm.church_auth_id = ep.user_id
JOIN
  church_locations cl ON cm.church_auth_id = cl.auth_id
LEFT JOIN
  accounts acc ON a.id = acc.user_id
WHERE
  a.is_suspended = false
  AND a.is_deleted = false;

  -- AND a.is_verified = true;


 
-- ! TODO: Add Order by join date
-- name: GetChurchMembers5 :many
WITH user_details AS (
    SELECT
        u.id AS user_id,
        u.first_name,
        a.id AS auth_id,
        a.username,
        a.is_verified,
        ep.image_url,
        ucm.join_date
    FROM user_church_membership ucm
    JOIN users u ON ucm.user_id = u.user_id
    JOIN authentications a ON u.user_id = a.id
    LEFT JOIN entity_profiles ep ON a.id = ep.user_id
    WHERE ucm.church_id = $1 AND ucm.active = true AND a.id <> $2
)
SELECT
    jsonb_agg(user_details) AS user_details,
    EXISTS (
        SELECT 1
        FROM user_church_membership ucm
        WHERE ucm.church_id = $1 AND ucm.user_id = $2 AND ucm.active = true
    ) AS is_member
FROM user_details
-- GROUP BY join_date
-- ORDER BY join_date
OFFSET $3
LIMIT $4;


-- name: GetChurchMembersUid :many
SELECT uc.user_id
FROM user_church_membership uc
JOIN churches c ON uc.church_id = c.id
WHERE c.auth_id = $1;
