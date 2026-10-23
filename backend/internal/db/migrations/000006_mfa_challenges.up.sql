CREATE TABLE mfa_challenges (
    jti UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    attempts SMALLINT NOT NULL DEFAULT 0 CHECK (attempts BETWEEN 0 AND 5),
    consumed BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_mfa_challenges_expiry ON mfa_challenges(expires_at);
