import { describe, expect, it } from 'vitest';
import { avgSatisfaction, isSLABreached, nextPriority } from '../src/cx';
import { conversationTransitions } from '../src/contracts';
import type { CXConversation, CXSurvey } from '../src/contracts';

function baseConversation(overrides: Partial<CXConversation> = {}): CXConversation {
  return {
    id: 'c1', tenant_id: 't', customer_id: 'cu', customer_name: 'Ana', channel: 'chat',
    subject: 'Ayuda', status: 'open', priority: 'medium', assigned_to: null, assignee_name: '',
    created_by: null, external_thread_id: '', sla_response_minutes: 15, sla_resolution_minutes: 240,
    response_due_at: null, resolution_due_at: null, first_response_at: null, escalated: false,
    created_at: '', updated_at: '', ...overrides,
  };
}

// The frontend only offers the transitions the Go state machines allow;
// the server remains the authority, but the UI must not offer dead ends.
describe('CX transitions (frontend mirror of the Go state machines)', () => {
  it('an open conversation cannot jump back to open, but can wait for the customer', () => {
    expect(conversationTransitions.open).toContain('waiting_customer');
    expect(conversationTransitions.open).not.toContain('open');
  });

  it('a resolved conversation reopens to in_progress only (plus closed)', () => {
    expect(conversationTransitions.resolved).toContain('in_progress');
    expect(conversationTransitions.resolved).not.toContain('waiting_customer');
  });

  it('every transition target is a known state', () => {
    for (const [states, transitions] of [[conversationTransitions, conversationTransitions]] as const) {
      for (const next of Object.values(transitions).flat()) {
        expect(Object.keys(states)).toContain(next);
      }
    }
  });
});

describe('CX frontend helpers', () => {
  it('a conversation inside its deadlines is not breached, and a resolved one never is', () => {
    const now = new Date();
    const soon = new Date(now.getTime() + 3600_000).toISOString();
    expect(isSLABreached(baseConversation({ response_due_at: soon, resolution_due_at: soon }), now)).toBe(false);
    expect(isSLABreached(baseConversation({ status: 'resolved', response_due_at: new Date(now.getTime() - 1000).toISOString(), resolution_due_at: new Date(now.getTime() - 1000).toISOString() }), now)).toBe(false);
  });

  it('a conversation past its first-response deadline with no agent reply is breached', () => {
    const now = new Date();
    const past = new Date(now.getTime() - 1000).toISOString();
    expect(isSLABreached(baseConversation({ response_due_at: past, resolution_due_at: new Date(now.getTime() + 3600_000).toISOString() }), now)).toBe(true);
  });

  it('escalation raises one priority level and critical stays critical', () => {
    expect(nextPriority('low')).toBe('medium');
    expect(nextPriority('medium')).toBe('high');
    expect(nextPriority('high')).toBe('critical');
    expect(nextPriority('critical')).toBe('critical');
  });

  it('average satisfaction is 0 without surveys and averages only the ratings given', () => {
    expect(avgSatisfaction([])).toBe(0);
    const surveys: CXSurvey[] = [
      { id: '1', tenant_id: 't', conversation_id: 'a', customer_id: 'cu', rating: 5, comment: '', created_at: '' },
      { id: '2', tenant_id: 't', conversation_id: 'b', customer_id: 'cu', rating: 3, comment: '', created_at: '' },
    ];
    expect(avgSatisfaction(surveys)).toBe(4);
  });
});
