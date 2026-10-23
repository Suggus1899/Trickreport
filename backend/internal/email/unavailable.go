package email

import "errors"

// UnavailableSender preserves queued mail when deployment SMTP is missing.
type UnavailableSender struct{}

func (UnavailableSender) Send(to, subject, body string) error {
	return errors.New("SMTP is not configured")
}
