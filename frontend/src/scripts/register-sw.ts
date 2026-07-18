// Registers the service worker (/sw.js) when supported by the browser.
// Also wires up a message listener so the page can react to background-sync
// events posted by the service worker (e.g. trigger syncPendingTickets).
export function registerServiceWorker(): void {
  if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
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
  });
}

registerServiceWorker();
