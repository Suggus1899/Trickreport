import { defineMiddleware } from 'astro:middleware';
import { requireAuth } from './lib/auth';

// Routes reachable without a session. Everything else redirects to /login
// when there's no valid token — new pages read the resolved session from
// `Astro.locals.user` instead of repeating the cookie+getMe+redirect block.
const PUBLIC_PATHS = new Set(['/', '/login', '/register', '/logout', '/offline', '/404']);

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
  const token = context.cookies.get('trickreport_token')?.value;
  context.locals.token = token;
  context.locals.user = token ? await requireAuth(token) : null;

  if (!isPublic(context.url.pathname) && !context.locals.user) {
    return context.redirect('/login');
  }

  return next();
});
