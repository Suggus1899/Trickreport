-- Existing tokens must be replaced by a newly persisted access/refresh pair.
UPDATE user_sessions SET revoked = TRUE;
CREATE UNIQUE INDEX idx_sessions_access_jti ON user_sessions(jti) WHERE revoked = FALSE;
CREATE UNIQUE INDEX idx_sessions_refresh_jti ON user_sessions(refresh_jti) WHERE revoked = FALSE;
