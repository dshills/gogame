package graphics

import (
	"fmt"

	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Font represents a loaded TTF font.
type Font struct {
	font *ttf.Font
	size int
}

// LoadFont loads a TTF font from file.
//
// Parameters:
//
//	path: Path to TTF font file
//	size: Font size in points
//
// Returns:
//
//	*Font: Loaded font
//	error: Non-nil if font loading fails
//
// Example:
//
//	font, err := graphics.LoadFont("/System/Library/Fonts/Helvetica.ttc", 24)
func LoadFont(path string, size int) (*Font, error) {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	return &Font{
		font: font,
		size: size,
	}, nil
}

// Close closes the font and frees resources.
func (f *Font) Close() {
	if f.font != nil {
		f.font.Close()
	}
}

// RenderText renders text to a texture.
//
// Parameters:
//
//	renderer: SDL renderer
//	text: Text to render
//	color: Text color
//
// Returns:
//
//	*sdl.Texture: Rendered text texture
//	int32: Texture width
//	int32: Texture height
//	error: Non-nil if rendering fails
func (f *Font) RenderText(renderer *sdl.Renderer, text string, color gamemath.Color) (*sdl.Texture, int32, int32, error) {
	// Create surface with text
	surface, err := f.font.RenderUTF8Blended(text, sdl.Color{
		R: color.R,
		G: color.G,
		B: color.B,
		A: color.A,
	})
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to render text surface: %w", err)
	}
	defer surface.Free()

	// Create texture from surface
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create texture from surface: %w", err)
	}

	return texture, surface.W, surface.H, nil
}

// TextRenderer provides high-level text rendering with caching.
type TextRenderer struct {
	renderer *sdl.Renderer
	font     *Font
}

// NewTextRenderer creates a new text renderer.
func NewTextRenderer(renderer *sdl.Renderer, font *Font) *TextRenderer {
	return &TextRenderer{
		renderer: renderer,
		font:     font,
	}
}

// DrawText renders text at a position.
//
// Parameters:
//
//	text: Text to render
//	x, y: Screen position (top-left corner)
//	color: Text color
//
// Returns:
//
//	error: Non-nil if rendering fails
//
// Example:
//
//	err := textRenderer.DrawText("Score: 100", 10, 10, gamemath.White)
func (tr *TextRenderer) DrawText(text string, x, y int, color gamemath.Color) error {
	if text == "" {
		return nil
	}

	// Render text to texture
	texture, width, height, err := tr.font.RenderText(tr.renderer, text, color)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	// Draw texture at position
	destRect := sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: width,
		H: height,
	}

	return tr.renderer.Copy(texture, nil, &destRect)
}

// MeasureText returns the dimensions of rendered text.
//
// Parameters:
//
//	text: Text to measure
//
// Returns:
//
//	width: Text width in pixels
//	height: Text height in pixels
//	error: Non-nil if measurement fails
func (tr *TextRenderer) MeasureText(text string) (int, int, error) {
	w, h, err := tr.font.font.SizeUTF8(text)
	return w, h, err
}
