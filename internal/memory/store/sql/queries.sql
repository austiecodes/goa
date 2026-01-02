-- Memory Store Queries
-- This file contains all DML (Data Manipulation Language) queries for the memory store.

-- ============================================================================
-- MEMORY QUERIES
-- ============================================================================

-- name: InsertMemory
INSERT INTO memories (id, text, tags, source, confidence, created_at, provider, model_id, dim, embedding)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: SelectAllMemories
SELECT id, text, tags, source, confidence, created_at, provider, model_id, dim, embedding
FROM memories
ORDER BY created_at DESC;

-- name: DeleteMemory
DELETE FROM memories WHERE id = ?;

-- name: UpdateMemoryEmbedding
UPDATE memories
SET embedding = ?, model_id = ?, dim = ?, provider = ?
WHERE id = ?;

-- name: SearchMemoriesFTS
SELECT m.id, m.text, m.tags, m.source, m.confidence, m.created_at,
       m.provider, m.model_id, m.dim, m.embedding,
       snippet(memories_fts, 0, '>>>', '<<<', '...', 32) as snippet,
       rank
FROM memories m
JOIN memories_fts fts ON m.rowid = fts.rowid
WHERE memories_fts MATCH ?
ORDER BY rank
LIMIT ?;

-- name: ClearMemories
DELETE FROM memories;

-- ============================================================================
-- HISTORY QUERIES
-- ============================================================================

-- name: InsertHistory
INSERT INTO history (id, role, content, created_at, session_id)
VALUES (?, ?, ?, ?, ?);

-- name: SearchHistoryFTS
SELECT h.id, h.role, h.content, h.created_at, h.session_id,
       snippet(history_fts, 0, '>>>', '<<<', '...', 32) as snippet,
       rank
FROM history h
JOIN history_fts fts ON h.rowid = fts.rowid
WHERE history_fts MATCH ?
ORDER BY rank
LIMIT ?;

-- name: SelectRecentHistory
SELECT id, role, content, created_at, session_id
FROM history
ORDER BY created_at DESC
LIMIT ?;

-- name: ClearHistory
DELETE FROM history;
