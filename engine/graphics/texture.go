package graphics

import "github.com/veandco/go-sdl2/sdl"

// Texture represents a loaded image texture
type Texture struct {
	sdlTexture *sdl.Texture // SDL texture handle (internal)
	Width      int          // Texture width in pixels
	Height     int          // Texture height in pixels
	Path       string       // Source file path
}

// NewTexture creates a new texture wrapper around an SDL texture
func NewTexture(sdlTexture *sdl.Texture, width, height int, path string) *Texture {
	return &Texture{
		sdlTexture: sdlTexture,
		Width:      width,
		Height:     height,
		Path:       path,
	}
}

// Destroy releases the SDL texture resources
func (t *Texture) Destroy() error {
	if t.sdlTexture != nil {
		return t.sdlTexture.Destroy()
	}
	return nil
}

// GetSDLTexture returns the underlying SDL texture (for internal use)
func (t *Texture) GetSDLTexture() *sdl.Texture {
	return t.sdlTexture
}
