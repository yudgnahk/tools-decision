Analyze this project (or project idea) and recommend MCP servers. $ARGUMENTS

If `$ARGUMENTS` is present, treat it as a NEW project idea and infer the stack.
If `$ARGUMENTS` is empty, inspect the existing repository files to detect the stack.

For ideas like:
`I want to create a API service with Golang, Gorm and Gin, about the micro-services with auth service first`

infer:
- Language: Go
- Frameworks: Gin, Gorm
- Architecture: Microservices
- First service focus: authentication
- Likely services: PostgreSQL, Redis, Docker

Recommend MCP servers with concise reasons, then generate config in the user's preferred format:
- Claude Desktop
- Cursor
- VS Code
- Generic (mcp.json)

Always include required environment variables for selected servers.
