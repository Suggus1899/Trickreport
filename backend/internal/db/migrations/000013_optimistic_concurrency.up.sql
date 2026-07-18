-- Migration: 013_optimistic_concurrency
-- Adds a version column to tickets for optimistic concurrency control.

ALTER TABLE tickets ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 0;
