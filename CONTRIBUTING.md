# Contributing to GoPunch

Thank you for your interest in contributing to GoPunch! This document provides guidelines for contributing.

## How to Contribute

### Reporting Bugs

1. Check existing issues to avoid duplicates
2. Use the bug report template
3. Include Go version, OS, and steps to reproduce

### Suggesting Features

1. Open an issue with the "enhancement" label
2. Describe the use case and proposed solution

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Write clean, documented code
4. Add tests if applicable
5. Run `go test ./...` before submitting
6. Submit a PR with a clear description

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Add comments for exported functions

## Development Setup

```bash
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch
go mod download
go build -o gopunch ./cmd/gopunch
```

## Testing

```bash
go test ./...
```

## License

By contributing, you agree that your contributions will be licensed under the project's license.
