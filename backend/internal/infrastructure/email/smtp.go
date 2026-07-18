package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/trickreport/backend/internal/email"
)

// SMTPSender sends emails over SMTP.
type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	useTLS   bool
}

// NewSMTPSender creates a new SMTPSender.
func NewSMTPSender(host string, port int, username, password, from string, useTLS bool) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		useTLS:   useTLS,
	}
}

// Compile-time assertion that SMTPSender implements email.Sender.
var _ email.Sender = (*SMTPSender)(nil)

// Send delivers an email to the given recipient.
func (s *SMTPSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Build the RFC 822 message.
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, to, subject, body,
	)

	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	if s.useTLS {
		// STARTTLS: connect plain, upgrade to TLS, then send.
		conn, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("smtp: dial: %w", err)
		}
		defer conn.Close()

		if err := conn.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
			return fmt.Errorf("smtp: starttls: %w", err)
		}

		if auth != nil {
			if err := conn.Auth(auth); err != nil {
				return fmt.Errorf("smtp: auth: %w", err)
			}
		}

		if err := conn.Mail(s.from); err != nil {
			return fmt.Errorf("smtp: mail: %w", err)
		}
		if err := conn.Rcpt(to); err != nil {
			return fmt.Errorf("smtp: rcpt: %w", err)
		}

		w, err := conn.Data()
		if err != nil {
			return fmt.Errorf("smtp: data: %w", err)
		}
		if _, err := w.Write([]byte(msg)); err != nil {
			return fmt.Errorf("smtp: write: %w", err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("smtp: close: %w", err)
		}

		return conn.Quit()
	}

	// Plaintext (or implicit TLS handled by the server).
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
