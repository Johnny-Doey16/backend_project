-- Drop refs
FOREIGN KEY (church_id) REFERENCES churches(id)
ALTER TABLE "church_projects" DROP CONSTRAINT IF EXISTS "church_projects_church_id_fkey";


ALTER TABLE "project_donations" DROP CONSTRAINT IF EXISTS "project_donations_project_id_fkey";
ALTER TABLE "project_donations" DROP CONSTRAINT IF EXISTS "project_donations_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "church_projects_church_id_idx";
DROP INDEX IF EXISTS "project_donations_user_id_project_id_idx";

-- Drop tables
DROP TABLE IF EXISTS church_projects;
DROP TABLE IF EXISTS project_donations;
