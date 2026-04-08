package types

// Skill represents an AI agent skill (instructions/prompts)
// Skills provide guidance on HOW an agent should approach tasks,
// complementing MCP servers which provide WHAT the agent can do.
type Skill struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Slug         string      `json:"slug"`
	Description  string      `json:"description"`
	Author       string      `json:"author"`
	Repository   string      `json:"repository,omitempty"`
	Category     string      `json:"category"`
	Instructions string      `json:"instructions"`
	Variables    []Variable  `json:"variables,omitempty"`
	Compat       SkillCompat `json:"compatibility"`
	Quality      Quality     `json:"quality"`
	Source       string      `json:"source"` // official, community, user

	// Tool synergy - which MCP servers work best with this skill
	RequiredTools    []string `json:"required_tools,omitempty"`
	RecommendedTools []string `json:"recommended_tools,omitempty"`
}

// Variable represents an input parameter for a skill
type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
}

// SkillCompat represents skill compatibility information
type SkillCompat struct {
	Languages    []string `json:"languages"`
	Frameworks   []string `json:"frameworks,omitempty"`
	ProjectTypes []string `json:"project_types,omitempty"`
	UseCases     []string `json:"use_cases,omitempty"`
}

// SkillCategory constants for skill categorization
const (
	SkillCategoryDebugging     = "debugging"
	SkillCategoryReview        = "review"
	SkillCategoryArchitecture  = "architecture"
	SkillCategoryTesting       = "testing"
	SkillCategoryDevOps        = "devops"
	SkillCategoryDocumentation = "documentation"
	SkillCategoryPerformance   = "performance"
	SkillCategoryBestPractices = "best-practices"
)

// UseCase represents a detected use case for skill matching
type UseCase struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// UseCaseCategory constants
const (
	UseCaseDebugging      = "debugging"
	UseCaseCodeReview     = "code-review"
	UseCaseArchitecture   = "architecture"
	UseCaseTesting        = "testing"
	UseCaseDevOps         = "devops"
	UseCaseDocumentation  = "documentation"
	UseCasePerformance    = "performance"
	UseCaseRefactoring    = "refactoring"
	UseCaseSecurity       = "security"
	UseCaseAPIDesign      = "api-design"
	UseCaseDatabaseDesign = "database-design"
)
