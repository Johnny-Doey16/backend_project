-- TODO: DROP this
-- CREATE TYPE participant_status AS ENUM ('invited', 'accepted', 'declined');

ALTER TABLE "prayer_room" DROP CONSTRAINT IF EXISTS "prayer_room_author_id_fkey";

ALTER TABLE "prayer_participants" DROP CONSTRAINT IF EXISTS "prayer_participants_user_id_fkey";
ALTER TABLE "prayer_participants" DROP CONSTRAINT IF EXISTS "prayer_participants_room_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "prayer_room_room_id_idx";
DROP INDEX IF EXISTS "prayer_room_author_id_idx";
DROP INDEX IF EXISTS "prayer_participants_user_id_idx";
DROP INDEX IF EXISTS "prayer_participants_room_id_idx";

-- Drop prayer_participants Table
DROP TABLE prayer_participants;

-- Drop prayer_room Table
DROP TABLE prayer_room;


-- Drop prayer_invite_notifications table
DROP TABLE IF EXISTS prayer_invite_notifications;