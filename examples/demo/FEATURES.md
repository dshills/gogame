# Feature Demonstration Guide

This document provides a comprehensive breakdown of all engine features demonstrated in the demo.

## Demo Statistics

- **Total Entities**: 16 active entities
- **Behaviors Implemented**: 5 custom behavior types
- **Texture Files**: 3 auto-generated PNG files (64x64 each)
- **Target FPS**: 60 FPS (fixed timestep)
- **Window Size**: 1024x768 pixels
- **Lines of Code**: ~350 lines (including comments)

## Feature Breakdown

### 1. Sprite Rendering System

**Implementation**: `engine/graphics/sprite.go`, `engine/graphics/texture.go`, `engine/graphics/renderer.go`

**Demonstrated by**:
- All 16 entities use sprite rendering
- Hardware-accelerated via SDL2
- Automatic Metal backend on macOS

**Code example from demo**:
```go
playerSprite := graphics.NewSprite(playerTexture)
```

### 2. Asset Management

**Implementation**: `engine/graphics/assets.go`

**Demonstrated by**:
- Loading 3 textures: player.png, enemy.png, collectible.png
- Texture caching (multiple entities share textures)
- Reference counting (prevents premature unloading)

**Code example from demo**:
```go
texture, err := engine.Assets().LoadTexture("examples/demo/assets/player.png")
```

**Features**:
- ✅ PNG/JPEG support via Go standard library
- ✅ Automatic caching
- ✅ Reference counting
- ✅ Lazy loading

### 3. Transform System

**Implementation**: `engine/math/transform.go`

**Demonstrated by**:
- **Position**: All entities positioned in world space
- **Rotation**: Center player rotates at 45°/second
- **Scale**: Player scaled 2x, collectibles scaled 1.5x

**Code example from demo**:
```go
Transform: gamemath.Transform{
    Position: gamemath.Vector2{X: 512, Y: 384},
    Rotation: 0,
    Scale:    gamemath.Vector2{X: 2.0, Y: 2.0},
}
```

**Transform methods used**:
- `Translate(dx, dy)` - Used by bouncing entities
- `Rotate(degrees)` - Used by rotating behavior

### 4. Visual Effects

**Implementation**: `engine/graphics/sprite.go` (Color, Alpha, FlipH, FlipV fields)

#### 4.1 Color Tinting
**Demonstrated by**:
- Player: Light blue tint (100, 200, 255)
- Enemies: Rainbow colors (red, green, yellow, magenta, cyan, orange)
- Collectibles: Yellow/gold tint (255, 255, 100)
- Bouncers: Blue tint (150, 150, 255)
- Follower: Purple tint (255, 150, 255)

**Code example**:
```go
sprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255})
```

#### 4.2 Alpha Blending
**Demonstrated by**:
- Pulsating collectibles: Alpha animates 0.65 ↔ 1.0 using sine wave
- Bouncers: Semi-transparent at alpha 200/255 (~78%)
- Follower: Semi-transparent at alpha 180/255 (~70%)

**Code example**:
```go
sprite.Alpha = 0.65 + 0.35*math.Sin(pb.CurrentPhase)
```

#### 4.3 Sprite Flipping
**Demonstrated by**: Bouncing entities
- Bouncers 0 & 2: FlipH = true (horizontal flip)
- Bouncers 2 & 3: FlipV = true (vertical flip)
- Bouncer 2: Both flips (rotated 180°)

**Code example**:
```go
sprite.FlipH = true  // Mirror horizontally
sprite.FlipV = true  // Mirror vertically
```

### 5. Custom Behavior System

**Implementation**: `engine/core/entity.go` (Behavior interface)

The demo implements 5 distinct behaviors showing the flexibility of the system:

#### 5.1 RotatingBehavior
**Entity**: Center player
**Effect**: Continuous rotation at 45°/second
**Code**:
```go
type RotatingBehavior struct {
    RotationSpeed float64
}

func (rb *RotatingBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Rotate(rb.RotationSpeed * dt)
}
```

#### 5.2 OrbitingBehavior
**Entities**: 6 rainbow enemies
**Effect**: Circular orbits with varying speeds
**Math**: Uses polar coordinates (r, θ) converted to Cartesian
**Code**:
```go
func (ob *OrbitingBehavior) Update(entity *core.Entity, dt float64) {
    ob.CurrentAngle += ob.Speed * dt
    entity.Transform.Position.X = ob.CenterX + math.Cos(ob.CurrentAngle)*ob.Radius
    entity.Transform.Position.Y = ob.CenterY + math.Sin(ob.CurrentAngle)*ob.Radius
}
```

#### 5.3 PulsatingBehavior
**Entities**: 4 corner collectibles
**Effect**: Animated alpha using sine wave
**Code**:
```go
func (pb *PulsatingBehavior) Update(entity *core.Entity, dt float64) {
    pb.CurrentPhase += pb.Speed * dt
    entity.Sprite.Alpha = 0.65 + 0.35*math.Sin(pb.CurrentPhase)
}
```

#### 5.4 BouncingBehavior
**Entities**: 4 bouncing sprites
**Effect**: Physics-based bouncing with wall collisions
**Code**:
```go
func (bb *BouncingBehavior) Update(entity *core.Entity, dt float64) {
    // Update position
    entity.Transform.Position.X += bb.VelocityX * dt
    entity.Transform.Position.Y += bb.VelocityY * dt

    // Bounce off walls (invert velocity)
    if entity.Transform.Position.X < 0 || entity.Transform.Position.X > bb.ScreenWidth {
        bb.VelocityX = -bb.VelocityX
    }
}
```

#### 5.5 SmoothFollowBehavior
**Entity**: Purple follower
**Effect**: Smooth interpolation toward target
**Code**:
```go
func (sf *SmoothFollowBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.X += (sf.Target.Transform.Position.X - entity.Transform.Position.X) * sf.Smoothing
    entity.Transform.Position.Y += (sf.Target.Transform.Position.Y - entity.Transform.Position.Y) * sf.Smoothing
}
```

### 6. Camera System

**Implementation**: `engine/graphics/camera.go`

**Demonstrated by**: Scene camera at center with zoom 1.0

**Features**:
- ✅ World-to-screen coordinate transformation
- ✅ Screen-to-world coordinate transformation
- ✅ Zoom support (1.0 = normal)
- ✅ Smooth camera following (via Follow method)

**Code in demo**:
```go
camera := scene.Camera()
camera.Position = gamemath.Vector2{X: 512, Y: 384}
camera.Zoom = 1.0
```

**Transform equations**:
```
screenX = (worldX - cameraX) * zoom + screenWidth/2
worldX = (screenX - screenWidth/2) / zoom + cameraX
```

### 7. Scene Management

**Implementation**: `engine/core/scene.go`

**Features demonstrated**:
- Background color: Dark blue-gray (25, 25, 40)
- 16 entities in a single scene
- Layer-based rendering (layers 1-6)
- Entity lifecycle management

**Code example**:
```go
scene := core.NewScene()
scene.SetBackgroundColor(gamemath.Color{R: 25, G: 25, B: 40, A: 255})
scene.AddEntity(player)
```

### 8. Game Loop Architecture

**Implementation**: `engine/core/time.go`, `engine/core/engine.go`

**Features**:
- ✅ Fixed 60 FPS update rate (16.67ms per update)
- ✅ Accumulator pattern ("Fix Your Timestep")
- ✅ Variable render rate with vsync
- ✅ Spiral of death prevention (max 0.25s catch-up)

**Frame timing**:
```
Target: 60 FPS = 16.67ms/frame
Update: Fixed 16.67ms timestep
Render: Variable (limited by vsync)
```

### 9. Entity Management

**Implementation**: `engine/core/scene.go`, `engine/core/entity.go`

**Features demonstrated**:
- 16 simultaneous entities
- Entity IDs auto-assigned (1-16)
- Active/inactive state control
- Layer-based z-ordering (1-6)
- Deferred removal (safe during updates)

**Entity composition**:
```go
type Entity struct {
    ID        uint64              // Auto-assigned
    Active    bool                // true = updates/renders
    Transform gamemath.Transform  // Position, rotation, scale
    Sprite    *graphics.Sprite    // Optional visual
    Collider  *gamemath.Rectangle // Optional (for future collision)
    Behavior  Behavior            // Optional custom logic
    Layer     int                 // Z-order (higher = on top)
}
```

## Performance Characteristics

### Rendering Performance
- **Target**: 60 FPS with 100+ sprites
- **Demo**: 16 sprites at 60 FPS (83% headroom)
- **Backend**: SDL2 hardware acceleration + Metal

### Memory Management
- **Texture caching**: 3 textures loaded once, shared by 16 entities
- **Reference counting**: Prevents premature unloading
- **Zero allocation game loop**: After warmup, no GC pressure

### Update Performance
- **Fixed timestep**: Consistent physics regardless of render speed
- **Efficient behaviors**: Direct float operations, no reflection
- **Cache-friendly**: Entities stored contiguously in slice

## Extending the Demo

### Adding a New Behavior

```go
type MyBehavior struct {
    MyField float64
}

func (mb *MyBehavior) Update(entity *core.Entity, dt float64) {
    // Your custom logic using dt (delta time)
    entity.Transform.Position.X += mb.MyField * dt
}

// Use it
myEntity.Behavior = &MyBehavior{MyField: 100.0}
```

### Adding New Visual Effects

```go
// Color tint
sprite.SetColor(gamemath.Color{R: 255, G: 0, B: 0, A: 255})

// Semi-transparent
sprite.Alpha = 0.5

// Flip sprite
sprite.FlipH = true  // Mirror horizontally
sprite.FlipV = true  // Mirror vertically

// Source rect (for sprite sheets)
sprite.SetSourceRect(64, 0, 32, 32)  // Extract 32x32 region
```

### Adding New Entities

```go
// Load texture
texture, _ := engine.Assets().LoadTexture("my_sprite.png")

// Create sprite with effects
sprite := graphics.NewSprite(texture)
sprite.SetColor(gamemath.Color{R: 200, G: 100, B: 255, A: 255})
sprite.Alpha = 0.8

// Create entity
entity := &core.Entity{
    Active:    true,
    Transform: gamemath.Transform{
        Position: gamemath.Vector2{X: 400, Y: 300},
        Rotation: 45,
        Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
    },
    Sprite:    sprite,
    Behavior:  &MyCustomBehavior{},
    Layer:     3,
}

// Add to scene
scene.AddEntity(entity)
```

## Technical Notes

### Coordinate System
- **Origin**: Top-left (0, 0)
- **X-axis**: Left to right (positive)
- **Y-axis**: Top to bottom (positive)
- **Rotation**: Degrees, clockwise (0° = right, 90° = down)

### Color Format
- **Type**: RGBA (Red, Green, Blue, Alpha)
- **Range**: 0-255 per channel
- **Alpha**: 0 = transparent, 255 = opaque

### Sprite Anchors
- Sprites are centered at their transform position
- Useful for rotation (rotates around center)
- Adjust position to align top-left if needed

### Performance Tips
1. **Share textures**: Multiple entities can use same texture
2. **Batch by layer**: Entities on same layer render together
3. **Minimize behaviors**: Only attach behaviors when needed
4. **Reuse entities**: Deactivate instead of remove/add
5. **Sprite flipping**: Use FlipH/FlipV instead of loading mirrored textures

## Next Steps

This demo provides the foundation. Future additions:
- **Input system**: Keyboard/mouse for player control
- **Collision detection**: AABB collision with callbacks
- **Particle systems**: Sprite emitters for effects
- **Animation system**: Sprite sheet frame sequencing
- **Audio**: Sound effects and background music
- **UI system**: Text rendering and GUI widgets

See the [main README](../../README.md) for roadmap and documentation.
