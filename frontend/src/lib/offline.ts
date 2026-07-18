// Minimal IndexedDB wrapper for offline ticket storage — no external deps.
// Stores pending ticket submissions in a 'pending_tickets' object store and
// syncs them back to the API when connectivity is restored.

import { createTicket } from './api';

export interface PendingTicket {
  id: string; // local UUID
  data: {
    title: string;
    description: string;
    priority: string;
    category: string;
  };
  createdAt: number;
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

// ─── Public API ──────────────────────────────────────────────────────

/** Persist a ticket to be submitted later when the connection is restored. */
export async function savePendingTicket(
  data: PendingTicket['data']
): Promise<string> {
  const ticket: PendingTicket = {
    id: crypto.randomUUID(),
    data,
    createdAt: Date.now(),
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

/** Remove a pending ticket by its local id. */
export async function deletePendingTicket(id: string): Promise<void> {
  const db = await openDB();
  await promisifyRequest(txStore(db, 'readwrite').delete(id));
  db.close();
}

/**
 * Attempt to submit every pending ticket to the API.
 * On success the ticket is removed from the store; failures are tolerated
 * so that one bad ticket doesn't block the rest. Returns the number of
 * tickets successfully synced.
 */
export async function syncPendingTickets(token: string): Promise<number> {
  const pending = await getPendingTickets();
  let synced = 0;
  for (const ticket of pending) {
    try {
      await createTicket(token, ticket.data);
      await deletePendingTicket(ticket.id);
      synced++;
    } catch (err) {
      // Stop on the first failure — likely a connectivity issue, so the
      // remaining tickets will be retried on the next sync attempt.
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
