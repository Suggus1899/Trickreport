// Minimal IndexedDB wrapper for offline ticket storage — no external deps.
// Pure storage layer: no network calls here, see ./sync.ts for that.

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
  serverVersion?: unknown;
  lastAttempt?: number;
}

const DB_NAME = 'trickreport-offline';
const DB_VERSION = 1;
const STORE_NAME = 'pending_tickets';

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

/** Persist a ticket to be submitted later when the connection is restored. */
export async function savePendingTicket(data: PendingTicket['data']): Promise<string> {
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

/** Overwrite an existing pending ticket record (used to update sync state). */
export async function putTicket(ticket: PendingTicket): Promise<void> {
  const db = await openDB();
  await promisifyRequest(txStore(db, 'readwrite').put(ticket));
  db.close();
}

/** Return all pending tickets, oldest first. */
export async function getPendingTickets(): Promise<PendingTicket[]> {
  const db = await openDB();
  const all = await promisifyRequest<PendingTicket[]>(
    txStore(db, 'readonly').getAll() as IDBRequest<PendingTicket[]>,
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
