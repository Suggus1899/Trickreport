# Registro de la Etapa 5 — Entregar atención al cliente

Estado: **completada en código**, con las verificaciones ejecutables en
este entorno en verde. Quedan como gates externos la suite Go con
PostgreSQL 16 y el E2E con Playwright (CI), porque en este entorno no
hay Go en el PATH ni PostgreSQL disponible: ningún test Go ni ninguna
migración pudo correrse aquí. La validación viva de correo y WhatsApp
queda pendiente de instancias y credenciales reales (ver
`docs/cx-integraciones.md`): los canales son configurables y no se
inventó ninguna conexión viva ni ninguna credencial.

## Condición de salida

Conversaciones completas, canales operativos y permisos comprobados en
código: cada consulta filtra por `tenant_id` y por `customer_id`, y un
cliente externo ve cualquier otro cliente —y sus conversaciones— como
"no encontrado" (nunca como dato parcial). La regla vive en Go
(`domain.CanAccessCustomer`, aplicada en el servicio); la interfaz solo
filtra lo que muestra. La asignación también se comprueba en el
servidor: tomar una conversación ya asignada a otro agente falla con
conflicto (`domain.CanTakeConversation`), aunque la UI mostrara un
estado anterior.

## Qué se construyó (backend + frontend contra API real, no maquetas)

### Clientes externos y contactos (`/cx/customers`)

- CRUD de clientes externos (nombre, correo, teléfono, empresa,
  estado, consentimiento de contacto y solicitud de eliminación de
  datos) y contactos por cliente, con rol `contact` o `client_admin`.
  Un contacto puede enlazarse a un usuario de la organización
  (`user_id`); sin enlace es solo un registro del cliente.
- Migración `000016_cx`. Los correos reales son únicos por empresa;
  el correo vacío queda fuera del índice único para que la ingesta
  pueda crear clientes solo con teléfono (WhatsApp).

### Bandeja y gestión de conversaciones (`/cx/conversations`)

- Conversaciones por canal (`email`, `portal`, `chat`, `whatsapp`)
  con asunto, estado, prioridad, asignación, SLA resuelto y
  vencimientos. Máquina de estados
  `open → in_progress → waiting_customer → resolved → closed`, con
  reapertura (`resolved → in_progress`, `closed → in_progress`).
- Tomar, liberar y escalar desde la bandeja. Escalar sube un nivel de
  prioridad (crítica se queda en crítica) y marca la conversación
  como escalada.
- Mensajes con remitente (`customer`, `agent`, `system`), notas
  internas (`is_internal`) y canal. Las notas internas las filtra el
  servicio para clientes externos: nunca llegan a su lista de
  mensajes.
- Reglas de historial mantenidas por el servicio: la primera respuesta
  visible de un agente marca `first_response_at`; una respuesta de
  agente a una conversación en espera la devuelve a en curso; una
  respuesta de cliente a una conversación resuelta la reabre
  (caso CX-037 del catálogo). Una conversación cerrada no acepta
  mensajes nuevos.

### Correo, portal y chat en vivo

- El cliente externo crea y sigue sus conversaciones desde el portal
  (misma bandeja, filtrada por sus clientes); el agente trabaja
  desde la bandeja con historial completo y contexto de canal.
- Cuentas de correo configurables (`/cx/email-accounts`): hosts y
  puertos SMTP/IMAP, usuarios y remitente. **Sin contraseñas en la
  base de datos**: viven en la configuración del operador.
- Ingesta autenticada de correo (`POST /cx/email/ingest`, solo
  personal): un correo entrante se añade a la conversación abierta
  de su hilo o abre una nueva, y el remitente se corresponde con un
  cliente por correo o se crea uno. La ingesta es idempotente por
  `external_message_id`: un mensaje repetido no se duplica.
- El chat en vivo es el canal `chat` de la misma bandeja, con el SLA
  de respuesta más corto por defecto (15 minutos).

### WhatsApp (`/cx/whatsapp-accounts`, webhook)

- Cuentas de WhatsApp Business Platform con la API oficial de Meta:
  identificadores de empresa y número, nombre visible y solo el hash
  SHA-256 del token de verificación del webhook. **Sin token de
  acceso en la base de datos.**
- Verificación del webhook (GET): modo `subscribe` y token que
  coincide con el hash guardado; el desafío se devuelve tal cual.
- Ingesta del webhook (POST): autenticada por la firma HMAC
  `X-Hub-Signature-256` con el secreto de aplicación del operador
  (`CX_WHATSAPP_APP_SECRET`). Sin secreto configurado, el endpoint
  rechaza todo: ningún mensaje puede inyectarse antes de la
  activación. Cada mensaje de texto se ingiere idempotentemente por
  su id de Meta; los estados de entrega y los mensajes no de texto
  se confirman y se ignoran.
- Limitación conocida, no oculta: el webhook todavía no resuelve la
  cuenta desde `metadata.phone_number_id`; la ingesta se enruta por
  empresa y teléfono del remitente, y el servicio no usa el id de
  cuenta para enrutar. Resolverlo es trabajo pendiente, no un
  comportamiento inventado.

### Asignación y automatización

- Políticas SLA por canal y prioridad (`/cx/sla-policies`, una por
  canal y prioridad): ganan a los valores por defecto del dominio
  (chat 15/240, WhatsApp 60/720, portal 120/1440, correo 240/1440
  minutos de respuesta/resolución; la prioridad crítica reduce a la
  mitad ambos plazos, con un mínimo de 15 minutos de respuesta).
- Reglas de asignación (`/cx/assignment-rules`): la primera regla
  activa que coincide con canal y prioridad, en orden de posición,
  asigna la conversación al crearla y la pasa a en curso.
- Respuestas automáticas (`/cx/auto-replies`) por palabra clave
  (insensible a mayúsculas) o por canal: un mensaje de cliente que
  coincide recibe una respuesta inmediata del sistema.

### SLA y escalaciones

- Cada conversación guarda sus minutos SLA y sus vencimientos de
  respuesta y resolución, calculados al crearla. Una conversación se
  considera incumplida cuando pasa su vencimiento de resolución sin
  estar resuelta o cerrada, o su vencimiento de primera respuesta
  sin respuesta de agente. El informe cuenta las incumplidas.

### Encuestas y satisfacción (`/cx/surveys`)

- Una encuesta por conversación, solo tras resolver o cerrar, solo
  del cliente propietario, con valoración 1–5 y comentario. El
  informe muestra el recuento y la satisfacción media.

### Administración e informes (`/cx/report`)

- Totales de clientes, conversaciones abiertas y resueltas,
  incumplimientos de SLA, encuestas y satisfacción media, y volumen
  por canal. Los clientes externos ven el informe limitado a sus
  clientes; el personal ve toda la empresa.

### Conocimiento y autoservicio

- Reutilizado, no duplicado: la bandeja enlaza con el módulo de
  conocimiento existente (`/articles`, `knowledge.tsx`), que ya
  cubre búsqueda y lectura de artículos para solicitantes.

## Verificado aquí

- `npm run check` (tsc del workspace): en verde tras corregir un
  cierre de etiqueta en `cx.tsx`.
- `npm run lint`: en verde (solo el aviso preexistente sobre el
  directorio `pages`, ajeno a esta etapa).
- Vitest del módulo CX (`packages/workspace/tests/cx.test.ts`):
  7 pruebas en verde (transiciones espejo de Go, incumplimiento de
  SLA, escalación por niveles y media de satisfacción).
- Compilación de producción de las 3 apps Next.js (`it`,
  `msp` y `customer-service`, con `--webpack`): en verde.

## Pendiente (no verificado, no presentado como verde)

- Suite Go con PostgreSQL 16 (`go test ./... -race`, incluidos los
  tests de dominio CX): Go no está en el PATH de este entorno.
- Migración `000016_cx` aplicada sobre una base vacía y sobre la
  versión anterior: pendiente de PostgreSQL.
- Playwright E2E con API real.
- Activación viva de correo (SMTP/IMAP reales, o Mailpit local) y de
  WhatsApp (credenciales empresariales reales de Meta y secreto de
  aplicación): ver `docs/cx-integraciones.md`.
- Resolución de la cuenta de WhatsApp desde
  `metadata.phone_number_id` en el webhook (hoy se enruta por
  empresa y teléfono del remitente).
- Envío saliente por WhatsApp (responder desde la bandeja por la
  API de Meta): la bandeja registra mensajes, pero el envío a Meta
  no está implementado; queda para una entrega posterior, no se
  presenta como operativo.
