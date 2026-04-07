# Contributing to tools-decision

Thanks for your interest in contributing! This project aims to make MCP server discovery effortless for developers.

## Ways to Contribute

### 1. Report Issues

Found a bug or have a feature request? [Open an issue](https://github.com/yudgnahk/tools-decision/issues).

**Good bug reports include:**
- What you expected to happen
- What actually happened
- Steps to reproduce
- Your environment (OS, Node version, project type)

### 2. Improve Detection

The analyzer detects languages, frameworks, and project types. Help us improve:

**Add a new language/framework detector:**

```typescript
// src/analyzer/detectors/ruby.ts
import { Detector, ProjectContext } from '../types';

export const rubyDetector: Detector = {
  name: 'ruby',
  
  async detect(projectPath: string): Promise<Partial<ProjectContext>> {
    const gemfile = await readFile(join(projectPath, 'Gemfile'));
    if (!gemfile) return {};
    
    const context: Partial<ProjectContext> = {
      languages: [{ name: 'ruby', confidence: 0.95 }],
      frameworks: [],
    };
    
    // Detect Rails
    if (gemfile.includes('rails')) {
      context.frameworks.push({ name: 'rails', confidence: 0.9 });
    }
    
    return context;
  }
};
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

```typescript
// src/config/formats/windsurf.ts
import { ConfigFormat } from '../types';

export const windsurfFormat: ConfigFormat = {
  name: 'windsurf',
  filename: '.windsurf/mcp.json',
  
  generate(servers: MCPServer[]): object {
    return {
      version: 1,
      servers: servers.map(s => ({
        name: s.slug,
        command: s.installation.command,
        args: s.installation.args,
      }))
    };
  }
};
```

## Development Setup

### Prerequisites

- Node.js 18+
- pnpm (recommended) or npm

### Getting Started

```bash
# Clone the repo
git clone https://github.com/yudgnahk/tools-decision.git
cd tools-decision

# Install dependencies
pnpm install

# Run in development mode
pnpm dev

# Run tests
pnpm test

# Build
pnpm build

# Test CLI locally
pnpm link
tools-decision --help
```

### Project Structure

```
src/
├── cli/          # Command-line interface
├── analyzer/     # Project detection
├── registry/     # MCP registry fetching/caching
├── matcher/      # Recommendation algorithm
└── config/       # Config file generation
```

### Running Tests

```bash
# Run all tests
pnpm test

# Run specific test file
pnpm test analyzer

# Run with coverage
pnpm test:coverage

# Run in watch mode
pnpm test:watch
```

### Test Fixtures

We have sample projects for testing detection:

```
tests/fixtures/
├── nextjs-app/        # Next.js + TypeScript
├── fastapi-app/       # Python FastAPI
├── go-api/            # Go with Gin
└── monorepo/          # Multi-language monorepo
```

Add new fixtures to test edge cases.

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
fix: handle missing package.json gracefully
docs: add installation instructions for homebrew
refactor: simplify matching algorithm
test: add fixtures for Python monorepos
```

### Code Style

- TypeScript strict mode
- Prettier for formatting
- ESLint for linting
- Run `pnpm lint` before committing

## Architecture Decisions

### Why Node.js/TypeScript?
- Most developers have Node.js installed
- TypeScript provides type safety
- Easy distribution via npm/npx

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
