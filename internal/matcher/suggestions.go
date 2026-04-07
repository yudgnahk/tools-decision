package matcher

import (
	"github.com/yudgnahk/tools-decision/pkg/types"
)

// ToolSuggestion maps services/features to recommended MCP tools
type ToolSuggestion struct {
	Service     string   // The detected service/need
	Tools       []string // Recommended tool slugs
	Description string   // Why these tools
}

// GetSuggestionsForContext returns tool suggestions based on project context
func GetSuggestionsForContext(ctx *types.ProjectContext) []ToolSuggestion {
	var suggestions []ToolSuggestion

	// Language-based suggestions
	for _, lang := range ctx.Languages {
		switch lang.Name {
		case "typescript", "javascript":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "javascript",
				Tools:       []string{"filesystem", "git"},
				Description: "Core tools for JS/TS development",
			})
		case "python":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "python",
				Tools:       []string{"filesystem", "git", "python-repl"},
				Description: "Core tools for Python development",
			})
		case "go":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "go",
				Tools:       []string{"filesystem", "git"},
				Description: "Core tools for Go development",
			})
		}
	}

	// Service-based suggestions
	for _, svc := range ctx.Services {
		switch svc.Name {
		case "postgresql":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "postgresql",
				Tools:       []string{"postgres-mcp", "database-tools"},
				Description: "PostgreSQL database operations and schema management",
			})
		case "mysql":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "mysql",
				Tools:       []string{"mysql-mcp", "database-tools"},
				Description: "MySQL database operations",
			})
		case "mongodb":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "mongodb",
				Tools:       []string{"mongodb-mcp"},
				Description: "MongoDB document operations",
			})
		case "redis":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "redis",
				Tools:       []string{"redis-mcp"},
				Description: "Redis cache operations",
			})
		case "docker":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "docker",
				Tools:       []string{"docker-mcp"},
				Description: "Docker container management",
			})
		case "kubernetes":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "kubernetes",
				Tools:       []string{"kubernetes-mcp", "kubectl"},
				Description: "Kubernetes cluster operations",
			})
		case "aws", "s3":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "aws",
				Tools:       []string{"aws-mcp", "s3-mcp"},
				Description: "AWS cloud services integration",
			})
		case "stripe":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "stripe",
				Tools:       []string{"stripe-mcp"},
				Description: "Stripe payment processing",
			})
		case "auth":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "auth",
				Tools:       []string{"auth0-mcp", "jwt-tools"},
				Description: "Authentication and authorization tools",
			})
		case "openai":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "openai",
				Tools:       []string{"openai-mcp"},
				Description: "OpenAI API integration",
			})
		case "anthropic":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "anthropic",
				Tools:       []string{"anthropic-mcp"},
				Description: "Anthropic Claude API integration",
			})
		case "graphql":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "graphql",
				Tools:       []string{"graphql-mcp"},
				Description: "GraphQL API tools",
			})
		case "grpc":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "grpc",
				Tools:       []string{"grpc-mcp", "protobuf-tools"},
				Description: "gRPC service tools",
			})
		case "email":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "email",
				Tools:       []string{"email-mcp", "sendgrid-mcp"},
				Description: "Email sending and management",
			})
		case "testing", "jest", "pytest":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "testing",
				Tools:       []string{"test-runner-mcp"},
				Description: "Test execution and coverage",
			})
		case "prisma":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "prisma",
				Tools:       []string{"prisma-mcp"},
				Description: "Prisma ORM schema and migrations",
			})
		}
	}

	// Framework-based suggestions
	for _, fw := range ctx.Frameworks {
		switch fw.Name {
		case "nextjs":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "nextjs",
				Tools:       []string{"nextjs-mcp", "vercel-mcp"},
				Description: "Next.js app development tools",
			})
		case "react":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "react",
				Tools:       []string{"react-devtools-mcp"},
				Description: "React component debugging",
			})
		case "fastapi":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "fastapi",
				Tools:       []string{"fastapi-mcp", "openapi-mcp"},
				Description: "FastAPI development tools",
			})
		case "django":
			suggestions = append(suggestions, ToolSuggestion{
				Service:     "django",
				Tools:       []string{"django-mcp"},
				Description: "Django development tools",
			})
		}
	}

	// Always include these core tools
	suggestions = append(suggestions, ToolSuggestion{
		Service:     "core",
		Tools:       []string{"filesystem", "git", "github"},
		Description: "Essential development tools",
	})

	return deduplicateSuggestions(suggestions)
}

// deduplicateSuggestions removes duplicate tool recommendations
func deduplicateSuggestions(suggestions []ToolSuggestion) []ToolSuggestion {
	seen := make(map[string]bool)
	var result []ToolSuggestion

	for _, s := range suggestions {
		key := s.Service
		if !seen[key] {
			seen[key] = true
			result = append(result, s)
		}
	}

	return result
}

// GetAllSuggestedTools returns a flat list of all suggested tool slugs
func GetAllSuggestedTools(suggestions []ToolSuggestion) []string {
	seen := make(map[string]bool)
	var tools []string

	for _, s := range suggestions {
		for _, tool := range s.Tools {
			if !seen[tool] {
				seen[tool] = true
				tools = append(tools, tool)
			}
		}
	}

	return tools
}
