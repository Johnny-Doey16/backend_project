-- Reversing the actions in the up migration

-- Drop foreign key constraints
ALTER TABLE "sessions" DROP CONSTRAINT "sessions_user_id_fkey";
ALTER TABLE "user_logins" DROP CONSTRAINT "user_logins_user_id_fkey";
ALTER TABLE "email_verification_requests" DROP CONSTRAINT IF EXISTS "email_verification_requests_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "sessions_user_id_expires_at_idx";
DROP INDEX IF EXISTS "user_logins_user_id_login_at_idx";
DROP INDEX IF EXISTS "email_verification_requests_user_id_token_email_created_at_updated_at_idx";

-- Drop tables
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "user_logins";
DROP TABLE IF EXISTS "email_verification_requests";