# Roles y permisos por producto

Los roles describen responsabilidades; **ningún rol concede acceso por sí solo en la interfaz**: cada acción se valida en el servidor contra el permiso y el ámbito (producto, organización, cliente) del recurso. Estos son los roles que usa el catálogo de casos de uso ([`docs/use-cases/`](../use-cases/README.md)).

## Principios

- Permiso mínimo necesario por rol y por ámbito.
- Separación de funciones en flujos sensibles (quien solicita no aprueba su propia solicitud restringida; quien registra horas no las aprueba).
- El rol de «auditor» es de solo lectura y su alcance también está delimitado por ámbito.
- Las integraciones actúan como identidad propia («servicio de integración autorizado»), con credenciales revocables y alcance acotado, nunca como un usuario humano compartido.

## Soporte TI

| Rol | Responsabilidad |
| --- | --- |
| Administrador de organización | Configura la organización, identidades, permisos, calendarios y integraciones. |
| Solicitante | Empleado que pide servicios e incidentes y sigue sus solicitudes. |
| Agente | Atiende incidentes y solicitudes en las colas de su equipo. |
| Responsable de servicio | Supervisa colas, SLA, catálogo y problemas de su ámbito. |
| Aprobador designado | Decide aprobaciones asignadas, con delegación temporal trazable. |
| Gestor de activos | Mantiene activos, CMDB, licencias y custodias. |
| Gestor de cambios | Tramita cambios, riesgos, ventanas y revisiones. |
| Auditor de solo lectura | Consulta actividad, informes y evidencias sin modificar. |

## MSP

| Rol | Responsabilidad |
| --- | --- |
| Administrador del proveedor | Configura el espacio del proveedor, clientes, técnicos e integraciones. |
| Responsable de servicio | Coordina la mesa multicliente, guardias y escalaciones. |
| Técnico del proveedor | Atiende tickets e intervenciones de los clientes a los que tiene acceso. |
| Gestor de contratos y costes | Contratos, SLA comerciales, horas, gastos y preparación de facturación. |
| Contacto del cliente | Persona autorizada del cliente: abre y sigue sus tickets y autoriza sesiones remotas. |
| Auditor de solo lectura | Revisa actividad e informes dentro de su alcance asignado. |

## Atención al cliente

| Rol | Responsabilidad |
| --- | --- |
| Administrador | Configura canales, colas, plantillas, privacidad e integraciones de la empresa. |
| Supervisor de atención | Supervisa la bandeja, asignaciones, calidad y SLA. |
| Agente | Atiende conversaciones de sus colas por todos los canales. |
| Cliente externo | Accede solo a sus conversaciones, casos y datos autorizados. |
| Auditor de solo lectura | Consulta auditoría e informes anonimizados dentro de su alcance. |
