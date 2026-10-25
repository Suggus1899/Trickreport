import { describe, expect, it } from 'vitest';
import { activeContractForClient, formatCents, isActiveContract, projectProgressPct } from '../src/msp';
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
