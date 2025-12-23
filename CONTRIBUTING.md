<div align="center">

# Contributing to GoPunch

**Help us build the best uptime monitor**

</div>

---

Thank you for your interest in contributing to GoPunch! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Bugs

1. Check existing issues to avoid duplicates.
2. Use the bug report template.
3. Include Go version, OS, and steps to reproduce.

### Suggesting Features

1. Open an issue with the "enhancement" label.
2. Describe the use case and proposed solution.

### Pull Requests

1. **Fork** the repository.
2. **Create a branch**: `git checkout -b feature/your-feature`.
3. **Write code**: Keep it clean, documented, and tested.
4. **Test**: Run `go test ./...` before submitting.
5. **Submit PR**: Provide a clear description of your changes.

## Code Style

- Follow standard Go conventions.
- Use `gofmt` to format your code.
- Keep functions focused and small.
- Add comments for exported functions.

## Development Setup

```bash
git clone https://github.com/TheRemyyy/gopunch.git
cd gopunch
go mod download
go build -o gopunch ./cmd/gopunch
```

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.

---

<div align="center">
<sub>Built with ❤️ and Go</sub>
</div>
