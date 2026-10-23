CREATE FUNCTION revoke_sessions_on_identity_change() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF OLD.password IS DISTINCT FROM NEW.password
        OR OLD.role IS DISTINCT FROM NEW.role
        OR OLD.active IS DISTINCT FROM NEW.active
        OR OLD.mfa_enabled IS DISTINCT FROM NEW.mfa_enabled THEN
        UPDATE user_sessions SET revoked = TRUE WHERE user_id = NEW.id AND NOT revoked;
    END IF;
    RETURN NEW;
END;
$$;
CREATE TRIGGER users_revoke_sessions
AFTER UPDATE OF password, role, active, mfa_enabled ON users
FOR EACH ROW EXECUTE FUNCTION revoke_sessions_on_identity_change();
