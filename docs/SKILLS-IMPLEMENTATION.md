# Skills Feature Implementation Plan

This document outlines the implementation plan for adding AI agent skills selection alongside MCP servers.

## Overview

**Goal**: Enable tools-decision to recommend AI agent skills (instructions/prompts) alongside MCP servers (tools), providing users with both capabilities AND expertise for their AI coding assistants.

## Key Concepts

| Concept | MCP Servers | Skills |
|---------|-------------|--------|
| **Nature** | Executable functions | Instruction templates |
| **Purpose** | What the agent CAN do | How the agent SHOULD think |
| **Examples** | `query`, `read_file`, `create_pr` | "Debug Go apps", "Security review" |
| **Output** | API results | Guided agent behavior |

## Implementation Phases

### Phase 1: Type Definitions

**Files to create/modify:**
- [ ] `pkg/types/skill.go` - Skill type definition
- [ ] `pkg/types/recommendation.go` - Add SkillRecommendation type

**Skill type structure:**
```go
type Skill struct {
    ID           string       `json:"id"`
    Name         string       `json:"name"`
    Slug         string       `json:"slug"`
    Description  string       `json:"description"`
    Author       string       `json:"author"`
    Category     string       `json:"category"`
    Instructions string       `json:"instructions"`
    Variables    []Variable   `json:"variables"`
    RequiredTools    []string `json:"required_tools"`
    RecommendedTools []string `json:"recommended_tools"`
    Compat       SkillCompat  `json:"compatibility"`
    Quality      Quality      `json:"quality"`
    Source       string       `json:"source"`
}
```

### Phase 2: Embedded Skills Registry

**Files to create:**
- [ ] `internal/registry/embedded_skills.go` - Curated skill library

**Initial skill categories:**
1. **Debugging** - Language-specific debugging workflows
2. **Code Review** - Security, performance, best practices
3. **Architecture** - Design patterns, microservices, API design
4. **Testing** - Test strategy, coverage, TDD workflows
5. **DevOps** - CI/CD, deployment, infrastructure

**Target: 15-20 embedded skills** covering major languages and use cases.

### Phase 3: Skill Matcher Algorithm

**Files to create:**
- [ ] `internal/matcher/skill_matcher.go` - Skill matching algorithm

**Matching weights:**
```go
type SkillMatchWeights struct {
    Language    float64  // 0.30 - Skills are language-specific
    UseCase     float64  // 0.25 - Match detected use case
    Framework   float64  // 0.20 - Framework compatibility
    ToolSynergy float64  // 0.15 - Works with recommended servers
    Quality     float64  // 0.10 - Quality metrics
}
```

**Use case detection patterns:**
- "debug", "fix", "error", "bug" → debugging skills
- "review", "audit", "check" → code review skills
- "design", "architect", "structure" → architecture skills
- "test", "coverage", "tdd" → testing skills
- "deploy", "ci", "cd", "pipeline" → devops skills

### Phase 4: Project Context Enhancement

**Files to modify:**
- [ ] `pkg/types/project.go` - Add UseCases field
- [ ] `internal/analyzer/intent.go` - Add use case detection patterns

**New field in ProjectContext:**
```go
type ProjectContext struct {
    // ... existing fields ...
    UseCases []UseCase `json:"use_cases"`
}

type UseCase struct {
    Name       string  `json:"name"`
    Confidence float64 `json:"confidence"`
}
```

### Phase 5: CLI Integration

**Files to modify:**
- [ ] `cmd/tools-decision/commands/root.go` - Add skill display/selection

**New output format:**
```
Recommended MCP Servers:
1. PostgreSQL - Store user data
2. Redis - Token caching

Recommended Skills:
1. Go Debugging Assistant - Systematic debugging for Go apps
2. Security Code Review - Security-focused review checklist

Synergies:
• "database-optimization" works best with postgres server
```

### Phase 6: Skill Config Generation

**Files to modify:**
- [ ] `internal/config/generator.go` - Add skill file generation

**Output structure:**
```
.claude/
├── mcp.json
└── commands/
    ├── go-debug.md
    └── security-review.md

.opencode/
├── mcp.json
└── commands/
    ├── go-debug.md
    └── security-review.md
```

---

## Detailed Task Checklist

### Phase 1: Type Definitions
- [ ] Create `pkg/types/skill.go`
  - [ ] Define `Skill` struct
  - [ ] Define `Variable` struct
  - [ ] Define `SkillCompat` struct
  - [ ] Define `SkillCategory` constants
- [ ] Update `pkg/types/recommendation.go`
  - [ ] Add `SkillRecommendation` struct
  - [ ] Add `CombinedRecommendation` struct

### Phase 2: Embedded Skills Registry
- [ ] Create `internal/registry/embedded_skills.go`
  - [ ] Implement `GetEmbeddedSkills()` function
  - [ ] Add debugging skills (go-debug, python-debug, js-debug)
  - [ ] Add code review skills (security-review, performance-review)
  - [ ] Add architecture skills (microservices-design, api-design)
  - [ ] Add testing skills (test-strategy, tdd-workflow)
  - [ ] Add devops skills (ci-cd-setup, docker-best-practices)

### Phase 3: Skill Matcher
- [ ] Create `internal/matcher/skill_matcher.go`
  - [ ] Implement `SkillMatcher` struct
  - [ ] Implement `SkillMatchWeights` configuration
  - [ ] Implement `Match()` method
  - [ ] Implement scoring functions:
    - [ ] `languageScore()` - Language compatibility
    - [ ] `useCaseScore()` - Use case matching
    - [ ] `frameworkScore()` - Framework compatibility
    - [ ] `toolSynergyScore()` - MCP server synergy
    - [ ] `qualityScore()` - Quality metrics
- [ ] Add tests `internal/matcher/skill_matcher_test.go`

### Phase 4: Use Case Detection
- [ ] Update `pkg/types/project.go`
  - [ ] Add `UseCase` struct
  - [ ] Add `UseCases` field to `ProjectContext`
- [ ] Update `internal/analyzer/intent.go`
  - [ ] Add use case patterns to `defaultPatterns()`
  - [ ] Extract use cases in `AnalyzeIdea()`
- [ ] Update project analyzers to detect use cases
  - [ ] Add use case inference from dependencies
  - [ ] Add use case inference from project structure

### Phase 5: CLI Integration
- [ ] Update `cmd/tools-decision/commands/root.go`
  - [ ] Add skill recommendations to analyze output
  - [ ] Add skill selection prompt
  - [ ] Add `--skills` flag for filtering
  - [ ] Update JSON output format

### Phase 6: Config Generation
- [ ] Update `internal/config/generator.go`
  - [ ] Add `GenerateSkillFile()` function
  - [ ] Add skill command file template
  - [ ] Handle variables in skill templates
  - [ ] Update `GenerateConfig()` to include skills

### Phase 7: Documentation & Tests
- [ ] Update `docs/ARCHITECTURE.md` with skills
- [ ] Update `docs/CHECKLIST.md` with skills tasks
- [ ] Add integration tests for skill flow
- [ ] Update slash commands to mention skills

---

## Skill Matching Algorithm Design

### Overview

The skill matcher scores each skill against the project context using weighted factors:

```
TotalScore = Σ (weight × factor_score)

Factors:
  Language:    0.30 × languageScore
  UseCase:     0.25 × useCaseScore
  Framework:   0.20 × frameworkScore
  ToolSynergy: 0.15 × toolSynergyScore
  Quality:     0.10 × qualityScore
```

### Factor Scoring

#### 1. Language Score (0.30 weight)

Skills are highly language-specific. A Go debugging skill is useless for a Python project.

```go
func (m *SkillMatcher) languageScore(projectLangs []types.Language, skillLangs []string) float64 {
    // Check if skill supports "all" languages
    if contains(skillLangs, "all") {
        return 0.7 // Generic skills get moderate score
    }
    
    // Find best matching language
    var maxScore float64
    for _, lang := range projectLangs {
        if contains(skillLangs, strings.ToLower(lang.Name)) {
            if lang.Confidence > maxScore {
                maxScore = lang.Confidence
            }
        }
    }
    return maxScore
}
```

#### 2. Use Case Score (0.25 weight)

Match detected use cases (debugging, review, design, etc.) to skill categories.

```go
var useCaseToCategories = map[string][]string{
    "debugging":     {"debugging", "troubleshooting", "error-handling"},
    "code-review":   {"review", "security", "performance", "best-practices"},
    "architecture":  {"architecture", "design", "patterns", "microservices"},
    "testing":       {"testing", "tdd", "coverage", "quality"},
    "devops":        {"devops", "ci-cd", "deployment", "infrastructure"},
    "documentation": {"documentation", "api-docs", "readme"},
}

func (m *SkillMatcher) useCaseScore(projectUseCases []types.UseCase, skillCategory string) float64 {
    for _, uc := range projectUseCases {
        if categories, ok := useCaseToCategories[uc.Name]; ok {
            if contains(categories, skillCategory) {
                return uc.Confidence
            }
        }
    }
    return 0
}
```

#### 3. Framework Score (0.20 weight)

Some skills are framework-specific (e.g., "Next.js Performance" vs generic "React Performance").

```go
func (m *SkillMatcher) frameworkScore(projectFrameworks []types.Framework, skillFrameworks []string) float64 {
    if contains(skillFrameworks, "all") {
        return 0.5
    }
    
    var maxScore float64
    for _, fw := range projectFrameworks {
        if contains(skillFrameworks, strings.ToLower(fw.Name)) {
            if fw.Confidence > maxScore {
                maxScore = fw.Confidence
            }
        }
    }
    return maxScore
}
```

#### 4. Tool Synergy Score (0.15 weight)

Skills work better when the recommended MCP servers are available.

```go
func (m *SkillMatcher) toolSynergyScore(recommendedServers []types.Recommendation, skill types.Skill) float64 {
    serverIDs := make(map[string]bool)
    for _, rec := range recommendedServers {
        serverIDs[rec.Server.ID] = true
    }
    
    // Check required tools
    requiredMet := 0
    for _, tool := range skill.RequiredTools {
        if serverIDs[tool] {
            requiredMet++
        }
    }
    
    if len(skill.RequiredTools) > 0 && requiredMet < len(skill.RequiredTools) {
        return 0.2 // Penalize if missing required tools
    }
    
    // Bonus for recommended tools
    recommendedMet := 0
    for _, tool := range skill.RecommendedTools {
        if serverIDs[tool] {
            recommendedMet++
        }
    }
    
    if len(skill.RecommendedTools) > 0 {
        return 0.5 + 0.5*(float64(recommendedMet)/float64(len(skill.RecommendedTools)))
    }
    
    return 0.5
}
```

#### 5. Quality Score (0.10 weight)

Reuse existing quality scoring from server matcher.

```go
func (m *SkillMatcher) qualityScore(quality types.Quality) float64 {
    if quality.Score > 0 {
        return quality.Score
    }
    
    score := 0.5
    if quality.Stars > 100 { score += 0.1 }
    if quality.Stars > 500 { score += 0.1 }
    if quality.Maintained { score += 0.2 }
    
    return min(score, 1.0)
}
```

### Minimum Threshold

Skills with score < 0.15 are filtered out (slightly higher than server threshold of 0.1 since skills should be more targeted).

### Synergy Detection

After matching, detect synergies between recommended servers and skills:

```go
type Synergy struct {
    SkillID   string
    ServerID  string
    Reason    string
}

func DetectSynergies(servers []types.Recommendation, skills []types.SkillRecommendation) []Synergy {
    var synergies []Synergy
    
    for _, skill := range skills {
        for _, tool := range skill.Skill.RecommendedTools {
            for _, server := range servers {
                if server.Server.ID == tool {
                    synergies = append(synergies, Synergy{
                        SkillID:  skill.Skill.ID,
                        ServerID: server.Server.ID,
                        Reason:   fmt.Sprintf("'%s' skill works best with %s server", skill.Skill.Name, server.Server.Name),
                    })
                }
            }
        }
    }
    
    return synergies
}
```

---

## Initial Embedded Skills

| ID | Name | Category | Languages | Required Tools |
|----|------|----------|-----------|----------------|
| `go-debug` | Go Debugging Assistant | debugging | go | filesystem, git |
| `python-debug` | Python Debugging Guide | debugging | python | filesystem, git |
| `js-debug` | JavaScript/TS Debugging | debugging | javascript, typescript | filesystem, git |
| `security-review` | Security Code Review | review | all | git, github |
| `performance-review` | Performance Review | review | all | filesystem |
| `api-design` | REST API Design Guide | architecture | all | filesystem |
| `microservices-design` | Microservices Architecture | architecture | go, typescript, python, java | docker, kubernetes |
| `database-optimization` | Database Query Optimization | performance | all | postgres, mysql |
| `test-strategy` | Test Strategy Guide | testing | all | filesystem, git |
| `tdd-workflow` | TDD Workflow | testing | all | filesystem, git |
| `ci-cd-setup` | CI/CD Pipeline Setup | devops | all | github, docker |
| `docker-best-practices` | Docker Best Practices | devops | all | docker |
| `go-project-structure` | Go Project Structure | architecture | go | filesystem |
| `react-patterns` | React Best Practices | architecture | typescript, javascript | filesystem |
| `error-handling` | Error Handling Patterns | best-practices | all | filesystem |

---

## Success Criteria

1. **Type system** compiles with new Skill types
2. **Embedded registry** has 15+ quality skills
3. **Skill matcher** scores skills with >80% relevance accuracy
4. **CLI** displays both servers and skills in recommendations
5. **Config generator** creates skill command files
6. **Tests** pass for all new components
7. **Documentation** updated with skills feature

---

## Timeline Estimate

| Phase | Effort | Dependencies |
|-------|--------|--------------|
| Phase 1: Types | 1 hour | None |
| Phase 2: Registry | 2-3 hours | Phase 1 |
| Phase 3: Matcher | 2 hours | Phase 1 |
| Phase 4: Use Cases | 1 hour | Phase 1 |
| Phase 5: CLI | 2 hours | Phases 2, 3, 4 |
| Phase 6: Config | 1-2 hours | Phase 5 |
| Phase 7: Docs/Tests | 1-2 hours | All phases |

**Total: ~10-13 hours**
