-- Migration: 013_optimistic_concurrency (down)
-- Removes the version column from tickets.

ALTER TABLE tickets DROP COLUMN IF EXISTS version;
