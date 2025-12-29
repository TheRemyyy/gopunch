# Command: check

The `check` command performs a one-time health check on one or more targets. It is highly flexible and suitable for both interactive use and automated scripts.

## Usage

```bash
gopunch check [flags] <url> [url...]
```

## Supported Schemes

| Scheme | Example | Description |
| :--- | :--- | :--- |
| `http://` / `https://` | `https://api.com` | Standard web request. |
| `tcp://` | `tcp://localhost:5432` | Checks if a TCP port is open. |
| `dns://` | `dns://google.com` | Verifies that a domain resolves to IPs. |
| `ssl://` | `ssl://example.com:443` | Inspects SSL certificate validity and expiry. |

*If no scheme is provided, `https://` is assumed by default.*

## Global Flags

| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--timeout` | `-t` | `10` | Timeout in seconds for each request. |
| `--concurrency` | `-c` | `10` | Max number of concurrent requests. |
| `--retries` | `-r` | `0` | Number of retries on failure (with backoff). |
| `--format` | `-f` | `table` | Output format: `table`, `json`, `csv`, `minimal`. |
| `--quiet` | `-q` | `false` | If set, suppresses output and uses exit codes only. |

## HTTP Specific Flags

| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--method` | `-m` | `GET` | HTTP verb to use. |
| `--header` | `-H` | - | Custom headers (e.g., `-H "Auth: Bearer x"`). Can be repeated. |
| `--data` | `-d` | - | Request body for POST/PUT requests. |
| `--expect` | `-e` | - | List of allowed status codes (e.g., `-e 200,201`). |
| `--insecure` | `-k` | `false` | Skip TLS certificate verification. |
| `--follow` | `-L` | `true` | Follow HTTP redirects. |

## Examples

### 📊 Beautiful Table Output (Default)
```bash
gopunch check https://google.com https://github.com
```

### 🤖 Automation with JSON & JQ
```bash
gopunch check -f json https://api.myservice.com | jq '.[0].success'
```

### 🔒 Checking SSL Expiry
```bash
gopunch check ssl://expired.badssl.com
```

### 🛠️ Advanced API Testing
```bash
gopunch check https://api.test.com -m POST -H "Content-Type: application/json" -d '{"ping": "pong"}' -e 201
```

## Exit Codes

- `0`: All checks passed.
- `1`: One or more checks failed, or a system error occurred.
