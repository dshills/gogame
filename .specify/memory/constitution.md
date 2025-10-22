<!--
  ============================================================================
  SYNC IMPACT REPORT
  ============================================================================
  Version Change: [NONE - Initial constitution] → 1.0.0

  Modified Principles: N/A (initial creation)

  Added Sections:
  - Core Principles (6 principles)
  - Development Workflow
  - Agent Coordination
  - Governance

  Removed Sections: N/A

  Templates Requiring Updates:
  - ✅ plan-template.md (Constitution Check section verified - compatible)
  - ✅ spec-template.md (User story prioritization aligns with P1/P2/P3)
  - ✅ tasks-template.md (Parallel task execution aligns with concurrency principle)
  - ✅ All command files verified (no agent-specific references)

  Follow-up TODOs: None

  Rationale for Initial Version 1.0.0:
  This is the first ratified constitution for the gogame project, establishing
  foundational principles for Go development, concurrent agent coordination,
  and autonomous execution aligned with Specify framework workflows.
  ============================================================================
-->

# gogame Constitution

## Core Principles

### I. Go Idiomatic Development (NON-NEGOTIABLE)

All code MUST follow Go idioms and best practices as defined by Effective Go and the Go Code Review Comments:

- Use `gofmt` and `goimports` for consistent formatting
- Follow Go naming conventions (MixedCaps, not snake_case)
- Prefer composition over inheritance
- Handle errors explicitly; never ignore returned errors
- Write clear, simple code that favors readability over cleverness
- Use interfaces to define behavior, not concrete types
- Keep packages focused with clear, single responsibilities

**Rationale**: Go's simplicity and consistency are core strengths. Idiomatic code is maintainable, reviewable, and integrates seamlessly with the ecosystem.

### II. Test-Driven Development

Testing is mandatory but flexible in approach:

- Unit tests MUST cover core game logic and business rules
- Integration tests SHOULD cover component interactions
- Table-driven tests are PREFERRED for comprehensive coverage
- Use `go test ./...` as the primary testing command
- Benchmark performance-critical paths with `go test -bench`
- Tests MUST be runnable independently and in any order
- When user stories explicitly request TDD: write tests first, ensure they fail, then implement

**Rationale**: Go's excellent testing tooling makes TDD natural. Tests serve as documentation and safety nets for refactoring.

### III. Concurrent Agent Execution

Development workflow MUST maximize parallel execution and autonomous operation:

- Use concurrent agents whenever tasks are independent
- Agents MUST complete maximum work without stopping for approval
- Minimize sequential dependencies in task planning
- Design features as independently testable, deployable slices
- Parallelize file operations, tests, and research activities
- Document blocking dependencies explicitly in task lists

**Rationale**: User requirement for "concurrent agents as much as possible" and "complete as much work as possible without stopping" drives this principle.

### IV. Independent User Stories

Every feature MUST be decomposed into prioritized, independently valuable user stories:

- Assign clear priorities: P1 (MVP), P2, P3, etc.
- Each story MUST be developable independently
- Each story MUST be testable independently
- Each story MUST deliver standalone value
- P1 stories define the Minimum Viable Product
- Use Given/When/Then format for acceptance criteria

**Rationale**: Independent stories enable parallel development, incremental delivery, and rapid validation of core value.

### V. Simplicity and YAGNI

Start simple and add complexity only when justified:

- Avoid premature optimization
- Reject over-engineering and speculative features
- Prefer standard library over external dependencies when possible
- Keep package dependencies minimal and well-justified
- Document complexity when unavoidable with clear rationale
- Challenge every abstraction: "Do we need this now?"

**Rationale**: Go's philosophy emphasizes simplicity. Complexity must earn its place through demonstrated need.

### VI. Structured Documentation

Project documentation MUST follow the Specify framework structure:

- Feature specifications in `specs/[###-feature]/spec.md`
- Implementation plans in `specs/[###-feature]/plan.md`
- Task lists in `specs/[###-feature]/tasks.md`
- Research artifacts in `specs/[###-feature]/research.md`
- Runtime guidance in `CLAUDE.md` for AI assistants
- Use Markdown for all documentation
- Link artifacts bidirectionally for traceability

**Rationale**: Consistent structure enables efficient collaboration between human developers and AI agents.

## Development Workflow

### Code Quality Gates

Before any commit:

1. Code MUST pass `go build ./...`
2. Code MUST pass `go test ./...`
3. Code MUST pass `gofmt -w .` (no formatting changes)
4. Code SHOULD pass `golangci-lint run` if configured
5. All errors MUST be handled explicitly

### Feature Development Cycle

1. **Specify** (`/speckit.specify`): Create feature specification with prioritized user stories
2. **Clarify** (`/speckit.clarify`): Resolve underspecified requirements
3. **Plan** (`/speckit.plan`): Generate technical design and architecture decisions
4. **Tasks** (`/speckit.tasks`): Create dependency-ordered, parallelizable task list
5. **Analyze** (`/speckit.analyze`): Verify cross-artifact consistency
6. **Implement** (`/speckit.implement`): Execute tasks with maximum concurrency
7. **Validate**: Test each user story independently before proceeding

### Commit Standards

- Commit messages MUST follow conventional commits format
- Commits SHOULD be atomic (one logical change per commit)
- Commit after completing each task or logical task group
- Include issue/story reference where applicable

### Module Management

- Keep `go.mod` dependencies minimal and justified
- Run `go mod tidy` before committing dependency changes
- Document why each external dependency is needed
- Prefer well-maintained, widely-used libraries

## Agent Coordination

When multiple agents (human or AI) work on this codebase:

- **Parallel Tasks**: Agents SHOULD work on different user stories concurrently
- **Shared Foundation**: Complete foundational phase before story work begins
- **Independence**: Each agent owns complete user stories, not partial implementations
- **Communication**: Update task status immediately upon completion
- **Conflict Avoidance**: Different files per agent; use `[P]` markers in task lists
- **Autonomy**: Agents MUST make decisions and proceed without constant approval

## Governance

### Constitutional Authority

- This constitution supersedes informal practices and preferences
- All feature specifications, plans, and tasks MUST comply with these principles
- Code reviews MUST verify constitutional compliance
- Complexity violations REQUIRE explicit justification in plan.md Complexity Tracking section

### Amendment Process

- Amendments REQUIRE documented rationale and impact analysis
- Version increments follow semantic versioning:
  - **MAJOR**: Breaking changes to core principles
  - **MINOR**: New principles or substantial expansions
  - **PATCH**: Clarifications, wording fixes, non-semantic refinements
- Update dependent templates within same commit as constitution changes
- Generate Sync Impact Report for each amendment

### Compliance

- Use `CLAUDE.md` for runtime AI agent guidance
- Reference this constitution in all `/speckit.*` command executions
- Challenge non-compliance immediately in code reviews
- Maintain traceability from user stories → tasks → code

**Version**: 1.0.0 | **Ratified**: 2025-10-22 | **Last Amended**: 2025-10-22
