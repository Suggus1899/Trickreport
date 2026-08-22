import { useEffect, useRef, useState } from 'react';
import {
  getNotifications,
  markNotificationRead,
  markAllNotificationsRead,
  type Notification,
} from '@/lib/api';
import { PUBLIC_API_URL } from '@/lib/config';

function timeAgo(iso: string): string {
  const s = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
  if (s < 60) return 'just now';
  if (s < 3600) return `${Math.floor(s / 60)}m ago`;
  if (s < 86400) return `${Math.floor(s / 3600)}h ago`;
  return `${Math.floor(s / 86400)}d ago`;
}

function wsUrl(): string {
  if (PUBLIC_API_URL) return `${PUBLIC_API_URL.replace(/^http/, 'ws')}/api/v1/ws`;
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${proto}//${window.location.host}/api/v1/ws`;
}

interface Toast extends Notification {
  toastId: string;
}

export function NotificationBell() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [open, setOpen] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const rootRef = useRef<HTMLDivElement>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  const unread = notifications.filter((n) => !n.read).length;

  async function load() {
    try {
      const data = await getNotifications();
      setNotifications(data);
    } catch {
      // silent — bell just stays at its last known state
    }
  }

  function pushToast(n: Notification) {
    const toast: Toast = { ...n, toastId: `${n.id}-${Date.now()}` };
    setToasts((t) => [...t, toast]);
    setTimeout(() => setToasts((t) => t.filter((x) => x.toastId !== toast.toastId)), 6000);
  }

  function handleRealtimeMessage(raw: string) {
    try {
      const msg = JSON.parse(raw);
      if (!msg?.type) return;
      if (msg.type === 'NOTIFICATION' && msg.data) {
        setNotifications((prev) => [msg.data, ...prev]);
        pushToast(msg.data);
        return;
      }
      if (msg.type === 'TICKET_CREATED' && msg.data) {
        const n: Notification = {
          id: `ws-${Date.now()}`,
          tenant_id: msg.data.tenant_id ?? '',
          user_id: '',
          title: `New ticket: ${msg.data.title || 'Untitled'}`,
          body: 'A new ticket has been created.',
          type: 'ticket_created',
          ref_id: msg.data.id,
          ref_type: 'ticket',
          read: false,
          created_at: new Date().toISOString(),
        };
        setNotifications((prev) => [n, ...prev]);
        pushToast(n);
        return;
      }
      if (msg.type === 'TICKET_UPDATED' && msg.data) {
        const n: Notification = {
          id: `ws-${Date.now()}`,
          tenant_id: '',
          user_id: '',
          title: `Ticket updated: status → ${msg.data.status || 'unknown'}`,
          body: 'A ticket status has been changed.',
          type: 'ticket_updated',
          ref_id: msg.data.ticket_id,
          ref_type: 'ticket',
          read: false,
          created_at: new Date().toISOString(),
        };
        setNotifications((prev) => [n, ...prev]);
        pushToast(n);
        return;
      }
      if (msg.type === 'NEW_COMMENT' && msg.data) {
        const n: Notification = {
          id: `ws-${Date.now()}`,
          tenant_id: '',
          user_id: '',
          title: 'New comment',
          body: String(msg.data.content || '').slice(0, 100),
          type: 'new_comment',
          ref_id: msg.data.ticket_id,
          ref_type: 'ticket',
          read: false,
          created_at: new Date().toISOString(),
        };
        setNotifications((prev) => [n, ...prev]);
        pushToast(n);
      }
    } catch {
      // ignore non-JSON frames
    }
  }

  useEffect(() => {
    load();

    function connect() {
      if (!('WebSocket' in window)) return;
      // The httpOnly auth cookie rides along automatically — same-origin in
      // production behind Caddy. No token is ever exposed to this script.
      const ws = new WebSocket(wsUrl());
      wsRef.current = ws;
      ws.addEventListener('open', load);
      ws.addEventListener('message', (e) => handleRealtimeMessage(e.data));
      ws.addEventListener('close', () => {
        wsRef.current = null;
        reconnectTimer.current = setTimeout(connect, 5000);
      });
      ws.addEventListener('error', () => ws.close());
    }
    connect();

    const poll = setInterval(load, 30000);
    return () => {
      clearInterval(poll);
      clearTimeout(reconnectTimer.current);
      wsRef.current?.close();
    };
  }, []);

  useEffect(() => {
    function onClickOutside(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener('click', onClickOutside);
    return () => document.removeEventListener('click', onClickOutside);
  }, []);

  async function handleMarkRead(id: string) {
    setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, read: true } : n)));
    try {
      await markNotificationRead(id);
    } catch {
      // best-effort
    }
  }

  async function handleMarkAll() {
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    try {
      await markAllNotificationsRead();
    } catch {
      // best-effort
    }
  }

  return (
    <>
      <div ref={rootRef} className="relative">
        <button
          type="button"
          aria-label="Notifications"
          aria-haspopup="true"
          aria-expanded={open}
          onClick={(e) => {
            e.stopPropagation();
            setOpen((v) => !v);
          }}
          className="relative inline-flex items-center justify-center w-10 h-10 rounded-full border bg-background hover:bg-muted"
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
            <path d="M13.73 21a2 2 0 0 1-3.46 0" />
          </svg>
          {unread > 0 && (
            <span className="absolute -top-0.5 -right-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-destructive px-1 text-[0.65rem] font-bold text-destructive-foreground">
              {unread > 9 ? '9+' : unread}
            </span>
          )}
        </button>

        {open && (
          <div
            role="dialog"
            aria-label="Notifications"
            className="absolute top-[calc(100%+0.5rem)] right-0 z-50 w-88 max-w-[90vw] max-h-96 overflow-y-auto rounded-lg border bg-popover p-2 shadow-md"
          >
            <div className="flex items-center justify-between px-2 py-1.5">
              <span className="text-sm font-bold">Notifications</span>
              <button type="button" onClick={handleMarkAll} className="text-xs font-semibold text-primary hover:underline">
                Mark all read
              </button>
            </div>
            {notifications.length === 0 ? (
              <div className="p-5 text-center text-sm text-muted-foreground">No notifications yet.</div>
            ) : (
              <div className="flex flex-col">
                {notifications.slice(0, 20).map((n) => (
                  <button
                    key={n.id}
                    type="button"
                    onClick={() => handleMarkRead(n.id)}
                    className={`flex flex-col items-start gap-0.5 rounded-md px-2.5 py-2 text-left text-sm hover:bg-accent ${
                      n.read ? '' : 'bg-primary/5'
                    }`}
                  >
                    <span className="font-semibold">{n.title}</span>
                    <span className="text-xs text-muted-foreground">{n.body}</span>
                    <span className="text-[0.7rem] text-muted-foreground">{timeAgo(n.created_at)}</span>
                  </button>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      <div className="fixed bottom-6 right-6 z-[200] flex flex-col gap-2" aria-live="polite" aria-atomic="true">
        {toasts.map((t) => (
          <div key={t.toastId} className="flex items-start gap-2.5 w-72 max-w-[24rem] rounded-lg border bg-card p-3.5 shadow-md animate-in slide-in-from-right">
            <div className="flex-1 min-w-0">
              <p className="text-sm font-bold">{t.title}</p>
              <p className="text-xs text-muted-foreground mt-0.5">{t.body}</p>
            </div>
            <button
              type="button"
              aria-label="Close"
              onClick={() => setToasts((prev) => prev.filter((x) => x.toastId !== t.toastId))}
              className="shrink-0 text-muted-foreground hover:text-foreground"
            >
              &times;
            </button>
          </div>
        ))}
      </div>
    </>
  );
}
