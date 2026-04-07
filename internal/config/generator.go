package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// Format represents an output configuration format
type Format string

const (
	FormatClaude  Format = "claude"
	FormatCursor  Format = "cursor"
	FormatVSCode  Format = "vscode"
	FormatGeneric Format = "generic"
)

// Generator generates MCP configuration files
type Generator struct{}

// New creates a new config generator
func New() *Generator {
	return &Generator{}
}

// Generate creates a configuration for the specified format
func (g *Generator) Generate(servers []types.MCPServer, format Format) (*types.ConfigOutput, error) {
	switch format {
	case FormatClaude:
		return g.generateClaude(servers)
	case FormatCursor:
		return g.generateCursor(servers)
	case FormatVSCode:
		return g.generateVSCode(servers)
	case FormatGeneric:
		return g.generateGeneric(servers)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// generateClaude generates Claude Desktop configuration
func (g *Generator) generateClaude(servers []types.MCPServer) (*types.ConfigOutput, error) {
	mcpServers := make(map[string]any)

	for _, server := range servers {
		serverConfig := map[string]any{
			"command": server.Install.Command,
			"args":    server.Install.Args,
		}

		// Add environment variables if any
		if len(server.Install.Env) > 0 {
			env := make(map[string]string)
			for _, e := range server.Install.Env {
				if e.Default != "" {
					env[e.Name] = e.Default
				} else {
					env[e.Name] = fmt.Sprintf("${%s}", e.Name)
				}
			}
			serverConfig["env"] = env
		}

		mcpServers[server.Slug] = serverConfig
	}

	content := map[string]any{
		"mcpServers": mcpServers,
	}

	// Determine config path based on OS
	configPath := getClaudeConfigPath()

	return &types.ConfigOutput{
		Format:   string(FormatClaude),
		Filename: "claude_desktop_config.json",
		Path:     configPath,
		Content:  content,
		EnvVars:  collectEnvVars(servers),
	}, nil
}

// generateCursor generates Cursor IDE configuration
func (g *Generator) generateCursor(servers []types.MCPServer) (*types.ConfigOutput, error) {
	mcpServers := make(map[string]any)

	for _, server := range servers {
		serverConfig := map[string]any{
			"command": server.Install.Command,
			"args":    server.Install.Args,
		}

		if len(server.Install.Env) > 0 {
			env := make(map[string]string)
			for _, e := range server.Install.Env {
				if e.Default != "" {
					env[e.Name] = e.Default
				} else {
					env[e.Name] = fmt.Sprintf("${%s}", e.Name)
				}
			}
			serverConfig["env"] = env
		}

		mcpServers[server.Slug] = serverConfig
	}

	content := map[string]any{
		"mcpServers": mcpServers,
	}

	return &types.ConfigOutput{
		Format:   string(FormatCursor),
		Filename: "mcp.json",
		Path:     ".cursor/mcp.json",
		Content:  content,
		EnvVars:  collectEnvVars(servers),
	}, nil
}

// generateVSCode generates VS Code configuration
func (g *Generator) generateVSCode(servers []types.MCPServer) (*types.ConfigOutput, error) {
	mcpServers := make(map[string]any)

	for _, server := range servers {
		mcpServers[server.Slug] = map[string]any{
			"command": server.Install.Command,
			"args":    server.Install.Args,
		}
	}

	content := map[string]any{
		"mcp.servers": mcpServers,
	}

	return &types.ConfigOutput{
		Format:   string(FormatVSCode),
		Filename: "settings.json",
		Path:     ".vscode/settings.json",
		Content:  content,
		EnvVars:  collectEnvVars(servers),
	}, nil
}

// generateGeneric generates a generic MCP configuration
func (g *Generator) generateGeneric(servers []types.MCPServer) (*types.ConfigOutput, error) {
	mcpServers := make(map[string]any)

	for _, server := range servers {
		mcpServers[server.Slug] = map[string]any{
			"command": server.Install.Command,
			"args":    server.Install.Args,
		}
	}

	content := map[string]any{
		"version": "1.0",
		"servers": mcpServers,
	}

	return &types.ConfigOutput{
		Format:   string(FormatGeneric),
		Filename: "mcp.json",
		Path:     "mcp.json",
		Content:  content,
		EnvVars:  collectEnvVars(servers),
	}, nil
}

// Write writes the configuration to disk
func (g *Generator) Write(output *types.ConfigOutput, projectPath string) error {
	fullPath := filepath.Join(projectPath, output.Path)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(output.Content, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getClaudeConfigPath returns the Claude Desktop config path for the current OS
func getClaudeConfigPath() string {
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "Claude", "claude_desktop_config.json")
	default: // Linux
		return filepath.Join(home, ".config", "claude", "claude_desktop_config.json")
	}
}

// collectEnvVars collects all required environment variables from servers
func collectEnvVars(servers []types.MCPServer) []types.EnvVar {
	seen := make(map[string]bool)
	var envVars []types.EnvVar

	for _, server := range servers {
		for _, env := range server.Install.Env {
			if !seen[env.Name] {
				seen[env.Name] = true
				envVars = append(envVars, env)
			}
		}
	}

	return envVars
}
