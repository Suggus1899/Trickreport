'use client';
import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { api, json } from './api';
import type { Approval, Asset, CatalogCategory, CatalogService, Change, CIType, CIRelationship, ITSMReport, Problem, Ticket, User } from './contracts';
import { changeTransitions, problemTransitions } from './contracts';
import { Button, Empty, Feedback, Field, Heading, LoadState, Pill, Select, Textarea, useLoad, useMutation, useText } from './ui';

// ITSM (Etapa 3): catálogo de servicios, aprobaciones, activos/CMDB,
// problemas, cambios e informes. Todas las pantallas consumen la API
// real (/api/v1); el servidor Go sigue siendo la autoridad de permisos
// (la ocultación de un botón aquí nunca sustituye esa comprobación).

export function Catalog({ user }: { user: User }) {
  const t = useText(), m = useMutation(), router = useRouter();
  const staff = user.role !== 'end_user';
  const cats = useLoad(() => api<CatalogCategory[]>('/catalog/categories'));
  const svcs = useLoad(() => api<CatalogService[]>('/catalog/services'));
  const [editingCat, setEditingCat] = useState<CatalogCategory | null>();
  const [editingSvc, setEditingSvc] = useState<CatalogService | null>();
  const [requesting, setRequesting] = useState<CatalogService | null>();
  return <>
    <Heading title={t('Catálogo de servicios', 'Service catalog')} description={t('Servicios que puedes solicitar y su estado.', 'Services you can request and their status.')} action={staff ? <Button onPress={() => setEditingSvc(null)}>{t('Nuevo servicio', 'New service')}</Button> : undefined} />
    <Feedback {...m} />
    {staff && editingSvc !== undefined && <section className="panel"><h2>{editingSvc ? t('Editar servicio', 'Edit service') : t('Nuevo servicio', 'New service')}</h2>
      <form key={editingSvc?.id || 'new'} action={form => m.run(async () => { await api(editingSvc ? `/catalog/services/${editingSvc.id}` : '/catalog/services', json(editingSvc ? 'PUT' : 'POST', { category_id: form.get('category') || null, name: form.get('name'), description: form.get('description'), requires_approval: form.get('approval') === 'on', active: form.get('active') === 'on' })); setEditingSvc(undefined); svcs.reload(); }, t('Servicio guardado.', 'Service saved.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required defaultValue={editingSvc?.name} /><Select label={t('Categoría', 'Category')} name="category" defaultValue={editingSvc?.category_id || ''}><option value="">{t('Sin categoría', 'No category')}</option>{(cats.data || []).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</Select></div>
        <Textarea label={t('Descripción', 'Description')} name="description" defaultValue={editingSvc?.description} />
        <label className="checkbox"><input type="checkbox" name="approval" defaultChecked={editingSvc?.requires_approval} />{t('Requiere aprobación', 'Requires approval')}</label>
        <label className="checkbox"><input type="checkbox" name="active" defaultChecked={editingSvc?.active ?? true} />{t('Activo en el catálogo', 'Active in the catalog')}</label>
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditingSvc(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    {staff && editingCat !== undefined && <section className="panel"><h2>{editingCat ? t('Editar categoría', 'Edit category') : t('Nueva categoría', 'New category')}</h2>
      <form key={editingCat?.id || 'new'} action={form => m.run(async () => { await api(editingCat ? `/catalog/categories/${editingCat.id}` : '/catalog/categories', json(editingCat ? 'PUT' : 'POST', { name: form.get('name'), description: form.get('description') })); setEditingCat(undefined); cats.reload(); }, t('Categoría guardada.', 'Category saved.'))}>
        <Field label={t('Nombre', 'Name')} name="name" required defaultValue={editingCat?.name} /><Field label={t('Descripción', 'Description')} name="description" defaultValue={editingCat?.description} />
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditingCat(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    {requesting && <section className="panel"><h2>{t('Solicitar servicio', 'Request service')}: {requesting.name}</h2>
      <form action={form => m.run(async () => { const res = await api<{ ticket_id: string }>(`/catalog/services/${requesting.id}/request`, json('POST', { title: form.get('title'), description: form.get('description'), priority: form.get('priority') })); setRequesting(null); router.push(`/tickets/${res.ticket_id}`); }, t('Solicitud creada.', 'Request created.'))}>
        <Field label={t('Título', 'Title')} name="title" required defaultValue={requesting.name} /><Textarea label={t('Descripción', 'Description')} name="description" required rows={4} />
        <Select label={t('Prioridad', 'Priority')} name="priority" defaultValue="medium"><option value="low">{t('Baja', 'Low')}</option><option value="medium">{t('Media', 'Medium')}</option><option value="high">{t('Alta', 'High')}</option><option value="critical">{t('Crítica', 'Critical')}</option></Select>
        <div className="form-actions"><Button variant="secondary" onPress={() => setRequesting(null)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Crear solicitud', 'Create request')}</Button></div>
      </form></section>}
    <section className="panel"><div className="section-header"><h2>{t('Servicios', 'Services')}</h2>{staff && <Button variant="secondary" onPress={() => setEditingCat(null)}>{t('Nueva categoría', 'New category')}</Button>}</div>
      <LoadState {...svcs} retry={svcs.reload}>{!(svcs.data || []).length ? <Empty>{t('No hay servicios en el catálogo.', 'No services in the catalog.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Servicio', 'Service')}</th><th>{t('Categoría', 'Category')}</th><th>{t('Aprobación', 'Approval')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{svcs.data!.map(s => <tr key={s.id}><td>{s.name}<small>{s.description}</small></td><td>{s.category_name || '—'}</td><td>{s.requires_approval ? t('Requerida', 'Required') : t('No requerida', 'Not required')}</td><td><Button size="sm" variant="secondary" isDisabled={!s.active} onPress={() => setRequesting(s)}>{t('Solicitar', 'Request')}</Button>{staff && <Button size="sm" variant="secondary" onPress={() => setEditingSvc(s)}>{t('Editar', 'Edit')}</Button>}</td></tr>)}</tbody></table></div>}</LoadState>
      <h2 className="mt-6">{t('Categorías', 'Categories')}</h2>
      <LoadState {...cats} retry={cats.reload}><div className="list">{(cats.data || []).map(c => <div className="list-item" key={c.id}><span>{c.name}</span><p>{c.description}</p>{staff && <Button size="sm" variant="secondary" onPress={() => setEditingCat(c)}>{t('Editar', 'Edit')}</Button>}</div>)}</div></LoadState>
    </section>
  </>;
}

export function Approvals({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const state = useLoad(() => api<Approval[]>('/approvals'));
  const [history, setHistory] = useState<{ id: string; events: ApprovalEvent[] } | null>(null);
  void user;
  return <>
    <Heading title={t('Aprobaciones', 'Approvals')} description={t('Decide las solicitudes pendientes y consulta su historial.', 'Decide pending requests and review their history.')} />
    <Feedback {...m} />
    <section className="panel"><LoadState {...state} retry={state.reload}>{!(state.data || []).length ? <Empty>{t('No hay aprobaciones.', 'No approvals.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Destino', 'Target')}</th><th>{t('Estado', 'Status')}</th><th>{t('Comentario', 'Comment')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{state.data!.map(a => <tr key={a.id}><td>{a.target_type}<small className="mono">#{a.target_id.slice(0, 8)}</small></td><td><Pill value={a.status} /></td><td>{a.comment || '—'}</td><td>{a.status === 'pending' && <><Button size="sm" onPress={() => m.run(async () => { await api(`/approvals/${a.id}/decide`, json('POST', { decision: 'approved', comment: '' })); state.reload(); }, t('Aprobación aceptada.', 'Approval accepted.'))}>{t('Aprobar', 'Approve')}</Button><Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/approvals/${a.id}/decide`, json('POST', { decision: 'rejected', comment: '' })); state.reload(); }, t('Aprobación rechazada.', 'Approval rejected.'))}>{t('Rechazar', 'Reject')}</Button></>}<Button size="sm" variant="secondary" onPress={async () => { const events = await api<ApprovalEvent[]>(`/approvals/${a.id}/history`); setHistory({ id: a.id, events }); }}>{t('Historial', 'History')}</Button></td></tr>)}</tbody></table></div>}</LoadState>
      {history && <section className="panel"><h2>{t('Historial de aprobación', 'Approval history')}</h2>{!history.events.length ? <Empty>{t('Sin decisiones registradas.', 'No decisions recorded.')}</Empty> : history.events.map(e => <div className="list-item" key={e.id}><span>{e.from_status} → {e.to_status}</span><p>{e.comment || '—'} · {new Date(e.created_at).toLocaleString()}</p></div>)}<Button variant="secondary" onPress={() => setHistory(null)}>{t('Cerrar', 'Close')}</Button></section>}
    </section>
  </>;
}

type ApprovalEvent = { id: string; from_status: string; to_status: string; comment: string; created_at: string };

export function Assets({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const types = useLoad(() => api<CIType[]>('/ci-types'));
  const assets = useLoad(() => api<Asset[]>('/assets'));
  const [editing, setEditing] = useState<Asset | null>();
  const [selected, setSelected] = useState<Asset | null>(null);
  const rels = useLoad(() => selected ? api<CIRelationship[]>(`/assets/${selected.id}/relationships`) : Promise.resolve([] as CIRelationship[]), [selected?.id]);
  return <>
    <Heading title={t('Activos y CMDB', 'Assets and CMDB')} description={t('Elementos de configuración, sus tipos y relaciones.', 'Configuration items, their types and relationships.')} action={staff ? <Button onPress={() => setEditing(null)}>{t('Nuevo activo', 'New asset')}</Button> : undefined} />
    <Feedback {...m} />
    {staff && editing !== undefined && <section className="panel"><h2>{editing ? t('Editar activo', 'Edit asset') : t('Nuevo activo', 'New asset')}</h2>
      <form key={editing?.id || 'new'} action={form => m.run(async () => { await api(editing ? `/assets/${editing.id}` : '/assets', json(editing ? 'PUT' : 'POST', { ci_type_id: form.get('ci_type'), name: form.get('name'), serial_number: form.get('serial'), status: form.get('status') })); setEditing(undefined); assets.reload(); }, t('Activo guardado.', 'Asset saved.'))}>
        <div className="form-grid"><Field label={t('Nombre', 'Name')} name="name" required defaultValue={editing?.name} /><Select label={t('Tipo de CI', 'CI type')} name="ci_type" required defaultValue={editing?.ci_type_id || ''}><option value="" disabled>{t('Selecciona un tipo', 'Choose a type')}</option>{(types.data || []).map(ct => <option key={ct.id} value={ct.id}>{ct.name}</option>)}</Select><Field label={t('Número de serie', 'Serial number')} name="serial" defaultValue={editing?.serial_number} /><Select label={t('Estado', 'Status')} name="status" defaultValue={editing?.status || 'active'}><option value="active">{t('Activo', 'Active')}</option><option value="maintenance">{t('En mantenimiento', 'Maintenance')}</option><option value="retired">{t('Retirado', 'Retired')}</option><option value="stock">{t('En almacén', 'In stock')}</option></Select></div>
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditing(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...assets} retry={assets.reload}>{!(assets.data || []).length ? <Empty>{t('No hay activos registrados.', 'No assets registered.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Activo', 'Asset')}</th><th>{t('Tipo', 'Type')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{assets.data!.map(a => <tr key={a.id}><td>{a.name}<small>{a.serial_number || a.ci_type_name || ''}</small></td><td>{a.ci_type_name || '—'}</td><td><Pill value={a.status} /></td><td><Button size="sm" variant="secondary" onPress={() => setSelected(a)}>{t('Relaciones', 'Relationships')}</Button>{staff && <Button size="sm" variant="secondary" onPress={() => setEditing(a)}>{t('Editar', 'Edit')}</Button>}</td></tr>)}</tbody></table></div>}</LoadState>
      {staff && <form className="mt-6" action={form => m.run(async () => { await api('/ci-types', json('POST', { name: form.get('name'), description: '' })); types.reload(); }, t('Tipo de CI creado.', 'CI type created.'))}><div className="form-grid"><Field label={t('Nuevo tipo de CI', 'New CI type')} name="name" required /><Button type="submit" isDisabled={m.busy}>{t('Crear tipo', 'Create type')}</Button></div></form>}
      {selected && <section className="panel"><h2>{t('Relaciones de', 'Relationships of')} {selected.name}</h2>
        <LoadState {...rels} retry={rels.reload}>{!(rels.data || []).length ? <Empty>{t('Sin relaciones.', 'No relationships.')}</Empty> : rels.data!.map(r => <div className="list-item" key={r.id}><span>{r.source_ci_id === selected.id ? '→' : '←'} {r.relationship_type}</span><p className="mono">{r.source_ci_id === selected.id ? r.target_ci_id.slice(0, 8) : r.source_ci_id.slice(0, 8)}</p>{staff && <Button size="sm" variant="danger" onPress={() => m.run(async () => { await api(`/assets/relationships/${r.id}`, { method: 'DELETE' }); rels.reload(); })}>{t('Eliminar', 'Delete')}</Button>}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api('/assets/relationships', json('POST', { source_ci_id: form.get('source'), target_ci_id: form.get('target'), relationship_type: form.get('type') })); rels.reload(); }, t('Relación creada.', 'Relationship created.'))}><div className="form-grid"><Select label={t('Origen', 'Source')} name="source" required defaultValue={selected.id}>{(assets.data || []).map(a => <option key={a.id} value={a.id}>{a.name}</option>)}</Select><Select label={t('Destino', 'Target')} name="target" required>{(assets.data || []).map(a => <option key={a.id} value={a.id}>{a.name}</option>)}</Select><Select label={t('Relación', 'Relationship')} name="type" defaultValue="depends_on"><option value="depends_on">{t('Depende de', 'Depends on')}</option><option value="hosts">{t('Hospeda', 'Hosts')}</option><option value="uses">{t('Usa', 'Uses')}</option></Select></div><Button type="submit" isDisabled={m.busy}>{t('Añadir relación', 'Add relationship')}</Button></form>}
      </section>}
    </section>
  </>;
}

export function Problems({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const state = useLoad(() => api<Problem[]>('/problems'));
  const tickets = useLoad(() => api<Ticket[]>('/tickets'));
  const [editing, setEditing] = useState<Problem | null>();
  const [selected, setSelected] = useState<Problem | null>(null);
  const incidents = useLoad(() => selected ? api<string[]>(`/problems/${selected.id}/incidents`) : Promise.resolve([] as string[]), [selected?.id]);
  return <>
    <Heading title={t('Problemas', 'Problems')} description={t('Averías recurrentes, causa raíz y solución temporal.', 'Recurring faults, root cause and workaround.')} action={staff ? <Button onPress={() => setEditing(null)}>{t('Nuevo problema', 'New problem')}</Button> : undefined} />
    <Feedback {...m} />
    {staff && editing !== undefined && <section className="panel"><h2>{editing ? t('Editar problema', 'Edit problem') : t('Nuevo problema', 'New problem')}</h2>
      <form key={editing?.id || 'new'} action={form => m.run(async () => { await api(editing ? `/problems/${editing.id}` : '/problems', json(editing ? 'PATCH' : 'POST', editing ? { status: editing.status, root_cause: form.get('root'), workaround: form.get('workaround') } : { title: form.get('title'), description: form.get('description'), root_cause: form.get('root'), workaround: form.get('workaround') })); setEditing(undefined); state.reload(); }, t('Problema guardado.', 'Problem saved.'))}>
        {!editing && <><Field label={t('Título', 'Title')} name="title" required /><Textarea label={t('Descripción', 'Description')} name="description" rows={4} /></>}
        <Textarea label={t('Causa raíz', 'Root cause')} name="root" defaultValue={editing?.root_cause} rows={3} /><Textarea label={t('Solución temporal', 'Workaround')} name="workaround" defaultValue={editing?.workaround} rows={3} />
        <div className="form-actions"><Button variant="secondary" onPress={() => setEditing(undefined)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...state} retry={state.reload}>{!(state.data || []).length ? <Empty>{t('No hay problemas registrados.', 'No problems registered.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Problema', 'Problem')}</th><th>{t('Estado', 'Status')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{state.data!.map(p => <tr key={p.id}><td>{p.title}<small>{p.root_cause || p.description}</small></td><td><Pill value={p.status} /></td><td><Button size="sm" variant="secondary" onPress={() => setSelected(p)}>{t('Incidentes', 'Incidents')}</Button>{staff && <><Button size="sm" variant="secondary" onPress={() => setEditing(p)}>{t('Causa y solución', 'Cause and workaround')}</Button>{(problemTransitions[p.status] || []).map(next => <Button key={next} size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/problems/${p.id}`, json('PATCH', { status: next })); state.reload(); })}>{next}</Button>)}</>}</td></tr>)}</tbody></table></div>}</LoadState>
      {selected && <section className="panel"><h2>{t('Incidentes vinculados', 'Linked incidents')}: {selected.title}</h2>
        <LoadState {...incidents} retry={incidents.reload}>{!(incidents.data || []).length ? <Empty>{t('Sin incidentes vinculados.', 'No linked incidents.')}</Empty> : incidents.data!.map(id => <div className="list-item" key={id}><Link href={`/tickets/${id}`}>#{id.slice(0, 8)}</Link>{staff && <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/problems/${selected.id}/incidents/${id}`, { method: 'DELETE' }); incidents.reload(); })}>{t('Desvincular', 'Unlink')}</Button>}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/problems/${selected.id}/incidents`, json('POST', { ticket_id: form.get('ticket') })); incidents.reload(); }, t('Incidente vinculado.', 'Incident linked.'))}><div className="form-grid"><Select label={t('Incidente', 'Incident')} name="ticket" required>{(tickets.data || []).map(tk => <option key={tk.id} value={tk.id}>{tk.title}</option>)}</Select><Button type="submit" isDisabled={m.busy}>{t('Vincular', 'Link')}</Button></div></form>}
      </section>}
    </section>
  </>;
}

export function Changes({ user }: { user: User }) {
  const t = useText(), m = useMutation();
  const staff = user.role !== 'end_user';
  const state = useLoad(() => api<Change[]>('/changes'));
  const problems = useLoad(() => api<Problem[]>('/problems'));
  const tickets = useLoad(() => api<Ticket[]>('/tickets'));
  const [creating, setCreating] = useState(false);
  const [selected, setSelected] = useState<Change | null>(null);
  const linked = useLoad(() => selected ? api<string[]>(`/changes/${selected.id}/tickets`) : Promise.resolve([] as string[]), [selected?.id]);
  return <>
    <Heading title={t('Cambios', 'Changes')} description={t('Planificación, aprobación, ventana, riesgo e implementación.', 'Planning, approval, window, risk and implementation.')} action={staff ? <Button onPress={() => setCreating(true)}>{t('Nuevo cambio', 'New change')}</Button> : undefined} />
    <Feedback {...m} />
    {staff && creating && <section className="panel"><h2>{t('Nuevo cambio', 'New change')}</h2>
      <form action={form => m.run(async () => { await api('/changes', json('POST', { title: form.get('title'), description: form.get('description'), problem_id: form.get('problem') || null, risk: form.get('risk'), planned_start: form.get('start') ? new Date(String(form.get('start'))).toISOString() : '', planned_end: form.get('end') ? new Date(String(form.get('end'))).toISOString() : '' })); setCreating(false); state.reload(); }, t('Cambio creado.', 'Change created.'))}>
        <Field label={t('Título', 'Title')} name="title" required /><Textarea label={t('Descripción', 'Description')} name="description" rows={4} />
        <div className="form-grid"><Select label={t('Problema relacionado', 'Related problem')} name="problem" defaultValue=""><option value="">{t('Ninguno', 'None')}</option>{(problems.data || []).map(p => <option key={p.id} value={p.id}>{p.title}</option>)}</Select><Select label={t('Riesgo', 'Risk')} name="risk" defaultValue="medium"><option value="low">{t('Bajo', 'Low')}</option><option value="medium">{t('Medio', 'Medium')}</option><option value="high">{t('Alto', 'High')}</option><option value="critical">{t('Crítico', 'Critical')}</option></Select><Field label={t('Inicio de ventana', 'Window start')} name="start" type="datetime-local" /><Field label={t('Fin de ventana', 'Window end')} name="end" type="datetime-local" /></div>
        <div className="form-actions"><Button variant="secondary" onPress={() => setCreating(false)}>{t('Cancelar', 'Cancel')}</Button><Button type="submit" isDisabled={m.busy}>{t('Guardar', 'Save')}</Button></div>
      </form></section>}
    <section className="panel"><LoadState {...state} retry={state.reload}>{!(state.data || []).length ? <Empty>{t('No hay cambios registrados.', 'No changes registered.')}</Empty> : <div className="table-wrap"><table className="table"><thead><tr><th>{t('Cambio', 'Change')}</th><th>{t('Estado', 'Status')}</th><th>{t('Riesgo', 'Risk')}</th><th>{t('Ventana', 'Window')}</th><th>{t('Acciones', 'Actions')}</th></tr></thead><tbody>{state.data!.map(c => <tr key={c.id}><td>{c.title}<small>{c.description}</small></td><td><Pill value={c.status} /></td><td><Pill value={c.risk} /></td><td>{c.planned_start ? new Date(c.planned_start).toLocaleString() : '—'}</td><td><Button size="sm" variant="secondary" onPress={() => setSelected(c)}>{t('Solicitudes', 'Requests')}</Button>{staff && (changeTransitions[c.status] || []).map(next => <Button key={next} size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/changes/${c.id}/status`, json('PATCH', { status: next })); state.reload(); }, t('Estado actualizado.', 'Status updated.'))}>{next}</Button>)}</td></tr>)}</tbody></table></div>}</LoadState>
      {selected && <p className="search-caption">{t('Pasar a «aprobado» exige una aprobación aceptada para el cambio; al pasar a «awaiting_approval» se crea una pendiente si no existe.', 'Moving to “approved” requires an accepted approval for the change; moving to “awaiting_approval” creates a pending one if none exists.')}</p>}
      {selected && <section className="panel"><h2>{t('Solicitudes vinculadas', 'Linked requests')}: {selected.title}</h2>
        <LoadState {...linked} retry={linked.reload}>{!(linked.data || []).length ? <Empty>{t('Sin solicitudes vinculadas.', 'No linked requests.')}</Empty> : linked.data!.map(id => <div className="list-item" key={id}><Link href={`/tickets/${id}`}>#{id.slice(0, 8)}</Link>{staff && <Button size="sm" variant="secondary" onPress={() => m.run(async () => { await api(`/changes/${selected.id}/tickets/${id}`, { method: 'DELETE' }); linked.reload(); })}>{t('Desvincular', 'Unlink')}</Button>}</div>)}</LoadState>
        {staff && <form className="mt-6" action={form => m.run(async () => { await api(`/changes/${selected.id}/tickets`, json('POST', { ticket_id: form.get('ticket') })); linked.reload(); }, t('Solicitud vinculada.', 'Request linked.'))}><div className="form-grid"><Select label={t('Solicitud', 'Request')} name="ticket" required>{(tickets.data || []).map(tk => <option key={tk.id} value={tk.id}>{tk.title}</option>)}</Select><Button type="submit" isDisabled={m.busy}>{t('Vincular', 'Link')}</Button></div></form>}
      </section>}
    </section>
  </>;
}

export function ITSMReportView() {
  const t = useText();
  const state = useLoad(() => api<ITSMReport>('/admin/analytics/itsm'));
  return <>
    <Heading title={t('Informes ITSM', 'ITSM reports')} description={t('Cambios, problemas, activos y aprobaciones de tu organización.', 'Changes, problems, assets and approvals in your organization.')} />
    <LoadState {...state} retry={state.reload}>{state.data && <>
      <div className="stats"><section className="stat"><div className="stat-top">{t('Aprobaciones pendientes', 'Pending approvals')}</div><strong>{state.data.pending_approvals}</strong></section><section className="stat"><div className="stat-top">{t('Servicios activos', 'Active services')}</div><strong>{state.data.active_services}</strong></section></div>
      <div className="two-columns">
        <section className="panel"><h2>{t('Cambios por estado', 'Changes by status')}</h2>{state.data.changes_by_status.map(row => <div className="list-item" key={row.status}><Pill value={row.status} /><strong>{row.count}</strong></div>)}</section>
        <section className="panel"><h2>{t('Problemas por estado', 'Problems by status')}</h2>{state.data.problems_by_status.map(row => <div className="list-item" key={row.status}><Pill value={row.status} /><strong>{row.count}</strong></div>)}</section>
        <section className="panel"><h2>{t('Activos por tipo', 'Assets by type')}</h2>{state.data.assets_by_type.map(row => <div className="list-item" key={row.type}><span>{row.type || '—'}</span><strong>{row.count}</strong></div>)}</section>
        <section className="panel"><h2>{t('Aprobaciones por estado', 'Approvals by status')}</h2>{state.data.approvals_by_status.map(row => <div className="list-item" key={row.status}><Pill value={row.status} /><strong>{row.count}</strong></div>)}</section>
      </div>
    </>}</LoadState>
  </>;
}
