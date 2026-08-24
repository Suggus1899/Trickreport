-- Mirrors idx_tickets_sla_breached: the escalation worker polls this exact
-- shape (WHERE escalated_at IS NULL) every tick, so without a partial index
-- it degrades to a full table scan as tickets grows.
CREATE INDEX idx_tickets_escalated_at ON tickets(escalated_at) WHERE escalated_at IS NULL;
