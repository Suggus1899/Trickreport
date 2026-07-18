# Frontend — Next.js + PWA

## Stack

| Propósito | Tecnología |
|-----------|-----------|
| Framework | Next.js 14+ (App Router) |
| UI | Tailwind CSS + shadcn/ui |
| Estado global | Zustand |
| Formularios | React Hook Form + Zod |
| PWA | next-pwa |
| WebSocket | `socket.io` o native WebSocket |

## Estructura de carpetas (sugerida)

```
src/
├── app/                    # App Router pages
│   ├── (auth)/            # Login, logout
│   ├── (dashboard)/       # Dashboard home
│   ├── tickets/           # CRUD tickets
│   ├── articles/          # Base de conocimientos
│   ├── admin/             # Admin panel
│   └── layout.tsx
├── components/
│   ├── ui/                # shadcn/ui components
│   ├── tickets/           # TicketCard, TicketForm, etc.
│   ├── dashboard/         # Charts, SummaryCards
│   └── layout/            # Sidebar, Header, etc.
├── lib/
│   ├── api/               # Cliente HTTP para el backend Go
│   ├── auth/              # Context de autenticación
│   └── utils.ts
├── hooks/                 # Custom hooks
└── types/                 # TypeScript types
```

## Vistas principales

Ver `08-views.md` para detalle de cada pantalla.

## PWA

- Service worker con estrategia cache-first para assets estáticos
- Manifest con iconos en varios tamaños
- Soporte offline parcial (lectura de tickets cacheados)
- Push notifications via Web Push API + backend
