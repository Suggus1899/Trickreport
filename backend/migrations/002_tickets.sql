-- Migration: 002_tickets
-- Creates: tickets, ticket_comments, ticket_history

-- ─────────────────────────────────────────────
-- Enums
-- ─────────────────────────────────────────────
CREATE TYPE ticket_status AS ENUM ('open', 'in_progress', 'waiting_client', 'resolved', 'closed');
CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high', 'critical');

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
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Indexes for frequent queries
CREATE INDEX idx_tickets_tenant_id ON tickets(tenant_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_assigned_to ON tickets(assigned_to);
CREATE INDEX idx_tickets_created_by ON tickets(created_by);

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
    field         VARCHAR(50) NOT NULL, -- e.g., 'status', 'priority', 'assigned_to'
    old_value     TEXT,
    new_value     TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ticket_history_ticket_id ON ticket_history(ticket_id);
