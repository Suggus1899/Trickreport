-- Migration: 003_knowledge_base
-- Creates: articles

CREATE TABLE articles (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID          NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title         VARCHAR(255)  NOT NULL,
    content       TEXT          NOT NULL,
    category      VARCHAR(100)  NOT NULL DEFAULT 'general',
    tags          TEXT[]        NOT NULL DEFAULT '{}',
    published     BOOLEAN       NOT NULL DEFAULT FALSE,
    created_by    UUID          NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Indexes for frequent queries
CREATE INDEX idx_articles_tenant_id ON articles(tenant_id);
CREATE INDEX idx_articles_published ON articles(published);
CREATE INDEX idx_articles_category ON articles(category);

-- Trigram index for simple text search (title + content) could be added here if needed,
-- but for Phase 3 we will rely on basic ILIKE or equality on tags/categories.
