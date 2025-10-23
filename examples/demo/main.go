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
	"github.com/dshills/gogame/engine/input"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// generateTestTextures creates simple colored PNG files for the demo.
func generateTestTextures() error {
	textures := map[string]color.RGBA{
		"examples/demo/assets/player.png":      {R: 100, G: 200, B: 255, A: 255}, // Light blue
		"examples/demo/assets/enemy.png":       {R: 200, G: 50, B: 50, A: 255},   // Red
		"examples/demo/assets/collectible.png": {R: 255, G: 215, B: 0, A: 255},   // Gold
		"examples/demo/assets/wall.png":        {R: 100, G: 100, B: 100, A: 255}, // Gray
	}

	for path, col := range textures {
		// Skip if file already exists
		if _, err := os.Stat(path); err == nil {
			continue
		}

		// Create a 32x32 image
		img := image.NewRGBA(image.Rect(0, 0, 32, 32))

		// Fill with color (with a border for visual interest)
		for y := 0; y < 32; y++ {
			for x := 0; x < 32; x++ {
				// Create a border
				if x < 2 || x >= 30 || y < 2 || y >= 30 {
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
			_ = file.Close()
			return fmt.Errorf("failed to encode PNG %s: %w", path, err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close texture file %s: %w", path, err)
		}

		fmt.Printf("  Generated texture: %s\n", path)
	}

	return nil
}

// PlayerController demonstrates input handling with WASD movement.
type PlayerController struct {
	Speed         float64
	InputMgr      *input.InputManager
	CollectCount  int
	LastCollision uint64
}

func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
	// Movement with input
	moveSpeed := pc.Speed * dt
	if pc.InputMgr.ActionHeld(input.ActionMoveUp) {
		entity.Transform.Position.Y -= moveSpeed
	}
	if pc.InputMgr.ActionHeld(input.ActionMoveDown) {
		entity.Transform.Position.Y += moveSpeed
	}
	if pc.InputMgr.ActionHeld(input.ActionMoveLeft) {
		entity.Transform.Position.X -= moveSpeed
	}
	if pc.InputMgr.ActionHeld(input.ActionMoveRight) {
		entity.Transform.Position.X += moveSpeed
	}

	// Rotate player sprite to face movement direction
	if pc.InputMgr.ActionHeld(input.ActionMoveUp) || pc.InputMgr.ActionHeld(input.ActionMoveDown) ||
		pc.InputMgr.ActionHeld(input.ActionMoveLeft) || pc.InputMgr.ActionHeld(input.ActionMoveRight) {
		entity.Transform.Rotation += 180 * dt // Spin while moving
	}
}

// EnemyPatrol makes enemies patrol back and forth.
type EnemyPatrol struct {
	Speed     float64
	MinX      float64
	MaxX      float64
	Direction float64
}

func (ep *EnemyPatrol) Update(entity *core.Entity, dt float64) {
	entity.Transform.Position.X += ep.Direction * ep.Speed * dt

	// Reverse direction at boundaries
	if entity.Transform.Position.X <= ep.MinX {
		ep.Direction = 1
		entity.Transform.Position.X = ep.MinX
	} else if entity.Transform.Position.X >= ep.MaxX {
		ep.Direction = -1
		entity.Transform.Position.X = ep.MaxX
	}

	// Rotate enemy
	entity.Transform.Rotation += 90 * dt
}

// CollectibleBehavior makes collectibles pulse and rotate.
type CollectibleBehavior struct {
	Speed        float64
	CurrentPhase float64
}

func (cb *CollectibleBehavior) Update(entity *core.Entity, dt float64) {
	cb.CurrentPhase += cb.Speed * dt

	// Pulse scale
	scaleFactor := 1.0 + 0.2*math.Sin(cb.CurrentPhase)
	entity.Transform.Scale = gamemath.Vector2{X: scaleFactor, Y: scaleFactor}

	// Rotate
	entity.Transform.Rotation += 120 * dt

	// Pulse alpha
	if entity.Sprite != nil {
		entity.Sprite.Alpha = 0.7 + 0.3*math.Sin(cb.CurrentPhase*2)
	}
}

// Global reference for collision detection
var (
	playerEntity      *core.Entity
	collectibles      []*core.Entity
	enemies           []*core.Entity
	frameCount        int
	lastCollisionLog  int
	collectiblesFound int
)

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║          gogame Engine - Comprehensive Feature Demo      ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("This demo showcases ALL engine features:")
	fmt.Println()
	fmt.Println("  ✓ Core Game Loop - Fixed 60 FPS with delta time")
	fmt.Println("  ✓ Entity/Scene Management - Dynamic add/remove")
	fmt.Println("  ✓ Input Handling - WASD movement with action mapping")
	fmt.Println("  ✓ Asset Loading - Texture caching and reference counting")
	fmt.Println("  ✓ Collision Detection - AABB with layer masks")
	fmt.Println("  ✓ Sprite Rendering - Textures, colors, alpha, transforms")
	fmt.Println("  ✓ Camera System - World space rendering")
	fmt.Println()
	fmt.Println("Controls:")
	fmt.Println("  WASD / Arrow Keys - Move player (blue square)")
	fmt.Println("  ESC - Print debug info")
	fmt.Println()
	fmt.Println("Objective:")
	fmt.Println("  Collect all 4 golden collectibles!")
	fmt.Println("  Avoid the red patrolling enemies!")
	fmt.Println()

	// Create test assets directory
	if err := os.MkdirAll("examples/demo/assets", 0755); err != nil {
		log.Fatal("Failed to create assets directory:", err)
	}

	// Generate test textures
	if err := generateTestTextures(); err != nil {
		log.Fatal("Failed to generate test textures:", err)
	}

	// Create engine
	engine, err := core.NewEngine("gogame Feature Demo - All Systems", 800, 600, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Setup input
	inputMgr := engine.Input()
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
	inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
	inputMgr.BindAction(input.ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
	inputMgr.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 30, G: 30, B: 50, A: 255})
	engine.SetScene(scene)

	// Load textures (demonstrates asset management with reference counting)
	assets := engine.Assets()
	playerTexture, _ := assets.LoadTexture("examples/demo/assets/player.png")
	enemyTexture, _ := assets.LoadTexture("examples/demo/assets/enemy.png")
	collectibleTexture, _ := assets.LoadTexture("examples/demo/assets/collectible.png")
	wallTexture, _ := assets.LoadTexture("examples/demo/assets/wall.png")

	// Create player entity (Layer 0 - Player)
	playerController := &PlayerController{
		Speed:    250,
		InputMgr: inputMgr,
	}
	playerEntity = &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 400, Y: 300},
			Rotation: 0,
			Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
		},
		Sprite:   graphics.NewSprite(playerTexture),
		Collider: physics.NewCollider(32, 32),
		Behavior: playerController,
		Layer:    2,
	}
	playerEntity.Collider.CollisionLayer = 0                  // Player layer
	playerEntity.Collider.CollisionMask = (1 << 1) | (1 << 2) // Collides with enemies and collectibles
	scene.AddEntity(playerEntity)
	fmt.Println("✓ Player created (blue, WASD controls, layer 0)")

	// Create walls (Layer 3 - Walls)
	wallPositions := []struct {
		x, y, w, h float64
	}{
		{400, 50, 700, 20},  // Top
		{400, 550, 700, 20}, // Bottom
		{50, 300, 20, 500},  // Left
		{750, 300, 20, 500}, // Right
		{400, 200, 200, 20}, // Obstacle 1
		{400, 400, 200, 20}, // Obstacle 2
	}

	for _, wall := range wallPositions {
		wallEntity := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: wall.x, Y: wall.y},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: wall.w / 32, Y: wall.h / 32},
			},
			Sprite:   graphics.NewSprite(wallTexture),
			Collider: physics.NewCollider(32, 32),
			Layer:    0,
		}
		wallEntity.Collider.CollisionLayer = 3                  // Wall layer
		wallEntity.Collider.CollisionMask = (1 << 0) | (1 << 1) // Collides with player and enemies
		scene.AddEntity(wallEntity)
	}
	fmt.Println("✓ Walls created (gray, collision layer 3)")

	// Create patrolling enemies (Layer 1 - Enemies)
	enemyConfigs := []struct {
		y          float64
		minX, maxX float64
		dir        float64
	}{
		{150, 150, 650, 1},
		{300, 200, 600, -1},
		{450, 150, 650, 1},
	}

	for i, cfg := range enemyConfigs {
		enemySprite := graphics.NewSprite(enemyTexture)
		enemy := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: cfg.minX, Y: cfg.y},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 1.2, Y: 1.2},
			},
			Sprite:   enemySprite,
			Collider: physics.NewCollider(32, 32),
			Behavior: &EnemyPatrol{
				Speed:     100 + float64(i)*20,
				MinX:      cfg.minX,
				MaxX:      cfg.maxX,
				Direction: cfg.dir,
			},
			Layer: 1,
		}
		enemy.Collider.CollisionLayer = 1     // Enemy layer
		enemy.Collider.CollisionMask = 1 << 0 // Collides with player
		scene.AddEntity(enemy)
		enemies = append(enemies, enemy)
	}
	fmt.Println("✓ Enemies created (red, patrolling, layer 1)")

	// Create collectibles (Layer 2 - Collectibles)
	collectiblePositions := []gamemath.Vector2{
		{X: 150, Y: 100},
		{X: 650, Y: 100},
		{X: 150, Y: 500},
		{X: 650, Y: 500},
	}

	for i, pos := range collectiblePositions {
		collectibleSprite := graphics.NewSprite(collectibleTexture)
		collectible := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: pos,
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 1.0, Y: 1.0},
			},
			Sprite:   collectibleSprite,
			Collider: physics.NewCollider(32, 32),
			Behavior: &CollectibleBehavior{
				Speed:        2.0 + float64(i)*0.3,
				CurrentPhase: float64(i) * math.Pi / 2,
			},
			Layer: 1,
		}
		collectible.Collider.CollisionLayer = 2     // Collectible layer
		collectible.Collider.CollisionMask = 1 << 0 // Collides with player
		collectible.Collider.IsTrigger = true       // Don't block movement
		scene.AddEntity(collectible)
		collectibles = append(collectibles, collectible)
	}
	fmt.Println("✓ Collectibles created (gold, pulsating, layer 2)")

	// Setup camera
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}
	camera.Zoom = 1.0

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Demo running! Move with WASD and collect golden items!")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Wrap player behavior to detect collisions
	originalBehavior := playerEntity.Behavior
	playerEntity.Behavior = BehaviorFunc(func(entity *core.Entity, dt float64) {
		// Update original behavior
		if originalBehavior != nil {
			originalBehavior.Update(entity, dt)
		}

		// Check collisions every frame
		frameCount++

		// Check ESC for debug info
		if inputMgr.KeyPressed(input.KeyEscape) {
			fmt.Println()
			fmt.Println("═══ DEBUG INFO ═══")
			fmt.Printf("Player Position: (%.0f, %.0f)\n", entity.Transform.Position.X, entity.Transform.Position.Y)
			fmt.Printf("Collectibles Found: %d / %d\n", collectiblesFound, len(collectiblePositions))
			fmt.Printf("Active Entities: %d\n", len(scene.GetAllEntities()))
			fmt.Println("═══════════════════")
			fmt.Println()
		}

		// Check collisions with collectibles
		for i, collectible := range collectibles {
			if !collectible.Active {
				continue
			}

			if entity.Collider.Intersects(collectible.GetCollider(), entity.Transform, collectible.Transform) {
				// Collect it!
				collectible.Active = false
				scene.RemoveEntity(collectible.ID)
				collectiblesFound++

				fmt.Printf("✨ Collectible %d/%d found! ", collectiblesFound, len(collectiblePositions))
				if collectiblesFound == len(collectiblePositions) {
					fmt.Println("🎉 ALL COLLECTIBLES FOUND! You win!")
				} else {
					fmt.Printf("%d remaining.\n", len(collectiblePositions)-collectiblesFound)
				}

				// Remove from tracking
				collectibles[i] = collectibles[len(collectibles)-1]
				collectibles = collectibles[:len(collectibles)-1]
			}
		}

		// Check collisions with enemies (every 30 frames to avoid spam)
		if frameCount%30 == 0 {
			for _, enemy := range enemies {
				if !enemy.Active {
					continue
				}

				if entity.Collider.Intersects(enemy.GetCollider(), entity.Transform, enemy.Transform) {
					if lastCollisionLog != frameCount {
						fmt.Printf("💥 COLLISION with enemy at (%.0f, %.0f)! Watch out!\n",
							enemy.Transform.Position.X, enemy.Transform.Position.Y)
						lastCollisionLog = frameCount
					}
				}
			}
		}
	})

	// Run game loop
	if err := engine.Run(); err != nil {
		log.Fatal("Engine error:", err)
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Demo complete!")
	fmt.Printf("Final Score: %d / %d collectibles found\n", collectiblesFound, len(collectiblePositions))
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
}

// BehaviorFunc adapter to wrap functions as Behaviors.
type BehaviorFunc func(entity *core.Entity, dt float64)

func (f BehaviorFunc) Update(entity *core.Entity, dt float64) {
	f(entity, dt)
}
