import './globals.css';
import type { Metadata } from 'next';
export const metadata: Metadata = { title: 'Trickreport · Atención al cliente', description: 'Gestión de atención y casos de cliente.', icons: { icon: '/trickreport-mark.svg' }, robots: { index: false, follow: false } };
export default function Layout({ children }: { children: React.ReactNode }) { return <html lang="es" suppressHydrationWarning><body>{children}</body></html>; }
