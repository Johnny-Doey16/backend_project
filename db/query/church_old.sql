-- name CheckDenominationLastChanged one
-- SELECT last_denomination_change FROM users WHERE user_id = $1;

-- name CheckChurchLastChanged one
-- SELECT last_church_change FROM users WHERE user_id = $1;

-- -- name: UpdateDenominationForUser :exec
-- UPDATE users
-- SET denomination_id = $1, last_denomination_change = NOW()
-- WHERE id = $2
--   AND (last_denomination_change IS NULL OR last_denomination_change < NOW() - INTERVAL '1 year');

-- -- name: UpdateChurchForUser :exec
-- UPDATE users
-- SET church_id = $1, last_church_change = NOW()
-- WHERE id = $2
--   AND (last_church_change IS NULL OR last_church_change < NOW() - INTERVAL '6 months');

-- -- name: CreateDenominationForUser :exec
-- INSERT INTO user_denomination_membership (user_id, denomination_id, join_date)
-- VALUES ($1, $2, now());

-- -- name: GetUserDenominationMembership :one
-- SELECT * FROM user_denomination_membership WHERE user_id = $1, denomination_id = $2;

  
-- -- name: LeaveDenomination :exec
-- UPDATE user_denomination_membership
-- SET join_date = now()
-- WHERE id = $1
--   AND (join_date < NOW() - INTERVAL '1 year');




-- !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

-- -- Join a Church
-- -- :user_id is the user's ID, :church_id is the church's ID
-- INSERT INTO memberships (user_id, church_id, join_date)
-- VALUES (:user_id, :church_id, NOW())
-- RETURNING id;

-- -- Leave a Church
-- -- :user_id is the user's ID, :church_id is the church's ID
-- DELETE FROM memberships
-- WHERE user_id = :user_id AND church_id = :church_id;


-- -- Update user's church information
-- UPDATE users
-- SET current_church_id = $1, last_church_change = NOW()
-- WHERE id = $2
--   AND (last_church_change IS NULL OR last_church_change < NOW() - INTERVAL '6 months');

-- -- Set the user's current church to NULL
-- UPDATE users
-- SET current_church_id = NULL, last_church_change = NOW()
-- WHERE id = $1;

-- -- Update user's denomination information
-- UPDATE users
-- SET current_denomination_id = $1, last_denomination_change = NOW()
-- WHERE id = $2
--   AND (last_denomination_change IS NULL OR last_denomination_change < NOW() - INTERVAL '1 year');

-- -- Set the user's current denomination to NULL
-- UPDATE users
-- SET current_denomination_id = NULL, last_denomination_change = NOW()
-- WHERE id = $1;

-- -- Search for churches by name or denomination, prioritizing user's denomination
-- SELECT c.* FROM churches c
-- JOIN denominations d ON c.denomination_id = d.id
-- WHERE LOWER(c.name) LIKE LOWER($1) OR d.id = $2
-- ORDER BY (d.id = $3) DESC;

-- -- Search churches nearby
-- SELECT id, name, vicar, latitude, longitude,
--     earth_distance(
--         ll_to_earth($1, $2),
--         ll_to_earth(latitude, longitude)
--     ) as distance
-- FROM churches
-- ORDER BY distance ASC
-- LIMIT 10;