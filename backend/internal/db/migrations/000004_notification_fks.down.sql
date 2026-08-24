ALTER TABLE notifications
    DROP CONSTRAINT IF EXISTS notifications_tenant_fk,
    DROP CONSTRAINT IF EXISTS notifications_user_fk;
