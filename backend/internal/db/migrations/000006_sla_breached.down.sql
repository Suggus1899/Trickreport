-- Rollback: 006_sla_breached
-- Drops the partial index on tickets.sla_breached

DROP INDEX IF EXISTS idx_tickets_sla_breached;
