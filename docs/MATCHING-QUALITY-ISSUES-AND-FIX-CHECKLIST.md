# Matching Quality Issues and Fix Checklist

Last updated: 2026-04-08

## Context

Real-world testing across 9 repositories showed over-recommendation of generic/high-quality MCP servers and weak intent matching for some projects. Example: `banana-playground` received `stripe` and `prisma` despite no payment or Prisma usage in the repo.

## Observed Issues

### 1) False positives from quality-only scoring

- Symptom: servers like `stripe` appear with reason only `High quality and well-maintained`.
- Impact: irrelevant tools are presented as top-10 recommendations.
- Evidence: `banana-playground` output had `stripe` with no `matched_on` signals.

### 2) Weak JS/TS compatibility causes ORM false positive

- Symptom: `prisma` appears in projects with any JS footprint even without Prisma files/deps.
- Impact: ORM tool recommended without actual ORM intent.
- Evidence: `banana-playground` had `prisma` because of JS compatibility match.

### 3) Missing project-intent gating

- Symptom: recommendations can pass threshold from compatibility + quality even when service/category intent is absent.
- Impact: top results become repetitive across unrelated repos.

### 4) Analyzer likely contaminated by non-source directories

- Symptom: browser artifacts / large folders (e.g. `chromium-data`, `node_modules`) can influence language/framework detection if not excluded aggressively.
- Impact: noisy context, less precise matching.

### 5) No penalty when required env vars imply high setup cost without relevance

- Symptom: sensitive/external integrations (`stripe`, `slack`) can rank without corresponding repo signals.
- Impact: poor recommendation precision and higher user friction.

## Root Causes (Current Logic)

- `internal/matcher/matcher.go`
  - Quality contributes directly even when no intent match exists.
  - Minimum threshold (`score > 0.1`) allows quality-only candidates to pass.
- `internal/registry/embedded.go`
  - Some tool compat metadata is broad (`all`, generic JS compatibility), enabling weak matches.
- Analyzer/detector pipeline
  - Needs stricter ignore rules for generated/cache/vendor directories.

## Fix Plan (Prioritized)

## Phase 1: Precision Guardrails (High Priority)

Status: partially completed

- [x] Require at least one non-quality signal for inclusion.
  - Rule: candidate must match **service/category/framework/language-with-high-confidence**; quality alone cannot include.
- [x] Add a hard penalty for candidates with zero `matched_on` terms.
- [ ] Raise minimum threshold or make threshold dynamic by project confidence.
- [x] Do not include payment/comms/cloud integrations unless intent signals exist.
  - Examples: `stripe`, `slack`, `aws` require explicit service/tag hits.

## Phase 2: Better Signal Weighting (High Priority)

Status: mostly completed

- [x] Re-balance weights in `MatchWeights`:
  - increase service/category importance
  - reduce quality contribution when intent is weak
- [x] Add signal classes:
  - Strong: explicit dependencies/files/configs (`package.json`, `go.mod`, `schema.prisma`, env keys)
  - Medium: framework/language
  - Weak: generic compatibility (`all`)
- [x] Sort ties stably by strong-signal count then score.

## Phase 3: Analyzer Noise Reduction (Medium Priority)

Status: completed

- [x] Exclude noisy directories from detection and intent analysis:
  - `node_modules`, `.git`, `chromium-data`, `dist`, `build`, `coverage`, `.next`, `vendor`, cache dirs
- [x] Ensure only project source/config files drive context.

## Phase 4: Registry Metadata Tightening (Medium Priority)

Status: completed

- [x] Refine compat metadata for sensitive integrations.
  - Example: `stripe` should not be effectively generic.
- [x] Add `intent_required`/`requires_explicit_signal` metadata for high-friction servers.
- [x] Improve categories/tags for stronger service matching.

## Phase 5: Regression Tests (High Priority)

- [x] Add matcher tests preventing known false positives:
  - `banana-playground` fixture should NOT recommend `stripe`.
  - JS-only project without Prisma should NOT recommend `prisma`.
- [ ] Add positive controls:
  - payment service fixture should recommend `stripe`.
  - Prisma project fixture should recommend `prisma`.
- [ ] Add end-to-end snapshot tests for at least 5 real repo fixtures.

## Acceptance Criteria

- [ ] In non-payment repos, `stripe` is absent from top-10 unless explicit payment signals exist.
- [ ] In non-Prisma repos, `prisma` is absent from top-10 unless Prisma signals exist.
- [ ] At least 30% reduction in repeated generic recommendations across the 9-repo suite.
- [ ] Every recommended server has at least one meaningful match reason beyond quality.
- [ ] Real-world suite passes with improved relevance and no regressions in build/tests.

## Implementation Checklist (Actionable)

- [x] Update matcher inclusion rule in `internal/matcher/matcher.go`.
- [x] Add `isQualityOnly` guard + exclusion path.
- [ ] Reweight `DefaultWeights()` and adjust threshold.
- [x] Add explicit-signal gating for sensitive servers (config/tag-based).
- [ ] Tighten analyzer ignore paths.
- [ ] Update embedded metadata for `stripe`, `prisma`, and other high-friction integrations.
- [x] Add unit tests for false-positive/true-positive scenarios.
- [x] Re-run real-world 9-repo suite and record before/after comparison.

## Suggested Metrics to Track

- Precision@10 on real-world suite
- Count of quality-only recommendations
- Recommendation diversity across repos
- % of recommendations with explicit matched signals
