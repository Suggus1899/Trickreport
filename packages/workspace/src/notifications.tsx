'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Bell } from 'lucide-react';
import { api } from './api';
import type { Notification } from './contracts';
import { Button, Empty, Notice, useText } from './ui';

export function Notifications({ tenant }: { tenant: string }) {
  const t = useText();
  const [open, setOpen] = useState(false);
  const [unread, setUnread] = useState(0);
  const [items, setItems] = useState<Notification[]>();
  const [error, setError] = useState('');

  useEffect(() => {
    let active = true;
    const load = async () => {
      try {
        const result = await api<{ unread: number }>('/notifications/unread-count', {}, tenant);
        if (active) setUnread(result.unread);
      } catch { /* El contador es informativo; la sesión sigue siendo válida. */ }
    };
    void load();
    const timer = setInterval(() => { void load(); }, 60_000);
    return () => { active = false; clearInterval(timer); };
  }, [tenant]);

  async function toggle() {
    const next = !open;
    setOpen(next);
    if (!next) return;
    setError('');
    try {
      setItems(await api<Notification[]>('/notifications', {}, tenant));
    } catch {
      setError(t('No pudimos cargar las notificaciones.', 'We could not load your notifications.'));
    }
  }

  async function markRead(item: Notification) {
    if (item.read) return;
    try {
      await api(`/notifications/${item.id}/read`, { method: 'POST' }, tenant);
      setItems(list => list?.map(entry => entry.id === item.id ? { ...entry, read: true } : entry));
      setUnread(count => Math.max(0, count - 1));
    } catch { /* Se reintentará al abrir de nuevo. */ }
  }

  async function markAllRead() {
    try {
      await api('/notifications/read-all', { method: 'POST' }, tenant);
      setItems(list => list?.map(entry => ({ ...entry, read: true })));
      setUnread(0);
    } catch { /* Sin cambios si el servidor no responde. */ }
  }

  return <div className="notifications">
    <Button variant="tertiary" className="icon-button" aria-label={t('Notificaciones', 'Notifications')} aria-expanded={open} onPress={() => { void toggle(); }}>
      <Bell size={17} />
      {unread > 0 && <span className="bell-badge" aria-label={`${unread} ${t('sin leer', 'unread')}`}>{unread > 99 ? '99+' : unread}</span>}
    </Button>
    {open && <>
      <button className="panel-backdrop" aria-hidden="true" tabIndex={-1} onClick={() => setOpen(false)} />
      <section className="notifications-panel" aria-label={t('Notificaciones', 'Notifications')}>
        <div className="section-header"><h2>{t('Notificaciones', 'Notifications')}</h2>{unread > 0 && <Button size="sm" variant="secondary" onPress={() => { void markAllRead(); }}>{t('Marcar todas como leídas', 'Mark all as read')}</Button>}</div>
        {error && <Notice error>{error}</Notice>}
        {!error && !items && <div className="loading" role="status"><span className="loading-dot" />{t('Cargando información…', 'Loading information…')}</div>}
        {items && (items.length === 0
          ? <Empty>{t('No tienes notificaciones.', 'You have no notifications.')}</Empty>
          : items.map(item => {
            const content = <><strong>{item.title}</strong><p>{item.body}</p><small>{new Date(item.created_at).toLocaleString()}</small></>;
            const className = `notification-item ${item.read ? '' : 'unread'}`;
            return item.ref_type === 'ticket' && item.ref_id
              ? <Link key={item.id} className={className} href={`/tickets/${item.ref_id}`} onClick={() => { void markRead(item); setOpen(false); }}>{content}</Link>
              : <button key={item.id} type="button" className={className} onClick={() => { void markRead(item); }}>{content}</button>;
          }))}
      </section>
    </>}
  </div>;
}
