-- -- name: SuggestUsersToFollowPaginated :many
-- SELECT DISTINCT
--     u.id,
--     u.username,
--     us.first_name,
--     us.last_name,
--     ep.image_url,
--     (SELECT COUNT(*) FROM follow WHERE following_user_id = recommended_users.id) AS follow_count,
--     (SELECT COUNT(*) FROM likes WHERE user_id = recommended_users.id) AS like_count
-- FROM
--     authentications u
-- JOIN entity_profiles ep ON u.id = ep.user_id
-- JOIN users us ON u.id = us.user_id
-- JOIN (
--     SELECT u2.id
--     FROM authentications u1
--     JOIN likes l1 ON u1.id = l1.user_id
--     JOIN likes l2 ON l1.post_id = l2.post_id AND l1.user_id != l2.user_id
--     JOIN authentications u2 ON l2.user_id = u2.id
--     WHERE u1.id = $1
--     UNION
--     SELECT u2.id
--     FROM follow f1
--     JOIN follow f2 ON f1.following_user_id = f2.follower_user_id
--     JOIN authentications u2 ON f2.following_user_id = u2.id
--     WHERE f1.follower_user_id = $1
-- ) recommended_users ON u.id = recommended_users.id
-- JOIN (
--     SELECT user_id, MAX(last_active_at) AS last_active
--     FROM sessions
--     GROUP BY user_id
-- ) recent_sessions ON u.id = recent_sessions.user_id
-- WHERE
--     recent_sessions.last_active > NOW() - INTERVAL '1 DAY' AND
--     u.id != $1
-- ORDER BY
--     (SELECT COUNT(*) FROM follow WHERE following_user_id = recommended_users.id) DESC,
--     (SELECT COUNT(*) FROM likes WHERE user_id = recommended_users.id) DESC,
--     follow_count DESC,
--     like_count DESC
-- LIMIT $2 OFFSET $3;

-- name: SuggestUsersToFollowPaginated :many
SELECT DISTINCT
    u.id,
    u.username,
    u.is_verified,
    us.first_name,
    us.last_name,
    ep.image_url,
    ep.following_count,
    ep.followers_count,
    (SELECT COUNT(*) FROM follow WHERE following_user_id = recommended_users.id) AS follow_count,
    (SELECT COUNT(*) FROM likes WHERE user_id = recommended_users.id) AS like_count,
    CASE 
        WHEN recommended_users.id IS NOT NULL THEN 0 -- Existing recommended users
        ELSE RANDOM() -- Random ordering for other users
    END AS random_order
FROM
    authentications u
JOIN entity_profiles ep ON u.id = ep.user_id
JOIN users us ON u.id = us.user_id
LEFT JOIN (
    SELECT u2.id
    FROM authentications u1
    JOIN likes l1 ON u1.id = l1.user_id
    JOIN likes l2 ON l1.post_id = l2.post_id AND l1.user_id != l2.user_id
    JOIN authentications u2 ON l2.user_id = u2.id
    WHERE u1.id = $1
    UNION
    SELECT u2.id
    FROM follow f1
    JOIN follow f2 ON f1.following_user_id = f2.follower_user_id
    JOIN authentications u2 ON f2.following_user_id = u2.id
    WHERE f1.follower_user_id = $1
) recommended_users ON u.id = recommended_users.id
JOIN (
    SELECT user_id, MAX(last_active_at) AS last_active
    FROM sessions
    GROUP BY user_id
) recent_sessions ON u.id = recent_sessions.user_id
WHERE
    (recent_sessions.last_active > NOW() - INTERVAL '1 DAY' OR recent_sessions.last_active IS NULL) AND
    u.id != $1
ORDER BY
    random_order, -- Random ordering for other users
    (SELECT COUNT(*) FROM follow WHERE following_user_id = recommended_users.id) DESC,
    (SELECT COUNT(*) FROM likes WHERE user_id = recommended_users.id) DESC,
    follow_count DESC,
    like_count DESC
LIMIT $2 OFFSET $3;
