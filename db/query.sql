-- name: GetUpdates :many
SELECT * FROM updates ORDER BY created_at DESC;

-- name: GetUpdatesSince :many
SELECT * FROM updates WHERE created_at > ? ORDER BY created_at DESC;

-- name: ArchiveUpdate :exec
UPDATE updates SET archived = TRUE WHERE id = ?;

INSERT INTO updates (description) VALUES (?)

