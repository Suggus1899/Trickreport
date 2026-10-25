# Trickreport: cierre de producción

Fecha: 2026-10-25 (America/Caracas). Estado: en ejecución. Complementa `2026-10-22-evolucion-trickreport.md`; para despliegue SaaS sustituye la sección "Operación objetivo" de ese documento. Ninguna casilla representa una capacidad terminada sin evidencia.

## Decisiones aprobadas (2026-10-25)

- Trabajo directo en `F:\Proyectos\Trickreport-remoto`, rama `main`, sin ramas auxiliares. Commits en Conventional Commits, sin atribución de IA ni `Co-Authored-By`.
- La worktree Delta antigua (`hdcgq4qnxfqz`, base `1fc2126`) no se toca: tiene 172 entradas sin confirmar, divergentes de `origin/main`.
- Sin Docker ni instalaciones nuevas en esta máquina (no Node 24, no navegadores nuevos, no CLIs, no toolchains). Lo que exige Docker o Node 24 queda en CI o marcado como no verificado.
- No hay push a `origin/main` sin autorización específica. Los 7 jobs de CI en el mismo commit requieren publicar.
- Despliegue SaaS: Render (API Go, worker, PostgreSQL gestionado) y Vercel (aplicaciones Next). Instalación propia: Compose + Caddy, sin cambios de diseño.
- Alcance: 600 casos (200 TI, 200 MSP, 200 CX) en `docs/use-cases/*.json`, todos `planned`. Es un programa de varias sesiones; cada caso se marca implementado solo con código, contrato, pruebas, documentación y operación.

## Capacidades verificadas (documentación pública, 2026-10-25)

- **Render Postgres HA:** standby en otra zona de la misma región; replicación **asíncrona** (retraso de segundos); conmutación automática tras ~30 s sin respuesta del primario, en segundos. Requiere al menos 1 CPU y PostgreSQL 13 o superior. El standby se factura igual que el primario.
- Consecuencia: no cumple por sí solo "sin pérdida de escrituras confirmadas" (ver decisión pendiente).
- **Render web services:** despliegue sin caída salvo con disco persistente; health check HTTP (2xx/3xx = sano); varias instancias en nodos distintos; el servicio debe escuchar en `0.0.0.0`.
- Las conexiones WebSocket activas se cortan al terminar la instancia anterior: el cliente debe reconectar.
- **Vercel Functions:** límite de cuerpo de petición de 4,5 MB (413 `FUNCTION_PAYLOAD_TOO_LARGE`). Los webhooks no deben pasar por funciones Next: van a la API Go.
- No verificado: si las rewrites hacia un origen externo aplican el mismo límite. Prueba pendiente antes de depender de ello.

## Decisión pendiente (usuario)

- **RPO del SaaS.** Con replicación asíncrona, una conmutación puede perder escrituras confirmadas de los últimos segundos.
- Opción A: aceptar ese RPO para el primer SaaS, documentarlo en la oferta y probar la recuperación (recomendada para empezar).
- Opción B: réplica síncrona con otro proveedor o infraestructura propia (más operación).
- Hasta decidir, el criterio "sin pérdida de escrituras confirmadas" queda marcado como no cumplido en SaaS.

## Estado de partida (verificado en lectura)

- `main` = `origin/main` = `9647012`, árbol limpio.
- Backend Go (`go 1.26.10`), módulo `github.com/trickreport/backend`; migraciones hasta `000016_cx`.
- CI: 7 jobs en `.github/workflows/ci.yml` (backend, lint-backend, frontend, security-backend, security-frontend, operations, container-scan).
- Hallazgos P1/P2 de la auditoría: ver E0.
- No hay índice GitNexus en este clon; el impacto se calcula con búsqueda y lectura de llamadores, y se anota en cada corrección.

## Entregas

- [ ] **E0 — base confiable:** corregir los hallazgos P1/P2 con regresión primero; baseline local documentado; 7 jobs de CI verdes en un mismo commit (requiere push autorizado).
- [ ] **E1 — aislamiento de productos:** composición Go por producto sin pool compartido; bases, usuarios PG, secretos, audiencias JWT y cookies separados; matriz de permisos; autorización antes de paginar o exportar; outbox y LISTEN/NOTIFY para WebSocket; dos réplicas.
- [ ] **E2 — ITSM TI en Next:** sesiones, MFA, recuperación, perfil, dashboard, tickets, adjuntos, conocimiento, administración, notificaciones y borradores. Retirada de Astro con cuarentena de borradores sin propietario.
- [ ] **E3 — ITSM TI avanzado:** catálogo de servicios, aprobaciones, CMDB, problemas y cambios.
- [ ] **E4 — MSP:** clientes, contratos y SLA, técnicos, horas y gastos, proyectos, rentabilidad, inventario, descubrimiento, monitorización (Zabbix) y remoto (MeshCentral).
- [ ] **E5 — atención al cliente:** portal, correo (IMAP y SMTP), chat, WhatsApp con API oficial, asignación, historial, autoservicio, SLA y satisfacción.
- [ ] **E6 — instalación, HA y SaaS:** Compose para instalación propia; SaaS en Render y Vercel; release reproducible; recuperación ensayada; cargas medidas; cobertura del catálogo comprobada.

## E0 — detalle

Cada corrección lleva primero su prueba de regresión que falla, luego el cambio y un commit propio.

| ID | Hallazgo | Corrección | Prueba |
| --- | --- | --- | --- |
| E0-a | Webhooks: `Caddyfile` envía `/api/*` a Next; `packages/workspace/src/proxy.ts` responde 403 a cualquier petición sin `Origin`; la lista de cabeceras no incluye `X-MSP-Webhook-Token` ni `X-Hub-Signature-256`. | Frontera propia para webhooks: sin exigir `Origin`, cuerpo y cabeceras intactos, directo a Go. La protección de origen del resto de `/api` no cambia. | Vitest sobre `proxy.ts`; `graphify-out/verify-webhook-boundary.mjs` con Node 25. |
| E0-b | MSP: `GET /{id}/tasks` y `RequestRemoteSession` solo comprueban tenant, no cliente. | Un único helper de autorización por cliente para proyectos, solicitudes remotas, tickets, activos, búsquedas y exportaciones. | Prueba Go de integración en PostgreSQL desechable con dos clientes del mismo proveedor. |
| E0-c | Contratos: el formulario envía `currency: 'EUR'` y `sla_overrides: {}` que no muestra; el UPDATE sobrescribe. Fechas: `slice(0,16)` más `toISOString()` desplaza `12:00Z` a `16:00Z`. | No sobrescribir campos no mostrados (cargar valores existentes o combinar en servidor); una sola conversión local/UTC. | Vitest de conversión de fechas; prueba Go de actualización parcial. |
| E0-d | SLA: inserts de tickets MSP y TI sin plazos; el worker no detecta vencimientos. | Un único cálculo de SLA para todos los orígenes; creación atómica de ticket, plazos e historial. | Prueba Go de creación atómica y de vencimiento en el worker. |
| E0-e | Accesibilidad: `ui.tsx` (Select) asocia el label con el texto de las opciones; el sidebar fijo oculta el cierre de sesión. | Asociación label/id o `aria-labelledby`; zona de navegación desplazable. | Pruebas E2E en `apps/it/tests/integration/service.spec.ts`. |
| E0-f | CX: el detalle no se recarga tras cambiar estado, tomar o escalar; los contactos no vinculan `user_id`. | Recargar el detalle tras cada acción; permitir vincular contacto con usuario. | E2E y vitest. |
| E0-g | Calidad: 132 incidencias de lint backend; PostCSS con alerta alta en las 3 apps; 5 altas en imagen; cobertura 0 % en `itsm`, `msp` y `cx` (paquetes de aplicación). | Corregir por grupos; actualizar dependencias sin cambiar contratos. | `golangci-lint run ./...` v2.15.0; `npm audit`; pruebas de aplicación nuevas. |

Documentación a actualizar con E0: `README.md` (migraciones hasta 000016; MSP y CX ya no son shells), `AGENTS.md` y `CLAUDE.md` (Go 1.26; Astro se retira en E2; migraciones hasta 000016).

## Verificación

- **Local, sin Docker:** `go build`, `go vet`, `gofmt -l`, `go test ./...` (con `-race` si hay CGO y compilador C); PostgreSQL 18.4 desechable (CI usa 16.15: la diferencia queda anotada); `npm run lint`, `npm run check`, `npm test`, `npm run build`; E2E con el navegador de Playwright ya presente si la configuración lo permite.
- **Solo CI:** validación de Compose, promtool v3.15.1, Trivy v0.76.1 (imágenes), `scripts/restore.test.sh` y govulncheck con la versión fijada.
- Resultados y limitaciones en `docs/operations/baseline.md` (addendum 2026-10-25, sin archivo nuevo).

## Reglas

- Español principal, inglés completo. Tema claro y oscuro. Marca: índigo `#4F46E5` y `#818CF8`; verde `#047857` y `#34D399`; fondo `#0B1220`; tipografías Public Sans y Space Mono.
- Administración comercial manual al inicio. Sin agente remoto propio, apps móviles nativas, telefonía, canales sociales adicionales ni cobros automáticos.
- Migraciones aditivas, aplicadas con `--migrate-up` y `--migrate-down`, no con scripts SQL sueltos.
- Pruebas con datos solo en bases desechables. Nunca contra los datos de `/d/DB/data`.
- Sin reset, sin force-push, sin publicar historia divergente sin reconciliarla.
