-- Etapa 3 (ITSM): activos y CMDB.
-- Un activo es un elemento de configuración (CI): su tipo vive en
-- ci_types y sus relaciones en ci_relationships. Todo aísla por
-- tenant_id y el vínculo activo-ticket vive en ticket_assets.

CREATE TABLE ci_types (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ci_types_tenant_name_unique UNIQUE (tenant_id, name)
);
CREATE INDEX idx_ci_types_tenant ON ci_types(tenant_id);

CREATE TABLE assets (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    ci_type_id   UUID        NOT NULL REFERENCES ci_types(id) ON DELETE RESTRICT,
    name         VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255) NOT NULL DEFAULT '',
    status       VARCHAR(50) NOT NULL DEFAULT 'active'
                 CHECK (status IN ('active', 'maintenance', 'retired', 'stock')),
    attributes   JSONB       NOT NULL DEFAULT '{}',
    owner_id     UUID        REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT assets_tenant_name_unique UNIQUE (tenant_id, name)
);
CREATE INDEX idx_assets_tenant ON assets(tenant_id);
CREATE INDEX idx_assets_ci_type ON assets(ci_type_id);
CREATE INDEX idx_assets_owner ON assets(owner_id);

-- Relaciones dirigidas entre CIs. El servicio valida que ambos
-- extremos pertenezcan al mismo tenant y rechaza auto-relaciones.
CREATE TABLE ci_relationships (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    source_ci_id      UUID        NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    target_ci_id      UUID        NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL
                      CHECK (relationship_type IN ('depends_on', 'hosts', 'uses')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ci_relationships_unique UNIQUE (tenant_id, source_ci_id, target_ci_id, relationship_type),
    CONSTRAINT ci_relationships_no_self CHECK (source_ci_id <> target_ci_id)
);
CREATE INDEX idx_ci_relationships_source ON ci_relationships(source_ci_id);
CREATE INDEX idx_ci_relationships_target ON ci_relationships(target_ci_id);

CREATE TABLE ticket_assets (
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    asset_id  UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (ticket_id, asset_id)
);
CREATE INDEX idx_ticket_assets_asset ON ticket_assets(asset_id);
