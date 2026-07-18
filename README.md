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
