package auth

import (
	"sync"
	"time"
)

// TokenStore is the abstraction for blacklisting JWTs by their jti claim.
// Implementations must be safe for concurrent use.
type TokenStore interface {
	// Blacklist adds the given jti to the blacklist until expiresAt.
	Blacklist(jti string, expiresAt time.Time)
	// IsBlacklisted returns true if the jti is currently blacklisted.
	IsBlacklisted(jti string) bool
}

// blacklistedEntry holds the expiration time of a blacklisted token.
type blacklistedEntry struct {
	expiresAt time.Time
}

// MemoryTokenStore is an in-memory implementation of TokenStore. It runs a
// background goroutine that periodically removes expired entries.
type MemoryTokenStore struct {
	mu      sync.RWMutex
	entries map[string]blacklistedEntry
	stopCh  chan struct{}
}

// NewMemoryTokenStore creates a new MemoryTokenStore and starts the cleanup
// goroutine. Call Close to stop the goroutine.
func NewMemoryTokenStore() *MemoryTokenStore {
	s := &MemoryTokenStore{
		entries: make(map[string]blacklistedEntry),
		stopCh:  make(chan struct{}),
	}
	go s.cleanup()
	return s
}

// Blacklist adds the given jti to the blacklist until expiresAt.
func (s *MemoryTokenStore) Blacklist(jti string, expiresAt time.Time) {
	s.mu.Lock()
	s.entries[jti] = blacklistedEntry{expiresAt: expiresAt}
	s.mu.Unlock()
}

// IsBlacklisted returns true if the jti is currently blacklisted and has not
// yet expired.
func (s *MemoryTokenStore) IsBlacklisted(jti string) bool {
	s.mu.RLock()
	entry, ok := s.entries[jti]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	return time.Now().Before(entry.expiresAt)
}

// Close stops the background cleanup goroutine.
func (s *MemoryTokenStore) Close() {
	close(s.stopCh)
}

// cleanup periodically removes expired blacklist entries.
func (s *MemoryTokenStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.purgeExpired()
		}
	}
}

// purgeExpired removes all entries whose expiration has passed.
func (s *MemoryTokenStore) purgeExpired() {
	now := time.Now()
	s.mu.Lock()
	for jti, entry := range s.entries {
		if !now.Before(entry.expiresAt) {
			delete(s.entries, jti)
		}
	}
	s.mu.Unlock()
}

// Compile-time interface assertion.
var _ TokenStore = (*MemoryTokenStore)(nil)
