# tools-decision Progress Checklist

Last updated: 2026-04-08

## Completed

- [x] Build passes (`go build`)
- [x] Root analyze flow works for existing projects
- [x] `init` flow works for project idea parsing
- [x] `search` command implemented with local keyword scoring
- [x] `update` command implemented with embedded registry status output
- [x] Added Rust detector (`Cargo.toml`)
- [x] Added Java/Kotlin detector (`pom.xml`, `build.gradle`, `build.gradle.kts`)
- [x] Extended intent patterns for Rust/Java/Kotlin ecosystems
- [x] Added analyzer unit tests (`internal/analyzer/intent_test.go`)
- [x] Added matcher unit tests (`internal/matcher/matcher_test.go`)
- [x] All tests passing (`go test ./...`)
- [x] Added OpenCode command `/tools-decision`
- [x] Kept `/mcp-setup` as legacy alias
- [x] Added Claude Code command `/tools-decision`
- [x] Updated documentation for OpenCode + Claude Code command setup

### Skills Feature

- [x] Created Skill type definition (`pkg/types/skill.go`)
- [x] Added SkillRecommendation and CombinedRecommendation types
- [x] Created embedded skills registry with 15 curated skills
- [x] Implemented SkillMatcher with weighted scoring algorithm
- [x] Added use case detection in intent analyzer
- [x] Integrated skills into CLI analyze and init commands
- [x] Added skill config generation (writes command files)
- [x] Added synergy detection between servers and skills
- [x] Added skill matcher tests (`internal/matcher/skill_matcher_test.go`)
- [x] Created implementation plan (`docs/SKILLS-IMPLEMENTATION.md`)

### Config Command

- [x] Implement `tools-decision config` command end-to-end
- [x] Select servers non-interactively via `--servers` flag
- [x] Select skills non-interactively via `--skills` flag
- [x] Support `all` and `recommended` special values
- [x] Generate and write config directly for selected format
- [x] Support JSON output for config command
- [x] Support `--dry-run` for preview without writing
- [x] Support `--list` to show available servers and skills

### Live Registry Update

- [x] Wire `Registry` + `Cache` into `update` command
- [x] Implement Smithery fetcher (`internal/registry/fetchers.go`)
- [x] Implement Glama fetcher structure (API not confirmed)
- [x] Merge/dedupe embedded + fetched sources with source metadata
- [x] Added fetcher unit tests (`internal/registry/fetchers_test.go`)

## Remaining (Medium Priority)

- [ ] Improve search ranking quality
  - [ ] Token normalization/synonyms
  - [ ] Better field weighting
  - [ ] Stable sorting for ties

- [ ] Add end-to-end CLI tests
  - [ ] Analyze existing project flow
  - [ ] Init idea flow
  - [ ] Search flow
  - [ ] Config generation flow
  - [ ] Update flow

- [ ] Expand embedded skills registry
  - [ ] Add more language-specific debugging skills
  - [ ] Add framework-specific skills (React, Vue, Django, etc.)
  - [ ] Add cloud/infrastructure skills (AWS, GCP, Azure)

## Remaining (Low Priority)

- [ ] Add more detectors beyond current set (as needed)
- [ ] Expand embedded registry coverage with stricter quality metadata
- [ ] Fetch skills from Smithery (131K+ skills available)
- [ ] Implement official MCP registry fetcher (GitHub parsing)

## Quick Resume Commands

```bash
go test ./...
go build -o tools-decision ./cmd/tools-decision
./tools-decision --help
./tools-decision init "REST API with PostgreSQL and Redis"
./tools-decision search database --json
./tools-decision update
./tools-decision config --list
./tools-decision config --servers postgres,redis --skills go-debug --dry-run
```

## Integration Smoke Examples

OpenCode:

```text
/tools-decision I want to create a API service with Golang, Gorm and Gin, about the micro-services with auth service first
```

Claude Code:

```text
/tools-decision I want to create a API service with Golang, Gorm and Gin, about the micro-services with auth service first
```

## Config Command Examples

```bash
# List all available servers and skills
./tools-decision config --list
./tools-decision config --list --json

# Generate config for specific servers
./tools-decision config --servers postgres,redis,github

# Generate config with skills
./tools-decision config --skills go-debug,security-review

# Preview without writing
./tools-decision config --servers postgres --dry-run

# Use all recommended
./tools-decision config --servers recommended --skills recommended
```

## Skills Feature Examples

```bash
# Analyze project with skills
./tools-decision

# Output includes both MCP servers and skills:
# Recommended MCP Servers:
#   1. PostgreSQL (score: 92%)
#   2. Redis (score: 85%)
#
# Recommended Skills:
#   1. Go Debugging Assistant (score: 88%)
#   2. API Design Guide (score: 82%)
#
# Synergies:
#   • 'Go Debugging Assistant' skill works best with Git server
```
