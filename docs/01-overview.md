# Trickreport — Visión General

Sistema de soporte IT para gestión de tickets e incidencias.

## Stack Tecnológico

| Capa | Tecnología |
|------|-----------|
| Backend | Go |
| Frontend Web | Next.js (React) |
| Móvil | PWA (Progressive Web App) |
| Base de datos | PostgreSQL |
| Autenticación | SSO / LDAP |

## Objetivos

- Gestión eficiente de tickets e incidencias de soporte IT
- Cumplimiento de SLA con escalamientos automáticos
- Base de conocimientos para reducir incidencias recurrentes
- Dashboard con métricas y reportes para toma de decisiones
- Acceso desde web y móvil (PWA) para agentes y usuarios finales
- Roles bien definidos: admin, agente, usuario final
- Multi-tenant: soporte para múltiples organizaciones
- Reglas de automatización para asignación y notificaciones
- Adjuntos en tickets (capturas, logs, documentos)

## Usuarios del sistema

- **Usuario final** — Crea tickets, consulta estado, accede a base de conocimientos
- **Agente de soporte** — Gestiona tickets, responde, escala, escribe artículos
- **Admin** — Configura SLA, roles, usuarios, reportes globales
