# Error handling

[← Docs index](./README.md)

Non-2xx responses are returned as `*APIError`. When the body matches Shopware's
error envelope, `Detail` holds the first error's message.

```go
resp, err := client.Get(ctx, "/search/does-not-exist")
if err != nil {
	var apiErr *shopware.APIError
	if errors.As(err, &apiErr) {
		fmt.Println("status:", apiErr.StatusCode)
		fmt.Println("detail:", apiErr.Detail) // e.g. "Entity does-not-exist not found"
		fmt.Println("raw body:", apiErr.Body)
	}
	return err
}
```

On a `401`, the client clears the cached token and retries the request **once**
with a fresh token before surfacing the error, so an expired token does not leak
out as a failure.

## See also

- [Token storage](./token-storage.md) — how the cached token is invalidated.
