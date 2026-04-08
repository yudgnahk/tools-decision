# Plan to Reach "All Good" Recommendations

Last updated: 2026-04-08

## Goal

Move from current real-world quality (`good: 4/9`) to `good: 9/9` across the same 9-repo suite, with stable precision and no return of prior false positives (`stripe`, `prisma` leakage).

## What "Good" Must Mean

A repo is `good` when:

- Top servers reflect the repo's actual stack and workflow.
- Top skills reflect the repo's work type (API, CLI, desktop, data-processing, document, media pipeline, etc.).
- At least 80% of top-5 skills and top-10 servers are clearly actionable for that repo.
- Recommendations are explainable by concrete signals (deps/files/config/project type), not generic fallback only.

## Current Gaps (from manual review)

`poor` repos:
- `diem_thi_2025`: utility/data-processing Go repo receives API/security-heavy skills.
- `goxkey`: Rust desktop app gets backend/API-oriented skills.
- `resume`: LaTeX/document repo gets API/database/testing skills.

`partial` repos:
- `banana-playground`: Go + Playwright/Gemini automation not getting browser automation/tooling intent.
- `lightnovel-anime-ai`: Python AI/media pipeline gets mostly generic API/review skills.

## Root Causes

1) Project archetype detection is too coarse.
- We detect language/framework/services well enough for backend projects, but not repo intent classes like desktop utility, document authoring, automation script, media pipeline.

2) Skill catalog is backend-biased.
- Embedded skills are mostly API/review/devops patterns; weak coverage for desktop, automation, docs, media/AI pipelines.

3) Fallback behavior is still too generic for non-backend repos.
- Core MCP baseline is okay, but skill fallback lacks archetype-sensitive defaults.

4) Missing strong negative evidence.
- No explicit demotion rules for API/database skills in non-API repos (e.g., no server runtime, no API framework).

## Deep Fix Strategy

## Track A: Add Project Archetype Layer (Highest Impact)

Introduce a new "archetype" output in analyzer (multi-label with confidence):

- `api_service`
- `cli_tool`
- `desktop_app`
- `automation_bot`
- `data_processing`
- `ai_content_pipeline`
- `document_authoring`
- `library`

Signals (examples):
- Desktop: `druid`, `tauri`, Cocoa/GTK, app bundle metadata.
- Document: `.tex` volume, LaTeX build scripts, no runtime deps.
- Data-processing: Excel/CSV/Pandas/openpyxl/excelize + batch scripts.
- Automation: Playwright/Puppeteer/Selenium/CDP scripts.
- AI pipeline: genai/openai + moviepy/pydub/ffmpeg + pipeline scripts.

Implementation targets:
- `pkg/types/project.go`: add `Archetypes []ArchetypeSignal`.
- `internal/analyzer/*`: add archetype detectors and merge logic.
- Add confidence calibration rules per archetype.

## Track B: Archetype-Aware Matching Policy

Match servers/skills in two stages:

1. **Eligibility gate** by archetype
- Ex: `api-design` only eligible when `api_service` confidence high.
- Ex: DB-heavy MCPs require `api_service` or explicit DB deps.

2. **Ranking** within eligible set
- Keep current weighted matching, but apply archetype multipliers.

Add negative demotions:
- If archetype = `document_authoring`, demote API/microservices/database skills heavily unless explicit evidence.
- If archetype = `desktop_app`, promote desktop debugging/perf UX skills; demote API architecture defaults.

Implementation targets:
- `internal/matcher/matcher.go`: eligibility + demotion rules.
- `internal/matcher/skill_matcher.go`: same archetype gate/demotion.

## Track C: Expand and Rebalance Skill Catalog

Add missing embedded skills (high priority):

- Rust desktop debugging
- GUI event-loop troubleshooting
- Automation/browser scripting (Playwright/Puppeteer)
- Data processing and ETL quality checks
- Python AI/media pipeline reliability
- Document/LaTeX authoring and build troubleshooting

Then tune compatibility metadata for each new skill:
- narrow `project_types`
- explicit `use_cases`
- language/framework constraints where needed

Implementation targets:
- `internal/registry/embedded_skills.go`

## Track D: MCP Catalog Metadata Tightening for Non-API Context

Current high-friction gating is good (`stripe`, `prisma`, etc.). Extend policy:

- Mark more servers with archetype suitability tags.
- Add `RecommendedArchetypes` and optional `ExcludedArchetypes` metadata.

Example:
- `postgres/mysql/mongodb/redis` strongly preferred for `api_service`, weak for `document_authoring` unless explicit dep signal.

Implementation targets:
- `pkg/types/server.go`: add archetype metadata fields.
- `internal/registry/embedded.go`: annotate top servers.

## Track E: Evaluation Harness and Quality Gates

Create repeatable eval that prevents regressions:

- Keep 9-repo suite as benchmark fixtures.
- For each repo, define expected allowlist/denylist for top servers and top skills.
- Add CI gate: fail if score drops below threshold.

Suggested metrics:
- `GoodRate` target: `>= 9/9`
- `Precision@5 skills` target: `>= 0.8`
- `Precision@10 servers` target: `>= 0.8`
- `GenericSkillRate` target: `< 0.25`

Implementation targets:
- `internal/matcher/*_test.go` + a new evaluation test package
- `docs/REAL-WORLD-MANUAL-REVIEW-2026-04-08.md` follow-up table for before/after

## Concrete Backlog (Execution Order)

1. Add archetype types + analyzer output.
2. Add archetype detectors for desktop/document/automation/data/AI-pipeline.
3. Add matcher eligibility + demotion rules (servers + skills).
4. Expand embedded skills for non-backend archetypes.
5. Add server archetype metadata and apply in ranking.
6. Add benchmark tests with allowlist/denylist per repo.
7. Re-run 9 repos, publish comparison report.

## Expected Outcome by Repo

- `diem_thi_2025` should shift to utility/data-processing skills, drop API-design defaults.
- `goxkey` should shift to Rust/desktop-focused skills.
- `resume` should shift to document/LaTeX-oriented skills and minimal core MCP tools.
- `banana-playground` should add browser automation relevance.
- `lightnovel-anime-ai` should add AI/media pipeline skills.

## Risks and Controls

- Risk: overfitting to 9 repos.
  - Control: keep broad unit fixtures for common archetypes.
- Risk: too strict gating hides useful options.
  - Control: always keep small core baseline and one exploratory slot.
- Risk: metadata maintenance burden.
  - Control: default-safe metadata + targeted overrides for high-impact tools.
