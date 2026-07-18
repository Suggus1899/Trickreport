-- Migration: 009_password_reset (rollback)
-- Drops: password_reset_tokens

DROP TABLE IF EXISTS password_reset_tokens;
