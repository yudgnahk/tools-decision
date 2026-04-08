package matcher

import (
	"testing"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

func TestSkillMatcher_Match(t *testing.T) {
	m := NewSkillMatcher()

	// Create a minimal skill registry for testing
	skills := []types.Skill{
		{
			ID:          "go-debug",
			Name:        "Go Debugging Assistant",
			Slug:        "go-debug",
			Description: "Systematic debugging for Go applications",
			Category:    types.SkillCategoryDebugging,
			Compat: types.SkillCompat{
				Languages:    []string{"go"},
				Frameworks:   []string{"gin", "echo"},
				ProjectTypes: []string{"api", "cli"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{"github"},
			Quality:          types.Quality{Score: 0.9, Maintained: true},
			Source:           "official",
		},
		{
			ID:          "security-review",
			Name:        "Security Code Review",
			Slug:        "security-review",
			Description: "Security-focused code review checklist",
			Category:    types.SkillCategoryReview,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseCodeReview, types.UseCaseSecurity},
			},
			RequiredTools:    []string{"git"},
			RecommendedTools: []string{"github"},
			Quality:          types.Quality{Score: 0.95, Maintained: true},
			Source:           "official",
		},
		{
			ID:          "api-design",
			Name:        "REST API Design Guide",
			Slug:        "api-design",
			Description: "Best practices for designing RESTful APIs",
			Category:    types.SkillCategoryArchitecture,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"api"},
				UseCases:     []string{types.UseCaseArchitecture, types.UseCaseAPIDesign},
			},
			RequiredTools:    []string{"filesystem"},
			RecommendedTools: []string{"fetch"},
			Quality:          types.Quality{Score: 0.92, Maintained: true},
			Source:           "official",
		},
		{
			ID:          "python-debug",
			Name:        "Python Debugging Guide",
			Slug:        "python-debug",
			Description: "Systematic debugging for Python applications",
			Category:    types.SkillCategoryDebugging,
			Compat: types.SkillCompat{
				Languages:    []string{"python"},
				Frameworks:   []string{"fastapi", "django", "flask"},
				ProjectTypes: []string{"api", "cli", "web_app"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{"github"},
			Quality:          types.Quality{Score: 0.9, Maintained: true},
			Source:           "official",
		},
	}

	// Create server recommendations for synergy scoring
	serverRecommendations := []types.Recommendation{
		{
			Server: types.MCPServer{ID: "filesystem", Name: "Filesystem"},
			Score:  0.9,
		},
		{
			Server: types.MCPServer{ID: "git", Name: "Git"},
			Score:  0.85,
		},
		{
			Server: types.MCPServer{ID: "github", Name: "GitHub"},
			Score:  0.8,
		},
	}

	tests := []struct {
		name           string
		context        *types.ProjectContext
		wantTopSkill   string
		wantMinResults int
	}{
		{
			name: "Go project with debugging use case",
			context: &types.ProjectContext{
				Languages: []types.Language{{Name: "go", Confidence: 0.95}},
				UseCases:  []types.UseCase{{Name: types.UseCaseDebugging, Confidence: 0.9}},
				Type:      types.ProjectTypeAPI,
			},
			wantTopSkill:   "go-debug",
			wantMinResults: 1,
		},
		{
			name: "Python project should rank python skills higher",
			context: &types.ProjectContext{
				Languages: []types.Language{{Name: "python", Confidence: 0.95}},
				UseCases:  []types.UseCase{{Name: types.UseCaseDebugging, Confidence: 0.9}},
			},
			wantTopSkill:   "python-debug",
			wantMinResults: 1,
		},
		{
			name: "API project should include api-design skill",
			context: &types.ProjectContext{
				Type:     types.ProjectTypeAPI,
				UseCases: []types.UseCase{{Name: types.UseCaseAPIDesign, Confidence: 0.85}},
			},
			wantMinResults: 1,
		},
		{
			name: "Security review use case should match security skill",
			context: &types.ProjectContext{
				Languages: []types.Language{{Name: "typescript", Confidence: 0.9}},
				UseCases:  []types.UseCase{{Name: types.UseCaseSecurity, Confidence: 0.9}},
			},
			wantTopSkill:   "security-review",
			wantMinResults: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := m.Match(tt.context, skills, serverRecommendations, 10)

			if len(results) < tt.wantMinResults {
				t.Errorf("expected at least %d results, got %d", tt.wantMinResults, len(results))
			}

			if tt.wantTopSkill != "" && len(results) > 0 {
				if results[0].Skill.Slug != tt.wantTopSkill {
					t.Errorf("expected top skill %s, got %s (score: %.2f)", tt.wantTopSkill, results[0].Skill.Slug, results[0].Score)
					// Print all results for debugging
					for i, r := range results {
						t.Logf("  %d. %s (score: %.2f)", i+1, r.Skill.Slug, r.Score)
					}
				}
			}
		})
	}
}

func TestSkillMatcher_LanguageScore(t *testing.T) {
	m := NewSkillMatcher()

	tests := []struct {
		name         string
		projectLangs []types.Language
		skillLangs   []string
		wantScore    float64
		wantMatched  int
	}{
		{
			name:         "Exact language match",
			projectLangs: []types.Language{{Name: "go", Confidence: 0.95}},
			skillLangs:   []string{"go"},
			wantScore:    0.95,
			wantMatched:  1,
		},
		{
			name:         "All languages support",
			projectLangs: []types.Language{{Name: "rust", Confidence: 0.9}},
			skillLangs:   []string{"all"},
			wantScore:    0.7,
			wantMatched:  1, // "all languages" is returned
		},
		{
			name:         "No match",
			projectLangs: []types.Language{{Name: "go", Confidence: 0.9}},
			skillLangs:   []string{"python", "javascript"},
			wantScore:    0,
			wantMatched:  0,
		},
		{
			name: "Multiple languages, partial match",
			projectLangs: []types.Language{
				{Name: "go", Confidence: 0.95},
				{Name: "python", Confidence: 0.5},
			},
			skillLangs:  []string{"go"},
			wantScore:   0.95,
			wantMatched: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, matched := m.languageScore(tt.projectLangs, tt.skillLangs)
			if score != tt.wantScore {
				t.Errorf("expected score %.2f, got %.2f", tt.wantScore, score)
			}
			if len(matched) != tt.wantMatched {
				t.Errorf("expected %d matched, got %d", tt.wantMatched, len(matched))
			}
		})
	}
}

func TestSkillMatcher_ToolSynergyScore(t *testing.T) {
	m := NewSkillMatcher()

	tests := []struct {
		name         string
		serverIDs    map[string]bool
		skill        types.Skill
		wantMinScore float64
		wantMaxScore float64
	}{
		{
			name:      "All required tools available",
			serverIDs: map[string]bool{"filesystem": true, "git": true},
			skill: types.Skill{
				RequiredTools:    []string{"filesystem", "git"},
				RecommendedTools: []string{},
			},
			wantMinScore: 0.5,
			wantMaxScore: 1.0,
		},
		{
			name:      "Missing required tool",
			serverIDs: map[string]bool{"filesystem": true},
			skill: types.Skill{
				RequiredTools:    []string{"filesystem", "git"},
				RecommendedTools: []string{},
			},
			wantMinScore: 0.0,
			wantMaxScore: 0.35,
		},
		{
			name:      "Recommended tools available",
			serverIDs: map[string]bool{"filesystem": true, "github": true},
			skill: types.Skill{
				RequiredTools:    []string{},
				RecommendedTools: []string{"github", "postgres"},
			},
			wantMinScore: 0.5,
			wantMaxScore: 0.8,
		},
		{
			name:      "No tool requirements",
			serverIDs: map[string]bool{"filesystem": true},
			skill: types.Skill{
				RequiredTools:    []string{},
				RecommendedTools: []string{},
			},
			wantMinScore: 0.5,
			wantMaxScore: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := m.toolSynergyScore(tt.serverIDs, tt.skill)
			if score < tt.wantMinScore || score > tt.wantMaxScore {
				t.Errorf("expected score between %.2f and %.2f, got %.2f", tt.wantMinScore, tt.wantMaxScore, score)
			}
		})
	}
}

func TestDetectSynergies(t *testing.T) {
	servers := []types.Recommendation{
		{Server: types.MCPServer{ID: "filesystem", Name: "Filesystem"}},
		{Server: types.MCPServer{ID: "git", Name: "Git"}},
		{Server: types.MCPServer{ID: "github", Name: "GitHub"}},
	}

	skills := []types.SkillRecommendation{
		{
			Skill: types.Skill{
				ID:               "go-debug",
				Name:             "Go Debugging",
				RequiredTools:    []string{"filesystem", "git"},
				RecommendedTools: []string{"github"},
			},
		},
	}

	synergies := DetectSynergies(servers, skills)

	if len(synergies) == 0 {
		t.Error("expected to detect synergies")
	}

	// Should detect synergy with github (recommended tool)
	foundGithub := false
	for _, s := range synergies {
		if s.ServerID == "github" {
			foundGithub = true
			break
		}
	}
	if !foundGithub {
		t.Error("expected to find GitHub synergy")
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		slice    []string
		str      string
		expected bool
	}{
		{[]string{"go", "python"}, "go", true},
		{[]string{"go", "python"}, "GO", true}, // case insensitive
		{[]string{"go", "python"}, "rust", false},
		{[]string{}, "go", false},
		{[]string{"all"}, "all", true},
	}

	for _, tt := range tests {
		result := containsString(tt.slice, tt.str)
		if result != tt.expected {
			t.Errorf("containsString(%v, %s) = %v, want %v", tt.slice, tt.str, result, tt.expected)
		}
	}
}

func TestSkillMatcher_ArchetypeGatingForDocumentRepo(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "api-design",
			Name:        "REST API Design Guide",
			Slug:        "api-design",
			Description: "Best practices for API architecture",
			Category:    types.SkillCategoryArchitecture,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"api"},
				UseCases:     []string{types.UseCaseAPIDesign},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
		{
			ID:          "latex-docs",
			Name:        "LaTeX Documentation Workflow",
			Slug:        "latex-docs",
			Description: "Troubleshooting and authoring guidance for LaTeX docs",
			Category:    types.SkillCategoryDocumentation,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"library"},
				UseCases:     []string{types.UseCaseDocumentation},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
	}

	ctx := &types.ProjectContext{
		Archetypes: []types.ArchetypeSignal{
			{Name: types.ArchetypeDocumentAuthor, Confidence: 0.95},
		},
		UseCases: []types.UseCase{{Name: types.UseCaseDocumentation, Confidence: 0.85}},
	}

	results := m.Match(ctx, skills, nil, 10)
	if len(results) == 0 {
		t.Fatalf("expected at least one matched skill")
	}
	if results[0].Skill.Slug != "latex-docs" {
		t.Fatalf("expected documentation skill first, got %s", results[0].Skill.Slug)
	}
	for _, r := range results {
		if r.Skill.Slug == "api-design" {
			t.Fatalf("expected api-design to be gated out for document archetype")
		}
	}
}

func TestSkillMatcher_GenericSkillPenaltyForStrongArchetype(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "generic-review",
			Name:        "Generic Review",
			Slug:        "generic-review",
			Description: "Generic code review",
			Category:    types.SkillCategoryReview,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseCodeReview},
			},
			Quality: types.Quality{Score: 0.95, Maintained: true},
		},
		{
			ID:          "latex-authoring-build",
			Name:        "Document and LaTeX Build Troubleshooting",
			Slug:        "latex-authoring-build",
			Description: "LaTeX build troubleshooting",
			Category:    types.SkillCategoryDocumentation,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"latex"},
				ProjectTypes: []string{"library"},
				UseCases:     []string{types.UseCaseDocumentation},
			},
			Quality: types.Quality{Score: 0.88, Maintained: true},
		},
	}

	ctx := &types.ProjectContext{
		Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeDocumentAuthor, Confidence: 0.95}},
		UseCases:   []types.UseCase{{Name: types.UseCaseDocumentation, Confidence: 0.9}},
	}

	results := m.Match(ctx, skills, nil, 5)
	if len(results) == 0 {
		t.Fatalf("expected at least one result")
	}
	if results[0].Skill.Slug != "latex-authoring-build" {
		t.Fatalf("expected latex skill ranked first, got %s", results[0].Skill.Slug)
	}
	for _, r := range results {
		if r.Skill.Slug == "generic-review" {
			t.Fatalf("expected generic-review to be filtered for document archetype without matching use case")
		}
	}
}

func TestSkillMatcher_ArchetypeSpecificSkillGating(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "browser-automation-scripting",
			Name:        "Browser Automation Scripting",
			Slug:        "browser-automation-scripting",
			Description: "Automation patterns",
			Category:    types.SkillCategoryTesting,
			Compat: types.SkillCompat{
				Languages:    []string{"javascript", "typescript", "python"},
				Frameworks:   []string{"playwright", "puppeteer", "selenium"},
				ProjectTypes: []string{"cli", "web_app"},
				UseCases:     []string{types.UseCaseTesting},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
	}

	t.Run("gated for non-automation repo", func(t *testing.T) {
		ctx := &types.ProjectContext{
			Languages:  []types.Language{{Name: "go", Confidence: 0.95}},
			UseCases:   []types.UseCase{{Name: types.UseCaseTesting, Confidence: 0.9}},
			Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeDataProcessing, Confidence: 0.95}},
		}
		results := m.Match(ctx, skills, nil, 5)
		if len(results) != 0 {
			t.Fatalf("expected no automation skills, got %d", len(results))
		}
	})

	t.Run("allowed for automation repo", func(t *testing.T) {
		ctx := &types.ProjectContext{
			Languages:  []types.Language{{Name: "typescript", Confidence: 0.95}},
			Frameworks: []types.Framework{{Name: "playwright", Confidence: 0.95}},
			UseCases:   []types.UseCase{{Name: types.UseCaseTesting, Confidence: 0.9}},
			Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeAutomationBot, Confidence: 0.95}},
			Type:       types.ProjectTypeCLI,
		}
		results := m.Match(ctx, skills, nil, 5)
		if len(results) == 0 {
			t.Fatalf("expected automation skill to be included")
		}
	})
}

func TestSkillMatcher_DesktopSkillsGatedWithoutDesktopSignals(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "gui-event-loop-troubleshooting",
			Name:        "GUI Event Loop Troubleshooting",
			Slug:        "gui-event-loop-troubleshooting",
			Description: "Desktop event-loop diagnostics",
			Category:    types.SkillCategoryPerformance,
			Compat: types.SkillCompat{
				Languages:    []string{"rust", "typescript"},
				Frameworks:   []string{"tauri", "electron"},
				ProjectTypes: []string{"desktop"},
				UseCases:     []string{types.UseCaseDebugging, types.UseCasePerformance},
			},
			Quality: types.Quality{Score: 0.88, Maintained: true},
		},
	}

	ctx := &types.ProjectContext{
		Languages: []types.Language{{Name: "python", Confidence: 0.98}},
		UseCases:  []types.UseCase{{Name: types.UseCaseDebugging, Confidence: 0.9}},
		Archetypes: []types.ArchetypeSignal{
			{Name: types.ArchetypeAIContentPipe, Confidence: 0.95},
		},
		Type: types.ProjectTypeCLI,
	}

	results := m.Match(ctx, skills, nil, 5)
	if len(results) != 0 {
		t.Fatalf("expected desktop-only skill to be gated, got %d", len(results))
	}
}

func TestSkillMatcher_APIDomainSkillOutranksGenericForAPIArchetype(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "generic-review",
			Name:        "Generic Review",
			Slug:        "generic-review",
			Description: "General review",
			Category:    types.SkillCategoryReview,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseCodeReview},
			},
			Quality: types.Quality{Score: 0.95, Maintained: true},
		},
		{
			ID:          "api-design",
			Name:        "API Design",
			Slug:        "api-design",
			Description: "API design practices",
			Category:    types.SkillCategoryArchitecture,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"api"},
				UseCases:     []string{types.UseCaseAPIDesign, types.UseCaseArchitecture},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
	}

	ctx := &types.ProjectContext{
		UseCases: []types.UseCase{{Name: types.UseCaseAPIDesign, Confidence: 0.9}},
		Archetypes: []types.ArchetypeSignal{
			{Name: types.ArchetypeAPIService, Confidence: 0.95},
		},
		Type: types.ProjectTypeAPI,
	}

	results := m.Match(ctx, skills, nil, 5)
	if len(results) == 0 {
		t.Fatalf("expected at least one skill, got %d", len(results))
	}
	if results[0].Skill.Slug != "api-design" {
		t.Fatalf("expected api-design to outrank generic skill, got %s", results[0].Skill.Slug)
	}
}

func TestSkillMatcher_GoAPIDomainOutranksGenericWithoutUseCases(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "security-review",
			Name:        "Security Review",
			Slug:        "security-review",
			Description: "Generic review",
			Category:    types.SkillCategoryReview,
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseSecurity, types.UseCaseCodeReview},
			},
			Quality: types.Quality{Score: 0.95, Maintained: true},
		},
		{
			ID:          "go-debug",
			Name:        "Go Debugging Assistant",
			Slug:        "go-debug",
			Description: "Go debugging",
			Category:    types.SkillCategoryDebugging,
			Compat: types.SkillCompat{
				Languages:    []string{"go"},
				Frameworks:   []string{"gin", "echo"},
				ProjectTypes: []string{"api", "cli", "library"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
	}

	ctx := &types.ProjectContext{
		Languages: []types.Language{{Name: "go", Confidence: 0.98}},
		Frameworks: []types.Framework{
			{Name: "gin", Confidence: 0.95},
		},
		Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeAPIService, Confidence: 0.95}},
		Type:       types.ProjectTypeAPI,
	}

	results := m.Match(ctx, skills, nil, 5)
	if len(results) == 0 {
		t.Fatalf("expected at least one skill, got %d", len(results))
	}
	if results[0].Skill.Slug != "go-debug" {
		t.Fatalf("expected go-debug to outrank generic skill, got %s", results[0].Skill.Slug)
	}
}

func TestSkillMatcher_MixedProjectTypeDebugSkillNotGatedByAPIArchetype(t *testing.T) {
	m := NewSkillMatcher()

	skills := []types.Skill{
		{
			ID:          "go-debug",
			Name:        "Go Debugging Assistant",
			Slug:        "go-debug",
			Description: "Go debugging",
			Category:    types.SkillCategoryDebugging,
			Compat: types.SkillCompat{
				Languages:    []string{"go"},
				Frameworks:   []string{"cobra"},
				ProjectTypes: []string{"api", "cli", "library"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			Quality: types.Quality{Score: 0.9, Maintained: true},
		},
	}

	ctx := &types.ProjectContext{
		Languages:  []types.Language{{Name: "go", Confidence: 0.98}},
		Frameworks: []types.Framework{{Name: "cobra", Confidence: 0.95}},
		UseCases:   []types.UseCase{{Name: types.UseCaseDebugging, Confidence: 0.9}},
		Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeAutomationBot, Confidence: 0.9}},
		Type:       types.ProjectTypeCLI,
	}

	results := m.Match(ctx, skills, nil, 5)
	if len(results) == 0 {
		t.Fatalf("expected go-debug to remain eligible in non-api context")
	}
	if results[0].Skill.Slug != "go-debug" {
		t.Fatalf("expected go-debug to be recommended, got %s", results[0].Skill.Slug)
	}
}
