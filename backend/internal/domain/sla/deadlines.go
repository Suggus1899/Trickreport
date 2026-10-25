package sla

import "time"

// Deadlines are the absolute SLA clocks of one ticket. A nil clock means the
// policy does not define it, so no breach is ever reported for that clock.
type Deadlines struct {
	Response   *time.Time
	Resolution *time.Time
	Escalation *time.Time
}

// NewDeadlines computes the clocks for a ticket created at createdAt. A
// non-positive minute count leaves its clock unset; escalation is optional.
func NewDeadlines(createdAt time.Time, responseMinutes, resolutionMinutes, escalationMinutes int) Deadlines {
	return Deadlines{
		Response:   deadlineAfter(createdAt, responseMinutes),
		Resolution: deadlineAfter(createdAt, resolutionMinutes),
		Escalation: deadlineAfter(createdAt, escalationMinutes),
	}
}

// Validate reports whether the policy can drive the ticket clocks. Response
// and resolution are required; escalation may be unset but not negative.
func (p Policy) Validate() error {
	if p.ResponseTimeMinutes <= 0 || p.ResolutionTimeMinutes <= 0 || p.EscalationMinutes < 0 {
		return ErrValidation
	}
	return nil
}

func deadlineAfter(createdAt time.Time, minutes int) *time.Time {
	if minutes <= 0 {
		return nil
	}
	d := createdAt.Add(time.Duration(minutes) * time.Minute)
	return &d
}
