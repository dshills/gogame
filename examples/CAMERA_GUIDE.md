# Camera Positioning Guide

## Understanding the Camera System

The gogame engine uses a **camera-based coordinate system** where:
- **World coordinates** are where you place entities
- **Screen coordinates** are where they appear on screen
- The **camera** transforms world → screen

## Default Camera Position

By default, the camera is positioned at **(0, 0)**:

```go
camera := scene.Camera()
// camera.Position = Vector2{X: 0, Y: 0}  // default
```

## Common Issue: Entities Not Visible

If you position entities at (400, 300) but the camera is at (0, 0):

```
Camera at (0, 0)
Screen is 800x600

Entity at world (400, 300)
→ appears at screen (400 + 400, 300 + 300) = (800, 600)
→ BOTTOM RIGHT CORNER!
```

## Solution: Position Camera at Screen Center

For an 800x600 window, position the camera at **(400, 300)**:

```go
camera := scene.Camera()
camera.Position = gamemath.Vector2{X: 400, Y: 300}
```

Now:
```
Entity at world (400, 300)
→ appears at screen (400, 300)
→ CENTER OF SCREEN ✓
```

## Camera Transform Formula

```
screenX = (worldX - cameraX) * zoom + screenWidth/2
screenY = (worldY - cameraY) * zoom + screenHeight/2
```

**Examples (800x600 screen, camera at 400,300, zoom 1.0):**

| World Position | Screen Position | Location |
|----------------|-----------------|----------|
| (400, 300) | (400, 300) | Center |
| (200, 300) | (200, 300) | Left of center |
| (600, 300) | (600, 300) | Right of center |
| (400, 100) | (400, 100) | Top of center |
| (400, 500) | (400, 500) | Bottom of center |

## Coordinate System

```
Screen: 800x600

    0,0 ─────────────────── 800,0
     │                        │
     │       400,300          │
     │         (center)       │
     │                        │
    0,600 ──────────────── 800,600
```

## Best Practices

### 1. Position Camera at Window Center

```go
engine, _ := core.NewEngine("Game", 800, 600, false)
scene := core.NewScene()

// Always position camera at screen center
camera := scene.Camera()
camera.Position = gamemath.Vector2{
    X: float64(engine.Width()) / 2,   // 400 for 800px width
    Y: float64(engine.Height()) / 2,  // 300 for 600px height
}
```

### 2. Or Position Entities Relative to (0, 0)

If you keep camera at (0, 0):

```go
// Entities positioned relative to origin
player := &core.Entity{
    Transform: gamemath.Transform{
        Position: gamemath.Vector2{X: 0, Y: 0},  // Center of screen
    },
}

enemy := &core.Entity{
    Transform: gamemath.Transform{
        Position: gamemath.Vector2{X: 100, Y: 0},  // Right of center
    },
}
```

### 3. Camera Following

To make camera follow the player:

```go
func (g *Game) Update(dt float64) {
    // Camera smoothly follows player
    camera.Follow(
        player.Transform.Position.X,
        player.Transform.Position.Y,
        0.1,  // Smoothing (0.0 = instant, 1.0 = no follow)
    )
}
```

### 4. Manual Camera Control

```go
// Pan camera with arrow keys
if input.KeyHeld(input.KeyArrowLeft) {
    camera.Position.X -= 200 * dt
}
if input.KeyHeld(input.KeyArrowRight) {
    camera.Position.X += 200 * dt
}
```

### 5. Zoom

```go
// Zoom in (larger objects)
camera.Zoom = 2.0  // 2x zoom

// Zoom out (smaller objects, see more)
camera.Zoom = 0.5  // 0.5x zoom

// Normal
camera.Zoom = 1.0
```

## Common Patterns

### Pattern 1: Fixed Camera (Menu/UI)
```go
camera.Position = gamemath.Vector2{X: 400, Y: 300}
// Never moves - entities move relative to screen
```

### Pattern 2: Following Camera (Action Game)
```go
// Camera follows player
camera.Follow(player.Transform.Position.X, player.Transform.Position.Y, 0.1)
```

### Pattern 3: Screen-Space UI
```go
// UI elements in world coordinates matching screen
uiButton := &core.Entity{
    Transform: gamemath.Transform{
        Position: camera.Position + gamemath.Vector2{X: 300, Y: -250},
    },
}
// Button stays at screen (700, 50) as camera moves
```

## Debugging Camera Issues

**Problem**: Entities not visible

**Check:**
1. What's the camera position?
2. What's the entity position?
3. Calculate screen position: `(entityX - cameraX) * zoom + screenW/2`

**Example Debug:**
```go
log.Printf("Camera: (%.0f, %.0f)", camera.Position.X, camera.Position.Y)
log.Printf("Entity: (%.0f, %.0f)", entity.Transform.Position.X, entity.Transform.Position.Y)

screenX, screenY := camera.WorldToScreen(
    entity.Transform.Position.X,
    entity.Transform.Position.Y,
)
log.Printf("Screen: (%d, %d)", screenX, screenY)
```

**Is it visible?**
- screenX should be between 0 and 800
- screenY should be between 0 and 600

## Examples in the Repo

All examples now correctly position the camera:

```go
// examples/simple/main.go
camera.Position = gamemath.Vector2{X: 400, Y: 300}

// examples/assets/main.go
camera.Position = gamemath.Vector2{X: 400, Y: 300}

// examples/moving/main.go
camera.Position = gamemath.Vector2{X: 400, Y: 300}

// examples/demo/main.go
camera.Position = gamemath.Vector2{X: 400, Y: 300}
```

## Quick Reference

| Window Size | Camera Position for Centered View |
|-------------|-----------------------------------|
| 800 x 600 | (400, 300) |
| 1024 x 768 | (512, 384) |
| 1280 x 720 | (640, 360) |
| 1920 x 1080 | (960, 540) |

**Formula**: `(width/2, height/2)`

---

**Remember**: The camera position is the **center** of what you're looking at in the world. Position it where you want to look!
