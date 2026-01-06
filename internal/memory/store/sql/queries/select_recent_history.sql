SELECT id, role, content, created_at, session_id
FROM history
ORDER BY created_at DESC
LIMIT ?;
