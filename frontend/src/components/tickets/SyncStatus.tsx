import { useCallback, useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { getPendingTickets, deletePending, type PendingTicket } from '@/lib/offline/db';
import { syncPendingTickets, retryPending, resolveConflict } from '@/lib/offline/sync';

const STATUS_BADGE: Record<string, { variant: 'default' | 'destructive' | 'secondary' | 'outline'; label: string }> = {
  pending: { variant: 'outline', label: 'Pending' },
  conflict: { variant: 'destructive', label: 'Conflict' },
  error: { variant: 'secondary', label: 'Error' },
  synced: { variant: 'default', label: 'Synced' },
};

export function SyncStatus() {
  const [tickets, setTickets] = useState<PendingTicket[]>([]);
  const [syncMessage, setSyncMessage] = useState('');
  const [syncing, setSyncing] = useState(false);
  const [busyId, setBusyId] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setTickets(await getPendingTickets());
  }, []);

  useEffect(() => {
    refresh();
    window.addEventListener('online', refresh);
    return () => window.removeEventListener('online', refresh);
  }, [refresh]);

  async function handleSyncAll() {
    if (!navigator.onLine) {
      setSyncMessage("You're offline. Try again when connected.");
      return;
    }
    setSyncing(true);
    setSyncMessage('Syncing…');
    try {
      const synced = await syncPendingTickets();
      setSyncMessage(synced > 0 ? `Synced ${synced} ticket${synced > 1 ? 's' : ''} successfully.` : 'No new tickets synced.');
      await refresh();
    } catch (err) {
      setSyncMessage(`Sync failed: ${err instanceof Error ? err.message : 'unknown error'}`);
    } finally {
      setSyncing(false);
    }
  }

  async function handleAction(id: string, action: 'retry' | 'local' | 'server' | 'delete') {
    setBusyId(id);
    try {
      if (action === 'delete') {
        await deletePending(id);
      } else if (action === 'retry') {
        await retryPending(id);
      } else {
        await resolveConflict(id, action);
      }
    } catch (err) {
      setSyncMessage(`Action failed: ${err instanceof Error ? err.message : 'unknown error'}`);
    } finally {
      setBusyId(null);
      await refresh();
    }
  }

  const pendingCount = tickets.filter((t) => (t.status || 'pending') === 'pending').length;
  const conflictCount = tickets.filter((t) => t.status === 'conflict').length;
  const errorCount = tickets.filter((t) => t.status === 'error').length;

  return (
    <div className="flex flex-col gap-6">
      <div className="rounded-xl border bg-card p-5">
        <div className="flex items-center justify-between flex-wrap gap-3">
          <div className="flex flex-col gap-1 text-sm text-muted-foreground">
            <p>
              Pending: <span className="font-bold text-foreground">{pendingCount}</span>
            </p>
            <p>
              Conflicts: <span className="font-bold text-destructive">{conflictCount}</span>
            </p>
            <p>
              Errors: <span className="font-bold text-amber-600">{errorCount}</span>
            </p>
          </div>
          <Button type="button" variant="outline" onClick={handleSyncAll} disabled={syncing}>
            {syncing ? 'Syncing…' : 'Sync all'}
          </Button>
        </div>
        {syncMessage && <div className="mt-4 rounded-lg bg-muted p-4 text-sm">{syncMessage}</div>}
      </div>

      <div className="rounded-xl border bg-card p-5">
        {tickets.length === 0 ? (
          <p className="text-center text-muted-foreground py-6">No pending offline tickets. Everything is in sync.</p>
        ) : (
          <div className="flex flex-col gap-4">
            {tickets.map((t) => {
              const badge = STATUS_BADGE[t.status || 'pending'] ?? STATUS_BADGE.pending;
              const busy = busyId === t.id;
              return (
                <div key={t.id} className="rounded-lg bg-muted/50 p-4 flex flex-col gap-3">
                  <div className="flex items-center justify-between gap-3 flex-wrap">
                    <div>
                      <p className="font-semibold">{t.data.title || '(untitled)'}</p>
                      <p className="text-xs text-muted-foreground">
                        Created {new Date(t.createdAt).toLocaleString()} · {t.data.priority || 'medium'} · {t.data.category || 'general'}
                      </p>
                    </div>
                    <Badge variant={badge.variant}>{badge.label}</Badge>
                  </div>
                  {t.error && <p className="text-xs text-destructive">{t.error}</p>}
                  {t.serverVersion !== undefined && (
                    <div className="rounded-md bg-background p-3 text-xs">
                      <p className="font-semibold mb-1">Server version:</p>
                      <pre className="whitespace-pre-wrap text-muted-foreground">{JSON.stringify(t.serverVersion, null, 2)}</pre>
                    </div>
                  )}
                  <div className="flex items-center gap-2 flex-wrap">
                    <Button size="sm" variant="outline" disabled={busy} onClick={() => handleAction(t.id, 'retry')}>
                      Retry
                    </Button>
                    {t.status === 'conflict' && (
                      <>
                        <Button size="sm" variant="outline" disabled={busy} onClick={() => handleAction(t.id, 'local')}>
                          Keep local
                        </Button>
                        <Button size="sm" variant="outline" disabled={busy} onClick={() => handleAction(t.id, 'server')}>
                          Keep server
                        </Button>
                      </>
                    )}
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={busy}
                      onClick={() => {
                        if (confirm('Delete this pending ticket permanently?')) handleAction(t.id, 'delete');
                      }}
                    >
                      Delete
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
