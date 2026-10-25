package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

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
var _ email.ContextSender = (*SMTPSender)(nil)

// Send delivers an email to the given recipient.
func (s *SMTPSender) Send(to, subject, body string) error {
	return s.SendContext(context.Background(), to, subject, body)
}

// SendContext bounds the entire exchange, including the greeting and TLS handshake.
func (s *SMTPSender) SendContext(ctx context.Context, to, subject, body string) error {
	for _, header := range []string{s.from, to, subject} {
		if strings.ContainsAny(header, "\r\n\x00") {
			return fmt.Errorf("smtp: invalid header")
		}
	}
	fromAddress, err := mail.ParseAddress(s.from)
	if err != nil {
		return fmt.Errorf("smtp: invalid from header: %w", err)
	}
	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("smtp: invalid to header: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	conn, err := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, "tcp", net.JoinHostPort(s.host, strconv.Itoa(s.port)))
	if err != nil {
		return fmt.Errorf("smtp: dial: %w", err)
	}
	defer func() { _ = conn.Close() }()
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		return err
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stop()
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("smtp: greeting: %w", err)
	}
	defer func() { _ = client.Close() }()
	if offered, _ := client.Extension("STARTTLS"); s.useTLS || offered {
		if err := client.StartTLS(&tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("smtp: starttls: %w", err)
		}
	}
	if s.username != "" {
		if err := client.Auth(smtp.PlainAuth("", s.username, s.password, s.host)); err != nil {
			return fmt.Errorf("smtp: auth: %w", err)
		}
	}
	if err := client.Mail(fromAddress.Address); err != nil {
		return fmt.Errorf("smtp: mail: %w", err)
	}
	if err := client.Rcpt(toAddress.Address); err != nil {
		return fmt.Errorf("smtp: rcpt: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp: data: %w", err)
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", fromAddress.String(), toAddress.String(), mime.QEncoding.Encode("UTF-8", subject), body)
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp: write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp: close: %w", err)
	}
	return client.Quit()
}
