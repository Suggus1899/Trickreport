-- Migration: 010_account_lockout
-- Alters: users (adds failed_login_attempts, locked_until)

ALTER TABLE users ADD COLUMN failed_login_attempts INT DEFAULT 0;
ALTER TABLE users ADD COLUMN locked_until TIMESTAMPTZ;
