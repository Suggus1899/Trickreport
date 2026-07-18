# Backend — API en Go

## Librerías sugeridas

| Propósito | Librería |
|-----------|----------|
| HTTP router | `chi` o `gorilla/mux` |
| ORM / DB | `pgx` + `sqlx` |
| WebSocket | `gorilla/websocket` |
| Auth | `casbin` (RBAC) + `ldap` (para LDAP) |
| Validación | `go-playground/validator` |
| Config | `viper` |
| Logs | `zerolog` |
| Tests | `testify` |

## Endpoints REST (versión 1)

### Auth
- `POST /api/v1/auth/login` — Login con LDAP/SSO
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me` — Perfil actual

### Usuarios (admin)
- `GET /api/v1/users`
- `GET /api/v1/users/:id`
- `POST /api/v1/users`
- `PUT /api/v1/users/:id`
- `DELETE /api/v1/users/:id`

### Tickets
- `GET /api/v1/tickets` — Lista (filtros: status, priority, assigned_to)
- `GET /api/v1/tickets/:id`
- `POST /api/v1/tickets` — Crear
- `PUT /api/v1/tickets/:id` — Actualizar
- `PATCH /api/v1/tickets/:id/status` — Cambiar estado
- `POST /api/v1/tickets/:id/assign` — Asignar agente

### Comentarios
- `GET /api/v1/tickets/:id/comments`
- `POST /api/v1/tickets/:id/comments`

### SLAs
- `GET /api/v1/sla-policies`
- `POST /api/v1/sla-policies`
- `PUT /api/v1/sla-policies/:id`

### Base de conocimientos
- `GET /api/v1/articles` — Lista pública
- `GET /api/v1/articles/:id`
- `POST /api/v1/articles` — Agente/admin
- `PUT /api/v1/articles/:id`
- `DELETE /api/v1/articles/:id`

### Dashboard
- `GET /api/v1/dashboard/summary` — Tickets abiertos, vencidos, por agente
- `GET /api/v1/dashboard/sla-compliance` — % cumplimiento SLA

### Adjuntos
- `POST /api/v1/tickets/:id/attachments` — Subir archivo (multipart)
- `GET /api/v1/tickets/:id/attachments` — Listar adjuntos
- `GET /api/v1/attachments/:id/download` — Descargar
- `DELETE /api/v1/attachments/:id`

### Automatizaciones (admin)
- `GET /api/v1/automation-rules`
- `POST /api/v1/automation-rules`
- `PUT /api/v1/automation-rules/:id`
- `DELETE /api/v1/automation-rules/:id`

### Notificaciones
- `GET /api/v1/notifications`
- `PATCH /api/v1/notifications/:id/read`
- `WebSocket /api/v1/ws` — Tiempo real

## WebSocket

El endpoint `/api/v1/ws` permite:
- Notificaciones en tiempo real al usuario autenticado
- Eventos: ticket_assigned, new_comment, sla_breach, status_change

## Multi-tenant

Cada request incluye el tenant (subdominio o header `X-Tenant-ID`). Todos los repositorios filtran por `tenant_id`. La tabla `tenants` almacena configuración por organización.

## Automation Engine

Un worker procesa reglas de automatización cuando se dispara un evento (creación de ticket, cambio de estado, etc.). Evalúa condiciones y ejecuta acciones (asignar, cambiar prioridad, notificar, etc.).

## Scheduler (escalamiento SLA)

Un goroutine en background corre cada N minutos y:
1. Busca tickets que superan el SLA de respuesta/resolución
2. Registra el breach en ticket_history
3. Escala automáticamente al agente superior (o reasigna)
4. Dispara notificación
