# Concurrency

[← Docs index](./README.md)

A `Client` is safe for concurrent use. Concurrent requests that all hit a cold
or expired token cache are collapsed (via `singleflight`) into a single token
fetch that they all share, so a burst of goroutines never stampedes the token
endpoint.

```go
var wg sync.WaitGroup
for _, path := range paths {
	wg.Add(1)
	go func(p string) {
		defer wg.Done()
		resp, err := client.Get(ctx, p) // one token fetch shared by all
		_ = resp
		_ = err
	}(path)
}
wg.Wait()
```

## See also

- [Token storage](./token-storage.md) — the cache the fetch collapses onto.
