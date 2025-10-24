# gogame - 2D Game Engine for macOS

> ⚠️ **ALPHA SOFTWARE**: This engine is in active development. APIs may change without notice. Not recommended for production use. Expect bugs, incomplete features, and breaking changes.

A simple, performant 2D game engine written in Go for macOS, designed for creating arcade-style games and platformers.

## Features

### ✅ Currently Implemented

**Core Engine**
- **60 FPS Rendering**: Hardware-accelerated sprite rendering with SDL2/Metal backend
- **Entity-Component System**: Hybrid architecture with Entity structs and Behavior interface
- **Fixed Timestep Loop**: Consistent 60 FPS updates with delta time for frame-rate independence
- **Scene Management**: Entity containers with background colors, layers, and camera system

**Graphics & Rendering**
- **Sprite Rendering**: PNG/JPEG texture loading with reference counting and caching
- **Text Rendering**: TTF font support with SDL2_ttf for in-game text display
- **Visual Effects**: Color tinting, alpha blending, sprite flipping (horizontal/vertical)
- **Camera System**: World-to-screen transforms with position, zoom, and smooth following
- **Transform System**: Position, rotation, and scale with interpolation support
- **Layer Rendering**: Z-ordering for proper sprite draw order

**Input System**
- **Keyboard Input**: Full keyboard support with pressed/held/released states
- **Mouse Input**: Position tracking, button states, and movement delta
- **Action Mapping**: Bind multiple keys to named actions (e.g., "Jump", "MoveLeft")
- **Frame-Perfect Input**: Double-buffered state for accurate edge detection

**Physics & Collision**
- **AABB Collision Detection**: Axis-aligned bounding box collision with efficient O(n²) broad phase
- **Collision Callbacks**: Event-driven collision handling with OnCollisionEnter/Stay/Exit
- **Collision Filtering**: Layer masks for selective collision detection
- **Collider Components**: Attach colliders to entities for automatic collision detection

**Asset Management**
- **Texture Loading**: PNG and JPEG support with automatic format detection
- **Reference Counting**: Efficient texture reuse with automatic cleanup
- **Asset Caching**: Single texture instance shared across multiple sprites

### 🚧 Planned Features

- **Audio System**: Sound effects and music playback with SDL_mixer
- **Sprite Animation**: Frame-based animation with sprite sheets
- **Particle Systems**: Configurable particle emitters for effects
- **Tilemap Support**: Efficient rendering of tile-based levels
- **Spatial Partitioning**: Quadtree/grid for optimized collision detection
- **Physics Engine**: Velocity, acceleration, and basic physics simulation
- **UI System**: Buttons, labels, and basic UI components

## Prerequisites

### macOS Requirements

- **macOS 12.0 (Monterey) or newer**
- **Homebrew** package manager
- **Go 1.25.3 or newer**

### Install SDL2

```bash
# Install SDL2, SDL2_image, and SDL2_ttf via Homebrew
brew install sdl2 sdl2_image sdl2_ttf pkg-config

# Verify installation
pkg-config --modversion sdl2
# Should output: 2.x.x
pkg-config --modversion sdl2_ttf
# Should output: 2.x.x
```

## Installation

```bash
# Clone the repository
git clone https://github.com/dshills/gogame.git
cd gogame

# Initialize Go module
go mod download

# Verify installation by running simple example
go run examples/simple/main.go
```

## Quick Start

### Minimal Example (35 lines)

```go
package main

import (
    "log"
    "runtime"
    "github.com/dshills/gogame/engine/core"
    gamemath "github.com/dshills/gogame/engine/math"
)

func main() {
    // CRITICAL: SDL requires running on the main OS thread
    runtime.LockOSThread()

    // Create engine with 800x600 window
    engine, err := core.NewEngine("My Game", 800, 600, false)
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Shutdown()

    // Create scene with sky blue background
    scene := core.NewScene()
    scene.SetBackgroundColor(gamemath.Color{R: 52, G: 152, B: 219, A: 255})

    // Position camera at center (IMPORTANT for visibility)
    camera := scene.Camera()
    camera.Position = gamemath.Vector2{X: 400, Y: 300}

    engine.SetScene(scene)

    // Run game loop (blocks until window closed)
    if err := engine.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### With Sprites and Movement

```go
package main

import (
    "log"
    "runtime"
    "github.com/dshills/gogame/engine/core"
    "github.com/dshills/gogame/engine/graphics"
    "github.com/dshills/gogame/engine/input"
    gamemath "github.com/dshills/gogame/engine/math"
)

// PlayerController moves entity with WASD/Arrow keys
type PlayerController struct {
    Speed    float64
    InputMgr *input.InputManager
}

func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
    moveSpeed := pc.Speed * dt
    if pc.InputMgr.ActionHeld(input.ActionMoveUp) {
        entity.Transform.Position.Y -= moveSpeed
    }
    if pc.InputMgr.ActionHeld(input.ActionMoveDown) {
        entity.Transform.Position.Y += moveSpeed
    }
    if pc.InputMgr.ActionHeld(input.ActionMoveLeft) {
        entity.Transform.Position.X -= moveSpeed
    }
    if pc.InputMgr.ActionHeld(input.ActionMoveRight) {
        entity.Transform.Position.X += moveSpeed
    }
}

func main() {
    runtime.LockOSThread()

    engine, err := core.NewEngine("Player Control", 800, 600, false)
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Shutdown()

    // Setup input bindings
    inputMgr := engine.Input()
    inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
    inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
    inputMgr.BindAction(input.ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
    inputMgr.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)

    scene := core.NewScene()
    scene.SetBackgroundColor(gamemath.Color{R: 30, G: 30, B: 50, A: 255})

    // Position camera at screen center
    camera := scene.Camera()
    camera.Position = gamemath.Vector2{X: 400, Y: 300}

    engine.SetScene(scene)

    // Load texture
    texture, err := engine.Assets().LoadTexture("player.png")
    if err != nil {
        log.Fatal(err)
    }

    // Create player sprite
    sprite := graphics.NewSprite(texture)
    sprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255})

    // Create player entity
    player := &core.Entity{
        Active: true,
        Transform: gamemath.Transform{
            Position: gamemath.Vector2{X: 400, Y: 300},
            Scale:    gamemath.Vector2{X: 2, Y: 2},
        },
        Sprite: sprite,
        Behavior: &PlayerController{
            Speed:    200,
            InputMgr: inputMgr,
        },
        Layer: 1,
    }
    scene.AddEntity(player)

    if err := engine.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### With Collision Detection

```go
// Add collision callbacks to entities
player.OnCollisionEnter = func(self, other *core.Entity) {
    log.Printf("Player hit entity %d!", other.ID)
    // Change color on collision
    if self.Sprite != nil {
        self.Sprite.SetColor(gamemath.Color{R: 255, G: 100, B: 100, A: 255})
    }
}

player.OnCollisionExit = func(self, other *core.Entity) {
    log.Printf("Player left entity %d", other.ID)
    // Restore original color
    if self.Sprite != nil {
        self.Sprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255})
    }
}

// Add collider to player
player.Collider = physics.NewCollider(64, 64) // Width, Height
player.Collider.Layer = 1 // Collision layer
```

## Examples

The `examples/` directory contains fully-functional demonstrations of engine features:

### 1. **simple/** - Minimal Example
Bare-bones 35-line example showing engine initialization and game loop.
```bash
go run examples/simple/main.go
```

### 2. **demo/** - ★ Full Feature Showcase
Comprehensive demo with 16+ entities demonstrating:
- Custom behaviors (rotating, orbiting, pulsating, bouncing)
- Color tinting and alpha blending
- Sprite flipping
- Camera following
- Layered rendering

```bash
go run examples/demo/main.go
```
See [examples/demo/README.md](examples/demo/README.md) for detailed explanation.

### 3. **moving/** - Movement Behaviors
Demonstrates 4 movement patterns:
- Linear velocity (constant speed)
- Bouncing (edge collision)
- Circular motion (orbiting)
- Wave pattern (sine wave)

```bash
go run examples/moving/main.go
```
See [examples/moving/README.md](examples/moving/README.md) for behavior implementation.

### 4. **player-control/** - Keyboard Input
Interactive player control with WASD/Arrow keys:
- Action mapping
- Input state management
- Delta-time movement

```bash
go run examples/player-control/main.go
```

### 5. **collision-callbacks/** - Collision Events
Demonstrates collision callback system:
- OnCollisionEnter (collision starts)
- OnCollisionStay (collision continues)
- OnCollisionExit (collision ends)

```bash
go run examples/collision-callbacks/main.go
```

### 6. **collision/** - Basic Collision Detection
Simple collision detection without callbacks.

```bash
go run examples/collision/main.go
```

### 7. **assets/** - Texture Loading
Demonstrates:
- PNG texture loading
- Reference counting
- Multiple sprites sharing textures
- Asset caching

```bash
go run examples/assets/main.go
```

### 8. **space-battle/** - ★ Complete Playable Game
Full-featured space shooter demonstrating production-ready game:
- Player input and control (WASD/Arrows + Space to shoot)
- Score tracking and game state management
- Text rendering for UI (score, game over)
- Entity spawning and lifecycle management
- Collision detection with game over logic
- Parallax background effects
- Restart functionality

```bash
go run examples/space-battle/main.go
```
See [examples/space-battle/README.md](examples/space-battle/README.md) for detailed gameplay guide.

**Note**: Examples that load textures require PNG files. The `assets/` and `space-battle/` examples include test textures.

## Camera System

⚠️ **Important**: The camera system uses world-to-screen coordinate transformation. By default, the camera is at (0, 0), which can cause entities to render off-screen.

**Always position the camera at screen center:**

```go
camera := scene.Camera()
camera.Position = gamemath.Vector2{
    X: float64(engine.Width()) / 2,   // 400 for 800px width
    Y: float64(engine.Height()) / 2,  // 300 for 600px height
}
```

For detailed explanation of the camera system, see [examples/CAMERA_GUIDE.md](examples/CAMERA_GUIDE.md).

## Project Structure

```
gogame/
├── engine/
│   ├── core/           # Engine, Scene, Entity, game loop
│   ├── graphics/       # Renderer, Sprite, Texture, Camera
│   ├── input/          # InputManager, actions, keycodes
│   ├── physics/        # Collision detection, Collider
│   └── math/           # Vector2, Rectangle, Transform, Color
├── examples/           # Example games and demos
│   ├── simple/         # Minimal 35-line example
│   ├── demo/           # Full feature showcase
│   ├── moving/         # Movement behaviors
│   ├── player-control/ # Keyboard input
│   ├── collision-callbacks/  # Collision events
│   ├── collision/      # Basic collision
│   ├── assets/         # Texture loading
│   └── CAMERA_GUIDE.md # Camera system documentation
├── specs/              # Technical specifications
└── tests/              # Unit, integration, and benchmark tests (planned)
```

## API Documentation

### Core Types

- **`core.Engine`** - Main engine instance managing window, rendering, and game loop
- **`core.Scene`** - Entity container with camera and background color
- **`core.Entity`** - Game object with transform, sprite, collider, and behavior
- **`core.Behavior`** - Interface for custom entity behavior (Update method)

### Key Interfaces

```go
// Behavior is called every frame to update entity logic
type Behavior interface {
    Update(entity *Entity, dt float64)
}

// CollisionCallback is called when collision events occur
type CollisionCallback func(self, other *Entity)
```

### Common Patterns

**Creating an Entity:**
```go
entity := &core.Entity{
    Active:    true,
    Transform: gamemath.Transform{Position: gamemath.Vector2{X: 400, Y: 300}},
    Sprite:    sprite,
    Collider:  physics.NewCollider(width, height),
    Behavior:  &MyBehavior{},
    Layer:     1,
}
scene.AddEntity(entity)
```

**Custom Behavior:**
```go
type MyBehavior struct {
    Speed float64
}

func (b *MyBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.X += b.Speed * dt
}
```

**Collision Detection:**
```go
entity.Collider = physics.NewCollider(64, 64)
entity.OnCollisionEnter = func(self, other *core.Entity) {
    log.Println("Collision detected!")
}
```

**Input Handling:**
```go
inputMgr := engine.Input()
inputMgr.BindAction(input.ActionJump, input.KeySpace)

if inputMgr.ActionPressed(input.ActionJump) {
    player.Jump()
}
```

**Camera Following:**
```go
func (g *Game) Update(dt float64) {
    camera := scene.Camera()
    camera.Follow(
        player.Transform.Position.X,
        player.Transform.Position.Y,
        0.1, // Smoothing (0.0 = instant, 1.0 = no follow)
    )
}
```

**Text Rendering:**
```go
// Load font
font, err := graphics.LoadFont("/System/Library/Fonts/Helvetica.ttc", 24)
if err != nil {
    log.Fatal(err)
}
defer font.Close()

// Create text renderer
textRenderer := graphics.NewTextRenderer(engine.Renderer(), font)

// Draw text on screen
err = textRenderer.DrawText("Score: 100", 10, 10, gamemath.Color{R: 255, G: 255, B: 255, A: 255})
```

## Performance

### Current Performance (on M4 Pro, 64GB RAM)

- **Rendering**: 60 FPS with 100+ sprites
- **Input Latency**: <16ms (single frame delay)
- **Asset Loading**: PNG textures load in <50ms
- **Collision Detection**: O(n²) handles ~1000 entities at 60 FPS

### Performance Targets

- 60 FPS with 500+ sprites
- <16ms input latency
- <100ms asset loading
- Memory stable over 1 hour runtime
- Collision optimization with spatial partitioning (planned)

**Note**: Performance metrics measured on M4 Pro with 64GB RAM. Performance may vary on different hardware configurations.

## Testing

```bash
# Run all tests
go test ./...

# Run unit tests with verbose output
go test ./tests/unit/... -v

# Run with coverage
go test -cover ./...

# Run benchmarks (when implemented)
go test -bench=. ./tests/benchmarks/
```

**Current Status**:
- ✅ **84 unit tests** for math components (Vector2, Rectangle, Transform, Color)
- ⏳ Integration tests for engine components (planned)
- ⏳ Benchmark tests for performance validation (planned)

## Known Issues

- **Alpha Software**: APIs are unstable and subject to change
- **No Audio**: Audio system not yet implemented
- **No Animation**: Sprite animation system not yet implemented
- **No Spatial Partitioning**: Collision detection is O(n²) - not suitable for 1000+ entities
- **No Layer Sorting**: Entities rendered in addition order, not by layer value (TODO in scene.go:311)
- **macOS Only**: Currently only supports macOS with SDL2/Metal backend

## Documentation

- **[Camera Guide](examples/CAMERA_GUIDE.md)** - Understanding the camera coordinate system
- **[Demo README](examples/demo/README.md)** - Full feature showcase explanation
- **[Space Battle README](examples/space-battle/README.md)** - Complete playable game guide
- **[Moving README](examples/moving/README.md)** - Movement behavior patterns
- **[Quickstart Guide](specs/001-macos-game-engine/quickstart.md)** - Getting started
- **[API Documentation](https://pkg.go.dev/github.com/dshills/gogame)** - Generated from code comments
- **[Architecture](specs/001-macos-game-engine/plan.md)** - Technical design document

## Troubleshooting

### SDL2 Not Found
```bash
pkg-config --modversion sdl2
```
If this fails, reinstall SDL2:
```bash
brew uninstall sdl2 sdl2_image sdl2_ttf
brew install sdl2 sdl2_image sdl2_ttf pkg-config
```

### Text Rendering Errors
If you get errors about missing SDL2_ttf:
```bash
# Install SDL2_ttf
brew install sdl2_ttf

# Verify installation
pkg-config --modversion sdl2_ttf
```

### Entities Not Visible
Check camera position! By default camera is at (0, 0). Position it at screen center:
```go
camera.Position = gamemath.Vector2{X: 400, Y: 300} // For 800x600 window
```
See [CAMERA_GUIDE.md](examples/CAMERA_GUIDE.md) for details.

### "cannot find package"
Ensure Go module is initialized:
```bash
go mod download
```

## Contributing

⚠️ **Alpha Stage**: This project is in early development. APIs are unstable.

Contributions are welcome, but be aware that:
- Breaking changes are frequent
- API design is still evolving
- Documentation may be incomplete

Please see CONTRIBUTING.md for guidelines (when available).

## License

MIT License - See LICENSE file for details.

## Credits

- Built with [SDL2](https://www.libsdl.org/) for cross-platform graphics
- Text rendering with [SDL2_ttf](https://github.com/libsdl-org/SDL_ttf)
- Uses [go-sdl2](https://github.com/veandco/go-sdl2) Go bindings
- Developed with Go 1.25.3

---

**Status**: Alpha Development | **Last Updated**: October 2025 | **Version**: 0.1.0-alpha
