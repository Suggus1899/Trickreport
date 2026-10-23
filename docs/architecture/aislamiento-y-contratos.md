# Aislamiento de datos y contratos

## Límites de datos

| Contexto | Regla de aislamiento |
| --- | --- |
| Entre productos | Bases de datos, credenciales, sesiones, claves de firma y permisos independientes. Una cuenta puede existir en varios productos con accesos distintos. |
| Soporte TI | La organización delimita empleados, equipos, activos y solicitudes. |
| MSP | El proveedor delimita la organización operadora; cada cliente delimita tickets, activos, contratos y accesos de sus contactos y técnicos. |
| Atención al cliente | La empresa delimita su operación; cada cliente externo accede únicamente a sus conversaciones y recursos autorizados. |

Reglas derivadas:

- Ninguna consulta cruza el límite de su ámbito: toda lectura o escritura valida el ámbito (producto, organización, cliente) en el servidor antes de tocar datos.
- Un identificador válido en un ámbito no revela existencia ni datos en otro: la respuesta ante un identificador ajeno es indistinguible de «no existe».
- Las exportaciones, los informes y las búsquedas respetan el mismo ámbito que la operación equivalente en pantalla.

## Contratos e interfaces públicas

- Se mantiene **`/api/v1`** para los contratos actuales de Soporte TI y se amplían recursos por módulos. MSP y Atención publican sus propios recursos bajo su propia API.
- Se publica **OpenAPI por producto**; los tipos que usa el frontend se generan desde ese OpenAPI, no se escriben a mano.
- La autorización vive en Go: la visibilidad de un botón nunca sustituye una comprobación de permisos en el servidor.
- Se conservan cookies de sesión seguras, CSRF, revocación de sesiones, MFA e idempotencia en los flujos que lo necesitan.
- Los eventos WebSocket se definen por producto y sus destinatarios se filtran por ámbito antes de enviarlos.
- Los webhooks de integraciones se tratan como entradas autenticadas, repetibles y potencialmente desordenadas: verificación de firma, clave idempotente y tolerancia al desorden.
- Los borradores offline se aíslan por producto, organización y usuario; los borradores antiguos sin propietario verificable se conservan en cuarentena, no se asignan por inferencia.

## Integraciones iniciales

| Necesidad | Elección inicial |
| --- | --- |
| Correo | SMTP para salida e IMAP para entrada; Mailpit para pruebas locales. |
| Directorio empresarial | Conservar y completar LDAP. |
| Monitorización y descubrimiento (MSP) | Zabbix mediante su API y webhooks. |
| Acceso remoto (MSP) | MeshCentral, con permisos, autorización y trazabilidad desde Trickreport. |
| WhatsApp (Atención) | WhatsApp Business Platform, API oficial de Meta; primero con recursos de prueba. |
| Automatización externa | API documentada y webhooks de entrada y salida. |
