package evaluation

import (
	"testing"

	"github.com/yudgnahk/tools-decision/internal/matcher"
	"github.com/yudgnahk/tools-decision/internal/registry"
	"github.com/yudgnahk/tools-decision/pkg/types"
)

const (
	// Initial CI gate thresholds. Tighten toward plan targets as Track E progresses.
	targetGoodRate         = 0.65
	targetSkillPrecision   = 0.80
	targetServerPrecision  = 0.90
	targetGenericSkillRate = 0.40

	// Repo-level "good" classification heuristic used for GoodRate.
	goodSkillPrecisionCutoff  = 0.80
	goodServerPrecisionCutoff = 0.80
)

type benchmarkFixture struct {
	name         string
	ctx          *types.ProjectContext
	allowServers []string
	denyServers  []string
	allowSkills  []string
	denySkills   []string
}

func TestNineRepoBenchmarkQualityGate(t *testing.T) {
	servers := registry.GetEmbeddedServers()
	skills := registry.GetEmbeddedSkills()

	m := matcher.New()
	sm := matcher.NewSkillMatcher()

	fixtures := []benchmarkFixture{
		{
			name: "banana-playground",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "go", Confidence: 0.95}, {Name: "typescript", Confidence: 0.85}},
				Frameworks: []types.Framework{
					{Name: "playwright", Confidence: 0.95},
					{Name: "cobra", Confidence: 0.8},
				},
				Services: []types.Service{{Name: "openai", Confidence: 0.85}},
				UseCases: []types.UseCase{{Name: types.UseCaseTesting, Confidence: 0.9}, {Name: types.UseCaseDebugging, Confidence: 0.8}},
				Type:     types.ProjectTypeCLI,
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeAutomationBot, Confidence: 0.95},
					{Name: types.ArchetypeCLITool, Confidence: 0.8},
				},
			},
			allowServers: []string{"filesystem", "git", "github", "fetch", "memory", "puppeteer", "openai"},
			denyServers:  []string{"prisma", "stripe", "postgres", "mysql", "mongodb"},
			allowSkills:  []string{"browser-automation-scripting", "go-debug", "js-debug", "test-strategy", "tdd-workflow"},
			denySkills:   []string{"api-design", "database-optimization", "microservices-design"},
		},
		{
			name: "diem_thi_2025",
			ctx: &types.ProjectContext{
				Languages:  []types.Language{{Name: "go", Confidence: 0.96}},
				Frameworks: []types.Framework{{Name: "excelize", Confidence: 0.9}},
				UseCases:   []types.UseCase{{Name: types.UseCaseTesting, Confidence: 0.85}, {Name: types.UseCasePerformance, Confidence: 0.75}},
				Type:       types.ProjectTypeCLI,
				Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeDataProcessing, Confidence: 0.95}, {Name: types.ArchetypeCLITool, Confidence: 0.8}},
			},
			allowServers: []string{"filesystem", "git", "github", "fetch", "memory", "sqlite"},
			denyServers:  []string{"prisma", "stripe", "kubernetes", "docker", "aws"},
			allowSkills:  []string{"etl-data-quality", "go-debug", "error-handling", "test-strategy", "code-refactoring"},
			denySkills:   []string{"api-design", "database-optimization", "microservices-design"},
		},
		{
			name: "gokit",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "go", Confidence: 0.98}},
				Frameworks: []types.Framework{
					{Name: "gin", Confidence: 0.95},
				},
				Services: []types.Service{
					{Name: "postgresql", Confidence: 0.95},
					{Name: "mysql", Confidence: 0.9},
					{Name: "redis", Confidence: 0.9},
					{Name: "database", Confidence: 0.9},
				},
				UseCases: []types.UseCase{{Name: types.UseCaseAPIDesign, Confidence: 0.9}, {Name: types.UseCasePerformance, Confidence: 0.8}},
				Type:     types.ProjectTypeAPI,
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeAPIService, Confidence: 0.97},
				},
			},
			allowServers: []string{"postgres", "mysql", "redis", "sqlite", "filesystem", "git", "github", "fetch", "memory", "mongodb"},
			denyServers:  []string{"stripe", "prisma"},
			allowSkills:  []string{"go-project-structure", "go-debug", "api-design", "performance-review", "database-optimization", "microservices-design"},
		},
		{
			name: "goxkey",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "rust", Confidence: 0.98}},
				Frameworks: []types.Framework{
					{Name: "tauri", Confidence: 0.95},
				},
				UseCases: []types.UseCase{{Name: types.UseCaseDebugging, Confidence: 0.9}, {Name: types.UseCasePerformance, Confidence: 0.85}},
				Type:     types.ProjectTypeDesktop,
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeDesktopApp, Confidence: 0.97},
				},
			},
			allowServers: []string{"filesystem", "git", "github", "fetch", "memory", "puppeteer"},
			denyServers:  []string{"prisma", "stripe", "postgres", "mysql", "mongodb", "kubernetes"},
			allowSkills:  []string{"rust-desktop-debug", "gui-event-loop-troubleshooting", "error-handling", "performance-review"},
			denySkills:   []string{"api-design", "database-optimization", "microservices-design"},
		},
		{
			name: "lightnovel-anime-ai",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "python", Confidence: 0.98}},
				Frameworks: []types.Framework{
					{Name: "langchain", Confidence: 0.9},
					{Name: "moviepy", Confidence: 0.9},
				},
				Services: []types.Service{{Name: "openai", Confidence: 0.95}, {Name: "pydub", Confidence: 0.9}, {Name: "ffmpeg", Confidence: 0.9}},
				UseCases: []types.UseCase{{Name: types.UseCaseDebugging, Confidence: 0.9}, {Name: types.UseCasePerformance, Confidence: 0.85}},
				Type:     types.ProjectTypeCLI,
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeAIContentPipe, Confidence: 0.97},
					{Name: types.ArchetypeAutomationBot, Confidence: 0.8},
				},
			},
			allowServers: []string{"openai", "fetch", "filesystem", "git", "github", "memory", "puppeteer"},
			denyServers:  []string{"prisma", "stripe"},
			allowSkills:  []string{"python-ai-media-pipeline", "python-debug", "browser-automation-scripting", "performance-review"},
			denySkills:   []string{"api-design", "microservices-design", "database-optimization"},
		},
		{
			name: "opencode",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "typescript", Confidence: 0.95}},
				Frameworks: []types.Framework{
					{Name: "nextjs", Confidence: 0.85},
				},
				Services:   []types.Service{{Name: "aws", Confidence: 0.95}},
				UseCases:   []types.UseCase{{Name: types.UseCaseDevOps, Confidence: 0.85}, {Name: types.UseCaseAPIDesign, Confidence: 0.8}},
				Type:       types.ProjectTypeAPI,
				Archetypes: []types.ArchetypeSignal{{Name: types.ArchetypeAPIService, Confidence: 0.9}},
			},
			allowServers: []string{"aws", "filesystem", "git", "github", "fetch", "memory", "docker", "kubernetes"},
			denyServers:  []string{"stripe"},
			allowSkills:  []string{"js-debug", "ci-cd-setup", "api-design", "microservices-design", "security-review"},
			denySkills:   []string{"latex-authoring-build"},
		},
		{
			name: "payment-service",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "go", Confidence: 0.98}},
				Frameworks: []types.Framework{
					{Name: "gin", Confidence: 0.95},
					{Name: "gorm", Confidence: 0.9},
				},
				Services: []types.Service{{Name: "postgresql", Confidence: 0.97}, {Name: "database", Confidence: 0.9}},
				UseCases: []types.UseCase{{Name: types.UseCaseAPIDesign, Confidence: 0.9}, {Name: types.UseCaseSecurity, Confidence: 0.85}},
				Type:     types.ProjectTypeAPI,
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeAPIService, Confidence: 0.97},
				},
			},
			allowServers: []string{"postgres", "mysql", "sqlite", "mongodb", "redis", "filesystem", "git", "github", "fetch", "memory"},
			denyServers:  []string{"stripe", "prisma"},
			allowSkills:  []string{"go-project-structure", "go-debug", "microservices-design", "api-design", "security-review", "database-optimization"},
		},
		{
			name: "resume",
			ctx: &types.ProjectContext{
				UseCases: []types.UseCase{{Name: types.UseCaseDocumentation, Confidence: 0.95}},
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeDocumentAuthor, Confidence: 0.99},
				},
			},
			allowServers: []string{"filesystem", "git", "github", "fetch", "memory"},
			denyServers:  []string{"prisma", "stripe", "postgres", "mysql", "mongodb", "redis", "kubernetes", "docker", "aws"},
			allowSkills:  []string{"latex-authoring-build", "security-review", "error-handling", "code-refactoring"},
			denySkills:   []string{"api-design", "database-optimization", "microservices-design"},
		},
		{
			name: "ticket-services",
			ctx: &types.ProjectContext{
				Languages: []types.Language{{Name: "go", Confidence: 0.98}},
				Frameworks: []types.Framework{
					{Name: "gin", Confidence: 0.95},
				},
				Services: []types.Service{{Name: "postgresql", Confidence: 0.97}, {Name: "database", Confidence: 0.9}},
				UseCases: []types.UseCase{{Name: types.UseCaseAPIDesign, Confidence: 0.9}, {Name: types.UseCaseArchitecture, Confidence: 0.85}},
				Type:     types.ProjectTypeAPI,
				Archetypes: []types.ArchetypeSignal{
					{Name: types.ArchetypeAPIService, Confidence: 0.97},
				},
			},
			allowServers: []string{"postgres", "sqlite", "mysql", "mongodb", "redis", "filesystem", "git", "github", "fetch", "memory"},
			denyServers:  []string{"stripe", "prisma"},
			allowSkills:  []string{"go-project-structure", "go-debug", "microservices-design", "api-design", "performance-review", "database-optimization"},
		},
	}

	type result struct {
		name            string
		serverPrecision float64
		skillPrecision  float64
		genericSkill    float64
		good            bool
	}

	var results []result
	for _, fx := range fixtures {
		serverRecs := m.Match(fx.ctx, servers, 10)
		skillRecs := sm.Match(fx.ctx, skills, serverRecs, 5)

		topServerSlugs := takeServerSlugs(serverRecs, 10)
		topSkillSlugs := takeSkillSlugs(skillRecs, 5)

		assertNoDenied(t, fx.name, "server", topServerSlugs, fx.denyServers)
		assertNoDenied(t, fx.name, "skill", topSkillSlugs, fx.denySkills)

		sPrecision := precisionAt(topServerSlugs, fx.allowServers, 10)
		kPrecision := precisionAt(topSkillSlugs, fx.allowSkills, 5)
		gRate := genericSkillRate(skillRecs, 5)

		isGood := sPrecision >= goodServerPrecisionCutoff && kPrecision >= goodSkillPrecisionCutoff
		results = append(results, result{
			name:            fx.name,
			serverPrecision: sPrecision,
			skillPrecision:  kPrecision,
			genericSkill:    gRate,
			good:            isGood,
		})

		t.Logf("%s: servers=%v skills=%v serverP=%.2f skillP=%.2f genericSkillRate=%.2f", fx.name, topServerSlugs, topSkillSlugs, sPrecision, kPrecision, gRate)
	}

	avgServerPrecision := 0.0
	avgSkillPrecision := 0.0
	avgGenericSkillRate := 0.0
	goodCount := 0

	for _, r := range results {
		avgServerPrecision += r.serverPrecision
		avgSkillPrecision += r.skillPrecision
		avgGenericSkillRate += r.genericSkill
		if r.good {
			goodCount++
		}
	}

	n := float64(len(results))
	avgServerPrecision /= n
	avgSkillPrecision /= n
	avgGenericSkillRate /= n
	goodRate := float64(goodCount) / n

	t.Logf("metrics: GoodRate=%.2f Precision@5(skills)=%.2f Precision@10(servers)=%.2f GenericSkillRate=%.2f", goodRate, avgSkillPrecision, avgServerPrecision, avgGenericSkillRate)

	if goodRate < targetGoodRate {
		t.Fatalf("GoodRate below threshold: got %.2f want >= %.2f", goodRate, targetGoodRate)
	}
	if avgSkillPrecision < targetSkillPrecision {
		t.Fatalf("Precision@5(skills) below threshold: got %.2f want >= %.2f", avgSkillPrecision, targetSkillPrecision)
	}
	if avgServerPrecision < targetServerPrecision {
		t.Fatalf("Precision@10(servers) below threshold: got %.2f want >= %.2f", avgServerPrecision, targetServerPrecision)
	}
	if avgGenericSkillRate >= targetGenericSkillRate {
		t.Fatalf("GenericSkillRate above threshold: got %.2f want < %.2f", avgGenericSkillRate, targetGenericSkillRate)
	}
}

func takeServerSlugs(recs []types.Recommendation, limit int) []string {
	if limit <= 0 || len(recs) == 0 {
		return nil
	}
	if len(recs) < limit {
		limit = len(recs)
	}
	out := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, recs[i].Server.Slug)
	}
	return out
}

func takeSkillSlugs(recs []types.SkillRecommendation, limit int) []string {
	if limit <= 0 || len(recs) == 0 {
		return nil
	}
	if len(recs) < limit {
		limit = len(recs)
	}
	out := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, recs[i].Skill.Slug)
	}
	return out
}

func precisionAt(got []string, allow []string, k int) float64 {
	if len(got) == 0 || len(allow) == 0 || k <= 0 {
		return 0
	}
	if len(got) < k {
		k = len(got)
	}
	allowSet := make(map[string]bool, len(allow))
	for _, s := range allow {
		allowSet[s] = true
	}
	matches := 0
	for i := 0; i < k; i++ {
		if allowSet[got[i]] {
			matches++
		}
	}
	return float64(matches) / float64(k)
}

func assertNoDenied(t *testing.T, repoName, kind string, got, deny []string) {
	t.Helper()
	if len(deny) == 0 || len(got) == 0 {
		return
	}
	denySet := make(map[string]bool, len(deny))
	for _, d := range deny {
		denySet[d] = true
	}
	for _, s := range got {
		if denySet[s] {
			t.Fatalf("%s: denied %s %q found in top results: %v", repoName, kind, s, got)
		}
	}
}

func genericSkillRate(recs []types.SkillRecommendation, k int) float64 {
	if len(recs) == 0 || k <= 0 {
		return 0
	}
	if len(recs) < k {
		k = len(recs)
	}
	generic := 0
	for i := 0; i < k; i++ {
		skill := recs[i].Skill
		if supportsAll(skill.Compat.Languages) && supportsAll(skill.Compat.Frameworks) && supportsAll(skill.Compat.ProjectTypes) {
			generic++
		}
	}
	return float64(generic) / float64(k)
}

func supportsAll(values []string) bool {
	for _, v := range values {
		if v == "all" {
			return true
		}
	}
	return false
}
