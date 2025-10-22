package graphics

import (
	"fmt"

	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/veandco/go-sdl2/sdl"
)

// Renderer wraps SDL2 rendering operations
type Renderer struct {
	sdlRenderer *sdl.Renderer
}

// NewRenderer creates a renderer from an SDL renderer
func NewRenderer(sdlRenderer *sdl.Renderer) *Renderer {
	return &Renderer{
		sdlRenderer: sdlRenderer,
	}
}

// Clear clears the screen with the specified color
func (r *Renderer) Clear(color gamemath.Color) error {
	if err := r.sdlRenderer.SetDrawColor(color.R, color.G, color.B, color.A); err != nil {
		return fmt.Errorf("failed to set draw color: %w", err)
	}
	if err := r.sdlRenderer.Clear(); err != nil {
		return fmt.Errorf("failed to clear screen: %w", err)
	}
	return nil
}

// Present presents the rendered frame to the screen
func (r *Renderer) Present() {
	r.sdlRenderer.Present()
}

// DrawSprite renders a sprite at the specified transform with camera transform applied
func (r *Renderer) DrawSprite(sprite *Sprite, transform gamemath.Transform, camera *Camera) error {
	if sprite == nil || sprite.Texture == nil {
		return nil // Nothing to render
	}

	// Convert world position to screen position via camera
	screenX, screenY := camera.WorldToScreen(transform.Position.X, transform.Position.Y)

	// Calculate final dimensions with scale
	finalWidth := int(sprite.SourceRect.Width * transform.Scale.X * camera.Zoom)
	finalHeight := int(sprite.SourceRect.Height * transform.Scale.Y * camera.Zoom)

	// Create source rectangle (region of texture to render)
	srcRect := &sdl.Rect{
		X: int32(sprite.SourceRect.X),
		Y: int32(sprite.SourceRect.Y),
		W: int32(sprite.SourceRect.Width),
		H: int32(sprite.SourceRect.Height),
	}

	// Create destination rectangle (where to render on screen)
	// Center the sprite at the screen position
	dstRect := &sdl.Rect{
		X: int32(screenX - finalWidth/2),
		Y: int32(screenY - finalHeight/2),
		W: int32(finalWidth),
		H: int32(finalHeight),
	}

	// Apply color tint
	texture := sprite.Texture.GetSDLTexture()
	if err := texture.SetColorMod(sprite.Color.R, sprite.Color.G, sprite.Color.B); err != nil {
		return fmt.Errorf("failed to set color mod: %w", err)
	}

	// Apply alpha
	alpha := uint8(sprite.Alpha * 255)
	if err := texture.SetAlphaMod(alpha); err != nil {
		return fmt.Errorf("failed to set alpha mod: %w", err)
	}

	// Determine flip mode
	flip := sdl.FLIP_NONE
	if sprite.FlipH && sprite.FlipV {
		flip = sdl.FLIP_HORIZONTAL | sdl.FLIP_VERTICAL
	} else if sprite.FlipH {
		flip = sdl.FLIP_HORIZONTAL
	} else if sprite.FlipV {
		flip = sdl.FLIP_VERTICAL
	}

	// Render the sprite
	if err := r.sdlRenderer.CopyEx(
		texture,
		srcRect,
		dstRect,
		transform.Rotation, // Rotation angle in degrees
		nil,                // Center point (nil = center of sprite)
		flip,
	); err != nil {
		return fmt.Errorf("failed to render sprite: %w", err)
	}

	return nil
}

// Destroy releases renderer resources
func (r *Renderer) Destroy() error {
	if r.sdlRenderer != nil {
		return r.sdlRenderer.Destroy()
	}
	return nil
}

// GetSDLRenderer returns the underlying SDL renderer (for internal use)
func (r *Renderer) GetSDLRenderer() *sdl.Renderer {
	return r.sdlRenderer
}
