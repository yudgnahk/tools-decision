package types

// Recommendation represents a recommended MCP server
type Recommendation struct {
	Server    MCPServer `json:"server"`
	Score     float64   `json:"score"`      // Overall match score 0-1
	Reasons   []string  `json:"reasons"`    // Why this was recommended
	MatchedOn []string  `json:"matched_on"` // What project features matched
}

// ConfigOutput represents a generated configuration
type ConfigOutput struct {
	Format   string         `json:"format"`
	Filename string         `json:"filename"`
	Path     string         `json:"path"`
	Content  map[string]any `json:"content"`
	EnvVars  []EnvVar       `json:"env_vars,omitempty"`
}
