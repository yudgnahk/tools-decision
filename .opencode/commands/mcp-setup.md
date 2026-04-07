---
description: Recommend and configure MCP servers for this project (legacy alias)
agent: build
---

Set up MCP (Model Context Protocol) servers for this project. $ARGUMENTS

Note: `/mcp-setup` is a legacy alias. Prefer `/tools-decision` for new usage.

## Determine Project State

First, check if this is a NEW project or an EXISTING project:

**If files exist** (package.json, go.mod, requirements.txt, etc.):
- This is an EXISTING project
- Analyze the actual files to understand the stack

**If no project files exist** (or user provided a description):
- This is a NEW project idea
- Use the description "$ARGUMENTS" to infer what tools are needed

---

## For EXISTING Projects

1. Examine the project structure:
   - Check for package.json, go.mod, requirements.txt, Cargo.toml, pyproject.toml
   - Identify languages, frameworks, and dependencies
   - Look for database configs, docker-compose, .env files

2. Based on detected stack, recommend relevant MCP servers.

---

## For NEW Projects (from idea/description)

Parse the project idea and identify:

### Languages (look for keywords)
- "typescript", "ts", "javascript", "node" → TypeScript/JavaScript
- "python", "django", "flask", "fastapi" → Python
- "go", "golang", "gin", "fiber" → Go
- "rust" → Rust

### Frameworks
- "next", "nextjs" → Next.js
- "react" → React
- "vue", "nuxt" → Vue/Nuxt
- "express", "nestjs" → Node.js backends
- "fastapi", "django", "flask" → Python backends
- "gin", "echo", "fiber" → Go backends

### Services/Features
- "postgres", "postgresql", "database" → postgres-mcp
- "mysql" → mysql-mcp
- "mongodb", "mongo" → mongodb-mcp
- "redis", "cache" → redis-mcp
- "docker", "container" → docker-mcp
- "kubernetes", "k8s" → kubernetes-mcp
- "aws", "s3" → aws-mcp, s3-mcp
- "stripe", "payment" → stripe-mcp
- "auth", "authentication" → auth tools
- "graphql" → graphql-mcp
- "openai", "ai", "llm" → openai-mcp

---

## MCP Server Registry

Here are the most common MCP servers to recommend:

### Databases
- `postgres-mcp` - PostgreSQL queries, schema inspection, migrations
- `mysql-mcp` - MySQL database operations
- `mongodb-mcp` - MongoDB document operations
- `sqlite-mcp` - SQLite database operations
- `redis-mcp` - Redis cache operations
- `prisma-mcp` - Prisma ORM schema and migrations

### Infrastructure
- `docker-mcp` - Docker container management
- `kubernetes-mcp` - Kubernetes cluster operations
- `aws-mcp` - AWS services integration
- `s3-mcp` - S3 file storage operations
- `vercel-mcp` - Vercel deployment tools
- `cloudflare-mcp` - Cloudflare services

### Development
- `filesystem` - File system operations (reading, writing, searching)
- `git` - Git version control operations
- `github` - GitHub API (issues, PRs, repos)
- `gitlab` - GitLab API operations

### Frameworks
- `nextjs-mcp` - Next.js specific tools
- `fastapi-mcp` - FastAPI development tools
- `django-mcp` - Django development tools

### External Services
- `stripe-mcp` - Stripe payment processing
- `openai-mcp` - OpenAI API integration
- `anthropic-mcp` - Anthropic Claude API
- `slack-mcp` - Slack messaging integration
- `email-mcp` - Email sending (SendGrid, SES)

### Testing
- `jest-mcp` - Jest test runner
- `pytest-mcp` - Pytest test runner

---

## Output

1. **List the recommended MCP servers** with brief explanations of why each is useful
2. **Ask the user** which config format they want:
   - Claude Desktop (writes to ~/Library/Application Support/Claude/claude_desktop_config.json)
   - Cursor (writes to .cursor/mcp.json)
   - VS Code (writes to .vscode/settings.json)
   - Generic (writes to mcp.json)
3. **Generate the configuration** in the selected format
4. **Write the file** to the appropriate location
5. **List any environment variables** that need to be set (DATABASE_URL, API keys, etc.)

---

## Example Output Format

```
## Detected/Inferred Stack
- Language: TypeScript
- Framework: Next.js 14
- Database: PostgreSQL
- Services: Stripe, Redis

## Recommended MCP Servers

1. **postgres-mcp** - Database queries and schema management
2. **redis-mcp** - Cache operations
3. **stripe-mcp** - Payment processing
4. **github** - Version control integration
5. **filesystem** - File operations

Which config format would you like?
1. Claude Desktop
2. Cursor
3. VS Code
4. Generic

[After user selects, generate and write the config]
```
