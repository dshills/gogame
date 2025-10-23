# Moving Sprite Example

Demonstrates various movement behaviors with velocity-based entity motion.

## Running

```bash
go run examples/moving/main.go
```

**Note**: This example reuses textures from `examples/assets/`. Run that example first to generate the PNG files, or the sprites will be invisible (but movement still works).

## Movement Patterns Demonstrated

### 1. Linear Velocity (Blue sprite)
- **Behavior**: `VelocityBehavior`
- **Movement**: Constant velocity in X and Y directions
- **Physics**: Simple displacement = velocity × time
- **Use case**: Bullets, projectiles, constant-speed enemies

```go
type VelocityBehavior struct {
    VelocityX float64 // Pixels per second
    VelocityY float64
}
```

### 2. Bouncing (Red sprite)
- **Behavior**: `BouncingBehavior`
- **Movement**: Velocity-based with edge collision
- **Physics**: Reverses velocity on boundary collision
- **Use case**: Bouncing balls, pong-style games, confined movement

```go
type BouncingBehavior struct {
    VelocityX, VelocityY float64
    ScreenWidth, ScreenHeight float64
    Margin float64 // Collision boundary
}
```

### 3. Circular Motion (Green sprite)
- **Behavior**: `CircularMotionBehavior`
- **Movement**: Orbits around a center point
- **Physics**: Parametric circle using cos/sin
- **Use case**: Orbiting satellites, circular patrols, orbital mechanics

```go
type CircularMotionBehavior struct {
    CenterX, CenterY float64
    Radius float64
    AngularSpeed float64 // Radians per second
    CurrentAngle float64
}
```

### 4. Wave Pattern (Gold sprite)
- **Behavior**: `WavingBehavior`
- **Movement**: Sine wave with forward velocity
- **Physics**: Y = BaseY + sin(time × frequency) × amplitude
- **Use case**: Swimming fish, floating objects, wave-like patterns

```go
type WavingBehavior struct {
    VelocityX float64 // Forward speed
    Amplitude float64 // Wave height
    Frequency float64 // Cycles per second
    BaseY float64     // Center line
}
```

### 5. Multiple Small Sprites
- Semi-transparent bouncing sprites
- Different velocities and colors
- Demonstrates multiple entities with same behavior type

## Key Concepts

### Delta Time (dt)
All behaviors use delta time for frame-rate independent movement:

```go
func (vb *VelocityBehavior) Update(entity *core.Entity, dt float64) {
    entity.Transform.Position.X += vb.VelocityX * dt
    entity.Transform.Position.Y += vb.VelocityY * dt
}
```

- `dt` is time since last frame in seconds (typically ~0.016 at 60 FPS)
- Movement is consistent regardless of frame rate
- Formula: `distance = velocity × time`

### Behavior Interface
All movement types implement the `Behavior` interface:

```go
type Behavior interface {
    Update(entity *core.Entity, dt float64)
}
```

This allows:
- Pluggable movement logic
- Reusable behaviors across entities
- Clean separation of concerns

### Transform Updates
Behaviors modify the entity's `Transform`:

```go
entity.Transform.Position.X += deltaX
entity.Transform.Position.Y += deltaY
entity.Transform.Rotation = angle
```

The engine automatically renders sprites at the updated positions.

## Physics Formulas Used

**Linear Motion:**
```
position(t) = position₀ + velocity × t
```

**Circular Motion:**
```
x(t) = centerX + radius × cos(angle)
y(t) = centerY + radius × sin(angle)
angle(t) = angle₀ + angularSpeed × t
```

**Harmonic Motion (Wave):**
```
y(t) = baseY + amplitude × sin(2π × frequency × t)
```

**Elastic Collision (Bouncing):**
```
if collision: velocity = -velocity
```

## Extending This Example

### Add Acceleration

```go
type AcceleratingBehavior struct {
    VelocityX, VelocityY float64
    AccelX, AccelY float64
}

func (ab *AcceleratingBehavior) Update(entity *core.Entity, dt float64) {
    ab.VelocityX += ab.AccelX * dt
    ab.VelocityY += ab.AccelY * dt
    entity.Transform.Position.X += ab.VelocityX * dt
    entity.Transform.Position.Y += ab.VelocityY * dt
}
```

### Add Friction/Damping

```go
type DampedVelocityBehavior struct {
    VelocityX, VelocityY float64
    Friction float64 // 0.0 to 1.0
}

func (dvb *DampedVelocityBehavior) Update(entity *core.Entity, dt float64) {
    dvb.VelocityX *= (1.0 - dvb.Friction*dt)
    dvb.VelocityY *= (1.0 - dvb.Friction*dt)
    entity.Transform.Position.X += dvb.VelocityX * dt
    entity.Transform.Position.Y += dvb.VelocityY * dt
}
```

### Add Seeking/Following

```go
type SeekingBehavior struct {
    Target *core.Entity
    Speed float64
}

func (sb *SeekingBehavior) Update(entity *core.Entity, dt float64) {
    if sb.Target == nil { return }

    // Calculate direction to target
    dx := sb.Target.Transform.Position.X - entity.Transform.Position.X
    dy := sb.Target.Transform.Position.Y - entity.Transform.Position.Y
    distance := math.Sqrt(dx*dx + dy*dy)

    if distance > 0 {
        // Move toward target
        entity.Transform.Position.X += (dx/distance) * sb.Speed * dt
        entity.Transform.Position.Y += (dy/distance) * sb.Speed * dt
    }
}
```

## Performance

- **10 entities** moving simultaneously
- **60 FPS** with vsync
- **Frame-independent** movement using delta time
- **Zero allocations** in update loop (after warmup)

## Related Examples

- `examples/simple/` - Basic rendering without movement
- `examples/player-control/` - Input-driven movement (WASD)
- `examples/collision/` - Movement with collision detection
- `examples/demo/` - Comprehensive demo with all features

## Technical Notes

**Why use dt (delta time)?**
- Ensures consistent speed across different frame rates
- A sprite moving at 100 px/s will move 100 pixels in 1 second, regardless of FPS
- Formula: `new_position = old_position + velocity * dt`

**Why behaviors?**
- Reusable movement logic
- Easy to swap behaviors at runtime
- Clean entity composition
- Testable in isolation

**Why separate position update from rendering?**
- Updates happen at fixed 60 FPS (physics)
- Rendering can vary (vsync, performance)
- Clean separation of game logic and graphics

Enjoy exploring different movement patterns! 🎮
