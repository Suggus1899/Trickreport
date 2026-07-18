-- Migration: 001_init
-- Creates: tenants, users

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

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
-- Role enum
-- ─────────────────────────────────────────────
CREATE TYPE user_role AS ENUM ('admin', 'agent', 'end_user');

-- ─────────────────────────────────────────────
-- Users
-- ─────────────────────────────────────────────
CREATE TABLE users (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    role        user_role    NOT NULL DEFAULT 'end_user',
    password    VARCHAR(255),           -- bcrypt hash; NULL when using SSO/LDAP only
    ldap_dn     VARCHAR(500),           -- DN in LDAP directory, if applicable
    avatar_url  VARCHAR(500),
    active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_tenant_unique UNIQUE (tenant_id, email)
);

-- ─────────────────────────────────────────────
-- Indexes
-- ─────────────────────────────────────────────
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email     ON users(email);
CREATE INDEX idx_users_active    ON users(active);

-- ─────────────────────────────────────────────
-- Seed: default tenant + superadmin
-- (override credentials via env before first run)
-- ─────────────────────────────────────────────
INSERT INTO tenants (id, name, slug, settings)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Default Organization',
    'default',
    '{}'
);

-- password: 'changeme' (bcrypt cost=10) — CHANGE IN PRODUCTION
INSERT INTO users (id, tenant_id, name, email, role, password)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'Super Admin',
    'admin@trickreport.local',
    'admin',
    '$2a$10$FDaOn5b0yomg..3IVVQGkepzPl8hek/LctHP0fNgWbCrK43UHs6ey'
);
