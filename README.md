# Trickreport — Multi-tenant Help Desk

Trickreport is a clean-architecture Help Desk and Ticketing System built with **Go** (Backend), **Astro** (Frontend), **Tailwind CSS** and **PostgreSQL**, featuring multi-tenancy, real-time WebSockets, SLA automations, and an analytics dashboard.

## Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.23+ (chi, pgx, gorilla/websocket) |
| Frontend | Astro 5.x (SSR with Node adapter) + Tailwind CSS |
| UI Style | Neumorphism 2.0 — soft shadows, depth, tactile interactions |
| Database | PostgreSQL 14+ |
| Auth | JWT (cookies + Authorization header), optional LDAP |

## Prerequisites

- Go 1.23+
- Node.js 20+ / npm 10+
- PostgreSQL 14+ (running locally on `localhost:5432`)

## Quick start (local, no Docker)

> **Local development ALWAYS runs without Docker.** Docker is only for production/staging deployments (see [Deployment](#deployment-docker--on-premise)).

### 1. Database

Start PostgreSQL and create the database:

```powershell
# Windows PowerShell
psql -U postgres -c "CREATE DATABASE trickreport;"
```

```bash
# Linux / macOS
createdb trickreport
```

Migrations run **automatically** on backend startup — no manual migration step needed.

### 2. Backend

Copy and customize the environment file:

```powershell
cd backend
copy .env.example .env
# Edit .env: set DATABASE_URL with your local PostgreSQL credentials
```

```bash
# Linux / macOS
cd backend
cp .env.example .env
# Edit .env: set DATABASE_URL with your local PostgreSQL credentials
```

Run the server:

```bash
cd backend
go run cmd/api/main.go
```

Or use the Makefile:

```bash
cd backend
make run
```

The backend runs at `http://localhost:8080`.
Migrations are applied automatically on startup.

Default credentials:

- Email: `admin@trickreport.local`
- Password: `changeme` (set `ADMIN_PASSWORD` in `.env`)

### 3. Frontend

Install dependencies and run the dev server:

```powershell
cd frontend
copy .env.example .env.local
npm install
npm run dev
```

```bash
# Linux / macOS
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

The frontend runs at `http://localhost:4321`.

### 4. One-command dev (both backend + frontend)

From the project root:

```powershell
# Windows PowerShell
.\scripts\dev.ps1
```

```bash
# Linux / macOS
./scripts/dev.sh
```

This starts the backend and frontend concurrently. Press `Ctrl+C` to stop both.

Or with the root Makefile:

```bash
make dev
```

## Project structure

```
backend/
  cmd/api/main.go          # Entry point
  internal/
    application/           # Use cases / services
    domain/                # Domain models and rules
    infrastructure/        # Repositories, auth, realtime
    interfaces/http/       # HTTP handlers, middleware, server, wire
  internal/db/migrations/  # Embedded SQL migrations (auto-applied)

frontend/
  src/
    components/            # NeuButton, NeuCard, NeuInput, NeuBadge, NotificationBell, ...
    layouts/               # App shell with sidebar + header
    pages/                 # Routes (Astro file-based routing)
    lib/api.ts             # API client for the Go backend
  src/styles/global.css    # Neumorphism 2.0 design tokens
```

## Features

1. **Multi-tenancy:** Tenant isolation via `tenant_id` (default tenant seeded).
2. **Authentication:** JWT-based Auth with RBAC (End User, Agent, Admin), refresh tokens, MFA/TOTP, password reset, account lockout.
3. **Tickets:** Full ticket lifecycle, status management, assignment, comments, history, attachments.
4. **Knowledge Base:** Markdown-supported articles with CRUD for agents/admins.
5. **SLAs & Automations:** SLA deadlines tracked by a background worker, visual automation builder.
6. **Analytics:** Dashboard with KPI cards, SVG charts, ticket metrics, CSV/PDF export.
7. **Real-Time:** WebSockets for live ticket notifications, toast popups, notification bell.
8. **Profiles:** Admin profile with system stats, user profile with MFA/sessions management.
9. **UX:** Dark mode, responsive sidebar, loading states, pagination, accessibility (WCAG 2.2), keyboard shortcuts, PWA with offline support.

## Development

### Running tests

```bash
# Backend
cd backend
make test          # go test ./... -v -race
make coverage      # generates coverage.out and prints per-function coverage

# Frontend
cd frontend
npm run build      # type-checks + builds
npx vitest run     # unit tests
```

### Linting & formatting

```bash
# Backend
cd backend
make lint          # golangci-lint
make fmt           # go fmt
make security      # govulncheck

# Frontend
cd frontend
npm run lint       # eslint
npm run format     # prettier
```

### Pre-commit hooks

**Frontend (husky + lint-staged):** after `npm install` in `frontend/`, husky
is wired automatically via the `prepare` script. Staged `.js/.ts/.astro` files
are linted and formatted; `.css/.json/.md` files are formatted.

**Backend (go vet + gofmt):** enable the backend hook with:

```bash
git config core.hooksPath backend/.githooks
```

This runs `go vet ./...` and `gofmt -l .` before each commit.

### API documentation

Interactive API docs are served at the `/swagger` endpoint when the Swagger
handler is enabled (run `make install-tools` then `swag init` in `backend/` to
generate the spec). See the backend `Makefile` `install-tools` target.

## Production build

Backend:

```bash
cd backend
go build -o dist/trickreport-api cmd/api/main.go
```

Frontend:

```bash
cd frontend
npm run build
node ./dist/server/entry.mjs
```

## Deployment (Docker / On-Premise)

> Docker is **only for production/staging deployments**. Local development always runs natively (see [Quick start](#quick-start-local-no-docker)).

### Docker Quick Start

1. Create a `.env` file with the required secrets:

   ```bash
   cp .env.onpremise.example .env
   # edit .env and set strong POSTGRES_PASSWORD, JWT_SECRET, ADMIN_PASSWORD
   ```

2. Build and start all services:

   ```bash
   docker compose up -d --build
   ```

3. Verify the backend is healthy:

   ```bash
   curl http://localhost:8080/health
   ```

The API is available at `http://localhost:8080`, PostgreSQL on `5432`.
All services run on a dedicated `trickreport_net` bridge network with resource
limits and rotated JSON logs.

### On-Premise Deployment

For self-hosted deployments behind a single Caddy reverse proxy on port 80:

1. Copy the on-premise env template and edit it:

   ```bash
   cp .env.onpremise.example .env.onpremise
   ```

2. Start the stack:

   ```bash
   docker compose -f docker-compose.onpremise.yml --env-file .env.onpremise up -d --build
   ```

The app is then served on `http://localhost` (port 80) with the frontend and
backend on the same origin.

### Staging Environment

A staging compose file mirrors on-premise but uses a separate database
(`trickreport_staging`) and a staging domain.

1. Copy the staging env template and edit it:

   ```bash
   cp .env.staging.example .env.staging
   ```

2. Start the staging stack:

   ```bash
   docker compose -f docker-compose.staging.yml --env-file .env.staging up -d --build
   ```

By default staging is exposed on host port `8081` (configurable via
`STAGING_PORT` in `.env.staging`).

### Monitoring

A Prometheus + Grafana + Alertmanager stack is provided as an overlay.

```bash
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml up -d
```

| Service      | URL                       | Notes                                  |
|--------------|---------------------------|----------------------------------------|
| Prometheus   | http://localhost:9090     | Scrapes backend `/metrics`             |
| Grafana      | http://localhost:3001     | admin / `${GRAFANA_PASSWORD:-admin}`   |
| Alertmanager | http://localhost:9093     | Routes service-down / error-rate alerts |

Scrape config lives in `monitoring/prometheus.yml`, alert rules in
`monitoring/rules.yml`, and Alertmanager routing in
`monitoring/alertmanager.yml`.

### Backup & Restore

Scheduled PostgreSQL backups run via a sidecar service that writes gzipped
dumps to `./backups` on a daily cron, with 7-day / 4-week / 6-month retention.

Start the backup service (overlay with the main compose file):

```bash
docker compose -f docker-compose.yml -f docker-compose.backup.yml up -d backup
```

Restore a backup:

```bash
./scripts/restore.sh backups/trickreport_2024-01-01T00:00:00Z.sql.gz
```

> The restore script drops and recreates the `public` schema before loading,
> so it should only be run against a target database you are willing to reset.

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `failed to connect to database` | Wrong `DATABASE_URL` in `.env` or PostgreSQL not running | Verify PostgreSQL is running on `localhost:5432`; check credentials in `backend/.env` |
| `JWT_SECRET must be set` on startup | `ENV=production` without a secret | Set a strong `JWT_SECRET` in `backend/.env` |
| `ADMIN_PASSWORD must be set` on startup | `ENV=production` without admin password | Set `ADMIN_PASSWORD` in `backend/.env` |
| Frontend cannot reach API | `CORS_ORIGINS` missing the frontend origin | Add `http://localhost:4321` to `CORS_ORIGINS` in `backend/.env` |
| `port is already allocated` (5432 / 8080) | Another process holds the port | Stop the conflicting process or change `PORT` in `.env` |
| WebSocket notifications not working | Token not passed or origin blocked | Check `CORS_ORIGINS` includes frontend origin; WebSocket uses `?token=` query param |
| `COOKIE_SECURE` warning in production | Secure cookies disabled in prod | Set `COOKIE_SECURE=true` (forced automatically in production) |

## Contributing

1. Install prerequisites (Go 1.23+, Node 20+, PostgreSQL 14+).
2. Fork and clone the repository.
3. Backend: `cd backend && go run cmd/api/main.go` (migrations auto-applied on startup).
4. Frontend: `cd frontend && npm install && npm run dev`.
5. Or use `make dev` from the root to start both.

### Submitting pull requests

1. Create a feature branch from `main` (`git checkout -b feat/my-feature`).
2. Keep commits focused and use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, `chore:`, …).
3. Ensure `make test` and `make lint` pass locally.
4. Open a PR describing the change, motivation, and any migration/deploy notes.

## License

Released under the **MIT License**. See `LICENSE` for details.
