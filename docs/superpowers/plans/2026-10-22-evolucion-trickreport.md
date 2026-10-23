# Trickreport: ejecución del plan maestro

Fecha: 2026-10-22 (America/Caracas). Estado: en ejecución; ninguna casilla representa una capacidad terminada sin evidencia.

## Objetivo y decisiones aprobadas

Tres productos bajo la marca Trickreport: soporte TI, MSP y atención al cliente. Un repositorio, código común reutilizable, tres aplicaciones/API/bases de datos/cuentas/versiones independientes. TI se entrega primero. Cada producto ofrece instalación propia y SaaS. Español principal e inglés. Catálogo completo de 200 casos por producto (600), entrega funcional por etapas.

Frontend: Next.js App Router, TypeScript, React 19, Tailwind CSS 4 y HeroUI 3; npm workspaces. Backend: conservar Go/chi/pgx y PostgreSQL 16. Migrar recorridos completos desde Astro, preservando contratos de seguridad. No reescribir servicios útiles.

Marca profesional: índigo #4F46E5 y verde #047857; oscuro #0B1220 con acentos #818CF8 y #34D399. Public Sans y Space Mono. Logo con T geométrica y seguimiento de solicitudes; SVG, PNG, variantes y favicon.

Integraciones: SMTP/IMAP, LDAP existente, Zabbix, MeshCentral y API oficial de WhatsApp. Administración comercial manual al inicio. Sin agente remoto propio, aplicaciones móviles nativas, telefonía, canales sociales adicionales ni cobros automáticos.

Trabajo sobre main, sin respaldos ni ramas auxiliares. Preservar cambios previos y .agents/. No reset, force-push ni publicación de historial divergente sin reconciliarlo. Desarrollo local nativo; contenedores para despliegue. Pruebas de datos únicamente en bases desechables.

## Entregas y dependencias

- [ ] E0: cerrar los 35 puntos de los seis bloques de la auditoría anterior, con estado por tarea y evidencia actual.
- [ ] E1: catálogo de 600 casos, permisos, ejemplos, aceptación, prioridades, navegación e identidad visual.
- [ ] E2: soporte TI en Next con sesiones, MFA, recuperación, perfil, dashboard, tickets/comentarios/adjuntos, conocimiento y administración; instalación sencilla y alta disponibilidad antes del primer SaaS.
- [ ] E3: catálogo de servicios, aprobaciones, activos/CMDB, problemas y cambios.
- [ ] E4: MSP: clientes, contratos/SLA, técnicos, horas/gastos, proyectos, rentabilidad, inventario, descubrimiento, monitorización y remoto integrado.
- [ ] E5: atención: portal, correo, chat, WhatsApp, asignación, historial, autoservicio, SLA y satisfacción.
- [ ] E6: cobertura comprobada del catálogo, documentación bilingüe, carga, recuperación y versiones reproducibles.

E0 conserva los bloques: base confiable; identidad/privacidad; flujos esenciales; calidad/CI; operación/entrega; producto/documentación. Los entregables preparatorios independientes pueden avanzar juntos; ninguna publicación salta la validación de la base.

## Interfaces y aislamiento

- Mantener /api/v1 en TI; OpenAPI por producto y tipos frontend derivados.
- Tokens, secretos y cookies independientes por producto; autorización y tenant comprobados en Go.
- En MSP, tenant es el proveedor y client es su cliente; contactos y técnicos se limitan a clientes autorizados.
- Persistir sesiones, outbox, notificaciones y trabajos; idempotencia ante reintentos y webhooks duplicados.
- WebSocket filtra destinatarios; desconexión y revocación no pueden recuperar acceso.
- Borradores offline por producto/organización/usuario; borradores antiguos sin dueño siguen en cuarentena.
- Páginas privadas sin caché compartida; imágenes de Next idénticas por versión y recursos previos durante transición.

## Operación objetivo

Instalación sencilla: Compose + Caddy + Next + Go + PostgreSQL. SaaS: tres dominios de fallo, nodo de aplicación y nodo de datos por dominio (seis VPS/servidores); Swarm con tres managers, dos réplicas por aplicación, HAProxy/Keepalived con red e IP virtual compatibles; PostgreSQL/Patroni/etcd con replicación síncrona estricta. Cada producto tiene base/usuario/secretos propios aunque la infraestructura sea compartida. Equipos definitivos dependen de la carga medida.

Objetivos de aceptación (no capacidad ya acreditada): 15.000 cuentas totales, 1.500 sesiones concurrentes, mezcla 80% lecturas/20% escrituras, 150 operaciones de negocio/s durante 45 minutos tras calentamiento. p95 lecturas ordinarias <500ms y escrituras <1s; errores inesperados <1%; sin pérdida de escrituras confirmadas ni duplicados; recuperación automática ante un nodo caído <=60s. Validar carga repartida y concentrada por producto, WebSocket y despliegue gradual.

## Evidencia y terminado

Backend: build, vet, lint, race y PostgreSQL 16. Frontend: tipos, lint, unitarias, build y E2E con API real. Seguridad: refresh concurrente, MFA, revocación, CSRF, permisos/tenant/cliente/producto, cache/offline. Migraciones desde vacío y versión previa. Fallos de SMTP, repeticiones/desorden de webhooks, reinicio/failover. Accesibilidad, móvil/escritorio, ES/EN, claro/oscuro. Guías verificadas desde instalación limpia.

Un caso sólo está operable cuando código, API, pruebas, documentación y operación están conectados. Los 600 casos no se marcan implementados por existir su catálogo. Documentación de arquitectura, API, instalación, uso por rol, integración, actualización, recuperación y troubleshooting acompaña cada entrega.

## Registro de ejecución

El registro y los resultados completos se conservan durante la ejecución en `.superpowers/sdd/2026-10-22-evolucion/`. La especificación vinculante es el plan aprobado en esta conversación; este documento conserva sus decisiones para continuidad.
