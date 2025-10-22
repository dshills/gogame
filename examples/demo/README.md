# gogame Feature Demo

A comprehensive demonstration of all current gogame engine features.

## Features Demonstrated

### 1. **Sprite Rendering**
- Loads PNG textures using the AssetManager
- Renders sprites with SDL2 hardware acceleration
- Demonstrates texture caching and reference counting

### 2. **Transform System**
- **Position**: Entities positioned anywhere in world space
- **Rotation**: Smooth continuous rotation (see center player)
- **Scale**: 2x scaled player entity demonstrates scaling

### 3. **Visual Effects**
- **Color Tinting**: Rainbow-colored orbiting enemies
- **Alpha Blending**: Pulsating collectibles with animated transparency
- **Sprite Flipping**: Bouncing entities with FlipH/FlipV

### 4. **Custom Behaviors**
The demo includes 5 different behavior implementations:

- **RotatingBehavior**: Continuous rotation (center player)
- **OrbitingBehavior**: Circular orbits with varying speeds (6 enemies)
- **PulsatingBehavior**: Alpha animation (4 corner collectibles)
- **BouncingBehavior**: Physics-based bouncing (4 bouncers)
- **SmoothFollowBehavior**: Interpolated following (purple follower)

### 5. **Camera System**
- World-to-screen coordinate transformation
- Camera positioning and zoom support
- Ready for smooth camera following (easily extensible)

### 6. **Entity Management**
- 16 active entities simultaneously
- Different layer values for z-ordering
- Entity lifecycle management (add/remove)

### 7. **Game Loop**
- Fixed 60 FPS update rate with accumulator pattern
- Variable render rate with vsync
- Independent update and render phases

### 8. **Scene Management**
- Custom background color
- Centralized entity container
- Deferred entity removal (safe during updates)

## Running the Demo

```bash
# From project root
go run examples/demo/main.go
```

The demo will automatically:
1. Create the `examples/demo/assets` directory
2. Generate simple colored test textures (64x64 PNG files)
3. Launch a 1024x768 window
4. Run the demo until you close the window

## What You'll See

- **Center**: A large rotating player sprite (blue tint, 2x scale)
- **Orbiting**: 6 rainbow-colored enemies orbiting the player at different speeds
- **Corners**: 4 gold collectibles pulsating with alpha animation
- **Bouncing**: 4 semi-transparent entities bouncing around the screen (some flipped)
- **Following**: 1 purple sprite smoothly following the player

**Total**: 16 entities demonstrating all engine features simultaneously at 60 FPS.

## Code Structure

The demo is self-contained in `main.go` and includes:

- **Texture Generation**: Creates test PNG files programmatically
- **Behavior Implementations**: 5 custom behavior types showing the Behavior interface
- **Entity Setup**: Demonstrates all entity configuration options
- **Scene Configuration**: Shows scene initialization and setup

## Performance

The demo runs at 60 FPS with 16 active entities, demonstrating:
- Hardware-accelerated rendering
- Efficient update loop
- Zero-allocation game loop (after warmup)
- Fixed timestep ensures consistent physics

## Extending the Demo

To add your own entities:

```go
// Create sprite
mySprite := graphics.NewSprite(myTexture)
mySprite.SetColor(gamemath.Color{R: 255, G: 0, B: 0, A: 255})
mySprite.Alpha = 0.8

// Create entity
myEntity := &core.Entity{
    Active: true,
    Transform: gamemath.Transform{
        Position: gamemath.Vector2{X: 100, Y: 100},
        Rotation: 0,
        Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
    },
    Sprite:   mySprite,
    Behavior: &MyCustomBehavior{}, // Optional
    Layer:    1,
}

// Add to scene
scene.AddEntity(myEntity)
```

## Next Steps

This demo shows the foundation. Future features to explore:
- Input handling (keyboard/mouse)
- Collision detection (AABB)
- Audio playback
- Particle systems
- Sprite animations

See the [main README](../../README.md) for full engine documentation.
