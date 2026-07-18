# Plan de Implementación

## Fase 1 — Fundación (Sprint 1-2)

- [ ] Setup del proyecto Go (estructura, router, config)
- [ ] Setup del proyecto Next.js + Tailwind + shadcn/ui
- [ ] Conexión a PostgreSQL + migraciones iniciales
- [ ] Modelo de datos (migraciones)
- [ ] Auth: login LDAP + JWT + middleware
- [ ] CRUD básico de usuarios (admin)
- [ ] Despliegue básico (Docker + docker-compose)

## Fase 2 — Tickets (Sprint 3-4)

- [ ] CRUD de tickets
- [ ] Comentarios + historial
- [ ] Cambio de estados
- [ ] Asignación de tickets
- [ ] Frontend: lista, detalle, crear ticket
- [ ] Filtros y búsqueda

## Fase 3 — SLA + Notificaciones + Automatizaciones (Sprint 5-6)

- [ ] Políticas de SLA
- [ ] Scheduler de escalamiento
- [ ] Notificaciones in-app + WebSocket
- [ ] Notificaciones email
- [ ] Notificaciones push (PWA)
- [ ] Frontend: indicadores SLA en tickets
- [ ] Preferencias de notificación
- [ ] Motor de automatización (reglas IF-THEN)
- [ ] CRUD de reglas + UI

## Fase 4 — Base de Conocimientos + Adjuntos (Sprint 7)

- [ ] CRUD de artículos
- [ ] Búsqueda y tags
- [ ] Sugerencias al crear ticket
- [ ] Frontend: lista y vista de artículos
- [ ] Subida de adjuntos en tickets
- [ ] Almacenamiento local/S3 + validación de tipos/tamaños

## Fase 5 — Dashboard + Reportes (Sprint 8)

- [ ] Dashboard de agente
- [ ] Dashboard de admin
- [ ] Gráficos y métricas
- [ ] Exportación CSV/PDF
- [ ] Frontend: vistas de dashboard

## Fase 6 — Multi-tenant + PWA + Mejoras (Sprint 9-10)

- [ ] Service worker y manifest
- [ ] Cache offline parcial
- [ ] Testing end-to-end
- [ ] Performance optimization
- [ ] Aislamiento multi-tenant por tenant_id
- [ ] Superadmin: gestión de tenants
- [ ] Documentación de usuario y admin
- [ ] Carga de datos de prueba

## Fase 7 — Producción (Sprint 11-12)

- [ ] CI/CD pipeline
- [ ] Monitoreo y logging
- [ ] Seguridad (audit, rate limiting, CORS)
- [ ] Backup y restore
- [ ] Deploy a producción
