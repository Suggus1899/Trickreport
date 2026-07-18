-- Migration: 006_sla_breached
-- Adds index for worker SLA breach tracking.
-- The sla_breached column is created in 005_automations.sql.

CREATE INDEX IF NOT EXISTS idx_tickets_sla_breached ON tickets(sla_breached) WHERE sla_breached = FALSE;
