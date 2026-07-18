-- Migration: 010_account_lockout (rollback)
-- Alters: users (removes failed_login_attempts, locked_until)

ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS failed_login_attempts;
