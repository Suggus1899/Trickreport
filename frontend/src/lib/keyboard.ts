// Global keyboard shortcut handler.
// Supports single-key shortcuts and two-key "g <key>" sequences (Vim-style).
// Registers default navigation shortcuts and a help toggle.

export interface Shortcut {
  keys: string[]; // sequence, e.g. ['g', 'd'] or ['?']
  description: string;
  action: () => void;
}

const isTypingTarget = (el: EventTarget | null): boolean => {
  if (!(el instanceof HTMLElement)) return false;
  const tag = el.tagName.toLowerCase();
  return tag === 'input' || tag === 'textarea' || tag === 'select' || el.isContentEditable;
};

export class KeyboardShortcuts {
  private shortcuts: Shortcut[] = [];
  private buffer: string[] = [];
  private bufferTimer: ReturnType<typeof setTimeout> | null = null;
  private helpToggle: (() => void) | null = null;

  constructor() {
    if (typeof window === 'undefined') return;
    window.addEventListener('keydown', this.onKeyDown);
  }

  register(shortcut: Shortcut): void {
    this.shortcuts.push(shortcut);
  }

  setHelpToggle(fn: () => void): void {
    this.helpToggle = fn;
  }

  private onKeyDown = (e: KeyboardEvent): void => {
    // Always allow Escape to close help / overlays handled elsewhere.
    if (e.key === 'Escape') {
      this.buffer = [];
      return;
    }

    // "/" focuses search — handled by GlobalSearch via its own listener,
    // but we avoid swallowing it when typing.
    if (isTypingTarget(e.target)) return;
    if (e.metaKey || e.ctrlKey || e.altKey) return;

    const key = e.key.toLowerCase();

    // Help toggle on '?'
    if (key === '?' || (e.shiftKey && e.key === '/')) {
      e.preventDefault();
      this.helpToggle?.();
      this.buffer = [];
      return;
    }

    // Buffer the key for sequence matching.
    this.buffer.push(key);
    if (this.buffer.length > 2) this.buffer = this.buffer.slice(-2);

    if (this.bufferTimer) clearTimeout(this.bufferTimer);
    this.bufferTimer = setTimeout(() => {
      this.buffer = [];
    }, 800);

    // Try to match a shortcut whose sequence is a prefix of / equals the buffer.
    for (const s of this.shortcuts) {
      if (this.matches(s.keys, this.buffer)) {
        if (this.buffer.length === s.keys.length) {
          e.preventDefault();
          s.action();
          this.buffer = [];
          return;
        }
        // Partial match — keep buffering.
        return;
      }
    }

    // No match — reset.
    this.buffer = [];
  };

  private matches(keys: string[], buffer: string[]): boolean {
    if (buffer.length > keys.length) return false;
    for (let i = 0; i < buffer.length; i++) {
      if (buffer[i] !== keys[i]) return false;
    }
    return true;
  }

  destroy(): void {
    if (typeof window === 'undefined') return;
    window.removeEventListener('keydown', this.onKeyDown);
  }
}

let instance: KeyboardShortcuts | null = null;

export function getKeyboardShortcuts(): KeyboardShortcuts {
  if (!instance) {
    instance = new KeyboardShortcuts();
  }
  return instance;
}

/** Register the default navigation shortcuts. */
export function registerDefaultShortcuts(): KeyboardShortcuts {
  const ks = getKeyboardShortcuts();
  const go = (path: string) => () => {
    window.location.href = path;
  };
  ks.register({ keys: ['g', 'd'], description: 'Go to Dashboard', action: go('/dashboard') });
  ks.register({ keys: ['g', 't'], description: 'Go to Tickets', action: go('/tickets') });
  ks.register({ keys: ['g', 'a'], description: 'Go to Knowledge Base', action: go('/articles') });
  ks.register({ keys: ['g', 'u'], description: 'Go to Users (admin)', action: go('/admin/users') });
  ks.register({ keys: ['n'], description: 'New Ticket', action: go('/tickets/new') });
  return ks;
}
