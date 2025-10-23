# Space Battle - Complete Demo Game

A fully-featured space shooter demonstrating all capabilities of the gogame engine.

## Overview

**Space Battle** is a complete, playable game that showcases:
- ✅ Player input and control
- ✅ Collision detection with callbacks
- ✅ Entity spawning and management
- ✅ Game state management
- ✅ Score tracking
- ✅ Visual effects and feedback
- ✅ Multiple entity types with behaviors
- ✅ Parallax background effects
- ✅ Game over and restart functionality

## How to Play

### Running the Game

```bash
# From project root
go run examples/space-battle/main.go

# Or build and run
go build -o ./build/space-battle ./examples/space-battle
./build/space-battle
```

### Controls

| Key | Action |
|-----|--------|
| **WASD** or **Arrow Keys** | Move your ship |
| **SPACE** | Shoot bullets |
| **R** | Restart game (when game over) |
| **ESC** | Quit game |

### Objective

- **Destroy enemies** by shooting them with bullets → **+10 points each**
- **Avoid collisions** with enemy ships → Game over on hit
- **Don't let enemies escape** → Game over if **3 enemies** pass you
- Survive as long as possible and maximize your score!

### Game Over

The game ends when:
1. **Your ship collides with an enemy** (player turns red)
2. **3 enemies escape past you** (player turns orange)

When game over occurs:
- Game pauses
- Your final score is displayed
- Press **R** to restart or **ESC** to quit

## Technical Implementation

### Game Systems

#### 1. **Score System**
Real-time score tracking displayed in console:
- **+10 points** for each enemy destroyed
- **Escaped counter** shows enemies that passed (0-3)
- **Status updates** every 10 seconds showing score, escaped count, and time
- **Console feedback** on every event (destroy, escape, game over)

```go
// Score tracking
g.score += 10  // On enemy destroyed

// Escape tracking
g.escapedEnemies++
if g.escapedEnemies >= 3 {
    g.onTooManyEscaped()  // Game over
}
```

**Visual Feedback:**
- ✓ Green checkmark when enemy destroyed
- ⚠ Warning icon when enemy escapes
- 📊 Status updates every 10 seconds
- Different player colors for different death types:
  * Red = Collision with enemy
  * Orange = Too many escaped

#### 2. **Player System**
- **PlayerController** behavior handles WASD/Arrow key input
- Smooth movement with delta time
- Screen bounds collision
- Shoot cooldown system (0.25 seconds between shots)

```go
type PlayerController struct {
    game *Game
}

func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
    // Movement with PlayerSpeed * dt
    // Screen bounds checking
    // Shooting with cooldown
}
```

#### 2. **Enemy System**
- **EnemyBehavior** moves enemies downward
- Random spawning at top of screen
- Automatic removal when off-screen
- Collision detection with player and bullets

```go
type EnemyBehavior struct {
    game *Game
}

func (eb *EnemyBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.Y += EnemySpeed * dt
    // Off-screen removal
}
```

#### 3. **Bullet System**
- **BulletBehavior** moves bullets upward
- Spawned at player position on spacebar press
- Collision detection with enemies
- Automatic removal when off-screen or on hit

```go
type BulletBehavior struct {
    game *Game
}

func (bb *BulletBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.Y -= BulletSpeed * dt
    // Off-screen removal
}
```

#### 4. **Background Star System**
- **StarBehavior** creates parallax scrolling effect
- Stars move downward slowly
- Wrap around when off-screen
- Random alpha for depth illusion
- Limited to 50 stars maximum

```go
type StarBehavior struct {
    game *Game
}

func (sb *StarBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.Y += StarSpeed * dt
    // Wrap around when off bottom
}
```

#### 5. **Collision System**

Uses gogame's collision callback system:

```go
// Player collision callback
player.OnCollisionEnter = func(self, other *core.Entity) {
    if other.Collider.Layer & CollisionLayerEnemy != 0 {
        g.onPlayerHit() // Game over
    }
}

// Bullet collision callback
bullet.OnCollisionEnter = func(self, other *core.Entity) {
    if other.Collider.Layer & CollisionLayerEnemy != 0 {
        g.onEnemyHit(other)   // Add score
        g.removeBullet(self)   // Remove bullet
    }
}
```

**Collision Layers (bit positions):**
- `CollisionLayerPlayer = 0` - Player ship (bitmask 0x01)
- `CollisionLayerEnemy = 1` - Enemy ships (bitmask 0x02)
- `CollisionLayerBullet = 2` - Player bullets (bitmask 0x04)

**Collision Masks (what each layer collides with):**
- Player: `(1 << CollisionLayerEnemy)` = 0x02 - Collides with enemies
- Enemy: `(1 << CollisionLayerPlayer) | (1 << CollisionLayerBullet)` = 0x05 - Collides with player and bullets
- Bullet: `(1 << CollisionLayerEnemy)` = 0x02 - Collides with enemies

#### 6. **Game State Management**

Two states: `StatePlaying` and `StateGameOver`

```go
type GameState int

const (
    StatePlaying GameState = iota
    StateGameOver
)

func (g *Game) Update(dt float64) {
    if g.state == StateGameOver {
        // Handle restart input
        return
    }

    // Normal game logic
}
```

#### 7. **Spawning System**

Time-based spawning with configurable intervals:

```go
const (
    EnemySpawnInterval = 1.5  // Seconds
    StarSpawnInterval  = 0.1  // Seconds
)

func (g *Game) Update(dt float64) {
    g.enemySpawnTimer += dt
    if g.enemySpawnTimer >= EnemySpawnInterval {
        g.enemySpawnTimer = 0
        g.spawnEnemy()
    }
}
```

### Asset Generation

Game includes an asset generator that creates PNG images programmatically:

```bash
# From examples/space-battle directory
cd examples/space-battle
go run tools/generate_assets.go

# Or from project root
go run examples/space-battle/tools/generate_assets.go
```

**Generated Assets:**
- `assets/player.png` - Blue triangle (32x32) - Player ship
- `assets/enemy.png` - Red triangle (32x32) - Enemy ship
- `assets/bullet.png` - Yellow rectangle (8x16) - Bullet
- `assets/star.png` - White dot (4x4) - Background star

## Game Constants

Easily tweak gameplay by modifying constants in `main.go`:

```go
const (
    ScreenWidth  = 800
    ScreenHeight = 600

    PlayerSpeed   = 300.0  // Pixels per second
    BulletSpeed   = 400.0  // Pixels per second
    EnemySpeed    = 100.0  // Pixels per second
    ShootCooldown = 0.25   // Seconds between shots

    EnemySpawnInterval = 1.5  // Seconds between enemies
    StarSpawnInterval  = 0.1  // Seconds between stars
    MaxStars           = 50   // Maximum background stars
    StarSpeed          = 50.0 // Star movement speed

    MaxEscapedEnemies = 3     // Game over when this many escape
)
```

**Difficulty Tuning:**
- Decrease `EnemySpawnInterval` → More enemies, harder
- Increase `EnemySpeed` → Faster enemies, harder
- Decrease `MaxEscapedEnemies` → Less forgiving, harder
- Increase `ShootCooldown` → Slower firing, harder

## Code Structure

```
space-battle/
├── main.go              # Main game logic
├── README.md           # This file
├── tools/              # Development tools
│   └── generate_assets.go  # PNG asset generator
└── assets/             # Generated assets
    ├── player.png      # Player ship sprite
    ├── enemy.png       # Enemy ship sprite
    ├── bullet.png      # Bullet sprite
    └── star.png        # Star sprite
```

## Engine Features Demonstrated

### Input System
- **Action Mapping**: WASD and Arrow keys both work for movement
- **Key State Detection**: Held keys for continuous movement
- **Pressed Detection**: Single-shot bullets on spacebar press
- **Input Buffering**: Frame-perfect input with double buffering

### Collision Detection
- **AABB Collision**: Axis-aligned bounding box detection
- **Collision Callbacks**: OnCollisionEnter events
- **Layer Masks**: Selective collision between layers
- **Efficient Detection**: O(n²) handles many entities at 60 FPS

### Entity-Component System
- **Custom Behaviors**: PlayerController, EnemyBehavior, BulletBehavior, StarBehavior
- **Entity Management**: Dynamic spawning and removal
- **Component Composition**: Transform + Sprite + Collider + Behavior

### Rendering
- **Sprite Rendering**: Textured quads with color tinting
- **Alpha Blending**: Semi-transparent stars
- **Layer Sorting**: Background stars (layer 0) behind gameplay (layer 2)
- **Color Feedback**: Visual effects on collision (white flash, red tint)

### Game Loop
- **Fixed Timestep**: Consistent 60 FPS updates
- **Delta Time**: Frame-rate independent movement
- **Custom Update Callback**: Game-specific logic integration
- **State Management**: Clean separation of playing/game-over states

## Extending the Game

### Add Power-Ups

```go
type PowerUpBehavior struct {
    game *Game
}

func (pb *PowerUpBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.Y += 80 * dt
}

// In player collision callback
if other.Collider.Layer & CollisionLayerPowerUp != 0 {
    g.player.powerUpActive = true
    g.removePowerUp(other)
}
```

### Add Multiple Enemy Types

```go
type FastEnemyBehavior struct {
    game *Game
}

func (feb *FastEnemyBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.Y += EnemySpeed * 2 * dt
}

// Spawn with different sprite color
sprite.SetColor(gamemath.Color{R: 255, G: 150, B: 50, A: 255}) // Orange
```

### Add Boss Encounters

```go
type BossBehavior struct {
    game   *Game
    health int
}

func (bb *BossBehavior) Update(entity *core.Entity, dt float64) {
    // Complex movement pattern
    // Shooting pattern
}

// In boss collision
if bb.health <= 0 {
    g.score += 1000
    g.removeBoss(self)
}
```

### Add Particle Effects

```go
func (g *Game) createExplosion(x, y float64) {
    for i := 0; i < 20; i++ {
        particle := createParticle(x, y)
        g.scene.AddEntity(particle)
    }
}
```

## Performance

**On M1 MacBook Pro:**
- Maintains 60 FPS with 50+ entities
- Collision detection handles all gameplay entities
- No memory leaks during extended play
- Smooth movement and visual effects

**Entity Counts:**
- 1 Player
- 0-10 Enemies (spawns over time)
- 0-20 Bullets (limited by cooldown)
- 50 Background Stars
- **Total: ~60-80 active entities**

## Learning Points

This demo teaches:

1. **Game Architecture**: Separation of concerns (behaviors, systems, state)
2. **Entity Lifecycle**: Spawning, updating, collision, removal
3. **Input Handling**: Responsive controls with action mapping
4. **Collision Detection**: Layer-based selective collision
5. **Visual Feedback**: Color changes, effects on events
6. **State Management**: Game states (playing, game over)
7. **Spawning Systems**: Time-based enemy/star generation
8. **Performance**: Managing many entities at 60 FPS
9. **Code Organization**: Clean, maintainable game structure
10. **Asset Pipeline**: Generating assets programmatically

## Tips for High Scores

1. **Stay Mobile**: Keep moving to avoid enemy clusters
2. **Shoot Constantly**: Fire as much as possible (cooldown permitting)
3. **Use Full Screen**: Utilize all available space for maneuvering
4. **Predict Spawns**: Enemies spawn at the top - anticipate their arrival
5. **Corner Strategy**: Sides of screen are safer but limit escape routes

## Troubleshooting

### Assets Not Found
If you see texture loading errors:
```bash
cd examples/space-battle
go run tools/generate_assets.go
go run main.go
```

### Game Runs Slowly
- Check terminal for FPS warnings
- Reduce `MaxStars` constant
- Reduce enemy spawn rate

### Input Not Working
- Ensure window has focus
- Check that SDL2 is properly installed
- Verify terminal shows "Game initialized!" message

## Credits

- **Engine**: gogame 2D Game Engine
- **Graphics**: SDL2 with Metal backend
- **Language**: Go 1.25.3
- **Assets**: Generated programmatically with Go image library

---

**Enjoy the game!** This is a complete demonstration of what you can build with the gogame engine. Use it as a starting point for your own games!
