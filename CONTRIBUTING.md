# Contributing to Go DDD Modular Monolith

Thank you for your interest in contributing to this project.

## Development Setup

1. Fork the repository
2. Clone your fork
3. Create a branch for your feature or fix
4. Make your changes
5. Run tests and linter
6. Submit a pull request

## Code Standards

### Go Code Style

Follow standard Go conventions:

- Use `gofmt`, `gofumpt`, and `goimports`
- Run `make fmt` before committing
- Follow effective Go guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types

### Architecture Principles

This project follows specific architectural patterns. When contributing:

1. Respect module boundaries
2. Keep domain logic in the service layer
3. Use ports (interfaces) for dependencies
4. Adapters should be thin
5. No business logic in controllers
6. No business logic in repositories

### Testing

- Write unit tests for new features
- Use mocks for external dependencies
- Write E2E tests for critical flows
- Aim for high test coverage
- Tests must pass before PR approval

Run tests:

```bash
make test-unit      # Unit tests
make test-e2e       # E2E tests
make test-all       # All tests
```

### Code Quality

Before submitting:

```bash
make fmt            # Format code
make lint           # Run linter
make tests          # Run tests
```

All checks must pass.

## Pull Request Process

1. Update tests for your changes
2. Update documentation if needed
3. Ensure all tests pass
4. Ensure linter passes
5. Write a clear PR description
6. Link related issues
7. Wait for review

## Project Structure

When adding new features, follow the existing structure:

### Adding a New Module

```
modules/
└── your-module/
    ├── factory/
    │   └── factory.go          # Dependency injection
    ├── internal/
    │   ├── adapter/
    │   │   ├── controller/     # REST endpoints
    │   │   ├── repository/     # Database access
    │   │   └── cron/           # Background jobs
    │   ├── ports/              # Interfaces
    │   └── service/            # Business logic
    ├── mocks/                  # Generated mocks
    ├── module.go               # Module registration
    └── domain.go               # Domain entities
```

### Adding a New Endpoint

1. Define request/response DTOs in `internal/adapter/controller/dto/`
2. Add handler method in controller
3. Register route in `NewController` function
4. Implement business logic in service
5. Add unit tests
6. Add E2E tests
7. Update documentation

### Adding a New Repository Method

1. Define interface in `internal/ports/`
2. Implement in `internal/adapter/repository/`
3. Add tests
4. Generate mocks: `make go-generate`

## Commit Messages

Use clear, descriptive commit messages:

```
feat: add payment refund endpoint
fix: correct currency validation
docs: update API documentation
test: add tests for payment creation
refactor: simplify repository query
```

Prefixes:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `test:` - Tests
- `refactor:` - Code refactoring
- `chore:` - Maintenance

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
