-- Migration: 005_automations
-- Creates: automation_rules and updates tickets table

CREATE TABLE automation_rules (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID          NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name          VARCHAR(255)  NOT NULL,
    description   TEXT,
    trigger_type  VARCHAR(50)   NOT NULL, -- 'sla_breach', 'ticket_created', 'status_changed'
    conditions    JSONB         NOT NULL DEFAULT '{}'::jsonb,
    actions       JSONB         NOT NULL DEFAULT '[]'::jsonb,
    is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_automation_rules_tenant_active ON automation_rules(tenant_id, is_active);

-- Add flag to tickets to prevent spamming notifications for the same breach
ALTER TABLE tickets ADD COLUMN sla_breached BOOLEAN NOT NULL DEFAULT FALSE;
