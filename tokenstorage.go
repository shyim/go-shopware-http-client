package shopware

import (
	"context"
	"sync"
	"time"
)

// TokenStorage stores and retrieves OAuth access tokens for a Client.
//
// Implementations must treat expired tokens as absent: Get returns an empty
// token once the stored expiry has passed, prompting the Client to fetch a
// fresh one. Implementations must be safe for concurrent use.
//
// Tokens are keyed by an opaque string chosen by the Client (its credentials'
// identity by default), so a single backend can serve many independent shops.
type TokenStorage interface {
	// Get retrieves the token for key. It returns an empty string and the zero
	// time when no valid (unexpired) token exists.
	Get(ctx context.Context, key string) (token string, expiry time.Time, err error)

	// Set stores token for key with the given expiry time.
	Set(ctx context.Context, key string, token string, expiry time.Time) error

	// Delete removes the token for key.
	Delete(ctx context.Context, key string) error
}

// NoOpTokenStorage is a TokenStorage that never caches. Every request triggers
// a fresh token fetch. Suitable for tests or callers that manage tokens
// elsewhere.
type NoOpTokenStorage struct{}

// NewNoOpTokenStorage creates a NoOpTokenStorage.
func NewNoOpTokenStorage() *NoOpTokenStorage { return &NoOpTokenStorage{} }

// Get always reports no cached token.
func (*NoOpTokenStorage) Get(context.Context, string) (string, time.Time, error) {
	return "", time.Time{}, nil
}

// Set discards the token.
func (*NoOpTokenStorage) Set(context.Context, string, string, time.Time) error { return nil }

// Delete is a no-op.
func (*NoOpTokenStorage) Delete(context.Context, string) error { return nil }

type tokenEntry struct {
	token  string
	expiry time.Time
}

// InMemoryTokenStorage is a process-local TokenStorage backed by a map. It is
// the default for the Client. Tokens are lost on restart and not shared across
// processes; for multi-instance deployments back the Client with a distributed
// store instead.
type InMemoryTokenStorage struct {
	mu     sync.RWMutex
	tokens map[string]tokenEntry
}

// NewInMemoryTokenStorage creates an empty in-memory token storage.
func NewInMemoryTokenStorage() *InMemoryTokenStorage {
	return &InMemoryTokenStorage{tokens: make(map[string]tokenEntry)}
}

// Get returns the cached token for key, or an empty token if absent or expired.
func (s *InMemoryTokenStorage) Get(_ context.Context, key string) (string, time.Time, error) {
	s.mu.RLock()
	e, ok := s.tokens[key]
	s.mu.RUnlock()

	if !ok || time.Now().After(e.expiry) {
		return "", time.Time{}, nil
	}
	return e.token, e.expiry, nil
}

// Set stores token for key with its expiry.
func (s *InMemoryTokenStorage) Set(_ context.Context, key, token string, expiry time.Time) error {
	s.mu.Lock()
	s.tokens[key] = tokenEntry{token: token, expiry: expiry}
	s.mu.Unlock()
	return nil
}

// Delete removes the token for key.
func (s *InMemoryTokenStorage) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	delete(s.tokens, key)
	s.mu.Unlock()
	return nil
}
