package types

// ProjectContext represents the analyzed context of a project
type ProjectContext struct {
	Path       string      `json:"path"`
	Languages  []Language  `json:"languages"`
	Frameworks []Framework `json:"frameworks"`
	Tools      []Tool      `json:"tools"`
	Services   []Service   `json:"services"`
	Type       ProjectType `json:"type"`
	UseCases   []UseCase   `json:"use_cases,omitempty"`
}

// Language represents a detected programming language
type Language struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
	FilesCount int     `json:"files_count,omitempty"`
}

// Framework represents a detected framework
type Framework struct {
	Name       string  `json:"name"`
	Version    string  `json:"version,omitempty"`
	Confidence float64 `json:"confidence"`
}

// Tool represents a detected development tool
type Tool struct {
	Name       string `json:"name"`
	ConfigFile string `json:"config_file,omitempty"`
}

// Service represents a detected external service
type Service struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// ProjectType represents the type of project
type ProjectType string

const (
	ProjectTypeWebApp  ProjectType = "web_app"
	ProjectTypeAPI     ProjectType = "api"
	ProjectTypeCLI     ProjectType = "cli"
	ProjectTypeLibrary ProjectType = "library"
	ProjectTypeDesktop ProjectType = "desktop"
	ProjectTypeMobile  ProjectType = "mobile"
	ProjectTypeUnknown ProjectType = "unknown"
)
