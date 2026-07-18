-- Migration: 001_init (consolidated)
-- Creates the full Trickreport schema in one migration.
-- Consolidated from the original 14 incremental migrations.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ─────────────────────────────────────────────
-- Enums
-- ─────────────────────────────────────────────
CREATE TYPE user_role AS ENUM ('admin', 'agent', 'end_user');
CREATE TYPE ticket_status AS ENUM ('open', 'in_progress', 'waiting_client', 'resolved', 'closed');
CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high', 'critical');

-- ─────────────────────────────────────────────
-- Tenants
-- ─────────────────────────────────────────────
CREATE TABLE tenants (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    logo_url   VARCHAR(500),
    settings   JSONB        NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ─────────────────────────────────────────────
-- Users
-- ─────────────────────────────────────────────
CREATE TABLE users (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name                   VARCHAR(255) NOT NULL,
    email                  VARCHAR(255) NOT NULL,
    role                   user_role    NOT NULL DEFAULT 'end_user',
    password               VARCHAR(255),           -- bcrypt hash; NULL when using SSO/LDAP only
    ldap_dn                VARCHAR(500),           -- DN in LDAP directory, if applicable
    avatar_url             VARCHAR(500),
    active                 BOOLEAN      NOT NULL DEFAULT TRUE,
    failed_login_attempts  INT          DEFAULT 0,
    locked_until           TIMESTAMPTZ,
    mfa_secret             VARCHAR(255),
    mfa_enabled            BOOLEAN      DEFAULT FALSE,
    created_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_tenant_unique UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email     ON users(email);
CREATE INDEX idx_users_active    ON users(active);

-- ─────────────────────────────────────────────
-- Tickets
-- ─────────────────────────────────────────────
CREATE TABLE tickets (
    id            UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID            NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title         VARCHAR(255)    NOT NULL,
    description   TEXT            NOT NULL,
    status        ticket_status   NOT NULL DEFAULT 'open',
    priority      ticket_priority NOT NULL DEFAULT 'low',
    category      VARCHAR(100)    NOT NULL DEFAULT 'general',
    created_by    UUID            NOT NULL REFERENCES users(id),
    assigned_to   UUID            REFERENCES users(id),
    sla_deadline  TIMESTAMPTZ,
    sla_breached  BOOLEAN         NOT NULL DEFAULT FALSE,
    version       INT             NOT NULL DEFAULT 0,
    search_vector tsvector,
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tickets_tenant_id    ON tickets(tenant_id);
CREATE INDEX idx_tickets_status       ON tickets(status);
CREATE INDEX idx_tickets_assigned_to  ON tickets(assigned_to);
CREATE INDEX idx_tickets_created_by   ON tickets(created_by);
CREATE INDEX idx_tickets_sla_breached ON tickets(sla_breached) WHERE sla_breached = FALSE;
CREATE INDEX tickets_search_idx       ON tickets USING gin(search_vector);

-- Full-text search trigger for tickets
CREATE OR REPLACE FUNCTION tickets_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tickets_search_vector_trigger BEFORE INSERT OR UPDATE ON tickets
FOR EACH ROW EXECUTE FUNCTION tickets_search_vector_update();

-- ─────────────────────────────────────────────
-- Ticket Comments
-- ─────────────────────────────────────────────
CREATE TABLE ticket_comments (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id     UUID        NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id       UUID        NOT NULL REFERENCES users(id),
    content       TEXT        NOT NULL,
    is_internal   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ticket_comments_ticket_id ON ticket_comments(ticket_id);

-- ─────────────────────────────────────────────
-- Ticket History (Audit Log)
-- ─────────────────────────────────────────────
CREATE TABLE ticket_history (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id     UUID        NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id       UUID        NOT NULL REFERENCES users(id),
    field         VARCHAR(50) NOT NULL,
    old_value     TEXT,
    new_value     TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ticket_history_ticket_id ON ticket_history(ticket_id);

-- ─────────────────────────────────────────────
-- Ticket Attachments
-- ─────────────────────────────────────────────
CREATE TABLE ticket_attachments (
    id           UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id    UUID          NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id      UUID          NOT NULL REFERENCES users(id),
    filename     VARCHAR(255)  NOT NULL,
    content_type VARCHAR(100)  NOT NULL,
    file_size    BIGINT        NOT NULL,
    file_data    BYTEA         NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX ticket_attachments_ticket_idx ON ticket_attachments(ticket_id);

-- ─────────────────────────────────────────────
-- Articles (Knowledge Base)
-- ─────────────────────────────────────────────
CREATE TABLE articles (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID          NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title         VARCHAR(255)  NOT NULL,
    content       TEXT          NOT NULL,
    category      VARCHAR(100)  NOT NULL DEFAULT 'general',
    tags          TEXT[]        NOT NULL DEFAULT '{}',
    published     BOOLEAN       NOT NULL DEFAULT FALSE,
    created_by    UUID          NOT NULL REFERENCES users(id),
    search_vector tsvector,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_articles_tenant_id  ON articles(tenant_id);
CREATE INDEX idx_articles_published ON articles(published);
CREATE INDEX idx_articles_category   ON articles(category);
CREATE INDEX articles_search_idx     ON articles USING gin(search_vector);

-- Full-text search trigger for articles
CREATE OR REPLACE FUNCTION articles_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER articles_search_vector_trigger BEFORE INSERT OR UPDATE ON articles
FOR EACH ROW EXECUTE FUNCTION articles_search_vector_update();

-- ─────────────────────────────────────────────
-- SLA Policies
-- ─────────────────────────────────────────────
CREATE TABLE sla_policies (
    id                      UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID            NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    priority                ticket_priority NOT NULL,
    response_time_minutes   INTEGER         NOT NULL DEFAULT 60,
    resolution_time_minutes INTEGER         NOT NULL DEFAULT 1440,
    escalation_minutes      INTEGER         NOT NULL DEFAULT 120,
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, priority)
);

-- ─────────────────────────────────────────────
-- Automation Rules
-- ─────────────────────────────────────────────
CREATE TABLE automation_rules (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID          NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name          VARCHAR(255)  NOT NULL,
    description   TEXT,
    trigger_type  VARCHAR(50)   NOT NULL,
    conditions    JSONB         NOT NULL DEFAULT '{}'::jsonb,
    actions       JSONB         NOT NULL DEFAULT '[]'::jsonb,
    is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_automation_rules_tenant_active ON automation_rules(tenant_id, is_active);

-- ─────────────────────────────────────────────
-- Password Reset Tokens
-- ─────────────────────────────────────────────
CREATE TABLE password_reset_tokens (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ  NOT NULL,
    used       BOOLEAN      DEFAULT FALSE,
    created_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tokens_token      ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_tokens_user_id    ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);

-- ─────────────────────────────────────────────
-- User Sessions
-- ─────────────────────────────────────────────
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

-- ─────────────────────────────────────────────
-- Notifications
-- ─────────────────────────────────────────────
CREATE TABLE notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    user_id     UUID NOT NULL,
    title       VARCHAR(255) NOT NULL,
    body        TEXT NOT NULL DEFAULT '',
    type        VARCHAR(50) NOT NULL DEFAULT 'info',
    ref_id      UUID,
    ref_type    VARCHAR(50),
    read        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_unread ON notifications (user_id, read) WHERE read = FALSE;
CREATE INDEX idx_notifications_tenant      ON notifications (tenant_id);
CREATE INDEX idx_notifications_created     ON notifications (created_at DESC);

-- ─────────────────────────────────────────────
-- Seed: default tenant only.
-- The initial admin user is created at server startup from
-- ADMIN_EMAIL / ADMIN_PASSWORD env vars (see internal/bootstrap).
-- ─────────────────────────────────────────────
INSERT INTO tenants (id, name, slug, settings)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Default Organization',
    'default',
    '{}'
)
ON CONFLICT (id) DO NOTHING;
