package matcher

import (
	"sort"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// Matcher matches project context to MCP servers
type Matcher struct {
	weights MatchWeights
}

type scoreBreakdown struct {
	LanguageMatched  bool
	FrameworkMatched bool
	ServiceMatched   bool
	LanguageStrong   bool
	FrameworkStrong  bool
	ServiceStrong    bool
	MatchedServices  []string
	MatchedFramework []string
}

// MatchWeights configures the scoring weights
type MatchWeights struct {
	Language  float64
	Framework float64
	Category  float64
	Service   float64
	Quality   float64
}

// DefaultWeights returns the default matching weights
func DefaultWeights() MatchWeights {
	return MatchWeights{
		Language:  0.20,
		Framework: 0.20,
		Category:  0.20,
		Service:   0.30,
		Quality:   0.10,
	}
}

// New creates a new Matcher with default weights
func New() *Matcher {
	return &Matcher{
		weights: DefaultWeights(),
	}
}

// NewWithWeights creates a new Matcher with custom weights
func NewWithWeights(weights MatchWeights) *Matcher {
	return &Matcher{
		weights: weights,
	}
}

// Match finds the best MCP servers for a project
func (m *Matcher) Match(ctx *types.ProjectContext, servers []types.MCPServer, limit int) []types.Recommendation {
	type candidate struct {
		rec       types.Recommendation
		breakdown scoreBreakdown
	}
	var candidates []candidate

	for _, server := range servers {
		score, reasons, matchedOn, breakdown := m.score(ctx, server)

		// Guardrail: do not recommend quality-only matches, except core baseline tools.
		if !breakdown.LanguageMatched && !breakdown.FrameworkMatched && !breakdown.ServiceMatched {
			if !isCoreBaseline(server.Slug) {
				continue
			}
		}

		// Guardrail: some integrations require explicit project intent.
		if !hasExplicitIntent(server, breakdown) {
			continue
		}

		minThreshold := 0.1
		if isCoreBaseline(server.Slug) && !breakdown.LanguageMatched && !breakdown.FrameworkMatched && !breakdown.ServiceMatched {
			minThreshold = 0.08
		}

		if score > minThreshold {
			candidates = append(candidates, candidate{rec: types.Recommendation{
				Server:    server,
				Score:     score,
				Reasons:   reasons,
				MatchedOn: matchedOn,
			}, breakdown: breakdown})
		}
	}

	// Stable sort by strong-signal count, then score, then slug.
	sort.SliceStable(candidates, func(i, j int) bool {
		si := strongSignalCount(candidates[i].breakdown)
		sj := strongSignalCount(candidates[j].breakdown)
		if si != sj {
			return si > sj
		}
		if candidates[i].rec.Score != candidates[j].rec.Score {
			return candidates[i].rec.Score > candidates[j].rec.Score
		}
		return candidates[i].rec.Server.Slug < candidates[j].rec.Server.Slug
	})

	var recommendations []types.Recommendation
	for _, c := range candidates {
		recommendations = append(recommendations, c.rec)
	}

	// Limit results
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

// score calculates the match score for a server
func (m *Matcher) score(ctx *types.ProjectContext, server types.MCPServer) (float64, []string, []string, scoreBreakdown) {
	var totalScore float64
	var reasons []string
	var matchedOn []string
	breakdown := scoreBreakdown{}

	// Language matching
	langScore, langMatched, langStrong := m.languageScore(ctx.Languages, server.Compat.Languages)
	if langScore > 0 {
		breakdown.LanguageMatched = true
		breakdown.LanguageStrong = langStrong
		totalScore += langScore * m.weights.Language
		matchedOn = append(matchedOn, langMatched...)
		reasons = append(reasons, "Compatible with "+strings.Join(langMatched, ", "))
	}

	// Framework matching
	fwScore, fwMatched, fwStrong := m.frameworkScore(ctx.Frameworks, server.Compat.Frameworks)
	if fwScore > 0 {
		breakdown.FrameworkMatched = true
		breakdown.FrameworkStrong = fwStrong
		breakdown.MatchedFramework = fwMatched
		totalScore += fwScore * m.weights.Framework
		matchedOn = append(matchedOn, fwMatched...)
		reasons = append(reasons, "Works with "+strings.Join(fwMatched, ", "))
	}

	// Category/service matching
	svcScore, svcMatched, svcStrong := m.serviceScore(ctx.Services, server.Categories, server.Tags)
	if svcScore > 0 {
		breakdown.ServiceMatched = true
		breakdown.ServiceStrong = svcStrong
		breakdown.MatchedServices = svcMatched
		totalScore += svcScore * m.weights.Service
		matchedOn = append(matchedOn, svcMatched...)
		reasons = append(reasons, "Provides "+strings.Join(svcMatched, ", "))
	}

	// Quality score
	qualityScore := m.qualityScore(server.Quality)
	totalScore += qualityScore * m.weights.Quality
	if qualityScore > 0.7 {
		reasons = append(reasons, "High quality and well-maintained")
	}

	return totalScore, reasons, matchedOn, breakdown
}

func requiresExplicitIntent(server types.MCPServer) bool {
	if server.RequiresExplicitSignal {
		return true
	}

	// Backward-compatible fallback for sources without metadata.
	slug := strings.ToLower(server.Slug)
	switch slug {
	case "stripe", "slack", "aws", "kubernetes", "docker", "prisma":
		return true
	default:
		return false
	}
}

func hasExplicitIntent(server types.MCPServer, b scoreBreakdown) bool {
	if !requiresExplicitIntent(server) {
		return true
	}

	matched := make(map[string]bool)
	for _, s := range b.MatchedServices {
		matched[strings.ToLower(s)] = true
	}
	for _, f := range b.MatchedFramework {
		matched[strings.ToLower(f)] = true
	}

	for _, signal := range server.ExplicitSignals {
		if matched[strings.ToLower(signal)] {
			return true
		}
	}

	// Fallback behavior if metadata not populated.
	if len(server.ExplicitSignals) == 0 {
		slug := strings.ToLower(server.Slug)
		switch slug {
		case "prisma":
			return matched["prisma"] || matched["nextjs"] || matched["nestjs"] || matched["express"]
		case "stripe":
			return matched["stripe"] || matched["payment"] || matched["payments"] || matched["billing"]
		case "aws":
			return matched["aws"] || matched["s3"]
		case "kubernetes":
			return matched["kubernetes"] || matched["k8s"]
		case "docker":
			return matched["docker"] || matched["container"]
		case "slack":
			return matched["slack"]
		default:
			return b.FrameworkMatched || b.ServiceMatched
		}
	}

	return false
}

func strongSignalCount(b scoreBreakdown) int {
	count := 0
	if b.ServiceStrong {
		count++
	}
	if b.FrameworkStrong {
		count++
	}
	if b.LanguageStrong {
		count++
	}
	return count
}

func isCoreBaseline(slug string) bool {
	switch slug {
	case "filesystem", "git", "github", "fetch", "memory":
		return true
	default:
		return false
	}
}

func isGenericServiceSignal(name string) bool {
	switch strings.ToLower(name) {
	case "database", "db", "sql", "storage", "cache", "testing", "test":
		return true
	default:
		return false
	}
}

// languageScore calculates the language match score
func (m *Matcher) languageScore(projectLangs []types.Language, serverLangs []string) (float64, []string, bool) {
	var matched []string
	var maxScore float64
	strong := false

	serverLangSet := make(map[string]bool)
	for _, l := range serverLangs {
		serverLangSet[strings.ToLower(l)] = true
	}

	for _, lang := range projectLangs {
		if serverLangSet[strings.ToLower(lang.Name)] {
			matched = append(matched, lang.Name)
			if lang.Confidence >= 0.93 {
				strong = true
			}
			if lang.Confidence > maxScore {
				maxScore = lang.Confidence
			}
		}
	}

	return maxScore, matched, strong
}

// frameworkScore calculates the framework match score
func (m *Matcher) frameworkScore(projectFrameworks []types.Framework, serverFrameworks []string) (float64, []string, bool) {
	var matched []string
	var maxScore float64
	strong := false

	serverFwSet := make(map[string]bool)
	for _, f := range serverFrameworks {
		serverFwSet[strings.ToLower(f)] = true
	}

	for _, fw := range projectFrameworks {
		if serverFwSet[strings.ToLower(fw.Name)] {
			matched = append(matched, fw.Name)
			if fw.Confidence >= 0.85 {
				strong = true
			}
			if fw.Confidence > maxScore {
				maxScore = fw.Confidence
			}
		}
	}

	return maxScore, matched, strong
}

// serviceScore calculates the service/category match score
func (m *Matcher) serviceScore(projectServices []types.Service, categories, tags []string) (float64, []string, bool) {
	var matched []string
	var maxScore float64
	strong := false

	// Combine categories and tags for matching
	serverTerms := make(map[string]bool)
	for _, c := range categories {
		serverTerms[strings.ToLower(c)] = true
	}
	for _, t := range tags {
		serverTerms[strings.ToLower(t)] = true
	}

	for _, svc := range projectServices {
		if serverTerms[strings.ToLower(svc.Name)] {
			matched = append(matched, svc.Name)
			score := svc.Confidence
			if isGenericServiceSignal(svc.Name) {
				score = svc.Confidence * 0.35
			} else if svc.Confidence >= 0.85 {
				strong = true
			}
			if score > maxScore {
				maxScore = score
			}
		}
	}

	return maxScore, matched, strong
}

// qualityScore calculates a quality score from 0-1
func (m *Matcher) qualityScore(quality types.Quality) float64 {
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

	return min(score, 1.0)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
