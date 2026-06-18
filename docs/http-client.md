# Custom HTTP client & headers

[← Docs index](./README.md)

Inject your own `*http.Client` to add timeouts, tracing, retries, or SSRF
protection. If you do not set one, a client with `DefaultTimeout` (30s) is used.

```go
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

shopware.Config{
	BaseURL:     shopURL,
	Credentials: shopware.NewIntegrationCredentials(id, secret),
	HTTPClient: &http.Client{
		Timeout:   30 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	},
}
```

`Config.Headers` attaches headers to **every** request (token requests
included) — handy for a reverse-proxy auth token or a deployment-specific
header:

```go
shopware.Config{
	Headers: map[string]string{"x-internal-proxy-token": proxyToken},
}
```

## See also

- [Raw requests](./raw-requests.md) — per-request headers.
- [Token storage](./token-storage.md) — the token request uses this client too.
