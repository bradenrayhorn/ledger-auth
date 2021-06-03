-- name: CreateActiveSession :exec
INSERT INTO active_sessions (
  session_id, user_id
) VALUES (?, ?);

-- name: GetActiveSessions :many
SELECT * FROM active_sessions WHERE user_id = ?;
