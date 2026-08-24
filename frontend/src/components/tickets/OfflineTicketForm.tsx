import { useEffect, useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { createTicket, ApiNetworkError } from '@/lib/api';
import { useOnlineStatus } from '@/lib/hooks/useOnlineStatus';
import { useOfflineTickets } from '@/lib/hooks/useOfflineTickets';

const PRIORITIES = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
  { value: 'critical', label: 'Critical' },
];

export function OfflineTicketForm() {
  const online = useOnlineStatus();
  const { pendingCount, save, sync } = useOfflineTickets();
  const [banner, setBanner] = useState('');
  const [syncMessage, setSyncMessage] = useState('');
  const [syncing, setSyncing] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  async function runSync() {
    if (!navigator.onLine) {
      setSyncMessage("You're still offline. Try again when connected.");
      return;
    }
    setSyncing(true);
    setSyncMessage('Syncing...');
    try {
      const synced = await sync();
      setSyncMessage(
        synced > 0
          ? `Synced ${synced} ticket${synced > 1 ? 's' : ''} successfully.`
          : 'No pending tickets to sync.',
      );
    } catch (err) {
      setSyncMessage(`Sync failed: ${err instanceof Error ? err.message : 'unknown error'}`);
    } finally {
      setSyncing(false);
    }
  }

  useEffect(() => {
    window.addEventListener('online', runSync);
    window.addEventListener('trickreport:sync-tickets', runSync);
    return () => {
      window.removeEventListener('online', runSync);
      window.removeEventListener('trickreport:sync-tickets', runSync);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const formData = new FormData(form);
    const data = {
      title: formData.get('title')?.toString() || '',
      description: formData.get('description')?.toString() || '',
      priority: formData.get('priority')?.toString() || 'medium',
      category: formData.get('category')?.toString() || 'general',
    };

    setSubmitting(true);
    setBanner('');
    try {
      if (navigator.onLine) {
        try {
          await createTicket(data);
          window.location.href = '/tickets';
          return;
        } catch (err) {
          // Only a real network failure falls back to the offline queue — a
          // validation error (4xx) would fail identically on every future
          // sync attempt, so it must surface to the user now instead.
          if (!(err instanceof ApiNetworkError)) {
            setBanner(err instanceof Error ? err.message : 'Failed to create ticket.');
            return;
          }
        }
      }
      await save(data);
      setBanner(
        navigator.onLine
          ? 'Network error — saved offline, will sync when connected.'
          : 'Saved offline, will sync when connected.',
      );
      form.reset();
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="rounded-xl border bg-card p-6">
      {banner && (
        <div className="rounded-lg bg-amber-500/10 text-amber-700 dark:text-amber-400 p-4 text-sm mb-5">
          {banner}
        </div>
      )}
      {syncMessage && <div className="rounded-lg bg-muted p-4 text-sm mb-5">{syncMessage}</div>}
      {pendingCount > 0 && (
        <div className="rounded-lg bg-muted p-4 text-sm text-muted-foreground mb-5">
          {pendingCount} pending ticket{pendingCount > 1 ? 's' : ''} waiting to sync.
        </div>
      )}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <div className="flex flex-col gap-2">
          <Label htmlFor="title">Title</Label>
          <Input id="title" name="title" required minLength={3} />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="description">Description</Label>
          <Textarea id="description" name="description" required minLength={10} rows={5} />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <div className="flex flex-col gap-2">
            <Label htmlFor="priority">Priority</Label>
            <select
              id="priority"
              name="priority"
              className="h-9 rounded-md border bg-background px-3 text-sm"
            >
              {PRIORITIES.map((p) => (
                <option key={p.value} value={p.value}>
                  {p.label}
                </option>
              ))}
            </select>
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="category">Category</Label>
            <Input id="category" name="category" defaultValue="general" />
          </div>
        </div>
        <div className="flex items-center justify-end gap-3 pt-2">
          <a
            href="/tickets"
            className="text-sm font-semibold text-muted-foreground hover:text-foreground px-4 py-2"
          >
            Cancel
          </a>
          <Button type="submit" disabled={submitting}>
            {submitting ? 'Saving…' : 'Create Ticket'}
          </Button>
        </div>
      </form>

      <div className="flex items-center justify-between gap-3 pt-6 mt-6 border-t">
        <span className="text-sm text-muted-foreground">
          Pending offline tickets will sync automatically when you're back online.
        </span>
        <Button type="button" variant="outline" onClick={runSync} disabled={syncing || !online}>
          {syncing ? 'Syncing…' : 'Sync now'}
        </Button>
      </div>
    </div>
  );
}
