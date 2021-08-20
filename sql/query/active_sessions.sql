-- name: CreateActiveSession :exec
INSERT INTO active_sessions (
  session_id, user_id
) VALUES ($1, $2);

-- name: GetActiveSessions :many
SELECT * FROM active_sessions WHERE user_id = $1;

-- name: DeleteActiveSessionsForUser :exec
DELETE FROM active_sessions WHERE user_id = $1;

-- name: DeleteActiveSessionsByID :exec
DELETE FROM active_sessions WHERE session_id = ANY($1::char(86)[]);
