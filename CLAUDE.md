# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go 1.25.3 project (`github.com/dshills/gogame`) in early development stage. The repository is set up with the Specify framework for structured feature development but contains no application code yet.

## Development Commands

### Go Standard Commands

```bash
# Run tests
go test ./...

# Run tests for specific package
go test ./path/to/package

# Run single test
go test -run TestName ./path/to/package

# Build
go build ./...

# Format code
gofmt -w .

# Run linter (if golangci-lint is available)
golangci-lint run
```

### Module Management

```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Specify Framework Integration

This project uses Specify for feature specification and development workflow. The framework provides structured templates and slash commands for managing feature development.

### Key Directories

- `.specify/` - Specify framework configuration
  - `memory/constitution.md` - Project principles and governance (currently template-only)
  - `templates/` - Feature specification, planning, and task templates
- `.claude/commands/` - Slash commands for Specify workflow
  - `/speckit.specify` - Create or update feature specifications
  - `/speckit.plan` - Execute implementation planning
  - `/speckit.tasks` - Generate dependency-ordered tasks
  - `/speckit.implement` - Execute implementation plan
  - `/speckit.clarify` - Identify underspecified areas
  - `/speckit.analyze` - Cross-artifact consistency analysis
  - `/speckit.checklist` - Generate custom checklists
  - `/speckit.constitution` - Create/update project constitution

### Specify Workflow

1. **Specify**: Use `/speckit.specify` to create feature specification from natural language description
2. **Clarify**: Use `/speckit.clarify` to identify and resolve underspecified requirements
3. **Plan**: Use `/speckit.plan` to generate design artifacts and implementation plan
4. **Tasks**: Use `/speckit.tasks` to create actionable, dependency-ordered task list
5. **Analyze**: Use `/speckit.analyze` to verify consistency across artifacts
6. **Implement**: Use `/speckit.implement` to execute the task list

### Project Constitution

The project constitution (`.specify/memory/constitution.md` v1.0.0) defines six core principles:

1. **Go Idiomatic Development** (NON-NEGOTIABLE) - Follow Effective Go, handle errors explicitly, prefer composition
2. **Test-Driven Development** - Mandatory testing with flexible approach, table-driven tests preferred
3. **Concurrent Agent Execution** - Maximize parallel work, complete maximum tasks without stopping
4. **Independent User Stories** - P1/P2/P3 priorities, each story independently testable and valuable
5. **Simplicity and YAGNI** - Start simple, avoid premature optimization, justify complexity
6. **Structured Documentation** - Follow Specify framework artifact structure

**Code Quality Gates**: All code must pass `go build ./...`, `go test ./...`, and `gofmt` before commit.

**Agent Coordination**: Multiple agents should work on different user stories concurrently after completing shared foundation. Agents must work autonomously without constant approval.

## Architecture Patterns

Since this is a new project with no code yet, future development should establish:
- Package structure and organization conventions
- Error handling patterns
- Testing strategy (unit, integration, end-to-end)
- Configuration management approach
- Logging and observability patterns

When implementing features, follow the user story prioritization approach defined in `.specify/templates/spec-template.md`:
- P1 (Priority 1) stories are most critical and should be independently testable MVPs
- Each story should be developable, testable, and deployable independently
- User scenarios must have clear acceptance criteria in Given/When/Then format

## Active Technologies
- File-based (PNG/JPEG image assets loaded from disk, no database required) (001-macos-game-engine)

## Recent Changes
- 001-macos-game-engine: Added Go 1.25.3
