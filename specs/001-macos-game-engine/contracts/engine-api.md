# Engine API Contract

**Feature**: macOS Game Engine
**Version**: 1.0.0
**Date**: 2025-10-22
**Related**: [data-model.md](../data-model.md)

## Overview

This document defines the public API contracts for the game engine. All types, methods, and behaviors are specified to enable parallel development and ensure consistency. API follows Go idioms with explicit error handling, composition over inheritance, and simple interfaces.

---

## Package Structure

```
github.com/dshills/gogame/
├── engine/             # Core engine types
├── graphics/           # Rendering and sprites
├── input/              # Input handling
├── physics/            # Collision detection
└── math/               # Vector and geometry types
```

---

## Core Engine API

### Package: `engine`

#### Type: `Engine`

The root game engine managing window, rendering, and game loop.

```go
type Engine struct {
    // Private fields (implementation detail)
}

// New creates a new game engine instance
//
// Parameters:
//   title: Window title
//   width: Window width in pixels
//   height: Window height in pixels
//   fullscreen: Start in fullscreen mode
//
// Returns:
//   *Engine: Initialized engine
//   error: Non-nil if window/renderer creation fails
//
// Example:
//   engine, err := engine.New("My Game", 800, 600, false)
//   if err != nil {
//       log.Fatal(err)
//   }
func New(title string, width, height int, fullscreen bool) (*Engine, error)

// SetScene sets the active scene
//
// Parameters:
//   scene: Scene to activate
//
// Behavior:
//   - Previous scene (if any) is not destroyed (developer must manage)
//   - New scene begins updating/rendering immediately
//
// Example:
//   menuScene := engine.NewScene()
//   engine.SetScene(menuScene)
func (e *Engine) SetScene(scene *Scene)

// GetScene returns the currently active scene
//
// Returns:
//   *Scene: Active scene, or nil if none set
func (e *Engine) GetScene() *Scene

// Run starts the game loop (blocking)
//
// Behavior:
//   - Runs until window closed or Stop() called
//   - Fixed 60 FPS update rate
//   - Variable rendering rate (vsync if enabled)
//   - Calls scene Update() and Render() each frame
//
// Example:
//   engine.Run()  // Blocks until game ends
func (e *Engine) Run()

// Stop signals the game loop to exit
//
// Behavior:
//   - Game loop exits after current frame completes
//   - Resources remain allocated (call Shutdown to cleanup)
//
// Example:
//   if input.ActionPressed(ActionQuit) {
//       engine.Stop()
//   }
func (e *Engine) Stop()

// Shutdown releases all engine resources
//
// Behavior:
//   - Destroys SDL window and renderer
//   - Unloads all textures
//   - Must be called before program exit
//   - Engine unusable after this call
//
// Example:
//   defer engine.Shutdown()
func (e *Engine) Shutdown()

// Input returns the input manager
//
// Returns:
//   *input.InputManager: Input subsystem
func (e *Engine) Input() *input.InputManager

// Assets returns the asset manager
//
// Returns:
//   *graphics.AssetManager: Asset loading subsystem
func (e *Engine) Assets() *graphics.AssetManager
```

---

#### Type: `Scene`

Container for entities representing a game level or screen.

```go
type Scene struct {
    // Private fields
}

// NewScene creates an empty scene
//
// Returns:
//   *Scene: New scene with no entities
//
// Example:
//   scene := engine.NewScene()
func NewScene() *Scene

// AddEntity adds an entity to the scene
//
// Parameters:
//   entity: Entity to add
//
// Returns:
//   uint64: Assigned entity ID (unique within scene)
//
// Behavior:
//   - Entity begins updating/rendering immediately if Active
//   - ID assigned sequentially starting from 1
//
// Example:
//   player := &engine.Entity{Active: true}
//   playerID := scene.AddEntity(player)
func (s *Scene) AddEntity(entity *Entity) uint64

// RemoveEntity removes an entity by ID
//
// Parameters:
//   id: Entity ID to remove
//
// Behavior:
//   - Entity removed immediately (doesn't update/render next frame)
//   - Safe to call during Update() (deferred removal)
//   - No-op if ID not found
//
// Example:
//   scene.RemoveEntity(enemyID)
func (s *Scene) RemoveEntity(id uint64)

// GetEntity retrieves an entity by ID
//
// Parameters:
//   id: Entity ID to query
//
// Returns:
//   *Entity: Entity with matching ID, or nil if not found
//
// Example:
//   entity := scene.GetEntity(playerID)
//   if entity != nil {
//       entity.Transform.Position.X += 10
//   }
func (s *Scene) GetEntity(id uint64) *Entity

// GetEntitiesAt finds all entities at a world position
//
// Parameters:
//   x, y: World coordinates
//
// Returns:
//   []*Entity: Entities whose colliders contain the point (may be empty)
//
// Behavior:
//   - Only checks entities with non-nil Collider
//   - Returns entities in arbitrary order
//   - Empty slice if no matches
//
// Example:
//   mouseWorldX, mouseWorldY := camera.ScreenToWorld(mouseX, mouseY)
//   entities := scene.GetEntitiesAt(mouseWorldX, mouseWorldY)
func (s *Scene) GetEntitiesAt(x, y float64) []*Entity

// Camera returns the scene's camera
//
// Returns:
//   *graphics.Camera: Scene camera for view transform
func (s *Scene) Camera() *graphics.Camera

// SetBackgroundColor sets the clear color
//
// Parameters:
//   color: RGBA color to clear screen with
//
// Example:
//   scene.SetBackgroundColor(math.Color{R: 135, G: 206, B: 235, A: 255})  // Sky blue
func (s *Scene) SetBackgroundColor(color math.Color)
```

---

#### Type: `Entity`

Game object with position, optional visuals, collision, and behavior.

```go
type Entity struct {
    ID        uint64                // Unique identifier (assigned by Scene)
    Active    bool                  // Update/render only if true
    Transform math.Transform        // Position, rotation, scale (required)
    Sprite    *graphics.Sprite      // Optional visual representation
    Collider  *physics.Collider     // Optional collision bounds
    Behavior  Behavior               // Optional custom update logic
    Layer     int                    // Z-order (higher renders on top)
}

// Update updates the entity's transform and behavior
//
// Parameters:
//   dt: Delta time in seconds
//
// Behavior:
//   - Calls Behavior.Update() if non-nil
//   - Called automatically by Scene during update phase
//
// Example:
//   // Typically called by engine, not user code
//   entity.Update(0.016)  // 16ms frame
func (e *Entity) Update(dt float64)

// Render draws the entity's sprite
//
// Parameters:
//   renderer: SDL renderer
//   camera: Camera for view transform
//
// Behavior:
//   - Renders Sprite if non-nil
//   - Applies transform (position, rotation, scale)
//   - Culls if outside camera view
//   - Called automatically by Scene during render phase
//
// Example:
//   // Typically called by engine, not user code
//   entity.Render(renderer, camera)
func (e *Entity) Render(renderer *graphics.Renderer, camera *graphics.Camera)

// GetBounds returns world-space bounding box
//
// Returns:
//   math.Rectangle: World-space bounds (or zero rect if no collider)
//
// Example:
//   bounds := entity.GetBounds()
//   if bounds.Contains(clickX, clickY) {
//       fmt.Println("Entity clicked!")
//   }
func (e *Entity) GetBounds() math.Rectangle
```

---

#### Interface: `Behavior`

Custom entity update logic.

```go
// Behavior defines custom per-frame logic for an entity
//
// Implemented by game developers to control entity behavior
type Behavior interface {
    // Update is called every frame
    //
    // Parameters:
    //   entity: The entity this behavior is attached to
    //   dt: Delta time in seconds (typically 0.016 at 60 FPS)
    //
    // Example:
    //   func (pc *PlayerController) Update(entity *Entity, dt float64) {
    //       if input.ActionHeld(ActionMoveRight) {
    //           entity.Transform.Position.X += pc.Speed * dt
    //       }
    //   }
    Update(entity *Entity, dt float64)
}
```

---

## Graphics API

### Package: `graphics`

#### Type: `Sprite`

Visual representation attached to entities.

```go
type Sprite struct {
    Texture    *Texture        // Loaded texture (via AssetManager)
    SourceRect math.Rectangle  // Region of texture to render (for sprite sheets)
    Color      math.Color      // Tint color (white = no tint)
    Alpha      float64         // Opacity (0.0 = transparent, 1.0 = opaque)
    FlipH      bool            // Flip horizontally
    FlipV      bool            // Flip vertically
}

// NewSprite creates a sprite from a texture
//
// Parameters:
//   texture: Loaded texture
//
// Returns:
//   *Sprite: Sprite rendering full texture
//
// Example:
//   texture, _ := assets.LoadTexture("player.png")
//   sprite := graphics.NewSprite(texture)
func NewSprite(texture *Texture) *Sprite

// SetSourceRect sets the sprite sheet region
//
// Parameters:
//   x, y: Top-left corner in texture
//   width, height: Region dimensions
//
// Example:
//   // Extract 32x32 sprite from sprite sheet
//   sprite.SetSourceRect(64, 0, 32, 32)
func (s *Sprite) SetSourceRect(x, y, width, height int)

// SetColor sets the tint color
//
// Parameters:
//   color: RGBA tint (white = no tint)
//
// Example:
//   sprite.SetColor(math.Color{R: 255, G: 0, B: 0, A: 255})  // Red tint
func (s *Sprite) SetColor(color math.Color)
```

---

#### Type: `AssetManager`

Manages texture loading and caching.

```go
type AssetManager struct {
    // Private fields
}

// LoadTexture loads a texture from disk or returns cached
//
// Parameters:
//   path: File path (PNG or JPEG)
//
// Returns:
//   *Texture: Loaded texture
//   error: Non-nil if file not found or decode fails
//
// Behavior:
//   - Returns existing texture if already loaded
//   - Increments reference count
//   - Caches in LRU cache
//
// Example:
//   texture, err := assets.LoadTexture("assets/player.png")
//   if err != nil {
//       log.Fatal(err)
//   }
func (am *AssetManager) LoadTexture(path string) (*Texture, error)

// UnloadTexture decrements reference count
//
// Parameters:
//   path: File path of texture to unload
//
// Behavior:
//   - Decrements reference count
//   - Unloads if count reaches zero and LRU evicts
//   - Safe to call multiple times
//   - No-op if texture not loaded
//
// Example:
//   assets.UnloadTexture("assets/player.png")
func (am *AssetManager) UnloadTexture(path string)
```

---

#### Type: `Camera`

Defines view transformation from world to screen space.

```go
type Camera struct {
    Position math.Vector2  // Camera center in world space
    Zoom     float64       // Zoom factor (1.0 = normal, >1.0 = zoomed in)
}

// NewCamera creates a camera at origin with no zoom
//
// Returns:
//   *Camera: Camera at (0,0) with zoom 1.0
func NewCamera() *Camera

// WorldToScreen transforms world coordinates to screen pixels
//
// Parameters:
//   worldX, worldY: World coordinates
//
// Returns:
//   screenX, screenY: Screen pixel coordinates
//
// Example:
//   screenX, screenY := camera.WorldToScreen(entity.Transform.Position.X, entity.Transform.Position.Y)
func (c *Camera) WorldToScreen(worldX, worldY float64) (screenX, screenY int)

// ScreenToWorld transforms screen pixels to world coordinates
//
// Parameters:
//   screenX, screenY: Screen pixel coordinates
//
// Returns:
//   worldX, worldY: World coordinates
//
// Example:
//   worldX, worldY := camera.ScreenToWorld(mouseX, mouseY)
//   entities := scene.GetEntitiesAt(worldX, worldY)
func (c *Camera) ScreenToWorld(screenX, screenY int) (worldX, worldY float64)

// Follow smoothly moves camera toward target
//
// Parameters:
//   target: Entity to follow
//   smoothing: Interpolation factor (0.0 = instant, 1.0 = no follow)
//
// Example:
//   camera.Follow(player, 0.1)  // Smooth follow with 10% interpolation
func (c *Camera) Follow(target *Entity, smoothing float64)
```

---

## Input API

### Package: `input`

#### Type: `InputManager`

Centralized keyboard and mouse input handling.

```go
type InputManager struct {
    // Private fields
}

// ActionPressed returns true if action just pressed this frame
//
// Parameters:
//   action: Action to query
//
// Returns:
//   bool: True if pressed this frame (was up last frame)
//
// Example:
//   if input.ActionPressed(ActionJump) {
//       player.Jump()
//   }
func (im *InputManager) ActionPressed(action Action) bool

// ActionReleased returns true if action just released this frame
//
// Parameters:
//   action: Action to query
//
// Returns:
//   bool: True if released this frame (was down last frame)
func (im *InputManager) ActionReleased(action Action) bool

// ActionHeld returns true if action currently held
//
// Parameters:
//   action: Action to query
//
// Returns:
//   bool: True if down this frame and last frame
//
// Example:
//   if input.ActionHeld(ActionMoveRight) {
//       player.Transform.Position.X += player.Speed * dt
//   }
func (im *InputManager) ActionHeld(action Action) bool

// MousePosition returns current mouse position in window coordinates
//
// Returns:
//   x, y: Mouse position (top-left = 0,0)
//
// Example:
//   mouseX, mouseY := input.MousePosition()
//   worldX, worldY := camera.ScreenToWorld(int(mouseX), int(mouseY))
func (im *InputManager) MousePosition() (x, y float64)

// MouseDelta returns mouse movement this frame
//
// Returns:
//   dx, dy: Pixel offset from last frame
//
// Example:
//   dx, dy := input.MouseDelta()
//   camera.Position.X -= dx
//   camera.Position.Y -= dy
func (im *InputManager) MouseDelta() (dx, dy float64)

// BindAction configures key bindings for an action
//
// Parameters:
//   action: Action to bind
//   keys: One or more keys that trigger this action
//
// Example:
//   input.BindAction(ActionJump, KeySpace, KeyW)  // Space or W jumps
func (im *InputManager) BindAction(action Action, keys ...KeyCode)
```

---

#### Type: `Action`

Game action enumeration (developer-defined).

```go
type Action int

// Example actions (developers define their own)
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

---

#### Type: `KeyCode`

Keyboard key and mouse button enumeration.

```go
type KeyCode uint32

// Keyboard keys
const (
    KeyA KeyCode = iota
    KeyB
    // ... (full alphabet)
    KeySpace
    KeyEnter
    KeyEscape
    KeyArrowUp
    KeyArrowDown
    KeyArrowLeft
    KeyArrowRight
    // ... (all keys)
)

// Mouse buttons
const (
    KeyMouseLeft KeyCode = 1000 + iota
    KeyMouseRight
    KeyMouseMiddle
)
```

---

## Physics API

### Package: `physics`

#### Type: `Collider`

Rectangular collision bounds.

```go
type Collider struct {
    Bounds         math.Rectangle  // Local-space bounds (relative to entity)
    Offset         math.Vector2    // Offset from entity position
    IsTrigger      bool            // Trigger = events only, no blocking
    CollisionLayer uint32          // Layer bitmask (which layer am I in?)
    CollisionMask  uint32          // Mask bitmask (which layers do I collide with?)
}

// NewCollider creates a collider with bounds
//
// Parameters:
//   width, height: Collision box dimensions
//
// Returns:
//   *Collider: Collider centered on entity
//
// Example:
//   collider := physics.NewCollider(32, 48)  // 32x48 hitbox
func NewCollider(width, height float64) *Collider

// GetWorldBounds returns world-space collision box
//
// Parameters:
//   transform: Entity's transform
//
// Returns:
//   math.Rectangle: World-space axis-aligned bounds
//
// Example:
//   worldBounds := collider.GetWorldBounds(&entity.Transform)
//   if worldBounds.Intersects(otherBounds) {
//       // Collision detected
//   }
func (c *Collider) GetWorldBounds(transform *math.Transform) math.Rectangle

// Intersects checks collision with another collider
//
// Parameters:
//   other: Other collider
//   transform: This entity's transform
//   otherTransform: Other entity's transform
//
// Returns:
//   bool: True if colliders overlap
//
// Example:
//   if collider.Intersects(enemy.Collider, &player.Transform, &enemy.Transform) {
//       player.TakeDamage()
//   }
func (c *Collider) Intersects(other *Collider, transform, otherTransform *math.Transform) bool
```

---

## Math Types API

### Package: `math`

#### Type: `Vector2`

2D vector for positions, velocities, offsets.

```go
type Vector2 struct {
    X float64
    Y float64
}

// Add returns vector sum
func (v Vector2) Add(other Vector2) Vector2

// Sub returns vector difference
func (v Vector2) Sub(other Vector2) Vector2

// Scale returns scaled vector
func (v Vector2) Scale(factor float64) Vector2

// Length returns magnitude
func (v Vector2) Length() float64

// Normalize returns unit vector
func (v Vector2) Normalize() Vector2

// Distance returns distance to other vector
func (v Vector2) Distance(other Vector2) float64
```

---

#### Type: `Rectangle`

Axis-aligned rectangle for bounds, regions, collision.

```go
type Rectangle struct {
    X      float64  // Left edge
    Y      float64  // Top edge
    Width  float64
    Height float64
}

// Intersects checks if rectangles overlap
func (r Rectangle) Intersects(other Rectangle) bool

// Contains checks if point is inside rectangle
func (r Rectangle) Contains(x, y float64) bool

// Center returns center point
func (r Rectangle) Center() Vector2
```

---

#### Type: `Transform`

Position, rotation, and scale for entity placement.

```go
type Transform struct {
    Position math.Vector2  // World position
    Rotation float64       // Angle in degrees (0° = right, 90° = down)
    Scale    math.Vector2  // Scale factors (1.0 = normal)
}

// Translate moves by offset
func (t *Transform) Translate(dx, dy float64)

// Rotate rotates by angle
func (t *Transform) Rotate(degrees float64)
```

---

#### Type: `Color`

RGBA color with 8-bit channels.

```go
type Color struct {
    R uint8  // Red (0-255)
    G uint8  // Green (0-255)
    B uint8  // Blue (0-255)
    A uint8  // Alpha (0-255, 255 = opaque)
}

// Predefined colors
var (
    White       = Color{255, 255, 255, 255}
    Black       = Color{0, 0, 0, 255}
    Red         = Color{255, 0, 0, 255}
    Green       = Color{0, 255, 0, 255}
    Blue        = Color{0, 0, 255, 255}
    Transparent = Color{0, 0, 0, 0}
)
```

---

## Usage Examples

### Minimal Game (P1 - Basic Rendering)

```go
package main

import (
    "github.com/dshills/gogame/engine"
    "github.com/dshills/gogame/graphics"
    "github.com/dshills/gogame/math"
)

func main() {
    // Create engine
    eng, _ := engine.New("My Game", 800, 600, false)
    defer eng.Shutdown()

    // Create scene
    scene := engine.NewScene()
    eng.SetScene(scene)

    // Load texture
    texture, _ := eng.Assets().LoadTexture("player.png")

    // Create entity with sprite
    player := &engine.Entity{
        Active:    true,
        Transform: math.Transform{Position: math.Vector2{X: 400, Y: 300}},
        Sprite:    graphics.NewSprite(texture),
        Layer:     1,
    }
    scene.AddEntity(player)

    // Run game
    eng.Run()
}
```

**Lines**: 28 (meets SC-001: "under 50 lines")

---

### Player Controller (P3 - Input Handling)

```go
type PlayerController struct {
    Speed float64
}

func (pc *PlayerController) Update(entity *engine.Entity, dt float64) {
    input := eng.Input()

    // Movement
    if input.ActionHeld(input.ActionMoveRight) {
        entity.Transform.Position.X += pc.Speed * dt
    }
    if input.ActionHeld(input.ActionMoveLeft) {
        entity.Transform.Position.X -= pc.Speed * dt
    }

    // Jump
    if input.ActionPressed(input.ActionJump) {
        entity.Transform.Position.Y -= 100  // Simple jump
    }
}

// Attach to entity
player.Behavior = &PlayerController{Speed: 200}
```

**Lines**: 5 for common tasks (meets SC-005)

---

## Error Handling

All functions that can fail return `error` as the last return value:

```go
engine, err := engine.New("Game", 800, 600, false)
if err != nil {
    log.Fatal("Failed to create engine:", err)
}

texture, err := assets.LoadTexture("missing.png")
if err != nil {
    log.Fatal("Failed to load texture:", err)  // Clear error message
}
```

**Error Messages** (FR-012):
- "Failed to create SDL window: <reason>"
- "Failed to load texture: file not found: <path>"
- "Failed to decode image: <reason>"

---

## Thread Safety

**Single-threaded design** (per spec assumptions):
- All engine methods called from main game loop thread
- No mutexes required
- SDL operations inherently thread-unsafe

**Future**: If multi-threading added, add `sync.RWMutex` to Scene and InputManager

---

## Performance Contracts

| Operation | Time Complexity | Notes |
|-----------|-----------------|-------|
| Scene.AddEntity() | O(1) | Append to slice |
| Scene.RemoveEntity() | O(n) | Linear search + swap-remove |
| Scene.GetEntity() | O(n) | Linear search (use sparingly) |
| Scene.GetEntitiesAt() | O(n) | Check all collidable entities |
| Collider.Intersects() | O(1) | AABB overlap test |
| AssetManager.LoadTexture() | O(1) cached, O(disk) first load | Hash map lookup |

**Optimization Notes**:
- For <1000 entities, linear search acceptable (cache-friendly)
- Spatial partitioning (quadtree) if exceeds 1000 entities
- Collision broad phase O(n²) acceptable for <100 collidable entities

---

**Contract Status**: ✅ Complete
**API Coverage**: All 15 functional requirements (FR-001 through FR-015)
**Ready for Implementation**: Yes
