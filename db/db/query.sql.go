// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const archiveUpdate = `-- name: ArchiveUpdate :exec
UPDATE updates SET archived = TRUE WHERE id = ?
`

func (q *Queries) ArchiveUpdate(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, archiveUpdate, id)
	return err
}

const getUpdates = `-- name: GetUpdates :many
SELECT id, description, archived, created_at, updated_at FROM updates ORDER BY created_at DESC
`

func (q *Queries) GetUpdates(ctx context.Context) ([]Update, error) {
	rows, err := q.db.QueryContext(ctx, getUpdates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Update
	for rows.Next() {
		var i Update
		if err := rows.Scan(
			&i.ID,
			&i.Description,
			&i.Archived,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getUpdatesSince = `-- name: GetUpdatesSince :many
SELECT id, description, archived, created_at, updated_at FROM updates WHERE created_at > ? ORDER BY created_at DESC
`

func (q *Queries) GetUpdatesSince(ctx context.Context, createdAt sql.NullTime) ([]Update, error) {
	rows, err := q.db.QueryContext(ctx, getUpdatesSince, createdAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Update
	for rows.Next() {
		var i Update
		if err := rows.Scan(
			&i.ID,
			&i.Description,
			&i.Archived,
			&i.CreatedAt,
			&i.UpdatedAt,
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
