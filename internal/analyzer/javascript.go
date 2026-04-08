package analyzer

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// JavaScriptDetector detects JavaScript/TypeScript projects
type JavaScriptDetector struct{}

// NewJavaScriptDetector creates a new JavaScript detector
func NewJavaScriptDetector() *JavaScriptDetector {
	return &JavaScriptDetector{}
}

// Name returns the detector name
func (d *JavaScriptDetector) Name() string {
	return "javascript"
}

// PackageJSON represents a package.json file
type PackageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

// Detect analyzes the project for JavaScript/TypeScript
func (d *JavaScriptDetector) Detect(ctx context.Context, projectPath string) (*DetectorResult, error) {
	result := &DetectorResult{}

	// Check for package.json
	pkgPath := filepath.Join(projectPath, "package.json")
	pkgData, err := os.ReadFile(pkgPath)
	if err != nil {
		// No package.json, not a JS project
		return nil, nil
	}

	var pkg PackageJSON
	if err := json.Unmarshal(pkgData, &pkg); err != nil {
		return nil, nil
	}

	// Detect TypeScript
	if _, hasTS := pkg.DevDependencies["typescript"]; hasTS {
		result.Languages = append(result.Languages, types.Language{
			Name:       "typescript",
			Confidence: 0.95,
		})
	}

	// Always add JavaScript if we found package.json
	result.Languages = append(result.Languages, types.Language{
		Name:       "javascript",
		Confidence: 0.9,
	})

	// Merge all dependencies
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	// Detect frameworks
	frameworkMap := map[string]string{
		"next":               "nextjs",
		"react":              "react",
		"vue":                "vue",
		"@angular/core":      "angular",
		"svelte":             "svelte",
		"express":            "express",
		"fastify":            "fastify",
		"@nestjs/core":       "nestjs",
		"nuxt":               "nuxt",
		"remix":              "remix",
		"astro":              "astro",
		"playwright":         "playwright",
		"puppeteer":          "puppeteer",
		"selenium-webdriver": "selenium",
	}

	for dep, framework := range frameworkMap {
		if version, ok := allDeps[dep]; ok {
			result.Frameworks = append(result.Frameworks, types.Framework{
				Name:       framework,
				Version:    cleanVersion(version),
				Confidence: 0.95,
			})
		}
	}

	// Detect tools
	toolMap := map[string]string{
		"eslint":      ".eslintrc",
		"prettier":    ".prettierrc",
		"jest":        "jest.config",
		"vitest":      "vitest.config",
		"webpack":     "webpack.config",
		"vite":        "vite.config",
		"tailwindcss": "tailwind.config",
	}

	for dep, configPrefix := range toolMap {
		if _, ok := allDeps[dep]; ok {
			result.Tools = append(result.Tools, types.Tool{
				Name: dep,
			})
		} else {
			// Check for config file
			matches, _ := filepath.Glob(filepath.Join(projectPath, configPrefix+"*"))
			if len(matches) > 0 {
				result.Tools = append(result.Tools, types.Tool{
					Name:       dep,
					ConfigFile: filepath.Base(matches[0]),
				})
			}
		}
	}

	// Detect services from dependencies
	serviceMap := map[string]string{
		"pg":                 "postgresql",
		"mysql2":             "mysql",
		"mongodb":            "mongodb",
		"redis":              "redis",
		"ioredis":            "redis",
		"@prisma/client":     "prisma",
		"prisma":             "prisma",
		"aws-sdk":            "aws",
		"@aws-sdk/client-s3": "s3",
		"openai":             "openai",
		"@anthropic-ai/sdk":  "anthropic",
		"langchain":          "langchain",
		"playwright":         "playwright",
		"puppeteer":          "puppeteer",
		"selenium-webdriver": "selenium",
		"fluent-ffmpeg":      "ffmpeg",
	}

	for dep, service := range serviceMap {
		if _, ok := allDeps[dep]; ok {
			result.Services = append(result.Services, types.Service{
				Name:       service,
				Confidence: 0.9,
			})
		}
	}

	// Detect project type
	if _, ok := allDeps["next"]; ok {
		result.Type = types.ProjectTypeWebApp
	} else if _, ok := allDeps["express"]; ok {
		result.Type = types.ProjectTypeAPI
	} else if pkg.Scripts["start"] != "" && strings.Contains(pkg.Scripts["start"], "node") {
		result.Type = types.ProjectTypeCLI
	}

	return result, nil
}

// cleanVersion removes common version prefixes
func cleanVersion(version string) string {
	version = strings.TrimPrefix(version, "^")
	version = strings.TrimPrefix(version, "~")
	return version
}
