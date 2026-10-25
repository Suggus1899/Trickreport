# MSP · Integraciones: Zabbix y MeshCentral

> Estado: **adaptadores configurables implementados; validación viva
> pendiente**. En este entorno no hay instancias ni credenciales
> reales de Zabbix ni de MeshCentral, así que no se inventó ninguna
> conexión viva: lo que sigue es lo que hay que configurar para
> activarlas y cómo se comprueba cada una.

## Zabbix: alertas por webhook

### Qué hace Trickreport

- `POST /api/v1/msp/monitoring/webhook` recibe una entrega de Zabbix.
- Autenticación: cabecera `X-MSP-Webhook-Token` con el token de la
  fuente, comparado contra el hash SHA-256 guardado en
  `msp_monitoring_sources.webhook_token_hash`. El token en claro no se
  guarda en ningún sitio y no se registra en logs.
- Tenant: cabecera `X-Tenant-ID` con el slug de la organización
  operadora (el middleware lo resuelve a UUID). Sin sesión de usuario:
  Zabbix no puede iniciar sesión en Trickreport, así que esta ruta
  está fuera del grupo autenticado y exenta de CSRF solo para este
  path (ver `middleware/csrf.go`).
- Idempotencia: si ya existe una alerta **abierta** para la misma
  fuente y el mismo recurso, la entrega actualiza esa alerta en lugar
  de crear otra; el índice parcial
  `idx_msp_alerts_open_dedup` lo garantiza también en base de datos
  ante entregas concurrentes. Solo la primera alerta abierta crea un
  ticket, con cliente, contrato activo y SLA resueltos en el servidor.
- Una entrega con `status: "resolved"` resuelve la alerta abierta del
  recurso; no crea ticket y no cierra el ticket existente (el cierre
  lo decide el técnico).

### Cuerpo esperado (JSON)

```json
{
  "client_id": "<uuid del cliente MSP>",
  "external_resource": "srv-db-01",
  "severity": "high",
  "status": "open",
  "title": "Servidor caído",
  "message": "Zabbix: host srv-db-01 no responde desde 2026-10-24T10:00Z"
}
```

`severity` admite `info`, `warning`, `average`, `high`, `disaster`.
La respuesta es `{"alert": {...}, "ticket_created": true|false}`.

### Activación en Zabbix

1. En Trickreport, crear la fuente en **Monitorización** con nombre,
   URL base y token de webhook. Guardar el token fuera de
   Trickreport: solo se conserva su hash.
2. En Zabbix, crear un tipo de medio webhook con URL
   `https://<trickreport>/api/v1/msp/monitoring/webhook` y cabeceras
   `X-Tenant-ID: <slug>` y `X-MSP-Webhook-Token: <token>`, con el
   cuerpo anterior. El `client_id` debe corresponder al cliente MSP
   que representa al host o grupo de hosts (el mapeo recurso→cliente
   se define en la plantilla del webhook de Zabbix, no en
   Trickreport).
3. Comprobación: disparar una alerta de prueba dos veces seguidas y
   verificar que solo se crea un ticket y que la segunda entrega
   devuelve `ticket_created: false`; después enviar el evento de
   resolución y verificar que la alerta queda `resolved`.

### Pendiente de validación viva

- Instancia Zabbix real, token real y entrega de prueba de extremo a
  extremo (pendiente de credenciales e instancia del operador).

## MeshCentral: acceso remoto con trazabilidad

### Qué hace Trickreport

Trickreport **no abre la sesión remota**: registra quién la pidió,
quién la autorizó y cuándo empezó y terminó, vinculada a cliente,
ticket y técnico (`msp_remote_sessions`).

- `POST /msp/remote-sessions`: cualquier usuario autenticado solicita
  una sesión indicando el `external_session_id` de MeshCentral.
- `POST /msp/remote-sessions/{id}/authorize` (`{"authorize": true|false}`):
  solo personal (admin/agent). Un técnico nunca se autoautoriza.
- `POST /msp/remote-sessions/{id}/start`: solo si la sesión está
  `authorized`, y solo el técnico propietario o el personal.
- `POST /msp/remote-sessions/{id}/end`: cierra la sesión y registra
  `ended_at`. Solo personal.

### Activación con MeshCentral

1. Desplegar MeshCentral según su documentación de integración
   (embedding) y obtener el identificador de sesión/dispositivo.
2. Al solicitar la sesión en Trickreport, pegar ese identificador en
   el campo "ID de sesión MeshCentral" junto al ticket y la nota de
   auditoría que justifica el acceso.
3. El personal autoriza la sesión en Trickreport y solo entonces el
   técnico la inicia desde MeshCentral; al terminar, se cierra en
   Trickreport para dejar la traza completa.

### Pendiente de validación viva

- Instancia MeshCentral real y sesión de prueba de extremo a extremo
  (pendiente de credenciales e instancia del operador). No se
  implementó inicio de sesión automático en MeshCentral desde
  Trickreport: el enlace es manual por diseño en esta etapa.

## Descubrimiento de inventario

El inventario por cliente reutiliza los activos de la Etapa 3
(`assets`, `ci_types`): `msp_client_assets` solo asigna cada activo a
un cliente. El descubrimiento es manual o por importación
documentada: se registra el hallazgo en
`msp_discovered_assets` (ID externo, nombre, tipo, payload) y se
importa eligiendo el tipo de CI; la importación crea el activo, lo
vincula al cliente y marca el hallazgo como importado en una sola
transacción. No hay agente de descubrimiento propio en esta etapa
(quedó fuera del plan maestro).
