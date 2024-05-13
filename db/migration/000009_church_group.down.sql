-- Drop the organization_forums table and its indexes
DROP INDEX IF EXISTS idx_organization_forums_created_at;
DROP INDEX IF EXISTS idx_organization_forums_user_id;
DROP INDEX IF EXISTS idx_organization_forums_organization_id;
DROP TABLE IF EXISTS organization_forums;

-- Drop the group_forums table and its indexes
DROP INDEX IF EXISTS idx_group_forums_created_at;
DROP INDEX IF EXISTS idx_group_forums_user_id;
DROP INDEX IF EXISTS idx_group_forums_group_id;
DROP TABLE IF EXISTS group_forums;

-- Drop the user_organization_membership table and its indexes
DROP INDEX IF EXISTS idx_user_organization_membership_organization_id;
DROP INDEX IF EXISTS idx_user_organization_membership_user_id;
DROP TABLE IF EXISTS user_organization_membership;

-- Drop the organizations table and its index
DROP INDEX IF EXISTS idx_organizations_church_id;
DROP TABLE IF EXISTS organizations;

-- Drop the user_group_membership table and its indexes
DROP INDEX IF EXISTS idx_user_group_membership_group_id;
DROP INDEX IF EXISTS idx_user_group_membership_user_id;
DROP TABLE IF EXISTS user_group_membership;

-- Drop the groups table and its index
DROP INDEX IF EXISTS idx_groups_denomination_id;
DROP TABLE IF EXISTS groups;