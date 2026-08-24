-- Tracks when a ticket was escalated (SLA policy's escalation_minutes elapsed
-- without resolution), so the worker only fires the "escalation" automation
-- trigger once per ticket. Mirrors sla_breached's own boolean guard, but as a
-- timestamp since "when" is also useful for reporting later.
ALTER TABLE tickets ADD COLUMN escalated_at TIMESTAMPTZ NULL;
