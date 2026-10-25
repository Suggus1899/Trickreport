# Registro de la Etapa 6 — Consolidación

Estado: **consolidación documentada y verificada en lo ejecutable en
este entorno; dos puntos quedan bloqueados y se declaran como tales
abajo, no como hechos.** Las Etapas 2–5 están publicadas en `main`
(ver SHAs). Esta etapa no cambia código de producto: su entrega es
la revisión de cobertura, la auditoría de despliegue y este registro.

## Cobertura del catálogo: 600 casos verificados en local

Revisión ejecutada contra `docs/use-cases/{ti,msp,cx}.json`:

- **TI: 200 casos, IDs TI-001…TI-200 únicos**, reparto por área
  idéntico al plan maestro (incidents 32, problems-changes 28,
  assets 24, service-catalog 22, sla 20, identity 16, automation 16,
  reporting 16, knowledge 14, approvals 12). Fases: 114 en fase 2
  (Soporte TI) y 86 en fase 3 (ITSM).
- **MSP: 200 casos, IDs MSP-001…MSP-200 únicos**, reparto del plan
  (service-desk 30, monitoring-remote 28, contracts 22,
  time-expenses 20, inventory-discovery 20, clients 18,
  technicians 16, profitability 16, reporting 16, projects 14).
  Fase: 200 en fase 4.
- **CX: 200 casos, IDs CX-001…CX-200 únicos**, reparto del plan
  (inbox 32, portal-chat 26, whatsapp 24, email 22,
  customer-identity 20, routing 20, sla 16, knowledge 14,
  admin-reporting 14, surveys 12). Fase: 200 en fase 5.
- Cada caso conserva los 15 campos obligatorios del plan
  (identificador, objetivo, actor, permisos, precondiciones, flujo,
  alternativas, ejemplo, pantallas, entidades, contratos,
  dependencias, aceptación, prioridad, etapa y estado).

**Sobre el campo `status`:** los 600 casos siguen en `planned`.
No se hizo un cambio masivo a `implemented`: la propia nota del
catálogo (`status_note`) advierte que `planned` no significa
capacidad implementada ni validada y que hay que enlazar evidencia
antes de avanzar de estado. La evidencia por área es la tabla de
etapas de abajo; el cambio de estado caso por caso queda pendiente
de ese enlazado, no se inventó.

## Etapas publicadas en `main` (SHAs verificados con `list_commits`)

- **Etapa 2 (Soporte TI):** `2035dcc`, `640610f`, `df77a17`,
  `66deafa`, `5b8f581`, `919367e`.
- **Etapa 3 (ITSM):** `905d7e5`, `6d15b08`, `1bd9aff`, `657e2ac`.
- **Etapa 4 (MSP):** `6096a3d`, `20f6d34`, `deae22a`, `f88e237`.
- **Etapa 5 (Atención):** `d40e823`, `90f5e53`, `69efe56`,
  `a079d93`.
- **Migraciones:** `000001`–`000016`, cada una con up y down
  presentes en `backend/internal/db/migrations/`.

## Auditoría de despliegue: `frontend/` Astro ya no alimenta nada activo

- El servicio `frontend` de `docker-compose.yml` construye
  `apps/it/Dockerfile` (Next.js), no `frontend/`.
- El job `frontend` de `.github/workflows/ci.yml` y el de
  `release.yml` también apuntan a `apps/it/` (Next.js).
- La referencia HA SaaS está versionada en `deploy/ha/`
  (Swarm, Patroni, HAProxy, Keepalived), como referencia de
  despliegue, no como infraestructura probada aquí.

## Checks ejecutables en este entorno (resultado real)

- `npm run check` (tsc): verde.
- `npm run lint`: verde (solo el aviso preexistente del
  directorio `pages`).
- Vitest completo del workspace: **21/21 en 4 archivos**
  (incluye itsm 5 y cx 7).
- Build de producción de las 3 apps Next.js (`it`, `msp`,
  `customer-service`): verde en la Etapa 5; la Etapa 6 no cambia
  código, así que esa evidencia sigue vigente para este árbol.
- Go, PostgreSQL y Docker **no están disponibles en este
  entorno**: `go build`/`go vet`, la suite Go con PostgreSQL 16
  y cualquier prueba de carga no se ejecutaron aquí.

## NO ejecutado / bloqueado (declarado, no oculto)

### 1. Eliminación física de `frontend/` Astro — BLOQUEADA por vía de publicación

`frontend/` sigue presente en `main` (106 archivos rastreados;
el resto del tamaño local es `node_modules` no rastreado).
No se eliminó por una limitación de herramientas, no por
decisión de producto:

- `push_files` (la vía usada en todas las etapas) **no puede
  borrar archivos**: solo crea o actualiza.
- `delete_file` del conector es **un archivo por llamada, cada
  una con aprobación fresca**: borrar `frontend/` exigiría 106
  aprobaciones, una por archivo.
- No hay herramienta de API de árboles Git en el conector para
  borrar un directorio en un solo commit.
- El clon local tiene remoto HTTPS sin credenciales
  configuradas y no existe una vía de `git push` autorizada
  (los flujos con PAT o device auth están prohibidos por la
  guía del conector).

Consecuencia: `frontend/` es código muerto documentado (fuera de
CI y de Compose, ver auditoría arriba), pero sigue en el repo.
El paso que falta es una decisión del operador: autorizar una
vía de borrado con Git (un solo commit que elimine el árbol)
o las 106 aprobaciones de `delete_file`. Hasta entonces, las
líneas de `README.md` que dicen que Astro "se mantiene como
referencia de migración" siguen siendo la descripción correcta
del estado real, y no se modificaron para no afirmar una
eliminación que no ocurrió.

### 2. Validación de capacidad — PENDIENTE de infraestructura real

Los objetivos del plan (15.000 cuentas, 1.500 concurrentes,
≥150 ops/s durante 45 min, p95 de lecturas <500 ms) **no se
probaron aquí**: no hay Go, PostgreSQL, Docker ni máquinas de
carga en este entorno. Quedan para una prueba en
infraestructura real, con máquinas, datos y carga publicados,
como exige el propio plan. No se presenta ningún número de
capacidad como medido.

### 3. Pendientes heredados de CI / operador

- Suite Go completa con PostgreSQL 16 y Playwright E2E con API
  real (Etapas 3–5): pendiente de CI.
- Activación viva con credenciales reales: Zabbix, MeshCentral,
  correo SMTP/IMAP y WhatsApp Business Platform (Etapas 4–5):
  pendiente del operador; no se inventó ninguna conexión viva.
- Protección de rama y ensayo de backup/restore (Etapa 0):
  pendientes (sin herramienta de protección de rama en el
  conector; sin infraestructura para el ensayo).
- WhatsApp: el webhook aún no resuelve la cuenta desde
  `metadata.phone_number_id` y el envío saliente por la API de
  Meta no está implementado (limitaciones ya documentadas en
  la Etapa 5).
