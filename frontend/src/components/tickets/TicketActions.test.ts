import { describe, it, expect } from 'vitest';
import { STATUS_TRANSITIONS } from './TicketActions';

const ALL_STATUSES = ['open', 'in_progress', 'waiting_client', 'resolved', 'closed'];

describe('TicketActions status transition guard', () => {
  it('never allows transitioning a status to itself', () => {
    for (const [from, targets] of Object.entries(STATUS_TRANSITIONS)) {
      expect(targets).not.toContain(from);
    }
  });

  it('only lists known ticket statuses as transition targets', () => {
    for (const targets of Object.values(STATUS_TRANSITIONS)) {
      for (const target of targets) {
        expect(ALL_STATUSES).toContain(target);
      }
    }
  });

  it('has an entry for every known status', () => {
    for (const status of ALL_STATUSES) {
      expect(STATUS_TRANSITIONS).toHaveProperty(status);
    }
  });

  it('closed tickets can only be reopened to in_progress', () => {
    expect(STATUS_TRANSITIONS.closed).toEqual(['in_progress']);
  });
});
