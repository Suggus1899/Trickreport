package email

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSMTPRejectsInjectedHeaders(t *testing.T) {
	for _, tc := range []struct{ name, from, to, subject string }{
		{"from", "from@example.test\r\nBcc: other@example.test", "to@example.test", "Hello"},
		{"to", "from@example.test", "to@example.test\nBcc: other@example.test", "Hello"},
		{"subject", "from@example.test", "to@example.test", "Hello\r\nBcc: other@example.test"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := NewSMTPSender("127.0.0.1", 1, "", "", tc.from, false).Send(tc.to, tc.subject, "body")
			if err == nil || !strings.Contains(err.Error(), "header") {
				t.Fatalf("expected header rejection before dialing, got %v", err)
			}
		})
	}
}

func TestSMTPContextCancelsStalledGreeting(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	host, port, _ := net.SplitHostPort(listener.Addr().String())
	portNumber, _ := strconv.Atoi(port)
	sender := NewSMTPSender(host, portNumber, "", "", "from@example.test", false)
	contextSender, ok := any(sender).(interface {
		SendContext(context.Context, string, string, string) error
	})
	if !ok {
		t.Fatal("SMTP must support context cancellation")
	}
	accepted := make(chan struct{})
	closed := make(chan struct{})
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			close(accepted)
			close(closed)
			return
		}
		defer conn.Close()
		close(accepted)
		_, _ = conn.Read(make([]byte, 1))
		close(closed)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- contextSender.SendContext(ctx, "to@example.test", "Hello", "body") }()
	select {
	case <-accepted:
	case <-time.After(2 * time.Second):
		t.Fatal("SMTP did not connect")
	}
	cancel()
	select {
	case sendErr := <-result:
		if sendErr == nil {
			t.Fatal("cancelled send succeeded")
		}
	case <-time.After(time.Second):
		t.Fatal("SMTP blocked after cancellation")
	}
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("SMTP left connection open")
	}
}
