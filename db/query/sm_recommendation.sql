-- Recommendations based on user activity

-- name: GetPostsRecommendationOLD :many
WITH UserActivityRecommendations AS (
    SELECT p.*
    FROM posts p
    WHERE p.user_id = $1 OR p.id IN (SELECT post_id FROM likes WHERE user_id = $1)
    ORDER BY p.created_at DESC
    LIMIT 3
),

-- Recommendations based on popularity
PopularPostsRecommendations AS (
    SELECT p.*
    FROM posts p
    JOIN posts_metrics pm ON p.id = pm.post_id
    ORDER BY pm.likes DESC
    LIMIT 3
),

-- Recommendations from followed users
FollowedUsersRecommendations AS (
    SELECT p.*
    FROM posts p
    JOIN follow f ON p.user_id = f.following_user_id
    WHERE f.follower_user_id = $1
    ORDER BY p.created_at DESC
    LIMIT 2
),

-- Collaborative Filtering Recommendations
CollaborativeFilteringRecommendations AS (
    SELECT p.*
    FROM posts p
    JOIN likes l ON p.id = l.post_id
    WHERE l.user_id IN (
        SELECT following_user_id
        FROM follow
        WHERE follower_user_id = $1
    )
    AND p.id NOT IN (
        -- Exclude the user's own posts
        SELECT id
        FROM posts
        WHERE user_id = $1
    )
    ORDER BY p.created_at DESC
    LIMIT 2
),

-- Combine all recommendations
CombinedRecommendations AS (
    SELECT * FROM UserActivityRecommendations
    UNION
    SELECT * FROM PopularPostsRecommendations
    UNION
    SELECT * FROM FollowedUsersRecommendations
    UNION
    SELECT * FROM CollaborativeFilteringRecommendations
),


-- SELECT * FROM CombinedRecommendations
-- ORDER BY created_at DESC
-- LIMIT 10;



POSTS_WITH_IMAGES AS (
    SELECT
        r.*,
        json_agg(pi.image_url) AS image_urls
    FROM
        CombinedRecommendations r
        LEFT JOIN posts_images pi ON r.id = pi.post_id
    GROUP BY
        r.id, r.content, r.created_at, r.user_id, r.total_images, r.updated_at, r.deleted_at
)

SELECT * FROM POSTS_WITH_IMAGES
ORDER BY created_at DESC;