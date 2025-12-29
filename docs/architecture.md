# Architecture & Internal Packages

GoPunch is designed as a modular toolset. Each internal package has a specific responsibility, allowing for easy expansion.

## Internal Package Map

### `internal/checker`
The orchestration layer.
- **Purpose**: Manages the high-level flow of health checks.
- **Key Logic**: Uses `golang.org/x/sync/semaphore` to strictly control the number of goroutines running at once.
- **Dispatch**: Maps URL schemes to the appropriate logic in `internal/health`.

### `internal/health`
The protocol implementation layer.
- **Purpose**: Pure logic for checking different services.
- **Checks**:
    - `CheckHTTP`: Uses `net/http` with a custom `Transport`.
    - `CheckTCP`: Uses `net.DialTimeout`.
    - `CheckDNS`: Uses `net.Resolver`.
    - `CheckSSL`: Uses `crypto/tls` and inspects `ConnectionState`.

### `internal/alerter`
The notification layer.
- **Purpose**: Handles stateful alerting.
- **Key Logic**: Implements a thread-safe map of `lastAlert` timestamps per URL to manage the cooldown period.

## Data Flow

1.  **CLI Init**: `cobra` parses flags and merges them with `gopunch.json`.
2.  **Options Prep**: Configuration is converted into a `checker.Options` struct.
3.  **Concurrency Pool**: `checker.CheckURLs` starts a goroutine for each target.
4.  **Semaphore Wait**: Each goroutine waits for a slot in the concurrency pool.
5.  **Execution**: The appropriate `health` function is called.
6.  **Reporting**: Results are collected and passed to the output formatter or the `alerter`.

## Performance Considerations

- **Memory**: GoPunch is extremely lean. It avoids large buffers and cleans up HTTP response bodies immediately.
- **Network**: By disabling Keep-Alives in monitoring mode, it ensures that DNS resolution and TLS handshakes are re-tested during every cycle.