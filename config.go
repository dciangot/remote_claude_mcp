package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	SSH struct {
		Host    string `json:"host"`
		User    string `json:"user"`
		KeyPath string `json:"key_path"`
		Port    int    `json:"port"`
	} `json:"ssh"`
	Remote struct {
		WorkingDir     string `json:"working_dir"`
		ClaudeCodePath string `json:"claude_code_path"`
	} `json:"remote"`
	Server struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"server"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Set defaults
	if config.SSH.Port == 0 {
		config.SSH.Port = 22
	}
	if config.Remote.ClaudeCodePath == "" {
		config.Remote.ClaudeCodePath = "claude"
	}
	if config.Server.Name == "" {
		config.Server.Name = "remote-claude-mcp"
	}
	if config.Server.Version == "" {
		config.Server.Version = "0.1.0"
	}

	return &config, nil
}