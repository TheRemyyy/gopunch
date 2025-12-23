<div align="center">

# ⚡ GoPunch

**Dead Simple Uptime Monitoring Tool**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/TheRemyyy/gopunch?style=flat-square&color=yellow)](https://github.com/TheRemyyy/gopunch/stargazers)
[![Build](https://img.shields.io/github/actions/workflow/status/TheRemyyy/gopunch/build.yml?style=flat-square&label=Build)](https://github.com/TheRemyyy/gopunch/actions)

A lightweight CLI utility for checking uptime and response times of websites.  
Supports one-off checks or recurring interval-based monitoring with detailed statistics.

[Installation](#-installation) • [Usage](#-usage) • [Configuration](#-configuration) • [Contributing](CONTRIBUTING.md)

</div>

---

## ✨ Features

- **Multi-URL Support** — Check multiple endpoints with customizable HTTP methods
- **Interval Monitoring** — Continuous monitoring with configurable intervals
- **Detailed Statistics** — Response times, error counts, status code distribution, histograms
- **Export Options** — Export to CSV or JSON format
- **Concurrent Requests** — Thread-safe parallel requests with configurable concurrency
- **JSON Config** — Streamlined setup via configuration file
- **Zero Dependencies** — Lightweight and easy to deploy

---

## 🚀 Quick Start

```bash
# Clone and run
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch

# Single check
go run ./cmd/gopunch https://postman-echo.com/get

# Build binary
go build -o gopunch ./cmd/gopunch
./gopunch https://postman-echo.com/get
```

---

## 📦 Installation

### From Source

```bash
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch
go build -o gopunch ./cmd/gopunch
```

### Dependencies

```bash
go get golang.org/x/sync
```

---

## 🔧 Usage

```bash
gopunch [flags] <url1> [url2 ...]
```

### Examples

| Purpose | Command |
|---------|---------|
| Single check | `gopunch https://example.com` |
| With config file | `gopunch --config=config.json` |
| Interval monitoring | `gopunch --interval=5 https://example.com` |
| POST request | `gopunch --method=POST --data='{"key":"value"}' https://api.example.com` |
| Export to CSV | `gopunch --export=stats.csv https://example.com` |
| Verbose output | `gopunch --verbose https://example.com` |

---

## ⚙️ Configuration

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to JSON config file | `""` |
| `--interval` | Interval between checks in seconds (0 = once) | `0` |
| `--timeout` | Request timeout in seconds | `5` |
| `--method` | HTTP method (GET, POST, HEAD, etc.) | `GET` |
| `--data` | Request body for POST/PUT | `""` |
| `--concurrency` | Maximum concurrent requests | `10` |
| `--verbose` | Enable verbose output | `false` |
| `--logfile` | Write output to file | `""` |
| `--export` | Export stats (`.csv` or `.json`) | `""` |

### JSON Config File

```json
{
  "urls": ["https://example.com", "https://api.example.com"],
  "interval": 5,
  "timeout": 15,
  "concurrency": 2,
  "verbose": true,
  "logfile": "results.log",
  "export": "stats.csv",
  "url_configs": [
    {
      "url": "https://api.example.com",
      "method": "POST",
      "data": "{\"test\":\"value\"}"
    }
  ]
}
```

---

## 📊 Output

### Console

```
✅ https://example.com - Status: 200 OK, Time: 150 ms, Size: 1256 bytes

--- Statistics ---
Total checks: 1
Total errors: 0
https://example.com - Checks: 1, Avg response time: 150.00 ms
Status codes:
  200 OK: 1 (100.00%)
Response time histogram:
  <100ms: 0 | 100-500ms: 1 | 500-1000ms: 0 | >1000ms: 0
```

### CSV Export

```csv
URL,Checks,AvgTime_ms,AvgSize_bytes,StatusCode,Count
https://example.com,1,150.00,1256.00,200 OK,1
```

### JSON Export

```json
{
  "total_checks": 1,
  "total_errors": 0,
  "urls": {
    "https://example.com": {
      "checks": 1,
      "avg_time_ms": 150,
      "status_codes": { "200 OK": 1 }
    }
  }
}
```

---

## 🛑 Graceful Shutdown

For interval-based monitoring, press `Ctrl+C` to stop. GoPunch will print final statistics and export to the specified file.

---

## 📝 Notes

- Use `https://postman-echo.com/get` for reliable testing
- Increase `--timeout` if encountering deadline errors
- Flags override config file values

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a PR.

---

## 🔒 Security

For security issues, please see our [Security Policy](SECURITY.md).

---

## 📄 License

**MIT License** © 2025 TheRemyyy
