<div align="center">

# ⚡ GoPunch v2.0

**Dead Simple Uptime Monitoring Tool**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/TheRemyyy/gopunch?style=flat-square&color=yellow)](https://github.com/TheRemyyy/gopunch/stargazers)
[![Build](https://img.shields.io/github/actions/workflow/status/TheRemyyy/gopunch/build.yml?style=flat-square&label=Build)](https://github.com/TheRemyyy/gopunch/actions)

A lightweight CLI utility for checking uptime and response times of websites.  
Supports one-off checks, continuous monitoring, custom headers, retries, and alerting.

[Quick Start](#-quick-start) • [Installation](#-installation) • [Commands](#-commands) • [Configuration](#-configuration)

</div>

---

## ✨ Features

- **🎯 Simple Commands** — `check` for one-off, `watch` for continuous monitoring
- **📊 Multiple Output Formats** — Table, JSON, CSV, minimal
- **🔄 Retry Logic** — Automatic retries with exponential backoff
- **🔐 TLS Options** — Skip verification, custom headers, follow redirects
- **🚨 Alerting** — Discord/Slack webhooks with cooldown
- **⚡ Concurrent** — Parallel requests with configurable concurrency
- **🏥 Health Checks** — HTTP, TCP port, DNS, SSL certificate expiry
- **📝 Config Files** — JSON configuration with `gopunch init`

---

## 🚀 Quick Start

```bash
# One-time check
gopunch check https://example.com

# Check multiple URLs with JSON output
gopunch check --format json https://example.com https://api.example.com

# Continuous monitoring every 5 seconds
gopunch watch https://example.com --interval 5

# Generate config file
gopunch init
```

---

## 📦 Installation

### From Releases

Download the latest binary for your platform from [Releases](https://github.com/TheRemyyy/gopunch/releases).

### From Source

```bash
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch
go build -o gopunch ./cmd/gopunch
```

---

## 🔧 Commands

### `check` — One-time health check

```bash
gopunch check [flags] <url> [url...]

Flags:
  -t, --timeout int       Request timeout in seconds (default 10)
  -m, --method string     HTTP method (default "GET")
  -H, --header strings    Custom headers (repeatable)
  -d, --data string       Request body for POST/PUT
  -k, --insecure          Skip TLS verification
  -L, --follow            Follow redirects (default true)
  -e, --expect ints       Expected status codes
  -r, --retries int       Number of retries on failure
  -f, --format string     Output: table, json, csv, minimal (default "table")
  -q, --quiet             Suppress output, exit code only
  -c, --concurrency int   Max concurrent requests (default 10)
```

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

Flags:
  -i, --interval int      Check interval in seconds (default 5)
  # ... same flags as check
```

**Example:**

```bash
gopunch watch https://example.com https://api.example.com --interval 10
```

Press `Ctrl+C` to stop and see summary statistics.

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

---

## ⚙️ Configuration

### JSON Config File

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

---

## 📊 Output Formats

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

---

## 🚨 Alerting

GoPunch supports Discord/Slack webhooks for alerts when endpoints go down.

```json
{
  "alerting": {
    "enabled": true,
    "cooldown_seconds": 300,
    "webhook": {
      "url": "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID"
    }
  }
}
```

- **Cooldown** prevents alert spam (default 5 minutes)
- **Recovery alerts** sent when endpoint comes back online

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## 🔒 Security

For security issues, see [SECURITY.md](SECURITY.md).

---

## 📄 License

**MIT License** © 2025 TheRemyyy
