# Trickreport: Plan de implementación de mejoras integrales

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convertir las capacidades actuales de Trickreport en un producto seguro, integrado, comprobable y operable, cerrando los hallazgos de la auditoría mediante 35 tareas en seis bloques.

**Architecture:** Conservar el monolito modular Go y el frontend Astro SSR con islas React. Mantener PostgreSQL como fuente de verdad para datos, sesiones y trabajos que deban sobrevivir reinicios. Reutilizar servicios, repositorios, cliente API, componentes UI y composición manual existentes.

**Tech Stack:** Go, chi, pgx, JWT, LDAP, Astro, React, TypeScript, Tailwind, PostgreSQL 16, Docker Compose y Caddy. Objetivo de runtime propuesto: Go 1.26.10 y Node 24.22.1 LTS, con revisión de parches al ejecutar.

**Spec:** [Auditoría integral del 21 de octubre](F:/Proyectos/Trickreport/graphify-out/AUDITORIA-2026-10-21.md), [baseline y verificaciones](F:/Proyectos/Trickreport/graphify-out/audit-evidence/summary.json). Este plan incorpora su alcance y mantiene una matriz de trazabilidad para que la especificación viaje con la ejecución.

Fecha de elaboración: 21 de octubre de 2026, America/Caracas. Baseline: `382b6417393f011ff41a1c8cb919a2385055fbe6`, rama local `dev`. **Documento de planificación: ninguna tarea de producto está ejecutada por haber escrito este archivo.**

## Global Constraints

- Desarrollo local nativo, sin Docker; contenedores para validación de imágenes y despliegue. PostgreSQL de pruebas separado del desarrollo y de producción.
- Preservar el cambio previo de formato de `ArticleEditor.tsx` y las tres mejoras locales de autenticación por tenant, recuperación y editor de artículos.
- Mantener PostgreSQL 16. Elevar runtimes de build, CI y ejecución conjuntamente; no usar el Node 25 local de la auditoría como referencia de producción.
- Reutilizar dependencias instaladas; añadir Playwright solo para el recorrido E2E que hoy falta. Outbox y sesiones usan PostgreSQL existente.
- Migraciones nuevas y aditivas, aplicadas mediante flags de la API. No reescribir las cuatro migraciones aplicadas ni volver a consolidarlas.
- Antes de editar símbolos: Codebase Memory, Graphify, cobertura de rutas y lectura de callers. Riesgo HIGH/CRITICAL se comunica antes de editar; UNKNOWN exige corroboración. Antes de cada commit: revisar el diff completo y validar los cambios relevantes, sin aceptar resultados parciales como definitivos.
- Commits Conventional Commits, pequeños y revisables; sin atribución AI. Un commit corresponde a un entregable verificable, no a todo un bloque.
- Cada defecto se convierte primero en una prueba que falla con el comportamiento esperado corregido. Los checks de auditoría que hacen PASS al reproducir defectos deben invertirse antes de incorporarlos como regresiones.
- Ninguna prueba destructiva usa la DB de trabajo. Migración, restore, expiraciones, fallos de SMTP y concurrencia se ejercitan en entornos desechables.
- Los bloques son áreas de trabajo; las pruebas acompañan cada arreglo. El bloque 3 consolida la barrera de entrega, no pospone pruebas hasta entonces.

## Review Focus

1. Dos refresh simultáneos, logout y reinicio: solo una rotación gana y ninguna credencial revocada recupera acceso. Pruebas B1.1.
2. Usuario A sale y usuario B entra en el mismo navegador: B no ve ni sincroniza datos de A; borradores antiguos no se reasignan automáticamente. Pruebas B1.5 y B2.7.
3. Escritura confirmada por DB seguida de fallo de respuesta: reintentar no crea otro ticket; cambiar el cuerpo con la misma clave devuelve conflicto. Pruebas B2.7.
4. SQL de automatización y worker: una asignación de otro tenant o un fallo a mitad de operación no rompe historial, versión o permisos. Pruebas B2.5 y B2.8.
5. Backup corrupto o DB con historial de migraciones anterior: el procedimiento se detiene antes de destruir datos y permite recuperar una copia consistente. Pruebas B4.2.

## Decisiones de trabajo propuestas

Estas decisiones hacen ejecutable el plan. Las que afectan al producto están identificadas como propuestas y se pueden ajustar antes de implementar su tarea; no se presentan como requisitos ya aprobados por el usuario.

| Tema | Decisión propuesta y alcance |
|---|---|
| Arquitectura | Mantener la estructura actual y corregir contratos; las excepciones SQL quedan documentadas y probadas |
| Despliegue inicial | Una réplica backend y una frontend, PostgreSQL y Caddy bajo un mismo origen HTTPS. La topología de desarrollo sigue siendo nativa |
| Sesiones | Access de 8 horas, refresh de 7 días como hoy; validación contra sesión activa en DB y estado/rol actuales del usuario. La rotación invalida el par anterior |
| MFA | Desafío separado de 5 minutos, un solo uso y máximo 5 intentos por desafío; nunca autoriza API ordinaria. Setup/disable requieren prueba del secreto almacenado |
| Contraseñas | Aplicar la política de dominio existente de 12 caracteres y complejidad a todas las altas/cambios locales; no cambiar la política LDAP |
| Offline | Cachear assets públicos y página offline estática. Mantener creación diferida de tickets con identidad explícita; retirar caché de HTML/API privados |
| Borradores existentes | Preservarlos; los que no tengan identidad quedan en cuarentena, sin visibilidad/sync automática en otra sesión. No inferir propietario del usuario recién conectado |
| Datos y adjuntos | Mantener adjuntos en PostgreSQL durante esta estabilización; corregir validación y descarga. S3 no forma parte de la primera entrega |
| Fiabilidad de email | Outbox PostgreSQL para eventos de negocio; entrega al menos una vez, con reintentos y fallos visibles. No prometer exactamente una entrega SMTP |
| Métricas | Usuario final ve sus tickets; admin/agente ven el tenant. Diferenciar creación, resolución y estado actual en KPIs y gráficas |
| Idioma | Mantener el inglés actual durante estabilización; traducción completa queda como decisión de producto posterior, sin introducir un sistema i18n ahora |
| Alcance posterior | SSO OAuth2, superadmin/gestión de tenants, envío Web Push, S3 y múltiples réplicas quedan en backlog explícito con requisitos propios |

**Impacto previsto para usuarios:** la corrección de sesiones requerirá volver a iniciar sesión con credenciales existentes; no se cambiarán sus contraseñas. Los borradores conservan datos, pero los antiguos sin propietario no se enviarán automáticamente. Los cambios de contraseña y desactivación revocarán sesiones; cambios de rol invalidarán credenciales que conserven el rol anterior.

## Orden de ejecución y entregas

| Bloque | Tareas | Prioridad | Puerta de salida |
|---|---:|---|---|
| 0. Base confiable | 3 | P0 | Rama y runtimes reproducibles; análisis de impacto disponible |
| 1. Identidad y privacidad | 7 | P0/P1 | Las pruebas negativas de identidad, tenant, serialización y caché pasan |
| 2. Flujos esenciales | 10 | P1 | Los recorridos principales terminan con API y DB reales |
| 3. Calidad y CI | 5 | P1 | Checks obligatorios verdes y contratos críticos sin skips |
| 4. Operación y entrega | 6 | P1/P2 | Arranque, recuperación, alertas y release ensayados |
| 5. Producto y documentación | 4 | P2 | Cada capacidad tiene definición, evidencia y alcance publicado |

**Ruta crítica:** B0 → B2.1 (config) + B3.3 (DB de prueba) + B3.1 (races) → B1 → resto de B2 → B3.4/B3.5 → B4 → cierre B5. B1.3 puede adelantarse tras B0: eliminar la serialización vulnerable no debe esperar al upgrade completo del frontend. B1.4 se entrega como unidad: transporte WebSocket y filtrado de destinatarios.

Cada tarea tiene un responsable funcional sugerido; no presupone personas adicionales ni autoriza crear chats, subagentes o mensajes. Las tareas se pueden implementar secuencialmente. Se estima duración después de B0 y de disponer de la DB de pruebas: hoy una fecha de cierre ocultaría incertidumbre sobre migraciones y compatibilidad de dependencias.

## Mapa de archivos

Los nombres cortos de las tareas remiten a estas rutas exactas; sus directorios son el alcance inicial, no permiso para una refactorización general.

| Código | Fuentes existentes |
|---|---|
| AUTH | [service.go](F:/Proyectos/Trickreport/backend/internal/application/auth/service.go), [mfa.go](F:/Proyectos/Trickreport/backend/internal/application/auth/mfa.go), [password_reset.go](F:/Proyectos/Trickreport/backend/internal/application/auth/password_reset.go) |
| AUTH-HTTP | [auth.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/handler/auth.go), [middleware.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/middleware/middleware.go), [csrf.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/middleware/csrf.go) |
| JWT/SES | [jwt.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/auth/jwt.go), [session_repo.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/session_repo.go), [user_repo.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/user_repo.go), [password_reset_repo.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/password_reset_repo.go) |
| WIRING | [server.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/server.go), [wire.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/wire.go) |
| CONFIG | [config.go](F:/Proyectos/Trickreport/backend/internal/config/config.go), [main.go](F:/Proyectos/Trickreport/backend/cmd/api/main.go), [bootstrap.go](F:/Proyectos/Trickreport/backend/internal/bootstrap/bootstrap.go) |
| TICKET | [service.go](F:/Proyectos/Trickreport/backend/internal/application/ticket/service.go), [ticket_repo.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/ticket_repo.go), [ticket.go handler](F:/Proyectos/Trickreport/backend/internal/interfaces/http/handler/ticket.go), [ticket.go dominio](F:/Proyectos/Trickreport/backend/internal/domain/ticket/ticket.go) |
| AUTO/SLA | [engine.go](F:/Proyectos/Trickreport/backend/internal/application/automation/engine.go), [service.go](F:/Proyectos/Trickreport/backend/internal/application/automation/service.go), [automation_executor.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/automation_executor.go), [worker.go](F:/Proyectos/Trickreport/backend/internal/worker/worker.go) |
| RT | [hub.go](F:/Proyectos/Trickreport/backend/internal/realtime/hub.go), [client.go](F:/Proyectos/Trickreport/backend/internal/realtime/client.go), [message.go](F:/Proyectos/Trickreport/backend/internal/realtime/message.go), [adapter.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/realtime/adapter.go), [logging.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/middleware/logging.go), [metrics.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/middleware/metrics.go) |
| API-UI | [api.ts](F:/Proyectos/Trickreport/frontend/src/lib/api.ts), [config.ts](F:/Proyectos/Trickreport/frontend/src/lib/config.ts), [auth.ts](F:/Proyectos/Trickreport/frontend/src/lib/auth.ts), [middleware.ts](F:/Proyectos/Trickreport/frontend/src/middleware.ts) |
| PWA | [sw.js](F:/Proyectos/Trickreport/frontend/public/sw.js), [register-sw.ts](F:/Proyectos/Trickreport/frontend/src/scripts/register-sw.ts), [db.ts](F:/Proyectos/Trickreport/frontend/src/lib/offline/db.ts), [sync.ts](F:/Proyectos/Trickreport/frontend/src/lib/offline/sync.ts), [useOfflineTickets.ts](F:/Proyectos/Trickreport/frontend/src/lib/hooks/useOfflineTickets.ts) |
| EMAIL | [queue.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/email/queue.go), [smtp.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/email/smtp.go), [sender.go](F:/Proyectos/Trickreport/backend/internal/email/sender.go) |
| DATA | [migraciones](F:/Proyectos/Trickreport/backend/internal/db/migrations), [tx.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/tx.go), [integration_test.go](F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/integration_test.go), [contract_test.go](F:/Proyectos/Trickreport/backend/internal/interfaces/http/handler/contract_test.go) |
| BUILD | [package.json](F:/Proyectos/Trickreport/frontend/package.json), [go.mod](F:/Proyectos/Trickreport/backend/go.mod), [Dockerfile frontend](F:/Proyectos/Trickreport/frontend/Dockerfile), [Dockerfile backend](F:/Proyectos/Trickreport/backend/Dockerfile), [ci.yml](F:/Proyectos/Trickreport/.github/workflows/ci.yml), [release.yml](F:/Proyectos/Trickreport/.github/workflows/release.yml) |
| OPS | [Compose principal](F:/Proyectos/Trickreport/docker-compose.yml), [on-premise](F:/Proyectos/Trickreport/docker-compose.onpremise.yml), [staging](F:/Proyectos/Trickreport/docker-compose.staging.yml), [monitoring](F:/Proyectos/Trickreport/docker-compose.monitoring.yml), [backup](F:/Proyectos/Trickreport/docker-compose.backup.yml), [Caddyfile](F:/Proyectos/Trickreport/Caddyfile), [restore.sh](F:/Proyectos/Trickreport/scripts/restore.sh) |

Los nuevos nombres de migración propuestos son `000005_session_security`, `000006_mfa_challenges`, `000007_password_reset_security`, `000008_ticket_idempotency`, `000009_sla_tracking` y `000010_email_outbox`, cada uno con `.up.sql` y `.down.sql` bajo DATA. B0 confirma que esos números siguen libres tras reconciliar Git; si cambian, se actualiza este registro y todos sus consumidores antes de escribir SQL.

## Bloque 0 — Recuperar una base confiable

### B0.1 — Preservar y reconciliar Git

**Responsable:** mantenimiento. **Depende:** ninguno. **Archivos:** trabajo local existente, refs Git y nuevo `F:/Proyectos/Trickreport/docs/operations/baseline.md`.

- [ ] Registrar estado y hashes; guardar una referencia de respaldo del HEAD y un parche del cambio local de ArticleEditor; preservar los archivos sin seguimiento. Comparar `git cherry`, `git diff origin/dev HEAD` y cambios de contenido de `main`.
- [ ] Crear la rama de estabilización con prefijo `codex/` desde el historial remoto reconciliado y trasladar solo los parches funcionales locales que falten, preservando el lookup por tenant, recuperación y editor. Resolver también la diferencia histórica de prueba; no asumir que todo hash distinto es trabajo nuevo.
- [ ] Verificar diff final por áreas, ejecutar builds y registrar los fallos de baseline conocidos; confirmar que no reaparecen secretos o archivos locales ignorados.
- [ ] Revisar y guardar el registro; un commit de reconciliación no incluye artefactos masivos ni mezclas de arreglos posteriores.

**Salida:** puede explicarse cada diferencia de la rama resultante y recuperarse la copia anterior. **Rollback:** referencia y parche conservados, sin force-push a ramas compartidas.

### B0.2 — Hacer fiable el análisis de alcance

**Responsable:** mantenimiento. **Depende:** B0.1. **Archivos:** configuración de índices y guías existentes; Codebase Memory y Graphify.

- [ ] Resolver por ruta el proyecto Codebase Memory y verificar `index_status`. Excluir artefactos de `graphify-out/` del análisis estructural del producto; registrar las exclusiones de SQL, scripts, tests y Astro, con fallback de fuente.
- [ ] Confirmar la disponibilidad de Codebase Memory y Graphify, y guardar la forma reproducible de invocarlos. No confundir las skills instaladas con un índice disponible.
- [ ] Ejecutar impacto de `Authenticate`, `Login`, `BroadcastEvent`, `matchConditions`, `api` y el flujo de sync; guardar callers, riesgo y huecos. Los resultados truncados se completan antes de editar esos símbolos.
- [ ] Revisar la evidencia y documentar el procedimiento en el baseline. Actualizar Graphify al cambiar código/especificaciones; no usarlo para certificar ausencia de llamadas.

**Salida:** impacto reproducible, ligado al checkout correcto. **Estado actual:** Codebase Memory fue indexado al preparar este plan; Graphify queda como apoyo para navegación y la cobertura todavía requiere corroboración de fuente.

### B0.3 — Unificar runtimes

**Responsable:** backend/frontend/infraestructura. **Depende:** B0.1. **Archivos:** BUILD, Makefiles, scripts de desarrollo, ejemplos de entorno y nuevo `F:/Proyectos/Trickreport/frontend/.node-version`.

- [ ] Confirmar en las fuentes oficiales que Go 1.26.10 y Node 24.22.1 siguen siendo parches adecuados; registrar la pareja elegida. Estos números son objetivos propuestos verificados al redactar, no requisitos históricos del repositorio.
- [ ] Alinear `go.mod`, CI, Dockerfiles, versión Node del proyecto y documentación; `engines.node` fija la línea 24 soportada. Conservar PostgreSQL 16.
- [ ] Ejecutar `go build ./...`, `go vet ./...`, `npm ci` y `npm run build` bajo esos runtimes. Registrar incompatibilidades de Astro/integraciones para B1.7, sin resolverlas mediante versiones forzadas incoherentes.
- [ ] Revisar lockfiles y registrar un commit de toolchain reproducible.

**Salida:** local, CI e imágenes apuntan al mismo runtime. **Rollback:** manifiestos/lockfile e imágenes previas, sujeto al criterio de seguridad de B1.7.

## Bloque 1 — Cerrar identidad y privacidad

### B1.1 — Sesiones persistentes, propósito de tokens y rotación atómica

**Responsable:** backend. **Depende:** B0, B2.1, B3.3. **Hallazgo:** AUTH-02. **Archivos:** AUTH, AUTH-HTTP, JWT/SES, WIRING; migración `000005_session_security`; pruebas de auth, jwt, middleware, servidor y contratos.

**Contrato:** `Service.AuthenticateAccess(ctx context.Context, token string) (*Claims, error)` es la única entrada del middleware para autorizar. Reutilizar `SessionInput` con `JTI` de access y `RefreshJTI` de refresh. Extender el repositorio con `FindActiveByAccessJTI(ctx, jti) (*SessionInput,error)`, `RotateSession(ctx, oldRefreshJTI, replacement SessionInput) (bool,error)` y `RevokeAllForUser(ctx,userID) error`.

- [ ] Añadir regresiones `TestAuthenticate_RejectsRefreshAndRevoked`, `TestRefresh_ConcurrentSingleWinner` y `TestSession_SurvivesRestart`: exigir 401 para refresh como access, logout, JWT sin sesión y rol/usuario cambiado; dos refresh producen un éxito y un 401. Ejecutarlas y observar el fallo inicial.
- [ ] Persistir el par al login; consultar sesión activa, usuario, tenant y rol actuales; rotar con compare-and-set SQL. Propagar fallos de persistencia. Logout revoca la sesión identificada; la blacklist en memoria deja de ser la autoridad.
- [ ] Actualizar todos los callers de `Authenticate`, mocks y `buildContractServer`; verificar paquetes auth/middleware/handler/http con `-race`, más integración PostgreSQL.
- [ ] Analizar cambios y revisar/commit. Incluir nota de despliegue: tokens anteriores sin sesión dejan de funcionar y requieren login.

### B1.2 — Desafío MFA separado y de un solo uso

**Responsable:** backend/frontend. **Depende:** B1.1. **Hallazgo:** AUTH-01. **Archivos:** AUTH, JWT, AUTH-HTTP, WIRING, `LoginForm.tsx`, `MfaSettings.tsx`; crear `F:/Proyectos/Trickreport/backend/internal/infrastructure/postgres/mfa_challenge_repo.go` y migración `000006_mfa_challenges`.

**Contrato:** `TokenGenerator.GenerateMFA(userID,tenantID uuid.UUID,role string) (string,error)` genera `token_type=mfa`, expiración de 5 minutos. Un puerto `MFAChallengeRepository` guarda JTI, usuario, tenant, vencimiento, intentos y consumo; los límites y consumo se actualizan atómicamente en PostgreSQL. Los access/refresh normales no completan MFA.

- [ ] Escribir `TestMFA_ChallengeCannotAuthorizeAPI`, `TestMFA_ExpiredReusedAndSixthAttemptRejected`, `TestMFA_EnableUsesStoredSecret` y `TestMFA_DisableRequiresCode`; comprobar fallos con 401 y ausencia de credenciales completas antes de TOTP.
- [ ] Crear desafío, aplicar máximo 5 intentos, consumirlo una sola vez y emitir sesión B1.1 solo tras validar TOTP. En enable verificar el secreto almacenado, sin confiar en uno arbitrario enviado por el cliente; exigir código válido para disable. Aplicar el rate limiter existente al segundo paso.
- [ ] Actualizar mocks de TokenGenerator y puertos; verificar rutas API reales y test UI del segundo paso, sin guardar el desafío en logs ni almacenamiento permanente del navegador.
- [ ] Revisar cambio de contrato, migración y commit conjunto backend/frontend.

### B1.3 — Retirar serialización de datos no confiables en scripts

**Responsable:** frontend. **Depende:** B0. **Hallazgo:** SEC-01. **Archivos:** `F:/Proyectos/Trickreport/frontend/src/pages/tickets.astro`, API-UI y exportación; crear `F:/Proyectos/Trickreport/frontend/src/lib/export.test.ts` y regresión E2E en B3.4.

- [ ] Reproducir con título que contiene cierre de script, comillas y marcador HTML inerte; la aserción corregida exige que siga siendo texto y no cree un nodo externo al contenedor esperado.
- [ ] Retirar `define:vars` del JSON de tickets. Usar un atributo HTML fijo `data-tickets` escapado por Astro y un script procesado que lee `dataset`; no usar `set:html` ni nombres de atributos controlados por datos. Conservar filtros y filas exportadas.
- [ ] Probar todos los `define:vars` restantes según procedencia de datos, build y `astro check`. La regresión debe ejecutar el render real, además del serializador aislado.
- [ ] Revisar/commit como parche independiente; B1.7 mantiene obligación de corregir el framework y demás advisories.

### B1.4 — WebSocket con destinatarios autorizados y transporte funcional

**Responsable:** backend/frontend. **Depende:** B1.1. **Hallazgos:** DATA-01 y RT-01. **Archivos:** RT, TICKET, WIRING, `NotificationBell.tsx` y sus tests.

**Contrato:** reemplazar broadcast tenant-global por `SendToUsers(tenantID uuid.UUID,userIDs []uuid.UUID,eventType string,data any)` en el puerto de tickets y adaptador. Resolver destinatarios con consultas por tenant: evento de ticket público → creador y personal activo autorizado; nota interna → personal; notificación → solo `notification.UserID`. Los payloads son DTO explícitos, no agregados Go completos.

- [ ] Escribir pruebas de handshake a través del router completo y matriz de dos usuarios finales, agente, admin y segundo tenant: 101 para cliente válido; cero frames privados para otros usuarios; rol/campos JSON compatibles con el frontend.
- [ ] Preservar `Hijacker`/`Flusher` requeridos en logging/métricas usando las utilidades de chi instaladas si encajan. Aplicar la selección de destinatarios en todos los emisores y limitar vida de socket a sesión vigente; logout, expiración y desactivación cierran conexiones sin necesitar un nuevo evento.
- [ ] Verificar conexión, reconexión, polling sin duplicados, nota interna y cierre por sesión revocada con API real. No entregar solo el arreglo de handshake.
- [ ] Revisar DTO, callers y `-race`; commit atómico de transporte y privacidad.

### B1.5 — Aislar caché y borradores offline por identidad

**Responsable:** frontend. **Depende:** B1.1. **Hallazgo:** DATA-02. **Archivos:** PWA, `Sidebar.tsx`, `logout.astro`, `OfflineTicketForm.tsx`, `SyncStatus.tsx`, `offline.test.ts`.

**Contrato:** `PendingTicket` añade `userId` y `tenantId`; lectura y sync requieren esa identidad explícita. Caché HTTP privada eliminada; solo assets públicos y fallback offline estático se cachean. Una migración IndexedDB preserva registros antiguos sin adjudicarles propietario.

- [ ] Escribir `it('never shows or syncs user A drafts as B')`, `it('does not return private cached HTML/API after logout')` y `it('quarantines unowned legacy drafts without deleting them')`; repetir con mismo usuario en otro tenant.
- [ ] Cambiar estrategias SW, subir versión y limpiar caches privadas antiguas durante activate; proteger la migración IDB ante abortos y fallos. Logout y cambio de identidad limpian la vista, cierran socket y cancelan sync; borradores propios conservados quedan inaccesibles para otra cuenta.
- [ ] Verificar login A → crear borrador → logout → login B → offline → volver a A. Los pendientes sin identidad se muestran solo como aviso de recuperación pendiente, sin autosync ni reasignación silenciosa.
- [ ] Revisar compatibilidad de SW/IDB y commit; documentar limitación de privacidad en dispositivo compartido y procedimiento de recuperación de borradores antiguos.

### B1.6 — Política común de credenciales y permisos

**Responsable:** backend/frontend. **Depende:** B1.1/B1.2. **Archivos:** AUTH, aplicación/handler user, dominio user/password, AUTH-HTTP, `RegisterForm.tsx`, `UsersManager.tsx`, `MfaSettings.tsx` y tests correspondientes.

- [ ] Añadir tabla de casos: alta local con 6/11 caracteres rechazada; contraseña de 12 con requisitos del dominio aceptada; agent no administra usuarios; end_user no asigna tickets ni crea notas internas; `X-Tenant-ID` no cambia el tenant firmado. Probar usuario desactivado y cambio de rol con JWT anterior.
- [ ] Reutilizar `ValidatePasswordComplexity`, eliminar validaciones divergentes y alinear textos/atributos frontend. Vincular cambios de contraseña/desactivación a revocación B1.1; nunca alterar el tenant de una cuenta desde DTO de edición común. Verificar redirects y control servidor de todas las páginas admin.
- [ ] Ejecutar pruebas de cada rol y tenant, incluidos LDAP opcional, lockout y endpoints MFA; respuestas 401/403/validación sin datos de otro usuario.
- [ ] Revisar matriz de permisos y commit; no cambiar la política del directorio LDAP.

### B1.7 — Actualizar y clasificar dependencias vulnerables

**Responsable:** backend/frontend/infraestructura. **Depende:** B0.3, B1.3. **Hallazgo:** DEP-01. **Archivos:** BUILD, lockfile frontend, go.sum, dependencias runtime de imágenes.

- [ ] Repetir audit completo y producción, govulncheck bajo Windows y Linux, y enumerar cada aviso con versión, uso real, corrección y prueba afectada. Añadir al registro la incompatibilidad Node 20/integración React.
- [ ] Actualizar Astro, adapter, integración React y Vite como conjunto compatible según las guías mayores oficiales; conservar React/Tailwind si sus contratos siguen válidos. Corregir Go/toolchain/x/text y dependencias transitivas. Retirar o mover fuera del runtime CLI como shadcn si ningún import productivo los necesita; no ejecutar `audit fix --force` sin revisar el árbol.
- [ ] Exigir builds, tipos, hidratación, auth, SW y render seguro B1.3 bajo el runtime objetivo. Cerrar avisos altos/críticos aplicables al runtime; cualquier falso positivo/exclusión necesita evidencia concreta y registro, no silenciar el scanner.
- [ ] Revisar diff de paquetes, imagen resultante y commit; no volver a una versión vulnerable como rollback funcional.

## Bloque 2 — Recuperar y completar flujos esenciales

### B2.1 — Configuración sin archivo local

**Responsable:** backend. **Depende:** B0.3. **Hallazgo:** DEPLOY-01. **Archivos:** CONFIG, `config_test.go`, ejemplos de entorno y pruebas de bootstrap.

- [ ] Escribir `TestLoad_EnvironmentOnly` con directorio vacío y env sintético: DATABASE_URL/JWT_SECRET/ADMIN_EMAIL/ADMIN_PASSWORD/SMTP/LDAP llegan íntegros; puerto y defaults siguen funcionando. Faltas obligatorias de producción causan salida con mensaje de clave, sin secreto.
- [ ] Registrar explícitamente las claves Viper y validar configuración antes de abrir DB. `staging` usa las mismas validaciones de seguridad que producción; mantener HTTP local solo en perfil de desarrollo/LAN explícito. Validar password bootstrap con la política acordada.
- [ ] Ejecutar subprocess tests development/staging/production, sin imprimir DSN ni credenciales. Verificar que flags de migración siguen funcionando bajo configuración controlada.
- [ ] Revisar ejemplos, tests y commit. Es la primera corrección de código propuesta porque permite ensayar las demás en procesos reales.

### B2.2 — Contrato CSRF común para navegador y SSR

**Responsable:** backend/frontend. **Depende:** B2.1 y B1.1. **Hallazgo:** FLOW-01. **Archivos:** AUTH-HTTP, WIRING, API-UI, páginas new ticket/admin SLA/logout y tests.

**Contrato:** extender opciones del cliente con `csrfToken?: string` para peticiones SSR; en servidor reenviar solo cookie CSRF y header correspondiente junto al Bearer. En navegador mantener cookies e inicializar con `GET /api/v1/auth/csrf` → 204 antes de una mutación que necesite token. No renombrar el cliente ni añadir wrappers por dominio.

- [ ] Crear pruebas de registro inicial, login → POST SSR ticket, PUT SSR SLA y logout; con token correcto se completa, sin él o con origen externo se rechaza. Dos GET concurrentes no invalidan el token de una pestaña.
- [ ] Añadir bootstrap CSRF, reutilizar cookie válida en GET y generar una nueva solo cuando falta/es inválida; sustituir fallback fijo por error seguro de generación. Pasar CSRF desde Astro.cookies en los formularios SSR. Mantener protección de los endpoints autenticados y validación de origen del frontend.
- [ ] Ejecutar contratos y browser real por el mismo origen; servidor no intenta leer `document`, navegador no intenta fijar `Cookie` manualmente.
- [ ] Revisar controles y commit backend/frontend; el arreglo no consiste en eximir toda petición Bearer.

### B2.3 — Sesión frontend y recuperación de estados inválidos

**Responsable:** frontend. **Depende:** B1.1/B2.2. **Hallazgo:** UX-01. **Archivos:** API-UI, páginas login/register/dashboard/admin/tickets/artículos y Layout.

- [ ] Crear `middleware.test.ts`: cookie inválida termina en login sin bucle; cookie válida permite dashboard; backend caído presenta indisponibilidad recuperable, no una redirección circular. Medir una resolución de usuario por request SSR.
- [ ] Basar login/register en `Astro.locals.user`, limpiar cookies inválidas y reutilizar locals en páginas. Distinguir 401 de error de transporte. En browser, un 401 lleva a reautenticación conservando borradores B1.5; en SSR no redirigir por simple existencia de cookie.
- [ ] Repetir cookies caducadas, logout con fallo de red, carga directa de rutas admin y navegación móvil; ninguna ruta requiere getMe duplicado en middleware y página.
- [ ] Revisar errores visibles y commit, conservando el cliente API central.

### B2.4 — Recuperación por email completa y de consumo atómico

**Responsable:** backend/frontend. **Depende:** B1.1/B1.6/B2.1. **Archivos:** AUTH/reset, JWT/SES reset repo, WIRING, CONFIG y formularios reset; migración `000007_password_reset_security`.

**Contrato:** URL pública absoluta configurada con `APP_PUBLIC_URL`; token aleatorio guardado como hash. Reutilizar `TxManager` para bloquear/validar token, cambiar password, marcar consumo y revocar sesiones en una transacción. Vigencia de una hora conservada.

- [ ] Crear `TestReset_AbsoluteLinkSingleUseAndAtomic`, con dos confirmaciones simultáneas: una gana; fallo SMTP deja evidencia operativa; DB fallida no deja password cambiado con token reutilizable. POST reset devuelve 202 y el mismo cuerpo genérico para cualquier email con formato válido, exista o no; no revelar el fallo SMTP específico en la respuesta pública.
- [ ] Construir enlace desde origen configurado, no desde Host no confiable. Persistir hashes y consumir bajo transacción; invalidar tokens anteriores al pedir uno nuevo y aplicar límite por cuenta/IP. Registrar y medir fallos operativos sin token ni enlace sensible; B4.3 añade su reintento durable, conservando la misma respuesta pública.
- [ ] Probar con capturador SMTP local: solicitar, recibir, abrir, cambiar, rechazar contraseña anterior y token reutilizado. Invalidar tokens legacy en despliegue con nota al usuario; no conservar secretos raw como fallback.
- [ ] Revisar migración/puertos/mocks y commit; migración no modifica contraseñas existentes.

### B2.5 — Unificar condiciones y acciones de automatización

**Responsable:** backend/frontend. **Depende:** B1.6/B3.3. **Hallazgo:** AUTO-01. **Archivos:** AUTO/SLA, dominio automation, handler automation, `AutomationsManager.tsx`, tests de motor/servicio/SQL/UI.

**Contrato:** forma canónica `{conditions:[{field,op,value}]}` con AND, campos status/priority/category/title/old_value/new_value y ops equals/not_equals/contains. Admitir lectura del mapa legacy como igualdades y normalizar sus alias `oldvalue`/`newvalue`; rechazar estructuras/campos/acciones desconocidos al guardar. API devuelve el formato canónico.

- [ ] Escribir pruebas equivalentes de mapa legacy y cláusulas; dos condiciones verdaderas ejecutan una acción, una falsa ninguna; payload inválido devuelve 400. Estado imposible, UUID de otro tenant y usuario inactivo no producen escrituras parciales.
- [ ] Normalizar y validar en el servicio compartido; ajustar UI y motor al mismo contrato. En ejecutor SQL aplicar transiciones, tenant/rol del asignado, versión e historial en transacción usando reglas/helpers existentes; no crear un segundo servicio de tickets.
- [ ] Verificar los disparadores reales creados por tickets y worker, las tres acciones visibles y errores de acciones. `add_tag` deja de aparentar éxito: se rechaza como no soportado hasta que exista una capacidad real.
- [ ] Revisar compatibilidad de reglas guardadas y commit; no migrar reglas válidas de forma destructiva.

### B2.6 — Completar consistencia de tickets y asignación

**Responsable:** backend/frontend. **Depende:** B1.6/B2.2/B3.3. **Archivos:** TICKET, `TicketActions.tsx`, `TicketFilters.tsx`, tickets.astro, API-UI y pruebas.

- [ ] Probar asignar y desasignar: `assigned_to=null` termina en NULL SQL; otro tenant, end_user como agente e inactivo se rechazan. Dos cambios de estado sobre versión anterior producen un conflicto explícito, no un 500 genérico.
- [ ] Alinear `AssignReq`/servicio/repositorio/cliente con ausencia real de asignado; mapear conflicto/validación a respuestas útiles. Enviar filtros soportados a la API y paginar sin que exportación o lista contradigan los criterios visibles; asegurar orden estable.
- [ ] Verificar crear → asignar → comentar → cambiar estado → cerrar y reabrir según matriz de dominio; notas internas no aparecen para usuario final por HTTP ni WS. Historial registra cambios una vez.
- [ ] Revisar contratos backend/frontend y commit; las transiciones existentes solo cambian si esta prueba revela un defecto.

### B2.7 — Idempotencia de alta y sincronización offline

**Responsable:** backend/frontend. **Depende:** B1.5/B2.6/B3.3. **Archivos:** TICKET, DATA, API-UI y PWA; migración `000008_ticket_idempotency`; tests PostgreSQL y offline.

**Contrato:** POST tickets acepta `Idempotency-Key` UUID, único por tenant+usuario+operación. En la misma transacción guarda hash de payload y ticket resultante; repetición idéntica devuelve el mismo recurso, cuerpo distinto devuelve 409. El UUID local del borrador es la clave estable.

- [ ] Crear `TestCreateTicket_IdempotencySurvivesResponseLoss` y concurrencia: después de persistir y perder la respuesta, reintento crea cero filas adicionales; otro payload con la clave provoca 409; otra identidad no recupera el resultado.
- [ ] Reutilizar transacciones existentes para reservar clave y crear ticket. Mantener claves mientras exista el ticket. Browser y SSR generan una clave por envío lógico; offline conserva la misma hasta confirmar. Reintentar automáticamente solo lecturas o escrituras idempotentes.
- [ ] Impedir sync simultánea de pestañas/contextos mediante estado IDB y backend idempotente; solo borrar pendiente tras respuesta confirmada. Detenerse en 401/403 y diferenciar validación de conflicto, sin reenviar bajo otra cuenta.
- [ ] Verificar desconexión, recarga, dos pestañas y respuesta 5xx; revisar/commit con migración y actualización del SW.

### B2.8 — SLA y worker consistentes ante concurrencia/fallos

**Responsable:** backend. **Depende:** B2.5/B2.6/B3.3. **Archivos:** AUTO/SLA, repositorios SLA/tickets, dominio ticket, DATA; migración `000009_sla_tracking` si faltan timestamps necesarios.

- [ ] Añadir integración para resolución/escalación una vez, dos workers simultáneos, ticket cerrado y error al insertar historial. Añadir primera respuesta de personal y su deadline; fijar que comentario interno no cuenta como respuesta al cliente.
- [ ] Respetar deadline persistido y no recalcular silenciosamente contratos de tickets antiguos al editar una política. Reclamar filas con locking y actualizar marca+comentario+historial en una transacción. Registrar `first_responded_at`/`resolved_at` para métricas reales; datos históricos desconocidos quedan NULL.
- [ ] Probar fallo intermedio seguido de retry, cambio de prioridad y cierre/reapertura. Regla propuesta de política: respuesta y resolución mayores que cero; escalación cero deshabilita ese umbral y un valor positivo lo activa; todos los negativos se rechazan. Reflejarlo en formulario, API y worker.
- [ ] Revisar migración y semántica de SLA; commit. La primera respuesta y la resolución se presentan como métricas separadas.

### B2.9 — Adjuntos con validación real y descarga segura

**Responsable:** backend/frontend. **Depende:** B1.6/B3.3. **Archivos:** attachment_service.go, handler attachment, attachment_repo.go, `TicketActions.tsx` y tests.

- [ ] Añadir tabla de MIME declarado vs bytes: PDF/imagen válidos, tipo falsificado, archivo vacío, tamaño sobre límite y MIME con parámetros. Probar acceso a adjunto mediante otro ticket/tenant y límites multipart incluidos overhead.
- [ ] Usar detección de contenido de la biblioteca estándar; validar formatos contenedor permitidos mediante su estructura usando dependencias existentes. Sanear nombre y Content-Disposition, aplicar `nosniff`, límites y streaming donde el repositorio actual lo permita; conservar DB como almacenamiento.
- [ ] Verificar upload/list/download reales para roles autorizados y rechazos 403/404/413/415. Un tipo válido no implica garantía de ausencia de malware: documentar esa limitación.
- [ ] Revisar lista soportada y commit sin conectar S3 o storage local innecesariamente.

### B2.10 — Cerrar recorridos de base de conocimiento

**Responsable:** frontend/backend. **Depende:** B1.3/B1.6/B2.3. **Archivos:** application/article/service.go, article_repo.go, handler article, `ArticleEditor.tsx`, páginas articles y API-UI.

- [ ] Añadir pruebas de crear borrador → editar → publicar → buscar: usuario final solo ve publicados; agente no puede borrar si el backend reserva borrado al admin; una URL directa no salta permisos. Contenido HTML se conserva como texto inerte.
- [ ] Corregir únicamente gaps de esos recorridos, validación de título/contenido, manejo 404/403 y conservación del formulario ante error. Añadir acción de borrado admin con confirmación clara si se mantiene CRUD como capacidad anunciada; conservar editor de texto plano.
- [ ] Verificar búsqueda, categorías/tags, recarga y permisos en los tres roles; no se pierde el cambio local previo de formato.
- [ ] Revisar tests de ArticleEditor y commit; editor enriquecido y sugerencias automáticas no se añaden a esta entrega.

## Bloque 3 — Calidad y CI como barrera de entrega

### B3.1 — Eliminar carreras de los dobles de prueba

**Responsable:** backend. **Depende:** B0. **Archivos:** `F:/Proyectos/Trickreport/backend/internal/application/ticket/service_test.go` y doubles realmente afectados.

- [ ] Ejecutar `go test ./internal/application/ticket -race -count=10` y conservar la reproducción de lecturas concurrentes de mockEmail/mockEngine.
- [ ] Reemplazar esperas `time.Sleep` por canales/completions y sincronizar el acceso a los doubles; no modificar la concurrencia productiva para acomodar un mock inseguro.
- [ ] Repetir el paquete con `-race -count=10` y después `go test ./... -race`. Cero carreras, cero intermitencias observadas en esa serie; documentar que no es prueba matemática de ausencia.
- [ ] Revisar/commit como cambio de pruebas separado de arreglos funcionales.

### B3.2 — Tipos y lint con alcance real

**Responsable:** frontend/backend. **Depende:** B0.3/B1.7. **Archivos:** eslint.config.js, package.json, vitest.config.ts, Makefiles, .golangci.yml y archivos con diagnósticos reales.

- [ ] Verificar que `eslint --print-config` devuelve configuración para api.ts y LoginForm.tsx, y que una infracción de muestra se detecta. Registrar los 33 hints y distinguir deprecación de defecto.
- [ ] Configurar el parser TypeScript instalado para ts/tsx y Astro; añadir scripts `check` y `test:ci`. Fijar workers CI a una cantidad reproducible según recursos, sin silenciar errores de arranque. Corregir diagnósticos pertinentes y actualizar APIs obsoletas cuando la sustitución exista.
- [ ] Ejecutar `npm run lint`, `npm run check`, `npm run test:ci`, `npm run build`, gofmt check, `go vet` y golangci-lint. No usar `--fix` sobre todo el árbol para producir un diff de formato ajeno.
- [ ] Revisar configuración y commit; no añadir tests que solo dupliquen clases CSS o implementación trivial.

### B3.3 — Harness PostgreSQL de pruebas y contratos

**Responsable:** backend/QA. **Depende:** B0.3/B2.1. **Archivos:** DATA tests, crear `F:/Proyectos/Trickreport/backend/internal/testutil/database.go`, fixtures y scripts nativos de prueba.

- [ ] Provisionar DB desechable dedicada con nombre/DSN de prueba explícito. Añadir guardia que impide usar una DB que no esté identificada como test; no imprimir password.
- [ ] Reutilizar helper para migrar por API/función de migración, crear tenants+usuarios de los tres roles y limpiar. Hacer contractSetup seguro en concurrencia. La suite de integración falla si falta TEST_DATABASE_URL cuando se invoca como suite obligatoria.
- [ ] Ejecutar los seis CRUD y cinco contratos hoy omitidos con DB real; después ampliar con casos B1/B2. Fixture de login debe crear sesión real, no fabricar JWT que el sistema seguro ya no acepte.
- [ ] Revisar aislamiento entre paquetes y commit del harness; su limpieza solo afecta la DB desechable.

### B3.4 — Recorridos E2E por rol y dispositivo

**Responsable:** QA/frontend/backend. **Depende:** B1/B2/B3.3. **Archivos nuevos:** `F:/Proyectos/Trickreport/frontend/playwright.config.ts` y `F:/Proyectos/Trickreport/frontend/e2e/{auth,tickets,articles,admin,offline,realtime}.spec.ts`; package.json.

- [ ] Añadir Playwright como herramienta de desarrollo; iniciar API, Astro y DB de prueba con el mismo origen y configuración soportada. Usuarios sintéticos admin/agent/end_user en dos tenants, sin contraseñas de desarrollo reales.
- [ ] Escribir recorridos de login/MFA/logout/reset; ticket/asignación/comentario/adjunto; borrador/publicación; edición SLA/regla; notificación por socket; offline A/B e idempotencia. Incluir contenido hostil inerte y solicitudes directas que la UI oculta.
- [ ] Ejecutar Chromium en 1440×900 y 390×844; exigir sin errores de consola de la aplicación, controles accesibles por nombre y sin desbordamiento horizontal en rutas críticas. Guardar trazas solo en fallo y redactar datos sensibles.
- [ ] Revisar y commit de casos; no sustituir API real por mocks en la suite de aceptación.

### B3.5 — CI completo y protección de ramas

**Responsable:** mantenimiento/infraestructura. **Depende:** B3.1–B3.4/B1.7. **Archivos:** BUILD, workflows y configuración de rama GitHub durante ejecución autorizada.

- [ ] Construir jobs backend build/vet/lint/race; frontend lint/check/tests/build; integración con servicio PostgreSQL 16; E2E; audit y scan de imágenes. Corregir ruta de govulncheck y fijar versiones/SHAs de herramientas y Actions compatibles con políticas del builder SLSA.
- [ ] Hacer fallar integración por DB ausente y verificar conteo de skips críticos cero; conservar reportes y perfil de cobertura de una ejecución satisfactoria. Añadir `permissions` mínimos y evitar secrets en jobs de PR no confiables.
- [ ] Probar una ejecución satisfactoria de la rama y una regresión controlada que bloquea el job correcto. Marcar checks requeridos y PR antes de merge en main/dev; no introducir bypass por `continue-on-error`.
- [ ] Revisar configuración remota y commit; retirar la regresión controlada. La protección no se da por aplicada solo por describirla en YAML.

## Bloque 4 — Operación, recuperación y entregas

### B4.1 — Una topología desplegable y configuración SSR correcta

**Responsable:** infraestructura/frontend/backend. **Depende:** B1/B2/B3.5. **Archivos:** BUILD, OPS, API-UI config y ejemplos de entorno.

- [ ] Definir Compose base con postgres/backend/frontend/caddy; overlays on-premise y staging reutilizan esa base. Separar `API_URL=http://backend:8080` del SSR y `PUBLIC_API_URL=''` del navegador; on-premise no convierte SSR a URL relativa.
- [ ] Alinear variables de build/runtime del framework seleccionado en B1.7, `HOST=0.0.0.0`, `PORT=4321`, cookies HTTPS y validación production/staging. Exigir secretos con errores explícitos; LDAP/SMTP declarados cuando se habilitan. Dar escritura solo a directorios necesarios y reducir herramientas del runtime Node.
- [ ] En un host con Docker disponible, construir sin caché ni `.env` dentro de imagen; comprobar health, rutas, login, escritura SSR y WS detrás de Caddy. El Compose principal sirve también la UI. Registrar imágenes, digests y recursos.
- [ ] Revisar/commit con comandos documentados. Desarrollo nativo conserva scripts funcionales y mismos contratos de origen.

### B4.2 — Migraciones, backup y restore ensayados

**Responsable:** backend/operación. **Depende:** B4.1/B3.3. **Archivos:** DATA, OPS backup/restore, crear `F:/Proyectos/Trickreport/docs/operations/recovery.md`.

- [ ] Ensayar instalación limpia y upgrade de una copia representativa del esquema anterior consolidado. Si la transición no es compatible, implementar puente explícito probado; no modificar retrospectivamente una migración aplicada.
- [ ] Corregir DB/usuario configurables de backup y descubrimiento de red/servicio en restore. Validar gzip/SQL y restaurar primero en DB nueva; comprobar conteos/relaciones antes de sustituir la DB objetivo. Mantener respaldo verificable de la anterior y excluir restauraciones destructivas sobre destino implícito.
- [ ] Probar backup íntegro, corrupto y vacío, falta de espacio/error SQL, y recuperación del esquema con nuevas sesiones/outbox. Medir pérdida máxima y duración reales del ensayo; no inventar RPO/RTO ni prometer que pg_dump recupera toda pérdida reciente.
- [ ] Revisar procedimiento, rollback de cada nueva migración y commit. La aprobación de un corte productivo posterior se basa en el ensayo, no en ejecutar restore sobre desarrollo.

### B4.3 — Email durable y fallos visibles

**Responsable:** backend. **Depende:** B2.4/B2.8/B3.3. **Archivos:** EMAIL, TICKET/AUTO, DATA y WIRING; migración `000010_email_outbox`; crear repositorio outbox y tests en infraestructura/postgres.

- [ ] Añadir `TestOutbox_SurvivesRestartAndRetries`, cola llena, SMTP intermitente y fallo antes/después de confirmar ticket. Una transacción fallida no envía email; un evento comprometido no desaparece al reiniciar.
- [ ] Insertar evento outbox en la misma transacción de negocio, reclamar con lease/locking PostgreSQL, registrar intentos, reintentar con límite y mantener fallos terminales consultables. Usar emisor SMTP existente y clave de evento para evitar inserciones duplicadas.
- [ ] Probar apagado, recuperación de lease y timeout SMTP. Configuración productiva no cae silenciosamente a ConsoleSender que imprime enlaces de reset. Documentar posible duplicado de entrega tras crash entre envío y confirmación.
- [ ] Revisar logs sin contenido sensible y commit. Retirar la cola en memoria como autoridad del email de negocio cuando el outbox esté conectado.

### B4.4 — Readiness, shutdown y logging útiles

**Responsable:** backend/infraestructura. **Depende:** B1.4/B2.8/B4.3. **Archivos:** CONFIG/main, WIRING, worker, hub, limiters, EMAIL y tests HTTP/worker.

- [ ] Añadir pruebas de worker detenido, último ciclo fallido/antiguo, SMTP inaccesible y cancelación durante request. Distinguir liveness del proceso, readiness DB/worker y degradación de email sin presentar todo como saludable.
- [ ] Exponer heartbeat/último error de componentes; inicializarlo correctamente para el primer ciclo. Esperar cierre acotado de worker/outbox, sockets y goroutines de limpieza; fijar timeout de 10 segundos del servidor como presupuesto total documentado.
- [ ] Unificar logs estructurados a stdout en contenedores; archivo local opcional con rotación explícita si se conserva. Verificar request_id consistente y ausencia de tokens, passwords, reset links y DSN completos en errores.
- [ ] Revisar tests con `-race` y commit; una caída de SMTP no tumba la API si la outbox puede preservar el trabajo.

### B4.5 — Métricas, alertas y controles operativos

**Responsable:** operación/backend. **Depende:** B4.1/B4.4. **Archivos:** monitoring/*.yml, OPS monitoring, metrics.go, crear `F:/Proyectos/Trickreport/monitoring/rules.test.yml`.

- [ ] Probar reglas con promtool: backend down, error rate y worker obsoleto; comprobar que las series corresponden a las métricas emitidas. Montar rules.yml y validar configuración antes del arranque.
- [ ] Cambiar etiqueta de path a patrón chi de cardinalidad acotada; sin IDs/email. Añadir series outbox pendientes/fallidos y progreso worker. Configurar Alertmanager con receptor real del entorno y credenciales externas, desactivar password admin por defecto y limitar puertos de administración.
- [ ] Ensayar caída/fallo y recuperación; comprobar notificación recibida en canal de prueba y resolución. Definir panel básico para requests, latencia, errores, DB, worker y email; fijar versiones de imágenes.
- [ ] Revisar ruido/sensibilidad y commit. Las pruebas de notificación usan un destino autorizado al ejecutar, sin enviar mensajes a personas desde la planificación.

### B4.6 — Release reproducible y vuelta atrás

**Responsable:** mantenimiento/infraestructura. **Depende:** B3.5/B4.1–B4.5. **Archivos:** release.yml, .slsa-goreleaser/*.yml, empaquetado frontend y documentación de operación.

- [ ] Verificar rutas `dir/main` de cada builder, input manual tag-name y asociación entre tag, commit y artefactos. La release exige checks del mismo commit, no el último CI de otra referencia.
- [ ] Empaquetar backend e imagen/frontend con sus dependencias de runtime; incluir checksums, procedencia, versión y contrato de migración. Usar actions/refs permitidos por el builder de procedencia, comprobando sus restricciones antes de fijarlos.
- [ ] Ensayar una release candidata y ejecutar su artefacto en entorno limpio. Hacer rollback compatible con DB migrada; si requiere restore, usar B4.2 y declarar el impacto de datos y sesiones.
- [ ] Revisar y registrar candidato, instrucciones y evidencia. Publicación definitiva o despliegue productivo son acciones posteriores a este plan, con resultado concreto revisable.

## Bloque 5 — Consolidar producto, arquitectura y documentación

### B5.1 — KPIs y exportaciones con significado común

**Responsable:** producto/backend/frontend. **Depende:** B2.8/B2.6. **Archivos:** aplicación/handler analytics, analytics_repo.go, dashboard-metrics.ts, DashboardCharts/AnalyticsCharts, dashboard/admin analytics, export.ts y tests.

- [ ] Definir y probar con dataset fijo: creación por created_at, resolución por resolved_at, abiertos por estado actual, primera respuesta y cumplimiento SLA separado. End_user solo cuenta sus tickets; agente/admin el tenant. NULL histórico no se presenta como tiempo cero.
- [ ] Alinear consultas y UI al contrato; usar agrupación de días en UTC y mostrar esa zona en gráficas. Consolidar CSV en el helper existente y neutralizar fórmulas de hoja de cálculo, manteniendo comillas/saltos. El PDF por impresión conserva texto seguro y estados de popup bloqueado.
- [ ] Comprobar que summary, gráfica, filtros y exportación cuadran; días sin actividad valen cero y divisiones sin base muestran estado comprensible.
- [ ] Revisar definiciones con ejemplos y commit. Búsqueda full-text conserva configuración English hasta una decisión de idioma propia, documentando su alcance.

### B5.2 — Pasada UX y accesibilidad por apartados

**Responsable:** frontend/producto. **Depende:** B2/B3.4. **Archivos:** Layout, chrome/*, ui/*, componentes auth/tickets/articles/admin y styles/tailwind.css.

- [ ] Recorrer por rol login, recuperación, perfil/MFA, dashboard, lista/detalle/alta tickets, adjuntos, conocimiento, usuarios, SLA, reglas, analíticas y offline en escritorio/móvil. Registrar loading, vacío, error, permisos y éxito de cada vista.
- [ ] Corregir foco y teclado de menús/diálogos, labels, anuncios de errores, contraste, reduced-motion y tablas móviles. Dar al banner SW el mismo acabado de los componentes existentes y evitar anunciar una actualización en cada instalación inicial.
- [ ] Ejecutar teclado completo y comprobaciones de accesibilidad de las rutas críticas; overlays de menú se cierran y devuelven foco. Validar a 390×844 y 1440×900, sin screenshots como único criterio.
- [ ] Revisar cambios visuales acotados y commit. Mantener copy consistente en inglés; registrar traducción completa como decisión posterior.

### B5.3 — Límites de arquitectura y piezas esbozadas

**Responsable:** backend/frontend. **Depende:** flujos estabilizados y B0.2. **Archivos:** tenant_resolver.go, cached_tenant_resolver.go, config/features.go, storage/*, domain/event/*, worker/job.go, WIRING, swagger.go y anotaciones de handlers.

- [ ] Revisar impacto/callers de cada candidato y clasificarlo: usado, necesario para alcance actual o backlog. Confirmar con fuente y tests; el grafo vacío no basta para borrar.
- [ ] Retirar el import de interfaz HTTP en el resolver PostgreSQL mediante contrato estructural o puerto existente. Eliminar stubs/flags/cola/bus sin consumidor del producto si el análisis confirma que solo aparentan capacidades; no conectarlos solo para justificar que existen. Conservar utilidades usadas por las entregas anteriores.
- [ ] Generar y registrar OpenAPI de rutas soportadas; comprobar `/swagger/doc.json` y auth/errores contra el router. Actualizar comentarios TODO que ya están resueltos y mantener excepciones SQL explícitas.
- [ ] Ejecutar build/tests, impacto y revisión; commit de limpieza acotada, con backlog S3/SSO/superadmin/push/multirréplica documentado.

### B5.4 — Roadmap, guías y cierre de auditoría

**Responsable:** mantenimiento/producto. **Depende:** todos los cierres anteriores. **Archivos:** docs/01–09, README, CLAUDE, AGENTS y nuevos `F:/Proyectos/Trickreport/docs/operations/{deployment,recovery,release}.md`, `F:/Proyectos/Trickreport/docs/architecture/decisions.md`.

- [ ] Actualizar stack, rutas, migraciones y comandos reales. Marcar cada capacidad como implementada/conectada/probada/operable con enlace a evidencia; sustituir casillas del roadmap que hoy no representan el estado.
- [ ] Publicar guías por rol y operación: alta, asignación, notas internas, publicación, SLA/reglas, MFA/reset, offline, recuperación y release. Documentar supuestos propuestos que finalmente se aceptaron y capacidades fuera de alcance.
- [ ] Reejecutar auditoría de cierre desde checkout limpio: test/race/integración/E2E/dep scan, arranque y restore. Cada hallazgo tiene commit, prueba, evidencia y estado cerrado o pendiente justificado; no usar una nueva puntuación subjetiva como sustituto.
- [ ] Revisar matriz completa, refrescar índices/grafo y guardar documento final de cierre. Esta tarea no se marca completa mientras exista un bloqueo de acceso o datos sin resolver.

## Matriz de trazabilidad

| Auditoría / apartado | Tareas que lo cierran |
|---|---|
| AUTH-01 MFA | B1.2, B1.6, B3.4 |
| AUTH-02 refresh/logout/sesiones | B1.1, B1.6, B2.3, B3.3 |
| SEC-01 serialización Astro | B1.3, B1.7, B3.4 |
| DATA-01 privacidad tiempo real y RT-01 handshake | B1.4, B4.4 |
| DATA-02 PWA y cruces de cuenta | B1.5, B2.7, B3.4 |
| DEPLOY-01 entorno sin .env | B2.1, B4.1 |
| FLOW-01 CSRF SSR | B2.2, B3.4 |
| AUTO-01 condiciones/invariantes | B2.5, B2.8 |
| CI-01 pipeline y checks ausentes | B3.1–B3.5 |
| DEP-01 dependencias/toolchain | B0.3, B1.7, B4.1 |
| DEPLOY-02 topologías | B4.1, B4.6 |
| OPS-01 monitoring/restore | B4.2–B4.5 |
| UX-01 redirects | B2.3, B3.4 |
| Usuarios/perfiles/roles | B1.1, B1.6, B5.2 |
| Tickets, filtros, concurrencia, asignación | B2.6, B2.7, B5.2 |
| Conocimiento | B2.10, B3.4, B5.2 |
| Adjuntos | B2.9, B4.2 |
| SLA/worker | B2.8, B4.4, B5.1 |
| Email/reset | B2.4, B4.3 |
| Analíticas/reportes/exportación | B5.1 |
| Arquitectura/componentes no conectados/Swagger | B5.3 |
| Ramas, commits, release/protecciones | B0.1, B3.5, B4.6 |
| Docs/roadmap/idioma | B5.2–B5.4 |
| Índices y confianza del análisis | B0.2, verificación por tarea |

## Definición de terminado

Una tarea queda terminada cuando su prueba del comportamiento corregido pasa, se actualizaron consumidores/mocks/fixtures/contratos, el análisis de impacto y cambios no está incompleto, el diff es revisable y la documentación afectada coincide. No basta con añadir un archivo o conseguir un build.

Una entrega se puede promover cuando no quedan regresiones respecto al baseline y sus dependencias están cerradas. La entrega candidata requiere: identidad/tenant correctos; contratos e integración críticos sin skips; suite completa `-race`; frontend lint/check/tests/build; E2E por rol; avisos de runtime altos/críticos aplicables cerrados; imágenes reproducibles; restore y rollback ensayados.

No hay objetivo arbitrario de cobertura global. Se exige cobertura de las rutas y fallos concretos de esta matriz; el perfil global satisfactorio se conserva para detectar retrocesos.

## Primera entrega recomendada

Ejecutar B0.1–B0.3, B2.1, B3.1 y B3.3, y después B1.1–B1.3. Así se preserva el trabajo local, se consigue un entorno de comprobación confiable y se cierran primero los límites más graves de autenticación y serialización. La entrega siguiente agrupa privacidad offline y tiempo real con sus pruebas de dos identidades.

## Fuentes y verificación de este plan

- [Node.js: líneas soportadas y recomendación LTS](https://nodejs.org/en/about/previous-releases) y [checksums oficiales Node 24.22.1](https://nodejs.org/download/release/latest-v24.x/SHASUMS256.txt).
- [Go: versiones estables publicadas, incluida 1.26.10](https://go.dev/dl/?mode=json).
- [Guía oficial de actualización de Astro](https://docs.astro.build/en/upgrade-astro/) y [migración a Astro 7](https://docs.astro.build/en/guides/upgrade-to/v7/).
- Codebase Memory volvió a estar disponible: proyecto `Trickreport`, raíz `F:/Proyectos/Trickreport`, generación `2026-10-22T02:36:21Z` (21 de octubre por la noche en Caracas). La consulta de Authenticate devolvió router, contratos y tests como callers. El modo fast excluye documentación, SQL, public/scripts y algunos tests; check_index_coverage señaló esas exclusiones, parseo parcial de Dockerfile y metadatos cambiados. Se corroboró con fuente; el índice no acredita cobertura total.
- Graphify mantiene el corpus previo a este plan: 2.189 nodos, 5.925 aristas, 131 comunidades. La consulta documental acotada confirmó divergencias de diseño; no se atribuyen estos planes futuros al grafo existente.
- El HEAD y el cambio local de ArticleEditor seguían iguales al redactar. No se han repetido pruebas de producto, ejecutado migraciones o aplicado mejoras para escribir el plan. La auditoría original permanece como snapshot histórico.

**Autorrevisión:** 35 tareas, seis bloques, todos los hallazgos identificados y apartados funcionales trazados; dependencias entre bloques explícitas; cinco escenarios de Review Focus asignados a pruebas. Las políticas de producto de la tabla inicial son propuestas de este documento, no decisiones silenciosas ejecutadas.
