// Syncs offline-queued tickets back to the API once connectivity is restored.
//
// Conflict handling: when the server returns HTTP 409, the ticket is marked
// as a conflict (serverVersion + reason attached) and kept in the store so
// the user can resolve it manually (Retry / Keep local / Keep server / Delete).

import { createTicket, ApiValidationError } from '../api';
import {
  type PendingTicket,
  getPendingTickets,
  deletePending,
  putTicket,
} from './db';

const SYNC_TAG = 'sync-tickets';

/**
 * Registers a Background Sync request so the service worker wakes the page
 * (via a postMessage, see public/sw.js) as soon as connectivity returns —
 * even if the tab isn't focused. Falls back silently where unsupported
 * (e.g. Safari); the `online` event + manual "Sync now" button still work.
 */
export async function registerBackgroundSync(): Promise<void> {
  if (!('serviceWorker' in navigator)) return;
  const registration = await navigator.serviceWorker.ready;
  if (!('sync' in registration)) return;
  try {
    await (registration as ServiceWorkerRegistration & { sync: { register(tag: string): Promise<void> } }).sync.register(
      SYNC_TAG,
    );
  } catch (err) {
    console.warn('Background sync registration failed:', err);
  }
}

async function attemptSync(ticket: PendingTicket, token?: string): Promise<void> {
  try {
    await createTicket(ticket.data, token);
    await deletePending(ticket.id);
  } catch (err) {
    if (err instanceof ApiValidationError && err.status === 409) {
      const body = (err.details ?? {}) as { ticket?: unknown; server_version?: unknown };
      ticket.status = 'conflict';
      ticket.serverVersion = body.ticket ?? body.server_version ?? err.details;
      ticket.error = err.message || 'Conflict detected — server has a conflicting record.';
      ticket.lastAttempt = Date.now();
      await putTicket(ticket);
      return;
    }
    ticket.status = 'error';
    ticket.error = err instanceof Error ? err.message : 'Sync failed';
    ticket.lastAttempt = Date.now();
    await putTicket(ticket);
    throw err;
  }
}

/** Retry syncing a single pending ticket. */
export async function retryPending(id: string, token?: string): Promise<void> {
  const pending = await getPendingTickets();
  const ticket = pending.find((t) => t.id === id);
  if (!ticket) return;
  await attemptSync(ticket, token);
}

/**
 * Resolve a conflict.
 * - resolution 'local'  → retry the submission, forcing the local version.
 * - resolution 'server' → discard the local version, keep the server's.
 * - resolution 'discard'→ delete the local pending ticket entirely.
 */
export async function resolveConflict(
  id: string,
  resolution: 'local' | 'server' | 'discard',
  token?: string,
): Promise<void> {
  const pending = await getPendingTickets();
  const ticket = pending.find((t) => t.id === id);
  if (!ticket) return;

  if (resolution === 'discard' || resolution === 'server') {
    await deletePending(id);
    return;
  }

  ticket.status = 'pending';
  ticket.error = undefined;
  ticket.serverVersion = undefined;
  await putTicket(ticket);
  await retryPending(id, token);
}

/**
 * Attempt to submit every pending ticket to the API.
 * On success the ticket is removed from the store; conflicts (409) are
 * recorded so they can be resolved manually. Returns the number of tickets
 * successfully synced.
 */
export async function syncPendingTickets(token?: string): Promise<number> {
  const pending = await getPendingTickets();
  let synced = 0;
  for (const ticket of pending) {
    // Skip already-conflicted tickets unless explicitly retried.
    if (ticket.status === 'conflict') continue;
    try {
      await attemptSync(ticket, token);
      synced++;
    } catch (err) {
      // Network/transport error — stop here; remaining tickets retry next time.
      console.warn(`Failed to sync ticket ${ticket.id}:`, err);
      break;
    }
  }
  return synced;
}
