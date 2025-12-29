# Protocols: TCP, DNS & SSL

Beyond HTTP, GoPunch supports specialized checks for other common network services.

## TCP Port Checks (`tcp://`)

Verifies that a specific port is reachable and accepting connections.

- **Usage**: `gopunch check tcp://db.example.com:5432`
- **Logic**: Performs a raw TCP dial with the specified timeout.
- **Success**: Connection established successfully.

## DNS Resolution Checks (`dns://`)

Ensures that a hostname can be resolved by the system's DNS resolver.

- **Usage**: `gopunch check dns://google.com`
- **Logic**: Performs a lookup using the system's default resolver.
- **Success**: Hostname resolves to one or more IP addresses.
- **Info Output**: Displays the number of IPs found (e.g., `3 IPs`).

## SSL Certificate Expiry Checks (`ssl://`)

Deeply inspects the SSL/TLS certificate of a target to check for validity and upcoming expiration.

- **Usage**: `gopunch check ssl://example.com:443`
- **Logic**: Performs a TLS handshake and parses the peer certificate.
- **Success**: Certificate is currently valid (not expired and already active).
- **Info Output**: Displays the days remaining until expiry (e.g., `Valid (245 days)`).

### Why use `ssl://` instead of `https://`?
While `https://` will fail if a certificate is expired, it won't tell you *when* it expires. The `ssl://` check provides proactive information about how many days you have left before renewal is needed.
