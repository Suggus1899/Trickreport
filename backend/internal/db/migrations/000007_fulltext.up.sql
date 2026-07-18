-- Migration: 007_fulltext
-- Adds full-text search vectors to tickets and articles

-- ─────────────────────────────────────────────
-- Tickets
-- ─────────────────────────────────────────────
ALTER TABLE tickets ADD COLUMN search_vector tsvector;
CREATE INDEX tickets_search_idx ON tickets USING gin(search_vector);

CREATE OR REPLACE FUNCTION tickets_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tickets_search_vector_trigger BEFORE INSERT OR UPDATE ON tickets
FOR EACH ROW EXECUTE FUNCTION tickets_search_vector_update();

-- Backfill existing rows
UPDATE tickets SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(description, ''));

-- ─────────────────────────────────────────────
-- Articles
-- ─────────────────────────────────────────────
ALTER TABLE articles ADD COLUMN search_vector tsvector;
CREATE INDEX articles_search_idx ON articles USING gin(search_vector);

CREATE OR REPLACE FUNCTION articles_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER articles_search_vector_trigger BEFORE INSERT OR UPDATE ON articles
FOR EACH ROW EXECUTE FUNCTION articles_search_vector_update();

-- Backfill existing rows
UPDATE articles SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content, ''));
