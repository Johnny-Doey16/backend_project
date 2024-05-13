-- name: GetPostsRecommendation :many
WITH USER_ENTITY AS (
  SELECT
    au.id AS auth_id,
    au.username,
    au.is_verified,
    CASE
      WHEN au.user_type = 'user' THEN 'user'
      WHEN au.user_type = 'churchAdmin' THEN 'churchAdmin'
    END AS entity_type,
    CASE
      WHEN au.user_type = 'user' THEN us.first_name
      WHEN au.user_type = 'churchAdmin' THEN c.name
    END AS name,
    ep.image_url AS user_image_url,
    ep.following_count,
    ep.followers_count,
    CASE
      WHEN au.user_type = 'user' THEN NULL
      WHEN au.user_type = 'churchAdmin' THEN c.members_count
    END AS members_count
  FROM authentications au
  LEFT JOIN users us ON au.id = us.user_id AND au.user_type = 'user'
  LEFT JOIN churches c ON au.id = c.auth_id AND au.user_type = 'churchAdmin'
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE au.is_deleted = FALSE
    AND au.is_suspended = FALSE
),
RECOMMENDATIONS AS (
  -- ! Recent interactions Get posts the user liked, repost or was mentioned in
  SELECT
    p.*,
    0 AS popularity_rank,
    'recent_interactions' AS recommendation_type,
    'Because you recently interacted with this post' AS reason,
    NULL AS rUsername,
    NULL AS title
  FROM posts p
  INNER JOIN (
    SELECT post_id
    FROM likes
    WHERE likes.user_id = $1
    UNION
    SELECT original_post_id AS post_id
    FROM repost
    WHERE user_id = $1
    UNION
    SELECT post_id
    FROM post_mentions
    WHERE mentioned_user_id = $1
  ) AS recent_activity ON p.id = recent_activity.post_id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = '558c60aa-977f-4d38-885b-e813656371ac' AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = '558c60aa-977f-4d38-885b-e813656371ac')
  WHERE p.deleted_at IS NULL
    AND p.suspended_at IS NULL
    AND bu.id IS NULL

  UNION ALL

  -- ! Follow-Based Recommendations: Recommend posts from users that the current user is following but has not interacted with recently.
  SELECT
    p.*,
    0 AS popularity_rank,
    'follow_based' AS recommendation_type,
    'Because you follow @' || ue.username AS reason,
    ue.username AS rUsername,
    NULL AS title
  FROM posts p
  JOIN follow f ON p.user_id = f.following_user_id
  JOIN USER_ENTITY ue ON p.user_id = ue.auth_id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = '558c60aa-977f-4d38-885b-e813656371ac' AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = '558c60aa-977f-4d38-885b-e813656371ac')
  WHERE f.follower_user_id = $1
    AND p.deleted_at IS NULL
    AND p.suspended_at IS NULL
    AND bu.id IS NULL
    -- AND p.id NOT IN (
    --   SELECT post_id FROM likes WHERE user_id = $1
    -- )
    -- AND p.id NOT IN (
    --   SELECT post_id FROM post_comments WHERE user_id = $1
    -- )

  
  UNION ALL

  -- TODO: Edit to also include the user's location
  -- ! Popularity based/Trending Recommendations: Recommend posts that are currently trending or popular among all users, based on recent engagement metrics.
  SELECT
    p.*,
    'popularity_based' AS recommendation_type,
    'This post is popular among users' AS reason,
    NULL AS rUsername,
    NULL AS title
  FROM (
    SELECT
      p.*,
      ROW_NUMBER() OVER (
        ORDER BY (
          (pm.likes * 0.5) + (pm.reposts * 0.5) + (pm.comments * 0.5) + (EXTRACT(EPOCH FROM NOW() - p.created_at) / 3600 / 24)
        ) DESC
      ) AS popularity_rank
    FROM posts p
    JOIN posts_metrics pm ON p.id = pm.post_id
    WHERE p.created_at >= NOW() - INTERVAL '7 days'
  ) p
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $1 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $1)
  WHERE p.popularity_rank <= 5
    AND p.deleted_at IS NULL
    AND p.suspended_at IS NULL
    AND bu.id IS NULL

  UNION ALL

  -- ! Collaborative Filtering: Get posts that the user's followers liked
  -- TODO: Recommend posts liked by users who have similar tastes as the current user, based on the posts they have both liked.
  SELECT
    p.*,
    0 AS popularity_rank,
    'collaborative_filtering' AS recommendation_type,
    'Because ' || ue.username || ' likes this post' AS reason,
    NULL AS rUsername,
    NULL AS title
  FROM posts p
  JOIN likes l ON p.id = l.post_id
  JOIN USER_ENTITY ue ON l.user_id = ue.auth_id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $1 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $1)
  WHERE l.user_id IN (
    SELECT following_user_id
    FROM follow
    WHERE follower_user_id = $1
  )
    AND p.deleted_at IS NULL
    AND p.suspended_at IS NULL
    AND bu.id IS NULL
    AND p.id NOT IN (
      SELECT id
      FROM posts
      WHERE user_id = $1
    )

  UNION ALL

  -- ! Church Announcements: Get announcements by user's church
  SELECT
    a.id,
    a.user_id,
    a.content,
    a.total_images,
    a.created_at,
    a.updated_at,
    a.suspended_at,
    a.deleted_at,
    0 AS popularity_rank,
    'church_announcement' AS recommendation_type,
    'Announcement from ' || c.name AS reason,
    NULL AS rUsername,
    a.title
  FROM announcements a
  JOIN churches c ON a.user_id = c.auth_id
  JOIN user_church_membership ucm ON ucm.church_id = c.id
    AND ucm.user_id = $1
    AND ucm.active = TRUE
    LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $1 AND bu.blocked_user_id = a.user_id)
    OR (bu.blocking_user_id = a.user_id AND bu.blocked_user_id = $1)
  WHERE a.created_at >= NOW() - INTERVAL '2 weeks'
    AND a.deleted_at IS NULL
    AND a.suspended_at IS NULL
    AND bu.id IS NULL

  UNION ALL

  -- ! Social Graph Recommendations: Recommend posts from users who are friends or connections of the user's friends or connections, leveraging the social graph data.
  SELECT
    p.*,
    0 AS popularity_rank,
    'social_graph' AS recommendation_type,
    string_agg(DISTINCT au.username || ' ' || interaction_type, ', ') AS reason,
    NULL AS rUsername,
    NULL AS title
  FROM (
    SELECT
      p.id,
      ROW_NUMBER() OVER (
        ORDER BY (
          COUNT(l.user_id) + COUNT(pc.user_id) + COUNT(pm.mentioned_user_id) + COUNT(v.user_id) + COUNT(r.user_id)
        ) DESC
      ) AS social_graph_rank,
      CASE
        WHEN l.user_id IS NOT NULL THEN 'likes this post'
        WHEN pc.user_id IS NOT NULL THEN 'commented on this post'
        WHEN pm.mentioned_user_id IS NOT NULL THEN 'was mentioned in this post'
        WHEN v.user_id IS NOT NULL THEN 'viewed this post'
        WHEN r.user_id IS NOT NULL THEN 'reposted this post'
      END AS interaction_type
    FROM posts p
    LEFT JOIN likes l ON p.id = l.post_id
    LEFT JOIN post_comments pc ON p.id = pc.post_id
    LEFT JOIN post_mentions pm ON p.id = pm.post_id
    LEFT JOIN post_views v ON p.id = v.post_id
    LEFT JOIN repost r ON p.id = r.original_post_id
    LEFT JOIN follow f ON l.user_id = f.following_user_id
      OR pc.user_id = f.following_user_id
      OR pm.mentioned_user_id = f.following_user_id
      OR v.user_id = f.following_user_id
      OR r.user_id = f.following_user_id
    LEFT JOIN authentications au ON l.user_id = au.id
      OR pc.user_id = au.id
      OR pm.mentioned_user_id = au.id
      OR v.user_id = au.id
      OR r.user_id = au.id
    WHERE EXISTS (
      SELECT 1
      FROM follow
      WHERE follower_user_id = $1
        AND following_user_id = f.following_user_id
    )
    GROUP BY p.id, interaction_type
  ) AS social_graph_posts
  JOIN posts p ON social_graph_posts.id = p.id
  LEFT JOIN authentications au ON p.user_id = au.id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $1 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $1)
  WHERE social_graph_rank <= 6
    AND p.user_id != $1
    AND p.deleted_at IS NULL
    AND p.suspended_at IS NULL
    AND bu.id IS NULL
  GROUP BY p.id

  UNION ALL

  -- ! Content-Based Recommendations: Recommend posts with similar hashtags or content that the user has interacted with before.
  SELECT
    p.*,
    0 AS popularity_rank,
    'content_based' AS recommendation_type,
    'Because you interacted with posts containing similar hashtags' AS reason,
    NULL AS rUsername,
    NULL AS title
  FROM posts p
  INNER JOIN post_hashtag ph ON ph.post_id = p.id
  INNER JOIN (
    SELECT DISTINCT hashtag_id
    FROM post_hashtag
    INNER JOIN likes ON post_hashtag.post_id = likes.post_id
    WHERE likes.user_id = $1
  ) AS user_liked_tags ON ph.hashtag_id = user_liked_tags.hashtag_id
  LEFT JOIN blocked_users bu ON (bu.blocking_user_id = $1 AND bu.blocked_user_id = p.user_id)
    OR (bu.blocking_user_id = p.user_id AND bu.blocked_user_id = $1)
  WHERE p.id NOT IN (
    SELECT id
    FROM posts
    WHERE user_id = $1
  )
    AND p.deleted_at IS NULL
    AND p.suspended_at IS NULL
    AND bu.id IS NULL
),
RANKED_RECOMMENDATIONS AS (
  SELECT
    recommendation_type,
    reason,
    rUsername,
    ROW_NUMBER() OVER (PARTITION BY recommendation_type ORDER BY id) AS rn
  FROM RECOMMENDATIONS
),
POSTS_WITH_IMAGES_AND_USER_DETAILS AS (
  SELECT
    r.*,
    ue.name,
    ue.username,
    ue.is_verified,
    ue.entity_type,
    ue.following_count,
    ue.followers_count,
    ue.user_image_url,
    json_agg(pi.image_url) AS post_image_urls,
    pm.views,
    pm.likes,
    pm.comments,
    pm.reposts
  FROM (
    SELECT DISTINCT ON (p.id)
      p.*
    FROM RANKED_RECOMMENDATIONS r
    INNER JOIN RECOMMENDATIONS p
      ON r.rn <= 20 / (
        SELECT COUNT(DISTINCT recommendation_type)
        FROM RANKED_RECOMMENDATIONS
      )
      AND p.recommendation_type = r.recommendation_type
    ORDER BY p.id, r.rn
    LIMIT $2 OFFSET $3
  ) r
  JOIN USER_ENTITY ue ON r.user_id = ue.auth_id
  LEFT JOIN posts_images pi ON r.id = pi.post_id
  LEFT JOIN posts_metrics pm ON r.id = pm.post_id
  GROUP BY
    r.id,
    r.content,
    r.created_at,
    r.user_id,
    r.recommendation_type,
    r.total_images,
    r.updated_at,
    r.deleted_at,
    r.suspended_at,
    ue.username,
    ue.is_verified,
    ue.following_count,
    ue.followers_count,
    ue.entity_type,
    ue.name,
    ue.user_image_url,
    r.popularity_rank,
    r.title,
    pm.views,
    pm.likes,
    pm.comments,
    pm.reposts,
    r.reason,
    r.rUsername
)
SELECT * FROM POSTS_WITH_IMAGES_AND_USER_DETAILS;




-- name: GetPostsRecommendationOld :many
WITH USER_ENTITY AS (
  SELECT
    au.id AS auth_id,
    au.username,
    au.is_verified,
    'user' AS entity_type,
    us.first_name AS name,
    ep.image_url AS user_image_url,
    ep.following_count,
    ep.followers_count,
    NULL AS members_count
  FROM authentications au
  JOIN users us ON au.id = us.user_id
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE au.user_type = 'user'
  AND au.is_deleted = FALSE
  AND au.is_suspended = FALSE
  UNION ALL
  SELECT
    au.id AS auth_id,
    au.username,
    au.is_verified,
    'churchAdmin' AS entity_type,
    c.name,
    ep.image_url AS user_image_url,
    ep.following_count,
    ep.followers_count,
    c.members_count
  FROM authentications au
  JOIN churches c ON au.id = c.auth_id
  JOIN entity_profiles ep ON au.id = ep.user_id
  WHERE au.user_type = 'churchAdmin'
  AND au.is_deleted = FALSE
  AND au.is_suspended = FALSE
),
RECOMMENDATIONS AS (
  
  -- ! Get posts the user liked, repost or was mentioned in
    SELECT p.*, 0 AS popularity_rank, 'recent_interactions' AS recommendation_type, 'Because you recently interacted with this post' AS reason, NULL AS rUsername, NULL AS title
    FROM posts p
    JOIN (
        SELECT post_id
        FROM likes
        WHERE likes.user_id = $1
        UNION
        SELECT original_post_id AS post_id
        FROM repost
        WHERE user_id = $1
        UNION
        SELECT post_id
        FROM post_mentions
        WHERE mentioned_user_id = $1
    ) AS recent_activity ON p.id = recent_activity.post_id
    WHERE p.deleted_at IS NULL AND p.suspended_at IS NULL
    UNION

    -- ! Get posts of the user's followers
    SELECT p.*, 0 AS popularity_rank, 'follow_based' AS recommendation_type, 'Because you follow @' || ue.username AS reason, ue.username AS rUsername, NULL AS title
    FROM posts p
    JOIN follow f ON p.user_id = f.following_user_id
    JOIN USER_ENTITY ue ON p.user_id = ue.auth_id
    WHERE f.follower_user_id = $1
    AND p.deleted_at IS NULL AND p.suspended_at IS NULL
    UNION
    

    -- ! Get popular posts
    SELECT p.*, 'popularity_based' AS recommendation_type, 'This post is popular among users' AS reason, NULL AS rUsername, NULL AS title
    FROM (
        -- SELECT p.*, ROW_NUMBER() OVER (ORDER BY (pm.likes + pm.reposts + pm.comments) DESC) AS popularity_rank
        SELECT p.*, ROW_NUMBER() OVER ( ORDER BY ((pm.likes * 0.5) + (pm.reposts * 0.5) + (pm.comments * 0.5) + (EXTRACT(EPOCH FROM NOW() - p.created_at) / 3600 / 24)) DESC) AS popularity_rank
        FROM posts p
        JOIN posts_metrics pm ON p.id = pm.post_id
        WHERE p.created_at >= NOW() - INTERVAL '7 days'
    ) p
    WHERE p.popularity_rank <= 5
    AND p.deleted_at IS NULL AND p.suspended_at IS NULL
    UNION

    -- TODO: Test suspended or deleted post is not included
    -- ! Get posts that the user's followers liked
    SELECT p.*, 0 AS popularity_rank, 'collaborative_filtering' AS recommendation_type, 'Because ' || ue.username || ' likes this post' AS reason, NULL AS rUsername, NULL AS title
    FROM posts p
    JOIN likes l ON p.id = l.post_id
    JOIN USER_ENTITY ue ON l.user_id = ue.auth_id
    WHERE l.user_id IN (
        SELECT following_user_id
        FROM follow
        WHERE follower_user_id = $1
    )
    AND p.deleted_at IS NULL AND p.suspended_at IS NULL
    AND p.id NOT IN (
        SELECT id
        FROM posts
        WHERE user_id = $1
    )
    UNION

    -- ! Get announcements by user's church
    SELECT
    a.id,
    a.user_id,
    a.content,
    a.total_images,
    a.created_at,
    a.updated_at,
    a.suspended_at,
    a.deleted_at,
    0 AS popularity_rank,
    'church_announcement' AS recommendation_type,
    'Announcement from ' || c.name AS reason,
    NULL AS rUsername,
    a.title
  FROM
    announcements a
    JOIN churches c ON a.user_id = c.auth_id
    JOIN user_church_membership ucm ON ucm.church_id = c.id AND ucm.user_id = $1 AND ucm.active = TRUE
  WHERE
    a.created_at >= NOW() - INTERVAL '2 weeks'
    AND a.deleted_at IS NULL AND a.suspended_at IS NULL
    UNION

    -- ! Get posts based on social graph
    SELECT
  p.*,
  0 AS popularity_rank,
  'social_graph' AS recommendation_type,
  'Because ' || string_agg(DISTINCT au.username || ' ' || interaction_type, ', ') AS reason,
  NULL AS rUsername, NULL AS title
FROM (
  SELECT
    p.id,
    ROW_NUMBER() OVER (ORDER BY (COUNT(l.user_id) + COUNT(pc.user_id) + COUNT(pm.mentioned_user_id) + COUNT(v.user_id) + COUNT(r.user_id)) DESC) AS social_graph_rank,
    CASE
      WHEN l.user_id IS NOT NULL THEN 'likes this post'
      WHEN pc.user_id IS NOT NULL THEN 'commented on this post'
      WHEN pm.mentioned_user_id IS NOT NULL THEN 'was mentioned in this post'
      WHEN v.user_id IS NOT NULL THEN 'viewed this post'
      WHEN r.user_id IS NOT NULL THEN 'reposted this post'
    END AS interaction_type
  FROM
    posts p
    LEFT JOIN likes l ON p.id = l.post_id
    LEFT JOIN post_comments pc ON p.id = pc.post_id
    LEFT JOIN post_mentions pm ON p.id = pm.post_id
    LEFT JOIN post_views v ON p.id = v.post_id
    LEFT JOIN repost r ON p.id = r.original_post_id
    LEFT JOIN follow f ON l.user_id = f.following_user_id OR pc.user_id = f.following_user_id OR pm.mentioned_user_id = f.following_user_id OR v.user_id = f.following_user_id OR r.user_id = f.following_user_id
    LEFT JOIN authentications au ON l.user_id = au.id OR pc.user_id = au.id OR pm.mentioned_user_id = au.id OR v.user_id = au.id OR r.user_id = au.id
  WHERE
    EXISTS (
      SELECT 1
      FROM follow
      WHERE follower_user_id = $1
        AND following_user_id = f.follower_user_id
    )
  GROUP BY
    p.id, interaction_type
) AS social_graph_posts
JOIN posts p ON social_graph_posts.id = p.id
LEFT JOIN authentications au ON p.user_id = au.id
WHERE
  social_graph_rank <= 6
  AND p.user_id != $1
  AND p.deleted_at IS NULL AND p.suspended_at IS NULL
GROUP BY
  p.id
    UNION

    -- TODO: Test suspended or deleted post is not included
    -- ! content based recommendation
    SELECT p.*, 0 AS popularity_rank, 'content_based' AS recommendation_type, 'Because you interacted with posts containing similar hashtags' AS reason, NULL AS rUsername, NULL AS title
    FROM posts p
    INNER JOIN post_hashtag ph ON ph.post_id = p.id
    INNER JOIN (
        SELECT hashtag_id
        FROM post_hashtag
        WHERE post_id IN (
            SELECT post_id
            FROM likes
            WHERE user_id = $1
        )
    ) AS user_liked_tags ON ph.hashtag_id = user_liked_tags.hashtag_id
    WHERE p.id NOT IN (
        SELECT id
        FROM posts
        WHERE user_id = $1
    )
    AND p.deleted_at IS NULL AND p.suspended_at IS NULL
),
RANKED_RECOMMENDATIONS AS (
    SELECT
        recommendation_type,
        reason,
        rUsername,
        ROW_NUMBER() OVER (PARTITION BY recommendation_type ORDER BY id) AS rn
    FROM
        RECOMMENDATIONS
),
POSTS_WITH_IMAGES_AND_USER_DETAILS AS (
  SELECT
    r.*,
    ue.name,
    ue.username,
    ue.is_verified,
    ue.entity_type,
    ue.following_count,
    ue.followers_count,
    ue.user_image_url,
    json_agg(pi.image_url) AS post_image_urls,
    pm.views,
    pm.likes,
    pm.comments,
    pm.reposts
  FROM
        (
            SELECT DISTINCT ON (p.id)
                p.*
            FROM
                RANKED_RECOMMENDATIONS r
            INNER JOIN
                RECOMMENDATIONS p ON r.rn <= 20 / (
                    SELECT
                        COUNT(DISTINCT recommendation_type)
                    FROM
                        RANKED_RECOMMENDATIONS
                )
                AND p.recommendation_type = r.recommendation_type
            ORDER BY p.id, r.rn
            LIMIT $2 OFFSET $3
        ) r
  JOIN USER_ENTITY ue ON r.user_id = ue.auth_id
  LEFT JOIN posts_images pi ON r.id = pi.post_id
  LEFT JOIN posts_metrics pm ON r.id = pm.post_id
  GROUP BY
    r.id, r.content, r.created_at, r.user_id, r.recommendation_type,
    r.total_images, r.updated_at, r.deleted_at, r.suspended_at,
    ue.username, ue.is_verified, ue.following_count, ue.followers_count,
    ue.entity_type, ue.name, ue.user_image_url, r.popularity_rank, r.title,
    pm.views, pm.likes, pm.comments, pm.reposts, r.reason, r.rUsername
)
SELECT * FROM POSTS_WITH_IMAGES_AND_USER_DETAILS;







-- -- * Collaborative Filtering: Recommend posts liked by users who have similar tastes as the current user, based on the posts they have both liked.
-- SELECT p.id, p.content, COUNT(*) AS common_likes
-- FROM likes l
-- INNER JOIN posts p ON l.post_id = p.id
-- WHERE l.user_id IN (
--   SELECT DISTINCT l2.user_id
--   FROM likes l1
--   INNER JOIN likes l2 ON l1.post_id = l2.post_id
--   WHERE l1.user_id = $1 AND l2.user_id != $1
-- )
-- AND l.post_id NOT IN (
--   SELECT post_id FROM likes WHERE user_id = $1
-- )
-- GROUP BY p.id, p.content
-- ORDER BY common_likes DESC
-- LIMIT 10;


-- ! Item-Based Collaborative Filtering: Recommend items based on similarity scores between items, which can be calculated using SQL aggregate functions and joins.
-- SELECT B.item_id, COUNT(*) as similarity_score
-- FROM user_likes A
-- JOIN user_likes B ON A.user_id = B.user_id
-- WHERE A.item_id = 'known_liked_item_id'
-- AND B.item_id != 'known_liked_item_id'
-- GROUP BY B.item_id
-- ORDER BY similarity_score DESC
-- LIMIT 10;


-- ! User Segmentation: Group users by certain characteristics (location, age, previous likes) and recommend popular items within the segment.
-- SELECT item_id, COUNT(*) as total_likes
-- FROM user_likes
-- WHERE user_id IN (
--     SELECT user_id
--     FROM authentications
--     WHERE age BETWEEN 18 AND 25
-- )
-- GROUP BY item_id
-- ORDER BY total_likes DESC
-- LIMIT 10;


-- ! Weighted Hybrid Techniques: Combine results from various SQL-based recommendation strategies with weights to account for different confidence levels.
-- SELECT item_id,
--        (similarity_score * {weight_similarity} +
--         frequency * {weight_frequency} +
--         total_likes * {weight_likes}) AS weighted_score
-- FROM (
--     -- Combine the results from the above three methods
--     -- Each sub-query would correspond to an individual method
--     -- You would join them on item_id and then calculate the weighted score
-- ) AS combined_results
-- ORDER BY weighted_score DESC
-- LIMIT 10;