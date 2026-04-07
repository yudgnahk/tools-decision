package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// SmitheryFetcher fetches MCP servers from Smithery.ai
type SmitheryFetcher struct {
	client  *http.Client
	baseURL string
}

// NewSmitheryFetcher creates a new Smithery fetcher
func NewSmitheryFetcher() *SmitheryFetcher {
	return &SmitheryFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://smithery.ai/api",
	}
}

// Name returns the fetcher name
func (f *SmitheryFetcher) Name() string {
	return "smithery"
}

// SmitheryServer represents a server from Smithery API
type SmitheryServer struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Repository  string   `json:"repository"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Stars       int      `json:"stars"`
	NPM         string   `json:"npm_package"`
}

// Fetch fetches servers from Smithery
func (f *SmitheryFetcher) Fetch(ctx context.Context) ([]types.MCPServer, error) {
	// TODO: Implement actual API call
	// This is a placeholder - need to check Smithery's actual API structure

	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL+"/servers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Smithery: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Smithery returned status %d", resp.StatusCode)
	}

	var smitheryServers []SmitheryServer
	if err := json.NewDecoder(resp.Body).Decode(&smitheryServers); err != nil {
		return nil, fmt.Errorf("failed to decode Smithery response: %w", err)
	}

	// Convert to our types
	var servers []types.MCPServer
	for _, s := range smitheryServers {
		servers = append(servers, types.MCPServer{
			ID:          fmt.Sprintf("smithery-%s", s.Slug),
			Name:        s.Name,
			Slug:        s.Slug,
			Description: s.Description,
			Author:      s.Author,
			Repository:  s.Repository,
			NPM:         s.NPM,
			Categories:  s.Categories,
			Tags:        s.Tags,
			Quality: types.Quality{
				Stars:      s.Stars,
				Maintained: true,
			},
			Source: "smithery",
		})
	}

	return servers, nil
}

// GlamaFetcher fetches MCP servers from Glama.ai
type GlamaFetcher struct {
	client  *http.Client
	baseURL string
}

// NewGlamaFetcher creates a new Glama fetcher
func NewGlamaFetcher() *GlamaFetcher {
	return &GlamaFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://glama.ai/api/mcp",
	}
}

// Name returns the fetcher name
func (f *GlamaFetcher) Name() string {
	return "glama"
}

// Fetch fetches servers from Glama
func (f *GlamaFetcher) Fetch(ctx context.Context) ([]types.MCPServer, error) {
	// TODO: Implement actual API call
	// This is a placeholder - need to check Glama's actual API structure
	return nil, fmt.Errorf("Glama fetcher not yet implemented")
}

// OfficialFetcher fetches from the official MCP registry
type OfficialFetcher struct {
	client  *http.Client
	baseURL string
}

// NewOfficialFetcher creates a new official registry fetcher
func NewOfficialFetcher() *OfficialFetcher {
	return &OfficialFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		// The official registry is on GitHub
		baseURL: "https://raw.githubusercontent.com/modelcontextprotocol/servers/main",
	}
}

// Name returns the fetcher name
func (f *OfficialFetcher) Name() string {
	return "official"
}

// Fetch fetches servers from the official registry
func (f *OfficialFetcher) Fetch(ctx context.Context) ([]types.MCPServer, error) {
	// TODO: Implement actual fetch from GitHub
	// The official registry structure needs to be checked
	return nil, fmt.Errorf("Official fetcher not yet implemented")
}
