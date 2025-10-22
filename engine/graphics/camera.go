package graphics

import gamemath "github.com/dshills/gogame/engine/math"

// Camera defines view transformation from world to screen space
type Camera struct {
	Position     gamemath.Vector2 // Camera center in world space
	Zoom         float64          // Zoom factor (1.0 = normal, >1.0 = zoomed in)
	screenWidth  int              // Cached screen dimensions
	screenHeight int              // Cached screen dimensions
}

// NewCamera creates a camera at origin with no zoom
//
// Returns:
//
//	*Camera: Camera at (0,0) with zoom 1.0
func NewCamera() *Camera {
	return &Camera{
		Position:     gamemath.Vector2{X: 0, Y: 0},
		Zoom:         1.0,
		screenWidth:  800, // Default, will be updated by engine
		screenHeight: 600,
	}
}

// SetScreenSize updates the camera's screen dimensions (called by engine on resize)
func (c *Camera) SetScreenSize(width, height int) {
	c.screenWidth = width
	c.screenHeight = height
}

// WorldToScreen transforms world coordinates to screen pixels
//
// Parameters:
//
//	worldX, worldY: World coordinates
//
// Returns:
//
//	screenX, screenY: Screen pixel coordinates
//
// Example:
//
//	screenX, screenY := camera.WorldToScreen(entity.Transform.Position.X, entity.Transform.Position.Y)
func (c *Camera) WorldToScreen(worldX, worldY float64) (screenX, screenY int) {
	// Transform: world position - camera position, then apply zoom, then add screen center
	relX := (worldX - c.Position.X) * c.Zoom
	relY := (worldY - c.Position.Y) * c.Zoom

	screenX = int(relX + float64(c.screenWidth)/2)
	screenY = int(relY + float64(c.screenHeight)/2)
	return
}

// ScreenToWorld transforms screen pixels to world coordinates
//
// Parameters:
//
//	screenX, screenY: Screen pixel coordinates
//
// Returns:
//
//	worldX, worldY: World coordinates
//
// Example:
//
//	worldX, worldY := camera.ScreenToWorld(mouseX, mouseY)
//	entities := scene.GetEntitiesAt(worldX, worldY)
func (c *Camera) ScreenToWorld(screenX, screenY int) (worldX, worldY float64) {
	// Inverse transform: remove screen center, reverse zoom, then add camera position
	relX := (float64(screenX) - float64(c.screenWidth)/2) / c.Zoom
	relY := (float64(screenY) - float64(c.screenHeight)/2) / c.Zoom

	worldX = relX + c.Position.X
	worldY = relY + c.Position.Y
	return
}

// Follow smoothly moves camera toward target
//
// Parameters:
//
//	targetX, targetY: Target world position
//	smoothing: Interpolation factor (0.0 = instant, 1.0 = no follow)
//
// Example:
//
//	camera.Follow(player.Transform.Position.X, player.Transform.Position.Y, 0.1)
func (c *Camera) Follow(targetX, targetY float64, smoothing float64) {
	// Linear interpolation: camera position moves toward target
	c.Position.X += (targetX - c.Position.X) * (1.0 - smoothing)
	c.Position.Y += (targetY - c.Position.Y) * (1.0 - smoothing)
}
