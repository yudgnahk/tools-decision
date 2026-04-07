# API Specification: Tools Decision

## Overview

This document specifies the REST API for Tools Decision. The API provides programmatic access to project analysis, tool recommendations, and configuration management.

**Base URL:** `https://api.tools-decision.dev/v1`

**API Version:** v1 (current)

---

## Authentication

### API Key Authentication

All API requests require authentication via API key in the header.

```http
Authorization: Bearer td_live_xxxxxxxxxxxxxxxxxxxx
```

### API Key Types

| Type | Prefix | Use Case |
|------|--------|----------|
| Live | `td_live_` | Production usage |
| Test | `td_test_` | Development/testing |

### Rate Limits

| Tier | Requests/min | Requests/day |
|------|--------------|--------------|
| Free | 10 | 100 |
| Pro | 60 | 5,000 |
| Team | 120 | 25,000 |
| Enterprise | Custom | Custom |

Rate limit headers are included in all responses:

```http
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1714521600
```

---

## Common Response Format

### Success Response

```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "request_id": "req_abc123",
    "timestamp": "2026-04-07T10:30:00Z",
    "version": "v1"
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The project_path field is required",
    "details": {
      "field": "project_path",
      "reason": "missing_required_field"
    }
  },
  "meta": {
    "request_id": "req_abc123",
    "timestamp": "2026-04-07T10:30:00Z"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Malformed request body |
| `UNAUTHORIZED` | 401 | Invalid or missing API key |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `RATE_LIMITED` | 429 | Rate limit exceeded |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Endpoints

### Project Analysis

#### POST /analyze

Analyze a project and return detected context.

**Request:**

```http
POST /v1/analyze
Content-Type: application/json
Authorization: Bearer td_live_xxxx

{
  "project": {
    "files": [
      {
        "path": "package.json",
        "content": "{\"name\": \"my-app\", \"dependencies\": {\"react\": \"^18.0.0\", \"next\": \"^14.0.0\"}}"
      },
      {
        "path": "tsconfig.json",
        "content": "{\"compilerOptions\": {\"target\": \"ES2022\"}}"
      }
    ]
  },
  "options": {
    "include_frameworks": true,
    "include_dependencies": true,
    "include_inferred_needs": true
  }
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "analysis_id": "anl_abc123xyz",
    "project": {
      "detected_languages": [
        {
          "language": "typescript",
          "confidence": 0.95,
          "files_count": 45
        },
        {
          "language": "javascript",
          "confidence": 0.85,
          "files_count": 12
        }
      ],
      "detected_frameworks": [
        {
          "framework": "next.js",
          "version": "14.x",
          "confidence": 0.98
        },
        {
          "framework": "react",
          "version": "18.x",
          "confidence": 0.98
        }
      ],
      "detected_tools": [
        {
          "tool": "eslint",
          "config_file": ".eslintrc.js"
        },
        {
          "tool": "prettier",
          "config_file": ".prettierrc"
        }
      ],
      "project_type": "web_application",
      "project_type_confidence": 0.92
    },
    "inferred_needs": [
      {
        "category": "database",
        "reason": "Prisma schema detected",
        "confidence": 0.88
      },
      {
        "category": "testing",
        "reason": "Jest configuration found",
        "confidence": 0.95
      },
      {
        "category": "deployment",
        "reason": "Vercel configuration detected",
        "confidence": 0.90
      }
    ]
  },
  "meta": {
    "request_id": "req_def456",
    "timestamp": "2026-04-07T10:30:00Z",
    "processing_time_ms": 234
  }
}
```

---

### Tool Recommendations

#### POST /recommend

Get tool recommendations based on project analysis.

**Request:**

```http
POST /v1/recommend
Content-Type: application/json
Authorization: Bearer td_live_xxxx

{
  "analysis_id": "anl_abc123xyz",
  "options": {
    "limit": 10,
    "min_quality_score": 0.7,
    "categories": ["database", "testing", "documentation"],
    "exclude_installed": true
  }
}
```

**Alternative Request (without prior analysis):**

```json
{
  "context": {
    "languages": ["typescript", "python"],
    "frameworks": ["fastapi", "react"],
    "project_type": "full_stack_web",
    "needs": ["database", "authentication", "file_storage"]
  },
  "options": {
    "limit": 10
  }
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "recommendations": [
      {
        "tool": {
          "id": "tool_postgres-mcp",
          "name": "PostgreSQL MCP Server",
          "slug": "postgres-mcp",
          "description": "MCP server for PostgreSQL database operations",
          "author": "mcp-community",
          "repository": "https://github.com/mcp-community/postgres-mcp",
          "version": "1.2.0",
          "license": "MIT"
        },
        "match": {
          "score": 0.94,
          "reasons": [
            "Matches detected Prisma/PostgreSQL usage",
            "High compatibility with TypeScript projects",
            "Recommended for Next.js applications"
          ],
          "categories": ["database"]
        },
        "quality": {
          "overall_score": 0.89,
          "maintenance_score": 0.92,
          "popularity_score": 0.85,
          "security_score": 0.91,
          "github_stars": 1245,
          "weekly_downloads": 15420,
          "last_updated": "2026-03-15T00:00:00Z",
          "open_issues": 12,
          "contributors": 28
        },
        "capabilities": [
          "query_execution",
          "schema_inspection",
          "migration_support",
          "connection_pooling"
        ],
        "rank": 1
      },
      {
        "tool": {
          "id": "tool_jest-runner",
          "name": "Jest Test Runner MCP",
          "slug": "jest-runner-mcp",
          "description": "Run and analyze Jest tests via MCP",
          "author": "testing-tools",
          "repository": "https://github.com/testing-tools/jest-runner-mcp",
          "version": "2.0.1",
          "license": "MIT"
        },
        "match": {
          "score": 0.91,
          "reasons": [
            "Jest configuration detected in project",
            "Optimized for TypeScript/React testing"
          ],
          "categories": ["testing"]
        },
        "quality": {
          "overall_score": 0.87,
          "maintenance_score": 0.88,
          "popularity_score": 0.82,
          "security_score": 0.95,
          "github_stars": 892,
          "weekly_downloads": 8750,
          "last_updated": "2026-03-28T00:00:00Z",
          "open_issues": 5,
          "contributors": 15
        },
        "capabilities": [
          "test_execution",
          "coverage_reporting",
          "watch_mode",
          "snapshot_testing"
        ],
        "rank": 2
      }
    ],
    "summary": {
      "total_matches": 47,
      "returned": 10,
      "categories_covered": ["database", "testing", "documentation"],
      "average_quality_score": 0.85
    }
  },
  "meta": {
    "request_id": "req_ghi789",
    "timestamp": "2026-04-07T10:31:00Z",
    "model_version": "rec-v2.1"
  }
}
```

---

### Tool Search

#### GET /tools/search

Search for MCP tools by query.

**Request:**

```http
GET /v1/tools/search?q=database&category=storage&limit=20&offset=0
Authorization: Bearer td_live_xxxx
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query |
| `category` | string | No | Filter by category |
| `language` | string | No | Filter by language compatibility |
| `min_quality` | float | No | Minimum quality score (0-1) |
| `limit` | int | No | Results per page (default: 20, max: 100) |
| `offset` | int | No | Pagination offset |
| `sort` | string | No | Sort by: `relevance`, `quality`, `popularity`, `recent` |

**Response:**

```json
{
  "success": true,
  "data": {
    "tools": [
      {
        "id": "tool_postgres-mcp",
        "name": "PostgreSQL MCP Server",
        "slug": "postgres-mcp",
        "description": "MCP server for PostgreSQL database operations",
        "short_description": "PostgreSQL database integration",
        "categories": ["database", "storage"],
        "languages": ["typescript", "javascript", "python"],
        "quality_score": 0.89,
        "popularity_rank": 12,
        "verified": true
      }
    ],
    "pagination": {
      "total": 156,
      "limit": 20,
      "offset": 0,
      "has_more": true
    }
  },
  "meta": {
    "request_id": "req_jkl012",
    "timestamp": "2026-04-07T10:32:00Z"
  }
}
```

---

### Tool Details

#### GET /tools/{tool_id}

Get detailed information about a specific tool.

**Request:**

```http
GET /v1/tools/tool_postgres-mcp
Authorization: Bearer td_live_xxxx
```

**Response:**

```json
{
  "success": true,
  "data": {
    "tool": {
      "id": "tool_postgres-mcp",
      "name": "PostgreSQL MCP Server",
      "slug": "postgres-mcp",
      "description": "Full-featured MCP server for PostgreSQL database operations including queries, schema management, and migrations.",
      "author": {
        "name": "mcp-community",
        "url": "https://github.com/mcp-community",
        "verified": true
      },
      "repository": "https://github.com/mcp-community/postgres-mcp",
      "homepage": "https://postgres-mcp.dev",
      "documentation": "https://postgres-mcp.dev/docs",
      "version": "1.2.0",
      "license": "MIT",
      "created_at": "2025-06-15T00:00:00Z",
      "updated_at": "2026-03-15T00:00:00Z",
      
      "categories": ["database", "storage", "sql"],
      "tags": ["postgresql", "sql", "database", "orm"],
      
      "capabilities": {
        "tools": [
          {
            "name": "execute_query",
            "description": "Execute SQL queries against the database",
            "parameters": ["query", "params"]
          },
          {
            "name": "list_tables",
            "description": "List all tables in the database",
            "parameters": ["schema"]
          },
          {
            "name": "describe_table",
            "description": "Get schema information for a table",
            "parameters": ["table_name"]
          }
        ],
        "resources": [
          {
            "name": "database_schema",
            "description": "Current database schema"
          }
        ],
        "prompts": []
      },
      
      "compatibility": {
        "languages": ["typescript", "javascript", "python", "go"],
        "frameworks": ["next.js", "express", "fastapi", "django"],
        "platforms": ["macos", "linux", "windows"],
        "mcp_version": ">=1.0.0"
      },
      
      "installation": {
        "npm": "npm install @mcp-community/postgres-mcp",
        "pip": "pip install postgres-mcp",
        "docker": "docker pull mcpcommunity/postgres-mcp:latest"
      },
      
      "configuration": {
        "required": [
          {
            "name": "DATABASE_URL",
            "description": "PostgreSQL connection string",
            "type": "string",
            "example": "postgresql://user:pass@localhost:5432/mydb"
          }
        ],
        "optional": [
          {
            "name": "POOL_SIZE",
            "description": "Connection pool size",
            "type": "integer",
            "default": 10
          }
        ]
      },
      
      "quality": {
        "overall_score": 0.89,
        "scores": {
          "maintenance": 0.92,
          "popularity": 0.85,
          "security": 0.91,
          "documentation": 0.88,
          "testing": 0.85
        },
        "metrics": {
          "github_stars": 1245,
          "github_forks": 234,
          "weekly_downloads": 15420,
          "monthly_downloads": 58000,
          "open_issues": 12,
          "closed_issues": 456,
          "contributors": 28,
          "commits_last_month": 34
        },
        "last_audit": "2026-03-01T00:00:00Z"
      },
      
      "related_tools": [
        {
          "id": "tool_mysql-mcp",
          "name": "MySQL MCP Server",
          "reason": "Similar database functionality"
        },
        {
          "id": "tool_prisma-mcp",
          "name": "Prisma MCP Server",
          "reason": "Often used together"
        }
      ]
    }
  },
  "meta": {
    "request_id": "req_mno345",
    "timestamp": "2026-04-07T10:33:00Z"
  }
}
```

---

### Configuration Generation

#### POST /config/generate

Generate configuration files for installing recommended tools.

**Request:**

```http
POST /v1/config/generate
Content-Type: application/json
Authorization: Bearer td_live_xxxx

{
  "tools": [
    {
      "tool_id": "tool_postgres-mcp",
      "config": {
        "DATABASE_URL": "${env:DATABASE_URL}"
      }
    },
    {
      "tool_id": "tool_jest-runner",
      "config": {}
    }
  ],
  "format": "claude_desktop",
  "options": {
    "include_comments": true,
    "env_template": true
  }
}
```

**Supported Formats:**

| Format | Description |
|--------|-------------|
| `claude_desktop` | Claude Desktop config format |
| `cursor` | Cursor IDE config format |
| `vscode` | VS Code MCP extension format |
| `generic` | Generic MCP config |
| `docker_compose` | Docker Compose service definitions |

**Response:**

```json
{
  "success": true,
  "data": {
    "configs": {
      "claude_desktop": {
        "filename": "claude_desktop_config.json",
        "content": {
          "mcpServers": {
            "postgres-mcp": {
              "command": "npx",
              "args": ["-y", "@mcp-community/postgres-mcp"],
              "env": {
                "DATABASE_URL": "${env:DATABASE_URL}"
              }
            },
            "jest-runner": {
              "command": "npx",
              "args": ["-y", "@testing-tools/jest-runner-mcp"]
            }
          }
        },
        "path": "~/Library/Application Support/Claude/claude_desktop_config.json"
      }
    },
    "env_template": {
      "filename": ".env.mcp",
      "content": "# MCP Server Environment Variables\n\n# PostgreSQL MCP Server\nDATABASE_URL=postgresql://user:password@localhost:5432/database\n",
      "variables": [
        {
          "name": "DATABASE_URL",
          "description": "PostgreSQL connection string",
          "required": true,
          "tool": "postgres-mcp"
        }
      ]
    },
    "install_commands": [
      {
        "tool": "postgres-mcp",
        "command": "npm install -g @mcp-community/postgres-mcp"
      },
      {
        "tool": "jest-runner",
        "command": "npm install -g @testing-tools/jest-runner-mcp"
      }
    ]
  },
  "meta": {
    "request_id": "req_pqr678",
    "timestamp": "2026-04-07T10:34:00Z"
  }
}
```

---

### User Preferences

#### GET /preferences

Get user preferences for recommendations.

**Request:**

```http
GET /v1/preferences
Authorization: Bearer td_live_xxxx
```

**Response:**

```json
{
  "success": true,
  "data": {
    "preferences": {
      "default_format": "claude_desktop",
      "min_quality_score": 0.7,
      "preferred_languages": ["typescript", "python"],
      "excluded_tools": ["tool_deprecated-server"],
      "favorite_tools": ["tool_postgres-mcp", "tool_github-mcp"],
      "notification_settings": {
        "new_recommendations": true,
        "tool_updates": true,
        "security_alerts": true
      }
    }
  }
}
```

#### PATCH /preferences

Update user preferences.

**Request:**

```http
PATCH /v1/preferences
Content-Type: application/json
Authorization: Bearer td_live_xxxx

{
  "min_quality_score": 0.8,
  "preferred_languages": ["typescript", "python", "go"]
}
```

---

### Feedback

#### POST /feedback

Submit feedback on recommendations.

**Request:**

```http
POST /v1/feedback
Content-Type: application/json
Authorization: Bearer td_live_xxxx

{
  "recommendation_id": "rec_abc123",
  "tool_id": "tool_postgres-mcp",
  "feedback_type": "helpful",
  "rating": 5,
  "comment": "Perfect match for my project needs",
  "context": {
    "installed": true,
    "used_successfully": true
  }
}
```

**Feedback Types:**

| Type | Description |
|------|-------------|
| `helpful` | Recommendation was useful |
| `not_helpful` | Recommendation wasn't relevant |
| `installed` | User installed the tool |
| `rejected` | User explicitly rejected |
| `issue` | Problem with the tool |

---

### Team Management (Team/Enterprise)

#### GET /teams/{team_id}/configs

Get shared team configurations.

**Request:**

```http
GET /v1/teams/team_abc123/configs
Authorization: Bearer td_live_xxxx
```

**Response:**

```json
{
  "success": true,
  "data": {
    "configs": [
      {
        "id": "cfg_xyz789",
        "name": "Frontend Development",
        "description": "Standard MCP tools for frontend projects",
        "tools": [
          {
            "tool_id": "tool_prettier-mcp",
            "config": {}
          },
          {
            "tool_id": "tool_eslint-mcp",
            "config": {}
          }
        ],
        "created_by": "user_john",
        "created_at": "2026-03-01T00:00:00Z",
        "updated_at": "2026-03-15T00:00:00Z",
        "usage_count": 45
      }
    ]
  }
}
```

#### POST /teams/{team_id}/configs

Create a shared team configuration.

---

### Webhooks (Enterprise)

#### POST /webhooks

Register a webhook endpoint.

**Request:**

```http
POST /v1/webhooks
Content-Type: application/json
Authorization: Bearer td_live_xxxx

{
  "url": "https://example.com/webhooks/tools-decision",
  "events": ["tool.new", "tool.updated", "tool.security_alert"],
  "secret": "whsec_xxxxxxxxxxxx"
}
```

**Available Events:**

| Event | Description |
|-------|-------------|
| `tool.new` | New tool added matching your criteria |
| `tool.updated` | Installed tool has updates |
| `tool.security_alert` | Security issue detected |
| `config.synced` | Team config was synced |
| `recommendation.generated` | New recommendations available |

---

## SDK Examples

### JavaScript/TypeScript

```typescript
import { ToolsDecision } from '@tools-decision/sdk';

const client = new ToolsDecision({
  apiKey: process.env.TOOLS_DECISION_API_KEY
});

// Analyze a project
const analysis = await client.analyze({
  files: [
    { path: 'package.json', content: packageJsonContent }
  ]
});

// Get recommendations
const recommendations = await client.recommend({
  analysisId: analysis.id,
  options: { limit: 10 }
});

// Generate configuration
const config = await client.generateConfig({
  tools: recommendations.slice(0, 3).map(r => ({
    toolId: r.tool.id,
    config: {}
  })),
  format: 'claude_desktop'
});

console.log(config.content);
```

### Python

```python
from tools_decision import ToolsDecision

client = ToolsDecision(api_key=os.environ["TOOLS_DECISION_API_KEY"])

# Analyze a project
analysis = client.analyze(
    files=[
        {"path": "requirements.txt", "content": requirements_content}
    ]
)

# Get recommendations
recommendations = client.recommend(
    analysis_id=analysis.id,
    options={"limit": 10, "categories": ["database", "testing"]}
)

# Generate configuration
for rec in recommendations[:5]:
    print(f"{rec.tool.name}: {rec.match.score}")
```

### Go

```go
package main

import (
    "context"
    td "github.com/tools-decision/sdk-go"
)

func main() {
    client := td.NewClient(os.Getenv("TOOLS_DECISION_API_KEY"))
    
    // Analyze project
    analysis, err := client.Analyze(context.Background(), &td.AnalyzeRequest{
        Files: []td.File{
            {Path: "go.mod", Content: goModContent},
        },
    })
    
    // Get recommendations
    recs, err := client.Recommend(context.Background(), &td.RecommendRequest{
        AnalysisID: analysis.ID,
        Options: td.RecommendOptions{Limit: 10},
    })
    
    for _, rec := range recs {
        fmt.Printf("%s: %.2f\n", rec.Tool.Name, rec.Match.Score)
    }
}
```

---

## OpenAPI Specification

The full OpenAPI 3.0 specification is available at:

- **Interactive Docs:** https://api.tools-decision.dev/docs
- **OpenAPI JSON:** https://api.tools-decision.dev/openapi.json
- **OpenAPI YAML:** https://api.tools-decision.dev/openapi.yaml

---

## Changelog

### v1 (Current)

- Initial API release
- Project analysis endpoints
- Tool recommendation engine
- Configuration generation
- User preferences
- Team configurations (Team/Enterprise)
- Webhooks (Enterprise)

### Planned: v2

- GraphQL API
- Real-time subscriptions
- Batch operations
- Enhanced analytics endpoints

---

*Document Version: 1.0*  
*Last Updated: April 2026*  
*API Version: v1*
