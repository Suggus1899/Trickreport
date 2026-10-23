import { readFileSync } from 'node:fs';
import { runInNewContext } from 'node:vm';
import { describe, expect, it } from 'vitest';

function worker() {
  const handlers = new Map<string, (event: Record<string, unknown>) => void>();
  const stores = new Map<string, Map<string, Response>>();
  let online = true;
  const cacheKey = (input: string | Request) => typeof input === 'string' ? input : new URL(input.url, 'https://app.test').pathname;
  const caches = {
    keys: async () => [...stores.keys()],
    delete: async (key: string) => stores.delete(key),
    open: async (key: string) => {
      if (!stores.has(key)) stores.set(key, new Map());
      const entries = stores.get(key)!;
      return {
        add: async (input: Request) => { entries.set(cacheKey(input), await fetch(input)); },
        addAll: async (urls: string[]) => { await Promise.all(urls.map(async (url) => { entries.set(url, await fetch(url)); })); },
        match: async (input: string | Request) => entries.get(cacheKey(input))?.clone(),
        put: async (input: Request, response: Response) => { entries.set(cacheKey(input), response); },
      };
    },
    match: async (input: string | Request) => {
      for (const entries of stores.values()) {
        const found = entries.get(cacheKey(input));
        if (found) return found.clone();
      }
    },
  };
  const self = {
    location: new URL('https://app.test/sw.js'),
    addEventListener: (name: string, handler: (event: Record<string, unknown>) => void) => handlers.set(name, handler),
    skipWaiting: () => {},
    clients: { claim: async () => {}, matchAll: async () => [] },
  };
  const fetch = async (input?: string | Request | { url: string }) => {
    if (!online) throw new TypeError('Offline');
    if (String(typeof input === 'string' ? input : input?.url).endsWith('/offline.html')) return new Response('Public offline fallback');
    return new Response('Private account A', { headers: { 'Content-Type': 'text/html' } });
  };
  class WorkerRequest extends Request {
    constructor(input: string, init?: RequestInit) { super(new URL(input, 'https://app.test'), init); }
  }
  runInNewContext(readFileSync(new URL('../../../public/sw.js', import.meta.url), 'utf8'), { self, caches, fetch, URL, Request: WorkerRequest, Response, Promise });
  const dispatch = async (type: string) => {
    let done = Promise.resolve();
    handlers.get(type)?.({ waitUntil: (promise: Promise<void>) => { done = promise; } });
    await done;
  };
  const visit = async (path: string, mode = 'navigate') => {
    let response: Promise<Response> | undefined;
    const req = { url: `https://app.test${path}`, method: 'GET', mode };
    handlers.get('fetch')?.({ request: req, respondWith: (promise: Promise<Response>) => { response = promise; } });
    return response ? response : fetch(req);
  };
  return { stores, dispatch, visit, offline: () => { online = false; } };
}

describe('service worker private data policy', () => {
  it('does not return private cached HTML or API after logout', async () => {
    const sw = worker();
    await sw.dispatch('install');
    expect(await (await sw.visit('/tickets')).text()).toBe('Private account A');
    await sw.visit('/api/v1/tickets', 'cors');
    await Promise.resolve();
    sw.offline();
    expect(await (await sw.visit('/tickets')).text()).toBe('Public offline fallback');
    await expect(sw.visit('/api/v1/tickets', 'cors')).rejects.toThrow('Offline');
    for (const entries of sw.stores.values()) {
      expect([...entries.keys()]).not.toContain('/tickets');
      expect([...entries.keys()]).not.toContain('/api/v1/tickets');
    }
  });

  it('purges legacy private caches during activation', async () => {
    const sw = worker();
    sw.stores.set('trickreport-v1-shell', new Map([['/tickets', new Response('Private A')]]));
    sw.stores.set('trickreport-v1-api', new Map([['/api/v1/tickets', new Response('Private A')]]));
    await sw.dispatch('install');
    await sw.dispatch('activate');
    expect(sw.stores.has('trickreport-v1-shell')).toBe(false);
    expect(sw.stores.has('trickreport-v1-api')).toBe(false);
  });
});
