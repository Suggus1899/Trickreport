-- Migration: 011_mfa (rollback)
-- Alters: users (removes mfa_secret, mfa_enabled)

ALTER TABLE users DROP COLUMN IF EXISTS mfa_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS mfa_secret;
