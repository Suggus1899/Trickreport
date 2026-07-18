-- Rollback: 005_automations
-- Drops: automation_rules, removes sla_breached column from tickets

ALTER TABLE tickets DROP COLUMN IF EXISTS sla_breached;
DROP TABLE IF EXISTS automation_rules;
