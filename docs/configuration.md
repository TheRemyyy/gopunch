# Configuration Reference

GoPunch is highly configurable. You can define your monitoring targets and global settings in a `gopunch.json` file.

## Configuration File Structure

The file uses a flat JSON structure for global settings and a nested object for alerting.

### Full Example (`gopunch.json`)

```json
{
  "urls": [
    "https://example.com",
    "https://api.example.com/health",
    "tcp://db.local:5432",
    "ssl://mysite.com:443"
  ],
  "interval": 30,
  "timeout": 10,
  "method": "GET",
  "headers": {
    "User-Agent": "GoPunch/2.0",
    "X-Monitor": "Internal"
  },
  "insecure": false,
  "follow": true,
  "concurrency": 10,
  "retries": 2,
  "expected_codes": [200, 201, 204],
  "alerting": {
    "enabled": true,
    "cooldown": 300,
    "webhook": {
      "url": "https://discord.com/api/webhooks/YOUR_ID",
      "method": "POST"
    }
  }
}
```

## Options Breakdown

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `urls` | `[]string` | `[]` | List of URLs/targets to check. |
| `interval` | `int` | `5` | Seconds between cycles in `watch` mode. |
| `timeout` | `int` | `10` | Request timeout in seconds. |
| `method` | `string` | `"GET"` | HTTP method for checks. |
| `headers` | `map` | `{}` | Custom HTTP headers. |
| `insecure` | `bool` | `false` | Skip TLS certificate verification. |
| `follow` | `bool` | `true` | Follow HTTP redirects. |
| `concurrency`| `int` | `10` | Max parallel requests. |
| `retries` | `int` | `0` | Retries on failure with exponential backoff. |
| `expected_codes` | `[]int` | `200-399`| Status codes treated as success. |

## Precedence Rules

GoPunch resolves settings in the following order:
1.  **Command Line Flags**: Highest priority. If you run `gopunch watch --interval 10`, it will ignore the 30s interval in your config file.
2.  **Configuration File**: Loaded if `gopunch.json` exists or is specified via `--config`.
3.  **Hardcoded Defaults**: Used if neither flags nor config define a value.

## Environment Variables
*(Planned for future versions: Support for GO_PUNCH_CONFIG path)*
