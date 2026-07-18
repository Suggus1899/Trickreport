# Trickreport — Multi-tenant Help Desk

Trickreport is a clean-architecture Help Desk and Ticketing System built with **Go** (Backend), **Astro** (Frontend), **Tailwind CSS** and **PostgreSQL**, featuring multi-tenancy, real-time WebSockets, SLA automations, and an analytics dashboard.

## Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.23+ (chi, pgx, gorilla/websocket, a-h/templ) |
| Frontend | Astro 5.x (SSR with Node adapter) + Tailwind CSS |
| UI Style | Neumorphism 2.0 — soft shadows, depth, tactile interactions |
| Database | PostgreSQL 14+ |
| Auth | JWT (cookies + Authorization header), optional LDAP |

## Prerequisites

- Go 1.23+
- Node.js 20+ / npm 10+
- PostgreSQL 14+
- `templ` CLI (`go install github.com/a-h/templ/cmd/templ@v0.3.1020`)

## Local setup without Docker

### 1. Database

Start PostgreSQL and create the database:

```powershell
# Windows PowerShell
psql -U postgres -c "CREATE DATABASE trickreport;"
```

Apply migrations:

```powershell
cd backend
$env:PGPASSWORD='YOUR_POSTGRES_PASSWORD'
psql -U postgres -d trickreport -f migrations/001_init.sql
psql -U postgres -d trickreport -f migrations/002_tickets.sql
psql -U postgres -d trickreport -f migrations/003_knowledge_base.sql
psql -U postgres -d trickreport -f migrations/004_sla_policies.sql
psql -U postgres -d trickreport -f migrations/005_automations.sql
psql -U postgres -d trickreport -f migrations/006_sla_breached.sql
```

### 2. Backend

Copy and customize the environment file:

```powershell
cd backend
copy .env.example .env
```

Run the server:

```powershell
templ generate ./internal/interfaces/http/views/
go run cmd/api/main.go
```

The backend runs at `http://localhost:8080`.

Default credentials:

- Email: `admin@trickreport.local`
- Password: `changeme`

### 3. Frontend

Install dependencies and run the dev server:

```powershell
cd frontend
copy .env.example .env.local
npm install
npm run dev
```

The frontend runs at `http://localhost:4321`.

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

## Project structure

```
backend/
  cmd/api/main.go          # Entry point
  internal/
    application/           # Use cases / services
    domain/                # Domain models and rules
    infrastructure/        # Repositories, auth, realtime
    interfaces/http/       # HTTP handlers, middleware, views (Templ)
  migrations/              # Plain SQL migrations

frontend/
  src/
    components/            # NeuButton, NeuCard, NeuInput, NeuBadge
    layouts/               # App shell
    pages/                 # Routes (Astro file-based routing)
    lib/api.ts             # API client for the Go backend
  src/styles/global.css    # Neumorphism 2.0 design tokens
```

## Features

1. **Multi-tenancy:** Tenant isolation via `tenant_id` (default tenant seeded).
2. **Authentication:** JWT-based Auth with RBAC (End User, Agent, Admin).
3. **Tickets:** Full ticket lifecycle, status management, assignment, comments and history.
4. **Knowledge Base:** Markdown-supported articles with CRUD for agents/admins.
5. **SLAs & Automations:** SLA deadlines tracked by a background worker.
6. **Analytics:** Dashboard with KPI cards and ticket metrics.
7. **Real-Time:** WebSockets integration for live UI updates.

## Docker Quick Start

The fastest way to run the whole stack is with Docker Compose.

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

## On-Premise Deployment

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

## Staging Environment

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

## Monitoring

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

## Backup & Restore

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
| `failed to connect to database` / `pg_isready` fails | Wrong `POSTGRES_PASSWORD` or DB not ready | Verify `.env` values match; wait for the `postgres` healthcheck to pass (`docker compose ps`). |
| `port is already allocated` (5432 / 8080 / 80) | Another process holds the port | Stop the conflicting process or remap the host port in the compose file. |
| `JWT_SECRET must be set` on startup | `ENV=production` without a secret | Set a strong `JWT_SECRET` in `.env`; it must not contain the word "change". |
| `ADMIN_PASSWORD must be set` on startup | `ENV=production` without admin password | Set `ADMIN_PASSWORD` in `.env`. |
| Migration errors on startup | Out-of-order or partially applied migrations | Inspect `schema_migrations` table; re-run with a clean volume (`docker compose down -v`) for dev only. |
| Frontend cannot reach API | `CORS_ORIGINS` missing the frontend origin | Add the frontend origin to `CORS_ORIGINS` (comma-separated). |
| `COOKIE_SECURE` warning in production | Secure cookies disabled in prod | Set `COOKIE_SECURE=true` (forced automatically in production). |

## Contributing

### Development environment

1. Install prerequisites (Go 1.23+, Node 20+, PostgreSQL 14+, `templ` CLI).
2. Fork and clone the repository.
3. Backend: `cd backend && go run cmd/api/main.go` (apply migrations first — see [Local setup without Docker](#local-setup-without-docker)).
4. Frontend: `cd frontend && npm install && npm run dev`.

### Pre-commit hooks

**Frontend (husky + lint-staged):** after `npm install` in `frontend/`, husky
is wired automatically via the `prepare` script. Staged `.js/.ts/.astro` files
are linted and formatted; `.css/.json/.md` files are formatted.

**Backend (go vet + gofmt):** enable the backend hook with:

```bash
git config core.hooksPath backend/.githooks
```

This runs `go vet ./...` and `gofmt -l .` before each commit.

### Running tests

```bash
# Backend
cd backend
make test          # go test ./... -v -race
make coverage      # generates coverage.out and prints per-func coverage

# Frontend
cd frontend
npm run build      # type-checks + builds
```

### Submitting pull requests

1. Create a feature branch from `main` (`git checkout -b feat/my-feature`).
2. Keep commits focused and use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, `chore:`, …).
3. Ensure `make test` and `make lint` pass locally.
4. Open a PR describing the change, motivation, and any migration/deploy notes.

## API documentation

Interactive API docs are served at the `/swagger` endpoint when the Swagger
handler is enabled (run `make install-tools` then `swag init` in `backend/` to
generate the spec). See the backend `Makefile` `install-tools` target.

## License

Released under the **MIT License**. See `LICENSE` for details.

