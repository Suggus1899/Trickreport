# Registro de la Etapa 2 — Entregar Soporte TI con Next.js

Estado: **completada en código y configuración**, con las verificaciones
ejecutables en este entorno en verde. Quedan como gates externos la suite
E2E/de integración (CI) y la aprobación de la alta disponibilidad antes
de publicar el SaaS, tal como exige el plan maestro.

## Punto de partida verificado

El monorepo ya traía de las etapas 0–1 el paquete compartido
`@trickreport/workspace` con los recorridos implementados contra la API
real (`/api/v1`) y las tres apps Next.js como envoltorios. En esta etapa
se revisó ese estado contra el código (no se asumió), se completó la base
común que faltaba y se verificó el conjunto. Las rutas del backend
revisadas en `backend/internal/interfaces/http/server.go` coinciden con
las que consume el cliente tipado (`packages/workspace/src/contracts.ts`
y `api.ts`): auth con CSRF y refresh, MFA, tickets con comentarios,
adjuntos e historial, artículos, administración de usuarios, SLA,
automatizaciones, analíticas y notificaciones.

## Recorridos de Soporte TI (`apps/it`)

| Recorrido | Estado | Notas |
|---|---|---|
| Acceso y registro | Operable | Login, autorregistro y selección de organización vía proxy con cookies HttpOnly; tokens eliminados del cuerpo en el proxy. |
| Recuperación de contraseña | Operable | Solicitud y confirmación con token (`/password-reset`, `/password-reset/confirm`). |
| Sesiones | Operable | Sesión validada en servidor (`/auth/me`), refresh con reintento único, logout, evento de fin de sesión; MFA TOTP con desafío en login y gestión en el perfil. |
| Dashboard | Operable | Resumen real (`/dashboard/summary`) + actividad reciente; los accesos rápidos enlazan a la lista ya filtrada (ver «filtros persistentes»). |
| Tickets: lista | Operable | Búsqueda y filtros por estado/prioridad **persistentes en la URL** (`q`, `status`, `priority`), añadidos en esta etapa. |
| Tickets: detalle | Operable | Comentarios (con notas internas para personal), historial, cambio de estado con control de versión, asignación, adjuntos (subida y descarga). |
| Tickets: crear y borradores | Operable | Creación con clave de idempotencia; borradores locales en IndexedDB aislados por producto/organización/usuario, con lease y sincronización verificada contra la sesión. |
| Artículos / conocimiento | Operable | Lista con búsqueda, lectura, crear/editar (personal), publicar y eliminar (admin). |
| Perfil | Operable | Identidad, rol, organización y MFA (alta con QR, activar/desactivar). |
| Administración: usuarios | Operable (admin) | Alta, edición de rol y activación, búsqueda. |
| Administración: SLA | Operable (admin) | Políticas por prioridad (respuesta, resolución, escalación). |
| Administración: automatizaciones | Operable (admin) | CRUD de reglas con disparadores validados por el servidor. |
| Administración: informes | Operable (admin) | Gráficos reales (`/admin/analytics/charts`): cumplimiento SLA, distribución, resolución media y volumen. |
| Notificaciones | Operable | Campana añadida en esta etapa: contador de no leídas (sondeo 60 s), lista, marcar leída y marcar todas, enlace al ticket relacionado. |

Experiencia común: navegación por rol (la administración solo aparece y
solo se sirve para `admin`; el servidor Go sigue siendo la autoridad de
permisos), estados de carga/vacío/error con reintento, estado de permiso
insuficiente, tablas utilizables en móvil, menú móvil, foco visible y
salto al contenido.

## Base común completada en esta etapa

- **Tokens de marca**: `packages/workspace/src/styles.css` usa ya los
  valores de `docs/brand/tokens.css` (fondo `#f8fafc`, texto `#0f172a`) y
  los tokens de estado (`--danger`, `--warning`, `--success`, `--info`)
  en avisos y píldoras, en tema claro y oscuro.
- **i18n y tema persistentes**: idioma (es inicial, en disponible) y tema
  se guardan en cookie + `localStorage`; el layout de `apps/it` renderiza
  `lang` y la clase `dark` en servidor desde esas cookies, sin parpadeo.
  Los layouts de MSP y atención adoptarán el mismo patrón en sus etapas
  (4 y 5); el cliente compartido ya persiste las preferencias en las tres.
- **Filtros persistentes** en la lista de tickets (URL), que además
  corrige los enlaces del dashboard a `/tickets?status=open`.
- **Notificaciones** en la barra superior contra la API real.
- Utilidad CSS `.mt-6` que el detalle de ticket ya referenciaba.

## Instalación de Soporte TI (dos modalidades preparadas)

- **Autohospedada sencilla**: el `docker-compose.yml` y el `Caddyfile`
  de la raíz ya sirven `apps/it` (servicio `frontend` desde
  `apps/it/Dockerfile`, salida standalone) con la API y PostgreSQL;
  Caddy enruta `/api/v1/ws` directamente a la API y cachea
  `/_next/static` como inmutable. Guía local en `docs/installation.md`.
- **SaaS con alta disponibilidad**: referencia versionada y comentada en
  [`deploy/ha/`](../deploy/ha/README.md) (pila Swarm, Patroni con
  replicación síncrona estricta, HAProxy por rol, Keepalived con VIP),
  enlazada desde `docs/architecture/despliegue.md`. **Pendiente de
  aprobación**: debe validarse en la topología real (incluida la prueba
  de capacidad) antes de publicar el SaaS; no se desplegó infraestructura.

## Verificaciones ejecutadas

- `npm run check` (TypeScript del paquete compartido): **verde**.
- `npm run lint` (ESLint en `apps` y `packages`): **verde**.
- `npm test` (Vitest, proxy): **verde** (2 pruebas).
- `npm run build` de las tres apps (webpack, Next.js 16.4.2): **verde**.
- Prueba de humo de runtime sobre el build de producción de `apps/it`:
  `/login` responde 200 en español por defecto; `/tickets` sin sesión
  redirige (307) a `/login`; con cookies de preferencia el HTML llega con
  `lang="en"` y clase `dark`.
- YAML de `deploy/ha/docker-stack.yml` y `deploy/ha/patroni/patroni.yml`
  validado sintácticamente.
- **No ejecutadas aquí** (sin Go, Docker ni descarga de navegadores en
  este entorno): la suite del backend y los Playwright E2E e integración
  de `apps/it`. Existen y son gate de CI: `.github/workflows/ci.yml` los
  corre en cada push a `main` con PostgreSQL 16 (job «frontend»:
  lint → check → unitarios → E2E → integración → build). La etapa no se
  considera verificada de punta a punta hasta ver ese CI en verde.

## Decisiones y desviaciones

- **Frontend Astro (`frontend/`)**: se deja **intacto** como legado. Ya
  no lo construye CI ni lo sirve ninguna topología versionada, y los
  recorridos de Soporte TI están completos en Next.js; su retirada
  física exige además actualizar las referencias documentales
  (`AGENTS.md` de la raíz y `docs/01`–`09`), por lo que se traslada a la
  Etapa 6 (consolidación) en vez de improvisarla aquí. MSP y atención
  nunca dependieron de él.
- **Tiempo real**: la campana usa sondeo porque el WebSocket del backend
  autentica por cabecera `Authorization`/token y la cookie HttpOnly de
  sesión no es accesible al navegador. Pendiente para una etapa de
  backend: autenticar el WS por cookie y conectar la campana.
- **Catálogo de casos**: los recorridos anteriores cubren las áreas TI
  de fase 2 del catálogo (identidad, incidentes, conocimiento, SLA,
  administración e informes). El volteo de `status` en los JSON
  generados (`docs/use-cases/`) queda pendiente de una pasada del
  generador, para no editar a mano archivos generados.
- **OpenAPI**: el backend ya publica Swagger (`/swagger/*` vía swaggo);
  el cliente del frontend mantiene contratos tipados a mano alineados
  con las rutas reales. La generación automática de tipos desde el
  OpenAPI por producto queda como mejora de la Etapa 6.

## Pasa a la Etapa 3

Completar el ITSM por módulos sobre esta base: catálogo de servicios,
aprobaciones, activos, CMDB, problemas y cambios; ampliar informes y
automatizaciones, con sus casos TI de fase 3, pruebas y documentación.
