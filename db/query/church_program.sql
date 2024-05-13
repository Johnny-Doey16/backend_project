-- name: CreateChurchProgram :exec
INSERT INTO church_programs (church_id, program_type, program_name, program_desc, program_day, program_start_time, program_end_time, program_freq, program_image_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetChurchProgramsByChurchId :many
SELECT * FROM church_programs
WHERE church_id = $1
-- ORDER BY id
OFFSET $2
LIMIT $3;

-- name: UpdateChurchProgram :exec
UPDATE church_programs
SET
    program_type = COALESCE(sqlc.narg('program_type'), program_type),
    program_name = COALESCE(sqlc.narg('program_name'), program_name),
    program_desc = COALESCE(sqlc.narg('program_desc'), program_desc),
    program_day = COALESCE(sqlc.narg('program_day'), program_day),
    program_start_time = COALESCE(sqlc.narg('program_start_time'), program_start_time),
    program_end_time = COALESCE(sqlc.narg('program_end_time'), program_end_time),
    program_freq = COALESCE(sqlc.narg('program_freq'), program_freq),
    program_image_url = COALESCE(sqlc.narg('program_image_url'), program_image_url),
    updated_at = NOW()
WHERE
    id = sqlc.arg('id');

-- name: DeleteChurchProgram :exec
DELETE FROM church_programs
WHERE id = $1;

-- name: GetChurchProgramById :many
SELECT * FROM church_programs
WHERE id = $1;

-- name: GetChurchProgramsByType :many
SELECT * FROM church_programs
WHERE program_type = $1;

-- name: GetChurchProgramsByDay :many
SELECT * FROM church_programs
WHERE program_day = $1;

-- name: GetChurchProgramsByFrequency :many
SELECT * FROM church_programs
WHERE program_freq = $1;

-- name: GetChurchProgramsInRange :many
SELECT * FROM church_programs
WHERE program_start_time >= $1 AND program_end_time <= $2;
