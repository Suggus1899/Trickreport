import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock the config module so API_URL is deterministic regardless of env.
vi.mock('./config', () => ({
  API_URL: 'http://localhost:8080',
  ON_PREMISE: false,
}));

// fetch is a global; cast to allow vi.fn assignment.
const fetchMock = vi.fn() as unknown as typeof fetch;
(globalThis as any).fetch = fetchMock;

import { api, login } from './api';

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
