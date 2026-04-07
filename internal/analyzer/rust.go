package analyzer

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// RustDetector detects Rust projects
type RustDetector struct{}

// NewRustDetector creates a new Rust detector
func NewRustDetector() *RustDetector {
	return &RustDetector{}
}

// Name returns the detector name
func (d *RustDetector) Name() string {
	return "rust"
}

// Detect analyzes the project for Rust
func (d *RustDetector) Detect(ctx context.Context, projectPath string) (*DetectorResult, error) {
	result := &DetectorResult{}

	// Check for Cargo.toml
	cargoPath := filepath.Join(projectPath, "Cargo.toml")
	data, err := os.ReadFile(cargoPath)
	if err != nil {
		return nil, nil
	}

	result.Languages = append(result.Languages, types.Language{
		Name:       "rust",
		Confidence: 0.98,
	})

	deps := parseCargoToml(string(data))
	cargoContent := string(data)

	// Detect frameworks
	frameworkMap := map[string]string{
		"actix-web": "actix",
		"axum":      "axum",
		"rocket":    "rocket",
		"warp":      "warp",
		"hyper":     "hyper",
		"tonic":     "tonic", // gRPC
		"clap":      "clap",
		"structopt": "structopt",
		"tokio":     "tokio",
		"async-std": "async-std",
		"tauri":     "tauri",
		"yew":       "yew",
		"leptos":    "leptos",
		"dioxus":    "dioxus",
		"bevy":      "bevy",
	}

	for dep, framework := range frameworkMap {
		if deps[dep] {
			result.Frameworks = append(result.Frameworks, types.Framework{
				Name:       framework,
				Confidence: 0.95,
			})
		}
	}

	// Detect services
	serviceMap := map[string]string{
		"sqlx":           "database",
		"diesel":         "database",
		"sea-orm":        "database",
		"tokio-postgres": "postgresql",
		"postgres":       "postgresql",
		"mysql":          "mysql",
		"mongodb":        "mongodb",
		"redis":          "redis",
		"aws-sdk":        "aws",
		"rusoto":         "aws",
		"lapin":          "rabbitmq",
		"rdkafka":        "kafka",
	}

	for dep, service := range serviceMap {
		if deps[dep] {
			result.Services = append(result.Services, types.Service{
				Name:       service,
				Confidence: 0.9,
			})
		}
	}

	// Detect project type based on Cargo.toml content and frameworks
	if deps["actix-web"] || deps["axum"] || deps["rocket"] || deps["warp"] || deps["tonic"] {
		result.Type = types.ProjectTypeAPI
	} else if deps["clap"] || deps["structopt"] {
		result.Type = types.ProjectTypeCLI
	} else if deps["tauri"] || deps["yew"] || deps["leptos"] || deps["dioxus"] {
		result.Type = types.ProjectTypeDesktop
	} else if deps["bevy"] {
		result.Type = types.ProjectTypeUnknown // game
	}

	// Check for [[bin]] section indicating CLI
	if strings.Contains(cargoContent, "[[bin]]") {
		if result.Type == types.ProjectTypeUnknown {
			result.Type = types.ProjectTypeCLI
		}
	}

	// Check for [lib] section indicating library
	if strings.Contains(cargoContent, "[lib]") && result.Type == types.ProjectTypeUnknown {
		result.Type = types.ProjectTypeLibrary
	}

	return result, nil
}

// parseCargoToml parses a Cargo.toml file and returns dependencies
func parseCargoToml(content string) map[string]bool {
	deps := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(content))
	inDeps := false
	inDevDeps := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Track sections
		if strings.HasPrefix(line, "[dependencies]") {
			inDeps = true
			inDevDeps = false
			continue
		}
		if strings.HasPrefix(line, "[dev-dependencies]") {
			inDeps = false
			inDevDeps = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDeps = false
			inDevDeps = false
			continue
		}

		// Parse dependencies
		if inDeps || inDevDeps {
			// Handle both formats:
			// crate_name = "version"
			// crate_name = { version = "x", features = [...] }
			if strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) >= 1 {
					dep := strings.TrimSpace(parts[0])
					// Normalize crate names (replace _ with -)
					dep = strings.ReplaceAll(dep, "_", "-")
					deps[dep] = true
				}
			}
		}
	}

	return deps
}
