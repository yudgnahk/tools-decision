package registry

import (
	"github.com/yudgnahk/tools-decision/pkg/types"
)

// GetEmbeddedServers returns the bundled MCP server registry
// This is a curated list of popular, well-maintained MCP servers
func GetEmbeddedServers() []types.MCPServer {
	return []types.MCPServer{
		// === Core Development Tools ===
		{
			ID:          "filesystem",
			Name:        "Filesystem",
			Slug:        "filesystem",
			Description: "Read, write, and manage files and directories",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"core", "filesystem"},
			Tags:        []string{"files", "directories", "read", "write"},
			Capabilities: types.Capabilities{
				Tools: []string{"read_file", "write_file", "list_directory", "create_directory", "move_file", "search_files"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      1000,
				Maintained: true,
				Score:      0.95,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"},
			},
			Source: "official",
		},
		{
			ID:          "git",
			Name:        "Git",
			Slug:        "git",
			Description: "Git version control operations - status, diff, commit, branch management",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"core", "git", "vcs"},
			Tags:        []string{"git", "version-control", "commit", "branch"},
			Capabilities: types.Capabilities{
				Tools: []string{"git_status", "git_diff", "git_commit", "git_log", "git_branch"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      800,
				Maintained: true,
				Score:      0.92,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-git"},
			},
			Source: "official",
		},
		{
			ID:          "github",
			Name:        "GitHub",
			Slug:        "github",
			Description: "GitHub API integration - issues, PRs, repos, actions",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"core", "github", "vcs"},
			Tags:        []string{"github", "issues", "pull-requests", "repos"},
			Capabilities: types.Capabilities{
				Tools: []string{"create_issue", "list_issues", "create_pr", "list_prs", "get_repo"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      900,
				Maintained: true,
				Score:      0.93,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-github"},
				Env: []types.EnvVar{
					{Name: "GITHUB_PERSONAL_ACCESS_TOKEN", Description: "GitHub personal access token", Required: true},
				},
			},
			Source: "official",
		},

		// === Databases ===
		{
			ID:          "postgres",
			Name:        "PostgreSQL",
			Slug:        "postgres",
			Description: "PostgreSQL database operations - queries, schema inspection, migrations",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"database", "postgresql", "sql"},
			Tags:        []string{"postgres", "postgresql", "database", "sql", "queries"},
			Capabilities: types.Capabilities{
				Tools:     []string{"query", "list_tables", "describe_table", "list_schemas"},
				Resources: []string{"schema"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      600,
				Maintained: true,
				Score:      0.90,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-postgres"},
				Env: []types.EnvVar{
					{Name: "POSTGRES_CONNECTION_STRING", Description: "PostgreSQL connection string", Required: true},
				},
			},
			Source: "official",
		},
		{
			ID:          "sqlite",
			Name:        "SQLite",
			Slug:        "sqlite",
			Description: "SQLite database operations for local databases",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"database", "sqlite", "sql"},
			Tags:        []string{"sqlite", "database", "sql", "local"},
			Capabilities: types.Capabilities{
				Tools: []string{"query", "list_tables", "describe_table"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      400,
				Maintained: true,
				Score:      0.85,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-sqlite", "--db-path", "./database.db"},
			},
			Source: "official",
		},
		{
			ID:          "mysql",
			Name:        "MySQL",
			Slug:        "mysql",
			Description: "MySQL/MariaDB database operations",
			Author:      "community",
			Repository:  "https://github.com/benborber/mcp-server-mysql",
			Categories:  []string{"database", "mysql", "sql"},
			Tags:        []string{"mysql", "mariadb", "database", "sql"},
			Capabilities: types.Capabilities{
				Tools: []string{"query", "list_tables", "describe_table"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      200,
				Maintained: true,
				Score:      0.80,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-mysql"},
				Env: []types.EnvVar{
					{Name: "MYSQL_HOST", Description: "MySQL host", Required: true},
					{Name: "MYSQL_USER", Description: "MySQL user", Required: true},
					{Name: "MYSQL_PASSWORD", Description: "MySQL password", Required: true},
					{Name: "MYSQL_DATABASE", Description: "MySQL database", Required: true},
				},
			},
			Source: "community",
		},
		{
			ID:          "mongodb",
			Name:        "MongoDB",
			Slug:        "mongodb",
			Description: "MongoDB document database operations",
			Author:      "community",
			Repository:  "https://github.com/kiliczsh/mcp-mongo-server",
			Categories:  []string{"database", "mongodb", "nosql"},
			Tags:        []string{"mongodb", "mongo", "nosql", "documents"},
			Capabilities: types.Capabilities{
				Tools: []string{"find", "insert", "update", "delete", "aggregate"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      150,
				Maintained: true,
				Score:      0.78,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-mongo-server"},
				Env: []types.EnvVar{
					{Name: "MONGODB_URI", Description: "MongoDB connection URI", Required: true},
				},
			},
			Source: "community",
		},
		{
			ID:          "redis",
			Name:        "Redis",
			Slug:        "redis",
			Description: "Redis cache and data structure operations",
			Author:      "community",
			Repository:  "https://github.com/gongrzhe/mcp-server-redis",
			Categories:  []string{"database", "redis", "cache"},
			Tags:        []string{"redis", "cache", "key-value"},
			Capabilities: types.Capabilities{
				Tools: []string{"get", "set", "delete", "keys", "hget", "hset"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      100,
				Maintained: true,
				Score:      0.75,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-redis"},
				Env: []types.EnvVar{
					{Name: "REDIS_URL", Description: "Redis connection URL", Required: true, Default: "redis://localhost:6379"},
				},
			},
			Source: "community",
		},

		// === Infrastructure ===
		{
			ID:          "docker",
			Name:        "Docker",
			Slug:        "docker",
			Description: "Docker container management - build, run, manage containers",
			Author:      "community",
			Repository:  "https://github.com/ckreiling/mcp-server-docker",
			Categories:  []string{"infrastructure", "docker", "containers"},
			Tags:        []string{"docker", "containers", "images", "compose"},
			Capabilities: types.Capabilities{
				Tools: []string{"list_containers", "run_container", "stop_container", "build_image", "list_images"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      300,
				Maintained: true,
				Score:      0.82,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-docker"},
			},
			Source: "community",
		},
		{
			ID:          "kubernetes",
			Name:        "Kubernetes",
			Slug:        "kubernetes",
			Description: "Kubernetes cluster management and operations",
			Author:      "community",
			Repository:  "https://github.com/Flux159/mcp-server-kubernetes",
			Categories:  []string{"infrastructure", "kubernetes", "k8s"},
			Tags:        []string{"kubernetes", "k8s", "pods", "deployments", "services"},
			Capabilities: types.Capabilities{
				Tools: []string{"list_pods", "list_deployments", "apply_manifest", "get_logs", "describe_resource"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      250,
				Maintained: true,
				Score:      0.80,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-kubernetes"},
			},
			Source: "community",
		},
		{
			ID:          "aws",
			Name:        "AWS",
			Slug:        "aws",
			Description: "AWS cloud services integration",
			Author:      "community",
			Repository:  "https://github.com/rishikavikondala/mcp-server-aws",
			Categories:  []string{"infrastructure", "aws", "cloud"},
			Tags:        []string{"aws", "s3", "ec2", "lambda", "cloud"},
			Capabilities: types.Capabilities{
				Tools: []string{"s3_list", "s3_get", "s3_put", "ec2_list", "lambda_invoke"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      180,
				Maintained: true,
				Score:      0.78,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-aws"},
				Env: []types.EnvVar{
					{Name: "AWS_ACCESS_KEY_ID", Description: "AWS access key", Required: true},
					{Name: "AWS_SECRET_ACCESS_KEY", Description: "AWS secret key", Required: true},
					{Name: "AWS_REGION", Description: "AWS region", Required: false, Default: "us-east-1"},
				},
			},
			Source: "community",
		},

		// === Web & APIs ===
		{
			ID:          "fetch",
			Name:        "Fetch",
			Slug:        "fetch",
			Description: "HTTP requests and web content fetching",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"web", "http", "api"},
			Tags:        []string{"http", "fetch", "api", "requests", "web"},
			Capabilities: types.Capabilities{
				Tools: []string{"fetch", "fetch_html", "fetch_json"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      500,
				Maintained: true,
				Score:      0.88,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-fetch"},
			},
			Source: "official",
		},
		{
			ID:          "puppeteer",
			Name:        "Puppeteer",
			Slug:        "puppeteer",
			Description: "Browser automation and web scraping with Puppeteer",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"web", "browser", "automation"},
			Tags:        []string{"puppeteer", "browser", "scraping", "automation", "testing"},
			Capabilities: types.Capabilities{
				Tools: []string{"navigate", "screenshot", "click", "type", "evaluate"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      450,
				Maintained: true,
				Score:      0.85,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-puppeteer"},
			},
			Source: "official",
		},

		// === AI & LLM ===
		{
			ID:          "openai",
			Name:        "OpenAI",
			Slug:        "openai",
			Description: "OpenAI API integration for GPT models and embeddings",
			Author:      "community",
			Repository:  "https://github.com/mcp-get-community/server-openai",
			Categories:  []string{"ai", "llm", "openai"},
			Tags:        []string{"openai", "gpt", "ai", "embeddings", "chat"},
			Capabilities: types.Capabilities{
				Tools: []string{"chat_completion", "create_embedding", "generate_image"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      200,
				Maintained: true,
				Score:      0.80,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-openai"},
				Env: []types.EnvVar{
					{Name: "OPENAI_API_KEY", Description: "OpenAI API key", Required: true},
				},
			},
			Source: "community",
		},

		// === Payments ===
		{
			ID:          "stripe",
			Name:        "Stripe",
			Slug:        "stripe",
			Description: "Stripe payment processing integration",
			Author:      "community",
			Repository:  "https://github.com/stripe/stripe-mcp",
			Categories:  []string{"payments", "stripe", "billing"},
			Tags:        []string{"stripe", "payments", "billing", "subscriptions"},
			Capabilities: types.Capabilities{
				Tools: []string{"create_payment_intent", "list_customers", "create_subscription"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      300,
				Maintained: true,
				Score:      0.85,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@stripe/mcp-server"},
				Env: []types.EnvVar{
					{Name: "STRIPE_SECRET_KEY", Description: "Stripe secret key", Required: true},
				},
			},
			Source: "community",
		},

		// === Communication ===
		{
			ID:          "slack",
			Name:        "Slack",
			Slug:        "slack",
			Description: "Slack messaging and channel management",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"communication", "slack", "messaging"},
			Tags:        []string{"slack", "messaging", "channels", "notifications"},
			Capabilities: types.Capabilities{
				Tools: []string{"send_message", "list_channels", "list_users"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      350,
				Maintained: true,
				Score:      0.83,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-slack"},
				Env: []types.EnvVar{
					{Name: "SLACK_BOT_TOKEN", Description: "Slack bot token", Required: true},
				},
			},
			Source: "official",
		},

		// === Memory & Context ===
		{
			ID:          "memory",
			Name:        "Memory",
			Slug:        "memory",
			Description: "Persistent memory and knowledge graph for context retention",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"memory", "context", "knowledge"},
			Tags:        []string{"memory", "context", "knowledge-graph", "persistence"},
			Capabilities: types.Capabilities{
				Tools:     []string{"store_memory", "recall_memory", "search_memories"},
				Resources: []string{"memories"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      700,
				Maintained: true,
				Score:      0.90,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-memory"},
			},
			Source: "official",
		},

		// === Search ===
		{
			ID:          "brave-search",
			Name:        "Brave Search",
			Slug:        "brave-search",
			Description: "Web search using Brave Search API",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"search", "web"},
			Tags:        []string{"search", "web", "brave"},
			Capabilities: types.Capabilities{
				Tools: []string{"search", "search_news"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      400,
				Maintained: true,
				Score:      0.85,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-brave-search"},
				Env: []types.EnvVar{
					{Name: "BRAVE_API_KEY", Description: "Brave Search API key", Required: true},
				},
			},
			Source: "official",
		},

		// === Time & Scheduling ===
		{
			ID:          "time",
			Name:        "Time",
			Slug:        "time",
			Description: "Time and timezone utilities",
			Author:      "modelcontextprotocol",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			Categories:  []string{"utilities", "time"},
			Tags:        []string{"time", "timezone", "date"},
			Capabilities: types.Capabilities{
				Tools: []string{"get_current_time", "convert_timezone"},
			},
			Compat: types.Compat{
				Languages:  []string{"all"},
				Frameworks: []string{"all"},
			},
			Quality: types.Quality{
				Stars:      200,
				Maintained: true,
				Score:      0.80,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-time"},
			},
			Source: "official",
		},

		// === Framework Specific ===
		{
			ID:          "prisma",
			Name:        "Prisma",
			Slug:        "prisma",
			Description: "Prisma ORM schema management and migrations",
			Author:      "community",
			Repository:  "https://github.com/nicholasgriffintn/mcp-server-prisma",
			Categories:  []string{"database", "orm", "prisma"},
			Tags:        []string{"prisma", "orm", "migrations", "schema"},
			Capabilities: types.Capabilities{
				Tools: []string{"introspect", "migrate", "generate", "studio"},
			},
			Compat: types.Compat{
				Languages:  []string{"typescript", "javascript"},
				Frameworks: []string{"nextjs", "express", "nestjs"},
			},
			Quality: types.Quality{
				Stars:      150,
				Maintained: true,
				Score:      0.78,
			},
			Install: types.Install{
				Command: "npx",
				Args:    []string{"-y", "mcp-server-prisma"},
			},
			Source: "community",
		},
	}
}
