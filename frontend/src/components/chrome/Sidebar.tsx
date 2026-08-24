import { useEffect, useState } from 'react';
import { LayoutDashboard, Ticket, BookOpen, ShieldCheck, Menu, X, Sun, Moon, LogOut } from 'lucide-react';
import { logout, type User } from '@/lib/api';
import { useOnlineStatus } from '@/lib/hooks/useOnlineStatus';

const NAV = [
  { label: 'Dashboard', href: '/dashboard', roles: ['admin', 'agent', 'end_user'], icon: LayoutDashboard },
  { label: 'Tickets', href: '/tickets', roles: ['admin', 'agent', 'end_user'], icon: Ticket },
  { label: 'Knowledge Base', href: '/articles', roles: ['admin', 'agent', 'end_user'], icon: BookOpen },
  { label: 'Admin', href: '/admin', roles: ['admin'], icon: ShieldCheck },
] as const;

export function Sidebar({ user, pathname }: { user: Pick<User, 'name' | 'role'>; pathname: string }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [dark, setDark] = useState(false);
  const online = useOnlineStatus();

  useEffect(() => {
    setDark(document.documentElement.classList.contains('dark'));
  }, []);

  function toggleTheme() {
    const next = !dark;
    setDark(next);
    document.documentElement.classList.toggle('dark', next);
    document.documentElement.classList.toggle('light', !next);
    localStorage.setItem('trickreport-theme', next ? 'dark' : 'light');
  }

  async function handleLogout() {
    try {
      await logout();
    } catch {
      // Best-effort — clear the client view regardless of API outcome.
    }
    window.location.href = '/login';
  }

  const visibleNav = NAV.filter((item) => (item.roles as readonly string[]).includes(user.role));

  return (
    <>
      <div className="md:hidden">
        <button
          type="button"
          className="fixed top-4 left-4 z-40 inline-flex items-center justify-center w-10 h-10 rounded-full border bg-card shadow-sm transition-transform active:scale-95"
          aria-label="Toggle navigation"
          aria-expanded={mobileOpen}
          aria-controls="app-sidebar"
          onClick={() => setMobileOpen((v) => !v)}
        >
          {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
        </button>
        {mobileOpen && (
          <div
            className="fixed inset-0 z-30 bg-black/40 animate-in fade-in-0 duration-150"
            onClick={() => setMobileOpen(false)}
            aria-hidden="true"
          />
        )}
      </div>

      <aside
        id="app-sidebar"
        className={`fixed md:sticky top-0 left-0 z-30 h-screen w-64 shrink-0 border-r bg-card p-6 flex flex-col gap-6 transition-transform md:translate-x-0 ${
          mobileOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="flex items-center justify-between px-1">
          <a href="/dashboard" className="flex items-center gap-2 no-underline">
            <span className="inline-flex items-center justify-center size-7 rounded-md bg-primary/10 text-primary">
              <Ticket className="size-4" aria-hidden="true" />
            </span>
            <span className="font-heading font-bold text-lg tracking-tight uppercase text-foreground">
              Trickreport
            </span>
          </a>
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-semibold transition-colors ${
              online ? 'bg-success/10 text-success' : 'bg-destructive/10 text-destructive'
            }`}
          >
            <span className="w-1.5 h-1.5 rounded-full bg-current" aria-hidden="true" />
            {online ? 'Online' : 'Offline'}
          </span>
        </div>

        <nav className="flex flex-col gap-1" aria-label="Main navigation">
          {visibleNav.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(`${item.href}/`);
            const Icon = item.icon;
            return (
              <a
                key={item.href}
                href={item.href}
                aria-current={isActive ? 'page' : undefined}
                className={`relative flex items-center gap-2.5 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  isActive ? 'bg-primary/10 text-primary' : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                }`}
              >
                {isActive && (
                  <span
                    aria-hidden="true"
                    className="absolute left-0 top-1/2 -translate-y-1/2 h-4 w-0.5 rounded-full bg-primary"
                  />
                )}
                <Icon className="size-4 shrink-0" aria-hidden="true" />
                {item.label}
              </a>
            );
          })}
        </nav>

        <div className="mt-auto flex flex-col gap-3">
          <div className="flex items-center justify-between px-1">
            <button
              type="button"
              onClick={handleLogout}
              className="inline-flex items-center gap-1.5 text-sm font-semibold text-muted-foreground hover:text-foreground transition-colors"
            >
              <LogOut className="size-3.5" aria-hidden="true" />
              Logout
            </button>
            <button
              type="button"
              onClick={toggleTheme}
              aria-label="Toggle dark mode"
              className="inline-flex items-center justify-center w-9 h-9 rounded-full border bg-background hover:bg-muted transition-colors active:scale-95"
            >
              {dark ? <Moon className="size-4" /> : <Sun className="size-4" />}
            </button>
          </div>
          <a
            href={user.role === 'admin' ? '/admin/profile' : '/profile'}
            className="rounded-md border bg-muted/50 p-3 hover:bg-muted transition-colors no-underline"
          >
            <p className="font-semibold text-sm">{user.name}</p>
            <p className="text-xs text-muted-foreground capitalize">{user.role.replace('_', ' ')}</p>
          </a>
        </div>
      </aside>
    </>
  );
}
