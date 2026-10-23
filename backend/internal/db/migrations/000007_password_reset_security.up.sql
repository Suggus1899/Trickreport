-- Invalidate legacy plaintext credentials without removing their audit records.
UPDATE password_reset_tokens SET used = TRUE WHERE used = FALSE;
CREATE INDEX idx_reset_unused ON password_reset_tokens(token) WHERE used = FALSE;
