-- Drop foreign key constraints
ALTER TABLE "posts_images" DROP CONSTRAINT IF EXISTS "posts_images_post_id_fkey";
ALTER TABLE "posts" DROP CONSTRAINT IF EXISTS "posts_user_id_fkey";
ALTER TABLE "announcements" DROP CONSTRAINT IF EXISTS "announcements_user_id_fkey";
ALTER TABLE "follow" DROP CONSTRAINT IF EXISTS "follow_follower_user_id_fkey";
ALTER TABLE "follow" DROP CONSTRAINT IF EXISTS "follow_following_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "follow_follower_user_id_following_user_id_idx";
DROP INDEX IF EXISTS "posts_images_post_id_id_image_url_idx";
DROP INDEX IF EXISTS "posts_user_id_id_created_at_idx";

-- Drop tables
DROP TABLE IF EXISTS "follow";
DROP TABLE IF EXISTS "posts_images";
DROP TABLE IF EXISTS "posts";

DROP TABLE IF EXISTS "announcements";
