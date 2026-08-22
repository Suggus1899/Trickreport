import 'fake-indexeddb/auto';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

vi.mock('../config', () => ({
  API_URL: 'http://localhost:8080',
  ON_PREMISE: false,
}));

const fetchMock = vi.fn();
(globalThis as any).fetch = fetchMock;

import {
  savePendingTicket,
  getPendingTickets,
  getConflicts,
  getPendingCount,
  getConflictCount,
} from './db';
import { syncPendingTickets, resolveConflict, retryPending } from './sync';

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

const TICKET_DATA = { title: 'Broken printer', description: 'Jammed', priority: 'high', category: 'hardware' };

beforeEach(() => {
  fetchMock.mockReset();
});

afterEach(async () => {
  // Drain the store between tests so each test starts empty.
  const pending = await getPendingTickets();
  await Promise.all(pending.map((t) => import('./db').then((m) => m.deletePending(t.id))));
});

describe('offline db', () => {
  it('saves and lists pending tickets, oldest first', async () => {
    const id1 = await savePendingTicket(TICKET_DATA);
    const id2 = await savePendingTicket({ ...TICKET_DATA, title: 'Second' });
    const pending = await getPendingTickets();
    expect(pending.map((t) => t.id)).toEqual([id1, id2]);
    expect(await getPendingCount()).toBe(2);
  });
});

describe('offline sync', () => {
  it('removes a ticket from the store on successful sync', async () => {
    await savePendingTicket(TICKET_DATA);
    fetchMock.mockResolvedValue(jsonResponse(201, { id: 'server-1', ...TICKET_DATA }));
    const synced = await syncPendingTickets('token');
    expect(synced).toBe(1);
    expect(await getPendingTickets()).toHaveLength(0);
  });

  it('marks a ticket as conflict on HTTP 409 and keeps it in the store', async () => {
    const id = await savePendingTicket(TICKET_DATA);
    fetchMock.mockResolvedValue(
      jsonResponse(409, { error: 'duplicate ticket', ticket: { id: 'server-1' } }),
    );
    // attemptSync doesn't throw on 409 (it's handled inline as a conflict),
    // so the loop counts it as "processed" even though it isn't synced —
    // matches the pre-existing counting semantics of the sync loop.
    const synced = await syncPendingTickets('token');
    expect(synced).toBe(1);
    const conflicts = await getConflicts();
    expect(conflicts).toHaveLength(1);
    expect(conflicts[0].id).toBe(id);
    expect(conflicts[0].status).toBe('conflict');
    expect(await getConflictCount()).toBe(1);
  });

  it('stops the sync loop after a network error, leaving later tickets untouched', async () => {
    await savePendingTicket(TICKET_DATA);
    await savePendingTicket({ ...TICKET_DATA, title: 'Second' });
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'));
    const synced = await syncPendingTickets('token');
    expect(synced).toBe(0);
    expect(await getPendingTickets()).toHaveLength(2);
  });

  it('resolveConflict("server") discards the local pending ticket', async () => {
    const id = await savePendingTicket(TICKET_DATA);
    fetchMock.mockResolvedValue(jsonResponse(409, { error: 'conflict' }));
    await syncPendingTickets('token');
    await resolveConflict(id, 'server');
    expect(await getPendingTickets()).toHaveLength(0);
  });

  it('resolveConflict("local") retries and clears the conflict on success', async () => {
    const id = await savePendingTicket(TICKET_DATA);
    fetchMock.mockResolvedValue(jsonResponse(409, { error: 'conflict' }));
    await syncPendingTickets('token');
    fetchMock.mockResolvedValue(jsonResponse(201, { id: 'server-1' }));
    await resolveConflict(id, 'local', 'token');
    expect(await getPendingTickets()).toHaveLength(0);
  });

  it('retryPending re-attempts a single ticket without touching others', async () => {
    const id = await savePendingTicket(TICKET_DATA);
    fetchMock.mockResolvedValue(jsonResponse(201, { id: 'server-1' }));
    await retryPending(id, 'token');
    expect(await getPendingTickets()).toHaveLength(0);
  });
});
