-- Migration: 004_sla_policies
-- Creates: sla_policies

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

-- Seed default SLAs for the initial tenant (optional, they can also be created via UI)
-- By default, inserting a set of basic policies for 'default' tenant if we want.
-- We'll just let the UI handle upserts for now.
