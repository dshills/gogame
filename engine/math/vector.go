package math

import "math"

// Vector2 represents a 2D vector for positions, velocities, and offsets
type Vector2 struct {
	X float64
	Y float64
}

// Add returns the vector sum of v and other
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub returns the vector difference of v and other
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Scale returns a scaled vector by the given factor
func (v Vector2) Scale(factor float64) Vector2 {
	return Vector2{
		X: v.X * factor,
		Y: v.Y * factor,
	}
}

// Length returns the magnitude of the vector
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a unit vector in the same direction
// If the vector has zero length, returns a zero vector
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{
		X: v.X / length,
		Y: v.Y / length,
	}
}

// Distance returns the distance between v and other
func (v Vector2) Distance(other Vector2) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}
