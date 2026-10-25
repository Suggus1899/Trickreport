'use client';

import { useState } from 'react';
import Link from 'next/link';
import { api, json } from './api';
import type { CXAssignmentRule, CXAutoReply, CXConversation, CXCustomer, CXContact, CXEmailAccount, CXMessage, CXReport, CXSLAPolicy, CXSurvey, CXWhatsAppAccount, User } from './contracts';
import { conversationTransitions } from './contracts';
import { Button, Empty, Feedback, Field, Heading, LoadState, Pill, Select, Textarea, useLoad, useMutation, useText } from './ui';

// CX (Etapa 5): atención al cliente. Todas las pantallas consumen la
// API real (/api/v1); el servidor Go sigue siendo la autoridad de
// permisos y de aislamiento entre clientes externos (la ocultación de
// un botón aquí nunca sustituye esa comprobación). La bandeja muestra
// la conversación completa con su historial y su canal; el cliente
// externo solo ve sus propias conversaciones y nunca las notas
// internas, porque el servidor las filtra. El conocimiento no se
// duplica aquí: el autoservicio reutiliza el módulo existente en
// /articles (knowledge.tsx).

// ── Helpers puros (cubiertos por tests Vitest) ──────────────────────

export function isSLABreached(c: CXConversation, at = new Date()): boolean {
  if (c.status === 'resolved' || c.status === 'closed') return false;
  if (c.resolution_due_at && new Date(c.resolution_due_at) < at) return true;
  if (!c.first_response_at && c.response_due_at && new Date(c.response_due_at) < at) return true;
  return false;
}

export function avgSatisfaction(surveys: CXSurvey[]): number {
  if (!surveys.length) return 0;
  return surveys.reduce((sum, s) => sum + s.rating, 0) / surveys.length;
}

export function nextPriority(priority: string): string {
  if (priority === 'low') return 'medium';
  if (priority === 'medium') return 'high';
  return 'critical';
}

const channels = ['email', 'portal', 'chat', 'whatsapp'] as const;
const priorities = ['low', 'medium', 'high', 'critical'] as const;

// ── Clientes externos y contactos ───────────────────────────────────

export function CXCustomers({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const customers = useLoad(() => api<CXCustomer[]>('/cx/customers'));
  const [editing, setEditing] = useState<CXCustomer | null>();
  const [selected, setSelected] = useState<CXCustomer | null>(null);
  const contacts = useLoad(() => selected ? api<CXContact[]>(`/cx/customers/${selected.id}/contacts`) : Promise.resolve([] as CXContact[]), [selected?.id]);
  return <>
    <Heading title={t('Clientes externos', 'External customers')} description={t('Personas y empresas que contactan por portal, correo, chat o WhatsApp. Cada cliente solo ve sus propias conversaciones.', 'People and companies contacting via portal, email, chat or WhatsApp. Each customer only sees their own conversations.')} action={staff ? <Button onPress={() => setEditing(null)}>{t('Nuevo cliente', 'New customer')}</Button> : undefined} />
    <Feedback {...m} />
    {staff && editing !== undefined && <section className="panel"><h2>{editing ? t('Editar cliente', 'Edit customer') : t('Nuevo cliente', 'New customer')}</h2>
      <form key={editing?.id || 'new'} action={form => m.run(async () => { await api(editing ? `/cx/customers/${editing.id}` : '/cx/customers', json(editing ? 'PUT' : 'POST', { name: form.get('name'), email: form.get('email'), phone: form.get('phone'), company: form.get('company'), status: form.get('status'), consent_given: form.get('consent') === 'on', erasure_requested: form.get('erasure') === 'on' })); setEditing(undefined); customers.reload(); }, t('Cliente guardado.', 'Customer saved.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required defaultValue={editing?.name} /><Field label={t('Correo', 'Email')} name="email" type="email" defaultValue={editing?.email} /><Field label={t('Teléfono', 'Phone')} name="phone" defaultValue={editing?.phone} /><Field label={t('Empresa', 'Company')} name="company" defaultValue={editing?.company} /><Select label={t('Estado', 'Status')} name="status" defaultValue={editing?.status || 'active'}><option value="active">{t('Activo', 'Active')}</option><option value="inactive">{t('Inactivo', 'Inactive')}</option></Select></div>
        <label className="checkbox"><input type="checkbox" name="consent" defaultChecked={editing?.consent_given} />{t('Consentimiento de contacto registrado', 'Contact consent recorded')}</label>
        <label className="checkbox"><input type="checkbox" name="erasure" defaultChecked={editing?.erasure_requested} />{t('Eliminación de datos solicitada', 'Data erasure requested')}</label>
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditing(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...customers} retry={customers.reload}>{!(customers.data || []).length ? <Empty>{t('No hay clientes registrados.', 'No customers registered.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Cliente', 'Customer')}</th><th>{t('Contacto', 'Contact')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{customers.data!.map(c => <tr key={c.id}><td>{c.name}<small>{c.company || ''}</small></td><td>{c.email || c.phone || '—'}</td><td><Pill value={c.status} /></td><td><Button size="sm" variant="secondary" onPress={() => setSelected(c)}>{t('Contactos', 'Contacts')}</Button>{staff && <><Button size="sm" variant="secondary" onPress={() => setEditing(c)}>{t('Editar', 'Edit')}</Button><Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/cx/customers/${c.id}`, { method: 'DELETE' }); customers.reload(); }, t('Cliente eliminado.', 'Customer deleted.'))}>{t('Eliminar', 'Delete')}</Button></>}</td></tr>)}</tbody></table></div>}</LoadState>
      {selected && <section className="panel"><h2>{t('Contactos de', 'Contacts of')} {selected.name}</h2>
        <LoadState {...contacts} retry={contacts.reload}>{!(contacts.data || []).length ? <Empty>{t('Sin contactos registrados.', 'No contacts registered.')}</Empty> : contacts.data!.map(ct => <div className="list-item" key={ct.id}><span>{ct.name}</span><p>{ct.email} · {ct.role}</p>{staff && <Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/cx/customers/${selected.id}/contacts/${ct.id}`, { method: 'DELETE' }); contacts.reload(); }, t('Contacto eliminado.', 'Contact deleted.'))}>{t('Eliminar', 'Delete')}</Button>}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/cx/customers/${selected.id}/contacts`, json('POST', { name: form.get('name'), email: form.get('email'), role: form.get('role') })); contacts.reload(); }, t('Contacto añadido.', 'Contact added.'))}><div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('Correo', 'Email')} name="email" type="email" required /><Select label={t('Rol', 'Role')} name="role" defaultValue="contact"><option value="contact">{t('Contacto', 'Contact')}</option><option value="client_admin">{t('Administrador del cliente', 'Customer admin')}</option></Select></div><Button type="submit" isDisabled={m.busy}>{t('Añadir contacto', 'Add contact')}</Button></form>}
      </section>}
    </section>
  </>;
}

// ── Bandeja y gestión de conversaciones ─────────────────────────────

export function CXInbox({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const customers = useLoad(() => api<CXCustomer[]>('/cx/customers'));
  const [channel, setChannel] = useState('');
  const [status, setStatus] = useState('');
  const [selected, setSelected] = useState<CXConversation | null>(null);
  const conversations = useLoad(() => api<CXConversation[]>(`/cx/conversations${channel ? `?channel=${channel}` : ''}${status ? `${channel ? '&' : '?'}status=${status}` : ''}`), [channel, status]);
  const messages = useLoad(() => selected ? api<CXMessage[]>(`/cx/conversations/${selected.id}/messages`) : Promise.resolve([] as CXMessage[]), [selected?.id]);
  return <>
    <Heading title={t('Bandeja de atención', 'Support inbox')} description={t('Conversaciones de todos los canales, con historial completo y contexto de canal. Las notas internas solo las ve el equipo.', 'Conversations from every channel, with full history and channel context. Internal notes are staff-only.')} action={<Link className="button button--secondary" href="/articles">{t('Conocimiento y autoservicio', 'Knowledge and self-service')}</Link>} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><Select label={t('Canal', 'Channel')} value={channel} onChange={event => setChannel(event.target.value)}><option value="">{t('Todos los canales', 'All channels')}</option>{channels.map(c => <option key={c} value={c}>{c}</option>)}</Select><Select label={t('Estado', 'Status')} value={status} onChange={event => setStatus(event.target.value)}><option value="">{t('Todos los estados', 'All statuses')}</option>{Object.keys(conversationTransitions).map(s => <option key={s} value={s}>{s}</option>)}</Select></div></section>
    <section className="panel"><LoadState {...conversations} retry={conversations.reload}>{!(conversations.data || []).length ? <Empty>{t('No hay conversaciones en esta bandeja.', 'No conversations in this inbox.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Asunto', 'Subject')}</th><th>{t('Cliente', 'Customer')}</th><th>{t('Canal', 'Channel')}</th><th>{t('Estado', 'Status')}</th><th>{t('SLA', 'SLA')}</th></tr></thead><tbody>{conversations.data!.map(c => <tr key={c.id} onClick={() => setSelected(c)} style={{ cursor: 'pointer' }}><td>{c.subject}<small>{c.priority}{c.escalated ? ` · ${t('escalada', 'escalated')}` : ''}</small></td><td>{c.customer_name || c.customer_id.slice(0, 8)}</td><td>{c.channel}</td><td><Pill value={c.status} /></td><td>{isSLABreached(c) ? t('Incumplido', 'Breached') : t('En plazo', 'Within SLA')}</td></tr>)}</tbody></table></div>}</LoadState></section>
    {selected && <CXConversationDetail conversation={selected} user={user} messages={messages.data || []} reload={() => { messages.reload(); conversations.reload(); }} />}
    <section className="panel"><h2>{t('Nueva conversación', 'New conversation')}</h2>
      <form action={form => m.run(async () => { const result = await api<CXConversation>('/cx/conversations', json('POST', { customer_id: form.get('customer'), channel: form.get('channel'), subject: form.get('subject'), priority: form.get('priority') })); setSelected(result); conversations.reload(); }, t('Conversación creada.', 'Conversation created.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Customer')} name="customer" required>{(customers.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Select label={t('Canal', 'Channel')} name="channel" defaultValue="portal">{channels.map(c => <option key={c} value={c}>{c}</option>)}</Select><Field label={t('Asunto', 'Subject')} name="subject" required /><Select label={t('Prioridad', 'Priority')} name="priority" defaultValue="medium">{priorities.map(p => <option key={p} value={p}>{p}</option>)}</Select></div>
        <Button type="submit" isDisabled={m.busy}>{t('Abrir conversación', 'Open conversation')}</Button>
      </form></section>
  </>;
}

function CXConversationDetail({ conversation, user, messages, reload }: { conversation: CXConversation; user: User; messages: CXMessage[]; reload: () => void }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const [surveyDone, setSurveyDone] = useState(false);
  return <section className="panel"><h2>{conversation.subject}</h2>
    <p>{conversation.channel} · <Pill value={conversation.status} /> · {conversation.priority} · {t('Asignada a', 'Assigned to')} {conversation.assignee_name || t('nadie', 'nobody')} · SLA {conversation.sla_response_minutes}/{conversation.sla_resolution_minutes} {t('min', 'min')} · {isSLABreached(conversation) ? t('SLA incumplido', 'SLA breached') : t('SLA en plazo', 'SLA within deadline')}</p>
    <Feedback {...m} />
    {staff && <div className="form-actions">
      <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/cx/conversations/${conversation.id}/take`, { method: 'POST' }); reload(); }, t('Conversación tomada.', 'Conversation taken.'))}>{t('Tomar', 'Take')}</Button>
      <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/cx/conversations/${conversation.id}/release`, { method: 'POST' }); reload(); }, t('Conversación liberada.', 'Conversation released.'))}>{t('Liberar', 'Release')}</Button>
      <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/cx/conversations/${conversation.id}/escalate`, { method: 'POST' }); reload(); }, t('Conversación escalada.', 'Conversation escalated.'))}>{t('Escalar', 'Escalate')} ({t('a', 'to')} {nextPriority(conversation.priority)})</Button>
      {(conversationTransitions[conversation.status] || []).map(next => <Button key={next} size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/cx/conversations/${conversation.id}/status`, json('PATCH', { status: next })); reload(); }, t('Estado actualizado.', 'Status updated.'))}>{next}</Button>)}
    </div>}
    <div className="mt-6">{!messages.length ? <Empty>{t('Sin mensajes todavía.', 'No messages yet.')}</Empty> : messages.map(msg => <div className="list-item" key={msg.id}><span>{msg.sender_type}{msg.is_internal ? ` · ${t('nota interna', 'internal note')}` : ''}</span><p>{msg.body}</p><small>{msg.sender_name || ''} · {new Date(msg.created_at).toLocaleString()}</small></div>)}</div>
    {conversation.status !== 'closed' && <form className="mt-6" action={form => m.run(async () => { await api(`/cx/conversations/${conversation.id}/messages`, json('POST', { body: form.get('body'), is_internal: staff && form.get('internal') === 'on' })); reload(); }, t('Mensaje enviado.', 'Message sent.'))}>
      <Textarea label={t('Mensaje', 'Message')} name="body" required rows={4} />
      {staff && <label className="checkbox"><input type="checkbox" name="internal" />{t('Nota interna (no la ve el cliente)', 'Internal note (hidden from the customer)')}</label>}
      <Button type="submit" isDisabled={m.busy}>{t('Enviar', 'Send')}</Button>
    </form>}
    {(conversation.status === 'resolved' || conversation.status === 'closed') && !surveyDone && <form className="mt-6" action={form => m.run(async () => { await api(`/cx/conversations/${conversation.id}/surveys`, json('POST', { rating: Number(form.get('rating')), comment: form.get('comment') })); setSurveyDone(true); }, t('Encuesta enviada. Gracias.', 'Survey sent. Thank you.'))}>
      <div className="form-grid"><Select label={t('Satisfacción (1–5)', 'Satisfaction (1–5)')} name="rating" defaultValue="5">{[1, 2, 3, 4, 5].map(n => <option key={n} value={n}>{n}</option>)}</Select><Field label={t('Comentario', 'Comment')} name="comment" /></div>
      <Button type="submit" isDisabled={m.busy}>{t('Enviar encuesta', 'Send survey')}</Button>
    </form>}
  </section>;
}

// ── Canales: correo y WhatsApp ───────────────────────────────────────

export function CXChannels({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const emailAccounts = useLoad(() => api<CXEmailAccount[]>('/cx/email-accounts'));
  const whatsappAccounts = useLoad(() => api<CXWhatsAppAccount[]>('/cx/whatsapp-accounts'));
  return <>
    <Heading title={t('Canales', 'Channels')} description={t('Correo con SMTP/IMAP configurables y WhatsApp Business Platform con la API oficial de Meta. Solo se guarda configuración: ninguna contraseña ni token de acceso vive en la base de datos.', 'Email with configurable SMTP/IMAP and WhatsApp Business Platform via Meta\'s official API. Only configuration is stored: no password or access token lives in the database.')} />
    <Feedback {...m} />
    <section className="panel"><h2>{t('Correo', 'Email')}</h2>
      <LoadState {...emailAccounts} retry={emailAccounts.reload}>{!(emailAccounts.data || []).length ? <Empty>{t('Sin cuentas de correo configuradas.', 'No email accounts configured.')}</Empty> : emailAccounts.data!.map(a => <div className="list-item" key={a.id}><span>{a.name}</span><p>{a.smtp_host}:{a.smtp_port} · {a.imap_host}:{a.imap_port} · {a.active ? t('activa', 'active') : t('inactiva', 'inactive')}</p></div>)}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/cx/email-accounts', json('POST', { name: form.get('name'), smtp_host: form.get('smtp_host'), smtp_port: Number(form.get('smtp_port') || 587), smtp_username: form.get('smtp_username'), smtp_from: form.get('smtp_from'), imap_host: form.get('imap_host'), imap_port: Number(form.get('imap_port') || 993), imap_username: form.get('imap_username'), active: form.get('active') === 'on' })); emailAccounts.reload(); }, t('Cuenta de correo guardada.', 'Email account saved.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('SMTP host', 'SMTP host')} name="smtp_host" /><Field label={t('SMTP puerto', 'SMTP port')} name="smtp_port" type="number" defaultValue="587" /><Field label={t('SMTP usuario', 'SMTP user')} name="smtp_username" /><Field label={t('Remitente', 'From')} name="smtp_from" type="email" /><Field label={t('IMAP host', 'IMAP host')} name="imap_host" /><Field label={t('IMAP puerto', 'IMAP port')} name="imap_port" type="number" defaultValue="993" /><Field label={t('IMAP usuario', 'IMAP user')} name="imap_username" /></div>
        <label className="checkbox"><input type="checkbox" name="active" />{t('Activa', 'Active')}</label>
        <p>{t('Las credenciales SMTP/IMAP no se guardan aquí: viven en la configuración del operador (ver docs/cx-integraciones.md).', 'SMTP/IMAP credentials are not stored here: they live in the operator configuration (see docs/cx-integraciones.md).')}</p>
        <Button type="submit" isDisabled={m.busy}>{t('Guardar cuenta', 'Save account')}</Button>
      </form>}
    </section>
    <section className="panel"><h2>WhatsApp</h2>
      <LoadState {...whatsappAccounts} retry={whatsappAccounts.reload}>{!(whatsappAccounts.data || []).length ? <Empty>{t('Sin cuentas de WhatsApp configuradas.', 'No WhatsApp accounts configured.')}</Empty> : whatsappAccounts.data!.map(a => <div className="list-item" key={a.id}><span>{a.name}</span><p>{a.display_name} · {a.phone_number_id} · {a.active ? t('activa', 'active') : t('inactiva', 'inactive')}</p></div>)}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/cx/whatsapp-accounts', json('POST', { name: form.get('name'), business_id: form.get('business_id'), phone_number_id: form.get('phone_number_id'), display_name: form.get('display_name'), verify_token: form.get('verify_token'), active: form.get('active') === 'on' })); whatsappAccounts.reload(); }, t('Cuenta de WhatsApp guardada.', 'WhatsApp account saved.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('ID de empresa', 'Business ID')} name="business_id" /><Field label={t('ID de número', 'Phone number ID')} name="phone_number_id" /><Field label={t('Nombre visible', 'Display name')} name="display_name" /><Field label={t('Token de verificación del webhook', 'Webhook verify token')} name="verify_token" type="password" /></div>
        <label className="checkbox"><input type="checkbox" name="active" />{t('Activa', 'Active')}</label>
        <p>{t('Solo se guarda el hash del token de verificación, nunca el token ni el token de acceso de Meta. La activación viva queda pendiente de credenciales empresariales reales.', 'Only the verify-token hash is stored, never the token or Meta\'s access token. Live activation awaits real business credentials.')}</p>
        <Button type="submit" isDisabled={m.busy}>{t('Guardar cuenta', 'Save account')}</Button>
      </form>}
    </section>
  </>;
}

// ── Automatización: SLA, asignación y respuestas automáticas ────────

export function CXAutomation({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const policies = useLoad(() => api<CXSLAPolicy[]>('/cx/sla-policies'));
  const rules = useLoad(() => api<CXAssignmentRule[]>('/cx/assignment-rules'));
  const replies = useLoad(() => api<CXAutoReply[]>('/cx/auto-replies'));
  const users = useLoad(() => staff ? api<User[]>('/users') : Promise.resolve([] as User[]));
  return <>
    <Heading title={t('Automatización', 'Automation')} description={t('SLA por canal y prioridad, reglas de asignación automática y respuestas automáticas por palabra clave o canal.', 'SLA per channel and priority, automatic assignment rules and auto-replies by keyword or channel.')} />
    <Feedback {...m} />
    <section className="panel"><h2>{t('Políticas SLA', 'SLA policies')}</h2>
      <LoadState {...policies} retry={policies.reload}>{!(policies.data || []).length ? <Empty>{t('Sin políticas: se usan los valores por defecto por canal.', 'No policies: channel defaults apply.')}</Empty> : policies.data!.map(p => <div className="list-item" key={p.id}><span>{p.channel} · {p.priority}</span><p>{p.response_minutes} / {p.resolution_minutes} {t('min', 'min')}</p></div>)}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/cx/sla-policies', json('PUT', { channel: form.get('channel'), priority: form.get('priority'), response_minutes: Number(form.get('resp') || 0), resolution_minutes: Number(form.get('res') || 0) })); policies.reload(); }, t('Política SLA guardada.', 'SLA policy saved.'))}>
        <div className="form-grid"><Select label={t('Canal', 'Channel')} name="channel" defaultValue="portal">{channels.map(c => <option key={c} value={c}>{c}</option>)}</Select><Select label={t('Prioridad', 'Priority')} name="priority" defaultValue="medium">{priorities.map(p => <option key={p} value={p}>{p}</option>)}</Select><Field label={t('Respuesta (min)', 'Response (min)')} name="resp" type="number" min="0" /><Field label={t('Resolución (min)', 'Resolution (min)')} name="res" type="number" min="0" /></div>
        <Button type="submit" isDisabled={m.busy}>{t('Guardar política', 'Save policy')}</Button>
      </form>}
    </section>
    <section className="panel"><h2>{t('Reglas de asignación', 'Assignment rules')}</h2>
      <LoadState {...rules} retry={rules.reload}>{!(rules.data || []).length ? <Empty>{t('Sin reglas de asignación.', 'No assignment rules.')}</Empty> : rules.data!.map(r => <div className="list-item" key={r.id}><span>{r.channel || t('cualquier canal', 'any channel')} · {r.priority || t('cualquier prioridad', 'any priority')}</span><p>{t('posición', 'position')} {r.position}</p>{staff && <Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/cx/assignment-rules/${r.id}`, { method: 'DELETE' }); rules.reload(); }, t('Regla eliminada.', 'Rule deleted.'))}>{t('Eliminar', 'Delete')}</Button>}</div>)}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/cx/assignment-rules', json('POST', { channel: form.get('channel'), priority: form.get('priority'), assignee_user_id: form.get('assignee'), position: Number(form.get('position') || 0), active: true })); rules.reload(); }, t('Regla guardada.', 'Rule saved.'))}>
        <div className="form-grid"><Select label={t('Canal', 'Channel')} name="channel"><option value="">{t('Cualquier canal', 'Any channel')}</option>{channels.map(c => <option key={c} value={c}>{c}</option>)}</Select><Select label={t('Prioridad', 'Priority')} name="priority"><option value="">{t('Cualquier prioridad', 'Any priority')}</option>{priorities.map(p => <option key={p} value={p}>{p}</option>)}</Select><Select label={t('Agente asignado', 'Assigned agent')} name="assignee" required>{(users.data || []).map(u => <option key={u.id} value={u.id}>{u.name}</option>)}</Select><Field label={t('Posición', 'Position')} name="position" type="number" defaultValue="0" /></div>
        <Button type="submit" isDisabled={m.busy}>{t('Guardar regla', 'Save rule')}</Button>
      </form>}
    </section>
    <section className="panel"><h2>{t('Respuestas automáticas', 'Auto-replies')}</h2>
      <LoadState {...replies} retry={replies.reload}>{!(replies.data || []).length ? <Empty>{t('Sin respuestas automáticas.', 'No auto-replies.')}</Empty> : replies.data!.map(r => <div className="list-item" key={r.id}><span>{r.trigger_kind}{r.keyword ? ` · ${r.keyword}` : ''}{r.channel ? ` · ${r.channel}` : ''}</span><p>{r.response_text.slice(0, 120)}</p>{staff && <Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/cx/auto-replies/${r.id}`, { method: 'DELETE' }); replies.reload(); }, t('Respuesta eliminada.', 'Reply deleted.'))}>{t('Eliminar', 'Delete')}</Button>}</div>)}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/cx/auto-replies', json('POST', { trigger_kind: form.get('trigger'), keyword: form.get('keyword'), channel: form.get('channel'), response_text: form.get('response'), active: true })); replies.reload(); }, t('Respuesta guardada.', 'Reply saved.'))}>
        <div className="form-grid"><Select label={t('Disparador', 'Trigger')} name="trigger" defaultValue="keyword"><option value="keyword">{t('Palabra clave', 'Keyword')}</option><option value="channel">{t('Canal', 'Channel')}</option></Select><Field label={t('Palabra clave', 'Keyword')} name="keyword" /><Select label={t('Canal', 'Channel')} name="channel"><option value="">{t('Cualquier canal', 'Any channel')}</option>{channels.map(c => <option key={c} value={c}>{c}</option>)}</Select></div>
        <Textarea label={t('Respuesta', 'Reply')} name="response" required rows={3} />
        <Button type="submit" isDisabled={m.busy}>{t('Guardar respuesta', 'Save reply')}</Button>
      </form>}
    </section>
  </>;
}

// ── Informes y satisfacción ──────────────────────────────────────────

export function CXReportView() {
  const t = useText();
  const report = useLoad(() => api<CXReport>('/cx/report'));
  const surveys = useLoad(() => api<CXSurvey[]>('/cx/surveys'));
  return <>
    <Heading title={t('Informes de atención', 'Customer service reports')} description={t('Volumen por canal, conversaciones abiertas y resueltas, SLA incumplido y satisfacción media.', 'Volume per channel, open and resolved conversations, breached SLA and average satisfaction.')} />
    <section className="panel"><LoadState {...report} retry={report.reload}>{report.data && <div className="form-grid">
      <p>{t('Clientes', 'Customers')}: {report.data.customers_total}</p>
      <p>{t('Conversaciones abiertas', 'Open conversations')}: {report.data.open_conversations}</p>
      <p>{t('Conversaciones resueltas', 'Resolved conversations')}: {report.data.resolved_conversations}</p>
      <p>{t('SLA incumplido', 'SLA breached')}: {report.data.sla_breached}</p>
      <p>{t('Encuestas', 'Surveys')}: {report.data.survey_count} · {t('satisfacción media', 'average satisfaction')}: {report.data.avg_satisfaction.toFixed(2)}</p>
      <p>{t('Por canal', 'By channel')}: {Object.entries(report.data.conversations_by_channel).map(([ch, n]) => `${ch}: ${n}`).join(' · ') || '—'}</p>
    </div>}</LoadState></section>
    <section className="panel"><h2>{t('Encuestas', 'Surveys')}</h2>
      <LoadState {...surveys} retry={surveys.reload}>{!(surveys.data || []).length ? <Empty>{t('Sin encuestas todavía.', 'No surveys yet.')}</Empty> : <>
        <p>{t('Satisfacción media calculada en la interfaz', 'Average satisfaction computed in the UI')}: {avgSatisfaction(surveys.data!).toFixed(2)}</p>
        {surveys.data!.map(s => <div className="list-item" key={s.id}><span>{s.rating} / 5</span><p>{s.comment || '—'}</p></div>)}
      </>}</LoadState>
    </section>
  </>;
}
