// Package main provides a collision detection example demonstrating AABB collision with layer masks.
package main

import (
	"log"
	"runtime"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	"github.com/dshills/gogame/engine/input"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// PlayerController with collision-aware movement.
type PlayerController struct {
	Speed float64
}

func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
	moveSpeed := pc.Speed * dt

	if inputMgr != nil {
		if inputMgr.ActionHeld(input.ActionMoveUp) {
			entity.Transform.Position.Y -= moveSpeed
		}
		if inputMgr.ActionHeld(input.ActionMoveDown) {
			entity.Transform.Position.Y += moveSpeed
		}
		if inputMgr.ActionHeld(input.ActionMoveLeft) {
			entity.Transform.Position.X -= moveSpeed
		}
		if inputMgr.ActionHeld(input.ActionMoveRight) {
			entity.Transform.Position.X += moveSpeed
		}
	}
}

// Global input manager reference (simplified for example)
var inputMgr *input.InputManager

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	log.Println("=== Collision Detection Example ===")
	log.Println("Controls:")
	log.Println("  WASD/Arrows - Move player (green)")
	log.Println("  Watch console for collision messages")
	log.Println()

	// Create engine
	engine, err := core.NewEngine("Collision Detection - gogame", 800, 600, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Get input manager
	inputMgr = engine.Input()
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
	inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
	inputMgr.BindAction(input.ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
	inputMgr.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 30, G: 30, B: 50, A: 255})
	engine.SetScene(scene)

	// Load textures
	assets := engine.Assets()
	playerTexture, _ := assets.LoadTexture("examples/assets/player.png")
	enemyTexture, _ := assets.LoadTexture("examples/assets/enemy.png")

	// Create player entity with collider (Layer 0 - Player)
	player := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 200, Y: 300},
			Rotation: 0,
			Scale:    gamemath.Vector2{X: 2, Y: 2},
		},
		Collider: physics.NewCollider(32, 32),
		Behavior: &PlayerController{Speed: 200},
		Layer:    1,
	}
	if playerTexture != nil {
		player.Sprite = graphics.NewSprite(playerTexture)
	}
	player.Collider.CollisionLayer = 0          // Player layer
	player.Collider.CollisionMask = 1<<1 | 1<<2 // Collides with enemies and walls
	scene.AddEntity(player)

	// Create enemy entities (Layer 1 - Enemies)
	enemyPositions := []gamemath.Vector2{
		{X: 400, Y: 200},
		{X: 500, Y: 300},
		{X: 400, Y: 400},
	}

	for i, pos := range enemyPositions {
		enemy := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: pos,
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 2, Y: 2},
			},
			Collider: physics.NewCollider(32, 32),
			Layer:    1,
		}
		if enemyTexture != nil {
			enemy.Sprite = graphics.NewSprite(enemyTexture)
		}
		enemy.Collider.CollisionLayer = 1     // Enemy layer
		enemy.Collider.CollisionMask = 1 << 0 // Collides with player
		scene.AddEntity(enemy)

		if i == 0 {
			log.Printf("Enemy created at (%.0f, %.0f)", pos.X, pos.Y)
		}
	}

	// Create walls (Layer 2 - Walls) - just colliders, no sprites
	walls := []gamemath.Rectangle{
		{X: 50, Y: 50, Width: 700, Height: 20},  // Top wall
		{X: 50, Y: 530, Width: 700, Height: 20}, // Bottom wall
		{X: 50, Y: 50, Width: 20, Height: 500},  // Left wall
		{X: 730, Y: 50, Width: 20, Height: 500}, // Right wall
	}

	for _, wallRect := range walls {
		wall := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{
					X: wallRect.X + wallRect.Width/2,
					Y: wallRect.Y + wallRect.Height/2,
				},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 1, Y: 1},
			},
			Collider: physics.NewCollider(wallRect.Width, wallRect.Height),
			Layer:    0,
		}
		wall.Collider.CollisionLayer = 2     // Wall layer
		wall.Collider.CollisionMask = 1 << 0 // Collides with player only
		scene.AddEntity(wall)
	}

	log.Println("Scene created:")
	log.Println("  - 1 player (green, layer 0)")
	log.Println("  - 3 enemies (red, layer 1)")
	log.Println("  - 4 walls (invisible, layer 2)")
	log.Println()
	log.Println("Collision detection running every frame!")
	log.Println("Move the player into enemies or walls to test")
	log.Println()

	// Track collision state for logging
	lastCollisionCheck := 0
	frameCount := 0

	// Custom behavior to check collisions
	playerBehavior := &PlayerController{Speed: 200}
	player.Behavior = playerBehavior

	// Before running, add a wrapper behavior that checks collisions
	originalBehavior := player.Behavior
	player.Behavior = BehaviorFunc(func(entity *core.Entity, dt float64) {
		// Call original behavior
		if originalBehavior != nil {
			originalBehavior.Update(entity, dt)
		}

		// Every 30 frames, check if player is colliding with anything
		frameCount++
		if frameCount%30 == 0 {
			// Get all entities and check collisions manually for demo
			scene := engine.GetScene()
			if scene != nil {
				allEntities := scene.GetAllEntities()
				colliding := false
				for _, other := range allEntities {
					if other.GetID() == entity.GetID() {
						continue
					}
					if entity.Collider != nil && other.GetCollider() != nil {
						if entity.Collider.Intersects(other.GetCollider(), entity.Transform, other.GetTransform()) {
							if lastCollisionCheck != frameCount {
								log.Printf("COLLISION: Player colliding with entity %d at (%.0f, %.0f)",
									other.GetID(), other.GetTransform().Position.X, other.GetTransform().Position.Y)
								colliding = true
							}
						}
					}
				}
				if colliding {
					lastCollisionCheck = frameCount
				}
			}
		}
	})

	log.Println("Running... Use WASD or arrow keys to move!")
	log.Println()

	// Run game loop
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Game closed.")
}

// BehaviorFunc adapter type.
type BehaviorFunc func(entity *core.Entity, dt float64)

// Update implements Behavior interface.
func (f BehaviorFunc) Update(entity *core.Entity, dt float64) {
	f(entity, dt)
}
