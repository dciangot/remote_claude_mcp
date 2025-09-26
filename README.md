# Remote Claude MCP Server

A Model Context Protocol (MCP) server written in Go that allows Claude Desktop to use Claude Code on a remote server via SSH.

## Features

- **SSH Connection**: Securely connect to remote servers using SSH key authentication
- **Claude Code Proxy**: Proxies Claude Code tools (bash, read, write, edit, glob, grep, webfetch) to remote server
- **Sandbox Support**: Supports sandboxed command execution for security
- **Configuration**: Supports both JSON config files and command line arguments
- **Tool Compatibility**: Provides the same tools as local Claude Code but executed remotely

## Installation

1. Ensure you have Go 1.23 or later installed
2. Clone this repository:

   ```bash
   git clone <repository-url>
   cd remote-claude-mcp
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Build the server:

   ```bash
   go build -o remote-claude-mcp
   ```

## Configuration

### Option 1: Config File (Recommended)

Create a configuration file based on `config.example.json`:

```json
{
  "ssh": {
    "host": "your-remote-server.com",
    "user": "your-username",
    "key_path": "/path/to/your/private/ssh/key",
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

### Option 2: Command Line Arguments

```bash
./remote-claude-mcp <ssh-host> <ssh-user> <ssh-key-path> <remote-working-dir> [ssh-port]
```

## Usage

### Running the MCP Server

With config file:

```bash
./remote-claude-mcp config.json
```

With command line arguments:

```bash
./remote-claude-mcp remote-server.com username ~/.ssh/id_rsa /home/username/workspace 22
```

### Integrating with Claude Desktop

Add the following to your Claude Desktop MCP settings:

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

## Available Tools

The server provides the following tools that mirror Claude Code functionality:

- **bash**: Execute bash commands on the remote server
  - Supports `sandbox` mode for secure execution
  - Supports `timeout` parameter
- **read**: Read files from the remote server
  - Supports `offset` and `limit` for large files
- **write**: Write files to the remote server
- **edit**: Edit files using string replacement
  - Supports `replace_all` option
- **glob**: Find files matching patterns
- **grep**: Search for text in files
  - Supports case insensitive search (`-i`)
  - Supports line numbers (`-n`)
- **webfetch**: Fetch web content via the remote server

## Security Considerations

1. **SSH Keys**: Use SSH key authentication instead of passwords
2. **Working Directory**: The server operates within the specified working directory
3. **Sandbox Mode**: Use sandbox mode for bash commands when possible
4. **File Permissions**: Ensure proper file permissions on the remote server
5. **Network Security**: Consider using VPN or SSH tunneling for additional security

## Prerequisites

### Remote Server Setup

1. **SSH Access**: Ensure SSH key-based authentication is configured
2. **Working Directory**: Create and set permissions for the working directory
3. **Claude Code (Optional)**: Install Claude Code on the remote server for enhanced functionality
4. **Basic Tools**: Ensure `bash`, `cat`, `grep`, `find`, `curl` are available

### Claude Desktop Setup

1. Install Claude Desktop
2. Configure MCP servers in Claude Desktop settings
3. Add this server to your MCP configuration

## Troubleshooting

### Connection Issues

- Verify SSH key permissions (`chmod 600 ~/.ssh/id_rsa`)
- Test SSH connection manually: `ssh -i ~/.ssh/id_rsa user@host`
- Check firewall and network connectivity

### Tool Execution Issues

- Verify working directory exists and is writable
- Check that required tools (`bash`, `cat`, etc.) are in PATH on remote server
- Review server logs for error messages

### Claude Desktop Integration

- Verify MCP configuration syntax in Claude Desktop settings
- Check Claude Desktop logs for connection errors
- Ensure the server binary path is correct and executable

## Development

### Building from Source

```bash
go mod tidy
go build -o remote-claude-mcp
```

### Running Tests

```bash
go test ./...
```

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues and questions:

1. Check the troubleshooting section
2. Review existing GitHub issues
3. Create a new issue with detailed information

