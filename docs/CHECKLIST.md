# tools-decision Progress Checklist

Last updated: 2026-04-07

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

## Remaining (High Priority)

- [ ] Implement `tools-decision config` command end-to-end
  - [ ] Select servers non-interactively or via flags
  - [ ] Generate and write config directly for selected format
  - [ ] Support JSON output for config command

- [ ] Implement live registry update flow
  - [ ] Wire `Registry` + `Cache` into `update` command
  - [ ] Implement Smithery fetcher (`internal/registry/fetchers.go`)
  - [ ] Implement Glama fetcher (`internal/registry/fetchers.go`)
  - [ ] Implement official registry fetcher (`internal/registry/fetchers.go`)
  - [ ] Merge/dedupe embedded + fetched sources with source metadata

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

- [ ] Improve update command UX
  - [ ] Sort category output alphabetically
  - [ ] Optional summary views (`--source`, `--category`)

## Remaining (Low Priority)

- [ ] Add more detectors beyond current set (as needed)
- [ ] Expand embedded registry coverage with stricter quality metadata

## Quick Resume Commands

```bash
go test ./...
go build -o tools-decision ./cmd/tools-decision
./tools-decision --help
./tools-decision init "REST API with PostgreSQL and Redis"
./tools-decision search database --json
./tools-decision update
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
