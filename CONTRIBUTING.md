# Contributing to Luno MCP Server

Thank you for your interest in contributing to the Luno MCP Server! This document provides guidelines and instructions to help you get started.

## Development Setup

### Prerequisites

- Go 1.24 or later
- [pre-commit](https://pre-commit.com/) for git hooks
- Luno account with API keys (for testing)

### Installing Pre-commit Hooks

We use pre-commit hooks to ensure code quality and consistency. The hooks will run automatically before each commit, checking for common issues and formatting your code.

1. Install pre-commit (if not already installed):
   ```bash
   # macOS
   brew install pre-commit

   # pip
   pip install pre-commit
   ```

2. Install the git hooks:
   ```bash
   pre-commit install
   ```

3. To manually run the hooks on all files:
   ```bash
   pre-commit run --all-files
   ```

## Development Workflow

1. Create a new branch for your feature or bugfix (on a fork if you don't have permission to write to the repo):
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and ensure they follow the project's coding conventions.

3. Run the tests to ensure your changes don't break existing functionality:
   ```bash
   make test
   ```

4. Build and test your changes locally:
   ```bash
   make build
   ./mcp-luno
   ```

5. Commit your changes. The pre-commit hooks will automatically run and may modify some files:
   ```bash
   git add .
   git commit -m "Your descriptive commit message"
   ```

6. Create a pull request.

## Project Structure

- `cmd/` - Application entry points
  - `server/` - Main server application
  - `debug/` - Debugging utilities
- `internal/` - Private application code
  - `config/` - Configuration handling
  - `logging/` - Logging infrastructure
  - `resources/` - MCP resource implementations
  - `server/` - Server implementation
  - `tools/` - MCP tool implementations
  - `tests/` - Integration tests
- `LICENSE` - Project license
- `Makefile` - Build and development commands
- `README.md` - Project documentation

## Testing

We strive for good test coverage. Please add tests for new features and ensure existing tests pass:

```bash
make test
```

## Code Style Guidelines

- Follow standard Go code conventions and idioms
- Use meaningful variable and function names
- Write clear comments for public APIs and complex logic
- Ensure error handling is comprehensive
- Keep functions focused and modular

## License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
