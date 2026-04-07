# MCP Setup Command for OpenCode

This directory contains custom OpenCode slash commands for MCP server discovery and configuration.

## Installation

Copy the command files to your OpenCode commands directory:

```bash
# Global installation (available in all projects)
mkdir -p ~/.config/opencode/commands
cp mcp-setup.md ~/.config/opencode/commands/
cp tools-decision.md ~/.config/opencode/commands/

# Or per-project installation
mkdir -p .opencode/commands
cp mcp-setup.md .opencode/commands/
cp tools-decision.md .opencode/commands/
```

## Usage

In OpenCode, run:

```
/tools-decision

# Legacy alias (still supported)
/mcp-setup
```

The command will:
1. Analyze your project structure
2. Recommend relevant MCP servers
3. Generate configuration for your preferred AI tool
4. Write the configuration file

You can also pass a new-project idea directly:

```
/tools-decision I want to create an API service with Golang, Gorm and Gin, about the micro-services with auth service first
```

## Customization

Edit `tools-decision.md` (or `mcp-setup.md`) to customize:
- The analysis criteria
- Which MCP servers to consider
- Output format preferences
- Agent or model to use

## Example

```
$ opencode
> /tools-decision I want to build a Next.js app with PostgreSQL and Redis

Analyzing project...

Detected:
  Language:   TypeScript
  Framework:  Next.js 14
  Database:   PostgreSQL (Prisma)
  Testing:    Jest

Recommended MCP Servers:
  1. postgres-mcp     - Database queries and schema inspection
  2. prisma-mcp       - Prisma schema management
  3. github-mcp       - GitHub integration
  4. jest-mcp         - Test execution

Which format? [claude/cursor/vscode/generic]: cursor

Configuration written to .cursor/mcp.json
```
