-- name: CreateGroupWithAdmin :one
WITH new_group AS (
    INSERT INTO groups (denomination_id, name, description)
    VALUES ($1, $2, $3)
    RETURNING id AS group_id -- Specify the column alias here
)
INSERT INTO user_group_membership (group_id, user_id, is_admin)
SELECT group_id, $4, TRUE FROM new_group -- Use the column alias
RETURNING group_id, user_id, is_admin;

-- name: JoinGroup :one
INSERT INTO user_group_membership (group_id, user_id)
SELECT $1, $2
WHERE NOT EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE group_id=$1 AND user_id=$2
)
RETURNING *;

-- name: GetGroup :one
SELECT * FROM groups WHERE id = $1;

-- name: ListGroups :many
SELECT * FROM groups;

-- name: GetMembership :one
SELECT * FROM user_group_membership WHERE id = $1;

-- name: ListMemberships :many
SELECT * FROM user_group_membership WHERE group_id = $1;

-- *** ADMIN FUNCTIONS
-- name: UpdateGroup :exec
UPDATE groups
SET name=$2, description=$3
WHERE groups.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE user_group_membership.group_id=groups.id AND user_group_membership.user_id=$4 AND user_group_membership.is_admin=true
);

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE groups.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE user_group_membership.group_id=groups.id AND user_group_membership.user_id=$2 AND user_group_membership.is_admin=true
);

-- name: CreateMembership :one
INSERT INTO user_group_membership (group_id, user_id, join_date, is_admin)
SELECT $1, $2, now(), $3
WHERE EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE user_group_membership.group_id=$1 AND user_group_membership.user_id=$4 AND user_group_membership.is_admin=true
)
RETURNING *;

-- name: UpdateMembership :exec
UPDATE user_group_membership
SET is_admin=$2
WHERE user_group_membership.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership AS ugm
    WHERE ugm.group_id=(SELECT group_id FROM user_group_membership WHERE id=$1) AND ugm.user_id=$3 AND ugm.is_admin=true
);

-- name: DeleteMembership :exec
DELETE FROM user_group_membership
WHERE user_group_membership.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership AS ugm
    WHERE ugm.group_id=(SELECT group_id FROM user_group_membership WHERE user_group_membership.id=$1)
    AND ugm.user_id=$2 AND ugm.is_admin=true
);