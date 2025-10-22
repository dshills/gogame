package math

// Transform represents position, rotation, and scale for entity placement.
type Transform struct {
	Position Vector2 // World position
	Rotation float64 // Angle in degrees (0° = right, 90° = down)
	Scale    Vector2 // Scale factors (1.0 = normal)
}

// Translate moves the transform by the given offset.
func (t *Transform) Translate(dx, dy float64) {
	t.Position.X += dx
	t.Position.Y += dy
}

// Rotate rotates the transform by the given angle in degrees.
func (t *Transform) Rotate(degrees float64) {
	t.Rotation += degrees
}
