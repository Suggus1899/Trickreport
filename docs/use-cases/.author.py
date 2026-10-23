import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent

ROLES = {
    'admin': 'Administrador de organización', 'user': 'Solicitante',
    'agent': 'Agente', 'manager': 'Responsable de servicio',
    'approver': 'Aprobador designado', 'asset': 'Gestor de activos',
    'change': 'Gestor de cambios', 'auditor': 'Auditor de solo lectura',
    'contact': 'Contacto del cliente', 'tech': 'Técnico del proveedor',
    'finance': 'Gestor de contratos y costes', 'customer': 'Cliente externo',
    'supervisor': 'Supervisor de atención', 'system': 'Servicio de integración autorizado',
}

def write(product, name, groups):
    cases = []
    allocations = []
    for area, label, expected, phase, prerequisites, entities, rows in groups:
        lines = [line.strip() for line in rows.strip().splitlines() if line.strip()]
        assert len(lines) == expected, (product, area, len(lines), expected)
        allocations.append({'id': area, 'name': label, 'count': expected})
        for line in lines:
            title, actor, example, outcome, exception = line.split('|')
            case_id = f'{product}-{len(cases)+1:03}'
            cases.append({
                'id': case_id, 'product': product, 'area': area, 'title': title,
                'actor': {'role': actor, 'name': ROLES[actor]},
                'objective': outcome,
                'permissions': [f'{product.lower()}.{area}:authorized-action',
                    'Comprobar permiso de la acción y alcance del recurso en el servidor; el rol visible no autoriza por sí solo.'],
                'preconditions': ['Organización activa y actor autenticado o integración verificada dentro de este producto.', prerequisites],
                'example': {'context': {'TI':'Organización ficticia Clínica Aurora; sede Caracas.', 'MSP':'Proveedor ficticio Andes IT; cliente autorizado Clínica Aurora, distinto de cliente Hotel Río.', 'CX':'Empresa ficticia Tienda Faro; cliente externo Ana Pérez, distinto de Luis Díaz.'}[product], 'input': example, 'expected_result': outcome},
                'main_flow': [f'El actor accede a {label.lower()} dentro de su organización y selecciona los recursos autorizados.',
                    f'Introduce o selecciona: {example}', f'Ejecuta «{title}»; el servidor valida permisos y reglas del módulo.',
                    f'El sistema presenta el resultado verificable: {outcome}'],
                'exception_flows': [exception, 'Si la sesión caduca o el recurso pertenece a otro ámbito, se deniega la operación sin revelar datos del recurso.'],
                'acceptance_criteria': [f'Dado el ejemplo, al completar el flujo se verifica: {outcome}',
                    f'Se reproduce y controla esta alternativa: {exception}',
                    'Una identidad sin el permiso o alcance necesario no obtiene el recurso ni produce efectos secundarios.'],
                'dependencies': ['product-identity-and-scope', area],
                'surfaces': [label, f'Detalle / formulario para {title.lower()}'],
                'entities': entities.split(','),
                'interfaces': {'status': 'proposed', 'contract': f'{product} / {label}: operaciones de {title.lower()}; definir en OpenAPI al implementar el caso.'},
                'priority': 'P1' if phase == 2 else 'P2', 'phase': phase,
                'status': 'planned', 'evidence': {'implementation': [], 'tests': [], 'operations': []},
                'documentation': f'{product.lower()}.md#{case_id.lower()}',
            })
    assert len(cases) == 200, (product, len(cases))
    data = {'schema_version': 1, 'product': product, 'name': name,
        'language': 'es', 'created_on': '2026-10-23',
        'status_note': 'Catálogo objetivo aprobado. planned no significa capacidad implementada ni validada; enlazar evidencia antes de avanzar estado.',
        'areas': allocations, 'use_cases': cases}
    (ROOT / f'{product.lower()}.json').write_text(json.dumps(data, ensure_ascii=False, indent=2)+'\n', encoding='utf-8')
    md = [f'# {name}', '', f'Catálogo de **{len(cases)} casos propuestos**. Todos comienzan en `planned`; su inclusión no afirma que la función exista.', '', '| Área | Casos |', '| --- | ---: |']
    md.extend(f'| {area["name"]} | {area["count"]} |' for area in allocations)
    for case in cases:
        md.extend(['', f'<a id="{case["id"].lower()}"></a>', f'## {case["id"]} · {case["title"]}', '',
            f'**Área:** {case["area"]} · **Actor:** {case["actor"]["name"]} · **Prioridad:** {case["priority"]} · **Estado:** `{case["status"]}`', '',
            f'**Objetivo:** {case["objective"]}', '', f'**Ejemplo:** {case["example"]["context"]} {case["example"]["input"]}', '',
            '**Flujo principal:**', *[f'{index}. {step}' for index, step in enumerate(case['main_flow'], 1)], '',
            '**Alternativas y errores:**', *[f'- {step}' for step in case['exception_flows']], '',
            '**Aceptación:**', *[f'- {criterion}' for criterion in case['acceptance_criteria']], '',
            f'**Contrato:** {case["interfaces"]["contract"]} `({case["interfaces"]["status"]})`',
            f'**Evidencia:** implementación {len(case["evidence"]["implementation"])} · pruebas {len(case["evidence"]["tests"])} · operación {len(case["evidence"]["operations"])}'])
    (ROOT / f'{product.lower()}.md').write_text('\n'.join(md)+'\n', encoding='utf-8')

write('TI', 'Trickreport · Soporte TI', [
('identity', 'Identidad, organización y permisos', 16, 2, 'Organización Aurora creada; métodos de acceso configurados.', 'organization,user,session,role', '''
Invitar empleado a la organización|admin|ana@aurora.example, departamento Radiología|Ana recibe una invitación limitada a Aurora|Correo ya vinculado: informar sin duplicar membresía.
Iniciar sesión local|user|Ana presenta su correo y contraseña válidos|Se crea una sesión revocable para Ana|Contraseña incorrecta: respuesta uniforme y límite de intentos.
Acceder mediante directorio LDAP|user|Usuario aperez pertenece al grupo TI-Aurora|Se verifica el directorio y se aplica el rol autorizado|Directorio inaccesible: informar indisponibilidad sin conceder acceso.
Activar segundo factor|user|Ana confirma un código TOTP de su dispositivo|Los accesos posteriores exigen el segundo factor|Código incorrecto: no activar el secreto propuesto.
Completar desafío MFA|user|Sesión pendiente y código TOTP de seis dígitos|Se emite una sesión tras consumir el desafío|Desafío usado o vencido: rechazar su reutilización.
Recuperar contraseña local|user|Ana solicita recuperación y recibe un enlace válido|Nueva contraseña aplicada y sesiones anteriores revocadas|Enlace vencido: permitir pedir otro sin cambiar credenciales.
Cerrar una sesión concreta|user|Ana identifica el navegador del equipo compartido|La sesión elegida deja de autorizar peticiones|Sesión de otra persona: denegar su revocación.
Cambiar permisos de un agente|admin|Pedro pasa de agente a solicitante|Nuevas peticiones pierden acceso a colas de agentes|Último administrador: impedir dejar la organización sin gestión.
Desactivar a un empleado|admin|Pedro termina su relación laboral el viernes|Acceso revocado conservando autoría de tickets históricos|Tickets asignados: mostrar pendientes para reasignación.
Crear un departamento|admin|Departamento Laboratorio con responsable Marta|Se dispone de agrupación organizativa para solicitudes|Nombre duplicado dentro de Aurora: pedir diferenciación.
Asignar una sede|admin|Ana trabaja en sede Valencia|Su perfil permite contextualizar servicios y activos|Sede inactiva: impedir nuevas asignaciones.
Delegar administración de un equipo|admin|Marta administra únicamente el equipo Redes|Marta puede gestionar miembros de Redes|Intento sobre otro equipo: denegar sin modificar miembros.
Actualizar preferencias personales|user|Ana selecciona inglés y zona America/Caracas|Idioma y horarios se muestran según su perfil|Zona inválida: conservar la preferencia previa.
Consultar actividad de acceso|user|Ana revisa sesiones abiertas durante siete días|Ve fecha y dispositivo de sus propios accesos|Datos de otro usuario: mantenerlos fuera del resultado.
Rotar credencial de una integración|admin|Conector inventario tiene una clave comprometida|La nueva clave sustituye a la revocada con trazabilidad|Cliente usa clave anterior: denegar sin fallback secreto.
Verificar aislamiento de organización|auditor|Dos organizaciones tienen un ticket con número local 42|Una identidad Aurora solo obtiene recursos de Aurora|Identificador ajeno enviado manualmente: respuesta sin datos ajenos.
'''),
('incidents', 'Incidentes y seguimiento', 32, 2, 'Categorías y equipos de atención configurados.', 'ticket,comment,attachment,notification', '''
Registrar incidente desde el portal|user|Ana no puede imprimir desde Radiología|Se crea un ticket con referencia y confirmación|Formulario incompleto: conservar texto y señalar campos necesarios.
Registrar incidente en nombre de un empleado|agent|Pedro recibe llamada de Ana por acceso bloqueado|Ticket identifica solicitante Ana y registrador Pedro|Empleado fuera de organización: impedir su selección.
Clasificar impacto y urgencia|agent|Historia clínica caída para 80 empleados durante consulta|La prioridad refleja impacto amplio y urgencia alta|Valores incompatibles: solicitar corrección antes de guardar.
Asignar incidente a un técnico|manager|Ticket 42 requiere a Pedro del equipo Redes|Pedro recibe responsabilidad y aviso del ticket|Técnico inactivo: impedir asignación.
Transferir incidente entre equipos|agent|Fallo de impresora corresponde a Infraestructura|Equipo receptor recibe contexto e historial íntegros|Transferencia sin permiso: conservar el equipo actual.
Aceptar trabajo asignado|agent|Pedro toma el ticket 42 pendiente en su cola|El ticket entra en atención con fecha de inicio|Otro agente lo tomó antes: mostrar conflicto actualizado.
Solicitar información al empleado|agent|Se necesita captura del error de impresión|Pregunta pública y estado de espera quedan registrados|Comentario vacío: no cambiar estado.
Responder una solicitud de información|user|Ana adjunta código de error E204 del ticket 42|El agente recibe la respuesta y reanuda revisión|Ticket cerrado: aplicar política de reapertura visible.
Añadir nota interna de diagnóstico|agent|Pedro documenta credenciales rotadas en servidor afectado|Nota interna visible únicamente para personal autorizado|Solicitante consulta el hilo: nunca recibe esa nota.
Adjuntar evidencia a un incidente|user|Ana carga captura PNG de 300 KB|Archivo queda vinculado y descargable con autorización|Tipo real prohibido o tamaño excesivo: rechazar archivo.
Retirar un adjunto incorrecto|agent|Captura adjunta contiene datos personales innecesarios|Se retira su acceso y queda registro de la retirada|Adjunto no pertenece al ticket: denegar operación.
Relacionar incidentes duplicados|agent|Tickets 42 y 45 describen la misma caída|Un ticket principal conserva referencias a los duplicados|Tickets de otra organización: impedir vinculación.
Dividir un reporte con varios problemas|agent|Un ticket mezcla impresión y acceso al correo|Se crean tareas de atención separadas enlazadas al original|Fallo al crear una parte: evitar un estado intermedio incoherente.
Declarar incidente mayor|manager|Caída general del sistema de turnos|Se activa coordinación y seguimiento de incidente mayor|Incidente ya declarado: evitar equipos y avisos duplicados.
Publicar actualización de incidente mayor|manager|Restablecimiento parcial a las 11:20 en Radiología|Afectados autorizados reciben una actualización común|Destinatario fuera de alcance: excluirlo del envío.
Suscribirse a un incidente|user|Ana sigue la caída de impresión de su sede|Recibe novedades públicas del incidente autorizado|Pérdida de acceso posterior: retirar futuras notificaciones.
Registrar solución aplicada|agent|Se reinstala controlador firmado de impresora R3|Ticket contiene pasos, responsable y hora de resolución|Solución vacía: impedir marcar como resuelto.
Confirmar resolución de un incidente|user|Ana imprime correctamente después de la reparación|Su confirmación queda asociada al cierre|Incidente aún activo: solicitar al agente resolución previa.
Reabrir un incidente resuelto|user|La impresora falla nuevamente una hora después|Ticket vuelve a atención conservando la solución anterior|Fuera de ventana: crear seguimiento vinculado según política.
Cancelar un incidente reportado por error|user|Ana descubre que eligió la impresora equivocada|Ticket cancelado conserva motivo y auditoría|Trabajo ya facturado o ligado a cambio: requerir revisión autorizada.
Consultar línea temporal del incidente|user|Ticket 42 cambió tres veces de estado|Se muestran eventos públicos ordenados con autor y fecha|Eventos internos: filtrarlos antes de responder al cliente.
Buscar incidentes por texto|agent|Consulta impresora R3 y sede Caracas|Resultados pertenecen al alcance autorizado y muestran coincidencias|Consulta vacía: usar listado paginado sin búsqueda global irrestricta.
Guardar una vista de trabajo|agent|Filtro mis tickets urgentes sin respuesta|La vista vuelve a aplicar el filtro al abrirla|Campo de filtro eliminado: mostrar ajuste necesario.
Atender cola de incidentes sin asignar|manager|Cinco tickets nuevos carecen de responsable|Se prioriza una cola con antigüedad y SLA visibles|Ticket asignado mientras se abre: actualizar la cola.
Escalar un incidente a especialista|agent|Diagnóstico de red requiere equipo de Seguridad|Se registra motivo y equipo de escalación|Sin equipo disponible: mantener propietario y alertar responsable.
Registrar tiempo de atención|agent|Pedro dedica 35 minutos al ticket 42|Duración queda atribuida al agente y actividad|Tiempo negativo o solapado prohibido: rechazar registro.
Crear borrador sin conexión|user|Ana redacta incidente desde una zona sin cobertura|Borrador persistido bajo su identidad y organización|Cambio de usuario: ocultar borradores de identidad anterior.
Sincronizar borrador pendiente|user|Vuelve conexión con un borrador y clave idempotente|Se crea exactamente un ticket y se marca sincronizado|Timeout posterior al alta: reintentar la misma clave sin duplicar.
Descartar un borrador offline|user|Reporte local quedó obsoleto antes de enviarse|Se elimina solo ese borrador y su sincronización pendiente|Sincronización ya en curso: coordinar cancelación y mostrar resultado.
Recibir actualización en tiempo real|agent|Otro agente asigna el ticket abierto por Pedro|La vista refleja nueva asignación sin recargar manualmente|WebSocket desconectado: recuperar estado desde API al reconectar.
Aplicar cambios masivos de cola|manager|Tres tickets de impresoras pasan a Infraestructura|Cada cambio autorizado informa éxito o conflicto individual|Ticket modificado simultáneamente: no sobrescribirlo silenciosamente.
Consultar incidentes relacionados con un activo|agent|Servidor PACS-01 presenta varias interrupciones|Se listan incidentes autorizados del mismo activo|Activo retirado: conservar acceso histórico con permisos vigentes.
'''),
('service-catalog', 'Solicitudes y catálogo', 22, 3, 'Servicios y responsables definidos por la organización.', 'service,request,form,fulfillment-task', '''
Publicar un servicio solicitables|manager|Servicio Alta de correo para nuevos empleados|Aparece en el catálogo con responsable y plazo|Falta responsable: impedir publicación.
Definir formulario de una solicitud|manager|Alta de correo requiere nombre, área y fecha|El formulario solicita los datos necesarios por tipo de servicio|Campo incompatible con datos existentes: conservar versión anterior.
Solicitar acceso a una aplicación|user|Ana necesita acceso de consulta a PACS|Se registra una solicitud con alcance de acceso explícito|Permiso privilegiado sin justificación: pedir justificación obligatoria.
Solicitar equipo de trabajo|user|Nueva enfermera necesita portátil para el lunes|Solicitud informa modelo requerido y fecha objetivo|Equipo sin disponibilidad: mostrar espera sin prometer entrega.
Solicitar instalación de software|user|Ana necesita visor DICOM autorizado|Se crea tarea con software y equipo destino|Software fuera del catálogo: enviar a revisión autorizada.
Solicitar incorporación de empleado|manager|Luis ingresa el 2 de noviembre a Laboratorio|Se agrupan tareas de cuenta, equipo y accesos|Fecha anterior al día actual: pedir confirmación de regularización.
Solicitar baja de empleado|manager|Pedro finaliza contrato el 30 de noviembre|Se coordinan revocación de accesos y devolución de activos|Identidad ambigua: impedir una baja sobre persona incorrecta.
Solicitar traslado de puesto|manager|Ana cambia de Radiología a Laboratorio|Se revisan accesos, activos y pertenencias del nuevo puesto|Acceso incompatible con nueva función: retirarlo según aprobación.
Solicitar ampliación de almacenamiento|user|Servidor de imágenes requiere 200 GB adicionales|Solicitud indica capacidad, motivo y activo afectado|Supera cuota aprobada: requerir aprobación de capacidad.
Solicitar acceso temporal|user|Consultor necesita VPN durante cinco días|La solicitud incorpora una fecha de caducidad|Fecha final anterior a inicial: rechazar periodo.
Solicitar devolución de un activo|user|Ana devuelve portátil AUR-021 al almacén|Se abre una tarea de recepción y verificación|Activo asignado a otra persona: pedir autorización del custodio.
Consultar estado de una solicitud|user|Alta de correo tiene dos tareas completas de tres|Se muestra progreso público por etapas|Tareas internas sensibles: resumir sin exponer instrucciones privadas.
Completar tarea de provisión|agent|Cuenta de correo creada para Luis|Quedan resultado y evidencia vinculados a la tarea|Dependencia previa sin completar: bloquear cierre prematuro.
Entregar recurso al solicitante|agent|Portátil AUR-021 queda preparado para Luis|Entrega queda registrada con receptor y fecha|Receptor rechaza equipo: mantener entrega pendiente.
Confirmar recepción del servicio|user|Luis verifica correo y portátil entregados|La solicitud conserva aceptación del destinatario|Servicio incompleto: permitir observaciones y devolución a atención.
Programar ejecución de solicitud|agent|Instalación de visor será viernes a las 17:00|Solicitante conoce ventana y técnico asignado|Ventana fuera del horario autorizado: pedir otra franja.
Cancelar solicitud antes de ejecución|user|La incorporación de Luis fue cancelada|Se cancelan tareas pendientes sin ocultar trabajo ya realizado|Provisión ya iniciada: derivar al responsable para compensación.
Duplicar solicitud como nueva petición|user|Otra enfermera requiere el mismo tipo de equipo|Se reutilizan campos permitidos con nuevo destinatario y referencia|Datos sensibles del original: no copiarlos al nuevo destinatario.
Retirar servicio del catálogo|manager|Visor antiguo deja de distribuirse|Se impiden nuevas solicitudes conservando las abiertas|Solicitud en curso: permitir completar o migrar con motivo.
Versionar condiciones de un servicio|manager|Entrega de portátil pasa de cinco a siete días|Nuevas solicitudes usan versión nueva; antiguas conservan condiciones|Edición concurrente: informar conflicto de versión.
Restringir catálogo por departamento|manager|Acceso PACS disponible solo para áreas clínicas|Otros departamentos no ven ni solicitan el servicio restringido|URL directa desde área ajena: denegar solicitud.
Consultar coste estimado de solicitud|user|Licencia anual cuesta 120 unidades monetarias|El solicitante ve estimación antes de enviarla|Precio ausente: mostrar pendiente de estimación, no cero ficticio.
'''),
('knowledge', 'Conocimiento y autoservicio', 14, 2, 'Autores autorizados y categorías de conocimiento disponibles.', 'article,revision,feedback', '''
Redactar artículo de solución|agent|Procedimiento para restablecer cola de impresión|Se guarda borrador con pasos y categoría|Contenido inválido: conservar borrador sin publicar.
Publicar artículo revisado|manager|Borrador revisado sobre impresora R3|La versión aprobada aparece en autoservicio|Revisión pendiente: impedir publicación directa.
Buscar solución antes de reportar|user|Ana busca error E204 de impresión|Encuentra artículos visibles relacionados con el error|Sin coincidencias: ofrecer crear ticket con búsqueda conservada.
Insertar solución de conocimiento en respuesta|agent|Artículo aprobado explica reinicio de cola|Respuesta enlaza la versión consultada|Artículo retirado: advertir y pedir contenido vigente.
Valorar utilidad de un artículo|user|Ana resolvió la impresión siguiendo el artículo|Se registra valoración de utilidad de esa versión|Envío repetido: actualizar una valoración por identidad.
Reportar información obsoleta|user|Ruta de configuración ya no existe en Windows|El autor recibe observación vinculada al artículo|Comentario ofensivo: aplicar reglas de moderación visibles.
Editar una nueva revisión|agent|Se actualiza procedimiento por cambio de controlador|Borrador nuevo no altera la publicación vigente|Edición concurrente: mostrar conflicto de versión.
Comparar revisiones de un artículo|manager|Revisión 3 cambia pasos de reinicio respecto a 2|Se distinguen cambios antes de aprobar|Revisión inexistente: informar sin sustituir por otra.
Restaurar revisión conocida|manager|Revisión 4 introdujo instrucciones equivocadas|Se publica una nueva revisión basada en la 3 con trazabilidad|Revisión incompatible con servicio actual: advertir antes de publicar.
Restringir artículo a agentes|manager|Diagnóstico incluye detalles de red interna|Solo personal autorizado puede buscar o leer el artículo|Enlace compartido con empleado: denegar contenido interno.
Organizar categorías de conocimiento|manager|Impresión se agrupa bajo Puesto de trabajo|Artículos mantienen enlaces y aparecen en categoría correcta|Eliminar categoría con artículos: exigir reasignación.
Mantener traducción de una guía|agent|Artículo español tiene versión inglesa revisada|Ambas versiones conservan relación y estado de actualización|Original cambia: marcar traducción pendiente de revisión.
Programar revisión de contenido|manager|Guía VPN requiere revisión cada 90 días|Responsable recibe tarea de revisión a su vencimiento|Responsable inactivo: escalar al gestor de conocimiento.
Retirar artículo peligroso|manager|Procedimiento permite borrar configuración necesaria|Se retira publicación conservando auditoría y referencia histórica|Ticket lo enlazaba: mostrar aviso de artículo retirado.
'''),
('sla', 'SLA, calendarios y escalaciones', 20, 2, 'Prioridades y calendario de servicio definidos.', 'sla-policy,calendar,ticket-deadline,escalation', '''
Definir objetivo de primera respuesta|manager|Prioridad alta exige respuesta en 30 minutos|Política asigna plazo persistente a tickets nuevos|Tiempo no positivo: rechazar la política.
Definir objetivo de resolución|manager|Incidente normal debe resolverse en ocho horas hábiles|Se calcula vencimiento de resolución según calendario|Calendario inexistente: impedir activar política.
Crear horario de atención|admin|Equipo atiende lunes a viernes de 08:00 a 17:00|El cómputo excluye horas fuera de atención|Franjas solapadas: solicitar corrección.
Registrar festivo de servicio|admin|25 de diciembre queda fuera del calendario laboral|Vencimientos nuevos excluyen el festivo|Calendario compartido: mostrar servicios afectados antes de guardar.
Aplicar SLA por prioridad|manager|Incidente de impacto alto usa política Urgente|La política seleccionada queda trazable en el ticket|Dos reglas empatan: usar precedencia explícita configurada.
Pausar SLA por espera autorizada|agent|Se espera respuesta de Ana desde las 10:00|Pausa queda fechada con motivo admitido|Estado no permite pausa: conservar el contador activo.
Reanudar SLA tras respuesta|system|Ana aporta información a las 11:15|Se reanuda cómputo sin contar pausa autorizada|Evento repetido: no descontar dos veces el intervalo.
Consultar cuenta regresiva de SLA|agent|Ticket 42 vence en 25 minutos|La vista muestra fecha y tiempo restante coherentes|Reloj cliente desviado: tomar vencimiento del servidor.
Notificar proximidad de incumplimiento|system|Quedan diez minutos para primera respuesta|El responsable recibe una alerta deduplicada|Canal temporalmente caído: conservar aviso para reintento.
Escalar incumplimiento a supervisor|system|Ticket 42 supera su plazo de respuesta|Supervisor recibe contexto y se registra escalación|Worker reinicia: no repetir escalación ya registrada.
Registrar primera respuesta efectiva|agent|Pedro envía primera respuesta pública a Ana|Se fija una sola marca de primera respuesta|Nota interna previa: no contarla como respuesta pública.
Calcular duración de resolución|auditor|Ticket abierto lunes y resuelto martes con pausa|El informe refleja duración según calendario y pausas|Datos anteriores sin marcas: mostrar desconocido, no duración cero.
Recalcular objetivo por cambio autorizado|manager|Prioridad cambia de normal a crítica por mayor impacto|Nuevo plazo y motivo quedan auditados|Cambio retroactivo borra incumplimiento: conservar evidencia original.
Versionar una política SLA|manager|Objetivo normal cambia de ocho a seis horas|Nuevos tickets usan versión nueva conservando historial|Tickets abiertos: exigir opción explícita antes de recalcularlos.
Excluir ticket del cómputo con motivo|manager|Ticket de prueba fue creado por error operativo|Se excluye del indicador con justificación auditable|Agente sin privilegio: impedir exclusión.
Asignar calendario por sede|admin|Sede Valencia tiene horario distinto a Caracas|Tickets usan calendario de la sede correspondiente|Sede sin horario: aplicar calendario organizativo documentado.
Detectar tickets sin política|manager|Cinco tickets heredados carecen de SLA|Se presenta lista de cobertura pendiente|Sin fechas históricas: no inventar vencimientos anteriores.
Definir escalación por etapas|manager|A 100% alerta supervisor; a 150% alerta gerente|Cada etapa genera un único registro de escalación|Misma persona en dos etapas: conservar etapas sin spam duplicado.
Auditar una pausa de SLA|auditor|Ticket 42 estuvo pausado 75 minutos por espera|Se muestran inicio, fin, causa y actor|Pausa sin cierre por caída: identificar inconsistencia para reparación.
Comparar cumplimiento por servicio|manager|Impresión y VPN durante octubre|Se muestran numerador, denominador y periodos comparables|Muestra vacía: indicar sin datos en vez de 100%.
'''),
('assets', 'Activos y CMDB', 24, 3, 'Clases de activos y custodios habilitados.', 'asset,configuration-item,relationship,assignment,license', '''
Registrar activo físico|asset|Portátil AUR-021 con serie SN-8942|Se crea activo único con clase y ubicación|Serie duplicada: resolver coincidencia antes de crear.
Importar inventario de equipos|asset|CSV con 40 portátiles y números de serie|Se valida por fila y se informa resultado sin duplicados|Fila inválida: informar causa sin descartar filas válidas silenciosamente.
Asignar activo a empleado|asset|AUR-021 se entrega a Ana|Custodia registra persona, fecha y condición|Activo ya asignado: exigir transferencia o devolución.
Transferir custodia de activo|asset|Ana entrega AUR-021 a Luis|Historial conserva ambas custodias y fecha de cambio|Transferencia concurrente: impedir dos custodios activos.
Registrar devolución a almacén|asset|AUR-021 retorna con cargador y buen estado|Activo pasa a disponible tras recepción registrada|Accesorio faltante: dejar observación y estado de revisión.
Reservar equipo para incorporación|asset|AUR-025 reservado para Luis el lunes|Otro proceso no puede asignarlo durante la reserva|Reserva vencida: liberar con evento auditable.
Registrar garantía de fabricante|asset|AUR-021 cubierto hasta 30 de junio|Ficha muestra proveedor y vencimiento de garantía|Fecha anterior a compra: solicitar corrección.
Programar mantenimiento preventivo|asset|Servidor PACS-01 requiere limpieza trimestral|Se crean actividades con próximo vencimiento|Equipo retirado antes del vencimiento: cancelar futuras actividades.
Registrar reparación externa|asset|AUR-021 se envía a proveedor por batería|Se conserva custodia externa y número de servicio|Equipo con datos sensibles: requerir tratamiento autorizado previo.
Retirar activo de servicio|asset|AUR-009 es irreparable|Activo deja de ser asignable y conserva historial|Activo enlazado a servicio crítico: exigir revisión de dependencias.
Documentar disposición final|asset|AUR-009 se entrega a reciclaje con certificado|Se registra destino y evidencia del proceso|Falta autorización de disposición: no completar el estado final.
Registrar licencia de software|asset|Diez licencias del visor DICOM vencen en enero|Inventario muestra capacidad y vigencia|Cantidad negativa o contrato repetido: rechazar alta inconsistente.
Asignar licencia a un equipo|asset|Visor DICOM se asigna a AUR-021|Se consume una plaza disponible|Sin plazas: impedir sobreasignación y mostrar déficit.
Recuperar licencia no utilizada|asset|AUR-021 retirado libera su licencia DICOM|La plaza vuelve a disponible con motivo|Software aún activo confirmado: requerir retirada verificada.
Consultar software por activo|agent|Diagnóstico de AUR-021 requiere versión de visor|Se presenta inventario y fecha de observación|Observación antigua: indicar antigüedad del dato.
Crear elemento de configuración lógico|asset|Servicio Historia clínica depende de aplicación PACS|Se registra CI con responsable y criticidad|CI homónimo: pedir identificador que evite ambigüedad.
Relacionar dependencias de servicio|asset|PACS depende de base de datos DB-01|Mapa muestra dirección de dependencia verificable|Relación consigo mismo: rechazar vínculo inválido.
Analizar impacto de caída de activo|agent|DB-01 deja de responder|Se muestran servicios y CIs dependientes autorizados|Dependencias incompletas: señalar cobertura conocida.
Vincular activo a incidente|agent|Ticket 42 corresponde a AUR-021|Historial del activo enlaza el incidente|Activo de otra organización: denegar vinculación.
Conciliar inventario descubierto|asset|Escaneo detecta SN-8942 con nueva dirección IP|Se actualiza observación sin duplicar el activo|Coincidencias ambiguas: pedir resolución manual.
Detectar activo sin custodio|asset|AUR-030 figura activo sin asignación ni almacén|Aparece en lista de inconsistencias del inventario|Activo en reparación: reconocer custodia externa válida.
Auditar inventario de una sede|asset|Se verifican 25 equipos de Valencia|Se registran encontrados, ausentes y diferencias|Equipo encontrado en otra sede: abrir discrepancia, no moverlo en silencio.
Consultar coste de ciclo de vida|manager|AUR-021 costó 900 y acumuló 120 en reparaciones|Se muestra desglose de coste conocido|Monedas distintas: separar importes hasta tener conversión explícita.
Planificar renovación de equipos|manager|Ocho equipos superan cinco años y no tienen garantía|Se produce lista priorizada para renovación|Fecha de compra desconocida: marcar dato pendiente.
'''),
('approvals', 'Aprobaciones', 12, 3, 'Responsables y matriz de aprobación definidos.', 'approval,approval-step,delegation,request', '''
Pedir aprobación al responsable|user|Ana solicita licencia de 120 unidades|Su responsable recibe solicitud con coste y motivo|Sin responsable asignado: enviar a cola de resolución administrativa.
Aprobar solicitud con observaciones|approver|Marta autoriza licencia para Ana|Aprobación registra identidad, fecha y condiciones|Solicitud cambió de coste: exigir revisar nueva versión.
Rechazar solicitud justificada|approver|Licencia solicitada duplica una asignación activa|Se registra motivo y se informa al solicitante|Motivo vacío: impedir rechazo sin explicación.
Solicitar aclaración antes de aprobar|approver|Falta duración del acceso solicitado|Solicitud vuelve a aclaración conservando revisión previa|Solicitante inactivo: informar al responsable del proceso.
Resolver aprobación en varios niveles|manager|Compra de 5000 requiere responsable y finanzas|Se ejecutan etapas en orden con decisiones separadas|Primer nivel rechaza: no ejecutar provisión posterior.
Exigir separación de funciones|admin|Solicitante Ana también integra grupo aprobador|Ana no puede aprobar su propia solicitud restringida|Sin aprobadores alternativos: detener y escalar.
Delegar aprobación durante ausencia|approver|Marta delega a Luis del 2 al 6 de noviembre|Luis actúa durante ese intervalo con identidad propia|Delegación vencida: retirar facultad inmediatamente.
Escalar aprobación vencida|system|Aprobación sin respuesta supera 48 horas|Se avisa al suplente configurado conservando trazabilidad|Sin suplente: alertar administrador sin autoaprobar.
Retirar una petición de aprobación|user|Ana cancela la necesidad antes de decidirse|Aprobación pendiente queda retirada|Decisión ya emitida: registrar cancelación posterior por otro flujo.
Aplicar aprobación paralela|manager|Cambio requiere aceptación de Seguridad y Operaciones|Solo avanza cuando ambos grupos aprueban|Un grupo rechaza: bloquear avance aunque el otro apruebe.
Reabrir revisión por cambio material|system|Coste aprobado aumenta de 120 a 300|La autorización anterior queda invalidada para el nuevo importe|Cambio cosmético: conservar aprobación según política documentada.
Consultar historial de decisiones|auditor|Solicitud 81 pasó por dos responsables y un delegado|Se muestran versiones, decisiones y delegaciones completas|Comentario reservado: aplicar restricción del auditor.
'''),
('problems-changes', 'Problemas y cambios', 28, 3, 'Gestores designados y servicios registrados.', 'problem,known-error,change,change-task,review', '''
Abrir problema a partir de recurrencia|agent|Cinco incidentes de impresión ocurren en una semana|Se crea problema enlazado a los incidentes|Problema equivalente activo: proponer relación al existente.
Agrupar incidentes bajo un problema|agent|Tickets 42, 45 y 49 comparten error E204|El problema conserva evidencias de los incidentes|Incidente ajeno: impedir relación cruzada.
Priorizar investigación de problema|manager|Falla recurrente afecta 80 empleados|Se fija prioridad y responsable de investigación|Impacto sin evidencia: dejarlo pendiente de validación.
Registrar hipótesis de causa|agent|Se sospecha desbordamiento de cola del controlador|La hipótesis queda separada de la causa confirmada|Evidencia contradictoria: marcar hipótesis descartada con motivo.
Documentar causa raíz verificada|agent|Registro demuestra fuga de memoria en versión 4.2|Problema contiene causa y evidencia reproducible|Evidencia ausente: no declarar causa confirmada.
Publicar solución temporal|manager|Reiniciar cola restablece impresión hasta actualizar|Agentes disponen de mitigación enlazada al problema|Mitigación con riesgo: exigir advertencia y autorización correspondiente.
Registrar error conocido|manager|Defecto de controlador confirmado sin parche aplicado|Se publica error conocido con síntomas y mitigación|Causa aún hipotética: conservar estado de investigación.
Proponer corrección definitiva|agent|Actualizar controlador a versión corregida 4.3|La propuesta se enlaza al problema y cambio necesario|No existe versión validada: mantener propuesta pendiente.
Revisar eficacia de corrección|manager|Durante 14 días no reaparece E204|Se registra ventana observada y conclusión de eficacia|Reaparece incidente relacionado: mantener problema abierto.
Cerrar problema con evidencias|manager|Corrección verificada y documentación actualizada|Problema cierra conservando vínculos y causa|Acciones pendientes: impedir cierre definitivo.
Reabrir problema recurrente|agent|E204 reaparece después del cierre|Se inicia nueva investigación preservando diagnóstico anterior|Síntomas diferentes: crear problema distinto relacionado.
Consultar tendencias de problemas|manager|Se comparan fallos de impresión de tres meses|Se distinguen recurrencia, duración y servicios afectados|Periodo sin datos: no inferir mejora.
Solicitar cambio normal|change|Actualizar controlador de impresión corporativo|Cambio contiene objetivo, alcance y propietario|Falta servicio afectado: impedir enviar a evaluación.
Clasificar riesgo de cambio|change|Actualización afecta 40 equipos durante consulta|Se registra riesgo según impacto y reversibilidad|Respuestas incompletas: mantener evaluación pendiente.
Definir plan de implementación|change|Actualizar primero dos equipos piloto y luego el resto|Plan ordenado identifica responsables y verificaciones|Paso destructivo sin mitigación: exigir revisión de riesgo.
Definir plan de retorno|change|Restaurar controlador 4.2 si falla la impresión|Existe criterio de retorno y responsable ejecutor|Retorno técnicamente imposible: declarar riesgo y autorización especial.
Evaluar dependencias de cambio|change|Actualización depende de reinicio del servidor R3|Se vinculan prerrequisitos y CIs afectados|Dependencia circular: impedir programación incoherente.
Solicitar autorización del comité|change|Cambio de alto riesgo requiere comité de cambios|Se registran decisiones de los aprobadores requeridos|Quorum incompleto: no autorizar ejecución.
Programar ventana de cambio|change|Domingo de 02:00 a 03:00 para R3|Calendario muestra ventana y servicios afectados|Ventana coincide con cambio incompatible: informar conflicto.
Notificar intervención programada|manager|Usuarios de Radiología tendrán interrupción de diez minutos|Afectados reciben horario y alcance públicos|Contacto no autorizado: excluir datos internos del aviso.
Ejecutar tarea de cambio|agent|Pedro actualiza los dos equipos piloto|Se registra inicio, evidencia y resultado de la tarea|Validación piloto falla: detener despliegue restante.
Activar retorno de cambio|change|Versión nueva impide imprimir formularios críticos|Se ejecuta retorno y se registra resultado fallido del cambio|Retorno falla: escalar incidente mayor asociado.
Validar servicio tras cambio|agent|Ana imprime formulario y etiqueta desde R3|Pruebas de servicio quedan registradas antes del cierre|Una prueba falla: mantener cambio en validación.
Cerrar cambio satisfactorio|change|Todas las pruebas y tareas están completas|Cambio cierra con tiempos reales y resultado|Falta confirmación de servicio: bloquear cierre satisfactorio.
Revisar cambio fallido|manager|Despliegue requirió retorno por incompatibilidad|Revisión documenta causa, impacto y acciones correctivas|Acción sin responsable: impedir dar revisión por concluida.
Ejecutar cambio estándar preautorizado|agent|Alta de impresora sigue procedimiento aprobado vigente|Se aplica plantilla con evidencia de ejecución|El alcance supera la plantilla: tramitar cambio normal.
Tramitar cambio de emergencia|change|Parche urgente corrige caída activa de PACS|Autorización de emergencia y revisión posterior quedan registradas|Sin autoridad de emergencia: escalar sin omitir decisión.
Cancelar cambio programado|change|Proveedor retira versión antes de la ventana|Se cancelan tareas y se avisa a afectados|Tareas ya ejecutadas: registrar estado real y compensación requerida.
'''),
('automation', 'Automatizaciones e integraciones', 16, 2, 'Reglas, credenciales y destinos autorizados configurados.', 'automation-rule,event,webhook,outbox', '''
Crear regla de asignación automática|admin|Tickets de impresoras van a Infraestructura|Regla asigna solo tickets que cumplen todas sus condiciones|Condición desconocida: rechazar regla en lugar de ignorarla.
Validar regla antes de activarla|admin|Regla compara prioridad con valor crítico|Se muestran errores de operadores y campos incompatibles|Operador no soportado: impedir activación.
Simular regla con ejemplo|admin|Ticket de impresora normal y regla urgente|La simulación informa que no coincide sin efectos reales|Acción externa: nunca enviarla durante simulación.
Activar automatización programada|admin|Revisar tickets sin respuesta cada quince minutos|Se ejecutan acciones sobre candidatos autorizados|Dos workers coinciden: coordinar para evitar ejecución duplicada.
Desactivar regla problemática|admin|Regla nueva asigna equipo incorrecto|Nuevos eventos dejan de ejecutar la regla|Acción ya ejecutada: conservar historial para corrección manual.
Consultar historial de ejecuciones|auditor|Regla 8 actuó sobre ticket 42 a las 10:00|Se ve versión, evento, condiciones y resultado|Registro contiene secreto: presentar datos redactados.
Enviar webhook al cambiar estado|system|Ticket 42 pasa a resuelto|Destino autorizado recibe evento firmado con identificador estable|Destino caído: reintentar con límites y registro.
Recibir evento externo idempotente|system|Sistema de inventario repite evento EVT-77|Se aplica una sola actualización por identificador|Mismo identificador y otro contenido: registrar conflicto.
Reprocesar fallo de integración|admin|Webhook EVT-77 falló por credencial vencida|Tras corregir credencial se reintenta el evento original|Evento ya completado: impedir duplicar efectos.
Sincronizar miembros de directorio|admin|Grupo LDAP añade dos agentes y retira uno|Se concilian miembros según mapeo autorizado|Directorio incompleto: no desactivar masivamente por error transitorio.
Notificar resolución por correo|system|Ticket 42 se resuelve para Ana|Mensaje durable se entrega o conserva fallo visible|SMTP no disponible: conservar tarea pendiente.
Consultar cola de correo fallido|admin|Tres correos agotaron intentos de entrega|Se muestran causas y destinatarios autorizados con datos mínimos|Contenido sensible: no exponer cuerpo en logs generales.
Reintentar correo tras corregir configuración|admin|SMTP vuelve a estar disponible para correo 91|Se agenda un intento trazable del mensaje pendiente|Mensaje ya enviado: advertir antes de reenvío explícito.
Aplicar límite de solicitudes API|admin|Integración supera su cuota por minuto|Se limita consumo sin afectar a otras organizaciones|Cambio de réplica: conservar límite compartido efectivo.
Revocar acceso de conector|admin|Inventario externo deja de estar autorizado|Peticiones posteriores del conector fallan y quedan auditadas|Trabajo en curso: detener nuevas acciones de la credencial revocada.
Consultar salud de integraciones|admin|LDAP responde y SMTP lleva diez minutos fallando|Panel diferencia conexión, último éxito y error actual|Nunca ejecutado: mostrar sin evidencia en vez de saludable.
'''),
('reporting', 'Informes y auditoría', 16, 2, 'Datos de prueba o producción disponibles dentro del alcance del actor.', 'report,audit-event,export', '''
Consultar volumen de tickets|manager|Octubre agrupado por día de Caracas|Totales coinciden con tickets visibles y fechas del periodo|Intervalo invertido: rechazar consulta.
Comparar estados de cola|manager|Tickets abiertos, en progreso y resueltos actuales|Gráfico y listado usan la misma definición de estado|Estado nuevo desconocido: mostrar categoría explícita.
Medir tiempo de primera respuesta|manager|Tickets respondidos durante octubre|Indicador usa primera respuesta pública real|Tickets sin respuesta: excluir del promedio y contar pendientes.
Medir tiempo de resolución|manager|Incidentes resueltos durante la última semana|Se calcula duración desde marcas válidas|Fechas heredadas ausentes: presentar registros no calculables.
Consultar carga por agente|manager|Pedro tiene ocho tickets y Marta tiene tres|Se ve distribución de trabajo dentro del equipo autorizado|Agente desactivado con tickets: incluirlo como pendiente de reasignación.
Identificar tickets envejecidos|manager|Tickets abiertos por más de diez días|Lista ordena antigüedad con responsable y siguiente acción|Ticket pausado: mostrar pausa sin ocultar antigüedad total.
Analizar categorías más frecuentes|manager|Impresión concentra 40 incidentes en octubre|Ranking enlaza a tickets de cada categoría|Categoría eliminada: conservar etiqueta histórica identificable.
Exportar resultados a CSV|manager|Vista filtrada de 42 tickets de Radiología|Archivo respeta filtros y neutraliza fórmulas de celdas|Título empieza con igual: exportar como texto seguro.
Exportar resumen a PDF|manager|Resumen mensual para comité de servicio|Documento incluye periodo, filtros y fecha de generación|Resultado demasiado grande: ofrecer exportación acotada.
Consultar cambios de permisos|auditor|Pedro pasó de agente a solicitante ayer|Auditoría muestra actor, valor anterior y nuevo|Actor técnico: identificar integración responsable.
Investigar acceso a recurso sensible|auditor|Se revisan accesos autorizados al ticket restringido 88|Eventos permiten atribuir acceso y resultado|Log no recoge lectura histórica: declarar ausencia de evidencia.
Auditar transiciones de tickets|auditor|Ticket 42 fue reabierto y reasignado|Se reconstruye secuencia sin alterar registros|Marcas coincidentes: conservar orden estable por identificador.
Revisar cumplimiento de catálogo|manager|Solicitudes de portátil entregadas durante octubre|Se compara plazo prometido y entrega real|Solicitud cancelada: mostrarla separada del cumplimiento.
Analizar uso de autoservicio|manager|Guías consultadas antes de abrir tickets de impresión|Se presenta interacción medible sin atribuir resolución supuesta|Usuario anónimo: agregar sin inventar identidad.
Detectar calidad insuficiente de datos|manager|Tickets sin categoría y activos sin sede|Informe enumera registros que requieren corrección|Dato opcional: no contarlo como error obligatorio.
Revisar historial de exportaciones|auditor|Marta exportó tickets de Radiología el lunes|Se registra quién exportó, alcance y resultado|Exportación fallida: distinguir intento de archivo entregado.
'''),
])

def write_generated(product, name, groups):
    generated = []
    for area, label, count, phase, prerequisites, entities, operations, samples in groups:
        combos = [(operation, sample) for operation in operations for sample in samples]
        assert len(combos) >= count, (product, area, len(combos), count)
        rows = []
        for (title, actor, action, outcome, exception), (resource, detail) in combos[:count]:
            rows.append('|'.join((
                f'{title}: {resource}', actor,
                f'{action} para {resource}: {detail}.',
                outcome.format(resource=resource),
                exception.format(resource=resource),
            )))
        generated.append((area, label, count, phase, prerequisites, entities, '\n'.join(rows)))
    write(product, name, generated)

CLIENTS_MSP = [('Clínica Aurora', 'dominio aurora.example, Caracas'), ('Hotel Río', 'dominio hotelrio.example, Valencia'), ('Colegio Sol', 'dominio colegiosol.example, Maracay'), ('Logística Norte', 'dominio lognorte.example, Barquisimeto')]

write_generated('MSP', 'Trickreport · Servicios MSP', [
('clients', 'Clientes, contactos y permisos', 18, 4, 'Proveedor Andes IT registrado; dominios y contactos verificados.', 'client,site,contact,membership,role', [
('Crear espacio independiente', 'admin', 'Registrar un espacio con su dominio y administrador', 'El espacio de {resource} queda separado de los demás clientes.', 'Dominio ya asociado a otro espacio: bloquear el alta y conservar la asociación actual.'),
('Añadir sede de cliente', 'admin', 'Vincular una sede y su zona horaria', 'La sede de {resource} se incorpora al ámbito correcto del cliente.', 'Código de sede duplicado: solicitar corrección sin crear una sede parcial.'),
('Invitar contacto autorizado', 'contact', 'Invitar una persona con permisos mínimos', 'El contacto de {resource} sólo accede a recursos permitidos.', 'Correo de otro cliente: impedir que herede acceso entre espacios.'),
('Cambiar el estado del cliente', 'manager', 'Marcar el estado operativo y su motivo', 'El estado de {resource} se actualiza con autoría y fecha.', 'Contrato suspendido: limitar cambios a un administrador autorizado.'),
('Exportar inventario de cliente', 'auditor', 'Solicitar una exportación dentro del alcance asignado', 'La exportación de {resource} contiene únicamente activos autorizados.', 'Identificador de otro cliente: denegar sin revelar si existe.'),
], CLIENTS_MSP),
('service-desk', 'Mesa de servicio multicliente', 30, 4, 'Colas segregadas por cliente y reglas de asignación configuradas.', 'ticket,queue,client,priority,assignment', [
('Registrar incidente del cliente', 'tech', 'Clasificar síntoma, impacto y cliente afectado', 'El incidente de {resource} se abre en la cola de su cliente.', 'Cliente desconocido o inactivo: no crear una solicitud huérfana.'),
('Asignar ticket a especialista', 'manager', 'Elegir técnico habilitado para el servicio requerido', 'La asignación de {resource} queda visible para la cola autorizada.', 'Técnico sin acceso a ese cliente: rechazar la asignación.'),
('Escalar incumplimiento de SLA', 'tech', 'Solicitar atención de segundo nivel antes del vencimiento', 'La escalación de {resource} registra regla, responsable y plazo.', 'SLA pausado según contrato: calcular con la pausa explícita.'),
('Coordinar incidente mayor', 'manager', 'Abrir sala de coordinación y actualización al contacto', 'El evento de {resource} reúne tareas sin mezclar clientes.', 'Impacto no confirmado: no notificar como incidente general.'),
('Cerrar solicitud con solución', 'tech', 'Documentar causa, acción y evidencia de resolución', 'El cierre de {resource} conserva el historial y notifica al contacto autorizado.', 'Trabajo incompleto: impedir cierre y conservar el estado actual.'),
('Reasignar ticket entre técnicos', 'manager', 'Mover el ticket a otro técnico con acceso al mismo cliente', 'La reasignación de {resource} conserva historial y responsable anterior.', 'Técnico destino sin acceso al cliente: rechazar el movimiento.'),
('Añadir nota interna de servicio', 'tech', 'Registrar diagnóstico visible sólo para el proveedor', 'La nota de {resource} no se expone al contacto del cliente.', 'Contacto consulta el ticket: filtrar la nota interna.'),
('Consultar cola consolidada por cliente', 'manager', 'Revisar tickets abiertos, vencidos y sin asignar del cliente', 'La cola de {resource} muestra prioridad y SLA de cada ticket.', 'Cliente fuera del alcance del gestor: excluirlo del resultado.'),
], [('caída de VPN', 'Cliente Aurora, sede Caracas'), ('correo bloqueado', 'Hotel Río, equipo de recepción'), ('respaldo retrasado', 'Colegio Sol, servidor académico'), ('latencia WAN', 'Logística Norte, centro de despacho')]),
('contracts', 'Contratos y SLA', 22, 4, 'Servicios publicables y reglas comerciales aprobadas.', 'contract,service,sla,calendar,renewal', [
('Crear contrato administrado', 'finance', 'Registrar vigencia, moneda y capacidad mensual', 'El contrato de {resource} establece límites y fechas verificables.', 'Vigencia final anterior a la inicial: rechazar el contrato.'),
('Asignar paquete de servicio', 'finance', 'Asociar cobertura, horario y canales incluidos', 'El paquete de {resource} refleja los servicios contratados.', 'Servicio fuera del catálogo vigente: impedir su activación.'),
('Definir objetivo de respuesta', 'manager', 'Configurar objetivos por prioridad y horario hábil', 'El SLA de {resource} calcula plazos con el calendario del cliente.', 'Objetivo nulo o negativo: no guardar la política.'),
('Renovar contrato con cambios', 'finance', 'Preparar una nueva versión para el próximo periodo', 'La renovación de {resource} mantiene el historial y aplica desde su fecha.', 'Dos versiones activas en el mismo periodo: exigir resolución del conflicto.'),
('Suspender servicio por vencimiento', 'finance', 'Registrar autorización y fecha efectiva de suspensión', 'La cobertura de {resource} cambia según la política publicada.', 'Solicitud sin autorización contractual: mantener el servicio activo.'),
('Registrar bonificación o penalización de SLA', 'finance', 'Calcular el ajuste del periodo a partir de los incumplimientos registrados', 'El ajuste de {resource} queda trazado a los tickets que lo originan.', 'Incumplimiento en disputa: excluirlo del cálculo hasta resolver la revisión.'),
], [('Contrato Aurora Gold', 'USD 4 000 al mes, 24x7'), ('Contrato Hotel Río', '80 horas al mes, horario extendido'), ('Contrato Colegio Sol', '40 endpoints, horario escolar'), ('Contrato Norte', 'soporte de red, lunes a sábado')]),
('technicians', 'Técnicos y planificación', 16, 4, 'Técnicos, habilidades, calendarios y zonas horarias verificados.', 'technician,shift,skill,capacity,schedule', [
('Crear turno de guardia', 'manager', 'Definir horario, zona y nivel de escalación', 'La guardia de {resource} queda asignada con cobertura visible.', 'Dos turnos incompatibles: pedir resolución antes de publicar.'),
('Asignar técnico por habilidad', 'manager', 'Filtrar personal por certificación y disponibilidad', 'La cobertura de {resource} corresponde a un técnico elegible.', 'Certificación vencida: excluir al técnico de la asignación.'),
('Registrar disponibilidad', 'tech', 'Indicar ausencia o cambio de disponibilidad', 'La agenda de {resource} actualiza la capacidad de planificación.', 'Periodo superpuesto con guardia confirmada: alertar y mantener el turno.'),
('Transferir guardia con contexto', 'tech', 'Entregar pendientes, riesgos y contactos de emergencia', 'El relevo de {resource} preserva responsables y decisiones.', 'Receptor sin permiso de cliente: redactar el contexto restringido.'),
], [('guardia de redes', 'zona America/Caracas, 18:00–06:00'), ('equipo nube', 'dos técnicos certificados AZ-104'), ('turno fin de semana', 'sábado 08:00–20:00'), ('equipo de campo', 'dos sedes y cuatro visitas pendientes')]),
('time-expenses', 'Horas, gastos y costes', 20, 4, 'Contratos con tarifas y reglas de aprobación de costes definidos.', 'timesheet,expense,cost,approval', [
('Registrar horas trabajadas', 'tech', 'Asociar duración, tarea y evidencia del trabajo', 'Las horas de {resource} quedan ligadas al técnico y contrato correctos.', 'Periodo cerrado: enviar a aprobación excepcional.'),
('Presentar gasto reembolsable', 'tech', 'Adjuntar recibo y centro de coste', 'El gasto de {resource} entra al circuito de aprobación asignado.', 'Moneda o recibo inválidos: rechazar el gasto sin duplicarlo.'),
('Aprobar partes de tiempo', 'approver', 'Revisar detalle y comparar con el alcance contratado', 'La aprobación de {resource} conserva autor y comentario.', 'Aprobador es quien registró las horas: pedir aprobación independiente.'),
('Imputar coste a un contrato', 'finance', 'Asignar el coste aprobado al contrato y periodo que lo originó', 'El coste de {resource} se refleja en la rentabilidad del contrato.', 'Contrato cerrado o de otro cliente: impedir la imputación.'),
('Revisar desvío de costes', 'finance', 'Comparar horas, gastos y avance con el límite acordado', 'El desvío de {resource} se informa antes de exceder el límite.', 'Datos importados sin conciliar: mostrar advertencia y no facturar.'),
], [('migración M365', '120 horas presupuestadas, fase de descubrimiento'), ('renovación de red', '40 switches en seis sedes'), ('puesta en marcha del SOC', 'hitos mensuales y guardias incluidas'), ('inventario de portátiles', '75 equipos, Caracas y Valencia')]),
('projects', 'Proyectos de servicio', 14, 4, 'Proyectos, presupuestos y reglas de aprobación definidos.', 'project,milestone,task,budget', [
('Crear proyecto de cliente', 'manager', 'Definir alcance, responsable, presupuesto y fechas', 'El proyecto de {resource} tiene control de alcance y presupuesto.', 'Presupuesto sin aprobador: impedir su publicación.'),
('Definir hitos del proyecto', 'manager', 'Establecer entregables, fechas y criterios de aceptación por hito', 'Los hitos de {resource} quedan ordenados y visibles para el equipo.', 'Hito sin responsable o fecha: mantener el proyecto en borrador.'),
('Asignar tareas del proyecto', 'manager', 'Repartir tareas entre técnicos con acceso al cliente', 'Las tareas de {resource} tienen responsable y fecha compromiso.', 'Técnico sin acceso al cliente del proyecto: rechazar la asignación.'),
('Cerrar proyecto con informe', 'manager', 'Recopilar hitos, horas, costes y lecciones del proyecto', 'El cierre de {resource} conserva el informe y el historial completo.', 'Hitos abiertos o costes sin aprobar: impedir el cierre definitivo.'),
], [('migración M365', '120 horas presupuestadas, fase de descubrimiento'), ('renovación de red', '40 switches en seis sedes'), ('puesta en marcha del SOC', 'hitos mensuales y guardias incluidas'), ('inventario de portátiles', '75 equipos, Caracas y Valencia')]),
('profitability', 'Rentabilidad y preparación de facturación', 16, 4, 'Horas, gastos, contratos y costos conciliados.', 'revenue,cost,margin,invoice-draft', [
('Calcular margen por cliente', 'finance', 'Comparar ingresos del contrato con costos aprobados', 'El margen de {resource} usa datos conciliados y periodo declarado.', 'Costos incompletos: rotular el margen como provisional.'),
('Comparar rentabilidad de servicios', 'finance', 'Agrupar ingresos y costos por paquete de servicio', 'La comparación de {resource} separa contratos y moneda.', 'Mezcla de monedas: exigir conversión con tipo de cambio registrado.'),
('Preparar borrador de facturación', 'finance', 'Reunir horas aprobadas, gastos y ajustes del periodo', 'El borrador de {resource} lista conceptos trazables a tickets y contratos.', 'Horas no aprobadas: excluirlas del borrador facturable.'),
('Revisar desviación de horas', 'manager', 'Comparar horas estimadas, aprobadas y consumidas', 'La desviación de {resource} identifica contratos que requieren revisión.', 'Horas no aprobadas: excluirlas del total facturable.'),
], [('octubre Aurora', 'ingreso USD 4 000, costos por conciliar'), ('servicio de guardia', 'tres clientes, periodo trimestral'), ('proyecto de migración', '120 h estimadas frente a 88 aprobadas'), ('renovaciones de diciembre', 'monedas USD y EUR en reportes separados')]),
('inventory-discovery', 'Inventario y descubrimiento integrado', 20, 4, 'Ubicaciones, categorías, proveedores y conectores de descubrimiento autorizados.', 'asset,model,warranty,location,discovery', [
('Dar de alta un activo', 'asset', 'Capturar serial, modelo, cliente y ubicación', 'El activo de {resource} se identifica de forma única y trazable.', 'Número de serie ya registrado: evitar inventario duplicado.'),
('Vincular activo a un usuario', 'asset', 'Registrar custodio y fecha de entrega', 'La custodia de {resource} refleja la persona y sede autorizadas.', 'Usuario de otro cliente: bloquear el vínculo.'),
('Actualizar garantía', 'asset', 'Adjuntar proveedor, cobertura y vencimiento', 'La garantía de {resource} se consulta con alertas de vencimiento.', 'Documento no verificable: dejar la actualización en revisión.'),
('Conciliar inventario físico', 'auditor', 'Comparar escaneo con inventario esperado de la sede', 'La conciliación de {resource} muestra diferencias para revisión.', 'Escaneo parcial: no marcar como verificada la ubicación completa.'),
('Importar descubrimiento aprobado', 'system', 'Conciliar dispositivos detectados con el cliente y segmento', 'Los hallazgos de {resource} se crean dentro del ámbito autorizado.', 'Segmento fuera de inventario: poner hallazgos en cuarentena.'),
], [('laptop TR-104', 'Dell Latitude, Clínica Aurora, Caracas'), ('switch SW-22', 'Cisco Catalyst, Hotel Río, MDF 2'), ('servidor SR-08', 'Linux, Colegio Sol, sala técnica'), ('tablet TB-31', 'Zebra, Logística Norte, almacén 4')]),
('monitoring-remote', 'Monitorización y acceso remoto integrado', 28, 4, 'Conectores Zabbix y MeshCentral autorizados; umbrales, consentimientos y ventanas revisados.', 'monitor,event,remote-session,device,consent,audit', [
('Convertir alerta en incidente', 'system', 'Aplicar umbral, ventana y regla de correlación', 'La alerta de {resource} genera un solo incidente asignable.', 'Alerta duplicada del mismo evento: actualizar el incidente existente.'),
('Silenciar mantenimiento autorizado', 'manager', 'Definir servicio, motivo y ventana temporal', 'El silencio de {resource} vence automáticamente al terminar la ventana.', 'Ventana sin aprobador: mantener la alerta activa.'),
('Correlacionar alertas relacionadas', 'tech', 'Agrupar eventos por servicio y causa probable', 'La correlación de {resource} reduce duplicados sin perder señales.', 'Causa no confirmada: mostrar relación sugerida, no consolidar automáticamente.'),
('Cerrar alerta tras recuperación', 'system', 'Comprobar recuperación y vincular evidencia de monitorización', 'La señal de {resource} cierra o actualiza el incidente asociado.', 'Sin lectura reciente: no inferir que el servicio se recuperó.'),
('Solicitar sesión remota', 'tech', 'Elegir dispositivo y motivo de intervención', 'La sesión de {resource} requiere consentimiento del cliente.', 'Dispositivo no inventariado: impedir conexión remota.'),
('Iniciar sesión registrada', 'tech', 'Conectar mediante el conector administrado y autorizado', 'El acceso a {resource} deja registro de inicio y técnico.', 'Credencial de integración inválida: fallar cerrado y alertar al administrador.'),
('Revocar sesión en curso', 'manager', 'Detener sesión por solicitud del cliente', 'La conexión a {resource} se cierra y queda motivo registrado.', 'Estado remoto desconocido: no reportar desconexión confirmada.'),
], [('servidor SR-08', 'Colegio Sol, sala técnica'), ('portátil TR-104', 'Clínica Aurora, usuaria de Radiología'), ('router SW-22', 'Hotel Río, MDF 2'), ('terminal TB-31', 'Logística Norte, almacén 4')]),
('reporting', 'Informes y auditoría', 16, 4, 'Datos operativos disponibles dentro del alcance del actor.', 'report,audit-event,export', [
('Consultar volumen de tickets por cliente', 'manager', 'Agrupar tickets por cliente, estado y periodo', 'El volumen de {resource} coincide con los tickets visibles del periodo.', 'Cliente fuera del alcance: excluirlo sin revelar su actividad.'),
('Comparar cumplimiento de SLA por contrato', 'manager', 'Calcular respuestas y resoluciones dentro de objetivo por contrato', 'El cumplimiento de {resource} separa pausas autorizadas de vencimientos.', 'Contrato sin política vigente: mostrarlo como sin cobertura, no como 100%.'),
('Auditar accesos a datos de clientes', 'auditor', 'Revisar lecturas y cambios con actor, cliente y alcance', 'La auditoría de {resource} permite atribuir cada acceso registrado.', 'Evento sin actor identificable: marcarlo como evidencia incompleta.'),
('Exportar informe operativo', 'manager', 'Descargar el informe con filtros y fecha de generación', 'La exportación de {resource} respeta filtros y permisos del solicitante.', 'Texto con fórmula de hoja de cálculo: exportarlo como texto seguro.'),
], [('octubre multicliente', 'tickets y SLA de los cuatro clientes activos'), ('guardias del trimestre', 'turnos, escalaciones y tiempos de respuesta'), ('activos por sede', 'inventario y garantías de Caracas y Valencia'), ('costes por contrato', 'horas y gastos aprobados del periodo')]),
])

write_generated('CX', 'Trickreport · Atención al cliente', [
('customer-identity', 'Clientes, contactos y permisos', 20, 5, 'Organización Tienda Faro y políticas de privacidad publicadas.', 'customer,profile,consent,preference', [
('Crear perfil de cliente', 'customer', 'Registrar datos de contacto mínimos y canal preferido', 'El perfil de {resource} conserva sólo datos necesarios para atender.', 'Correo ya registrado: ofrecer recuperación sin crear un duplicado.'),
('Verificar identidad del solicitante', 'agent', 'Confirmar el factor permitido antes de revelar datos', 'La identidad de {resource} se valida antes de consultar su caso.', 'Verificación fallida: no confirmar si existe una compra o ticket.'),
('Actualizar preferencias de contacto', 'customer', 'Elegir canal, idioma y horario de comunicación', 'Las preferencias de {resource} se aplican a mensajes futuros.', 'Canal no disponible: mantener preferencia anterior y pedir elección.'),
('Registrar consentimiento', 'admin', 'Guardar texto, versión, propósito y fecha de aceptación', 'El consentimiento de {resource} queda ligado a la versión mostrada.', 'Propósito obligatorio ausente: no iniciar comunicaciones comerciales.'),
('Solicitar eliminación de datos', 'customer', 'Abrir solicitud de privacidad con verificación de identidad', 'La solicitud de {resource} entra a revisión y deja trazabilidad.', 'Obligación de conservación activa: informar alcance y retención aplicable.'),
], [('Ana Pérez', 'correo ana@example.test, español'), ('Luis Díaz', 'teléfono terminado en 4421, WhatsApp opt-in'), ('Marta Gil', 'compradora recurrente, idioma inglés'), ('Pablo León', 'cuenta compartida que requiere verificación')]),
('inbox', 'Bandeja y gestión de conversaciones', 32, 5, 'Estados de conversación y visibilidad pública/interna configurados.', 'conversation,message,note,status,assignment', [
('Responder consulta del cliente', 'agent', 'Redactar explicación pública con siguiente paso claro', 'La respuesta de {resource} se entrega y queda fechada.', 'Nota marcada como interna: impedir su envío al cliente.'),
('Añadir nota privada al caso', 'agent', 'Registrar diagnóstico interno sin incluirlo en mensaje público', 'La nota de {resource} sólo se muestra a personal autorizado.', 'Solicitante consulta la conversación: filtrar la nota interna.'),
('Solicitar información faltante', 'agent', 'Pedir el dato necesario y pausar el plazo según política', 'La solicitud de {resource} deja explícito qué dato falta.', 'Caso cerrado: validar reapertura antes de enviar el pedido.'),
('Resolver caso con causa', 'agent', 'Documentar causa, acción y respuesta ofrecida', 'La resolución de {resource} conserva el historial previo.', 'Solución vacía: no permitir cerrar el caso.'),
('Reabrir por respuesta nueva', 'customer', 'Responder al mismo hilo dentro de la ventana permitida', 'El caso de {resource} vuelve al flujo definido sin duplicarse.', 'Ventana de reapertura vencida: crear seguimiento relacionado.'),
('Responder por canal de origen', 'agent', 'Seleccionar respuesta pública y canal permitido', 'La respuesta a {resource} mantiene el hilo y preferencia del cliente.', 'Canal no verificado: no enviar y conservar el borrador.'),
('Fusionar conversaciones duplicadas', 'supervisor', 'Elegir la conversación principal y conservar ambas líneas temporales', 'La fusión de {resource} deja una sola conversación activa con historial completo.', 'Conversaciones de clientes distintos: impedir la fusión.'),
('Marcar conversación como spam', 'agent', 'Clasificar el mensaje y detener respuestas automáticas', 'La conversación de {resource} sale de la bandeja activa conservando auditoría.', 'Cliente legítimo marcado por error: permitir restaurarla desde la bandeja.'),
], [('consulta sobre ORD-1042', 'se explica ventana estimada de entrega'), ('devolución incompleta', 'nota de almacén restringida al personal'), ('garantía solicitada', 'falta foto del número de serie'), ('reembolso ORD-882', 'cliente informa que el abono no aparece')]),
('email', 'Correo electrónico', 22, 5, 'Buzón IMAP/SMTP verificado; remitentes, firmas y retención aprobados.', 'email,message,thread,attachment,outbox', [
('Recibir correo de soporte', 'system', 'Validar remitente y asociar mensaje a la conversación', 'El correo de {resource} aparece una sola vez en la línea temporal.', 'Mensaje duplicado o con firma inválida: no abrir otro caso.'),
('Crear caso desde correo entrante', 'system', 'Extraer asunto, remitente y adjuntos permitidos del mensaje', 'El caso de {resource} nace vinculado al hilo de correo original.', 'Remitente no verificable: crear consulta sin datos de cuenta asociados.'),
('Responder correo conservando el hilo', 'agent', 'Redactar la respuesta con referencias del mensaje original', 'La respuesta de {resource} llega al mismo hilo del cliente.', 'Dirección del cliente inválida: conservar el borrador y avisar al agente.'),
('Detectar rebote y actualizar contacto', 'system', 'Interpretar el aviso de rebote y marcar la dirección afectada', 'El rebote de {resource} actualiza el contacto sin reenviar en bucle.', 'Rebote temporal: programar revisión sin marcar la dirección como muerta.'),
('Adjuntar archivo permitido al correo', 'agent', 'Subir comprobante dentro del límite y tipo aceptado', 'La evidencia de {resource} queda vinculada con acceso controlado.', 'Archivo sobredimensionado o tipo falso: rechazar sin publicarlo.'),
('Reintentar envío fallido', 'system', 'Reenviar el mensaje pendiente cuando SMTP vuelve a responder', 'El envío de {resource} se completa una sola vez o queda fallo visible.', 'Mensaje ya entregado: no duplicar el correo al cliente.'),
], [('correo de Ana', 'asunto pedido ORD-1042, hilo con Message-ID único'), ('correo de Luis', 'consulta de horario, responde al hilo de bienvenida'), ('correo rebotado', 'dirección obsoleta, no reenviar automáticamente'), ('correo de Marta', 'factura adjunta del pedido ORD-2030')]),
('portal-chat', 'Portal web y chat en vivo', 26, 5, 'Formularios, categorías, campos mínimos y reglas de acceso definidos.', 'case,form,category,attachment,chat', [
('Abrir consulta desde portal', 'customer', 'Describir el motivo y adjuntar evidencia opcional', 'El caso de {resource} recibe confirmación y número de seguimiento.', 'Campos inválidos: conservar el borrador y señalar el dato necesario.'),
('Enviar reporte en nombre del cliente', 'agent', 'Confirmar identidad y registrar canal de origen', 'El caso de {resource} identifica al cliente y al agente registrador.', 'Identidad no verificada: guardar como consulta interna sin datos de cuenta.'),
('Clasificar tipo de consulta', 'supervisor', 'Elegir categoría que corresponde al asunto recibido', 'La clasificación de {resource} dirige el caso a la cola adecuada.', 'Categoría desactivada: pedir una alternativa vigente.'),
('Guardar caso incompleto', 'customer', 'Continuar más tarde con los campos ya capturados', 'El borrador de {resource} queda privado a su propietario.', 'Sesión expirada: no guardar información en un perfil distinto.'),
('Adjuntar archivo permitido', 'customer', 'Subir comprobante dentro del límite y tipo aceptado', 'La evidencia de {resource} queda vinculada con acceso controlado.', 'Archivo sobredimensionado o tipo falso: rechazar sin publicarlo.'),
('Iniciar chat con agente', 'customer', 'Solicitar ayuda y aceptar el aviso de privacidad', 'La conversación de {resource} se asigna sin compartir datos entre clientes.', 'Sin agentes disponibles: informar espera y ofrecer alternativa.'),
('Continuar conversación del portal al chat', 'customer', 'Abrir el chat desde su caso conservando el contexto autorizado', 'La conversación de {resource} continúa en chat sin crear un caso nuevo.', 'Identidad no verificada en el chat: pedir verificación antes de mostrar el historial.'),
], [('pedido ORD-1042', 'consulta de entrega, factura PDF opcional'), ('devolución ORD-882', 'producto recibido con embalaje dañado'), ('acceso a la cuenta', 'error al cambiar correo, sin incluir contraseña'), ('garantía de cafetera', 'número de serie y foto del defecto')]),
('whatsapp', 'WhatsApp', 24, 5, 'Número de WhatsApp Business verificado; plantillas aprobadas por Meta y opt-in registrado.', 'whatsapp,conversation,template,opt-in', [
('Iniciar conversación de WhatsApp', 'customer', 'Escribir al número oficial y aceptar el aviso de privacidad', 'La conversación de {resource} se crea con opt-in y canal registrados.', 'Número sin opt-in: no enviar mensajes proactivos.'),
('Enviar plantilla aprobada', 'agent', 'Elegir plantilla vigente y completar sus variables verificadas', 'La plantilla de {resource} se envía con el contenido aprobado por Meta.', 'Plantilla pausada o rechazada: impedir el envío y ofrecer alternativa.'),
('Recibir mensaje de WhatsApp', 'system', 'Verificar firma del webhook y asociar el número a la conversación', 'El mensaje de {resource} entra una sola vez en la bandeja.', 'Evento repetido del webhook: procesarlo una única vez por identificador.'),
('Vincular conversación de WhatsApp al portal', 'agent', 'Verificar la identidad y relacionar el chat con el caso del portal', 'La conversación de {resource} comparte historial autorizado entre canales.', 'Identidad no verificada: no mostrar datos del portal en el chat.'),
('Respetar la ventana de 24 horas', 'system', 'Comprobar la última interacción del cliente antes de responder en texto libre', 'La respuesta de {resource} usa texto libre sólo dentro de la ventana permitida.', 'Ventana vencida: exigir plantilla aprobada en lugar de texto libre.'),
('Gestionar baja de WhatsApp', 'customer', 'Enviar la palabra de baja o revocar el canal desde el portal', 'La baja de {resource} detiene mensajes de WhatsApp conservando otros canales.', 'Baja ambigua: confirmar la interpretación antes de aplicarla.'),
], [('mensaje de Marta', 'WhatsApp opt-in sobre estado de envío'), ('mensaje de Luis', 'teléfono terminado en 4421, primera consulta'), ('mensaje de Ana', 'continúa su caso del portal por WhatsApp'), ('mensaje de Pablo', 'solicita baja del canal')]),
('routing', 'Asignación y automatización', 20, 5, 'Horarios de atención, idiomas, competencias y cargas definidos.', 'queue,agent,skill,routing-rule', [
('Asignar cola por categoría', 'system', 'Aplicar categoría, idioma y horario de servicio', 'El caso de {resource} llega a una cola habilitada.', 'Regla sin destino válido: enviar a triage y alertar al supervisor.'),
('Asignar agente por competencia', 'supervisor', 'Comparar producto, idioma y disponibilidad del agente', 'La asignación de {resource} cumple competencias mínimas.', 'Agente fuera de turno: dejar pendiente de asignación.'),
('Reequilibrar cola saturada', 'supervisor', 'Mover casos según prioridad y capacidad disponible', 'La carga de {resource} se reparte respetando prioridad y propiedad.', 'Caso en atención activa: solicitar transferencia acordada.'),
('Priorizar cliente con impacto alto', 'system', 'Aplicar política publicada a la prioridad del caso', 'La prioridad de {resource} refleja impacto y urgencia verificables.', 'Atributo de cliente no verificado: no elevar privilegios automáticamente.'),
('Asignar caso al volver de espera', 'system', 'Reactivar caso tras respuesta entrante del cliente', 'La respuesta de {resource} vuelve a la cola apropiada.', 'Caso cerrado por política: mostrar estado y ofrecer reapertura controlada.'),
], [('devolución de producto', 'cola Garantías, español, horario hábil'), ('facturación empresarial', 'cola Cuentas, idioma inglés'), ('chat de entrega', 'equipo Envíos con dos agentes disponibles'), ('reclamo de seguridad', 'supervisor especializado, atención urgente')]),
('knowledge', 'Conocimiento y autoservicio', 14, 5, 'Artículos aprobados, búsqueda e idioma de contenido identificados.', 'article,search,feedback,deflection', [
('Buscar solución antes de crear caso', 'customer', 'Consultar frases relevantes en el idioma elegido', 'Los resultados de {resource} muestran artículos publicados.', 'Artículo borrador o restringido: no incluirlo en resultados públicos.'),
('Publicar artículo de ayuda', 'supervisor', 'Revisar exactitud, enlaces y público destinatario', 'La guía de {resource} queda disponible tras aprobación.', 'Instrucciones inseguras o desactualizadas: devolver a revisión.'),
('Valorar utilidad de artículo', 'customer', 'Indicar si la guía resolvió la consulta', 'La valoración de {resource} se registra sin exponer datos personales.', 'Voto repetido: actualizar respuesta anterior sin inflar conteos.'),
('Vincular artículo a caso', 'agent', 'Adjuntar una guía pertinente como respuesta sugerida', 'La recomendación de {resource} mantiene versión y enlace visibles.', 'Guía retirada: no enviar un enlace roto al cliente.'),
], [('restablecer contraseña', 'guía vigente en español e inglés'), ('devolver una compra', 'artículo aprobado con pasos y condiciones'), ('consultar garantía', 'contenido publicado para pequeños electrodomésticos'), ('seguimiento de envío', 'guía con enlace al portal de pedidos')]),
('sla', 'SLA y escalaciones', 16, 5, 'Calendario, horario de cobertura, feriados y niveles de prioridad publicados.', 'sla,calendar,escalation,case', [
('Iniciar reloj de primera respuesta', 'system', 'Usar hora de recepción y calendario efectivo', 'El reloj de {resource} empieza con el primer evento confiable.', 'Zona horaria ausente: usar calendario configurado y registrar la elección.'),
('Pausar SLA esperando al cliente', 'agent', 'Marcar la pregunta pública que requiere respuesta', 'El plazo de {resource} pausa según la regla aprobada.', 'Nota interna sola: no pausar el SLA del cliente.'),
('Escalar caso próximo a vencer', 'system', 'Notificar supervisor antes del umbral configurado', 'La alerta de {resource} se emite una sola vez por umbral.', 'Canal de aviso caído: registrar y usar destino alternativo autorizado.'),
('Reportar incumplimiento de SLA', 'auditor', 'Consultar casos del periodo y motivo de retraso', 'El informe de {resource} separa pausas y vencimientos reales.', 'Datos de periodo incompletos: mostrar conteo pendiente de conciliación.'),
], [('respuesta de ORD-1042', 'prioridad normal, calendario Caracas'), ('reclamo de seguridad', 'prioridad crítica, escalación inmediata'), ('consulta recibida sábado', 'el horario inicia lunes 08:00'), ('caso en espera', 'pausado desde solicitud pública de evidencia')]),
('surveys', 'Encuestas y satisfacción', 12, 5, 'Encuestas, métricas y políticas de muestreo aprobadas.', 'survey,response,quality', [
('Enviar encuesta tras resolución', 'system', 'Esperar cierre confirmado y respetar frecuencia máxima', 'La encuesta de {resource} se envía una vez al contacto autorizado.', 'Cliente sin consentimiento o con envío reciente: omitir encuesta.'),
('Registrar satisfacción del cliente', 'customer', 'Guardar puntuación y comentario opcional', 'La valoración de {resource} se asocia al caso sin alterar su cierre.', 'Comentario con secreto o datos excesivos: ocultar y solicitar ajuste.'),
('Revisar calidad de respuesta', 'supervisor', 'Muestrear respuestas según criterios publicados', 'La revisión de {resource} registra criterios y acción de mejora.', 'Muestra sesgada o pequeña: rotular el resultado con su límite.'),
], [('encuesta ORD-1042', 'cliente confirmó entrega y dejó 5/5'), ('queja por devolución', 'baja satisfacción, revisión de supervisor'), ('chat de compra', 'respuesta inicial en horario nocturno'), ('muestra de octubre', 'respuestas agrupadas por canal y semana')]),
('admin-reporting', 'Administración e informes', 14, 5, 'Acceso por rol, retención y conectores configurados y revisados.', 'privacy-request,api-key,audit,report,webhook', [
('Procesar solicitud de acceso a datos', 'admin', 'Verificar titular y preparar datos dentro del alcance legal', 'La solicitud de {resource} queda auditada y limitada al titular.', 'Identidad no verificada: no exportar datos de la cuenta.'),
('Configurar webhook firmado', 'admin', 'Guardar destino HTTPS y secreto de forma protegida', 'El conector de {resource} firma eventos y permite rotar claves.', 'URL privada o no verificable: impedir envío para evitar acceso interno.'),
('Eliminar credencial comprometida', 'admin', 'Revocar clave y comprobar trabajos pendientes', 'La clave de {resource} deja de autorizar nuevas llamadas.', 'Revocación afecta tareas críticas: detener nuevas acciones y alertar responsable.'),
('Consultar tiempos por canal', 'manager', 'Comparar recepción y respuesta dentro del mismo periodo', 'El informe de {resource} separa canales y zonas horarias.', 'Canal sin marcas fiables: excluirlo y mostrar el dato faltante.'),
], [('perfil de Ana', 'solicitud de exportación verificada por correo'), ('CRM externo', 'webhook HTTPS firmado de actualización de caso'), ('clave de chat', 'rotación por sospecha de exposición'), ('muestra de octubre', 'respuestas agrupadas por canal y semana')]),
])
(ROOT / 'README.md').write_text('''# Catálogo de casos de uso Trickreport

El catálogo define alcance futuro para los tres productos. Tener un caso documentado no significa que esté implementado ni comprobado. Los casos sólo pueden salir de `planned` cuando se enlacen código, contrato, pruebas y evidencia operativa.

| Producto | Casos | Datos | Lectura |
| --- | ---: | --- | --- |
| Soporte TI | 200 | [ti.json](ti.json) | [ti.md](ti.md) |
| Servicios MSP | 200 | [msp.json](msp.json) | [msp.md](msp.md) |
| Atención al cliente | 200 | [cx.json](cx.json) | [cx.md](cx.md) |

Los ejemplos usan organizaciones, clientes, personas y pedidos ficticios. Los permisos del catálogo describen validaciones requeridas en el servidor; el rol de interfaz nunca concede acceso por sí solo. Los contratos API están marcados como propuestos y se definirán al implementar cada caso.

Para regenerar los JSON y documentos Markdown desde la fuente, ejecutar `python .author.py` dentro de este directorio.
''', encoding='utf-8')

