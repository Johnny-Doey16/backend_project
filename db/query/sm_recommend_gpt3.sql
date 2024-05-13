-- name: GetPostsRecommendationFUCKED :many

WITH UserActivityRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    WHERE p.user_id = $1 OR p.id IN (SELECT post_id FROM likes WHERE user_id = $1)
    ORDER BY p.created_at DESC
    LIMIT 3
),

PopularPostsRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN posts_metrics pm ON p.id = pm.post_id
    ORDER BY pm.likes DESC
    LIMIT 3
),

FollowedUsersRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN follow f ON p.user_id = f.following_user_id
    WHERE f.follower_user_id = $1
    ORDER BY p.created_at DESC
    LIMIT 2
),

CollaborativeFilteringRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN likes l ON p.id = l.post_id
    WHERE l.user_id IN (
        SELECT following_user_id
        FROM follow
        WHERE follower_user_id = $1
    )
    AND p.id NOT IN (
        SELECT id
        FROM posts
        WHERE user_id = $1
    )
    ORDER BY p.created_at DESC
    LIMIT 2
),

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


-- name: GetPostsRecommendationCheck :many

WITH UserActivityRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    WHERE p.user_id = $1 OR p.id IN (SELECT post_id FROM likes WHERE user_id = $1)
    ORDER BY p.created_at DESC
    LIMIT 3
),

PopularPostsRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN posts_metrics pm ON p.id = pm.post_id
    ORDER BY pm.likes DESC
    LIMIT 3
),

FollowedUsersRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN follow f ON p.user_id = f.following_user_id
    WHERE f.follower_user_id = $1
    ORDER BY p.created_at DESC
    LIMIT 2
),

CollaborativeFilteringRecommendations AS (
    SELECT p.*, pi.*
    FROM posts p
    LEFT JOIN posts_images pi ON p.id = pi.post_id
    JOIN likes l ON p.id = l.post_id
    WHERE l.user_id IN (
        SELECT following_user_id
        FROM follow
        WHERE follower_user_id = $1
    )
    AND p.id NOT IN (
        SELECT id
        FROM posts
        WHERE user_id = $1
    )
    ORDER BY p.created_at DESC
    LIMIT 2
),

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
