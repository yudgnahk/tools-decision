package matcher

import (
	"sort"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// SkillMatcher matches project context to AI agent skills
type SkillMatcher struct {
	weights SkillMatchWeights
}

// SkillMatchWeights configures the scoring weights for skill matching
type SkillMatchWeights struct {
	Language    float64 // How well skill matches project language
	UseCase     float64 // How well skill matches detected use case
	Framework   float64 // How well skill matches project framework
	ToolSynergy float64 // How well skill works with recommended MCP servers
	Quality     float64 // Quality metrics of the skill
}

// DefaultSkillWeights returns the default skill matching weights
func DefaultSkillWeights() SkillMatchWeights {
	return SkillMatchWeights{
		Language:    0.30, // Skills are highly language-specific
		UseCase:     0.25, // Use case matching is important
		Framework:   0.20, // Framework compatibility
		ToolSynergy: 0.15, // Works well with recommended servers
		Quality:     0.10, // Quality metrics
	}
}

// NewSkillMatcher creates a new SkillMatcher with default weights
func NewSkillMatcher() *SkillMatcher {
	return &SkillMatcher{
		weights: DefaultSkillWeights(),
	}
}

// NewSkillMatcherWithWeights creates a new SkillMatcher with custom weights
func NewSkillMatcherWithWeights(weights SkillMatchWeights) *SkillMatcher {
	return &SkillMatcher{
		weights: weights,
	}
}

// Match finds the best skills for a project
func (m *SkillMatcher) Match(
	ctx *types.ProjectContext,
	skills []types.Skill,
	recommendedServers []types.Recommendation,
	limit int,
) []types.SkillRecommendation {
	var recommendations []types.SkillRecommendation

	// Build server ID set for synergy scoring
	serverIDs := make(map[string]bool)
	for _, rec := range recommendedServers {
		serverIDs[rec.Server.ID] = true
	}

	for _, skill := range skills {
		score, reasons, matchedOn := m.score(ctx, skill, serverIDs)
		score, eligible := m.applyArchetypePolicy(ctx, skill, score)
		if !eligible {
			continue
		}

		// Higher threshold for skills (0.15) since they should be more targeted
		if score > 0.15 {
			recommendations = append(recommendations, types.SkillRecommendation{
				Skill:     skill,
				Score:     score,
				Reasons:   reasons,
				MatchedOn: matchedOn,
			})
		}
	}

	// Sort by score descending
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

func (m *SkillMatcher) applyArchetypePolicy(
	ctx *types.ProjectContext,
	skill types.Skill,
	score float64,
) (float64, bool) {
	if len(ctx.Archetypes) == 0 {
		return score, true
	}

	apiConf := archetypeConfidence(ctx.Archetypes, types.ArchetypeAPIService)
	docConf := archetypeConfidence(ctx.Archetypes, types.ArchetypeDocumentAuthor)
	desktopConf := archetypeConfidence(ctx.Archetypes, types.ArchetypeDesktopApp)
	dataConf := archetypeConfidence(ctx.Archetypes, types.ArchetypeDataProcessing)
	autoConf := archetypeConfidence(ctx.Archetypes, types.ArchetypeAutomationBot)
	aiPipeConf := archetypeConfidence(ctx.Archetypes, types.ArchetypeAIContentPipe)

	hasAutomationSignal := hasFrameworkOrService(ctx, "playwright", "puppeteer", "selenium")
	hasDataSignal := hasFrameworkOrService(ctx, "pandas", "openpyxl", "excelize")
	hasAIPipelineSignal := hasFrameworkOrService(ctx, "openai", "anthropic", "langchain", "moviepy", "pydub", "ffmpeg")

	requiresAPI := containsString(skill.Compat.UseCases, types.UseCaseAPIDesign) ||
		containsString(skill.Compat.UseCases, types.UseCaseDatabaseDesign) ||
		containsString(skill.Compat.ProjectTypes, string(types.ProjectTypeAPI))

	explicitAPIUseCase := hasProjectUseCase(ctx.UseCases, types.UseCaseAPIDesign, 0.8) ||
		hasProjectUseCase(ctx.UseCases, types.UseCaseDatabaseDesign, 0.8)

	if requiresAPI && apiConf < 0.65 && !explicitAPIUseCase {
		return 0, false
	}

	slug := strings.ToLower(skill.Slug)
	if strings.Contains(slug, "browser-automation") && autoConf < 0.6 && !hasAutomationSignal {
		return 0, false
	}
	if strings.Contains(slug, "etl-data-quality") && dataConf < 0.6 && !hasDataSignal {
		return 0, false
	}
	if strings.Contains(slug, "ai-media-pipeline") && aiPipeConf < 0.6 && !hasAIPipelineSignal {
		return 0, false
	}
	if strings.Contains(slug, "latex") && docConf < 0.6 && !hasProjectUseCase(ctx.UseCases, types.UseCaseDocumentation, 0.75) {
		return 0, false
	}
	if (strings.Contains(slug, "rust-desktop") || strings.Contains(slug, "gui-event-loop")) && desktopConf < 0.6 {
		score *= 0.6
	}
	if strings.Contains(slug, "etl-data-quality") && dataConf < 0.6 {
		score *= 0.75
	}
	if strings.Contains(slug, "ai-media-pipeline") && aiPipeConf < 0.6 {
		score *= 0.6
	}

	if docConf >= 0.7 && isBackendOrInfraSkill(skill) && !explicitAPIUseCase {
		return 0, false
	}

	if ctx.Type != types.ProjectTypeUnknown && len(skill.Compat.ProjectTypes) > 0 {
		if !containsString(skill.Compat.ProjectTypes, "all") && !containsString(skill.Compat.ProjectTypes, string(ctx.Type)) {
			score *= 0.65
		}
	}

	if isGenericSkill(skill) {
		if docConf >= 0.7 || desktopConf >= 0.7 || dataConf >= 0.7 || autoConf >= 0.7 || aiPipeConf >= 0.7 {
			if !hasUseCaseIntersection(ctx.UseCases, skill.Compat.UseCases, 0.75) {
				return 0, false
			}
		}
		if docConf >= 0.7 || desktopConf >= 0.7 || dataConf >= 0.7 || autoConf >= 0.7 || aiPipeConf >= 0.7 {
			score *= 0.45
		} else if apiConf < 0.6 {
			score *= 0.7
		}
	}

	if desktopConf >= 0.7 && isBackendOrInfraSkill(skill) && !explicitAPIUseCase {
		score *= 0.6
	}

	if dataConf >= 0.7 && isBackendOrInfraSkill(skill) && !explicitAPIUseCase {
		score *= 0.7
	}

	if apiConf >= 0.7 && containsString(skill.Compat.ProjectTypes, string(types.ProjectTypeAPI)) {
		score *= 1.12
	}

	if docConf >= 0.7 && (skill.Category == types.SkillCategoryDocumentation || hasProjectType(skill.Compat.ProjectTypes, string(types.ProjectTypeLibrary))) {
		score *= 1.2
	}
	if docConf >= 0.7 && strings.Contains(slug, "latex") {
		score *= 1.35
	}

	if autoConf >= 0.7 && isAutomationSkill(skill) {
		score *= 1.15
	}
	if autoConf >= 0.7 && strings.Contains(slug, "browser-automation") {
		score *= 1.2
	}

	if aiPipeConf >= 0.7 && isAIPipelineSkill(skill) {
		score *= 1.15
	}
	if aiPipeConf >= 0.7 && strings.Contains(slug, "ai-media-pipeline") {
		score *= 1.25
	}
	if dataConf >= 0.7 && strings.Contains(slug, "etl-data-quality") {
		score *= 1.25
	}

	if strings.Contains(slug, "go-debug") && hasLanguage(ctx, "go", 0.9) {
		score *= 1.15
	}
	if strings.Contains(slug, "python-debug") && hasLanguage(ctx, "python", 0.9) {
		score *= 1.15
	}

	if score > 1.0 {
		score = 1.0
	}

	return score, true
}

func hasProjectUseCase(useCases []types.UseCase, name string, minConfidence float64) bool {
	for _, uc := range useCases {
		if strings.EqualFold(uc.Name, name) && uc.Confidence >= minConfidence {
			return true
		}
	}
	return false
}

func hasProjectType(projectTypes []string, p string) bool {
	for _, pt := range projectTypes {
		if strings.EqualFold(pt, p) {
			return true
		}
	}
	return false
}

func skillTerms(skill types.Skill) string {
	return strings.ToLower(skill.Slug + " " + skill.Name + " " + skill.Description + " " + skill.Category + " " + strings.Join(skill.Compat.UseCases, " "))
}

func isBackendOrInfraSkill(skill types.Skill) bool {
	if containsString(skill.Compat.UseCases, types.UseCaseAPIDesign) || containsString(skill.Compat.UseCases, types.UseCaseDatabaseDesign) {
		return true
	}
	if containsString(skill.Compat.ProjectTypes, string(types.ProjectTypeAPI)) {
		return true
	}
	terms := skillTerms(skill)
	for _, token := range []string{"api", "rest", "graphql", "grpc", "database", "sql", "microservice", "kubernetes", "docker", "devops"} {
		if strings.Contains(terms, token) {
			return true
		}
	}
	return false
}

func isAutomationSkill(skill types.Skill) bool {
	terms := skillTerms(skill)
	for _, token := range []string{"automation", "playwright", "puppeteer", "selenium", "browser", "e2e"} {
		if strings.Contains(terms, token) {
			return true
		}
	}
	return false
}

func isAIPipelineSkill(skill types.Skill) bool {
	terms := skillTerms(skill)
	for _, token := range []string{"ai", "llm", "openai", "anthropic", "media", "audio", "video", "pipeline", "ffmpeg"} {
		if strings.Contains(terms, token) {
			return true
		}
	}
	return false
}

func isGenericSkill(skill types.Skill) bool {
	return supportsAllCompat(skill.Compat.Languages) && supportsAllCompat(skill.Compat.Frameworks) && supportsAllCompat(skill.Compat.ProjectTypes)
}

func supportsAllCompat(values []string) bool {
	for _, v := range values {
		if strings.EqualFold(v, "all") {
			return true
		}
	}
	return false
}

func hasUseCaseIntersection(projectUseCases []types.UseCase, skillUseCases []string, minConfidence float64) bool {
	if len(projectUseCases) == 0 || len(skillUseCases) == 0 {
		return false
	}
	skillSet := make(map[string]bool, len(skillUseCases))
	for _, s := range skillUseCases {
		skillSet[strings.ToLower(s)] = true
	}
	for _, uc := range projectUseCases {
		if uc.Confidence < minConfidence {
			continue
		}
		if skillSet[strings.ToLower(uc.Name)] {
			return true
		}
	}
	return false
}

func hasFrameworkOrService(ctx *types.ProjectContext, names ...string) bool {
	for _, fw := range ctx.Frameworks {
		for _, n := range names {
			if strings.EqualFold(fw.Name, n) {
				return true
			}
		}
	}
	for _, svc := range ctx.Services {
		for _, n := range names {
			if strings.EqualFold(svc.Name, n) {
				return true
			}
		}
	}
	return false
}

func hasLanguage(ctx *types.ProjectContext, language string, minConfidence float64) bool {
	for _, l := range ctx.Languages {
		if strings.EqualFold(l.Name, language) && l.Confidence >= minConfidence {
			return true
		}
	}
	return false
}

// score calculates the match score for a skill
func (m *SkillMatcher) score(
	ctx *types.ProjectContext,
	skill types.Skill,
	serverIDs map[string]bool,
) (float64, []string, []string) {
	var totalScore float64
	var reasons []string
	var matchedOn []string

	// Language matching
	langScore, langMatched := m.languageScore(ctx.Languages, skill.Compat.Languages)
	if langScore > 0 {
		totalScore += langScore * m.weights.Language
		matchedOn = append(matchedOn, langMatched...)
		if len(langMatched) > 0 {
			reasons = append(reasons, "Designed for "+strings.Join(langMatched, ", "))
		}
	} else if !containsString(skill.Compat.Languages, "all") {
		// If skill requires specific language and none match, heavily penalize
		totalScore -= 0.2
	}

	// Use case matching
	useCaseScore, useCaseMatched := m.useCaseScore(ctx.UseCases, skill.Compat.UseCases, skill.Category)
	if useCaseScore > 0 {
		totalScore += useCaseScore * m.weights.UseCase
		matchedOn = append(matchedOn, useCaseMatched...)
		if len(useCaseMatched) > 0 {
			reasons = append(reasons, "Matches use case: "+strings.Join(useCaseMatched, ", "))
		}
	}

	// Framework matching
	fwScore, fwMatched := m.frameworkScore(ctx.Frameworks, skill.Compat.Frameworks)
	if fwScore > 0 {
		totalScore += fwScore * m.weights.Framework
		matchedOn = append(matchedOn, fwMatched...)
		if len(fwMatched) > 0 {
			reasons = append(reasons, "Works with "+strings.Join(fwMatched, ", "))
		}
	}

	// Tool synergy scoring
	synergyScore := m.toolSynergyScore(serverIDs, skill)
	totalScore += synergyScore * m.weights.ToolSynergy
	if synergyScore > 0.7 {
		reasons = append(reasons, "Synergizes with recommended MCP servers")
	}

	// Quality score
	qualityScore := m.qualityScore(skill.Quality)
	totalScore += qualityScore * m.weights.Quality
	if qualityScore > 0.8 {
		reasons = append(reasons, "High quality and well-maintained")
	}

	// Project type matching bonus
	if len(skill.Compat.ProjectTypes) > 0 && ctx.Type != types.ProjectTypeUnknown {
		if containsString(skill.Compat.ProjectTypes, string(ctx.Type)) || containsString(skill.Compat.ProjectTypes, "all") {
			totalScore += 0.05
			matchedOn = append(matchedOn, string(ctx.Type)+" project")
		}
	}

	return totalScore, reasons, matchedOn
}

// languageScore calculates the language match score
func (m *SkillMatcher) languageScore(projectLangs []types.Language, skillLangs []string) (float64, []string) {
	// If skill supports all languages, give moderate score
	if containsString(skillLangs, "all") {
		return 0.7, []string{"all languages"}
	}

	var matched []string
	var maxScore float64

	skillLangSet := make(map[string]bool)
	for _, l := range skillLangs {
		skillLangSet[strings.ToLower(l)] = true
	}

	for _, lang := range projectLangs {
		if skillLangSet[strings.ToLower(lang.Name)] {
			matched = append(matched, lang.Name)
			if lang.Confidence > maxScore {
				maxScore = lang.Confidence
			}
		}
	}

	return maxScore, matched
}

// useCaseScore calculates the use case match score
func (m *SkillMatcher) useCaseScore(projectUseCases []types.UseCase, skillUseCases []string, skillCategory string) (float64, []string) {
	var matched []string
	var maxScore float64

	// Build set of skill use cases including category
	skillUseCaseSet := make(map[string]bool)
	for _, uc := range skillUseCases {
		skillUseCaseSet[strings.ToLower(uc)] = true
	}
	// Add category as an implicit use case
	if skillCategory != "" {
		skillUseCaseSet[strings.ToLower(skillCategory)] = true
	}

	// Also map related use cases
	useCaseMappings := map[string][]string{
		"debugging":       {"debugging", "troubleshooting", "error-handling"},
		"code-review":     {"review", "security", "performance", "best-practices"},
		"architecture":    {"architecture", "design", "patterns", "microservices", "api-design"},
		"testing":         {"testing", "tdd", "coverage", "quality"},
		"devops":          {"devops", "ci-cd", "deployment", "infrastructure"},
		"performance":     {"performance", "optimization", "profiling"},
		"security":        {"security", "audit", "vulnerability"},
		"refactoring":     {"refactoring", "cleanup", "improvement"},
		"documentation":   {"documentation", "docs", "readme"},
		"database-design": {"database", "schema", "sql", "query"},
	}

	for _, uc := range projectUseCases {
		ucLower := strings.ToLower(uc.Name)

		// Direct match
		if skillUseCaseSet[ucLower] {
			matched = append(matched, uc.Name)
			if uc.Confidence > maxScore {
				maxScore = uc.Confidence
			}
			continue
		}

		// Check mapped use cases
		if mappedCases, ok := useCaseMappings[ucLower]; ok {
			for _, mc := range mappedCases {
				if skillUseCaseSet[mc] {
					matched = append(matched, uc.Name)
					if uc.Confidence > maxScore {
						maxScore = uc.Confidence
					}
					break
				}
			}
		}
	}

	return maxScore, matched
}

// frameworkScore calculates the framework match score
func (m *SkillMatcher) frameworkScore(projectFrameworks []types.Framework, skillFrameworks []string) (float64, []string) {
	// If skill supports all frameworks, give moderate score
	if containsString(skillFrameworks, "all") {
		return 0.5, nil
	}

	var matched []string
	var maxScore float64

	skillFwSet := make(map[string]bool)
	for _, f := range skillFrameworks {
		skillFwSet[strings.ToLower(f)] = true
	}

	for _, fw := range projectFrameworks {
		if skillFwSet[strings.ToLower(fw.Name)] {
			matched = append(matched, fw.Name)
			if fw.Confidence > maxScore {
				maxScore = fw.Confidence
			}
		}
	}

	return maxScore, matched
}

// toolSynergyScore calculates how well the skill works with recommended MCP servers
func (m *SkillMatcher) toolSynergyScore(serverIDs map[string]bool, skill types.Skill) float64 {
	// Check required tools
	if len(skill.RequiredTools) > 0 {
		requiredMet := 0
		for _, tool := range skill.RequiredTools {
			if serverIDs[tool] {
				requiredMet++
			}
		}
		// If not all required tools are available, penalize
		if requiredMet < len(skill.RequiredTools) {
			return 0.3 // Significant penalty
		}
	}

	// Bonus for recommended tools
	if len(skill.RecommendedTools) > 0 {
		recommendedMet := 0
		for _, tool := range skill.RecommendedTools {
			if serverIDs[tool] {
				recommendedMet++
			}
		}
		// Score based on how many recommended tools are available
		return 0.5 + 0.5*(float64(recommendedMet)/float64(len(skill.RecommendedTools)))
	}

	// No tool requirements, neutral score
	return 0.5
}

// qualityScore calculates a quality score from 0-1
func (m *SkillMatcher) qualityScore(quality types.Quality) float64 {
	// If pre-calculated score exists, use it
	if quality.Score > 0 {
		return quality.Score
	}

	// Calculate from metrics
	score := 0.5 // Base score

	// Stars contribution (up to 0.2)
	if quality.Stars > 100 {
		score += 0.1
	}
	if quality.Stars > 500 {
		score += 0.1
	}

	// Maintenance contribution (up to 0.2)
	if quality.Maintained {
		score += 0.2
	}

	// Downloads contribution (up to 0.1)
	if quality.Downloads > 1000 {
		score += 0.05
	}
	if quality.Downloads > 10000 {
		score += 0.05
	}

	if score > 1.0 {
		return 1.0
	}
	return score
}

// DetectSynergies finds beneficial combinations of servers and skills
func DetectSynergies(servers []types.Recommendation, skills []types.SkillRecommendation) []types.Synergy {
	var synergies []types.Synergy

	// Build server lookup
	serverByID := make(map[string]types.MCPServer)
	for _, rec := range servers {
		serverByID[rec.Server.ID] = rec.Server
	}

	for _, skillRec := range skills {
		skill := skillRec.Skill

		// Check recommended tools
		for _, toolID := range skill.RecommendedTools {
			if server, ok := serverByID[toolID]; ok {
				synergies = append(synergies, types.Synergy{
					SkillID:    skill.ID,
					ServerID:   server.ID,
					SkillName:  skill.Name,
					ServerName: server.Name,
					Reason:     "'" + skill.Name + "' skill works best with " + server.Name + " server",
				})
			}
		}

		// Check required tools (stronger synergy)
		for _, toolID := range skill.RequiredTools {
			if server, ok := serverByID[toolID]; ok {
				// Avoid duplicates
				isDuplicate := false
				for _, s := range synergies {
					if s.SkillID == skill.ID && s.ServerID == server.ID {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					synergies = append(synergies, types.Synergy{
						SkillID:    skill.ID,
						ServerID:   server.ID,
						SkillName:  skill.Name,
						ServerName: server.Name,
						Reason:     "'" + skill.Name + "' requires " + server.Name + " server",
					})
				}
			}
		}
	}

	return synergies
}

// containsString checks if a slice contains a string (case-insensitive)
func containsString(slice []string, str string) bool {
	strLower := strings.ToLower(str)
	for _, s := range slice {
		if strings.ToLower(s) == strLower {
			return true
		}
	}
	return false
}
