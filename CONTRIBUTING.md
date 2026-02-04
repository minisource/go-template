# Contributing to Go Microservice Template

Thank you for your interest in contributing! This document provides guidelines for contributing to this template.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/template_golang.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit with a clear message
7. Push and create a Pull Request

## Development Setup

```bash
# Install dependencies
cd src && go mod download

# Install development tools
make install-tools

# Run the application
make run

# Run tests
make test

# Run linter
make lint
```

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `make fmt` before committing
- Run `make lint` to check for issues
- Use meaningful variable and function names
- Add comments for exported functions

## Project Structure

When adding new features, follow the existing architecture:

```
domain/     → Entities and repository interfaces (no dependencies)
usecase/    → Business logic (depends on domain)
infra/      → External implementations (depends on domain)
api/        → HTTP layer (depends on usecase)
```

## Testing

- Add unit tests for new business logic in `tests/unit/`
- Add integration tests for API endpoints in `tests/integration/`
- Maintain test coverage above 70%

## Commit Messages

Follow conventional commits:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks

Example: `feat: add product entity with CRUD operations`

## Pull Request Process

1. Update README if needed
2. Update documentation for new features
3. Ensure all tests pass
4. Request review from maintainers

## Questions?

Open an issue for questions or suggestions.
