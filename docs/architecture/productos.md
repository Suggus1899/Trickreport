# Arquitectura por producto

Cada producto es una aplicación independiente que comparte marca, paquetes y patrones, pero no datos ni sesiones. El backend Go se compone por producto a partir de módulos reutilizables: un mismo módulo (tickets, SLA, conocimiento) puede formar parte de varios productos con la configuración y los permisos de cada uno.

## Soporte TI (`apps/it`)

**Propósito:** gestionar los servicios tecnológicos de una organización.

- **Base:** la API actual de Trickreport (`/api/v1`) y sus módulos existentes: tickets, artículos de conocimiento, SLA, automatizaciones, autenticación con MFA y LDAP.
- **Módulos del plan:** identidad y organización, incidentes y seguimiento, solicitudes y catálogo de servicios, conocimiento y autoservicio, SLA y escalaciones, activos y CMDB, aprobaciones, problemas y cambios, informes y auditoría.
- **Recorridos:** portal del solicitante (empleados), consola de agentes, administración, informes.
- **Etapas:** Etapa 2 (migración a Next.js de los recorridos actuales) y Etapa 3 (ITSM completo por módulos).
- **Casos:** `TI-001` a `TI-200`.

## MSP (`apps/msp`)

**Propósito:** gestionar la operación de un proveedor que atiende a varias empresas cliente.

- **Composición propia** sobre los módulos compartidos de mesa de servicio, SLA, activos e informes, con el ámbito «cliente» como límite de datos de primer nivel.
- **Módulos del plan:** clientes y contactos, mesa de servicio multicliente, contratos y SLA, técnicos y planificación, horas, gastos y costes, proyectos de servicio, rentabilidad y preparación de facturación, inventario y descubrimiento integrado, monitorización (Zabbix) y acceso remoto integrado (MeshCentral), informes y auditoría.
- **Recorridos:** consola del proveedor con cambio explícito de cliente, contrato visible y contexto económico de la intervención; vista del contacto del cliente limitada a su ámbito.
- **Etapa:** Etapa 4.
- **Casos:** `MSP-001` a `MSP-200`.

## Atención al cliente (`apps/customer-service`)

**Propósito:** gestionar las conversaciones y solicitudes de clientes externos de una empresa.

- **Composición propia** centrada en conversaciones omnicanal: portal web, correo (SMTP/IMAP), chat en vivo y WhatsApp Business Platform (API oficial de Meta).
- **Módulos del plan:** clientes y contactos, bandeja de conversaciones, correo electrónico, portal y chat, WhatsApp, asignación y automatización, conocimiento y autoservicio, SLA y escalaciones, encuestas y satisfacción, administración e informes.
- **Recorridos:** bandeja del agente centrada en conversaciones, canal e historial; portal del cliente externo limitado a sus conversaciones autorizadas.
- **Etapa:** Etapa 5.
- **Casos:** `CX-001` a `CX-200`.

## Paquetes compartidos

- **Diseño:** tokens y componentes base de la marca (ver [`docs/brand/`](../brand/README.md)).
- **Traducciones:** catálogos es/en con los mismos identificadores.
- **Contratos:** tipos generados desde el OpenAPI de cada producto.
- **Utilidades:** código con el mismo comportamiento y responsabilidad en los tres productos; lo que difiere entre productos no se comparte.
