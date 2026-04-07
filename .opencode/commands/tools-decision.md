---
description: Analyze project and recommend MCP servers using tools-decision
---

Analyze this project (or project idea) and recommend MCP servers. $ARGUMENTS

## Instructions

You are helping the user set up MCP (Model Context Protocol) servers for their project.

### If the user provided a project description:

The user said: "$ARGUMENTS"

Parse this description to identify:

**Technologies mentioned:**
- Languages: Go/Golang, Python, TypeScript/JavaScript, Rust, Java/Kotlin
- Frameworks: Gin, Gorm, Echo, Fiber (Go), FastAPI, Django, Flask (Python), Next.js, Express, NestJS (JS/TS), Axum, Actix (Rust), Spring Boot (Java)
- Databases: PostgreSQL, MySQL, MongoDB, Redis, SQLite
- Services: Docker, Kubernetes, AWS, Stripe, Auth/OAuth

**Architecture patterns:**
- "microservice", "micro-service" → Multiple services, likely needs docker, kubernetes
- "auth service" → Authentication/authorization, likely needs jwt, oauth tools
- "API", "REST", "GraphQL" → API development tools

### If no description provided:

Analyze the current project directory:
1. Check for go.mod, package.json, requirements.txt, Cargo.toml, pom.xml
2. Identify languages, frameworks, and dependencies
3. Look for docker-compose.yml, .env files, database configs

---

## Recommended MCP Servers by Stack

### Go Projects (Gin, Gorm, Echo, Fiber)
- `postgres` or `mysql` - Database operations
- `redis` - Caching (if mentioned)
- `docker` - Container management
- `github` - Version control
- `filesystem` - File operations

### For Microservices Architecture
- `docker` - Container management
- `kubernetes` - K8s operations (if using K8s)
- `postgres`/`mysql`/`mongodb` - Database per service
- `redis` - Shared cache/message queue

### For Auth Service
- `postgres` - User data storage
- `redis` - Session/token cache
- `github` - OAuth integration reference

---

## Output Format

1. **Summarize what you detected/inferred:**
```
## Project Analysis

**Type:** Microservices API
**Language:** Go
**Frameworks:** Gin, Gorm
**Services:** Auth service (first), more planned
**Databases:** PostgreSQL (inferred for user data)
**Infrastructure:** Docker, potentially Kubernetes
```

2. **Recommend MCP servers with reasons:**
```
## Recommended MCP Servers

1. **PostgreSQL** - Store user data, credentials, sessions
   - Needed for: Auth service user management
   
2. **Redis** - Token caching, rate limiting
   - Needed for: JWT token blacklist, session cache
   
3. **Docker** - Container management
   - Needed for: Microservices deployment
   
4. **GitHub** - Issue tracking, PR management
   - Needed for: Development workflow
   
5. **Filesystem** - File operations
   - Needed for: Config files, migrations
```

3. **Ask which format they want:**
- Claude Desktop
- Cursor  
- VS Code
- Generic (mcp.json)

4. **Generate the config** and write to the appropriate location

5. **List required environment variables:**
```
## Required Environment Variables

POSTGRES_CONNECTION_STRING=postgresql://user:pass@localhost:5432/auth
REDIS_URL=redis://localhost:6379
GITHUB_PERSONAL_ACCESS_TOKEN=ghp_xxxxx
```

---

## Example for user's request

For: "I want to create an API service with Golang, Gorm and Gin, about the micro-services with auth service first"

**Detected:**
- Language: Go
- Frameworks: Gin (HTTP), Gorm (ORM)
- Architecture: Microservices
- First service: Auth/Authentication
- Databases: PostgreSQL (implied by Gorm + auth)

**Recommended servers:**
1. postgres - User/auth data storage
2. redis - Session/token cache
3. docker - Container orchestration  
4. github - Code management
5. filesystem - File operations
