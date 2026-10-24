-- Etapa 3 (ITSM): problemas y cambios.
-- Un problema agrupa varios incidentes (problem_incidents); un cambio
-- puede corregir un problema y tocar varios tickets (change_tickets).
-- La aprobación de un cambio reutiliza la tabla approvals con
-- target_type = 'change' (migración 000012).

CREATE TABLE problems (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    status      VARCHAR(50) NOT NULL DEFAULT 'open'
                CHECK (status IN ('open', 'investigating', 'resolved', 'closed')),
    root_cause  TEXT        NOT NULL DEFAULT '',
    workaround  TEXT        NOT NULL DEFAULT '',
    created_by  UUID        NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_problems_tenant ON problems(tenant_id);
CREATE INDEX idx_problems_status ON problems(status);

CREATE TABLE problem_incidents (
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    ticket_id  UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (problem_id, ticket_id)
);
CREATE INDEX idx_problem_incidents_ticket ON problem_incidents(ticket_id);

CREATE TABLE changes (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title         VARCHAR(255) NOT NULL,
    description   TEXT        NOT NULL DEFAULT '',
    problem_id    UUID        REFERENCES problems(id) ON DELETE SET NULL,
    status        VARCHAR(50) NOT NULL DEFAULT 'draft'
                  CHECK (status IN ('draft', 'planned', 'awaiting_approval', 'approved', 'in_progress', 'completed', 'closed', 'rejected')),
    risk          VARCHAR(20) NOT NULL DEFAULT 'medium'
                  CHECK (risk IN ('low', 'medium', 'high', 'critical')),
    planned_start TIMESTAMPTZ,
    planned_end   TIMESTAMPTZ,
    implemented_at TIMESTAMPTZ,
    closed_at     TIMESTAMPTZ,
    created_by    UUID        NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT changes_window_order CHECK (planned_end IS NULL OR planned_start IS NULL OR planned_end >= planned_start)
);
CREATE INDEX idx_changes_tenant ON changes(tenant_id);
CREATE INDEX idx_changes_status ON changes(status);
CREATE INDEX idx_changes_problem ON changes(problem_id);

CREATE TABLE change_tickets (
    change_id  UUID NOT NULL REFERENCES changes(id) ON DELETE CASCADE,
    ticket_id  UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (change_id, ticket_id)
);
CREATE INDEX idx_change_tickets_ticket ON change_tickets(ticket_id);
