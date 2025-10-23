// Package main provides a player control example demonstrating keyboard input with WASD movement.
package main

import (
	"log"
	"runtime"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/input"
	gamemath "github.com/dshills/gogame/engine/math"
)

// PlayerController implements player movement with WASD keys.
type PlayerController struct {
	Speed float64 // Pixels per second
}

// Update moves the player based on input.
func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
	// Get input from the engine - this would normally be passed in,
	// but for this example we'll access it through a global reference
	// In a real game, you'd pass the engine/input manager to behaviors
	moveSpeed := pc.Speed * dt

	// Check action bindings
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

		// ESC to quit (just for demonstration)
		if inputMgr.KeyPressed(input.KeyEscape) {
			log.Println("ESC pressed - close the window to exit")
		}
	}
}

// Global input manager reference (simplified for example)
var inputMgr *input.InputManager

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	log.Println("=== Player Control Example ===")
	log.Println("Controls:")
	log.Println("  W/Up Arrow    - Move up")
	log.Println("  S/Down Arrow  - Move down")
	log.Println("  A/Left Arrow  - Move left")
	log.Println("  D/Right Arrow - Move right")
	log.Println("  ESC           - Print message")
	log.Println()

	// Create engine
	engine, err := core.NewEngine("Player Control - gogame", 800, 600, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Get input manager reference
	inputMgr = engine.Input()

	// Bind actions to keys
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
	inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
	inputMgr.BindAction(input.ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
	inputMgr.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 30, G: 30, B: 50, A: 255}) // Dark blue
	engine.SetScene(scene)

	// Create player entity with controller behavior
	player := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 400, Y: 300}, // Center of screen
			Rotation: 0,
			Scale:    gamemath.Vector2{X: 1, Y: 1},
		},
		Behavior: &PlayerController{Speed: 200}, // 200 pixels/second
		Layer:    1,
	}

	// For visual representation, we'd normally load a texture here
	// For this example, the player is invisible but movement still works
	// You can add: texture, _ := engine.Assets().LoadTexture("player.png")
	//              player.Sprite = graphics.NewSprite(texture)

	scene.AddEntity(player)

	log.Println("Player created at center of screen (400, 300)")
	log.Println("Note: Player is invisible in this example - add a sprite to see it!")
	log.Println()
	log.Println("Running... Use WASD or arrow keys to move!")
	log.Println()

	// Run game loop
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Game closed.")
}
