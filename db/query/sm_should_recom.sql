-- name: IsPostRecommended :one
SELECT EXISTS (
    SELECT 1
    FROM (
        -- User Activity (Recent Interactions) Recommendations
        SELECT p.id
        FROM posts p
        JOIN (
            SELECT post_id FROM likes WHERE likes.user_id = $1
            UNION
            SELECT original_post_id AS post_id FROM repost WHERE user_id = $1
            UNION
            SELECT post_id FROM post_mentions WHERE mentioned_user_id = $1
        ) AS recent_activity ON p.id = recent_activity.post_id

        UNION

        -- Followed Users Recommendations
        SELECT p.id
        FROM posts p
        JOIN follow f ON p.user_id = f.following_user_id
        WHERE f.follower_user_id = $1

        UNION

        -- Popularity-Based Recommendations
        SELECT p.id
        FROM posts p
        JOIN posts_metrics pm ON p.id = pm.post_id

        UNION

        -- Collaborative Filtering Recommendations
        SELECT p.id
        FROM posts p
        JOIN likes l ON p.id = l.post_id
        WHERE l.user_id IN (
            SELECT following_user_id FROM follow WHERE follower_user_id = $1
        )
        AND p.id NOT IN (
            SELECT id FROM posts WHERE user_id = $1
        )

        UNION

        -- Content-Based Recommendations (not defined in initial context, so this is a placeholder)
        SELECT p.id
        FROM posts p
        WHERE p.user_id = $1
    ) AS RECOMMENDATIONS
) AS should_send_post;