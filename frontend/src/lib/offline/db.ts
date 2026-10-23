// Minimal IndexedDB wrapper for offline ticket storage — no external deps.
// Pure storage layer: no network calls here, see ./sync.ts for that.

export interface OfflineIdentity {
  userId: string;
  tenantId: string;
}

export interface PendingTicket extends OfflineIdentity {
  id: string; // local UUID
  data: {
    title: string;
    description: string;
    priority: string;
    category: string;
  };
  createdAt: number;
  // Sync state — absent means "pending" (never tried).
  status?: 'pending' | 'syncing' | 'synced' | 'conflict' | 'invalid' | 'error';
  error?: string;
  // Present when the server returned 409 — contains the conflicting server record.
  serverVersion?: unknown;
  lastAttempt?: number;
  leaseExpiresAt?: number;
}

const DB_NAME = 'trickreport-offline';
const DB_VERSION = 2;
const STORE_NAME = 'pending_tickets';
const LEGACY_STORE = 'legacy_pending_tickets';

function openDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);
    request.onupgradeneeded = () => {
      const db = request.result;
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME, { keyPath: 'id' });
      }
      if (!db.objectStoreNames.contains(LEGACY_STORE)) db.createObjectStore(LEGACY_STORE, { keyPath: 'id' });
      const transaction = request.transaction!;
      const cursor = transaction.objectStore(STORE_NAME).openCursor();
      cursor.onsuccess = () => {
        const row = cursor.result;
        if (!row) return;
        if (!row.value.userId || !row.value.tenantId) {
          transaction.objectStore(LEGACY_STORE).put(row.value);
          row.delete();
        }
        row.continue();
      };
    };
    let blocked = false;
    request.onblocked = () => { blocked = true; reject(new Error('Close other Trickreport tabs to upgrade offline storage.')); };
    request.onsuccess = () => {
      if (blocked) { request.result.close(); return; }
      request.result.onversionchange = () => request.result.close();
      resolve(request.result);
    };
    request.onerror = () => reject(request.error);
  });
}

async function withStore<T>(mode: IDBTransactionMode, action: (store: IDBObjectStore) => Promise<T>, name = STORE_NAME): Promise<T> {
  const db = await openDB();
  const transaction = db.transaction(name, mode);
  const completion = new Promise<void>((resolve, reject) => {
    transaction.oncomplete = () => resolve();
    transaction.onabort = () => reject(transaction.error ?? new Error('Offline storage transaction aborted.'));
    transaction.onerror = () => reject(transaction.error);
  });
  try {
    const result = await action(transaction.objectStore(name));
    await completion;
    return result;
  } catch (error) {
    await completion.catch(() => {});
    throw error;
  } finally { db.close(); }
}

function promisifyRequest<T>(request: IDBRequest<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

/** Persist a ticket to be submitted later when the connection is restored. */
export async function savePendingTicket(data: PendingTicket['data'], identity: OfflineIdentity, id = crypto.randomUUID()): Promise<string> {
  if (!identity?.userId || !identity?.tenantId) throw new Error('A verified user and tenant are required to save an offline draft.');
  const ticket: PendingTicket = {
    id,
    userId: identity.userId,
    tenantId: identity.tenantId,
    data,
    createdAt: Date.now(),
    status: 'pending',
  };
  await withStore('readwrite', (store) => promisifyRequest(store.add(ticket)));
  return ticket.id;
}

/** Overwrite an existing pending ticket record (used to update sync state). */
export async function putTicket(ticket: PendingTicket): Promise<void> {
  await withStore('readwrite', async (store) => {
    const current = await promisifyRequest<PendingTicket | undefined>(store.get(ticket.id));
    if (current?.userId === ticket.userId && current.tenantId === ticket.tenantId) {
      await promisifyRequest(store.put(ticket));
    }
  });
}

/** Return all pending tickets, oldest first. */
export async function getPendingTickets(identity: OfflineIdentity): Promise<PendingTicket[]> {
  const all = await withStore('readonly', (store) => promisifyRequest<PendingTicket[]>(store.getAll()));
  return all.filter((ticket) => ticket.userId === identity.userId && ticket.tenantId === identity.tenantId).sort((a, b) => a.createdAt - b.createdAt);
}

/** Return only tickets that are in a conflict state. */
export async function getConflicts(identity: OfflineIdentity): Promise<PendingTicket[]> {
  const pending = await getPendingTickets(identity);
  return pending.filter((t) => t.status === 'conflict');
}

/** Remove a pending ticket by its local id. */
export async function deletePending(id: string, identity: OfflineIdentity): Promise<void> {
  await withStore('readwrite', async (store) => {
    const ticket = await promisifyRequest<PendingTicket | undefined>(store.get(id));
    if (ticket?.userId === identity.userId && ticket.tenantId === identity.tenantId) await promisifyRequest(store.delete(id));
  });
}

/** A transaction prevents two tabs from submitting the same draft concurrently. */
export async function claimPendingTicket(id: string, identity: OfflineIdentity, retry = false): Promise<PendingTicket | undefined> {
  return withStore('readwrite', async (store) => {
    const ticket = await promisifyRequest<PendingTicket | undefined>(store.get(id));
    if (!ticket || ticket.userId !== identity.userId || ticket.tenantId !== identity.tenantId) return;
    if ((ticket.leaseExpiresAt ?? 0) > Date.now() || (!retry && (ticket.status === 'conflict' || ticket.status === 'invalid'))) return;
    ticket.status = 'syncing';
    // Recover a killed tab after the timeout and retry window has elapsed.
    ticket.leaseExpiresAt = Date.now() + 180_000;
    await promisifyRequest(store.put(ticket));
    return ticket;
  });
}

export async function getLegacyPendingCount(): Promise<number> {
  return withStore('readonly', (store) => promisifyRequest(store.count()), LEGACY_STORE);
}

/** Alias kept for clarity on the sync-status page. */
export const deletePendingTicket = deletePending;

/** Convenience helper for the UI. */
export async function getPendingCount(identity: OfflineIdentity): Promise<number> {
  const pending = await getPendingTickets(identity);
  return pending.length;
}

/** Count of conflicts only — useful for badges. */
export async function getConflictCount(identity: OfflineIdentity): Promise<number> {
  const conflicts = await getConflicts(identity);
  return conflicts.length;
}
