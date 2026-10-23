# Instalación y desarrollo local

## Requisitos

- Node.js **24.22.1** y npm incluidos con Node.
- Go **1.26.10**.
- PostgreSQL **16.15** en local, o Docker Engine con Docker Compose para el modo on-premise.

Los productos MSP y atención al cliente aún no tienen API y base de datos propias. Esta guía instala el backend TI y su frontend Next.js.

## Primera instalación local

1. Instala las dependencias del workspace desde la raíz:

   ```sh
   npm ci
   ```

2. Crea una base de datos PostgreSQL vacía. En una instalación local:

   ```sh
   createdb trickreport
   ```

3. Copia `backend/.env.example` a `backend/.env` y configura, como mínimo:

   - `DATABASE_URL` con la base local.
   - `JWT_SECRET` con un secreto aleatorio de al menos 32 bytes, por ejemplo `openssl rand -hex 32`.
   - `ADMIN_EMAIL` y `ADMIN_PASSWORD` para crear el administrador inicial.
   - `CORS_ORIGINS=http://localhost:3000`.

4. Inicia el backend desde una terminal:

   ```sh
   cd backend
   go run ./cmd/api
   ```

   El servidor aplica las migraciones incluidas al iniciar. No apliques los archivos SQL manualmente.

5. Inicia el frontend desde otra terminal, en la raíz:

   ```sh
   npm run dev:it
   ```

   Abre `http://localhost:3000`. `API_URL` usa `http://127.0.0.1:8080` por defecto; `APP_PUBLIC_URL` toma el origen de la petición local.

En Windows también puedes ejecutar `scripts/dev.ps1` o `start-local.ps1`; en macOS/Linux usa `scripts/dev.sh`.

## Pruebas y compilación

Desde la raíz, con Node 24.22.1 activo:

```sh
npm run lint
npm run check
npm test
npm run test:e2e
npm run build
```

En `backend/`:

```sh
go test ./... -race
go vet ./...
go build ./...
```

Las pruebas de PostgreSQL que requieren datos de prueba usan una base desechable; nunca apuntes sus variables de prueba a una base de producción.

## Instalación on-premise con Compose

1. Copia `.env.onpremise.example` a `.env.onpremise`.
2. Asigna contraseñas aleatorias a `POSTGRES_PASSWORD`, `JWT_SECRET` y `ADMIN_PASSWORD`. Define `APP_PUBLIC_URL` con el origen HTTPS y `CADDY_SITE` con su hostname.
3. Verifica y arranca:

   ```sh
   docker compose --env-file .env.onpremise config --quiet
   docker compose --env-file .env.onpremise up -d --build
   ```

   El hostname debe resolver al servidor y Caddy necesita acceso a los puertos configurados para obtener TLS. Guarda el archivo de entorno fuera del control de versiones.

Para una instalación existente, conserva el nombre de proyecto y volumen de PostgreSQL. Sigue [`operations/recovery.md`](operations/recovery.md) antes de restaurar o cambiar volúmenes. No elimines volúmenes para resolver un error de arranque.

## Rutas y límites actuales

- `/health` comprueba la API y la base de datos; `/ready` indica si la instancia puede recibir tráfico.
- `/api/*` pasa por el proxy de Next.js, que conserva cookies HttpOnly y elimina tokens de los JSON de autenticación. `/api/v1/ws` conserva la actualización WebSocket hacia Go.
- `apps/it` es el único producto desplegado por el Compose base actual. Las apps MSP y atención al cliente son interfaces iniciales y no deben exponerse como servicios operativos hasta contar con sus propias API, bases de datos, secretos y pruebas de aislamiento.
- La alta disponibilidad y la capacidad de 15.000 cuentas y 1.500 usuarios concurrentes son objetivos. No se han demostrado mediante pruebas de carga y failover.

Lee la [guía en inglés](installation.en.md), la [guía de despliegue](operations/deployment.md) y el [estado del proyecto](../README.md).
