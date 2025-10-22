package graphics

import gamemath "github.com/dshills/gogame/engine/math"

// Sprite represents a visual representation attached to entities.
type Sprite struct {
	Texture    *Texture           // Loaded texture (via AssetManager)
	SourceRect gamemath.Rectangle // Region of texture to render (for sprite sheets)
	Color      gamemath.Color     // Tint color (white = no tint)
	Alpha      float64            // Opacity (0.0 = transparent, 1.0 = opaque)
	FlipH      bool               // Flip horizontally
	FlipV      bool               // Flip vertically
}

// NewSprite creates a sprite from a texture
//
// Parameters:
//
//	texture: Loaded texture
//
// Returns:
//
//	*Sprite: Sprite rendering full texture
//
// Example:
//
//	texture, _ := assets.LoadTexture("player.png")
//	sprite := graphics.NewSprite(texture)
func NewSprite(texture *Texture) *Sprite {
	return &Sprite{
		Texture: texture,
		SourceRect: gamemath.Rectangle{
			X:      0,
			Y:      0,
			Width:  float64(texture.Width),
			Height: float64(texture.Height),
		},
		Color: gamemath.White,
		Alpha: 1.0,
		FlipH: false,
		FlipV: false,
	}
}

// SetSourceRect sets the sprite sheet region
//
// Parameters:
//
//	x, y: Top-left corner in texture
//	width, height: Region dimensions
//
// Example:
//
//	// Extract 32x32 sprite from sprite sheet
//	sprite.SetSourceRect(64, 0, 32, 32)
func (s *Sprite) SetSourceRect(x, y, width, height int) {
	s.SourceRect = gamemath.Rectangle{
		X:      float64(x),
		Y:      float64(y),
		Width:  float64(width),
		Height: float64(height),
	}
}

// SetColor sets the tint color
//
// Parameters:
//
//	color: RGBA tint (white = no tint)
//
// Example:
//
//	sprite.SetColor(math.Color{R: 255, G: 0, B: 0, A: 255})  // Red tint
func (s *Sprite) SetColor(color gamemath.Color) {
	s.Color = color
}
