ALTER TABLE "user_roles" DROP CONSTRAINT IF EXISTS "user_roles_user_id_fkey";
ALTER TABLE "user_roles" DROP CONSTRAINT IF EXISTS "user_roles_role_id_fkey";
ALTER TABLE "entity_profiles" DROP CONSTRAINT IF EXISTS "user_profiles_user_id_fkey";

-- ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "user_profiles_user_id_fkey"; TODO: Fix this

DROP INDEX IF EXISTS "user_id_email_phone_username_created_at_updated_at_idx";

DROP INDEX IF EXISTS "user_roles_role_id_idx";

DROP INDEX IF EXISTS "user_profiles_user_id_following_count_follower_count_idx";

DROP INDEX IF EXISTS "users_email_username_created_at_updated_at";
DROP INDEX IF EXISTS "accounts_user_id_idx";


DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "user_roles";
DROP TABLE IF EXISTS "entity_profiles";
DROP TABLE IF EXISTS "authentications";
DROP TABLE accounts;
