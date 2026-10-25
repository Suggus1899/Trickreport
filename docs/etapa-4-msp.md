# Registro de la Etapa 4 — Entregar MSP

Estado: **completada en código**, con las verificaciones ejecutables en
este entorno en verde. Quedan como gates externos la suite Go con
PostgreSQL 16 y el E2E con Playwright (CI), porque en este entorno no
hay PostgreSQL ni Docker disponibles: ningún test que toque la base de
datos pudo correrse aquí. La validación viva de Zabbix y MeshCentral
queda pendiente de instancias y credenciales reales (ver
`docs/msp-integraciones.md`): los adaptadores son configurables y no
se inventó ninguna conexión viva.

## Condición de salida

Operación multicliente comprobada en código y sin filtraciones entre
clientes: cada consulta filtra por `tenant_id` y por `client_id`, y
un contacto de un cliente ve cualquier otro cliente como "no
encontrado" (nunca como dato parcial). La regla vive en Go
(`domain.CanAccessClient`, aplicada en el servicio); la interfaz solo
filtra lo que muestra.

## Qué se construyó (backend + frontend contra API real, no maquetas)

### Clientes y contactos (`/msp/clients`)

- CRUD de clientes del proveedor (nombre, dominio único por
  organización, ciudad, estado, contacto principal) y contactos por
  cliente, con rol `contact` o `client_admin`. Un contacto puede
  enlazarse a un usuario de la organización (`user_id`); sin enlace
  es solo un registro del cliente.
- Cambio explícito de cliente en la UI: el selector persiste en
  `localStorage` (`trickreport_msp_client`) y todas las pantallas MSP
  filtran por él. La mesa multicliente muestra el contrato vigente
  junto al cliente activo.
- Migración `000015_msp`.

### Contratos y SLA (`/msp/contracts`)

- Contratos por cliente con vigencia, cobertura, estado, cuota
  mensual, tarifa horaria, horas incluidas y SLA base (respuesta y
  resolución) más anulaciones por prioridad (`sla_overrides`).
- El ticket multicliente resuelve cliente + contrato + SLA en el
  servidor: contrato explícito (verificado del mismo cliente) o el
  contrato activo dentro de su ventana; sin contrato activo, el
  ticket se crea sin SLA (minutos a 0 y `contract_id` nulo), nunca
  con el SLA de otro cliente.

### Mesa multicliente (`/msp/tickets`)

- Tickets con cliente explícito y contexto persistido en
  `msp_ticket_context` (ticket, cliente, contrato, SLA resuelto).
- Filtro por cliente en lista; un contacto solo ve los tickets de
  sus clientes, y un filtro por otro cliente devuelve lista vacía,
  no error revelador ni datos ajenos.
- `GET /msp/tickets/{id}/context` expone el contexto resuelto.

### Técnicos y planificación (`/msp/technicians`, `/msp/interventions`)

- Técnicos vinculados a usuarios, con habilidades, tarifa horaria y
  estado; disponibilidad semanal por día y franja.
- Intervenciones por cliente y técnico con ventana validada; el
  servidor rechaza solapes del mismo técnico (intervenciones no
  canceladas) y tickets de otro cliente. Máquina de estados
  `scheduled → in_progress → done`, con `cancelled` como salida.

### Horas, gastos y costes (`/msp/time-entries`, `/msp/expenses`)

- Registro de tiempo por ticket o intervención, con minutos (1–1440),
  fecha, marca facturable y coste calculado en el servidor desde la
  tarifa del técnico, con redondeo a medio céntimo hacia arriba.
- Gastos por cliente, ticket o proyecto, con importe positivo.
- El servidor rechaza tiempo y gastos cuyo ticket o proyecto
  pertenezca a otro cliente.

### Proyectos (`/msp/projects`)

- Proyectos por cliente con presupuesto y ventana, y tareas con
  estado (`todo → doing → done`), técnico asignado y minutos
  estimados. El avance es el porcentaje de tareas terminadas (0 %
  sin tareas).

### Rentabilidad y facturación (`/msp/billing-drafts`, `/msp/profitability`)

- Resumen por cliente: horas, ingresos por tiempo (tarifa de
  contrato cuando la entrada la lleva, si no tarifa de técnico),
  coste de tiempo, gastos, borradores emitidos y beneficio.
- `POST /msp/billing-drafts` genera un borrador para un periodo:
  horas facturables no facturadas valoradas a tarifa de contrato,
  gastos no facturados del periodo y cuota mensual del contrato
  activo. **El borrador nunca cobra nada**; solo registra lo que se
  facturaría. Las horas y gastos incluidos quedan marcados como
  facturados en la misma transacción que crea el borrador, y hay un
  único borrador por cliente y periodo, así que un segundo borrador
  no puede facturar dos veces lo mismo. Estados: `draft → issued` o
  `void` (un borrador anulado no se reemite).

### Inventario y descubrimiento (`/msp/clients/{id}/assets`, `.../discovered-assets`)

- Reutiliza los activos de la Etapa 3 sin tocarlos:
  `msp_client_assets` asigna cada activo a un cliente.
- Descubrimiento manual/importación documentada: se registra el
  hallazgo (ID externo, nombre, tipo, payload) y se importa
  eligiendo el tipo de CI; la importación crea el activo, lo vincula
  al cliente y marca el hallazgo como importado en una transacción.
  No hay agente de descubrimiento propio en esta etapa.

### Monitorización Zabbix (`/msp/monitoring`, webhook)

- Fuentes configurables (`/msp/monitoring/sources`) con token de
  webhook: solo se guarda el hash SHA-256, nunca el token.
- `POST /api/v1/msp/monitoring/webhook` (sin sesión; tenant por
  `X-Tenant-ID`, token por `X-MSP-Webhook-Token`, exento de CSRF
  solo para este path): ingesta autenticada e idempotente. Una
  entrega repetida para un recurso con alerta abierta actualiza esa
  alerta; solo la primera crea ticket (con cliente, contrato y SLA
  resueltos, atribuido al primer admin/agente de la organización,
  porque un webhook no trae usuario). El índice parcial de alertas
  abiertas impide duplicados incluso ante entregas concurrentes.
- Activación y comprobación paso a paso en
  `docs/msp-integraciones.md`. Validación viva pendiente.

### Acceso remoto MeshCentral (`/msp/remote-sessions`)

- Sesiones vinculadas a cliente, ticket y técnico, con proveedor
  `meshcentral` e ID de sesión externo. Flujo:
  `requested → authorized → active → ended`, con `denied` como
  salida. Solo el personal autoriza (un técnico nunca se
  autoautoriza) y solo el técnico propietario o el personal inicia
  una sesión autorizada. Trickreport registra y audita; MeshCentral
  ejecuta la sesión. Activación en `docs/msp-integraciones.md`.
  Validación viva pendiente.

### Informes (`/admin/msp`, `/admin/analytics` existente)

- Informe MSP nuevo (solo admin): clientes totales y activos,
  contratos activos, alertas abiertas, sesiones activas o
  autorizadas y tickets MSP abiertos.

## Capas y archivos principales

| Capa | Archivos |
|---|---|
| Migración | `backend/internal/db/migrations/000015_msp.{up,down}.sql` (desde base vacía, sin tocar 000001–000014) |
| Dominio | `backend/internal/domain/msp/msp.go` (tipos, validación, SLA por contrato, coste de tiempo, avance de proyecto, no-duplicación de alertas, reglas de sesión remota, aislamiento por cliente) + `msp_test.go` |
| Aplicación | `backend/internal/application/msp/service.go` (aislamiento por cliente, resolución de contrato/SLA, borradores de facturación, ingesta de alertas, autorización de sesiones) |
| Infraestructura | `backend/internal/infrastructure/postgres/msp_repo.go` (todo el SQL, siempre filtrado por `tenant_id` y `client_id`) |
| HTTP | `backend/internal/interfaces/http/handler/msp.go`, rutas en `server.go`, wiring en `wire.go`, exención CSRF del webhook en `middleware/csrf.go` |
| Frontend | `packages/workspace/src/msp.tsx` (clientes, contratos, mesa, técnicos, horas/gastos, proyectos, facturación, inventario, monitorización, acceso remoto, informe), tipos y transiciones en `contracts.ts`, rutas y navegación en `workspace.tsx` (`/clients`, `/msp-tickets`, `/contracts`, `/technicians`, `/time`, `/projects`, `/billing`, `/inventory`, `/monitoring`, `/remote`, `/admin/msp`) |
| Docs | `docs/etapa-4-msp.md`, `docs/msp-integraciones.md` |

## Verificación local (resultado real)

- `go build ./...` en `backend/`: OK.
- `go vet ./...` en `backend/`: OK.
- `go test ./internal/domain/msp/`: 8 tests en verde (aislamiento,
  ventana de contrato y SLA, coste de tiempo, avance de proyecto,
  no-duplicación de alertas, solapes, reglas de sesión remota,
  total de borrador y beneficio).
- `npm run check` (tsc): OK. `npm run lint` (eslint): OK.
- `npm run test` (Vitest): 14/14 en verde (7 nuevos de MSP:
  transiciones, formato de céntimos, avance, contrato vigente).
- `npm run build` (3 apps: it, msp, customer-service): OK (exit 0;
  el churn de `apps/it/next-env.d.ts` generado por el build se
  revirtió en local al contenido remoto y no se sube).
- No hay PostgreSQL 16 ni Docker en este entorno: los tests Go con
  base de datos, la suite completa con `-race` y el E2E con
  Playwright quedan para CI.

## Pendiente

- CI: suite Go con PostgreSQL 16 (incluidas las pruebas de
  aislamiento entre clientes contra base de datos), E2E Playwright
  con API real, y migraciones desde base vacía y desde la versión
  anterior (000014 → 000015).
- Activación viva: Zabbix (webhook de prueba, doble entrega,
  resolución) y MeshCentral (sesión de prueba) con instancias y
  credenciales reales del operador.
- Cobros automáticos: fuera de alcance por decisión del plan
  maestro; los borradores de esta etapa no cobran.
