# Vistas y Flujos de Pantalla

## Web (Next.js)

### Layout principal
- Sidebar con navegación (Dashboard, Tickets, Conocimientos, Admin)
- Header con avatar, notificaciones, búsqueda global
- Contenido central

### 1. Login (`/login`)
- Input usuario/contraseña + botón "Login con SSO"
- Recuperación de contraseña (si no es SSO puro)

### 2. Dashboard (`/dashboard`)
- Cards resumen: tickets abiertos, en progreso, vencidos hoy
- Timeline de actividad reciente
- Gráfico de tickets por día/semana
- Lista de tickets vencidos próximos (agente/admin)

### 3. Lista de Tickets (`/tickets`)
- Tabla con columnas: ID, título, prioridad, estado, asignado, SLA
- Filtros combinados: estado, prioridad, categoría, agente, fecha
- Búsqueda por texto
- Botón "Nuevo ticket"

### 4. Detalle de Ticket (`/tickets/:id`)
- Información completa del ticket
- Historial de cambios cronológico
- Comentarios (con toggle interno/solo agente para agentes)
- Botones de acción: cambiar estado, asignar, escalar
- Sección de adjuntos con upload drag-and-drop y vista previa
- Indicador de SLA con barra de progreso

### 5. Nuevo Ticket (`/tickets/new`)
- Formulario: título, descripción, categoría, prioridad
- Sugerencias de artículos de KB mientras escribe

### 6. Base de Conocimientos (`/articles`)
- Grid/tarjetas de artículos
- Búsqueda y filtros por categoría/tags
- Vista de artículo individual (`/articles/:id`)

### 7. Admin — SLA Policies (`/admin/sla`)
- Tabla de políticas por prioridad
- Formulario para editar tiempos

### 8. Admin — Usuarios (`/admin/users`)
- Tabla de usuarios con roles
- Crear/editar/desactivar usuarios

### 9. Admin — Automatizaciones (`/admin/automation`)
- Lista de reglas con toggle on/off
- Formulario: disparador + condiciones + acciones
- Vista previa de a qué tickets afectaría

### 10. Admin — Tenants (`/admin/tenants`) (solo superadmin)
- Lista de organizaciones
- Crear/editar/desactivar tenant
- Configuración por tenant (SLA defaults, branding, storage)

### 11. Admin — Reportes (`/admin/reports`)
- Selector de fechas
- Gráficos de métricas globales
- Botón exportar CSV/PDF

### 10. Perfil (`/profile`)
- Datos del usuario
- Preferencias de notificación
- Cambiar avatar

## PWA (móvil)

Las mismas vistas se renderizan responsive con Tailwind CSS.
El manifest y service worker permiten:

- Instalación en home screen
- Notificaciones push
- Cache offline para artículos de KB y tickets recientes
- Navegación tipo app (sin barra de navegador)
