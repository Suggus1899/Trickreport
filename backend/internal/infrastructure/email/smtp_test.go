package email

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNewSMTPSender(t *testing.T) {
	s := NewSMTPSender("smtp.example.com", 587, "user", "pass", "noreply@example.com", true)
	if s == nil {
		t.Fatal("expected non-nil sender")
	}
	if s.host != "smtp.example.com" {
		t.Errorf("host = %q", s.host)
	}
	if s.port != 587 {
		t.Errorf("port = %d", s.port)
	}
	if s.username != "user" {
		t.Errorf("username = %q", s.username)
	}
	if s.password != "pass" {
		t.Errorf("password = %q", s.password)
	}
	if s.from != "noreply@example.com" {
		t.Errorf("from = %q", s.from)
	}
	if !s.useTLS {
		t.Error("useTLS should be true")
	}
}

func TestNewSMTPSender_NoAuth(t *testing.T) {
	s := NewSMTPSender("localhost", 25, "", "", "noreply@localhost", false)
	if s == nil {
		t.Fatal("expected non-nil sender")
	}
	if s.username != "" {
		t.Errorf("username = %q, want empty", s.username)
	}
	if s.useTLS {
		t.Error("useTLS should be false")
	}
}

func TestSMTPSender_Send_DialError(t *testing.T) {
	// Use a non-existent host to trigger a dial error.
	s := NewSMTPSender("nonexistent.invalid", 1, "", "", "noreply@localhost", true)
	err := s.Send("to@example.com", "Test", "Body")
	if err == nil {
		t.Error("expected error for non-existent host")
	}
}

func TestSMTPSender_Send_NoTLS_SendMailError(t *testing.T) {
	// Use a non-existent host with useTLS=false to exercise the smtp.SendMail path.
	s := NewSMTPSender("nonexistent.invalid", 1, "", "", "noreply@localhost", false)
	err := s.Send("to@example.com", "Test", "Body")
	if err == nil {
		t.Error("expected error for non-existent host")
	}
}

// startMockSMTP starts a minimal SMTP server that responds to basic commands.
// It returns the listener's address and a cleanup function.
func startMockSMTP(t *testing.T) (string, func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return // listener closed
			}
			go handleMockSMTP(conn)
		}
	}()

	return ln.Addr().String(), func() { ln.Close() }
}

func handleMockSMTP(conn net.Conn) {
	defer conn.Close()
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	// Greeting
	fmt.Fprintf(w, "220 mock.smtp ESMTP\r\n")
	w.Flush()

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "EHLO") || strings.HasPrefix(line, "HELO"):
			// Multi-line EHLO response with extensions
			fmt.Fprintf(w, "250-mock.smtp\r\n")
			fmt.Fprintf(w, "250-AUTH PLAIN\r\n")
			fmt.Fprintf(w, "250 OK\r\n")
			w.Flush()

		case strings.HasPrefix(line, "AUTH"):
			// Accept any AUTH
			fmt.Fprintf(w, "235 Authentication successful\r\n")
			w.Flush()

		case strings.HasPrefix(line, "MAIL FROM"):
			fmt.Fprintf(w, "250 OK\r\n")
			w.Flush()

		case strings.HasPrefix(line, "RCPT TO"):
			fmt.Fprintf(w, "250 OK\r\n")
			w.Flush()

		case strings.HasPrefix(line, "DATA"):
			fmt.Fprintf(w, "354 Start mail input\r\n")
			w.Flush()
			// Read until ".\r\n"
			for {
				bodyLine, err := r.ReadString('\n')
				if err != nil {
					return
				}
				if bodyLine == ".\r\n" || bodyLine == ".\n" || strings.TrimSpace(bodyLine) == "." {
					break
				}
			}
			fmt.Fprintf(w, "250 OK message accepted\r\n")
			w.Flush()

		case strings.HasPrefix(line, "QUIT"):
			fmt.Fprintf(w, "221 Bye\r\n")
			w.Flush()
			return

		case strings.HasPrefix(line, "RSET"):
			fmt.Fprintf(w, "250 OK\r\n")
			w.Flush()

		case strings.HasPrefix(line, "NOOP"):
			fmt.Fprintf(w, "250 OK\r\n")
			w.Flush()

		default:
			fmt.Fprintf(w, "250 OK\r\n")
			w.Flush()
		}
	}
}

func TestSMTPSender_Send_NoTLS_Success(t *testing.T) {
	addr, cleanup := startMockSMTP(t)
	defer cleanup()

	host, portStr, _ := net.SplitHostPort(addr)
	port := 0
	fmt.Sscanf(portStr, "%d", &port)

	s := NewSMTPSender(host, port, "", "", "noreply@localhost", false)
	err := s.Send("to@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSMTPSender_Send_TLS_StartTLSError(t *testing.T) {
	addr, cleanup := startMockSMTP(t)
	defer cleanup()

	host, portStr, _ := net.SplitHostPort(addr)
	port := 0
	fmt.Sscanf(portStr, "%d", &port)

	// useTLS=true but the mock server doesn't support STARTTLS,
	// so StartTLS should fail.
	s := NewSMTPSender(host, port, "", "", "noreply@localhost", true)
	err := s.Send("to@example.com", "Test Subject", "Test Body")
	if err == nil {
		t.Error("expected error for StartTLS failure")
	}
}

func TestSMTPSender_Send_NoTLS_WithAuth_Success(t *testing.T) {
	addr, cleanup := startMockSMTP(t)
	defer cleanup()

	// Give the server a moment to start
	time.Sleep(10 * time.Millisecond)

	host, portStr, _ := net.SplitHostPort(addr)
	port := 0
	fmt.Sscanf(portStr, "%d", &port)

	s := NewSMTPSender(host, port, "user", "pass", "noreply@localhost", false)
	err := s.Send("to@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
