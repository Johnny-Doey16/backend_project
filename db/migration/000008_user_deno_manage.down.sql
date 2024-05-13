-- Remove the indexes created
DROP INDEX IF EXISTS idx_user_church_membership_user_id;
DROP INDEX IF EXISTS idx_user_denomination_membership_user_id;
DROP INDEX IF EXISTS idx_user_profiles_user_id;
DROP INDEX IF EXISTS idx_user_active_membership;
DROP INDEX IF EXISTS idx_user_church_active_membership;


-- Drop the foreign key constraints from users before dropping columns
ALTER TABLE users
DROP CONSTRAINT IF EXISTS user_profiles_denomination_id_fkey,
DROP CONSTRAINT IF EXISTS user_profiles_church_id_fkey;

-- TODO: Drop fks

-- Remove the columns added to the users table
ALTER TABLE users
DROP COLUMN IF EXISTS denomination_id,
DROP COLUMN IF EXISTS church_id,
DROP COLUMN IF EXISTS last_denomination_change,
DROP COLUMN IF EXISTS last_church_change;

-- Drop the User Church Membership table
DROP TABLE IF EXISTS user_church_membership;

-- Drop the User Denomination Membership table
DROP TABLE IF EXISTS user_denomination_membership;

-- Drop the Churches table
DROP TABLE IF EXISTS church_locations;

-- Drop the Churches table
DROP TABLE IF EXISTS churches;

-- Drop the Denominations table
DROP TABLE IF EXISTS denominations;