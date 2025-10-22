# Implementation Plan: macOS Game Engine

**Branch**: `001-macos-game-engine` | **Date**: 2025-10-22 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-macos-game-engine/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a 2D game engine for macOS that enables Go developers to create arcade-style games and platformers with minimal code. The engine provides core capabilities: scene rendering at 60 FPS, entity management with game loop, input handling, asset loading, and basic collision detection. Technical approach will use native macOS graphics APIs with Go bindings for optimal performance.

## Technical Context

**Language/Version**: Go 1.25.3
**Primary Dependencies**: NEEDS CLARIFICATION (macOS graphics API choice - Metal, OpenGL, or SDL2 bindings)
**Storage**: File-based (PNG/JPEG image assets loaded from disk, no database required)
**Testing**: Go standard testing (`go test`), table-driven tests, benchmarks for performance validation
**Target Platform**: macOS 12.0+ (Monterey or newer), M1/Intel architecture support
**Project Type**: Single project (game engine library + example games)
**Performance Goals**: 60 FPS with 100 sprites, <16ms input latency, <100ms asset loading
**Constraints**: Single-threaded game loop, no audio support, 2D only, memory stable over 1 hour
**Scale/Scope**: Support indie game development, 100-1000 entities per scene, simple to medium game complexity

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Initial Check (Pre-Research)

#### ✅ I. Go Idiomatic Development (NON-NEGOTIABLE)
- **Status**: PASS
- **Compliance**: Using Go 1.25.3, will follow gofmt/goimports, MixedCaps naming, explicit error handling
- **Validation**: All code will pass `go build` and `gofmt -w .` before commit

#### ✅ II. Test-Driven Development
- **Status**: PASS
- **Compliance**: Will use `go test ./...`, table-driven tests for game logic, benchmarks for performance
- **Validation**: Unit tests for core engine components, integration tests for complete game scenarios

#### ✅ III. Concurrent Agent Execution
- **Status**: PASS
- **Compliance**: 5 independent user stories (P1-P5) enable parallel development
- **Validation**: Each priority can be developed by separate agents after foundation complete

#### ✅ IV. Independent User Stories
- **Status**: PASS
- **Compliance**: Spec defines P1 (rendering) as MVP, P2-P5 as incremental value additions
- **Validation**: Each story has independent acceptance criteria and test plans

#### ✅ V. Simplicity and YAGNI
- **Status**: PASS with justification required
- **Compliance**: Prefer Go standard library, minimize dependencies
- **Justification Needed**: Graphics API dependency unavoidable - will research minimal binding approach
- **Note**: External dependency for macOS graphics is necessary; will justify choice in research.md

#### ✅ VI. Structured Documentation
- **Status**: PASS
- **Compliance**: Following Specify structure (spec.md, plan.md, research.md, data-model.md, contracts/, tasks.md)
- **Validation**: All artifacts in `specs/001-macos-game-engine/` with bidirectional links

---

### Final Check (Post-Design)

**Re-evaluation Date**: 2025-10-22
**Phase 1 Artifacts**: research.md, data-model.md, contracts/engine-api.md, quickstart.md

#### ✅ I. Go Idiomatic Development (NON-NEGOTIABLE)
- **Status**: PASS (Confirmed)
- **Design Compliance**:
  - API contracts use Go interfaces, composition over inheritance
  - Error handling: All fallible operations return `(T, error)`
  - Naming: Entity, Transform, Sprite, Collider follow MixedCaps
  - Package organization: Flat hierarchy (engine/core, engine/graphics, not deep nesting)
- **Validation**: Code structure in plan.md follows Go project layout best practices

#### ✅ II. Test-Driven Development
- **Status**: PASS (Confirmed)
- **Design Compliance**:
  - Test directory structure defined: unit/, integration/, benchmarks/
  - API designed for testability (dependency injection via interfaces)
  - Benchmark targets: <16ms per frame, zero allocations in game loop
  - Example test cases in contracts: vector operations, collision detection, input state
- **Validation**: Tests organized by user story, enabling independent validation per acceptance criteria

#### ✅ III. Concurrent Agent Execution
- **Status**: PASS (Confirmed)
- **Design Compliance**:
  - Package structure enables parallel development: graphics/, input/, physics/ are independent
  - P1-P5 user stories map to separate packages (P1→core+graphics, P3→input, P5→physics)
  - No circular dependencies: math← physics/graphics/input← core
  - Clear contracts enable interface-driven development (agents can implement in parallel)
- **Validation**: Agent can work on input/ while another works on physics/ without conflicts

#### ✅ IV. Independent User Stories
- **Status**: PASS (Confirmed)
- **Design Compliance**:
  - P1 (Basic Rendering): Deliverable with only core/ + graphics/ + math/
  - P2 (Entity Management): Adds entity/ package, P1 still works
  - P3 (Input Handling): Adds input/ package, P1+P2 still work
  - P4 (Asset Loading): Extends graphics/assets.go, prior work unaffected
  - P5 (Collision): Adds physics/ package, all prior features functional
  - Examples directory mirrors user stories (examples/simple/, examples/player-control/, etc.)
- **Validation**: Each package can be tested independently via unit tests

#### ✅ V. Simplicity and YAGNI
- **Status**: PASS (Justified)
- **Design Compliance**:
  - **Single External Dependency**: SDL2 (github.com/veandco/go-sdl2)
    - **Justification**: Research validated no pure Go alternative for macOS Metal access
    - **Alternatives Considered**: Raw Metal bindings (too complex), Ebiten (full engine, not library), OpenGL (deprecated)
    - **Why SDL2**: Battle-tested, hardware-accelerated, cross-platform, automatic Metal backend
  - **Go Standard Library Preferred**: `image/png`, `image/jpeg`, `time`, `sync`, `testing` used where possible
  - **No Premature Optimization**: OOP entities chosen over ECS (ECS only beneficial >10,000 entities, spec targets 100-1000)
  - **No Speculative Features**: Audio, networking, 3D explicitly out of scope per assumptions
- **Validation**: Complexity Tracking section empty (no violations beyond SDL2)

#### ✅ VI. Structured Documentation
- **Status**: PASS (Confirmed)
- **Design Compliance**:
  - ✅ spec.md: Feature specification with P1-P5 user stories
  - ✅ plan.md: This file with technical context, constitution check, project structure
  - ✅ research.md: Graphics API, game loop, input handling research consolidated
  - ✅ data-model.md: All entities, components, subsystems defined
  - ✅ contracts/engine-api.md: Complete API specification for all packages
  - ✅ quickstart.md: Developer guide with P1-P5 examples
  - Bidirectional links: All documents reference each other at top
- **Validation**: All 6 artifacts present in `specs/001-macos-game-engine/` with consistent structure

---

### Summary

**All constitution principles: PASS**

**Key Validations**:
1. **Go Idiomatic**: API contracts follow Go conventions, explicit error handling, composition-based design
2. **TDD**: Test structure defined, API testable, benchmarks specified
3. **Concurrent Agents**: Package independence enables parallel work on P1-P5
4. **Independent Stories**: Each user story deliverable independently via package additions
5. **Simplicity**: Single justified external dependency (SDL2), no premature optimization
6. **Documentation**: All 6 Specify artifacts complete with bidirectional traceability

**Ready for Phase 2**: Yes - proceed with `/speckit.tasks` to generate implementation tasks

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
gogame/
├── engine/
│   ├── core/
│   │   ├── engine.go       # Main engine, game loop, SDL window/renderer
│   │   ├── scene.go        # Scene management, entity container
│   │   └── time.go         # Frame timing, fixed timestep accumulator
│   ├── graphics/
│   │   ├── renderer.go     # SDL2 rendering abstraction
│   │   ├── sprite.go       # Sprite component with texture reference
│   │   ├── texture.go      # Texture data structure
│   │   ├── camera.go       # Camera/viewport transform
│   │   └── assets.go       # AssetManager with reference counting
│   ├── entity/
│   │   ├── entity.go       # Entity struct with components
│   │   ├── transform.go    # Transform component (position, rotation, scale)
│   │   └── behavior.go     # Behavior interface for custom logic
│   ├── input/
│   │   ├── input.go        # InputManager with state tracking
│   │   ├── actions.go      # Action enumeration and mapping
│   │   └── keycodes.go     # KeyCode constants (SDL scancodes)
│   ├── physics/
│   │   ├── collision.go    # AABB collision detection
│   │   └── collider.go     # Collider component with bounds
│   └── math/
│       ├── vector.go       # Vector2 with operations
│       ├── rectangle.go    # Rectangle for bounds/regions
│       ├── transform.go    # Transform operations
│       └── color.go        # Color type with RGBA
├── examples/
│   ├── simple/             # P1: Basic rendering example
│   │   └── main.go         # 28-line minimal example
│   ├── moving/             # P2: Moving sprite example
│   │   └── main.go         # Entity with behavior
│   ├── player-control/     # P3: Input handling example
│   │   └── main.go         # Player controller with WASD
│   ├── assets/             # P4: Asset loading example
│   │   └── main.go         # Multiple textures
│   └── collision/          # P5: Collision detection example
│       └── main.go         # Player-enemy collision
├── tests/
│   ├── unit/
│   │   ├── vector_test.go      # Math operations
│   │   ├── collision_test.go   # AABB tests
│   │   └── input_test.go       # Action mapping
│   ├── integration/
│   │   ├── gameloop_test.go    # Full update/render cycle
│   │   └── scene_test.go       # Entity lifecycle
│   └── benchmarks/
│       ├── render_bench.go     # Sprite rendering performance
│       └── collision_bench.go  # Collision detection scaling
├── assets/                 # Example game assets
│   ├── player.png
│   ├── enemy.png
│   └── ...
├── go.mod
├── go.sum
├── CLAUDE.md               # Runtime AI agent guidance
└── README.md               # Project overview and quickstart
```

**Structure Decision**: Single project structure chosen. This is a game engine library with internal packages (engine/*) and example programs (examples/*). Not a web or mobile application. Engine packages follow Go package organization: `engine/core` for core types, `engine/graphics` for rendering, `engine/input` for input handling, `engine/physics` for collision, `engine/math` for utility types.

**Rationale**:
- Library structure enables external projects to import `github.com/dshills/gogame/engine`
- Examples directory demonstrates each user story (P1-P5)
- Tests organized by type (unit, integration, benchmarks)
- Flat package hierarchy avoids deep nesting (Go best practice)

## Complexity Tracking

No constitutional violations requiring justification.

**SDL2 Dependency**: External dependency for graphics API approved in Constitution Check (Principle V). Research phase validated SDL2 as minimal viable approach—simpler alternatives (pure Go graphics) don't exist for macOS Metal access.

**No Other Complexity**: Single-threaded design, OOP entities (not ECS), straightforward package organization all align with Simplicity principle.
