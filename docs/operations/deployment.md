# Despliegue

## Versiones y desarrollo

El runtime activo usa Go 1.26.10, Node 24.22.1 y PostgreSQL 16.15. El frontend TI usa Next.js 16, React 19, Tailwind CSS 4 y HeroUI 3. Para instalarlo en local, sigue [`../installation.md`](../installation.md). El árbol `frontend/` conserva el cliente Astro anterior como referencia y no forma parte del despliegue Compose actual.

## On-premise

El Compose base ejecuta PostgreSQL, la API Go, el frontend Next.js TI y Caddy. Solo Caddy publica puertos. El BFF de Next.js reenvía `/api/*` a la API sin exponer en JSON los tokens de autenticación; `/api/v1/ws` se actualiza directamente en Caddy para permitir WebSocket.

Copia `.env.onpremise.example` a `.env.onpremise`, configura secretos aleatorios y el hostname HTTPS, y valida antes de arrancar:

```sh
docker compose --env-file .env.onpremise config --quiet
docker compose --env-file .env.onpremise up -d --build
```

Caddy necesita DNS público y acceso a los puertos de validación ACME. Para redes cerradas, instala una CA y certificados compatibles antes del despliegue. Una instalación que ya usa `pgdata` debe conservar su nombre de proyecto y volumen; staging utiliza otro volumen.

## Salud y límites conocidos

`/health` comprueba PostgreSQL y `/ready` comprueba si la API puede recibir peticiones. El proxy y Caddy no convierten estas rutas en una prueba de carga ni de alta disponibilidad.

El Compose actual despliega únicamente TI con una instancia de API, una de Next y un PostgreSQL primario. MSP y atención al cliente aún no disponen de APIs, bases de datos, secretos ni despliegues independientes. Los objetivos de 15.000 cuentas, 1.500 conexiones y failover en 60 segundos requieren dimensionamiento, una topología redundante y pruebas antes de ofrecerlos.

La configuración fue revisada estáticamente. No se comprobó el arranque real, TLS ni el failover en Docker durante esta ejecución; valida esos pasos en un host con Docker antes de usarla en producción.
