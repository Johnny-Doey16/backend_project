-- TODO: Drop unique constraints on tables
-- Drop index on type column
DROP INDEX IF EXISTS idx_notifications_type;

-- Drop mentions_notifications table
DROP TABLE IF EXISTS mentions_notifications;

-- Drop follow_notifications table
DROP TABLE IF EXISTS follow_notifications;

-- Drop notifications table
DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS post_comment_notifications;

DROP TABLE IF EXISTS post_like_notifications;

DROP TABLE IF EXISTS announcement_notifications;