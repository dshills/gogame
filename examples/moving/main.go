// Package main demonstrates moving sprites with velocity-based behaviors.
package main

import (
	"log"
	"math"
	"runtime"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
)

// VelocityBehavior moves an entity with constant velocity.
type VelocityBehavior struct {
	VelocityX float64 // Pixels per second in X direction
	VelocityY float64 // Pixels per second in Y direction
}

// Update moves the entity based on velocity and delta time.
func (vb *VelocityBehavior) Update(entity *core.Entity, dt float64) {
	entity.Transform.Position.X += vb.VelocityX * dt
	entity.Transform.Position.Y += vb.VelocityY * dt
}

// BouncingBehavior moves an entity and bounces it off screen edges.
type BouncingBehavior struct {
	VelocityX    float64
	VelocityY    float64
	ScreenWidth  float64
	ScreenHeight float64
	Margin       float64 // Distance from edge to bounce
}

// Update moves entity and bounces off edges.
func (bb *BouncingBehavior) Update(entity *core.Entity, dt float64) {
	// Update position
	entity.Transform.Position.X += bb.VelocityX * dt
	entity.Transform.Position.Y += bb.VelocityY * dt

	// Bounce off left/right edges
	if entity.Transform.Position.X < bb.Margin {
		entity.Transform.Position.X = bb.Margin
		bb.VelocityX = -bb.VelocityX
	} else if entity.Transform.Position.X > bb.ScreenWidth-bb.Margin {
		entity.Transform.Position.X = bb.ScreenWidth - bb.Margin
		bb.VelocityX = -bb.VelocityX
	}

	// Bounce off top/bottom edges
	if entity.Transform.Position.Y < bb.Margin {
		entity.Transform.Position.Y = bb.Margin
		bb.VelocityY = -bb.VelocityY
	} else if entity.Transform.Position.Y > bb.ScreenHeight-bb.Margin {
		entity.Transform.Position.Y = bb.ScreenHeight - bb.Margin
		bb.VelocityY = -bb.VelocityY
	}
}

// CircularMotionBehavior moves entity in a circular path.
type CircularMotionBehavior struct {
	CenterX      float64 // Center of circular path
	CenterY      float64 // Center of circular path
	Radius       float64 // Radius of circle
	AngularSpeed float64 // Radians per second
	CurrentAngle float64 // Current angle in radians
}

// Update moves entity along circular path.
func (cm *CircularMotionBehavior) Update(entity *core.Entity, dt float64) {
	cm.CurrentAngle += cm.AngularSpeed * dt

	// Keep angle in [0, 2π] range
	if cm.CurrentAngle > 2*math.Pi {
		cm.CurrentAngle -= 2 * math.Pi
	}

	// Update position
	entity.Transform.Position.X = cm.CenterX + math.Cos(cm.CurrentAngle)*cm.Radius
	entity.Transform.Position.Y = cm.CenterY + math.Sin(cm.CurrentAngle)*cm.Radius

	// Rotate sprite to face direction of movement
	entity.Transform.Rotation = cm.CurrentAngle * (180 / math.Pi)
}

// WavingBehavior moves entity in a sine wave pattern.
type WavingBehavior struct {
	VelocityX float64 // Forward speed (pixels/second)
	Amplitude float64 // Height of wave
	Frequency float64 // Wave cycles per second
	BaseY     float64 // Center line Y position
	Time      float64 // Accumulated time
}

// Update moves entity in wave pattern.
func (wb *WavingBehavior) Update(entity *core.Entity, dt float64) {
	wb.Time += dt

	// Move forward
	entity.Transform.Position.X += wb.VelocityX * dt

	// Calculate wave Y position
	entity.Transform.Position.Y = wb.BaseY + math.Sin(wb.Time*wb.Frequency*2*math.Pi)*wb.Amplitude

	// Wrap around screen
	if entity.Transform.Position.X > 850 {
		entity.Transform.Position.X = -50
	}
}

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	log.Println("=== Moving Sprite Example ===")
	log.Println("Demonstrates various movement behaviors:")
	log.Println("  - Linear velocity (constant speed)")
	log.Println("  - Bouncing (collision with screen edges)")
	log.Println("  - Circular motion (orbiting)")
	log.Println("  - Wave pattern (sine wave movement)")
	log.Println()

	// Create engine
	engine, err := core.NewEngine("Moving Sprites - gogame", 800, 600, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 20, G: 30, B: 40, A: 255})

	// Position camera at center of screen to see entities properly
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}

	engine.SetScene(scene)

	// Load texture (reuse from assets example)
	assets := engine.Assets()
	texture, err := assets.LoadTexture("examples/assets/player.png")
	if err != nil {
		log.Printf("Note: Could not load texture (%v)", err)
		log.Println("The sprites will be invisible but movement still works!")
		log.Println("Run the assets example first to generate textures.")
	}

	// 1. Linear Velocity - Moves diagonally across screen
	if texture != nil {
		sprite1 := graphics.NewSprite(texture)
		sprite1.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255}) // Light blue
		entity1 := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 100, Y: 100},
				Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
			},
			Sprite: sprite1,
			Behavior: &VelocityBehavior{
				VelocityX: 50, // 50 pixels/second right
				VelocityY: 30, // 30 pixels/second down
			},
			Layer: 1,
		}
		scene.AddEntity(entity1)
		log.Println("✓ Created linear velocity sprite (blue, moving diagonal)")
	}

	// 2. Bouncing - Bounces off screen edges
	if texture != nil {
		sprite2 := graphics.NewSprite(texture)
		sprite2.SetColor(gamemath.Color{R: 255, G: 100, B: 100, A: 255}) // Red
		entity2 := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 400, Y: 300},
				Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
			},
			Sprite: sprite2,
			Behavior: &BouncingBehavior{
				VelocityX:    150,
				VelocityY:    -100,
				ScreenWidth:  800,
				ScreenHeight: 600,
				Margin:       32, // Half sprite size
			},
			Layer: 1,
		}
		scene.AddEntity(entity2)
		log.Println("✓ Created bouncing sprite (red, bounces off edges)")
	}

	// 3. Circular Motion - Orbits around center
	if texture != nil {
		sprite3 := graphics.NewSprite(texture)
		sprite3.SetColor(gamemath.Color{R: 100, G: 255, B: 100, A: 255}) // Green
		entity3 := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 400, Y: 150},
				Scale:    gamemath.Vector2{X: 1.2, Y: 1.2},
			},
			Sprite: sprite3,
			Behavior: &CircularMotionBehavior{
				CenterX:      400,
				CenterY:      300,
				Radius:       150,
				AngularSpeed: 1.0, // 1 radian per second
				CurrentAngle: 0,
			},
			Layer: 1,
		}
		scene.AddEntity(entity3)
		log.Println("✓ Created circular motion sprite (green, orbits center)")
	}

	// 4. Wave Pattern - Moves in sine wave
	if texture != nil {
		sprite4 := graphics.NewSprite(texture)
		sprite4.SetColor(gamemath.Color{R: 255, G: 215, B: 0, A: 255}) // Gold
		entity4 := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 50, Y: 450},
				Scale:    gamemath.Vector2{X: 1.0, Y: 1.0},
			},
			Sprite: sprite4,
			Behavior: &WavingBehavior{
				VelocityX: 80,  // Move right
				Amplitude: 50,  // Wave height
				Frequency: 1.0, // 1 cycle per second
				BaseY:     450,
				Time:      0,
			},
			Layer: 1,
		}
		scene.AddEntity(entity4)
		log.Println("✓ Created wave motion sprite (gold, sine wave)")
	}

	// Create multiple small bouncing sprites for visual interest
	if texture != nil {
		for i := 0; i < 5; i++ {
			sprite := graphics.NewSprite(texture)
			// Random colors
			colors := []gamemath.Color{
				{R: 255, G: 150, B: 150, A: 200},
				{R: 150, G: 255, B: 150, A: 200},
				{R: 150, G: 150, B: 255, A: 200},
				{R: 255, G: 255, B: 150, A: 200},
				{R: 255, G: 150, B: 255, A: 200},
			}
			sprite.SetColor(colors[i])
			sprite.Alpha = 0.8

			entity := &core.Entity{
				Active: true,
				Transform: gamemath.Transform{
					Position: gamemath.Vector2{
						X: 100 + float64(i)*150,
						Y: 200 + float64(i)*50,
					},
					Scale: gamemath.Vector2{X: 0.8, Y: 0.8},
				},
				Sprite: sprite,
				Behavior: &BouncingBehavior{
					VelocityX:    80 + float64(i)*20,
					VelocityY:    -60 + float64(i)*15,
					ScreenWidth:  800,
					ScreenHeight: 600,
					Margin:       25,
				},
				Layer: 0,
			}
			scene.AddEntity(entity)
		}
		log.Println("✓ Created 5 small bouncing sprites (semi-transparent)")
	}

	log.Println()
	log.Println("Movement patterns active:")
	log.Println("  - Linear: Moving diagonally")
	log.Println("  - Bouncing: Bouncing off edges")
	log.Println("  - Circular: Orbiting the center")
	log.Println("  - Wave: Sine wave pattern")
	log.Println("  - Multiple small sprites bouncing")
	log.Println()
	log.Println("Close window to exit.")
	log.Println()

	// Run game loop
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Example complete!")
}
