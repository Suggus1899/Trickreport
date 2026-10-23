ALTER TABLE tickets
    ADD COLUMN response_deadline TIMESTAMPTZ,
    ADD COLUMN escalation_deadline TIMESTAMPTZ,
    ADD COLUMN first_responded_at TIMESTAMPTZ,
    ADD COLUMN resolved_at TIMESTAMPTZ,
    ADD COLUMN response_sla_breached BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX tickets_pending_response_idx ON tickets (response_deadline)
    WHERE first_responded_at IS NULL AND response_sla_breached = FALSE;
CREATE INDEX tickets_pending_resolution_idx ON tickets (sla_deadline)
    WHERE sla_breached = FALSE;
CREATE INDEX tickets_pending_escalation_deadline_idx ON tickets (escalation_deadline)
    WHERE escalated_at IS NULL;

ALTER TABLE sla_policies ADD CONSTRAINT sla_policy_positive_durations
    CHECK (response_time_minutes > 0 AND resolution_time_minutes > 0 AND escalation_minutes >= 0) NOT VALID;
