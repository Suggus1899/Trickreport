// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('./config', () => ({ API_URL: 'http://localhost:8080', PUBLIC_API_URL: '', resolveTenantSlug: () => 'default' }));
import { api } from './api';

beforeEach(() => {
  document.cookie = 'trickreport_csrf=; Max-Age=0; path=/';
});
afterEach(() => vi.unstubAllGlobals());

describe('browser CSRF bootstrap', () => {
  it('initializes CSRF once before concurrent first mutations without setting Cookie manually', async () => {
    const requests: Array<{ url: string; init: RequestInit }> = [];
    vi.stubGlobal('fetch', async (url: string, init: RequestInit) => {
      requests.push({ url, init });
      if (url.endsWith('/auth/csrf')) {
        await Promise.resolve();
        document.cookie = 'trickreport_csrf=stable-token; path=/';
        return new Response(null, { status: 204 });
      }
      return new Response('{}', { status: 200 });
    });
    await Promise.all([api('/auth/login', { method: 'POST', body: '{}' }), api('/auth/register', { method: 'POST', body: '{}' })]);
    expect(requests.map((r) => r.url)).toEqual(['/api/v1/auth/csrf', '/api/v1/auth/login', '/api/v1/auth/register']);
    for (const { init } of requests.slice(1)) {
      expect(new Headers(init.headers).get('X-CSRF-Token')).toBe('stable-token');
      expect(new Headers(init.headers).has('Cookie')).toBe(false);
      expect(init.credentials).toBe('include');
    }
  });

  it('does not submit when the bootstrap cannot establish a CSRF cookie', async () => {
    const paths: string[] = [];
    vi.stubGlobal('fetch', async (url: string) => {
      paths.push(url);
      return new Response(null, { status: 204 });
    });
    await expect(api('/tickets', { method: 'POST', body: '{}' })).rejects.toThrow(/CSRF/i);
    expect(paths).toEqual(['/api/v1/auth/csrf']);
  });
});
