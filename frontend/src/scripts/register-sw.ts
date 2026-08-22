// Registers the service worker (/sw.js) when supported by the browser.
// Also wires up a message listener so the page can react to background-sync
// events posted by the service worker (e.g. trigger syncPendingTickets).
export function registerServiceWorker(): void {
  if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
    return;
  }

  // El SW cachea JS same-origin con estrategia cache-first, pensada para
  // assets con hash estable. Contra el dev server sirve módulos de Vite viejos
  // — se lo vio devolver HTML donde se esperaba un módulo ("Failed to load
  // module script"), rompiendo la hidratación. Fuera de producción no
  // registramos nada y desinstalamos el que haya quedado de una sesión previa.
  if (!import.meta.env.PROD) {
    void unregisterInDev();
    return;
  }

  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js', { scope: '/' })
      .then((registration) => {
        // Listen for updates and activate the new SW immediately.
        registration.addEventListener('updatefound', () => {
          const installing = registration.installing;
          if (!installing) return;
          installing.addEventListener('statechange', () => {
            if (installing.state === 'installed' && navigator.serviceWorker.controller) {
              // A new version is available — let it take over on next navigation.
              installing.postMessage({ type: 'SKIP_WAITING' });
            }
          });
        });
      })
      .catch((err) => {
        // SW registration failed — app still works online-only.
        console.warn('Service worker registration failed:', err);
      });
  });

  // Relay background-sync messages from the SW to the page.
  navigator.serviceWorker.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'SYNC_TICKETS') {
      window.dispatchEvent(new CustomEvent('trickreport:sync-tickets'));
    }
    if (event.data && event.data.type === 'NEW_VERSION_CACHED') {
      window.dispatchEvent(new CustomEvent('trickreport:update-available'));
      showUpdateBanner();
    }
  });

  // Also detect controller change (new SW took over).
  navigator.serviceWorker.addEventListener('controllerchange', () => {
    showUpdateBanner();
  });
}

async function unregisterInDev(): Promise<void> {
  const registrations = await navigator.serviceWorker.getRegistrations();
  await Promise.all(registrations.map((r) => r.unregister()));
  if ('caches' in window) {
    const keys = await caches.keys();
    await Promise.all(keys.filter((k) => k.startsWith('trickreport-')).map((k) => caches.delete(k)));
  }
}

function showUpdateBanner(): void {
  if (document.getElementById('trickreport-update-banner')) return;
  const banner = document.createElement('div');
  banner.id = 'trickreport-update-banner';
  banner.className = 'update-banner';
  banner.innerHTML =
    '<span>A new version of Trickreport is available.</span>' +
    '<button type="button">Reload</button>';
  document.body.appendChild(banner);
  banner.querySelector('button')?.addEventListener('click', () => {
    window.location.reload();
  });
}

registerServiceWorker();
