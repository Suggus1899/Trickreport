import type { Ticket } from './api';

export interface DailyPoint {
  date: string;
  count: number;
}

export interface StatusSlice {
  status: string;
  count: number;
}

const DAY_MS = 24 * 60 * 60 * 1000;
const STATUS_ORDER = ['open', 'in_progress', 'waiting_client', 'resolved', 'closed'];

function dayKey(iso: string): string {
  return iso.slice(0, 10);
}

/** Daily count of tickets created within the trailing `days` days, oldest first. */
export function dailyVolume(tickets: Ticket[], days: number): DailyPoint[] {
  const buckets = new Map<string, number>();
  const now = Date.now();
  for (let i = days - 1; i >= 0; i--) {
    buckets.set(dayKey(new Date(now - i * DAY_MS).toISOString()), 0);
  }
  for (const t of tickets) {
    const key = dayKey(t.created_at);
    if (buckets.has(key)) buckets.set(key, (buckets.get(key) ?? 0) + 1);
  }
  return [...buckets.entries()].map(([date, count]) => ({ date, count }));
}

/** Current snapshot count per status, in a fixed display order. Any status
 * not in STATUS_ORDER (unexpected/legacy value) is appended at the end
 * instead of being silently dropped from the total. */
export function statusDistribution(tickets: Ticket[]): StatusSlice[] {
  const counts = new Map(STATUS_ORDER.map((s) => [s, 0]));
  for (const t of tickets) counts.set(t.status, (counts.get(t.status) ?? 0) + 1);
  return [...counts.entries()].map(([status, count]) => ({ status, count }));
}

/** % change between the two halves of a daily series (e.g. this week vs last week). */
export function weekOverWeekDelta(points: DailyPoint[]): number {
  const half = Math.floor(points.length / 2);
  const prev = points.slice(0, half).reduce((s, p) => s + p.count, 0);
  const curr = points.slice(half).reduce((s, p) => s + p.count, 0);
  if (prev === 0) return curr > 0 ? 100 : 0;
  return Math.round(((curr - prev) / prev) * 100);
}

/** Tickets assigned to a user, most recently updated first. */
export function myTickets(tickets: Ticket[], userId: string, limit = 5): Ticket[] {
  return tickets
    .filter((t) => t.assigned_to === userId)
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, limit);
}
