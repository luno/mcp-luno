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

### Setting up API Credentials for Development

Set these either through:

#### A shell file

Either set this through your shell file or terminal with:
Set the following environment variables:

```bash
export LUNO_API_KEY_ID=your_api_key_id
export LUNO_API_SECRET=your_api_secret
# Optional: Enable debug mode
export LUNO_API_DEBUG=true
```

#### An .env file

Copy the .env.example file and name it .env (this always should be .gitignored), and paste your keys in there.

Depending on your setup, you might need an additional step to load these vars for your application. E.g. [godotenv](https://github.com/joho/godotenv)

### Running the server

#### Standard I/O mode (default)

```bash
luno-mcp
```

#### Server-Sent Events (SSE) mode

```bash
luno-mcp --transport sse --sse-address localhost:8080
```

#### Using Docker

Build the Docker image:

```bash
docker build -t luno-mcp .
```

Run the Docker container:

```bash
docker run -e LUNO_API_KEY_ID=$LUNO_API_KEY_ID -e LUNO_API_SECRET=$LUNO_API_SECRET luno-mcp
```

Alternatively, you can use an `.env` file to provide these environment variables. This simplifies the command and prevents your API key and secret from being stored in your shell history. First, ensure you have an `.env` file (you can copy `.env.example` and fill in your details) with your credentials, for example:

```env
LUNO_API_KEY_ID=your_api_key_id
LUNO_API_SECRET=your_api_secret
```

Then, run the container using:

```bash
docker run --env-file .env luno-mcp
```

You can also use the `--transport sse` and `--sse-address` flags with Docker when using an `.env` file:

```bash
docker run --env-file .env -p 8080:8080 luno-mcp --transport sse --sse-address 0.0.0.0:8080
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
   ./luno-mcp
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
