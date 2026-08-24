import { useEffect, useRef, useState } from 'react';
import { getTickets, getArticles, type Ticket, type Article } from '@/lib/api';

interface ResultItem {
  href: string;
  title: string;
  sub: string;
}

function isTyping(el: EventTarget | null): boolean {
  if (!(el instanceof HTMLElement)) return false;
  const tag = el.tagName.toLowerCase();
  return tag === 'input' || tag === 'textarea' || tag === 'select' || el.isContentEditable;
}

function ResultRow({
  item,
  active,
  onHover,
}: {
  item: ResultItem;
  active: boolean;
  onHover: () => void;
}) {
  return (
    <a
      href={item.href}
      className={`flex flex-col gap-0.5 rounded-md px-2.5 py-1.5 text-sm no-underline ${
        active ? 'bg-accent text-accent-foreground' : 'hover:bg-accent hover:text-accent-foreground'
      }`}
      onMouseEnter={onHover}
    >
      <span className="font-semibold">{item.title}</span>
      <span className="text-xs text-muted-foreground">{item.sub}</span>
    </a>
  );
}

export function GlobalSearch() {
  const [query, setQuery] = useState('');
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [articles, setArticles] = useState<Article[]>([]);
  const [open, setOpen] = useState(false);
  const [activeIndex, setActiveIndex] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const rootRef = useRef<HTMLDivElement>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const searchSeqRef = useRef(0);

  useEffect(() => {
    function onSlash(e: KeyboardEvent) {
      if (e.key === '/' && !isTyping(e.target)) {
        e.preventDefault();
        inputRef.current?.focus();
        inputRef.current?.select();
      }
    }
    function onClickOutside(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener('keydown', onSlash);
    document.addEventListener('click', onClickOutside);
    return () => {
      document.removeEventListener('keydown', onSlash);
      document.removeEventListener('click', onClickOutside);
    };
  }, []);

  function handleChange(value: string) {
    setQuery(value);
    setActiveIndex(-1);
    clearTimeout(debounceRef.current);
    if (!value.trim()) {
      setOpen(false);
      return;
    }
    debounceRef.current = setTimeout(() => runSearch(value.trim()), 300);
  }

  async function runSearch(q: string) {
    const seq = ++searchSeqRef.current;
    const ql = q.toLowerCase();
    const [t, a] = await Promise.all([
      getTickets().catch(() => [] as Ticket[]),
      getArticles(q).catch(() => [] as Article[]),
    ]);
    // A newer search may have started (and possibly already resolved) while
    // this one was in flight — drop this stale response instead of
    // overwriting more recent results.
    if (seq !== searchSeqRef.current) return;
    setTickets(t.filter((x) => x.title.toLowerCase().includes(ql)).slice(0, 6));
    setArticles(a.slice(0, 6));
    setOpen(true);
  }

  const items: ResultItem[] = [
    ...tickets.map((t) => ({
      href: `/tickets/${t.id}`,
      title: t.title,
      sub: `${t.status} · ${t.priority}`,
    })),
    ...articles.map((a) => ({ href: `/articles/${a.id}`, title: a.title, sub: a.category })),
  ];

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActiveIndex((i) => (items.length ? (i + 1) % items.length : -1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActiveIndex((i) => (items.length ? (i - 1 + items.length) % items.length : -1));
    } else if (e.key === 'Enter') {
      if (activeIndex >= 0 && items[activeIndex]) {
        e.preventDefault();
        window.location.href = items[activeIndex].href;
      }
    } else if (e.key === 'Escape') {
      setOpen(false);
      inputRef.current?.blur();
    }
  }

  return (
    <div ref={rootRef} className="relative flex-1 min-w-[14rem] max-w-md">
      <svg
        className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <circle cx="11" cy="11" r="8" />
        <line x1="21" y1="21" x2="16.65" y2="16.65" />
      </svg>
      <input
        ref={inputRef}
        type="text"
        role="combobox"
        aria-expanded={open}
        aria-controls="gs-results"
        autoComplete="off"
        placeholder="Search tickets, articles…  (press /)"
        value={query}
        onChange={(e) => handleChange(e.target.value)}
        onKeyDown={handleKeyDown}
        onFocus={() => query.trim() && setOpen(true)}
        className="w-full h-9 rounded-md border bg-background pl-9 pr-10 text-sm shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
      />
      <kbd className="absolute right-2 top-1/2 -translate-y-1/2 rounded border bg-muted px-1.5 py-0.5 text-[0.65rem] text-muted-foreground">
        /
      </kbd>

      {open && (
        <div
          id="gs-results"
          role="listbox"
          className="absolute top-[calc(100%+0.5rem)] left-0 right-0 z-50 max-h-96 overflow-y-auto rounded-lg border bg-popover p-2 shadow-md"
        >
          {items.length === 0 ? (
            <div className="p-4 text-center text-sm text-muted-foreground">No results found.</div>
          ) : (
            <>
              {tickets.length > 0 && (
                <div className="px-2 pt-1 pb-1 text-[0.65rem] font-bold uppercase tracking-wide text-muted-foreground">
                  Tickets
                </div>
              )}
              {items.slice(0, tickets.length).map((item, i) => (
                <ResultRow
                  key={item.href}
                  item={item}
                  active={i === activeIndex}
                  onHover={() => setActiveIndex(i)}
                />
              ))}
              {articles.length > 0 && (
                <div className="px-2 pt-2 pb-1 text-[0.65rem] font-bold uppercase tracking-wide text-muted-foreground">
                  Articles
                </div>
              )}
              {items.slice(tickets.length).map((item, i) => {
                const idx = tickets.length + i;
                return (
                  <ResultRow
                    key={item.href}
                    item={item}
                    active={idx === activeIndex}
                    onHover={() => setActiveIndex(idx)}
                  />
                );
              })}
            </>
          )}
        </div>
      )}
    </div>
  );
}
