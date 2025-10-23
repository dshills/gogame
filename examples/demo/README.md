# gogame Engine - Comprehensive Feature Demo

This demo showcases **all** current features of the gogame engine in an interactive game experience.

## Running the Demo

```bash
go run examples/demo/main.go
```

The demo will automatically generate placeholder textures if they don't exist.

## Game Objective

🎯 **Collect all 4 golden collectibles while avoiding the red patrolling enemies!**

## Controls

- **WASD** or **Arrow Keys** - Move the blue player
- **ESC** - Print debug information (position, score, entity count)

## Features Demonstrated

### ✅ Core Game Loop (User Story 1)
- Fixed 60 FPS update rate with delta time
- Vsync rendering
- Frame-independent physics and movement

### ✅ Entity/Scene Management (User Story 2)
- Dynamic entity creation and removal
- Scene-based entity organization
- Entity behaviors for custom logic
- Layer-based rendering

### ✅ Input Handling (User Story 3)
- **Action mapping** - Multiple keys bound to actions
- **State tracking** - Pressed, held, released detection
- **WASD movement** - Smooth player control
- **Keyboard input** - ESC for debug info

### ✅ Asset Loading (User Story 4)
- **Texture loading** - PNG support with caching
- **Reference counting** - Efficient memory management
- **Asset manager** - Centralized texture management
- Multiple textures loaded and shared across entities

### ✅ Collision Detection (User Story 5)
- **AABB collision** - Axis-aligned bounding box testing
- **Layer masks** - Selective collision filtering
  - Player (layer 0) collides with enemies (layer 1) and collectibles (layer 2)
  - Walls (layer 3) block player and enemies
- **Trigger colliders** - Collectibles don't block movement
- **Dynamic collision** - Real-time detection during gameplay
- **Collision feedback** - Console messages on collision events

### ✅ Sprite Rendering
- **Texture mapping** - Sprites from loaded textures
- **Transform support** - Position, rotation, scale
- **Color tinting** - Different colors per sprite
- **Alpha blending** - Semi-transparent and pulsating effects
- **Sprite flipping** - Horizontal and vertical

### ✅ Camera System
- **World-space rendering** - Camera transforms all entities
- **Camera positioning** - Configurable view center
- **Zoom support** - Scalable view (currently 1.0)

### ✅ Behavior System
- **PlayerController** - Input-driven movement and rotation
- **EnemyPatrol** - Autonomous back-and-forth movement
- **CollectibleBehavior** - Pulsating scale and alpha animation
- **BehaviorFunc** - Function adapter for inline behaviors

## Entity Types in Demo

| Entity Type | Count | Features | Collision Layer |
|-------------|-------|----------|----------------|
| Player | 1 | WASD control, collision detection, rotation | 0 |
| Enemies | 3 | Patrol movement, collision detection, rotation | 1 |
| Collectibles | 4 | Pulsating animation, trigger collision, removal on collect | 2 |
| Walls | 6 | Static obstacles, blocking collision | 3 |

## Collision Matrix

```
           Player  Enemy  Collectible  Wall
Player        -      ✓         ✓        ✓
Enemy         ✓      -         -        ✓
Collectible   ✓      -         -        -
Wall          ✓      ✓         -        -
```

## Code Highlights

### Input Action Mapping
```go
inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
```

### Collision Detection with Layer Masks
```go
playerEntity.Collider.CollisionLayer = 0                    // Player layer
playerEntity.Collider.CollisionMask = (1 << 1) | (1 << 2)  // Collides with enemies and collectibles

if entity.Collider.Intersects(collectible.GetCollider(), entity.Transform, collectible.Transform) {
    // Collision detected!
}
```

### Asset Loading with Reference Counting
```go
playerTexture, _ := assets.LoadTexture("examples/demo/assets/player.png")
// Texture is cached and reference counted automatically
```

### Custom Behaviors
```go
type EnemyPatrol struct {
    Speed, MinX, MaxX, Direction float64
}

func (ep *EnemyPatrol) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.X += ep.Direction * ep.Speed * dt
    // Boundary checking and direction reversal
}
```

## Performance

- **60 FPS** - Smooth fixed timestep updates
- **~20 entities** - Player, enemies, collectibles, walls
- **O(n²) collision** - Simple broad-phase (suitable for <1000 entities)
- **Efficient rendering** - Hardware-accelerated SDL2 with vsync

## Architecture Patterns

- **Entity-Component System** - Entities with optional Sprite, Collider, Behavior
- **Behavior Pattern** - Interface for custom entity logic
- **Scene Management** - Deferred entity removal for safe iteration
- **Asset Caching** - Reference-counted texture management
- **Layer Masks** - Bitfield-based collision filtering

## Next Steps

This demo is a complete example of the engine's capabilities. To build your own game:

1. Study the behavior implementations (PlayerController, EnemyPatrol, etc.)
2. Examine the collision detection setup with layer masks
3. Look at how input actions are bound and used
4. See how entities are created with sprites, colliders, and behaviors
5. Understand the scene/entity lifecycle

Refer to the individual example directories for focused demonstrations:
- `examples/simple/` - Minimal engine setup
- `examples/player-control/` - Input handling focus
- `examples/assets/` - Asset loading focus
- `examples/collision/` - Collision detection focus

## Technical Details

**Engine Version:** 0.1.0
**Go Version:** 1.25.3+
**Dependencies:** SDL2 (via go-sdl2)
**Platform:** macOS (with Metal backend)

Enjoy exploring the gogame engine! 🎮
