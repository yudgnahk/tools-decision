package types

// MCPServer represents an MCP server from the registry
type MCPServer struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	Slug                   string       `json:"slug"`
	Description            string       `json:"description"`
	Author                 string       `json:"author"`
	Repository             string       `json:"repository,omitempty"`
	NPM                    string       `json:"npm,omitempty"`
	PyPI                   string       `json:"pypi,omitempty"`
	Categories             []string     `json:"categories"`
	Tags                   []string     `json:"tags"`
	Capabilities           Capabilities `json:"capabilities"`
	Compat                 Compat       `json:"compatibility"`
	Quality                Quality      `json:"quality"`
	Install                Install      `json:"installation"`
	RequiresExplicitSignal bool         `json:"requires_explicit_signal,omitempty"`
	ExplicitSignals        []string     `json:"explicit_signals,omitempty"`
	RecommendedArchetypes  []Archetype  `json:"recommended_archetypes,omitempty"`
	ExcludedArchetypes     []Archetype  `json:"excluded_archetypes,omitempty"`
	Source                 string       `json:"source"` // Which registry this came from
}

// Capabilities represents what an MCP server can do
type Capabilities struct {
	Tools     []string `json:"tools"`
	Resources []string `json:"resources"`
	Prompts   []string `json:"prompts"`
}

// Compat represents compatibility information
type Compat struct {
	Languages  []string `json:"languages"`
	Frameworks []string `json:"frameworks"`
}

// Quality represents quality metrics
type Quality struct {
	Stars      int     `json:"stars,omitempty"`
	Downloads  int     `json:"downloads,omitempty"`
	LastUpdate string  `json:"last_update,omitempty"`
	Maintained bool    `json:"maintained"`
	Score      float64 `json:"score"` // Calculated quality score 0-1
}

// Install represents installation information
type Install struct {
	Command string   `json:"command"` // npx, pip, go install, etc.
	Args    []string `json:"args"`
	Env     []EnvVar `json:"env,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
}
