import type { Product } from './contracts';

const ports = { it: 8080, msp: 8081, 'customer-service': 8082 };
export function backendURL(product: Product): string {
  return process.env.API_URL || `http://127.0.0.1:${ports[product]}`;
}

export function stripTokens(body: Record<string, unknown>): Record<string, unknown> {
  const publicBody = { ...body };
  delete publicBody.token;
  delete publicBody.refresh_token;
  return publicBody;
}

const forwardedHeaders = ['cookie', 'content-type', 'accept', 'x-csrf-token', 'idempotency-key', 'x-tenant-id', 'user-agent'];
// Provider callbacks are authenticated by their own token or HMAC signature, not by a browser Origin.
const providerWebhooks = new Set(['msp/monitoring/webhook', 'cx/whatsapp/webhook']);
const webhookHeaders = ['x-msp-webhook-token', 'x-hub-signature-256'];
const isProviderWebhook = (path: string[]) => providerWebhooks.has(path.join('/'));

/** Same-origin browser boundary. API credentials stay in HttpOnly cookies. */
export async function proxy(request: Request, context: { params: Promise<{ path: string[] }> }, product: Product): Promise<Response> {
  const url = new URL(request.url);
  const safe = ['GET', 'HEAD', 'OPTIONS'].includes(request.method);
  const origin = request.headers.get('origin');
  const expectedOrigin = process.env.APP_PUBLIC_URL ? new URL(process.env.APP_PUBLIC_URL).origin : url.origin;
  const { path } = await context.params;
  const webhook = isProviderWebhook(path);
  if (!safe && !webhook && origin !== expectedOrigin) return Response.json({ error: 'Invalid request origin' }, { status: 403 });
  if (!path.length || path.some(segment => segment === '.' || segment === '..' || /[\\/]/.test(segment))) return new Response(null, { status: 400 });
  const headers = new Headers();
  for (const key of webhook ? [...forwardedHeaders, ...webhookHeaders] : forwardedHeaders) {
    const value = request.headers.get(key);
    if (value) headers.set(key, value);
  }
  if (origin) headers.set('origin', origin);
  const target = `${backendURL(product)}/api/v1/${path.map(encodeURIComponent).join('/')}${url.search}`;
  let upstream: Response;
  try {
    upstream = await fetch(target, { method: request.method, headers, body: safe ? undefined : await request.arrayBuffer(), cache: 'no-store', redirect: 'manual', signal: AbortSignal.timeout(30_000) });
  } catch {
    return Response.json({ error: 'API unavailable. Please retry shortly.' }, { status: 503, headers: { 'Cache-Control': 'no-store', 'Retry-After': '5' } });
  }
  const responseHeaders = new Headers({ 'Cache-Control': 'no-store', 'X-Content-Type-Options': 'nosniff' });
  for (const name of ['content-type', 'content-disposition', 'retry-after']) {
    const value = upstream.headers.get(name);
    if (value) responseHeaders.set(name, value);
  }
  for (const cookie of upstream.headers.getSetCookie()) responseHeaders.append('set-cookie', cookie);
  if (path[0] === 'auth' && upstream.headers.get('content-type')?.includes('application/json')) {
    return Response.json(stripTokens(await upstream.json()), { status: upstream.status, headers: responseHeaders });
  }
  return new Response(upstream.body, { status: upstream.status, headers: responseHeaders });
}
