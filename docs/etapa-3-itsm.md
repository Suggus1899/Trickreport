# Registro de la Etapa 3 — Completar el ITSM por módulos

Estado: **completada en código**, con las verificaciones ejecutables en
este entorno en verde. Quedan como gates externos la suite Go con
PostgreSQL 16 y el E2E con Playwright (CI), porque en este entorno no
hay PostgreSQL ni Docker disponibles: ningún test que toque la base de
datos pudo correrse aquí.

## Qué se construyó (backend + frontend contra API real, no maquetas)

### Catálogo de servicios (`/catalog`)
- Categorías y servicios del catálogo con CRUD (personal: admin/agent;
  lectura: cualquier usuario autenticado de la organización).
- Solicitar un servicio desde el catálogo crea una solicitud (ticket)
  vinculada: `POST /catalog/services/{id}/request`. Si el servicio
  exige aprobación, se crea además una aprobación pendiente para ese
  ticket, con el aprobador indicado o sin asignar.
- Migración `000012_catalog_approvals`.

### Aprobaciones (`/approvals`)
- Aprobaciones ligadas a tickets, servicios del catálogo y cambios
  (`target_type` + `target_id`), con aprobar/rechazar con comentario e
  historial inmutable (`approval_history`).
- La regla de permiso vive en Go (`domain.CanDecide`, aplicada en el
  servicio): decide el aprobador nombrado, cualquier admin, o —solo
  cuando no hay aprobador asignado— cualquier agente. Un solicitante
  que no sea aprobador ni personal no puede decidir, aunque la interfaz
  le mostrara el botón. La decisión es atómica (actualización +
  historial en una transacción, solo si sigue pendiente).
- Visibilidad en lista: el personal ve todas las de su organización;
  un solicitante solo ve las que pidió o las que debe aprobar.
- Migración `000012_catalog_approvals`.

### Activos y CMDB (`/assets`, `/ci-types`, `/ticket-assets`)
- CRUD de activos y de tipos de CI, con estado
  (active/maintenance/retired/stock) y atributos JSONB.
- Relaciones dirigidas entre CIs: `depends_on`, `hosts`, `uses`, sin
  auto-relaciones y con ambos extremos verificados en la misma
  organización. Un tipo de CI con activos no se puede borrar
  (RESTRICT → conflicto).
- Vínculo activo↔solicitud (`/ticket-assets`), con aislamiento por
  organización en todas las consultas.
- Migración `000013_assets_cmdb`.

### Problemas (`/problems`)
- Problemas con estado (open/investigating/resolved/closed), causa raíz
  y solución temporal, y vínculo con varios incidentes
  (`/problems/{id}/incidents`). La máquina de estados se valida en el
  dominio; un problema cerrado solo puede reabrirse a investigación.
- Migración `000014_problems_changes`.

### Cambios (`/changes`)
- Cambios con planificación (ventana inicio/fin validada), riesgo,
  problema relacionado, solicitudes vinculadas e implementación y
  cierre con marca temporal.
- Máquina de estados completa
  (draft → planned → awaiting_approval → approved → in_progress →
  completed → closed, con rejected como salida). Dos puertas aplicadas
  en el servicio: pasar a `approved` exige una aprobación aceptada
  para el cambio, y pasar a `awaiting_approval` crea una aprobación
  pendiente si no existe.
- Migración `000014_problems_changes`.

### Informes y automatizaciones (`/admin/analytics/itsm`, `/admin/itsm`)
- Informe ITSM nuevo (solo admin): cambios por estado, problemas por
  estado, activos por tipo, aprobaciones por estado, aprobaciones
  pendientes y servicios activos. Se suma a los informes existentes
  (analytics/charts) y a las automatizaciones de la Etapa 2, sin
  modificarlos.

## Capas y archivos principales

| Capa | Archivos |
|---|---|
| Migraciones | `backend/internal/db/migrations/000012_catalog_approvals.{up,down}.sql`, `000013_assets_cmdb.{up,down}.sql`, `000014_problems_changes.{up,down}.sql` (desde base vacía, sin tocar 000001–000011) |
| Dominio | `backend/internal/domain/itsm/itsm.go` (tipos, validación, máquinas de estado, regla `CanDecide`) |
| Aplicación | `backend/internal/application/itsm/service.go` (permisos, puertas de cambio, solicitud desde catálogo) |
| Infraestructura | `backend/internal/infrastructure/postgres/itsm_repo.go` (todo el SQL, siempre filtrado por `tenant_id`) |
| HTTP | `backend/internal/interfaces/http/handler/itsm.go`, rutas en `server.go`, wiring en `wire.go` |
| Frontend | `packages/workspace/src/itsm.tsx` (catálogo, aprobaciones, activos/CMDB, problemas, cambios, informe), tipos y transiciones en `contracts.ts`, rutas y navegación en `workspace.tsx` (`/catalog`, `/approvals`, `/assets`, `/problems`, `/changes`, `/admin/itsm`) |

## Verificación local (real, en este entorno)

- `go build ./...` y `go vet ./...` en `backend/`: en verde.
- `go test ./internal/domain/itsm/`: 7/7 pruebas en verde (regla de
  decisión, transiciones de problemas y cambios, ventana y riesgo de
  cambio, relaciones CMDB, validación de activos y destinos de
  aprobación).
- `npm run check`, `npm run lint`, `npm run test` y `npm run build`
  (3 apps): en verde. Vitest: 7/7 (2 archivos: proxy preexistente +
  transiciones ITSM nuevas).

## Pendiente de CI (no verificado aquí, no presentado como verificado)

- Suite Go completa con PostgreSQL 16 (`go test ./... -race`): sin
  PostgreSQL local, las migraciones 000012–000014 y el repositorio
  ITSM no se han ejercitado contra una base real. La sintaxis SQL
  sigue el patrón de los repositorios existentes, pero la prueba de
  migración desde base vacía y desde 000011 queda para CI.
- Playwright E2E contra API real: recorridos de catálogo, aprobación,
  activo, problema y cambio.
- Alta disponibilidad: fuera del alcance de esta etapa (Etapa 2 la
  dejó como gate antes de publicar el SaaS).

## Casos del catálogo cubiertos (núcleo, no cada variante)

Según `docs/use-cases/ti.md`: catálogo y solicitudes TI-049–TI-070,
activos TI-105–TI-128, aprobaciones TI-129–TI-140, problemas
TI-141–TI-152, cambios TI-153–TI-168, y el informe ITSM como parte de
TI-185–TI-200. Esta etapa implementa el flujo central de cada uno
(CRUD, vínculo, decisión, máquina de estados, informe). No cubre
todavía variantes avanzadas del catálogo como aprobaciones en varios
niveles, delegación o escalación de aprobaciones vencidas, gestión de
licencias, mantenimiento programado ni planes de retorno formales de
cambio: quedan como trabajo futuro dentro de estos mismos módulos, sin
cambiar las rutas ya publicadas.
