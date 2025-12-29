# Command: watch

The `watch` command is used for continuous monitoring. It executes checks at regular intervals and provides a live view of the status of your services.

## Usage

```bash
gopunch watch [flags] [url...]
```

## Features

- **Live TUI Updates**: See results as they happen with timestamps.
- **Summary Statistics**: When stopped (Ctrl+C), it displays a comprehensive table with uptime percentage and latency stats.
- **Alert Integration**: Automatically sends alerts via configured webhooks on failure and recovery.

## Flags

| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--interval` | `-i` | `5` | Time between check cycles in seconds. |
| `--quiet` | `-q` | `false` | Only print errors to the console. |
| *...all `check` flags* | | | Inherits all flags from the `check` command. |

## Interactive Output

When running, `watch` provides a color-coded log:
- `✓` (Green): Check succeeded.
- `✗` (Red): Check failed (displays error message).

```text
⚡ Watching 2 URL(s) every 5s (Ctrl+C to stop)

[15:30:00] ✓ https://google.com - 200 145ms
[15:30:00] ✓ https://github.com - 200 89ms
[15:30:05] ✓ https://google.com - 200 152ms
[15:30:05] ✗ https://github.com - request failed: timeout
```

## Watch Summary

Upon exiting, GoPunch calculates:
- **Uptime %**: Ratio of successful checks to total checks.
- **Average Latency**: Mean response time across all checks.
- **Min/Max Latency**: Peak performance and worst-case response times.

## Alerting in Watch Mode

If `alerting` is enabled in your configuration, `watch` will:
1.  Send a **Failure Alert** the moment a service goes down.
2.  Maintain **Cooldown** to avoid spamming your notification channel.
3.  Send a **Recovery Alert** when the service is back online (status 200-399).
