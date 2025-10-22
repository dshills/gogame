# Data Model: macOS Game Engine

**Feature**: macOS Game Engine
**Date**: 2025-10-22
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md) | [research.md](./research.md)

## Overview

This document defines all entities, components, and data structures for the game engine. The architecture follows a hybrid OOP approach with composable components, as determined in the research phase. Entities are structured for 100-1000 concurrent instances with minimal memory overhead and GC pressure.

---

## Core Entities

### Engine

The root orchestrator that manages the game loop, scene, and subsystems.

**Attributes**:
- `window`: SDL window handle for rendering
- `renderer`: SDL hardware-accelerated renderer
- `scene`: Currently active game scene
- `input`: InputManager for keyboard/mouse handling
- `assets`: AssetManager for texture loading and caching
- `running`: Boolean flag controlling game loop
- `targetFPS`: Target frame rate (default 60)
- `fixedTimestep`: Fixed update interval (1/60 seconds)

**Lifecycle**:
- `Initialize()`: Create window, renderer, subsystems
- `Run()`: Execute main game loop (update/render cycle)
- `Shutdown()`: Clean up resources, close window

**Relationships**:
- Contains one Scene
- Contains one InputManager
- Contains one AssetManager
- Manages game loop timing

---

### Scene

Represents a game level or screen, containing all entities and managing their lifecycle.

**Attributes**:
- `entities`: Slice of all active entities in the scene
- `camera`: Camera defining view transform
- `backgroundColor`: RGBA background color
- `entityIDCounter`: Auto-increment counter for entity IDs

**Operations**:
- `AddEntity(entity *Entity) uint64`: Add entity, return assigned ID
- `RemoveEntity(id uint64)`: Remove entity by ID
- `GetEntity(id uint64) *Entity`: Query entity by ID
- `GetEntitiesAt(x, y float64) []*Entity`: Spatial query for entities at position
- `Update(dt float64)`: Update all entities
- `Render(renderer *Renderer)`: Render all visible entities

**Lifecycle**:
- Entities added/removed dynamically during gameplay
- Update phase iterates all entities
- Render phase culls off-screen entities

**Relationships**:
- Contains many Entities (0 to 1000)
- Contains one Camera
- Owned by Engine

---

### Entity

A game object that exists in the scene with position, appearance, and behavior.

**Attributes**:
- `ID`: Unique identifier (uint64)
- `Active`: Boolean flag (inactive entities skip update/render)
- `Transform`: Required transform component (position, rotation, scale)
- `Sprite`: Optional sprite component (visual representation)
- `Collider`: Optional collision component (bounding box)
- `Behavior`: Optional behavior interface (custom update logic)
- `Layer`: Z-order for rendering (higher renders on top)

**Operations**:
- `Update(dt float64)`: Update transform, execute behavior
- `Render(renderer *Renderer, camera *Camera)`: Render sprite if present
- `SetBehavior(behavior Behavior)`: Attach custom update logic
- `HasCollider() bool`: Check if collision-enabled
- `GetBounds() Rectangle`: Return world-space bounding box

**Lifecycle**:
- Created via `Scene.AddEntity()`
- Updated every frame if Active
- Rendered every frame if Active and Sprite non-nil
- Removed via `Scene.RemoveEntity()`

**Relationships**:
- Owned by Scene
- Has one Transform (required)
- Has zero or one Sprite (optional)
- Has zero or one Collider (optional)
- Has zero or one Behavior (optional)

---

## Component Types

### Transform

Position, rotation, and scale information determining entity's location and orientation in world space.

**Attributes**:
- `Position`: Vector2 (X, Y in world coordinates)
- `Rotation`: float64 (angle in degrees, 0° = right, 90° = down)
- `Scale`: Vector2 (X scale, Y scale, 1.0 = no scaling)
- `Pivot`: Vector2 (rotation/scale origin point, relative to sprite)

**Operations**:
- `Translate(dx, dy float64)`: Move by offset
- `Rotate(degrees float64)`: Rotate by angle
- `SetScale(sx, sy float64)`: Set scale factors
- `GetWorldMatrix() Matrix3x3`: Compute transformation matrix for rendering

**Usage**:
- Required for all entities
- Modified during update phase for movement
- Read during render phase for positioning

**Coordinate System**:
- Origin (0,0) = top-left of screen
- X-axis increases right
- Y-axis increases down
- Rotation: 0° right, 90° down, 180° left, 270° up

---

### Sprite

Visual representation component that renders a texture at the entity's transform.

**Attributes**:
- `Texture`: Reference to loaded texture (via AssetManager)
- `SourceRect`: Rectangle defining region of texture to render (for sprite sheets)
- `Color`: RGBA tint color (white = no tint)
- `Alpha`: Opacity (0.0 = transparent, 1.0 = opaque)
- `FlipH`: Boolean horizontal flip
- `FlipV`: Boolean vertical flip

**Operations**:
- `SetTexture(texture *Texture)`: Assign texture from AssetManager
- `SetSourceRect(x, y, w, h int)`: Define sprite sheet region
- `SetColor(r, g, b, a uint8)`: Apply color tint
- `Render(renderer *Renderer, transform *Transform)`: Draw sprite at transform

**Rendering**:
- Uses texture coordinates for sprite sheet support
- Applies transform (position, rotation, scale)
- Rendered in layer order (lower layers first)
- Culled if off-screen

**Memory**:
- Textures are reference-counted via AssetManager (shared across sprites)
- Multiple sprites can reference same texture
- SourceRect enables sprite sheets (single texture, multiple sprites)

---

### Collider

Rectangular bounding box for collision detection.

**Attributes**:
- `Bounds`: Rectangle (X, Y, Width, Height) in local space
- `Offset`: Vector2 offset from entity position
- `IsTrigger`: Boolean (trigger = no physics response, just events)
- `CollisionLayer`: Bit mask defining which layers this collider belongs to
- `CollisionMask`: Bit mask defining which layers this collider checks

**Operations**:
- `GetWorldBounds(transform *Transform) Rectangle`: Compute world-space bounds
- `Intersects(other *Collider, transform, otherTransform *Transform) bool`: Check overlap
- `Contains(x, y float64, transform *Transform) bool`: Point-in-box test

**Collision Detection**:
- Axis-Aligned Bounding Box (AABB) collision
- O(n²) broad phase (acceptable for 50-100 collidable entities per spec)
- Layer masks enable selective collision (e.g., player vs enemies, not player vs player)
- Trigger colliders generate events without blocking movement

**Events**:
- OnCollisionEnter: First frame of overlap
- OnCollisionStay: Continuing overlap
- OnCollisionExit: First frame of separation

---

### Behavior

Interface for custom entity update logic.

**Interface**:
```go
type Behavior interface {
    Update(entity *Entity, dt float64)
}
```

**Usage**:
- Developers implement this interface for custom logic
- Called every frame during update phase
- Has full access to entity's transform, sprite, collider

**Examples**:
- Player controller: Read input, move entity
- Enemy AI: Follow player, patrol waypoints
- Projectile: Move forward, destroy on collision

**Implementation**:
```go
type PlayerController struct {
    Speed float64
}

func (pc *PlayerController) Update(entity *Entity, dt float64) {
    if input.ActionHeld(ActionMoveRight) {
        entity.Transform.Position.X += pc.Speed * dt
    }
}
```

---

## Subsystems

### AssetManager

Manages texture loading, caching, and reference counting.

**Attributes**:
- `textures`: Map of path → Texture
- `refCounts`: Map of path → int (reference count)
- `lruCache`: LRU cache (keeps 50 most recent)
- `textureAtlas`: Internal atlas for batch rendering

**Operations**:
- `LoadTexture(path string) (*Texture, error)`: Load from disk or return cached
- `UnloadTexture(path string)`: Decrement ref count, unload if zero
- `GetTexture(path string) *Texture`: Query loaded texture
- `Clear()`: Unload all textures

**Caching Strategy**:
- Lazy loading: Load on first reference
- Reference counting: Track active users
- LRU cache: Keep 50 most recent at zero refs
- Texture atlas: Auto-pack for batch rendering

**Memory**:
- Typical 2D game: 50-100 textures, 200-400MB
- Pre-allocated decode buffers via sync.Pool (GC optimization)

---

### InputManager

Centralized keyboard and mouse input handling with action mapping.

**Attributes**:
- `current`: Map of KeyCode → bool (this frame state)
- `previous`: Map of KeyCode → bool (last frame state)
- `actions`: Map of Action → []KeyCode (action bindings)
- `mouseX, mouseY`: Current mouse position (window coords)
- `mouseDeltaX, mouseDeltaY`: Mouse movement this frame

**Operations**:
- `ActionPressed(action Action) bool`: Just pressed this frame
- `ActionReleased(action Action) bool`: Just released this frame
- `ActionHeld(action Action) bool`: Held this frame
- `MousePosition() (float64, float64)`: Get mouse window coords
- `MouseWorldPosition(camera *Camera) (float64, float64)`: Get mouse world coords
- `BindAction(action Action, keys ...KeyCode)`: Configure key binding
- `Update()`: Shift current → previous at frame end

**Input Processing**:
1. SDL event pump updates current state
2. Game logic queries actions
3. `Update()` shifts state for next frame

**Action Mapping**:
```go
const (
    ActionMoveUp Action = iota
    ActionMoveDown
    ActionMoveLeft
    ActionMoveRight
    ActionJump
    ActionFire
)

// Default bindings
ActionMoveUp: [KeyW, KeyArrowUp]
ActionMoveRight: [KeyD, KeyArrowRight]
```

---

### Camera

Defines view transform from world space to screen space.

**Attributes**:
- `Position`: Vector2 (camera center in world space)
- `Zoom`: float64 (1.0 = normal, >1.0 = zoomed in, <1.0 = zoomed out)
- `Rotation`: float64 (camera rotation in degrees)
- `ViewportWidth, ViewportHeight`: int (screen dimensions)

**Operations**:
- `WorldToScreen(worldX, worldY float64) (screenX, screenY int)`: Transform position
- `ScreenToWorld(screenX, screenY int) (worldX, worldY float64)`: Inverse transform
- `GetVisibleBounds() Rectangle`: Compute world-space rectangle visible on screen
- `Follow(target *Entity, smoothing float64)`: Smooth camera follow

**Rendering**:
- All entities rendered relative to camera
- Entities outside visible bounds are culled
- Zoom affects sprite sizes
- Rotation affects entire scene (rarely used in 2D)

---

## Data Structures

### Vector2

2D vector for positions, velocities, offsets.

**Fields**:
- `X float64`
- `Y float64`

**Operations**:
- `Add(other Vector2) Vector2`
- `Sub(other Vector2) Vector2`
- `Scale(factor float64) Vector2`
- `Length() float64`
- `Normalize() Vector2`
- `Distance(other Vector2) float64`

---

### Rectangle

Axis-aligned rectangle for bounds, colliders, source regions.

**Fields**:
- `X float64` (left edge)
- `Y float64` (top edge)
- `Width float64`
- `Height float64`

**Operations**:
- `Intersects(other Rectangle) bool`: Check overlap
- `Contains(x, y float64) bool`: Point-in-rect test
- `Union(other Rectangle) Rectangle`: Bounding rectangle
- `Center() Vector2`: Compute center point

---

### Color

RGBA color with 8-bit channels.

**Fields**:
- `R uint8` (red, 0-255)
- `G uint8` (green, 0-255)
- `B uint8` (blue, 0-255)
- `A uint8` (alpha, 0-255, 255 = opaque)

**Predefined Colors**:
- White: (255, 255, 255, 255)
- Black: (0, 0, 0, 255)
- Transparent: (0, 0, 0, 0)

---

### Texture

Loaded image data from disk.

**Attributes**:
- `sdlTexture`: SDL texture handle (internal)
- `Width int`: Texture width in pixels
- `Height int`: Texture height in pixels
- `Path string`: Source file path

**Lifecycle**:
- Created by AssetManager.LoadTexture()
- Shared across multiple sprites
- Reference-counted
- Freed when all references released and LRU evicts

---

### Action

Enumeration of game actions (developer-defined).

**Type**: `int` (Go const)

**Examples**:
```go
const (
    ActionMoveUp Action = iota
    ActionMoveDown
    ActionMoveLeft
    ActionMoveRight
    ActionJump
    ActionFire
    ActionPause
)
```

**Usage**:
- Decouples game logic from hardware keys
- Enables player-configurable bindings
- Future gamepad support without code changes

---

### KeyCode

Enumeration of keyboard keys and mouse buttons.

**Type**: `uint32` (SDL scancode)

**Common Values**:
- `KeyW`, `KeyA`, `KeyS`, `KeyD`
- `KeySpace`, `KeyEnter`, `KeyEscape`
- `KeyArrowUp`, `KeyArrowDown`, `KeyArrowLeft`, `KeyArrowRight`
- `KeyMouseLeft`, `KeyMouseRight`, `KeyMouseMiddle`

**Mapping**: Wraps SDL scancodes for platform-independent key identification

---

## Entity State Transitions

### Entity Lifecycle States

```
[Created] → [Active] → [Destroyed]
              ↓
           [Inactive] ← (can toggle Active flag)
```

**Created**: Entity allocated, components initialized, ID assigned
**Active**: Entity updates and renders each frame
**Inactive**: Entity skips update/render (still in scene)
**Destroyed**: Entity removed from scene, components released

---

## Collision Event Flow

```
Frame N: Entity A and B not overlapping
Frame N+1: Entity A and B overlap
  → OnCollisionEnter(A, B) fires

Frame N+2: Entity A and B still overlapping
  → OnCollisionStay(A, B) fires

Frame N+3: Entity A and B no longer overlapping
  → OnCollisionExit(A, B) fires
```

**Event Data**:
- `Entity *Entity`: The other entity in collision
- `Collider *Collider`: The other collider
- `Normal Vector2`: Collision surface normal (for physics response)

---

## Memory Layout Optimization

### Entity Storage

**Option 1: Slice of Structs** (chosen for simplicity)
```go
entities []Entity  // Cache-friendly iteration
```

**Option 2: Struct of Arrays** (considered but rejected as premature optimization)
```go
transforms []Transform
sprites []Sprite
colliders []*Collider  // Optional components use pointers
```

**Rationale**: For 100-1000 entities, slice-of-structs simplicity outweighs SoA cache benefits. Reconsider if scale exceeds 10,000 entities.

### Component Optionality

- **Transform**: Always present (embedded in Entity)
- **Sprite**: Optional (pointer, nil if no visual)
- **Collider**: Optional (pointer, nil if no collision)
- **Behavior**: Optional (interface, nil if no custom logic)

**Memory Savings**: Empty entities (Transform only) use ~64 bytes vs 256+ bytes with all components

---

## Validation Against Requirements

| Requirement | Data Model Support |
|-------------|-------------------|
| FR-001: Configurable window | Engine.window, Engine.renderer |
| FR-002: 60 FPS with 100 entities | Fixed timestep, optimized entity iteration |
| FR-003: 2D sprite rendering | Sprite component with transform |
| FR-004: Game loop | Engine.Run() with update/render phases |
| FR-005: Add/remove entities | Scene.AddEntity(), Scene.RemoveEntity() |
| FR-006: Load PNG/JPEG | AssetManager.LoadTexture() |
| FR-007: Keyboard/mouse input | InputManager with action mapping |
| FR-008: Collision detection | Collider component with AABB |
| FR-009: Custom behavior | Behavior interface |
| FR-010: Coordinate system | Transform with top-left origin |
| FR-011: Graceful shutdown | Engine.Shutdown() |
| FR-012: Clear error reporting | Texture loading returns error |
| FR-013: Z-ordering | Entity.Layer for render sorting |
| FR-014: Frame timing | Engine.fixedTimestep, delta time |
| FR-015: RGBA alpha | Sprite.Alpha, Color with A channel |

---

## Extensibility Points

For future features beyond initial scope:

- **Audio**: Add AudioSource component, AudioManager subsystem
- **Animation**: Add Animation component with keyframes
- **Particle Systems**: Add ParticleEmitter component
- **Tilemaps**: Add Tilemap entity type
- **Physics**: Add RigidBody component with velocity/forces
- **Networking**: Add NetworkIdentity component for sync
- **UI**: Add Canvas entity type with UI elements

All extensions follow same component pattern: optional components on entities, managed by specialized subsystems.

---

**Data Model Status**: ✅ Complete
**Entities Defined**: 8 core + 6 components + 2 subsystems
**Ready for Contracts Phase**: Yes
