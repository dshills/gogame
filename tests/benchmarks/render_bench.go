package benchmarks

import (
	"runtime"
	"testing"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
)

// BenchmarkRendering100Sprites benchmarks rendering performance with 100 sprites.
// Target: <16ms per frame (60 FPS)
func BenchmarkRendering100Sprites(b *testing.B) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Create engine (headless would be ideal, but SDL requires window)
	engine, err := core.NewEngine("Benchmark", 800, 600, false)
	if err != nil {
		b.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 0, G: 0, B: 0, A: 255})

	// Position camera
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}

	engine.SetScene(scene)

	// Create 100 sprite entities
	for i := 0; i < 100; i++ {
		// Create simple colored sprite (no texture needed for benchmark)
		sprite := graphics.NewSprite(nil)
		sprite.SetColor(gamemath.Color{
			R: uint8(i % 255),
			G: uint8((i * 2) % 255),
			B: uint8((i * 3) % 255),
			A: 255,
		})

		entity := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{
					X: float64((i % 10) * 80),
					Y: float64((i / 10) * 60),
				},
				Scale: gamemath.Vector2{X: 1, Y: 1},
			},
			Sprite: sprite,
			Layer:  0,
		}
		scene.AddEntity(entity)
	}

	// Reset timer after setup
	b.ResetTimer()

	// Benchmark rendering
	for i := 0; i < b.N; i++ {
		// Simulate one frame render
		scene.Render(engine.Renderer(), camera)
	}
}

// BenchmarkRenderingWithTextures benchmarks rendering with actual textures.
func BenchmarkRenderingWithTextures(b *testing.B) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	engine, err := core.NewEngine("Benchmark Textures", 800, 600, false)
	if err != nil {
		b.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	// Load a test texture if available (skip if not found)
	texture, err := engine.Assets().LoadTexture("examples/assets/player.png")
	if err != nil {
		b.Skip("Test texture not available, skipping texture benchmark")
		return
	}

	scene := core.NewScene()
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}
	engine.SetScene(scene)

	// Create 50 textured sprites
	for i := 0; i < 50; i++ {
		sprite := graphics.NewSprite(texture)
		entity := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{
					X: float64((i % 10) * 80),
					Y: float64((i / 10) * 80),
				},
				Scale: gamemath.Vector2{X: 1, Y: 1},
			},
			Sprite: sprite,
			Layer:  0,
		}
		scene.AddEntity(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		scene.Render(engine.Renderer(), camera)
	}
}
