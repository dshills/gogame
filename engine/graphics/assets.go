package graphics

import (
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// AssetManager manages texture loading and caching
type AssetManager struct {
	renderer *sdl.Renderer
	textures map[string]*Texture // Cache of loaded textures
	refCount map[string]int      // Reference counting
}

// NewAssetManager creates a new asset manager
func NewAssetManager(renderer *sdl.Renderer) *AssetManager {
	return &AssetManager{
		renderer: renderer,
		textures: make(map[string]*Texture),
		refCount: make(map[string]int),
	}
}

// LoadTexture loads a texture from disk or returns cached
//
// Parameters:
//
//	path: File path (PNG or JPEG)
//
// Returns:
//
//	*Texture: Loaded texture
//	error: Non-nil if file not found or decode fails
//
// Behavior:
//   - Returns existing texture if already loaded
//   - Increments reference count
//   - Caches texture
//
// Example:
//
//	texture, err := assets.LoadTexture("assets/player.png")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (am *AssetManager) LoadTexture(path string) (*Texture, error) {
	// Check if already loaded
	if texture, exists := am.textures[path]; exists {
		am.refCount[path]++
		return texture, nil
	}

	// Load image file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: file not found: %s: %w", path, err)
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %s: %w", path, err)
	}

	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create SDL surface from image data
	surface, err := sdl.CreateRGBSurface(
		0,
		int32(width),
		int32(height),
		32,
		0x000000ff,
		0x0000ff00,
		0x00ff0000,
		0xff000000,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create surface: %w", err)
	}
	defer surface.Free()

	// Copy image data to surface
	pixels := surface.Pixels()
	pixelIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Convert from 16-bit to 8-bit
			pixels[pixelIndex] = uint8(r >> 8)   // R
			pixels[pixelIndex+1] = uint8(g >> 8) // G
			pixels[pixelIndex+2] = uint8(b >> 8) // B
			pixels[pixelIndex+3] = uint8(a >> 8) // A
			pixelIndex += 4
		}
	}

	// Create SDL texture from surface
	sdlTexture, err := am.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture: %w", err)
	}

	// Set blend mode for alpha transparency
	sdlTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	// Wrap in our Texture type
	texture := NewTexture(sdlTexture, width, height, path)

	// Cache texture
	am.textures[path] = texture
	am.refCount[path] = 1

	_ = format // Suppress unused variable warning
	return texture, nil
}

// UnloadTexture decrements reference count
//
// Parameters:
//
//	path: File path of texture to unload
//
// Behavior:
//   - Decrements reference count
//   - Unloads if count reaches zero
//   - Safe to call multiple times
//   - No-op if texture not loaded
//
// Example:
//
//	assets.UnloadTexture("assets/player.png")
func (am *AssetManager) UnloadTexture(path string) {
	if _, exists := am.textures[path]; !exists {
		return // Not loaded
	}

	am.refCount[path]--

	// Unload if no more references
	if am.refCount[path] <= 0 {
		if texture, exists := am.textures[path]; exists {
			texture.Destroy()
			delete(am.textures, path)
			delete(am.refCount, path)
		}
	}
}

// Destroy unloads all textures
func (am *AssetManager) Destroy() {
	for path, texture := range am.textures {
		texture.Destroy()
		delete(am.textures, path)
		delete(am.refCount, path)
	}
}
