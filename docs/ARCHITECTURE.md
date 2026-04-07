# Architecture

This document describes the technical architecture of tools-decision, a local-first CLI tool for MCP server discovery and configuration.

## Design Principles

1. **Local-first** - Everything runs on your machine, no cloud dependency
2. **Fast** - Analysis completes in seconds, not minutes
3. **Offline-capable** - Works with cached data when offline
4. **Zero config** - Works out of the box with sensible defaults
5. **Composable** - Can be used standalone or integrated into other tools

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        tools-decision                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   CLI       │  │  Analyzer   │  │  Matcher    │             │
│  │   Layer     │──│  Engine     │──│  Engine     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│         │                                   │                   │
│         │         ┌─────────────┐          │                   │
│         │         │  Registry   │◀─────────┘                   │
│         │         │  Cache      │                              │
│         │         └─────────────┘                              │
│         │                │                                      │
│         ▼                ▼                                      │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │  Config     │  │  Registry   │                              │
│  │  Generator  │  │  Fetcher    │                              │
│  └─────────────┘  └─────────────┘                              │
│                          │                                      │
└──────────────────────────│──────────────────────────────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │   External Registries   │
              │  (Smithery, Glama, etc) │
              └─────────────────────────┘
```

## Components

### 1. CLI Layer

Entry point for the tool. Handles:
- Command parsing (`analyze`, `search`, `config`, `install`)
- Interactive prompts
- Output formatting (table, JSON, plain text)
- Progress indicators

```
tools-decision [command] [options]

Commands:
  (default)     Analyze current project and recommend tools
  search        Search for MCP servers by keyword
  config        Generate configuration file
  install       Install recommended servers
  update        Update registry cache
  list          List installed MCP servers

Options:
  --format      Output format (claude, cursor, vscode, generic)
  --json        Output as JSON
  --quiet       Minimal output
  --no-cache    Force fresh registry fetch
```

### 2. Analyzer Engine

Scans the project directory to build a context profile.

**Inputs analyzed:**
- `package.json` - Node.js dependencies, scripts
- `tsconfig.json` / `jsconfig.json` - TypeScript/JavaScript config
- `requirements.txt` / `pyproject.toml` - Python dependencies
- `go.mod` - Go modules
- `Cargo.toml` - Rust crates
- `pom.xml` / `build.gradle` - Java dependencies
- `.env` / `.env.example` - Environment variables (names only, not values)
- `docker-compose.yml` - Services used
- `Makefile` / `Taskfile` - Build tasks
- Directory structure - `src/`, `tests/`, `docs/`, etc.

**Output: Project Context**
```typescript
interface ProjectContext {
  languages: Language[];           // typescript, python, go, etc.
  frameworks: Framework[];         // nextjs, fastapi, gin, etc.
  dependencies: Dependency[];      // parsed from package managers
  projectType: ProjectType;        // web_app, api, cli, library, etc.
  detectedServices: Service[];     // postgres, redis, s3, etc.
  testingFramework?: string;       // jest, pytest, go test, etc.
  buildTool?: string;              // webpack, vite, esbuild, etc.
  features: string[];              // auth, payments, file_upload, etc.
}
```

### 3. Registry Cache

Local cache of MCP server metadata from multiple registries.

**Cache location:** `~/.cache/tools-decision/`

```
~/.cache/tools-decision/
├── registries/
│   ├── official.json      # Official MCP registry
│   ├── smithery.json      # Smithery.ai
│   └── glama.json         # Glama.ai
├── tools/
│   └── index.json         # Merged, deduplicated tool index
└── meta.json              # Cache timestamps, versions
```

**Cache behavior:**
- Auto-refresh if older than 24 hours
- Manual refresh with `tools-decision update`
- Offline mode uses stale cache with warning
- Cache size: ~5-10MB typical

### 4. Registry Fetcher

Fetches and normalizes data from multiple MCP registries.

**Supported registries:**

| Registry | API | Data |
|----------|-----|------|
| Official MCP | GitHub API | servers.json from repo |
| Smithery | REST API | Server listings |
| Glama | REST API | Indexed servers |
| mcp.so | REST API | Community servers |

**Normalization:**
All registries are normalized to a common schema:

```typescript
interface MCPServer {
  id: string;                    // Unique identifier
  name: string;                  // Display name
  slug: string;                  // URL-safe name
  description: string;           // Short description
  author: string;                // Author/org name
  repository?: string;           // GitHub URL
  npm?: string;                  // npm package name
  pypi?: string;                 // PyPI package name
  
  categories: string[];          // database, testing, docs, etc.
  tags: string[];                // More specific tags
  
  capabilities: {
    tools: string[];             // Tool names exposed
    resources: string[];         // Resource types
    prompts: string[];           // Prompt templates
  };
  
  compatibility: {
    languages: string[];         // Works with these languages
    frameworks: string[];        // Optimized for these frameworks
  };
  
  quality: {
    stars?: number;              // GitHub stars
    downloads?: number;          // Package downloads
    lastUpdate?: string;         // Last commit/publish date
    maintained: boolean;         // Active maintenance
  };
  
  installation: {
    command: string;             // npx, pip, etc.
    args: string[];              // Package name, flags
  };
  
  source: string;                // Which registry this came from
}
```

### 5. Matcher Engine

Matches project context to relevant MCP servers.

**Matching algorithm:**

```
Score = Σ (weight × match_score)

Factors:
  - Language match:     0.25 weight
  - Framework match:    0.25 weight
  - Category relevance: 0.20 weight
  - Service match:      0.15 weight (e.g., postgres → postgres-mcp)
  - Quality score:      0.15 weight
```

**Quality scoring:**

```typescript
function qualityScore(server: MCPServer): number {
  let score = 0.5; // Base score
  
  // Maintenance (up to +0.2)
  if (server.quality.lastUpdate > sixMonthsAgo) score += 0.1;
  if (server.quality.lastUpdate > oneMonthAgo) score += 0.1;
  
  // Popularity (up to +0.2)
  if (server.quality.stars > 100) score += 0.1;
  if (server.quality.stars > 500) score += 0.1;
  
  // Package health (up to +0.1)
  if (server.quality.downloads > 1000) score += 0.05;
  if (server.quality.maintained) score += 0.05;
  
  return Math.min(score, 1.0);
}
```

**Output: Ranked recommendations**

```typescript
interface Recommendation {
  server: MCPServer;
  score: number;              // 0-1 overall match score
  reasons: string[];          // Why this was recommended
  matchedOn: string[];        // What project features matched
}
```

### 6. Config Generator

Generates configuration files for various AI tools.

**Supported formats:**

```typescript
type ConfigFormat = 'claude' | 'cursor' | 'vscode' | 'generic';

interface GeneratedConfig {
  format: ConfigFormat;
  filename: string;
  path: string;               // Where to write the file
  content: object;            // The configuration object
  envVars?: EnvVar[];         // Required environment variables
}
```

**Example outputs:**

Claude Desktop:
```json
{
  "mcpServers": {
    "postgres-mcp": {
      "command": "npx",
      "args": ["-y", "@mcp/postgres-server"],
      "env": {
        "DATABASE_URL": "${DATABASE_URL}"
      }
    }
  }
}
```

Cursor:
```json
{
  "mcpServers": {
    "postgres-mcp": {
      "command": "npx",
      "args": ["-y", "@mcp/postgres-server"],
      "env": {
        "DATABASE_URL": "${DATABASE_URL}"
      }
    }
  }
}
```

## Data Flow

### Typical Flow: Project Analysis

```
User runs: tools-decision

1. CLI parses arguments, loads config
2. Check cache freshness
   └─ If stale: Registry Fetcher updates cache
3. Analyzer Engine scans project directory
   └─ Returns: ProjectContext
4. Matcher Engine queries cache with context
   └─ Returns: Recommendation[]
5. CLI displays recommendations interactively
6. User selects tools to install
7. Config Generator creates config file
8. CLI writes config to appropriate location
```

### Typical Flow: Slash Command (OpenCode)

```
User types: /mcp-setup

1. OpenCode invokes: npx tools-decision --json
2. tools-decision analyzes project (steps 2-4 above)
3. Returns JSON with recommendations
4. OpenCode displays recommendations to user
5. User confirms selection
6. OpenCode invokes: npx tools-decision config --format claude --tools postgres-mcp,jest-mcp
7. Config file is generated and written
```

## File Structure

```
tools-decision/
├── src/
│   ├── cli/
│   │   ├── index.ts           # Entry point
│   │   ├── commands/
│   │   │   ├── analyze.ts     # Default command
│   │   │   ├── search.ts      # Search command
│   │   │   ├── config.ts      # Config generation
│   │   │   └── install.ts     # Install helper
│   │   └── ui/
│   │       ├── prompts.ts     # Interactive prompts
│   │       └── display.ts     # Output formatting
│   │
│   ├── analyzer/
│   │   ├── index.ts           # Analyzer orchestrator
│   │   ├── detectors/
│   │   │   ├── javascript.ts  # JS/TS detection
│   │   │   ├── python.ts      # Python detection
│   │   │   ├── go.ts          # Go detection
│   │   │   └── ...
│   │   └── types.ts           # ProjectContext types
│   │
│   ├── registry/
│   │   ├── cache.ts           # Cache management
│   │   ├── fetcher.ts         # Registry fetching
│   │   ├── sources/
│   │   │   ├── official.ts    # Official MCP registry
│   │   │   ├── smithery.ts    # Smithery API
│   │   │   └── glama.ts       # Glama API
│   │   └── types.ts           # MCPServer types
│   │
│   ├── matcher/
│   │   ├── index.ts           # Matching algorithm
│   │   ├── scoring.ts         # Quality scoring
│   │   └── types.ts           # Recommendation types
│   │
│   └── config/
│       ├── generator.ts       # Config generation
│       └── formats/
│           ├── claude.ts      # Claude Desktop format
│           ├── cursor.ts      # Cursor format
│           └── vscode.ts      # VS Code format
│
├── tests/
│   ├── analyzer/
│   ├── matcher/
│   └── fixtures/              # Sample projects for testing
│
├── package.json
├── tsconfig.json
└── README.md
```

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | TypeScript | Type safety, ecosystem |
| Runtime | Node.js | Ubiquitous, easy install |
| CLI Framework | Commander.js | Simple, well-documented |
| Prompts | Inquirer.js | Interactive CLI |
| HTTP | undici/fetch | Built into Node 18+ |
| File parsing | Various | JSON, YAML, TOML parsers |
| Testing | Vitest | Fast, modern |
| Build | tsup | Simple bundling |
| Package | npm | Standard distribution |

## Performance Targets

| Operation | Target | Notes |
|-----------|--------|-------|
| Cold start | <500ms | First run, no cache |
| Warm analysis | <200ms | With cached registry |
| Registry update | <5s | Full refresh |
| Config generation | <50ms | After selection |

## Future Considerations

### Potential Enhancements

1. **Plugin system** - Custom detectors for niche frameworks
2. **Config profiles** - Save/share configurations
3. **Update checker** - Notify when better tools available
4. **Dry run mode** - Preview without writing files
5. **Uninstall command** - Clean removal of MCP configs

### Non-Goals (Out of Scope)

- Cloud sync / accounts
- Team features
- Telemetry / analytics
- Paid features
- GUI application
