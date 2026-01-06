SELECT h.id, h.role, h.content, h.created_at, h.session_id,
       snippet(history_fts, 0, '>>>', '<<<', '...', 32) as snippet,
       rank
FROM history h
JOIN history_fts fts ON h.rowid = fts.rowid
WHERE history_fts MATCH ?
ORDER BY rank
LIMIT ?;
