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
		Language:  0.25,
		Framework: 0.25,
		Category:  0.20,
		Service:   0.15,
		Quality:   0.15,
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
	var recommendations []types.Recommendation

	for _, server := range servers {
		score, reasons, matchedOn := m.score(ctx, server)

		if score > 0.1 { // Minimum threshold
			recommendations = append(recommendations, types.Recommendation{
				Server:    server,
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

// score calculates the match score for a server
func (m *Matcher) score(ctx *types.ProjectContext, server types.MCPServer) (float64, []string, []string) {
	var totalScore float64
	var reasons []string
	var matchedOn []string

	// Language matching
	langScore, langMatched := m.languageScore(ctx.Languages, server.Compat.Languages)
	if langScore > 0 {
		totalScore += langScore * m.weights.Language
		matchedOn = append(matchedOn, langMatched...)
		reasons = append(reasons, "Compatible with "+strings.Join(langMatched, ", "))
	}

	// Framework matching
	fwScore, fwMatched := m.frameworkScore(ctx.Frameworks, server.Compat.Frameworks)
	if fwScore > 0 {
		totalScore += fwScore * m.weights.Framework
		matchedOn = append(matchedOn, fwMatched...)
		reasons = append(reasons, "Works with "+strings.Join(fwMatched, ", "))
	}

	// Category/service matching
	svcScore, svcMatched := m.serviceScore(ctx.Services, server.Categories, server.Tags)
	if svcScore > 0 {
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

	return totalScore, reasons, matchedOn
}

// languageScore calculates the language match score
func (m *Matcher) languageScore(projectLangs []types.Language, serverLangs []string) (float64, []string) {
	var matched []string
	var maxScore float64

	serverLangSet := make(map[string]bool)
	for _, l := range serverLangs {
		serverLangSet[strings.ToLower(l)] = true
	}

	for _, lang := range projectLangs {
		if serverLangSet[strings.ToLower(lang.Name)] {
			matched = append(matched, lang.Name)
			if lang.Confidence > maxScore {
				maxScore = lang.Confidence
			}
		}
	}

	return maxScore, matched
}

// frameworkScore calculates the framework match score
func (m *Matcher) frameworkScore(projectFrameworks []types.Framework, serverFrameworks []string) (float64, []string) {
	var matched []string
	var maxScore float64

	serverFwSet := make(map[string]bool)
	for _, f := range serverFrameworks {
		serverFwSet[strings.ToLower(f)] = true
	}

	for _, fw := range projectFrameworks {
		if serverFwSet[strings.ToLower(fw.Name)] {
			matched = append(matched, fw.Name)
			if fw.Confidence > maxScore {
				maxScore = fw.Confidence
			}
		}
	}

	return maxScore, matched
}

// serviceScore calculates the service/category match score
func (m *Matcher) serviceScore(projectServices []types.Service, categories, tags []string) (float64, []string) {
	var matched []string
	var maxScore float64

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
			if svc.Confidence > maxScore {
				maxScore = svc.Confidence
			}
		}
	}

	return maxScore, matched
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
