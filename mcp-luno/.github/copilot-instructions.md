# GitHub Copilot Instructions

This document provides guidelines for using GitHub Copilot with the Luno MCP server. These instructions help Copilot understand our codebase structure, standards, and best practices.

## Repository Structure

The Luno MCP repository follows a structured organization:

- `/cmd` - contains main application entry points
  - `/server` - main MCP server implementation
- `/internal` - internal packages not meant for external use
  - `/config` - configuration handling
  - `/resources` - MCP resources implementation
  - `/server` - server setup and initialization
  - `/tools` - tools implementation for Luno API interaction
  - `/tests` - testing utilities

## Coding Standards

When generating code with Copilot, ensure it follows these principles:

### Go Standards

- Follow Go idioms and best practices
- Use error handling with proper context
- Always either return errors or log them, but never both
- We very rarely should be ignoring errors. If we do, explicitly ignore it with `_ =`
- Implement interfaces only when needed
- Prefer simple, readable code over clever solutions
- Use the Go standard library when appropriate

### MCP Server Guidelines

- Follow the Model Context Protocol (MCP) standards
- Separate concerns between resources and tools
- Implement proper error handling and reporting
- Use the mcp-go library patterns correctly
- Ensure proper validation of inputs from clients

## Testing Requirements

- Write unit tests for all business logic
- Don't just write tests for the happy path, consider boundary/edge and error cases
- Tests should be deterministic whenever possible
- Always write table-driven tests
  - Use spaces rather than `_`s in test names
  - Structure test cases clearly with descriptive fields
- Test functions independently when possible to avoid external dependencies
- Use testify or standard library assertions to check results

## Documentation

- Document APIs, interfaces, and non-obvious implementations
- Reference existing documentation when extending functionality
- Only add comments when extra context is required, i.e. if it's obvious what the code is doing and why, a comment isn't necessary

## Security Considerations

- Never generate code that stores credentials
- Consider potential edge cases and failure modes
- Validate inputs, especially from external sources
- Don't log sensitive information like API keys or full responses
- Only log information needed for debugging

## API Integration Best Practices

- Normalize and validate inputs before sending to the Luno API
- Handle API errors gracefully and provide helpful messages to users
- Cache responses when appropriate to reduce API calls
- Implement proper rate limiting and backoff mechanisms
- Validate trading pairs and other user inputs before submitting orders

## Context-Specific Instructions

When working with:

- Cryptocurrency trading: ensure proper validation of amounts, prices, and currency pairs
- API responses: provide clear error messages and normalize data for consistency
- Authentication: securely handle API keys and secrets
- Market data: format and present data in a user-friendly way
- Tools: follow the MCP tool patterns and provide helpful error messages
