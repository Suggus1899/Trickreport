import { afterEach, describe, expect, it, vi } from 'vitest';
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import reactRenderer from '@astrojs/react/server.js';
import { JSDOM } from 'jsdom';
import Tickets from './tickets.astro';
import SLA from './admin/sla.astro';
import type { Ticket, User } from '../lib/api';

const user: User = { id: 'user-a', tenant_id: 'tenant-a', name: 'Admin', email: 'admin@test.test', role: 'admin', active: true, created_at: '2026-10-01T00:00:00Z' };
const hostile = 'Printer </script><div id="outside-marker">inert</div> " & <text>';
const ticket: Ticket = { id: 'ticket-a', tenant_id: 'tenant-a', title: hostile, description: 'Details for printer', status: 'open', priority: 'high', category: 'hardware', created_by: 'user-a', sla_breached: false, created_at: '2026-10-21T10:00:00Z', updated_at: '2026-10-21T10:00:00Z' };
afterEach(() => vi.unstubAllGlobals());

async function render(component: Parameters<InstanceType<typeof AstroContainer>['renderToString']>[0], url: string) {
  const container = await AstroContainer.create();
  container.addServerRenderer({ name: '@astrojs/react', renderer: reactRenderer });
  container.addClientRenderer({ name: '@astrojs/react', entrypoint: '@astrojs/react/client.js' });
  return container.renderToString(component, { locals: { user, token: 'session' }, request: new Request(url, { headers: { Cookie: 'trickreport_token=session; trickreport_csrf=csrf' } }) });
}

describe('untrusted data in real Astro output', () => {
  it('keeps the hostile ticket inert and exposes only filtered rows for export through an escaped HTML attribute', async () => {
    vi.stubGlobal('fetch', async (url: string) => new Response(JSON.stringify(url.endsWith('/auth/me') ? user : [ticket, { ...ticket, id: 'hidden', title: 'Monitor', description: 'Other category' }])));
    const html = await render(Tickets, 'https://app.test/tickets?q=printer');
    const doc = new JSDOM(html).window.document;
    expect(doc.querySelector('#outside-marker')).toBeNull();
    expect(doc.querySelector('tbody')?.textContent).toContain(hostile);
    const rows = JSON.parse(doc.querySelector<HTMLElement>('[data-tickets]')?.dataset.tickets ?? 'null');
    expect(rows).toEqual([ticket]);
  });

  it('keeps an untrusted SLA success parameter inert in the rendered page', async () => {
    vi.stubGlobal('fetch', async (url: string) => new Response(JSON.stringify(url.endsWith('/auth/me') ? user : [])));
    const html = await render(SLA, `https://app.test/admin/sla?saved=${encodeURIComponent(hostile)}`);
    const doc = new JSDOM(html).window.document;
    expect(doc.querySelector('#outside-marker')).toBeNull();
    expect(doc.querySelector<HTMLElement>('[data-saved]')?.dataset.saved).toBe(hostile);
  });
});
