package email

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type Sender interface {
	Send(to, subject, body string) error
}

// ConsoleSender logs the email to stdout instead of sending it. Good for dev.
type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) Send(to, subject, body string) error {
	divider := strings.Repeat("-", 60)
	fmt.Printf("\n%s\n[MOCK EMAIL SENT]\nTo: %s\nSubject: %s\nBody: \n%s\n%s\n", divider, to, subject, body, divider)
	log.Info().Str("to", to).Str("subject", subject).Msg("Mock email sent")
	return nil
}
