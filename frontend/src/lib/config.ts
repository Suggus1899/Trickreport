/**
 * On-premise / self-hosted configuration.
 *
 * When ON_PREMISE is true the frontend expects to be served from the same
 * origin as the API, so API_URL is empty and all requests use relative paths
 * (e.g. `/api/v1/...`). No external resources (fonts, CDNs, analytics) are
 * loaded, keeping the deployment fully air-gapped.
 */
export const ON_PREMISE =
  import.meta.env.ON_PREMISE === 'true' ||
  import.meta.env.PUBLIC_ON_PREMISE === 'true';

/**
 * Base URL for the backend API.
 *
 * - On-premise: empty string → requests go to the same origin (`/api/v1/...`).
 * - Development/standalone: defaults to `http://localhost:8080`.
 */
// `??`, not `||`: an empty value is a deliberate choice meaning "same origin,
// use relative paths", and `||` would discard it for the fallback.
export const API_URL = import.meta.env.API_URL ?? (ON_PREMISE ? '' : 'http://localhost:8080');

/**
 * Client-safe copy of API_URL for React islands and other browser-side code.
 * Astro only inlines `PUBLIC_`-prefixed env vars into the client bundle, so
 * this must be set to the same value as API_URL in every environment.
 */
export const PUBLIC_API_URL =
  import.meta.env.PUBLIC_API_URL ?? (ON_PREMISE ? '' : 'http://localhost:8080');

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
