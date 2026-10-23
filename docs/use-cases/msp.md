# Trickreport · Servicios MSP

Catálogo de **200 casos propuestos**. Todos comienzan en `planned`; su inclusión no afirma que la función exista.

| Área | Casos |
| --- | ---: |
| Clientes y espacios aislados | 20 |
| Contratos, catálogo y SLA | 20 |
| Mesa de servicio para varios clientes | 20 |
| Técnicos, turnos y capacidad | 20 |
| Proyectos, horas y gastos | 20 |
| Activos, garantías e inventario | 20 |
| Descubrimiento y monitoreo | 20 |
| Soporte remoto integrado | 20 |
| Rentabilidad e informes comerciales | 20 |
| Integraciones, API y seguridad | 20 |

<a id="msp-001"></a>
## MSP-001 · Crear espacio independiente: Clínica Aurora

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El espacio de Clínica Aurora queda separado de los demás clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar un espacio con su dominio y administrador para Clínica Aurora: dominio aurora.example, Caracas.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar un espacio con su dominio y administrador para Clínica Aurora: dominio aurora.example, Caracas.
3. Ejecuta «Crear espacio independiente: Clínica Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El espacio de Clínica Aurora queda separado de los demás clientes.

**Alternativas y errores:**
- Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El espacio de Clínica Aurora queda separado de los demás clientes.
- Se reproduce y controla esta alternativa: Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de crear espacio independiente: clínica aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-002"></a>
## MSP-002 · Crear espacio independiente: Hotel Río

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El espacio de Hotel Río queda separado de los demás clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar un espacio con su dominio y administrador para Hotel Río: dominio hotelrio.example, Valencia.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar un espacio con su dominio y administrador para Hotel Río: dominio hotelrio.example, Valencia.
3. Ejecuta «Crear espacio independiente: Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El espacio de Hotel Río queda separado de los demás clientes.

**Alternativas y errores:**
- Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El espacio de Hotel Río queda separado de los demás clientes.
- Se reproduce y controla esta alternativa: Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de crear espacio independiente: hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-003"></a>
## MSP-003 · Crear espacio independiente: Colegio Sol

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El espacio de Colegio Sol queda separado de los demás clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar un espacio con su dominio y administrador para Colegio Sol: dominio colegiosol.example, Maracay.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar un espacio con su dominio y administrador para Colegio Sol: dominio colegiosol.example, Maracay.
3. Ejecuta «Crear espacio independiente: Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El espacio de Colegio Sol queda separado de los demás clientes.

**Alternativas y errores:**
- Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El espacio de Colegio Sol queda separado de los demás clientes.
- Se reproduce y controla esta alternativa: Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de crear espacio independiente: colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-004"></a>
## MSP-004 · Crear espacio independiente: Logística Norte

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El espacio de Logística Norte queda separado de los demás clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar un espacio con su dominio y administrador para Logística Norte: dominio lognorte.example, Barquisimeto.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar un espacio con su dominio y administrador para Logística Norte: dominio lognorte.example, Barquisimeto.
3. Ejecuta «Crear espacio independiente: Logística Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El espacio de Logística Norte queda separado de los demás clientes.

**Alternativas y errores:**
- Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El espacio de Logística Norte queda separado de los demás clientes.
- Se reproduce y controla esta alternativa: Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de crear espacio independiente: logística norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-005"></a>
## MSP-005 · Añadir sede de cliente: Clínica Aurora

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sede de Clínica Aurora se incorpora al ámbito correcto del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Vincular una sede y su zona horaria para Clínica Aurora: dominio aurora.example, Caracas.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Vincular una sede y su zona horaria para Clínica Aurora: dominio aurora.example, Caracas.
3. Ejecuta «Añadir sede de cliente: Clínica Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sede de Clínica Aurora se incorpora al ámbito correcto del cliente.

**Alternativas y errores:**
- Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sede de Clínica Aurora se incorpora al ámbito correcto del cliente.
- Se reproduce y controla esta alternativa: Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de añadir sede de cliente: clínica aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-006"></a>
## MSP-006 · Añadir sede de cliente: Hotel Río

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sede de Hotel Río se incorpora al ámbito correcto del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Vincular una sede y su zona horaria para Hotel Río: dominio hotelrio.example, Valencia.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Vincular una sede y su zona horaria para Hotel Río: dominio hotelrio.example, Valencia.
3. Ejecuta «Añadir sede de cliente: Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sede de Hotel Río se incorpora al ámbito correcto del cliente.

**Alternativas y errores:**
- Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sede de Hotel Río se incorpora al ámbito correcto del cliente.
- Se reproduce y controla esta alternativa: Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de añadir sede de cliente: hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-007"></a>
## MSP-007 · Añadir sede de cliente: Colegio Sol

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sede de Colegio Sol se incorpora al ámbito correcto del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Vincular una sede y su zona horaria para Colegio Sol: dominio colegiosol.example, Maracay.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Vincular una sede y su zona horaria para Colegio Sol: dominio colegiosol.example, Maracay.
3. Ejecuta «Añadir sede de cliente: Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sede de Colegio Sol se incorpora al ámbito correcto del cliente.

**Alternativas y errores:**
- Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sede de Colegio Sol se incorpora al ámbito correcto del cliente.
- Se reproduce y controla esta alternativa: Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de añadir sede de cliente: colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-008"></a>
## MSP-008 · Añadir sede de cliente: Logística Norte

**Área:** clients · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sede de Logística Norte se incorpora al ámbito correcto del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Vincular una sede y su zona horaria para Logística Norte: dominio lognorte.example, Barquisimeto.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Vincular una sede y su zona horaria para Logística Norte: dominio lognorte.example, Barquisimeto.
3. Ejecuta «Añadir sede de cliente: Logística Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sede de Logística Norte se incorpora al ámbito correcto del cliente.

**Alternativas y errores:**
- Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sede de Logística Norte se incorpora al ámbito correcto del cliente.
- Se reproduce y controla esta alternativa: Código de sede duplicado: solicitar corrección sin crear una sede parcial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de añadir sede de cliente: logística norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-009"></a>
## MSP-009 · Invitar contacto autorizado: Clínica Aurora

**Área:** clients · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contacto de Clínica Aurora sólo accede a recursos permitidos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Invitar una persona con permisos mínimos para Clínica Aurora: dominio aurora.example, Caracas.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Invitar una persona con permisos mínimos para Clínica Aurora: dominio aurora.example, Caracas.
3. Ejecuta «Invitar contacto autorizado: Clínica Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contacto de Clínica Aurora sólo accede a recursos permitidos.

**Alternativas y errores:**
- Correo de otro cliente: impedir que herede acceso entre espacios.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contacto de Clínica Aurora sólo accede a recursos permitidos.
- Se reproduce y controla esta alternativa: Correo de otro cliente: impedir que herede acceso entre espacios.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de invitar contacto autorizado: clínica aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-010"></a>
## MSP-010 · Invitar contacto autorizado: Hotel Río

**Área:** clients · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contacto de Hotel Río sólo accede a recursos permitidos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Invitar una persona con permisos mínimos para Hotel Río: dominio hotelrio.example, Valencia.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Invitar una persona con permisos mínimos para Hotel Río: dominio hotelrio.example, Valencia.
3. Ejecuta «Invitar contacto autorizado: Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contacto de Hotel Río sólo accede a recursos permitidos.

**Alternativas y errores:**
- Correo de otro cliente: impedir que herede acceso entre espacios.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contacto de Hotel Río sólo accede a recursos permitidos.
- Se reproduce y controla esta alternativa: Correo de otro cliente: impedir que herede acceso entre espacios.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de invitar contacto autorizado: hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-011"></a>
## MSP-011 · Invitar contacto autorizado: Colegio Sol

**Área:** clients · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contacto de Colegio Sol sólo accede a recursos permitidos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Invitar una persona con permisos mínimos para Colegio Sol: dominio colegiosol.example, Maracay.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Invitar una persona con permisos mínimos para Colegio Sol: dominio colegiosol.example, Maracay.
3. Ejecuta «Invitar contacto autorizado: Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contacto de Colegio Sol sólo accede a recursos permitidos.

**Alternativas y errores:**
- Correo de otro cliente: impedir que herede acceso entre espacios.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contacto de Colegio Sol sólo accede a recursos permitidos.
- Se reproduce y controla esta alternativa: Correo de otro cliente: impedir que herede acceso entre espacios.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de invitar contacto autorizado: colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-012"></a>
## MSP-012 · Invitar contacto autorizado: Logística Norte

**Área:** clients · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contacto de Logística Norte sólo accede a recursos permitidos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Invitar una persona con permisos mínimos para Logística Norte: dominio lognorte.example, Barquisimeto.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Invitar una persona con permisos mínimos para Logística Norte: dominio lognorte.example, Barquisimeto.
3. Ejecuta «Invitar contacto autorizado: Logística Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contacto de Logística Norte sólo accede a recursos permitidos.

**Alternativas y errores:**
- Correo de otro cliente: impedir que herede acceso entre espacios.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contacto de Logística Norte sólo accede a recursos permitidos.
- Se reproduce y controla esta alternativa: Correo de otro cliente: impedir que herede acceso entre espacios.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de invitar contacto autorizado: logística norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-013"></a>
## MSP-013 · Cambiar el estado del cliente: Clínica Aurora

**Área:** clients · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de Clínica Aurora se actualiza con autoría y fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Marcar el estado operativo y su motivo para Clínica Aurora: dominio aurora.example, Caracas.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marcar el estado operativo y su motivo para Clínica Aurora: dominio aurora.example, Caracas.
3. Ejecuta «Cambiar el estado del cliente: Clínica Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de Clínica Aurora se actualiza con autoría y fecha.

**Alternativas y errores:**
- Contrato suspendido: limitar cambios a un administrador autorizado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de Clínica Aurora se actualiza con autoría y fecha.
- Se reproduce y controla esta alternativa: Contrato suspendido: limitar cambios a un administrador autorizado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de cambiar el estado del cliente: clínica aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-014"></a>
## MSP-014 · Cambiar el estado del cliente: Hotel Río

**Área:** clients · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de Hotel Río se actualiza con autoría y fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Marcar el estado operativo y su motivo para Hotel Río: dominio hotelrio.example, Valencia.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marcar el estado operativo y su motivo para Hotel Río: dominio hotelrio.example, Valencia.
3. Ejecuta «Cambiar el estado del cliente: Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de Hotel Río se actualiza con autoría y fecha.

**Alternativas y errores:**
- Contrato suspendido: limitar cambios a un administrador autorizado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de Hotel Río se actualiza con autoría y fecha.
- Se reproduce y controla esta alternativa: Contrato suspendido: limitar cambios a un administrador autorizado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de cambiar el estado del cliente: hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-015"></a>
## MSP-015 · Cambiar el estado del cliente: Colegio Sol

**Área:** clients · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de Colegio Sol se actualiza con autoría y fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Marcar el estado operativo y su motivo para Colegio Sol: dominio colegiosol.example, Maracay.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marcar el estado operativo y su motivo para Colegio Sol: dominio colegiosol.example, Maracay.
3. Ejecuta «Cambiar el estado del cliente: Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de Colegio Sol se actualiza con autoría y fecha.

**Alternativas y errores:**
- Contrato suspendido: limitar cambios a un administrador autorizado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de Colegio Sol se actualiza con autoría y fecha.
- Se reproduce y controla esta alternativa: Contrato suspendido: limitar cambios a un administrador autorizado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de cambiar el estado del cliente: colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-016"></a>
## MSP-016 · Cambiar el estado del cliente: Logística Norte

**Área:** clients · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de Logística Norte se actualiza con autoría y fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Marcar el estado operativo y su motivo para Logística Norte: dominio lognorte.example, Barquisimeto.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Marcar el estado operativo y su motivo para Logística Norte: dominio lognorte.example, Barquisimeto.
3. Ejecuta «Cambiar el estado del cliente: Logística Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de Logística Norte se actualiza con autoría y fecha.

**Alternativas y errores:**
- Contrato suspendido: limitar cambios a un administrador autorizado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de Logística Norte se actualiza con autoría y fecha.
- Se reproduce y controla esta alternativa: Contrato suspendido: limitar cambios a un administrador autorizado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de cambiar el estado del cliente: logística norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-017"></a>
## MSP-017 · Exportar inventario de cliente: Clínica Aurora

**Área:** clients · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de Clínica Aurora contiene únicamente activos autorizados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar una exportación dentro del alcance asignado para Clínica Aurora: dominio aurora.example, Caracas.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar una exportación dentro del alcance asignado para Clínica Aurora: dominio aurora.example, Caracas.
3. Ejecuta «Exportar inventario de cliente: Clínica Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de Clínica Aurora contiene únicamente activos autorizados.

**Alternativas y errores:**
- Identificador de otro cliente: denegar sin revelar si existe.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de Clínica Aurora contiene únicamente activos autorizados.
- Se reproduce y controla esta alternativa: Identificador de otro cliente: denegar sin revelar si existe.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de exportar inventario de cliente: clínica aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-018"></a>
## MSP-018 · Exportar inventario de cliente: Hotel Río

**Área:** clients · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de Hotel Río contiene únicamente activos autorizados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar una exportación dentro del alcance asignado para Hotel Río: dominio hotelrio.example, Valencia.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar una exportación dentro del alcance asignado para Hotel Río: dominio hotelrio.example, Valencia.
3. Ejecuta «Exportar inventario de cliente: Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de Hotel Río contiene únicamente activos autorizados.

**Alternativas y errores:**
- Identificador de otro cliente: denegar sin revelar si existe.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de Hotel Río contiene únicamente activos autorizados.
- Se reproduce y controla esta alternativa: Identificador de otro cliente: denegar sin revelar si existe.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de exportar inventario de cliente: hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-019"></a>
## MSP-019 · Exportar inventario de cliente: Colegio Sol

**Área:** clients · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de Colegio Sol contiene únicamente activos autorizados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar una exportación dentro del alcance asignado para Colegio Sol: dominio colegiosol.example, Maracay.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar una exportación dentro del alcance asignado para Colegio Sol: dominio colegiosol.example, Maracay.
3. Ejecuta «Exportar inventario de cliente: Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de Colegio Sol contiene únicamente activos autorizados.

**Alternativas y errores:**
- Identificador de otro cliente: denegar sin revelar si existe.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de Colegio Sol contiene únicamente activos autorizados.
- Se reproduce y controla esta alternativa: Identificador de otro cliente: denegar sin revelar si existe.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de exportar inventario de cliente: colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-020"></a>
## MSP-020 · Exportar inventario de cliente: Logística Norte

**Área:** clients · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de Logística Norte contiene únicamente activos autorizados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar una exportación dentro del alcance asignado para Logística Norte: dominio lognorte.example, Barquisimeto.

**Flujo principal:**
1. El actor accede a clientes y espacios aislados dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar una exportación dentro del alcance asignado para Logística Norte: dominio lognorte.example, Barquisimeto.
3. Ejecuta «Exportar inventario de cliente: Logística Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de Logística Norte contiene únicamente activos autorizados.

**Alternativas y errores:**
- Identificador de otro cliente: denegar sin revelar si existe.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de Logística Norte contiene únicamente activos autorizados.
- Se reproduce y controla esta alternativa: Identificador de otro cliente: denegar sin revelar si existe.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Clientes y espacios aislados: operaciones de exportar inventario de cliente: logística norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-021"></a>
## MSP-021 · Crear contrato administrado: Contrato Aurora Gold

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contrato de Contrato Aurora Gold establece límites y fechas verificables.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar vigencia, moneda y capacidad mensual para Contrato Aurora Gold: USD 4 000 al mes, 24x7.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar vigencia, moneda y capacidad mensual para Contrato Aurora Gold: USD 4 000 al mes, 24x7.
3. Ejecuta «Crear contrato administrado: Contrato Aurora Gold»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contrato de Contrato Aurora Gold establece límites y fechas verificables.

**Alternativas y errores:**
- Vigencia final anterior a la inicial: rechazar el contrato.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contrato de Contrato Aurora Gold establece límites y fechas verificables.
- Se reproduce y controla esta alternativa: Vigencia final anterior a la inicial: rechazar el contrato.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de crear contrato administrado: contrato aurora gold; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-022"></a>
## MSP-022 · Crear contrato administrado: Contrato Hotel Río

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contrato de Contrato Hotel Río establece límites y fechas verificables.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar vigencia, moneda y capacidad mensual para Contrato Hotel Río: 80 horas al mes, horario extendido.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar vigencia, moneda y capacidad mensual para Contrato Hotel Río: 80 horas al mes, horario extendido.
3. Ejecuta «Crear contrato administrado: Contrato Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contrato de Contrato Hotel Río establece límites y fechas verificables.

**Alternativas y errores:**
- Vigencia final anterior a la inicial: rechazar el contrato.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contrato de Contrato Hotel Río establece límites y fechas verificables.
- Se reproduce y controla esta alternativa: Vigencia final anterior a la inicial: rechazar el contrato.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de crear contrato administrado: contrato hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-023"></a>
## MSP-023 · Crear contrato administrado: Contrato Colegio Sol

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contrato de Contrato Colegio Sol establece límites y fechas verificables.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar vigencia, moneda y capacidad mensual para Contrato Colegio Sol: 40 endpoints, horario escolar.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar vigencia, moneda y capacidad mensual para Contrato Colegio Sol: 40 endpoints, horario escolar.
3. Ejecuta «Crear contrato administrado: Contrato Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contrato de Contrato Colegio Sol establece límites y fechas verificables.

**Alternativas y errores:**
- Vigencia final anterior a la inicial: rechazar el contrato.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contrato de Contrato Colegio Sol establece límites y fechas verificables.
- Se reproduce y controla esta alternativa: Vigencia final anterior a la inicial: rechazar el contrato.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de crear contrato administrado: contrato colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-024"></a>
## MSP-024 · Crear contrato administrado: Contrato Norte

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El contrato de Contrato Norte establece límites y fechas verificables.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar vigencia, moneda y capacidad mensual para Contrato Norte: soporte de red, lunes a sábado.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar vigencia, moneda y capacidad mensual para Contrato Norte: soporte de red, lunes a sábado.
3. Ejecuta «Crear contrato administrado: Contrato Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El contrato de Contrato Norte establece límites y fechas verificables.

**Alternativas y errores:**
- Vigencia final anterior a la inicial: rechazar el contrato.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El contrato de Contrato Norte establece límites y fechas verificables.
- Se reproduce y controla esta alternativa: Vigencia final anterior a la inicial: rechazar el contrato.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de crear contrato administrado: contrato norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-025"></a>
## MSP-025 · Asignar paquete de servicio: Contrato Aurora Gold

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El paquete de Contrato Aurora Gold refleja los servicios contratados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar cobertura, horario y canales incluidos para Contrato Aurora Gold: USD 4 000 al mes, 24x7.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar cobertura, horario y canales incluidos para Contrato Aurora Gold: USD 4 000 al mes, 24x7.
3. Ejecuta «Asignar paquete de servicio: Contrato Aurora Gold»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El paquete de Contrato Aurora Gold refleja los servicios contratados.

**Alternativas y errores:**
- Servicio fuera del catálogo vigente: impedir su activación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El paquete de Contrato Aurora Gold refleja los servicios contratados.
- Se reproduce y controla esta alternativa: Servicio fuera del catálogo vigente: impedir su activación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de asignar paquete de servicio: contrato aurora gold; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-026"></a>
## MSP-026 · Asignar paquete de servicio: Contrato Hotel Río

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El paquete de Contrato Hotel Río refleja los servicios contratados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar cobertura, horario y canales incluidos para Contrato Hotel Río: 80 horas al mes, horario extendido.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar cobertura, horario y canales incluidos para Contrato Hotel Río: 80 horas al mes, horario extendido.
3. Ejecuta «Asignar paquete de servicio: Contrato Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El paquete de Contrato Hotel Río refleja los servicios contratados.

**Alternativas y errores:**
- Servicio fuera del catálogo vigente: impedir su activación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El paquete de Contrato Hotel Río refleja los servicios contratados.
- Se reproduce y controla esta alternativa: Servicio fuera del catálogo vigente: impedir su activación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de asignar paquete de servicio: contrato hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-027"></a>
## MSP-027 · Asignar paquete de servicio: Contrato Colegio Sol

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El paquete de Contrato Colegio Sol refleja los servicios contratados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar cobertura, horario y canales incluidos para Contrato Colegio Sol: 40 endpoints, horario escolar.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar cobertura, horario y canales incluidos para Contrato Colegio Sol: 40 endpoints, horario escolar.
3. Ejecuta «Asignar paquete de servicio: Contrato Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El paquete de Contrato Colegio Sol refleja los servicios contratados.

**Alternativas y errores:**
- Servicio fuera del catálogo vigente: impedir su activación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El paquete de Contrato Colegio Sol refleja los servicios contratados.
- Se reproduce y controla esta alternativa: Servicio fuera del catálogo vigente: impedir su activación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de asignar paquete de servicio: contrato colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-028"></a>
## MSP-028 · Asignar paquete de servicio: Contrato Norte

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El paquete de Contrato Norte refleja los servicios contratados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar cobertura, horario y canales incluidos para Contrato Norte: soporte de red, lunes a sábado.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar cobertura, horario y canales incluidos para Contrato Norte: soporte de red, lunes a sábado.
3. Ejecuta «Asignar paquete de servicio: Contrato Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El paquete de Contrato Norte refleja los servicios contratados.

**Alternativas y errores:**
- Servicio fuera del catálogo vigente: impedir su activación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El paquete de Contrato Norte refleja los servicios contratados.
- Se reproduce y controla esta alternativa: Servicio fuera del catálogo vigente: impedir su activación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de asignar paquete de servicio: contrato norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-029"></a>
## MSP-029 · Definir objetivo de respuesta: Contrato Aurora Gold

**Área:** contracts · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El SLA de Contrato Aurora Gold calcula plazos con el calendario del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar objetivos por prioridad y horario hábil para Contrato Aurora Gold: USD 4 000 al mes, 24x7.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar objetivos por prioridad y horario hábil para Contrato Aurora Gold: USD 4 000 al mes, 24x7.
3. Ejecuta «Definir objetivo de respuesta: Contrato Aurora Gold»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El SLA de Contrato Aurora Gold calcula plazos con el calendario del cliente.

**Alternativas y errores:**
- Objetivo nulo o negativo: no guardar la política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El SLA de Contrato Aurora Gold calcula plazos con el calendario del cliente.
- Se reproduce y controla esta alternativa: Objetivo nulo o negativo: no guardar la política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de definir objetivo de respuesta: contrato aurora gold; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-030"></a>
## MSP-030 · Definir objetivo de respuesta: Contrato Hotel Río

**Área:** contracts · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El SLA de Contrato Hotel Río calcula plazos con el calendario del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar objetivos por prioridad y horario hábil para Contrato Hotel Río: 80 horas al mes, horario extendido.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar objetivos por prioridad y horario hábil para Contrato Hotel Río: 80 horas al mes, horario extendido.
3. Ejecuta «Definir objetivo de respuesta: Contrato Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El SLA de Contrato Hotel Río calcula plazos con el calendario del cliente.

**Alternativas y errores:**
- Objetivo nulo o negativo: no guardar la política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El SLA de Contrato Hotel Río calcula plazos con el calendario del cliente.
- Se reproduce y controla esta alternativa: Objetivo nulo o negativo: no guardar la política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de definir objetivo de respuesta: contrato hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-031"></a>
## MSP-031 · Definir objetivo de respuesta: Contrato Colegio Sol

**Área:** contracts · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El SLA de Contrato Colegio Sol calcula plazos con el calendario del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar objetivos por prioridad y horario hábil para Contrato Colegio Sol: 40 endpoints, horario escolar.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar objetivos por prioridad y horario hábil para Contrato Colegio Sol: 40 endpoints, horario escolar.
3. Ejecuta «Definir objetivo de respuesta: Contrato Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El SLA de Contrato Colegio Sol calcula plazos con el calendario del cliente.

**Alternativas y errores:**
- Objetivo nulo o negativo: no guardar la política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El SLA de Contrato Colegio Sol calcula plazos con el calendario del cliente.
- Se reproduce y controla esta alternativa: Objetivo nulo o negativo: no guardar la política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de definir objetivo de respuesta: contrato colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-032"></a>
## MSP-032 · Definir objetivo de respuesta: Contrato Norte

**Área:** contracts · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El SLA de Contrato Norte calcula plazos con el calendario del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar objetivos por prioridad y horario hábil para Contrato Norte: soporte de red, lunes a sábado.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar objetivos por prioridad y horario hábil para Contrato Norte: soporte de red, lunes a sábado.
3. Ejecuta «Definir objetivo de respuesta: Contrato Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El SLA de Contrato Norte calcula plazos con el calendario del cliente.

**Alternativas y errores:**
- Objetivo nulo o negativo: no guardar la política.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El SLA de Contrato Norte calcula plazos con el calendario del cliente.
- Se reproduce y controla esta alternativa: Objetivo nulo o negativo: no guardar la política.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de definir objetivo de respuesta: contrato norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-033"></a>
## MSP-033 · Renovar contrato con cambios: Contrato Aurora Gold

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La renovación de Contrato Aurora Gold mantiene el historial y aplica desde su fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Preparar una nueva versión para el próximo periodo para Contrato Aurora Gold: USD 4 000 al mes, 24x7.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Preparar una nueva versión para el próximo periodo para Contrato Aurora Gold: USD 4 000 al mes, 24x7.
3. Ejecuta «Renovar contrato con cambios: Contrato Aurora Gold»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La renovación de Contrato Aurora Gold mantiene el historial y aplica desde su fecha.

**Alternativas y errores:**
- Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La renovación de Contrato Aurora Gold mantiene el historial y aplica desde su fecha.
- Se reproduce y controla esta alternativa: Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de renovar contrato con cambios: contrato aurora gold; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-034"></a>
## MSP-034 · Renovar contrato con cambios: Contrato Hotel Río

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La renovación de Contrato Hotel Río mantiene el historial y aplica desde su fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Preparar una nueva versión para el próximo periodo para Contrato Hotel Río: 80 horas al mes, horario extendido.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Preparar una nueva versión para el próximo periodo para Contrato Hotel Río: 80 horas al mes, horario extendido.
3. Ejecuta «Renovar contrato con cambios: Contrato Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La renovación de Contrato Hotel Río mantiene el historial y aplica desde su fecha.

**Alternativas y errores:**
- Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La renovación de Contrato Hotel Río mantiene el historial y aplica desde su fecha.
- Se reproduce y controla esta alternativa: Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de renovar contrato con cambios: contrato hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-035"></a>
## MSP-035 · Renovar contrato con cambios: Contrato Colegio Sol

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La renovación de Contrato Colegio Sol mantiene el historial y aplica desde su fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Preparar una nueva versión para el próximo periodo para Contrato Colegio Sol: 40 endpoints, horario escolar.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Preparar una nueva versión para el próximo periodo para Contrato Colegio Sol: 40 endpoints, horario escolar.
3. Ejecuta «Renovar contrato con cambios: Contrato Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La renovación de Contrato Colegio Sol mantiene el historial y aplica desde su fecha.

**Alternativas y errores:**
- Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La renovación de Contrato Colegio Sol mantiene el historial y aplica desde su fecha.
- Se reproduce y controla esta alternativa: Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de renovar contrato con cambios: contrato colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-036"></a>
## MSP-036 · Renovar contrato con cambios: Contrato Norte

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La renovación de Contrato Norte mantiene el historial y aplica desde su fecha.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Preparar una nueva versión para el próximo periodo para Contrato Norte: soporte de red, lunes a sábado.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Preparar una nueva versión para el próximo periodo para Contrato Norte: soporte de red, lunes a sábado.
3. Ejecuta «Renovar contrato con cambios: Contrato Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La renovación de Contrato Norte mantiene el historial y aplica desde su fecha.

**Alternativas y errores:**
- Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La renovación de Contrato Norte mantiene el historial y aplica desde su fecha.
- Se reproduce y controla esta alternativa: Dos versiones activas en el mismo periodo: exigir resolución del conflicto.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de renovar contrato con cambios: contrato norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-037"></a>
## MSP-037 · Suspender servicio por vencimiento: Contrato Aurora Gold

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de Contrato Aurora Gold cambia según la política publicada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar autorización y fecha efectiva de suspensión para Contrato Aurora Gold: USD 4 000 al mes, 24x7.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar autorización y fecha efectiva de suspensión para Contrato Aurora Gold: USD 4 000 al mes, 24x7.
3. Ejecuta «Suspender servicio por vencimiento: Contrato Aurora Gold»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de Contrato Aurora Gold cambia según la política publicada.

**Alternativas y errores:**
- Solicitud sin autorización contractual: mantener el servicio activo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de Contrato Aurora Gold cambia según la política publicada.
- Se reproduce y controla esta alternativa: Solicitud sin autorización contractual: mantener el servicio activo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de suspender servicio por vencimiento: contrato aurora gold; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-038"></a>
## MSP-038 · Suspender servicio por vencimiento: Contrato Hotel Río

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de Contrato Hotel Río cambia según la política publicada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar autorización y fecha efectiva de suspensión para Contrato Hotel Río: 80 horas al mes, horario extendido.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar autorización y fecha efectiva de suspensión para Contrato Hotel Río: 80 horas al mes, horario extendido.
3. Ejecuta «Suspender servicio por vencimiento: Contrato Hotel Río»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de Contrato Hotel Río cambia según la política publicada.

**Alternativas y errores:**
- Solicitud sin autorización contractual: mantener el servicio activo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de Contrato Hotel Río cambia según la política publicada.
- Se reproduce y controla esta alternativa: Solicitud sin autorización contractual: mantener el servicio activo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de suspender servicio por vencimiento: contrato hotel río; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-039"></a>
## MSP-039 · Suspender servicio por vencimiento: Contrato Colegio Sol

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de Contrato Colegio Sol cambia según la política publicada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar autorización y fecha efectiva de suspensión para Contrato Colegio Sol: 40 endpoints, horario escolar.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar autorización y fecha efectiva de suspensión para Contrato Colegio Sol: 40 endpoints, horario escolar.
3. Ejecuta «Suspender servicio por vencimiento: Contrato Colegio Sol»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de Contrato Colegio Sol cambia según la política publicada.

**Alternativas y errores:**
- Solicitud sin autorización contractual: mantener el servicio activo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de Contrato Colegio Sol cambia según la política publicada.
- Se reproduce y controla esta alternativa: Solicitud sin autorización contractual: mantener el servicio activo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de suspender servicio por vencimiento: contrato colegio sol; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-040"></a>
## MSP-040 · Suspender servicio por vencimiento: Contrato Norte

**Área:** contracts · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de Contrato Norte cambia según la política publicada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar autorización y fecha efectiva de suspensión para Contrato Norte: soporte de red, lunes a sábado.

**Flujo principal:**
1. El actor accede a contratos, catálogo y sla dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar autorización y fecha efectiva de suspensión para Contrato Norte: soporte de red, lunes a sábado.
3. Ejecuta «Suspender servicio por vencimiento: Contrato Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de Contrato Norte cambia según la política publicada.

**Alternativas y errores:**
- Solicitud sin autorización contractual: mantener el servicio activo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de Contrato Norte cambia según la política publicada.
- Se reproduce y controla esta alternativa: Solicitud sin autorización contractual: mantener el servicio activo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Contratos, catálogo y SLA: operaciones de suspender servicio por vencimiento: contrato norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-041"></a>
## MSP-041 · Registrar incidente del cliente: caída de VPN

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El incidente de caída de VPN se abre en la cola de su cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Clasificar síntoma, impacto y cliente afectado para caída de VPN: Cliente Aurora, sede Caracas.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Clasificar síntoma, impacto y cliente afectado para caída de VPN: Cliente Aurora, sede Caracas.
3. Ejecuta «Registrar incidente del cliente: caída de VPN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El incidente de caída de VPN se abre en la cola de su cliente.

**Alternativas y errores:**
- Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El incidente de caída de VPN se abre en la cola de su cliente.
- Se reproduce y controla esta alternativa: Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de registrar incidente del cliente: caída de vpn; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-042"></a>
## MSP-042 · Registrar incidente del cliente: correo bloqueado

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El incidente de correo bloqueado se abre en la cola de su cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Clasificar síntoma, impacto y cliente afectado para correo bloqueado: Hotel Río, equipo de recepción.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Clasificar síntoma, impacto y cliente afectado para correo bloqueado: Hotel Río, equipo de recepción.
3. Ejecuta «Registrar incidente del cliente: correo bloqueado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El incidente de correo bloqueado se abre en la cola de su cliente.

**Alternativas y errores:**
- Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El incidente de correo bloqueado se abre en la cola de su cliente.
- Se reproduce y controla esta alternativa: Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de registrar incidente del cliente: correo bloqueado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-043"></a>
## MSP-043 · Registrar incidente del cliente: respaldo retrasado

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El incidente de respaldo retrasado se abre en la cola de su cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Clasificar síntoma, impacto y cliente afectado para respaldo retrasado: Colegio Sol, servidor académico.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Clasificar síntoma, impacto y cliente afectado para respaldo retrasado: Colegio Sol, servidor académico.
3. Ejecuta «Registrar incidente del cliente: respaldo retrasado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El incidente de respaldo retrasado se abre en la cola de su cliente.

**Alternativas y errores:**
- Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El incidente de respaldo retrasado se abre en la cola de su cliente.
- Se reproduce y controla esta alternativa: Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de registrar incidente del cliente: respaldo retrasado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-044"></a>
## MSP-044 · Registrar incidente del cliente: latencia WAN

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El incidente de latencia WAN se abre en la cola de su cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Clasificar síntoma, impacto y cliente afectado para latencia WAN: Logística Norte, centro de despacho.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Clasificar síntoma, impacto y cliente afectado para latencia WAN: Logística Norte, centro de despacho.
3. Ejecuta «Registrar incidente del cliente: latencia WAN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El incidente de latencia WAN se abre en la cola de su cliente.

**Alternativas y errores:**
- Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El incidente de latencia WAN se abre en la cola de su cliente.
- Se reproduce y controla esta alternativa: Cliente desconocido o inactivo: no crear una solicitud huérfana.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de registrar incidente del cliente: latencia wan; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-045"></a>
## MSP-045 · Asignar ticket a especialista: caída de VPN

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La asignación de caída de VPN queda visible para la cola autorizada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir técnico habilitado para el servicio requerido para caída de VPN: Cliente Aurora, sede Caracas.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir técnico habilitado para el servicio requerido para caída de VPN: Cliente Aurora, sede Caracas.
3. Ejecuta «Asignar ticket a especialista: caída de VPN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La asignación de caída de VPN queda visible para la cola autorizada.

**Alternativas y errores:**
- Técnico sin acceso a ese cliente: rechazar la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La asignación de caída de VPN queda visible para la cola autorizada.
- Se reproduce y controla esta alternativa: Técnico sin acceso a ese cliente: rechazar la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de asignar ticket a especialista: caída de vpn; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-046"></a>
## MSP-046 · Asignar ticket a especialista: correo bloqueado

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La asignación de correo bloqueado queda visible para la cola autorizada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir técnico habilitado para el servicio requerido para correo bloqueado: Hotel Río, equipo de recepción.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir técnico habilitado para el servicio requerido para correo bloqueado: Hotel Río, equipo de recepción.
3. Ejecuta «Asignar ticket a especialista: correo bloqueado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La asignación de correo bloqueado queda visible para la cola autorizada.

**Alternativas y errores:**
- Técnico sin acceso a ese cliente: rechazar la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La asignación de correo bloqueado queda visible para la cola autorizada.
- Se reproduce y controla esta alternativa: Técnico sin acceso a ese cliente: rechazar la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de asignar ticket a especialista: correo bloqueado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-047"></a>
## MSP-047 · Asignar ticket a especialista: respaldo retrasado

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La asignación de respaldo retrasado queda visible para la cola autorizada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir técnico habilitado para el servicio requerido para respaldo retrasado: Colegio Sol, servidor académico.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir técnico habilitado para el servicio requerido para respaldo retrasado: Colegio Sol, servidor académico.
3. Ejecuta «Asignar ticket a especialista: respaldo retrasado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La asignación de respaldo retrasado queda visible para la cola autorizada.

**Alternativas y errores:**
- Técnico sin acceso a ese cliente: rechazar la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La asignación de respaldo retrasado queda visible para la cola autorizada.
- Se reproduce y controla esta alternativa: Técnico sin acceso a ese cliente: rechazar la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de asignar ticket a especialista: respaldo retrasado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-048"></a>
## MSP-048 · Asignar ticket a especialista: latencia WAN

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La asignación de latencia WAN queda visible para la cola autorizada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir técnico habilitado para el servicio requerido para latencia WAN: Logística Norte, centro de despacho.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir técnico habilitado para el servicio requerido para latencia WAN: Logística Norte, centro de despacho.
3. Ejecuta «Asignar ticket a especialista: latencia WAN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La asignación de latencia WAN queda visible para la cola autorizada.

**Alternativas y errores:**
- Técnico sin acceso a ese cliente: rechazar la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La asignación de latencia WAN queda visible para la cola autorizada.
- Se reproduce y controla esta alternativa: Técnico sin acceso a ese cliente: rechazar la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de asignar ticket a especialista: latencia wan; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-049"></a>
## MSP-049 · Escalar incumplimiento de SLA: caída de VPN

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La escalación de caída de VPN registra regla, responsable y plazo.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar atención de segundo nivel antes del vencimiento para caída de VPN: Cliente Aurora, sede Caracas.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar atención de segundo nivel antes del vencimiento para caída de VPN: Cliente Aurora, sede Caracas.
3. Ejecuta «Escalar incumplimiento de SLA: caída de VPN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La escalación de caída de VPN registra regla, responsable y plazo.

**Alternativas y errores:**
- SLA pausado según contrato: calcular con la pausa explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La escalación de caída de VPN registra regla, responsable y plazo.
- Se reproduce y controla esta alternativa: SLA pausado según contrato: calcular con la pausa explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de escalar incumplimiento de sla: caída de vpn; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-050"></a>
## MSP-050 · Escalar incumplimiento de SLA: correo bloqueado

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La escalación de correo bloqueado registra regla, responsable y plazo.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar atención de segundo nivel antes del vencimiento para correo bloqueado: Hotel Río, equipo de recepción.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar atención de segundo nivel antes del vencimiento para correo bloqueado: Hotel Río, equipo de recepción.
3. Ejecuta «Escalar incumplimiento de SLA: correo bloqueado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La escalación de correo bloqueado registra regla, responsable y plazo.

**Alternativas y errores:**
- SLA pausado según contrato: calcular con la pausa explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La escalación de correo bloqueado registra regla, responsable y plazo.
- Se reproduce y controla esta alternativa: SLA pausado según contrato: calcular con la pausa explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de escalar incumplimiento de sla: correo bloqueado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-051"></a>
## MSP-051 · Escalar incumplimiento de SLA: respaldo retrasado

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La escalación de respaldo retrasado registra regla, responsable y plazo.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar atención de segundo nivel antes del vencimiento para respaldo retrasado: Colegio Sol, servidor académico.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar atención de segundo nivel antes del vencimiento para respaldo retrasado: Colegio Sol, servidor académico.
3. Ejecuta «Escalar incumplimiento de SLA: respaldo retrasado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La escalación de respaldo retrasado registra regla, responsable y plazo.

**Alternativas y errores:**
- SLA pausado según contrato: calcular con la pausa explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La escalación de respaldo retrasado registra regla, responsable y plazo.
- Se reproduce y controla esta alternativa: SLA pausado según contrato: calcular con la pausa explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de escalar incumplimiento de sla: respaldo retrasado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-052"></a>
## MSP-052 · Escalar incumplimiento de SLA: latencia WAN

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La escalación de latencia WAN registra regla, responsable y plazo.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Solicitar atención de segundo nivel antes del vencimiento para latencia WAN: Logística Norte, centro de despacho.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Solicitar atención de segundo nivel antes del vencimiento para latencia WAN: Logística Norte, centro de despacho.
3. Ejecuta «Escalar incumplimiento de SLA: latencia WAN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La escalación de latencia WAN registra regla, responsable y plazo.

**Alternativas y errores:**
- SLA pausado según contrato: calcular con la pausa explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La escalación de latencia WAN registra regla, responsable y plazo.
- Se reproduce y controla esta alternativa: SLA pausado según contrato: calcular con la pausa explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de escalar incumplimiento de sla: latencia wan; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-053"></a>
## MSP-053 · Coordinar incidente mayor: caída de VPN

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de caída de VPN reúne tareas sin mezclar clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Abrir sala de coordinación y actualización al contacto para caída de VPN: Cliente Aurora, sede Caracas.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Abrir sala de coordinación y actualización al contacto para caída de VPN: Cliente Aurora, sede Caracas.
3. Ejecuta «Coordinar incidente mayor: caída de VPN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de caída de VPN reúne tareas sin mezclar clientes.

**Alternativas y errores:**
- Impacto no confirmado: no notificar como incidente general.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de caída de VPN reúne tareas sin mezclar clientes.
- Se reproduce y controla esta alternativa: Impacto no confirmado: no notificar como incidente general.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de coordinar incidente mayor: caída de vpn; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-054"></a>
## MSP-054 · Coordinar incidente mayor: correo bloqueado

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de correo bloqueado reúne tareas sin mezclar clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Abrir sala de coordinación y actualización al contacto para correo bloqueado: Hotel Río, equipo de recepción.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Abrir sala de coordinación y actualización al contacto para correo bloqueado: Hotel Río, equipo de recepción.
3. Ejecuta «Coordinar incidente mayor: correo bloqueado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de correo bloqueado reúne tareas sin mezclar clientes.

**Alternativas y errores:**
- Impacto no confirmado: no notificar como incidente general.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de correo bloqueado reúne tareas sin mezclar clientes.
- Se reproduce y controla esta alternativa: Impacto no confirmado: no notificar como incidente general.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de coordinar incidente mayor: correo bloqueado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-055"></a>
## MSP-055 · Coordinar incidente mayor: respaldo retrasado

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de respaldo retrasado reúne tareas sin mezclar clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Abrir sala de coordinación y actualización al contacto para respaldo retrasado: Colegio Sol, servidor académico.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Abrir sala de coordinación y actualización al contacto para respaldo retrasado: Colegio Sol, servidor académico.
3. Ejecuta «Coordinar incidente mayor: respaldo retrasado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de respaldo retrasado reúne tareas sin mezclar clientes.

**Alternativas y errores:**
- Impacto no confirmado: no notificar como incidente general.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de respaldo retrasado reúne tareas sin mezclar clientes.
- Se reproduce y controla esta alternativa: Impacto no confirmado: no notificar como incidente general.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de coordinar incidente mayor: respaldo retrasado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-056"></a>
## MSP-056 · Coordinar incidente mayor: latencia WAN

**Área:** service-desk · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de latencia WAN reúne tareas sin mezclar clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Abrir sala de coordinación y actualización al contacto para latencia WAN: Logística Norte, centro de despacho.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Abrir sala de coordinación y actualización al contacto para latencia WAN: Logística Norte, centro de despacho.
3. Ejecuta «Coordinar incidente mayor: latencia WAN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de latencia WAN reúne tareas sin mezclar clientes.

**Alternativas y errores:**
- Impacto no confirmado: no notificar como incidente general.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de latencia WAN reúne tareas sin mezclar clientes.
- Se reproduce y controla esta alternativa: Impacto no confirmado: no notificar como incidente general.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de coordinar incidente mayor: latencia wan; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-057"></a>
## MSP-057 · Cerrar solicitud con solución: caída de VPN

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El cierre de caída de VPN conserva el historial y notifica al contacto autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Documentar causa, acción y evidencia de resolución para caída de VPN: Cliente Aurora, sede Caracas.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Documentar causa, acción y evidencia de resolución para caída de VPN: Cliente Aurora, sede Caracas.
3. Ejecuta «Cerrar solicitud con solución: caída de VPN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El cierre de caída de VPN conserva el historial y notifica al contacto autorizado.

**Alternativas y errores:**
- Trabajo incompleto: impedir cierre y conservar el estado actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El cierre de caída de VPN conserva el historial y notifica al contacto autorizado.
- Se reproduce y controla esta alternativa: Trabajo incompleto: impedir cierre y conservar el estado actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de cerrar solicitud con solución: caída de vpn; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-058"></a>
## MSP-058 · Cerrar solicitud con solución: correo bloqueado

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El cierre de correo bloqueado conserva el historial y notifica al contacto autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Documentar causa, acción y evidencia de resolución para correo bloqueado: Hotel Río, equipo de recepción.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Documentar causa, acción y evidencia de resolución para correo bloqueado: Hotel Río, equipo de recepción.
3. Ejecuta «Cerrar solicitud con solución: correo bloqueado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El cierre de correo bloqueado conserva el historial y notifica al contacto autorizado.

**Alternativas y errores:**
- Trabajo incompleto: impedir cierre y conservar el estado actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El cierre de correo bloqueado conserva el historial y notifica al contacto autorizado.
- Se reproduce y controla esta alternativa: Trabajo incompleto: impedir cierre y conservar el estado actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de cerrar solicitud con solución: correo bloqueado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-059"></a>
## MSP-059 · Cerrar solicitud con solución: respaldo retrasado

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El cierre de respaldo retrasado conserva el historial y notifica al contacto autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Documentar causa, acción y evidencia de resolución para respaldo retrasado: Colegio Sol, servidor académico.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Documentar causa, acción y evidencia de resolución para respaldo retrasado: Colegio Sol, servidor académico.
3. Ejecuta «Cerrar solicitud con solución: respaldo retrasado»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El cierre de respaldo retrasado conserva el historial y notifica al contacto autorizado.

**Alternativas y errores:**
- Trabajo incompleto: impedir cierre y conservar el estado actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El cierre de respaldo retrasado conserva el historial y notifica al contacto autorizado.
- Se reproduce y controla esta alternativa: Trabajo incompleto: impedir cierre y conservar el estado actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de cerrar solicitud con solución: respaldo retrasado; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-060"></a>
## MSP-060 · Cerrar solicitud con solución: latencia WAN

**Área:** service-desk · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El cierre de latencia WAN conserva el historial y notifica al contacto autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Documentar causa, acción y evidencia de resolución para latencia WAN: Logística Norte, centro de despacho.

**Flujo principal:**
1. El actor accede a mesa de servicio para varios clientes dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Documentar causa, acción y evidencia de resolución para latencia WAN: Logística Norte, centro de despacho.
3. Ejecuta «Cerrar solicitud con solución: latencia WAN»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El cierre de latencia WAN conserva el historial y notifica al contacto autorizado.

**Alternativas y errores:**
- Trabajo incompleto: impedir cierre y conservar el estado actual.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El cierre de latencia WAN conserva el historial y notifica al contacto autorizado.
- Se reproduce y controla esta alternativa: Trabajo incompleto: impedir cierre y conservar el estado actual.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Mesa de servicio para varios clientes: operaciones de cerrar solicitud con solución: latencia wan; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-061"></a>
## MSP-061 · Crear turno de guardia: guardia de redes

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La guardia de guardia de redes queda asignada con cobertura visible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir horario, zona y nivel de escalación para guardia de redes: zona America/Caracas, 18:00–06:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir horario, zona y nivel de escalación para guardia de redes: zona America/Caracas, 18:00–06:00.
3. Ejecuta «Crear turno de guardia: guardia de redes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La guardia de guardia de redes queda asignada con cobertura visible.

**Alternativas y errores:**
- Dos turnos incompatibles: pedir resolución antes de publicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La guardia de guardia de redes queda asignada con cobertura visible.
- Se reproduce y controla esta alternativa: Dos turnos incompatibles: pedir resolución antes de publicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de crear turno de guardia: guardia de redes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-062"></a>
## MSP-062 · Crear turno de guardia: equipo nube

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La guardia de equipo nube queda asignada con cobertura visible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir horario, zona y nivel de escalación para equipo nube: dos técnicos certificados AZ-104.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir horario, zona y nivel de escalación para equipo nube: dos técnicos certificados AZ-104.
3. Ejecuta «Crear turno de guardia: equipo nube»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La guardia de equipo nube queda asignada con cobertura visible.

**Alternativas y errores:**
- Dos turnos incompatibles: pedir resolución antes de publicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La guardia de equipo nube queda asignada con cobertura visible.
- Se reproduce y controla esta alternativa: Dos turnos incompatibles: pedir resolución antes de publicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de crear turno de guardia: equipo nube; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-063"></a>
## MSP-063 · Crear turno de guardia: turno fin de semana

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La guardia de turno fin de semana queda asignada con cobertura visible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir horario, zona y nivel de escalación para turno fin de semana: sábado 08:00–20:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir horario, zona y nivel de escalación para turno fin de semana: sábado 08:00–20:00.
3. Ejecuta «Crear turno de guardia: turno fin de semana»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La guardia de turno fin de semana queda asignada con cobertura visible.

**Alternativas y errores:**
- Dos turnos incompatibles: pedir resolución antes de publicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La guardia de turno fin de semana queda asignada con cobertura visible.
- Se reproduce y controla esta alternativa: Dos turnos incompatibles: pedir resolución antes de publicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de crear turno de guardia: turno fin de semana; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-064"></a>
## MSP-064 · Crear turno de guardia: equipo de campo

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La guardia de equipo de campo queda asignada con cobertura visible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir horario, zona y nivel de escalación para equipo de campo: dos sedes y cuatro visitas pendientes.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir horario, zona y nivel de escalación para equipo de campo: dos sedes y cuatro visitas pendientes.
3. Ejecuta «Crear turno de guardia: equipo de campo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La guardia de equipo de campo queda asignada con cobertura visible.

**Alternativas y errores:**
- Dos turnos incompatibles: pedir resolución antes de publicar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La guardia de equipo de campo queda asignada con cobertura visible.
- Se reproduce y controla esta alternativa: Dos turnos incompatibles: pedir resolución antes de publicar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de crear turno de guardia: equipo de campo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-065"></a>
## MSP-065 · Asignar técnico por habilidad: guardia de redes

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de guardia de redes corresponde a un técnico elegible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar personal por certificación y disponibilidad para guardia de redes: zona America/Caracas, 18:00–06:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar personal por certificación y disponibilidad para guardia de redes: zona America/Caracas, 18:00–06:00.
3. Ejecuta «Asignar técnico por habilidad: guardia de redes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de guardia de redes corresponde a un técnico elegible.

**Alternativas y errores:**
- Certificación vencida: excluir al técnico de la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de guardia de redes corresponde a un técnico elegible.
- Se reproduce y controla esta alternativa: Certificación vencida: excluir al técnico de la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de asignar técnico por habilidad: guardia de redes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-066"></a>
## MSP-066 · Asignar técnico por habilidad: equipo nube

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de equipo nube corresponde a un técnico elegible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar personal por certificación y disponibilidad para equipo nube: dos técnicos certificados AZ-104.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar personal por certificación y disponibilidad para equipo nube: dos técnicos certificados AZ-104.
3. Ejecuta «Asignar técnico por habilidad: equipo nube»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de equipo nube corresponde a un técnico elegible.

**Alternativas y errores:**
- Certificación vencida: excluir al técnico de la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de equipo nube corresponde a un técnico elegible.
- Se reproduce y controla esta alternativa: Certificación vencida: excluir al técnico de la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de asignar técnico por habilidad: equipo nube; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-067"></a>
## MSP-067 · Asignar técnico por habilidad: turno fin de semana

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de turno fin de semana corresponde a un técnico elegible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar personal por certificación y disponibilidad para turno fin de semana: sábado 08:00–20:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar personal por certificación y disponibilidad para turno fin de semana: sábado 08:00–20:00.
3. Ejecuta «Asignar técnico por habilidad: turno fin de semana»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de turno fin de semana corresponde a un técnico elegible.

**Alternativas y errores:**
- Certificación vencida: excluir al técnico de la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de turno fin de semana corresponde a un técnico elegible.
- Se reproduce y controla esta alternativa: Certificación vencida: excluir al técnico de la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de asignar técnico por habilidad: turno fin de semana; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-068"></a>
## MSP-068 · Asignar técnico por habilidad: equipo de campo

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La cobertura de equipo de campo corresponde a un técnico elegible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar personal por certificación y disponibilidad para equipo de campo: dos sedes y cuatro visitas pendientes.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar personal por certificación y disponibilidad para equipo de campo: dos sedes y cuatro visitas pendientes.
3. Ejecuta «Asignar técnico por habilidad: equipo de campo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La cobertura de equipo de campo corresponde a un técnico elegible.

**Alternativas y errores:**
- Certificación vencida: excluir al técnico de la asignación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La cobertura de equipo de campo corresponde a un técnico elegible.
- Se reproduce y controla esta alternativa: Certificación vencida: excluir al técnico de la asignación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de asignar técnico por habilidad: equipo de campo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-069"></a>
## MSP-069 · Registrar disponibilidad: guardia de redes

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La agenda de guardia de redes actualiza la capacidad de planificación.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Indicar ausencia o cambio de disponibilidad para guardia de redes: zona America/Caracas, 18:00–06:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Indicar ausencia o cambio de disponibilidad para guardia de redes: zona America/Caracas, 18:00–06:00.
3. Ejecuta «Registrar disponibilidad: guardia de redes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La agenda de guardia de redes actualiza la capacidad de planificación.

**Alternativas y errores:**
- Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La agenda de guardia de redes actualiza la capacidad de planificación.
- Se reproduce y controla esta alternativa: Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de registrar disponibilidad: guardia de redes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-070"></a>
## MSP-070 · Registrar disponibilidad: equipo nube

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La agenda de equipo nube actualiza la capacidad de planificación.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Indicar ausencia o cambio de disponibilidad para equipo nube: dos técnicos certificados AZ-104.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Indicar ausencia o cambio de disponibilidad para equipo nube: dos técnicos certificados AZ-104.
3. Ejecuta «Registrar disponibilidad: equipo nube»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La agenda de equipo nube actualiza la capacidad de planificación.

**Alternativas y errores:**
- Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La agenda de equipo nube actualiza la capacidad de planificación.
- Se reproduce y controla esta alternativa: Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de registrar disponibilidad: equipo nube; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-071"></a>
## MSP-071 · Registrar disponibilidad: turno fin de semana

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La agenda de turno fin de semana actualiza la capacidad de planificación.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Indicar ausencia o cambio de disponibilidad para turno fin de semana: sábado 08:00–20:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Indicar ausencia o cambio de disponibilidad para turno fin de semana: sábado 08:00–20:00.
3. Ejecuta «Registrar disponibilidad: turno fin de semana»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La agenda de turno fin de semana actualiza la capacidad de planificación.

**Alternativas y errores:**
- Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La agenda de turno fin de semana actualiza la capacidad de planificación.
- Se reproduce y controla esta alternativa: Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de registrar disponibilidad: turno fin de semana; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-072"></a>
## MSP-072 · Registrar disponibilidad: equipo de campo

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La agenda de equipo de campo actualiza la capacidad de planificación.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Indicar ausencia o cambio de disponibilidad para equipo de campo: dos sedes y cuatro visitas pendientes.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Indicar ausencia o cambio de disponibilidad para equipo de campo: dos sedes y cuatro visitas pendientes.
3. Ejecuta «Registrar disponibilidad: equipo de campo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La agenda de equipo de campo actualiza la capacidad de planificación.

**Alternativas y errores:**
- Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La agenda de equipo de campo actualiza la capacidad de planificación.
- Se reproduce y controla esta alternativa: Periodo superpuesto con guardia confirmada: alertar y mantener el turno.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de registrar disponibilidad: equipo de campo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-073"></a>
## MSP-073 · Distribuir carga entre equipos: guardia de redes

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La carga de guardia de redes se reparte sin exceder límites publicados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar tickets abiertos y capacidad comprometida para guardia de redes: zona America/Caracas, 18:00–06:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar tickets abiertos y capacidad comprometida para guardia de redes: zona America/Caracas, 18:00–06:00.
3. Ejecuta «Distribuir carga entre equipos: guardia de redes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La carga de guardia de redes se reparte sin exceder límites publicados.

**Alternativas y errores:**
- Métrica incompleta: marcar la estimación como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La carga de guardia de redes se reparte sin exceder límites publicados.
- Se reproduce y controla esta alternativa: Métrica incompleta: marcar la estimación como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de distribuir carga entre equipos: guardia de redes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-074"></a>
## MSP-074 · Distribuir carga entre equipos: equipo nube

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La carga de equipo nube se reparte sin exceder límites publicados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar tickets abiertos y capacidad comprometida para equipo nube: dos técnicos certificados AZ-104.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar tickets abiertos y capacidad comprometida para equipo nube: dos técnicos certificados AZ-104.
3. Ejecuta «Distribuir carga entre equipos: equipo nube»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La carga de equipo nube se reparte sin exceder límites publicados.

**Alternativas y errores:**
- Métrica incompleta: marcar la estimación como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La carga de equipo nube se reparte sin exceder límites publicados.
- Se reproduce y controla esta alternativa: Métrica incompleta: marcar la estimación como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de distribuir carga entre equipos: equipo nube; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-075"></a>
## MSP-075 · Distribuir carga entre equipos: turno fin de semana

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La carga de turno fin de semana se reparte sin exceder límites publicados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar tickets abiertos y capacidad comprometida para turno fin de semana: sábado 08:00–20:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar tickets abiertos y capacidad comprometida para turno fin de semana: sábado 08:00–20:00.
3. Ejecuta «Distribuir carga entre equipos: turno fin de semana»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La carga de turno fin de semana se reparte sin exceder límites publicados.

**Alternativas y errores:**
- Métrica incompleta: marcar la estimación como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La carga de turno fin de semana se reparte sin exceder límites publicados.
- Se reproduce y controla esta alternativa: Métrica incompleta: marcar la estimación como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de distribuir carga entre equipos: turno fin de semana; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-076"></a>
## MSP-076 · Distribuir carga entre equipos: equipo de campo

**Área:** workforce · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La carga de equipo de campo se reparte sin exceder límites publicados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar tickets abiertos y capacidad comprometida para equipo de campo: dos sedes y cuatro visitas pendientes.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar tickets abiertos y capacidad comprometida para equipo de campo: dos sedes y cuatro visitas pendientes.
3. Ejecuta «Distribuir carga entre equipos: equipo de campo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La carga de equipo de campo se reparte sin exceder límites publicados.

**Alternativas y errores:**
- Métrica incompleta: marcar la estimación como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La carga de equipo de campo se reparte sin exceder límites publicados.
- Se reproduce y controla esta alternativa: Métrica incompleta: marcar la estimación como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de distribuir carga entre equipos: equipo de campo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-077"></a>
## MSP-077 · Transferir guardia con contexto: guardia de redes

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El relevo de guardia de redes preserva responsables y decisiones.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Entregar pendientes, riesgos y contactos de emergencia para guardia de redes: zona America/Caracas, 18:00–06:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Entregar pendientes, riesgos y contactos de emergencia para guardia de redes: zona America/Caracas, 18:00–06:00.
3. Ejecuta «Transferir guardia con contexto: guardia de redes»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El relevo de guardia de redes preserva responsables y decisiones.

**Alternativas y errores:**
- Receptor sin permiso de cliente: redactar el contexto restringido.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El relevo de guardia de redes preserva responsables y decisiones.
- Se reproduce y controla esta alternativa: Receptor sin permiso de cliente: redactar el contexto restringido.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de transferir guardia con contexto: guardia de redes; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-078"></a>
## MSP-078 · Transferir guardia con contexto: equipo nube

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El relevo de equipo nube preserva responsables y decisiones.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Entregar pendientes, riesgos y contactos de emergencia para equipo nube: dos técnicos certificados AZ-104.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Entregar pendientes, riesgos y contactos de emergencia para equipo nube: dos técnicos certificados AZ-104.
3. Ejecuta «Transferir guardia con contexto: equipo nube»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El relevo de equipo nube preserva responsables y decisiones.

**Alternativas y errores:**
- Receptor sin permiso de cliente: redactar el contexto restringido.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El relevo de equipo nube preserva responsables y decisiones.
- Se reproduce y controla esta alternativa: Receptor sin permiso de cliente: redactar el contexto restringido.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de transferir guardia con contexto: equipo nube; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-079"></a>
## MSP-079 · Transferir guardia con contexto: turno fin de semana

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El relevo de turno fin de semana preserva responsables y decisiones.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Entregar pendientes, riesgos y contactos de emergencia para turno fin de semana: sábado 08:00–20:00.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Entregar pendientes, riesgos y contactos de emergencia para turno fin de semana: sábado 08:00–20:00.
3. Ejecuta «Transferir guardia con contexto: turno fin de semana»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El relevo de turno fin de semana preserva responsables y decisiones.

**Alternativas y errores:**
- Receptor sin permiso de cliente: redactar el contexto restringido.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El relevo de turno fin de semana preserva responsables y decisiones.
- Se reproduce y controla esta alternativa: Receptor sin permiso de cliente: redactar el contexto restringido.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de transferir guardia con contexto: turno fin de semana; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-080"></a>
## MSP-080 · Transferir guardia con contexto: equipo de campo

**Área:** workforce · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El relevo de equipo de campo preserva responsables y decisiones.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Entregar pendientes, riesgos y contactos de emergencia para equipo de campo: dos sedes y cuatro visitas pendientes.

**Flujo principal:**
1. El actor accede a técnicos, turnos y capacidad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Entregar pendientes, riesgos y contactos de emergencia para equipo de campo: dos sedes y cuatro visitas pendientes.
3. Ejecuta «Transferir guardia con contexto: equipo de campo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El relevo de equipo de campo preserva responsables y decisiones.

**Alternativas y errores:**
- Receptor sin permiso de cliente: redactar el contexto restringido.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El relevo de equipo de campo preserva responsables y decisiones.
- Se reproduce y controla esta alternativa: Receptor sin permiso de cliente: redactar el contexto restringido.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Técnicos, turnos y capacidad: operaciones de transferir guardia con contexto: equipo de campo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-081"></a>
## MSP-081 · Crear proyecto de cliente: migración M365

**Área:** projects-time-expenses · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El proyecto de migración M365 tiene control de alcance y presupuesto.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir alcance, responsable, presupuesto y fechas para migración M365: 120 horas presupuestadas, fase de descubrimiento.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir alcance, responsable, presupuesto y fechas para migración M365: 120 horas presupuestadas, fase de descubrimiento.
3. Ejecuta «Crear proyecto de cliente: migración M365»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El proyecto de migración M365 tiene control de alcance y presupuesto.

**Alternativas y errores:**
- Presupuesto sin aprobador: impedir su publicación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El proyecto de migración M365 tiene control de alcance y presupuesto.
- Se reproduce y controla esta alternativa: Presupuesto sin aprobador: impedir su publicación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de crear proyecto de cliente: migración m365; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-082"></a>
## MSP-082 · Crear proyecto de cliente: renovación de red

**Área:** projects-time-expenses · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El proyecto de renovación de red tiene control de alcance y presupuesto.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir alcance, responsable, presupuesto y fechas para renovación de red: 40 switches en seis sedes.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir alcance, responsable, presupuesto y fechas para renovación de red: 40 switches en seis sedes.
3. Ejecuta «Crear proyecto de cliente: renovación de red»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El proyecto de renovación de red tiene control de alcance y presupuesto.

**Alternativas y errores:**
- Presupuesto sin aprobador: impedir su publicación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El proyecto de renovación de red tiene control de alcance y presupuesto.
- Se reproduce y controla esta alternativa: Presupuesto sin aprobador: impedir su publicación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de crear proyecto de cliente: renovación de red; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-083"></a>
## MSP-083 · Crear proyecto de cliente: puesta en marcha del SOC

**Área:** projects-time-expenses · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El proyecto de puesta en marcha del SOC tiene control de alcance y presupuesto.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir alcance, responsable, presupuesto y fechas para puesta en marcha del SOC: hitos mensuales y guardias incluidas.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir alcance, responsable, presupuesto y fechas para puesta en marcha del SOC: hitos mensuales y guardias incluidas.
3. Ejecuta «Crear proyecto de cliente: puesta en marcha del SOC»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El proyecto de puesta en marcha del SOC tiene control de alcance y presupuesto.

**Alternativas y errores:**
- Presupuesto sin aprobador: impedir su publicación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El proyecto de puesta en marcha del SOC tiene control de alcance y presupuesto.
- Se reproduce y controla esta alternativa: Presupuesto sin aprobador: impedir su publicación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de crear proyecto de cliente: puesta en marcha del soc; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-084"></a>
## MSP-084 · Crear proyecto de cliente: inventario de portátiles

**Área:** projects-time-expenses · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El proyecto de inventario de portátiles tiene control de alcance y presupuesto.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir alcance, responsable, presupuesto y fechas para inventario de portátiles: 75 equipos, Caracas y Valencia.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir alcance, responsable, presupuesto y fechas para inventario de portátiles: 75 equipos, Caracas y Valencia.
3. Ejecuta «Crear proyecto de cliente: inventario de portátiles»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El proyecto de inventario de portátiles tiene control de alcance y presupuesto.

**Alternativas y errores:**
- Presupuesto sin aprobador: impedir su publicación.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El proyecto de inventario de portátiles tiene control de alcance y presupuesto.
- Se reproduce y controla esta alternativa: Presupuesto sin aprobador: impedir su publicación.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de crear proyecto de cliente: inventario de portátiles; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-085"></a>
## MSP-085 · Registrar horas trabajadas: migración M365

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Las horas de migración M365 quedan ligadas al técnico y contrato correctos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar duración, tarea y evidencia del trabajo para migración M365: 120 horas presupuestadas, fase de descubrimiento.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar duración, tarea y evidencia del trabajo para migración M365: 120 horas presupuestadas, fase de descubrimiento.
3. Ejecuta «Registrar horas trabajadas: migración M365»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Las horas de migración M365 quedan ligadas al técnico y contrato correctos.

**Alternativas y errores:**
- Periodo cerrado: enviar a aprobación excepcional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Las horas de migración M365 quedan ligadas al técnico y contrato correctos.
- Se reproduce y controla esta alternativa: Periodo cerrado: enviar a aprobación excepcional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de registrar horas trabajadas: migración m365; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-086"></a>
## MSP-086 · Registrar horas trabajadas: renovación de red

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Las horas de renovación de red quedan ligadas al técnico y contrato correctos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar duración, tarea y evidencia del trabajo para renovación de red: 40 switches en seis sedes.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar duración, tarea y evidencia del trabajo para renovación de red: 40 switches en seis sedes.
3. Ejecuta «Registrar horas trabajadas: renovación de red»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Las horas de renovación de red quedan ligadas al técnico y contrato correctos.

**Alternativas y errores:**
- Periodo cerrado: enviar a aprobación excepcional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Las horas de renovación de red quedan ligadas al técnico y contrato correctos.
- Se reproduce y controla esta alternativa: Periodo cerrado: enviar a aprobación excepcional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de registrar horas trabajadas: renovación de red; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-087"></a>
## MSP-087 · Registrar horas trabajadas: puesta en marcha del SOC

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Las horas de puesta en marcha del SOC quedan ligadas al técnico y contrato correctos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar duración, tarea y evidencia del trabajo para puesta en marcha del SOC: hitos mensuales y guardias incluidas.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar duración, tarea y evidencia del trabajo para puesta en marcha del SOC: hitos mensuales y guardias incluidas.
3. Ejecuta «Registrar horas trabajadas: puesta en marcha del SOC»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Las horas de puesta en marcha del SOC quedan ligadas al técnico y contrato correctos.

**Alternativas y errores:**
- Periodo cerrado: enviar a aprobación excepcional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Las horas de puesta en marcha del SOC quedan ligadas al técnico y contrato correctos.
- Se reproduce y controla esta alternativa: Periodo cerrado: enviar a aprobación excepcional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de registrar horas trabajadas: puesta en marcha del soc; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-088"></a>
## MSP-088 · Registrar horas trabajadas: inventario de portátiles

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Las horas de inventario de portátiles quedan ligadas al técnico y contrato correctos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asociar duración, tarea y evidencia del trabajo para inventario de portátiles: 75 equipos, Caracas y Valencia.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asociar duración, tarea y evidencia del trabajo para inventario de portátiles: 75 equipos, Caracas y Valencia.
3. Ejecuta «Registrar horas trabajadas: inventario de portátiles»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Las horas de inventario de portátiles quedan ligadas al técnico y contrato correctos.

**Alternativas y errores:**
- Periodo cerrado: enviar a aprobación excepcional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Las horas de inventario de portátiles quedan ligadas al técnico y contrato correctos.
- Se reproduce y controla esta alternativa: Periodo cerrado: enviar a aprobación excepcional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de registrar horas trabajadas: inventario de portátiles; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-089"></a>
## MSP-089 · Presentar gasto reembolsable: migración M365

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El gasto de migración M365 entra al circuito de aprobación asignado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar recibo y centro de coste para migración M365: 120 horas presupuestadas, fase de descubrimiento.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar recibo y centro de coste para migración M365: 120 horas presupuestadas, fase de descubrimiento.
3. Ejecuta «Presentar gasto reembolsable: migración M365»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El gasto de migración M365 entra al circuito de aprobación asignado.

**Alternativas y errores:**
- Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El gasto de migración M365 entra al circuito de aprobación asignado.
- Se reproduce y controla esta alternativa: Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de presentar gasto reembolsable: migración m365; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-090"></a>
## MSP-090 · Presentar gasto reembolsable: renovación de red

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El gasto de renovación de red entra al circuito de aprobación asignado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar recibo y centro de coste para renovación de red: 40 switches en seis sedes.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar recibo y centro de coste para renovación de red: 40 switches en seis sedes.
3. Ejecuta «Presentar gasto reembolsable: renovación de red»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El gasto de renovación de red entra al circuito de aprobación asignado.

**Alternativas y errores:**
- Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El gasto de renovación de red entra al circuito de aprobación asignado.
- Se reproduce y controla esta alternativa: Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de presentar gasto reembolsable: renovación de red; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-091"></a>
## MSP-091 · Presentar gasto reembolsable: puesta en marcha del SOC

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El gasto de puesta en marcha del SOC entra al circuito de aprobación asignado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar recibo y centro de coste para puesta en marcha del SOC: hitos mensuales y guardias incluidas.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar recibo y centro de coste para puesta en marcha del SOC: hitos mensuales y guardias incluidas.
3. Ejecuta «Presentar gasto reembolsable: puesta en marcha del SOC»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El gasto de puesta en marcha del SOC entra al circuito de aprobación asignado.

**Alternativas y errores:**
- Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El gasto de puesta en marcha del SOC entra al circuito de aprobación asignado.
- Se reproduce y controla esta alternativa: Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de presentar gasto reembolsable: puesta en marcha del soc; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-092"></a>
## MSP-092 · Presentar gasto reembolsable: inventario de portátiles

**Área:** projects-time-expenses · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El gasto de inventario de portátiles entra al circuito de aprobación asignado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar recibo y centro de coste para inventario de portátiles: 75 equipos, Caracas y Valencia.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar recibo y centro de coste para inventario de portátiles: 75 equipos, Caracas y Valencia.
3. Ejecuta «Presentar gasto reembolsable: inventario de portátiles»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El gasto de inventario de portátiles entra al circuito de aprobación asignado.

**Alternativas y errores:**
- Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El gasto de inventario de portátiles entra al circuito de aprobación asignado.
- Se reproduce y controla esta alternativa: Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de presentar gasto reembolsable: inventario de portátiles; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-093"></a>
## MSP-093 · Aprobar partes de tiempo: migración M365

**Área:** projects-time-expenses · **Actor:** Aprobador designado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La aprobación de migración M365 conserva autor y comentario.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Revisar detalle y comparar con el alcance contratado para migración M365: 120 horas presupuestadas, fase de descubrimiento.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisar detalle y comparar con el alcance contratado para migración M365: 120 horas presupuestadas, fase de descubrimiento.
3. Ejecuta «Aprobar partes de tiempo: migración M365»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La aprobación de migración M365 conserva autor y comentario.

**Alternativas y errores:**
- Aprobador es quien registró las horas: pedir aprobación independiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La aprobación de migración M365 conserva autor y comentario.
- Se reproduce y controla esta alternativa: Aprobador es quien registró las horas: pedir aprobación independiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de aprobar partes de tiempo: migración m365; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-094"></a>
## MSP-094 · Aprobar partes de tiempo: renovación de red

**Área:** projects-time-expenses · **Actor:** Aprobador designado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La aprobación de renovación de red conserva autor y comentario.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Revisar detalle y comparar con el alcance contratado para renovación de red: 40 switches en seis sedes.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisar detalle y comparar con el alcance contratado para renovación de red: 40 switches en seis sedes.
3. Ejecuta «Aprobar partes de tiempo: renovación de red»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La aprobación de renovación de red conserva autor y comentario.

**Alternativas y errores:**
- Aprobador es quien registró las horas: pedir aprobación independiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La aprobación de renovación de red conserva autor y comentario.
- Se reproduce y controla esta alternativa: Aprobador es quien registró las horas: pedir aprobación independiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de aprobar partes de tiempo: renovación de red; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-095"></a>
## MSP-095 · Aprobar partes de tiempo: puesta en marcha del SOC

**Área:** projects-time-expenses · **Actor:** Aprobador designado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La aprobación de puesta en marcha del SOC conserva autor y comentario.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Revisar detalle y comparar con el alcance contratado para puesta en marcha del SOC: hitos mensuales y guardias incluidas.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisar detalle y comparar con el alcance contratado para puesta en marcha del SOC: hitos mensuales y guardias incluidas.
3. Ejecuta «Aprobar partes de tiempo: puesta en marcha del SOC»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La aprobación de puesta en marcha del SOC conserva autor y comentario.

**Alternativas y errores:**
- Aprobador es quien registró las horas: pedir aprobación independiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La aprobación de puesta en marcha del SOC conserva autor y comentario.
- Se reproduce y controla esta alternativa: Aprobador es quien registró las horas: pedir aprobación independiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de aprobar partes de tiempo: puesta en marcha del soc; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-096"></a>
## MSP-096 · Aprobar partes de tiempo: inventario de portátiles

**Área:** projects-time-expenses · **Actor:** Aprobador designado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La aprobación de inventario de portátiles conserva autor y comentario.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Revisar detalle y comparar con el alcance contratado para inventario de portátiles: 75 equipos, Caracas y Valencia.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Revisar detalle y comparar con el alcance contratado para inventario de portátiles: 75 equipos, Caracas y Valencia.
3. Ejecuta «Aprobar partes de tiempo: inventario de portátiles»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La aprobación de inventario de portátiles conserva autor y comentario.

**Alternativas y errores:**
- Aprobador es quien registró las horas: pedir aprobación independiente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La aprobación de inventario de portátiles conserva autor y comentario.
- Se reproduce y controla esta alternativa: Aprobador es quien registró las horas: pedir aprobación independiente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de aprobar partes de tiempo: inventario de portátiles; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-097"></a>
## MSP-097 · Revisar desvío de presupuesto: migración M365

**Área:** projects-time-expenses · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El desvío de migración M365 se informa antes de exceder el límite.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas, gastos y avance con el límite acordado para migración M365: 120 horas presupuestadas, fase de descubrimiento.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas, gastos y avance con el límite acordado para migración M365: 120 horas presupuestadas, fase de descubrimiento.
3. Ejecuta «Revisar desvío de presupuesto: migración M365»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El desvío de migración M365 se informa antes de exceder el límite.

**Alternativas y errores:**
- Datos importados sin conciliar: mostrar advertencia y no facturar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El desvío de migración M365 se informa antes de exceder el límite.
- Se reproduce y controla esta alternativa: Datos importados sin conciliar: mostrar advertencia y no facturar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de revisar desvío de presupuesto: migración m365; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-098"></a>
## MSP-098 · Revisar desvío de presupuesto: renovación de red

**Área:** projects-time-expenses · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El desvío de renovación de red se informa antes de exceder el límite.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas, gastos y avance con el límite acordado para renovación de red: 40 switches en seis sedes.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas, gastos y avance con el límite acordado para renovación de red: 40 switches en seis sedes.
3. Ejecuta «Revisar desvío de presupuesto: renovación de red»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El desvío de renovación de red se informa antes de exceder el límite.

**Alternativas y errores:**
- Datos importados sin conciliar: mostrar advertencia y no facturar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El desvío de renovación de red se informa antes de exceder el límite.
- Se reproduce y controla esta alternativa: Datos importados sin conciliar: mostrar advertencia y no facturar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de revisar desvío de presupuesto: renovación de red; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-099"></a>
## MSP-099 · Revisar desvío de presupuesto: puesta en marcha del SOC

**Área:** projects-time-expenses · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El desvío de puesta en marcha del SOC se informa antes de exceder el límite.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas, gastos y avance con el límite acordado para puesta en marcha del SOC: hitos mensuales y guardias incluidas.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas, gastos y avance con el límite acordado para puesta en marcha del SOC: hitos mensuales y guardias incluidas.
3. Ejecuta «Revisar desvío de presupuesto: puesta en marcha del SOC»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El desvío de puesta en marcha del SOC se informa antes de exceder el límite.

**Alternativas y errores:**
- Datos importados sin conciliar: mostrar advertencia y no facturar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El desvío de puesta en marcha del SOC se informa antes de exceder el límite.
- Se reproduce y controla esta alternativa: Datos importados sin conciliar: mostrar advertencia y no facturar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de revisar desvío de presupuesto: puesta en marcha del soc; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-100"></a>
## MSP-100 · Revisar desvío de presupuesto: inventario de portátiles

**Área:** projects-time-expenses · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El desvío de inventario de portátiles se informa antes de exceder el límite.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas, gastos y avance con el límite acordado para inventario de portátiles: 75 equipos, Caracas y Valencia.

**Flujo principal:**
1. El actor accede a proyectos, horas y gastos dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas, gastos y avance con el límite acordado para inventario de portátiles: 75 equipos, Caracas y Valencia.
3. Ejecuta «Revisar desvío de presupuesto: inventario de portátiles»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El desvío de inventario de portátiles se informa antes de exceder el límite.

**Alternativas y errores:**
- Datos importados sin conciliar: mostrar advertencia y no facturar.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El desvío de inventario de portátiles se informa antes de exceder el límite.
- Se reproduce y controla esta alternativa: Datos importados sin conciliar: mostrar advertencia y no facturar.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Proyectos, horas y gastos: operaciones de revisar desvío de presupuesto: inventario de portátiles; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-101"></a>
## MSP-101 · Dar de alta un activo: laptop TR-104

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El activo de laptop TR-104 se identifica de forma única y trazable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Capturar serial, modelo, cliente y ubicación para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Capturar serial, modelo, cliente y ubicación para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.
3. Ejecuta «Dar de alta un activo: laptop TR-104»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El activo de laptop TR-104 se identifica de forma única y trazable.

**Alternativas y errores:**
- Número de serie ya registrado: evitar inventario duplicado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El activo de laptop TR-104 se identifica de forma única y trazable.
- Se reproduce y controla esta alternativa: Número de serie ya registrado: evitar inventario duplicado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de dar de alta un activo: laptop tr-104; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-102"></a>
## MSP-102 · Dar de alta un activo: switch SW-22

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El activo de switch SW-22 se identifica de forma única y trazable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Capturar serial, modelo, cliente y ubicación para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Capturar serial, modelo, cliente y ubicación para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.
3. Ejecuta «Dar de alta un activo: switch SW-22»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El activo de switch SW-22 se identifica de forma única y trazable.

**Alternativas y errores:**
- Número de serie ya registrado: evitar inventario duplicado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El activo de switch SW-22 se identifica de forma única y trazable.
- Se reproduce y controla esta alternativa: Número de serie ya registrado: evitar inventario duplicado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de dar de alta un activo: switch sw-22; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-103"></a>
## MSP-103 · Dar de alta un activo: servidor SR-08

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El activo de servidor SR-08 se identifica de forma única y trazable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Capturar serial, modelo, cliente y ubicación para servidor SR-08: Linux, Colegio Sol, sala técnica.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Capturar serial, modelo, cliente y ubicación para servidor SR-08: Linux, Colegio Sol, sala técnica.
3. Ejecuta «Dar de alta un activo: servidor SR-08»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El activo de servidor SR-08 se identifica de forma única y trazable.

**Alternativas y errores:**
- Número de serie ya registrado: evitar inventario duplicado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El activo de servidor SR-08 se identifica de forma única y trazable.
- Se reproduce y controla esta alternativa: Número de serie ya registrado: evitar inventario duplicado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de dar de alta un activo: servidor sr-08; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-104"></a>
## MSP-104 · Dar de alta un activo: tablet TB-31

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El activo de tablet TB-31 se identifica de forma única y trazable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Capturar serial, modelo, cliente y ubicación para tablet TB-31: Zebra, Logística Norte, almacén 4.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Capturar serial, modelo, cliente y ubicación para tablet TB-31: Zebra, Logística Norte, almacén 4.
3. Ejecuta «Dar de alta un activo: tablet TB-31»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El activo de tablet TB-31 se identifica de forma única y trazable.

**Alternativas y errores:**
- Número de serie ya registrado: evitar inventario duplicado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El activo de tablet TB-31 se identifica de forma única y trazable.
- Se reproduce y controla esta alternativa: Número de serie ya registrado: evitar inventario duplicado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de dar de alta un activo: tablet tb-31; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-105"></a>
## MSP-105 · Vincular activo a un usuario: laptop TR-104

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La custodia de laptop TR-104 refleja la persona y sede autorizadas.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar custodio y fecha de entrega para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar custodio y fecha de entrega para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.
3. Ejecuta «Vincular activo a un usuario: laptop TR-104»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La custodia de laptop TR-104 refleja la persona y sede autorizadas.

**Alternativas y errores:**
- Usuario de otro cliente: bloquear el vínculo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La custodia de laptop TR-104 refleja la persona y sede autorizadas.
- Se reproduce y controla esta alternativa: Usuario de otro cliente: bloquear el vínculo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de vincular activo a un usuario: laptop tr-104; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-106"></a>
## MSP-106 · Vincular activo a un usuario: switch SW-22

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La custodia de switch SW-22 refleja la persona y sede autorizadas.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar custodio y fecha de entrega para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar custodio y fecha de entrega para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.
3. Ejecuta «Vincular activo a un usuario: switch SW-22»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La custodia de switch SW-22 refleja la persona y sede autorizadas.

**Alternativas y errores:**
- Usuario de otro cliente: bloquear el vínculo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La custodia de switch SW-22 refleja la persona y sede autorizadas.
- Se reproduce y controla esta alternativa: Usuario de otro cliente: bloquear el vínculo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de vincular activo a un usuario: switch sw-22; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-107"></a>
## MSP-107 · Vincular activo a un usuario: servidor SR-08

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La custodia de servidor SR-08 refleja la persona y sede autorizadas.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar custodio y fecha de entrega para servidor SR-08: Linux, Colegio Sol, sala técnica.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar custodio y fecha de entrega para servidor SR-08: Linux, Colegio Sol, sala técnica.
3. Ejecuta «Vincular activo a un usuario: servidor SR-08»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La custodia de servidor SR-08 refleja la persona y sede autorizadas.

**Alternativas y errores:**
- Usuario de otro cliente: bloquear el vínculo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La custodia de servidor SR-08 refleja la persona y sede autorizadas.
- Se reproduce y controla esta alternativa: Usuario de otro cliente: bloquear el vínculo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de vincular activo a un usuario: servidor sr-08; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-108"></a>
## MSP-108 · Vincular activo a un usuario: tablet TB-31

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La custodia de tablet TB-31 refleja la persona y sede autorizadas.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar custodio y fecha de entrega para tablet TB-31: Zebra, Logística Norte, almacén 4.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar custodio y fecha de entrega para tablet TB-31: Zebra, Logística Norte, almacén 4.
3. Ejecuta «Vincular activo a un usuario: tablet TB-31»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La custodia de tablet TB-31 refleja la persona y sede autorizadas.

**Alternativas y errores:**
- Usuario de otro cliente: bloquear el vínculo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La custodia de tablet TB-31 refleja la persona y sede autorizadas.
- Se reproduce y controla esta alternativa: Usuario de otro cliente: bloquear el vínculo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de vincular activo a un usuario: tablet tb-31; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-109"></a>
## MSP-109 · Actualizar garantía: laptop TR-104

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La garantía de laptop TR-104 se consulta con alertas de vencimiento.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar proveedor, cobertura y vencimiento para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar proveedor, cobertura y vencimiento para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.
3. Ejecuta «Actualizar garantía: laptop TR-104»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La garantía de laptop TR-104 se consulta con alertas de vencimiento.

**Alternativas y errores:**
- Documento no verificable: dejar la actualización en revisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La garantía de laptop TR-104 se consulta con alertas de vencimiento.
- Se reproduce y controla esta alternativa: Documento no verificable: dejar la actualización en revisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de actualizar garantía: laptop tr-104; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-110"></a>
## MSP-110 · Actualizar garantía: switch SW-22

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La garantía de switch SW-22 se consulta con alertas de vencimiento.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar proveedor, cobertura y vencimiento para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar proveedor, cobertura y vencimiento para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.
3. Ejecuta «Actualizar garantía: switch SW-22»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La garantía de switch SW-22 se consulta con alertas de vencimiento.

**Alternativas y errores:**
- Documento no verificable: dejar la actualización en revisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La garantía de switch SW-22 se consulta con alertas de vencimiento.
- Se reproduce y controla esta alternativa: Documento no verificable: dejar la actualización en revisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de actualizar garantía: switch sw-22; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-111"></a>
## MSP-111 · Actualizar garantía: servidor SR-08

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La garantía de servidor SR-08 se consulta con alertas de vencimiento.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar proveedor, cobertura y vencimiento para servidor SR-08: Linux, Colegio Sol, sala técnica.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar proveedor, cobertura y vencimiento para servidor SR-08: Linux, Colegio Sol, sala técnica.
3. Ejecuta «Actualizar garantía: servidor SR-08»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La garantía de servidor SR-08 se consulta con alertas de vencimiento.

**Alternativas y errores:**
- Documento no verificable: dejar la actualización en revisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La garantía de servidor SR-08 se consulta con alertas de vencimiento.
- Se reproduce y controla esta alternativa: Documento no verificable: dejar la actualización en revisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de actualizar garantía: servidor sr-08; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-112"></a>
## MSP-112 · Actualizar garantía: tablet TB-31

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La garantía de tablet TB-31 se consulta con alertas de vencimiento.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Adjuntar proveedor, cobertura y vencimiento para tablet TB-31: Zebra, Logística Norte, almacén 4.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Adjuntar proveedor, cobertura y vencimiento para tablet TB-31: Zebra, Logística Norte, almacén 4.
3. Ejecuta «Actualizar garantía: tablet TB-31»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La garantía de tablet TB-31 se consulta con alertas de vencimiento.

**Alternativas y errores:**
- Documento no verificable: dejar la actualización en revisión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La garantía de tablet TB-31 se consulta con alertas de vencimiento.
- Se reproduce y controla esta alternativa: Documento no verificable: dejar la actualización en revisión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de actualizar garantía: tablet tb-31; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-113"></a>
## MSP-113 · Retirar activo con borrado seguro: laptop TR-104

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La baja de laptop TR-104 conserva auditoría y certificado disponible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar aprobación, respaldo y disposición final para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar aprobación, respaldo y disposición final para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.
3. Ejecuta «Retirar activo con borrado seguro: laptop TR-104»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La baja de laptop TR-104 conserva auditoría y certificado disponible.

**Alternativas y errores:**
- Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La baja de laptop TR-104 conserva auditoría y certificado disponible.
- Se reproduce y controla esta alternativa: Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de retirar activo con borrado seguro: laptop tr-104; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-114"></a>
## MSP-114 · Retirar activo con borrado seguro: switch SW-22

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La baja de switch SW-22 conserva auditoría y certificado disponible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar aprobación, respaldo y disposición final para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar aprobación, respaldo y disposición final para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.
3. Ejecuta «Retirar activo con borrado seguro: switch SW-22»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La baja de switch SW-22 conserva auditoría y certificado disponible.

**Alternativas y errores:**
- Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La baja de switch SW-22 conserva auditoría y certificado disponible.
- Se reproduce y controla esta alternativa: Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de retirar activo con borrado seguro: switch sw-22; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-115"></a>
## MSP-115 · Retirar activo con borrado seguro: servidor SR-08

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La baja de servidor SR-08 conserva auditoría y certificado disponible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar aprobación, respaldo y disposición final para servidor SR-08: Linux, Colegio Sol, sala técnica.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar aprobación, respaldo y disposición final para servidor SR-08: Linux, Colegio Sol, sala técnica.
3. Ejecuta «Retirar activo con borrado seguro: servidor SR-08»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La baja de servidor SR-08 conserva auditoría y certificado disponible.

**Alternativas y errores:**
- Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La baja de servidor SR-08 conserva auditoría y certificado disponible.
- Se reproduce y controla esta alternativa: Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de retirar activo con borrado seguro: servidor sr-08; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-116"></a>
## MSP-116 · Retirar activo con borrado seguro: tablet TB-31

**Área:** assets · **Actor:** Gestor de activos · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La baja de tablet TB-31 conserva auditoría y certificado disponible.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Registrar aprobación, respaldo y disposición final para tablet TB-31: Zebra, Logística Norte, almacén 4.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Registrar aprobación, respaldo y disposición final para tablet TB-31: Zebra, Logística Norte, almacén 4.
3. Ejecuta «Retirar activo con borrado seguro: tablet TB-31»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La baja de tablet TB-31 conserva auditoría y certificado disponible.

**Alternativas y errores:**
- Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La baja de tablet TB-31 conserva auditoría y certificado disponible.
- Se reproduce y controla esta alternativa: Activo con trabajo pendiente: detener la baja hasta transferir tareas.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de retirar activo con borrado seguro: tablet tb-31; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-117"></a>
## MSP-117 · Conciliar inventario físico: laptop TR-104

**Área:** assets · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conciliación de laptop TR-104 muestra diferencias para revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar escaneo con inventario esperado de la sede para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar escaneo con inventario esperado de la sede para laptop TR-104: Dell Latitude, Clínica Aurora, Caracas.
3. Ejecuta «Conciliar inventario físico: laptop TR-104»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conciliación de laptop TR-104 muestra diferencias para revisión.

**Alternativas y errores:**
- Escaneo parcial: no marcar como verificada la ubicación completa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conciliación de laptop TR-104 muestra diferencias para revisión.
- Se reproduce y controla esta alternativa: Escaneo parcial: no marcar como verificada la ubicación completa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de conciliar inventario físico: laptop tr-104; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-118"></a>
## MSP-118 · Conciliar inventario físico: switch SW-22

**Área:** assets · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conciliación de switch SW-22 muestra diferencias para revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar escaneo con inventario esperado de la sede para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar escaneo con inventario esperado de la sede para switch SW-22: Cisco Catalyst, Hotel Río, MDF 2.
3. Ejecuta «Conciliar inventario físico: switch SW-22»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conciliación de switch SW-22 muestra diferencias para revisión.

**Alternativas y errores:**
- Escaneo parcial: no marcar como verificada la ubicación completa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conciliación de switch SW-22 muestra diferencias para revisión.
- Se reproduce y controla esta alternativa: Escaneo parcial: no marcar como verificada la ubicación completa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de conciliar inventario físico: switch sw-22; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-119"></a>
## MSP-119 · Conciliar inventario físico: servidor SR-08

**Área:** assets · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conciliación de servidor SR-08 muestra diferencias para revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar escaneo con inventario esperado de la sede para servidor SR-08: Linux, Colegio Sol, sala técnica.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar escaneo con inventario esperado de la sede para servidor SR-08: Linux, Colegio Sol, sala técnica.
3. Ejecuta «Conciliar inventario físico: servidor SR-08»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conciliación de servidor SR-08 muestra diferencias para revisión.

**Alternativas y errores:**
- Escaneo parcial: no marcar como verificada la ubicación completa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conciliación de servidor SR-08 muestra diferencias para revisión.
- Se reproduce y controla esta alternativa: Escaneo parcial: no marcar como verificada la ubicación completa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de conciliar inventario físico: servidor sr-08; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-120"></a>
## MSP-120 · Conciliar inventario físico: tablet TB-31

**Área:** assets · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conciliación de tablet TB-31 muestra diferencias para revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar escaneo con inventario esperado de la sede para tablet TB-31: Zebra, Logística Norte, almacén 4.

**Flujo principal:**
1. El actor accede a activos, garantías e inventario dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar escaneo con inventario esperado de la sede para tablet TB-31: Zebra, Logística Norte, almacén 4.
3. Ejecuta «Conciliar inventario físico: tablet TB-31»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conciliación de tablet TB-31 muestra diferencias para revisión.

**Alternativas y errores:**
- Escaneo parcial: no marcar como verificada la ubicación completa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conciliación de tablet TB-31 muestra diferencias para revisión.
- Se reproduce y controla esta alternativa: Escaneo parcial: no marcar como verificada la ubicación completa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Activos, garantías e inventario: operaciones de conciliar inventario físico: tablet tb-31; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-121"></a>
## MSP-121 · Importar descubrimiento aprobado: uso de CPU

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Los hallazgos de uso de CPU se crean dentro del ámbito autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conciliar dispositivos detectados con el cliente y segmento para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conciliar dispositivos detectados con el cliente y segmento para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.
3. Ejecuta «Importar descubrimiento aprobado: uso de CPU»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Los hallazgos de uso de CPU se crean dentro del ámbito autorizado.

**Alternativas y errores:**
- Segmento fuera de inventario: poner hallazgos en cuarentena.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Los hallazgos de uso de CPU se crean dentro del ámbito autorizado.
- Se reproduce y controla esta alternativa: Segmento fuera de inventario: poner hallazgos en cuarentena.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de importar descubrimiento aprobado: uso de cpu; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-122"></a>
## MSP-122 · Importar descubrimiento aprobado: pérdida de enlace

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Los hallazgos de pérdida de enlace se crean dentro del ámbito autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conciliar dispositivos detectados con el cliente y segmento para pérdida de enlace: Zabbix, router hotel-rio-wan1.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conciliar dispositivos detectados con el cliente y segmento para pérdida de enlace: Zabbix, router hotel-rio-wan1.
3. Ejecuta «Importar descubrimiento aprobado: pérdida de enlace»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Los hallazgos de pérdida de enlace se crean dentro del ámbito autorizado.

**Alternativas y errores:**
- Segmento fuera de inventario: poner hallazgos en cuarentena.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Los hallazgos de pérdida de enlace se crean dentro del ámbito autorizado.
- Se reproduce y controla esta alternativa: Segmento fuera de inventario: poner hallazgos en cuarentena.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de importar descubrimiento aprobado: pérdida de enlace; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-123"></a>
## MSP-123 · Importar descubrimiento aprobado: disco casi lleno

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Los hallazgos de disco casi lleno se crean dentro del ámbito autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conciliar dispositivos detectados con el cliente y segmento para disco casi lleno: Zabbix, servidor colegio-sol-db.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conciliar dispositivos detectados con el cliente y segmento para disco casi lleno: Zabbix, servidor colegio-sol-db.
3. Ejecuta «Importar descubrimiento aprobado: disco casi lleno»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Los hallazgos de disco casi lleno se crean dentro del ámbito autorizado.

**Alternativas y errores:**
- Segmento fuera de inventario: poner hallazgos en cuarentena.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Los hallazgos de disco casi lleno se crean dentro del ámbito autorizado.
- Se reproduce y controla esta alternativa: Segmento fuera de inventario: poner hallazgos en cuarentena.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de importar descubrimiento aprobado: disco casi lleno; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-124"></a>
## MSP-124 · Importar descubrimiento aprobado: servicio de respaldo

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** Los hallazgos de servicio de respaldo se crean dentro del ámbito autorizado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conciliar dispositivos detectados con el cliente y segmento para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conciliar dispositivos detectados con el cliente y segmento para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.
3. Ejecuta «Importar descubrimiento aprobado: servicio de respaldo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: Los hallazgos de servicio de respaldo se crean dentro del ámbito autorizado.

**Alternativas y errores:**
- Segmento fuera de inventario: poner hallazgos en cuarentena.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: Los hallazgos de servicio de respaldo se crean dentro del ámbito autorizado.
- Se reproduce y controla esta alternativa: Segmento fuera de inventario: poner hallazgos en cuarentena.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de importar descubrimiento aprobado: servicio de respaldo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-125"></a>
## MSP-125 · Convertir alerta en incidente: uso de CPU

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La alerta de uso de CPU genera un solo incidente asignable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Aplicar umbral, ventana y regla de correlación para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Aplicar umbral, ventana y regla de correlación para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.
3. Ejecuta «Convertir alerta en incidente: uso de CPU»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La alerta de uso de CPU genera un solo incidente asignable.

**Alternativas y errores:**
- Alerta duplicada del mismo evento: actualizar el incidente existente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La alerta de uso de CPU genera un solo incidente asignable.
- Se reproduce y controla esta alternativa: Alerta duplicada del mismo evento: actualizar el incidente existente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de convertir alerta en incidente: uso de cpu; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-126"></a>
## MSP-126 · Convertir alerta en incidente: pérdida de enlace

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La alerta de pérdida de enlace genera un solo incidente asignable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Aplicar umbral, ventana y regla de correlación para pérdida de enlace: Zabbix, router hotel-rio-wan1.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Aplicar umbral, ventana y regla de correlación para pérdida de enlace: Zabbix, router hotel-rio-wan1.
3. Ejecuta «Convertir alerta en incidente: pérdida de enlace»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La alerta de pérdida de enlace genera un solo incidente asignable.

**Alternativas y errores:**
- Alerta duplicada del mismo evento: actualizar el incidente existente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La alerta de pérdida de enlace genera un solo incidente asignable.
- Se reproduce y controla esta alternativa: Alerta duplicada del mismo evento: actualizar el incidente existente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de convertir alerta en incidente: pérdida de enlace; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-127"></a>
## MSP-127 · Convertir alerta en incidente: disco casi lleno

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La alerta de disco casi lleno genera un solo incidente asignable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Aplicar umbral, ventana y regla de correlación para disco casi lleno: Zabbix, servidor colegio-sol-db.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Aplicar umbral, ventana y regla de correlación para disco casi lleno: Zabbix, servidor colegio-sol-db.
3. Ejecuta «Convertir alerta en incidente: disco casi lleno»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La alerta de disco casi lleno genera un solo incidente asignable.

**Alternativas y errores:**
- Alerta duplicada del mismo evento: actualizar el incidente existente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La alerta de disco casi lleno genera un solo incidente asignable.
- Se reproduce y controla esta alternativa: Alerta duplicada del mismo evento: actualizar el incidente existente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de convertir alerta en incidente: disco casi lleno; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-128"></a>
## MSP-128 · Convertir alerta en incidente: servicio de respaldo

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La alerta de servicio de respaldo genera un solo incidente asignable.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Aplicar umbral, ventana y regla de correlación para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Aplicar umbral, ventana y regla de correlación para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.
3. Ejecuta «Convertir alerta en incidente: servicio de respaldo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La alerta de servicio de respaldo genera un solo incidente asignable.

**Alternativas y errores:**
- Alerta duplicada del mismo evento: actualizar el incidente existente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La alerta de servicio de respaldo genera un solo incidente asignable.
- Se reproduce y controla esta alternativa: Alerta duplicada del mismo evento: actualizar el incidente existente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de convertir alerta en incidente: servicio de respaldo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-129"></a>
## MSP-129 · Silenciar mantenimiento autorizado: uso de CPU

**Área:** monitoring · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El silencio de uso de CPU vence automáticamente al terminar la ventana.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir servicio, motivo y ventana temporal para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir servicio, motivo y ventana temporal para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.
3. Ejecuta «Silenciar mantenimiento autorizado: uso de CPU»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El silencio de uso de CPU vence automáticamente al terminar la ventana.

**Alternativas y errores:**
- Ventana sin aprobador: mantener la alerta activa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El silencio de uso de CPU vence automáticamente al terminar la ventana.
- Se reproduce y controla esta alternativa: Ventana sin aprobador: mantener la alerta activa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de silenciar mantenimiento autorizado: uso de cpu; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-130"></a>
## MSP-130 · Silenciar mantenimiento autorizado: pérdida de enlace

**Área:** monitoring · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El silencio de pérdida de enlace vence automáticamente al terminar la ventana.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir servicio, motivo y ventana temporal para pérdida de enlace: Zabbix, router hotel-rio-wan1.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir servicio, motivo y ventana temporal para pérdida de enlace: Zabbix, router hotel-rio-wan1.
3. Ejecuta «Silenciar mantenimiento autorizado: pérdida de enlace»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El silencio de pérdida de enlace vence automáticamente al terminar la ventana.

**Alternativas y errores:**
- Ventana sin aprobador: mantener la alerta activa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El silencio de pérdida de enlace vence automáticamente al terminar la ventana.
- Se reproduce y controla esta alternativa: Ventana sin aprobador: mantener la alerta activa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de silenciar mantenimiento autorizado: pérdida de enlace; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-131"></a>
## MSP-131 · Silenciar mantenimiento autorizado: disco casi lleno

**Área:** monitoring · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El silencio de disco casi lleno vence automáticamente al terminar la ventana.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir servicio, motivo y ventana temporal para disco casi lleno: Zabbix, servidor colegio-sol-db.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir servicio, motivo y ventana temporal para disco casi lleno: Zabbix, servidor colegio-sol-db.
3. Ejecuta «Silenciar mantenimiento autorizado: disco casi lleno»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El silencio de disco casi lleno vence automáticamente al terminar la ventana.

**Alternativas y errores:**
- Ventana sin aprobador: mantener la alerta activa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El silencio de disco casi lleno vence automáticamente al terminar la ventana.
- Se reproduce y controla esta alternativa: Ventana sin aprobador: mantener la alerta activa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de silenciar mantenimiento autorizado: disco casi lleno; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-132"></a>
## MSP-132 · Silenciar mantenimiento autorizado: servicio de respaldo

**Área:** monitoring · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El silencio de servicio de respaldo vence automáticamente al terminar la ventana.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Definir servicio, motivo y ventana temporal para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Definir servicio, motivo y ventana temporal para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.
3. Ejecuta «Silenciar mantenimiento autorizado: servicio de respaldo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El silencio de servicio de respaldo vence automáticamente al terminar la ventana.

**Alternativas y errores:**
- Ventana sin aprobador: mantener la alerta activa.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El silencio de servicio de respaldo vence automáticamente al terminar la ventana.
- Se reproduce y controla esta alternativa: Ventana sin aprobador: mantener la alerta activa.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de silenciar mantenimiento autorizado: servicio de respaldo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-133"></a>
## MSP-133 · Correlacionar alertas relacionadas: uso de CPU

**Área:** monitoring · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La correlación de uso de CPU reduce duplicados sin perder señales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar eventos por servicio y causa probable para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar eventos por servicio y causa probable para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.
3. Ejecuta «Correlacionar alertas relacionadas: uso de CPU»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La correlación de uso de CPU reduce duplicados sin perder señales.

**Alternativas y errores:**
- Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La correlación de uso de CPU reduce duplicados sin perder señales.
- Se reproduce y controla esta alternativa: Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de correlacionar alertas relacionadas: uso de cpu; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-134"></a>
## MSP-134 · Correlacionar alertas relacionadas: pérdida de enlace

**Área:** monitoring · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La correlación de pérdida de enlace reduce duplicados sin perder señales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar eventos por servicio y causa probable para pérdida de enlace: Zabbix, router hotel-rio-wan1.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar eventos por servicio y causa probable para pérdida de enlace: Zabbix, router hotel-rio-wan1.
3. Ejecuta «Correlacionar alertas relacionadas: pérdida de enlace»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La correlación de pérdida de enlace reduce duplicados sin perder señales.

**Alternativas y errores:**
- Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La correlación de pérdida de enlace reduce duplicados sin perder señales.
- Se reproduce y controla esta alternativa: Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de correlacionar alertas relacionadas: pérdida de enlace; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-135"></a>
## MSP-135 · Correlacionar alertas relacionadas: disco casi lleno

**Área:** monitoring · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La correlación de disco casi lleno reduce duplicados sin perder señales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar eventos por servicio y causa probable para disco casi lleno: Zabbix, servidor colegio-sol-db.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar eventos por servicio y causa probable para disco casi lleno: Zabbix, servidor colegio-sol-db.
3. Ejecuta «Correlacionar alertas relacionadas: disco casi lleno»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La correlación de disco casi lleno reduce duplicados sin perder señales.

**Alternativas y errores:**
- Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La correlación de disco casi lleno reduce duplicados sin perder señales.
- Se reproduce y controla esta alternativa: Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de correlacionar alertas relacionadas: disco casi lleno; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-136"></a>
## MSP-136 · Correlacionar alertas relacionadas: servicio de respaldo

**Área:** monitoring · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La correlación de servicio de respaldo reduce duplicados sin perder señales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar eventos por servicio y causa probable para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar eventos por servicio y causa probable para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.
3. Ejecuta «Correlacionar alertas relacionadas: servicio de respaldo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La correlación de servicio de respaldo reduce duplicados sin perder señales.

**Alternativas y errores:**
- Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La correlación de servicio de respaldo reduce duplicados sin perder señales.
- Se reproduce y controla esta alternativa: Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de correlacionar alertas relacionadas: servicio de respaldo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-137"></a>
## MSP-137 · Cerrar alerta tras recuperación: uso de CPU

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La señal de uso de CPU cierra o actualiza el incidente asociado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comprobar recuperación y vincular evidencia de monitorización para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comprobar recuperación y vincular evidencia de monitorización para uso de CPU: Zabbix, host srv-aurora-02, umbral 90%.
3. Ejecuta «Cerrar alerta tras recuperación: uso de CPU»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La señal de uso de CPU cierra o actualiza el incidente asociado.

**Alternativas y errores:**
- Sin lectura reciente: no inferir que el servicio se recuperó.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La señal de uso de CPU cierra o actualiza el incidente asociado.
- Se reproduce y controla esta alternativa: Sin lectura reciente: no inferir que el servicio se recuperó.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de cerrar alerta tras recuperación: uso de cpu; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-138"></a>
## MSP-138 · Cerrar alerta tras recuperación: pérdida de enlace

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La señal de pérdida de enlace cierra o actualiza el incidente asociado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comprobar recuperación y vincular evidencia de monitorización para pérdida de enlace: Zabbix, router hotel-rio-wan1.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comprobar recuperación y vincular evidencia de monitorización para pérdida de enlace: Zabbix, router hotel-rio-wan1.
3. Ejecuta «Cerrar alerta tras recuperación: pérdida de enlace»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La señal de pérdida de enlace cierra o actualiza el incidente asociado.

**Alternativas y errores:**
- Sin lectura reciente: no inferir que el servicio se recuperó.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La señal de pérdida de enlace cierra o actualiza el incidente asociado.
- Se reproduce y controla esta alternativa: Sin lectura reciente: no inferir que el servicio se recuperó.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de cerrar alerta tras recuperación: pérdida de enlace; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-139"></a>
## MSP-139 · Cerrar alerta tras recuperación: disco casi lleno

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La señal de disco casi lleno cierra o actualiza el incidente asociado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comprobar recuperación y vincular evidencia de monitorización para disco casi lleno: Zabbix, servidor colegio-sol-db.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comprobar recuperación y vincular evidencia de monitorización para disco casi lleno: Zabbix, servidor colegio-sol-db.
3. Ejecuta «Cerrar alerta tras recuperación: disco casi lleno»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La señal de disco casi lleno cierra o actualiza el incidente asociado.

**Alternativas y errores:**
- Sin lectura reciente: no inferir que el servicio se recuperó.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La señal de disco casi lleno cierra o actualiza el incidente asociado.
- Se reproduce y controla esta alternativa: Sin lectura reciente: no inferir que el servicio se recuperó.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de cerrar alerta tras recuperación: disco casi lleno; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-140"></a>
## MSP-140 · Cerrar alerta tras recuperación: servicio de respaldo

**Área:** monitoring · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La señal de servicio de respaldo cierra o actualiza el incidente asociado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comprobar recuperación y vincular evidencia de monitorización para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.

**Flujo principal:**
1. El actor accede a descubrimiento y monitoreo dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comprobar recuperación y vincular evidencia de monitorización para servicio de respaldo: Zabbix, nodo lognorte-bkp, tarea nocturna.
3. Ejecuta «Cerrar alerta tras recuperación: servicio de respaldo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La señal de servicio de respaldo cierra o actualiza el incidente asociado.

**Alternativas y errores:**
- Sin lectura reciente: no inferir que el servicio se recuperó.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La señal de servicio de respaldo cierra o actualiza el incidente asociado.
- Se reproduce y controla esta alternativa: Sin lectura reciente: no inferir que el servicio se recuperó.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Descubrimiento y monitoreo: operaciones de cerrar alerta tras recuperación: servicio de respaldo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-141"></a>
## MSP-141 · Solicitar sesión remota: portátil del contacto

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sesión de portátil del contacto requiere consentimiento del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir dispositivo y motivo de intervención para portátil del contacto: TR-104, consentimiento de la usuaria.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir dispositivo y motivo de intervención para portátil del contacto: TR-104, consentimiento de la usuaria.
3. Ejecuta «Solicitar sesión remota: portátil del contacto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sesión de portátil del contacto requiere consentimiento del cliente.

**Alternativas y errores:**
- Dispositivo no inventariado: impedir conexión remota.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sesión de portátil del contacto requiere consentimiento del cliente.
- Se reproduce y controla esta alternativa: Dispositivo no inventariado: impedir conexión remota.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de solicitar sesión remota: portátil del contacto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-142"></a>
## MSP-142 · Solicitar sesión remota: servidor de aplicación

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sesión de servidor de aplicación requiere consentimiento del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir dispositivo y motivo de intervención para servidor de aplicación: SR-08, ventana 20:00–21:00.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir dispositivo y motivo de intervención para servidor de aplicación: SR-08, ventana 20:00–21:00.
3. Ejecuta «Solicitar sesión remota: servidor de aplicación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sesión de servidor de aplicación requiere consentimiento del cliente.

**Alternativas y errores:**
- Dispositivo no inventariado: impedir conexión remota.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sesión de servidor de aplicación requiere consentimiento del cliente.
- Se reproduce y controla esta alternativa: Dispositivo no inventariado: impedir conexión remota.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de solicitar sesión remota: servidor de aplicación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-143"></a>
## MSP-143 · Solicitar sesión remota: terminal de almacén

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sesión de terminal de almacén requiere consentimiento del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir dispositivo y motivo de intervención para terminal de almacén: TB-31, intervención supervisada.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir dispositivo y motivo de intervención para terminal de almacén: TB-31, intervención supervisada.
3. Ejecuta «Solicitar sesión remota: terminal de almacén»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sesión de terminal de almacén requiere consentimiento del cliente.

**Alternativas y errores:**
- Dispositivo no inventariado: impedir conexión remota.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sesión de terminal de almacén requiere consentimiento del cliente.
- Se reproduce y controla esta alternativa: Dispositivo no inventariado: impedir conexión remota.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de solicitar sesión remota: terminal de almacén; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-144"></a>
## MSP-144 · Solicitar sesión remota: router de sede

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La sesión de router de sede requiere consentimiento del cliente.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Elegir dispositivo y motivo de intervención para router de sede: SW-22, cambio previamente aprobado.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Elegir dispositivo y motivo de intervención para router de sede: SW-22, cambio previamente aprobado.
3. Ejecuta «Solicitar sesión remota: router de sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La sesión de router de sede requiere consentimiento del cliente.

**Alternativas y errores:**
- Dispositivo no inventariado: impedir conexión remota.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La sesión de router de sede requiere consentimiento del cliente.
- Se reproduce y controla esta alternativa: Dispositivo no inventariado: impedir conexión remota.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de solicitar sesión remota: router de sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-145"></a>
## MSP-145 · Autorizar acceso temporal: portátil del contacto

**Área:** remote-support · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La autorización de portátil del contacto expira en el límite aprobado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Confirmar técnico, duración y alcance de la sesión para portátil del contacto: TR-104, consentimiento de la usuaria.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Confirmar técnico, duración y alcance de la sesión para portátil del contacto: TR-104, consentimiento de la usuaria.
3. Ejecuta «Autorizar acceso temporal: portátil del contacto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La autorización de portátil del contacto expira en el límite aprobado.

**Alternativas y errores:**
- Aprobación vencida o revocada: negar la conexión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La autorización de portátil del contacto expira en el límite aprobado.
- Se reproduce y controla esta alternativa: Aprobación vencida o revocada: negar la conexión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de autorizar acceso temporal: portátil del contacto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-146"></a>
## MSP-146 · Autorizar acceso temporal: servidor de aplicación

**Área:** remote-support · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La autorización de servidor de aplicación expira en el límite aprobado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Confirmar técnico, duración y alcance de la sesión para servidor de aplicación: SR-08, ventana 20:00–21:00.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Confirmar técnico, duración y alcance de la sesión para servidor de aplicación: SR-08, ventana 20:00–21:00.
3. Ejecuta «Autorizar acceso temporal: servidor de aplicación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La autorización de servidor de aplicación expira en el límite aprobado.

**Alternativas y errores:**
- Aprobación vencida o revocada: negar la conexión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La autorización de servidor de aplicación expira en el límite aprobado.
- Se reproduce y controla esta alternativa: Aprobación vencida o revocada: negar la conexión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de autorizar acceso temporal: servidor de aplicación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-147"></a>
## MSP-147 · Autorizar acceso temporal: terminal de almacén

**Área:** remote-support · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La autorización de terminal de almacén expira en el límite aprobado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Confirmar técnico, duración y alcance de la sesión para terminal de almacén: TB-31, intervención supervisada.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Confirmar técnico, duración y alcance de la sesión para terminal de almacén: TB-31, intervención supervisada.
3. Ejecuta «Autorizar acceso temporal: terminal de almacén»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La autorización de terminal de almacén expira en el límite aprobado.

**Alternativas y errores:**
- Aprobación vencida o revocada: negar la conexión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La autorización de terminal de almacén expira en el límite aprobado.
- Se reproduce y controla esta alternativa: Aprobación vencida o revocada: negar la conexión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de autorizar acceso temporal: terminal de almacén; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-148"></a>
## MSP-148 · Autorizar acceso temporal: router de sede

**Área:** remote-support · **Actor:** Contacto del cliente · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La autorización de router de sede expira en el límite aprobado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Confirmar técnico, duración y alcance de la sesión para router de sede: SW-22, cambio previamente aprobado.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Confirmar técnico, duración y alcance de la sesión para router de sede: SW-22, cambio previamente aprobado.
3. Ejecuta «Autorizar acceso temporal: router de sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La autorización de router de sede expira en el límite aprobado.

**Alternativas y errores:**
- Aprobación vencida o revocada: negar la conexión.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La autorización de router de sede expira en el límite aprobado.
- Se reproduce y controla esta alternativa: Aprobación vencida o revocada: negar la conexión.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de autorizar acceso temporal: router de sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-149"></a>
## MSP-149 · Iniciar sesión registrada: portátil del contacto

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El acceso a portátil del contacto deja registro de inicio y técnico.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conectar mediante el conector administrado y autorizado para portátil del contacto: TR-104, consentimiento de la usuaria.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conectar mediante el conector administrado y autorizado para portátil del contacto: TR-104, consentimiento de la usuaria.
3. Ejecuta «Iniciar sesión registrada: portátil del contacto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El acceso a portátil del contacto deja registro de inicio y técnico.

**Alternativas y errores:**
- Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El acceso a portátil del contacto deja registro de inicio y técnico.
- Se reproduce y controla esta alternativa: Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de iniciar sesión registrada: portátil del contacto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-150"></a>
## MSP-150 · Iniciar sesión registrada: servidor de aplicación

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El acceso a servidor de aplicación deja registro de inicio y técnico.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conectar mediante el conector administrado y autorizado para servidor de aplicación: SR-08, ventana 20:00–21:00.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conectar mediante el conector administrado y autorizado para servidor de aplicación: SR-08, ventana 20:00–21:00.
3. Ejecuta «Iniciar sesión registrada: servidor de aplicación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El acceso a servidor de aplicación deja registro de inicio y técnico.

**Alternativas y errores:**
- Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El acceso a servidor de aplicación deja registro de inicio y técnico.
- Se reproduce y controla esta alternativa: Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de iniciar sesión registrada: servidor de aplicación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-151"></a>
## MSP-151 · Iniciar sesión registrada: terminal de almacén

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El acceso a terminal de almacén deja registro de inicio y técnico.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conectar mediante el conector administrado y autorizado para terminal de almacén: TB-31, intervención supervisada.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conectar mediante el conector administrado y autorizado para terminal de almacén: TB-31, intervención supervisada.
3. Ejecuta «Iniciar sesión registrada: terminal de almacén»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El acceso a terminal de almacén deja registro de inicio y técnico.

**Alternativas y errores:**
- Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El acceso a terminal de almacén deja registro de inicio y técnico.
- Se reproduce y controla esta alternativa: Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de iniciar sesión registrada: terminal de almacén; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-152"></a>
## MSP-152 · Iniciar sesión registrada: router de sede

**Área:** remote-support · **Actor:** Técnico del proveedor · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El acceso a router de sede deja registro de inicio y técnico.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Conectar mediante el conector administrado y autorizado para router de sede: SW-22, cambio previamente aprobado.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Conectar mediante el conector administrado y autorizado para router de sede: SW-22, cambio previamente aprobado.
3. Ejecuta «Iniciar sesión registrada: router de sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El acceso a router de sede deja registro de inicio y técnico.

**Alternativas y errores:**
- Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El acceso a router de sede deja registro de inicio y técnico.
- Se reproduce y controla esta alternativa: Credencial de integración inválida: fallar cerrado y alertar al administrador.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de iniciar sesión registrada: router de sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-153"></a>
## MSP-153 · Revocar sesión en curso: portátil del contacto

**Área:** remote-support · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conexión a portátil del contacto se cierra y queda motivo registrado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Detener sesión por solicitud del cliente para portátil del contacto: TR-104, consentimiento de la usuaria.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Detener sesión por solicitud del cliente para portátil del contacto: TR-104, consentimiento de la usuaria.
3. Ejecuta «Revocar sesión en curso: portátil del contacto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conexión a portátil del contacto se cierra y queda motivo registrado.

**Alternativas y errores:**
- Estado remoto desconocido: no reportar desconexión confirmada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conexión a portátil del contacto se cierra y queda motivo registrado.
- Se reproduce y controla esta alternativa: Estado remoto desconocido: no reportar desconexión confirmada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de revocar sesión en curso: portátil del contacto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-154"></a>
## MSP-154 · Revocar sesión en curso: servidor de aplicación

**Área:** remote-support · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conexión a servidor de aplicación se cierra y queda motivo registrado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Detener sesión por solicitud del cliente para servidor de aplicación: SR-08, ventana 20:00–21:00.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Detener sesión por solicitud del cliente para servidor de aplicación: SR-08, ventana 20:00–21:00.
3. Ejecuta «Revocar sesión en curso: servidor de aplicación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conexión a servidor de aplicación se cierra y queda motivo registrado.

**Alternativas y errores:**
- Estado remoto desconocido: no reportar desconexión confirmada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conexión a servidor de aplicación se cierra y queda motivo registrado.
- Se reproduce y controla esta alternativa: Estado remoto desconocido: no reportar desconexión confirmada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de revocar sesión en curso: servidor de aplicación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-155"></a>
## MSP-155 · Revocar sesión en curso: terminal de almacén

**Área:** remote-support · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conexión a terminal de almacén se cierra y queda motivo registrado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Detener sesión por solicitud del cliente para terminal de almacén: TB-31, intervención supervisada.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Detener sesión por solicitud del cliente para terminal de almacén: TB-31, intervención supervisada.
3. Ejecuta «Revocar sesión en curso: terminal de almacén»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conexión a terminal de almacén se cierra y queda motivo registrado.

**Alternativas y errores:**
- Estado remoto desconocido: no reportar desconexión confirmada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conexión a terminal de almacén se cierra y queda motivo registrado.
- Se reproduce y controla esta alternativa: Estado remoto desconocido: no reportar desconexión confirmada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de revocar sesión en curso: terminal de almacén; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-156"></a>
## MSP-156 · Revocar sesión en curso: router de sede

**Área:** remote-support · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La conexión a router de sede se cierra y queda motivo registrado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Detener sesión por solicitud del cliente para router de sede: SW-22, cambio previamente aprobado.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Detener sesión por solicitud del cliente para router de sede: SW-22, cambio previamente aprobado.
3. Ejecuta «Revocar sesión en curso: router de sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La conexión a router de sede se cierra y queda motivo registrado.

**Alternativas y errores:**
- Estado remoto desconocido: no reportar desconexión confirmada.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La conexión a router de sede se cierra y queda motivo registrado.
- Se reproduce y controla esta alternativa: Estado remoto desconocido: no reportar desconexión confirmada.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de revocar sesión en curso: router de sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-157"></a>
## MSP-157 · Consultar auditoría remota: portátil del contacto

**Área:** remote-support · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La auditoría de portátil del contacto muestra metadatos sin credenciales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar por cliente, técnico y periodo permitido para portátil del contacto: TR-104, consentimiento de la usuaria.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar por cliente, técnico y periodo permitido para portátil del contacto: TR-104, consentimiento de la usuaria.
3. Ejecuta «Consultar auditoría remota: portátil del contacto»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La auditoría de portátil del contacto muestra metadatos sin credenciales.

**Alternativas y errores:**
- Rol sin alcance de auditoría: denegar consulta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La auditoría de portátil del contacto muestra metadatos sin credenciales.
- Se reproduce y controla esta alternativa: Rol sin alcance de auditoría: denegar consulta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de consultar auditoría remota: portátil del contacto; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-158"></a>
## MSP-158 · Consultar auditoría remota: servidor de aplicación

**Área:** remote-support · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La auditoría de servidor de aplicación muestra metadatos sin credenciales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar por cliente, técnico y periodo permitido para servidor de aplicación: SR-08, ventana 20:00–21:00.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar por cliente, técnico y periodo permitido para servidor de aplicación: SR-08, ventana 20:00–21:00.
3. Ejecuta «Consultar auditoría remota: servidor de aplicación»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La auditoría de servidor de aplicación muestra metadatos sin credenciales.

**Alternativas y errores:**
- Rol sin alcance de auditoría: denegar consulta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La auditoría de servidor de aplicación muestra metadatos sin credenciales.
- Se reproduce y controla esta alternativa: Rol sin alcance de auditoría: denegar consulta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de consultar auditoría remota: servidor de aplicación; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-159"></a>
## MSP-159 · Consultar auditoría remota: terminal de almacén

**Área:** remote-support · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La auditoría de terminal de almacén muestra metadatos sin credenciales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar por cliente, técnico y periodo permitido para terminal de almacén: TB-31, intervención supervisada.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar por cliente, técnico y periodo permitido para terminal de almacén: TB-31, intervención supervisada.
3. Ejecuta «Consultar auditoría remota: terminal de almacén»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La auditoría de terminal de almacén muestra metadatos sin credenciales.

**Alternativas y errores:**
- Rol sin alcance de auditoría: denegar consulta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La auditoría de terminal de almacén muestra metadatos sin credenciales.
- Se reproduce y controla esta alternativa: Rol sin alcance de auditoría: denegar consulta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de consultar auditoría remota: terminal de almacén; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-160"></a>
## MSP-160 · Consultar auditoría remota: router de sede

**Área:** remote-support · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La auditoría de router de sede muestra metadatos sin credenciales.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Filtrar por cliente, técnico y periodo permitido para router de sede: SW-22, cambio previamente aprobado.

**Flujo principal:**
1. El actor accede a soporte remoto integrado dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Filtrar por cliente, técnico y periodo permitido para router de sede: SW-22, cambio previamente aprobado.
3. Ejecuta «Consultar auditoría remota: router de sede»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La auditoría de router de sede muestra metadatos sin credenciales.

**Alternativas y errores:**
- Rol sin alcance de auditoría: denegar consulta.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La auditoría de router de sede muestra metadatos sin credenciales.
- Se reproduce y controla esta alternativa: Rol sin alcance de auditoría: denegar consulta.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Soporte remoto integrado: operaciones de consultar auditoría remota: router de sede; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-161"></a>
## MSP-161 · Calcular margen por cliente: octubre Aurora

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El margen de octubre Aurora usa datos conciliados y periodo declarado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar ingresos del contrato con costos aprobados para octubre Aurora: ingreso USD 4 000, costos por conciliar.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar ingresos del contrato con costos aprobados para octubre Aurora: ingreso USD 4 000, costos por conciliar.
3. Ejecuta «Calcular margen por cliente: octubre Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El margen de octubre Aurora usa datos conciliados y periodo declarado.

**Alternativas y errores:**
- Costos incompletos: rotular el margen como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El margen de octubre Aurora usa datos conciliados y periodo declarado.
- Se reproduce y controla esta alternativa: Costos incompletos: rotular el margen como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de calcular margen por cliente: octubre aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-162"></a>
## MSP-162 · Calcular margen por cliente: servicio de guardia

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El margen de servicio de guardia usa datos conciliados y periodo declarado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar ingresos del contrato con costos aprobados para servicio de guardia: tres clientes, periodo trimestral.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar ingresos del contrato con costos aprobados para servicio de guardia: tres clientes, periodo trimestral.
3. Ejecuta «Calcular margen por cliente: servicio de guardia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El margen de servicio de guardia usa datos conciliados y periodo declarado.

**Alternativas y errores:**
- Costos incompletos: rotular el margen como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El margen de servicio de guardia usa datos conciliados y periodo declarado.
- Se reproduce y controla esta alternativa: Costos incompletos: rotular el margen como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de calcular margen por cliente: servicio de guardia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-163"></a>
## MSP-163 · Calcular margen por cliente: proyecto de migración

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El margen de proyecto de migración usa datos conciliados y periodo declarado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar ingresos del contrato con costos aprobados para proyecto de migración: 120 h estimadas frente a 88 aprobadas.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar ingresos del contrato con costos aprobados para proyecto de migración: 120 h estimadas frente a 88 aprobadas.
3. Ejecuta «Calcular margen por cliente: proyecto de migración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El margen de proyecto de migración usa datos conciliados y periodo declarado.

**Alternativas y errores:**
- Costos incompletos: rotular el margen como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El margen de proyecto de migración usa datos conciliados y periodo declarado.
- Se reproduce y controla esta alternativa: Costos incompletos: rotular el margen como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de calcular margen por cliente: proyecto de migración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-164"></a>
## MSP-164 · Calcular margen por cliente: renovaciones de diciembre

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El margen de renovaciones de diciembre usa datos conciliados y periodo declarado.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar ingresos del contrato con costos aprobados para renovaciones de diciembre: monedas USD y EUR en reportes separados.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar ingresos del contrato con costos aprobados para renovaciones de diciembre: monedas USD y EUR en reportes separados.
3. Ejecuta «Calcular margen por cliente: renovaciones de diciembre»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El margen de renovaciones de diciembre usa datos conciliados y periodo declarado.

**Alternativas y errores:**
- Costos incompletos: rotular el margen como provisional.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El margen de renovaciones de diciembre usa datos conciliados y periodo declarado.
- Se reproduce y controla esta alternativa: Costos incompletos: rotular el margen como provisional.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de calcular margen por cliente: renovaciones de diciembre; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-165"></a>
## MSP-165 · Comparar rentabilidad de servicios: octubre Aurora

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La comparación de octubre Aurora separa contratos y moneda.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar ingresos y costos por paquete de servicio para octubre Aurora: ingreso USD 4 000, costos por conciliar.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar ingresos y costos por paquete de servicio para octubre Aurora: ingreso USD 4 000, costos por conciliar.
3. Ejecuta «Comparar rentabilidad de servicios: octubre Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La comparación de octubre Aurora separa contratos y moneda.

**Alternativas y errores:**
- Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La comparación de octubre Aurora separa contratos y moneda.
- Se reproduce y controla esta alternativa: Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de comparar rentabilidad de servicios: octubre aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-166"></a>
## MSP-166 · Comparar rentabilidad de servicios: servicio de guardia

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La comparación de servicio de guardia separa contratos y moneda.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar ingresos y costos por paquete de servicio para servicio de guardia: tres clientes, periodo trimestral.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar ingresos y costos por paquete de servicio para servicio de guardia: tres clientes, periodo trimestral.
3. Ejecuta «Comparar rentabilidad de servicios: servicio de guardia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La comparación de servicio de guardia separa contratos y moneda.

**Alternativas y errores:**
- Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La comparación de servicio de guardia separa contratos y moneda.
- Se reproduce y controla esta alternativa: Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de comparar rentabilidad de servicios: servicio de guardia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-167"></a>
## MSP-167 · Comparar rentabilidad de servicios: proyecto de migración

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La comparación de proyecto de migración separa contratos y moneda.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar ingresos y costos por paquete de servicio para proyecto de migración: 120 h estimadas frente a 88 aprobadas.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar ingresos y costos por paquete de servicio para proyecto de migración: 120 h estimadas frente a 88 aprobadas.
3. Ejecuta «Comparar rentabilidad de servicios: proyecto de migración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La comparación de proyecto de migración separa contratos y moneda.

**Alternativas y errores:**
- Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La comparación de proyecto de migración separa contratos y moneda.
- Se reproduce y controla esta alternativa: Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de comparar rentabilidad de servicios: proyecto de migración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-168"></a>
## MSP-168 · Comparar rentabilidad de servicios: renovaciones de diciembre

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La comparación de renovaciones de diciembre separa contratos y moneda.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Agrupar ingresos y costos por paquete de servicio para renovaciones de diciembre: monedas USD y EUR en reportes separados.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Agrupar ingresos y costos por paquete de servicio para renovaciones de diciembre: monedas USD y EUR en reportes separados.
3. Ejecuta «Comparar rentabilidad de servicios: renovaciones de diciembre»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La comparación de renovaciones de diciembre separa contratos y moneda.

**Alternativas y errores:**
- Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La comparación de renovaciones de diciembre separa contratos y moneda.
- Se reproduce y controla esta alternativa: Mezcla de monedas: exigir conversión con tipo de cambio registrado.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de comparar rentabilidad de servicios: renovaciones de diciembre; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-169"></a>
## MSP-169 · Preparar resumen mensual: octubre Aurora

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El resumen de octubre Aurora informa métricas y fecha de corte.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Seleccionar periodo y clientes autorizados para octubre Aurora: ingreso USD 4 000, costos por conciliar.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Seleccionar periodo y clientes autorizados para octubre Aurora: ingreso USD 4 000, costos por conciliar.
3. Ejecuta «Preparar resumen mensual: octubre Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El resumen de octubre Aurora informa métricas y fecha de corte.

**Alternativas y errores:**
- Periodo abierto: distinguir valores preliminares de cerrados.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El resumen de octubre Aurora informa métricas y fecha de corte.
- Se reproduce y controla esta alternativa: Periodo abierto: distinguir valores preliminares de cerrados.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de preparar resumen mensual: octubre aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-170"></a>
## MSP-170 · Preparar resumen mensual: servicio de guardia

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El resumen de servicio de guardia informa métricas y fecha de corte.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Seleccionar periodo y clientes autorizados para servicio de guardia: tres clientes, periodo trimestral.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Seleccionar periodo y clientes autorizados para servicio de guardia: tres clientes, periodo trimestral.
3. Ejecuta «Preparar resumen mensual: servicio de guardia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El resumen de servicio de guardia informa métricas y fecha de corte.

**Alternativas y errores:**
- Periodo abierto: distinguir valores preliminares de cerrados.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El resumen de servicio de guardia informa métricas y fecha de corte.
- Se reproduce y controla esta alternativa: Periodo abierto: distinguir valores preliminares de cerrados.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de preparar resumen mensual: servicio de guardia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-171"></a>
## MSP-171 · Preparar resumen mensual: proyecto de migración

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El resumen de proyecto de migración informa métricas y fecha de corte.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Seleccionar periodo y clientes autorizados para proyecto de migración: 120 h estimadas frente a 88 aprobadas.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Seleccionar periodo y clientes autorizados para proyecto de migración: 120 h estimadas frente a 88 aprobadas.
3. Ejecuta «Preparar resumen mensual: proyecto de migración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El resumen de proyecto de migración informa métricas y fecha de corte.

**Alternativas y errores:**
- Periodo abierto: distinguir valores preliminares de cerrados.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El resumen de proyecto de migración informa métricas y fecha de corte.
- Se reproduce y controla esta alternativa: Periodo abierto: distinguir valores preliminares de cerrados.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de preparar resumen mensual: proyecto de migración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-172"></a>
## MSP-172 · Preparar resumen mensual: renovaciones de diciembre

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El resumen de renovaciones de diciembre informa métricas y fecha de corte.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Seleccionar periodo y clientes autorizados para renovaciones de diciembre: monedas USD y EUR en reportes separados.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Seleccionar periodo y clientes autorizados para renovaciones de diciembre: monedas USD y EUR en reportes separados.
3. Ejecuta «Preparar resumen mensual: renovaciones de diciembre»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El resumen de renovaciones de diciembre informa métricas y fecha de corte.

**Alternativas y errores:**
- Periodo abierto: distinguir valores preliminares de cerrados.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El resumen de renovaciones de diciembre informa métricas y fecha de corte.
- Se reproduce y controla esta alternativa: Periodo abierto: distinguir valores preliminares de cerrados.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de preparar resumen mensual: renovaciones de diciembre; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-173"></a>
## MSP-173 · Exportar datos comerciales: octubre Aurora

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de octubre Aurora aplica permisos y filtros usados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Descargar importes mínimos para contabilidad para octubre Aurora: ingreso USD 4 000, costos por conciliar.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Descargar importes mínimos para contabilidad para octubre Aurora: ingreso USD 4 000, costos por conciliar.
3. Ejecuta «Exportar datos comerciales: octubre Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de octubre Aurora aplica permisos y filtros usados.

**Alternativas y errores:**
- Fórmula en texto exportado: neutralizar contenido ejecutable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de octubre Aurora aplica permisos y filtros usados.
- Se reproduce y controla esta alternativa: Fórmula en texto exportado: neutralizar contenido ejecutable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de exportar datos comerciales: octubre aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-174"></a>
## MSP-174 · Exportar datos comerciales: servicio de guardia

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de servicio de guardia aplica permisos y filtros usados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Descargar importes mínimos para contabilidad para servicio de guardia: tres clientes, periodo trimestral.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Descargar importes mínimos para contabilidad para servicio de guardia: tres clientes, periodo trimestral.
3. Ejecuta «Exportar datos comerciales: servicio de guardia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de servicio de guardia aplica permisos y filtros usados.

**Alternativas y errores:**
- Fórmula en texto exportado: neutralizar contenido ejecutable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de servicio de guardia aplica permisos y filtros usados.
- Se reproduce y controla esta alternativa: Fórmula en texto exportado: neutralizar contenido ejecutable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de exportar datos comerciales: servicio de guardia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-175"></a>
## MSP-175 · Exportar datos comerciales: proyecto de migración

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de proyecto de migración aplica permisos y filtros usados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Descargar importes mínimos para contabilidad para proyecto de migración: 120 h estimadas frente a 88 aprobadas.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Descargar importes mínimos para contabilidad para proyecto de migración: 120 h estimadas frente a 88 aprobadas.
3. Ejecuta «Exportar datos comerciales: proyecto de migración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de proyecto de migración aplica permisos y filtros usados.

**Alternativas y errores:**
- Fórmula en texto exportado: neutralizar contenido ejecutable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de proyecto de migración aplica permisos y filtros usados.
- Se reproduce y controla esta alternativa: Fórmula en texto exportado: neutralizar contenido ejecutable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de exportar datos comerciales: proyecto de migración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-176"></a>
## MSP-176 · Exportar datos comerciales: renovaciones de diciembre

**Área:** profitability · **Actor:** Gestor de contratos y costes · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La exportación de renovaciones de diciembre aplica permisos y filtros usados.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Descargar importes mínimos para contabilidad para renovaciones de diciembre: monedas USD y EUR en reportes separados.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Descargar importes mínimos para contabilidad para renovaciones de diciembre: monedas USD y EUR en reportes separados.
3. Ejecuta «Exportar datos comerciales: renovaciones de diciembre»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La exportación de renovaciones de diciembre aplica permisos y filtros usados.

**Alternativas y errores:**
- Fórmula en texto exportado: neutralizar contenido ejecutable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La exportación de renovaciones de diciembre aplica permisos y filtros usados.
- Se reproduce y controla esta alternativa: Fórmula en texto exportado: neutralizar contenido ejecutable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de exportar datos comerciales: renovaciones de diciembre; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-177"></a>
## MSP-177 · Revisar desviación de horas: octubre Aurora

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La desviación de octubre Aurora identifica contratos que requieren revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas estimadas, aprobadas y consumidas para octubre Aurora: ingreso USD 4 000, costos por conciliar.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas estimadas, aprobadas y consumidas para octubre Aurora: ingreso USD 4 000, costos por conciliar.
3. Ejecuta «Revisar desviación de horas: octubre Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La desviación de octubre Aurora identifica contratos que requieren revisión.

**Alternativas y errores:**
- Horas no aprobadas: excluirlas del total facturable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La desviación de octubre Aurora identifica contratos que requieren revisión.
- Se reproduce y controla esta alternativa: Horas no aprobadas: excluirlas del total facturable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de revisar desviación de horas: octubre aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-178"></a>
## MSP-178 · Revisar desviación de horas: servicio de guardia

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La desviación de servicio de guardia identifica contratos que requieren revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas estimadas, aprobadas y consumidas para servicio de guardia: tres clientes, periodo trimestral.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas estimadas, aprobadas y consumidas para servicio de guardia: tres clientes, periodo trimestral.
3. Ejecuta «Revisar desviación de horas: servicio de guardia»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La desviación de servicio de guardia identifica contratos que requieren revisión.

**Alternativas y errores:**
- Horas no aprobadas: excluirlas del total facturable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La desviación de servicio de guardia identifica contratos que requieren revisión.
- Se reproduce y controla esta alternativa: Horas no aprobadas: excluirlas del total facturable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de revisar desviación de horas: servicio de guardia; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-179"></a>
## MSP-179 · Revisar desviación de horas: proyecto de migración

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La desviación de proyecto de migración identifica contratos que requieren revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas estimadas, aprobadas y consumidas para proyecto de migración: 120 h estimadas frente a 88 aprobadas.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas estimadas, aprobadas y consumidas para proyecto de migración: 120 h estimadas frente a 88 aprobadas.
3. Ejecuta «Revisar desviación de horas: proyecto de migración»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La desviación de proyecto de migración identifica contratos que requieren revisión.

**Alternativas y errores:**
- Horas no aprobadas: excluirlas del total facturable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La desviación de proyecto de migración identifica contratos que requieren revisión.
- Se reproduce y controla esta alternativa: Horas no aprobadas: excluirlas del total facturable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de revisar desviación de horas: proyecto de migración; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-180"></a>
## MSP-180 · Revisar desviación de horas: renovaciones de diciembre

**Área:** profitability · **Actor:** Responsable de servicio · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La desviación de renovaciones de diciembre identifica contratos que requieren revisión.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Comparar horas estimadas, aprobadas y consumidas para renovaciones de diciembre: monedas USD y EUR en reportes separados.

**Flujo principal:**
1. El actor accede a rentabilidad e informes comerciales dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Comparar horas estimadas, aprobadas y consumidas para renovaciones de diciembre: monedas USD y EUR en reportes separados.
3. Ejecuta «Revisar desviación de horas: renovaciones de diciembre»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La desviación de renovaciones de diciembre identifica contratos que requieren revisión.

**Alternativas y errores:**
- Horas no aprobadas: excluirlas del total facturable.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La desviación de renovaciones de diciembre identifica contratos que requieren revisión.
- Se reproduce y controla esta alternativa: Horas no aprobadas: excluirlas del total facturable.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Rentabilidad e informes comerciales: operaciones de revisar desviación de horas: renovaciones de diciembre; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-181"></a>
## MSP-181 · Crear credencial de integración: Zabbix Aurora

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La credencial de Zabbix Aurora se muestra una vez y queda acotada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asignar permisos mínimos y fecha de expiración para Zabbix Aurora: token de lectura, expiración 30 días.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asignar permisos mínimos y fecha de expiración para Zabbix Aurora: token de lectura, expiración 30 días.
3. Ejecuta «Crear credencial de integración: Zabbix Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La credencial de Zabbix Aurora se muestra una vez y queda acotada.

**Alternativas y errores:**
- Alcance administrativo no justificado: no emitir la credencial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La credencial de Zabbix Aurora se muestra una vez y queda acotada.
- Se reproduce y controla esta alternativa: Alcance administrativo no justificado: no emitir la credencial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de crear credencial de integración: zabbix aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-182"></a>
## MSP-182 · Crear credencial de integración: MeshCentral Norte

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La credencial de MeshCentral Norte se muestra una vez y queda acotada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asignar permisos mínimos y fecha de expiración para MeshCentral Norte: alcance para dispositivos administrados.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asignar permisos mínimos y fecha de expiración para MeshCentral Norte: alcance para dispositivos administrados.
3. Ejecuta «Crear credencial de integración: MeshCentral Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La credencial de MeshCentral Norte se muestra una vez y queda acotada.

**Alternativas y errores:**
- Alcance administrativo no justificado: no emitir la credencial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La credencial de MeshCentral Norte se muestra una vez y queda acotada.
- Se reproduce y controla esta alternativa: Alcance administrativo no justificado: no emitir la credencial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de crear credencial de integración: meshcentral norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-183"></a>
## MSP-183 · Crear credencial de integración: contabilidad MSP

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La credencial de contabilidad MSP se muestra una vez y queda acotada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asignar permisos mínimos y fecha de expiración para contabilidad MSP: webhook firmado de facturas.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asignar permisos mínimos y fecha de expiración para contabilidad MSP: webhook firmado de facturas.
3. Ejecuta «Crear credencial de integración: contabilidad MSP»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La credencial de contabilidad MSP se muestra una vez y queda acotada.

**Alternativas y errores:**
- Alcance administrativo no justificado: no emitir la credencial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La credencial de contabilidad MSP se muestra una vez y queda acotada.
- Se reproduce y controla esta alternativa: Alcance administrativo no justificado: no emitir la credencial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de crear credencial de integración: contabilidad msp; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-184"></a>
## MSP-184 · Crear credencial de integración: inventario externo

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La credencial de inventario externo se muestra una vez y queda acotada.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Asignar permisos mínimos y fecha de expiración para inventario externo: solo consulta de activos del Hotel Río.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Asignar permisos mínimos y fecha de expiración para inventario externo: solo consulta de activos del Hotel Río.
3. Ejecuta «Crear credencial de integración: inventario externo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La credencial de inventario externo se muestra una vez y queda acotada.

**Alternativas y errores:**
- Alcance administrativo no justificado: no emitir la credencial.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La credencial de inventario externo se muestra una vez y queda acotada.
- Se reproduce y controla esta alternativa: Alcance administrativo no justificado: no emitir la credencial.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de crear credencial de integración: inventario externo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-185"></a>
## MSP-185 · Verificar firma de webhook: Zabbix Aurora

**Área:** integrations · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de Zabbix Aurora se acepta sólo si su origen es íntegro.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Validar firma, marca temporal e identificador del evento para Zabbix Aurora: token de lectura, expiración 30 días.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Validar firma, marca temporal e identificador del evento para Zabbix Aurora: token de lectura, expiración 30 días.
3. Ejecuta «Verificar firma de webhook: Zabbix Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de Zabbix Aurora se acepta sólo si su origen es íntegro.

**Alternativas y errores:**
- Firma inválida o repetida: rechazar sin aplicar efectos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de Zabbix Aurora se acepta sólo si su origen es íntegro.
- Se reproduce y controla esta alternativa: Firma inválida o repetida: rechazar sin aplicar efectos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de verificar firma de webhook: zabbix aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-186"></a>
## MSP-186 · Verificar firma de webhook: MeshCentral Norte

**Área:** integrations · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de MeshCentral Norte se acepta sólo si su origen es íntegro.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Validar firma, marca temporal e identificador del evento para MeshCentral Norte: alcance para dispositivos administrados.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Validar firma, marca temporal e identificador del evento para MeshCentral Norte: alcance para dispositivos administrados.
3. Ejecuta «Verificar firma de webhook: MeshCentral Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de MeshCentral Norte se acepta sólo si su origen es íntegro.

**Alternativas y errores:**
- Firma inválida o repetida: rechazar sin aplicar efectos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de MeshCentral Norte se acepta sólo si su origen es íntegro.
- Se reproduce y controla esta alternativa: Firma inválida o repetida: rechazar sin aplicar efectos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de verificar firma de webhook: meshcentral norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-187"></a>
## MSP-187 · Verificar firma de webhook: contabilidad MSP

**Área:** integrations · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de contabilidad MSP se acepta sólo si su origen es íntegro.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Validar firma, marca temporal e identificador del evento para contabilidad MSP: webhook firmado de facturas.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Validar firma, marca temporal e identificador del evento para contabilidad MSP: webhook firmado de facturas.
3. Ejecuta «Verificar firma de webhook: contabilidad MSP»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de contabilidad MSP se acepta sólo si su origen es íntegro.

**Alternativas y errores:**
- Firma inválida o repetida: rechazar sin aplicar efectos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de contabilidad MSP se acepta sólo si su origen es íntegro.
- Se reproduce y controla esta alternativa: Firma inválida o repetida: rechazar sin aplicar efectos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de verificar firma de webhook: contabilidad msp; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-188"></a>
## MSP-188 · Verificar firma de webhook: inventario externo

**Área:** integrations · **Actor:** Servicio de integración autorizado · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El evento de inventario externo se acepta sólo si su origen es íntegro.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Validar firma, marca temporal e identificador del evento para inventario externo: solo consulta de activos del Hotel Río.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Validar firma, marca temporal e identificador del evento para inventario externo: solo consulta de activos del Hotel Río.
3. Ejecuta «Verificar firma de webhook: inventario externo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El evento de inventario externo se acepta sólo si su origen es íntegro.

**Alternativas y errores:**
- Firma inválida o repetida: rechazar sin aplicar efectos.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El evento de inventario externo se acepta sólo si su origen es íntegro.
- Se reproduce y controla esta alternativa: Firma inválida o repetida: rechazar sin aplicar efectos.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de verificar firma de webhook: inventario externo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-189"></a>
## MSP-189 · Rotar secreto de conector: Zabbix Aurora

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La rotación de Zabbix Aurora registra la transición sin exponer secretos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Reemplazar clave y revocar la anterior para Zabbix Aurora: token de lectura, expiración 30 días.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Reemplazar clave y revocar la anterior para Zabbix Aurora: token de lectura, expiración 30 días.
3. Ejecuta «Rotar secreto de conector: Zabbix Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La rotación de Zabbix Aurora registra la transición sin exponer secretos.

**Alternativas y errores:**
- Conector no validado: conservar clave previa hasta prueba explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La rotación de Zabbix Aurora registra la transición sin exponer secretos.
- Se reproduce y controla esta alternativa: Conector no validado: conservar clave previa hasta prueba explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de rotar secreto de conector: zabbix aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-190"></a>
## MSP-190 · Rotar secreto de conector: MeshCentral Norte

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La rotación de MeshCentral Norte registra la transición sin exponer secretos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Reemplazar clave y revocar la anterior para MeshCentral Norte: alcance para dispositivos administrados.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Reemplazar clave y revocar la anterior para MeshCentral Norte: alcance para dispositivos administrados.
3. Ejecuta «Rotar secreto de conector: MeshCentral Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La rotación de MeshCentral Norte registra la transición sin exponer secretos.

**Alternativas y errores:**
- Conector no validado: conservar clave previa hasta prueba explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La rotación de MeshCentral Norte registra la transición sin exponer secretos.
- Se reproduce y controla esta alternativa: Conector no validado: conservar clave previa hasta prueba explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de rotar secreto de conector: meshcentral norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-191"></a>
## MSP-191 · Rotar secreto de conector: contabilidad MSP

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La rotación de contabilidad MSP registra la transición sin exponer secretos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Reemplazar clave y revocar la anterior para contabilidad MSP: webhook firmado de facturas.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Reemplazar clave y revocar la anterior para contabilidad MSP: webhook firmado de facturas.
3. Ejecuta «Rotar secreto de conector: contabilidad MSP»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La rotación de contabilidad MSP registra la transición sin exponer secretos.

**Alternativas y errores:**
- Conector no validado: conservar clave previa hasta prueba explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La rotación de contabilidad MSP registra la transición sin exponer secretos.
- Se reproduce y controla esta alternativa: Conector no validado: conservar clave previa hasta prueba explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de rotar secreto de conector: contabilidad msp; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-192"></a>
## MSP-192 · Rotar secreto de conector: inventario externo

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** La rotación de inventario externo registra la transición sin exponer secretos.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Reemplazar clave y revocar la anterior para inventario externo: solo consulta de activos del Hotel Río.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Reemplazar clave y revocar la anterior para inventario externo: solo consulta de activos del Hotel Río.
3. Ejecuta «Rotar secreto de conector: inventario externo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: La rotación de inventario externo registra la transición sin exponer secretos.

**Alternativas y errores:**
- Conector no validado: conservar clave previa hasta prueba explícita.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: La rotación de inventario externo registra la transición sin exponer secretos.
- Se reproduce y controla esta alternativa: Conector no validado: conservar clave previa hasta prueba explícita.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de rotar secreto de conector: inventario externo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-193"></a>
## MSP-193 · Limitar acceso por cliente: Zabbix Aurora

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El conector de Zabbix Aurora no puede consultar otros clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar clientes y operaciones permitidas para Zabbix Aurora: token de lectura, expiración 30 días.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar clientes y operaciones permitidas para Zabbix Aurora: token de lectura, expiración 30 días.
3. Ejecuta «Limitar acceso por cliente: Zabbix Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El conector de Zabbix Aurora no puede consultar otros clientes.

**Alternativas y errores:**
- Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El conector de Zabbix Aurora no puede consultar otros clientes.
- Se reproduce y controla esta alternativa: Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de limitar acceso por cliente: zabbix aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-194"></a>
## MSP-194 · Limitar acceso por cliente: MeshCentral Norte

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El conector de MeshCentral Norte no puede consultar otros clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar clientes y operaciones permitidas para MeshCentral Norte: alcance para dispositivos administrados.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar clientes y operaciones permitidas para MeshCentral Norte: alcance para dispositivos administrados.
3. Ejecuta «Limitar acceso por cliente: MeshCentral Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El conector de MeshCentral Norte no puede consultar otros clientes.

**Alternativas y errores:**
- Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El conector de MeshCentral Norte no puede consultar otros clientes.
- Se reproduce y controla esta alternativa: Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de limitar acceso por cliente: meshcentral norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-195"></a>
## MSP-195 · Limitar acceso por cliente: contabilidad MSP

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El conector de contabilidad MSP no puede consultar otros clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar clientes y operaciones permitidas para contabilidad MSP: webhook firmado de facturas.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar clientes y operaciones permitidas para contabilidad MSP: webhook firmado de facturas.
3. Ejecuta «Limitar acceso por cliente: contabilidad MSP»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El conector de contabilidad MSP no puede consultar otros clientes.

**Alternativas y errores:**
- Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El conector de contabilidad MSP no puede consultar otros clientes.
- Se reproduce y controla esta alternativa: Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de limitar acceso por cliente: contabilidad msp; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-196"></a>
## MSP-196 · Limitar acceso por cliente: inventario externo

**Área:** integrations · **Actor:** Administrador de organización · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El conector de inventario externo no puede consultar otros clientes.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Configurar clientes y operaciones permitidas para inventario externo: solo consulta de activos del Hotel Río.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Configurar clientes y operaciones permitidas para inventario externo: solo consulta de activos del Hotel Río.
3. Ejecuta «Limitar acceso por cliente: inventario externo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El conector de inventario externo no puede consultar otros clientes.

**Alternativas y errores:**
- Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El conector de inventario externo no puede consultar otros clientes.
- Se reproduce y controla esta alternativa: Solicitud fuera del alcance: denegar y registrar evento mínimo.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de limitar acceso por cliente: inventario externo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-197"></a>
## MSP-197 · Revisar fallos de integración: Zabbix Aurora

**Área:** integrations · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de Zabbix Aurora distingue éxito, fallo y falta de ejecución.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Consultar estado, reintentos y última entrega para Zabbix Aurora: token de lectura, expiración 30 días.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Consultar estado, reintentos y última entrega para Zabbix Aurora: token de lectura, expiración 30 días.
3. Ejecuta «Revisar fallos de integración: Zabbix Aurora»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de Zabbix Aurora distingue éxito, fallo y falta de ejecución.

**Alternativas y errores:**
- Payload con datos sensibles: ocultarlo en logs de soporte.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de Zabbix Aurora distingue éxito, fallo y falta de ejecución.
- Se reproduce y controla esta alternativa: Payload con datos sensibles: ocultarlo en logs de soporte.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de revisar fallos de integración: zabbix aurora; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-198"></a>
## MSP-198 · Revisar fallos de integración: MeshCentral Norte

**Área:** integrations · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de MeshCentral Norte distingue éxito, fallo y falta de ejecución.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Consultar estado, reintentos y última entrega para MeshCentral Norte: alcance para dispositivos administrados.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Consultar estado, reintentos y última entrega para MeshCentral Norte: alcance para dispositivos administrados.
3. Ejecuta «Revisar fallos de integración: MeshCentral Norte»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de MeshCentral Norte distingue éxito, fallo y falta de ejecución.

**Alternativas y errores:**
- Payload con datos sensibles: ocultarlo en logs de soporte.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de MeshCentral Norte distingue éxito, fallo y falta de ejecución.
- Se reproduce y controla esta alternativa: Payload con datos sensibles: ocultarlo en logs de soporte.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de revisar fallos de integración: meshcentral norte; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-199"></a>
## MSP-199 · Revisar fallos de integración: contabilidad MSP

**Área:** integrations · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de contabilidad MSP distingue éxito, fallo y falta de ejecución.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Consultar estado, reintentos y última entrega para contabilidad MSP: webhook firmado de facturas.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Consultar estado, reintentos y última entrega para contabilidad MSP: webhook firmado de facturas.
3. Ejecuta «Revisar fallos de integración: contabilidad MSP»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de contabilidad MSP distingue éxito, fallo y falta de ejecución.

**Alternativas y errores:**
- Payload con datos sensibles: ocultarlo en logs de soporte.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de contabilidad MSP distingue éxito, fallo y falta de ejecución.
- Se reproduce y controla esta alternativa: Payload con datos sensibles: ocultarlo en logs de soporte.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de revisar fallos de integración: contabilidad msp; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0

<a id="msp-200"></a>
## MSP-200 · Revisar fallos de integración: inventario externo

**Área:** integrations · **Actor:** Auditor de solo lectura · **Prioridad:** P1 · **Estado:** `planned`

**Objetivo:** El estado de inventario externo distingue éxito, fallo y falta de ejecución.

**Ejemplo:** Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río. Consultar estado, reintentos y última entrega para inventario externo: solo consulta de activos del Hotel Río.

**Flujo principal:**
1. El actor accede a integraciones, api y seguridad dentro de su organización y selecciona los recursos autorizados.
2. Introduce o selecciona: Consultar estado, reintentos y última entrega para inventario externo: solo consulta de activos del Hotel Río.
3. Ejecuta «Revisar fallos de integración: inventario externo»; el servidor valida permisos y reglas del módulo.
4. El sistema presenta el resultado verificable: El estado de inventario externo distingue éxito, fallo y falta de ejecución.

**Alternativas y errores:**
- Payload con datos sensibles: ocultarlo en logs de soporte.
- Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.

**Aceptación:**
- Dado el ejemplo, al completar el flujo se verifica: El estado de inventario externo distingue éxito, fallo y falta de ejecución.
- Se reproduce y controla esta alternativa: Payload con datos sensibles: ocultarlo en logs de soporte.
- Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.

**Contrato:** MSP / Integraciones, API y seguridad: operaciones de revisar fallos de integración: inventario externo; definir en OpenAPI al implementar el caso. `(proposed)`
**Evidencia:** implementación 0 · pruebas 0 · operación 0
