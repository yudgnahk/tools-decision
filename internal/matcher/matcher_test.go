package matcher

import (
	"testing"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

func TestMatcher_Match(t *testing.T) {
	m := New()

	// Create a minimal server registry for testing
	servers := []types.MCPServer{
		{
			ID:          "postgres",
			Name:        "PostgreSQL",
			Slug:        "postgres",
			Description: "PostgreSQL database operations",
			Categories:  []string{"database", "postgresql"},
			Tags:        []string{"postgres", "sql"},
			Quality:     types.Quality{Score: 0.9, Maintained: true},
			Source:      "official",
		},
		{
			ID:          "mysql",
			Name:        "MySQL",
			Slug:        "mysql",
			Description: "MySQL database operations",
			Categories:  []string{"database", "mysql"},
			Tags:        []string{"mysql", "sql"},
			Quality:     types.Quality{Score: 0.8, Maintained: true},
			Source:      "community",
		},
		{
			ID:          "redis",
			Name:        "Redis",
			Slug:        "redis",
			Description: "Redis cache operations",
			Categories:  []string{"database", "redis", "cache"},
			Tags:        []string{"redis", "cache"},
			Quality:     types.Quality{Score: 0.75, Maintained: true},
			Source:      "community",
		},
		{
			ID:          "filesystem",
			Name:        "Filesystem",
			Slug:        "filesystem",
			Description: "File system operations",
			Categories:  []string{"core", "filesystem"},
			Tags:        []string{"files", "directories"},
			Quality:     types.Quality{Score: 0.95, Maintained: true},
			Source:      "official",
		},
	}

	tests := []struct {
		name           string
		context        *types.ProjectContext
		wantTopServer  string
		wantMinResults int
	}{
		{
			name: "PostgreSQL service should rank postgres first",
			context: &types.ProjectContext{
				Services: []types.Service{{Name: "postgresql", Confidence: 0.9}},
			},
			wantTopServer:  "postgres",
			wantMinResults: 1,
		},
		{
			name: "Redis service should rank redis highly",
			context: &types.ProjectContext{
				Services: []types.Service{{Name: "redis", Confidence: 0.9}},
			},
			wantTopServer:  "redis",
			wantMinResults: 1,
		},
		{
			name:           "Empty context returns only core baseline results",
			context:        &types.ProjectContext{},
			wantMinResults: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := m.Match(tt.context, servers, 10)

			if len(results) < tt.wantMinResults {
				t.Errorf("expected at least %d results, got %d", tt.wantMinResults, len(results))
			}

			if tt.wantTopServer != "" && len(results) > 0 {
				if results[0].Server.Slug != tt.wantTopServer {
					t.Errorf("expected top server %s, got %s", tt.wantTopServer, results[0].Server.Slug)
				}
			}
		})
	}
}

func TestMatcher_Guardrails(t *testing.T) {
	m := New()

	servers := []types.MCPServer{
		{
			ID:          "stripe",
			Name:        "Stripe",
			Slug:        "stripe",
			Description: "Stripe payment integration",
			Categories:  []string{"payments", "stripe"},
			Tags:        []string{"stripe", "billing"},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
		{
			ID:          "prisma",
			Name:        "Prisma",
			Slug:        "prisma",
			Description: "Prisma ORM",
			Categories:  []string{"database", "orm", "prisma"},
			Tags:        []string{"prisma", "orm"},
			Compat: types.Compat{
				Languages:  []string{"javascript", "typescript"},
				Frameworks: []string{"nextjs", "nestjs"},
			},
			Quality: types.Quality{Score: 0.8, Maintained: true},
		},
	}

	t.Run("stripe excluded without explicit payment signals", func(t *testing.T) {
		ctx := &types.ProjectContext{
			Languages: []types.Language{{Name: "javascript", Confidence: 0.9}},
		}
		results := m.Match(ctx, servers, 10)
		for _, r := range results {
			if r.Server.Slug == "stripe" {
				t.Fatalf("expected stripe to be excluded without explicit intent")
			}
		}
	})

	t.Run("prisma excluded on language-only match", func(t *testing.T) {
		ctx := &types.ProjectContext{
			Languages: []types.Language{{Name: "javascript", Confidence: 0.95}},
		}
		results := m.Match(ctx, servers, 10)
		for _, r := range results {
			if r.Server.Slug == "prisma" {
				t.Fatalf("expected prisma to be excluded without framework/service signals")
			}
		}
	})

	t.Run("prisma included with explicit framework signal", func(t *testing.T) {
		ctx := &types.ProjectContext{
			Languages:  []types.Language{{Name: "javascript", Confidence: 0.95}},
			Frameworks: []types.Framework{{Name: "nextjs", Confidence: 0.9}},
		}
		results := m.Match(ctx, servers, 10)
		found := false
		for _, r := range results {
			if r.Server.Slug == "prisma" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected prisma to be included with explicit framework signal")
		}
	})
}

func TestGetSuggestionsForContext(t *testing.T) {
	tests := []struct {
		name              string
		context           *types.ProjectContext
		wantMinSuggCount  int
		wantServiceExists string
	}{
		{
			name: "PostgreSQL service",
			context: &types.ProjectContext{
				Services: []types.Service{{Name: "postgresql"}},
			},
			wantMinSuggCount:  2, // postgresql + core
			wantServiceExists: "postgresql",
		},
		{
			name: "Multiple services",
			context: &types.ProjectContext{
				Services: []types.Service{
					{Name: "postgresql"},
					{Name: "redis"},
				},
			},
			wantMinSuggCount:  3, // postgresql + redis + core
			wantServiceExists: "redis",
		},
		{
			name:             "Empty context returns core suggestions",
			context:          &types.ProjectContext{},
			wantMinSuggCount: 1, // at least core
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := GetSuggestionsForContext(tt.context)
			if len(suggestions) < tt.wantMinSuggCount {
				t.Errorf("expected at least %d suggestions, got %d", tt.wantMinSuggCount, len(suggestions))
			}

			if tt.wantServiceExists != "" {
				found := false
				for _, s := range suggestions {
					if s.Service == tt.wantServiceExists {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected to find service %s in suggestions", tt.wantServiceExists)
				}
			}
		})
	}
}
