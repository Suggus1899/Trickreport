-- Etapa 3 (ITSM): catálogo de servicios y aprobaciones.
-- Todas las tablas aíslan por tenant_id; el borrado de un tenant
-- arrastra sus filas (ON DELETE CASCADE), igual que en 000001.

CREATE TABLE catalog_categories (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT catalog_categories_tenant_name_unique UNIQUE (tenant_id, name)
);
CREATE INDEX idx_catalog_categories_tenant ON catalog_categories(tenant_id);

CREATE TABLE catalog_services (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    category_id        UUID        REFERENCES catalog_categories(id) ON DELETE SET NULL,
    name               VARCHAR(255) NOT NULL,
    description        TEXT        NOT NULL DEFAULT '',
    requires_approval  BOOLEAN     NOT NULL DEFAULT FALSE,
    active             BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT catalog_services_tenant_name_unique UNIQUE (tenant_id, name)
);
CREATE INDEX idx_catalog_services_tenant ON catalog_services(tenant_id);
CREATE INDEX idx_catalog_services_category ON catalog_services(category_id);

-- Aprobaciones genéricas: un mismo flujo sirve a solicitudes del
-- catálogo, tickets y cambios (target_type + target_id). El aprobador
-- puede quedar sin asignar (NULL): en ese caso solo personal
-- (admin/agent) puede decidir; nunca un solicitante cualquiera.
CREATE TABLE approvals (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    target_type  VARCHAR(50) NOT NULL CHECK (target_type IN ('ticket', 'catalog_service', 'change')),
    target_id    UUID        NOT NULL,
    requester_id UUID        NOT NULL REFERENCES users(id),
    approver_id  UUID        REFERENCES users(id),
    status       VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    comment      TEXT        NOT NULL DEFAULT '',
    decided_by   UUID        REFERENCES users(id),
    decided_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_approvals_tenant ON approvals(tenant_id);
CREATE INDEX idx_approvals_target ON approvals(tenant_id, target_type, target_id);
CREATE INDEX idx_approvals_approver ON approvals(approver_id) WHERE status = 'pending';
-- Como máximo una aprobación pendiente por (destino, aprobador
-- asignado). Las no asignadas (approver_id NULL) no entran en el
-- índice parcial porque NULL nunca es igual a NULL.
CREATE UNIQUE INDEX idx_approvals_pending_unique
    ON approvals(tenant_id, target_type, target_id, approver_id)
    WHERE status = 'pending' AND approver_id IS NOT NULL;

-- Historial inmutable de decisiones sobre cada aprobación.
CREATE TABLE approval_history (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    approval_id UUID        NOT NULL REFERENCES approvals(id) ON DELETE CASCADE,
    actor_id    UUID        NOT NULL REFERENCES users(id),
    from_status VARCHAR(20) NOT NULL,
    to_status   VARCHAR(20) NOT NULL,
    comment     TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_approval_history_approval ON approval_history(approval_id);
