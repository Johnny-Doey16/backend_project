-- name: GetPostsRecommendationDisFunc :many
WITH UserActivityRecommendations AS (
    SELECT p.*, ARRAY_AGG(pi.image_url) AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    WHERE p.user_id = $1 OR p.id IN (SELECT post_id FROM likes WHERE user_id = $1)
    GROUP BY p.id
    ORDER BY p.created_at DESC
    LIMIT 3
),

-- Recommendations based on popularity
PopularPostsRecommendations AS (
    SELECT p.*, ARRAY_AGG(pi.image_url) AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN posts_metrics pm ON p.id = pm.post_id
    GROUP BY p.id, pm.likes
    ORDER BY pm.likes DESC
    LIMIT 3
),

-- Recommendations from followed users
FollowedUsersRecommendations AS (
    SELECT p.*, ARRAY_AGG(pi.image_url) AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN follow f ON p.user_id = f.following_user_id
    WHERE f.follower_user_id = $1
    GROUP BY p.id
    ORDER BY p.created_at DESC
    LIMIT 2
),

-- Collaborative Filtering Recommendations
CollaborativeFilteringRecommendations AS (
    SELECT p.*, ARRAY_AGG(pi.image_url) AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
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
    GROUP BY p.id
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
)
SELECT * FROM CombinedRecommendations
ORDER BY created_at DESC
LIMIT 10;


-- name: GetPostsRecommendationNewWrong :many
WITH UserActivityRecommendations AS (
    SELECT p.*, COALESCE(ARRAY_AGG(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    WHERE p.user_id = $1 OR p.id IN (SELECT post_id FROM likes WHERE user_id = $1)
    GROUP BY p.id
    ORDER BY p.created_at DESC
    LIMIT 3
),

-- Recommendations based on popularity
PopularPostsRecommendations AS (
    SELECT p.*, COALESCE(ARRAY_AGG(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN posts_metrics pm ON p.id = pm.post_id
    GROUP BY p.id, pm.likes
    ORDER BY pm.likes DESC
    LIMIT 3
),

-- Recommendations from followed users
FollowedUsersRecommendations AS (
    SELECT p.*, COALESCE(ARRAY_AGG(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN follow f ON p.user_id = f.following_user_id
    WHERE f.follower_user_id = $1
    GROUP BY p.id
    ORDER BY p.created_at DESC
    LIMIT 2
),

-- Collaborative Filtering Recommendations
CollaborativeFilteringRecommendations AS (
    SELECT p.*, COALESCE(ARRAY_AGG(pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL), '{}') AS image_urls
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
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
    GROUP BY p.id
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
)
SELECT * FROM CombinedRecommendations
ORDER BY created_at DESC
LIMIT 10;