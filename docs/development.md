# Development Guide

This guide helps you set up the development environment for GoPunch and explains how to contribute.

## Prerequisites

- **Go**: Version 1.22 or newer.
- **Git**: For version control.

## Setup

1.  Clone the repository:
    ```bash
    git clone https://github.com/TheRemyyy/gopunch.git
    cd gopunch
    ```
2.  Install dependencies:
    ```bash
    go mod download
    ```

## Building

To build the executable for your current OS/Architecture:

```bash
go build -o gopunch ./cmd/gopunch
```

To cross-compile for other platforms:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o gopunch.exe ./cmd/gopunch

# Linux
GOOS=linux GOARCH=amd64 go build -o gopunch ./cmd/gopunch
```

## Adding a New Protocol

To add a new health check protocol:

1.  Add the logic to `internal/health/`.
2.  Register the protocol prefix and update the `CheckURLs` switch in `internal/checker/checker.go`.
3.  Update the documentation in `docs/configuration.md`.

## Testing

Run the full test suite:

```bash
go test ./...
```

## Code Quality

Please run `gofmt` before submitting any pull requests to ensure consistent code style.
