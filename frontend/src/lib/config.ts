/**
 * On-premise / self-hosted configuration.
 *
 * Browser requests use the frontend origin; SSR uses the internal API origin.
 * Fonts and other public resources are bundled locally.
 */
export const ON_PREMISE =
  import.meta.env.ON_PREMISE === 'true' ||
  import.meta.env.PUBLIC_ON_PREMISE === 'true';

/**
 * Base URL for the backend API.
 *
 * The server reads its internal origin at runtime so one build can be promoted
 * between environments. Node fetch requires an absolute URL, also on-premise.
 */
export const API_URL = (import.meta.env.SSR ? process.env.API_URL : undefined) || import.meta.env.API_URL || 'http://localhost:8080';

/**
 * Empty means same-origin behind the reverse proxy or the development proxy.
 * This public build setting is independent of the server's private API_URL.
 */
export const PUBLIC_API_URL = import.meta.env.PUBLIC_API_URL ?? '';

/** Tenant used when the host carries no organization subdomain. */
export const DEFAULT_TENANT_SLUG = 'default';

/** Hostnames that never carry a tenant subdomain. */
const NON_TENANT_HOSTS = new Set(['www', 'app', 'localhost']);

/**
 * Organization the browser is acting on behalf of, taken from the subdomain
 * (`acme.trickreport.com` → `acme`).
 *
 * Only the public auth routes read this — once a JWT exists the backend takes
 * the tenant from the token claim and ignores the header entirely, so a wrong
 * value here can never move an authenticated request into another tenant.
 */
export function resolveTenantSlug(): string {
  if (typeof window === 'undefined') return DEFAULT_TENANT_SLUG;

  const [first, ...rest] = window.location.hostname.split('.');
  // Needs at least two labels after the first, otherwise there is no subdomain.
  if (rest.length < 2 || NON_TENANT_HOSTS.has(first)) return DEFAULT_TENANT_SLUG;
  return first;
}
