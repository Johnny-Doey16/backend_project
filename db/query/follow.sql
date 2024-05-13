-- Example: User with user_id 1 follows user with user_id 2

-- name: CreateFollow :exec
INSERT INTO follow (
    follower_user_id, following_user_id, created_at
    )
VALUES ($1, $2, now());

-- name: GetFollowers :many
SELECT
    follow.*,
    authentications.id,
    authentications.username,
    authentications.is_verified,
    authentications.created_at,
    users.first_name,
    entity_profiles.image_url
FROM
    follow
JOIN
    authentications ON follow.following_user_id = authentications.id
JOIN
    users ON follow.following_user_id = users.user_id
JOIN
    entity_profiles ON follow.following_user_id = entity_profiles.user_id
WHERE
    follow.following_user_id = $1;


-- name: GetFollowing :many
SELECT
    follow.*,
    authentications.id,
    authentications.username,
    authentications.is_verified,
    authentications.created_at,
    users.first_name,
    entity_profiles.image_url
FROM
    follow
JOIN
    authentications ON follow.following_user_id = authentications.id
JOIN
    users ON follow.following_user_id = users.user_id
JOIN
    entity_profiles ON follow.following_user_id = entity_profiles.user_id
WHERE
    follow.follower_user_id = $1;

-- name: UnFollow :exec
DELETE FROM follow WHERE follower_user_id = $1 AND following_user_id = $2;


-- name: GetMutualFollowers :many
SELECT u.username, u.id, e.image_url
FROM authentications u
JOIN follow f1 ON u.id = f1.following_user_id
JOIN follow f2 ON u.id = f2.following_user_id
JOIN entity_profiles e ON u.id = e.user_id
WHERE f1.follower_user_id = $1
AND f2.follower_user_id = $2;



--  GetMutualFollowers - problematic
-- SELECT u.*
-- FROM authentications u
-- JOIN follow f1 ON u.id = f1.following_user_id -- Users A is following
-- JOIN follow f2 ON f1.following_user_id = f2.follower_user_id -- Who are also following B
-- WHERE f1.follower_user_id = 'e7679d8b-0eac-4ea2-93cd-0018ab995922' -- A is the current user
-- AND f2.following_user_id = '558c60aa-977f-4d38-885b-e813656371ac'; -- B is the user A is viewing