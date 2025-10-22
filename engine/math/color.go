package math

// Color represents an RGBA color with 8-bit channels
type Color struct {
	R uint8 // Red (0-255)
	G uint8 // Green (0-255)
	B uint8 // Blue (0-255)
	A uint8 // Alpha (0-255, 255 = opaque)
}

// Predefined colors
var (
	White       = Color{255, 255, 255, 255}
	Black       = Color{0, 0, 0, 255}
	Red         = Color{255, 0, 0, 255}
	Green       = Color{0, 255, 0, 255}
	Blue        = Color{0, 0, 255, 255}
	Transparent = Color{0, 0, 0, 0}
)
