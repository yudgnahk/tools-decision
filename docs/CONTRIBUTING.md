# Contributing to tools-decision

Thanks for your interest in contributing! This project aims to make MCP server discovery effortless for developers.

## Ways to Contribute

### 1. Report Issues

Found a bug or have a feature request? [Open an issue](https://github.com/yudgnahk/tools-decision/issues).

**Good bug reports include:**
- What you expected to happen
- What actually happened
- Steps to reproduce
- Your environment (OS, Go version, project type)

### 2. Improve Detection

The analyzer detects languages, frameworks, and project types. Help us improve:

**Add a new language/framework detector:**

```go
// internal/analyzer/ruby.go
package analyzer

import (
    "context"
    "os"
    "path/filepath"
    "strings"

    "github.com/yudgnahk/tools-decision/pkg/types"
)

type RubyDetector struct{}

func NewRubyDetector() *RubyDetector {
    return &RubyDetector{}
}

func (d *RubyDetector) Name() string {
    return "ruby"
}

func (d *RubyDetector) Detect(_ context.Context, projectPath string) (*DetectorResult, error) {
    gemfilePath := filepath.Join(projectPath, "Gemfile")
    content, err := os.ReadFile(gemfilePath)
    if err != nil {
        return nil, nil
    }

    result := &DetectorResult{
        Languages: []types.Language{{Name: "ruby", Confidence: 0.95}},
    }

    if strings.Contains(string(content), "rails") {
        result.Frameworks = append(result.Frameworks, types.Framework{Name: "rails", Confidence: 0.9})
    }

    return result, nil
}
```

### 3. Add MCP Server Profiles

Help us better understand what each MCP server does:

```yaml
# data/profiles/postgres-mcp.yaml
id: postgres-mcp
name: PostgreSQL MCP Server
description: Database queries and schema inspection for PostgreSQL

# What this server is good for
use_cases:
  - database queries
  - schema inspection
  - migrations

# When to recommend this
triggers:
  dependencies:
    - pg
    - postgres
    - prisma  # with postgresql provider
  services:
    - postgresql
    - postgres

# Related servers
see_also:
  - prisma-mcp
  - mysql-mcp
```

### 4. Improve Matching

The matcher scores how relevant each server is. Ideas:
- Better weighting of factors
- New matching signals
- Edge case handling

### 5. Add Output Formats

Support more AI tools:

```go
// internal/config/generator.go
func (g *Generator) generateWindsurf(servers []types.MCPServer) (*types.ConfigOutput, error) {
    mcpServers := make(map[string]any)

    for _, server := range servers {
        mcpServers[server.Slug] = map[string]any{
            "command": server.Install.Command,
            "args":    server.Install.Args,
        }
    }

    content := map[string]any{
        "mcpServers": mcpServers,
    }

    return &types.ConfigOutput{
        Format:   "windsurf",
        Filename: "mcp.json",
        Path:     ".windsurf/mcp.json",
        Content:  content,
    }, nil
}
```

## Development Setup

### Prerequisites

- Go 1.22+
- Make (recommended)
- golangci-lint (optional, for linting)

### Getting Started

```bash
# Clone the repo
git clone https://github.com/yudgnahk/tools-decision.git
cd tools-decision

# Install dependencies
go mod download

# Run in development mode
make dev

# Run tests
make test

# Build
make build

# Test CLI locally
./tools-decision --help
```

### Project Structure

```
cmd/                   # CLI entrypoint
└── tools-decision/
internal/
├── analyzer/          # Project detection
├── registry/          # MCP registry fetching/caching
├── matcher/           # Recommendation algorithm
└── config/            # Config file generation
pkg/
└── types/             # Shared domain models
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test file
go test ./internal/analyzer -run TestIntent

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run in watch mode (requires watchexec)
watchexec -r go test ./...
```

### Testing Notes

- Prefer table-driven tests in `*_test.go` files.
- Keep tests deterministic (no network and no external services).
- Add targeted unit tests when adding detectors, matcher logic, or config formats.

## Pull Request Process

1. **Fork & branch** - Create a feature branch from `main`
2. **Make changes** - Keep commits focused and atomic
3. **Test** - Ensure all tests pass, add new tests for new features
4. **Document** - Update README if adding user-facing features
5. **PR** - Open a pull request with a clear description

### Commit Messages

Follow conventional commits:

```
feat: add Ruby/Rails detector
fix: handle missing go.mod gracefully
docs: add installation instructions for homebrew
refactor: simplify matching algorithm
test: add analyzer coverage for mixed-language repos
```

### Code Style

- Use `gofmt` for formatting
- Keep package boundaries clear (`cmd`, `internal`, `pkg`)
- Prefer small, focused functions and explicit error handling
- Run `make lint` before committing

## Architecture Decisions

### Why Go?
- Fast startup and great performance for CLI tools
- Single static binaries simplify installation and distribution
- Strong standard library and explicit error handling improve reliability

### Why local-first?
- Privacy: Your code never leaves your machine
- Speed: No network latency for analysis
- Reliability: Works offline

### Why cache registries locally?
- Speed: Instant recommendations
- Offline support: Works without internet
- Rate limits: Avoid API throttling

## Questions?

- Open a [Discussion](https://github.com/yudgnahk/tools-decision/discussions)
- Check existing [Issues](https://github.com/yudgnahk/tools-decision/issues)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
