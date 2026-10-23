// Syncs offline-queued tickets back to the API once connectivity is restored.
//
// Conflict handling: when the server returns HTTP 409, the ticket is marked
// as a conflict (serverVersion + reason attached) and kept in the store so
// the user can resolve it manually (Retry / Keep local / Keep server / Delete).

import { createTicket, getMe, ApiAuthError, ApiValidationError } from '../api';
import {
  type PendingTicket,
  getPendingTickets,
  deletePending,
  putTicket,
  claimPendingTicket,
  type OfflineIdentity,
} from './db';

const SYNC_TAG = 'sync-tickets';
const activeSyncs = new Set<AbortController>();

export function cancelPendingSync(): void {
  for (const controller of activeSyncs) controller.abort();
}

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

async function attemptSync(id: string, identity: OfflineIdentity, signal: AbortSignal, token?: string, retry = false): Promise<boolean> {
  const ticket = await claimPendingTicket(id, identity, retry);
  if (!ticket) return false;
  try {
    signal.throwIfAborted();
    await createTicket(ticket.data, token, { idempotencyKey: ticket.id, signal });
    signal.throwIfAborted();
    await deletePending(ticket.id, identity);
    return true;
  } catch (err) {
    ticket.leaseExpiresAt = undefined;
    ticket.lastAttempt = Date.now();
    if (err instanceof ApiValidationError && err.status === 409) {
      const body = (err.details ?? {}) as { ticket?: unknown; server_version?: unknown };
      ticket.status = 'conflict';
      ticket.serverVersion = body.ticket ?? body.server_version ?? err.details;
      ticket.error = err.message || 'Conflict detected — server has a conflicting record.';
      await putTicket(ticket);
      return false;
    }
    ticket.status = err instanceof ApiValidationError && err.status !== 403 ? 'invalid' : 'error';
    ticket.error = err instanceof Error ? err.message : 'Sync failed';
    await putTicket(ticket);
    throw err;
  }
}

/** Retry syncing a single pending ticket. */
export async function retryPending(id: string, identity: OfflineIdentity, token?: string): Promise<boolean> {
  return (await syncPendingTickets(identity, token, id)) === 1;
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
  identity: OfflineIdentity,
  token?: string,
): Promise<void> {
  const pending = await getPendingTickets(identity);
  const ticket = pending.find((t) => t.id === id);
  if (!ticket) return;

  if (resolution === 'discard' || resolution === 'server') {
    await deletePending(id, identity);
    return;
  }

  ticket.status = 'pending';
  ticket.error = undefined;
  ticket.serverVersion = undefined;
  await putTicket(ticket);
  await retryPending(id, identity, token);
}

/**
 * Attempt to submit every pending ticket to the API.
 * On success the ticket is removed from the store; conflicts (409) are
 * recorded so they can be resolved manually. Returns the number of tickets
 * successfully synced.
 */
export async function syncPendingTickets(identity: OfflineIdentity, token?: string, onlyId?: string): Promise<number> {
  const controller = new AbortController();
  activeSyncs.add(controller);
  try {
    const user = await getMe(token);
    if (user.id !== identity.userId || user.tenant_id !== identity.tenantId) throw new ApiAuthError('The current session does not match this draft identity.');
    const pending = await getPendingTickets(identity);
    let synced = 0;
    for (const ticket of pending) {
      if (controller.signal.aborted) break;
      if (onlyId && ticket.id !== onlyId) continue;
      try {
        if (await attemptSync(ticket.id, identity, controller.signal, token, Boolean(onlyId))) synced++;
      } catch (err) {
        if (err instanceof ApiAuthError || (err instanceof ApiValidationError && err.status === 403)) throw err;
        // A transport failure leaves the UUID and payload intact for the next attempt.
        break;
      }
    }
    return synced;
  } finally { activeSyncs.delete(controller); }
}
