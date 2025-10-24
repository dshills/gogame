package integration

import (
	"runtime"
	"testing"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
)

// TestTextureLoadingInGameLoop tests loading textures during game loop.
func TestTextureLoadingInGameLoop(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	engine, err := core.NewEngine("Asset Test", 800, 600, false)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	// Try to load test texture
	texture, err := engine.Assets().LoadTexture("examples/assets/player.png")
	if err != nil {
		t.Skip("Test texture not available, skipping")
		return
	}

	// Create sprite with texture
	sprite := graphics.NewSprite(texture)

	// Create scene and entity
	scene := core.NewScene()
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}

	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 400, Y: 300}},
		Sprite:    sprite,
		Layer:     0,
	}
	scene.AddEntity(entity)

	// Verify rendering doesn't crash
	scene.Render(engine.Renderer(), camera)
}

// TestMultipleSpritesSameTexture tests texture sharing.
func TestMultipleSpritesSameTexture(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	engine, err := core.NewEngine("Texture Sharing", 800, 600, false)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	texture, err := engine.Assets().LoadTexture("examples/assets/player.png")
	if err != nil {
		t.Skip("Test texture not available")
		return
	}

	// Create multiple sprites with same texture
	sprite1 := graphics.NewSprite(texture)
	sprite2 := graphics.NewSprite(texture)

	if sprite1 == nil || sprite2 == nil {
		t.Error("Failed to create sprites")
	}

	// Both sprites should reference same underlying texture
	// (actual verification would need access to internal texture refs)
}
