// Package main demonstrates collision callbacks with OnCollisionEnter, OnCollisionStay, and OnCollisionExit.
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

// PlayerController with input-driven movement.
type PlayerController struct {
	Speed    float64
	InputMgr *input.InputManager
}

func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
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
}

// Global counters for collision events
var (
	enterCount = 0
	stayCount  = 0
	exitCount  = 0
)

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	log.Println("╔═══════════════════════════════════════════════════════════╗")
	log.Println("║          Collision Callbacks Demonstration               ║")
	log.Println("╚═══════════════════════════════════════════════════════════╝")
	log.Println()
	log.Println("This example demonstrates collision event callbacks:")
	log.Println("  • OnCollisionEnter - Called when collision starts")
	log.Println("  • OnCollisionStay  - Called while collision continues")
	log.Println("  • OnCollisionExit  - Called when collision ends")
	log.Println()
	log.Println("Controls:")
	log.Println("  WASD / Arrow Keys - Move player (blue)")
	log.Println("  Move into the red target to trigger callbacks")
	log.Println()

	// Create engine
	engine, err := core.NewEngine("Collision Callbacks - gogame", 800, 600, false)
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

	// Load textures
	assets := engine.Assets()
	playerTexture, _ := assets.LoadTexture("examples/assets/player.png")
	enemyTexture, _ := assets.LoadTexture("examples/assets/enemy.png")

	// Create player entity with collision callbacks
	playerSprite := graphics.NewSprite(playerTexture)
	playerSprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255})

	player := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 200, Y: 300},
			Scale:    gamemath.Vector2{X: 2, Y: 2},
		},
		Sprite:   playerSprite,
		Collider: physics.NewCollider(32, 32),
		Behavior: &PlayerController{
			Speed:    200,
			InputMgr: inputMgr,
		},
		Layer: 1,
	}

	// Setup collision callbacks on player
	player.OnCollisionEnter = func(self, other *core.Entity) {
		enterCount++
		log.Printf("🟢 ENTER: Player collided with entity %d (Total enters: %d)", other.ID, enterCount)

		// Change player color when entering collision
		if self.Sprite != nil {
			self.Sprite.SetColor(gamemath.Color{R: 255, G: 100, B: 100, A: 255}) // Red
		}
	}

	player.OnCollisionStay = func(self, other *core.Entity) {
		stayCount++
		// Log every 30th frame to avoid spam
		if stayCount%30 == 0 {
			log.Printf("🟡 STAY: Player still colliding with entity %d (Total stays: %d)", other.ID, stayCount)
		}
	}

	player.OnCollisionExit = func(self, other *core.Entity) {
		exitCount++
		log.Printf("🔴 EXIT: Player stopped colliding with entity %d (Total exits: %d)", other.ID, exitCount)

		// Restore player color when exiting collision
		if self.Sprite != nil {
			self.Sprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255}) // Blue
		}
	}

	player.Collider.CollisionLayer = 0
	player.Collider.CollisionMask = 1 << 1 // Collides with layer 1
	scene.AddEntity(player)

	// Create target entity with collision callbacks
	targetSprite := graphics.NewSprite(enemyTexture)
	targetSprite.SetColor(gamemath.Color{R: 200, G: 50, B: 50, A: 255})

	target := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 600, Y: 300},
			Scale:    gamemath.Vector2{X: 2.5, Y: 2.5},
		},
		Sprite:   targetSprite,
		Collider: physics.NewCollider(32, 32),
		Layer:    1,
	}

	// Setup collision callbacks on target
	target.OnCollisionEnter = func(self, other *core.Entity) {
		log.Printf("🎯 TARGET: Detected player entering collision zone")
		// Make target pulse when hit
		if self.Sprite != nil {
			self.Sprite.Alpha = 0.5
		}
	}

	target.OnCollisionStay = func(self, other *core.Entity) {
		// Keep target semi-transparent while colliding
		if self.Sprite != nil {
			self.Sprite.Alpha = 0.5
		}
	}

	target.OnCollisionExit = func(self, other *core.Entity) {
		log.Printf("🎯 TARGET: Player left collision zone")
		// Restore target alpha
		if self.Sprite != nil {
			self.Sprite.Alpha = 1.0
		}
	}

	target.Collider.CollisionLayer = 1
	target.Collider.CollisionMask = 1 << 0 // Collides with player
	scene.AddEntity(target)

	// Create walls for additional collision testing
	wallPositions := []gamemath.Vector2{
		{X: 400, Y: 100},
		{X: 400, Y: 500},
	}

	for i, pos := range wallPositions {
		wallSprite := graphics.NewSprite(enemyTexture)
		wallSprite.SetColor(gamemath.Color{R: 150, G: 150, B: 150, A: 255})

		wall := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: pos,
				Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
			},
			Sprite:   wallSprite,
			Collider: physics.NewCollider(32, 32),
			Layer:    1,
		}

		// Walls also have callbacks
		wallID := i + 1
		wall.OnCollisionEnter = func(self, other *core.Entity) {
			log.Printf("🧱 WALL %d: Collision started", wallID)
		}

		wall.OnCollisionExit = func(self, other *core.Entity) {
			log.Printf("🧱 WALL %d: Collision ended", wallID)
		}

		wall.Collider.CollisionLayer = 1
		wall.Collider.CollisionMask = 1 << 0
		scene.AddEntity(wall)
	}

	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("Scene created:")
	log.Println("  • 1 Player (blue, WASD control)")
	log.Println("  • 1 Target (red, large)")
	log.Println("  • 2 Walls (gray)")
	log.Println()
	log.Println("Move the player around to test collision callbacks!")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println()

	// Run game loop
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println()
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("Final Statistics:")
	log.Printf("  OnCollisionEnter called: %d times\n", enterCount)
	log.Printf("  OnCollisionStay called:  %d times\n", stayCount)
	log.Printf("  OnCollisionExit called:  %d times\n", exitCount)
	log.Println("═══════════════════════════════════════════════════════════")
}
