import { describe, expect, it } from 'vitest';
import { activeContractForClient, contractPayload, formatCents, isActiveContract, projectProgressPct, toISO, toLocalInput } from '../src/msp';
import { interventionTransitions, projectTransitions, taskTransitions } from '../src/contracts';
import type { MSPContract } from '../src/contracts';

// The frontend only offers the transitions the Go state machines allow;
// the server remains the authority, but the UI must not offer dead ends.
describe('MSP transitions (frontend mirror of the Go state machines)', () => {
  it('a scheduled intervention cannot jump straight to done', () => {
    expect(interventionTransitions.scheduled).toContain('in_progress');
    expect(interventionTransitions.scheduled).not.toContain('done');
  });

  it('a planned project cannot jump straight to done', () => {
    expect(projectTransitions.planned).toContain('active');
    expect(projectTransitions.planned).not.toContain('done');
  });

  it('a todo task cannot jump straight to done', () => {
    expect(taskTransitions.todo).toContain('doing');
    expect(taskTransitions.todo).not.toContain('done');
  });

  it('every transition target is a known state', () => {
    for (const [states, transitions] of [[interventionTransitions, interventionTransitions], [projectTransitions, projectTransitions], [taskTransitions, taskTransitions]] as const) {
      for (const next of Object.values(transitions).flat()) {
        expect(Object.keys(states)).toContain(next);
      }
    }
  });
});

describe('MSP frontend helpers', () => {
  it('formats cents as euros', () => {
    expect(formatCents(3000)).toContain('30');
    expect(formatCents(0)).toContain('0');
  });

  it('project progress is 0 without tasks and counts only done tasks', () => {
    expect(projectProgressPct([])).toBe(0);
    expect(projectProgressPct([{ status: 'done' }, { status: 'doing' }, { status: 'done' }, { status: 'todo' }])).toBe(50);
  });

  it('only an active contract inside its window counts as current', () => {
    const now = new Date();
    const base: MSPContract = {
      id: 'c1', tenant_id: 't', client_id: 'cl', name: 'Soporte', status: 'active',
      starts_at: new Date(now.getTime() - 86400000).toISOString(),
      ends_at: new Date(now.getTime() + 86400000).toISOString(),
      coverage: '', currency: 'EUR', monthly_fee_cents: 10000, hourly_rate_cents: 5000,
      included_hours_monthly: 10, sla_response_minutes: 240, sla_resolution_minutes: 1440,
      sla_overrides: {}, created_at: '', updated_at: '',
    };
    expect(isActiveContract(base, now)).toBe(true);
    expect(isActiveContract({ ...base, status: 'suspended' }, now)).toBe(false);
    expect(isActiveContract({ ...base, ends_at: new Date(now.getTime() - 1000).toISOString() }, now)).toBe(false);
    expect(activeContractForClient([base, { ...base, id: 'c2', name: 'Viejo', starts_at: new Date(now.getTime() - 2 * 86400000).toISOString() }], 'cl', now)?.id).toBe('c1');
    expect(activeContractForClient([base], 'other', now)).toBeUndefined();
  });
});

// Pin the zone: the datetime-local round trip only shows up when the local
// offset is not zero. America/Caracas is UTC-4 with no DST since 2016.
process.env.TZ = 'America/Caracas';

describe('MSP contract form', () => {
  const stored: MSPContract = {
    id: 'c1', tenant_id: 't', client_id: 'cl', name: 'Anterior', status: 'active',
    starts_at: '2026-10-25T12:00:00.000Z', ends_at: null, coverage: '24x7 remoto', currency: 'USD',
    monthly_fee_cents: 10000, hourly_rate_cents: 5000, included_hours_monthly: 10,
    sla_response_minutes: 60, sla_resolution_minutes: 480,
    sla_overrides: { critical: { response_minutes: 15, resolution_minutes: 120 } },
    created_at: '', updated_at: '',
  };
  const form = (values: Record<string, string>) => ({ get: (name: string) => values[name] ?? null });
  const shown = { client: 'cl', name: 'Nuevo nombre', status: 'active', starts: '2026-10-25T08:00', fee: '150', rate: '60', included: '5', resp: '30', res: '240' };

  it('an edit keeps the stored fields the form does not show', () => {
    expect(contractPayload(form(shown), stored)).toMatchObject({
      name: 'Nuevo nombre', coverage: '24x7 remoto', currency: 'USD', sla_overrides: stored.sla_overrides,
    });
  });

  it('a new contract starts with the defaults the API expects', () => {
    expect(contractPayload(form(shown))).toMatchObject({ coverage: '', currency: 'EUR', sla_overrides: {} });
  });

  it('a datetime-local value is local wall-clock time and is stored as the same instant', () => {
    expect(toLocalInput(stored.starts_at)).toBe('2026-10-25T08:00');
    expect(toISO('2026-10-25T08:00')).toBe(stored.starts_at);
    expect(contractPayload(form(shown)).starts_at).toBe(stored.starts_at);
    expect(contractPayload(form({ ...shown, ends: '' })).ends_at).toBeNull();
  });
});
