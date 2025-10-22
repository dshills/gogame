package main

import (
	"log"
	"runtime"

	"github.com/dshills/gogame/engine/core"
	gamemath "github.com/dshills/gogame/engine/math"
)

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	// Create engine with 800x600 window
	engine, err := core.NewEngine("Simple Example - gogame", 800, 600, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 52, G: 152, B: 219, A: 255}) // Sky blue
	engine.SetScene(scene)

	// Create a simple colored square (we'll just use a texture-less sprite for now)
	// Note: For a real sprite, you'd load a texture with engine.Assets().LoadTexture("path.png")
	// For this minimal example, we'll create entities without sprites to test the engine loop

	// Create a centered entity (will be invisible without a sprite, but tests the system)
	player := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 400, Y: 300},
			Rotation: 0,
			Scale:    gamemath.Vector2{X: 1, Y: 1},
		},
		Layer: 1,
	}
	scene.AddEntity(player)

	// Run game loop (blocks until window closed)
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}
}
