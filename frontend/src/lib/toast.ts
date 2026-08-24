/**
 * Minimal framework-agnostic toast: plain DOM, no React state needed, so the
 * same function works from a vanilla `<script>` in an Astro page and from
 * inside a React island. Separate from NotificationBell's toasts, which
 * render real-time push notifications (title/body/link) — this one is for
 * quick "your action worked/failed" confirmation, replacing the silent
 * post-action re-renders (and the one native alert()) that used to be the
 * only feedback in the app.
 */

export type ToastVariant = 'success' | 'error' | 'info';

const VARIANT_CLASSES: Record<ToastVariant, string> = {
  success: 'border-success/30 bg-success/10 text-success',
  error: 'border-destructive/30 bg-destructive/10 text-destructive',
  info: 'border-border bg-card text-foreground',
};

function getRoot(): HTMLElement {
  let root = document.getElementById('trickreport-toast-root');
  if (!root) {
    root = document.createElement('div');
    root.id = 'trickreport-toast-root';
    root.className = 'fixed bottom-4 right-4 z-[300] flex flex-col gap-2 items-end pointer-events-none';
    document.body.appendChild(root);
  }
  return root;
}

export function toast(message: string, variant: ToastVariant = 'info', durationMs = 3200): void {
  if (typeof document === 'undefined') return;
  const root = getRoot();
  const el = document.createElement('div');
  el.className = `pointer-events-auto animate-in slide-in-from-bottom-2 fade-in-0 duration-200 rounded-md border px-4 py-2.5 text-sm font-medium shadow-lg ${VARIANT_CLASSES[variant]}`;
  el.textContent = message;
  root.appendChild(el);

  const remove = () => {
    el.classList.add('animate-out', 'fade-out-0', 'slide-out-to-right-2');
    el.addEventListener('animationend', () => el.remove(), { once: true });
    setTimeout(() => el.remove(), 300);
  };
  const timer = setTimeout(remove, durationMs);
  el.addEventListener('click', () => {
    clearTimeout(timer);
    remove();
  });
}
