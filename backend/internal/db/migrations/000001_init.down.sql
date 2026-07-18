-- Rollback: 001_init
-- Drops: users, tenants, user_role enum

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
DROP TABLE IF EXISTS tenants;
