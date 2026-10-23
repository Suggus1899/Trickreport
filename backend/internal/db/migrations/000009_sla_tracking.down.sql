DROP INDEX tickets_pending_response_idx;
DROP INDEX tickets_pending_resolution_idx;
DROP INDEX tickets_pending_escalation_deadline_idx;
ALTER TABLE sla_policies DROP CONSTRAINT sla_policy_positive_durations;
ALTER TABLE tickets DROP COLUMN response_deadline, DROP COLUMN escalation_deadline,
    DROP COLUMN first_responded_at, DROP COLUMN resolved_at, DROP COLUMN response_sla_breached;
