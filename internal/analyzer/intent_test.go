package analyzer

import (
	"testing"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

func TestIntentAnalyzer_AnalyzeIdea(t *testing.T) {
	analyzer := NewIntentAnalyzer()

	tests := []struct {
		name          string
		idea          string
		wantLanguage  string
		wantFramework string
		wantService   string
		wantType      types.ProjectType
	}{
		{
			name:        "REST API with PostgreSQL",
			idea:        "REST API with PostgreSQL and Redis",
			wantService: "postgresql",
			wantType:    types.ProjectTypeAPI,
		},
		{
			name:          "Next.js app",
			idea:          "Next.js app with authentication",
			wantLanguage:  "typescript",
			wantFramework: "nextjs",
			wantType:      types.ProjectTypeUnknown, // web_app depends on keyword matching
		},
		{
			name:          "Go CLI with Cobra",
			idea:          "CLI tool in Go using Cobra",
			wantLanguage:  "go",
			wantFramework: "cobra",
			wantType:      types.ProjectTypeCLI,
		},
		{
			name:          "Python FastAPI",
			idea:          "FastAPI backend with MongoDB",
			wantLanguage:  "python",
			wantFramework: "fastapi",
			wantService:   "mongodb",
			wantType:      types.ProjectTypeAPI, // "backend" triggers API type
		},
		{
			name:          "Rust web API",
			idea:          "Rust web API with Axum and PostgreSQL",
			wantLanguage:  "rust",
			wantFramework: "axum",
			wantService:   "postgresql",
			wantType:      types.ProjectTypeAPI, // "api" triggers API type
		},
		{
			name:          "Spring Boot API",
			idea:          "Spring Boot REST API with MySQL",
			wantLanguage:  "java",
			wantFramework: "spring-boot",
			wantService:   "mysql",
			wantType:      types.ProjectTypeAPI, // "rest" and "api" trigger API type
		},
		{
			name:        "Docker and Kubernetes",
			idea:        "Deploy with Docker and Kubernetes",
			wantService: "docker",
			// No specific project type expected
		},
		{
			name:        "AI/LLM integration",
			idea:        "Chatbot with OpenAI GPT",
			wantService: "openai",
			// No specific project type expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := analyzer.AnalyzeIdea(tt.idea)

			if tt.wantLanguage != "" {
				found := false
				for _, l := range ctx.Languages {
					if l.Name == tt.wantLanguage {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected language %s, got %v", tt.wantLanguage, ctx.Languages)
				}
			}

			if tt.wantFramework != "" {
				found := false
				for _, f := range ctx.Frameworks {
					if f.Name == tt.wantFramework {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected framework %s, got %v", tt.wantFramework, ctx.Frameworks)
				}
			}

			if tt.wantService != "" {
				found := false
				for _, s := range ctx.Services {
					if s.Name == tt.wantService {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected service %s, got %v", tt.wantService, ctx.Services)
				}
			}

			if tt.wantType != "" && ctx.Type != tt.wantType {
				t.Errorf("expected type %s, got %s", tt.wantType, ctx.Type)
			}
		})
	}
}

func TestMatchesAny(t *testing.T) {
	tests := []struct {
		text     string
		keywords []string
		want     bool
	}{
		{"rest api with postgresql", []string{"api", "rest"}, true},
		{"build a web app", []string{"web app"}, true},
		{"go programming", []string{"golang", "go"}, true},
		{"python script", []string{"ruby"}, false},
		{"", []string{"api"}, false},
		{"hello world", []string{}, false},
	}

	for _, tt := range tests {
		got := matchesAny(tt.text, tt.keywords)
		if got != tt.want {
			t.Errorf("matchesAny(%q, %v) = %v, want %v", tt.text, tt.keywords, got, tt.want)
		}
	}
}
