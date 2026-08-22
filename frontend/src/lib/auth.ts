import { getMe, type User } from './api';

/** Resolves the current session's user from the auth cookie, or null. */
export async function requireAuth(token: string | undefined): Promise<User | null> {
  if (!token) return null;
  try {
    return await getMe(token);
  } catch {
    return null;
  }
}
