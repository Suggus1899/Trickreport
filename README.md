<div align="center">

# 🎫 Trickreport

### Multi-tenant Help Desk & Ticketing Platform — Soporte técnico neumórfico

Plataforma de mesa de ayuda y ticketing **multi-tenant** construida con arquitectura limpia (Clean Architecture). Soporte técnico con WebSockets en tiempo real, SLA automatizados, base de conocimiento, analíticas y notificaciones en vivo.

</div>

<br>

<div align="center">

## 🛠️ Tech Stack

</div>

<table align="center">
<tr>
<th colspan="5" align="center" width="600"><sub><b>Frontend</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://astro.build/" target="_blank"><img src="https://cdn.simpleicons.org/astro/FF5D01" width="48" height="48" alt="Astro" /></a>
<br><sub><b><a href="https://astro.build/" target="_blank">Astro 5</a></b></sub>
<br><sub>SSR (Node adapter)</sub>
</td>
<td align="center" width="120">
<a href="https://www.typescriptlang.org/" target="_blank"><img src="https://cdn.simpleicons.org/typescript/3178C6" width="48" height="48" alt="TypeScript" /></a>
<br><sub><b><a href="https://www.typescriptlang.org/" target="_blank">TypeScript 5</a></b></sub>
<br><sub>Type-safe</sub>
</td>
<td align="center" width="120">
<a href="https://tailwindcss.com/" target="_blank"><img src="https://cdn.simpleicons.org/tailwindcss/06B6D4" width="48" height="48" alt="Tailwind CSS" /></a>
<br><sub><b><a href="https://tailwindcss.com/" target="_blank">Tailwind CSS</a></b></sub>
<br><sub>Utility-first</sub>
</td>
<td align="center" width="120">
<a href="https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps" target="_blank"><img src="https://cdn.simpleicons.org/pwa/5A0FC8" width="48" height="48" alt="PWA" /></a>
<br><sub><b><a href="https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps" target="_blank">PWA</a></b></sub>
<br><sub>Offline + installable</sub>
</td>
<td align="center" width="120">
<a href="https://vitest.dev/" target="_blank"><img src="https://cdn.simpleicons.org/vitest/6E9F18" width="48" height="48" alt="Vitest" /></a>
<br><sub><b><a href="https://vitest.dev/" target="_blank">Vitest</a></b></sub>
<br><sub>Unit testing</sub>
</td>
</tr>
<tr>
<th colspan="5" align="center" width="600"><sub><b>Backend</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://go.dev/" target="_blank"><img src="https://cdn.simpleicons.org/go/00ADD8" width="48" height="48" alt="Go" /></a>
<br><sub><b><a href="https://go.dev/" target="_blank">Go 1.25</a></b></sub>
<br><sub>Backend API</sub>
</td>
<td align="center" width="120">
<a href="https://github.com/go-chi/chi" target="_blank"><img src="https://raw.githubusercontent.com/go-chi/docs/master/assets/chi.png" width="48" height="48" alt="go-chi" /></a>
<br><sub><b><a href="https://github.com/go-chi/chi" target="_blank">go-chi v5</a></b></sub>
<br><sub>HTTP router</sub>
</td>
<td align="center" width="120">
<a href="https://www.postgresql.org/" target="_blank"><img src="https://cdn.simpleicons.org/postgresql/4169E1" width="48" height="48" alt="PostgreSQL" /></a>
<br><sub><b><a href="https://www.postgresql.org/" target="_blank">PostgreSQL 14+</a></b></sub>
<br><sub>Primary DB</sub>
</td>
<td align="center" width="120">
<a href="https://jwt.io/" target="_blank"><img src="https://cdn.simpleicons.org/jsonwebtokens/000000" width="48" height="48" alt="JWT" /></a>
<br><sub><b><a href="https://jwt.io/" target="_blank">JWT</a></b></sub>
<br><sub>HS256 · 8h + refresh</sub>
</td>
<td align="center" width="120">
<a href="https://github.com/gorilla/websocket" target="_blank"><img src="https://cdn.simpleicons.org/socketdotio/010101" width="48" height="48" alt="WebSocket" /></a>
<br><sub><b><a href="https://github.com/gorilla/websocket" target="_blank">WebSocket</a></b></sub>
<br><sub>Real-time notifications</sub>
</td>
</tr>
<tr>
<th colspan="5" align="center" width="600"><sub><b>Tooling & Infra</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://nodejs.org/" target="_blank"><img src="https://cdn.simpleicons.org/nodedotjs/339933" width="48" height="48" alt="Node.js" /></a>
<br><sub><b><a href="https://nodejs.org/" target="_blank">npm</a></b></sub>
<br><sub>Package manager</sub>
</td>
<td align="center" width="120">
<a href="https://eslint.org/" target="_blank"><img src="https://cdn.simpleicons.org/eslint/4B32C3" width="48" height="48" alt="ESLint" /></a>
<br><sub><b><a href="https://eslint.org/" target="_blank">ESLint 9</a></b></sub>
<br><sub>Linting</sub>
</td>
<td align="center" width="120">
<a href="https://prettier.io/" target="_blank"><img src="https://cdn.simpleicons.org/prettier/F7B93E" width="48" height="48" alt="Prettier" /></a>
<br><sub><b><a href="https://prettier.io/" target="_blank">Prettier</a></b></sub>
<br><sub>Formatting</sub>
</td>
<td align="center" width="120">
<a href="https://github.com/zerolog/zerolog" target="_blank"><img src="https://img.shields.io/badge/zerolog-000000?style=flat&logo=go&logoColor=white" width="48" height="20" alt="zerolog" /></a>
<br><sub><b><a href="https://github.com/zerolog/zerolog" target="_blank">zerolog</a></b></sub>
<br><sub>Structured logging</sub>
</td>
<td align="center" width="120">
<a href="https://www.docker.com/" target="_blank"><img src="https://cdn.simpleicons.org/docker/2496ED" width="48" height="48" alt="Docker" /></a>
<br><sub><b><a href="https://www.docker.com/" target="_blank">Docker</a></b></sub>
<br><sub>Prod deployment</sub>
</td>
</tr>
</table>

<br>

## 📐 Arquitectura

```
                    ┌─────────────────────────────────────────────────┐
                    │              Caddy / nginx (port 80)            │
                    │           trickreport.local / api               │
                    └──────┬──────────────────┬────────────────────────┘
                           │                  │
            ┌──────────────▼──────────┐  ┌────▼──────────────────────┐
            │  Backend Go (port 8080) │  │  Astro Frontend (SSR)     │
            │  go-chi · pgx · JWT     │  │  Neumorphism 2.0 UI       │
            │  WebSocket · zerolog    │  │  PWA · Dark mode · A11y   │
            └──────────┬──────────────┘  └───────────────────────────┘
                       │
            ┌──────────▼──────────┐
            │  PostgreSQL 14+     │
            │  Multi-tenant       │
            └─────────────────────┘
```

### Clean Architecture (Backend)

```
backend/internal/
├── domain/          ← Entities + business rules (no dependencies)
├── application/     ← Use cases / services (ports defined here)
├── infrastructure/  ← Adapters: PostgreSQL repos, JWT, email, WebSocket
└── interfaces/http/ ← HTTP handlers, middleware, server, wire (DI)
```

## 📱 Aplicaciones

| App | Descripción | Puerto | Rol |
|-----|-------------|--------|-----|
| **Frontend (Astro SSR)** | App única — login, dashboard, tickets, knowledge base, admin, profile, notificaciones en vivo | 4321 | Todos los roles |
| **Backend API (Go)** | REST API + WebSocket — auth, tickets, SLA, automations, analytics, notifications | 8080 | API server |

### Roles

| Rol | Permisos |
|-----|----------|
| **admin** | Todo: usuarios, SLA, automatizaciones, analytics, tickets, artículos |
| **agent** | Tickets, comentarios, artículos, asignación, historial |
| **end_user** | Crear tickets, ver los propios, knowledge base, perfil |

## 📦 Estructura del Proyecto

```
trickreport/
├── backend/                    ← Go API (go-chi + pgx + JWT)
│   ├── cmd/api/                ← Entry point + main
│   ├── internal/
│   │   ├── application/        ← Use cases (auth, ticket, user, article, sla, automation, analytics)
│   │   ├── bootstrap/          ← Initial admin seeding
│   │   ├── config/             ← env loading + feature flags
│   │   ├── db/                 ← pgxpool + embedded migrations (auto-applied)
│   │   ├── domain/             ← Entities + business rules
│   │   │   ├── article/        ← Knowledge base domain
│   │   │   ├── automation/     ← Automation rules domain
│   │   │   ├── event/          ← Domain events + event bus
│   │   │   ├── sla/            ← SLA policy domain
│   │   │   ├── ticket/         ← Ticket, comment, history, notification
│   │   │   └── user/           ← User + role domain
│   │   ├── infrastructure/     ← Adapters
│   │   │   ├── auth/           ← JWT, bcrypt, LDAP, token store
│   │   │   ├── email/          ← SMTP + queue + templates
│   │   │   ├── postgres/       ← Repositories + TxManager + cached tenant resolver
│   │   │   ├── realtime/       ← WebSocket hub adapter
│   │   │   └── storage/        ← Object storage (local + S3)
│   │   ├── interfaces/http/    ← HTTP layer
│   │   │   ├── handler/        ← HTTP handlers (auth, ticket, user, article, sla, automation, analytics, attachment, notification, metrics, swagger)
│   │   │   ├── middleware/     ← Auth, tenant, rate limit, CSRF, compress, metrics, request ID, max body, login rate limit
│   │   │   ├── response/       ← JSON response helpers
│   │   │   ├── validator/      ← Input validation
│   │   │   ├── server.go       ← Chi router + middleware chain
│   │   │   └── wire.go         ← Dependency injection
│   │   ├── realtime/           ← WebSocket hub + client
│   │   └── worker/             ← Background worker (SLA scanner, automation engine)
│   ├── internal/db/migrations/ ← Embedded SQL migrations (000001-000014)
│   ├── .golangci.yml           ← Linter config
│   ├── .env.example            ← Environment template
│   ├── Dockerfile              ← Production container
│   └── Makefile                ← Backend dev targets
│
├── frontend/                   ← Astro 5 SSR frontend
│   ├── src/
│   │   ├── components/         ← NeuButton, NeuCard, NeuInput, NeuBadge, NotificationBell, GlobalSearch, ThemeToggle, Pagination, Skeleton, Spinner, NeuModal, NeuToast, KeyboardHelp, OnlineStatus, charts
│   │   ├── layouts/            ← App shell (sidebar + header + nav)
│   │   ├── lib/                ← API client, export, keyboard shortcuts, config
│   │   ├── pages/              ← Routes (login, dashboard, tickets, articles, admin/*, profile, offline)
│   │   ├── scripts/            ← Service worker registration
│   │   └── styles/global.css   ← Neumorphism 2.0 design tokens
│   ├── public/                 ← PWA manifest, icons, service worker
│   ├── eslint.config.js        ← ESLint 9 flat config
│   ├── .prettierrc.json        ← Prettier config
│   ├── vitest.config.ts        ← Vitest config
│   └── package.json
│
├── scripts/                    ← Dev scripts
│   ├── dev.ps1                 ← Windows: start backend + frontend
│   └── dev.sh                  ← Linux/macOS: start backend + frontend
├── monitoring/                 ← Prometheus + Grafana + Alertmanager configs
├── .github/workflows/          ← CI (Go build/test/vet + frontend build + security scanning)
├── docker-compose*.yml         ← Production / staging / monitoring / backup
├── Makefile                    ← Root dev targets
└── README.md
```

## 🚀 Quick Start (Local, sin Docker)

> **El desarrollo local SIEMPRE se ejecuta sin Docker.** Docker es solo para despliegue productivo/staging.

### Prerrequisitos

| Herramienta | Versión | Instalación |
|-------------|---------|-------------|
| Go | 1.25+ | [go.dev/dl](https://go.dev/dl/) |
| Node.js | 20+ | [nodejs.org](https://nodejs.org/) |
| PostgreSQL | 14+ | [postgresql.org](https://www.postgresql.org/download/) |

### 1. Base de datos

```bash
# Crear base de datos
createdb trickreport

# Las migraciones se ejecutan automáticamente al iniciar el backend.
# No necesitas correr SQL manualmente.
```

### 2. Variables de entorno

Ver sección [🔧 Variables de Entorno](#-variables-de-entorno) abajo.

### 3. Backend (Go)

```bash
cd backend
go run cmd/api/main.go
# ✅ Servidor en http://localhost:8080
# ✅ Migraciones se ejecutan automáticamente
# ✅ Logs en backend/trickreport.log + stdout (dev)
```

### 4. Frontend (Astro)

```bash
cd frontend
npm install
npm run dev
# ✅ Frontend en http://localhost:4321
```

### 5. Un solo comando (backend + frontend)

```powershell
# Windows PowerShell
.\scripts\dev.ps1
```

```bash
# Linux / macOS
./scripts/dev.sh
```

O con el Makefile raíz:

```bash
make dev
```

### Credenciales por defecto

| Campo | Valor |
|-------|-------|
| Email | `admin@trickreport.local` |
| Password | `changeme` (configurable via `ADMIN_PASSWORD` en `.env`) |

## 📊 Features

### Core

| Feature | Descripción |
|---------|-------------|
| **Multi-tenancy** | Aislamiento por `tenant_id` con tenant resolver cacheado |
| **Auth** | JWT (cookie + header), refresh tokens con rotación, MFA/TOTP, password reset, account lockout (5 intentos), password complexity |
| **Tickets** | Ciclo completo: crear, asignar, cambiar estado, comentarios, historial, attachments |
| **Knowledge Base** | Artículos con soporte Markdown, CRUD para agents/admins |
| **SLA** | Policies por prioridad, deadline calculation automática, worker de background escanea mora |
| **Automations** | Reglas visuales (condiciones + acciones), motor de evaluación, test de reglas |
| **Analytics** | Dashboard con KPIs, SVG charts (line + donut), export CSV/PDF |
| **Real-time** | WebSocket: notificaciones de tickets en vivo, toast popups, notification bell con badge |
| **Notifications** | Persistencia en DB + broadcast WebSocket, mark read / mark all read, unread count |

### Frontend UX

| Feature | Descripción |
|---------|-------------|
| **Dark mode** | System preference + toggle manual, persistido en localStorage |
| **Responsive** | Sidebar colapsable con hamburger menu, mobile-first |
| **Loading states** | Skeletons + spinners |
| **Pagination + filtering** | Tickets y artículos con filtros por estado/prioridad/búsqueda |
| **Form validation** | Validación en tiempo real con error states |
| **Accessibility** | WCAG 2.2: ARIA, skip-link, sr-only, focus-visible, keyboard nav |
| **Toasts** | Notificaciones temporales con auto-dismiss |
| **Modals** | Diálogos accesibles |
| **Keyboard shortcuts** | `g d` dashboard, `g t` tickets, `g a` articles, `g u` users, `n` new ticket, `?` help |
| **Global search** | Búsqueda con navegación por teclado |
| **PWA** | Manifest + Service Worker + offline + install prompt + update toast + push notifications |
| **Profiles** | Admin profile (stats + quick actions), user profile (edit, password, MFA, sessions) |

### Backend Reliability

| Feature | Descripción |
|---------|-------------|
| **Transactions** | TxManager con propagación por contexto |
| **Domain events** | Event bus in-memory para desacoplar side effects |
| **Email queue** | Cola con retries + templates HTML |
| **Worker** | Background con retry logic, metrics, health check |
| **Connection pool** | Configurable via env |
| **Rate limiting** | Global por IP + login rate limiting por email+IP |
| **Security** | CSRF, max body size, compression, security headers, file upload validation |
| **Observability** | Prometheus metrics (`/metrics`), request ID propagation, Swagger (`/swagger`) |

## 🔧 Variables de Entorno

### Backend Go

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `DATABASE_URL` | ✅ | — | Connection string PostgreSQL |
| `JWT_SECRET` | ✅ | — | Secret para firmar JWT |
| `JWT_EXP_HOURS` | ❌ | `8` | Expiración del token en horas |
| `PORT` | ❌ | `8080` | Puerto del servidor |
| `ENV` | ❌ | `development` | `development` o `production` |
| `CORS_ORIGINS` | ✅ prod | `http://localhost:4321` | Origins permitidos (separados por coma) |
| `COOKIE_SECURE` | ❌ | `false` | Forzado a `true` en producción |
| `ADMIN_EMAIL` | ❌ | `admin@trickreport.local` | Email del admin inicial |
| `ADMIN_PASSWORD` | ❌ | `changeme` | Password del admin inicial |
| `LDAP_ENABLED` | ❌ | `false` | Habilitar autenticación LDAP |

### Frontend (Astro)

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `API_URL` | ✅ | — | URL base del backend Go |

### Ejemplo `.env`

```env
# Backend Go
DATABASE_URL=postgres://postgres:postgres@localhost:5432/trickreport?sslmode=disable
JWT_SECRET=tu-secreto-super-seguro-cambiar-en-produccion
JWT_EXP_HOURS=8
PORT=8080
ENV=development
CORS_ORIGINS=http://localhost:4321
ADMIN_EMAIL=admin@trickreport.local
ADMIN_PASSWORD=changeme
```

```env
# Frontend
API_URL=http://localhost:8080
```

## 🧪 Testing & CI

```bash
# Backend
cd backend
make test          # go test ./... -v -race
make coverage      # coverage.out + per-func report
make lint          # golangci-lint
make security      # govulncheck

# Frontend
cd frontend
npm run build      # type-check + build
npx vitest run     # unit tests
npm run lint       # eslint
```

### CI Pipeline (GitHub Actions)

| Job | Descripción |
|-----|-------------|
| **Backend (Go)** | `go build` → `go vet` → `go test` → `govulncheck` → coverage report |
| **Frontend (Astro)** | `npm ci` → `npm run build` → `npm run lint` → `npm audit` |
| **Security** | Trivy container scan + dependency audit |

## 🚢 Despliegue (Docker / On-Premise)

> Docker es **solo para despliegue productivo/staging**. El desarrollo local siempre es nativo.

### Docker Quick Start

```bash
cp .env.onpremise.example .env
# Editar .env con secrets fuertes
docker compose up -d --build
curl http://localhost:8080/health
```

### On-Premise (Caddy, puerto 80)

```bash
cp .env.onpremise.example .env.onpremise
docker compose -f docker-compose.onpremise.yml --env-file .env.onpremise up -d --build
```

### Staging

```bash
cp .env.staging.example .env.staging
docker compose -f docker-compose.staging.yml --env-file .env.staging up -d --build
```

### Monitoring (Prometheus + Grafana + Alertmanager)

```bash
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml up -d
```

| Service | URL | Notes |
|---------|-----|-------|
| Prometheus | http://localhost:9090 | Scrapes `/metrics` |
| Grafana | http://localhost:3001 | admin / `${GRAFANA_PASSWORD:-admin}` |
| Alertmanager | http://localhost:9093 | Alert routing |

### Backup

```bash
docker compose -f docker-compose.yml -f docker-compose.backup.yml up -d backup
./scripts/restore.sh backups/trickreport_2024-01-01.sql.gz
```

## 🐛 Troubleshooting

| Síntoma | Causa | Solución |
|---------|-------|----------|
| `failed to connect to database` | `DATABASE_URL` incorrecta o PostgreSQL no corre | Verificar PostgreSQL en `localhost:5432`; revisar credenciales en `backend/.env` |
| `JWT_SECRET must be set` | `ENV=production` sin secret | Setear `JWT_SECRET` fuerte en `backend/.env` |
| Frontend no llega al API | `CORS_ORIGINS` no incluye el origen | Agregar `http://localhost:4321` a `CORS_ORIGINS` |
| WebSocket no funciona | Origin bloqueado o token no pasado | Verificar `CORS_ORIGINS`; WS usa `?token=` query param |
| Puerto en uso (5432/8080/4321) | Otro proceso ocupa el puerto | Detener el proceso o cambiar `PORT` en `.env` |

## 🤝 Contributing

1. Instalar prerrequisitos (Go 1.25+, Node 20+, PostgreSQL 14+).
2. Fork + clone del repo.
3. `make dev` desde la raíz (arranca backend + frontend).
4. Crear branch desde `dev` (`git checkout -b feat/my-feature`).
5. Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`).
6. `make test` + `make lint` deben pasar.
7. Abrir PR a `dev`.

### Pre-commit hooks

**Frontend (husky + lint-staged):** se instala automáticamente con `npm install`.

**Backend (go vet + gofmt):**

```bash
git config core.hooksPath backend/.githooks
```

## 📜 Licencia

**© 2026 Gustavo Colina (@Suggus1899). Todos los derechos reservados.**

Este software y su código fuente son **propiedad exclusiva** de Gustavo Colina (@Suggus1899). 

- **No** está permitido copiar, modificar, distribuir, sublicenciar ni usar este código, total o parcialmente, sin autorización expresa y por escrito del autor.
- **No** está permitido usar este código con fines comerciales ni privados sin una licencia válida.
- Cualquier uso no autorizado constituye una violación de los derechos de autor y será perseguido conforme a la ley.

**Este es un software propietario. No es código abierto (open source) ni software libre.**

---

<div align="center">

<sub>Hecho con ❤️ para mesa de ayuda y soporte técnico</sub>

</div>
