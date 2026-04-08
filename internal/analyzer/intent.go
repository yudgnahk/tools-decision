package analyzer

import (
	"regexp"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// IntentAnalyzer analyzes a project idea/description to infer the tech stack
type IntentAnalyzer struct {
	patterns []IntentPattern
}

// IntentPattern maps keywords to inferred technologies
type IntentPattern struct {
	Keywords    []string
	Language    string
	Framework   string
	Service     string
	ProjectType types.ProjectType
	Confidence  float64
}

// NewIntentAnalyzer creates a new intent analyzer with default patterns
func NewIntentAnalyzer() *IntentAnalyzer {
	return &IntentAnalyzer{
		patterns: defaultPatterns(),
	}
}

// AnalyzeIdea analyzes a project idea and returns inferred context
func (a *IntentAnalyzer) AnalyzeIdea(idea string) *types.ProjectContext {
	idea = strings.ToLower(idea)

	ctx := &types.ProjectContext{
		Type: types.ProjectTypeUnknown,
	}

	languageScores := make(map[string]float64)
	frameworkScores := make(map[string]float64)
	serviceScores := make(map[string]float64)

	for _, pattern := range a.patterns {
		if matchesAny(idea, pattern.Keywords) {
			if pattern.Language != "" {
				languageScores[pattern.Language] += pattern.Confidence
			}
			if pattern.Framework != "" {
				frameworkScores[pattern.Framework] += pattern.Confidence
			}
			if pattern.Service != "" {
				serviceScores[pattern.Service] += pattern.Confidence
			}
			if pattern.ProjectType != "" && pattern.ProjectType != types.ProjectTypeUnknown {
				ctx.Type = pattern.ProjectType
			}
		}
	}

	// Convert scores to context
	for lang, score := range languageScores {
		ctx.Languages = append(ctx.Languages, types.Language{
			Name:       lang,
			Confidence: min(score, 1.0),
		})
	}

	for fw, score := range frameworkScores {
		ctx.Frameworks = append(ctx.Frameworks, types.Framework{
			Name:       fw,
			Confidence: min(score, 1.0),
		})
	}

	for svc, score := range serviceScores {
		ctx.Services = append(ctx.Services, types.Service{
			Name:       svc,
			Confidence: min(score, 1.0),
		})
	}

	// Detect use cases from the idea
	ctx.UseCases = detectUseCases(idea)

	return ctx
}

// detectUseCases analyzes text to detect intended use cases
func detectUseCases(text string) []types.UseCase {
	useCasePatterns := []struct {
		UseCase    string
		Keywords   []string
		Confidence float64
	}{
		// Debugging use cases
		{types.UseCaseDebugging, []string{"debug", "fix", "error", "bug", "issue", "problem", "crash", "failing"}, 0.9},
		{types.UseCaseDebugging, []string{"troubleshoot", "investigate", "diagnose"}, 0.85},

		// Code review use cases
		{types.UseCaseCodeReview, []string{"review", "audit", "check", "inspect"}, 0.85},
		{types.UseCaseCodeReview, []string{"pr review", "pull request", "code quality"}, 0.9},

		// Architecture use cases
		{types.UseCaseArchitecture, []string{"design", "architect", "structure", "organize"}, 0.85},
		{types.UseCaseArchitecture, []string{"microservice", "monolith", "modular"}, 0.9},

		// Testing use cases
		{types.UseCaseTesting, []string{"test", "testing", "coverage", "tdd"}, 0.9},
		{types.UseCaseTesting, []string{"unit test", "integration test", "e2e"}, 0.95},

		// DevOps use cases
		{types.UseCaseDevOps, []string{"deploy", "ci", "cd", "pipeline", "devops"}, 0.9},
		{types.UseCaseDevOps, []string{"kubernetes", "docker", "container", "infrastructure"}, 0.85},

		// Performance use cases
		{types.UseCasePerformance, []string{"performance", "optimize", "speed", "slow", "fast"}, 0.85},
		{types.UseCasePerformance, []string{"profiling", "benchmark", "latency"}, 0.9},

		// Security use cases
		{types.UseCaseSecurity, []string{"security", "secure", "vulnerability", "auth"}, 0.9},
		{types.UseCaseSecurity, []string{"penetration", "exploit", "hack"}, 0.85},

		// Refactoring use cases
		{types.UseCaseRefactoring, []string{"refactor", "cleanup", "improve", "modernize"}, 0.85},
		{types.UseCaseRefactoring, []string{"rewrite", "restructure", "simplify"}, 0.8},

		// API Design use cases
		{types.UseCaseAPIDesign, []string{"api", "rest", "graphql", "endpoint"}, 0.85},
		{types.UseCaseAPIDesign, []string{"openapi", "swagger", "grpc"}, 0.9},

		// Database use cases
		{types.UseCaseDatabaseDesign, []string{"database", "schema", "migration"}, 0.85},
		{types.UseCaseDatabaseDesign, []string{"query", "sql", "orm"}, 0.8},

		// Documentation use cases
		{types.UseCaseDocumentation, []string{"document", "documentation", "readme", "docs"}, 0.9},
		{types.UseCaseDocumentation, []string{"comment", "explain", "describe"}, 0.75},
	}

	useCaseScores := make(map[string]float64)

	for _, pattern := range useCasePatterns {
		if matchesAny(text, pattern.Keywords) {
			if pattern.Confidence > useCaseScores[pattern.UseCase] {
				useCaseScores[pattern.UseCase] = pattern.Confidence
			}
		}
	}

	var useCases []types.UseCase
	for uc, score := range useCaseScores {
		useCases = append(useCases, types.UseCase{
			Name:       uc,
			Confidence: score,
		})
	}

	return useCases
}

// matchesAny checks if the text contains any of the keywords
func matchesAny(text string, keywords []string) bool {
	for _, kw := range keywords {
		// Use word boundary matching for more accurate results
		pattern := `\b` + regexp.QuoteMeta(strings.ToLower(kw)) + `\b`
		if matched, _ := regexp.MatchString(pattern, text); matched {
			return true
		}
	}
	return false
}

// defaultPatterns returns the default intent patterns
func defaultPatterns() []IntentPattern {
	return []IntentPattern{
		// Project Types
		{Keywords: []string{"api", "rest", "backend", "server", "microservice"}, ProjectType: types.ProjectTypeAPI, Confidence: 0.8},
		{Keywords: []string{"web app", "webapp", "website", "frontend"}, ProjectType: types.ProjectTypeWebApp, Confidence: 0.8},
		{Keywords: []string{"cli", "command line", "terminal", "console"}, ProjectType: types.ProjectTypeCLI, Confidence: 0.9},
		{Keywords: []string{"library", "package", "sdk", "module"}, ProjectType: types.ProjectTypeLibrary, Confidence: 0.8},

		// Languages
		{Keywords: []string{"typescript", "ts"}, Language: "typescript", Confidence: 0.95},
		{Keywords: []string{"javascript", "js", "node", "nodejs"}, Language: "javascript", Confidence: 0.9},
		{Keywords: []string{"python", "py", "django", "flask", "fastapi"}, Language: "python", Confidence: 0.9},
		{Keywords: []string{"go", "golang"}, Language: "go", Confidence: 0.95},
		{Keywords: []string{"rust"}, Language: "rust", Confidence: 0.95},
		{Keywords: []string{"java", "spring", "kotlin"}, Language: "java", Confidence: 0.9},

		// Frameworks - JavaScript/TypeScript
		{Keywords: []string{"next", "nextjs", "next.js"}, Framework: "nextjs", Language: "typescript", Confidence: 0.95},
		{Keywords: []string{"react"}, Framework: "react", Language: "typescript", Confidence: 0.9},
		{Keywords: []string{"vue", "vuejs", "nuxt"}, Framework: "vue", Language: "typescript", Confidence: 0.9},
		{Keywords: []string{"express", "expressjs"}, Framework: "express", Language: "javascript", Confidence: 0.9},
		{Keywords: []string{"nest", "nestjs"}, Framework: "nestjs", Language: "typescript", Confidence: 0.95},
		{Keywords: []string{"svelte", "sveltekit"}, Framework: "svelte", Language: "typescript", Confidence: 0.9},
		{Keywords: []string{"astro"}, Framework: "astro", Language: "typescript", Confidence: 0.9},

		// Frameworks - Python
		{Keywords: []string{"fastapi", "fast api"}, Framework: "fastapi", Language: "python", Confidence: 0.95},
		{Keywords: []string{"django"}, Framework: "django", Language: "python", Confidence: 0.95},
		{Keywords: []string{"flask"}, Framework: "flask", Language: "python", Confidence: 0.95},
		{Keywords: []string{"streamlit"}, Framework: "streamlit", Language: "python", Confidence: 0.9},

		// Frameworks - Go
		{Keywords: []string{"gin"}, Framework: "gin", Language: "go", Confidence: 0.95},
		{Keywords: []string{"echo"}, Framework: "echo", Language: "go", Confidence: 0.95},
		{Keywords: []string{"fiber"}, Framework: "fiber", Language: "go", Confidence: 0.95},
		{Keywords: []string{"chi"}, Framework: "chi", Language: "go", Confidence: 0.9},
		{Keywords: []string{"cobra"}, Framework: "cobra", Language: "go", ProjectType: types.ProjectTypeCLI, Confidence: 0.9},

		// Frameworks - Rust
		{Keywords: []string{"rust"}, Language: "rust", Confidence: 0.95},
		{Keywords: []string{"actix", "actix-web"}, Framework: "actix", Language: "rust", Confidence: 0.95},
		{Keywords: []string{"axum"}, Framework: "axum", Language: "rust", Confidence: 0.95},
		{Keywords: []string{"rocket"}, Framework: "rocket", Language: "rust", Confidence: 0.95},
		{Keywords: []string{"tokio"}, Framework: "tokio", Language: "rust", Confidence: 0.9},
		{Keywords: []string{"warp"}, Framework: "warp", Language: "rust", Confidence: 0.9},
		{Keywords: []string{"tauri"}, Framework: "tauri", Language: "rust", ProjectType: types.ProjectTypeDesktop, Confidence: 0.95},
		{Keywords: []string{"clap"}, Framework: "clap", Language: "rust", ProjectType: types.ProjectTypeCLI, Confidence: 0.9},

		// Frameworks - Java/Kotlin
		{Keywords: []string{"spring boot", "springboot"}, Framework: "spring-boot", Language: "java", Confidence: 0.95},
		{Keywords: []string{"spring"}, Framework: "spring", Language: "java", Confidence: 0.9},
		{Keywords: []string{"quarkus"}, Framework: "quarkus", Language: "java", Confidence: 0.95},
		{Keywords: []string{"micronaut"}, Framework: "micronaut", Language: "java", Confidence: 0.95},
		{Keywords: []string{"kotlin"}, Language: "kotlin", Confidence: 0.95},
		{Keywords: []string{"ktor"}, Framework: "ktor", Language: "kotlin", Confidence: 0.95},
		{Keywords: []string{"android"}, Language: "kotlin", ProjectType: types.ProjectTypeMobile, Confidence: 0.9},

		// Databases
		{Keywords: []string{"postgres", "postgresql", "pg"}, Service: "postgresql", Confidence: 0.95},
		{Keywords: []string{"mysql", "mariadb"}, Service: "mysql", Confidence: 0.95},
		{Keywords: []string{"mongodb", "mongo"}, Service: "mongodb", Confidence: 0.95},
		{Keywords: []string{"sqlite"}, Service: "sqlite", Confidence: 0.9},
		{Keywords: []string{"redis", "cache", "caching"}, Service: "redis", Confidence: 0.85},
		{Keywords: []string{"database", "db", "sql"}, Service: "database", Confidence: 0.7},
		{Keywords: []string{"prisma"}, Service: "prisma", Framework: "prisma", Confidence: 0.9},
		{Keywords: []string{"drizzle"}, Service: "drizzle", Framework: "drizzle", Confidence: 0.9},

		// Cloud/Infrastructure
		{Keywords: []string{"docker", "container", "containerized"}, Service: "docker", Confidence: 0.9},
		{Keywords: []string{"kubernetes", "k8s"}, Service: "kubernetes", Confidence: 0.9},
		{Keywords: []string{"aws", "amazon"}, Service: "aws", Confidence: 0.85},
		{Keywords: []string{"s3", "storage", "file upload", "uploads"}, Service: "s3", Confidence: 0.8},
		{Keywords: []string{"gcp", "google cloud"}, Service: "gcp", Confidence: 0.85},
		{Keywords: []string{"azure"}, Service: "azure", Confidence: 0.85},
		{Keywords: []string{"vercel"}, Service: "vercel", Confidence: 0.9},
		{Keywords: []string{"cloudflare"}, Service: "cloudflare", Confidence: 0.9},

		// Features
		{Keywords: []string{"auth", "authentication", "login", "oauth", "jwt"}, Service: "auth", Confidence: 0.85},
		{Keywords: []string{"stripe", "payment", "payments", "billing"}, Service: "stripe", Confidence: 0.9},
		{Keywords: []string{"email", "mail", "sendgrid", "ses"}, Service: "email", Confidence: 0.8},
		{Keywords: []string{"websocket", "realtime", "real-time", "socket"}, Service: "websocket", Confidence: 0.85},
		{Keywords: []string{"graphql"}, Service: "graphql", Framework: "graphql", Confidence: 0.9},
		{Keywords: []string{"grpc"}, Service: "grpc", Framework: "grpc", Confidence: 0.9},

		// AI/ML
		{Keywords: []string{"ai", "llm", "gpt", "openai", "chatgpt"}, Service: "openai", Confidence: 0.85},
		{Keywords: []string{"anthropic", "claude"}, Service: "anthropic", Confidence: 0.9},
		{Keywords: []string{"langchain"}, Service: "langchain", Framework: "langchain", Confidence: 0.9},
		{Keywords: []string{"vector", "embedding", "pinecone", "qdrant"}, Service: "vectordb", Confidence: 0.85},

		// Testing
		{Keywords: []string{"jest", "vitest"}, Service: "jest", Confidence: 0.9},
		{Keywords: []string{"pytest"}, Service: "pytest", Confidence: 0.9},
		{Keywords: []string{"testing", "tests", "test"}, Service: "testing", Confidence: 0.7},
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
