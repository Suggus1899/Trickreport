import { defineMiddleware } from 'astro:middleware';
import { requireAuth } from './lib/auth';

// Routes reachable without a session. Everything else redirects to /login
// when there's no valid token — new pages read the resolved session from
// `Astro.locals.user` instead of repeating the cookie+getMe+redirect block.
const PUBLIC_PATHS = new Set([
  '/',
  '/login',
  '/register',
  '/logout',
  '/offline',
  '/404',
  // Password reset is reached precisely when the user cannot sign in.
  '/password-reset',
  '/password-reset/confirm',
]);

function isPublic(pathname: string): boolean {
  if (PUBLIC_PATHS.has(pathname)) return true;
  return (
    pathname.startsWith('/_astro/') ||
    pathname === '/sw.js' ||
    pathname === '/manifest.json' ||
    pathname === '/icon.svg' ||
    pathname === '/favicon.svg'
  );
}

export const onRequest = defineMiddleware(async (context, next) => {
  if (!['GET', 'HEAD', 'OPTIONS'].includes(context.request.method)) {
    if (context.request.headers.get('Origin') !== context.url.origin) return new Response('Invalid request origin', { status: 403 });
  }
  const token = context.cookies.get('trickreport_token')?.value;
  context.locals.token = token;
  context.locals.user = null;
  const pathname = context.url.pathname;
  if (pathname === '/logout' || pathname === '/offline.html' || pathname.startsWith('/_astro/') || /\.(?:js|json|svg|ico|woff2?|css)$/.test(pathname)) return next();
  try {
    context.locals.user = token ? await requireAuth(token) : null;
  } catch {
    return new Response('<!doctype html><html lang="en"><meta charset="utf-8"><title>Service unavailable</title><body><main><h1>Service temporarily unavailable</h1><p>Your session has been kept. Please try again shortly.</p><a href="">Try again</a></main></body></html>', {
      status: 503,
      headers: { 'Content-Type': 'text/html; charset=utf-8', 'Retry-After': '5', 'Cache-Control': 'no-store' },
    });
  }
  if (token && !context.locals.user) {
    context.cookies.delete('trickreport_token', { path: '/' });
    context.cookies.delete('trickreport_refresh', { path: '/' });
    context.locals.token = undefined;
  }

  if (!isPublic(context.url.pathname) && !context.locals.user) {
    return context.redirect('/login');
  }

  const response = await next();
  response.headers.set('Cache-Control', 'no-store');
  return response;
});
