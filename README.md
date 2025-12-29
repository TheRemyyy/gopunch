<div align="center">

# GoPunch

**Dead Simple Uptime Monitoring Tool**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/TheRemyyy/gopunch?style=flat-square&color=yellow)](https://github.com/TheRemyyy/gopunch/stargazers)
[![Build Status](https://img.shields.io/github/actions/workflow/status/TheRemyyy/gopunch/build.yml?style=flat-square&label=Build)](https://github.com/TheRemyyy/gopunch/actions)

*A high-performance CLI utility for checking uptime and response times of web services.*

[Quick Start](#quick-start) • [Installation](#installation) • [Commands](#commands) • [Documentation](#documentation)

</div>

---

## Overview

GoPunch is a lightweight, concurrent tool designed for monitoring availability and performance. Whether you need a quick one-off check or continuous monitoring with Discord/Slack alerts, GoPunch delivers with simplicity and speed.

### Key Features

- **🎯 Simple Commands** — `check` for one-off inspection, `watch` for continuous monitoring.
- **📊 Multiple Protocols** — Support for **HTTP**, **TCP**, **DNS**, and **SSL** expiry checks.
- **🚨 Instant Alerting** — Integrated **Discord** & **Slack** webhooks with cooldown management.
- **🔄 Smart Retries** — Automatic retry logic with exponential backoff.
- **⚡ High Concurrency** — Parallel execution using Go routines and semaphores.
- **📝 Exportable** — Output data to **Table**, **JSON**, **CSV**, or **Minimal** formats.

## <a id="installation"></a>📦 Installation

### From Source

```bash
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch
go build -o gopunch ./cmd/gopunch
```

## <a id="quick-start"></a>🚀 Quick Start

### Check a website
```bash
gopunch check https://example.com
```

### Continuous monitoring every 10 seconds
```bash
gopunch watch https://example.com --interval 10
```

### Generate a config template
```bash
gopunch init
```

## <a id="commands"></a>🔧 Commands Summary

- **`check`**: Performs a one-time health check. Supports various flags for methods, headers, and formats.
- **`watch`**: Starts a continuous monitoring loop with live updates and summary stats.
- **`init`**: Generates a sample `gopunch.json` configuration file.
- **`version`**: Displays the current version and build information.

---

## <a id="documentation"></a>📄 Documentation

For deep-dive information, please refer to the specialized guides in the `docs/` directory:

### Core Guides
- 📖 **[Documentation Overview](docs/overview.md)** — The starting point for all documentation.
- 🏗️ **[Architecture & Internals](docs/architecture.md)** — How the engine works and package structure.
- ⚙️ **[Configuration Reference](docs/configuration.md)** — Detailed look at `gopunch.json` and precedence.
- 🚨 **[Alerting System](docs/alerting.md)** — Setting up webhooks, cooldowns, and recovery notifications.
- 📊 **[Output Formats](docs/output-formats.md)** — Detailed examples of Table, JSON, CSV, and Minimal outputs.

### Command Manuals
- 🛠️ **[check command](docs/commands/check.md)** — Complete flag reference and examples for one-time checks.
- 🕒 **[watch command](docs/commands/watch.md)** — Detailed guide on real-time monitoring and statistics.
- 📝 **[init command](docs/commands/init.md)** — How to use and customize the configuration template.

### Protocol Details
- 🌐 **[HTTP/HTTPS](docs/protocols/http.md)** — Headers, body, redirects, and TLS settings.
- 🔌 **[TCP, DNS & SSL](docs/protocols/tcp-dns-ssl.md)** — Monitoring non-HTTP services and certificate expiry.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
<sub>Built with ❤️ and Go</sub>
</div>
