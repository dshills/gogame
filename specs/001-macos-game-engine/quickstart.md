# Quick Start Guide: macOS Game Engine

**Feature**: macOS Game Engine
**Version**: 1.0.0
**Target Audience**: Go developers building 2D games on macOS
**Related**: [spec.md](./spec.md) | [data-model.md](./data-model.md) | [contracts/engine-api.md](./contracts/engine-api.md)

## Overview

This guide demonstrates how to build a complete 2D game using the gogame engine. Examples progress from displaying a single sprite (5 minutes) to a playable game with player control, collision detection, and asset loading (30 minutes).

---

## Prerequisites

### System Requirements

- macOS 12.0 (Monterey) or newer
- Go 1.25.3+
- Homebrew package manager
- 50MB disk space for dependencies

### Installation

```bash
# Install SDL2 libraries via Homebrew
brew install sdl2 sdl2_image pkg-config

# Verify installation
pkg-config --modversion sdl2  # Should output 2.x.x

# Get the game engine
go get github.com/dshills/gogame
```

**Installation Time**: ~5 minutes

---

## Example 1: Hello Sprite (P1 - Basic Rendering)

**Goal**: Display a sprite on screen at 60 FPS
**Time**: 5 minutes
**Lines of Code**: 28

### Step 1: Create Project

```bash
mkdir my-game
cd my-game
go mod init my-game
go get github.com/dshills/gogame
```

### Step 2: Create Asset

Create `assets/player.png` (or download any 32x32 PNG image)

### Step 3: Write Game Code

`main.go`:
```go
package main

import (
    "log"

    "github.com/dshills/gogame/engine"
    "github.com/dshills/gogame/graphics"
    "github.com/dshills/gogame/math"
)

func main() {
    // Create engine with 800x600 window
    eng, err := engine.New("Hello Sprite", 800, 600, false)
    if err != nil {
        log.Fatal(err)
    }
    defer eng.Shutdown()

    // Create empty scene
    scene := engine.NewScene()
    eng.SetScene(scene)

    // Load player texture
    texture, err := eng.Assets().LoadTexture("assets/player.png")
    if err != nil {
        log.Fatal(err)
    }

    // Create entity with sprite at center
    player := &engine.Entity{
        Active:    true,
        Transform: math.Transform{
            Position: math.Vector2{X: 400, Y: 300},
            Scale:    math.Vector2{X: 1.0, Y: 1.0},
        },
        Sprite: graphics.NewSprite(texture),
        Layer:  1,
    }
    scene.AddEntity(player)

    // Run game loop (blocks until window closed)
    eng.Run()
}
```

### Step 4: Run

```bash
go run main.go
```

**Expected Result**: Window displays sprite at center, maintains 60 FPS

---

## Example 2: Moving Sprite (P2 - Entity Management)

**Goal**: Sprite moves smoothly across screen
**Time**: +5 minutes
**New Lines**: +15

### Step 1: Create Behavior

Add to `main.go` before `main()`:

```go
// MovingBehavior makes entity move horizontally
type MovingBehavior struct {
    Speed float64  // Pixels per second
}

func (mb *MovingBehavior) Update(entity *engine.Entity, dt float64) {
    // Move right
    entity.Transform.Position.X += mb.Speed * dt

    // Wrap around screen edges
    if entity.Transform.Position.X > 800 {
        entity.Transform.Position.X = 0
    }
}
```

### Step 2: Attach Behavior

Replace player creation with:

```go
player := &engine.Entity{
    Active:    true,
    Transform: math.Transform{
        Position: math.Vector2{X: 100, Y: 300},
        Scale:    math.Vector2{X: 1.0, Y: 1.0},
    },
    Sprite:   graphics.NewSprite(texture),
    Behavior: &MovingBehavior{Speed: 200},  // 200 pixels/second
    Layer:    1,
}
```

### Step 3: Run

```bash
go run main.go
```

**Expected Result**: Sprite moves smoothly right at 200 px/s, wraps at screen edge

**Frame-Rate Independence**: Movement speed consistent regardless of FPS

---

## Example 3: Player Control (P3 - Input Handling)

**Goal**: Control sprite with arrow keys or WASD
**Time**: +10 minutes
**New Lines**: +25

### Step 1: Define Actions

Add after imports:

```go
const (
    ActionMoveUp input.Action = iota
    ActionMoveDown
    ActionMoveLeft
    ActionMoveRight
)
```

### Step 2: Create Player Controller

Replace `MovingBehavior` with:

```go
type PlayerController struct {
    Speed float64
    input *input.InputManager
}

func (pc *PlayerController) Update(entity *engine.Entity, dt float64) {
    // Horizontal movement
    if pc.input.ActionHeld(ActionMoveRight) {
        entity.Transform.Position.X += pc.Speed * dt
    }
    if pc.input.ActionHeld(ActionMoveLeft) {
        entity.Transform.Position.X -= pc.Speed * dt
    }

    // Vertical movement
    if pc.input.ActionHeld(ActionMoveDown) {
        entity.Transform.Position.Y += pc.Speed * dt
    }
    if pc.input.ActionHeld(ActionMoveUp) {
        entity.Transform.Position.Y -= pc.Speed * dt
    }

    // Clamp to screen bounds
    if entity.Transform.Position.X < 0 {
        entity.Transform.Position.X = 0
    }
    if entity.Transform.Position.X > 800 {
        entity.Transform.Position.X = 800
    }
    if entity.Transform.Position.Y < 0 {
        entity.Transform.Position.Y = 0
    }
    if entity.Transform.Position.Y > 600 {
        entity.Transform.Position.Y = 600
    }
}
```

### Step 3: Configure Input Bindings

Add after engine creation:

```go
// Bind keys to actions
inputMgr := eng.Input()
inputMgr.BindAction(ActionMoveUp, input.KeyW, input.KeyArrowUp)
inputMgr.BindAction(ActionMoveDown, input.KeyS, input.KeyArrowDown)
inputMgr.BindAction(ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
inputMgr.BindAction(ActionMoveRight, input.KeyD, input.KeyArrowRight)
```

### Step 4: Update Player Entity

```go
player := &engine.Entity{
    Active:    true,
    Transform: math.Transform{
        Position: math.Vector2{X: 400, Y: 300},
        Scale:    math.Vector2{X: 1.0, Y: 1.0},
    },
    Sprite:   graphics.NewSprite(texture),
    Behavior: &PlayerController{Speed: 300, input: inputMgr},
    Layer:    1,
}
```

### Step 5: Run

```bash
go run main.go
```

**Expected Result**: Control sprite with WASD or arrow keys, clamped to screen

**Input Latency**: <16ms (one frame) from key press to movement

---

## Example 4: Multiple Assets (P4 - Asset Loading)

**Goal**: Add enemy sprite with different texture
**Time**: +5 minutes
**New Lines**: +20

### Step 1: Create Enemy Asset

Create `assets/enemy.png` (or different colored 32x32 PNG)

### Step 2: Load Enemy Texture

Add after player texture load:

```go
enemyTexture, err := eng.Assets().LoadTexture("assets/enemy.png")
if err != nil {
    log.Fatal(err)
}
```

### Step 3: Create Enemy Entity

Add after player creation:

```go
enemy := &engine.Entity{
    Active:    true,
    Transform: math.Transform{
        Position: math.Vector2{X: 600, Y: 200},
        Scale:    math.Vector2{X: 1.0, Y: 1.0},
    },
    Sprite:   graphics.NewSprite(enemyTexture),
    Behavior: &MovingBehavior{Speed: -150},  // Move left
    Layer:    1,
}
scene.AddEntity(enemy)
```

### Step 4: Run

```bash
go run main.go
```

**Expected Result**: Player (controllable) and enemy (moving) both on screen

**Asset Sharing**: If multiple entities use same texture, loaded once in memory

---

## Example 5: Collision Detection (P5 - Basic Physics)

**Goal**: Detect when player collides with enemy
**Time**: +10 minutes
**New Lines**: +30

### Step 1: Add Colliders

Update player creation:

```go
player := &engine.Entity{
    Active:    true,
    Transform: math.Transform{
        Position: math.Vector2{X: 400, Y: 300},
        Scale:    math.Vector2{X: 1.0, Y: 1.0},
    },
    Sprite:   graphics.NewSprite(texture),
    Collider: physics.NewCollider(32, 32),  // 32x32 hitbox
    Behavior: &PlayerController{Speed: 300, input: inputMgr},
    Layer:    1,
}
```

Update enemy creation:

```go
enemy := &engine.Entity{
    Active:    true,
    Transform: math.Transform{
        Position: math.Vector2{X: 600, Y: 200},
        Scale:    math.Vector2{X: 1.0, Y: 1.0},
    },
    Sprite:   graphics.NewSprite(enemyTexture),
    Collider: physics.NewCollider(32, 32),
    Behavior: &MovingBehavior{Speed: -150},
    Layer:    1,
}
```

### Step 2: Create Collision Checker

Add to `PlayerController.Update()` at end:

```go
// Check collision with all entities
scene := eng.GetScene()
for _, other := range scene.GetEntitiesAt(
    entity.Transform.Position.X,
    entity.Transform.Position.Y,
) {
    if other.ID != entity.ID && other.Collider != nil {
        if entity.Collider.Intersects(
            other.Collider,
            &entity.Transform,
            &other.Transform,
        ) {
            // Collision detected!
            fmt.Println("Player hit enemy!")

            // Change sprite color to red
            entity.Sprite.SetColor(math.Color{R: 255, G: 0, B: 0, A: 255})
        }
    }
}
```

### Step 3: Run

```bash
go run main.go
```

**Expected Result**: Player turns red when touching enemy, "Player hit enemy!" printed

**Collision Performance**: <1ms for 50 entities (well within 16ms frame budget)

---

## Common Patterns

### Camera Following Player

```go
func (pc *PlayerController) Update(entity *engine.Entity, dt float64) {
    // ... movement code ...

    // Camera smoothly follows player
    camera := scene.Camera()
    camera.Follow(entity, 0.1)  // 10% interpolation per frame
}
```

### Spawning Entities Dynamically

```go
// Spawn enemy on mouse click
if inputMgr.ActionPressed(ActionFire) {
    mouseX, mouseY := inputMgr.MousePosition()
    worldX, worldY := camera.ScreenToWorld(int(mouseX), int(mouseY))

    newEnemy := &engine.Entity{
        Active:    true,
        Transform: math.Transform{Position: math.Vector2{X: worldX, Y: worldY}},
        Sprite:    graphics.NewSprite(enemyTexture),
        Collider:  physics.NewCollider(32, 32),
        Layer:     1,
    }
    scene.AddEntity(newEnemy)
}
```

### Removing Entities

```go
// Remove enemy when out of bounds
func (mb *MovingBehavior) Update(entity *engine.Entity, dt float64) {
    entity.Transform.Position.X += mb.Speed * dt

    if entity.Transform.Position.X < -100 {
        scene.RemoveEntity(entity.ID)
    }
}
```

### Sprite Sheets (Animation)

```go
// Extract specific frame from sprite sheet
sprite := graphics.NewSprite(spriteSheetTexture)
sprite.SetSourceRect(
    frameIndex*32,  // X offset (32px per frame)
    0,              // Y offset
    32,             // Frame width
    32,             // Frame height
)
```

### Background Color

```go
scene.SetBackgroundColor(math.Color{R: 135, G: 206, B: 235, A: 255})  // Sky blue
```

---

## Performance Tips

### Memory Management

```go
// Pre-allocate entity slice if count known
const MaxEntities = 100
entities := make([]Entity, 0, MaxEntities)  // No reallocations up to 100

// Reuse entities via object pooling
var entityPool sync.Pool
entity := entityPool.Get().(*Entity)
defer entityPool.Put(entity)
```

### Reduce Allocations in Game Loop

```go
// Bad: Creates slice every frame
func Update(dt float64) {
    nearby := []Entity{}  // Allocation!
    // ...
}

// Good: Reuse slice
var nearby []Entity
func Update(dt float64) {
    nearby = nearby[:0]  // Reset length, keep capacity
    // ...
}
```

### Benchmark Performance

```go
go test -bench=. -benchmem ./...
```

**Target**: <16ms per frame, zero allocations in game loop

---

## Troubleshooting

### SDL2 Linking Errors

**Error**: `ld: library not found for -lSDL2`

**Solution**:
```bash
brew reinstall sdl2
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"
```

### Window Not Appearing

**Cause**: Engine.Run() called without setting scene

**Solution**:
```go
scene := engine.NewScene()
eng.SetScene(scene)  // Must set before Run()
eng.Run()
```

### Sprite Not Rendering

**Checklist**:
1. ✅ Entity.Active = true?
2. ✅ Entity.Sprite != nil?
3. ✅ Sprite.Texture loaded successfully?
4. ✅ Entity added to scene?
5. ✅ Entity on-screen (within camera view)?

### Low Frame Rate

**Causes**:
- Too many entities (>1000)
- Heavy computation in Update()
- Allocations in game loop

**Solutions**:
```go
// Profile to find bottleneck
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

// Check GC pauses
GODEBUG=gctrace=1 go run main.go
```

### Input Not Working

**Checklist**:
1. ✅ Actions defined as constants?
2. ✅ BindAction() called before Run()?
3. ✅ InputManager queried in Update(), not init?
4. ✅ Window has focus?

---

## Next Steps

### Build Complete Game

**Platformer Example**: `examples/platformer/`
- Player with jump/run
- Multiple levels
- Collectible items
- Enemy AI
- Score tracking

**Arcade Shooter Example**: `examples/shooter/`
- Player projectiles
- Enemy waves
- Particle effects
- Power-ups
- High score

### Advanced Topics

**Spatial Partitioning**: For >1000 entities, use quadtree
**Animation Systems**: Frame-based sprite animation
**Tilemaps**: Efficient large level rendering
**Audio Integration**: Add sound effects and music
**Save/Load Systems**: Persist game state

### Community

- **Repository**: https://github.com/dshills/gogame
- **Documentation**: https://pkg.go.dev/github.com/dshills/gogame
- **Examples**: https://github.com/dshills/gogame/tree/main/examples
- **Issues**: https://github.com/dshills/gogame/issues

---

## Success Metrics

Validate your game meets spec requirements:

| Metric | Target | How to Measure |
|--------|--------|----------------|
| Code Simplicity | <50 lines for basic game | Count lines in main.go |
| Frame Rate | 60 FPS with 100 sprites | Add FPS counter, spawn 100 entities |
| Asset Load Time | <100ms for <5MB images | Measure LoadTexture() duration |
| Input Latency | <16ms | Timestamp key press to entity response |
| Memory Stability | No leaks over 1 hour | Monitor memory usage: `top -pid $(pgrep my-game)` |

---

## Reference

### Complete Minimal Example

Full `main.go` for P1-P5 features (< 100 lines):

```go
package main

import (
    "fmt"
    "log"

    "github.com/dshills/gogame/engine"
    "github.com/dshills/gogame/graphics"
    "github.com/dshills/gogame/input"
    "github.com/dshills/gogame/math"
    "github.com/dshills/gogame/physics"
)

// Actions
const (
    ActionMoveUp input.Action = iota
    ActionMoveDown
    ActionMoveLeft
    ActionMoveRight
)

// Player controller
type PlayerController struct {
    Speed float64
    input *input.InputManager
    scene *engine.Scene
}

func (pc *PlayerController) Update(entity *engine.Entity, dt float64) {
    // Movement
    if pc.input.ActionHeld(ActionMoveRight) {
        entity.Transform.Position.X += pc.Speed * dt
    }
    if pc.input.ActionHeld(ActionMoveLeft) {
        entity.Transform.Position.X -= pc.Speed * dt
    }
    if pc.input.ActionHeld(ActionMoveDown) {
        entity.Transform.Position.Y += pc.Speed * dt
    }
    if pc.input.ActionHeld(ActionMoveUp) {
        entity.Transform.Position.Y -= pc.Speed * dt
    }

    // Collision check
    for _, other := range pc.scene.GetEntitiesAt(
        entity.Transform.Position.X,
        entity.Transform.Position.Y,
    ) {
        if other.ID != entity.ID && other.Collider != nil {
            if entity.Collider.Intersects(
                other.Collider,
                &entity.Transform,
                &other.Transform,
            ) {
                fmt.Println("Collision!")
                entity.Sprite.SetColor(math.Color{R: 255, G: 0, B: 0, A: 255})
            }
        }
    }
}

// Enemy AI
type MovingBehavior struct {
    Speed float64
}

func (mb *MovingBehavior) Update(entity *engine.Entity, dt float64) {
    entity.Transform.Position.X += mb.Speed * dt
    if entity.Transform.Position.X > 800 || entity.Transform.Position.X < 0 {
        mb.Speed = -mb.Speed  // Bounce
    }
}

func main() {
    // Engine
    eng, err := engine.New("Complete Example", 800, 600, false)
    if err != nil {
        log.Fatal(err)
    }
    defer eng.Shutdown()

    // Scene
    scene := engine.NewScene()
    scene.SetBackgroundColor(math.Color{R: 135, G: 206, B: 235, A: 255})
    eng.SetScene(scene)

    // Input
    inputMgr := eng.Input()
    inputMgr.BindAction(ActionMoveUp, input.KeyW, input.KeyArrowUp)
    inputMgr.BindAction(ActionMoveDown, input.KeyS, input.KeyArrowDown)
    inputMgr.BindAction(ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
    inputMgr.BindAction(ActionMoveRight, input.KeyD, input.KeyArrowRight)

    // Assets
    playerTex, _ := eng.Assets().LoadTexture("assets/player.png")
    enemyTex, _ := eng.Assets().LoadTexture("assets/enemy.png")

    // Player
    player := &engine.Entity{
        Active:    true,
        Transform: math.Transform{Position: math.Vector2{X: 400, Y: 300}},
        Sprite:    graphics.NewSprite(playerTex),
        Collider:  physics.NewCollider(32, 32),
        Behavior:  &PlayerController{Speed: 300, input: inputMgr, scene: scene},
        Layer:     1,
    }
    scene.AddEntity(player)

    // Enemy
    enemy := &engine.Entity{
        Active:    true,
        Transform: math.Transform{Position: math.Vector2{X: 600, Y: 300}},
        Sprite:    graphics.NewSprite(enemyTex),
        Collider:  physics.NewCollider(32, 32),
        Behavior:  &MovingBehavior{Speed: -150},
        Layer:     1,
    }
    scene.AddEntity(enemy)

    // Run
    eng.Run()
}
```

---

**Quickstart Status**: ✅ Complete
**Examples**: P1 through P5 covered
**Target Audience**: Go developers new to game development
**Estimated Time**: 30 minutes for all examples
