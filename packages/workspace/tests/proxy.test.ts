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

describe('provider webhooks through the same-origin proxy', () => {
  beforeEach(() => {
    vi.stubEnv('APP_PUBLIC_URL', 'https://support.example.test');
    vi.stubEnv('API_URL', 'http://api:8080');
  });
  afterEach(() => vi.unstubAllGlobals());

  it('forwards a Zabbix delivery without Origin, with its token and the raw body', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 202 }));
    vi.stubGlobal('fetch', fetchMock);
    const body = '{"host":"srv-1","trigger":"CPU high"}';
    const request = new Request('https://support.example.test/api/v1/msp/monitoring/webhook', {
      method: 'POST',
      headers: { 'content-type': 'application/json', 'x-msp-webhook-token': 'source-token', 'x-tenant-id': 'acme' },
      body,
    });

    const response = await proxy(request, { params: Promise.resolve({ path: ['msp', 'monitoring', 'webhook'] }) }, 'msp');

    expect(response.status).toBe(202);
    const [target, init] = fetchMock.mock.calls[0];
    expect(target).toBe('http://api:8080/api/v1/msp/monitoring/webhook');
    const headers = new Headers(init.headers);
    expect(headers.get('x-msp-webhook-token')).toBe('source-token');
    expect(headers.get('x-tenant-id')).toBe('acme');
    expect(headers.has('origin')).toBe(false);
    expect(new TextDecoder().decode(init.body)).toBe(body);
  });

  it('forwards a WhatsApp delivery with its HMAC signature header intact', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 200 }));
    vi.stubGlobal('fetch', fetchMock);
    const request = new Request('https://support.example.test/api/v1/cx/whatsapp/webhook', {
      method: 'POST',
      headers: { 'content-type': 'application/json', 'x-hub-signature-256': 'sha256=abc123' },
      body: '{"entry":[]}',
    });

    const response = await proxy(request, { params: Promise.resolve({ path: ['cx', 'whatsapp', 'webhook'] }) }, 'customer-service');

    expect(response.status).toBe(200);
    const [target, init] = fetchMock.mock.calls[0];
    expect(target).toBe('http://api:8080/api/v1/cx/whatsapp/webhook');
    expect(new Headers(init.headers).get('x-hub-signature-256')).toBe('sha256=abc123');
  });

  it('keeps rejecting a write without Origin outside the webhook paths', async () => {
    const fetchMock = vi.fn();
    vi.stubGlobal('fetch', fetchMock);
    const request = new Request('https://support.example.test/api/v1/msp/monitoring/webhook/extra', {
      method: 'POST',
      body: '{}',
    });

    const response = await proxy(request, { params: Promise.resolve({ path: ['msp', 'monitoring', 'webhook', 'extra'] }) }, 'msp');

    expect(response.status).toBe(403);
    expect(fetchMock).not.toHaveBeenCalled();
  });
});
