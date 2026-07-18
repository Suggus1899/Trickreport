-- Drop all tables, functions, triggers, and enums in reverse dependency order.

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS automation_rules;
DROP TABLE IF EXISTS sla_policies;

DROP TRIGGER IF EXISTS articles_search_vector_trigger ON articles;
DROP FUNCTION IF EXISTS articles_search_vector_update();
DROP TABLE IF EXISTS articles;

DROP TABLE IF EXISTS ticket_attachments;
DROP TABLE IF EXISTS ticket_history;
DROP TABLE IF EXISTS ticket_comments;

DROP TRIGGER IF EXISTS tickets_search_vector_trigger ON tickets;
DROP FUNCTION IF EXISTS tickets_search_vector_update();
DROP TABLE IF EXISTS tickets;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;

DROP TYPE IF EXISTS ticket_priority;
DROP TYPE IF EXISTS ticket_status;
DROP TYPE IF EXISTS user_role;
