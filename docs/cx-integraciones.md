# CX · Integraciones: correo y WhatsApp

> Estado: **adaptadores configurables implementados; validación viva
> pendiente**. En este entorno no hay instancias ni credenciales
> reales de correo ni de WhatsApp, así que no se inventó ninguna
> conexión viva: lo que sigue es lo que hay que configurar para
> activarlas y cómo se comprueba cada una.

## Correo: SMTP para salida e IMAP para entrada

### Qué hace Trickreport

- `GET /api/v1/cx/email-accounts` y `POST /api/v1/cx/email-accounts`
  (personal) guardan la configuración de cada cuenta: hosts y
  puertos SMTP/IMAP, usuarios y remitente, y si está activa.
- **No se guarda ninguna contraseña.** Las credenciales SMTP/IMAP
  viven en la configuración del operador (variables de entorno o
  gestor de secretos de la instalación), nunca en
  `cx_email_accounts` ni en logs.
- `POST /api/v1/cx/email/ingest` (personal, autenticado) es el punto
  de ingesta de un correo entrante: lo añade a la conversación
  abierta de su hilo (`thread_id`) o abre una nueva, corresponde el
  remitente con un cliente por correo o crea uno, y es idempotente
  por `external_message_id` (un mensaje repetido devuelve conflicto
  y no se duplica).

### Cuerpo esperado (JSON)

```json
{
  "account_id": "<uuid de la cuenta de correo>",
  "from_email": "ana@ejemplo.es",
  "subject": "Problema con mi pedido",
  "body": "El pedido ORD-1042 llegó incompleto.",
  "thread_id": "hilo-del-servidor-de-correo",
  "external_message_id": "<id único del mensaje>"
}
```

La respuesta es la conversación resultante. El envío saliente por
SMTP y el sondeo IMAP periódico son trabajo del operador: el
sondeo puede implementarse como un trabajo que lea el buzón (o
Mailpit en local, como fija el plan maestro) y llame a este
endpoint con cada mensaje nuevo. Ese trabajo todavía no está
implementado en Trickreport; el endpoint es el contrato que usará.

### Activación

1. En Trickreport, crear la cuenta en **Canales → Correo** con
   host y puerto SMTP (587 por defecto), host y puerto IMAP (993
   por defecto), usuarios y remitente.
2. Configurar las credenciales SMTP/IMAP en la configuración del
   operador, fuera de la base de datos.
3. Comprobación: enviar un correo de prueba al buzón, verificar
   que la ingesta crea (o reutiliza) la conversación correcta y
   que una segunda ingesta con el mismo `external_message_id` no
   duplica el mensaje.

### Pendiente de validación viva

- Servidor SMTP/IMAP real (o Mailpit local), credenciales reales y
  prueba de extremo a extremo de entrada y salida.

## WhatsApp: WhatsApp Business Platform (API oficial de Meta)

### Qué hace Trickreport

- `GET /api/v1/cx/whatsapp-accounts` y
  `POST /api/v1/cx/whatsapp-accounts` (personal) guardan la
  configuración de cada cuenta: nombre, `business_id`,
  `phone_number_id`, nombre visible y si está activa.
- Al crear la cuenta se puede entregar un `verify_token`: solo se
  guarda su hash SHA-256 en
  `cx_whatsapp_accounts.webhook_verify_token_hash`. El token en
  claro no se guarda en ningún sitio y no se registra en logs.
  **El token de acceso de Meta no se guarda en absoluto.**
- `GET /api/v1/cx/whatsapp/webhook` (sin sesión; empresa por
  `X-Tenant-ID`, exento de CSRF solo para este path) responde a la
  verificación de Meta: `hub.mode=subscribe` y
  `hub.verify_token` que coincida con el hash guardado devuelven
  `hub.challenge` tal cual.
- `POST /api/v1/cx/whatsapp/webhook` (sin sesión; empresa por
  `X-Tenant-ID`, exento de CSRF solo para este path) recibe los
  mensajes. Autenticación: firma HMAC `X-Hub-Signature-256` sobre
  el cuerpo con el secreto de aplicación del operador, en la
  variable `CX_WHATSAPP_APP_SECRET`. Sin secreto configurado, el
  endpoint responde 503 y no ingiere nada.
- Ingesta: cada mensaje de texto crea o reutiliza la conversación
  de WhatsApp del remitente (hilo `whatsapp:<teléfono>`) y se
  guarda idempotentemente por su id de Meta. Los estados de
  entrega y los mensajes no de texto se confirman y se ignoran.

### Activación

1. En Meta (aplicación de WhatsApp Business Platform), crear la
   cuenta y el número, y anotar `business_id` y `phone_number_id`.
2. En Trickreport, crear la cuenta en **Canales → WhatsApp** con
   esos identificadores, el nombre visible y un token de
   verificación. Guardar el token fuera de Trickreport: solo se
   conserva su hash.
3. Configurar `CX_WHATSAPP_APP_SECRET` con el secreto de aplicación
   de Meta en la configuración del operador, fuera de la base de
   datos.
4. En Meta, suscribir el webhook a
   `https://<trickreport>/api/v1/cx/whatsapp/webhook` con la
   cabecera `X-Tenant-ID: <slug>` y el token de verificación del
   paso 2, y suscribir el campo `messages`.
5. Comprobación: Meta verifica el webhook (GET); después enviar un
   mensaje de prueba desde el número de prueba y verificar que
   aparece en la bandeja como conversación de WhatsApp del cliente
   correcto, y que reenviar la misma entrega no duplica el mensaje.

### Pendiente de validación viva

- Cuenta empresarial real de Meta, número real, secreto de
  aplicación real y entrega de prueba de extremo a extremo.
- Resolución de la cuenta desde `metadata.phone_number_id` en el
  webhook: hoy la ingesta se enruta por empresa y teléfono del
  remitente, no por la cuenta de destino.
- Envío saliente por la API de Meta (responder desde la bandeja):
  no implementado en esta etapa; la bandeja registra los mensajes,
  pero no los envía a Meta.

## Canales que no necesitan credenciales

- **Portal web y chat en vivo:** son la misma bandeja con canal
  `portal` o `chat`. No hay ninguna integración externa que
  activar: el cliente abre la conversación en el portal y el agente
  la atiende desde la bandeja, con el historial completo.
