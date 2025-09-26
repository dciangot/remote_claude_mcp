# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based Model Context Protocol (MCP) server that allows Claude Desktop to use Claude Code tools on remote servers via SSH. The server acts as a proxy, forwarding Claude Code tool calls (bash, read, write, edit, glob, grep, webfetch) to a remote server through SSH connections.

## Core Architecture

- **RemoteClaudeServer struct**: Main server implementation that manages SSH connections and tool proxying
- **main.go**: Entry point, argument parsing, SSH connection setup, and MCP server initialization
- **handlers.go**: Tool handlers that implement the MCP tool interface (read, write, edit, bash, glob, grep, webfetch)
- **config.go**: JSON configuration file parsing with defaults

The server establishes SSH connections using key-based authentication and executes commands within a specified remote working directory. Each tool handler translates MCP tool requests into appropriate shell commands executed on the remote server.

## Development Commands

### Build and Test
```bash
make build          # Build the binary
make test           # Run tests
make lint           # Run linting and formatting (gofmt, go vet, staticcheck)
make clean          # Remove build artifacts
make deps           # Install/update dependencies (go mod tidy)
```

### Development Setup
```bash
make dev-setup      # Install development tools (gopls, staticcheck)
```

### Cross-platform Builds
```bash
make build-all      # Build for Linux, Darwin, Windows (amd64/arm64)
```

## Configuration

The server supports two configuration methods:
1. JSON config file (recommended) - see `config.example.json`
2. Command line arguments: `<ssh-host> <ssh-user> <ssh-key-path> <remote-working-dir> [ssh-port]`

## Key Dependencies

- `github.com/mark3labs/mcp-go`: MCP protocol implementation
- `golang.org/x/crypto/ssh`: SSH client functionality

## Tool Implementation Pattern

Each MCP tool handler follows this pattern:
1. Extract parameters from `mcp.CallToolRequest`
2. Validate required parameters
3. Construct appropriate shell command for remote execution
4. Execute via `executeRemoteCommand()`
5. Return `mcp.CallToolResult` with output or error

The bash tool includes special handling for sandbox mode and timeout parameters.

## MCP Client Configuration Examples

### Claude Desktop Configuration

Add to your Claude Desktop MCP settings (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "remote-claude": {
      "command": "/path/to/remote-claude-mcp",
      "args": ["/path/to/your/config.json"],
      "env": {}
    }
  }
}
```

### LM Studio Configuration

1. Open LM Studio and go to the **Developer** tab
2. Enable **MCP Server** support
3. Add a new MCP server with these settings:
   - **Name**: `remote-claude`
   - **Command**: `/path/to/remote-claude-mcp`
   - **Arguments**: `["/path/to/your/config.json"]`
   - **Working Directory**: (optional, defaults to current directory)
   - **Environment Variables**: (leave empty unless needed)

4. Your LM Studio MCP configuration should look like:
```json
{
  "mcpServers": {
    "remote-claude": {
      "command": "/path/to/remote-claude-mcp",
      "args": ["/path/to/your/config.json"]
    }
  }
}
```

### Configuration File Setup

Before configuring clients, create your `config.json` based on `config.example.json`:

```json
{
  "ssh": {
    "host": "your-remote-server.com",
    "user": "your-username",
    "key_path": "/Users/yourname/.ssh/id_rsa",
    "port": 22
  },
  "remote": {
    "working_dir": "/home/your-username/claude-workspace",
    "claude_code_path": "claude"
  },
  "server": {
    "name": "remote-claude-mcp",
    "version": "0.1.0"
  }
}
```

### Verification

After configuration, both clients should show the remote-claude server as connected and provide access to these tools:
- `bash` - Execute commands on remote server
- `read` - Read files from remote server
- `write` - Write files to remote server
- `edit` - Edit files using string replacement
- `glob` - Find files matching patterns
- `grep` - Search for text in files
- `webfetch` - Fetch web content via remote server