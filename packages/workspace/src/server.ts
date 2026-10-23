import { cookies } from 'next/headers';
import { backendURL } from './proxy';
import type { Product, User } from './contracts';

export async function getSession(product: Product): Promise<{ user: User | null; unavailable: boolean }> {
  const jar = await cookies();
  const token = jar.get('trickreport_token');
  if (!token) return { user: null, unavailable: false };
  try {
    const response = await fetch(`${backendURL(product)}/api/v1/auth/me`, { headers: { Cookie: `trickreport_token=${encodeURIComponent(token.value)}` }, cache: 'no-store', signal: AbortSignal.timeout(8_000) });
    if (response.status === 401 || response.status === 403) return { user: null, unavailable: false };
    if (!response.ok) return { user: null, unavailable: true };
    return { user: await response.json(), unavailable: false };
  } catch { return { user: null, unavailable: true }; }
}
