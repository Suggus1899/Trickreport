'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { api, json } from './api';
import type { Asset, CIType, MSPAlert, MSPAvailability, MSPBillingDraft, MSPClient, MSPClientAsset, MSPContact, MSPContract, MSPDiscoveredAsset, MSPExpense, MSPIntervention, MSPMonitoringSource, MSPProfitability, MSPProject, MSPProjectTask, MSPRemoteSession, MSPReport, MSPTechnician, MSPTicket, MSPTimeEntry, User } from './contracts';
import { interventionTransitions, projectTransitions, taskTransitions } from './contracts';
import { Button, Empty, Feedback, Field, Heading, LoadState, Pill, Select, Textarea, useLoad, useMutation, useText } from './ui';

// MSP (Etapa 4): operación multicliente del proveedor. Todas las
// pantallas consumen la API real (/api/v1); el servidor Go sigue
// siendo la autoridad de permisos y de aislamiento entre clientes (la
// ocultación de un botón aquí nunca sustituye esa comprobación). El
// cambio de cliente es explícito: el selector persiste en localStorage
// y cada pantalla filtra por él; el contrato vigente se muestra junto
// al cliente en la mesa multicliente.

// ── Helpers puros (cubiertos por tests Vitest) ──────────────────────

export function formatCents(cents: number, currency = 'EUR'): string {
  try {
    return new Intl.NumberFormat('es-ES', { style: 'currency', currency }).format(cents / 100);
  } catch {
    return `${(cents / 100).toFixed(2)} ${currency}`;
  }
}

export function projectProgressPct(tasks: { status: string }[]): number {
  if (!tasks.length) return 0;
  return Math.floor((tasks.filter(t => t.status === 'done').length * 100) / tasks.length);
}

export function isActiveContract(c: MSPContract, at = new Date()): boolean {
  if (c.status !== 'active') return false;
  if (new Date(c.starts_at) > at) return false;
  if (c.ends_at && new Date(c.ends_at) < at) return false;
  return true;
}

export function activeContractForClient(contracts: MSPContract[], clientId: string, at = new Date()): MSPContract | undefined {
  return contracts.filter(c => c.client_id === clientId && isActiveContract(c, at)).sort((a, b) => +new Date(b.starts_at) - +new Date(a.starts_at))[0];
}

const MSP_CLIENT_KEY = 'trickreport_msp_client';

export function readMSPClient(): string {
  if (typeof window === 'undefined') return '';
  try { return window.localStorage.getItem(MSP_CLIENT_KEY) || ''; } catch { return ''; }
}

export function writeMSPClient(value: string) {
  try { window.localStorage.setItem(MSP_CLIENT_KEY, value); } catch { /* storage unavailable */ }
}

function useMSPClientFilter() {
  const [clientId, setClientId] = useState(() => readMSPClient());
  function change(value: string) { setClientId(value); writeMSPClient(value); }
  return { clientId, change };
}

export function ClientFilter({ clients, value, onChange }: { clients: MSPClient[]; value: string; onChange: (value: string) => void }) {
  const t = useText();
  return <Select label={t('Cliente activo', 'Active client')} value={value} onChange={event => onChange(event.target.value)}>
    <option value="">{t('Todos los clientes', 'All clients')}</option>
    {clients.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
  </Select>;
}

function clientQuery(clientId: string) { return clientId ? `?client_id=${clientId}` : ''; }

// datetime-local values are wall-clock time in the user's zone; the API stores
// UTC instants. These two functions are the only conversion pair.
export function toISO(local: FormDataEntryValue | null): string {
  if (!local) return new Date().toISOString();
  const d = new Date(String(local));
  return Number.isNaN(d.getTime()) ? new Date().toISOString() : d.toISOString();
}

export function toLocalInput(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return '';
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

// The contract form shows only part of the stored contract. A PUT replaces the
// whole row, so fields the form does not show are carried over from the stored
// contract instead of being filled with defaults.
export function contractPayload(form: { get(name: string): FormDataEntryValue | null }, editing?: MSPContract | null) {
  return {
    client_id: form.get('client'),
    name: form.get('name'),
    status: form.get('status'),
    starts_at: toISO(form.get('starts')),
    ends_at: form.get('ends') ? toISO(form.get('ends')) : null,
    coverage: editing?.coverage ?? '',
    currency: editing?.currency || 'EUR',
    monthly_fee_cents: Math.round(Number(form.get('fee') || 0) * 100),
    hourly_rate_cents: Math.round(Number(form.get('rate') || 0) * 100),
    included_hours_monthly: Number(form.get('included') || 0),
    sla_response_minutes: Number(form.get('resp') || 240),
    sla_resolution_minutes: Number(form.get('res') || 1440),
    sla_overrides: editing?.sla_overrides ?? {},
  };
}

// ── Clientes y contactos ────────────────────────────────────────────

export function Clients({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const [editing, setEditing] = useState<MSPClient | null>();
  const [selected, setSelected] = useState<MSPClient | null>(null);
  const contacts = useLoad(() => selected ? api<MSPContact[]>(`/msp/clients/${selected.id}/contacts`) : Promise.resolve([] as MSPContact[]), [selected?.id]);
  const contracts = useLoad(() => selected ? api<MSPContract[]>(`/msp/contracts?client_id=${selected.id}`) : Promise.resolve([] as MSPContract[]), [selected?.id]);
  return <>
    <Heading title={t('Clientes MSP', 'MSP clients')} description={t('Empresas atendidas por el proveedor, cada una con su espacio separado.', 'Companies served by the provider, each in its own separated space.')} action={staff ? <Button onPress={() => setEditing(null)}>{t('Nuevo cliente', 'New client')}</Button> : undefined} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      <p>{clientId ? t('Cliente activo: ', 'Active client: ') + (clients.data || []).find(c => c.id === clientId)?.name : t('Sin cliente activo: las pantallas muestran todos los clientes visibles.', 'No active client: screens show every visible client.')}</p>
    </section>
    {staff && editing !== undefined && <section className="panel"><h2>{editing ? t('Editar cliente', 'Edit client') : t('Nuevo cliente', 'New client')}</h2>
      <form key={editing?.id || 'new'} action={form => m.run(async () => { await api(editing ? `/msp/clients/${editing.id}` : '/msp/clients', json(editing ? 'PUT' : 'POST', { name: form.get('name'), domain: form.get('domain'), city: form.get('city'), status: form.get('status'), contact_name: form.get('contact'), contact_email: form.get('email') })); setEditing(undefined); clients.reload(); }, t('Cliente guardado.', 'Client saved.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required defaultValue={editing?.name} /><Field label={t('Dominio', 'Domain')} name="domain" defaultValue={editing?.domain} /><Field label={t('Ciudad', 'City')} name="city" defaultValue={editing?.city} /><Select label={t('Estado', 'Status')} name="status" defaultValue={editing?.status || 'active'}><option value="active">{t('Activo', 'Active')}</option><option value="inactive">{t('Inactivo', 'Inactive')}</option><option value="suspended">{t('Suspendido', 'Suspended')}</option></Select><Field label={t('Contacto principal', 'Main contact')} name="contact" defaultValue={editing?.contact_name} /><Field label={t('Correo del contacto', 'Contact email')} name="email" type="email" defaultValue={editing?.contact_email} /></div>
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditing(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...clients} retry={clients.reload}>{!(clients.data || []).length ? <Empty>{t('No hay clientes registrados.', 'No clients registered.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Cliente', 'Client')}</th><th>{t('Dominio', 'Domain')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{clients.data!.map(c => <tr key={c.id}><td>{c.name}<small>{c.city || c.contact_email || ''}</small></td><td>{c.domain || '—'}</td><td><Pill value={c.status} /></td><td><Button size="sm" variant="secondary" onPress={() => setSelected(c)}>{t('Contactos y contratos', 'Contacts and contracts')}</Button><Button size="sm" variant="secondary" onPress={() => change(c.id)}>{t('Usar como activo', 'Set as active')}</Button>{staff && <><Button size="sm" variant="secondary" onPress={() => setEditing(c)}>{t('Editar', 'Edit')}</Button><Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/msp/clients/${c.id}`, { method: 'DELETE' }); clients.reload(); }, t('Cliente eliminado.', 'Client deleted.'))}>{t('Eliminar', 'Delete')}</Button></>}</td></tr>)}</tbody></table></div>}</LoadState>
      {selected && <section className="panel"><h2>{t('Contactos y contratos de', 'Contacts and contracts of')} {selected.name}</h2>
        <LoadState {...contacts} retry={contacts.reload}>{!(contacts.data || []).length ? <Empty>{t('Sin contactos registrados.', 'No contacts registered.')}</Empty> : contacts.data!.map(ct => <div className="list-item" key={ct.id}><span>{ct.name}</span><p>{ct.email} · {ct.role}</p>{staff && <Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/msp/clients/${selected.id}/contacts/${ct.id}`, { method: 'DELETE' }); contacts.reload(); }, t('Contacto eliminado.', 'Contact deleted.'))}>{t('Eliminar', 'Delete')}</Button>}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/msp/clients/${selected.id}/contacts`, json('POST', { name: form.get('name'), email: form.get('email'), role: form.get('role') })); contacts.reload(); }, t('Contacto añadido.', 'Contact added.'))}><div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('Correo', 'Email')} name="email" type="email" required /><Select label={t('Rol', 'Role')} name="role" defaultValue="contact"><option value="contact">{t('Contacto', 'Contact')}</option><option value="client_admin">{t('Administrador del cliente', 'Client admin')}</option></Select></div><Button type="submit" isDisabled={m.busy}>{t('Añadir contacto', 'Add contact')}</Button></form>}
        <h2 className="mt-6">{t('Contratos', 'Contracts')}</h2>
        <LoadState {...contracts} retry={contracts.reload}>{!(contracts.data || []).length ? <Empty>{t('Sin contratos para este cliente.', 'No contracts for this client.')}</Empty> : contracts.data!.map(ct => <div className="list-item" key={ct.id}><span>{ct.name}</span><p>{ct.status} · {formatCents(ct.monthly_fee_cents, ct.currency)} / {t('mes', 'month')} · SLA {ct.sla_response_minutes}/{ct.sla_resolution_minutes} {t('min', 'min')}</p></div>)}</LoadState>
      </section>}
    </section>
  </>;
}

// ── Contratos y SLA ─────────────────────────────────────────────────

export function MSPContracts({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const contracts = useLoad(() => api<MSPContract[]>(`/msp/contracts${clientQuery(clientId)}`), [clientId]);
  const [editing, setEditing] = useState<MSPContract | null>();
  return <>
    <Heading title={t('Contratos y SLA', 'Contracts and SLA')} description={t('Vigencia, cobertura y SLA aplicable por cliente.', 'Validity, coverage and applicable SLA per client.')} action={staff ? <Button onPress={() => setEditing(null)}>{t('Nuevo contrato', 'New contract')}</Button> : undefined} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div></section>
    {staff && editing !== undefined && <section className="panel"><h2>{editing ? t('Editar contrato', 'Edit contract') : t('Nuevo contrato', 'New contract')}</h2>
      <form key={editing?.id || 'new'} action={form => m.run(async () => { await api(editing ? `/msp/contracts/${editing.id}` : '/msp/contracts', json(editing ? 'PUT' : 'POST', contractPayload(form, editing))); setEditing(undefined); contracts.reload(); }, t('Contrato guardado.', 'Contract saved.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={editing?.client_id || clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Field label={t('Nombre', 'Name')} name="name" required defaultValue={editing?.name} /><Select label={t('Estado', 'Status')} name="status" defaultValue={editing?.status || 'draft'}><option value="draft">{t('Borrador', 'Draft')}</option><option value="active">{t('Activo', 'Active')}</option><option value="suspended">{t('Suspendido', 'Suspended')}</option><option value="expired">{t('Vencido', 'Expired')}</option><option value="cancelled">{t('Cancelado', 'Cancelled')}</option></Select><Field label={t('Inicio', 'Start')} name="starts" type="datetime-local" required defaultValue={editing ? toLocalInput(editing.starts_at) : ''} /><Field label={t('Fin', 'End')} name="ends" type="datetime-local" defaultValue={editing?.ends_at ? toLocalInput(editing.ends_at) : ''} /><Field label={t('Cuota mensual (€)', 'Monthly fee (€)')} name="fee" type="number" step="0.01" min="0" defaultValue={editing ? editing.monthly_fee_cents / 100 : ''} /><Field label={t('Tarifa horaria (€)', 'Hourly rate (€)')} name="rate" type="number" step="0.01" min="0" defaultValue={editing ? editing.hourly_rate_cents / 100 : ''} /><Field label={t('Horas incluidas / mes', 'Included hours / month')} name="included" type="number" step="0.5" min="0" defaultValue={editing?.included_hours_monthly} /><Field label={t('SLA respuesta (min)', 'SLA response (min)')} name="resp" type="number" min="0" defaultValue={editing?.sla_response_minutes ?? 240} /><Field label={t('SLA resolución (min)', 'SLA resolution (min)')} name="res" type="number" min="0" defaultValue={editing?.sla_resolution_minutes ?? 1440} /></div>
        <Textarea label={t('Cobertura', 'Coverage')} name="coverage" defaultValue={editing?.coverage} rows={3} />
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditing(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...contracts} retry={contracts.reload}>{!(contracts.data || []).length ? <Empty>{t('No hay contratos.', 'No contracts.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Contrato', 'Contract')}</th><th>{t('Estado', 'Status')}</th><th>{t('Cuota', 'Fee')}</th><th>{t('SLA', 'SLA')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{contracts.data!.map(c => <tr key={c.id}><td>{c.name}<small>{(clients.data || []).find(cl => cl.id === c.client_id)?.name || c.client_id.slice(0, 8)}</small></td><td><Pill value={c.status} /></td><td>{formatCents(c.monthly_fee_cents, c.currency)}</td><td>{c.sla_response_minutes} / {c.sla_resolution_minutes} {t('min', 'min')}</td><td>{staff && <Button size="sm" variant="secondary" onPress={() => setEditing(c)}>{t('Editar', 'Edit')}</Button>}</td></tr>)}</tbody></table></div>}</LoadState></section>
  </>;
}

// ── Mesa multicliente ───────────────────────────────────────────────

export function MSPTickets({ user }: { user: User }) {
  const t = useText(), m = useMutation(), router = useRouter();
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const tickets = useLoad(() => api<MSPTicket[]>(`/msp/tickets${clientQuery(clientId)}`), [clientId]);
  const contracts = useLoad(() => api<MSPContract[]>(`/msp/contracts${clientQuery(clientId)}`), [clientId]);
  const active = clientId ? activeContractForClient(contracts.data || [], clientId) : undefined;
  const [creating, setCreating] = useState(false);
  void user;
  return <>
    <Heading title={t('Mesa multicliente', 'Multi-client desk')} description={t('Tickets con cliente explícito y su contrato visible en contexto.', 'Tickets with an explicit client and its contract visible in context.')} action={<Button onPress={() => setCreating(true)}>{t('Nuevo ticket multicliente', 'New multi-client ticket')}</Button>} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      {clientId && <p>{active ? t('Contrato vigente: ', 'Current contract: ') + `${active.name} · SLA ${active.sla_response_minutes}/${active.sla_resolution_minutes} min · ${formatCents(active.monthly_fee_cents, active.currency)} / ${t('mes', 'month')}` : t('Este cliente no tiene contrato activo: el ticket se creará sin SLA.', 'This client has no active contract: the ticket will be created without an SLA.')}</p>}
    </section>
    {creating && <section className="panel"><h2>{t('Nuevo ticket multicliente', 'New multi-client ticket')}</h2>
      <form action={form => m.run(async () => { const res = await api<{ ticket_id: string }>('/msp/tickets', json('POST', { client_id: form.get('client'), title: form.get('title'), description: form.get('description'), priority: form.get('priority'), category: form.get('category') })); setCreating(false); tickets.reload(); router.push(`/tickets/${res.ticket_id}`); }, t('Ticket creado.', 'Ticket created.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Field label={t('Título', 'Title')} name="title" required /><Select label={t('Prioridad', 'Priority')} name="priority" defaultValue="medium"><option value="low">{t('Baja', 'Low')}</option><option value="medium">{t('Media', 'Medium')}</option><option value="high">{t('Alta', 'High')}</option><option value="critical">{t('Crítica', 'Critical')}</option></Select><Field label={t('Categoría', 'Category')} name="category" defaultValue="general" /></div>
        <Textarea label={t('Descripción', 'Description')} name="description" rows={4} />
        <div className="form-actions"><Button variant="secondary" onPress={() => setCreating(false)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Crear ticket', 'Create ticket')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...tickets} retry={tickets.reload}>{!(tickets.data || []).length ? <Empty>{t('No hay tickets multicliente.', 'No multi-client tickets.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Ticket', 'Ticket')}</th><th>{t('Cliente', 'Client')}</th><th>{t('Estado', 'Status')}</th><th>{t('SLA (min)', 'SLA (min)')}</th></tr></thead><tbody>{tickets.data!.map(tk => <tr key={tk.id}><td><Link href={`/tickets/${tk.id}`}>{tk.title}</Link><small>{tk.priority} · {tk.category}</small></td><td>{tk.client_name}</td><td><Pill value={tk.status} /></td><td>{tk.contract_id ? `${tk.sla_response_minutes} / ${tk.sla_resolution_minutes}` : t('Sin contrato', 'No contract')}</td></tr>)}</tbody></table></div>}</LoadState></section>
  </>;
}

// ── Técnicos y planificación ────────────────────────────────────────

export function Technicians({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const techs = useLoad(() => api<MSPTechnician[]>('/msp/technicians'));
  const users = useLoad(() => staff ? api<User[]>('/users') : Promise.resolve([] as User[]));
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const interventions = useLoad(() => api<MSPIntervention[]>(`/msp/interventions${clientQuery(clientId)}`), [clientId]);
  const [selected, setSelected] = useState<MSPTechnician | null>(null);
  const availability = useLoad(() => selected ? api<MSPAvailability[]>(`/msp/technicians/${selected.id}/availability`) : Promise.resolve([] as MSPAvailability[]), [selected?.id]);
  const days = [t('Domingo', 'Sunday'), t('Lunes', 'Monday'), t('Martes', 'Tuesday'), t('Miércoles', 'Wednesday'), t('Jueves', 'Thursday'), t('Viernes', 'Friday'), t('Sábado', 'Saturday')];
  return <>
    <Heading title={t('Técnicos y planificación', 'Technicians and planning')} description={t('Técnicos, disponibilidad e intervenciones programadas.', 'Technicians, availability and scheduled interventions.')} />
    <Feedback {...m} />
    <section className="panel"><LoadState {...techs} retry={techs.reload}>{!(techs.data || []).length ? <Empty>{t('No hay técnicos registrados.', 'No technicians registered.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Técnico', 'Technician')}</th><th>{t('Habilidades', 'Skills')}</th><th>{t('Tarifa', 'Rate')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{techs.data!.map(tech => <tr key={tech.id}><td>{tech.user_name}<small>{tech.active ? t('Activo', 'Active') : t('Inactivo', 'Inactive')}</small></td><td>{tech.skills.join(', ') || '—'}</td><td>{formatCents(tech.hourly_rate_cents)} / h</td><td><Button size="sm" variant="secondary" onPress={() => setSelected(tech)}>{t('Disponibilidad', 'Availability')}</Button></td></tr>)}</tbody></table></div>}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/msp/technicians', json('POST', { user_id: form.get('user'), skills: String(form.get('skills') || '').split(',').map(s => s.trim()).filter(Boolean), hourly_rate_cents: Math.round(Number(form.get('rate') || 0) * 100), active: true })); techs.reload(); }, t('Técnico creado.', 'Technician created.'))}><div className="form-grid"><Select label={t('Usuario', 'User')} name="user" required>{(users.data || []).map(u => <option key={u.id} value={u.id}>{u.name} ({u.role})</option>)}</Select><Field label={t('Habilidades (separadas por comas)', 'Skills (comma separated)')} name="skills" /><Field label={t('Tarifa horaria (€)', 'Hourly rate (€)')} name="rate" type="number" step="0.01" min="0" /></div><Button type="submit" isDisabled={m.busy}>{t('Crear técnico', 'Create technician')}</Button></form>}
      {selected && <section className="panel"><h2>{t('Disponibilidad de', 'Availability of')} {selected.user_name}</h2>
        <LoadState {...availability} retry={availability.reload}>{!(availability.data || []).length ? <Empty>{t('Sin disponibilidad registrada.', 'No availability registered.')}</Empty> : availability.data!.map(a => <div className="list-item" key={a.id}><span>{days[a.day_of_week]}</span><p>{a.start_time} – {a.end_time}</p></div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/msp/technicians/${selected.id}/availability`, json('POST', { day_of_week: Number(form.get('day')), start_time: form.get('start'), end_time: form.get('end') })); availability.reload(); }, t('Disponibilidad guardada.', 'Availability saved.'))}><div className="form-grid"><Select label={t('Día', 'Day')} name="day" defaultValue="1">{days.map((d, i) => <option key={i} value={i}>{d}</option>)}</Select><Field label={t('Desde', 'From')} name="start" type="time" required /><Field label={t('Hasta', 'To')} name="end" type="time" required /></div><Button type="submit" isDisabled={m.busy}>{t('Añadir disponibilidad', 'Add availability')}</Button></form>}
      </section>}
    </section>
    <section className="panel"><div className="section-header"><h2>{t('Intervenciones', 'Interventions')}</h2></div>
      <div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      <LoadState {...interventions} retry={interventions.reload}>{!(interventions.data || []).length ? <Empty>{t('No hay intervenciones.', 'No interventions.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Intervención', 'Intervention')}</th><th>{t('Ventana', 'Window')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{interventions.data!.map(v => <tr key={v.id}><td>{v.title}</td><td>{new Date(v.starts_at).toLocaleString()} – {new Date(v.ends_at).toLocaleString()}</td><td><Pill value={v.status} /></td><td>{staff && (interventionTransitions[v.status] || []).map(next => <Button key={next} size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/interventions/${v.id}/status`, json('PATCH', { status: next })); interventions.reload(); })}>{next}</Button>)}</td></tr>)}</tbody></table></div>}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/msp/interventions', json('POST', { client_id: form.get('client'), technician_id: form.get('tech'), title: form.get('title'), starts_at: toISO(form.get('starts')), ends_at: toISO(form.get('ends')) })); interventions.reload(); }, t('Intervención programada.', 'Intervention scheduled.'))}><div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Select label={t('Técnico', 'Technician')} name="tech" required>{(techs.data || []).map(tech => <option key={tech.id} value={tech.id}>{tech.user_name}</option>)}</Select><Field label={t('Título', 'Title')} name="title" required /><Field label={t('Inicio', 'Start')} name="starts" type="datetime-local" required /><Field label={t('Fin', 'End')} name="ends" type="datetime-local" required /></div><Button type="submit" isDisabled={m.busy}>{t('Programar intervención', 'Schedule intervention')}</Button></form>}
    </section>
  </>;
}

// ── Horas y gastos ──────────────────────────────────────────────────

export function TimeAndExpenses({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const techs = useLoad(() => api<MSPTechnician[]>('/msp/technicians'));
  const entries = useLoad(() => api<MSPTimeEntry[]>(`/msp/time-entries${clientQuery(clientId)}`), [clientId]);
  const expenses = useLoad(() => api<MSPExpense[]>(`/msp/expenses${clientQuery(clientId)}`), [clientId]);
  const totalCost = (entries.data || []).reduce((sum, e) => sum + e.cost_cents, 0);
  const totalExpenses = (expenses.data || []).reduce((sum, e) => sum + e.amount_cents, 0);
  return <>
    <Heading title={t('Horas, gastos y costes', 'Hours, expenses and costs')} description={t('Tiempo imputable por ticket o intervención, con coste calculado desde la tarifa del técnico.', 'Billable time per ticket or intervention, costed from the technician rate.')} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div><p>{t('Coste de tiempo registrado: ', 'Recorded time cost: ')}{formatCents(totalCost)} · {t('Gastos: ', 'Expenses: ')}{formatCents(totalExpenses)}</p></section>
    {staff && <section className="panel"><h2>{t('Registrar tiempo', 'Log time')}</h2>
      <form action={form => m.run(async () => { await api('/msp/time-entries', json('POST', { client_id: form.get('client'), technician_id: form.get('tech'), minutes: Number(form.get('minutes')), description: form.get('description'), billable: form.get('billable') === 'on', occurred_at: form.get('date') || undefined })); entries.reload(); }, t('Tiempo registrado.', 'Time logged.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Select label={t('Técnico', 'Technician')} name="tech" required>{(techs.data || []).map(tech => <option key={tech.id} value={tech.id}>{tech.user_name}</option>)}</Select><Field label={t('Minutos', 'Minutes')} name="minutes" type="number" min="1" max="1440" required /><Field label={t('Fecha', 'Date')} name="date" type="date" /></div>
        <Textarea label={t('Descripción', 'Description')} name="description" rows={2} />
        <label className="checkbox"><input type="checkbox" name="billable" defaultChecked />{t('Facturable', 'Billable')}</label>
        <div className="form-actions"><Button type="submit" isDisabled={m.busy}>{t('Registrar tiempo', 'Log time')}</Button></div>
      </form>
      <h2 className="mt-6">{t('Registrar gasto', 'Log expense')}</h2>
      <form action={form => m.run(async () => { await api('/msp/expenses', json('POST', { client_id: form.get('client'), description: form.get('description'), category: form.get('category'), amount_cents: Math.round(Number(form.get('amount') || 0) * 100), incurred_at: form.get('date') || undefined })); expenses.reload(); }, t('Gasto registrado.', 'Expense logged.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Field label={t('Importe (€)', 'Amount (€)')} name="amount" type="number" step="0.01" min="0.01" required /><Field label={t('Categoría', 'Category')} name="category" defaultValue="other" /><Field label={t('Fecha', 'Date')} name="date" type="date" /></div>
        <Textarea label={t('Descripción', 'Description')} name="description" rows={2} />
        <div className="form-actions"><Button type="submit" isDisabled={m.busy}>{t('Registrar gasto', 'Log expense')}</Button></div>
      </form></section>}
    <section className="panel"><h2>{t('Tiempo registrado', 'Logged time')}</h2>
      <LoadState {...entries} retry={entries.reload}>{!(entries.data || []).length ? <Empty>{t('Sin registros de tiempo.', 'No time entries.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Descripción', 'Description')}</th><th>{t('Minutos', 'Minutes')}</th><th>{t('Coste', 'Cost')}</th><th>{t('Facturado', 'Billed')}</th></tr></thead><tbody>{entries.data!.map(e => <tr key={e.id}><td>{e.description || '—'}<small>{new Date(e.occurred_at).toLocaleDateString()}</small></td><td>{e.minutes}</td><td>{formatCents(e.cost_cents)}</td><td>{e.billed ? t('Sí', 'Yes') : t('No', 'No')}</td></tr>)}</tbody></table></div>}</LoadState>
      <h2 className="mt-6">{t('Gastos', 'Expenses')}</h2>
      <LoadState {...expenses} retry={expenses.reload}>{!(expenses.data || []).length ? <Empty>{t('Sin gastos.', 'No expenses.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Descripción', 'Description')}</th><th>{t('Categoría', 'Category')}</th><th>{t('Importe', 'Amount')}</th></tr></thead><tbody>{expenses.data!.map(e => <tr key={e.id}><td>{e.description || '—'}<small>{new Date(e.incurred_at).toLocaleDateString()}</small></td><td>{e.category}</td><td>{formatCents(e.amount_cents)}</td></tr>)}</tbody></table></div>}</LoadState>
    </section>
  </>;
}

// ── Proyectos ───────────────────────────────────────────────────────

export function Projects({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const projects = useLoad(() => api<MSPProject[]>(`/msp/projects${clientQuery(clientId)}`), [clientId]);
  const [selected, setSelected] = useState<MSPProject | null>(null);
  const tasks = useLoad(() => selected ? api<MSPProjectTask[]>(`/msp/projects/${selected.id}/tasks`) : Promise.resolve([] as MSPProjectTask[]), [selected?.id]);
  const progress = projectProgressPct(tasks.data || []);
  return <>
    <Heading title={t('Proyectos de servicio', 'Service projects')} description={t('Proyectos por cliente, con fases en tareas y avance medible.', 'Projects per client, with phases as tasks and measurable progress.')} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div></section>
    {staff && <section className="panel"><h2>{t('Nuevo proyecto', 'New project')}</h2>
      <form action={form => m.run(async () => { await api('/msp/projects', json('POST', { client_id: form.get('client'), name: form.get('name'), budget_cents: Math.round(Number(form.get('budget') || 0) * 100) })); projects.reload(); }, t('Proyecto creado.', 'Project created.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('Presupuesto (€)', 'Budget (€)')} name="budget" type="number" step="0.01" min="0" /></div>
        <div className="form-actions"><Button type="submit" isDisabled={m.busy}>{t('Crear proyecto', 'Create project')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...projects} retry={projects.reload}>{!(projects.data || []).length ? <Empty>{t('No hay proyectos.', 'No projects.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Proyecto', 'Project')}</th><th>{t('Estado', 'Status')}</th><th>{t('Presupuesto', 'Budget')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{projects.data!.map(p => <tr key={p.id}><td>{p.name}<small>{(clients.data || []).find(c => c.id === p.client_id)?.name || ''}</small></td><td><Pill value={p.status} /></td><td>{formatCents(p.budget_cents)}</td><td><Button size="sm" variant="secondary" onPress={() => setSelected(p)}>{t('Tareas', 'Tasks')}</Button>{staff && (projectTransitions[p.status] || []).map(next => <Button key={next} size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/projects/${p.id}/status`, json('PATCH', { status: next })); projects.reload(); })}>{next}</Button>)}</td></tr>)}</tbody></table></div>}</LoadState>
      {selected && <section className="panel"><h2>{t('Tareas de', 'Tasks of')} {selected.name} · {t('Avance', 'Progress')}: {progress}%</h2>
        <LoadState {...tasks} retry={tasks.reload}>{!(tasks.data || []).length ? <Empty>{t('Sin tareas.', 'No tasks.')}</Empty> : tasks.data!.map(task => <div className="list-item" key={task.id}><span>{task.title}</span><p>{task.status} · {task.estimated_minutes} {t('min estimados', 'min estimated')}</p>{staff && (taskTransitions[task.status] || []).map(next => <Button key={next} size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/projects/${selected.id}/tasks/${task.id}/status`, json('PATCH', { status: next })); tasks.reload(); })}>{next}</Button>)}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/msp/projects/${selected.id}/tasks`, json('POST', { title: form.get('title'), estimated_minutes: Number(form.get('estimated') || 0) })); tasks.reload(); }, t('Tarea creada.', 'Task created.'))}><div className="form-grid"><Field label={t('Título', 'Title')} name="title" required /><Field label={t('Minutos estimados', 'Estimated minutes')} name="estimated" type="number" min="0" /></div><Button type="submit" isDisabled={m.busy}>{t('Añadir tarea', 'Add task')}</Button></form>}
      </section>}
    </section>
  </>;
}

// ── Rentabilidad y facturación ──────────────────────────────────────

export function Billing({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const drafts = useLoad(() => api<MSPBillingDraft[]>(`/msp/billing-drafts${clientQuery(clientId)}`), [clientId]);
  const profit = useLoad(() => clientId ? api<MSPProfitability>(`/msp/profitability?client_id=${clientId}`) : Promise.resolve(null as unknown as MSPProfitability), [clientId]);
  const staff = user.role !== 'end_user';
  return <>
    <Heading title={t('Rentabilidad y facturación', 'Profitability and billing')} description={t('Resumen por cliente y borradores de facturación. Un borrador nunca cobra nada: solo registra lo que se facturaría.', 'Summary per client and billing drafts. A draft never charges anything: it only records what would be billed.')} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      {clientId && <LoadState {...profit} retry={profit.reload}>{profit.data ? <div className="table-wrap"><table className="table"><tbody>
        <tr><td>{t('Horas imputadas', 'Logged hours')}</td><td>{((profit.data.minutes || 0) / 60).toFixed(1)} h</td></tr>
        <tr><td>{t('Ingresos por tiempo (tarifa de contrato)', 'Time revenue (contract rate)')}</td><td>{formatCents(profit.data.time_revenue_cents)}</td></tr>
        <tr><td>{t('Coste de tiempo (tarifa de técnico)', 'Time cost (technician rate)')}</td><td>{formatCents(profit.data.time_cost_cents)}</td></tr>
        <tr><td>{t('Gastos', 'Expenses')}</td><td>{formatCents(profit.data.expenses_cents)}</td></tr>
        <tr><td>{t('Borradores emitidos', 'Issued drafts')}</td><td>{formatCents(profit.data.billed_drafts_cents)}</td></tr>
        <tr><td><strong>{t('Beneficio', 'Profit')}</strong></td><td><strong>{formatCents(profit.data.profit_cents)}</strong></td></tr>
      </tbody></table></div> : <Empty>{t('Elige un cliente para ver su rentabilidad.', 'Choose a client to see its profitability.')}</Empty>}</LoadState>}
    </section>
    {staff && <section className="panel"><h2>{t('Generar borrador de facturación', 'Generate billing draft')}</h2>
      <form action={form => m.run(async () => { await api('/msp/billing-drafts', json('POST', { client_id: form.get('client'), period_start: form.get('from'), period_end: form.get('to'), notes: form.get('notes') })); drafts.reload(); }, t('Borrador generado. Las horas facturables del periodo quedan marcadas como facturadas.', 'Draft generated. Billable hours of the period are marked as billed.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Field label={t('Desde', 'From')} name="from" type="date" required /><Field label={t('Hasta', 'To')} name="to" type="date" required /></div>
        <Textarea label={t('Notas', 'Notes')} name="notes" rows={2} />
        <div className="form-actions"><Button type="submit" isDisabled={m.busy}>{t('Generar borrador', 'Generate draft')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...drafts} retry={drafts.reload}>{!(drafts.data || []).length ? <Empty>{t('No hay borradores de facturación.', 'No billing drafts.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Periodo', 'Period')}</th><th>{t('Total', 'Total')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{drafts.data!.map(d => <tr key={d.id}><td>{new Date(d.period_start).toLocaleDateString()} – {new Date(d.period_end).toLocaleDateString()}<small>{d.notes || ''}</small></td><td>{formatCents(d.total_cents)}</td><td><Pill value={d.status} /></td><td>{staff && d.status === 'draft' && <><Button size="sm" onPress={() => m.run(async () => { await api(`/msp/billing-drafts/${d.id}/status`, json('PATCH', { status: 'issued' })); drafts.reload(); }, t('Borrador emitido.', 'Draft issued.'))}>{t('Emitir', 'Issue')}</Button><Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/billing-drafts/${d.id}/status`, json('PATCH', { status: 'void' })); drafts.reload(); }, t('Borrador anulado.', 'Draft voided.'))}>{t('Anular', 'Void')}</Button></>}</td></tr>)}</tbody></table></div>}</LoadState></section>
  </>;
}

// ── Inventario por cliente ──────────────────────────────────────────

export function Inventory({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const assets = useLoad(() => api<Asset[]>('/assets'));
  const types = useLoad(() => api<CIType[]>('/ci-types'));
  const clientAssets = useLoad(() => clientId ? api<MSPClientAsset[]>(`/msp/clients/${clientId}/assets`) : Promise.resolve([] as MSPClientAsset[]), [clientId]);
  const discovered = useLoad(() => clientId ? api<MSPDiscoveredAsset[]>(`/msp/clients/${clientId}/discovered-assets`) : Promise.resolve([] as MSPDiscoveredAsset[]), [clientId]);
  return <>
    <Heading title={t('Inventario por cliente', 'Inventory per client')} description={t('Activos de la Etapa 3 asignados a cada cliente, más descubrimientos pendientes de importar.', 'Stage 3 assets assigned to each client, plus discoveries pending import.')} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      {!clientId && <Empty>{t('Elige un cliente para ver su inventario.', 'Choose a client to see its inventory.')}</Empty>}
      {clientId && <><LoadState {...clientAssets} retry={clientAssets.reload}>{!(clientAssets.data || []).length ? <Empty>{t('Este cliente no tiene activos asignados.', 'This client has no assigned assets.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Activo', 'Asset')}</th><th>{t('Tipo', 'Type')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{clientAssets.data!.map(a => <tr key={a.asset_id}><td>{a.name}<small>{a.serial_number || ''}</small></td><td>{a.ci_type_name}</td><td><Pill value={a.status} /></td><td>{staff && <Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/msp/clients/${clientId}/assets/${a.asset_id}`, { method: 'DELETE' }); clientAssets.reload(); }, t('Activo desvinculado.', 'Asset unlinked.'))}>{t('Desvincular', 'Unlink')}</Button>}</td></tr>)}</tbody></table></div>}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/msp/clients/${clientId}/assets`, json('POST', { asset_id: form.get('asset') })); clientAssets.reload(); }, t('Activo vinculado al cliente.', 'Asset linked to the client.'))}><div className="form-grid"><Select label={t('Activo existente', 'Existing asset')} name="asset" required>{(assets.data || []).map(a => <option key={a.id} value={a.id}>{a.name}</option>)}</Select></div><Button type="submit" isDisabled={m.busy}>{t('Vincular activo', 'Link asset')}</Button></form>}
        <h2 className="mt-6">{t('Descubrimiento e importación', 'Discovery and import')}</h2>
        <p>{t('El descubrimiento es manual o por importación documentada: registra aquí lo hallado (por ejemplo, desde Zabbix) y luego impórtalo como activo del cliente.', 'Discovery is manual or documented import: record here what was found (for example, from Zabbix) and then import it as a client asset.')}</p>
        <LoadState {...discovered} retry={discovered.reload}>{!(discovered.data || []).length ? <Empty>{t('Sin descubrimientos registrados.', 'No discoveries registered.')}</Empty> : discovered.data!.map(d => <div className="list-item" key={d.id}><span>{d.name}</span><p>{d.kind} · {d.external_id} · {d.status}</p>{staff && d.status === 'pending' && <form action={form => m.run(async () => { await api(`/msp/clients/${clientId}/discovered-assets/${d.id}/import`, json('POST', { ci_type_id: form.get('type'), asset_name: form.get('name') })); discovered.reload(); clientAssets.reload(); }, t('Activo importado.', 'Asset imported.'))}><div className="form-grid"><Select label={t('Tipo de CI', 'CI type')} name="type" required>{(types.data || []).map(ct => <option key={ct.id} value={ct.id}>{ct.name}</option>)}</Select><Field label={t('Nombre del activo', 'Asset name')} name="name" required defaultValue={d.name} /></div><Button size="sm" type="submit" isDisabled={m.busy}>{t('Importar', 'Import')}</Button></form>}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/msp/clients/${clientId}/discovered-assets`, json('POST', { external_id: form.get('external'), name: form.get('name'), kind: form.get('kind') })); discovered.reload(); }, t('Descubrimiento registrado.', 'Discovery registered.'))}><div className="form-grid"><Field label={t('ID externo', 'External ID')} name="external" required /><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('Tipo', 'Kind')} name="kind" defaultValue="host" /></div><Button type="submit" isDisabled={m.busy}>{t('Registrar descubrimiento', 'Register discovery')}</Button></form>}
      </>}
    </section>
  </>;
}

// ── Monitorización y acceso remoto ────────────────────────────────

export function Monitoring({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const sources = useLoad(() => api<MSPMonitoringSource[]>('/msp/monitoring/sources'));
  const alerts = useLoad(() => api<MSPAlert[]>(`/msp/monitoring/alerts${clientQuery(clientId)}`), [clientId]);
  return <>
    <Heading title={t('Monitorización (Zabbix)', 'Monitoring (Zabbix)')} description={t('Fuentes configurables y alertas que crean ticket sin duplicar la alerta abierta del mismo recurso.', 'Configurable sources and alerts that create a ticket without duplicating the open alert for the same resource.')} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      <p>{t('Activación: crea una fuente con su token de webhook y configura en Zabbix un webhook hacia POST /api/v1/msp/monitoring/webhook con las cabeceras X-Tenant-ID (slug de la organización) y X-MSP-Webhook-Token. Sin instancias ni credenciales reales, la validación viva queda pendiente (ver docs/msp-integraciones.md).', 'Activation: create a source with its webhook token and configure a Zabbix webhook to POST /api/v1/msp/monitoring/webhook with X-Tenant-ID (organization slug) and X-MSP-Webhook-Token headers. Without real instances or credentials, live validation is pending (see docs/msp-integraciones.md).')}</p>
    </section>
    {staff && <section className="panel"><h2>{t('Nueva fuente Zabbix', 'New Zabbix source')}</h2>
      <form action={form => m.run(async () => { await api('/msp/monitoring/sources', json('POST', { name: form.get('name'), base_url: form.get('url'), webhook_token: form.get('token') })); sources.reload(); }, t('Fuente creada. Guarda el token fuera de Trickreport: solo se conserva su hash.', 'Source created. Keep the token outside Trickreport: only its hash is stored.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required /><Field label={t('URL base de Zabbix', 'Zabbix base URL')} name="url" /><Field label={t('Token del webhook', 'Webhook token')} name="token" type="password" /></div>
        <div className="form-actions"><Button type="submit" isDisabled={m.busy}>{t('Crear fuente', 'Create source')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...sources} retry={sources.reload}>{!(sources.data || []).length ? <Empty>{t('Sin fuentes de monitorización.', 'No monitoring sources.')}</Empty> : sources.data!.map(s => <div className="list-item" key={s.id}><span>{s.name}</span><p>{s.kind} · {s.base_url || '—'} · {s.configured ? t('configurada', 'configured') : t('sin token', 'no token')}</p></div>)}</LoadState>
      <h2 className="mt-6">{t('Alertas', 'Alerts')}</h2>
      <LoadState {...alerts} retry={alerts.reload}>{!(alerts.data || []).length ? <Empty>{t('Sin alertas.', 'No alerts.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Recurso', 'Resource')}</th><th>{t('Gravedad', 'Severity')}</th><th>{t('Estado', 'Status')}</th><th>{t('Ticket', 'Ticket')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{alerts.data!.map(a => <tr key={a.id}><td>{a.external_resource}<small>{a.title || a.message}</small></td><td>{a.severity}</td><td><Pill value={a.status} /></td><td>{a.ticket_id ? <Link href={`/tickets/${a.ticket_id}`}>#{a.ticket_id.slice(0, 8)}</Link> : '—'}</td><td>{staff && a.status === 'open' && <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/monitoring/alerts/${a.id}/resolve`, { method: 'POST' }); alerts.reload(); }, t('Alerta resuelta.', 'Alert resolved.'))}>{t('Resolver', 'Resolve')}</Button>}</td></tr>)}</tbody></table></div>}</LoadState>
    </section>
  </>;
}

export function RemoteSessions({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const { clientId, change } = useMSPClientFilter();
  const clients = useLoad(() => api<MSPClient[]>('/msp/clients'));
  const techs = useLoad(() => api<MSPTechnician[]>('/msp/technicians'));
  const sessions = useLoad(() => api<MSPRemoteSession[]>(`/msp/remote-sessions${clientQuery(clientId)}`), [clientId]);
  return <>
    <Heading title={t('Acceso remoto (MeshCentral)', 'Remote access (MeshCentral)')} description={t('Sesiones con permiso, autorización y trazabilidad desde Trickreport. Trickreport registra la sesión; MeshCentral la ejecuta.', 'Sessions with permission, authorization and traceability from Trickreport. Trickreport records the session; MeshCentral runs it.')} />
    <Feedback {...m} />
    <section className="panel"><div className="form-grid"><ClientFilter clients={clients.data || []} value={clientId} onChange={change} /></div>
      <p>{t('Activación: enlaza el ID de sesión de MeshCentral al solicitarla aquí; la sesión solo puede iniciarse tras la autorización del personal. Sin instancia real, la validación viva queda pendiente (ver docs/msp-integraciones.md).', 'Activation: link the MeshCentral session ID when requesting it here; the session can only start after staff authorization. Without a real instance, live validation is pending (see docs/msp-integraciones.md).')}</p>
    </section>
    <section className="panel"><h2>{t('Solicitar sesión', 'Request session')}</h2>
      <form action={form => m.run(async () => { await api('/msp/remote-sessions', json('POST', { client_id: form.get('client'), technician_id: form.get('tech'), external_session_id: form.get('external'), audit_note: form.get('note') })); sessions.reload(); }, t('Sesión solicitada. Pendiente de autorización del personal.', 'Session requested. Awaiting staff authorization.'))}>
        <div className="form-grid"><Select label={t('Cliente', 'Client')} name="client" required defaultValue={clientId}>{(clients.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select><Select label={t('Técnico', 'Technician')} name="tech" required>{(techs.data || []).map(tech => <option key={tech.id} value={tech.id}>{tech.user_name}</option>)}</Select><Field label={t('ID de sesión MeshCentral', 'MeshCentral session ID')} name="external" /><Field label={t('Nota de auditoría', 'Audit note')} name="note" /></div>
        <div className="form-actions"><Button type="submit" isDisabled={m.busy}>{t('Solicitar sesión', 'Request session')}</Button></div>
      </form></section>
    <section className="panel"><LoadState {...sessions} retry={sessions.reload}>{!(sessions.data || []).length ? <Empty>{t('No hay sesiones remotas.', 'No remote sessions.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Sesión', 'Session')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{sessions.data!.map(s => <tr key={s.id}><td>{s.external_session_id || s.id.slice(0, 8)}<small>{s.audit_note || s.provider}</small></td><td><Pill value={s.status} /></td><td>{staff && s.status === 'requested' && <><Button size="sm" onPress={() => m.run(async () => { await api(`/msp/remote-sessions/${s.id}/authorize`, json('POST', { authorize: true })); sessions.reload(); }, t('Sesión autorizada.', 'Session authorized.'))}>{t('Autorizar', 'Authorize')}</Button><Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/remote-sessions/${s.id}/authorize`, json('POST', { authorize: false })); sessions.reload(); }, t('Sesión denegada.', 'Session denied.'))}>{t('Denegar', 'Deny')}</Button></>}{s.status === 'authorized' && <Button size="sm" onPress={() => m.run(async () => { await api(`/msp/remote-sessions/${s.id}/start`, { method: 'POST' }); sessions.reload(); }, t('Sesión iniciada.', 'Session started.'))}>{t('Iniciar', 'Start')}</Button>}{staff && s.status === 'active' && <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/msp/remote-sessions/${s.id}/end`, { method: 'POST' }); sessions.reload(); }, t('Sesión finalizada.', 'Session ended.'))}>{t('Finalizar', 'End')}</Button>}</td></tr>)}</tbody></table></div>}</LoadState></section>
  </>;
}

export function MSPReportView() {
  const t = useText();
  const report = useLoad(() => api<MSPReport>('/admin/msp'));
  return <>
    <Heading title={t('Informes MSP', 'MSP reports')} description={t('Operación multicliente de un vistazo.', 'Multi-client operation at a glance.')} />
    <section className="panel"><LoadState {...report} retry={report.reload}>{report.data ? <div className="table-wrap"><table className="table"><tbody>
      <tr><td>{t('Clientes totales', 'Total clients')}</td><td>{report.data.clients_total}</td></tr>
      <tr><td>{t('Clientes activos', 'Active clients')}</td><td>{report.data.active_clients}</td></tr>
      <tr><td>{t('Contratos activos', 'Active contracts')}</td><td>{report.data.active_contracts}</td></tr>
      <tr><td>{t('Alertas abiertas', 'Open alerts')}</td><td>{report.data.open_alerts}</td></tr>
      <tr><td>{t('Sesiones activas o autorizadas', 'Active or authorized sessions')}</td><td>{report.data.active_sessions}</td></tr>
      <tr><td>{t('Tickets MSP abiertos', 'Open MSP tickets')}</td><td>{report.data.open_msp_tickets}</td></tr>
    </tbody></table></div> : <Empty>{t('Sin datos.', 'No data.')}</Empty>}</LoadState></section>
  </>;
}
