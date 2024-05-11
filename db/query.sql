-- name: GetActiveUpdates :many
SELECT * FROM updates WHERE archived = FALSE ORDER BY created_at DESC;

-- name: GetUpdatesSince :many
SELECT * FROM updates WHERE created_at > ? ORDER BY created_at DESC;

-- name: ArchiveUpdate :exec
UPDATE updates SET archived = TRUE WHERE id = ?;

-- name: CreateUpdate :one
INSERT INTO updates (description) VALUES (?) RETURNING *;

