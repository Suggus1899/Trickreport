-- Etapa 5 (Atención al cliente): portal, correo, chat y WhatsApp;
-- bandeja, asignación, historial, conocimiento y satisfacción.
-- Todo aísla por tenant_id (la empresa) y, dentro de ella, por
-- customer_id: un cliente externo solo ve sus propias conversaciones,
-- porque cada consulta filtra por ambos. Las notas internas
-- (cx_messages.is_internal) nunca se exponen al cliente.
-- Correo y WhatsApp guardan solo configuración: hosts, direcciones e
-- identificadores, y el hash SHA-256 del token de verificación del
-- webhook. Ninguna contraseña SMTP/IMAP ni token de acceso de Meta se
-- guarda en estas tablas; viven en la configuración del operador
-- (ver docs/cx-integraciones.md). La activación viva de ambos canales
-- queda pendiente de credenciales e instancias reales.

CREATE TABLE cx_customers (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name              VARCHAR(255) NOT NULL,
    email             VARCHAR(255) NOT NULL DEFAULT '',
    phone             VARCHAR(100) NOT NULL DEFAULT '',
    company           VARCHAR(255) NOT NULL DEFAULT '',
    status            VARCHAR(50)  NOT NULL DEFAULT 'active'
                      CHECK (status IN ('active', 'inactive')),
    consent_given     BOOLEAN      NOT NULL DEFAULT FALSE,
    erasure_requested BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Correos reales únicos por empresa; el correo vacío queda fuera del
-- índice para que la ingesta pueda crear clientes solo con teléfono.
CREATE UNIQUE INDEX idx_cx_customers_tenant_email_unique
    ON cx_customers(tenant_id, email) WHERE email <> '';
CREATE INDEX idx_cx_customers_tenant ON cx_customers(tenant_id);
CREATE INDEX idx_cx_customers_phone ON cx_customers(tenant_id, phone);

-- Contactos de un cliente externo. user_id enlaza con un usuario de la
-- organización cuando el contacto accede a Trickreport; sin user_id
-- el contacto es solo un registro (nombre y correo) del cliente.
CREATE TABLE cx_customer_contacts (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id UUID        NOT NULL REFERENCES cx_customers(id) ON DELETE CASCADE,
    user_id     UUID        REFERENCES users(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    role        VARCHAR(50)  NOT NULL DEFAULT 'contact'
                CHECK (role IN ('contact', 'client_admin')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT cx_customer_contacts_customer_email_unique UNIQUE (customer_id, email)
);
CREATE INDEX idx_cx_customer_contacts_tenant ON cx_customer_contacts(tenant_id);
CREATE INDEX idx_cx_customer_contacts_customer ON cx_customer_contacts(customer_id);
CREATE INDEX idx_cx_customer_contacts_user ON cx_customer_contacts(user_id);

CREATE TABLE cx_conversations (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id            UUID        NOT NULL REFERENCES cx_customers(id) ON DELETE CASCADE,
    channel                VARCHAR(50)  NOT NULL
                           CHECK (channel IN ('email', 'portal', 'chat', 'whatsapp')),
    subject                VARCHAR(255) NOT NULL,
    status                 VARCHAR(50)  NOT NULL DEFAULT 'open'
                           CHECK (status IN ('open', 'in_progress', 'waiting_customer', 'resolved', 'closed')),
    priority               VARCHAR(50)  NOT NULL DEFAULT 'medium'
                           CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    assigned_to            UUID        REFERENCES users(id) ON DELETE SET NULL,
    created_by             UUID        REFERENCES users(id) ON DELETE SET NULL,
    external_thread_id     VARCHAR(255) NOT NULL DEFAULT '',
    sla_response_minutes   INT         NOT NULL DEFAULT 0 CHECK (sla_response_minutes >= 0),
    sla_resolution_minutes INT         NOT NULL DEFAULT 0 CHECK (sla_resolution_minutes >= 0),
    response_due_at        TIMESTAMPTZ,
    resolution_due_at      TIMESTAMPTZ,
    first_response_at      TIMESTAMPTZ,
    escalated              BOOLEAN      NOT NULL DEFAULT FALSE,
    escalated_at           TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_conversations_tenant ON cx_conversations(tenant_id);
CREATE INDEX idx_cx_conversations_customer ON cx_conversations(customer_id);
CREATE INDEX idx_cx_conversations_assigned ON cx_conversations(assigned_to);
CREATE INDEX idx_cx_conversations_status ON cx_conversations(status);
-- Una conversación abierta por hilo externo y cliente: la ingesta de
-- correo/WhatsApp reutiliza la conversación abierta en vez de duplicarla.
CREATE UNIQUE INDEX idx_cx_conversations_open_thread
    ON cx_conversations(tenant_id, customer_id, external_thread_id)
    WHERE status NOT IN ('resolved', 'closed') AND external_thread_id <> '';

CREATE TABLE cx_messages (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    conversation_id     UUID        NOT NULL REFERENCES cx_conversations(id) ON DELETE CASCADE,
    sender_type         VARCHAR(50)  NOT NULL
                        CHECK (sender_type IN ('customer', 'agent', 'system')),
    sender_id           UUID        REFERENCES users(id) ON DELETE SET NULL,
    body                TEXT        NOT NULL,
    is_internal         BOOLEAN      NOT NULL DEFAULT FALSE,
    channel             VARCHAR(50)  NOT NULL DEFAULT 'portal'
                        CHECK (channel IN ('email', 'portal', 'chat', 'whatsapp')),
    external_message_id VARCHAR(255) NOT NULL DEFAULT '',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_messages_conversation ON cx_messages(conversation_id);
-- Idempotencia de la ingesta: un mismo mensaje externo no se guarda dos veces.
CREATE UNIQUE INDEX idx_cx_messages_external_unique
    ON cx_messages(tenant_id, external_message_id)
    WHERE external_message_id <> '';

-- Configuración de correo: hosts y direcciones; sin contraseñas.
CREATE TABLE cx_email_accounts (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    smtp_host     VARCHAR(500) NOT NULL DEFAULT '',
    smtp_port     INT         NOT NULL DEFAULT 587 CHECK (smtp_port BETWEEN 0 AND 65535),
    smtp_username VARCHAR(255) NOT NULL DEFAULT '',
    smtp_from     VARCHAR(255) NOT NULL DEFAULT '',
    imap_host     VARCHAR(500) NOT NULL DEFAULT '',
    imap_port     INT         NOT NULL DEFAULT 993 CHECK (imap_port BETWEEN 0 AND 65535),
    imap_username VARCHAR(255) NOT NULL DEFAULT '',
    active        BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_email_accounts_tenant ON cx_email_accounts(tenant_id);

-- Configuración de WhatsApp Business Platform: identificadores y solo
-- el hash del token de verificación del webhook; sin token de acceso.
CREATE TABLE cx_whatsapp_accounts (
    id                        UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                 UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name                      VARCHAR(255) NOT NULL,
    business_id               VARCHAR(255) NOT NULL DEFAULT '',
    phone_number_id           VARCHAR(255) NOT NULL DEFAULT '',
    display_name              VARCHAR(255) NOT NULL DEFAULT '',
    webhook_verify_token_hash VARCHAR(255) NOT NULL DEFAULT '',
    active                    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_whatsapp_accounts_tenant ON cx_whatsapp_accounts(tenant_id);

-- Una política por canal y prioridad: gana sobre los valores por
-- defecto del dominio.
CREATE TABLE cx_sla_policies (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    channel            VARCHAR(50)  NOT NULL
                       CHECK (channel IN ('email', 'portal', 'chat', 'whatsapp')),
    priority           VARCHAR(50)  NOT NULL
                       CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    response_minutes   INT         NOT NULL DEFAULT 0 CHECK (response_minutes >= 0),
    resolution_minutes INT         NOT NULL DEFAULT 0 CHECK (resolution_minutes >= 0),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT cx_sla_policies_tenant_channel_priority_unique UNIQUE (tenant_id, channel, priority)
);
CREATE INDEX idx_cx_sla_policies_tenant ON cx_sla_policies(tenant_id);

-- Reglas de asignación automática, evaluadas en orden de posición:
-- la primera que coincide con canal/prioridad asigna la conversación.
CREATE TABLE cx_assignment_rules (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    channel          VARCHAR(50)  NOT NULL DEFAULT '',
    priority         VARCHAR(50)  NOT NULL DEFAULT '',
    assignee_user_id UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position         INT         NOT NULL DEFAULT 0,
    active           BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_assignment_rules_tenant ON cx_assignment_rules(tenant_id);

-- Respuestas automáticas por palabra clave o por canal.
CREATE TABLE cx_auto_replies (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    trigger_kind  VARCHAR(50)  NOT NULL CHECK (trigger_kind IN ('keyword', 'channel')),
    keyword       VARCHAR(255) NOT NULL DEFAULT '',
    channel       VARCHAR(50)  NOT NULL DEFAULT '',
    response_text TEXT        NOT NULL,
    active        BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_auto_replies_tenant ON cx_auto_replies(tenant_id);

-- Una encuesta por conversación, solo tras resolver o cerrar.
CREATE TABLE cx_surveys (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    conversation_id UUID        NOT NULL REFERENCES cx_conversations(id) ON DELETE CASCADE,
    customer_id     UUID        NOT NULL REFERENCES cx_customers(id) ON DELETE CASCADE,
    rating          INT         NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment         TEXT        NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT cx_surveys_conversation_unique UNIQUE (conversation_id)
);
CREATE INDEX idx_cx_surveys_tenant ON cx_surveys(tenant_id);
