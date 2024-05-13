-- Drop foreign key constraints first to avoid dependency issues
ALTER TABLE "moderation_queue" DROP CONSTRAINT IF EXISTS moderation_queue_post_id_fkey;
ALTER TABLE "moderation_queue" DROP CONSTRAINT IF EXISTS moderation_queue_user_id_fkey;

ALTER TABLE "reported_posts" DROP CONSTRAINT IF EXISTS reported_posts_user_id_fkey;
ALTER TABLE "reported_posts" DROP CONSTRAINT IF EXISTS reported_posts_post_id_fkey;

ALTER TABLE "blocked_users" DROP CONSTRAINT IF EXISTS unique_blocked_user;
ALTER TABLE "blocked_users" DROP CONSTRAINT IF EXISTS blocked_users_blocked_user_id_fkey;
ALTER TABLE "blocked_users" DROP CONSTRAINT IF EXISTS blocked_users_blocking_user_id_fkey;

ALTER TABLE "bookmarks" DROP CONSTRAINT IF EXISTS unique_user_post;
ALTER TABLE "bookmarks" DROP CONSTRAINT IF EXISTS bookmarks_post_id_fkey;
ALTER TABLE "bookmarks" DROP CONSTRAINT IF EXISTS bookmarks_user_id_fkey;

ALTER TABLE "post_hashtag" DROP CONSTRAINT IF EXISTS post_hashtag_hashtag_id_fkey;
ALTER TABLE "post_hashtag" DROP CONSTRAINT IF EXISTS post_hashtag_post_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS moderation_queue_user_id_post_id_idx;
DROP INDEX IF EXISTS reported_posts_post_id_user_id_idx;
DROP INDEX IF EXISTS blocked_users_blocked_user_id_blocking_user_id_idx;
DROP INDEX IF EXISTS bookmarks_user_id_post_id_idx;
DROP INDEX IF EXISTS post_hashtag_hashtag_id_id_post_id_idx;
DROP INDEX IF EXISTS hashtag_hash_tag_id_idx;

-- Drop tables
DROP TABLE IF EXISTS "moderation_queue";
DROP TABLE IF EXISTS "reported_posts";
DROP TABLE IF EXISTS "blocked_users";
DROP TABLE IF EXISTS "bookmarks";
DROP TABLE IF EXISTS "post_hashtag";
DROP TABLE IF EXISTS "hashtag";