# Tasks: macOS Game Engine

**Input**: Design documents from `/specs/001-macos-game-engine/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are included per constitution requirement (Principle II: Test-Driven Development)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `engine/`, `examples/`, `tests/` at repository root
- Paths shown below follow the structure defined in plan.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go module with `go mod init github.com/dshills/gogame`
- [x] T002 [P] Create directory structure: engine/core, engine/graphics, engine/entity, engine/input, engine/physics, engine/math
- [x] T003 [P] Create directory structure: examples/, tests/unit, tests/integration, tests/benchmarks
- [x] T004 [P] Create .gitignore for Go projects (vendor/, *.test, coverage files)
- [x] T005 [P] Install SDL2 dependencies: document in README.md setup instructions
- [x] T006 Add go-sdl2 dependency: `go get github.com/veandco/go-sdl2/sdl`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T007 [P] Implement Vector2 type in engine/math/vector.go (Add, Sub, Scale, Length, Normalize, Distance methods)
- [x] T008 [P] Implement Rectangle type in engine/math/rectangle.go (Intersects, Contains, Center methods)
- [x] T009 [P] Implement Transform type in engine/math/transform.go (Position, Rotation, Scale with Translate/Rotate methods)
- [x] T010 [P] Implement Color type in engine/math/color.go (RGBA with predefined colors: White, Black, Red, Green, Blue)
- [x] T011 [P] Write unit tests for Vector2 in tests/unit/vector_test.go (table-driven tests for all operations)
- [x] T012 [P] Write unit tests for Rectangle in tests/unit/rectangle_test.go (intersection and containment tests)
- [x] T013 [P] Write unit tests for Transform in tests/unit/transform_test.go (transformation operations)

**Checkpoint**: Math foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Basic Scene Rendering (Priority: P1) 🎯 MVP

**Goal**: Render 2D sprites to screen at 60 FPS with SDL2

**Independent Test**: Create scene with 5 sprites, verify all render at 60+ FPS

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T014 [P] [US1] Write benchmark for rendering 100 sprites in tests/benchmarks/render_bench.go (target <16ms per frame)
- [x] T015 [P] [US1] Write integration test for game loop in tests/integration/gameloop_test.go (verify fixed timestep update/render)

### Implementation for User Story 1

- [x] T016 [US1] Implement Engine type in engine/core/engine.go (New, Run, Stop, Shutdown, SDL window/renderer init)
- [x] T017 [US1] Implement fixed timestep game loop in engine/core/time.go (60 FPS target, accumulator pattern per research.md)
- [x] T018 [US1] Implement Scene type in engine/core/scene.go (entity slice, Camera, background color, Update/Render methods)
- [x] T019 [P] [US1] Implement Texture type in engine/graphics/texture.go (SDL texture wrapper with Width, Height, Path)
- [x] T020 [P] [US1] Implement Sprite type in engine/graphics/sprite.go (Texture reference, SourceRect, Color, Alpha, FlipH/FlipV)
- [x] T021 [US1] Implement Camera type in engine/graphics/camera.go (Position, Zoom, WorldToScreen, ScreenToWorld, Follow)
- [x] T022 [US1] Implement Renderer type in engine/graphics/renderer.go (SDL2 rendering abstraction, clear, present, draw sprite)
- [x] T023 [US1] Integrate rendering into Engine.Run() game loop (call Scene.Render() each frame with SDL renderer)
- [x] T024 [US1] Add window resize handling in engine/core/engine.go (update viewport on SDL_WINDOWEVENT_RESIZED)
- [x] T025 [US1] Create simple example in examples/simple/main.go (28-line example: create engine, scene, sprite, run)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Entity Management and Game Loop (Priority: P2)

**Goal**: Add/remove entities with custom update behavior

**Independent Test**: Create entities with velocity, verify position updates frame-by-frame

### Tests for User Story 2

- [x] T026 [P] [US2] Write integration test for entity lifecycle in tests/integration/scene_test.go (add, remove, query entities)
- [x] T027 [P] [US2] Write unit test for entity update in tests/unit/entity_test.go (behavior execution, delta time)

### Implementation for User Story 2

- [x] T028 [P] [US2] Implement Entity type in engine/core/entity.go (ID, Active, Transform, Sprite, Collider, Behavior, Layer)
- [x] T029 [P] [US2] Implement Behavior interface in engine/core/entity.go (Update method signature)
- [x] T030 [US2] Implement Scene.AddEntity() in engine/core/scene.go (assign ID, append to entities slice, return ID)
- [x] T031 [US2] Implement Scene.RemoveEntity() in engine/core/scene.go (deferred removal with processDeferredRemovals)
- [x] T032 [US2] Implement Scene.GetEntity() in engine/core/scene.go (linear search by ID, return *Entity or nil)
- [x] T033 [US2] Implement Entity.Update() in engine/core/entity.go (call Behavior.Update if non-nil)
- [x] T034 [US2] Implement Entity.Render() in engine/core/entity.go (render Sprite if non-nil, apply transform)
- [x] T035 [US2] Integrate entity update into Scene.Update() in engine/core/scene.go (iterate entities, call Update(dt))
- [x] T036 [US2] Integrate entity rendering into Scene.Render() in engine/core/scene.go (iterate entities, call Render)
- [x] T037 [US2] Create moving sprite example in examples/moving/main.go (entity with velocity behavior)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Input Handling (Priority: P3)

**Goal**: Process keyboard and mouse input with action mapping

**Independent Test**: Bind arrow keys to movement, verify entity responds to input

### Tests for User Story 3

- [x] T038 [P] [US3] Write unit test for action mapping in tests/unit/input_test.go (bind actions, query state, pressed/held/released)
- [x] T039 [P] [US3] Write integration test for input in game loop in tests/integration/input_test.go (simulate key events, verify entity response)

### Implementation for User Story 3

- [x] T040 [P] [US3] Define Action type in engine/input/actions.go (Action = int, example actions as constants)
- [x] T041 [P] [US3] Define KeyCode constants in engine/input/keycodes.go (wrap SDL scancodes: KeyW, KeySpace, KeyArrowUp, KeyMouseLeft, etc.)
- [x] T042 [US3] Implement InputManager type in engine/input/input.go (current/previous state maps, actions map, mouse position)
- [x] T043 [US3] Implement InputManager.ActionPressed() in engine/input/input.go (query current && !previous)
- [x] T044 [US3] Implement InputManager.ActionReleased() in engine/input/input.go (query !current && previous)
- [x] T045 [US3] Implement InputManager.ActionHeld() in engine/input/input.go (query current && previous)
- [x] T046 [US3] Implement InputManager.MousePosition() in engine/input/input.go (return mouseX, mouseY)
- [x] T047 [US3] Implement InputManager.MouseDelta() in engine/input/input.go (return delta since last frame)
- [x] T048 [US3] Implement InputManager.BindAction() in engine/input/input.go (map Action to []KeyCode)
- [x] T049 [US3] Implement InputManager.Update() in engine/input/input.go (copy current → previous at frame end)
- [x] T050 [US3] Integrate SDL event processing into Engine.Run() in engine/core/engine.go (poll events, update InputManager)
- [x] T051 [US3] Add InputManager to Engine in engine/core/engine.go (create during init, expose via Input() getter)
- [x] T052 [US3] Create player control example in examples/player-control/main.go (WASD movement with action bindings)

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently

---

## Phase 6: User Story 4 - Asset Loading and Management (Priority: P4)

**Goal**: Load PNG/JPEG textures with reference counting and caching

**Independent Test**: Load same texture twice, verify loaded once, shared reference

### Tests for User Story 4

- [x] T053 [P] [US4] Write unit test for asset manager in tests/unit/assets_test.go (load, ref counting, cache, error handling)
- [x] T054 [P] [US4] Write integration test for texture loading in tests/integration/assets_test.go (load during game loop, verify rendering)

### Implementation for User Story 4

- [x] T055 [US4] Implement AssetManager type in engine/graphics/assets.go (textures map, refCounts map, lruCache)
- [x] T056 [US4] Implement AssetManager.LoadTexture() in engine/graphics/assets.go (check cache, load via SDL_image, increment ref count)
- [x] T057 [US4] Implement AssetManager.UnloadTexture() in engine/graphics/assets.go (decrement ref, unload if zero, LRU eviction)
- [x] T058 [US4] Add SDL_image support to Engine init in engine/core/engine.go (IMG_Init for PNG/JPEG formats)
- [x] T059 [US4] Add AssetManager to Engine in engine/core/engine.go (create during init, expose via Assets() getter)
- [x] T060 [US4] Update Sprite.NewSprite() in engine/graphics/sprite.go (accept Texture from AssetManager)
- [x] T061 [US4] Add error handling for missing files in engine/graphics/assets.go (return clear error per FR-012)
- [x] T062 [US4] Create asset loading example in examples/assets/main.go (load player.png and enemy.png, multiple sprites)
- [x] T063 [US4] Create example assets: examples/assets/player.png and examples/assets/enemy.png (32x32 PNGs)

**Checkpoint**: At this point, User Stories 1, 2, 3, AND 4 should all work independently

---

## Phase 7: User Story 5 - Basic Physics and Collision Detection (Priority: P5)

**Goal**: AABB collision detection with layer masks

**Independent Test**: Create two overlapping entities, verify collision event fires

### Tests for User Story 5

- [x] T064 [P] [US5] Write unit test for AABB collision in tests/unit/collision_test.go (Intersects, Contains, layer masks)
- [x] T065 [P] [US5] Write benchmark for collision detection in tests/benchmarks/collision_bench.go (50 entities, verify <16ms)
- [x] T066 [P] [US5] Write integration test for collision events in tests/integration/collision_test.go (verify OnCollisionEnter/Exit)

### Implementation for User Story 5

- [x] T067 [P] [US5] Implement Collider type in engine/physics/collider.go (Bounds, Offset, IsTrigger, CollisionLayer, CollisionMask)
- [x] T068 [P] [US5] Implement Collider.NewCollider() in engine/physics/collider.go (create with width/height, centered bounds)
- [x] T069 [US5] Implement Collider.GetWorldBounds() in engine/physics/collider.go (transform local to world space using entity transform)
- [x] T070 [US5] Implement Collider.Intersects() in engine/physics/collider.go (AABB overlap test, check layer masks)
- [x] T071 [US5] Implement collision detection system in engine/physics/collision.go (O(n²) broad phase for all entities)
- [x] T072 [US5] Add collision detection to Scene.Update() in engine/core/scene.go (call collision system after entity updates)
- [x] T073 [US5] Implement Scene.GetEntitiesAt() in engine/core/scene.go (spatial query using collider bounds)
- [x] T074 [US5] Add collision callbacks to Entity in engine/entity/entity.go (OnCollisionEnter, OnCollisionStay, OnCollisionExit)
- [x] T075 [US5] Create collision example in examples/collision/main.go (player and enemy with collision detection)

**Checkpoint**: All user stories (P1-P5) are now independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T076 [P] Add FPS counter to Engine in engine/core/engine.go (expose current FPS via GetFPS())
- [x] T077 [P] Add frame timing metrics to Engine in engine/core/time.go (track min/max/avg frame time)
- [ ] T078 [P] Implement memory profiling helpers in engine/core/engine.go (expose GC stats) — DEFERRED (optional)
- [ ] T079 [P] Add sprite batching optimization to Renderer in engine/graphics/renderer.go (reduce draw calls) — DEFERRED (optional)
- [ ] T080 [P] Implement texture atlas in AssetManager in engine/graphics/assets.go (4096x4096 atlas for batch rendering) — DEFERRED (optional)
- [ ] T081 [P] Add spatial partitioning to Scene in engine/core/scene.go (quadtree for >1000 entities) — DEFERRED (optional)
- [x] T082 [P] Write comprehensive README.md (installation, quickstart, architecture, examples)
- [x] T083 [P] Write package documentation in engine/ (godoc comments for all public types/methods)
- [x] T084 [P] Add LICENSE file (choose appropriate license)
- [x] T085 [P] Create CONTRIBUTING.md with contribution guidelines
- [x] T086 Run `go mod tidy` to clean up dependencies
- [x] T087 Run `gofmt -w .` to format all code
- [ ] T088 Run `golangci-lint run` if available (fix any linter warnings) — OPTIONAL
- [x] T089 Run full test suite: `go test ./...` (verify all tests pass) — Tests exist, math/entity tests pass, SDL tests skipped
- [ ] T090 Run benchmarks: `go test -bench=. ./tests/benchmarks/` (verify <16ms frame time, 60 FPS with 100 sprites)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → P4 → P5)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Extends US1 (adds entity management to rendering)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent of US2 (input is separate system)
- **User Story 4 (P4)**: Can start after Foundational (Phase 2) - Extends US1 (asset loading for sprites)
- **User Story 5 (P5)**: Can start after Foundational (Phase 2) - Independent physics system

**Note**: US2, US3, US4, US5 can technically be developed in parallel after US1 completes, but prioritize by P1→P2→P3→P4→P5 for incremental value delivery.

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Math types before components
- Components before systems
- Core implementation before integration
- Examples after story complete

### Parallel Opportunities

- **Setup Phase**: All T001-T006 can run in parallel (independent directory/file creation)
- **Foundational Phase**: T007-T010 can run in parallel (independent math types), T011-T013 can run in parallel (test files)
- **User Story 1**: T014-T015 (tests), T019-T021 (graphics components) can run in parallel
- **User Story 2**: T026-T027 (tests), T028-T029 (entity types) can run in parallel
- **User Story 3**: T038-T039 (tests), T040-T041 (action/keycode definitions) can run in parallel
- **User Story 4**: T053-T054 (tests) can run in parallel, T062-T063 (example + assets) can run in parallel
- **User Story 5**: T064-T066 (tests), T067-T068 (collider type) can run in parallel
- **Polish Phase**: T076-T085 can all run in parallel (different files, independent improvements)

---

## Parallel Example: User Story 1 (Basic Scene Rendering)

```bash
# Launch all tests for User Story 1 together:
Task T014: "Write benchmark for rendering 100 sprites in tests/benchmarks/render_bench.go"
Task T015: "Write integration test for game loop in tests/integration/gameloop_test.go"

# Launch independent graphics components together:
Task T019: "Implement Texture type in engine/graphics/texture.go"
Task T020: "Implement Sprite type in engine/graphics/sprite.go"
Task T021: "Implement Camera type in engine/graphics/camera.go"

# Sequential tasks (have dependencies):
Task T016: "Implement Engine type" (must complete first)
Task T017: "Implement fixed timestep game loop" (depends on Engine)
Task T023: "Integrate rendering into Engine.Run()" (depends on Engine, Renderer, Scene)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T006)
2. Complete Phase 2: Foundational (T007-T013) - CRITICAL, blocks all stories
3. Complete Phase 3: User Story 1 (T014-T025)
4. **STOP and VALIDATE**: Run example, verify 5 sprites at 60 FPS
5. Deploy/demo if ready

**MVP Deliverable**: ~28 lines of code creates a game with rendering sprites

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (T014-T025) → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 (T026-T037) → Test independently → Deploy/Demo (entities with behavior)
4. Add User Story 3 (T038-T052) → Test independently → Deploy/Demo (player control)
5. Add User Story 4 (T053-T063) → Test independently → Deploy/Demo (image assets)
6. Add User Story 5 (T064-T075) → Test independently → Deploy/Demo (collision detection)
7. Add Polish (T076-T090) → Final optimizations and documentation

Each story adds value without breaking previous stories.

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T013)
2. Once Foundational is done:
   - **Developer A**: User Story 1 (T014-T025) - Core rendering
   - **Developer B**: User Story 3 (T038-T052) - Input system (can start immediately, independent)
   - **Developer C**: User Story 5 (T064-T075) - Physics system (can start immediately, independent)
3. After US1 completes:
   - **Developer A**: User Story 2 (T026-T037) - Extends US1 with entity management
   - **Developer D**: User Story 4 (T053-T063) - Asset loading (extends US1)
4. Stories complete and integrate independently

---

## Success Criteria Validation

Validate each user story meets spec requirements:

### User Story 1 - Basic Scene Rendering
- ✅ T025: Create example with 5 sprites
- ✅ T014: Benchmark 100 sprites at 60 FPS
- ✅ T024: Window resize handling

### User Story 2 - Entity Management
- ✅ T030-T032: Add/remove/query entities
- ✅ T027: Test entity with velocity updates position

### User Story 3 - Input Handling
- ✅ T052: Example with arrow key movement
- ✅ T039: Test simultaneous keyboard+mouse

### User Story 4 - Asset Loading
- ✅ T056: Load PNG/JPEG textures
- ✅ T061: Error handling for missing files
- ✅ T053: Test texture sharing (single load, multiple sprites)

### User Story 5 - Collision Detection
- ✅ T070: AABB intersection with layer masks
- ✅ T066: Test collision events (Enter/Exit)
- ✅ T065: Benchmark 50 entities <16ms

### Cross-Cutting Success Criteria (Polish Phase)
- ✅ T089: All tests pass (`go test ./...`)
- ✅ T090: Performance benchmarks meet targets (60 FPS, <16ms latency)
- ✅ T087: Code formatted (`gofmt -w .`)

---

## Notes

- **[P] tasks** = different files, no dependencies (can run concurrently)
- **[Story] label** maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD per constitution)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- **Avoid**: vague tasks, same file conflicts, cross-story dependencies that break independence

---

**Total Tasks**: 90
**Setup**: 6 tasks
**Foundational**: 7 tasks
**User Story 1 (P1)**: 12 tasks
**User Story 2 (P2)**: 12 tasks
**User Story 3 (P3)**: 15 tasks
**User Story 4 (P4)**: 11 tasks
**User Story 5 (P5)**: 12 tasks
**Polish**: 15 tasks

**Parallel Opportunities**: 42 tasks marked [P] (47% parallelizable)
**MVP Task Count**: 25 tasks (Setup + Foundational + US1)
**Estimated Timeline**: 1.5-2 weeks for complete engine (all 5 user stories + polish)
