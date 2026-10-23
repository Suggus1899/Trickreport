# Arquitectura de la familia Trickreport

Documentos de definición de la Etapa 1 del plan maestro. Describen la arquitectura objetivo de los tres productos; la arquitectura del sistema actual está en [`docs/02-architecture.md`](../02-architecture.md) y el catálogo funcional en [`docs/use-cases/`](../use-cases/README.md).

| Documento | Contenido |
| --- | --- |
| [productos.md](productos.md) | Arquitectura por producto: aplicaciones, módulos del backend y paquetes compartidos. |
| [aislamiento-y-contratos.md](aislamiento-y-contratos.md) | Límites de datos entre productos y ámbitos, contratos `/api/v1`, OpenAPI por producto y reglas de interfaces públicas. |
| [despliegue.md](despliegue.md) | Topologías de despliegue: autohospedada sencilla y SaaS con alta disponibilidad. |
| [roles.md](roles.md) | Roles por producto y principios de permisos. |

## Decisiones fijadas

- **Un repositorio** con componentes y utilidades compartidos (npm workspaces: `apps/*`, `packages/*`).
- **Tres aplicaciones independientes**: Soporte TI (`apps/it`), MSP (`apps/msp`) y Atención al cliente (`apps/customer-service`), cada una con su API, cuentas, base de datos y versiones.
- **Backend Go como monolito modular por producto**, reutilizando los servicios y patrones actuales; la API existente es la base de Soporte TI.
- **Frontend Next.js** (App Router, TypeScript, React 19, Tailwind CSS 4, HeroUI 3) con paquetes compartidos de diseño, traducciones, contratos y utilidades.
- **Dos modalidades** para los tres productos: autohospedada y SaaS con alta disponibilidad desde el primer lanzamiento en producción.
- **Bilingüe** español e inglés, con español como idioma inicial.
