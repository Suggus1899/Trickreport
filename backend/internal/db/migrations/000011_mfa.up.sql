-- Migration: 011_mfa
-- Alters: users (adds mfa_secret, mfa_enabled)

ALTER TABLE users ADD COLUMN mfa_secret VARCHAR(255);
ALTER TABLE users ADD COLUMN mfa_enabled BOOLEAN DEFAULT FALSE;
