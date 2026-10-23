export class APIError extends Error { constructor(message: string, public status: number) { super(message); } }
const safeMethods = new Set(['GET', 'HEAD', 'OPTIONS']);
let csrfRequest: Promise<void> | undefined;
let refreshRequest: Promise<unknown> | undefined;
function csrf(): string { return document.cookie.split('; ').find(value => value.startsWith('trickreport_csrf='))?.slice('trickreport_csrf='.length) || ''; }

export async function api<T>(path: string, init: RequestInit = {}, tenant = 'default', refresh = true): Promise<T> {
  const method = (init.method || 'GET').toUpperCase();
  if (!safeMethods.has(method) && !csrf()) {
    csrfRequest ??= api<void>('/auth/csrf', {}, tenant, false).finally(() => { csrfRequest = undefined; });
    await csrfRequest;
    if (!csrf()) throw new APIError('CSRF protection unavailable. Reload the page.', 403);
  }
  const headers = new Headers(init.headers);
  headers.set('X-Tenant-ID', tenant);
  if (!safeMethods.has(method)) headers.set('X-CSRF-Token', decodeURIComponent(csrf()));
  if (init.body && !(init.body instanceof FormData)) headers.set('Content-Type', 'application/json');
  const response = await fetch(`/api/v1${path}`, { ...init, headers, credentials: 'same-origin', cache: 'no-store', signal: init.signal || AbortSignal.timeout(30_000) });
  if (response.status === 401 && refresh && !path.startsWith('/auth/')) {
    refreshRequest ??= api('/auth/refresh', { method: 'POST', body: '{}' }, tenant, false).finally(() => { refreshRequest = undefined; });
    try { await refreshRequest; } catch { window.dispatchEvent(new Event('trickreport:session-ended')); throw new APIError('Session expired', 401); }
    return api<T>(path, init, tenant, false);
  }
  if (response.status === 204) return undefined as T;
  const body = await response.json().catch(() => ({}));
  if (!response.ok) throw new APIError(body.error || `Request failed (${response.status})`, response.status);
  return body as T;
}

export const json = (method: string, body: unknown, idempotencyKey?: string): RequestInit => ({ method, body: JSON.stringify(body), headers: idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : undefined });
