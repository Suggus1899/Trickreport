import './globals.css';
import type { Metadata } from 'next';
export const metadata: Metadata = { title: 'Trickreport · MSP', description: 'Gestión de servicios administrados.', icons: { icon: '/trickreport-mark.svg' }, robots: { index: false, follow: false } };
export default function Layout({ children }: { children: React.ReactNode }) { return <html lang="es" suppressHydrationWarning><body>{children}</body></html>; }
