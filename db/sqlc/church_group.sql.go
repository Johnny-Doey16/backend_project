// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: church_group.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createGroupWithAdmin = `-- name: CreateGroupWithAdmin :one
WITH new_group AS (
    INSERT INTO groups (denomination_id, name, description)
    VALUES ($1, $2, $3)
    RETURNING id AS group_id -- Specify the column alias here
)
INSERT INTO user_group_membership (group_id, user_id, is_admin)
SELECT group_id, $4, TRUE FROM new_group -- Use the column alias
RETURNING group_id, user_id, is_admin
`

type CreateGroupWithAdminParams struct {
	DenominationID int32          `json:"denomination_id"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	UserID         uuid.UUID      `json:"user_id"`
}

type CreateGroupWithAdminRow struct {
	GroupID int32        `json:"group_id"`
	UserID  uuid.UUID    `json:"user_id"`
	IsAdmin sql.NullBool `json:"is_admin"`
}

func (q *Queries) CreateGroupWithAdmin(ctx context.Context, arg CreateGroupWithAdminParams) (CreateGroupWithAdminRow, error) {
	row := q.db.QueryRowContext(ctx, createGroupWithAdmin,
		arg.DenominationID,
		arg.Name,
		arg.Description,
		arg.UserID,
	)
	var i CreateGroupWithAdminRow
	err := row.Scan(&i.GroupID, &i.UserID, &i.IsAdmin)
	return i, err
}

const createMembership = `-- name: CreateMembership :one
INSERT INTO user_group_membership (group_id, user_id, join_date, is_admin)
SELECT $1, $2, now(), $3
WHERE EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE user_group_membership.group_id=$1 AND user_group_membership.user_id=$4 AND user_group_membership.is_admin=true
)
RETURNING id, group_id, user_id, join_date, is_admin
`

type CreateMembershipParams struct {
	GroupID  int32        `json:"group_id"`
	UserID   uuid.UUID    `json:"user_id"`
	IsAdmin  sql.NullBool `json:"is_admin"`
	UserID_2 uuid.UUID    `json:"user_id_2"`
}

func (q *Queries) CreateMembership(ctx context.Context, arg CreateMembershipParams) (UserGroupMembership, error) {
	row := q.db.QueryRowContext(ctx, createMembership,
		arg.GroupID,
		arg.UserID,
		arg.IsAdmin,
		arg.UserID_2,
	)
	var i UserGroupMembership
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.UserID,
		&i.JoinDate,
		&i.IsAdmin,
	)
	return i, err
}

const deleteGroup = `-- name: DeleteGroup :exec
DELETE FROM groups
WHERE groups.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE user_group_membership.group_id=groups.id AND user_group_membership.user_id=$2 AND user_group_membership.is_admin=true
)
`

type DeleteGroupParams struct {
	ID     int32     `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteGroup(ctx context.Context, arg DeleteGroupParams) error {
	_, err := q.db.ExecContext(ctx, deleteGroup, arg.ID, arg.UserID)
	return err
}

const deleteMembership = `-- name: DeleteMembership :exec
DELETE FROM user_group_membership
WHERE user_group_membership.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership AS ugm
    WHERE ugm.group_id=(SELECT group_id FROM user_group_membership WHERE user_group_membership.id=$1)
    AND ugm.user_id=$2 AND ugm.is_admin=true
)
`

type DeleteMembershipParams struct {
	ID     int32     `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteMembership(ctx context.Context, arg DeleteMembershipParams) error {
	_, err := q.db.ExecContext(ctx, deleteMembership, arg.ID, arg.UserID)
	return err
}

const getGroup = `-- name: GetGroup :one
SELECT id, denomination_id, name, description FROM groups WHERE id = $1
`

func (q *Queries) GetGroup(ctx context.Context, id int32) (Group, error) {
	row := q.db.QueryRowContext(ctx, getGroup, id)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.DenominationID,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const getMembership = `-- name: GetMembership :one
SELECT id, group_id, user_id, join_date, is_admin FROM user_group_membership WHERE id = $1
`

func (q *Queries) GetMembership(ctx context.Context, id int32) (UserGroupMembership, error) {
	row := q.db.QueryRowContext(ctx, getMembership, id)
	var i UserGroupMembership
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.UserID,
		&i.JoinDate,
		&i.IsAdmin,
	)
	return i, err
}

const joinGroup = `-- name: JoinGroup :one
INSERT INTO user_group_membership (group_id, user_id)
SELECT $1, $2
WHERE NOT EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE group_id=$1 AND user_id=$2
)
RETURNING id, group_id, user_id, join_date, is_admin
`

type JoinGroupParams struct {
	GroupID int32     `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) JoinGroup(ctx context.Context, arg JoinGroupParams) (UserGroupMembership, error) {
	row := q.db.QueryRowContext(ctx, joinGroup, arg.GroupID, arg.UserID)
	var i UserGroupMembership
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.UserID,
		&i.JoinDate,
		&i.IsAdmin,
	)
	return i, err
}

const listGroups = `-- name: ListGroups :many
SELECT id, denomination_id, name, description FROM groups
`

func (q *Queries) ListGroups(ctx context.Context) ([]Group, error) {
	rows, err := q.db.QueryContext(ctx, listGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Group{}
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.ID,
			&i.DenominationID,
			&i.Name,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMemberships = `-- name: ListMemberships :many
SELECT id, group_id, user_id, join_date, is_admin FROM user_group_membership WHERE group_id = $1
`

func (q *Queries) ListMemberships(ctx context.Context, groupID int32) ([]UserGroupMembership, error) {
	rows, err := q.db.QueryContext(ctx, listMemberships, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserGroupMembership{}
	for rows.Next() {
		var i UserGroupMembership
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.UserID,
			&i.JoinDate,
			&i.IsAdmin,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGroup = `-- name: UpdateGroup :exec
UPDATE groups
SET name=$2, description=$3
WHERE groups.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership
    WHERE user_group_membership.group_id=groups.id AND user_group_membership.user_id=$4 AND user_group_membership.is_admin=true
)
`

type UpdateGroupParams struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	UserID      uuid.UUID      `json:"user_id"`
}

// *** ADMIN FUNCTIONS
func (q *Queries) UpdateGroup(ctx context.Context, arg UpdateGroupParams) error {
	_, err := q.db.ExecContext(ctx, updateGroup,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.UserID,
	)
	return err
}

const updateMembership = `-- name: UpdateMembership :exec
UPDATE user_group_membership
SET is_admin=$2
WHERE user_group_membership.id=$1 AND EXISTS (
    SELECT 1 FROM user_group_membership AS ugm
    WHERE ugm.group_id=(SELECT group_id FROM user_group_membership WHERE id=$1) AND ugm.user_id=$3 AND ugm.is_admin=true
)
`

type UpdateMembershipParams struct {
	ID      int32        `json:"id"`
	IsAdmin sql.NullBool `json:"is_admin"`
	UserID  uuid.UUID    `json:"user_id"`
}

func (q *Queries) UpdateMembership(ctx context.Context, arg UpdateMembershipParams) error {
	_, err := q.db.ExecContext(ctx, updateMembership, arg.ID, arg.IsAdmin, arg.UserID)
	return err
}
