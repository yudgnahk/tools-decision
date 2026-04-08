package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// SmitheryFetcher fetches MCP servers from Smithery.ai
type SmitheryFetcher struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// NewSmitheryFetcher creates a new Smithery fetcher
func NewSmitheryFetcher() *SmitheryFetcher {
	return &SmitheryFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.smithery.ai",
	}
}

// NewSmitheryFetcherWithKey creates a Smithery fetcher with an API key
func NewSmitheryFetcherWithKey(apiKey string) *SmitheryFetcher {
	f := NewSmitheryFetcher()
	f.apiKey = apiKey
	return f
}

// Name returns the fetcher name
func (f *SmitheryFetcher) Name() string {
	return "smithery"
}

// SmitheryServerResponse represents the API response for servers
type SmitheryServerResponse struct {
	Servers    []SmitheryServer   `json:"servers"`
	Pagination SmitheryPagination `json:"pagination"`
}

// SmitheryPagination represents pagination info
type SmitheryPagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalPages  int `json:"totalPages"`
	TotalCount  int `json:"totalCount"`
}

// SmitheryServer represents a server from Smithery API
type SmitheryServer struct {
	ID            string   `json:"id"`
	QualifiedName string   `json:"qualifiedName"`
	Namespace     *string  `json:"namespace"`
	Slug          *string  `json:"slug"`
	DisplayName   string   `json:"displayName"`
	Description   string   `json:"description"`
	IconURL       *string  `json:"iconUrl"`
	Verified      bool     `json:"verified"`
	UseCount      int      `json:"useCount"`
	Remote        *bool    `json:"remote"`
	IsDeployed    bool     `json:"isDeployed"`
	CreatedAt     string   `json:"createdAt"`
	Homepage      string   `json:"homepage"`
	Owner         *string  `json:"owner"`
	Score         *float64 `json:"score"`
}

// Fetch fetches servers from Smithery (with pagination)
func (f *SmitheryFetcher) Fetch(ctx context.Context) ([]types.MCPServer, error) {
	var allServers []types.MCPServer
	page := 1
	pageSize := 100 // Max allowed
	maxPages := 10  // Limit to avoid too many requests

	for page <= maxPages {
		servers, pagination, err := f.fetchPage(ctx, page, pageSize)
		if err != nil {
			if page == 1 {
				return nil, err // Fail if first page fails
			}
			break // Stop on error for subsequent pages
		}

		allServers = append(allServers, servers...)

		if page >= pagination.TotalPages {
			break
		}
		page++
	}

	return allServers, nil
}

// fetchPage fetches a single page of servers
func (f *SmitheryFetcher) fetchPage(ctx context.Context, page, pageSize int) ([]types.MCPServer, *SmitheryPagination, error) {
	u, _ := url.Parse(f.baseURL + "/servers")
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("isDeployed", "true") // Only get deployed servers
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if f.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+f.apiKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch from Smithery: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Smithery returned status %d", resp.StatusCode)
	}

	var response SmitheryServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode Smithery response: %w", err)
	}

	// Convert to our types
	var servers []types.MCPServer
	for _, s := range response.Servers {
		slug := s.QualifiedName
		if s.Slug != nil {
			slug = *s.Slug
		}

		server := types.MCPServer{
			ID:          s.ID,
			Name:        s.DisplayName,
			Slug:        slug,
			Description: s.Description,
			Source:      "smithery",
			Categories:  []string{},
			Tags:        []string{},
			Quality: types.Quality{
				Score:      0.5, // Default score
				Maintained: true,
			},
		}

		// Mark verified servers as higher quality
		if s.Verified {
			server.Quality.Score = 0.8
		}

		// Use count as a quality signal
		if s.UseCount > 1000 {
			server.Quality.Score += 0.1
		} else if s.UseCount > 100 {
			server.Quality.Score += 0.05
		}

		// Remote servers need different install config
		if s.Remote != nil && *s.Remote {
			server.Install = types.Install{
				Command: "npx",
				Args:    []string{"-y", "@smithery/cli", "connect", s.QualifiedName},
			}
		}

		servers = append(servers, server)
	}

	return servers, &response.Pagination, nil
}

// FetchSkills fetches skills from Smithery
func (f *SmitheryFetcher) FetchSkills(ctx context.Context) ([]types.Skill, error) {
	var allSkills []types.Skill
	page := 1
	pageSize := 100
	maxPages := 10

	for page <= maxPages {
		skills, pagination, err := f.fetchSkillsPage(ctx, page, pageSize)
		if err != nil {
			if page == 1 {
				return nil, err
			}
			break
		}

		allSkills = append(allSkills, skills...)

		if page >= pagination.TotalPages {
			break
		}
		page++
	}

	return allSkills, nil
}

// SmitherySkillResponse represents the API response for skills
type SmitherySkillResponse struct {
	Skills     []SmitherySkill    `json:"skills"`
	Pagination SmitheryPagination `json:"pagination"`
}

// SmitherySkill represents a skill from Smithery API
type SmitherySkill struct {
	ID               string   `json:"id"`
	Namespace        string   `json:"namespace"`
	Slug             string   `json:"slug"`
	DisplayName      string   `json:"displayName"`
	Description      string   `json:"description"`
	Prompt           *string  `json:"prompt"`
	QualityScore     float64  `json:"qualityScore"`
	ExternalStars    *int     `json:"externalStars,omitempty"`
	TotalActivations *int     `json:"totalActivations,omitempty"`
	UniqueUsers      *int     `json:"uniqueUsers,omitempty"`
	Categories       []string `json:"categories,omitempty"`
	Servers          []string `json:"servers,omitempty"`
	GitURL           *string  `json:"gitUrl,omitempty"`
	Verified         bool     `json:"verified"`
	Listed           bool     `json:"listed"`
	CreatedAt        string   `json:"createdAt"`
}

func (f *SmitheryFetcher) fetchSkillsPage(ctx context.Context, page, pageSize int) ([]types.Skill, *SmitheryPagination, error) {
	u, _ := url.Parse(f.baseURL + "/skills")
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if f.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+f.apiKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch skills from Smithery: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Smithery skills returned status %d", resp.StatusCode)
	}

	var response SmitherySkillResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode Smithery skills response: %w", err)
	}

	var skills []types.Skill
	for _, s := range response.Skills {
		instructions := ""
		if s.Prompt != nil {
			instructions = *s.Prompt
		}

		skill := types.Skill{
			ID:           s.ID,
			Name:         s.DisplayName,
			Slug:         s.Namespace + "-" + s.Slug,
			Description:  s.Description,
			Instructions: instructions,
			Source:       "smithery",
			Compat: types.SkillCompat{
				Languages: []string{}, // Smithery doesn't provide this yet
			},
			Quality: types.Quality{
				Score: s.QualityScore,
			},
		}

		if s.Servers != nil {
			skill.RequiredTools = s.Servers
		}

		skills = append(skills, skill)
	}

	return skills, &response.Pagination, nil
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

// GlamaServerResponse represents the Glama API response
type GlamaServerResponse struct {
	Servers []GlamaServer `json:"servers"`
}

// GlamaServer represents a server from Glama
type GlamaServer struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Repository  string   `json:"repository"`
	NPM         string   `json:"npm"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
}

// Fetch fetches servers from Glama
func (f *GlamaFetcher) Fetch(ctx context.Context) ([]types.MCPServer, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL+"/servers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Glama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Glama returned status %d", resp.StatusCode)
	}

	var response GlamaServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode Glama response: %w", err)
	}

	var servers []types.MCPServer
	for _, s := range response.Servers {
		servers = append(servers, types.MCPServer{
			ID:          "glama-" + s.Slug,
			Name:        s.Name,
			Slug:        s.Slug,
			Description: s.Description,
			Author:      s.Author,
			Repository:  s.Repository,
			NPM:         s.NPM,
			Categories:  s.Categories,
			Tags:        s.Tags,
			Source:      "glama",
			Quality: types.Quality{
				Score:      0.5,
				Maintained: true,
			},
		})
	}

	return servers, nil
}

// OfficialFetcher fetches from the official MCP registry on GitHub
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
		baseURL: "https://raw.githubusercontent.com/modelcontextprotocol/servers/main",
	}
}

// Name returns the fetcher name
func (f *OfficialFetcher) Name() string {
	return "official"
}

// Fetch fetches servers from the official registry
func (f *OfficialFetcher) Fetch(ctx context.Context) ([]types.MCPServer, error) {
	// The official MCP servers repo doesn't have a registry JSON file
	// We would need to parse the README or use GitHub API to list directories
	// For now, return the embedded servers as "official"
	return nil, fmt.Errorf("official registry parsing not yet implemented - use embedded servers")
}
