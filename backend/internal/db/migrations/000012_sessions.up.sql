-- Migration: 012_sessions
-- Creates: user_sessions

CREATE TABLE user_sessions (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    jti          VARCHAR(255) NOT NULL,
    refresh_jti  VARCHAR(255),
    ip_address   VARCHAR(45),
    user_agent   TEXT,
    created_at   TIMESTAMPTZ  DEFAULT NOW(),
    expires_at   TIMESTAMPTZ  NOT NULL,
    revoked      BOOLEAN      DEFAULT FALSE
);

CREATE INDEX idx_user_sessions_user_id      ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_jti          ON user_sessions(jti);
CREATE INDEX idx_user_sessions_refresh_jti  ON user_sessions(refresh_jti);
CREATE INDEX idx_user_sessions_expires_at   ON user_sessions(expires_at);
