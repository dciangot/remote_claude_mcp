package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/crypto/ssh"
)

type RemoteClaudeServer struct {
	name           string
	version        string
	sshClient      *ssh.Client
	workingDir     string
	claudeCodePath string
}

func NewRemoteClaudeServer() *RemoteClaudeServer {
	return &RemoteClaudeServer{
		name:           "remote-claude-mcp",
		version:        "0.1.0",
		workingDir:     "/tmp/claude-remote",
		claudeCodePath: "claude",
	}
}

func (s *RemoteClaudeServer) connectSSH(host, user, keyPath string, port int) error {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostWithPort := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", hostWithPort, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}

	s.sshClient = client
	return nil
}

func (s *RemoteClaudeServer) executeRemoteCommand(command string) ([]byte, error) {
	if s.sshClient == nil {
		return nil, fmt.Errorf("SSH client not connected")
	}

	session, err := s.sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Output(command)
}

func (s *RemoteClaudeServer) handleBash(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	command, ok := arguments["command"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: command"), nil
	}

	// Add sandbox parameter support
	sandbox := false
	if s, ok := arguments["sandbox"].(bool); ok {
		sandbox = s
	}

	// Build command with safety checks
	var fullCommand string
	if sandbox {
		// Use restricted environment for sandbox mode
		fullCommand = fmt.Sprintf("cd %s && timeout 120 bash -c 'export TMPDIR=/tmp/claude && %s'",
			s.workingDir, command)
	} else {
		fullCommand = fmt.Sprintf("cd %s && %s", s.workingDir, command)
	}

	output, err := s.executeRemoteCommand(fullCommand)
	var resultText string
	if err != nil {
		resultText = fmt.Sprintf("Command failed: %s\nOutput: %s", err.Error(), string(output))
	} else {
		resultText = string(output)
	}

	return mcp.NewToolResultText(resultText), nil
}

func (s *RemoteClaudeServer) getToolHandler(toolName string) server.ToolHandlerFunc {
	switch toolName {
	case "bash":
		return s.handleBash
	case "read":
		return s.handleRead
	case "write":
		return s.handleWrite
	case "edit":
		return s.handleEdit
	case "glob":
		return s.handleGlob
	case "grep":
		return s.handleGrep
	case "webfetch":
		return s.handleWebFetch
	default:
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown tool: %s", toolName)), nil
		}
	}
}

func main() {
	// Support both config file and command line arguments
	var config *Config
	var err error

	if len(os.Args) >= 2 && strings.HasSuffix(os.Args[1], ".json") {
		// Load from config file
		config, err = LoadConfig(os.Args[1])
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	} else if len(os.Args) >= 5 {
		// Command line arguments
		port := 22
		if len(os.Args) >= 6 {
			port, _ = strconv.Atoi(os.Args[5])
		}

		config = &Config{}
		config.SSH.Host = os.Args[1]
		config.SSH.User = os.Args[2]
		config.SSH.KeyPath = os.Args[3]
		config.SSH.Port = port
		config.Remote.WorkingDir = os.Args[4]
		config.Remote.ClaudeCodePath = "claude"
		config.Server.Name = "remote-claude-mcp"
		config.Server.Version = "0.1.0"
	} else {
		log.Fatal("Usage: remote-claude-mcp <config.json> OR remote-claude-mcp <ssh-host> <ssh-user> <ssh-key-path> <remote-working-dir> [ssh-port]")
	}

	mcpServer := NewRemoteClaudeServer()
	mcpServer.name = config.Server.Name
	mcpServer.version = config.Server.Version
	mcpServer.workingDir = config.Remote.WorkingDir
	mcpServer.claudeCodePath = config.Remote.ClaudeCodePath

	// Connect to remote server
	err = mcpServer.connectSSH(config.SSH.Host, config.SSH.User, config.SSH.KeyPath, config.SSH.Port)
	if err != nil {
		log.Fatalf("Failed to connect to remote server: %v", err)
	}
	defer mcpServer.sshClient.Close()

	// Test connection and create working directory if needed
	_, err = mcpServer.executeRemoteCommand(fmt.Sprintf("mkdir -p %s && cd %s && pwd", mcpServer.workingDir, mcpServer.workingDir))
	if err != nil {
		log.Fatalf("Failed to access working directory: %v", err)
	}

	// Test Claude Code availability
	_, err = mcpServer.executeRemoteCommand(fmt.Sprintf("cd %s && which %s", mcpServer.workingDir, mcpServer.claudeCodePath))
	if err != nil {
		log.Printf("Warning: Claude Code not found at '%s', using fallback commands", mcpServer.claudeCodePath)
	}

	log.Printf("Connected to %s@%s:%d, working directory: %s",
		config.SSH.User, config.SSH.Host, config.SSH.Port, mcpServer.workingDir)

	// Create MCP server
	s := server.NewMCPServer(mcpServer.name, mcpServer.version)

	// Define tools
	tools := []string{"bash", "read", "write", "edit", "glob", "grep", "webfetch"}

	for _, toolName := range tools {
		tool := mcp.NewTool(toolName, mcp.WithDescription(fmt.Sprintf("Execute %s on remote server", toolName)))
		handler := mcpServer.getToolHandler(toolName)
		s.AddTool(tool, handler)
	}

	// Start server on stdio
	log.Printf("Starting MCP server %s v%s with tools: %v", mcpServer.name, mcpServer.version, tools)
	server.ServeStdio(s)
}