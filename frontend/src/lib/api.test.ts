import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock the config module so API_URL is deterministic regardless of env.
vi.mock('./config', () => ({
  API_URL: 'http://localhost:8080',
  PUBLIC_API_URL: 'http://localhost:8080',
  ON_PREMISE: false,
  DEFAULT_TENANT_SLUG: 'default',
  resolveTenantSlug: () => 'default',
}));

// fetch is a global; vi.fn() is untyped so it can stand in for any call shape.
const fetchMock = vi.fn();
(globalThis as any).fetch = fetchMock;

import { api, login, updateTicketStatus, assignTicket, createTicket, upsertSLAPolicy, logout } from './api';

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

beforeEach(() => {
  fetchMock.mockReset();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('api client - URL construction', () => {
  it('builds the full URL from API_URL and path', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { ok: true }));
    await api('/tickets');
    const [url] = (fetchMock as any).mock.calls[0];
    expect(url).toBe('http://localhost:8080/api/v1/tickets');
  });

  it('preserves query strings in the path', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, []));
    await api('/articles?q=hello');
    const [url] = (fetchMock as any).mock.calls[0];
    expect(url).toBe('http://localhost:8080/api/v1/articles?q=hello');
  });

  it('login posts to /auth/login', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { token: 'abc', user: { id: '1' } }));
    await login({ email: 'a@b.com', password: 'secret123' });
    const [url, init] = (fetchMock as any).mock.calls[0];
    expect(url).toBe('http://localhost:8080/api/v1/auth/login');
    expect(init.method).toBe('POST');
  });
});

describe('api client - token handling', () => {
  it('sets Authorization header when token provided', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { ok: true }));
    await api('/tickets', { token: 'my-jwt' });
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.headers['Authorization']).toBe('Bearer my-jwt');
  });

  it('omits Authorization header when no token', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { ok: true }));
    await api('/tickets');
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.headers['Authorization']).toBeUndefined();
  });

  it('always sets X-Tenant-ID and Content-Type headers', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { ok: true }));
    await api('/tickets', { token: 't' });
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.headers['X-Tenant-ID']).toBe('default');
    expect(init.headers['Content-Type']).toBe('application/json');
  });

  it('includes credentials', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { ok: true }));
    await api('/tickets');
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.credentials).toBe('include');
  });
});

describe('api client - error handling', () => {
  it('throws with server error message on non-ok response', async () => {
    fetchMock.mockResolvedValue(jsonResponse(400, { error: 'bad request' }));
    await expect(api('/tickets')).rejects.toThrow('bad request');
  });

  it('throws with status text when body has no error field', async () => {
    fetchMock.mockResolvedValue(jsonResponse(500, {}));
    await expect(api('/tickets')).rejects.toThrow('Request failed with status 500');
  });

  it('throws with status text when body is not JSON', async () => {
    fetchMock.mockResolvedValue(
      new Response('not json', { status: 502, headers: { 'Content-Type': 'text/plain' } }),
    );
    await expect(api('/tickets')).rejects.toThrow('Request failed with status 502');
  });

  it('throws on 401 unauthorized', async () => {
    fetchMock.mockResolvedValue(jsonResponse(401, { error: 'invalid token' }));
    await expect(api('/tickets', { token: 'bad' })).rejects.toThrow('invalid token');
  });
});

describe('api client - ticket routes', () => {
  it('sends the displayed ticket version to detect concurrent status edits', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { status: 'resolved', version: 4 }));
    await updateTicketStatus('1', 'resolved', 't', 3);
    expect(JSON.parse(fetchMock.mock.calls[0][1].body)).toEqual({ status: 'resolved', version: 3 });
  });
  it('updateTicketStatus PATCHes /tickets/{id}/status', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { id: '1', status: 'resolved' }));
    await updateTicketStatus('1', 'resolved', 't');
    const [url, init] = (fetchMock as any).mock.calls[0];
    expect(url).toBe('http://localhost:8080/api/v1/tickets/1/status');
    expect(init.method).toBe('PATCH');
  });

  it('assignTicket POSTs /tickets/{id}/assign', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { id: '1', assigned_to: 'u1' }));
    await assignTicket('1', 'u1', 't');
    const [url, init] = (fetchMock as any).mock.calls[0];
    expect(url).toBe('http://localhost:8080/api/v1/tickets/1/assign');
    expect(init.method).toBe('POST');
  });
});

describe('api client - CSRF header', () => {
  afterEach(() => {
    (globalThis as any).document = undefined;
  });

  it('attaches X-CSRF-Token on mutating requests when the cookie is present', async () => {
    (globalThis as any).document = { cookie: 'trickreport_csrf=abc123; other=1' };
    fetchMock.mockResolvedValue(jsonResponse(201, { id: '1' }));
    await api('/tickets', { token: 't', method: 'POST', body: JSON.stringify({ title: 'x' }) });
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.headers['X-CSRF-Token']).toBe('abc123');
  });

  it('does not attach X-CSRF-Token on GET requests', async () => {
    (globalThis as any).document = { cookie: 'trickreport_csrf=abc123' };
    fetchMock.mockResolvedValue(jsonResponse(200, []));
    await api('/tickets', { token: 't' });
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.headers['X-CSRF-Token']).toBeUndefined();
  });

  it('omits X-CSRF-Token when there is no cookie to read (SSR context)', async () => {
    fetchMock.mockResolvedValue(jsonResponse(201, { id: '1' }));
    await api('/tickets', { token: 't', method: 'POST', body: JSON.stringify({ title: 'x' }) });
    const [, init] = (fetchMock as any).mock.calls[0];
    expect(init.headers['X-CSRF-Token']).toBeUndefined();
  });
});

describe('api client - response parsing', () => {
  it('returns parsed JSON on success', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, [{ id: '1' }, { id: '2' }]));
    const result = await api<any[]>('/tickets');
    expect(result).toHaveLength(2);
    expect(result[0].id).toBe('1');
  });

  it('returns undefined on 204 No Content', async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 204 }));
    const result = await api<void>('/notifications/1/read', { token: 't', method: 'POST' });
    expect(result).toBeUndefined();
  });
});

describe('SSR mutations and idempotency', () => {
  it('forwards only the explicit CSRF cookie and header alongside the Bearer session', async () => {
    fetchMock.mockResolvedValue(jsonResponse(201, { id: '1' }));
    await createTicket(
      { title: 'Printer', description: 'Paper is stuck', priority: 'high', category: 'hardware' },
      'session',
      { csrfToken: 'csrf-value', idempotencyKey: 'a1045a83-6d38-4c52-a435-91709dad4381' },
    );
    const [, init] = fetchMock.mock.calls[0];
    expect(init.headers).toMatchObject({
      Authorization: 'Bearer session',
      'X-CSRF-Token': 'csrf-value',
      Cookie: 'trickreport_csrf=csrf-value',
      'Idempotency-Key': 'a1045a83-6d38-4c52-a435-91709dad4381',
    });
  });

  it('carries explicit CSRF for SLA changes and server logout', async () => {
    fetchMock.mockImplementation(() => Promise.resolve(new Response(null, { status: 204 })));
    await upsertSLAPolicy('high', { response_time_minutes: 30, resolution_time_minutes: 60, escalation_minutes: 10 }, 'session', { csrfToken: 'csrf' });
    await logout('session', { csrfToken: 'csrf' });
    expect(fetchMock.mock.calls.map(([, init]) => [init.method, init.headers.Cookie, init.headers['X-CSRF-Token']])).toEqual([
      ['PUT', 'trickreport_csrf=csrf', 'csrf'],
      ['POST', 'trickreport_csrf=csrf', 'csrf'],
    ]);
  });

  it('does not repeat a non-idempotent mutation after a server failure', async () => {
    fetchMock.mockImplementation(() => Promise.resolve(jsonResponse(503, { error: 'Unavailable' })));
    await expect(api('/articles', { method: 'POST', body: '{}' })).rejects.toThrow('Unavailable');
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });
});
