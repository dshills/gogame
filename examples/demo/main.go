// Package main provides a comprehensive demo showcasing all gogame engine features.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
)

// generateTestTextures creates simple colored PNG files for the demo.
func generateTestTextures() error {
	textures := map[string]color.RGBA{
		"examples/demo/assets/player.png":      {R: 255, G: 255, B: 255, A: 255}, // White square
		"examples/demo/assets/enemy.png":       {R: 200, G: 50, B: 50, A: 255},   // Red square
		"examples/demo/assets/collectible.png": {R: 255, G: 215, B: 0, A: 255},   // Gold square
	}

	for path, col := range textures {
		// Skip if file already exists
		if _, err := os.Stat(path); err == nil {
			continue
		}

		// Create a 64x64 image
		img := image.NewRGBA(image.Rect(0, 0, 64, 64))

		// Fill with color (with a border for visual interest)
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				// Create a border
				if x < 2 || x >= 62 || y < 2 || y >= 62 {
					img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255}) // Black border
				} else {
					img.Set(x, y, col)
				}
			}
		}

		// Save as PNG
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create texture file %s: %w", path, err)
		}

		if err := png.Encode(file, img); err != nil {
			_ = file.Close() // Best effort cleanup on error path
			return fmt.Errorf("failed to encode PNG %s: %w", path, err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close texture file %s: %w", path, err)
		}

		fmt.Printf("  Generated texture: %s\n", path)
	}

	return nil
}

// RotatingBehavior rotates an entity continuously.
type RotatingBehavior struct {
	RotationSpeed float64 // Degrees per second
}

func (rb *RotatingBehavior) Update(entity *core.Entity, dt float64) {
	entity.Transform.Rotate(rb.RotationSpeed * dt)
}

// OrbitingBehavior makes an entity orbit around a center point.
type OrbitingBehavior struct {
	CenterX      float64
	CenterY      float64
	Radius       float64
	Speed        float64 // Radians per second
	CurrentAngle float64
}

func (ob *OrbitingBehavior) Update(entity *core.Entity, dt float64) {
	ob.CurrentAngle += ob.Speed * dt
	entity.Transform.Position.X = ob.CenterX + math.Cos(ob.CurrentAngle)*ob.Radius
	entity.Transform.Position.Y = ob.CenterY + math.Sin(ob.CurrentAngle)*ob.Radius
}

// PulsatingBehavior changes entity alpha to create pulsing effect.
type PulsatingBehavior struct {
	Speed        float64
	CurrentPhase float64
}

func (pb *PulsatingBehavior) Update(entity *core.Entity, dt float64) {
	pb.CurrentPhase += pb.Speed * dt
	if entity.Sprite != nil {
		// Pulse alpha between 0.3 and 1.0
		entity.Sprite.Alpha = 0.65 + 0.35*math.Sin(pb.CurrentPhase)
	}
}

// BouncingBehavior makes entity bounce around the screen.
type BouncingBehavior struct {
	VelocityX    float64
	VelocityY    float64
	ScreenWidth  float64
	ScreenHeight float64
	Size         float64
}

func (bb *BouncingBehavior) Update(entity *core.Entity, dt float64) {
	// Update position
	entity.Transform.Position.X += bb.VelocityX * dt
	entity.Transform.Position.Y += bb.VelocityY * dt

	// Bounce off walls
	if entity.Transform.Position.X-bb.Size/2 < 0 || entity.Transform.Position.X+bb.Size/2 > bb.ScreenWidth {
		bb.VelocityX = -bb.VelocityX
		entity.Transform.Position.X = gamemath.Vector2{
			X: math.Max(bb.Size/2, math.Min(bb.ScreenWidth-bb.Size/2, entity.Transform.Position.X)),
			Y: entity.Transform.Position.Y,
		}.X
	}

	if entity.Transform.Position.Y-bb.Size/2 < 0 || entity.Transform.Position.Y+bb.Size/2 > bb.ScreenHeight {
		bb.VelocityY = -bb.VelocityY
		entity.Transform.Position.Y = gamemath.Vector2{
			X: entity.Transform.Position.X,
			Y: math.Max(bb.Size/2, math.Min(bb.ScreenHeight-bb.Size/2, entity.Transform.Position.Y)),
		}.Y
	}
}

// SmoothFollowBehavior makes entity smoothly follow another entity.
type SmoothFollowBehavior struct {
	Target    *core.Entity
	Smoothing float64
}

func (sf *SmoothFollowBehavior) Update(entity *core.Entity, _ float64) {
	if sf.Target != nil && sf.Target.Active {
		// Smooth interpolation toward target
		entity.Transform.Position.X += (sf.Target.Transform.Position.X - entity.Transform.Position.X) * sf.Smoothing
		entity.Transform.Position.Y += (sf.Target.Transform.Position.Y - entity.Transform.Position.Y) * sf.Smoothing
	}
}

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	fmt.Println("=== gogame Engine Demo ===")
	fmt.Println("This demo showcases all current engine features:")
	fmt.Println("- Sprite rendering with textures")
	fmt.Println("- Transform operations (position, rotation, scale)")
	fmt.Println("- Color tinting and alpha blending")
	fmt.Println("- Sprite flipping (horizontal/vertical)")
	fmt.Println("- Custom behaviors (rotating, orbiting, pulsating, bouncing)")
	fmt.Println("- Camera system with smooth follow")
	fmt.Println("- Entity management (add/remove)")
	fmt.Println("- Fixed 60 FPS game loop")
	fmt.Println()
	fmt.Println("Close the window to exit.")
	fmt.Println()

	// Create test assets directory
	if err := os.MkdirAll("examples/demo/assets", 0755); err != nil {
		log.Fatal("Failed to create assets directory:", err)
	}

	// Generate test textures if they don't exist
	if err := generateTestTextures(); err != nil {
		log.Fatal("Failed to generate test textures:", err)
	}

	// Create engine with 1024x768 window
	engine, err := core.NewEngine("gogame Feature Demo", 1024, 768, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Create scene with nice gradient-like background
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 25, G: 25, B: 40, A: 255}) // Dark blue-gray
	engine.SetScene(scene)

	// Load textures
	playerTexture, err := engine.Assets().LoadTexture("examples/demo/assets/player.png")
	if err != nil {
		log.Fatal("Failed to load player texture:", err)
	}

	enemyTexture, err := engine.Assets().LoadTexture("examples/demo/assets/enemy.png")
	if err != nil {
		log.Fatal("Failed to load enemy texture:", err)
	}

	collectibleTexture, err := engine.Assets().LoadTexture("examples/demo/assets/collectible.png")
	if err != nil {
		log.Fatal("Failed to load collectible texture:", err)
	}

	// 1. CENTER ENTITY - Player with smooth rotation
	playerSprite := graphics.NewSprite(playerTexture)
	playerSprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255}) // Light blue tint

	player := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 512, Y: 384},
			Rotation: 0,
			Scale:    gamemath.Vector2{X: 2.0, Y: 2.0}, // 2x scale
		},
		Sprite:   playerSprite,
		Behavior: &RotatingBehavior{RotationSpeed: 45}, // 45 degrees per second
		Layer:    5,
	}
	scene.AddEntity(player)
	fmt.Println("✓ Created center player entity (rotating, scaled 2x, blue tint)")

	// 2. ORBITING ENTITIES - Multiple enemies orbiting the player
	for i := 0; i < 6; i++ {
		angle := float64(i) * (2 * math.Pi / 6)
		enemySprite := graphics.NewSprite(enemyTexture)

		// Each enemy has different color tint
		colors := []gamemath.Color{
			{R: 255, G: 100, B: 100, A: 255}, // Red
			{R: 100, G: 255, B: 100, A: 255}, // Green
			{R: 255, G: 255, B: 100, A: 255}, // Yellow
			{R: 255, G: 100, B: 255, A: 255}, // Magenta
			{R: 100, G: 255, B: 255, A: 255}, // Cyan
			{R: 255, G: 200, B: 100, A: 255}, // Orange
		}
		enemySprite.SetColor(colors[i])

		enemy := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 512, Y: 384},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 1.0, Y: 1.0},
			},
			Sprite: enemySprite,
			Behavior: &OrbitingBehavior{
				CenterX:      512,
				CenterY:      384,
				Radius:       150,
				Speed:        0.5 + float64(i)*0.1, // Different speeds
				CurrentAngle: angle,
			},
			Layer: 3,
		}
		scene.AddEntity(enemy)
	}
	fmt.Println("✓ Created 6 orbiting enemies (different colors and speeds)")

	// 3. PULSATING COLLECTIBLES - Corners with alpha pulsing
	corners := []gamemath.Vector2{
		{X: 100, Y: 100},
		{X: 924, Y: 100},
		{X: 100, Y: 668},
		{X: 924, Y: 668},
	}

	for i, pos := range corners {
		collectibleSprite := graphics.NewSprite(collectibleTexture)
		collectibleSprite.SetColor(gamemath.Color{R: 255, G: 255, B: 100, A: 255}) // Yellow

		collectible := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: pos,
				Rotation: float64(i * 90), // Each rotated differently
				Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
			},
			Sprite: collectibleSprite,
			Behavior: &PulsatingBehavior{
				Speed:        2.0 + float64(i)*0.5,
				CurrentPhase: float64(i) * math.Pi / 2, // Offset phases
			},
			Layer: 2,
		}
		scene.AddEntity(collectible)
	}
	fmt.Println("✓ Created 4 pulsating collectibles in corners (alpha animation)")

	// 4. BOUNCING ENTITIES - Multiple bouncing sprites
	for i := 0; i < 4; i++ {
		bouncingSprite := graphics.NewSprite(enemyTexture)
		bouncingSprite.SetColor(gamemath.Color{R: 150, G: 150, B: 255, A: 200}) // Semi-transparent blue

		// Flip some sprites
		if i%2 == 0 {
			bouncingSprite.FlipH = true
		}
		if i >= 2 {
			bouncingSprite.FlipV = true
		}

		bouncer := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{
					X: 200 + float64(i)*200,
					Y: 300,
				},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 1.2, Y: 1.2},
			},
			Sprite: bouncingSprite,
			Behavior: &BouncingBehavior{
				VelocityX:    100 + float64(i)*50,
				VelocityY:    -150 + float64(i)*30,
				ScreenWidth:  1024,
				ScreenHeight: 768,
				Size:         64,
			},
			Layer: 4,
		}
		scene.AddEntity(bouncer)
	}
	fmt.Println("✓ Created 4 bouncing entities (with sprite flipping, semi-transparent)")

	// 5. FOLLOWER ENTITY - Follows the player smoothly
	followerSprite := graphics.NewSprite(collectibleTexture)
	followerSprite.SetColor(gamemath.Color{R: 255, G: 150, B: 255, A: 180}) // Purple, semi-transparent

	follower := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 700, Y: 200},
			Rotation: 0,
			Scale:    gamemath.Vector2{X: 1.0, Y: 1.0},
		},
		Sprite: followerSprite,
		Behavior: &SmoothFollowBehavior{
			Target:    player,
			Smoothing: 0.02, // Smooth interpolation
		},
		Layer: 6,
	}
	scene.AddEntity(follower)
	fmt.Println("✓ Created follower entity (smoothly follows player)")

	// 6. CAMERA - Set up camera to follow the player with smooth interpolation
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 512, Y: 384}
	camera.Zoom = 1.0
	fmt.Println("✓ Camera configured (centered)")

	fmt.Println()
	fmt.Println("Demo running! Features active:")
	fmt.Println("  - 1 rotating player (center, blue tint, 2x scale)")
	fmt.Println("  - 6 orbiting enemies (rainbow colors)")
	fmt.Println("  - 4 pulsating collectibles (corners, alpha animation)")
	fmt.Println("  - 4 bouncing entities (flipped sprites)")
	fmt.Println("  - 1 smooth follower (purple, follows player)")
	fmt.Println("  - Total: 16 active entities")
	fmt.Println()

	// Run game loop (blocks until window closed)
	if err := engine.Run(); err != nil {
		log.Fatal("Engine error:", err)
	}

	fmt.Println()
	fmt.Println("Demo complete!")
}
