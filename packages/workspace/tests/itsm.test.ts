import { describe, expect, it } from 'vitest';
import { changeTransitions, problemTransitions } from '../src/contracts';

// The frontend only offers the transitions the Go state machines allow;
// the server remains the authority, but the UI must not offer dead ends.
describe('ITSM transitions (frontend mirror of the Go state machines)', () => {
  it('a draft change cannot jump straight to approved', () => {
    expect(changeTransitions.draft).not.toContain('approved');
    expect(changeTransitions.draft).toContain('planned');
  });

  it('a change must pass through awaiting_approval before approved', () => {
    expect(changeTransitions.planned).toContain('awaiting_approval');
    expect(changeTransitions.awaiting_approval).toContain('approved');
  });

  it('a completed change can only close', () => {
    expect(changeTransitions.completed).toEqual(['closed']);
  });

  it('a closed problem can only reopen to investigating', () => {
    expect(problemTransitions.closed).toEqual(['investigating']);
  });

  it('every transition target is a known state', () => {
    for (const next of Object.values(changeTransitions).flat()) {
      expect(Object.keys(changeTransitions)).toContain(next);
    }
    for (const next of Object.values(problemTransitions).flat()) {
      expect(Object.keys(problemTransitions)).toContain(next);
    }
  });
});
