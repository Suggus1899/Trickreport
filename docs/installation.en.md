# Installation and local development

## Requirements

- Node.js **24.22.1** and its bundled npm.
- Go **1.26.10**.
- PostgreSQL **16.15**, or Docker Engine with Docker Compose for an on-premise deployment.

The MSP and customer-service products do not have their own APIs or databases yet. This guide installs the IT backend and its Next.js frontend.

## Fresh local setup

1. Install the workspace dependencies from the repository root:

   ```sh
   npm ci
   ```

2. Create an empty PostgreSQL database. For a local installation:

   ```sh
   createdb trickreport
   ```

3. Copy `backend/.env.example` to `backend/.env` and configure at least:

   - `DATABASE_URL` for the local database.
   - `JWT_SECRET` with a random secret of at least 32 bytes, for example `openssl rand -hex 32`.
   - `ADMIN_EMAIL` and `ADMIN_PASSWORD` for the initial administrator.
   - `CORS_ORIGINS=http://localhost:3000`.

4. Start the backend in one terminal:

   ```sh
   cd backend
   go run ./cmd/api
   ```

   The server applies its embedded migrations at startup. Do not apply the SQL files manually.

5. Start the frontend in a second terminal from the repository root:

   ```sh
   npm run dev:it
   ```

   Open `http://localhost:3000`. `API_URL` defaults to `http://127.0.0.1:8080`; `APP_PUBLIC_URL` uses the local request origin by default.

On Windows, you can also run `scripts/dev.ps1` or `start-local.ps1`; on macOS/Linux, use `scripts/dev.sh`.

## Tests and builds

From the repository root with Node 24.22.1 active:

```sh
npm run lint
npm run check
npm test
npm run test:e2e
npm run build
```

From `backend/`:

```sh
go test ./... -race
go vet ./...
go build ./...
```

PostgreSQL tests that need data use a disposable test database. Never point test database variables at production.

## On-premise Compose installation

1. Copy `.env.onpremise.example` to `.env.onpremise`.
2. Set random values for `POSTGRES_PASSWORD`, `JWT_SECRET`, and `ADMIN_PASSWORD`. Set `APP_PUBLIC_URL` to the HTTPS origin and `CADDY_SITE` to its hostname.
3. Validate and start the stack:

   ```sh
   docker compose --env-file .env.onpremise config --quiet
   docker compose --env-file .env.onpremise up -d --build
   ```

   The hostname must resolve to the server, and Caddy needs access to the configured ports to obtain TLS. Keep the environment file outside version control.

For an existing installation, preserve its Compose project name and PostgreSQL volume. Follow [`operations/recovery.md`](operations/recovery.md) before restoring or changing volumes. Do not remove volumes to resolve a startup error.

## Current routes and limits

- `/health` checks the API and database; `/ready` reports whether the instance can receive traffic.
- `/api/*` goes through the Next.js proxy, which keeps HttpOnly cookies and removes tokens from authentication JSON. `/api/v1/ws` retains the WebSocket upgrade to Go.
- `apps/it` is the only product deployed by the current base Compose file. The MSP and customer-service apps are initial interfaces; do not expose them as operational services until they have separate APIs, databases, secrets, and isolation tests.
- High availability and capacity for 15,000 accounts and 1,500 concurrent users are goals. Load and failover tests have not demonstrated them.

See the [Spanish guide](installation.md), [deployment guide](operations/deployment.md), and [project status](../README.md).
