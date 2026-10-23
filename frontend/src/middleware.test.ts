import { afterEach, describe, expect, it, vi } from 'vitest';
import { onRequest } from './middleware';

vi.mock('./lib/config', () => ({ API_URL: 'http://api.test', PUBLIC_API_URL: '', resolveTenantSlug: () => 'default' }));
const user = { id: 'user-a', tenant_id: 'tenant-a', name: 'Agent', role: 'agent', email: 'a@test.test', active: true, created_at: '2026-10-01T00:00:00Z' };

function context(path: string, token?: string, init?: RequestInit) {
  const jar = new Map(token ? [['trickreport_token', token], ['trickreport_refresh', 'refresh']] : []);
  return {
    url: new URL(path, 'https://app.test'),
    request: new Request(new URL(path, 'https://app.test'), init),
    locals: {} as App.Locals,
    cookies: {
      get: (name: string) => jar.has(name) ? { value: jar.get(name) } : undefined,
      delete: (name: string) => { jar.delete(name); },
    },
    redirect: (location: string) => new Response(null, { status: 302, headers: { Location: location } }),
    jar,
  };
}

afterEach(() => vi.unstubAllGlobals());

describe('SSR session recovery', () => {
  it('clears an invalid cookie before login can redirect back into a loop', async () => {
    vi.stubGlobal('fetch', async () => new Response('{"error":"invalid session"}', { status: 401 }));
    const ctx = context('/login', 'expired');
    const response = (await onRequest(ctx as never, async () => new Response(ctx.locals.user ? 'dashboard' : 'login'))) as Response;
    expect(await response.text()).toBe('login');
    expect(ctx.jar.size).toBe(0);
    expect(ctx.locals.token).toBeUndefined();
  });

  it('resolves a valid session once and makes it available to the dashboard', async () => {
    let resolutions = 0;
    vi.stubGlobal('fetch', async () => { resolutions++; return new Response(JSON.stringify(user)); });
    const ctx = context('/dashboard', 'valid');
    const response = (await onRequest(ctx as never, async () => new Response(ctx.locals.user?.name))) as Response;
    expect(await response.text()).toBe('Agent');
    expect(resolutions).toBe(1);
  });

  it('presents a recoverable 503 on a backend outage without clearing the session', async () => {
    vi.stubGlobal('fetch', async () => { throw new TypeError('network unavailable'); });
    const ctx = context('/dashboard', 'valid');
    const response = (await onRequest(ctx as never, async () => new Response('unexpected'))) as Response;
    expect(response.status).toBe(503);
    expect(response.headers.get('Location')).toBeNull();
    expect(await response.text()).toMatch(/try again|retry/i);
    expect(ctx.jar.get('trickreport_token')).toBe('valid');
  });

  it('rejects a cross-origin SSR form before forwarding its cookie-derived CSRF token', async () => {
    vi.stubGlobal('fetch', async () => new Response(JSON.stringify(user)));
    const ctx = context('/tickets/new', 'valid', { method: 'POST', headers: { Origin: 'https://outside.test' } });
    const response = (await onRequest(ctx as never, async () => new Response('mutation reached'))) as Response;
    expect(response.status).toBe(403);
  });
});
