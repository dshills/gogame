# gogame - 2D Game Engine for macOS

A simple, performant 2D game engine written in Go for macOS, designed for creating arcade-style games and platformers.

## Features

✅ **Currently Implemented (MVP)**:
- **60 FPS Rendering**: Hardware-accelerated sprite rendering with SDL2/Metal backend
- **Entity Management**: Component-based entity system with custom behaviors
- **Transform System**: Position, rotation, and scale with smooth interpolation
- **Asset Loading**: PNG/JPEG texture loading with reference counting and caching
- **Visual Effects**: Color tinting, alpha blending, sprite flipping
- **Camera System**: World-to-screen transforms with zoom and follow
- **Scene Management**: Entity containers with background colors and layers
- **Simple API**: Create games in under 50 lines of code

🚧 **Planned**:
- **Input Handling**: Keyboard and mouse input with action mapping
- **Collision Detection**: AABB collision detection with layer masks
- **Audio System**: Sound effects and music playback

## Prerequisites

### macOS Requirements

- macOS 12.0 (Monterey) or newer
- Homebrew package manager
- Go 1.25.3 or newer

### Install SDL2

```bash
# Install SDL2 and SDL2_image via Homebrew
brew install sdl2 sdl2_image pkg-config

# Verify installation
pkg-config --modversion sdl2
```

## Installation

```bash
go get github.com/dshills/gogame
```

## Quick Start

```go
package main

import (
    "log"
    "runtime"
    "github.com/dshills/gogame/engine/core"
    "github.com/dshills/gogame/engine/graphics"
    "github.com/dshills/gogame/engine/math"
)

func main() {
    // CRITICAL: SDL requires running on the main OS thread
    runtime.LockOSThread()

    // Create engine with 800x600 window
    eng, err := core.NewEngine("My Game", 800, 600, false)
    if err != nil {
        log.Fatal(err)
    }
    defer eng.Shutdown()

    // Create scene
    scene := core.NewScene()
    eng.SetScene(scene)

    // Load texture and create sprite
    texture, _ := eng.Assets().LoadTexture("player.png")
    sprite := graphics.NewSprite(texture)

    // Create entity at center of screen
    player := &core.Entity{
        Active:    true,
        Transform: math.Transform{Position: math.Vector2{X: 400, Y: 300}},
        Sprite:    sprite,
        Layer:     1,
    }
    scene.AddEntity(player)

    // Run game loop
    if err := eng.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Project Structure

```
gogame/
├── engine/
│   ├── core/        # Engine, Scene, game loop
│   ├── graphics/    # Rendering, sprites, textures, camera
│   ├── entity/      # Entity, Transform, Behavior
│   ├── input/       # InputManager, actions, keycodes
│   ├── physics/     # Collision detection
│   └── math/        # Vector2, Rectangle, Transform, Color
├── examples/        # Example games
└── tests/           # Unit, integration, and benchmark tests
```

## Examples

### Feature Demo (Recommended Starting Point)

Run the comprehensive demo showcasing all engine features:

```bash
go run examples/demo/main.go
```

**Features demonstrated**: 16 entities with rotating, orbiting, pulsating, bouncing behaviors, color tinting, alpha blending, sprite flipping, camera system, and smooth following. See [examples/demo/README.md](examples/demo/README.md) for details.

### Additional Examples

- `simple/` - Minimal 35-line example showing engine initialization
- `demo/` - **★ Full feature showcase with custom behaviors** (see above)

🚧 **Coming Soon**:
- `moving/` - Moving sprites with velocity
- `player-control/` - Keyboard input handling
- `assets/` - Loading multiple textures
- `collision/` - Collision detection

## Running Tests

```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=. ./tests/benchmarks/

# Run with coverage
go test -cover ./...
```

## Performance Targets

- 60 FPS with 100+ sprites
- <16ms input latency
- <100ms asset loading
- Memory stable over 1 hour runtime

## Documentation

- [Quickstart Guide](specs/001-macos-game-engine/quickstart.md)
- [API Documentation](https://pkg.go.dev/github.com/dshills/gogame)
- [Architecture](specs/001-macos-game-engine/plan.md)

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please see CONTRIBUTING.md for guidelines.
