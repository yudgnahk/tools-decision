package analyzer

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// PythonDetector detects Python projects
type PythonDetector struct{}

// NewPythonDetector creates a new Python detector
func NewPythonDetector() *PythonDetector {
	return &PythonDetector{}
}

// Name returns the detector name
func (d *PythonDetector) Name() string {
	return "python"
}

// Detect analyzes the project for Python
func (d *PythonDetector) Detect(ctx context.Context, projectPath string) (*DetectorResult, error) {
	result := &DetectorResult{}

	// Check for requirements.txt
	reqPath := filepath.Join(projectPath, "requirements.txt")
	deps := make(map[string]bool)

	if data, err := os.ReadFile(reqPath); err == nil {
		result.Languages = append(result.Languages, types.Language{
			Name:       "python",
			Confidence: 0.95,
		})
		deps = parseRequirements(string(data))
	}

	// Check for pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if data, err := os.ReadFile(pyprojectPath); err == nil {
		if len(result.Languages) == 0 {
			result.Languages = append(result.Languages, types.Language{
				Name:       "python",
				Confidence: 0.95,
			})
		}
		for dep := range parsePyprojectDependencies(string(data)) {
			deps[dep] = true
		}
	}

	// Check for setup.py
	setupPath := filepath.Join(projectPath, "setup.py")
	if _, err := os.Stat(setupPath); err == nil {
		if len(result.Languages) == 0 {
			result.Languages = append(result.Languages, types.Language{
				Name:       "python",
				Confidence: 0.9,
			})
		}
	}

	if len(result.Languages) == 0 {
		return nil, nil
	}

	// Detect frameworks
	frameworkMap := map[string]string{
		"fastapi":   "fastapi",
		"django":    "django",
		"flask":     "flask",
		"starlette": "starlette",
		"tornado":   "tornado",
		"aiohttp":   "aiohttp",
		"streamlit": "streamlit",
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
		"pytest": "pytest",
		"black":  "black",
		"ruff":   "ruff",
		"mypy":   "mypy",
		"pylint": "pylint",
		"flake8": "flake8",
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
		"psycopg2":        "postgresql",
		"psycopg2-binary": "postgresql",
		"asyncpg":         "postgresql",
		"pymysql":         "mysql",
		"pymongo":         "mongodb",
		"redis":           "redis",
		"sqlalchemy":      "database",
		"boto3":           "aws",
		"openai":          "openai",
		"anthropic":       "anthropic",
		"langchain":       "langchain",
		"playwright":      "playwright",
		"selenium":        "selenium",
		"pandas":          "pandas",
		"openpyxl":        "openpyxl",
		"moviepy":         "moviepy",
		"pydub":           "pydub",
		"ffmpeg-python":   "ffmpeg",
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
	if deps["fastapi"] || deps["django"] || deps["flask"] {
		result.Type = types.ProjectTypeAPI
	} else if deps["click"] || deps["typer"] {
		result.Type = types.ProjectTypeCLI
	}

	return result, nil
}

// parseRequirements parses a requirements.txt file
func parseRequirements(content string) map[string]bool {
	deps := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Extract package name (before ==, >=, etc.)
		name := line
		for _, sep := range []string{"==", ">=", "<=", ">", "<", "~=", "!="} {
			if idx := strings.Index(line, sep); idx > 0 {
				name = line[:idx]
				break
			}
		}

		// Handle extras like package[extra]
		if idx := strings.Index(name, "["); idx > 0 {
			name = name[:idx]
		}

		deps[strings.ToLower(strings.TrimSpace(name))] = true
	}

	return deps
}

func parsePyprojectDependencies(content string) map[string]bool {
	deps := make(map[string]bool)
	lines := strings.Split(content, "\n")
	inProjectDeps := false
	inProjectDepList := false
	inPoetryDeps := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.ToLower(strings.Trim(line, "[]"))
			inProjectDeps = false
			inProjectDepList = false
			inPoetryDeps = false

			switch {
			case section == "project":
				inProjectDeps = true
			case section == "tool.poetry.dependencies":
				inPoetryDeps = true
			case strings.HasPrefix(section, "tool.poetry.group.") && strings.HasSuffix(section, ".dependencies"):
				inPoetryDeps = true
			}
			continue
		}

		if inProjectDeps && strings.HasPrefix(strings.ToLower(line), "dependencies") {
			inProjectDepList = true
			for _, token := range quotedTokens(line) {
				name := normalizeDepName(token)
				if name != "" {
					deps[name] = true
				}
			}
			if strings.Contains(line, "]") {
				inProjectDepList = false
			}
			continue
		}

		if inProjectDepList {
			for _, token := range quotedTokens(line) {
				name := normalizeDepName(token)
				if name != "" {
					deps[name] = true
				}
			}
			if strings.Contains(line, "]") {
				inProjectDepList = false
			}
			continue
		}

		if inPoetryDeps {
			if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "}") {
				continue
			}
			if idx := strings.Index(line, "="); idx > 0 {
				name := strings.ToLower(strings.TrimSpace(line[:idx]))
				if name != "" && name != "python" {
					deps[name] = true
				}
			}
		}
	}

	return deps
}

func quotedTokens(line string) []string {
	var tokens []string
	for {
		start := strings.Index(line, "\"")
		if start < 0 {
			return tokens
		}
		line = line[start+1:]
		end := strings.Index(line, "\"")
		if end < 0 {
			return tokens
		}
		tokens = append(tokens, line[:end])
		line = line[end+1:]
	}
}

func normalizeDepName(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		switch r {
		case ' ', '=', '>', '<', '!', '~', '[', ';':
			return true
		default:
			return false
		}
	})
	if len(parts) == 0 {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(parts[0]))
}
