# All-Good Recommendation Implementation Checklist

Last updated: 2026-04-08
Source plan: `docs/ALL-GOOD-RECOMMENDATION-PLAN.md`

## Track A: Archetype Layer (Highest Impact)

- [x] Add archetype model to project context (`pkg/types/project.go`)
- [x] Wire analyzer output to return archetype signals (`internal/analyzer/analyzer.go`)
- [x] Add initial archetype detection heuristics for desktop/document/automation/data/AI (`internal/analyzer/archetype.go`)
- [ ] Calibrate confidence rules for all archetypes against 9-repo suite
- [ ] Add detector-level archetype evidence where available

## Track B: Archetype-Aware Matching Policy

- [x] Add server eligibility gate by archetype (`internal/matcher/matcher.go`)
- [x] Add skill eligibility gate by archetype (`internal/matcher/skill_matcher.go`)
- [x] Add negative demotions for non-API repos (document/desktop/data)
- [x] Add archetype multipliers during ranking
- [x] Add matcher tests for gate + demotion behavior

## Track C: Skill Catalog Rebalance

- [x] Add Rust desktop debugging skill
- [x] Add GUI event-loop troubleshooting skill
- [x] Add automation/browser scripting skill set (Playwright/Puppeteer)
- [x] Add data-processing/ETL quality skill
- [x] Add Python AI/media pipeline reliability skill
- [x] Add document/LaTeX authoring troubleshooting skill
- [x] Tighten compatibility metadata for all new skills

## Track D: MCP Catalog Archetype Metadata

- [x] Add archetype metadata fields on server type (`pkg/types/server.go`)
- [x] Annotate embedded servers with recommended/excluded archetypes (`internal/registry/embedded.go`)
- [x] Apply archetype metadata in server ranking/gating logic
- [x] Add tests for archetype metadata enforcement

## Track E: Evaluation Harness and Quality Gates

- [x] Add benchmark fixture harness for 9-repo suite
- [x] Define allowlist/denylist expectations per repo (skills + servers)
- [x] Add CI gate for quality thresholds (`GoodRate`, `Precision@5`, `Precision@10`, `GenericSkillRate`)
- [x] Publish before/after report in docs
- [x] Tighten baseline gates (intermediate): `GoodRate>=0.65`, `Precision@5>=0.80`, `Precision@10>=0.90`, `GenericSkillRate<0.40`
- [x] Add generic-skill filtering for strong non-backend archetypes in skill matcher
- [x] Re-run real 9-repo suite with current binary (`tools-decision --json --quiet`)
- [ ] Tighten baseline gates to final target thresholds after additional matcher tuning

## Latest 9-Repo Re-Run Findings (2026-04-08)

- [x] `resume` now receives document-first skill output (`latex-authoring-build`)
- [ ] `banana-playground`: still missing automation server/tooling (`puppeteer`) and still API-heavy skills
- [ ] `diem_thi_2025`: still dominated by generic review/testing/API skills; needs data-processing skills
- [ ] `lightnovel-anime-ai`: still missing AI/media pipeline skill in top results; too generic
- [ ] API repos (`gokit`, `payment-service`, `ticket-services`) now over-index on generic review/test skills instead of domain skills (`go-debug`, `api-design`, `microservices-design`)

## In Progress Now

- [x] Start Track A implementation with analyzer archetype output and baseline tests
- [x] Implement first-pass Track B archetype gating/demotion in server + skill matchers
- [x] Implement first-pass Track C skill catalog rebalance for non-backend archetypes
- [x] Implement Track D server archetype metadata and enforce it in matcher policy
- [x] Implement Track E benchmark test harness + quality gates
- [x] Run and stabilize `go test ./...` after matcher integration work
