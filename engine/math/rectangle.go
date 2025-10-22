package math

// Rectangle represents an axis-aligned rectangle for bounds, regions, and collision.
type Rectangle struct {
	X      float64 // Left edge
	Y      float64 // Top edge
	Width  float64
	Height float64
}

// Intersects checks if this rectangle overlaps with another.
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// Contains checks if a point is inside this rectangle.
func (r Rectangle) Contains(x, y float64) bool {
	return x >= r.X &&
		x <= r.X+r.Width &&
		y >= r.Y &&
		y <= r.Y+r.Height
}

// Center returns the center point of the rectangle.
func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}
