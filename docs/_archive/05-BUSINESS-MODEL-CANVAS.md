# Business Model Canvas: Tools Decision

## Executive Summary

Tools Decision operates on a **freemium SaaS model** with enterprise upsells. The core value proposition—intelligent MCP server discovery and auto-configuration—is delivered through a CLI tool with cloud-backed intelligence.

---

## 1. Customer Segments

### Primary Segments

| Segment | Description | Size Estimate | Willingness to Pay |
|---------|-------------|---------------|-------------------|
| **Individual Developers** | Solo devs using AI coding assistants (Cursor, Claude, Copilot) | 5M+ globally | Low ($0-20/mo) |
| **Developer Teams** | Small to mid-size teams (2-50 devs) building AI-enhanced workflows | 500K+ teams | Medium ($50-500/mo) |
| **Enterprise Engineering** | Large organizations with security/compliance requirements | 50K+ orgs | High ($1K-10K+/mo) |
| **AI/ML Teams** | Teams building AI agents and autonomous systems | 100K+ teams | High ($200-2K/mo) |

### Segment Prioritization (Phase 1)

1. **Individual Developers** - Volume driver, community builders, word-of-mouth
2. **AI/ML Teams** - High engagement, early adopters, feedback loop
3. **Developer Teams** - Revenue bridge to enterprise

---

## 2. Value Propositions

### By Segment

| Segment | Primary Value | Secondary Value |
|---------|---------------|-----------------|
| **Individual Devs** | Save 30+ min/project on tool setup | Discover tools you didn't know existed |
| **Dev Teams** | Consistent tooling across team members | Shared configurations, reduced onboarding |
| **Enterprise** | Security-approved tool catalog, audit trails | Compliance, centralized management |
| **AI/ML Teams** | Optimal tool selection for agent capabilities | Performance analytics, cost optimization |

### Unique Value Proposition

> "The right MCP tools for your project, configured in seconds—not hours."

**Differentiation Matrix:**

| Capability | Tools Decision | Smithery | Glama | Manual |
|------------|----------------|----------|-------|--------|
| Project-aware recommendations | ✅ | ❌ | ❌ | ❌ |
| Auto-configuration | ✅ | ❌ | ❌ | ❌ |
| Multi-registry search | ✅ | ❌ | ❌ | ✅ |
| Quality scoring | ✅ | ⚠️ | ✅ | ❌ |
| Zero-config install | ✅ | ❌ | ❌ | ❌ |

---

## 3. Channels

### Acquisition Channels

| Channel | Type | Cost | Expected Impact |
|---------|------|------|-----------------|
| **GitHub/npm** | Organic | Low | High - Primary discovery |
| **Dev Communities** | Organic | Low | High - Reddit, Discord, HN |
| **Content Marketing** | Owned | Medium | Medium - Blog, tutorials |
| **Conference Talks** | Earned | Medium | Medium - Credibility |
| **IDE Marketplace** | Partner | Low | High - VS Code, JetBrains |
| **AI Tool Partnerships** | Partner | Low | High - Cursor, Claude integrations |

### Distribution Strategy

```
Phase 1: CLI tool (npm/homebrew)
     ↓
Phase 2: IDE extensions (VS Code, JetBrains)
     ↓
Phase 3: CI/CD integrations (GitHub Actions, GitLab CI)
     ↓
Phase 4: Platform APIs (embedded in other tools)
```

---

## 4. Customer Relationships

### Relationship Types by Segment

| Segment | Relationship Type | Touchpoints |
|---------|-------------------|-------------|
| **Free Users** | Self-service | Docs, community forum, GitHub issues |
| **Pro Users** | Automated + Community | Email support, Discord priority |
| **Team Users** | Dedicated support | Slack channel, monthly check-ins |
| **Enterprise** | Account management | Dedicated CSM, SLA, custom onboarding |

### Community Strategy

- **Open Source Core**: CLI tool is open source (MIT license)
- **Community Contributions**: Accept community-submitted tool profiles
- **Ambassador Program**: Power users get early access, swag, recognition
- **Feedback Loop**: Monthly community calls, public roadmap voting

---

## 5. Revenue Streams

### Pricing Tiers

| Tier | Price | Target Segment | Key Features |
|------|-------|----------------|--------------|
| **Free** | $0 | Individual Devs | 10 analyses/mo, basic recommendations |
| **Pro** | $19/mo | Power Users | Unlimited analyses, custom profiles, priority support |
| **Team** | $49/user/mo | Dev Teams | Shared configs, team analytics, SSO |
| **Enterprise** | Custom | Large Orgs | On-prem, audit logs, SLA, custom integrations |

### Revenue Model Details

```
Revenue Mix (Year 2 Target):
├── Pro Subscriptions:     40% ($19 × users)
├── Team Subscriptions:    35% ($49 × seats)
├── Enterprise Contracts:  20% (custom pricing)
└── Marketplace Cut:       5%  (future: tool recommendations)
```

### Unit Economics

| Metric | Target |
|--------|--------|
| **CAC (Customer Acquisition Cost)** | $50 (Pro), $500 (Team), $5K (Enterprise) |
| **LTV (Lifetime Value)** | $400 (Pro), $3K (Team), $50K (Enterprise) |
| **LTV:CAC Ratio** | 8:1 (Pro), 6:1 (Team), 10:1 (Enterprise) |
| **Monthly Churn** | <5% (Pro), <3% (Team), <1% (Enterprise) |
| **Gross Margin** | 85%+ |

---

## 6. Key Resources

### Technical Resources

| Resource | Description | Build vs Buy |
|----------|-------------|--------------|
| **Recommendation Engine** | ML model for tool matching | Build |
| **Project Analyzer** | Static analysis for project context | Build |
| **Registry Aggregator** | Multi-source tool index | Build |
| **Configuration Generator** | Auto-config for various formats | Build |
| **Cloud Infrastructure** | API hosting, analytics | Buy (AWS/GCP) |

### Human Resources (Year 1)

| Role | Count | Focus |
|------|-------|-------|
| **Founding Engineers** | 2-3 | Core product, ML, infrastructure |
| **DevRel** | 1 | Community, content, partnerships |
| **Designer** | 0.5 | CLI UX, branding, docs |

### Intellectual Property

- **Proprietary**: Recommendation algorithm, quality scoring model
- **Open Source**: CLI tool, configuration schemas
- **Data Asset**: Usage patterns, tool effectiveness metrics

---

## 7. Key Activities

### Core Activities

| Activity | Priority | Frequency |
|----------|----------|-----------|
| **Product Development** | Critical | Continuous |
| **Registry Indexing** | Critical | Daily |
| **Quality Scoring Updates** | High | Weekly |
| **Community Engagement** | High | Daily |
| **Content Creation** | Medium | Weekly |
| **Partnership Development** | Medium | Monthly |

### Development Priorities

```
Q1: MVP Launch
    - Project analyzer
    - Basic recommendations
    - CLI tool

Q2: Intelligence Layer
    - ML recommendations
    - Quality scoring
    - Usage analytics

Q3: Team Features
    - Shared configurations
    - Team dashboard
    - SSO integration

Q4: Enterprise & Scale
    - On-prem deployment
    - Audit logging
    - API marketplace
```

---

## 8. Key Partnerships

### Strategic Partnerships

| Partner Type | Examples | Value Exchange |
|--------------|----------|----------------|
| **AI Tool Vendors** | Anthropic, Cursor, Continue | Integration, co-marketing |
| **MCP Registries** | Smithery, Glama | Data access, referrals |
| **IDE Vendors** | VS Code, JetBrains | Marketplace distribution |
| **Cloud Providers** | AWS, GCP, Azure | Credits, co-selling |
| **Tool Authors** | Popular MCP server maintainers | Quality data, promotion |

### Partnership Strategy

```
Phase 1: Data Partnerships
         ├── Registry API access (Smithery, Glama)
         └── Tool author relationships

Phase 2: Distribution Partnerships
         ├── IDE marketplace listings
         └── AI tool integrations

Phase 3: Revenue Partnerships
         ├── Enterprise reseller agreements
         └── Tool marketplace revenue share
```

---

## 9. Cost Structure

### Year 1 Cost Breakdown

| Category | Monthly Cost | % of Total |
|----------|--------------|------------|
| **Personnel** | $50K | 70% |
| **Infrastructure** | $5K | 7% |
| **Tools & Services** | $2K | 3% |
| **Marketing** | $10K | 14% |
| **Legal & Admin** | $3K | 4% |
| **Contingency** | $1.5K | 2% |
| **Total** | ~$71.5K | 100% |

### Cost Scaling Model

```
Users        Infrastructure    Support    Total Marginal
10K          $2K/mo           $0         $0.20/user
100K         $10K/mo          $5K/mo     $0.15/user
1M           $50K/mo          $20K/mo    $0.07/user
```

### Path to Profitability

| Milestone | Users | MRR | Timeline |
|-----------|-------|-----|----------|
| **Break-even** | 5K paid | $75K | Month 18 |
| **Profitable** | 15K paid | $200K | Month 24 |
| **Scale** | 50K paid | $500K | Month 36 |

---

## 10. Financial Projections

### 3-Year Revenue Forecast

| Year | Free Users | Paid Users | ARR | Growth |
|------|------------|------------|-----|--------|
| Y1 | 50K | 2K | $300K | - |
| Y2 | 200K | 10K | $1.5M | 400% |
| Y3 | 500K | 30K | $5M | 233% |

### Funding Requirements

| Stage | Amount | Use of Funds | Timeline |
|-------|--------|--------------|----------|
| **Pre-seed** | $250K | MVP development, initial marketing | Months 1-6 |
| **Seed** | $1.5M | Team growth, scale infrastructure | Months 7-18 |
| **Series A** | $5M | Enterprise features, international | Months 19-36 |

---

## Business Model Canvas Summary

```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│  KEY PARTNERS   │ KEY ACTIVITIES  │VALUE PROPOSITION│   CUSTOMER      │    CUSTOMER     │
│                 │                 │                 │  RELATIONSHIPS  │    SEGMENTS     │
│ • AI tool       │ • Product dev   │ "Right MCP tools│ • Self-service  │ • Individual    │
│   vendors       │ • Registry      │  for your       │   (free)        │   developers    │
│ • MCP registries│   indexing      │  project,       │ • Community     │ • Dev teams     │
│ • IDE vendors   │ • Community     │  configured in  │   support       │ • Enterprise    │
│ • Tool authors  │   engagement    │  seconds"       │ • Dedicated CSM │ • AI/ML teams   │
│                 │                 │                 │   (enterprise)  │                 │
├─────────────────┼─────────────────┤                 ├─────────────────┼─────────────────┤
│  KEY RESOURCES  │                 │ • Project-aware │    CHANNELS     │                 │
│                 │                 │   recommendations│                │                 │
│ • Recommendation│                 │ • Auto-config   │ • npm/homebrew  │                 │
│   engine        │                 │ • Multi-registry│ • GitHub        │                 │
│ • Project       │                 │ • Quality scores│ • IDE markets   │                 │
│   analyzer      │                 │ • Zero friction │ • Dev community │                 │
│ • Engineering   │                 │                 │ • Partnerships  │                 │
│   team          │                 │                 │                 │                 │
├─────────────────┴─────────────────┴─────────────────┴─────────────────┴─────────────────┤
│                                    COST STRUCTURE                                        │
│  Personnel (70%) | Infrastructure (7%) | Marketing (14%) | Tools (3%) | Legal (4%)      │
├──────────────────────────────────────────────────────────────────────────────────────────┤
│                                    REVENUE STREAMS                                       │
│  Pro: $19/mo (40%) | Team: $49/user/mo (35%) | Enterprise: Custom (20%) | Market (5%)   │
└──────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Key Success Metrics

| Metric | Year 1 Target | Year 2 Target |
|--------|---------------|---------------|
| **Total Users** | 50K | 200K |
| **Paid Conversion** | 4% | 5% |
| **MRR** | $25K | $125K |
| **NPS** | 50+ | 60+ |
| **Tool Coverage** | 5K servers | 15K servers |
| **Accuracy (rec quality)** | 80% | 90% |

---

## Risk Mitigation

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| **Registry API changes** | High | Medium | Multi-source, maintain own index |
| **Large competitor entry** | High | Medium | Speed to market, community lock-in |
| **MCP adoption stalls** | Critical | Low | Diversify to other tool protocols |
| **Low conversion rate** | High | Medium | Iterate on value prop, pricing |
| **Technical complexity** | Medium | Medium | Start simple, iterate based on feedback |

---

*Document Version: 1.0*  
*Last Updated: April 2026*  
*Next Review: Quarterly*
