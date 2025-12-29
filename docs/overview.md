# GoPunch Documentation

Welcome to the comprehensive documentation for **GoPunch** — a high-performance, dead-simple uptime monitoring tool written in Go.

## Introduction

GoPunch is designed to be the only tool you need for quick health checks and continuous monitoring of your network services. Whether you are debugging a REST API, monitoring a database port, or keeping an eye on SSL certificate expiration, GoPunch has you covered.

## Documentation Structure

### 🛠️ Core Commands
- **[check](commands/check.md)**: One-time health checks with rich output formats.
- **[watch](commands/watch.md)**: Real-time monitoring with uptime statistics.
- **[init](commands/init.md)**: Quick start with configuration templates.

### 🌐 Supported Protocols
- **[HTTP/HTTPS](protocols/http.md)**: Custom methods, headers, body, and redirect logic.
- **[TCP, DNS & SSL](protocols/tcp-dns-ssl.md)**: Beyond simple HTTP checks.

### ⚙️ Configuration & Alerts
- **[Configuration Reference](configuration.md)**: Detailed `gopunch.json` documentation.
- **[Alerting System](alerting.md)**: Setting up Discord/Slack notifications and recovery alerts.

### 🏗️ Technical Details
- **[Architecture](architecture.md)**: Under the hood of the concurrency engine.
- **[Development Guide](development.md)**: Building, testing, and contributing.

## Why GoPunch?

1.  **Fast**: Built with Go's native concurrency, it can check hundreds of targets in milliseconds.
2.  **Versatile**: Automatically detects what kind of check to perform based on the URL scheme.
3.  **Beautiful**: Provides clean, color-coded CLI tables and a real-time monitoring dashboard.
4.  **Integrated**: Export to JSON/CSV for pipeline integration, or use Webhooks for instant team alerts.