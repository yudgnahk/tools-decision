package analyzer

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// GoDetector detects Go projects
type GoDetector struct{}

// NewGoDetector creates a new Go detector
func NewGoDetector() *GoDetector {
	return &GoDetector{}
}

// Name returns the detector name
func (d *GoDetector) Name() string {
	return "go"
}

// Detect analyzes the project for Go
func (d *GoDetector) Detect(ctx context.Context, projectPath string) (*DetectorResult, error) {
	result := &DetectorResult{}

	// Check for go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, nil
	}

	result.Languages = append(result.Languages, types.Language{
		Name:       "go",
		Confidence: 0.98,
	})

	deps := parseGoMod(string(data))

	// Detect frameworks
	frameworkMap := map[string]string{
		"github.com/gin-gonic/gin":            "gin",
		"github.com/labstack/echo":            "echo",
		"github.com/gofiber/fiber":            "fiber",
		"github.com/gorilla/mux":              "gorilla",
		"github.com/go-chi/chi":               "chi",
		"github.com/julienschmidt/httprouter": "httprouter",
		"google.golang.org/grpc":              "grpc",
		"github.com/spf13/cobra":              "cobra",
		"github.com/urfave/cli":               "cli",
	}

	for dep, framework := range frameworkMap {
		if deps[dep] {
			result.Frameworks = append(result.Frameworks, types.Framework{
				Name:       framework,
				Confidence: 0.95,
			})
		}
	}

	// Detect tools
	toolDeps := map[string]string{
		"github.com/stretchr/testify": "testify",
		"github.com/onsi/ginkgo":      "ginkgo",
		"github.com/onsi/gomega":      "gomega",
	}

	for dep, tool := range toolDeps {
		if deps[dep] {
			result.Tools = append(result.Tools, types.Tool{
				Name: tool,
			})
		}
	}

	// Detect services
	serviceMap := map[string]string{
		"github.com/lib/pq":                             "postgresql",
		"github.com/jackc/pgx":                          "postgresql",
		"github.com/go-sql-driver/mysql":                "mysql",
		"go.mongodb.org/mongo-driver":                   "mongodb",
		"github.com/go-redis/redis":                     "redis",
		"github.com/redis/go-redis":                     "redis",
		"github.com/aws/aws-sdk-go":                     "aws",
		"github.com/aws/aws-sdk-go-v2":                  "aws",
		"gorm.io/gorm":                                  "database",
		"github.com/jmoiron/sqlx":                       "database",
		"github.com/playwright-community/playwright-go": "playwright",
		"github.com/tebeka/selenium":                    "selenium",
		"github.com/xuri/excelize":                      "excelize",
		"github.com/xuri/excelize/v2":                   "excelize",
		"github.com/sashabaranov/go-openai":             "openai",
		"github.com/tmc/langchaingo":                    "langchain",
	}

	for dep, service := range serviceMap {
		if deps[dep] {
			result.Services = append(result.Services, types.Service{
				Name:       service,
				Confidence: 0.9,
			})
		}
	}

	// Detect project type
	if deps["github.com/gin-gonic/gin"] || deps["github.com/labstack/echo"] ||
		deps["github.com/gofiber/fiber"] || deps["google.golang.org/grpc"] {
		result.Type = types.ProjectTypeAPI
	} else if deps["github.com/spf13/cobra"] || deps["github.com/urfave/cli"] {
		result.Type = types.ProjectTypeCLI
	}

	// Check for cmd directory (often CLI entrypoints), but do not override API.
	cmdPath := filepath.Join(projectPath, "cmd")
	if info, err := os.Stat(cmdPath); err == nil && info.IsDir() && result.Type == types.ProjectTypeUnknown {
		result.Type = types.ProjectTypeCLI
	}

	return result, nil
}

// parseGoMod parses a go.mod file and returns dependencies
func parseGoMod(content string) map[string]bool {
	deps := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(content))
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Handle require block
		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}

		// Handle single require
		if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps[parts[1]] = true
			}
			continue
		}

		// Inside require block
		if inRequire {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				// Remove version suffix for matching
				dep := parts[0]
				deps[dep] = true

				// Also add base path for matching
				// e.g., github.com/jackc/pgx/v5 -> github.com/jackc/pgx
				for _, suffix := range []string{"/v2", "/v3", "/v4", "/v5"} {
					if strings.HasSuffix(dep, suffix) {
						deps[strings.TrimSuffix(dep, suffix)] = true
					}
				}
			}
		}
	}

	return deps
}
