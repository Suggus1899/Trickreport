# Trickreport · Soporte TI

Catálogo de **200 casos propuestos**. Todos comienzan en `planned`; su inclusión no afirma que la función exista.

| Área | Casos |
| --- | ---: |
| Identidad, organización y permisos | 16 |
| Incidentes y seguimiento | 32 |
| Solicitudes y catálogo | 22 |
| Conocimiento y autoservicio | 14 |
| SLA, calendarios y escalaciones | 20 |
| Activos y CMDB | 24 |
| Aprobaciones | 12 |
| Problemas y cambios | 28 |
| Automatizaciones e integraciones | 16 |
| Informes y auditoría | 16 |

<a id="ti-001"></a>
## TI-001 · Invitar empleado a la organización

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ana recibe una invitación limitada a Aurora

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. ana@aurora.example, departamento Radiología

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: ana@aurora.example, departamento Radiología
3. Ejecuta «Invitar empleado a la organización»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ana recibe una invitación limitada a Aurora

**Alternativas y errores:**
- Correo ya vinculado: informar sin duplicar membresía.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ana recibe una invitación limitada a Aurora
- Se reproduce y controla esta alternativa: Correo ya vinculado: informar sin duplicar membresía.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de invitar empleado a la organización; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-002"></a>
## TI-002 · Iniciar sesión local

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se crea una sesión revocable para Ana

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana presenta su correo y contraseña válidos

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana presenta su correo y contraseña válidos
3. Ejecuta «Iniciar sesión local»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crea una sesión revocable para Ana

**Alternativas y errores:**
- Contraseña incorrecta: respuesta uniforme y límite de intentos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crea una sesión revocable para Ana
- Se reproduce y controla esta alternativa: Contraseña incorrecta: respuesta uniforme y límite de intentos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de iniciar sesión local; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-003"></a>
## TI-003 · Acceder mediante directorio LDAP

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se verifica el directorio y se aplica el rol autorizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Usuario aperez pertenece al grupo TI-Aurora

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Usuario aperez pertenece al grupo TI-Aurora
3. Ejecuta «Acceder mediante directorio LDAP»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se verifica el directorio y se aplica el rol autorizado

**Alternativas y errores:**
- Directorio inaccesible: informar indisponibilidad sin conceder acceso.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se verifica el directorio y se aplica el rol autorizado
- Se reproduce y controla esta alternativa: Directorio inaccesible: informar indisponibilidad sin conceder acceso.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de acceder mediante directorio ldap; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-004"></a>
## TI-004 · Activar segundo factor

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Los accesos posteriores exigen el segundo factor

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana confirma un código TOTP de su dispositivo

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana confirma un código TOTP de su dispositivo
3. Ejecuta «Activar segundo factor»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Los accesos posteriores exigen el segundo factor

**Alternativas y errores:**
- Código incorrecto: no activar el secreto propuesto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Los accesos posteriores exigen el segundo factor
- Se reproduce y controla esta alternativa: Código incorrecto: no activar el secreto propuesto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de activar segundo factor; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-005"></a>
## TI-005 · Completar desafío MFA

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se emite una sesión tras consumir el desafío

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Sesión pendiente y código TOTP de seis dígitos

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Sesión pendiente y código TOTP de seis dígitos
3. Ejecuta «Completar desafío MFA»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se emite una sesión tras consumir el desafío

**Alternativas y errores:**
- Desafío usado o vencido: rechazar su reutilización.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se emite una sesión tras consumir el desafío
- Se reproduce y controla esta alternativa: Desafío usado o vencido: rechazar su reutilización.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de completar desafío mfa; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-006"></a>
## TI-006 · Recuperar contraseña local

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Nueva contraseña aplicada y sesiones anteriores revocadas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana solicita recuperación y recibe un enlace válido

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana solicita recuperación y recibe un enlace válido
3. Ejecuta «Recuperar contraseña local»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nueva contraseña aplicada y sesiones anteriores revocadas

**Alternativas y errores:**
- Enlace vencido: permitir pedir otro sin cambiar credenciales.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nueva contraseña aplicada y sesiones anteriores revocadas
- Se reproduce y controla esta alternativa: Enlace vencido: permitir pedir otro sin cambiar credenciales.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de recuperar contraseña local; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-007"></a>
## TI-007 · Cerrar una sesión concreta

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sesión elegida deja de autorizar peticiones

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana identifica el navegador del equipo compartido

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana identifica el navegador del equipo compartido
3. Ejecuta «Cerrar una sesión concreta»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sesión elegida deja de autorizar peticiones

**Alternativas y errores:**
- Sesión de otra persona: denegar su revocación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sesión elegida deja de autorizar peticiones
- Se reproduce y controla esta alternativa: Sesión de otra persona: denegar su revocación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de cerrar una sesión concreta; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-008"></a>
## TI-008 · Cambiar permisos de un agente

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Nuevas peticiones pierden acceso a colas de agentes

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro pasa de agente a solicitante

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro pasa de agente a solicitante
3. Ejecuta «Cambiar permisos de un agente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nuevas peticiones pierden acceso a colas de agentes

**Alternativas y errores:**
- Último administrador: impedir dejar la organización sin gestión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nuevas peticiones pierden acceso a colas de agentes
- Se reproduce y controla esta alternativa: Último administrador: impedir dejar la organización sin gestión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de cambiar permisos de un agente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-009"></a>
## TI-009 · Desactivar a un empleado

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Acceso revocado conservando autoría de tickets históricos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro termina su relación laboral el viernes

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro termina su relación laboral el viernes
3. Ejecuta «Desactivar a un empleado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Acceso revocado conservando autoría de tickets históricos

**Alternativas y errores:**
- Tickets asignados: mostrar pendientes para reasignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Acceso revocado conservando autoría de tickets históricos
- Se reproduce y controla esta alternativa: Tickets asignados: mostrar pendientes para reasignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de desactivar a un empleado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-010"></a>
## TI-010 · Crear un departamento

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se dispone de agrupación organizativa para solicitudes

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Departamento Laboratorio con responsable Marta

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Departamento Laboratorio con responsable Marta
3. Ejecuta «Crear un departamento»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se dispone de agrupación organizativa para solicitudes

**Alternativas y errores:**
- Nombre duplicado dentro de Aurora: pedir diferenciación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se dispone de agrupación organizativa para solicitudes
- Se reproduce y controla esta alternativa: Nombre duplicado dentro de Aurora: pedir diferenciación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de crear un departamento; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-011"></a>
## TI-011 · Asignar una sede

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Su perfil permite contextualizar servicios y activos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana trabaja en sede Valencia

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana trabaja en sede Valencia
3. Ejecuta «Asignar una sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Su perfil permite contextualizar servicios y activos

**Alternativas y errores:**
- Sede inactiva: impedir nuevas asignaciones.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Su perfil permite contextualizar servicios y activos
- Se reproduce y controla esta alternativa: Sede inactiva: impedir nuevas asignaciones.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de asignar una sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-012"></a>
## TI-012 · Delegar administración de un equipo

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Marta puede gestionar miembros de Redes

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Marta administra únicamente el equipo Redes

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marta administra únicamente el equipo Redes
3. Ejecuta «Delegar administración de un equipo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Marta puede gestionar miembros de Redes

**Alternativas y errores:**
- Intento sobre otro equipo: denegar sin modificar miembros.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Marta puede gestionar miembros de Redes
- Se reproduce y controla esta alternativa: Intento sobre otro equipo: denegar sin modificar miembros.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de delegar administración de un equipo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-013"></a>
## TI-013 · Actualizar preferencias personales

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Idioma y horarios se muestran según su perfil

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana selecciona inglés y zona America/Caracas

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana selecciona inglés y zona America/Caracas
3. Ejecuta «Actualizar preferencias personales»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Idioma y horarios se muestran según su perfil

**Alternativas y errores:**
- Zona inválida: conservar la preferencia previa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Idioma y horarios se muestran según su perfil
- Se reproduce y controla esta alternativa: Zona inválida: conservar la preferencia previa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de actualizar preferencias personales; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-014"></a>
## TI-014 · Consultar actividad de acceso

**Área:** identity · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ve fecha y dispositivo de sus propios accesos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana revisa sesiones abiertas durante siete días

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana revisa sesiones abiertas durante siete días
3. Ejecuta «Consultar actividad de acceso»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ve fecha y dispositivo de sus propios accesos

**Alternativas y errores:**
- Datos de otro usuario: mantenerlos fuera del resultado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ve fecha y dispositivo de sus propios accesos
- Se reproduce y controla esta alternativa: Datos de otro usuario: mantenerlos fuera del resultado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de consultar actividad de acceso; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-015"></a>
## TI-015 · Rotar credencial de una integración

**Área:** identity · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La nueva clave sustituye a la revocada con trazabilidad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Conector inventario tiene una clave comprometida

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conector inventario tiene una clave comprometida
3. Ejecuta «Rotar credencial de una integración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La nueva clave sustituye a la revocada con trazabilidad

**Alternativas y errores:**
- Cliente usa clave anterior: denegar sin fallback secreto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La nueva clave sustituye a la revocada con trazabilidad
- Se reproduce y controla esta alternativa: Cliente usa clave anterior: denegar sin fallback secreto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de rotar credencial de una integración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-016"></a>
## TI-016 · Verificar aislamiento de organización

**Área:** identity · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Una identidad Aurora solo obtiene recursos de Aurora

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Dos organizaciones tienen un ticket con número local 42

**Flujo principal:**
1. El actor accede a identidad, organización y permisos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Dos organizaciones tienen un ticket con número local 42
3. Ejecuta «Verificar aislamiento de organización»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Una identidad Aurora solo obtiene recursos de Aurora

**Alternativas y errores:**
- Identificador ajeno enviado manualmente: respuesta sin datos ajenos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Una identidad Aurora solo obtiene recursos de Aurora
- Se reproduce y controla esta alternativa: Identificador ajeno enviado manualmente: respuesta sin datos ajenos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Identidad, organización y permisos: operaciones de verificar aislamiento de organización; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-017"></a>
## TI-017 · Registrar incidente desde el portal

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se crea un ticket con referencia y confirmación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana no puede imprimir desde Radiología

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana no puede imprimir desde Radiología
3. Ejecuta «Registrar incidente desde el portal»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crea un ticket con referencia y confirmación

**Alternativas y errores:**
- Formulario incompleto: conservar texto y señalar campos necesarios.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crea un ticket con referencia y confirmación
- Se reproduce y controla esta alternativa: Formulario incompleto: conservar texto y señalar campos necesarios.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de registrar incidente desde el portal; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-018"></a>
## TI-018 · Registrar incidente en nombre de un empleado

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ticket identifica solicitante Ana y registrador Pedro

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro recibe llamada de Ana por acceso bloqueado

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro recibe llamada de Ana por acceso bloqueado
3. Ejecuta «Registrar incidente en nombre de un empleado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ticket identifica solicitante Ana y registrador Pedro

**Alternativas y errores:**
- Empleado fuera de organización: impedir su selección.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ticket identifica solicitante Ana y registrador Pedro
- Se reproduce y controla esta alternativa: Empleado fuera de organización: impedir su selección.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de registrar incidente en nombre de un empleado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-019"></a>
## TI-019 · Clasificar impacto y urgencia

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La prioridad refleja impacto amplio y urgencia alta

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Historia clínica caída para 80 empleados durante consulta

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Historia clínica caída para 80 empleados durante consulta
3. Ejecuta «Clasificar impacto y urgencia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La prioridad refleja impacto amplio y urgencia alta

**Alternativas y errores:**
- Valores incompatibles: solicitar corrección antes de guardar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La prioridad refleja impacto amplio y urgencia alta
- Se reproduce y controla esta alternativa: Valores incompatibles: solicitar corrección antes de guardar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de clasificar impacto y urgencia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-020"></a>
## TI-020 · Asignar incidente a un técnico

**Área:** incidents · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Pedro recibe responsabilidad y aviso del ticket

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 requiere a Pedro del equipo Redes

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 requiere a Pedro del equipo Redes
3. Ejecuta «Asignar incidente a un técnico»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Pedro recibe responsabilidad y aviso del ticket

**Alternativas y errores:**
- Técnico inactivo: impedir asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Pedro recibe responsabilidad y aviso del ticket
- Se reproduce y controla esta alternativa: Técnico inactivo: impedir asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de asignar incidente a un técnico; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-021"></a>
## TI-021 · Transferir incidente entre equipos

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Equipo receptor recibe contexto e historial íntegros

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Fallo de impresora corresponde a Infraestructura

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Fallo de impresora corresponde a Infraestructura
3. Ejecuta «Transferir incidente entre equipos»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Equipo receptor recibe contexto e historial íntegros

**Alternativas y errores:**
- Transferencia sin permiso: conservar el equipo actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Equipo receptor recibe contexto e historial íntegros
- Se reproduce y controla esta alternativa: Transferencia sin permiso: conservar el equipo actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de transferir incidente entre equipos; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-022"></a>
## TI-022 · Aceptar trabajo asignado

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El ticket entra en atención con fecha de inicio

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro toma el ticket 42 pendiente en su cola

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro toma el ticket 42 pendiente en su cola
3. Ejecuta «Aceptar trabajo asignado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El ticket entra en atención con fecha de inicio

**Alternativas y errores:**
- Otro agente lo tomó antes: mostrar conflicto actualizado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El ticket entra en atención con fecha de inicio
- Se reproduce y controla esta alternativa: Otro agente lo tomó antes: mostrar conflicto actualizado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de aceptar trabajo asignado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-023"></a>
## TI-023 · Solicitar información al empleado

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Pregunta pública y estado de espera quedan registrados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se necesita captura del error de impresión

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se necesita captura del error de impresión
3. Ejecuta «Solicitar información al empleado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Pregunta pública y estado de espera quedan registrados

**Alternativas y errores:**
- Comentario vacío: no cambiar estado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Pregunta pública y estado de espera quedan registrados
- Se reproduce y controla esta alternativa: Comentario vacío: no cambiar estado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de solicitar información al empleado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-024"></a>
## TI-024 · Responder una solicitud de información

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El agente recibe la respuesta y reanuda revisión

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana adjunta código de error E204 del ticket 42

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana adjunta código de error E204 del ticket 42
3. Ejecuta «Responder una solicitud de información»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El agente recibe la respuesta y reanuda revisión

**Alternativas y errores:**
- Ticket cerrado: aplicar política de reapertura visible.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El agente recibe la respuesta y reanuda revisión
- Se reproduce y controla esta alternativa: Ticket cerrado: aplicar política de reapertura visible.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de responder una solicitud de información; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-025"></a>
## TI-025 · Añadir nota interna de diagnóstico

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Nota interna visible únicamente para personal autorizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro documenta credenciales rotadas en servidor afectado

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro documenta credenciales rotadas en servidor afectado
3. Ejecuta «Añadir nota interna de diagnóstico»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nota interna visible únicamente para personal autorizado

**Alternativas y errores:**
- Solicitante consulta el hilo: nunca recibe esa nota.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nota interna visible únicamente para personal autorizado
- Se reproduce y controla esta alternativa: Solicitante consulta el hilo: nunca recibe esa nota.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de añadir nota interna de diagnóstico; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-026"></a>
## TI-026 · Adjuntar evidencia a un incidente

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Archivo queda vinculado y descargable con autorización

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana carga captura PNG de 300 KB

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana carga captura PNG de 300 KB
3. Ejecuta «Adjuntar evidencia a un incidente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Archivo queda vinculado y descargable con autorización

**Alternativas y errores:**
- Tipo real prohibido o tamaño excesivo: rechazar archivo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Archivo queda vinculado y descargable con autorización
- Se reproduce y controla esta alternativa: Tipo real prohibido o tamaño excesivo: rechazar archivo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de adjuntar evidencia a un incidente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-027"></a>
## TI-027 · Retirar un adjunto incorrecto

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se retira su acceso y queda registro de la retirada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Captura adjunta contiene datos personales innecesarios

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Captura adjunta contiene datos personales innecesarios
3. Ejecuta «Retirar un adjunto incorrecto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se retira su acceso y queda registro de la retirada

**Alternativas y errores:**
- Adjunto no pertenece al ticket: denegar operación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se retira su acceso y queda registro de la retirada
- Se reproduce y controla esta alternativa: Adjunto no pertenece al ticket: denegar operación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de retirar un adjunto incorrecto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-028"></a>
## TI-028 · Relacionar incidentes duplicados

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Un ticket principal conserva referencias a los duplicados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets 42 y 45 describen la misma caída

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets 42 y 45 describen la misma caída
3. Ejecuta «Relacionar incidentes duplicados»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Un ticket principal conserva referencias a los duplicados

**Alternativas y errores:**
- Tickets de otra organización: impedir vinculación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Un ticket principal conserva referencias a los duplicados
- Se reproduce y controla esta alternativa: Tickets de otra organización: impedir vinculación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de relacionar incidentes duplicados; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-029"></a>
## TI-029 · Dividir un reporte con varios problemas

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se crean tareas de atención separadas enlazadas al original

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Un ticket mezcla impresión y acceso al correo

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Un ticket mezcla impresión y acceso al correo
3. Ejecuta «Dividir un reporte con varios problemas»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crean tareas de atención separadas enlazadas al original

**Alternativas y errores:**
- Fallo al crear una parte: evitar un estado intermedio incoherente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crean tareas de atención separadas enlazadas al original
- Se reproduce y controla esta alternativa: Fallo al crear una parte: evitar un estado intermedio incoherente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de dividir un reporte con varios problemas; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-030"></a>
## TI-030 · Declarar incidente mayor

**Área:** incidents · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se activa coordinación y seguimiento de incidente mayor

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Caída general del sistema de turnos

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Caída general del sistema de turnos
3. Ejecuta «Declarar incidente mayor»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se activa coordinación y seguimiento de incidente mayor

**Alternativas y errores:**
- Incidente ya declarado: evitar equipos y avisos duplicados.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se activa coordinación y seguimiento de incidente mayor
- Se reproduce y controla esta alternativa: Incidente ya declarado: evitar equipos y avisos duplicados.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de declarar incidente mayor; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-031"></a>
## TI-031 · Publicar actualización de incidente mayor

**Área:** incidents · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Afectados autorizados reciben una actualización común

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Restablecimiento parcial a las 11:20 en Radiología

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Restablecimiento parcial a las 11:20 en Radiología
3. Ejecuta «Publicar actualización de incidente mayor»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Afectados autorizados reciben una actualización común

**Alternativas y errores:**
- Destinatario fuera de alcance: excluirlo del envío.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Afectados autorizados reciben una actualización común
- Se reproduce y controla esta alternativa: Destinatario fuera de alcance: excluirlo del envío.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de publicar actualización de incidente mayor; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-032"></a>
## TI-032 · Suscribirse a un incidente

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Recibe novedades públicas del incidente autorizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana sigue la caída de impresión de su sede

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana sigue la caída de impresión de su sede
3. Ejecuta «Suscribirse a un incidente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Recibe novedades públicas del incidente autorizado

**Alternativas y errores:**
- Pérdida de acceso posterior: retirar futuras notificaciones.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Recibe novedades públicas del incidente autorizado
- Se reproduce y controla esta alternativa: Pérdida de acceso posterior: retirar futuras notificaciones.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de suscribirse a un incidente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-033"></a>
## TI-033 · Registrar solución aplicada

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ticket contiene pasos, responsable y hora de resolución

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se reinstala controlador firmado de impresora R3

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se reinstala controlador firmado de impresora R3
3. Ejecuta «Registrar solución aplicada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ticket contiene pasos, responsable y hora de resolución

**Alternativas y errores:**
- Solución vacía: impedir marcar como resuelto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ticket contiene pasos, responsable y hora de resolución
- Se reproduce y controla esta alternativa: Solución vacía: impedir marcar como resuelto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de registrar solución aplicada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-034"></a>
## TI-034 · Confirmar resolución de un incidente

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Su confirmación queda asociada al cierre

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana imprime correctamente después de la reparación

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana imprime correctamente después de la reparación
3. Ejecuta «Confirmar resolución de un incidente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Su confirmación queda asociada al cierre

**Alternativas y errores:**
- Incidente aún activo: solicitar al agente resolución previa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Su confirmación queda asociada al cierre
- Se reproduce y controla esta alternativa: Incidente aún activo: solicitar al agente resolución previa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de confirmar resolución de un incidente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-035"></a>
## TI-035 · Reabrir un incidente resuelto

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ticket vuelve a atención conservando la solución anterior

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. La impresora falla nuevamente una hora después

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: La impresora falla nuevamente una hora después
3. Ejecuta «Reabrir un incidente resuelto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ticket vuelve a atención conservando la solución anterior

**Alternativas y errores:**
- Fuera de ventana: crear seguimiento vinculado según política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ticket vuelve a atención conservando la solución anterior
- Se reproduce y controla esta alternativa: Fuera de ventana: crear seguimiento vinculado según política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de reabrir un incidente resuelto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-036"></a>
## TI-036 · Cancelar un incidente reportado por error

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ticket cancelado conserva motivo y auditoría

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana descubre que eligió la impresora equivocada

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana descubre que eligió la impresora equivocada
3. Ejecuta «Cancelar un incidente reportado por error»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ticket cancelado conserva motivo y auditoría

**Alternativas y errores:**
- Trabajo ya facturado o ligado a cambio: requerir revisión autorizada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ticket cancelado conserva motivo y auditoría
- Se reproduce y controla esta alternativa: Trabajo ya facturado o ligado a cambio: requerir revisión autorizada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de cancelar un incidente reportado por error; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-037"></a>
## TI-037 · Consultar línea temporal del incidente

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se muestran eventos públicos ordenados con autor y fecha

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 cambió tres veces de estado

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 cambió tres veces de estado
3. Ejecuta «Consultar línea temporal del incidente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran eventos públicos ordenados con autor y fecha

**Alternativas y errores:**
- Eventos internos: filtrarlos antes de responder al cliente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran eventos públicos ordenados con autor y fecha
- Se reproduce y controla esta alternativa: Eventos internos: filtrarlos antes de responder al cliente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de consultar línea temporal del incidente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-038"></a>
## TI-038 · Buscar incidentes por texto

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Resultados pertenecen al alcance autorizado y muestran coincidencias

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Consulta impresora R3 y sede Caracas

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Consulta impresora R3 y sede Caracas
3. Ejecuta «Buscar incidentes por texto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Resultados pertenecen al alcance autorizado y muestran coincidencias

**Alternativas y errores:**
- Consulta vacía: usar listado paginado sin búsqueda global irrestricta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Resultados pertenecen al alcance autorizado y muestran coincidencias
- Se reproduce y controla esta alternativa: Consulta vacía: usar listado paginado sin búsqueda global irrestricta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de buscar incidentes por texto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-039"></a>
## TI-039 · Guardar una vista de trabajo

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La vista vuelve a aplicar el filtro al abrirla

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Filtro mis tickets urgentes sin respuesta

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtro mis tickets urgentes sin respuesta
3. Ejecuta «Guardar una vista de trabajo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La vista vuelve a aplicar el filtro al abrirla

**Alternativas y errores:**
- Campo de filtro eliminado: mostrar ajuste necesario.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La vista vuelve a aplicar el filtro al abrirla
- Se reproduce y controla esta alternativa: Campo de filtro eliminado: mostrar ajuste necesario.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de guardar una vista de trabajo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-040"></a>
## TI-040 · Atender cola de incidentes sin asignar

**Área:** incidents · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se prioriza una cola con antigüedad y SLA visibles

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Cinco tickets nuevos carecen de responsable

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Cinco tickets nuevos carecen de responsable
3. Ejecuta «Atender cola de incidentes sin asignar»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se prioriza una cola con antigüedad y SLA visibles

**Alternativas y errores:**
- Ticket asignado mientras se abre: actualizar la cola.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se prioriza una cola con antigüedad y SLA visibles
- Se reproduce y controla esta alternativa: Ticket asignado mientras se abre: actualizar la cola.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de atender cola de incidentes sin asignar; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-041"></a>
## TI-041 · Escalar un incidente a especialista

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se registra motivo y equipo de escalación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Diagnóstico de red requiere equipo de Seguridad

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Diagnóstico de red requiere equipo de Seguridad
3. Ejecuta «Escalar un incidente a especialista»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra motivo y equipo de escalación

**Alternativas y errores:**
- Sin equipo disponible: mantener propietario y alertar responsable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra motivo y equipo de escalación
- Se reproduce y controla esta alternativa: Sin equipo disponible: mantener propietario y alertar responsable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de escalar un incidente a especialista; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-042"></a>
## TI-042 · Registrar tiempo de atención

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Duración queda atribuida al agente y actividad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro dedica 35 minutos al ticket 42

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro dedica 35 minutos al ticket 42
3. Ejecuta «Registrar tiempo de atención»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Duración queda atribuida al agente y actividad

**Alternativas y errores:**
- Tiempo negativo o solapado prohibido: rechazar registro.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Duración queda atribuida al agente y actividad
- Se reproduce y controla esta alternativa: Tiempo negativo o solapado prohibido: rechazar registro.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de registrar tiempo de atención; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-043"></a>
## TI-043 · Crear borrador sin conexión

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Borrador persistido bajo su identidad y organización

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana redacta incidente desde una zona sin cobertura

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana redacta incidente desde una zona sin cobertura
3. Ejecuta «Crear borrador sin conexión»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Borrador persistido bajo su identidad y organización

**Alternativas y errores:**
- Cambio de usuario: ocultar borradores de identidad anterior.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Borrador persistido bajo su identidad y organización
- Se reproduce y controla esta alternativa: Cambio de usuario: ocultar borradores de identidad anterior.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de crear borrador sin conexión; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-044"></a>
## TI-044 · Sincronizar borrador pendiente

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se crea exactamente un ticket y se marca sincronizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Vuelve conexión con un borrador y clave idempotente

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Vuelve conexión con un borrador y clave idempotente
3. Ejecuta «Sincronizar borrador pendiente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crea exactamente un ticket y se marca sincronizado

**Alternativas y errores:**
- Timeout posterior al alta: reintentar la misma clave sin duplicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crea exactamente un ticket y se marca sincronizado
- Se reproduce y controla esta alternativa: Timeout posterior al alta: reintentar la misma clave sin duplicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de sincronizar borrador pendiente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-045"></a>
## TI-045 · Descartar un borrador offline

**Área:** incidents · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se elimina solo ese borrador y su sincronización pendiente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Reporte local quedó obsoleto antes de enviarse

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Reporte local quedó obsoleto antes de enviarse
3. Ejecuta «Descartar un borrador offline»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se elimina solo ese borrador y su sincronización pendiente

**Alternativas y errores:**
- Sincronización ya en curso: coordinar cancelación y mostrar resultado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se elimina solo ese borrador y su sincronización pendiente
- Se reproduce y controla esta alternativa: Sincronización ya en curso: coordinar cancelación y mostrar resultado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de descartar un borrador offline; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-046"></a>
## TI-046 · Recibir actualización en tiempo real

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La vista refleja nueva asignación sin recargar manualmente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Otro agente asigna el ticket abierto por Pedro

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Otro agente asigna el ticket abierto por Pedro
3. Ejecuta «Recibir actualización en tiempo real»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La vista refleja nueva asignación sin recargar manualmente

**Alternativas y errores:**
- WebSocket desconectado: recuperar estado desde API al reconectar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La vista refleja nueva asignación sin recargar manualmente
- Se reproduce y controla esta alternativa: WebSocket desconectado: recuperar estado desde API al reconectar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de recibir actualización en tiempo real; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-047"></a>
## TI-047 · Aplicar cambios masivos de cola

**Área:** incidents · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Cada cambio autorizado informa éxito o conflicto individual

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tres tickets de impresoras pasan a Infraestructura

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tres tickets de impresoras pasan a Infraestructura
3. Ejecuta «Aplicar cambios masivos de cola»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Cada cambio autorizado informa éxito o conflicto individual

**Alternativas y errores:**
- Ticket modificado simultáneamente: no sobrescribirlo silenciosamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Cada cambio autorizado informa éxito o conflicto individual
- Se reproduce y controla esta alternativa: Ticket modificado simultáneamente: no sobrescribirlo silenciosamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de aplicar cambios masivos de cola; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-048"></a>
## TI-048 · Consultar incidentes relacionados con un activo

**Área:** incidents · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se listan incidentes autorizados del mismo activo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Servidor PACS-01 presenta varias interrupciones

**Flujo principal:**
1. El actor accede a incidentes y seguimiento dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Servidor PACS-01 presenta varias interrupciones
3. Ejecuta «Consultar incidentes relacionados con un activo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se listan incidentes autorizados del mismo activo

**Alternativas y errores:**
- Activo retirado: conservar acceso histórico con permisos vigentes.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se listan incidentes autorizados del mismo activo
- Se reproduce y controla esta alternativa: Activo retirado: conservar acceso histórico con permisos vigentes.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Incidentes y seguimiento: operaciones de consultar incidentes relacionados con un activo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-049"></a>
## TI-049 · Publicar un servicio solicitables

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Aparece en el catálogo con responsable y plazo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Servicio Alta de correo para nuevos empleados

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Servicio Alta de correo para nuevos empleados
3. Ejecuta «Publicar un servicio solicitables»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Aparece en el catálogo con responsable y plazo

**Alternativas y errores:**
- Falta responsable: impedir publicación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Aparece en el catálogo con responsable y plazo
- Se reproduce y controla esta alternativa: Falta responsable: impedir publicación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de publicar un servicio solicitables; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-050"></a>
## TI-050 · Definir formulario de una solicitud

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** El formulario solicita los datos necesarios por tipo de servicio

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Alta de correo requiere nombre, área y fecha

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Alta de correo requiere nombre, área y fecha
3. Ejecuta «Definir formulario de una solicitud»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El formulario solicita los datos necesarios por tipo de servicio

**Alternativas y errores:**
- Campo incompatible con datos existentes: conservar versión anterior.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El formulario solicita los datos necesarios por tipo de servicio
- Se reproduce y controla esta alternativa: Campo incompatible con datos existentes: conservar versión anterior.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de definir formulario de una solicitud; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-051"></a>
## TI-051 · Solicitar acceso a una aplicación

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra una solicitud con alcance de acceso explícito

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana necesita acceso de consulta a PACS

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana necesita acceso de consulta a PACS
3. Ejecuta «Solicitar acceso a una aplicación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra una solicitud con alcance de acceso explícito

**Alternativas y errores:**
- Permiso privilegiado sin justificación: pedir justificación obligatoria.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra una solicitud con alcance de acceso explícito
- Se reproduce y controla esta alternativa: Permiso privilegiado sin justificación: pedir justificación obligatoria.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar acceso a una aplicación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-052"></a>
## TI-052 · Solicitar equipo de trabajo

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Solicitud informa modelo requerido y fecha objetivo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Nueva enfermera necesita portátil para el lunes

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Nueva enfermera necesita portátil para el lunes
3. Ejecuta «Solicitar equipo de trabajo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Solicitud informa modelo requerido y fecha objetivo

**Alternativas y errores:**
- Equipo sin disponibilidad: mostrar espera sin prometer entrega.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Solicitud informa modelo requerido y fecha objetivo
- Se reproduce y controla esta alternativa: Equipo sin disponibilidad: mostrar espera sin prometer entrega.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar equipo de trabajo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-053"></a>
## TI-053 · Solicitar instalación de software

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se crea tarea con software y equipo destino

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana necesita visor DICOM autorizado

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana necesita visor DICOM autorizado
3. Ejecuta «Solicitar instalación de software»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crea tarea con software y equipo destino

**Alternativas y errores:**
- Software fuera del catálogo: enviar a revisión autorizada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crea tarea con software y equipo destino
- Se reproduce y controla esta alternativa: Software fuera del catálogo: enviar a revisión autorizada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar instalación de software; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-054"></a>
## TI-054 · Solicitar incorporación de empleado

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se agrupan tareas de cuenta, equipo y accesos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Luis ingresa el 2 de noviembre a Laboratorio

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Luis ingresa el 2 de noviembre a Laboratorio
3. Ejecuta «Solicitar incorporación de empleado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se agrupan tareas de cuenta, equipo y accesos

**Alternativas y errores:**
- Fecha anterior al día actual: pedir confirmación de regularización.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se agrupan tareas de cuenta, equipo y accesos
- Se reproduce y controla esta alternativa: Fecha anterior al día actual: pedir confirmación de regularización.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar incorporación de empleado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-055"></a>
## TI-055 · Solicitar baja de empleado

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se coordinan revocación de accesos y devolución de activos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro finaliza contrato el 30 de noviembre

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro finaliza contrato el 30 de noviembre
3. Ejecuta «Solicitar baja de empleado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se coordinan revocación de accesos y devolución de activos

**Alternativas y errores:**
- Identidad ambigua: impedir una baja sobre persona incorrecta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se coordinan revocación de accesos y devolución de activos
- Se reproduce y controla esta alternativa: Identidad ambigua: impedir una baja sobre persona incorrecta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar baja de empleado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-056"></a>
## TI-056 · Solicitar traslado de puesto

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se revisan accesos, activos y pertenencias del nuevo puesto

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana cambia de Radiología a Laboratorio

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana cambia de Radiología a Laboratorio
3. Ejecuta «Solicitar traslado de puesto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se revisan accesos, activos y pertenencias del nuevo puesto

**Alternativas y errores:**
- Acceso incompatible con nueva función: retirarlo según aprobación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se revisan accesos, activos y pertenencias del nuevo puesto
- Se reproduce y controla esta alternativa: Acceso incompatible con nueva función: retirarlo según aprobación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar traslado de puesto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-057"></a>
## TI-057 · Solicitar ampliación de almacenamiento

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Solicitud indica capacidad, motivo y activo afectado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Servidor de imágenes requiere 200 GB adicionales

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Servidor de imágenes requiere 200 GB adicionales
3. Ejecuta «Solicitar ampliación de almacenamiento»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Solicitud indica capacidad, motivo y activo afectado

**Alternativas y errores:**
- Supera cuota aprobada: requerir aprobación de capacidad.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Solicitud indica capacidad, motivo y activo afectado
- Se reproduce y controla esta alternativa: Supera cuota aprobada: requerir aprobación de capacidad.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar ampliación de almacenamiento; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-058"></a>
## TI-058 · Solicitar acceso temporal

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** La solicitud incorpora una fecha de caducidad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Consultor necesita VPN durante cinco días

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Consultor necesita VPN durante cinco días
3. Ejecuta «Solicitar acceso temporal»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La solicitud incorpora una fecha de caducidad

**Alternativas y errores:**
- Fecha final anterior a inicial: rechazar periodo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La solicitud incorpora una fecha de caducidad
- Se reproduce y controla esta alternativa: Fecha final anterior a inicial: rechazar periodo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar acceso temporal; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-059"></a>
## TI-059 · Solicitar devolución de un activo

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se abre una tarea de recepción y verificación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana devuelve portátil AUR-021 al almacén

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana devuelve portátil AUR-021 al almacén
3. Ejecuta «Solicitar devolución de un activo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se abre una tarea de recepción y verificación

**Alternativas y errores:**
- Activo asignado a otra persona: pedir autorización del custodio.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se abre una tarea de recepción y verificación
- Se reproduce y controla esta alternativa: Activo asignado a otra persona: pedir autorización del custodio.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de solicitar devolución de un activo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-060"></a>
## TI-060 · Consultar estado de una solicitud

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se muestra progreso público por etapas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Alta de correo tiene dos tareas completas de tres

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Alta de correo tiene dos tareas completas de tres
3. Ejecuta «Consultar estado de una solicitud»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestra progreso público por etapas

**Alternativas y errores:**
- Tareas internas sensibles: resumir sin exponer instrucciones privadas.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestra progreso público por etapas
- Se reproduce y controla esta alternativa: Tareas internas sensibles: resumir sin exponer instrucciones privadas.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de consultar estado de una solicitud; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-061"></a>
## TI-061 · Completar tarea de provisión

**Área:** service-catalog · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Quedan resultado y evidencia vinculados a la tarea

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Cuenta de correo creada para Luis

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Cuenta de correo creada para Luis
3. Ejecuta «Completar tarea de provisión»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Quedan resultado y evidencia vinculados a la tarea

**Alternativas y errores:**
- Dependencia previa sin completar: bloquear cierre prematuro.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Quedan resultado y evidencia vinculados a la tarea
- Se reproduce y controla esta alternativa: Dependencia previa sin completar: bloquear cierre prematuro.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de completar tarea de provisión; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-062"></a>
## TI-062 · Entregar recurso al solicitante

**Área:** service-catalog · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Entrega queda registrada con receptor y fecha

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Portátil AUR-021 queda preparado para Luis

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Portátil AUR-021 queda preparado para Luis
3. Ejecuta «Entregar recurso al solicitante»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Entrega queda registrada con receptor y fecha

**Alternativas y errores:**
- Receptor rechaza equipo: mantener entrega pendiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Entrega queda registrada con receptor y fecha
- Se reproduce y controla esta alternativa: Receptor rechaza equipo: mantener entrega pendiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de entregar recurso al solicitante; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-063"></a>
## TI-063 · Confirmar recepción del servicio

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** La solicitud conserva aceptación del destinatario

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Luis verifica correo y portátil entregados

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Luis verifica correo y portátil entregados
3. Ejecuta «Confirmar recepción del servicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La solicitud conserva aceptación del destinatario

**Alternativas y errores:**
- Servicio incompleto: permitir observaciones y devolución a atención.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La solicitud conserva aceptación del destinatario
- Se reproduce y controla esta alternativa: Servicio incompleto: permitir observaciones y devolución a atención.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de confirmar recepción del servicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-064"></a>
## TI-064 · Programar ejecución de solicitud

**Área:** service-catalog · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Solicitante conoce ventana y técnico asignado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Instalación de visor será viernes a las 17:00

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Instalación de visor será viernes a las 17:00
3. Ejecuta «Programar ejecución de solicitud»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Solicitante conoce ventana y técnico asignado

**Alternativas y errores:**
- Ventana fuera del horario autorizado: pedir otra franja.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Solicitante conoce ventana y técnico asignado
- Se reproduce y controla esta alternativa: Ventana fuera del horario autorizado: pedir otra franja.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de programar ejecución de solicitud; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-065"></a>
## TI-065 · Cancelar solicitud antes de ejecución

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se cancelan tareas pendientes sin ocultar trabajo ya realizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. La incorporación de Luis fue cancelada

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: La incorporación de Luis fue cancelada
3. Ejecuta «Cancelar solicitud antes de ejecución»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se cancelan tareas pendientes sin ocultar trabajo ya realizado

**Alternativas y errores:**
- Provisión ya iniciada: derivar al responsable para compensación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se cancelan tareas pendientes sin ocultar trabajo ya realizado
- Se reproduce y controla esta alternativa: Provisión ya iniciada: derivar al responsable para compensación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de cancelar solicitud antes de ejecución; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-066"></a>
## TI-066 · Duplicar solicitud como nueva petición

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se reutilizan campos permitidos con nuevo destinatario y referencia

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Otra enfermera requiere el mismo tipo de equipo

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Otra enfermera requiere el mismo tipo de equipo
3. Ejecuta «Duplicar solicitud como nueva petición»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se reutilizan campos permitidos con nuevo destinatario y referencia

**Alternativas y errores:**
- Datos sensibles del original: no copiarlos al nuevo destinatario.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se reutilizan campos permitidos con nuevo destinatario y referencia
- Se reproduce y controla esta alternativa: Datos sensibles del original: no copiarlos al nuevo destinatario.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de duplicar solicitud como nueva petición; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-067"></a>
## TI-067 · Retirar servicio del catálogo

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se impiden nuevas solicitudes conservando las abiertas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Visor antiguo deja de distribuirse

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Visor antiguo deja de distribuirse
3. Ejecuta «Retirar servicio del catálogo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se impiden nuevas solicitudes conservando las abiertas

**Alternativas y errores:**
- Solicitud en curso: permitir completar o migrar con motivo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se impiden nuevas solicitudes conservando las abiertas
- Se reproduce y controla esta alternativa: Solicitud en curso: permitir completar o migrar con motivo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de retirar servicio del catálogo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-068"></a>
## TI-068 · Versionar condiciones de un servicio

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Nuevas solicitudes usan versión nueva; antiguas conservan condiciones

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Entrega de portátil pasa de cinco a siete días

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Entrega de portátil pasa de cinco a siete días
3. Ejecuta «Versionar condiciones de un servicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nuevas solicitudes usan versión nueva; antiguas conservan condiciones

**Alternativas y errores:**
- Edición concurrente: informar conflicto de versión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nuevas solicitudes usan versión nueva; antiguas conservan condiciones
- Se reproduce y controla esta alternativa: Edición concurrente: informar conflicto de versión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de versionar condiciones de un servicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-069"></a>
## TI-069 · Restringir catálogo por departamento

**Área:** service-catalog · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Otros departamentos no ven ni solicitan el servicio restringido

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Acceso PACS disponible solo para áreas clínicas

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Acceso PACS disponible solo para áreas clínicas
3. Ejecuta «Restringir catálogo por departamento»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Otros departamentos no ven ni solicitan el servicio restringido

**Alternativas y errores:**
- URL directa desde área ajena: denegar solicitud.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Otros departamentos no ven ni solicitan el servicio restringido
- Se reproduce y controla esta alternativa: URL directa desde área ajena: denegar solicitud.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de restringir catálogo por departamento; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-070"></a>
## TI-070 · Consultar coste estimado de solicitud

**Área:** service-catalog · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** El solicitante ve estimación antes de enviarla

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Licencia anual cuesta 120 unidades monetarias

**Flujo principal:**
1. El actor accede a solicitudes y catálogo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Licencia anual cuesta 120 unidades monetarias
3. Ejecuta «Consultar coste estimado de solicitud»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El solicitante ve estimación antes de enviarla

**Alternativas y errores:**
- Precio ausente: mostrar pendiente de estimación, no cero ficticio.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El solicitante ve estimación antes de enviarla
- Se reproduce y controla esta alternativa: Precio ausente: mostrar pendiente de estimación, no cero ficticio.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Solicitudes y catálogo: operaciones de consultar coste estimado de solicitud; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-071"></a>
## TI-071 · Redactar artículo de solución

**Área:** knowledge · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se guarda borrador con pasos y categoría

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Procedimiento para restablecer cola de impresión

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Procedimiento para restablecer cola de impresión
3. Ejecuta «Redactar artículo de solución»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se guarda borrador con pasos y categoría

**Alternativas y errores:**
- Contenido inválido: conservar borrador sin publicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se guarda borrador con pasos y categoría
- Se reproduce y controla esta alternativa: Contenido inválido: conservar borrador sin publicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de redactar artículo de solución; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-072"></a>
## TI-072 · Publicar artículo revisado

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La versión aprobada aparece en autoservicio

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Borrador revisado sobre impresora R3

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Borrador revisado sobre impresora R3
3. Ejecuta «Publicar artículo revisado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La versión aprobada aparece en autoservicio

**Alternativas y errores:**
- Revisión pendiente: impedir publicación directa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La versión aprobada aparece en autoservicio
- Se reproduce y controla esta alternativa: Revisión pendiente: impedir publicación directa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de publicar artículo revisado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-073"></a>
## TI-073 · Buscar solución antes de reportar

**Área:** knowledge · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Encuentra artículos visibles relacionados con el error

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana busca error E204 de impresión

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana busca error E204 de impresión
3. Ejecuta «Buscar solución antes de reportar»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Encuentra artículos visibles relacionados con el error

**Alternativas y errores:**
- Sin coincidencias: ofrecer crear ticket con búsqueda conservada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Encuentra artículos visibles relacionados con el error
- Se reproduce y controla esta alternativa: Sin coincidencias: ofrecer crear ticket con búsqueda conservada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de buscar solución antes de reportar; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-074"></a>
## TI-074 · Insertar solución de conocimiento en respuesta

**Área:** knowledge · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Respuesta enlaza la versión consultada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Artículo aprobado explica reinicio de cola

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Artículo aprobado explica reinicio de cola
3. Ejecuta «Insertar solución de conocimiento en respuesta»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Respuesta enlaza la versión consultada

**Alternativas y errores:**
- Artículo retirado: advertir y pedir contenido vigente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Respuesta enlaza la versión consultada
- Se reproduce y controla esta alternativa: Artículo retirado: advertir y pedir contenido vigente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de insertar solución de conocimiento en respuesta; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-075"></a>
## TI-075 · Valorar utilidad de un artículo

**Área:** knowledge · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se registra valoración de utilidad de esa versión

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana resolvió la impresión siguiendo el artículo

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana resolvió la impresión siguiendo el artículo
3. Ejecuta «Valorar utilidad de un artículo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra valoración de utilidad de esa versión

**Alternativas y errores:**
- Envío repetido: actualizar una valoración por identidad.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra valoración de utilidad de esa versión
- Se reproduce y controla esta alternativa: Envío repetido: actualizar una valoración por identidad.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de valorar utilidad de un artículo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-076"></a>
## TI-076 · Reportar información obsoleta

**Área:** knowledge · **Actor:** Solicitante · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El autor recibe observación vinculada al artículo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ruta de configuración ya no existe en Windows

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ruta de configuración ya no existe en Windows
3. Ejecuta «Reportar información obsoleta»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El autor recibe observación vinculada al artículo

**Alternativas y errores:**
- Comentario ofensivo: aplicar reglas de moderación visibles.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El autor recibe observación vinculada al artículo
- Se reproduce y controla esta alternativa: Comentario ofensivo: aplicar reglas de moderación visibles.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de reportar información obsoleta; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-077"></a>
## TI-077 · Editar una nueva revisión

**Área:** knowledge · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Borrador nuevo no altera la publicación vigente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se actualiza procedimiento por cambio de controlador

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se actualiza procedimiento por cambio de controlador
3. Ejecuta «Editar una nueva revisión»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Borrador nuevo no altera la publicación vigente

**Alternativas y errores:**
- Edición concurrente: mostrar conflicto de versión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Borrador nuevo no altera la publicación vigente
- Se reproduce y controla esta alternativa: Edición concurrente: mostrar conflicto de versión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de editar una nueva revisión; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-078"></a>
## TI-078 · Comparar revisiones de un artículo

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se distinguen cambios antes de aprobar

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Revisión 3 cambia pasos de reinicio respecto a 2

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisión 3 cambia pasos de reinicio respecto a 2
3. Ejecuta «Comparar revisiones de un artículo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se distinguen cambios antes de aprobar

**Alternativas y errores:**
- Revisión inexistente: informar sin sustituir por otra.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se distinguen cambios antes de aprobar
- Se reproduce y controla esta alternativa: Revisión inexistente: informar sin sustituir por otra.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de comparar revisiones de un artículo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-079"></a>
## TI-079 · Restaurar revisión conocida

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se publica una nueva revisión basada en la 3 con trazabilidad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Revisión 4 introdujo instrucciones equivocadas

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisión 4 introdujo instrucciones equivocadas
3. Ejecuta «Restaurar revisión conocida»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se publica una nueva revisión basada en la 3 con trazabilidad

**Alternativas y errores:**
- Revisión incompatible con servicio actual: advertir antes de publicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se publica una nueva revisión basada en la 3 con trazabilidad
- Se reproduce y controla esta alternativa: Revisión incompatible con servicio actual: advertir antes de publicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de restaurar revisión conocida; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-080"></a>
## TI-080 · Restringir artículo a agentes

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Solo personal autorizado puede buscar o leer el artículo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Diagnóstico incluye detalles de red interna

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Diagnóstico incluye detalles de red interna
3. Ejecuta «Restringir artículo a agentes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Solo personal autorizado puede buscar o leer el artículo

**Alternativas y errores:**
- Enlace compartido con empleado: denegar contenido interno.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Solo personal autorizado puede buscar o leer el artículo
- Se reproduce y controla esta alternativa: Enlace compartido con empleado: denegar contenido interno.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de restringir artículo a agentes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-081"></a>
## TI-081 · Organizar categorías de conocimiento

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Artículos mantienen enlaces y aparecen en categoría correcta

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Impresión se agrupa bajo Puesto de trabajo

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Impresión se agrupa bajo Puesto de trabajo
3. Ejecuta «Organizar categorías de conocimiento»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Artículos mantienen enlaces y aparecen en categoría correcta

**Alternativas y errores:**
- Eliminar categoría con artículos: exigir reasignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Artículos mantienen enlaces y aparecen en categoría correcta
- Se reproduce y controla esta alternativa: Eliminar categoría con artículos: exigir reasignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de organizar categorías de conocimiento; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-082"></a>
## TI-082 · Mantener traducción de una guía

**Área:** knowledge · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ambas versiones conservan relación y estado de actualización

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Artículo español tiene versión inglesa revisada

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Artículo español tiene versión inglesa revisada
3. Ejecuta «Mantener traducción de una guía»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ambas versiones conservan relación y estado de actualización

**Alternativas y errores:**
- Original cambia: marcar traducción pendiente de revisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ambas versiones conservan relación y estado de actualización
- Se reproduce y controla esta alternativa: Original cambia: marcar traducción pendiente de revisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de mantener traducción de una guía; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-083"></a>
## TI-083 · Programar revisión de contenido

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Responsable recibe tarea de revisión a su vencimiento

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Guía VPN requiere revisión cada 90 días

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Guía VPN requiere revisión cada 90 días
3. Ejecuta «Programar revisión de contenido»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Responsable recibe tarea de revisión a su vencimiento

**Alternativas y errores:**
- Responsable inactivo: escalar al gestor de conocimiento.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Responsable recibe tarea de revisión a su vencimiento
- Se reproduce y controla esta alternativa: Responsable inactivo: escalar al gestor de conocimiento.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de programar revisión de contenido; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-084"></a>
## TI-084 · Retirar artículo peligroso

**Área:** knowledge · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se retira publicación conservando auditoría y referencia histórica

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Procedimiento permite borrar configuración necesaria

**Flujo principal:**
1. El actor accede a conocimiento y autoservicio dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Procedimiento permite borrar configuración necesaria
3. Ejecuta «Retirar artículo peligroso»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se retira publicación conservando auditoría y referencia histórica

**Alternativas y errores:**
- Ticket lo enlazaba: mostrar aviso de artículo retirado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se retira publicación conservando auditoría y referencia histórica
- Se reproduce y controla esta alternativa: Ticket lo enlazaba: mostrar aviso de artículo retirado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Conocimiento y autoservicio: operaciones de retirar artículo peligroso; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-085"></a>
## TI-085 · Definir objetivo de primera respuesta

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Política asigna plazo persistente a tickets nuevos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Prioridad alta exige respuesta en 30 minutos

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Prioridad alta exige respuesta en 30 minutos
3. Ejecuta «Definir objetivo de primera respuesta»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Política asigna plazo persistente a tickets nuevos

**Alternativas y errores:**
- Tiempo no positivo: rechazar la política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Política asigna plazo persistente a tickets nuevos
- Se reproduce y controla esta alternativa: Tiempo no positivo: rechazar la política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de definir objetivo de primera respuesta; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-086"></a>
## TI-086 · Definir objetivo de resolución

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se calcula vencimiento de resolución según calendario

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Incidente normal debe resolverse en ocho horas hábiles

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Incidente normal debe resolverse en ocho horas hábiles
3. Ejecuta «Definir objetivo de resolución»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se calcula vencimiento de resolución según calendario

**Alternativas y errores:**
- Calendario inexistente: impedir activar política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se calcula vencimiento de resolución según calendario
- Se reproduce y controla esta alternativa: Calendario inexistente: impedir activar política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de definir objetivo de resolución; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-087"></a>
## TI-087 · Crear horario de atención

**Área:** sla · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El cómputo excluye horas fuera de atención

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Equipo atiende lunes a viernes de 08:00 a 17:00

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Equipo atiende lunes a viernes de 08:00 a 17:00
3. Ejecuta «Crear horario de atención»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El cómputo excluye horas fuera de atención

**Alternativas y errores:**
- Franjas solapadas: solicitar corrección.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El cómputo excluye horas fuera de atención
- Se reproduce y controla esta alternativa: Franjas solapadas: solicitar corrección.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de crear horario de atención; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-088"></a>
## TI-088 · Registrar festivo de servicio

**Área:** sla · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Vencimientos nuevos excluyen el festivo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. 25 de diciembre queda fuera del calendario laboral

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: 25 de diciembre queda fuera del calendario laboral
3. Ejecuta «Registrar festivo de servicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Vencimientos nuevos excluyen el festivo

**Alternativas y errores:**
- Calendario compartido: mostrar servicios afectados antes de guardar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Vencimientos nuevos excluyen el festivo
- Se reproduce y controla esta alternativa: Calendario compartido: mostrar servicios afectados antes de guardar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de registrar festivo de servicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-089"></a>
## TI-089 · Aplicar SLA por prioridad

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La política seleccionada queda trazable en el ticket

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Incidente de impacto alto usa política Urgente

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Incidente de impacto alto usa política Urgente
3. Ejecuta «Aplicar SLA por prioridad»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La política seleccionada queda trazable en el ticket

**Alternativas y errores:**
- Dos reglas empatan: usar precedencia explícita configurada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La política seleccionada queda trazable en el ticket
- Se reproduce y controla esta alternativa: Dos reglas empatan: usar precedencia explícita configurada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de aplicar sla por prioridad; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-090"></a>
## TI-090 · Pausar SLA por espera autorizada

**Área:** sla · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Pausa queda fechada con motivo admitido

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se espera respuesta de Ana desde las 10:00

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se espera respuesta de Ana desde las 10:00
3. Ejecuta «Pausar SLA por espera autorizada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Pausa queda fechada con motivo admitido

**Alternativas y errores:**
- Estado no permite pausa: conservar el contador activo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Pausa queda fechada con motivo admitido
- Se reproduce y controla esta alternativa: Estado no permite pausa: conservar el contador activo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de pausar sla por espera autorizada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-091"></a>
## TI-091 · Reanudar SLA tras respuesta

**Área:** sla · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se reanuda cómputo sin contar pausa autorizada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana aporta información a las 11:15

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana aporta información a las 11:15
3. Ejecuta «Reanudar SLA tras respuesta»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se reanuda cómputo sin contar pausa autorizada

**Alternativas y errores:**
- Evento repetido: no descontar dos veces el intervalo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se reanuda cómputo sin contar pausa autorizada
- Se reproduce y controla esta alternativa: Evento repetido: no descontar dos veces el intervalo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de reanudar sla tras respuesta; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-092"></a>
## TI-092 · Consultar cuenta regresiva de SLA

**Área:** sla · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La vista muestra fecha y tiempo restante coherentes

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 vence en 25 minutos

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 vence en 25 minutos
3. Ejecuta «Consultar cuenta regresiva de SLA»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La vista muestra fecha y tiempo restante coherentes

**Alternativas y errores:**
- Reloj cliente desviado: tomar vencimiento del servidor.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La vista muestra fecha y tiempo restante coherentes
- Se reproduce y controla esta alternativa: Reloj cliente desviado: tomar vencimiento del servidor.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de consultar cuenta regresiva de sla; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-093"></a>
## TI-093 · Notificar proximidad de incumplimiento

**Área:** sla · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El responsable recibe una alerta deduplicada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Quedan diez minutos para primera respuesta

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Quedan diez minutos para primera respuesta
3. Ejecuta «Notificar proximidad de incumplimiento»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El responsable recibe una alerta deduplicada

**Alternativas y errores:**
- Canal temporalmente caído: conservar aviso para reintento.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El responsable recibe una alerta deduplicada
- Se reproduce y controla esta alternativa: Canal temporalmente caído: conservar aviso para reintento.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de notificar proximidad de incumplimiento; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-094"></a>
## TI-094 · Escalar incumplimiento a supervisor

**Área:** sla · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Supervisor recibe contexto y se registra escalación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 supera su plazo de respuesta

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 supera su plazo de respuesta
3. Ejecuta «Escalar incumplimiento a supervisor»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Supervisor recibe contexto y se registra escalación

**Alternativas y errores:**
- Worker reinicia: no repetir escalación ya registrada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Supervisor recibe contexto y se registra escalación
- Se reproduce y controla esta alternativa: Worker reinicia: no repetir escalación ya registrada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de escalar incumplimiento a supervisor; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-095"></a>
## TI-095 · Registrar primera respuesta efectiva

**Área:** sla · **Actor:** Agente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se fija una sola marca de primera respuesta

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro envía primera respuesta pública a Ana

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro envía primera respuesta pública a Ana
3. Ejecuta «Registrar primera respuesta efectiva»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se fija una sola marca de primera respuesta

**Alternativas y errores:**
- Nota interna previa: no contarla como respuesta pública.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se fija una sola marca de primera respuesta
- Se reproduce y controla esta alternativa: Nota interna previa: no contarla como respuesta pública.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de registrar primera respuesta efectiva; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-096"></a>
## TI-096 · Calcular duración de resolución

**Área:** sla · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El informe refleja duración según calendario y pausas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket abierto lunes y resuelto martes con pausa

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket abierto lunes y resuelto martes con pausa
3. Ejecuta «Calcular duración de resolución»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El informe refleja duración según calendario y pausas

**Alternativas y errores:**
- Datos anteriores sin marcas: mostrar desconocido, no duración cero.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El informe refleja duración según calendario y pausas
- Se reproduce y controla esta alternativa: Datos anteriores sin marcas: mostrar desconocido, no duración cero.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de calcular duración de resolución; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-097"></a>
## TI-097 · Recalcular objetivo por cambio autorizado

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Nuevo plazo y motivo quedan auditados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Prioridad cambia de normal a crítica por mayor impacto

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Prioridad cambia de normal a crítica por mayor impacto
3. Ejecuta «Recalcular objetivo por cambio autorizado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nuevo plazo y motivo quedan auditados

**Alternativas y errores:**
- Cambio retroactivo borra incumplimiento: conservar evidencia original.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nuevo plazo y motivo quedan auditados
- Se reproduce y controla esta alternativa: Cambio retroactivo borra incumplimiento: conservar evidencia original.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de recalcular objetivo por cambio autorizado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-098"></a>
## TI-098 · Versionar una política SLA

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Nuevos tickets usan versión nueva conservando historial

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Objetivo normal cambia de ocho a seis horas

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Objetivo normal cambia de ocho a seis horas
3. Ejecuta «Versionar una política SLA»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nuevos tickets usan versión nueva conservando historial

**Alternativas y errores:**
- Tickets abiertos: exigir opción explícita antes de recalcularlos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nuevos tickets usan versión nueva conservando historial
- Se reproduce y controla esta alternativa: Tickets abiertos: exigir opción explícita antes de recalcularlos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de versionar una política sla; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-099"></a>
## TI-099 · Excluir ticket del cómputo con motivo

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se excluye del indicador con justificación auditable

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket de prueba fue creado por error operativo

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket de prueba fue creado por error operativo
3. Ejecuta «Excluir ticket del cómputo con motivo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se excluye del indicador con justificación auditable

**Alternativas y errores:**
- Agente sin privilegio: impedir exclusión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se excluye del indicador con justificación auditable
- Se reproduce y controla esta alternativa: Agente sin privilegio: impedir exclusión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de excluir ticket del cómputo con motivo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-100"></a>
## TI-100 · Asignar calendario por sede

**Área:** sla · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Tickets usan calendario de la sede correspondiente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Sede Valencia tiene horario distinto a Caracas

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Sede Valencia tiene horario distinto a Caracas
3. Ejecuta «Asignar calendario por sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Tickets usan calendario de la sede correspondiente

**Alternativas y errores:**
- Sede sin horario: aplicar calendario organizativo documentado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Tickets usan calendario de la sede correspondiente
- Se reproduce y controla esta alternativa: Sede sin horario: aplicar calendario organizativo documentado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de asignar calendario por sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-101"></a>
## TI-101 · Detectar tickets sin política

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se presenta lista de cobertura pendiente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Cinco tickets heredados carecen de SLA

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Cinco tickets heredados carecen de SLA
3. Ejecuta «Detectar tickets sin política»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se presenta lista de cobertura pendiente

**Alternativas y errores:**
- Sin fechas históricas: no inventar vencimientos anteriores.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se presenta lista de cobertura pendiente
- Se reproduce y controla esta alternativa: Sin fechas históricas: no inventar vencimientos anteriores.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de detectar tickets sin política; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-102"></a>
## TI-102 · Definir escalación por etapas

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Cada etapa genera un único registro de escalación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. A 100% alerta supervisor; a 150% alerta gerente

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: A 100% alerta supervisor; a 150% alerta gerente
3. Ejecuta «Definir escalación por etapas»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Cada etapa genera un único registro de escalación

**Alternativas y errores:**
- Misma persona en dos etapas: conservar etapas sin spam duplicado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Cada etapa genera un único registro de escalación
- Se reproduce y controla esta alternativa: Misma persona en dos etapas: conservar etapas sin spam duplicado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de definir escalación por etapas; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-103"></a>
## TI-103 · Auditar una pausa de SLA

**Área:** sla · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se muestran inicio, fin, causa y actor

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 estuvo pausado 75 minutos por espera

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 estuvo pausado 75 minutos por espera
3. Ejecuta «Auditar una pausa de SLA»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran inicio, fin, causa y actor

**Alternativas y errores:**
- Pausa sin cierre por caída: identificar inconsistencia para reparación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran inicio, fin, causa y actor
- Se reproduce y controla esta alternativa: Pausa sin cierre por caída: identificar inconsistencia para reparación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de auditar una pausa de sla; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-104"></a>
## TI-104 · Comparar cumplimiento por servicio

**Área:** sla · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se muestran numerador, denominador y periodos comparables

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Impresión y VPN durante octubre

**Flujo principal:**
1. El actor accede a sla, calendarios y escalaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Impresión y VPN durante octubre
3. Ejecuta «Comparar cumplimiento por servicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran numerador, denominador y periodos comparables

**Alternativas y errores:**
- Muestra vacía: indicar sin datos en vez de 100%.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran numerador, denominador y periodos comparables
- Se reproduce y controla esta alternativa: Muestra vacía: indicar sin datos en vez de 100%.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / SLA, calendarios y escalaciones: operaciones de comparar cumplimiento por servicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-105"></a>
## TI-105 · Registrar activo físico

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se crea activo único con clase y ubicación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Portátil AUR-021 con serie SN-8942

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Portátil AUR-021 con serie SN-8942
3. Ejecuta «Registrar activo físico»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crea activo único con clase y ubicación

**Alternativas y errores:**
- Serie duplicada: resolver coincidencia antes de crear.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crea activo único con clase y ubicación
- Se reproduce y controla esta alternativa: Serie duplicada: resolver coincidencia antes de crear.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de registrar activo físico; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-106"></a>
## TI-106 · Importar inventario de equipos

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se valida por fila y se informa resultado sin duplicados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. CSV con 40 portátiles y números de serie

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: CSV con 40 portátiles y números de serie
3. Ejecuta «Importar inventario de equipos»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se valida por fila y se informa resultado sin duplicados

**Alternativas y errores:**
- Fila inválida: informar causa sin descartar filas válidas silenciosamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se valida por fila y se informa resultado sin duplicados
- Se reproduce y controla esta alternativa: Fila inválida: informar causa sin descartar filas válidas silenciosamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de importar inventario de equipos; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-107"></a>
## TI-107 · Asignar activo a empleado

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Custodia registra persona, fecha y condición

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-021 se entrega a Ana

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-021 se entrega a Ana
3. Ejecuta «Asignar activo a empleado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Custodia registra persona, fecha y condición

**Alternativas y errores:**
- Activo ya asignado: exigir transferencia o devolución.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Custodia registra persona, fecha y condición
- Se reproduce y controla esta alternativa: Activo ya asignado: exigir transferencia o devolución.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de asignar activo a empleado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-108"></a>
## TI-108 · Transferir custodia de activo

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Historial conserva ambas custodias y fecha de cambio

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana entrega AUR-021 a Luis

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana entrega AUR-021 a Luis
3. Ejecuta «Transferir custodia de activo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Historial conserva ambas custodias y fecha de cambio

**Alternativas y errores:**
- Transferencia concurrente: impedir dos custodios activos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Historial conserva ambas custodias y fecha de cambio
- Se reproduce y controla esta alternativa: Transferencia concurrente: impedir dos custodios activos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de transferir custodia de activo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-109"></a>
## TI-109 · Registrar devolución a almacén

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Activo pasa a disponible tras recepción registrada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-021 retorna con cargador y buen estado

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-021 retorna con cargador y buen estado
3. Ejecuta «Registrar devolución a almacén»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Activo pasa a disponible tras recepción registrada

**Alternativas y errores:**
- Accesorio faltante: dejar observación y estado de revisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Activo pasa a disponible tras recepción registrada
- Se reproduce y controla esta alternativa: Accesorio faltante: dejar observación y estado de revisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de registrar devolución a almacén; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-110"></a>
## TI-110 · Reservar equipo para incorporación

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Otro proceso no puede asignarlo durante la reserva

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-025 reservado para Luis el lunes

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-025 reservado para Luis el lunes
3. Ejecuta «Reservar equipo para incorporación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Otro proceso no puede asignarlo durante la reserva

**Alternativas y errores:**
- Reserva vencida: liberar con evento auditable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Otro proceso no puede asignarlo durante la reserva
- Se reproduce y controla esta alternativa: Reserva vencida: liberar con evento auditable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de reservar equipo para incorporación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-111"></a>
## TI-111 · Registrar garantía de fabricante

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Ficha muestra proveedor y vencimiento de garantía

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-021 cubierto hasta 30 de junio

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-021 cubierto hasta 30 de junio
3. Ejecuta «Registrar garantía de fabricante»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ficha muestra proveedor y vencimiento de garantía

**Alternativas y errores:**
- Fecha anterior a compra: solicitar corrección.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ficha muestra proveedor y vencimiento de garantía
- Se reproduce y controla esta alternativa: Fecha anterior a compra: solicitar corrección.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de registrar garantía de fabricante; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-112"></a>
## TI-112 · Programar mantenimiento preventivo

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se crean actividades con próximo vencimiento

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Servidor PACS-01 requiere limpieza trimestral

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Servidor PACS-01 requiere limpieza trimestral
3. Ejecuta «Programar mantenimiento preventivo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crean actividades con próximo vencimiento

**Alternativas y errores:**
- Equipo retirado antes del vencimiento: cancelar futuras actividades.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crean actividades con próximo vencimiento
- Se reproduce y controla esta alternativa: Equipo retirado antes del vencimiento: cancelar futuras actividades.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de programar mantenimiento preventivo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-113"></a>
## TI-113 · Registrar reparación externa

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se conserva custodia externa y número de servicio

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-021 se envía a proveedor por batería

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-021 se envía a proveedor por batería
3. Ejecuta «Registrar reparación externa»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se conserva custodia externa y número de servicio

**Alternativas y errores:**
- Equipo con datos sensibles: requerir tratamiento autorizado previo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se conserva custodia externa y número de servicio
- Se reproduce y controla esta alternativa: Equipo con datos sensibles: requerir tratamiento autorizado previo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de registrar reparación externa; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-114"></a>
## TI-114 · Retirar activo de servicio

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Activo deja de ser asignable y conserva historial

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-009 es irreparable

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-009 es irreparable
3. Ejecuta «Retirar activo de servicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Activo deja de ser asignable y conserva historial

**Alternativas y errores:**
- Activo enlazado a servicio crítico: exigir revisión de dependencias.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Activo deja de ser asignable y conserva historial
- Se reproduce y controla esta alternativa: Activo enlazado a servicio crítico: exigir revisión de dependencias.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de retirar activo de servicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-115"></a>
## TI-115 · Documentar disposición final

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra destino y evidencia del proceso

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-009 se entrega a reciclaje con certificado

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-009 se entrega a reciclaje con certificado
3. Ejecuta «Documentar disposición final»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra destino y evidencia del proceso

**Alternativas y errores:**
- Falta autorización de disposición: no completar el estado final.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra destino y evidencia del proceso
- Se reproduce y controla esta alternativa: Falta autorización de disposición: no completar el estado final.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de documentar disposición final; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-116"></a>
## TI-116 · Registrar licencia de software

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Inventario muestra capacidad y vigencia

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Diez licencias del visor DICOM vencen en enero

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Diez licencias del visor DICOM vencen en enero
3. Ejecuta «Registrar licencia de software»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Inventario muestra capacidad y vigencia

**Alternativas y errores:**
- Cantidad negativa o contrato repetido: rechazar alta inconsistente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Inventario muestra capacidad y vigencia
- Se reproduce y controla esta alternativa: Cantidad negativa o contrato repetido: rechazar alta inconsistente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de registrar licencia de software; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-117"></a>
## TI-117 · Asignar licencia a un equipo

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se consume una plaza disponible

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Visor DICOM se asigna a AUR-021

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Visor DICOM se asigna a AUR-021
3. Ejecuta «Asignar licencia a un equipo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se consume una plaza disponible

**Alternativas y errores:**
- Sin plazas: impedir sobreasignación y mostrar déficit.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se consume una plaza disponible
- Se reproduce y controla esta alternativa: Sin plazas: impedir sobreasignación y mostrar déficit.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de asignar licencia a un equipo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-118"></a>
## TI-118 · Recuperar licencia no utilizada

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** La plaza vuelve a disponible con motivo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-021 retirado libera su licencia DICOM

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-021 retirado libera su licencia DICOM
3. Ejecuta «Recuperar licencia no utilizada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La plaza vuelve a disponible con motivo

**Alternativas y errores:**
- Software aún activo confirmado: requerir retirada verificada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La plaza vuelve a disponible con motivo
- Se reproduce y controla esta alternativa: Software aún activo confirmado: requerir retirada verificada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de recuperar licencia no utilizada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-119"></a>
## TI-119 · Consultar software por activo

**Área:** assets · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se presenta inventario y fecha de observación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Diagnóstico de AUR-021 requiere versión de visor

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Diagnóstico de AUR-021 requiere versión de visor
3. Ejecuta «Consultar software por activo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se presenta inventario y fecha de observación

**Alternativas y errores:**
- Observación antigua: indicar antigüedad del dato.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se presenta inventario y fecha de observación
- Se reproduce y controla esta alternativa: Observación antigua: indicar antigüedad del dato.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de consultar software por activo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-120"></a>
## TI-120 · Crear elemento de configuración lógico

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra CI con responsable y criticidad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Servicio Historia clínica depende de aplicación PACS

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Servicio Historia clínica depende de aplicación PACS
3. Ejecuta «Crear elemento de configuración lógico»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra CI con responsable y criticidad

**Alternativas y errores:**
- CI homónimo: pedir identificador que evite ambigüedad.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra CI con responsable y criticidad
- Se reproduce y controla esta alternativa: CI homónimo: pedir identificador que evite ambigüedad.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de crear elemento de configuración lógico; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-121"></a>
## TI-121 · Relacionar dependencias de servicio

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Mapa muestra dirección de dependencia verificable

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. PACS depende de base de datos DB-01

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: PACS depende de base de datos DB-01
3. Ejecuta «Relacionar dependencias de servicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Mapa muestra dirección de dependencia verificable

**Alternativas y errores:**
- Relación consigo mismo: rechazar vínculo inválido.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Mapa muestra dirección de dependencia verificable
- Se reproduce y controla esta alternativa: Relación consigo mismo: rechazar vínculo inválido.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de relacionar dependencias de servicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-122"></a>
## TI-122 · Analizar impacto de caída de activo

**Área:** assets · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se muestran servicios y CIs dependientes autorizados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. DB-01 deja de responder

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: DB-01 deja de responder
3. Ejecuta «Analizar impacto de caída de activo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran servicios y CIs dependientes autorizados

**Alternativas y errores:**
- Dependencias incompletas: señalar cobertura conocida.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran servicios y CIs dependientes autorizados
- Se reproduce y controla esta alternativa: Dependencias incompletas: señalar cobertura conocida.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de analizar impacto de caída de activo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-123"></a>
## TI-123 · Vincular activo a incidente

**Área:** assets · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Historial del activo enlaza el incidente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 corresponde a AUR-021

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 corresponde a AUR-021
3. Ejecuta «Vincular activo a incidente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Historial del activo enlaza el incidente

**Alternativas y errores:**
- Activo de otra organización: denegar vinculación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Historial del activo enlaza el incidente
- Se reproduce y controla esta alternativa: Activo de otra organización: denegar vinculación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de vincular activo a incidente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-124"></a>
## TI-124 · Conciliar inventario descubierto

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se actualiza observación sin duplicar el activo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Escaneo detecta SN-8942 con nueva dirección IP

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Escaneo detecta SN-8942 con nueva dirección IP
3. Ejecuta «Conciliar inventario descubierto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se actualiza observación sin duplicar el activo

**Alternativas y errores:**
- Coincidencias ambiguas: pedir resolución manual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se actualiza observación sin duplicar el activo
- Se reproduce y controla esta alternativa: Coincidencias ambiguas: pedir resolución manual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de conciliar inventario descubierto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-125"></a>
## TI-125 · Detectar activo sin custodio

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Aparece en lista de inconsistencias del inventario

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-030 figura activo sin asignación ni almacén

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-030 figura activo sin asignación ni almacén
3. Ejecuta «Detectar activo sin custodio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Aparece en lista de inconsistencias del inventario

**Alternativas y errores:**
- Activo en reparación: reconocer custodia externa válida.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Aparece en lista de inconsistencias del inventario
- Se reproduce y controla esta alternativa: Activo en reparación: reconocer custodia externa válida.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de detectar activo sin custodio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-126"></a>
## TI-126 · Auditar inventario de una sede

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registran encontrados, ausentes y diferencias

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se verifican 25 equipos de Valencia

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se verifican 25 equipos de Valencia
3. Ejecuta «Auditar inventario de una sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registran encontrados, ausentes y diferencias

**Alternativas y errores:**
- Equipo encontrado en otra sede: abrir discrepancia, no moverlo en silencio.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registran encontrados, ausentes y diferencias
- Se reproduce y controla esta alternativa: Equipo encontrado en otra sede: abrir discrepancia, no moverlo en silencio.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de auditar inventario de una sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-127"></a>
## TI-127 · Consultar coste de ciclo de vida

**Área:** assets · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se muestra desglose de coste conocido

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. AUR-021 costó 900 y acumuló 120 en reparaciones

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: AUR-021 costó 900 y acumuló 120 en reparaciones
3. Ejecuta «Consultar coste de ciclo de vida»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestra desglose de coste conocido

**Alternativas y errores:**
- Monedas distintas: separar importes hasta tener conversión explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestra desglose de coste conocido
- Se reproduce y controla esta alternativa: Monedas distintas: separar importes hasta tener conversión explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de consultar coste de ciclo de vida; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-128"></a>
## TI-128 · Planificar renovación de equipos

**Área:** assets · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se produce lista priorizada para renovación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ocho equipos superan cinco años y no tienen garantía

**Flujo principal:**
1. El actor accede a activos y cmdb dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ocho equipos superan cinco años y no tienen garantía
3. Ejecuta «Planificar renovación de equipos»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se produce lista priorizada para renovación

**Alternativas y errores:**
- Fecha de compra desconocida: marcar dato pendiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se produce lista priorizada para renovación
- Se reproduce y controla esta alternativa: Fecha de compra desconocida: marcar dato pendiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Activos y CMDB: operaciones de planificar renovación de equipos; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-129"></a>
## TI-129 · Pedir aprobación al responsable

**Área:** approvals · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Su responsable recibe solicitud con coste y motivo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana solicita licencia de 120 unidades

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana solicita licencia de 120 unidades
3. Ejecuta «Pedir aprobación al responsable»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Su responsable recibe solicitud con coste y motivo

**Alternativas y errores:**
- Sin responsable asignado: enviar a cola de resolución administrativa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Su responsable recibe solicitud con coste y motivo
- Se reproduce y controla esta alternativa: Sin responsable asignado: enviar a cola de resolución administrativa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de pedir aprobación al responsable; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-130"></a>
## TI-130 · Aprobar solicitud con observaciones

**Área:** approvals · **Actor:** Aprobador designado · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Aprobación registra identidad, fecha y condiciones

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Marta autoriza licencia para Ana

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marta autoriza licencia para Ana
3. Ejecuta «Aprobar solicitud con observaciones»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Aprobación registra identidad, fecha y condiciones

**Alternativas y errores:**
- Solicitud cambió de coste: exigir revisar nueva versión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Aprobación registra identidad, fecha y condiciones
- Se reproduce y controla esta alternativa: Solicitud cambió de coste: exigir revisar nueva versión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de aprobar solicitud con observaciones; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-131"></a>
## TI-131 · Rechazar solicitud justificada

**Área:** approvals · **Actor:** Aprobador designado · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra motivo y se informa al solicitante

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Licencia solicitada duplica una asignación activa

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Licencia solicitada duplica una asignación activa
3. Ejecuta «Rechazar solicitud justificada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra motivo y se informa al solicitante

**Alternativas y errores:**
- Motivo vacío: impedir rechazo sin explicación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra motivo y se informa al solicitante
- Se reproduce y controla esta alternativa: Motivo vacío: impedir rechazo sin explicación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de rechazar solicitud justificada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-132"></a>
## TI-132 · Solicitar aclaración antes de aprobar

**Área:** approvals · **Actor:** Aprobador designado · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Solicitud vuelve a aclaración conservando revisión previa

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Falta duración del acceso solicitado

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Falta duración del acceso solicitado
3. Ejecuta «Solicitar aclaración antes de aprobar»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Solicitud vuelve a aclaración conservando revisión previa

**Alternativas y errores:**
- Solicitante inactivo: informar al responsable del proceso.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Solicitud vuelve a aclaración conservando revisión previa
- Se reproduce y controla esta alternativa: Solicitante inactivo: informar al responsable del proceso.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de solicitar aclaración antes de aprobar; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-133"></a>
## TI-133 · Resolver aprobación en varios niveles

**Área:** approvals · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se ejecutan etapas en orden con decisiones separadas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Compra de 5000 requiere responsable y finanzas

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Compra de 5000 requiere responsable y finanzas
3. Ejecuta «Resolver aprobación en varios niveles»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se ejecutan etapas en orden con decisiones separadas

**Alternativas y errores:**
- Primer nivel rechaza: no ejecutar provisión posterior.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se ejecutan etapas en orden con decisiones separadas
- Se reproduce y controla esta alternativa: Primer nivel rechaza: no ejecutar provisión posterior.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de resolver aprobación en varios niveles; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-134"></a>
## TI-134 · Exigir separación de funciones

**Área:** approvals · **Actor:** Administrador de organización · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Ana no puede aprobar su propia solicitud restringida

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Solicitante Ana también integra grupo aprobador

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitante Ana también integra grupo aprobador
3. Ejecuta «Exigir separación de funciones»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ana no puede aprobar su propia solicitud restringida

**Alternativas y errores:**
- Sin aprobadores alternativos: detener y escalar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ana no puede aprobar su propia solicitud restringida
- Se reproduce y controla esta alternativa: Sin aprobadores alternativos: detener y escalar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de exigir separación de funciones; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-135"></a>
## TI-135 · Delegar aprobación durante ausencia

**Área:** approvals · **Actor:** Aprobador designado · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Luis actúa durante ese intervalo con identidad propia

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Marta delega a Luis del 2 al 6 de noviembre

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marta delega a Luis del 2 al 6 de noviembre
3. Ejecuta «Delegar aprobación durante ausencia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Luis actúa durante ese intervalo con identidad propia

**Alternativas y errores:**
- Delegación vencida: retirar facultad inmediatamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Luis actúa durante ese intervalo con identidad propia
- Se reproduce y controla esta alternativa: Delegación vencida: retirar facultad inmediatamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de delegar aprobación durante ausencia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-136"></a>
## TI-136 · Escalar aprobación vencida

**Área:** approvals · **Actor:** Servicio de integración autorizado · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se avisa al suplente configurado conservando trazabilidad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Aprobación sin respuesta supera 48 horas

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Aprobación sin respuesta supera 48 horas
3. Ejecuta «Escalar aprobación vencida»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se avisa al suplente configurado conservando trazabilidad

**Alternativas y errores:**
- Sin suplente: alertar administrador sin autoaprobar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se avisa al suplente configurado conservando trazabilidad
- Se reproduce y controla esta alternativa: Sin suplente: alertar administrador sin autoaprobar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de escalar aprobación vencida; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-137"></a>
## TI-137 · Retirar una petición de aprobación

**Área:** approvals · **Actor:** Solicitante · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Aprobación pendiente queda retirada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana cancela la necesidad antes de decidirse

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana cancela la necesidad antes de decidirse
3. Ejecuta «Retirar una petición de aprobación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Aprobación pendiente queda retirada

**Alternativas y errores:**
- Decisión ya emitida: registrar cancelación posterior por otro flujo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Aprobación pendiente queda retirada
- Se reproduce y controla esta alternativa: Decisión ya emitida: registrar cancelación posterior por otro flujo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de retirar una petición de aprobación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-138"></a>
## TI-138 · Aplicar aprobación paralela

**Área:** approvals · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Solo avanza cuando ambos grupos aprueban

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Cambio requiere aceptación de Seguridad y Operaciones

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Cambio requiere aceptación de Seguridad y Operaciones
3. Ejecuta «Aplicar aprobación paralela»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Solo avanza cuando ambos grupos aprueban

**Alternativas y errores:**
- Un grupo rechaza: bloquear avance aunque el otro apruebe.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Solo avanza cuando ambos grupos aprueban
- Se reproduce y controla esta alternativa: Un grupo rechaza: bloquear avance aunque el otro apruebe.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de aplicar aprobación paralela; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-139"></a>
## TI-139 · Reabrir revisión por cambio material

**Área:** approvals · **Actor:** Servicio de integración autorizado · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** La autorización anterior queda invalidada para el nuevo importe

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Coste aprobado aumenta de 120 a 300

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Coste aprobado aumenta de 120 a 300
3. Ejecuta «Reabrir revisión por cambio material»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La autorización anterior queda invalidada para el nuevo importe

**Alternativas y errores:**
- Cambio cosmético: conservar aprobación según política documentada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La autorización anterior queda invalidada para el nuevo importe
- Se reproduce y controla esta alternativa: Cambio cosmético: conservar aprobación según política documentada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de reabrir revisión por cambio material; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-140"></a>
## TI-140 · Consultar historial de decisiones

**Área:** approvals · **Actor:** Auditor de solo lectura · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se muestran versiones, decisiones y delegaciones completas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Solicitud 81 pasó por dos responsables y un delegado

**Flujo principal:**
1. El actor accede a aprobaciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitud 81 pasó por dos responsables y un delegado
3. Ejecuta «Consultar historial de decisiones»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran versiones, decisiones y delegaciones completas

**Alternativas y errores:**
- Comentario reservado: aplicar restricción del auditor.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran versiones, decisiones y delegaciones completas
- Se reproduce y controla esta alternativa: Comentario reservado: aplicar restricción del auditor.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Aprobaciones: operaciones de consultar historial de decisiones; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-141"></a>
## TI-141 · Abrir problema a partir de recurrencia

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se crea problema enlazado a los incidentes

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Cinco incidentes de impresión ocurren en una semana

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Cinco incidentes de impresión ocurren en una semana
3. Ejecuta «Abrir problema a partir de recurrencia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se crea problema enlazado a los incidentes

**Alternativas y errores:**
- Problema equivalente activo: proponer relación al existente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se crea problema enlazado a los incidentes
- Se reproduce y controla esta alternativa: Problema equivalente activo: proponer relación al existente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de abrir problema a partir de recurrencia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-142"></a>
## TI-142 · Agrupar incidentes bajo un problema

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** El problema conserva evidencias de los incidentes

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets 42, 45 y 49 comparten error E204

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets 42, 45 y 49 comparten error E204
3. Ejecuta «Agrupar incidentes bajo un problema»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El problema conserva evidencias de los incidentes

**Alternativas y errores:**
- Incidente ajeno: impedir relación cruzada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El problema conserva evidencias de los incidentes
- Se reproduce y controla esta alternativa: Incidente ajeno: impedir relación cruzada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de agrupar incidentes bajo un problema; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-143"></a>
## TI-143 · Priorizar investigación de problema

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se fija prioridad y responsable de investigación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Falla recurrente afecta 80 empleados

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Falla recurrente afecta 80 empleados
3. Ejecuta «Priorizar investigación de problema»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se fija prioridad y responsable de investigación

**Alternativas y errores:**
- Impacto sin evidencia: dejarlo pendiente de validación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se fija prioridad y responsable de investigación
- Se reproduce y controla esta alternativa: Impacto sin evidencia: dejarlo pendiente de validación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de priorizar investigación de problema; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-144"></a>
## TI-144 · Registrar hipótesis de causa

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** La hipótesis queda separada de la causa confirmada

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se sospecha desbordamiento de cola del controlador

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se sospecha desbordamiento de cola del controlador
3. Ejecuta «Registrar hipótesis de causa»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La hipótesis queda separada de la causa confirmada

**Alternativas y errores:**
- Evidencia contradictoria: marcar hipótesis descartada con motivo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La hipótesis queda separada de la causa confirmada
- Se reproduce y controla esta alternativa: Evidencia contradictoria: marcar hipótesis descartada con motivo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de registrar hipótesis de causa; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-145"></a>
## TI-145 · Documentar causa raíz verificada

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Problema contiene causa y evidencia reproducible

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Registro demuestra fuga de memoria en versión 4.2

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registro demuestra fuga de memoria en versión 4.2
3. Ejecuta «Documentar causa raíz verificada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Problema contiene causa y evidencia reproducible

**Alternativas y errores:**
- Evidencia ausente: no declarar causa confirmada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Problema contiene causa y evidencia reproducible
- Se reproduce y controla esta alternativa: Evidencia ausente: no declarar causa confirmada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de documentar causa raíz verificada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-146"></a>
## TI-146 · Publicar solución temporal

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Agentes disponen de mitigación enlazada al problema

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Reiniciar cola restablece impresión hasta actualizar

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Reiniciar cola restablece impresión hasta actualizar
3. Ejecuta «Publicar solución temporal»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Agentes disponen de mitigación enlazada al problema

**Alternativas y errores:**
- Mitigación con riesgo: exigir advertencia y autorización correspondiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Agentes disponen de mitigación enlazada al problema
- Se reproduce y controla esta alternativa: Mitigación con riesgo: exigir advertencia y autorización correspondiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de publicar solución temporal; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-147"></a>
## TI-147 · Registrar error conocido

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se publica error conocido con síntomas y mitigación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Defecto de controlador confirmado sin parche aplicado

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Defecto de controlador confirmado sin parche aplicado
3. Ejecuta «Registrar error conocido»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se publica error conocido con síntomas y mitigación

**Alternativas y errores:**
- Causa aún hipotética: conservar estado de investigación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se publica error conocido con síntomas y mitigación
- Se reproduce y controla esta alternativa: Causa aún hipotética: conservar estado de investigación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de registrar error conocido; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-148"></a>
## TI-148 · Proponer corrección definitiva

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** La propuesta se enlaza al problema y cambio necesario

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Actualizar controlador a versión corregida 4.3

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Actualizar controlador a versión corregida 4.3
3. Ejecuta «Proponer corrección definitiva»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La propuesta se enlaza al problema y cambio necesario

**Alternativas y errores:**
- No existe versión validada: mantener propuesta pendiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La propuesta se enlaza al problema y cambio necesario
- Se reproduce y controla esta alternativa: No existe versión validada: mantener propuesta pendiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de proponer corrección definitiva; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-149"></a>
## TI-149 · Revisar eficacia de corrección

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra ventana observada y conclusión de eficacia

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Durante 14 días no reaparece E204

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Durante 14 días no reaparece E204
3. Ejecuta «Revisar eficacia de corrección»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra ventana observada y conclusión de eficacia

**Alternativas y errores:**
- Reaparece incidente relacionado: mantener problema abierto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra ventana observada y conclusión de eficacia
- Se reproduce y controla esta alternativa: Reaparece incidente relacionado: mantener problema abierto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de revisar eficacia de corrección; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-150"></a>
## TI-150 · Cerrar problema con evidencias

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Problema cierra conservando vínculos y causa

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Corrección verificada y documentación actualizada

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Corrección verificada y documentación actualizada
3. Ejecuta «Cerrar problema con evidencias»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Problema cierra conservando vínculos y causa

**Alternativas y errores:**
- Acciones pendientes: impedir cierre definitivo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Problema cierra conservando vínculos y causa
- Se reproduce y controla esta alternativa: Acciones pendientes: impedir cierre definitivo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de cerrar problema con evidencias; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-151"></a>
## TI-151 · Reabrir problema recurrente

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se inicia nueva investigación preservando diagnóstico anterior

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. E204 reaparece después del cierre

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: E204 reaparece después del cierre
3. Ejecuta «Reabrir problema recurrente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se inicia nueva investigación preservando diagnóstico anterior

**Alternativas y errores:**
- Síntomas diferentes: crear problema distinto relacionado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se inicia nueva investigación preservando diagnóstico anterior
- Se reproduce y controla esta alternativa: Síntomas diferentes: crear problema distinto relacionado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de reabrir problema recurrente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-152"></a>
## TI-152 · Consultar tendencias de problemas

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se distinguen recurrencia, duración y servicios afectados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se comparan fallos de impresión de tres meses

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se comparan fallos de impresión de tres meses
3. Ejecuta «Consultar tendencias de problemas»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se distinguen recurrencia, duración y servicios afectados

**Alternativas y errores:**
- Periodo sin datos: no inferir mejora.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se distinguen recurrencia, duración y servicios afectados
- Se reproduce y controla esta alternativa: Periodo sin datos: no inferir mejora.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de consultar tendencias de problemas; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-153"></a>
## TI-153 · Solicitar cambio normal

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Cambio contiene objetivo, alcance y propietario

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Actualizar controlador de impresión corporativo

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Actualizar controlador de impresión corporativo
3. Ejecuta «Solicitar cambio normal»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Cambio contiene objetivo, alcance y propietario

**Alternativas y errores:**
- Falta servicio afectado: impedir enviar a evaluación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Cambio contiene objetivo, alcance y propietario
- Se reproduce y controla esta alternativa: Falta servicio afectado: impedir enviar a evaluación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de solicitar cambio normal; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-154"></a>
## TI-154 · Clasificar riesgo de cambio

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra riesgo según impacto y reversibilidad

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Actualización afecta 40 equipos durante consulta

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Actualización afecta 40 equipos durante consulta
3. Ejecuta «Clasificar riesgo de cambio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra riesgo según impacto y reversibilidad

**Alternativas y errores:**
- Respuestas incompletas: mantener evaluación pendiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra riesgo según impacto y reversibilidad
- Se reproduce y controla esta alternativa: Respuestas incompletas: mantener evaluación pendiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de clasificar riesgo de cambio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-155"></a>
## TI-155 · Definir plan de implementación

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Plan ordenado identifica responsables y verificaciones

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Actualizar primero dos equipos piloto y luego el resto

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Actualizar primero dos equipos piloto y luego el resto
3. Ejecuta «Definir plan de implementación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Plan ordenado identifica responsables y verificaciones

**Alternativas y errores:**
- Paso destructivo sin mitigación: exigir revisión de riesgo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Plan ordenado identifica responsables y verificaciones
- Se reproduce y controla esta alternativa: Paso destructivo sin mitigación: exigir revisión de riesgo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de definir plan de implementación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-156"></a>
## TI-156 · Definir plan de retorno

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Existe criterio de retorno y responsable ejecutor

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Restaurar controlador 4.2 si falla la impresión

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Restaurar controlador 4.2 si falla la impresión
3. Ejecuta «Definir plan de retorno»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Existe criterio de retorno y responsable ejecutor

**Alternativas y errores:**
- Retorno técnicamente imposible: declarar riesgo y autorización especial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Existe criterio de retorno y responsable ejecutor
- Se reproduce y controla esta alternativa: Retorno técnicamente imposible: declarar riesgo y autorización especial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de definir plan de retorno; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-157"></a>
## TI-157 · Evaluar dependencias de cambio

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se vinculan prerrequisitos y CIs afectados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Actualización depende de reinicio del servidor R3

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Actualización depende de reinicio del servidor R3
3. Ejecuta «Evaluar dependencias de cambio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se vinculan prerrequisitos y CIs afectados

**Alternativas y errores:**
- Dependencia circular: impedir programación incoherente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se vinculan prerrequisitos y CIs afectados
- Se reproduce y controla esta alternativa: Dependencia circular: impedir programación incoherente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de evaluar dependencias de cambio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-158"></a>
## TI-158 · Solicitar autorización del comité

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registran decisiones de los aprobadores requeridos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Cambio de alto riesgo requiere comité de cambios

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Cambio de alto riesgo requiere comité de cambios
3. Ejecuta «Solicitar autorización del comité»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registran decisiones de los aprobadores requeridos

**Alternativas y errores:**
- Quorum incompleto: no autorizar ejecución.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registran decisiones de los aprobadores requeridos
- Se reproduce y controla esta alternativa: Quorum incompleto: no autorizar ejecución.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de solicitar autorización del comité; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-159"></a>
## TI-159 · Programar ventana de cambio

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Calendario muestra ventana y servicios afectados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Domingo de 02:00 a 03:00 para R3

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Domingo de 02:00 a 03:00 para R3
3. Ejecuta «Programar ventana de cambio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Calendario muestra ventana y servicios afectados

**Alternativas y errores:**
- Ventana coincide con cambio incompatible: informar conflicto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Calendario muestra ventana y servicios afectados
- Se reproduce y controla esta alternativa: Ventana coincide con cambio incompatible: informar conflicto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de programar ventana de cambio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-160"></a>
## TI-160 · Notificar intervención programada

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Afectados reciben horario y alcance públicos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Usuarios de Radiología tendrán interrupción de diez minutos

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Usuarios de Radiología tendrán interrupción de diez minutos
3. Ejecuta «Notificar intervención programada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Afectados reciben horario y alcance públicos

**Alternativas y errores:**
- Contacto no autorizado: excluir datos internos del aviso.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Afectados reciben horario y alcance públicos
- Se reproduce y controla esta alternativa: Contacto no autorizado: excluir datos internos del aviso.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de notificar intervención programada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-161"></a>
## TI-161 · Ejecutar tarea de cambio

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se registra inicio, evidencia y resultado de la tarea

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro actualiza los dos equipos piloto

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro actualiza los dos equipos piloto
3. Ejecuta «Ejecutar tarea de cambio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra inicio, evidencia y resultado de la tarea

**Alternativas y errores:**
- Validación piloto falla: detener despliegue restante.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra inicio, evidencia y resultado de la tarea
- Se reproduce y controla esta alternativa: Validación piloto falla: detener despliegue restante.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de ejecutar tarea de cambio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-162"></a>
## TI-162 · Activar retorno de cambio

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se ejecuta retorno y se registra resultado fallido del cambio

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Versión nueva impide imprimir formularios críticos

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Versión nueva impide imprimir formularios críticos
3. Ejecuta «Activar retorno de cambio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se ejecuta retorno y se registra resultado fallido del cambio

**Alternativas y errores:**
- Retorno falla: escalar incidente mayor asociado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se ejecuta retorno y se registra resultado fallido del cambio
- Se reproduce y controla esta alternativa: Retorno falla: escalar incidente mayor asociado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de activar retorno de cambio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-163"></a>
## TI-163 · Validar servicio tras cambio

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Pruebas de servicio quedan registradas antes del cierre

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ana imprime formulario y etiqueta desde R3

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ana imprime formulario y etiqueta desde R3
3. Ejecuta «Validar servicio tras cambio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Pruebas de servicio quedan registradas antes del cierre

**Alternativas y errores:**
- Una prueba falla: mantener cambio en validación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Pruebas de servicio quedan registradas antes del cierre
- Se reproduce y controla esta alternativa: Una prueba falla: mantener cambio en validación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de validar servicio tras cambio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-164"></a>
## TI-164 · Cerrar cambio satisfactorio

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Cambio cierra con tiempos reales y resultado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Todas las pruebas y tareas están completas

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Todas las pruebas y tareas están completas
3. Ejecuta «Cerrar cambio satisfactorio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Cambio cierra con tiempos reales y resultado

**Alternativas y errores:**
- Falta confirmación de servicio: bloquear cierre satisfactorio.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Cambio cierra con tiempos reales y resultado
- Se reproduce y controla esta alternativa: Falta confirmación de servicio: bloquear cierre satisfactorio.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de cerrar cambio satisfactorio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-165"></a>
## TI-165 · Revisar cambio fallido

**Área:** problems-changes · **Actor:** Responsable de servicio · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Revisión documenta causa, impacto y acciones correctivas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Despliegue requirió retorno por incompatibilidad

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Despliegue requirió retorno por incompatibilidad
3. Ejecuta «Revisar cambio fallido»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Revisión documenta causa, impacto y acciones correctivas

**Alternativas y errores:**
- Acción sin responsable: impedir dar revisión por concluida.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Revisión documenta causa, impacto y acciones correctivas
- Se reproduce y controla esta alternativa: Acción sin responsable: impedir dar revisión por concluida.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de revisar cambio fallido; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-166"></a>
## TI-166 · Ejecutar cambio estándar preautorizado

**Área:** problems-changes · **Actor:** Agente · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se aplica plantilla con evidencia de ejecución

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Alta de impresora sigue procedimiento aprobado vigente

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Alta de impresora sigue procedimiento aprobado vigente
3. Ejecuta «Ejecutar cambio estándar preautorizado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se aplica plantilla con evidencia de ejecución

**Alternativas y errores:**
- El alcance supera la plantilla: tramitar cambio normal.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se aplica plantilla con evidencia de ejecución
- Se reproduce y controla esta alternativa: El alcance supera la plantilla: tramitar cambio normal.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de ejecutar cambio estándar preautorizado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-167"></a>
## TI-167 · Tramitar cambio de emergencia

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Autorización de emergencia y revisión posterior quedan registradas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Parche urgente corrige caída activa de PACS

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Parche urgente corrige caída activa de PACS
3. Ejecuta «Tramitar cambio de emergencia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Autorización de emergencia y revisión posterior quedan registradas

**Alternativas y errores:**
- Sin autoridad de emergencia: escalar sin omitir decisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Autorización de emergencia y revisión posterior quedan registradas
- Se reproduce y controla esta alternativa: Sin autoridad de emergencia: escalar sin omitir decisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de tramitar cambio de emergencia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-168"></a>
## TI-168 · Cancelar cambio programado

**Área:** problems-changes · **Actor:** Gestor de cambios · **Prioridad:** P2 · **Estado:** `planned`

**Objetivo:** Se cancelan tareas y se avisa a afectados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Proveedor retira versión antes de la ventana

**Flujo principal:**
1. El actor accede a problemas y cambios dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Proveedor retira versión antes de la ventana
3. Ejecuta «Cancelar cambio programado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se cancelan tareas y se avisa a afectados

**Alternativas y errores:**
- Tareas ya ejecutadas: registrar estado real y compensación requerida.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se cancelan tareas y se avisa a afectados
- Se reproduce y controla esta alternativa: Tareas ya ejecutadas: registrar estado real y compensación requerida.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Problemas y cambios: operaciones de cancelar cambio programado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-169"></a>
## TI-169 · Crear regla de asignación automática

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Regla asigna solo tickets que cumplen todas sus condiciones

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets de impresoras van a Infraestructura

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets de impresoras van a Infraestructura
3. Ejecuta «Crear regla de asignación automática»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Regla asigna solo tickets que cumplen todas sus condiciones

**Alternativas y errores:**
- Condición desconocida: rechazar regla en lugar de ignorarla.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Regla asigna solo tickets que cumplen todas sus condiciones
- Se reproduce y controla esta alternativa: Condición desconocida: rechazar regla en lugar de ignorarla.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de crear regla de asignación automática; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-170"></a>
## TI-170 · Validar regla antes de activarla

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se muestran errores de operadores y campos incompatibles

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Regla compara prioridad con valor crítico

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Regla compara prioridad con valor crítico
3. Ejecuta «Validar regla antes de activarla»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran errores de operadores y campos incompatibles

**Alternativas y errores:**
- Operador no soportado: impedir activación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran errores de operadores y campos incompatibles
- Se reproduce y controla esta alternativa: Operador no soportado: impedir activación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de validar regla antes de activarla; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-171"></a>
## TI-171 · Simular regla con ejemplo

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La simulación informa que no coincide sin efectos reales

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket de impresora normal y regla urgente

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket de impresora normal y regla urgente
3. Ejecuta «Simular regla con ejemplo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La simulación informa que no coincide sin efectos reales

**Alternativas y errores:**
- Acción externa: nunca enviarla durante simulación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La simulación informa que no coincide sin efectos reales
- Se reproduce y controla esta alternativa: Acción externa: nunca enviarla durante simulación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de simular regla con ejemplo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-172"></a>
## TI-172 · Activar automatización programada

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se ejecutan acciones sobre candidatos autorizados

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Revisar tickets sin respuesta cada quince minutos

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisar tickets sin respuesta cada quince minutos
3. Ejecuta «Activar automatización programada»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se ejecutan acciones sobre candidatos autorizados

**Alternativas y errores:**
- Dos workers coinciden: coordinar para evitar ejecución duplicada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se ejecutan acciones sobre candidatos autorizados
- Se reproduce y controla esta alternativa: Dos workers coinciden: coordinar para evitar ejecución duplicada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de activar automatización programada; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-173"></a>
## TI-173 · Desactivar regla problemática

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Nuevos eventos dejan de ejecutar la regla

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Regla nueva asigna equipo incorrecto

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Regla nueva asigna equipo incorrecto
3. Ejecuta «Desactivar regla problemática»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Nuevos eventos dejan de ejecutar la regla

**Alternativas y errores:**
- Acción ya ejecutada: conservar historial para corrección manual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Nuevos eventos dejan de ejecutar la regla
- Se reproduce y controla esta alternativa: Acción ya ejecutada: conservar historial para corrección manual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de desactivar regla problemática; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-174"></a>
## TI-174 · Consultar historial de ejecuciones

**Área:** automation · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se ve versión, evento, condiciones y resultado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Regla 8 actuó sobre ticket 42 a las 10:00

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Regla 8 actuó sobre ticket 42 a las 10:00
3. Ejecuta «Consultar historial de ejecuciones»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se ve versión, evento, condiciones y resultado

**Alternativas y errores:**
- Registro contiene secreto: presentar datos redactados.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se ve versión, evento, condiciones y resultado
- Se reproduce y controla esta alternativa: Registro contiene secreto: presentar datos redactados.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de consultar historial de ejecuciones; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-175"></a>
## TI-175 · Enviar webhook al cambiar estado

**Área:** automation · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Destino autorizado recibe evento firmado con identificador estable

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 pasa a resuelto

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 pasa a resuelto
3. Ejecuta «Enviar webhook al cambiar estado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Destino autorizado recibe evento firmado con identificador estable

**Alternativas y errores:**
- Destino caído: reintentar con límites y registro.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Destino autorizado recibe evento firmado con identificador estable
- Se reproduce y controla esta alternativa: Destino caído: reintentar con límites y registro.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de enviar webhook al cambiar estado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-176"></a>
## TI-176 · Recibir evento externo idempotente

**Área:** automation · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se aplica una sola actualización por identificador

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Sistema de inventario repite evento EVT-77

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Sistema de inventario repite evento EVT-77
3. Ejecuta «Recibir evento externo idempotente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se aplica una sola actualización por identificador

**Alternativas y errores:**
- Mismo identificador y otro contenido: registrar conflicto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se aplica una sola actualización por identificador
- Se reproduce y controla esta alternativa: Mismo identificador y otro contenido: registrar conflicto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de recibir evento externo idempotente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-177"></a>
## TI-177 · Reprocesar fallo de integración

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Tras corregir credencial se reintenta el evento original

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Webhook EVT-77 falló por credencial vencida

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Webhook EVT-77 falló por credencial vencida
3. Ejecuta «Reprocesar fallo de integración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Tras corregir credencial se reintenta el evento original

**Alternativas y errores:**
- Evento ya completado: impedir duplicar efectos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Tras corregir credencial se reintenta el evento original
- Se reproduce y controla esta alternativa: Evento ya completado: impedir duplicar efectos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de reprocesar fallo de integración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-178"></a>
## TI-178 · Sincronizar miembros de directorio

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se concilian miembros según mapeo autorizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Grupo LDAP añade dos agentes y retira uno

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Grupo LDAP añade dos agentes y retira uno
3. Ejecuta «Sincronizar miembros de directorio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se concilian miembros según mapeo autorizado

**Alternativas y errores:**
- Directorio incompleto: no desactivar masivamente por error transitorio.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se concilian miembros según mapeo autorizado
- Se reproduce y controla esta alternativa: Directorio incompleto: no desactivar masivamente por error transitorio.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de sincronizar miembros de directorio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-179"></a>
## TI-179 · Notificar resolución por correo

**Área:** automation · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Mensaje durable se entrega o conserva fallo visible

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 se resuelve para Ana

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 se resuelve para Ana
3. Ejecuta «Notificar resolución por correo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Mensaje durable se entrega o conserva fallo visible

**Alternativas y errores:**
- SMTP no disponible: conservar tarea pendiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Mensaje durable se entrega o conserva fallo visible
- Se reproduce y controla esta alternativa: SMTP no disponible: conservar tarea pendiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de notificar resolución por correo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-180"></a>
## TI-180 · Consultar cola de correo fallido

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se muestran causas y destinatarios autorizados con datos mínimos

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tres correos agotaron intentos de entrega

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tres correos agotaron intentos de entrega
3. Ejecuta «Consultar cola de correo fallido»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se muestran causas y destinatarios autorizados con datos mínimos

**Alternativas y errores:**
- Contenido sensible: no exponer cuerpo en logs generales.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se muestran causas y destinatarios autorizados con datos mínimos
- Se reproduce y controla esta alternativa: Contenido sensible: no exponer cuerpo en logs generales.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de consultar cola de correo fallido; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-181"></a>
## TI-181 · Reintentar correo tras corregir configuración

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se agenda un intento trazable del mensaje pendiente

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. SMTP vuelve a estar disponible para correo 91

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: SMTP vuelve a estar disponible para correo 91
3. Ejecuta «Reintentar correo tras corregir configuración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se agenda un intento trazable del mensaje pendiente

**Alternativas y errores:**
- Mensaje ya enviado: advertir antes de reenvío explícito.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se agenda un intento trazable del mensaje pendiente
- Se reproduce y controla esta alternativa: Mensaje ya enviado: advertir antes de reenvío explícito.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de reintentar correo tras corregir configuración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-182"></a>
## TI-182 · Aplicar límite de solicitudes API

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se limita consumo sin afectar a otras organizaciones

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Integración supera su cuota por minuto

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Integración supera su cuota por minuto
3. Ejecuta «Aplicar límite de solicitudes API»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se limita consumo sin afectar a otras organizaciones

**Alternativas y errores:**
- Cambio de réplica: conservar límite compartido efectivo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se limita consumo sin afectar a otras organizaciones
- Se reproduce y controla esta alternativa: Cambio de réplica: conservar límite compartido efectivo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de aplicar límite de solicitudes api; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-183"></a>
## TI-183 · Revocar acceso de conector

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Peticiones posteriores del conector fallan y quedan auditadas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Inventario externo deja de estar autorizado

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Inventario externo deja de estar autorizado
3. Ejecuta «Revocar acceso de conector»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Peticiones posteriores del conector fallan y quedan auditadas

**Alternativas y errores:**
- Trabajo en curso: detener nuevas acciones de la credencial revocada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Peticiones posteriores del conector fallan y quedan auditadas
- Se reproduce y controla esta alternativa: Trabajo en curso: detener nuevas acciones de la credencial revocada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de revocar acceso de conector; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-184"></a>
## TI-184 · Consultar salud de integraciones

**Área:** automation · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Panel diferencia conexión, último éxito y error actual

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. LDAP responde y SMTP lleva diez minutos fallando

**Flujo principal:**
1. El actor accede a automatizaciones e integraciones dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: LDAP responde y SMTP lleva diez minutos fallando
3. Ejecuta «Consultar salud de integraciones»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Panel diferencia conexión, último éxito y error actual

**Alternativas y errores:**
- Nunca ejecutado: mostrar sin evidencia en vez de saludable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Panel diferencia conexión, último éxito y error actual
- Se reproduce y controla esta alternativa: Nunca ejecutado: mostrar sin evidencia en vez de saludable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Automatizaciones e integraciones: operaciones de consultar salud de integraciones; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-185"></a>
## TI-185 · Consultar volumen de tickets

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Totales coinciden con tickets visibles y fechas del periodo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Octubre agrupado por día de Caracas

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Octubre agrupado por día de Caracas
3. Ejecuta «Consultar volumen de tickets»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Totales coinciden con tickets visibles y fechas del periodo

**Alternativas y errores:**
- Intervalo invertido: rechazar consulta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Totales coinciden con tickets visibles y fechas del periodo
- Se reproduce y controla esta alternativa: Intervalo invertido: rechazar consulta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de consultar volumen de tickets; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-186"></a>
## TI-186 · Comparar estados de cola

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Gráfico y listado usan la misma definición de estado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets abiertos, en progreso y resueltos actuales

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets abiertos, en progreso y resueltos actuales
3. Ejecuta «Comparar estados de cola»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Gráfico y listado usan la misma definición de estado

**Alternativas y errores:**
- Estado nuevo desconocido: mostrar categoría explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Gráfico y listado usan la misma definición de estado
- Se reproduce y controla esta alternativa: Estado nuevo desconocido: mostrar categoría explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de comparar estados de cola; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-187"></a>
## TI-187 · Medir tiempo de primera respuesta

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Indicador usa primera respuesta pública real

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets respondidos durante octubre

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets respondidos durante octubre
3. Ejecuta «Medir tiempo de primera respuesta»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Indicador usa primera respuesta pública real

**Alternativas y errores:**
- Tickets sin respuesta: excluir del promedio y contar pendientes.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Indicador usa primera respuesta pública real
- Se reproduce y controla esta alternativa: Tickets sin respuesta: excluir del promedio y contar pendientes.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de medir tiempo de primera respuesta; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-188"></a>
## TI-188 · Medir tiempo de resolución

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se calcula duración desde marcas válidas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Incidentes resueltos durante la última semana

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Incidentes resueltos durante la última semana
3. Ejecuta «Medir tiempo de resolución»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se calcula duración desde marcas válidas

**Alternativas y errores:**
- Fechas heredadas ausentes: presentar registros no calculables.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se calcula duración desde marcas válidas
- Se reproduce y controla esta alternativa: Fechas heredadas ausentes: presentar registros no calculables.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de medir tiempo de resolución; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-189"></a>
## TI-189 · Consultar carga por agente

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se ve distribución de trabajo dentro del equipo autorizado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro tiene ocho tickets y Marta tiene tres

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro tiene ocho tickets y Marta tiene tres
3. Ejecuta «Consultar carga por agente»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se ve distribución de trabajo dentro del equipo autorizado

**Alternativas y errores:**
- Agente desactivado con tickets: incluirlo como pendiente de reasignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se ve distribución de trabajo dentro del equipo autorizado
- Se reproduce y controla esta alternativa: Agente desactivado con tickets: incluirlo como pendiente de reasignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de consultar carga por agente; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-190"></a>
## TI-190 · Identificar tickets envejecidos

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Lista ordena antigüedad con responsable y siguiente acción

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets abiertos por más de diez días

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets abiertos por más de diez días
3. Ejecuta «Identificar tickets envejecidos»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Lista ordena antigüedad con responsable y siguiente acción

**Alternativas y errores:**
- Ticket pausado: mostrar pausa sin ocultar antigüedad total.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Lista ordena antigüedad con responsable y siguiente acción
- Se reproduce y controla esta alternativa: Ticket pausado: mostrar pausa sin ocultar antigüedad total.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de identificar tickets envejecidos; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-191"></a>
## TI-191 · Analizar categorías más frecuentes

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Ranking enlaza a tickets de cada categoría

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Impresión concentra 40 incidentes en octubre

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Impresión concentra 40 incidentes en octubre
3. Ejecuta «Analizar categorías más frecuentes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Ranking enlaza a tickets de cada categoría

**Alternativas y errores:**
- Categoría eliminada: conservar etiqueta histórica identificable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Ranking enlaza a tickets de cada categoría
- Se reproduce y controla esta alternativa: Categoría eliminada: conservar etiqueta histórica identificable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de analizar categorías más frecuentes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-192"></a>
## TI-192 · Exportar resultados a CSV

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Archivo respeta filtros y neutraliza fórmulas de celdas

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Vista filtrada de 42 tickets de Radiología

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Vista filtrada de 42 tickets de Radiología
3. Ejecuta «Exportar resultados a CSV»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Archivo respeta filtros y neutraliza fórmulas de celdas

**Alternativas y errores:**
- Título empieza con igual: exportar como texto seguro.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Archivo respeta filtros y neutraliza fórmulas de celdas
- Se reproduce y controla esta alternativa: Título empieza con igual: exportar como texto seguro.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de exportar resultados a csv; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-193"></a>
## TI-193 · Exportar resumen a PDF

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Documento incluye periodo, filtros y fecha de generación

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Resumen mensual para comité de servicio

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Resumen mensual para comité de servicio
3. Ejecuta «Exportar resumen a PDF»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Documento incluye periodo, filtros y fecha de generación

**Alternativas y errores:**
- Resultado demasiado grande: ofrecer exportación acotada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Documento incluye periodo, filtros y fecha de generación
- Se reproduce y controla esta alternativa: Resultado demasiado grande: ofrecer exportación acotada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de exportar resumen a pdf; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-194"></a>
## TI-194 · Consultar cambios de permisos

**Área:** reporting · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Auditoría muestra actor, valor anterior y nuevo

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Pedro pasó de agente a solicitante ayer

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Pedro pasó de agente a solicitante ayer
3. Ejecuta «Consultar cambios de permisos»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Auditoría muestra actor, valor anterior y nuevo

**Alternativas y errores:**
- Actor técnico: identificar integración responsable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Auditoría muestra actor, valor anterior y nuevo
- Se reproduce y controla esta alternativa: Actor técnico: identificar integración responsable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de consultar cambios de permisos; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-195"></a>
## TI-195 · Investigar acceso a recurso sensible

**Área:** reporting · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Eventos permiten atribuir acceso y resultado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Se revisan accesos autorizados al ticket restringido 88

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Se revisan accesos autorizados al ticket restringido 88
3. Ejecuta «Investigar acceso a recurso sensible»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Eventos permiten atribuir acceso y resultado

**Alternativas y errores:**
- Log no recoge lectura histórica: declarar ausencia de evidencia.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Eventos permiten atribuir acceso y resultado
- Se reproduce y controla esta alternativa: Log no recoge lectura histórica: declarar ausencia de evidencia.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de investigar acceso a recurso sensible; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-196"></a>
## TI-196 · Auditar transiciones de tickets

**Área:** reporting · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se reconstruye secuencia sin alterar registros

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Ticket 42 fue reabierto y reasignado

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Ticket 42 fue reabierto y reasignado
3. Ejecuta «Auditar transiciones de tickets»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se reconstruye secuencia sin alterar registros

**Alternativas y errores:**
- Marcas coincidentes: conservar orden estable por identificador.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se reconstruye secuencia sin alterar registros
- Se reproduce y controla esta alternativa: Marcas coincidentes: conservar orden estable por identificador.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de auditar transiciones de tickets; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-197"></a>
## TI-197 · Revisar cumplimiento de catálogo

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se compara plazo prometido y entrega real

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Solicitudes de portátil entregadas durante octubre

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitudes de portátil entregadas durante octubre
3. Ejecuta «Revisar cumplimiento de catálogo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se compara plazo prometido y entrega real

**Alternativas y errores:**
- Solicitud cancelada: mostrarla separada del cumplimiento.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se compara plazo prometido y entrega real
- Se reproduce y controla esta alternativa: Solicitud cancelada: mostrarla separada del cumplimiento.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de revisar cumplimiento de catálogo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-198"></a>
## TI-198 · Analizar uso de autoservicio

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se presenta interacción medible sin atribuir resolución supuesta

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Guías consultadas antes de abrir tickets de impresión

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Guías consultadas antes de abrir tickets de impresión
3. Ejecuta «Analizar uso de autoservicio»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se presenta interacción medible sin atribuir resolución supuesta

**Alternativas y errores:**
- Usuario anónimo: agregar sin inventar identidad.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se presenta interacción medible sin atribuir resolución supuesta
- Se reproduce y controla esta alternativa: Usuario anónimo: agregar sin inventar identidad.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de analizar uso de autoservicio; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-199"></a>
## TI-199 · Detectar calidad insuficiente de datos

**Área:** reporting · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Informe enumera registros que requieren corrección

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Tickets sin categoría y activos sin sede

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Tickets sin categoría y activos sin sede
3. Ejecuta «Detectar calidad insuficiente de datos»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Informe enumera registros que requieren corrección

**Alternativas y errores:**
- Dato opcional: no contarlo como error obligatorio.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Informe enumera registros que requieren corrección
- Se reproduce y controla esta alternativa: Dato opcional: no contarlo como error obligatorio.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de detectar calidad insuficiente de datos; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="ti-200"></a>
## TI-200 · Revisar historial de exportaciones

**Área:** reporting · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Se registra quién exportó, alcance y resultado

**Ejemplo:** Organización ficticia Clínica Aurora; sede Caracas. Marta exportó tickets de Radiología el lunes

**Flujo principal:**
1. El actor accede a informes y auditoría dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marta exportó tickets de Radiología el lunes
3. Ejecuta «Revisar historial de exportaciones»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Se registra quién exportó, alcance y resultado

**Alternativas y errores:**
- Exportación fallida: distinguir intento de archivo entregado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Se registra quién exportó, alcance y resultado
- Se reproduce y controla esta alternativa: Exportación fallida: distinguir intento de archivo entregado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** TI / Informes y auditoría: operaciones de revisar historial de exportaciones; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0
