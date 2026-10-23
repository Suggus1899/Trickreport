import 'fake-indexeddb/auto';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { savePendingTicket, getPendingTickets, getConflicts, deletePending } from './db';
import { syncPendingTickets, retryPending, resolveConflict, cancelPendingSync } from './sync';

vi.mock('../config', () => ({ API_URL: 'http://api.test', PUBLIC_API_URL: '', resolveTenantSlug: () => 'default' }));

const A = { userId: 'user-a', tenantId: 'tenant-a' };
const B = { userId: 'user-b', tenantId: 'tenant-a' };
const OTHER_TENANT = { userId: 'user-a', tenantId: 'tenant-b' };
const DATA = { title: 'Broken printer', description: 'Paper is jammed', priority: 'high', category: 'hardware' };
const writes: Array<{ key: string | null; title: string }> = [];
let session = A;
let fail: number | undefined;

function request<T>(req: IDBRequest<T>): Promise<T> {
  return new Promise((resolve, reject) => { req.onsuccess = () => resolve(req.result); req.onerror = () => reject(req.error); });
}

beforeEach(async () => {
  await request(indexedDB.deleteDatabase('trickreport-offline'));
  session = A;
  fail = undefined;
  writes.length = 0;
  vi.stubGlobal('fetch', async (url: string, init?: RequestInit) => {
    if (url.endsWith('/auth/me')) return new Response(JSON.stringify({ id: session.userId, tenant_id: session.tenantId, role: 'end_user', name: 'User', email: 'u@test.test', active: true, created_at: '2026-10-01T00:00:00Z' }));
    writes.push({ key: new Headers(init?.headers).get('Idempotency-Key'), title: JSON.parse(String(init?.body)).title });
    return new Response(JSON.stringify(fail ? { error: 'Failed request' } : { id: 'server-ticket' }), { status: fail ?? 201 });
  });
});
afterEach(() => vi.unstubAllGlobals());

describe('identity-scoped offline drafts', () => {
  it('never shows or syncs user A drafts as B', async () => {
    await savePendingTicket(DATA, A);
    session = B;
    expect(await getPendingTickets(B)).toEqual([]);
    expect(await syncPendingTickets(B)).toBe(0);
    await retryPending((await getPendingTickets(A))[0].id, B);
    expect(writes).toEqual([]);
    expect(await getPendingTickets(A)).toHaveLength(1);
  });

  it('keeps the same user isolated in another tenant', async () => {
    await savePendingTicket(DATA, A);
    session = OTHER_TENANT;
    expect(await getPendingTickets(OTHER_TENANT)).toEqual([]);
    expect(await syncPendingTickets(OTHER_TENANT)).toBe(0);
    expect(writes).toEqual([]);
  });

  it('refuses sync when the validated current session no longer matches the draft owner', async () => {
    await savePendingTicket(DATA, A);
    session = B;
    await expect(syncPendingTickets(A)).rejects.toThrow(/session|identity/i);
    expect(writes).toEqual([]);
    expect(await getPendingTickets(A)).toHaveLength(1);
  });

  it('quarantines unowned legacy drafts without deleting them', async () => {
    const opening = indexedDB.open('trickreport-offline', 1);
    opening.onupgradeneeded = () => opening.result.createObjectStore('pending_tickets', { keyPath: 'id' });
    const old = await request(opening);
    await request(old.transaction('pending_tickets', 'readwrite').objectStore('pending_tickets').add({ id: 'legacy-id', data: DATA, createdAt: 1 }));
    old.close();
    expect(await getPendingTickets(A)).toEqual([]);
    const db = await request(indexedDB.open('trickreport-offline'));
    const legacy = await request(db.transaction('legacy_pending_tickets').objectStore('legacy_pending_tickets').getAll());
    db.close();
    expect(legacy).toEqual([{ id: 'legacy-id', data: DATA, createdAt: 1 }]);
    expect(await getPendingTickets(B)).toEqual([]);
  });

  it('does not allow another identity to delete an owned draft', async () => {
    const id = await savePendingTicket(DATA, A);
    await deletePending(id, B);
    expect(await getPendingTickets(A)).toHaveLength(1);
  });
});

describe('idempotent offline sync', () => {
  it('keeps the local UUID after response loss and removes only a confirmed result', async () => {
    const id = await savePendingTicket(DATA, A);
    const server = new Map<string, string>();
    let submissions = 0;
    vi.stubGlobal('fetch', async (url: string, init?: RequestInit) => {
      if (url.endsWith('/auth/me')) return new Response(JSON.stringify({ id: A.userId, tenant_id: A.tenantId }));
      const key = new Headers(init?.headers).get('Idempotency-Key')!;
      if (!server.has(key)) server.set(key, 'server-ticket');
      if (++submissions === 1) throw new TypeError('Response lost');
      return new Response(JSON.stringify({ id: server.get(key) }), { status: 201 });
    });
    expect(await syncPendingTickets(A)).toBe(0);
    expect((await getPendingTickets(A))[0].id).toBe(id);
    expect(await syncPendingTickets(A)).toBe(1);
    expect([...server.keys()]).toEqual([id]);
    expect(await getPendingTickets(A)).toEqual([]);
  });

  it('claims a draft atomically so two concurrent sync contexts submit it once', async () => {
    await savePendingTicket(DATA, A);
    const results = await Promise.all([syncPendingTickets(A), syncPendingTickets(A)]);
    expect(results.reduce((sum, n) => sum + n, 0)).toBe(1);
    expect(writes).toHaveLength(1);
  });

  it('counts a conflict as unsynced and preserves it for manual resolution', async () => {
    await savePendingTicket(DATA, A);
    fail = 409;
    expect(await syncPendingTickets(A)).toBe(0);
    const conflicts = await getConflicts(A);
    expect(conflicts).toHaveLength(1);
    fail = undefined;
    await resolveConflict(conflicts[0].id, 'local', A);
    expect(await getPendingTickets(A)).toEqual([]);
  });

  it('stops on unauthorized or forbidden responses and preserves later drafts', async () => {
    await savePendingTicket(DATA, A);
    await savePendingTicket({ ...DATA, title: 'Second ticket' }, A);
    fail = 403;
    await expect(syncPendingTickets(A)).rejects.toThrow('Failed request');
    expect(writes).toHaveLength(1);
    expect(await getPendingTickets(A)).toHaveLength(2);
  });

  it('cancels in-flight synchronization when the session ends without losing the draft', async () => {
    await savePendingTicket(DATA, A);
    let submitted!: () => void;
    const started = new Promise<void>((resolve) => { submitted = resolve; });
    vi.stubGlobal('fetch', async (url: string, init?: RequestInit) => {
      if (url.endsWith('/auth/me')) return new Response(JSON.stringify({ id: A.userId, tenant_id: A.tenantId }));
      submitted();
      return new Promise<Response>((_resolve, reject) => init?.signal?.addEventListener('abort', () => reject(new DOMException('Aborted', 'AbortError'))));
    });
    const running = syncPendingTickets(A);
    await started;
    cancelPendingSync();
    await running;
    expect(await getPendingTickets(A)).toHaveLength(1);
  });
});
