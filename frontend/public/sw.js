// Trickreport Service Worker — vanilla JS (not processed by Vite/Astro)
// Strategies:
//   - App shell (HTML/CSS/JS/fonts): cache-first with network update (stale-while-revalidate)
//   - API calls (/api/v1): network-first, fallback to cache
//   - Static assets (images, icons, fonts files): cache-only
//   - Background sync: 'sync-tickets' triggers a message to the page to run syncPendingTickets

const CACHE_VERSION = 'trickreport-v1';
const APP_SHELL_CACHE = `${CACHE_VERSION}-shell`;
const API_CACHE = `${CACHE_VERSION}-api`;
const STATIC_CACHE = `${CACHE_VERSION}-static`;

// Resources that form the app shell. Relative to the SW scope (root).
const APP_SHELL = [
  '/',
  '/tickets',
  '/tickets/new',
  '/tickets/new-offline',
  '/dashboard',
  '/login',
  '/manifest.json',
  '/icon.svg',
];

// ─── Install: pre-cache the app shell ────────────────────────────────
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(APP_SHELL_CACHE).then((cache) => {
      // Use addAll but tolerate individual failures (some routes may 404 in dev)
      return Promise.all(
        APP_SHELL.map((url) =>
          cache.add(new Request(url, { cache: 'reload' })).catch(() => {})
        )
      );
    })
  );
  self.skipWaiting();
});

// ─── Activate: clean up old caches ───────────────────────────────────
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys
          .filter((key) => !key.startsWith(CACHE_VERSION))
          .map((key) => caches.delete(key))
      )
    )
  );
  self.clients.claim();
});

// ─── Helpers ─────────────────────────────────────────────────────────
function isApiRequest(url) {
  return url.pathname.startsWith('/api/v1') || url.pathname.includes('/api/v1');
}

function isStaticAsset(url) {
  return /\.(?:png|jpg|jpeg|gif|webp|svg|ico|woff|woff2|ttf|eot|otf)$/i.test(url.pathname);
}

// ─── Fetch handler ───────────────────────────────────────────────────
self.addEventListener('fetch', (event) => {
  const request = event.request;
  const url = new URL(request.url);

  // Only handle GET for caching; let everything else pass through.
  if (request.method !== 'GET') {
    return;
  }

  // 1. API calls → network-first, fallback to cache
  if (isApiRequest(url)) {
    event.respondWith(
      fetch(request)
        .then((response) => {
          // Cache a copy of successful responses
          if (response.ok) {
            const clone = response.clone();
            caches.open(API_CACHE).then((cache) => cache.put(request, clone));
          }
          return response;
        })
        .catch(() => caches.match(request).then((cached) => cached || new Response('Offline', { status: 503, statusText: 'Offline' })))
    );
    return;
  }

  // 2. Static assets → cache-only (populated on first fetch via below fallback)
  if (isStaticAsset(url)) {
    event.respondWith(
      caches.match(request).then((cached) => {
        if (cached) return cached;
        // Not in cache yet — fetch and store, then return.
        return fetch(request)
          .then((response) => {
            if (response.ok) {
              const clone = response.clone();
              caches.open(STATIC_CACHE).then((cache) => cache.put(request, clone));
            }
            return response;
          })
          .catch(() => cached || Response.error());
      })
    );
    return;
  }

  // 3. Same-origin navigation / app shell → stale-while-revalidate
  if (url.origin === self.location.origin && request.mode === 'navigate') {
    event.respondWith(
      caches.match(request).then((cached) => {
        const networkFetch = fetch(request)
          .then((response) => {
            if (response.ok) {
              const clone = response.clone();
              caches.open(APP_SHELL_CACHE).then((cache) => cache.put(request, clone));
            }
            return response;
          })
          .catch(() => cached || caches.match('/'));
        return cached || networkFetch;
      })
    );
    return;
  }

  // 4. Other same-origin GET (CSS/JS chunks) → cache-first with background update
  if (url.origin === self.location.origin) {
    event.respondWith(
      caches.match(request).then((cached) => {
        const networkFetch = fetch(request)
          .then((response) => {
            if (response.ok) {
              const clone = response.clone();
              caches.open(APP_SHELL_CACHE).then((cache) => cache.put(request, clone));
            }
            return response;
          })
          .catch(() => cached);
        return cached || networkFetch;
      })
    );
    return;
  }

  // Cross-origin (e.g. Google Fonts) → cache-first, network fallback, cache the result.
  event.respondWith(
    caches.match(request).then((cached) => {
      if (cached) return cached;
      return fetch(request)
        .then((response) => {
          if (response.ok) {
            const clone = response.clone();
            caches.open(STATIC_CACHE).then((cache) => cache.put(request, clone));
          }
          return response;
        })
        .catch(() => cached);
    })
  );
});

// ─── Background Sync ─────────────────────────────────────────────────
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-tickets') {
    event.waitUntil(
      (async () => {
        const allClients = await self.clients.matchAll({ includeUncontrolled: true });
        allClients.forEach((client) => {
          client.postMessage({ type: 'SYNC_TICKETS' });
        });
      })()
    );
  }
});

// ─── Message handler (e.g. manual sync trigger from page) ────────────
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});
