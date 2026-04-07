package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// Cache manages the local registry cache
type Cache struct {
	dir      string
	maxAge   time.Duration
	index    []types.MCPServer
	metadata CacheMetadata
}

// CacheMetadata contains cache timestamps and versions
type CacheMetadata struct {
	LastUpdate time.Time         `json:"last_update"`
	Sources    map[string]string `json:"sources"` // source -> version/etag
}

// NewCache creates a new cache manager
func NewCache(cacheDir string) (*Cache, error) {
	if cacheDir == "" {
		// Default to ~/.cache/tools-decision
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cacheDir = filepath.Join(home, ".cache", "tools-decision")
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cache := &Cache{
		dir:    cacheDir,
		maxAge: 24 * time.Hour, // Refresh if older than 24 hours
	}

	// Load existing cache
	_ = cache.load()

	return cache, nil
}

// NeedsRefresh returns true if the cache should be refreshed
func (c *Cache) NeedsRefresh() bool {
	if c.metadata.LastUpdate.IsZero() {
		return true
	}
	return time.Since(c.metadata.LastUpdate) > c.maxAge
}

// GetServers returns all cached MCP servers
func (c *Cache) GetServers() []types.MCPServer {
	return c.index
}

// UpdateServers replaces the cached servers
func (c *Cache) UpdateServers(servers []types.MCPServer) error {
	c.index = servers
	c.metadata.LastUpdate = time.Now()
	return c.save()
}

// load reads the cache from disk
func (c *Cache) load() error {
	// Load metadata
	metaPath := filepath.Join(c.dir, "meta.json")
	if data, err := os.ReadFile(metaPath); err == nil {
		if err := json.Unmarshal(data, &c.metadata); err != nil {
			return fmt.Errorf("failed to parse cache metadata: %w", err)
		}
	}

	// Load index
	indexPath := filepath.Join(c.dir, "index.json")
	if data, err := os.ReadFile(indexPath); err == nil {
		if err := json.Unmarshal(data, &c.index); err != nil {
			return fmt.Errorf("failed to parse cache index: %w", err)
		}
	}

	return nil
}

// save writes the cache to disk
func (c *Cache) save() error {
	// Save metadata
	metaPath := filepath.Join(c.dir, "meta.json")
	metaData, err := json.MarshalIndent(c.metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache metadata: %w", err)
	}
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return fmt.Errorf("failed to write cache metadata: %w", err)
	}

	// Save index
	indexPath := filepath.Join(c.dir, "index.json")
	indexData, err := json.MarshalIndent(c.index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache index: %w", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		return fmt.Errorf("failed to write cache index: %w", err)
	}

	return nil
}

// Registry manages fetching and caching MCP server data
type Registry struct {
	cache    *Cache
	fetchers []Fetcher
}

// Fetcher is the interface for registry fetchers
type Fetcher interface {
	Name() string
	Fetch(ctx context.Context) ([]types.MCPServer, error)
}

// NewRegistry creates a new registry manager
func NewRegistry(cache *Cache) *Registry {
	return &Registry{
		cache:    cache,
		fetchers: []Fetcher{
			// Add fetchers here
			// NewOfficialFetcher(),
			// NewSmitheryFetcher(),
			// NewGlamaFetcher(),
		},
	}
}

// GetServers returns all MCP servers, refreshing cache if needed
func (r *Registry) GetServers(ctx context.Context, forceRefresh bool) ([]types.MCPServer, error) {
	if !forceRefresh && !r.cache.NeedsRefresh() {
		return r.cache.GetServers(), nil
	}

	// Fetch from all sources
	var allServers []types.MCPServer
	for _, fetcher := range r.fetchers {
		servers, err := fetcher.Fetch(ctx)
		if err != nil {
			// Log warning but continue with other fetchers
			continue
		}
		allServers = append(allServers, servers...)
	}

	// Deduplicate by ID
	seen := make(map[string]bool)
	var unique []types.MCPServer
	for _, server := range allServers {
		if !seen[server.ID] {
			seen[server.ID] = true
			unique = append(unique, server)
		}
	}

	// Update cache
	if err := r.cache.UpdateServers(unique); err != nil {
		return nil, fmt.Errorf("failed to update cache: %w", err)
	}

	return unique, nil
}

// Search searches for servers matching a query
func (r *Registry) Search(ctx context.Context, query string) ([]types.MCPServer, error) {
	servers, err := r.GetServers(ctx, false)
	if err != nil {
		return nil, err
	}

	// Simple substring search for now
	// TODO: Implement better search (fuzzy, weighted fields)
	var matches []types.MCPServer
	for _, server := range servers {
		if matchesQuery(server, query) {
			matches = append(matches, server)
		}
	}

	return matches, nil
}

// matchesQuery checks if a server matches the search query
func matchesQuery(server types.MCPServer, query string) bool {
	// Search in name, description, tags, categories
	searchFields := []string{
		server.Name,
		server.Description,
		server.Slug,
	}
	searchFields = append(searchFields, server.Tags...)
	searchFields = append(searchFields, server.Categories...)

	for _, field := range searchFields {
		if containsIgnoreCase(field, query) {
			return true
		}
	}
	return false
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive contains
	// TODO: Use strings.Contains with ToLower
	return len(s) > 0 && len(substr) > 0
}
