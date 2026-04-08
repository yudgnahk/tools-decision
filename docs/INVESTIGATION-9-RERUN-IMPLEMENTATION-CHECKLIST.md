# Investigation 9 Re-Run Implementation Checklist

Date: 2026-04-08
Source: `docs/INVESTIGATION-9-REPO-RERUN-2026-04-08.md`

## Phase 1: Analyzer Signal Strengthening

- [x] Add stronger automation/AI dependency detection in JavaScript detector (`internal/analyzer/javascript.go`)
- [x] Add data/automation/AI dependency detection in Go detector (`internal/analyzer/go.go`)
- [x] Parse `pyproject.toml` dependencies to recover Python service signals (`internal/analyzer/python.go`)
- [x] Add archetype carry-over boosts when multiple weak signals align (`internal/analyzer/archetype.go`)
- [x] Add analyzer regression tests for JS automation, Go excelize, Python pyproject (`internal/analyzer/detector_signal_test.go`)

## Phase 2: Archetype-to-Matcher Confidence Carry-Over

- [x] Tune skill gating thresholds for weak-but-aligned archetype signals (`internal/matcher/skill_matcher.go`)
- [x] Add archetype-specific positive boosts in real contexts (automation/data/AI/API)
- [x] Add tests ensuring domain skills outrank generic skills for rerun gap repos (`internal/matcher/skill_matcher_test.go`)

## Phase 3: Generic Skill Pressure Reduction

- [x] Further demote `all`-compat generic skills in strong non-backend archetypes
- [x] Reduce generic dominance in API repos when domain/API-specific skills are available
- [x] Add benchmark checks for generic-heavy regressions in top-5 skill output

## Phase 4: Server Promotion and Ranking

- [x] Promote `puppeteer` when Playwright/Puppeteer/Selenium intent exists
- [x] Recalibrate server ranking for automation and AI/media archetypes
- [x] Add matcher tests for automation server promotion behavior

## Phase 5: Validation and Documentation

- [x] Re-run 9 repos with `tools-decision --json --quiet`
- [x] Publish before/after diff table for servers and skills
- [ ] Tighten quality gates to final thresholds after rerun passes

## Latest Re-Run Delta (after tuning)

Command: `tools-decision --json --quiet` (using rebuilt local binary)

| Repo | Before (key issue) | After (top change) | Status |
|---|---|---|---|
| `banana-playground` | Missing automation server/skill emphasis | `puppeteer` now top server; `js-debug` + `browser-automation-scripting` included | Improved |
| `diem_thi_2025` | Generic/API-heavy skill stack | `etl-data-quality` now top skill | Improved |
| `gokit` | Generic review/test skills dominating | Still generic-heavy skill top-5 | Unresolved |
| `goxkey` | Desktop/domain skill under-ranked | `rust-desktop-debug` appears in top-5, but generic skills still dominate | Partial |
| `lightnovel-anime-ai` | Missing AI/media pipeline skills | `python-ai-media-pipeline` now appears in top-2 | Improved |
| `opencode` | Generic-heavy skills | Largely unchanged (still generic-heavy in top-5) | Unresolved |
| `payment-service` | API repos skewed to generic review/testing | Still generic-heavy skill top-5 | Unresolved |
| `resume` | Needed document-first skill | `latex-authoring-build` remains top recommendation | Good |
| `ticket-services` | API repos skewed to generic review/testing | Still generic-heavy skill top-5 | Unresolved |

## Remaining Priority

1. Add benchmark fixture checks that explicitly fail when generic skills occupy too many top-5 slots on API archetypes.
2. Re-run all 9 repos again after API-skill rebalance and then tighten final quality gates.

## API Rebalance Spot Check (latest)

- `gokit`: now domain-first (`go-debug`, `go-project-structure`, `microservices-design`, `api-design`, `database-optimization`)
- `payment-service`: now domain-first (`go-debug`, `go-project-structure`, `microservices-design`, `api-design`, `database-optimization`)
- `ticket-services`: now domain-first (`go-debug`, `js-debug`, `go-project-structure`, `api-design`, `microservices-design`)
- `opencode`: still mixed, but less generic-heavy than earlier runs

## Latest Benchmark Gate Result

- `go test ./internal/evaluation -run TestNineRepoBenchmarkQualityGate -v` now passes
- Metrics: `GoodRate=0.78`, `Precision@5=0.87`, `Precision@10=0.98`, `GenericSkillRate=0.19`
- Key non-API recovery:
  - `banana-playground` skill precision improved to `0.80` with `go-debug` and `js-debug` returning to top-5
  - `lightnovel-anime-ai` skill precision improved to `0.80` with `python-debug` in top-5
