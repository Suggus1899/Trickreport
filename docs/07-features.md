# Funcionalidades Detalladas

## 1. CRUD de Tickets + Estados/Prioridades

- Creación con título, descripción, categoría, prioridad
- Estados: open → in_progress → waiting_client → resolved → closed
- Campos calculados: tiempo en estado actual, tiempo total abierto
- Historial de cambios auditable

## 2. SLA y Escalamientos

- Políticas configurables por prioridad (tiempo de respuesta y resolución)
- Colores de alerta: verde (OK), amarillo (50% tiempo), rojo (cerca del breach), gris (breach)
- Escalamiento automático: si un agente no responde en X minutos, se reasigna al superior
- Reporte de cumplimiento de SLA mensual

## 3. Notificaciones

| Tipo | Canal |
|------|-------|
| Ticket asignado | In-app + email |
| Nuevo comentario | In-app (WebSocket) + email |
| SLA próximo a vencer | In-app + email |
| Breach de SLA | In-app + email + escalamiento |
| Ticket resuelto | In-app + email |

- Preferencias de notificación por usuario
- Notificaciones push vía Web Push (PWA) cuando la app está en background

## 4. Base de Conocimientos

- Artículos categorizados con tags
- Búsqueda por título, contenido y tags
- Agentes y admin pueden crear/editar
- Usuarios finales solo lectura
- Opción de sugerir artículos relacionados al crear un ticket

## 5. Dashboard y Reportes

### Dashboard de agente
- Tickets asignados pendientes
- Próximos SLA a vencer
- Actividad reciente

### Dashboard de admin
- Todas las métricas de agente
- Tickets abiertos vs cerrados (gráfico)
- Cumplimiento SLA por agente/equipo
- Tiempo promedio de resolución
- Tickets por categoría
- Usuarios más activos

### Exportación
- CSV/PDF de reportes
- Tickets cerrados en un período

## 6. Adjuntos / Subida de Archivos

- Usuarios pueden adjuntar archivos a tickets (capturas, logs, PDFs)
- Límite de tamaño configurable por tenant
- Almacenamiento en disco local o S3 (configurable)
- Vista previa de imágenes en el detalle del ticket
- Historial de archivos (quién subió, cuándo)

## 7. Reglas de Automatización

- Admin crea reglas IF-THEN sin código
- **Triggers**: al crear ticket, al cambiar estado, al cambiar prioridad
- **Condiciones**: categoría, prioridad, texto contiene, email del creador
- **Acciones**: asignar a agente, cambiar prioridad, cambiar categoría, notificar, añadir tag
- Evaluación en tiempo real (sin polling)

## 8. Roles y Permisos

Ver `06-auth.md` — RBAC detallado.

## 9. Multi-tenant

- Cada tenant es una organización independiente
- Datos completamente aislados por `tenant_id`
- Cada tenant tiene su propia configuración (SLA, categorías, branding)
- Login aísla usuarios por tenant
- Admin global (superadmin) puede gestionar tenants
- URL por subdominio: `tenant1.trickreport.app`

## 10. Tags y Categorías

- Categorías configurables por admin (hardware, software, red, etc.)
- Tags libres en tickets y artículos
- Filtrado combinado por categoría + tags
