package registry

import (
	"github.com/yudgnahk/tools-decision/pkg/types"
)

// GetEmbeddedSkills returns the bundled AI agent skills registry
// This is a curated list of high-quality skills for AI coding assistants
func GetEmbeddedSkills() []types.Skill {
	return []types.Skill{
		// === Debugging Skills ===
		{
			ID:          "go-debug",
			Name:        "Go Debugging Assistant",
			Slug:        "go-debug",
			Description: "Systematic approach to debugging Go applications",
			Author:      "tools-decision",
			Category:    types.SkillCategoryDebugging,
			Instructions: `## Go Debugging Workflow

### 1. Reproduce the Issue
- Get the exact error message or unexpected behavior
- Identify minimal reproduction steps
- Check if the issue is consistent or intermittent

### 2. Check Common Go Issues
- **Nil pointer dereferences**: Check all pointer/interface values before use
- **Race conditions**: Run with ` + "`-race`" + ` flag: ` + "`go test -race ./...`" + `
- **Resource leaks**: Check for unclosed channels, goroutines, connections
- **Shadowed variables**: Look for ` + "`:=`" + ` accidentally shadowing outer scope

### 3. Use Debugging Tools
- **Delve debugger**: ` + "`dlv debug ./cmd/myapp`" + `
- **Add strategic logging**: Use ` + "`slog`" + ` or ` + "`zerolog`" + ` for structured logs
- **Static analysis**: Run ` + "`go vet`" + ` and ` + "`staticcheck`" + `
- **Stack traces**: Use ` + "`runtime/debug.PrintStack()`" + `

### 4. Analyze Stack Traces
- Follow goroutine traces from bottom to top
- Check for deadlocks (multiple goroutines waiting)
- Look for the originating call site

### 5. Fix and Verify
- Make minimal changes to fix the issue
- Add regression test
- Run full test suite: ` + "`go test ./...`" + `

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Error message or issue description", Required: false},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{"github"},
			Compat: types.SkillCompat{
				Languages:    []string{"go"},
				Frameworks:   []string{"gin", "echo", "fiber", "chi", "cobra"},
				ProjectTypes: []string{"api", "cli", "library"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.90,
			},
			Source: "official",
		},
		{
			ID:          "python-debug",
			Name:        "Python Debugging Guide",
			Slug:        "python-debug",
			Description: "Systematic debugging for Python applications",
			Author:      "tools-decision",
			Category:    types.SkillCategoryDebugging,
			Instructions: `## Python Debugging Workflow

### 1. Understand the Error
- Read the full traceback from bottom to top
- Identify the exception type and message
- Note the file and line number

### 2. Common Python Issues
- **ImportError/ModuleNotFoundError**: Check virtual environment and PYTHONPATH
- **AttributeError**: Verify object type and available attributes
- **TypeError**: Check function signatures and argument types
- **KeyError/IndexError**: Validate data structure contents

### 3. Debugging Tools
- **pdb/ipdb**: ` + "`import pdb; pdb.set_trace()`" + ` or ` + "`breakpoint()`" + `
- **Rich traceback**: ` + "`from rich import traceback; traceback.install()`" + `
- **Print debugging**: Use f-strings with ` + "`{var=}`" + ` syntax
- **Type checking**: Run ` + "`mypy`" + ` for static type analysis

### 4. Interactive Debugging with pdb
- ` + "`n`" + ` - next line
- ` + "`s`" + ` - step into function
- ` + "`c`" + ` - continue to next breakpoint
- ` + "`p <expr>`" + ` - print expression
- ` + "`l`" + ` - list source code

### 5. Async Debugging
- Use ` + "`asyncio.run()`" + ` for testing async code
- Check for missing ` + "`await`" + ` keywords
- Use ` + "`aiomonitor`" + ` for runtime inspection

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Error message or issue description", Required: false},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{"github"},
			Compat: types.SkillCompat{
				Languages:    []string{"python"},
				Frameworks:   []string{"fastapi", "django", "flask", "streamlit"},
				ProjectTypes: []string{"api", "cli", "library", "web_app"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.90,
			},
			Source: "official",
		},
		{
			ID:          "js-debug",
			Name:        "JavaScript/TypeScript Debugging",
			Slug:        "js-debug",
			Description: "Debugging guide for JavaScript and TypeScript applications",
			Author:      "tools-decision",
			Category:    types.SkillCategoryDebugging,
			Instructions: `## JavaScript/TypeScript Debugging Workflow

### 1. Identify the Error Type
- **Runtime errors**: Check browser/Node console for stack traces
- **Type errors** (TS): Read compiler error messages carefully
- **Promise rejections**: Look for unhandled rejection warnings

### 2. Common Issues
- **undefined is not a function**: Check if method exists, verify ` + "`this`" + ` binding
- **Cannot read property of undefined**: Validate object chain with optional chaining ` + "`?.`" + `
- **CORS errors**: Check server headers and request origin
- **Memory leaks**: Use Chrome DevTools Memory tab

### 3. Debugging Tools
- **Browser DevTools**: F12 → Sources → Set breakpoints
- **VS Code debugger**: ` + "`launch.json`" + ` configuration
- **Node.js inspect**: ` + "`node --inspect app.js`" + `
- **console methods**: ` + "`console.table()`" + `, ` + "`console.trace()`" + `, ` + "`console.time()`" + `

### 4. TypeScript-Specific
- Enable ` + "`strict: true`" + ` in tsconfig.json
- Use ` + "`// @ts-expect-error`" + ` to document expected issues
- Check discriminated unions for exhaustiveness

### 5. Async Debugging
- Use ` + "`async/await`" + ` over callbacks for clarity
- Add ` + "`.catch()`" + ` handlers to all promises
- Use ` + "`Promise.allSettled()`" + ` for parallel operations

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Error message or issue description", Required: false},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{"github", "puppeteer"},
			Compat: types.SkillCompat{
				Languages:    []string{"javascript", "typescript"},
				Frameworks:   []string{"nextjs", "react", "vue", "express", "nestjs"},
				ProjectTypes: []string{"web_app", "api", "library"},
				UseCases:     []string{types.UseCaseDebugging},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.90,
			},
			Source: "official",
		},

		// === Code Review Skills ===
		{
			ID:          "security-review",
			Name:        "Security Code Review",
			Slug:        "security-review",
			Description: "Security-focused code review checklist",
			Author:      "tools-decision",
			Category:    types.SkillCategoryReview,
			Instructions: `## Security Code Review Checklist

### 1. Input Validation
- [ ] All user inputs sanitized and validated
- [ ] SQL injection prevention (parameterized queries only)
- [ ] XSS prevention (output encoding, CSP headers)
- [ ] Path traversal prevention (no user input in file paths)
- [ ] Command injection prevention (avoid shell execution)

### 2. Authentication & Authorization
- [ ] Secure password handling (bcrypt, argon2 with proper cost)
- [ ] Session management secure (HttpOnly, Secure, SameSite cookies)
- [ ] Authorization checks on ALL endpoints
- [ ] Rate limiting on auth endpoints
- [ ] Account lockout after failed attempts

### 3. Secrets Management
- [ ] No hardcoded secrets, API keys, or passwords
- [ ] Secrets loaded from environment variables
- [ ] ` + "`.env`" + ` files in ` + "`.gitignore`" + `
- [ ] Sensitive data not logged
- [ ] Secrets rotated regularly

### 4. Data Protection
- [ ] Sensitive data encrypted at rest
- [ ] TLS/HTTPS enforced for all connections
- [ ] PII handled according to regulations (GDPR, etc.)
- [ ] Proper data retention and deletion

### 5. Dependencies
- [ ] No known vulnerable dependencies (` + "`npm audit`" + `, ` + "`safety check`" + `)
- [ ] Dependencies pinned to specific versions
- [ ] Regular dependency updates scheduled

### 6. Error Handling
- [ ] Errors don't leak sensitive information
- [ ] Generic error messages for users
- [ ] Detailed errors only in secure logs

Review this code with security focus: $ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Code to review or specific security concern", Required: false},
			},
			RequiredTools:    []string{"git", "filesystem"},
			RecommendedTools: []string{"github"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseCodeReview, types.UseCaseSecurity},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.95,
			},
			Source: "official",
		},
		{
			ID:          "performance-review",
			Name:        "Performance Code Review",
			Slug:        "performance-review",
			Description: "Performance-focused code review for optimization",
			Author:      "tools-decision",
			Category:    types.SkillCategoryPerformance,
			Instructions: `## Performance Code Review Checklist

### 1. Database Queries
- [ ] No N+1 query problems (use eager loading)
- [ ] Proper indexes on frequently queried columns
- [ ] Pagination for large result sets
- [ ] Query complexity analyzed (EXPLAIN)
- [ ] Connection pooling configured

### 2. Memory & CPU
- [ ] No memory leaks (unclosed resources, event listeners)
- [ ] Large objects processed in streams/chunks
- [ ] Expensive computations cached appropriately
- [ ] No unnecessary object allocations in loops

### 3. Network & I/O
- [ ] Responses compressed (gzip, brotli)
- [ ] Static assets cached with proper headers
- [ ] Parallel requests where possible
- [ ] Timeout configured for external calls
- [ ] Circuit breakers for failing services

### 4. Frontend Performance (if applicable)
- [ ] Bundle size optimized (code splitting)
- [ ] Images optimized (WebP, lazy loading)
- [ ] Critical CSS inlined
- [ ] JavaScript deferred or async
- [ ] Core Web Vitals metrics acceptable

### 5. Caching Strategy
- [ ] Appropriate cache invalidation strategy
- [ ] Cache keys include all relevant parameters
- [ ] TTL values appropriate for data freshness needs
- [ ] Distributed cache for scaled deployments

### 6. Profiling & Monitoring
- [ ] Performance metrics instrumented
- [ ] Slow queries logged
- [ ] Resource usage monitored
- [ ] Alerts configured for degradation

Review for performance issues: $ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Code to review or specific performance concern", Required: false},
			},
			RequiredTools:    []string{"filesystem"},
			RecommendedTools: []string{"postgres", "redis"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseCodeReview, types.UseCasePerformance},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.90,
			},
			Source: "official",
		},

		// === Architecture Skills ===
		{
			ID:          "api-design",
			Name:        "REST API Design Guide",
			Slug:        "api-design",
			Description: "Best practices for designing RESTful APIs",
			Author:      "tools-decision",
			Category:    types.SkillCategoryArchitecture,
			Instructions: `## REST API Design Guide

### 1. Resource Naming
- Use nouns, not verbs: ` + "`/users`" + ` not ` + "`/getUsers`" + `
- Plural for collections: ` + "`/users`" + `, ` + "`/orders`" + `
- Hierarchical for relationships: ` + "`/users/{id}/orders`" + `
- Use kebab-case: ` + "`/user-profiles`" + ` not ` + "`/userProfiles`" + `

### 2. HTTP Methods
| Method | Use Case | Idempotent |
|--------|----------|------------|
| GET | Read resource(s) | Yes |
| POST | Create resource | No |
| PUT | Replace resource | Yes |
| PATCH | Partial update | No |
| DELETE | Remove resource | Yes |

### 3. Status Codes
- ` + "`200`" + ` OK - Successful GET/PUT/PATCH
- ` + "`201`" + ` Created - Successful POST
- ` + "`204`" + ` No Content - Successful DELETE
- ` + "`400`" + ` Bad Request - Invalid input
- ` + "`401`" + ` Unauthorized - Missing/invalid auth
- ` + "`403`" + ` Forbidden - Valid auth, no permission
- ` + "`404`" + ` Not Found - Resource doesn't exist
- ` + "`422`" + ` Unprocessable Entity - Validation error
- ` + "`500`" + ` Internal Server Error - Server fault

### 4. Request/Response Format
` + "```json" + `
// Successful response
{
  "data": { "id": 1, "name": "John" },
  "meta": { "total": 100, "page": 1 }
}

// Error response
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": [{ "field": "email", "issue": "required" }]
  }
}
` + "```" + `

### 5. Pagination
- Use cursor-based for large/real-time data
- Use offset-based for simple cases
- Always include total count and next page link

### 6. Versioning
- URL path: ` + "`/api/v1/users`" + ` (recommended)
- Header: ` + "`Accept: application/vnd.api.v1+json`" + `
- Query param: ` + "`?version=1`" + ` (not recommended)

### 7. Filtering, Sorting, Fields
- Filter: ` + "`?status=active&role=admin`" + `
- Sort: ` + "`?sort=-created_at,name`" + ` (- for desc)
- Fields: ` + "`?fields=id,name,email`" + ` (sparse fieldsets)

Design API for: $ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "API requirements or endpoints to design", Required: false},
			},
			RequiredTools:    []string{"filesystem"},
			RecommendedTools: []string{"fetch"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"api"},
				UseCases:     []string{types.UseCaseArchitecture, types.UseCaseAPIDesign},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.92,
			},
			Source: "official",
		},
		{
			ID:          "microservices-design",
			Name:        "Microservices Architecture",
			Slug:        "microservices-design",
			Description: "Guide for designing microservices architecture",
			Author:      "tools-decision",
			Category:    types.SkillCategoryArchitecture,
			Instructions: `## Microservices Architecture Guide

### 1. Service Boundaries
**Domain-Driven Design Approach:**
- Identify bounded contexts from business domains
- Each service owns its domain logic completely
- Services communicate via well-defined APIs
- Avoid shared databases between services

**Signs of good boundaries:**
- Service can be deployed independently
- Team can work on service without coordination
- Changes are localized to single service
- Service has clear, single responsibility

### 2. Communication Patterns

**Synchronous (Request/Response):**
- REST/HTTP for simple queries
- gRPC for high-performance, internal calls
- Use for: reads, real-time requirements

**Asynchronous (Event-Driven):**
- Message queues (RabbitMQ, SQS) for commands
- Event streaming (Kafka) for event sourcing
- Use for: writes, eventual consistency OK

### 3. Data Management

**Database per Service:**
- Each service has private data store
- No direct database access between services
- Use API calls or events for data sharing

**Patterns:**
- Saga pattern for distributed transactions
- CQRS for read/write separation
- Event sourcing for audit trails

### 4. Resilience Patterns

**Circuit Breaker:**
` + "```go" + `
// Prevent cascade failures
if circuitBreaker.IsOpen() {
    return fallbackResponse()
}
` + "```" + `

**Retry with Backoff:**
- Exponential backoff for transient failures
- Jitter to prevent thundering herd
- Max retry limit

**Timeouts:**
- Always set timeouts on external calls
- Propagate deadline through call chain

### 5. Observability

**Three Pillars:**
1. **Logs**: Structured, correlated by trace ID
2. **Metrics**: RED (Rate, Errors, Duration)
3. **Traces**: Distributed tracing (Jaeger, Zipkin)

### 6. Deployment

**Container Orchestration:**
- Kubernetes for production
- Docker Compose for local development
- Service mesh (Istio, Linkerd) for advanced networking

Design microservices for: $ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "System requirements or services to design", Required: false},
			},
			RequiredTools:    []string{"docker", "filesystem"},
			RecommendedTools: []string{"kubernetes", "postgres", "redis"},
			Compat: types.SkillCompat{
				Languages:    []string{"go", "typescript", "python", "java"},
				Frameworks:   []string{"gin", "nestjs", "fastapi", "spring-boot"},
				ProjectTypes: []string{"api"},
				UseCases:     []string{types.UseCaseArchitecture},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.93,
			},
			Source: "official",
		},
		{
			ID:          "go-project-structure",
			Name:        "Go Project Structure",
			Slug:        "go-project-structure",
			Description: "Idiomatic Go project organization patterns",
			Author:      "tools-decision",
			Category:    types.SkillCategoryArchitecture,
			Instructions: `## Go Project Structure

### Standard Layout
` + "```" + `
myproject/
├── cmd/                    # Entry points
│   └── myapp/
│       └── main.go         # Minimal - just wires things together
├── internal/               # Private application code
│   ├── handler/            # HTTP/gRPC handlers
│   ├── service/            # Business logic
│   ├── repository/         # Data access layer
│   ├── model/              # Domain types
│   └── config/             # Configuration
├── pkg/                    # Public library code (optional)
├── api/                    # OpenAPI specs, protobuf definitions
├── migrations/             # Database migrations
├── scripts/                # Build/deployment scripts
├── Makefile                # Build automation
├── go.mod
└── go.sum
` + "```" + `

### Key Principles

**1. Dependency Direction:**
` + "```" + `
handler → service → repository
   ↓         ↓          ↓
  HTTP    Business    Database
` + "```" + `
- Dependencies flow inward
- Inner layers don't know about outer layers

**2. Interface Location:**
- Define interfaces where they're USED, not implemented
- Repository interface in service package, not repository package

` + "```go" + `
// internal/service/user.go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*model.User, error)
}

type UserService struct {
    repo UserRepository
}
` + "```" + `

**3. Constructor Pattern:**
` + "```go" + `
func NewUserService(repo UserRepository, logger *slog.Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}
` + "```" + `

**4. Error Handling:**
- Wrap errors with context
- Use sentinel errors for expected conditions
- Return errors, don't panic

` + "```go" + `
var ErrUserNotFound = errors.New("user not found")

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get user %s: %w", id, err)
    }
    return user, nil
}
` + "```" + `

### Anti-Patterns to Avoid
- ` + "`utils`" + ` or ` + "`helpers`" + ` packages (be specific)
- Global variables for dependencies
- Putting everything in ` + "`pkg/`" + ` (use ` + "`internal/`" + `)
- Deep nesting (prefer flat structure)

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Project requirements or structure questions", Required: false},
			},
			RequiredTools:    []string{"filesystem"},
			RecommendedTools: []string{"git"},
			Compat: types.SkillCompat{
				Languages:    []string{"go"},
				Frameworks:   []string{"gin", "echo", "fiber", "chi", "cobra"},
				ProjectTypes: []string{"api", "cli", "library"},
				UseCases:     []string{types.UseCaseArchitecture},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.92,
			},
			Source: "official",
		},

		// === Testing Skills ===
		{
			ID:          "test-strategy",
			Name:        "Test Strategy Guide",
			Slug:        "test-strategy",
			Description: "Comprehensive testing strategy for applications",
			Author:      "tools-decision",
			Category:    types.SkillCategoryTesting,
			Instructions: `## Test Strategy Guide

### Testing Pyramid
` + "```" + `
        /\
       /  \     E2E Tests (few)
      /----\    
     /      \   Integration Tests (some)
    /--------\  
   /          \ Unit Tests (many)
  /------------\
` + "```" + `

### 1. Unit Tests
**What to test:**
- Pure functions with business logic
- Edge cases and error handling
- Input validation

**Best practices:**
- One assertion per test (mostly)
- Test behavior, not implementation
- Use table-driven tests for variations
- Mock external dependencies

` + "```go" + `
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        amount   float64
        expected float64
    }{
        {"no discount under 100", 50, 50},
        {"10% discount over 100", 200, 180},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculateDiscount(tt.amount)
            if got != tt.expected {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
` + "```" + `

### 2. Integration Tests
**What to test:**
- Database operations
- API endpoint handlers
- External service interactions

**Best practices:**
- Use test containers for databases
- Clean up test data after each test
- Test happy path and error scenarios

### 3. End-to-End Tests
**What to test:**
- Critical user journeys
- Cross-service workflows
- Production-like environment

**Best practices:**
- Keep E2E tests minimal (expensive to maintain)
- Run in CI but not blocking
- Use realistic test data

### 4. Test Coverage
**Targets:**
- 80% line coverage for critical paths
- 100% coverage for security-sensitive code
- Don't chase coverage for coverage's sake

**What NOT to test:**
- Generated code
- Simple getters/setters
- Third-party library behavior

### 5. Test Organization
` + "```" + `
service/
├── user.go
├── user_test.go          # Unit tests
└── user_integration_test.go  # Integration tests
` + "```" + `

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Testing requirements or specific scenarios", Required: false},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{"docker"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseTesting},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.90,
			},
			Source: "official",
		},
		{
			ID:          "tdd-workflow",
			Name:        "TDD Workflow",
			Slug:        "tdd-workflow",
			Description: "Test-Driven Development workflow guide",
			Author:      "tools-decision",
			Category:    types.SkillCategoryTesting,
			Instructions: `## Test-Driven Development Workflow

### The TDD Cycle (Red-Green-Refactor)

` + "```" + `
    ┌─────────────────┐
    │  1. RED         │
    │  Write failing  │
    │  test           │
    └────────┬────────┘
             │
    ┌────────▼────────┐
    │  2. GREEN       │
    │  Write minimal  │
    │  code to pass   │
    └────────┬────────┘
             │
    ┌────────▼────────┐
    │  3. REFACTOR    │
    │  Improve code   │
    │  quality        │
    └────────┬────────┘
             │
             └──────────► Repeat
` + "```" + `

### Step 1: RED - Write a Failing Test

**Before writing any production code:**
1. Think about the desired behavior
2. Write a test that describes it
3. Run the test - it should FAIL

` + "```go" + `
func TestUserService_Create_ValidUser(t *testing.T) {
    // Arrange
    svc := NewUserService(mockRepo)
    input := CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // Act
    user, err := svc.Create(context.Background(), input)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, input.Email, user.Email)
    assert.NotEmpty(t, user.ID)
}
` + "```" + `

### Step 2: GREEN - Make it Pass

**Write the minimum code to make the test pass:**
- Don't over-engineer
- Don't add features not tested
- It's OK if the code is ugly

` + "```go" + `
func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
    user := &User{
        ID:    uuid.New().String(),
        Email: input.Email,
        Name:  input.Name,
    }
    return user, nil
}
` + "```" + `

### Step 3: REFACTOR - Improve the Code

**Now clean up:**
- Remove duplication
- Improve naming
- Extract methods if needed
- All tests must still pass!

### TDD Best Practices

**Write tests first for:**
- New features
- Bug fixes (write test that reproduces bug first)
- Refactoring (ensure behavior preserved)

**Test naming:**
` + "```" + `
Test[Unit]_[Scenario]_[ExpectedResult]
TestUserService_Create_WithInvalidEmail_ReturnsError
` + "```" + `

**One test, one behavior:**
- Each test should test ONE thing
- If test name has "and", split it

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Feature to implement with TDD", Required: false},
			},
			RequiredTools:    []string{"filesystem"},
			RecommendedTools: []string{"git"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseTesting},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.88,
			},
			Source: "official",
		},

		// === DevOps Skills ===
		{
			ID:          "ci-cd-setup",
			Name:        "CI/CD Pipeline Setup",
			Slug:        "ci-cd-setup",
			Description: "Guide for setting up CI/CD pipelines",
			Author:      "tools-decision",
			Category:    types.SkillCategoryDevOps,
			Instructions: `## CI/CD Pipeline Setup Guide

### Pipeline Stages

` + "```" + `
┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│  Build  │─▶│  Test   │─▶│  Scan   │─▶│ Package │─▶│ Deploy  │
└─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
` + "```" + `

### 1. Build Stage
- Install dependencies
- Compile/transpile code
- Generate artifacts

### 2. Test Stage
- Run unit tests
- Run integration tests
- Generate coverage reports

### 3. Security Scan Stage
- Dependency vulnerability scan
- Static code analysis (SAST)
- Secret detection

### 4. Package Stage
- Build Docker image
- Tag with version/commit SHA
- Push to container registry

### 5. Deploy Stage
- Deploy to staging
- Run smoke tests
- Deploy to production (manual gate)

### GitHub Actions Example

` + "```yaml" + `
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Docker image
        run: docker build -t myapp:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          docker tag myapp:${{ github.sha }} ghcr.io/${{ github.repository }}:${{ github.sha }}
          docker push ghcr.io/${{ github.repository }}:${{ github.sha }}

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy to Kubernetes
        run: kubectl set image deployment/myapp myapp=ghcr.io/${{ github.repository }}:${{ github.sha }}
` + "```" + `

### Best Practices
- Fail fast: run quick tests first
- Cache dependencies
- Use matrix builds for multiple versions
- Require PR reviews before merge
- Use environment protection rules

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "CI/CD requirements or platform specifics", Required: false},
			},
			RequiredTools:    []string{"github", "docker"},
			RecommendedTools: []string{"kubernetes"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseDevOps},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.90,
			},
			Source: "official",
		},
		{
			ID:          "docker-best-practices",
			Name:        "Docker Best Practices",
			Slug:        "docker-best-practices",
			Description: "Best practices for Docker images and containers",
			Author:      "tools-decision",
			Category:    types.SkillCategoryDevOps,
			Instructions: `## Docker Best Practices

### 1. Use Multi-Stage Builds

` + "```dockerfile" + `
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server /server
USER nonroot:nonroot
ENTRYPOINT ["/server"]
` + "```" + `

### 2. Minimize Image Size
- Use Alpine or distroless base images
- Remove package manager caches
- Don't install unnecessary tools
- Use ` + "`.dockerignore`" + `

### 3. Security Best Practices

**Don't run as root:**
` + "```dockerfile" + `
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
` + "```" + `

**Scan for vulnerabilities:**
` + "```bash" + `
docker scout cves myimage:latest
trivy image myimage:latest
` + "```" + `

**Pin base image versions:**
` + "```dockerfile" + `
# Good
FROM node:20.11.0-alpine3.19

# Bad
FROM node:latest
` + "```" + `

### 4. Layer Caching
Order Dockerfile instructions from least to most frequently changing:

` + "```dockerfile" + `
# 1. Base image (rarely changes)
FROM node:20-alpine

# 2. Install dependencies (changes with package.json)
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# 3. Copy source code (changes frequently)
COPY . .

# 4. Build (depends on source)
RUN npm run build
` + "```" + `

### 5. Health Checks
` + "```dockerfile" + `
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
` + "```" + `

### 6. Environment Variables
- Don't bake secrets into images
- Use ARG for build-time, ENV for runtime
- Document required env vars

` + "```dockerfile" + `
ARG APP_VERSION=dev
ENV APP_VERSION=$APP_VERSION
ENV PORT=8080
EXPOSE $PORT
` + "```" + `

### 7. Docker Compose for Development
` + "```yaml" + `
services:
  app:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - .:/app  # Hot reload
    environment:
      - DATABASE_URL=postgres://postgres:postgres@db:5432/app
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: postgres
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
` + "```" + `

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Docker requirements or specific use case", Required: false},
			},
			RequiredTools:    []string{"docker"},
			RecommendedTools: []string{"kubernetes", "github"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseDevOps},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.92,
			},
			Source: "official",
		},

		// === Database Skills ===
		{
			ID:          "database-optimization",
			Name:        "Database Query Optimization",
			Slug:        "database-optimization",
			Description: "Optimize database queries and schema design",
			Author:      "tools-decision",
			Category:    types.SkillCategoryPerformance,
			Instructions: `## Database Query Optimization Guide

### 1. Analyze Query Performance

**PostgreSQL EXPLAIN:**
` + "```sql" + `
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM users WHERE email = 'test@example.com';
` + "```" + `

**Key metrics to watch:**
- Seq Scan vs Index Scan
- Rows estimated vs actual
- Buffer hits vs reads

### 2. Indexing Strategy

**When to add indexes:**
- Columns in WHERE clauses
- Columns in JOIN conditions
- Columns in ORDER BY

**Index types:**
` + "```sql" + `
-- B-tree (default, most cases)
CREATE INDEX idx_users_email ON users(email);

-- Partial index (filtered queries)
CREATE INDEX idx_active_users ON users(email) WHERE active = true;

-- Composite index (multiple columns)
CREATE INDEX idx_users_name_created ON users(name, created_at);

-- GIN index (arrays, JSONB, full-text)
CREATE INDEX idx_users_tags ON users USING GIN(tags);
` + "```" + `

### 3. Query Optimization Patterns

**Avoid SELECT *:**
` + "```sql" + `
-- Bad
SELECT * FROM users WHERE id = 1;

-- Good
SELECT id, name, email FROM users WHERE id = 1;
` + "```" + `

**Use EXISTS instead of IN for large sets:**
` + "```sql" + `
-- Slower
SELECT * FROM orders WHERE user_id IN (SELECT id FROM users WHERE active = true);

-- Faster
SELECT * FROM orders o WHERE EXISTS (
    SELECT 1 FROM users u WHERE u.id = o.user_id AND u.active = true
);
` + "```" + `

**Pagination:**
` + "```sql" + `
-- Offset-based (slow for large offsets)
SELECT * FROM posts ORDER BY created_at DESC LIMIT 20 OFFSET 1000;

-- Cursor-based (fast, consistent)
SELECT * FROM posts 
WHERE created_at < '2024-01-15T10:00:00Z'
ORDER BY created_at DESC 
LIMIT 20;
` + "```" + `

### 4. N+1 Query Problem

**Problem:**
` + "```go" + `
users := db.Query("SELECT * FROM users")
for _, u := range users {
    orders := db.Query("SELECT * FROM orders WHERE user_id = ?", u.ID)  // N queries!
}
` + "```" + `

**Solution - Eager Loading:**
` + "```go" + `
// Single query with JOIN
rows := db.Query(` + "`" + `
    SELECT u.*, o.* 
    FROM users u 
    LEFT JOIN orders o ON o.user_id = u.id
` + "`" + `)
` + "```" + `

### 5. Connection Pooling

` + "```go" + `
db.SetMaxOpenConns(25)                 // Max connections
db.SetMaxIdleConns(10)                 // Keep-alive connections  
db.SetConnMaxLifetime(5 * time.Minute) // Recycle connections
` + "```" + `

### 6. Schema Design Tips

- Use appropriate data types (don't VARCHAR(255) everything)
- Normalize to 3NF, denormalize for read performance
- Use JSONB sparingly (harder to index/query)
- Add created_at, updated_at to all tables

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Query or schema to optimize", Required: false},
			},
			RequiredTools:    []string{},
			RecommendedTools: []string{"postgres", "mysql"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"api", "web_app"},
				UseCases:     []string{types.UseCasePerformance, types.UseCaseDatabaseDesign},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.92,
			},
			Source: "official",
		},

		// === Best Practices Skills ===
		{
			ID:          "error-handling",
			Name:        "Error Handling Patterns",
			Slug:        "error-handling",
			Description: "Best practices for error handling across languages",
			Author:      "tools-decision",
			Category:    types.SkillCategoryBestPractices,
			Instructions: `## Error Handling Patterns

### Principles

1. **Be explicit**: Errors are values, handle them explicitly
2. **Add context**: Wrap errors with information about where they occurred
3. **Fail fast**: Return early on error, don't nest error handling
4. **User vs Developer**: Different error messages for each audience

### Go Error Handling

` + "```go" + `
// Wrap errors with context
func GetUser(id string) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("get user %s: %w", id, err)
    }
    return user, nil
}

// Sentinel errors for expected conditions
var ErrNotFound = errors.New("not found")

// Check error types
if errors.Is(err, ErrNotFound) {
    // Handle not found
}

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
` + "```" + `

### TypeScript Error Handling

` + "```typescript" + `
// Custom error classes
class AppError extends Error {
  constructor(
    message: string,
    public code: string,
    public statusCode: number = 500
  ) {
    super(message);
    this.name = 'AppError';
  }
}

class NotFoundError extends AppError {
  constructor(resource: string) {
    super(` + "`${resource} not found`" + `, 'NOT_FOUND', 404);
  }
}

// Result type pattern
type Result<T, E = Error> = 
  | { ok: true; value: T }
  | { ok: false; error: E };

function parseJSON<T>(json: string): Result<T> {
  try {
    return { ok: true, value: JSON.parse(json) };
  } catch (e) {
    return { ok: false, error: e as Error };
  }
}
` + "```" + `

### Python Error Handling

` + "```python" + `
# Custom exceptions
class AppError(Exception):
    def __init__(self, message: str, code: str):
        self.message = message
        self.code = code
        super().__init__(self.message)

class NotFoundError(AppError):
    def __init__(self, resource: str):
        super().__init__(f"{resource} not found", "NOT_FOUND")

# Context manager for cleanup
@contextmanager
def database_connection():
    conn = create_connection()
    try:
        yield conn
    except Exception as e:
        conn.rollback()
        raise
    finally:
        conn.close()
` + "```" + `

### HTTP Error Responses

` + "```json" + `
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {"field": "email", "message": "Invalid email format"},
      {"field": "age", "message": "Must be positive"}
    ]
  }
}
` + "```" + `

### Anti-Patterns to Avoid

- Silent error swallowing: ` + "`catch (e) {}`" + `
- Generic error messages: "Something went wrong"
- Logging error, then ignoring it
- Throwing strings instead of Error objects
- Not cleaning up resources on error

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Error handling scenario or code to review", Required: false},
			},
			RequiredTools:    []string{"filesystem"},
			RecommendedTools: []string{"git"},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseDebugging, types.UseCaseCodeReview},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.88,
			},
			Source: "official",
		},
		{
			ID:          "code-refactoring",
			Name:        "Code Refactoring Guide",
			Slug:        "code-refactoring",
			Description: "Systematic approach to refactoring code",
			Author:      "tools-decision",
			Category:    types.SkillCategoryBestPractices,
			Instructions: `## Code Refactoring Guide

### When to Refactor

**Code Smells:**
- Long methods (> 20 lines)
- Large classes (> 200 lines)
- Duplicate code
- Deep nesting (> 3 levels)
- Primitive obsession (too many primitives)
- Feature envy (method uses another class's data more than its own)

### Refactoring Safely

1. **Ensure test coverage first**
2. Make small, incremental changes
3. Run tests after each change
4. Commit frequently

### Common Refactoring Patterns

**Extract Method:**
` + "```go" + `
// Before
func processOrder(order Order) error {
    // validate order
    if order.Items == nil || len(order.Items) == 0 {
        return errors.New("empty order")
    }
    if order.Total < 0 {
        return errors.New("invalid total")
    }
    // ... more code
}

// After
func processOrder(order Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    // ... more code
}

func validateOrder(order Order) error {
    if order.Items == nil || len(order.Items) == 0 {
        return errors.New("empty order")
    }
    if order.Total < 0 {
        return errors.New("invalid total")
    }
    return nil
}
` + "```" + `

**Replace Conditional with Polymorphism:**
` + "```go" + `
// Before
func calculatePrice(product Product) float64 {
    switch product.Type {
    case "book":
        return product.Price * 0.9  // 10% discount
    case "electronics":
        return product.Price * 1.1  // 10% markup
    default:
        return product.Price
    }
}

// After
type PricingStrategy interface {
    Calculate(price float64) float64
}

type BookPricing struct{}
func (b BookPricing) Calculate(price float64) float64 {
    return price * 0.9
}

type ElectronicsPricing struct{}
func (e ElectronicsPricing) Calculate(price float64) float64 {
    return price * 1.1
}
` + "```" + `

**Extract Interface:**
` + "```go" + `
// Before: concrete dependency
type UserService struct {
    db *sql.DB
}

// After: interface dependency
type UserRepository interface {
    FindByID(id string) (*User, error)
    Save(user *User) error
}

type UserService struct {
    repo UserRepository
}
` + "```" + `

**Introduce Parameter Object:**
` + "```go" + `
// Before: too many parameters
func createUser(name, email, phone, address, city, country string) error

// After: parameter object
type CreateUserInput struct {
    Name    string
    Email   string
    Phone   string
    Address Address
}

func createUser(input CreateUserInput) error
` + "```" + `

### Refactoring Workflow

1. Identify the code smell
2. Write tests if missing
3. Apply refactoring pattern
4. Run tests
5. Review changes
6. Commit with clear message

$ARGUMENTS`,
			Variables: []types.Variable{
				{Name: "ARGUMENTS", Description: "Code to refactor or specific concern", Required: false},
			},
			RequiredTools:    []string{"filesystem", "git"},
			RecommendedTools: []string{},
			Compat: types.SkillCompat{
				Languages:    []string{"all"},
				Frameworks:   []string{"all"},
				ProjectTypes: []string{"all"},
				UseCases:     []string{types.UseCaseRefactoring, types.UseCaseCodeReview},
			},
			Quality: types.Quality{
				Maintained: true,
				Score:      0.88,
			},
			Source: "official",
		},
	}
}
