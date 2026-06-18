# Token storage

[← Docs index](./README.md)

By default the client caches the OAuth token in-process for the lifetime of the
`Client`. You can swap the backend via `Config.TokenStorage` — for example to
share tokens across instances (Redis) or to disable caching entirely.

```go
// Default: process-local cache (same as not setting it).
shopware.Config{ TokenStorage: shopware.NewInMemoryTokenStorage() }

// No caching — fetch a fresh token on every request (tests / low traffic).
shopware.Config{ TokenStorage: shopware.NewNoOpTokenStorage() }
```

Implement `shopware.TokenStorage` for a distributed cache:

```go
type TokenStorage interface {
	Get(ctx context.Context, key string) (token string, expiry time.Time, err error)
	Set(ctx context.Context, key string, token string, expiry time.Time) error
	Delete(ctx context.Context, key string) error
}
```

## Storage keys

Tokens are keyed by the credentials' identity by default (e.g.
`integration:<client-id>` vs `password:<username>`), so distinct principals never
share a cached token. When several clients share one storage but talk to
**different shops** with the same principal, give each its own key:

```go
shopware.Config{
	TokenStorage:    sharedRedisStorage,
	TokenStorageKey: "shop-" + shopID,
}
```

## Expiry

The client stores each token's real expiry and applies a 30s safety margin when
reading it back, so a token is never used in the window where the server might
already reject it.

## See also

- [Authentication](./authentication.md) — how the cached token is obtained.
- [Concurrency](./concurrency.md) — concurrent token fetches are collapsed.
