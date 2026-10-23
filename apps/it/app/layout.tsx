import './globals.css';
import type { Metadata } from 'next';
export const metadata: Metadata = { title: 'Trickreport · Soporte TI', description: 'Mesa de servicio para equipos de soporte TI.', icons: { icon: '/trickreport-mark.svg' }, robots: { index: false, follow: false } };
export default function Layout({ children }: { children: React.ReactNode }) { return <html lang="es" suppressHydrationWarning><body>{children}</body></html>; }
