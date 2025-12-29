# Protocol: HTTP/HTTPS

HTTP and HTTPS are the primary protocols supported by GoPunch. It provides fine-grained control over how requests are made.

## Schema Detection
Any target starting with `http://` or `https://` is treated as an HTTP check. If no scheme is provided, `https://` is prepended automatically.

## Success Criteria
By default, a request is considered **successful** if the status code is in the `2xx` or `3xx` range (200-399).

You can override this using the `--expect` flag or the `expected_codes` config option.

## Customization Options

### Methods
Supports all standard HTTP methods: `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS`.

### Headers
You can provide multiple custom headers.
- **CLI**: `-H "Authorization: Bearer my-token" -H "X-Custom: value"`
- **Config**: `"headers": { "User-Agent": "GoPunch" }`

### Request Body
For `POST` or `PUT` requests, you can send a string body using `--data` (`-d`). If data is provided, the `Content-Type` is set to `application/json` by default unless overridden in headers.

### Redirects
GoPunch follows up to 10 redirects by default. You can disable this behavior with `--follow=false` (`-L=false`), in which case a 3xx status code will still be considered a success (as it is < 400).

### TLS / SSL
For internal or development environments with self-signed certificates, use the `--insecure` (`-k`) flag to skip certificate verification.

## Performance
HTTP checks utilize a tuned `http.Transport` with:
- Disabled Keep-Alives (to ensure each check is a fresh connection).
- Optimized Dial and TLS Handshake timeouts.
- Dedicated concurrency control via semaphores.
