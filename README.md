# tools-decision

Intelligent MCP server discovery and configuration for your projects.

Analyzes your project and recommends the best MCP (Model Context Protocol) servers, then generates ready-to-use configuration for Claude, Cursor, VS Code, and other AI coding tools.

## Why?

There are **21,000+ MCP servers** available. Finding the right ones for your project means:
- Searching multiple registries manually
- Reading documentation for each server
- Figuring out configuration formats
- Copy-pasting and hoping it works

**tools-decision** does this in seconds:

```bash
$ tools-decision

Analyzing project...

Detected:
  Language:   TypeScript
  Framework:  Next.js 14, React 18
  Database:   PostgreSQL (Prisma)
  Testing:    Jest

Recommended MCP Servers:
  1. postgres-mcp        Database queries and schema inspection
  2. prisma-mcp          Prisma schema and migrations
  3. nextjs-mcp          Next.js app router utilities
  4. jest-runner-mcp     Test execution and coverage

Install these tools? [Y/n]
```

## Installation

### Binary (Recommended)

Download the latest binary for your platform from [Releases](https://github.com/yudgnahk/tools-decision/releases).

```bash
# macOS (Apple Silicon)
curl -L https://github.com/yudgnahk/tools-decision/releases/latest/download/tools-decision-darwin-arm64 -o tools-decision
chmod +x tools-decision
sudo mv tools-decision /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/yudgnahk/tools-decision/releases/latest/download/tools-decision-darwin-amd64 -o tools-decision
chmod +x tools-decision
sudo mv tools-decision /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/yudgnahk/tools-decision/releases/latest/download/tools-decision-linux-amd64 -o tools-decision
chmod +x tools-decision
sudo mv tools-decision /usr/local/bin/
```

### Homebrew (macOS/Linux)

```bash
brew install yudgnahk/tap/tools-decision
```

### Go Install

```bash
go install github.com/yudgnahk/tools-decision@latest
```

### Build from Source

```bash
git clone https://github.com/yudgnahk/tools-decision.git
cd tools-decision
go build -o tools-decision ./cmd/tools-decision
```

## Usage

### As a CLI

Run in your project directory:

```bash
# Analyze and get recommendations
tools-decision

# Search for specific tools
tools-decision search "database"

# Generate config for a specific IDE
tools-decision config --format cursor
tools-decision config --format claude
tools-decision config --format vscode
```

### As a Slash Command (OpenCode)

Copy the command file to your OpenCode commands directory:

```bash
# Global (available in all projects)
mkdir -p ~/.config/opencode/commands
cp .opencode/commands/mcp-setup.md ~/.config/opencode/commands/
cp .opencode/commands/tools-decision.md ~/.config/opencode/commands/

# Or per-project
mkdir -p .opencode/commands
cp .opencode/commands/mcp-setup.md .opencode/commands/
cp .opencode/commands/tools-decision.md .opencode/commands/
```

Then use it in OpenCode:

```
/tools-decision

# Legacy alias (still supported)
/mcp-setup
```

You can also pass a project idea directly:

```
/tools-decision I want to create a API service with Golang, Gorm and Gin, about the micro-services with auth service first
```

The command will analyze your project or idea, recommend servers, and generate configuration.

### As a Slash Command (Claude Code)

Copy the command file to your Claude Code commands directory:

```bash
# Global
mkdir -p ~/.claude/commands
cp .claude/commands/tools-decision.md ~/.claude/commands/

# Or per-project
mkdir -p .claude/commands
cp .claude/commands/tools-decision.md .claude/commands/
```

Then use it:

```
/tools-decision I want to create a API service with Golang, Gorm and Gin, about the micro-services with auth service first
```

## How It Works

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Your Project   │────▶│  tools-decision  │────▶│  MCP Config     │
│                 │     │                  │     │                 │
│ - package.json  │     │ 1. Analyze       │     │ claude_desktop_ │
│ - tsconfig.json │     │ 2. Match tools   │     │ config.json     │
│ - .env          │     │ 3. Score quality │     │                 │
│ - etc.          │     │ 4. Generate cfg  │     │ .cursor/mcp.json│
└─────────────────┘     └──────────────────┘     └─────────────────┘
```

1. **Analyze** - Detects languages, frameworks, dependencies, project type
2. **Match** - Queries MCP registries for relevant servers
3. **Score** - Ranks by quality (maintenance, popularity, security)
4. **Configure** - Generates config for your preferred AI tool

## Supported Project Types

| Language | Frameworks | Detection |
|----------|------------|-----------|
| TypeScript/JavaScript | React, Next.js, Vue, Express, Nest.js | package.json, tsconfig |
| Python | FastAPI, Django, Flask | requirements.txt, pyproject.toml |
| Go | Gin, Echo, Fiber | go.mod |
| Rust | Actix, Axum, Rocket | Cargo.toml |
| Java | Spring Boot, Quarkus | pom.xml, build.gradle |

## Supported Output Formats

| Format | File | Used By |
|--------|------|---------|
| `claude` | `claude_desktop_config.json` | Claude Desktop |
| `cursor` | `.cursor/mcp.json` | Cursor IDE |
| `vscode` | `.vscode/mcp.json` | VS Code + MCP extension |
| `generic` | `mcp.json` | Generic MCP config |

## Configuration

Create `~/.config/tools-decision/config.yaml` for defaults:

```yaml
# Default output format
format: claude

# Minimum quality score (0-1)
min_quality: 0.7

# Preferred categories
prefer:
  - database
  - testing
  - documentation

# Tools to always exclude
exclude:
  - deprecated-server
```

## Data Sources

tools-decision aggregates from multiple MCP registries:

- [Official MCP Registry](https://github.com/modelcontextprotocol/servers)
- [Smithery](https://smithery.ai)
- [Glama](https://glama.ai/mcp/servers)
- [mcp.so](https://mcp.so)

All data is fetched and cached locally. No account required. No data sent anywhere.

## Privacy

- **Fully local** - Your code never leaves your machine
- **No telemetry** - We don't track anything
- **No account** - Just install and use
- **Offline capable** - Works with cached registry data

## Contributing

Contributions welcome! See [CONTRIBUTING.md](docs/CONTRIBUTING.md).

### Ideas for Contributions

- Add detection for more frameworks/languages
- Improve tool matching algorithm
- Add more output formats
- Better quality scoring heuristics
- Tool profile corrections

## License

MIT License - Use it however you want.

## Credits

Built for the MCP ecosystem. Thanks to:
- [Anthropic](https://anthropic.com) for creating MCP
- [Smithery](https://smithery.ai), [Glama](https://glama.ai) for registry APIs
- The MCP server authors who build amazing tools
