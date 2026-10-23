import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { proxy } from '../src/proxy';

const params = { params: Promise.resolve({ path: ['auth', 'login'] }) };

describe('same-origin API proxy', () => {
  beforeEach(() => {
    vi.stubEnv('APP_PUBLIC_URL', 'https://support.example.test');
    vi.stubEnv('API_URL', 'http://api:8080');
  });
  afterEach(() => vi.unstubAllGlobals());

  it('removes credentials from auth JSON while preserving HttpOnly cookies', async () => {
    const upstream = new Response(JSON.stringify({
      token: 'access-secret',
      refresh_token: 'refresh-secret',
      user: { id: 'u-1', name: 'Ana' },
    }), { headers: { 'content-type': 'application/json' } });
    upstream.headers.append('set-cookie', 'trickreport_token=access-secret; HttpOnly; Path=/');
    upstream.headers.append('set-cookie', 'trickreport_refresh=refresh-secret; HttpOnly; Path=/');
    const fetchMock = vi.fn().mockResolvedValue(upstream);
    vi.stubGlobal('fetch', fetchMock);

    const request = new Request('https://support.example.test/api/v1/auth/login', {
      method: 'POST',
      headers: { origin: 'https://support.example.test', 'content-type': 'application/json' },
      body: JSON.stringify({ email: 'ana@example.test', password: 'never-forward-to-browser' }),
    });
    const response = await proxy(request, params, 'it');

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual({ user: { id: 'u-1', name: 'Ana' } });
    expect(response.headers.getSetCookie()).toHaveLength(2);
    expect(fetchMock.mock.calls[0][0]).toBe('http://api:8080/api/v1/auth/login');
  });

  it('rejects a cross-origin write before contacting the API', async () => {
    const fetchMock = vi.fn();
    vi.stubGlobal('fetch', fetchMock);
    const request = new Request('https://support.example.test/api/v1/auth/login', {
      method: 'POST',
      headers: { origin: 'https://attacker.example' },
      body: '{}',
    });

    const response = await proxy(request, params, 'it');
    expect(response.status).toBe(403);
    expect(fetchMock).not.toHaveBeenCalled();
  });
});
