ALTER TABLE "change_identifier_requests" DROP CONSTRAINT IF EXISTS "change_identifier_requests_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "change_identifier_requests_user_id_identifier_token_expires_at_idx";

-- Drop tables
DROP TABLE IF EXISTS "change_identifier_requests";