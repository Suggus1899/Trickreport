import { Workspace } from '@trickreport/workspace';
import { getSession } from '@trickreport/workspace/server';
import { redirect } from 'next/navigation';
export const dynamic = 'force-dynamic';
export default async function Page({ params }: { params: Promise<{ path?: string[] }> }) {
 const { path = [] } = await params;
 const route = '/' + path.join('/');
 const session = await getSession('customer-service');
 const publicRoute = ['/', '/login', '/register', '/password-reset', '/password-reset/confirm'].includes(route);
 if (!publicRoute && !session.user && !session.unavailable) redirect('/login');
 if ((route === '/' || route === '/login') && session.user) redirect('/dashboard');
 return <Workspace product="customer-service" route={route} user={session.user} unavailable={session.unavailable} />;
}
