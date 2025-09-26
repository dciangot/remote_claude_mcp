package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *RemoteClaudeServer) handleRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	filePath, ok := arguments["file_path"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: file_path"), nil
	}

	// Handle offset and limit parameters
	offset := 0
	if o, ok := arguments["offset"].(float64); ok {
		offset = int(o)
	}

	limit := 2000 // default limit
	if l, ok := arguments["limit"].(float64); ok {
		limit = int(l)
	}

	var command string
	if offset > 0 || limit != 2000 {
		command = fmt.Sprintf("cd %s && tail -n +%d '%s' | head -n %d | cat -n",
			s.workingDir, offset+1, filePath, limit)
	} else {
		command = fmt.Sprintf("cd %s && cat -n '%s'", s.workingDir, filePath)
	}

	output, err := s.executeRemoteCommand(command)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func (s *RemoteClaudeServer) handleWrite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	filePath, ok := arguments["file_path"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: file_path"), nil
	}

	content, ok := arguments["content"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: content"), nil
	}

	// Escape content for safe shell execution
	escapedContent := strings.ReplaceAll(content, "'", "'\"'\"'")
	command := fmt.Sprintf("cd %s && echo '%s' > '%s'",
		s.workingDir, escapedContent, filePath)

	_, err := s.executeRemoteCommand(command)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error writing file: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File written successfully: %s", filePath)), nil
}

func (s *RemoteClaudeServer) handleEdit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	filePath, ok := arguments["file_path"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: file_path"), nil
	}

	oldString, ok := arguments["old_string"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: old_string"), nil
	}

	newString, ok := arguments["new_string"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: new_string"), nil
	}

	replaceAll := false
	if r, ok := arguments["replace_all"].(bool); ok {
		replaceAll = r
	}

	// Use sed for text replacement
	sedFlag := ""
	if replaceAll {
		sedFlag = "g"
	}

	// Escape strings for sed
	escapedOld := strings.ReplaceAll(oldString, "/", "\\/")
	escapedNew := strings.ReplaceAll(newString, "/", "\\/")

	command := fmt.Sprintf("cd %s && sed -i 's/%s/%s/%s' '%s'",
		s.workingDir, escapedOld, escapedNew, sedFlag, filePath)

	_, err := s.executeRemoteCommand(command)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error editing file: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File edited successfully: %s", filePath)), nil
}

func (s *RemoteClaudeServer) handleGlob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	pattern, ok := arguments["pattern"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: pattern"), nil
	}

	path := "."
	if p, ok := arguments["path"].(string); ok {
		path = p
	}

	command := fmt.Sprintf("cd %s && find %s -name '%s' -type f | sort",
		s.workingDir, path, pattern)

	output, err := s.executeRemoteCommand(command)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error searching files: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func (s *RemoteClaudeServer) handleGrep(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	pattern, ok := arguments["pattern"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: pattern"), nil
	}

	// Build grep command with parameters
	var grepArgs []string

	if i, ok := arguments["-i"].(bool); ok && i {
		grepArgs = append(grepArgs, "-i")
	}

	if n, ok := arguments["-n"].(bool); ok && n {
		grepArgs = append(grepArgs, "-n")
	}

	path := "."
	if p, ok := arguments["path"].(string); ok {
		path = p
	}

	command := fmt.Sprintf("cd %s && grep -r %s '%s' %s",
		s.workingDir, strings.Join(grepArgs, " "), pattern, path)

	output, err := s.executeRemoteCommand(command)
	var resultText string
	if err != nil {
		resultText = fmt.Sprintf("Search completed with status: %s\nOutput: %s", err.Error(), string(output))
	} else {
		resultText = string(output)
	}

	return mcp.NewToolResultText(resultText), nil
}

func (s *RemoteClaudeServer) handleWebFetch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	url, ok := arguments["url"].(string)
	if !ok {
		return mcp.NewToolResultError("Error: missing required parameter: url"), nil
	}

	// Use curl to fetch web content
	command := fmt.Sprintf("cd %s && curl -s -L '%s'", s.workingDir, url)

	output, err := s.executeRemoteCommand(command)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching URL: %s", err.Error())), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}