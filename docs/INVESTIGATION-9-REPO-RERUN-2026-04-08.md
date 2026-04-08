# Investigation: 9-Repo Re-Run Snapshot

Date: 2026-04-08
Command used per repo: `tools-decision --json --quiet`
Binary: `/Users/kelvin/Projects/yudgnahk/tools-decision/tools-decision`

## Scope

This document records the latest real-world re-run across the same 9 repositories and summarizes what still needs optimization.

## Re-Run Output (Top Recommendations)

1) `banana-playground`
- Servers: `fetch`, `filesystem`, `github`, `git`, `memory`
- Skills: `go-debug`, `api-design`, `go-project-structure`, `security-review`, `error-handling`

2) `diem_thi_2025`
- Servers: `filesystem`, `github`, `git`, `memory`, `fetch`
- Skills: `security-review`, `tdd-workflow`, `error-handling`, `api-design`, `go-project-structure`

3) `gokit`
- Servers: `postgres`, `mysql`, `sqlite`, `mongodb`, `redis`, `filesystem`, `github`, `git`, `memory`
- Skills: `security-review`, `performance-review`, `tdd-workflow`, `error-handling`, `test-strategy`

4) `goxkey`
- Servers: `filesystem`, `github`, `git`, `memory`, `fetch`
- Skills: `security-review`, `error-handling`, `tdd-workflow`, `api-design`, `rust-desktop-debug`

5) `lightnovel-anime-ai`
- Servers: `filesystem`, `github`, `git`, `memory`
- Skills: `security-review`, `tdd-workflow`, `error-handling`, `performance-review`, `test-strategy`

6) `opencode`
- Servers: `aws`, `filesystem`, `github`, `git`, `memory`, `fetch`
- Skills: `security-review`, `tdd-workflow`, `error-handling`, `api-design`, `test-strategy`

7) `payment-service`
- Servers: `postgres`, `mysql`, `sqlite`, `mongodb`, `redis`, `filesystem`, `github`, `git`, `memory`
- Skills: `security-review`, `performance-review`, `tdd-workflow`, `error-handling`, `test-strategy`

8) `resume`
- Servers: `filesystem`, `github`, `git`, `memory`
- Skills: `latex-authoring-build`

9) `ticket-services`
- Servers: `postgres`, `sqlite`, `mysql`, `mongodb`, `redis`, `filesystem`, `github`, `git`, `memory`
- Skills: `security-review`, `performance-review`, `tdd-workflow`, `error-handling`, `test-strategy`

## What Improved

- `resume` now gets a strongly relevant document skill: `latex-authoring-build`.
- API repos continue to get correct DB/server recommendations.
- Prior high-friction leak regressions (`stripe`, broad `prisma` leakage) remain controlled.

## Remaining Gaps

- **Archetype misses in live repos**:
  - `banana-playground` still does not surface automation tooling (expected `puppeteer`/automation emphasis).
  - `diem_thi_2025` still not reflecting data-processing intent in top skills.
  - `lightnovel-anime-ai` still misses AI/media pipeline skills in top skills.
- **Skill ranking still generic-heavy** across multiple repos:
  - `security-review`, `tdd-workflow`, `error-handling`, `test-strategy` dominate top-5.
  - Domain-specific skills (`browser-automation-scripting`, `etl-data-quality`, `python-ai-media-pipeline`, `go-debug`, `api-design`) are under-ranked in real runs.
- **API repos quality drift**:
  - `gokit`, `payment-service`, `ticket-services` currently skew toward generic review/testing skills instead of domain/API architecture and Go-specific guidance.

## Root Causes Observed

- Real-world archetype extraction is still weaker than fixture-based benchmark contexts.
- Generic skill compatibility (`all`) remains competitive unless strongly counter-weighted.
- Server ranking for automation/AI scenarios still favors core baseline without enough intent boosts.

## Optimization Backlog (for later pass)

1. Strengthen analyzer signals from repository files/scripts (automation/pipeline/media/build clues).
2. Increase archetype confidence carry-over into matcher when multiple weak signals align.
3. Further demote generic skills in strong non-backend archetypes unless use-case confidence is explicit.
4. Add positive score boosts for archetype-specific skills in real contexts (not only synthetic fixtures).
5. Tune automation server promotion (`puppeteer`) when Playwright/Puppeteer/Selenium intent is detected.
6. Re-run 9 repos and update before/after quality table.

## Current Status

- Implementation tracks A-E are in place.
- Final quality target is not yet reached in live rerun; further tuning is required.
