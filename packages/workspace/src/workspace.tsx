'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { BarChart3, BookOpen, Clock3, LayoutDashboard, LogOut, Menu, Moon, Sun, Ticket, Users, Workflow, X, type LucideIcon } from 'lucide-react';
import type { Locale, Product, User } from './contracts';
import { Analytics, Automations, SLAPolicies, Users as UserManagement } from './admin';
import { AuthPage, Profile } from './auth';
import { ArticleForm, ArticlePage, Knowledge } from './knowledge';
import { api } from './api';
import { Notifications } from './notifications';
import { Dashboard, Drafts, NewTicket, TicketDetail, TicketList } from './tickets';
import { Button, Empty, Feedback, Language, Notice, useText } from './ui';

type WorkspaceProps = { product: Product; route: string; user: User | null; unavailable: boolean };
type NavigationItem = { href: string; label: [string, string]; icon: LucideIcon; admin?: boolean };

const LOCALE_COOKIE = 'trickreport_locale';
const THEME_COOKIE = 'trickreport_theme';

function readPreference(name: string): string | null {
  if (typeof document === 'undefined') return null;
  const cookie = document.cookie.split('; ').find(value => value.startsWith(`${name}=`))?.slice(name.length + 1);
  if (cookie) return cookie;
  try { return window.localStorage.getItem(name); } catch { return null; }
}

function writePreference(name: string, value: string) {
  try { window.localStorage.setItem(name, value); } catch { /* storage unavailable */ }
  document.cookie = `${name}=${value}; path=/; max-age=31536000; samesite=lax`;
}

const productLabels: Record<Product, [string, string]> = {
  it: ['Soporte TI', 'IT support'],
  msp: ['Servicios MSP', 'MSP services'],
  'customer-service': ['Atención al cliente', 'Customer service'],
};

const navigation: NavigationItem[] = [
  { href: '/dashboard', label: ['Inicio', 'Overview'], icon: LayoutDashboard },
  { href: '/tickets', label: ['Solicitudes', 'Requests'], icon: Ticket },
  { href: '/articles', label: ['Conocimiento', 'Knowledge'], icon: BookOpen },
  { href: '/admin/users', label: ['Personas', 'People'], icon: Users, admin: true },
  { href: '/admin/sla', label: ['Acuerdos SLA', 'Service targets'], icon: Clock3, admin: true },
  { href: '/admin/automations', label: ['Automatizaciones', 'Automations'], icon: Workflow, admin: true },
  { href: '/admin/analytics', label: ['Informes', 'Reports'], icon: BarChart3, admin: true },
];

function PageContent({ product, route, user }: Required<Pick<WorkspaceProps, 'product' | 'route'>> & { user: User }) {
  const t = useText();
  if (route === '/dashboard') return <Dashboard user={user} />;
  if (route === '/tickets') return <TicketList />;
  if (route === '/tickets/new') return <NewTicket user={user} product={product} />;
  if (route === '/tickets/drafts') return <Drafts user={user} product={product} />;
  if (route.startsWith('/tickets/')) return <TicketDetail id={route.slice('/tickets/'.length)} user={user} />;
  if (route === '/articles') return <Knowledge user={user} />;
  if (route === '/articles/new') return user.role === 'end_user' ? <AccessDenied /> : <ArticleForm />;
  if (route.startsWith('/articles/')) {
    const parts = route.slice('/articles/'.length).split('/');
    if (parts.length === 2 && parts[1] === 'edit') return user.role === 'end_user' ? <AccessDenied /> : <ArticlePage id={parts[0]} user={user} edit />;
    if (parts.length === 1) return <ArticlePage id={parts[0]} user={user} />;
  }
  if (route === '/profile') return <Profile user={user} />;
  if (route === '/admin/users') return user.role === 'admin' ? <UserManagement /> : <AccessDenied />;
  if (route === '/admin/sla') return user.role === 'admin' ? <SLAPolicies /> : <AccessDenied />;
  if (route === '/admin/automations') return user.role === 'admin' ? <Automations /> : <AccessDenied />;
  if (route === '/admin/analytics') return user.role === 'admin' ? <Analytics /> : <AccessDenied />;
  return <section className="error-page"><Empty>{t('No encontramos esta página.','We could not find this page.')}</Empty></section>;
}

function AccessDenied() {
  const t = useText();
  return <section className="error-page"><Notice error>{t('No tienes acceso a esta sección.','You do not have access to this section.')}</Notice></section>;
}

export function Workspace({ product, route, user, unavailable }: WorkspaceProps) {
  const [locale, setLocale] = useState<Locale>(() => readPreference(LOCALE_COOKIE) === 'en' ? 'en' : 'es');
  const [dark, setDark] = useState(() => readPreference(THEME_COOKIE) === 'dark' || (typeof document !== 'undefined' && document.documentElement.classList.contains('dark')));
  const [menuState, setMenuState] = useState({ route, open: false });
  const [logoutError, setLogoutError] = useState('');
  const router = useRouter();
  const productName = productLabels[product];
  const menuOpen = menuState.route === route && menuState.open;
  const setMenuOpen = (open: boolean) => setMenuState({ route, open });

  useEffect(() => { document.documentElement.classList.toggle('dark', dark); }, [dark]);
  useEffect(() => { document.documentElement.lang = locale; }, [locale]);

  function changeLocale(value: Locale) {
    setLocale(value);
    writePreference(LOCALE_COOKIE, value);
  }
  function changeTheme() {
    const next = !dark;
    setDark(next);
    writePreference(THEME_COOKIE, next ? 'dark' : 'light');
  }
  async function logout() {
    setLogoutError('');
    try {
      await api('/auth/logout', { method: 'POST' }, user?.tenant_id || 'default', false);
      router.replace('/login');
      router.refresh();
    } catch (error) {
      setLogoutError(error instanceof Error ? error.message : 'Could not sign out');
    }
  }

  if (!user) {
    return <Language.Provider value={locale}>
      <div className="public-topbar"><Link className="brand" href="/"><Image className="brand-symbol" src="/trickreport-mark.svg" alt="" aria-hidden="true" width={34} height={34} unoptimized/>trickreport</Link>
        <div className="topbar-tools"><span className="product-tag">{locale === 'es' ? productName[0] : productName[1]}</span><LocalePicker locale={locale} changeLocale={changeLocale} /></div>
      </div>
      {unavailable && <div className="content"><Notice error>{locale === 'es' ? 'La API no está disponible ahora. Puedes reintentar en unos minutos.' : 'The API is unavailable right now. Please try again in a few minutes.'}</Notice></div>}
      <AuthPage route={route} product={product} />
    </Language.Provider>;
  }

  return <Language.Provider value={locale}>
    <a href="#main" className="skip-link">{locale === 'es' ? 'Saltar al contenido' : 'Skip to content'}</a>
    <div className="shell">
      {menuOpen && <button className="menu-backdrop" aria-label={locale === 'es' ? 'Cerrar menú' : 'Close menu'} onClick={() => setMenuOpen(false)} />}
      <aside className={`sidebar ${menuOpen ? 'is-open' : ''}`} id="primary-navigation">
        <div><Link className="brand" href="/dashboard" aria-label="Trickreport"><Image className="brand-symbol" src="/trickreport-mark.svg" alt="" aria-hidden="true" width={34} height={34} unoptimized/>trickreport</Link><p className="product-tag">{locale === 'es' ? productName[0] : productName[1]}</p></div>
        <nav className="nav" aria-label={locale === 'es' ? 'Navegación principal' : 'Primary navigation'}>
          <p className="nav-label">{locale === 'es' ? 'ESPACIO DE TRABAJO' : 'WORKSPACE'}</p>
          {navigation.filter(item => !item.admin || user.role === 'admin').map(({ href, label, icon: Icon }) => <Link key={href} href={href} onClick={() => setMenuOpen(false)} className={route === href || (href !== '/dashboard' && route.startsWith(`${href}/`)) ? 'active' : ''} aria-current={route === href ? 'page' : undefined}><Icon size={17} aria-hidden="true"/><span>{locale === 'es' ? label[0] : label[1]}</span></Link>)}
          <p className="nav-label nav-label-secondary">{locale === 'es' ? 'CUENTA' : 'ACCOUNT'}</p>
          <Link href="/profile" onClick={() => setMenuOpen(false)} className={route === '/profile' ? 'active' : ''}><Users size={17} aria-hidden="true"/><span>{locale === 'es' ? 'Mi perfil' : 'My profile'}</span></Link>
        </nav>
        <div className="sidebar-bottom"><div className="user-row"><span className="avatar" aria-hidden="true">{user.name.trim().slice(0, 1).toUpperCase()}</span><div><strong>{user.name}</strong><small>{user.email}</small></div></div><Button variant="secondary" className="full" onPress={logout}><LogOut size={15}/>{locale === 'es' ? 'Cerrar sesión' : 'Sign out'}</Button>{logoutError && <Feedback error={logoutError} success=""/>}</div>
      </aside>
      <div className="main-area">
        <header className="topbar"><div className="topbar-label"><span>{locale === 'es' ? 'Organización' : 'Organization'}</span><span>{user.tenant_id}</span></div>
          <div className="topbar-tools"><Button variant="tertiary" className="icon-button mobile-menu" aria-label={menuOpen ? 'Close navigation' : 'Open navigation'} aria-expanded={menuOpen} aria-controls="primary-navigation" onPress={() => setMenuOpen(!menuOpen)}>{menuOpen ? <X size={18}/> : <Menu size={18}/>}</Button><Notifications tenant={user.tenant_id}/><LocalePicker locale={locale} changeLocale={changeLocale}/><Button variant="tertiary" className="icon-button" onPress={changeTheme} aria-label={locale === 'es' ? (dark ? 'Cambiar a tema claro' : 'Cambiar a tema oscuro') : (dark ? 'Switch to light theme' : 'Switch to dark theme')}>{dark ? <Sun size={17}/> : <Moon size={17}/>}</Button><span className="avatar" aria-hidden="true">{user.name.trim().slice(0, 1).toUpperCase()}</span></div>
        </header>
        <main id="main" className="content"><PageContent product={product} route={route} user={user}/></main>
      </div>
    </div>
  </Language.Provider>;
}

function LocalePicker({ locale, changeLocale }: { locale: Locale; changeLocale: (locale: Locale) => void }) {
  return <label className="locale-picker"><span className="sr-only">Idioma / Language</span><select className="locale" value={locale} onChange={event => changeLocale(event.target.value as Locale)}><option value="es">ES</option><option value="en">EN</option></select></label>;
}
