# MCP Setup Command for OpenCode

This directory contains a custom OpenCode slash command for MCP server discovery and configuration.

## Installation

Copy the `mcp-setup.md` file to your OpenCode commands directory:

```bash
# Global installation (available in all projects)
mkdir -p ~/.config/opencode/commands
cp mcp-setup.md ~/.config/opencode/commands/

# Or per-project installation
mkdir -p .opencode/commands
cp mcp-setup.md .opencode/commands/
```

## Usage

In OpenCode, run:

```
/mcp-setup
```

The command will:
1. Analyze your project structure
2. Recommend relevant MCP servers
3. Generate configuration for your preferred AI tool
4. Write the configuration file

## Customization

Edit the `mcp-setup.md` file to customize:
- The analysis criteria
- Which MCP servers to consider
- Output format preferences
- Agent or model to use

## Example

```
$ opencode
> /mcp-setup

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
