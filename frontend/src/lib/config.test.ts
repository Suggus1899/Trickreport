import { describe, it, expect, afterEach, vi } from 'vitest';
import { resolveTenantSlug, DEFAULT_TENANT_SLUG } from './config';

/**
 * Runs under the `node` environment, so there is no window unless a test
 * installs one. That is exactly the SSR case, and it is asserted first.
 */
function withHostname(hostname: string) {
  vi.stubGlobal('window', { location: { hostname } });
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('resolveTenantSlug', () => {
  it('falls back to the default tenant server-side, where there is no window', () => {
    expect(resolveTenantSlug()).toBe(DEFAULT_TENANT_SLUG);
  });

  it('reads the organization from the subdomain', () => {
    withHostname('acme.trickreport.com');
    expect(resolveTenantSlug()).toBe('acme');
  });

  it.each(['www', 'app'])('treats %s as a non-tenant host', (host) => {
    withHostname(`${host}.trickreport.com`);
    expect(resolveTenantSlug()).toBe(DEFAULT_TENANT_SLUG);
  });

  it('falls back on an apex domain with no subdomain', () => {
    withHostname('trickreport.com');
    expect(resolveTenantSlug()).toBe(DEFAULT_TENANT_SLUG);
  });

  it('falls back on localhost during development', () => {
    withHostname('localhost');
    expect(resolveTenantSlug()).toBe(DEFAULT_TENANT_SLUG);
  });

  it('keeps only the first label on a deeper host', () => {
    withHostname('acme.eu.trickreport.com');
    expect(resolveTenantSlug()).toBe('acme');
  });
});
