-- Drop foreign key constraints
-- TODO: Drop post_views index, constraints and table
ALTER TABLE "post_comments" DROP CONSTRAINT IF EXISTS "post_comments_user_id_fkey";
ALTER TABLE "post_comments" DROP CONSTRAINT IF EXISTS "post_comments_post_id_fkey";

ALTER TABLE "repost" DROP CONSTRAINT IF EXISTS "unique_uid_org_post_id";
ALTER TABLE "repost" DROP CONSTRAINT IF EXISTS "repost_user_id_fkey";
ALTER TABLE "repost" DROP CONSTRAINT IF EXISTS "repost_original_post_id_fkey";

ALTER TABLE "likes" DROP CONSTRAINT IF EXISTS "likes_user_id_fkey";
ALTER TABLE "likes" DROP CONSTRAINT IF EXISTS "likes_post_id_fkey";

ALTER TABLE "post_mentions" DROP CONSTRAINT IF EXISTS "post_mentions_post_id_fkey";
ALTER TABLE "post_mentions" DROP CONSTRAINT IF EXISTS "post_mentions_mentioned_user_id_fkey";
ALTER TABLE "post_mentions" DROP CONSTRAINT IF EXISTS "unique_user_per_post_mentions";

ALTER TABLE "posts_metrics" DROP CONSTRAINT IF EXISTS "posts_metrics_post_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "post_comments_user_id_id_post_id_idx";
DROP INDEX IF EXISTS "repost_user_id_id_original_post_id_idx";
DROP INDEX IF EXISTS "post_mentions_post_id_mentioned_user_id_idx";
DROP INDEX IF EXISTS "posts_metrics_post_id_idx";

DROP INDEX IF EXISTS likes_user_id_post_id_idx;-- ON likes (user_id, post_id);

-- Drop tables
DROP TABLE IF EXISTS "post_comments";
DROP TABLE IF EXISTS "posts_metrics";
DROP TABLE IF EXISTS "repost";
DROP TABLE IF EXISTS "likes";
DROP TABLE IF EXISTS "post_mentions";