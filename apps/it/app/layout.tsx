import './globals.css';
import type { Metadata } from 'next';
import { cookies } from 'next/headers';
export const metadata: Metadata = { title: 'Trickreport · Soporte TI', description: 'Mesa de servicio para equipos de soporte TI.', icons: { icon: '/trickreport-mark.svg' }, robots: { index: false, follow: false } };
export default async function Layout({ children }: { children: React.ReactNode }) {
  const jar = await cookies();
  const locale = jar.get('trickreport_locale')?.value === 'en' ? 'en' : 'es';
  const dark = jar.get('trickreport_theme')?.value === 'dark';
  return <html lang={locale} className={dark ? 'dark' : undefined} suppressHydrationWarning><body>{children}</body></html>;
}
