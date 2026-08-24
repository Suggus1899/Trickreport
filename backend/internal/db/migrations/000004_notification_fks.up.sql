-- notifications was the only tenant-owned table without FK constraints,
-- leaving orphaned rows behind whenever a tenant or user was deleted.
ALTER TABLE notifications
    ADD CONSTRAINT notifications_tenant_fk FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    ADD CONSTRAINT notifications_user_fk   FOREIGN KEY (user_id)   REFERENCES users(id)   ON DELETE CASCADE;
