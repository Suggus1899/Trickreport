import { afterEach, describe, expect, it, vi } from 'vitest';
import { dailyVolume, weekOverWeekDelta } from './dashboard-metrics';
import type { Ticket } from './api';

const ticket = (created_at: string, resolved_at: string | null = null) => ({ id: created_at, tenant_id: 'tenant', title: 'Ticket', description: 'Details', status: 'resolved', priority: 'high', category: 'general', created_by: 'user', sla_breached: false, created_at, updated_at: created_at, resolved_at }) as Ticket;
afterEach(() => vi.useRealTimers());

describe('dashboard metric meaning', () => {
  it('groups timestamps by UTC date and fills days with no activity', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-10-22T12:00:00Z'));
    expect(dailyVolume([ticket('2026-10-21T23:30:00-04:00')], 3)).toEqual([
      { date: '2026-10-20', count: 0 }, { date: '2026-10-21', count: 0 }, { date: '2026-10-22', count: 1 },
    ]);
  });

  it('counts resolutions on resolved_at and leaves historical NULL dates unknown', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-10-22T12:00:00Z'));
    expect(dailyVolume([ticket('2026-10-20T10:00:00Z', '2026-10-22T10:00:00Z'), ticket('2026-10-21T10:00:00Z')], 3, 'resolved_at')).toEqual([
      { date: '2026-10-20', count: 0 }, { date: '2026-10-21', count: 0 }, { date: '2026-10-22', count: 1 },
    ]);
  });

  it('does not manufacture a percentage when the comparison period has no base', () => {
    expect(weekOverWeekDelta([{ date: '2026-10-21', count: 0 }, { date: '2026-10-22', count: 3 }])).toBeNull();
    expect(weekOverWeekDelta([{ date: '2026-10-21', count: 0 }, { date: '2026-10-22', count: 0 }])).toBeNull();
    expect(weekOverWeekDelta([{ date: '2026-10-21', count: 2 }, { date: '2026-10-22', count: 3 }])).toBe(50);
  });
});
