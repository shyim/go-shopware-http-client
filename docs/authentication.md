# Authentication

[← Docs index](./README.md)

The grant type is selected by the `Credentials` you pass in `Config`. This
mirrors how a Shopware integration and the admin SPA authenticate differently
against the same `/api/oauth/token` endpoint.

The client authenticates lazily on the first request, caches the token, and
refreshes it transparently — you never manage tokens by hand. Call
`client.Authenticate(ctx)` if you want to verify the credentials up front.

## Integration (client credentials)

Use this for server-to-server access with an API integration's ID and secret.

```go
client := shopware.NewClient(shopware.Config{
	BaseURL:     "https://my-shop.example.com",
	Credentials: shopware.NewIntegrationCredentials("CLIENT_ID", "CLIENT_SECRET"),
})
```

`ClientID` / `ClientSecret` on `Config` are a shorthand for the same thing — if
you leave `Credentials` nil they are used as integration credentials:

```go
// Equivalent to the call above.
client := shopware.NewClient(shopware.Config{
	BaseURL:      "https://my-shop.example.com",
	ClientID:     "CLIENT_ID",
	ClientSecret: "CLIENT_SECRET",
})
```

## Admin user (username / password)

Use this to act as a real admin user — for example a CLI tool that logs in with
the same credentials a person types into the Shopware administration. The
`client_id` is fixed to `administration` and the `write` scope is requested,
matching the admin SPA.

```go
client := shopware.NewClient(shopware.Config{
	BaseURL:     "https://my-shop.example.com",
	Credentials: shopware.NewPasswordCredentials("admin", "shopware"),
})
```

## Custom credentials

Implement the `Credentials` interface to support any other token body. The
interface is intentionally small:

```go
type Credentials interface {
	tokenRequest() map[string]string // JSON body POSTed to /api/oauth/token
	identity() string                // stable key, used as the default token-storage key
}
```

> The interface methods are unexported, so custom implementations must live in
> this package. If you need a third grant from outside the package, open a PR
> adding it next to `IntegrationCredentials` / `PasswordCredentials`.

## See also

- [Token storage](./token-storage.md) — caching tokens, distributed backends.
- [Error handling](./error-handling.md) — the one-shot 401 re-auth retry.
