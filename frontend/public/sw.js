// Cache public assets only. Private HTML and API responses always use the network.
const CACHE_VERSION = 'trickreport-v2-public';
const PUBLIC_ASSETS = ['/offline.html', '/manifest.json', '/icon.svg'];

self.addEventListener('install', (event) => {
  event.waitUntil(caches.open(CACHE_VERSION).then((cache) => cache.addAll(PUBLIC_ASSETS)));
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil((async () => {
    const keys = await caches.keys();
    await Promise.all(keys.filter((key) => key.startsWith('trickreport-') && key !== CACHE_VERSION).map((key) => caches.delete(key)));
    await self.clients.claim();
    await notifyClientsOfUpdate();
  })());
});

function isStaticAsset(url) {
  return url.origin === self.location.origin && (url.pathname.startsWith('/_astro/') || PUBLIC_ASSETS.includes(url.pathname) || url.pathname === '/favicon.svg');
}

self.addEventListener('fetch', (event) => {
  const request = event.request;
  const url = new URL(request.url);
  if (request.method !== 'GET' || url.origin !== self.location.origin || url.pathname.startsWith('/api/')) return;

  if (request.mode === 'navigate') {
    event.respondWith(fetch(request).catch(async () => {
      const cache = await caches.open(CACHE_VERSION);
      return await cache.match('/offline.html') || new Response('Offline', { status: 503 });
    }));
    return;
  }

  if (isStaticAsset(url)) {
    event.respondWith((async () => {
      const cache = await caches.open(CACHE_VERSION);
      const cached = await cache.match(request);
      if (cached) return cached;
      const response = await fetch(request);
      if (response.ok) {
        try { await cache.put(request, response.clone()); } catch { /* Storage may be full; the public asset still loads. */ }
      }
      return response;
    })());
  }
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

// ─── Update notification ─────────────────────────────────────────────
// Notify pages when a new service worker has activated so the page can
// show a "New version available — reload" banner.
async function notifyClientsOfUpdate() {
  const allClients = await self.clients.matchAll({ includeUncontrolled: true });
  allClients.forEach((client) => {
    client.postMessage({ type: 'NEW_VERSION_CACHED' });
  });
}

// ─── Push notifications ──────────────────────────────────────────────
self.addEventListener('push', (event) => {
  let payload = { title: 'Trickreport', body: 'You have a new notification.' };
  try {
    if (event.data) {
      const text = event.data.text();
      payload = text ? JSON.parse(text) : payload;
    }
  } catch {
    // Non-JSON payload — use the raw text as the body.
    if (event.data) payload.body = event.data.text();
  }

  const title = payload.title || 'Trickreport';
  const options = {
    body: payload.body || '',
    icon: '/icon.svg',
    badge: '/icon.svg',
    data: payload.data || payload.url || { url: '/dashboard' },
    tag: payload.tag || 'trickreport-notification',
    renotify: !!payload.renotify,
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

// ─── Notification click ──────────────────────────────────────────────
self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  const data = event.notification.data || {};
  const targetUrl = (typeof data === 'string' ? data : data.url) || '/dashboard';

  event.waitUntil(
    (async () => {
      const allClients = await self.clients.matchAll({ type: 'window', includeUncontrolled: true });
      // Focus an existing tab if one is open to the target URL.
      for (const client of allClients) {
        if (client.url.includes(targetUrl) && 'focus' in client) {
          return client.focus();
        }
      }
      // Otherwise open a new window.
      if (self.clients.openWindow) {
        return self.clients.openWindow(targetUrl);
      }
    })()
  );
});
