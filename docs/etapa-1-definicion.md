# Registro de la Etapa 1 — Definir producto y diseño

Estado: **completada**. Este documento deja constancia de lo definido en la Etapa 1 del plan maestro y de lo que pasa a la Etapa 2.

## Definido en esta etapa

### Catálogo de casos de uso (600)

- 200 por producto en [`docs/use-cases/`](../use-cases/README.md) (`ti.json`/`msp.json`/`cx.json` con su lectura en Markdown), con identificadores estables `TI-001`…`TI-200`, `MSP-001`…`MSP-200` y `CX-001`…`CX-200`.
- Reparto por áreas **exacto** al del plan maestro. Verificado por conteo: TI 16/32/22/14/20/24/12/28/16/16; MSP 18/30/22/16/20/14/16/20/28/16; CX 20/32/22/26/24/20/14/16/12/14.
- Cada caso incluye objetivo, actor, permisos, condiciones previas, flujo principal, alternativas y fallos, ejemplo con datos concretos y resultado esperado, pantallas, entidades, contrato de API (marcado `proposed`), criterios de aceptación, prioridad, etapa de entrega (`phase`: TI 2–3, MSP 4, CX 5) y estado (`planned`).
- En esta etapa se corrigieron MSP y CX, que tenían un reparto plano de 10×20 con áreas distintas a las del plan; TI ya era conforme y no cambió. El generador es `docs/use-cases/.author.py`.

### Arquitectura y roles

- [`docs/architecture/`](../architecture/README.md): arquitectura por producto, reglas de aislamiento de datos entre productos y ámbitos, reglas de contratos (`/api/v1`, OpenAPI por producto, autorización en el servidor, WebSocket, webhooks, borradores offline), las dos topologías de despliegue (Compose sencilla y SaaS HA con Swarm + Patroni) y el modelo de roles por producto que usa el catálogo.

### Marca e identidad

- [`docs/brand/`](../brand/README.md): símbolo «T geométrica» en versiones color, clara, oscura y monocroma; logotipo horizontal para fondos claros y oscuros; favicon (también en `apps/*/public/`); guía de uso.
- Tokens de diseño en `docs/brand/tokens.md` y `docs/brand/tokens.css`, alineados con los nombres que ya usa `packages/workspace` (`--accent`, `--green`, `--background`, `--foreground` y tokens de estado).

## Condición de salida de la etapa

- Cada caso tiene resultado esperado y criterios de aceptación: **cumplido**.
- Cada módulo tiene una etapa asignada: **cumplido** (campo `phase` por caso y documento de arquitectura por producto).

## Pasa a la Etapa 2

- Migrar los recorridos de Soporte TI a Next.js (acceso, sesiones, dashboard, tickets, artículos, perfil, administración, SLA y automatizaciones) contra la API real, por recorridos completos.
- Definir el OpenAPI de Soporte TI al implementar cada caso y generar desde ahí los tipos del frontend; los contratos del catálogo están marcados `proposed` hasta entonces.
- Preparar las dos modalidades de instalación de Soporte TI; la alta disponibilidad debe estar aprobada antes de publicar el SaaS.
- Generar las exportaciones PNG de marca desde los SVG si un canal las requiere (procedimiento en la guía de marca).
- Versiones en inglés de la interfaz y las guías, conservando los identificadores del catálogo.
