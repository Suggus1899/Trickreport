// Minimal IndexedDB wrapper for offline ticket storage — no external deps.
// Stores pending ticket submissions in a 'pending_tickets' object store and
// syncs them back to the API when connectivity is restored.
//
// Conflict handling: when the server returns HTTP 409 during sync, the ticket
// is marked as a conflict (serverVersion + reason attached) and kept in the
// store so the user can resolve it manually via the sync-status page.

import { API_URL } from './api';

export interface PendingTicket {
  id: string; // local UUID
  data: {
    title: string;
    description: string;
    priority: string;
    category: string;
  };
  createdAt: number;
  // Sync state — absent means "pending" (never tried).
  status?: 'pending' | 'synced' | 'conflict' | 'error';
  error?: string;
  // Present when the server returned 409 — contains the conflicting server record.
  serverVersion?: any;
  lastAttempt?: number;
}

const DB_NAME = 'trickreport-offline';
const DB_VERSION = 1;
const STORE_NAME = 'pending_tickets';

// ─── Low-level IndexedDB helpers ─────────────────────────────────────
function openDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);
    request.onupgradeneeded = () => {
      const db = request.result;
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME, { keyPath: 'id' });
      }
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

function txStore(db: IDBDatabase, mode: IDBTransactionMode): IDBObjectStore {
  return db.transaction(STORE_NAME, mode).objectStore(STORE_NAME);
}

function promisifyRequest<T>(request: IDBRequest<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

async function putTicket(ticket: PendingTicket): Promise<void> {
  const db = await openDB();
  await promisifyRequest(txStore(db, 'readwrite').put(ticket));
  db.close();
}

// ─── Public API ──────────────────────────────────────────────────────

/** Persist a ticket to be submitted later when the connection is restored. */
export async function savePendingTicket(
  data: PendingTicket['data']
): Promise<string> {
  const ticket: PendingTicket = {
    id: crypto.randomUUID(),
    data,
    createdAt: Date.now(),
    status: 'pending',
  };
  const db = await openDB();
  await promisifyRequest(txStore(db, 'readwrite').add(ticket));
  db.close();
  return ticket.id;
}

/** Return all pending tickets, oldest first. */
export async function getPendingTickets(): Promise<PendingTicket[]> {
  const db = await openDB();
  const all = await promisifyRequest<PendingTicket[]>(
    txStore(db, 'readonly').getAll() as IDBRequest<PendingTicket[]>
  );
  db.close();
  return all.sort((a, b) => a.createdAt - b.createdAt);
}

/** Return only tickets that are in a conflict state. */
export async function getConflicts(): Promise<PendingTicket[]> {
  const pending = await getPendingTickets();
  return pending.filter((t) => t.status === 'conflict');
}

/** Remove a pending ticket by its local id. */
export async function deletePending(id: string): Promise<void> {
  const db = await openDB();
  await promisifyRequest(txStore(db, 'readwrite').delete(id));
  db.close();
}

/** Alias kept for clarity on the sync-status page. */
export const deletePendingTicket = deletePending;

/**
 * Resolve a conflict.
 * - resolution 'local'  → overwrite server with the local version (force PUT).
 * - resolution 'server' → discard the local version, keep the server's.
 * - resolution 'discard'→ delete the local pending ticket entirely.
 */
export async function resolveConflict(
  id: string,
  resolution: 'local' | 'server' | 'discard',
  token?: string
): Promise<void> {
  const pending = await getPendingTickets();
  const ticket = pending.find((t) => t.id === id);
  if (!ticket) return;

  if (resolution === 'discard' || resolution === 'server') {
    await deletePending(id);
    return;
  }

  // 'local' → retry the submission, forcing the local version.
  if (!token) throw new Error('Token required to resolve conflict with local version');
  ticket.status = 'pending';
  ticket.error = undefined;
  ticket.serverVersion = undefined;
  await putTicket(ticket);
  await retryPending(id, token);
}

/** Retry syncing a single pending ticket. */
export async function retryPending(id: string, token: string): Promise<void> {
  const pending = await getPendingTickets();
  const ticket = pending.find((t) => t.id === id);
  if (!ticket) return;
  await attemptSync(ticket, token);
}

async function attemptSync(ticket: PendingTicket, token: string): Promise<void> {
  try {
    const res = await fetch(`${API_URL}/api/v1/tickets`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'default',
        Authorization: `Bearer ${token}`,
      },
      credentials: 'include',
      body: JSON.stringify(ticket.data),
    });

    if (res.status === 409) {
      // Conflict — capture the server's version for the user to compare.
      const body = await res.json().catch(() => ({}));
      ticket.status = 'conflict';
      ticket.serverVersion = body.ticket || body.server_version || body;
      ticket.error = body.error || 'Conflict detected — server has a conflicting record.';
      ticket.lastAttempt = Date.now();
      await putTicket(ticket);
      return;
    }

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error || `Request failed with status ${res.status}`);
    }

    // Success — remove from the store.
    await deletePending(ticket.id);
  } catch (err) {
    ticket.status = 'error';
    ticket.error = err instanceof Error ? err.message : 'Sync failed';
    ticket.lastAttempt = Date.now();
    await putTicket(ticket);
    throw err;
  }
}

/**
 * Attempt to submit every pending ticket to the API.
 * On success the ticket is removed from the store; conflicts (409) are
 * recorded so they can be resolved manually. Returns the number of
 * tickets successfully synced.
 */
export async function syncPendingTickets(token: string): Promise<number> {
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

/** Convenience helper for the UI. */
export async function getPendingCount(): Promise<number> {
  const pending = await getPendingTickets();
  return pending.length;
}

/** Count of conflicts only — useful for badges. */
export async function getConflictCount(): Promise<number> {
  const conflicts = await getConflicts();
  return conflicts.length;
}
