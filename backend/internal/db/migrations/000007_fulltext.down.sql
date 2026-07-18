-- Migration: 007_fulltext (down)

DROP TRIGGER IF EXISTS articles_search_vector_trigger ON articles;
DROP FUNCTION IF EXISTS articles_search_vector_update();
DROP INDEX IF EXISTS articles_search_idx;
ALTER TABLE articles DROP COLUMN IF EXISTS search_vector;

DROP TRIGGER IF EXISTS tickets_search_vector_trigger ON tickets;
DROP FUNCTION IF EXISTS tickets_search_vector_update();
DROP INDEX IF EXISTS tickets_search_idx;
ALTER TABLE tickets DROP COLUMN IF EXISTS search_vector;
