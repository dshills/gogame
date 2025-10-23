// Package main provides an asset loading example demonstrating texture loading with reference counting.
package main

import (
	"log"
	"runtime"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
)

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	log.Println("=== Asset Loading Example ===")
	log.Println("Demonstrates texture loading with reference counting and caching")
	log.Println()

	// Create engine
	engine, err := core.NewEngine("Asset Loading - gogame", 800, 600, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Get asset manager
	assets := engine.Assets()

	log.Println("Loading player.png...")
	playerTexture, err := assets.LoadTexture("examples/assets/player.png")
	if err != nil {
		log.Printf("WARNING: Could not load player.png: %v", err)
		log.Println("Note: This example requires actual PNG files to display sprites.")
		log.Println("Create examples/assets/player.png and enemy.png (32x32 PNGs) to see them render.")
	} else {
		log.Printf("✓ Loaded player.png (%dx%d)", playerTexture.Width, playerTexture.Height)
	}

	log.Println("Loading enemy.png...")
	enemyTexture, err := assets.LoadTexture("examples/assets/enemy.png")
	if err != nil {
		log.Printf("WARNING: Could not load enemy.png: %v", err)
	} else {
		log.Printf("✓ Loaded enemy.png (%dx%d)", enemyTexture.Width, enemyTexture.Height)
	}

	// Demonstrate reference counting - load player texture again
	if playerTexture != nil {
		log.Println("Loading player.png again to test caching...")
		playerTexture2, _ := assets.LoadTexture("examples/assets/player.png")
		if playerTexture == playerTexture2 {
			log.Println("✓ Reference counting working - same texture instance returned")
		}
	}

	// Create scene
	scene := core.NewScene()
	scene.SetBackgroundColor(gamemath.Color{R: 50, G: 50, B: 70, A: 255}) // Dark blue-gray

	// Position camera at center of screen to see entities properly
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}

	engine.SetScene(scene)

	// Create player entity with sprite (if texture loaded successfully)
	if playerTexture != nil {
		playerSprite := graphics.NewSprite(playerTexture)
		player := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 300, Y: 300},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 2, Y: 2}, // 2x scale for visibility
			},
			Sprite: playerSprite,
			Layer:  1,
		}
		scene.AddEntity(player)
		log.Println("✓ Created player entity at (300, 300)")
	}

	// Create enemy entity with sprite (if texture loaded successfully)
	if enemyTexture != nil {
		enemySprite := graphics.NewSprite(enemyTexture)
		enemy := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{X: 500, Y: 300},
				Rotation: 0,
				Scale:    gamemath.Vector2{X: 2, Y: 2},
			},
			Sprite: enemySprite,
			Layer:  1,
		}
		scene.AddEntity(enemy)
		log.Println("✓ Created enemy entity at (500, 300)")
	}

	// Create multiple player sprites to demonstrate sharing
	if playerTexture != nil {
		for i := 0; i < 3; i++ {
			sprite := graphics.NewSprite(playerTexture)
			entity := &core.Entity{
				Active: true,
				Transform: gamemath.Transform{
					Position: gamemath.Vector2{X: 200 + float64(i*100), Y: 450},
					Rotation: 0,
					Scale:    gamemath.Vector2{X: 1.5, Y: 1.5},
				},
				Sprite: sprite,
				Layer:  1,
			}
			scene.AddEntity(entity)
		}
		log.Println("✓ Created 3 additional player sprites (shared texture)")
	}

	log.Println()
	log.Println("Running... Close window to exit")
	log.Println()

	// Run game loop
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}

	// Cleanup - unload textures
	if playerTexture != nil {
		log.Println("Unloading player.png...")
		assets.UnloadTexture("examples/assets/player.png")
		// Note: Texture is still loaded because we loaded it twice (ref count = 1)
		assets.UnloadTexture("examples/assets/player.png")
		// Now ref count = 0, texture actually unloaded
	}
	if enemyTexture != nil {
		assets.UnloadTexture("examples/assets/enemy.png")
	}

	log.Println("Game closed.")
}
