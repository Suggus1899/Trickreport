import { useEffect, useState } from 'react';

const SHORTCUTS = [
  { keys: ['g', 'd'], desc: 'Go to Dashboard' },
  { keys: ['g', 't'], desc: 'Go to Tickets' },
  { keys: ['g', 'a'], desc: 'Go to Knowledge Base' },
  { keys: ['g', 'u'], desc: 'Go to Users (admin)' },
  { keys: ['n'], desc: 'New Ticket' },
  { keys: ['/'], desc: 'Focus search' },
  { keys: ['?'], desc: 'Show this help' },
  { keys: ['Esc'], desc: 'Close dialogs / help' },
];

const GOTO: Record<string, string> = {
  'g,d': '/dashboard',
  'g,t': '/tickets',
  'g,a': '/articles',
  'g,u': '/admin/users',
};

function isTyping(el: EventTarget | null): boolean {
  if (!(el instanceof HTMLElement)) return false;
  const tag = el.tagName.toLowerCase();
  return tag === 'input' || tag === 'textarea' || tag === 'select' || el.isContentEditable;
}

export function KeyboardHelp() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    let buffer: string[] = [];
    let timer: ReturnType<typeof setTimeout> | undefined;

    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        buffer = [];
        setOpen(false);
        return;
      }
      if (isTyping(e.target) || e.metaKey || e.ctrlKey || e.altKey) return;
      const key = e.key.toLowerCase();

      if (key === '?' || (e.shiftKey && e.key === '/')) {
        e.preventDefault();
        setOpen((v) => !v);
        buffer = [];
        return;
      }

      buffer.push(key);
      if (buffer.length > 2) buffer = buffer.slice(-2);
      clearTimeout(timer);
      timer = setTimeout(() => {
        buffer = [];
      }, 800);

      if (buffer.length === 1 && key === 'n') {
        e.preventDefault();
        window.location.href = '/tickets/new';
        buffer = [];
        return;
      }
      const dest = GOTO[buffer.join(',')];
      if (dest) {
        e.preventDefault();
        window.location.href = dest;
        buffer = [];
      }
    }

    document.addEventListener('keydown', onKeyDown);
    return () => {
      document.removeEventListener('keydown', onKeyDown);
      clearTimeout(timer);
    };
  }, []);

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/40 p-4"
      onClick={() => setOpen(false)}
    >
      <div
        className="w-full max-w-lg rounded-xl border bg-card p-6 shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Keyboard shortcuts</h2>
          <button
            type="button"
            onClick={() => setOpen(false)}
            aria-label="Close"
            className="rounded-md border px-2.5 py-1 text-xs font-semibold text-muted-foreground hover:text-foreground"
          >
            Esc
          </button>
        </div>
        <div className="flex flex-col gap-1">
          {SHORTCUTS.map((s) => (
            <div key={s.desc} className="flex items-center justify-between gap-4 py-2 border-b last:border-0">
              <span className="text-sm">{s.desc}</span>
              <span className="flex items-center gap-1">
                {s.keys.map((k, i) => (
                  <span key={i} className="flex items-center gap-1">
                    <kbd className="inline-flex min-w-6 items-center justify-center rounded border bg-muted px-1.5 py-0.5 font-mono text-xs font-semibold">
                      {k}
                    </kbd>
                    {i < s.keys.length - 1 && <span className="text-xs text-muted-foreground">+</span>}
                  </span>
                ))}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
