# VS Code settings.json configuration for Luno MCP

To integrate the Luno MCP server with VS Code Copilot, add the following configuration
to your VS Code settings.json file. You can access this file by pressing `Cmd+Shift+P`
(macOS) or `Ctrl+Shift+P` (Windows/Linux) and typing "Preferences: Open Settings (JSON)".

## For Stdio Mode (default)

```json
"mcp": {
  "servers": {
    "luno": {
      "command": "mcp-luno",
      "args": [],
      "env": {
        "LUNO_API_KEY_ID": "${env:LUNO_API_KEY_ID}",
        "LUNO_API_SECRET": "${env:LUNO_API_SECRET}"
      }
    }
  }
}
```

## For SSE Mode

```json
"mcp": {
  "servers": {
    "luno": {
      "type": "sse",
      "url": "http://localhost:8080/sse"
    }
  }
}
```

## For Docker Mode

```json
"mcp": {
  "servers": {
    "luno": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-e",
        "LUNO_API_KEY_ID",
        "-e",
        "LUNO_API_SECRET",
        "mcp-luno:latest"
      ],
      "env": {
        "LUNO_API_KEY_ID": "${env:LUNO_API_KEY_ID}",
        "LUNO_API_SECRET": "${env:LUNO_API_SECRET}"
      }
    }
  }
}
```

Make sure to set the environment variables `LUNO_API_KEY_ID` and `LUNO_API_SECRET` 
with your Luno API credentials before starting VS Code.
