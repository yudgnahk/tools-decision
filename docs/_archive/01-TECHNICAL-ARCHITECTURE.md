# Technical Architecture Document

## Tools Decision - MCP Server & Tool Selection Platform

**Version:** 1.0  
**Date:** April 2026  
**Status:** Draft

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [System Overview](#2-system-overview)
3. [Architecture Principles](#3-architecture-principles)
4. [High-Level Architecture](#4-high-level-architecture)
5. [Core Components](#5-core-components)
6. [Data Architecture](#6-data-architecture)
7. [Integration Architecture](#7-integration-architecture)
8. [Security Architecture](#8-security-architecture)
9. [Scalability & Performance](#9-scalability--performance)
10. [Technology Stack](#10-technology-stack)
11. [Deployment Architecture](#11-deployment-architecture)
12. [API Design](#12-api-design)

---

## 1. Executive Summary

Tools Decision is an intelligent platform that helps developers and AI agents select the optimal MCP (Model Context Protocol) servers and tools for their specific projects. The system analyzes project context, aggregates data from multiple MCP registries, and provides AI-powered recommendations.

### Key Capabilities

- **Project Analysis**: Automatically detect tech stack, dependencies, and integration needs
- **Registry Aggregation**: Unified index of 20,000+ MCP servers from multiple sources
- **Intelligent Matching**: AI-powered recommendation engine
- **Configuration Generation**: Auto-generate MCP configs for various clients
- **Usage Analytics**: Track tool performance and optimize recommendations

---

## 2. System Overview

```
┌────────────────────────────────────────────────────────────────────────────┐
│                              TOOLS DECISION                                 │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │    CLI      │    │   Web UI    │    │   VS Code   │    │    API      │ │
│  │   Client    │    │  Dashboard  │    │  Extension  │    │   Clients   │ │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘ │
│         │                  │                  │                  │        │
│         └──────────────────┴──────────────────┴──────────────────┘        │
│                                    │                                       │
│                            ┌───────▼───────┐                              │
│                            │   API Gateway  │                              │
│                            │   (REST/gRPC)  │                              │
│                            └───────┬───────┘                              │
│                                    │                                       │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐│
│  │                         Core Services                                  ││
│  │                                                                        ││
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ ││
│  │  │   Project    │  │  Registry    │  │ Recommendation│  │   Config   │ ││
│  │  │   Analyzer   │  │  Aggregator  │  │    Engine     │  │  Generator │ ││
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └────────────┘ ││
│  │                                                                        ││
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ ││
│  │  │   Search     │  │   Analytics  │  │    Auth      │  │   Webhook  │ ││
│  │  │   Service    │  │   Service    │  │   Service    │  │   Service  │ ││
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └────────────┘ ││
│  └───────────────────────────────────────────────────────────────────────┘│
│                                    │                                       │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐│
│  │                         Data Layer                                     ││
│  │                                                                        ││
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ ││
│  │  │  PostgreSQL  │  │    Redis     │  │  Meilisearch │  │    S3      │ ││
│  │  │  (Primary)   │  │   (Cache)    │  │   (Search)   │  │  (Storage) │ ││
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └────────────┘ ││
│  └───────────────────────────────────────────────────────────────────────┘│
│                                                                            │
└────────────────────────────────────────────────────────────────────────────┘

                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
            ┌───────▼───────┐ ┌─────▼─────┐ ┌───────▼───────┐
            │  Official MCP │ │ Smithery  │ │    Glama      │
            │   Registry    │ │    API    │ │     API       │
            └───────────────┘ └───────────┘ └───────────────┘
```

---

## 3. Architecture Principles

### 3.1 Design Principles

| Principle | Description |
|-----------|-------------|
| **Modularity** | Loosely coupled services that can evolve independently |
| **Extensibility** | Easy to add new registry sources, analyzers, and output formats |
| **Resilience** | Graceful degradation when external services are unavailable |
| **Performance** | Sub-second response times for recommendations |
| **Privacy** | Project analysis happens locally; only metadata sent to cloud |

### 3.2 Key Constraints

- Must support offline mode for project analysis
- Must handle 20,000+ MCP servers efficiently
- Must integrate with existing MCP ecosystem (not replace it)
- Must work across Claude Desktop, Cursor, VS Code, and custom clients

---

## 4. High-Level Architecture

### 4.1 Architecture Style

**Hybrid Architecture**: Combination of:
- **Microservices** for backend services
- **Monolithic CLI** for local operations
- **Event-driven** for async operations (registry sync, analytics)

### 4.2 Component Interaction Flow

```
┌──────────────────────────────────────────────────────────────────────┐
│                        USER WORKFLOW                                  │
└──────────────────────────────────────────────────────────────────────┘

    Developer                                              AI Agent
        │                                                      │
        ▼                                                      ▼
┌───────────────┐                                    ┌───────────────┐
│ tools-decision│                                    │   API Call    │
│     init      │                                    │   /recommend  │
└───────┬───────┘                                    └───────┬───────┘
        │                                                    │
        ▼                                                    ▼
┌───────────────────────────────────────────────────────────────────────┐
│                         PROJECT ANALYZER                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │  File       │  │ Dependency  │  │  Framework  │  │ Integration │  │
│  │  Scanner    │  │  Parser     │  │  Detector   │  │  Detector   │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │
└───────────────────────────────────┬───────────────────────────────────┘
                                    │
                                    ▼ Project Context
┌───────────────────────────────────────────────────────────────────────┐
│                      RECOMMENDATION ENGINE                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │  Semantic   │  │   Rule      │  │    ML       │  │   Ranking   │  │
│  │  Matcher    │  │   Engine    │  │   Model     │  │   System    │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │
└───────────────────────────────────┬───────────────────────────────────┘
                                    │
                                    ▼ Ranked Recommendations
┌───────────────────────────────────────────────────────────────────────┐
│                       CONFIG GENERATOR                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │   Claude    │  │   Cursor    │  │   VS Code   │  │   Custom    │  │
│  │  Desktop    │  │   Config    │  │   Config    │  │   Format    │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │
└───────────────────────────────────┬───────────────────────────────────┘
                                    │
                                    ▼
                            Generated Config
                         (mcp.json / settings.json)
```

---

## 5. Core Components

### 5.1 Project Analyzer

**Purpose**: Extract project context to inform tool recommendations.

```typescript
interface ProjectContext {
  // Basic Info
  name: string;
  path: string;
  type: ProjectType; // 'web' | 'api' | 'cli' | 'library' | 'monorepo'
  
  // Languages & Frameworks
  languages: Language[];
  frameworks: Framework[];
  buildTools: BuildTool[];
  
  // Dependencies
  dependencies: Dependency[];
  devDependencies: Dependency[];
  
  // Integrations Detected
  integrations: Integration[];
  
  // Infrastructure
  infrastructure: {
    cloud?: 'aws' | 'gcp' | 'azure' | 'vercel' | 'other';
    database?: string[];
    messaging?: string[];
    monitoring?: string[];
  };
  
  // MCP Context
  existingMCPConfig?: MCPConfig;
  mcpClient?: 'claude-desktop' | 'cursor' | 'vscode' | 'other';
}
```

**Analyzers**:

| Analyzer | Files Parsed | Extracts |
|----------|--------------|----------|
| PackageJsonAnalyzer | package.json | Dependencies, scripts, engines |
| RequirementsAnalyzer | requirements.txt, Pipfile, pyproject.toml | Python deps |
| GoModAnalyzer | go.mod | Go modules |
| DockerAnalyzer | Dockerfile, docker-compose.yml | Services, bases |
| EnvAnalyzer | .env, .env.example | Integration hints |
| GitAnalyzer | .git/config | Remote repos, CI/CD |
| ConfigAnalyzer | Various config files | Framework detection |

### 5.2 Registry Aggregator

**Purpose**: Maintain a unified, searchable index of all MCP servers.

```typescript
interface MCPServer {
  // Identity
  id: string;                    // Unique across all registries
  sourceRegistry: RegistrySource;
  externalId: string;            // ID in source registry
  
  // Basic Info
  name: string;
  slug: string;
  description: string;
  longDescription?: string;
  
  // Categorization
  categories: Category[];
  tags: string[];
  
  // Technical Details
  transport: 'stdio' | 'http' | 'sse' | 'websocket';
  runtime: 'node' | 'python' | 'go' | 'rust' | 'binary';
  installCommand?: string;
  
  // Capabilities
  tools: MCPTool[];
  resources: MCPResource[];
  prompts: MCPPrompt[];
  
  // Quality Signals
  quality: {
    score: number;              // 0-100
    official: boolean;
    verified: boolean;
    lastUpdated: Date;
    maintenanceStatus: 'active' | 'maintained' | 'stale' | 'abandoned';
  };
  
  // Usage Stats
  stats: {
    weeklyDownloads?: number;
    githubStars?: number;
    totalInstalls?: number;
    successRate?: number;
  };
  
  // Requirements
  requirements: {
    authRequired: boolean;
    authType?: 'api_key' | 'oauth' | 'token';
    envVars?: string[];
    permissions?: string[];
  };
  
  // Matching Metadata
  matching: {
    keywords: string[];
    useCases: string[];
    compatibleWith: string[];    // frameworks, languages
    alternatives: string[];      // similar servers
  };
}
```

**Registry Sources**:

```typescript
interface RegistrySource {
  id: string;
  name: string;
  baseUrl: string;
  apiVersion: string;
  syncFrequency: string;        // cron expression
  priority: number;              // for deduplication
  
  // Methods
  fetchServers(): AsyncGenerator<MCPServer>;
  fetchServerDetails(id: string): Promise<MCPServer>;
  searchServers(query: string): Promise<MCPServer[]>;
}

// Implementations
class OfficialRegistrySource implements RegistrySource { }
class SmitherySource implements RegistrySource { }
class GlamaSource implements RegistrySource { }
class GitHubSource implements RegistrySource { }  // Direct from repos
```

### 5.3 Recommendation Engine

**Purpose**: Match project context to optimal MCP servers.

```typescript
interface RecommendationEngine {
  recommend(
    context: ProjectContext,
    options: RecommendationOptions
  ): Promise<Recommendation[]>;
}

interface RecommendationOptions {
  limit?: number;                // Max recommendations
  categories?: Category[];       // Filter by category
  excludeInstalled?: boolean;    // Skip already configured
  includeAlternatives?: boolean; // Include similar options
  qualityThreshold?: number;     // Min quality score
}

interface Recommendation {
  server: MCPServer;
  score: number;                 // 0-1 relevance score
  confidence: number;            // 0-1 confidence in recommendation
  reasons: RecommendationReason[];
  alternatives: MCPServer[];
  config: GeneratedConfig;
}

interface RecommendationReason {
  type: 'framework_match' | 'dependency_match' | 'integration_detected' | 
        'use_case_match' | 'popularity' | 'quality';
  description: string;
  weight: number;
}
```

**Matching Algorithm**:

```
┌─────────────────────────────────────────────────────────────────┐
│                    RECOMMENDATION PIPELINE                       │
└─────────────────────────────────────────────────────────────────┘

Step 1: Candidate Generation
┌─────────────────────────────────────────────────────────────────┐
│  Project Context ──► Keyword Extraction ──► Semantic Search     │
│                                                                 │
│  Inputs:                        Outputs:                        │
│  - Frameworks: [react, next]    - notion-mcp (content)          │
│  - Integrations: [stripe]       - stripe-mcp (payments)         │
│  - Languages: [typescript]      - github-mcp (version control)  │
│  - Use case hints               - supabase-mcp (database)       │
│                                 - ... (top 50 candidates)       │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
Step 2: Rule-Based Filtering
┌─────────────────────────────────────────────────────────────────┐
│  Apply hard constraints:                                        │
│  - Remove incompatible (e.g., Python-only for Node project)     │
│  - Remove below quality threshold                               │
│  - Remove deprecated/abandoned                                  │
│  - Remove already installed                                     │
│                                 ──► 30 candidates remain        │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
Step 3: Scoring & Ranking
┌─────────────────────────────────────────────────────────────────┐
│  Score = Σ (weight_i × signal_i)                                │
│                                                                 │
│  Signals:                                                       │
│  ┌─────────────────────┬────────┬──────────────────────────┐   │
│  │ Signal              │ Weight │ Calculation               │   │
│  ├─────────────────────┼────────┼──────────────────────────┤   │
│  │ Framework Match     │ 0.25   │ Exact match = 1.0         │   │
│  │ Integration Match   │ 0.25   │ Detected in .env/code     │   │
│  │ Semantic Similarity │ 0.20   │ Embedding cosine sim      │   │
│  │ Quality Score       │ 0.15   │ From registry             │   │
│  │ Popularity          │ 0.10   │ Normalized downloads      │   │
│  │ Freshness           │ 0.05   │ Days since update         │   │
│  └─────────────────────┴────────┴──────────────────────────┘   │
│                                 ──► Ranked list                 │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
Step 4: Diversification
┌─────────────────────────────────────────────────────────────────┐
│  Ensure variety in recommendations:                             │
│  - Max 2 per category                                           │
│  - Include at least 1 from each detected integration            │
│  - Balance official vs. community                               │
│                                 ──► Final 10 recommendations    │
└─────────────────────────────────────────────────────────────────┘
```

### 5.4 Config Generator

**Purpose**: Generate ready-to-use MCP configurations.

```typescript
interface ConfigGenerator {
  generate(
    recommendations: Recommendation[],
    target: ConfigTarget,
    options: ConfigOptions
  ): GeneratedConfig;
}

type ConfigTarget = 
  | 'claude-desktop'
  | 'cursor' 
  | 'vscode-copilot'
  | 'opencode'
  | 'custom';

interface ConfigOptions {
  includeComments?: boolean;
  includeEnvTemplate?: boolean;
  autoDetectExisting?: boolean;
  mergeStrategy?: 'replace' | 'merge' | 'append';
}

interface GeneratedConfig {
  format: 'json' | 'yaml';
  content: string;
  filePath: string;
  envTemplate?: string;
  setupInstructions: string[];
}
```

**Output Examples**:

```json
// Claude Desktop (~/.config/claude/mcp.json)
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_TOKEN": "${GITHUB_TOKEN}"
      }
    },
    "stripe": {
      "command": "npx", 
      "args": ["-y", "stripe-mcp"],
      "env": {
        "STRIPE_API_KEY": "${STRIPE_API_KEY}"
      }
    }
  }
}
```

```json
// Cursor (.cursor/mcp.json)
{
  "servers": {
    "github": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"]
    }
  }
}
```

---

## 6. Data Architecture

### 6.1 Database Schema

```sql
-- Core Tables

CREATE TABLE mcp_servers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_registry VARCHAR(50) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    long_description TEXT,
    transport VARCHAR(50) NOT NULL,
    runtime VARCHAR(50),
    install_command TEXT,
    
    -- Quality
    quality_score INTEGER CHECK (quality_score >= 0 AND quality_score <= 100),
    is_official BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    maintenance_status VARCHAR(50),
    
    -- Stats
    weekly_downloads INTEGER,
    github_stars INTEGER,
    total_installs INTEGER,
    success_rate DECIMAL(5,4),
    
    -- Requirements
    auth_required BOOLEAN DEFAULT FALSE,
    auth_type VARCHAR(50),
    env_vars JSONB,
    
    -- Matching
    keywords TEXT[],
    use_cases TEXT[],
    compatible_with TEXT[],
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_synced_at TIMESTAMP WITH TIME ZONE,
    
    UNIQUE(source_registry, external_id)
);

CREATE TABLE mcp_tools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    server_id UUID REFERENCES mcp_servers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    input_schema JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE mcp_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES mcp_categories(id),
    icon VARCHAR(50)
);

CREATE TABLE server_categories (
    server_id UUID REFERENCES mcp_servers(id) ON DELETE CASCADE,
    category_id UUID REFERENCES mcp_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (server_id, category_id)
);

-- User & Analytics Tables

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    github_id VARCHAR(100) UNIQUE,
    plan VARCHAR(50) DEFAULT 'free',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(255),
    context_hash VARCHAR(64),  -- Hash of project context for caching
    context JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id),
    server_id UUID REFERENCES mcp_servers(id),
    score DECIMAL(5,4),
    reasons JSONB,
    accepted BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    user_id UUID REFERENCES users(id),
    server_id UUID REFERENCES mcp_servers(id),
    properties JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_servers_quality ON mcp_servers(quality_score DESC);
CREATE INDEX idx_servers_downloads ON mcp_servers(weekly_downloads DESC);
CREATE INDEX idx_servers_keywords ON mcp_servers USING GIN(keywords);
CREATE INDEX idx_servers_compatible ON mcp_servers USING GIN(compatible_with);
CREATE INDEX idx_tools_server ON mcp_tools(server_id);
CREATE INDEX idx_recommendations_project ON recommendations(project_id);
CREATE INDEX idx_analytics_type_date ON analytics_events(event_type, created_at);
```

### 6.2 Search Index (Meilisearch)

```json
{
  "index": "mcp_servers",
  "primaryKey": "id",
  "searchableAttributes": [
    "name",
    "description",
    "keywords",
    "use_cases",
    "tools.name",
    "tools.description"
  ],
  "filterableAttributes": [
    "categories",
    "runtime",
    "transport",
    "is_official",
    "auth_required",
    "quality_score",
    "compatible_with"
  ],
  "sortableAttributes": [
    "quality_score",
    "weekly_downloads",
    "github_stars",
    "updated_at"
  ],
  "rankingRules": [
    "words",
    "typo",
    "proximity",
    "attribute",
    "sort",
    "exactness",
    "quality_score:desc"
  ]
}
```

### 6.3 Cache Strategy

```typescript
interface CacheConfig {
  // Registry data - cache aggressively
  serverList: {
    ttl: '1h',
    strategy: 'stale-while-revalidate'
  },
  
  serverDetails: {
    ttl: '15m',
    strategy: 'cache-first'
  },
  
  // Project analysis - cache by hash
  projectContext: {
    ttl: '24h',
    key: 'project:{contextHash}',
    strategy: 'cache-first'
  },
  
  // Recommendations - short cache
  recommendations: {
    ttl: '5m',
    key: 'recs:{projectHash}:{options}',
    strategy: 'network-first'
  },
  
  // Search results
  search: {
    ttl: '10m',
    key: 'search:{query}:{filters}',
    strategy: 'stale-while-revalidate'
  }
}
```

---

## 7. Integration Architecture

### 7.1 External Registry Integration

```typescript
// Registry Sync Pipeline
class RegistrySyncPipeline {
  private sources: RegistrySource[];
  private queue: JobQueue;
  
  async sync() {
    for (const source of this.sources) {
      await this.queue.add('sync-registry', {
        sourceId: source.id,
        fullSync: false
      });
    }
  }
  
  async processSync(job: SyncJob) {
    const source = this.getSource(job.sourceId);
    const servers = await source.fetchServers();
    
    for await (const server of servers) {
      // Normalize to our schema
      const normalized = this.normalize(server, source);
      
      // Deduplicate across registries
      const existing = await this.findDuplicate(normalized);
      if (existing) {
        await this.merge(existing, normalized);
      } else {
        await this.insert(normalized);
      }
      
      // Update search index
      await this.searchIndex.upsert(normalized);
    }
  }
}
```

### 7.2 Webhook Integration

```typescript
// Receive updates from registries
app.post('/webhooks/registry/:source', async (req, res) => {
  const { source } = req.params;
  const event = req.body;
  
  switch (event.type) {
    case 'server.created':
    case 'server.updated':
      await syncService.syncServer(source, event.serverId);
      break;
    case 'server.deleted':
      await syncService.removeServer(source, event.serverId);
      break;
  }
  
  res.status(200).send('OK');
});
```

---

## 8. Security Architecture

### 8.1 Security Principles

| Principle | Implementation |
|-----------|----------------|
| **Privacy First** | Project analysis runs locally; only anonymized metadata sent to cloud |
| **Minimal Permissions** | CLI only needs read access to project files |
| **Secure Secrets** | Never store or transmit actual API keys/secrets |
| **Auth Options** | Support API keys, OAuth, and anonymous usage |

### 8.2 Data Flow Security

```
┌─────────────────────────────────────────────────────────────────┐
│                    LOCAL (Developer Machine)                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌───────────────┐         ┌───────────────────────────────┐   │
│  │ Source Code   │──scan──►│    Project Analyzer (Local)   │   │
│  │ .env files    │         │    - Extracts metadata only   │   │
│  │ Configs       │         │    - No secrets transmitted   │   │
│  └───────────────┘         └───────────────┬───────────────┘   │
│                                            │                    │
│                                            ▼ (metadata only)    │
│                            ┌───────────────────────────────┐   │
│                            │  {                            │   │
│                            │    frameworks: ["react"],     │   │
│                            │    integrations: ["stripe"],  │   │
│                            │    languages: ["typescript"]  │   │
│                            │  }                            │   │
│                            └───────────────┬───────────────┘   │
└────────────────────────────────────────────┼────────────────────┘
                                             │ HTTPS
                                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    CLOUD (Tools Decision API)                    │
├─────────────────────────────────────────────────────────────────┤
│  - Receives only project metadata                               │
│  - No access to actual code or secrets                          │
│  - Returns recommendations + config templates                   │
│  - Config templates use ${VARIABLE} placeholders               │
└─────────────────────────────────────────────────────────────────┘
```

### 8.3 Authentication

```typescript
// API Authentication Options
interface AuthConfig {
  // Anonymous - rate limited
  anonymous: {
    rateLimit: '10 requests/minute',
    features: ['search', 'basic-recommendations']
  },
  
  // API Key - for developers
  apiKey: {
    rateLimit: '100 requests/minute',
    features: ['all']
  },
  
  // OAuth - for integrations
  oauth: {
    providers: ['github'],
    scopes: ['read:user'],
    features: ['all', 'sync-preferences']
  }
}
```

---

## 9. Scalability & Performance

### 9.1 Performance Targets

| Operation | Target | P99 |
|-----------|--------|-----|
| Project Analysis (local) | < 2s | 5s |
| Search Query | < 100ms | 200ms |
| Recommendation Request | < 500ms | 1s |
| Config Generation | < 50ms | 100ms |
| Registry Sync (full) | < 10min | 30min |

### 9.2 Scaling Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                      LOAD BALANCER                               │
│                    (Cloudflare / AWS ALB)                        │
└─────────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            │                 │                 │
            ▼                 ▼                 ▼
    ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
    │   API Pod 1   │ │   API Pod 2   │ │   API Pod N   │
    │   (Stateless) │ │   (Stateless) │ │   (Stateless) │
    └───────┬───────┘ └───────┬───────┘ └───────┬───────┘
            │                 │                 │
            └─────────────────┼─────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│     Redis     │     │  PostgreSQL   │     │  Meilisearch  │
│   (Cluster)   │     │   (Primary    │     │   (Cluster)   │
│               │     │   + Replicas) │     │               │
└───────────────┘     └───────────────┘     └───────────────┘
```

### 9.3 Caching Layers

1. **CDN Cache** - Static assets, documentation
2. **API Response Cache** - Redis, common queries
3. **Search Cache** - Meilisearch internal + Redis
4. **Local Cache** - CLI caches registry data locally

---

## 10. Technology Stack

### 10.1 Backend

| Component | Technology | Rationale |
|-----------|------------|-----------|
| API Server | **Go** | Performance, concurrency, single binary |
| Search Engine | **Meilisearch** | Fast, typo-tolerant, easy to deploy |
| Primary Database | **PostgreSQL** | Reliability, JSON support, full-text |
| Cache | **Redis** | Speed, pub/sub for real-time |
| Job Queue | **Redis + Asynq** | Simple, Go-native |
| Object Storage | **S3/R2** | Large file storage |

### 10.2 CLI

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | **Go** | Cross-platform, single binary, fast |
| CLI Framework | **Cobra** | Industry standard for Go CLIs |
| Config Parser | **Viper** | Multi-format support |
| UI | **Bubble Tea** | Rich terminal UI |

### 10.3 Frontend (Dashboard)

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Framework | **Next.js 14** | SSR, App Router |
| Styling | **Tailwind CSS** | Rapid development |
| Components | **shadcn/ui** | High quality, customizable |
| State | **Zustand** | Simple, lightweight |

### 10.4 Infrastructure

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Container | **Docker** | Standard containerization |
| Orchestration | **Kubernetes** | Scalability, self-healing |
| CI/CD | **GitHub Actions** | Integration with repo |
| Monitoring | **Grafana + Prometheus** | Observability |
| Logging | **Loki** | Log aggregation |

---

## 11. Deployment Architecture

### 11.1 Environment Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                        ENVIRONMENTS                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────────────┐   │
│  │   Local     │   │   Staging   │   │    Production       │   │
│  │   Dev       │   │             │   │                     │   │
│  ├─────────────┤   ├─────────────┤   ├─────────────────────┤   │
│  │ Docker      │   │ K8s (Small) │   │ K8s (Auto-scaling)  │   │
│  │ Compose     │   │ 1 replica   │   │ 2-10 replicas       │   │
│  │             │   │ Shared DB   │   │ Dedicated DB        │   │
│  │             │   │             │   │ Multi-region        │   │
│  └─────────────┘   └─────────────┘   └─────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 11.2 Kubernetes Resources

```yaml
# api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tools-decision-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: tools-decision-api
  template:
    spec:
      containers:
      - name: api
        image: tools-decision/api:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-credentials
              key: url
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: tools-decision-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: tools-decision-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## 12. API Design

### 12.1 REST API Endpoints

```yaml
openapi: 3.0.0
info:
  title: Tools Decision API
  version: 1.0.0

paths:
  # Project Analysis
  /v1/analyze:
    post:
      summary: Analyze project context
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProjectContext'
      responses:
        200:
          description: Analysis result
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AnalysisResult'

  # Recommendations
  /v1/recommend:
    post:
      summary: Get MCP server recommendations
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                context:
                  $ref: '#/components/schemas/ProjectContext'
                options:
                  $ref: '#/components/schemas/RecommendationOptions'
      responses:
        200:
          description: Recommendations
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Recommendation'

  # Search
  /v1/servers/search:
    get:
      summary: Search MCP servers
      parameters:
        - name: q
          in: query
          schema:
            type: string
        - name: category
          in: query
          schema:
            type: string
        - name: runtime
          in: query
          schema:
            type: string
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
      responses:
        200:
          description: Search results
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SearchResults'

  # Server Details
  /v1/servers/{id}:
    get:
      summary: Get server details
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        200:
          description: Server details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MCPServer'

  # Config Generation
  /v1/config/generate:
    post:
      summary: Generate MCP configuration
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                servers:
                  type: array
                  items:
                    type: string
                target:
                  type: string
                  enum: [claude-desktop, cursor, vscode, opencode]
      responses:
        200:
          description: Generated config
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GeneratedConfig'

  # Categories
  /v1/categories:
    get:
      summary: List all categories
      responses:
        200:
          description: Categories
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Category'
```

### 12.2 CLI Commands

```bash
# Initialize for current project
tools-decision init
# Analyzes current directory, recommends MCPs, generates config

# Interactive mode
tools-decision init -i
# Step-by-step wizard with choices

# Search servers
tools-decision search "github integration"
tools-decision search --category databases
tools-decision search --runtime python

# Get recommendations without generating config
tools-decision recommend
tools-decision recommend --limit 5 --json

# Generate config for specific servers
tools-decision config github stripe notion
tools-decision config github --target cursor

# Add a server to existing config
tools-decision add github
tools-decision add stripe --env STRIPE_API_KEY=sk_test_xxx

# Remove a server
tools-decision remove github

# List configured servers
tools-decision list
tools-decision list --detailed

# Update registry cache
tools-decision update

# Validate configuration
tools-decision validate

# Show server details
tools-decision info github
tools-decision info @modelcontextprotocol/server-github
```

---

## Appendix A: Directory Structure

```
tools-decision/
├── cmd/
│   ├── api/              # API server
│   │   └── main.go
│   └── cli/              # CLI tool
│       └── main.go
├── internal/
│   ├── analyzer/         # Project analysis
│   │   ├── analyzer.go
│   │   ├── nodejs.go
│   │   ├── python.go
│   │   ├── go.go
│   │   └── docker.go
│   ├── registry/         # Registry aggregation
│   │   ├── source.go
│   │   ├── official.go
│   │   ├── smithery.go
│   │   ├── glama.go
│   │   └── sync.go
│   ├── recommend/        # Recommendation engine
│   │   ├── engine.go
│   │   ├── scorer.go
│   │   └── ranker.go
│   ├── config/           # Config generation
│   │   ├── generator.go
│   │   ├── claude.go
│   │   ├── cursor.go
│   │   └── vscode.go
│   ├── search/           # Search service
│   │   └── search.go
│   ├── api/              # API handlers
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   └── routes.go
│   └── storage/          # Data access
│       ├── postgres.go
│       ├── redis.go
│       └── meilisearch.go
├── pkg/
│   ├── mcp/              # MCP types
│   │   └── types.go
│   └── project/          # Project types
│       └── types.go
├── web/                  # Next.js dashboard
│   ├── app/
│   ├── components/
│   └── lib/
├── scripts/
│   ├── sync-registries.sh
│   └── seed-data.sh
├── deployments/
│   ├── docker/
│   └── kubernetes/
├── docs/
└── README.md
```

---

## Appendix B: Data Flow Diagrams

### B.1 Recommendation Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │────►│  Analyze │────►│  Match   │────►│  Rank    │
│  Request │     │  Project │     │  Servers │     │  Results │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                      │                │                │
                      ▼                ▼                ▼
                 ┌──────────┐    ┌──────────┐    ┌──────────┐
                 │  Extract │    │  Query   │    │  Apply   │
                 │  Context │    │  Index   │    │  Weights │
                 └──────────┘    └──────────┘    └──────────┘
                                                      │
                                                      ▼
                                               ┌──────────┐
                                               │  Return  │
                                               │  Top N   │
                                               └──────────┘
```

### B.2 Registry Sync Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Cron    │────►│  Fetch   │────►│ Normalize│────►│  Upsert  │
│  Trigger │     │  Source  │     │   Data   │     │    DB    │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                      │                │                │
                      ▼                ▼                ▼
                 ┌──────────┐    ┌──────────┐    ┌──────────┐
                 │  Handle  │    │  Detect  │    │  Update  │
                 │  Errors  │    │  Dupes   │    │  Index   │
                 └──────────┘    └──────────┘    └──────────┘
```

---

*Document End*
