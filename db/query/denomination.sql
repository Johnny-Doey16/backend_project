-- name: CreateDenomination :exec
INSERT INTO "denominations" ("name") VALUES ($1);

-- name: GetDenominationList :many
SELECT id, name FROM denominations;


-- name: GetUserAndMembership :one
SELECT u.*, m.*
FROM users u
LEFT JOIN user_denomination_membership m ON u.user_id = m.user_id AND m.active
WHERE u.user_id = $1
ORDER BY m.join_date DESC
LIMIT 1;


-- name: UpdateDenominationForUser :exec
WITH deactivated AS (
  UPDATE user_denomination_membership
  SET active = FALSE, leave_date = NOW()
  WHERE user_denomination_membership.user_id = $2 AND active
  RETURNING user_denomination_membership.user_id, denomination_id, join_date
),
can_change AS (
  SELECT user_id, denomination_id, join_date
  FROM deactivated
  WHERE join_date <= NOW() - INTERVAL '1 year'
),
updated_users AS (
  UPDATE users
  SET denomination_id = $1, last_denomination_change = NOW()
  FROM can_change
  WHERE users.user_id = can_change.user_id
    AND users.denomination_id != can_change.denomination_id
  RETURNING users.user_id
)
INSERT INTO user_denomination_membership (user_id, denomination_id, join_date, active)
SELECT user_id, $1, NOW(), TRUE
FROM can_change
WHERE NOT EXISTS (
  SELECT 1
  FROM updated_users
  WHERE updated_users.user_id = can_change.user_id
);



-- name: UpdateDenominationForUserOld :exec
WITH deactivated AS (
  UPDATE user_denomination_membership
  SET active = FALSE, leave_date = NOW()
  WHERE user_id = $2 AND active
  RETURNING user_id
),
updated_users AS (
  UPDATE users
  SET denomination_id = $1, last_denomination_change = NOW()
  FROM deactivated
  WHERE users.user_id = $2
    AND users.denomination_id IS DISTINCT FROM $1
    AND (users.last_denomination_change IS NULL OR users.last_denomination_change < NOW() - INTERVAL '1 year')
  RETURNING id
)
INSERT INTO user_denomination_membership (user_id, denomination_id, join_date, active)
SELECT $2, $1, NOW(), TRUE
FROM updated_users;



-- name: GetUserDenominationMembership :one
SELECT * FROM user_denomination_membership
WHERE user_id = $1 AND denomination_id = $2 AND active
ORDER BY join_date DESC
LIMIT 1;



-- name: LeaveDenomination :exec
UPDATE user_denomination_membership
SET active = FALSE, leave_date = NOW()
WHERE user_id = $1 AND active
  AND (join_date < NOW() - INTERVAL '1 year');



-- name: CreateDenominationForUser :exec
INSERT INTO "user_denomination_membership" (user_id, denomination_id, join_date) VALUES ($1, $2, NOW());


-- name: CreateDenominationForUserOld :exec
WITH updated_users AS (
  UPDATE users
  SET denomination_id = $2, last_denomination_change = NOW()
  WHERE user_id = $1
    AND users.denomination_id IS DISTINCT FROM $2
    AND (users.last_denomination_change IS NULL OR users.last_denomination_change < NOW() - INTERVAL '1 year')
  RETURNING id
)
INSERT INTO user_denomination_membership (user_id, denomination_id, join_date, active)
SELECT $1, $2, NOW(), TRUE
FROM updated_users
WHERE NOT EXISTS (
  SELECT 1 FROM user_denomination_membership
  WHERE user_id = $1 AND denomination_id = $2 AND active
)
AND (
  SELECT COUNT(*)
  FROM user_denomination_membership
  WHERE user_id = $1 AND active AND join_date >= NOW() - INTERVAL '1 year'
) = 0;