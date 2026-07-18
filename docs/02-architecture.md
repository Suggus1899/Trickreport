# Arquitectura del Sistema

```
┌─────────────────────────────────────────────────┐
│                   Cliente Web                    │
│            Next.js (React) + PWA                 │
└────────────────────┬────────────────────────────┘
                     │ HTTPS / REST + WebSocket
                     ▼
┌──────────────────────────────────────────────────┐
│            API Gateway (Go)                       │
│  ┌───────────┐ ┌──────────┐ ┌────────────────┐   │
│  │ REST API  │ │ WebSocket│ │ Auth Middleware │   │
│  └───────────┘ └──────────┘ └────────────────┘   │
└────────────────────┬──────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
┌──────────────┐ ┌────────┐ ┌──────────────┐
│ PostgreSQL   │ │ Redis  │ │  LDAP / SSO  │
│ (datos)      │ │ (cache)│ │  (auth)      │
└──────────────┘ └────────┘ └──────────────┘
```

## Principios

- **API-first** — Backend Go expone API REST, frontend consume
- **WebSocket** para notificaciones en tiempo real
- **Modular** — Separación clara por dominios (tickets, usuarios, conocimientos, SLA)
- **PWA** — El frontend Next.js funciona como PWA para experiencia móvil sin app store

## Capas del backend (Go)

1. **Handlers** — Entrada HTTP, validación básica
2. **Services** — Lógica de negocio
3. **Repositories** — Acceso a datos (PostgreSQL)
4. **Middleware** — Auth, logging, rate limiting
