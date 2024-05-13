-- name: CreatePrayerRoom :exec
INSERT INTO prayer_room (room_id, author_id, name, start_time, end_time)
VALUES ($1, $2, $3, $4, $5);

-- name: CreatePrayerParticipants :exec
INSERT INTO prayer_participants (user_id, room_id)
SELECT user_id, $2
FROM UNNEST($1::uuid[]) AS user_id;


-- name: UpdatePrayerRoom :exec
UPDATE prayer_room
SET
    name = COALESCE(sqlc.narg('name'), name),
    start_time = COALESCE(sqlc.narg('start_time'), start_time),
    end_time = COALESCE(sqlc.narg('end_time'), end_time)
WHERE
    room_id = sqlc.arg('room_id');

-- name: DeletePrayerRoom :exec
DELETE FROM prayer_room
WHERE room_id = $1;

-- name: DeleteParticipant :exec
DELETE FROM prayer_participants
WHERE room_id = $1;

-- name: UpdatePrayerInvitation :exec
UPDATE prayer_participants
SET status = $3
WHERE room_id = $1 AND user_id = $2;


-- name: GetPrayerRooms :many
SELECT pr.*, pp.status AS participant_status
FROM prayer_room pr
JOIN prayer_participants pp ON pr.room_id = pp.room_id
WHERE pp.user_id = $1
ORDER BY pr.start_time DESC -- Optional: you can order by any column you prefer
LIMIT $2 OFFSET $3;


-- name: GetPrayerRoomById :one
SELECT * FROM prayer_room
WHERE room_id = $1;