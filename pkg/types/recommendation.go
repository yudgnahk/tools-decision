package types

// Recommendation represents a recommended MCP server
type Recommendation struct {
	Server    MCPServer `json:"server"`
	Score     float64   `json:"score"`      // Overall match score 0-1
	Reasons   []string  `json:"reasons"`    // Why this was recommended
	MatchedOn []string  `json:"matched_on"` // What project features matched
}

// SkillRecommendation represents a recommended AI agent skill
type SkillRecommendation struct {
	Skill     Skill    `json:"skill"`
	Score     float64  `json:"score"`      // Overall match score 0-1
	Reasons   []string `json:"reasons"`    // Why this was recommended
	MatchedOn []string `json:"matched_on"` // What project features matched
}

// Synergy represents a beneficial combination of server and skill
type Synergy struct {
	SkillID    string `json:"skill_id"`
	ServerID   string `json:"server_id"`
	SkillName  string `json:"skill_name"`
	ServerName string `json:"server_name"`
	Reason     string `json:"reason"`
}

// CombinedRecommendation contains both server and skill recommendations
type CombinedRecommendation struct {
	Servers   []Recommendation      `json:"servers"`
	Skills    []SkillRecommendation `json:"skills"`
	Synergies []Synergy             `json:"synergies,omitempty"`
}

// ConfigOutput represents a generated configuration
type ConfigOutput struct {
	Format   string         `json:"format"`
	Filename string         `json:"filename"`
	Path     string         `json:"path"`
	Content  map[string]any `json:"content"`
	EnvVars  []EnvVar       `json:"env_vars,omitempty"`
}

// SkillOutput represents a generated skill command file
type SkillOutput struct {
	Skill    Skill  `json:"skill"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Content  string `json:"content"`
}
