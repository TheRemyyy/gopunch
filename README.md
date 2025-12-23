<div align="center">

# GoPunch

**Dead Simple Uptime Monitoring Tool**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/TheRemyyy/gopunch?style=flat-square&color=yellow)](https://github.com/TheRemyyy/gopunch/stargazers)
[![Build Status](https://img.shields.io/github/actions/workflow/status/TheRemyyy/gopunch/build.yml?style=flat-square&label=Build)](https://github.com/TheRemyyy/gopunch/actions)

*A lightweight CLI utility for checking uptime and response times of websites.*
*Supports one-off checks, continuous monitoring, custom headers, retries, and alerting.*

[Quick Start](#-quick-start) • [Installation](#-installation) • [Commands](#-commands) • [Configuration](#-configuration)

</div>

---

## Overview

GoPunch is a high-performance CLI tool designed for checking the availability and response times of web services. Whether you need a quick one-off check or continuous monitoring with alerts, GoPunch delivers with simplicity and speed.

### Key Features

- **🎯 Simple Commands** — `check` for one-off inspection, `watch` for continuous monitoring.
- **📊 Multiple Outputs** — Export data to **Table**, **JSON**, **CSV**, or keep it **Minimal**.
- **🔄 Smart Retries** — Automatic retry logic with exponential backoff for transient failures.
- **🔐 Advanced TLS** — Skip verification, customizable headers, and redirect handling.
- **🚨 Instant Alerting** — Integrated **Discord** & **Slack** webhooks with configurable cooldowns.
- **⚡ High Concurrency** — Parallel execution with adjustable concurrency levels.
- **🏥 Multi-Protocol** — Support for **HTTP**, **TCP**, **DNS**, and **SSL** expiry checks.
- **📝 Configuration** — Easy setup via JSON config generated with `gopunch init`.

## Installation

### From Releases

Download the latest binary for your platform from [Releases](https://github.com/TheRemyyy/gopunch/releases).

### From Source

```bash
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch
go build -o gopunch ./cmd/gopunch
```

## Quick Start

### One-time Check

```bash
gopunch check https://example.com
```

### JSON Output

```bash
gopunch check --format json https://example.com https://api.example.com
```

### Continuous Monitoring

```bash
# Monitor every 5 seconds
gopunch watch https://example.com --interval 5
```

### Generate Config

```bash
gopunch init
```

## Commands

### `check` — One-time health check

```bash
gopunch check [flags] <url> [url...]
```

**Flags:**

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-t, --timeout` | Request timeout in seconds | `10` |
| `-m, --method` | HTTP method | `"GET"` |
| `-H, --header` | Custom headers (repeatable) | - |
| `-d, --data` | Request body for POST/PUT | - |
| `-k, --insecure` | Skip TLS verification | `false` |
| `-L, --follow` | Follow redirects | `true` |
| `-e, --expect` | Expected status codes | - |
| `-r, --retries` | Number of retries on failure | `0` |
| `-f, --format` | Output: table, json, csv, minimal | `"table"` |
| `-q, --quiet` | Suppress output, exit code only | `false` |
| `-c, --concurrency` | Max concurrent requests | `10` |

**Examples:**

```bash
# POST with custom header
gopunch check https://api.example.com \
  --method POST \
  --header "Authorization: Bearer token" \
  --data '{"test": true}'

# Check expecting specific status codes
gopunch check https://example.com --expect 200 --expect 201

# JSON output for scripting
gopunch check --format json https://example.com | jq .
```

---

### `watch` — Continuous monitoring

```bash
gopunch watch [flags] <url> [url...]
```

**Flags:**

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-i, --interval` | Check interval in seconds | `5` |
| *...plus all `check` flags* | | |

**Example:**

```bash
gopunch watch https://example.com https://api.example.com --interval 10
```

*Press `Ctrl+C` to stop and see summary statistics.*

---

### `init` — Generate config file

```bash
gopunch init [filename]
```

Creates a `gopunch.json` configuration file with all options documented.

---

### `version` — Show version info

```bash
gopunch version
```

## Configuration

GoPunch uses a simple JSON configuration format:

```json
{
  "urls": ["https://example.com", "https://api.example.com"],
  "interval": 30,
  "timeout": 10,
  "method": "GET",
  "headers": {
    "Authorization": "Bearer token"
  },
  "insecure": false,
  "follow_redirects": true,
  "concurrency": 10,
  "retries": 2,
  "alerting": {
    "enabled": true,
    "cooldown_seconds": 300,
    "webhook": {
      "url": "https://discord.com/api/webhooks/YOUR_ID",
      "method": "POST"
    }
  }
}
```

## Output Formats

### Table (default)

```
STATUS  URL                      CODE  TIME   SIZE
✓       https://example.com      200   150ms  1.2KB
✗       https://api.example.com  500   89ms   0B
```

### JSON

```json
[
  {"url":"https://example.com","status_code":200,"duration_ms":150,"size":1234,"success":true,"error":null}
]
```

### CSV

```csv
url,status_code,duration_ms,size,success,error
https://example.com,200,150,1234,true,
```

## Alerting

GoPunch supports Discord/Slack webhooks for alerts when endpoints go down.

- **Cooldown**: Prevents alert spam (default 5 minutes).
- **Recovery**: Alerts are sent when an endpoint comes back online.

## Security

For security issues, see [SECURITY.md](SECURITY.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
<sub>Built with ❤️ and Go</sub>
</div>
