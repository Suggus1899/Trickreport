-- Etapa 4 (MSP): operación multicliente del proveedor.
-- Todo aísla por tenant_id (la organización operadora) y, dentro de
-- ella, por client_id: un contacto de un cliente nunca ve datos de
-- otro, porque cada consulta filtra por ambos. El dinero se guarda en
-- céntimos (enteros) para evitar errores de redondeo. No hay cobros:
-- la facturación genera borradores; los pagos automáticos quedan para
-- una etapa posterior, como fija el plan maestro.
-- Zabbix y MeshCentral son adaptadores configurables: aquí solo se
-- guarda la configuración y el hash del token del webhook, nunca
-- credenciales en claro ni tokens de sesión de terceros.

CREATE TABLE msp_clients (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    domain        VARCHAR(255) NOT NULL DEFAULT '',
    city          VARCHAR(255) NOT NULL DEFAULT '',
    status        VARCHAR(50)  NOT NULL DEFAULT 'active'
                  CHECK (status IN ('active', 'inactive', 'suspended')),
    contact_name  VARCHAR(255) NOT NULL DEFAULT '',
    contact_email VARCHAR(255) NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_clients_tenant_name_unique UNIQUE (tenant_id, name),
    CONSTRAINT msp_clients_tenant_domain_unique UNIQUE (tenant_id, domain)
);
CREATE INDEX idx_msp_clients_tenant ON msp_clients(tenant_id);
CREATE INDEX idx_msp_clients_status ON msp_clients(status);

-- Contactos de un cliente. user_id enlaza con un usuario de la
-- organización cuando el contacto accede a Trickreport; sin user_id
-- el contacto es solo un registro (nombre y correo) del cliente.
CREATE TABLE msp_client_contacts (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id  UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    user_id    UUID        REFERENCES users(id) ON DELETE SET NULL,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    role       VARCHAR(50)  NOT NULL DEFAULT 'contact'
               CHECK (role IN ('contact', 'client_admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_client_contacts_client_email_unique UNIQUE (client_id, email)
);
CREATE INDEX idx_msp_client_contacts_tenant ON msp_client_contacts(tenant_id);
CREATE INDEX idx_msp_client_contacts_client ON msp_client_contacts(client_id);
CREATE INDEX idx_msp_client_contacts_user ON msp_client_contacts(user_id);

CREATE TABLE msp_contracts (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id              UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    name                   VARCHAR(255) NOT NULL,
    status                 VARCHAR(50)  NOT NULL DEFAULT 'draft'
                           CHECK (status IN ('draft', 'active', 'suspended', 'expired', 'cancelled')),
    starts_at              TIMESTAMPTZ NOT NULL,
    ends_at                TIMESTAMPTZ,
    coverage               TEXT        NOT NULL DEFAULT '',
    currency               VARCHAR(10) NOT NULL DEFAULT 'EUR',
    monthly_fee_cents      BIGINT      NOT NULL DEFAULT 0 CHECK (monthly_fee_cents >= 0),
    hourly_rate_cents      BIGINT      NOT NULL DEFAULT 0 CHECK (hourly_rate_cents >= 0),
    included_hours_monthly NUMERIC(8,2) NOT NULL DEFAULT 0 CHECK (included_hours_monthly >= 0),
    sla_response_minutes   INT         NOT NULL DEFAULT 240 CHECK (sla_response_minutes >= 0),
    sla_resolution_minutes INT         NOT NULL DEFAULT 1440 CHECK (sla_resolution_minutes >= 0),
    sla_overrides          JSONB       NOT NULL DEFAULT '{}',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_contracts_tenant ON msp_contracts(tenant_id);
CREATE INDEX idx_msp_contracts_client ON msp_contracts(client_id);
CREATE INDEX idx_msp_contracts_status ON msp_contracts(status);

-- Contexto multicliente de un ticket: qué cliente, qué contrato y qué
-- SLA se resolvió al crearlo. Un ticket, como mucho, un contexto.
CREATE TABLE msp_ticket_context (
    ticket_id              UUID PRIMARY KEY REFERENCES tickets(id) ON DELETE CASCADE,
    tenant_id              UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id              UUID NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    contract_id            UUID REFERENCES msp_contracts(id) ON DELETE SET NULL,
    sla_response_minutes   INT  NOT NULL DEFAULT 0 CHECK (sla_response_minutes >= 0),
    sla_resolution_minutes INT  NOT NULL DEFAULT 0 CHECK (sla_resolution_minutes >= 0),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_ticket_context_client ON msp_ticket_context(client_id);
CREATE INDEX idx_msp_ticket_context_contract ON msp_ticket_context(contract_id);

CREATE TABLE msp_technicians (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id           UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skills            TEXT[]      NOT NULL DEFAULT '{}',
    hourly_rate_cents BIGINT      NOT NULL DEFAULT 0 CHECK (hourly_rate_cents >= 0),
    active            BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_technicians_tenant_user_unique UNIQUE (tenant_id, user_id)
);
CREATE INDEX idx_msp_technicians_tenant ON msp_technicians(tenant_id);

CREATE TABLE msp_availability (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    technician_id  UUID NOT NULL REFERENCES msp_technicians(id) ON DELETE CASCADE,
    day_of_week    INT  NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time     TIME NOT NULL,
    end_time       TIME NOT NULL,
    CONSTRAINT msp_availability_no_empty CHECK (start_time < end_time),
    CONSTRAINT msp_availability_unique UNIQUE (tenant_id, technician_id, day_of_week, start_time, end_time)
);
CREATE INDEX idx_msp_availability_technician ON msp_availability(technician_id);

CREATE TABLE msp_interventions (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id     UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    technician_id UUID        NOT NULL REFERENCES msp_technicians(id) ON DELETE CASCADE,
    ticket_id     UUID        REFERENCES tickets(id) ON DELETE SET NULL,
    title         VARCHAR(255) NOT NULL,
    starts_at     TIMESTAMPTZ NOT NULL,
    ends_at       TIMESTAMPTZ NOT NULL,
    status        VARCHAR(50)  NOT NULL DEFAULT 'scheduled'
                  CHECK (status IN ('scheduled', 'in_progress', 'done', 'cancelled')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_interventions_window CHECK (ends_at > starts_at)
);
CREATE INDEX idx_msp_interventions_tenant ON msp_interventions(tenant_id);
CREATE INDEX idx_msp_interventions_client ON msp_interventions(client_id);
CREATE INDEX idx_msp_interventions_technician ON msp_interventions(technician_id);

CREATE TABLE msp_time_entries (
    id                        UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                 UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id                 UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    ticket_id                 UUID        REFERENCES tickets(id) ON DELETE SET NULL,
    intervention_id           UUID        REFERENCES msp_interventions(id) ON DELETE SET NULL,
    technician_id             UUID        NOT NULL REFERENCES msp_technicians(id) ON DELETE CASCADE,
    minutes                   INT         NOT NULL CHECK (minutes > 0 AND minutes <= 1440),
    description               TEXT        NOT NULL DEFAULT '',
    billable                  BOOLEAN     NOT NULL DEFAULT TRUE,
    billed                    BOOLEAN     NOT NULL DEFAULT FALSE,
    hourly_rate_cents         BIGINT      NOT NULL DEFAULT 0 CHECK (hourly_rate_cents >= 0),
    contract_hourly_rate_cents BIGINT     CHECK (contract_hourly_rate_cents IS NULL OR contract_hourly_rate_cents >= 0),
    cost_cents                BIGINT      NOT NULL DEFAULT 0 CHECK (cost_cents >= 0),
    occurred_at               DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_time_entries_tenant ON msp_time_entries(tenant_id);
CREATE INDEX idx_msp_time_entries_client ON msp_time_entries(client_id);
CREATE INDEX idx_msp_time_entries_ticket ON msp_time_entries(ticket_id);
CREATE INDEX idx_msp_time_entries_unbilled ON msp_time_entries(client_id, billed) WHERE billed = FALSE;

CREATE TABLE msp_expenses (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id   UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    ticket_id   UUID        REFERENCES tickets(id) ON DELETE SET NULL,
    project_id  UUID,
    description TEXT        NOT NULL DEFAULT '',
    category    VARCHAR(100) NOT NULL DEFAULT 'other',
    amount_cents BIGINT     NOT NULL CHECK (amount_cents > 0),
    billed      BOOLEAN     NOT NULL DEFAULT FALSE,
    incurred_at DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_expenses_tenant ON msp_expenses(tenant_id);
CREATE INDEX idx_msp_expenses_client ON msp_expenses(client_id);

CREATE TABLE msp_projects (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id    UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    status       VARCHAR(50)  NOT NULL DEFAULT 'planned'
                 CHECK (status IN ('planned', 'active', 'on_hold', 'done', 'cancelled')),
    budget_cents BIGINT      NOT NULL DEFAULT 0 CHECK (budget_cents >= 0),
    starts_at    TIMESTAMPTZ,
    ends_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_projects_window CHECK (ends_at IS NULL OR starts_at IS NULL OR ends_at >= starts_at)
);
CREATE INDEX idx_msp_projects_tenant ON msp_projects(tenant_id);
CREATE INDEX idx_msp_projects_client ON msp_projects(client_id);

CREATE TABLE msp_project_tasks (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    project_id        UUID        NOT NULL REFERENCES msp_projects(id) ON DELETE CASCADE,
    title             VARCHAR(255) NOT NULL,
    status            VARCHAR(50)  NOT NULL DEFAULT 'todo'
                      CHECK (status IN ('todo', 'doing', 'done')),
    technician_id     UUID        REFERENCES msp_technicians(id) ON DELETE SET NULL,
    estimated_minutes INT         NOT NULL DEFAULT 0 CHECK (estimated_minutes >= 0),
    position          INT         NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_project_tasks_project ON msp_project_tasks(project_id);

-- La FK de msp_expenses.project_id se añade aquí porque msp_projects
-- se crea después; así el orden del archivo no condiciona el modelo.
ALTER TABLE msp_expenses
    ADD CONSTRAINT msp_expenses_project_fk
    FOREIGN KEY (project_id) REFERENCES msp_projects(id) ON DELETE SET NULL;
CREATE INDEX idx_msp_expenses_project ON msp_expenses(project_id);

CREATE TABLE msp_billing_drafts (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id    UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    contract_id  UUID        REFERENCES msp_contracts(id) ON DELETE SET NULL,
    period_start DATE        NOT NULL,
    period_end   DATE        NOT NULL,
    status       VARCHAR(50) NOT NULL DEFAULT 'draft'
                 CHECK (status IN ('draft', 'issued', 'void')),
    total_cents  BIGINT      NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    notes        TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_billing_drafts_period CHECK (period_end >= period_start),
    CONSTRAINT msp_billing_drafts_client_period_unique UNIQUE (tenant_id, client_id, period_start, period_end)
);
CREATE INDEX idx_msp_billing_drafts_client ON msp_billing_drafts(client_id);

CREATE TABLE msp_billing_items (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    draft_id     UUID        NOT NULL REFERENCES msp_billing_drafts(id) ON DELETE CASCADE,
    kind         VARCHAR(50) NOT NULL CHECK (kind IN ('time', 'expense', 'fee')),
    description  TEXT        NOT NULL DEFAULT '',
    amount_cents BIGINT      NOT NULL CHECK (amount_cents >= 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_billing_items_draft ON msp_billing_items(draft_id);

-- Inventario por cliente reutilizando los activos de la Etapa 3: el
-- activo sigue viviendo en assets; esta tabla solo dice de qué
-- cliente es dentro del proveedor.
CREATE TABLE msp_client_assets (
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id  UUID NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    asset_id   UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (client_id, asset_id)
);
CREATE INDEX idx_msp_client_assets_asset ON msp_client_assets(asset_id);

CREATE TABLE msp_discovered_assets (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id          UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    external_id        VARCHAR(255) NOT NULL,
    name               VARCHAR(255) NOT NULL,
    kind               VARCHAR(100) NOT NULL DEFAULT 'host',
    payload            JSONB       NOT NULL DEFAULT '{}',
    status             VARCHAR(50)  NOT NULL DEFAULT 'pending'
                       CHECK (status IN ('pending', 'imported', 'ignored')),
    imported_asset_id  UUID        REFERENCES assets(id) ON DELETE SET NULL,
    discovered_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_discovered_assets_unique UNIQUE (tenant_id, client_id, external_id)
);
CREATE INDEX idx_msp_discovered_assets_client ON msp_discovered_assets(client_id);

-- Fuentes de monitorización. webhook_token_hash es el SHA-256
-- hexadecimal del token que Zabbix envía en X-MSP-Webhook-Token; el
-- token en claro no se guarda en ningún sitio.
CREATE TABLE msp_monitoring_sources (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name               VARCHAR(255) NOT NULL,
    kind               VARCHAR(50)  NOT NULL DEFAULT 'zabbix' CHECK (kind = 'zabbix'),
    base_url           VARCHAR(500) NOT NULL DEFAULT '',
    webhook_token_hash VARCHAR(128) NOT NULL DEFAULT '',
    configured         BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT msp_monitoring_sources_tenant_name_unique UNIQUE (tenant_id, name)
);
CREATE INDEX idx_msp_monitoring_sources_tenant ON msp_monitoring_sources(tenant_id);

-- Alertas de monitorización. El índice parcial garantiza la regla de
-- no-duplicación: como mucho una alerta abierta por fuente y recurso.
CREATE TABLE msp_monitoring_alerts (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    source_id         UUID        NOT NULL REFERENCES msp_monitoring_sources(id) ON DELETE CASCADE,
    client_id         UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    external_resource VARCHAR(255) NOT NULL,
    severity          VARCHAR(50)  NOT NULL DEFAULT 'warning'
                      CHECK (severity IN ('info', 'warning', 'average', 'high', 'disaster')),
    status            VARCHAR(50)  NOT NULL DEFAULT 'open'
                      CHECK (status IN ('open', 'resolved')),
    title             VARCHAR(255) NOT NULL DEFAULT '',
    message           TEXT        NOT NULL DEFAULT '',
    ticket_id         UUID        REFERENCES tickets(id) ON DELETE SET NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_msp_alerts_open_dedup
    ON msp_monitoring_alerts(tenant_id, source_id, external_resource)
    WHERE status = 'open';
CREATE INDEX idx_msp_monitoring_alerts_client ON msp_monitoring_alerts(client_id);

-- Sesiones de acceso remoto (MeshCentral). Trickreport no abre la
-- sesión: registra quién la pidió, quién la autorizó y cuándo empezó
-- y terminó, vinculada a cliente, ticket y técnico.
CREATE TABLE msp_remote_sessions (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id           UUID        NOT NULL REFERENCES msp_clients(id) ON DELETE CASCADE,
    ticket_id           UUID        REFERENCES tickets(id) ON DELETE SET NULL,
    technician_id       UUID        NOT NULL REFERENCES msp_technicians(id) ON DELETE CASCADE,
    provider            VARCHAR(50)  NOT NULL DEFAULT 'meshcentral' CHECK (provider = 'meshcentral'),
    external_session_id VARCHAR(255) NOT NULL DEFAULT '',
    status              VARCHAR(50)  NOT NULL DEFAULT 'requested'
                        CHECK (status IN ('requested', 'authorized', 'active', 'ended', 'denied')),
    authorized_by       UUID        REFERENCES users(id) ON DELETE SET NULL,
    audit_note          TEXT        NOT NULL DEFAULT '',
    started_at          TIMESTAMPTZ,
    ended_at            TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_msp_remote_sessions_client ON msp_remote_sessions(client_id);
CREATE INDEX idx_msp_remote_sessions_technician ON msp_remote_sessions(technician_id);
