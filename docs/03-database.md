# Modelo de Datos — PostgreSQL

## Multi-tenant

Todas las tablas principales tienen `tenant_id` (UUID, FK → tenants) excepto la propia tabla `tenants`.

### tenants
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| name | VARCHAR | |
| slug | VARCHAR | Identificador único para URL/subdominio |
| logo_url | VARCHAR | |
| settings | JSONB | Configuración por tenant |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

## Entidades Principales

### users
| Columna | Tipo | Descripción |
|---------|------|-------------|
| tenant_id | UUID | FK → tenants |
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| name | VARCHAR | |
| email | VARCHAR | Único |
| role | ENUM | admin, agent, end_user |
| ldap_dn | VARCHAR | DN en LDAP (si aplica) |
| avatar_url | VARCHAR | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### tickets
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| tenant_id | UUID | FK → tenants |
| title | VARCHAR | |
| description | TEXT | |
| status | ENUM | open, in_progress, waiting_client, resolved, closed |
| priority | ENUM | low, medium, high, critical |
| category | VARCHAR | |
| created_by | UUID | FK → users |
| assigned_to | UUID | FK → users (nullable) |
| sla_deadline | TIMESTAMP | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### ticket_comments
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| ticket_id | UUID | FK → tickets |
| user_id | UUID | FK → users |
| content | TEXT | |
| is_internal | BOOLEAN | Solo visible para agentes |
| created_at | TIMESTAMP | |

### ticket_history
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| ticket_id | UUID | FK → tickets |
| user_id | UUID | FK → users |
| field | VARCHAR | Qué cambió |
| old_value | TEXT | |
| new_value | TEXT | |
| created_at | TIMESTAMP | |

### sla_policies
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| priority | ENUM | low, medium, high, critical |
| response_time_minutes | INT | |
| resolution_time_minutes | INT | |
| escalation_minutes | INT | |

### ticket_attachments
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| ticket_id | UUID | FK → tickets |
| file_name | VARCHAR | |
| file_size | INT | Bytes |
| mime_type | VARCHAR | |
| storage_path | VARCHAR | Ruta en disco/S3 |
| uploaded_by | UUID | FK → users |
| created_at | TIMESTAMP | |

### automation_rules
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| tenant_id | UUID | FK → tenants |
| name | VARCHAR | |
| trigger_type | ENUM | on_create, on_status_change, on_priority_change |
| conditions | JSONB | {field: "category", operator: "equals", value: "network"} |
| actions | JSONB | [{action: "assign_to", value: "user_id"}, {action: "set_priority", value: "high"}] |
| enabled | BOOLEAN | |
| created_by | UUID | FK → users |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### knowledge_articles
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| tenant_id | UUID | FK → tenants |
| title | VARCHAR | |
| content | TEXT | |
| tags | TEXT[] | |
| category | VARCHAR | |
| created_by | UUID | FK → users |
| published | BOOLEAN | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### notifications
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID | PK |
| user_id | UUID | FK → users |
| type | VARCHAR | ticket_assigned, sla_breach, comment, etc. |
| title | VARCHAR | |
| body | TEXT | |
| read | BOOLEAN | |
| created_at | TIMESTAMP | |

## Índices clave

- tickets (status, priority, assigned_to)
- tickets (created_by)
- tickets (sla_deadline) — para scheduler de escalamiento
- knowledge_articles (tags, category)
- notifications (user_id, read)
