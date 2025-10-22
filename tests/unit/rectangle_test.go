package unit

import (
	"testing"

	gamemath "github.com/dshills/gogame/engine/math"
)

func TestRectangle_Intersects(t *testing.T) {
	tests := []struct {
		name     string
		r1       gamemath.Rectangle
		r2       gamemath.Rectangle
		expected bool
	}{
		{
			name:     "overlapping rectangles",
			r1:       gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       gamemath.Rectangle{X: 50, Y: 50, Width: 100, Height: 100},
			expected: true,
		},
		{
			name:     "non-overlapping rectangles",
			r1:       gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       gamemath.Rectangle{X: 200, Y: 200, Width: 100, Height: 100},
			expected: false,
		},
		{
			name:     "touching edges (left-right)",
			r1:       gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       gamemath.Rectangle{X: 100, Y: 0, Width: 100, Height: 100},
			expected: false,
		},
		{
			name:     "touching edges (top-bottom)",
			r1:       gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       gamemath.Rectangle{X: 0, Y: 100, Width: 100, Height: 100},
			expected: false,
		},
		{
			name:     "one rectangle inside another",
			r1:       gamemath.Rectangle{X: 0, Y: 0, Width: 200, Height: 200},
			r2:       gamemath.Rectangle{X: 50, Y: 50, Width: 50, Height: 50},
			expected: true,
		},
		{
			name:     "partial overlap top-left",
			r1:       gamemath.Rectangle{X: 50, Y: 50, Width: 100, Height: 100},
			r2:       gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			expected: true,
		},
		{
			name:     "partial overlap bottom-right",
			r1:       gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       gamemath.Rectangle{X: 50, Y: 50, Width: 100, Height: 100},
			expected: true,
		},
		{
			name:     "identical rectangles",
			r1:       gamemath.Rectangle{X: 10, Y: 20, Width: 30, Height: 40},
			r2:       gamemath.Rectangle{X: 10, Y: 20, Width: 30, Height: 40},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.r1.Intersects(tt.r2)
			if result != tt.expected {
				t.Errorf("Intersects() = %v, want %v", result, tt.expected)
			}
			// Test symmetry: r1.Intersects(r2) should equal r2.Intersects(r1)
			reverseResult := tt.r2.Intersects(tt.r1)
			if result != reverseResult {
				t.Errorf("Intersects not symmetric: r1.Intersects(r2)=%v, r2.Intersects(r1)=%v", result, reverseResult)
			}
		})
	}
}

func TestRectangle_Contains(t *testing.T) {
	r := gamemath.Rectangle{X: 10, Y: 20, Width: 100, Height: 80}

	tests := []struct {
		name     string
		x        float64
		y        float64
		expected bool
	}{
		{
			name:     "point inside",
			x:        50,
			y:        50,
			expected: true,
		},
		{
			name:     "point at top-left corner",
			x:        10,
			y:        20,
			expected: true,
		},
		{
			name:     "point at bottom-right corner",
			x:        110,
			y:        100,
			expected: true,
		},
		{
			name:     "point on left edge",
			x:        10,
			y:        50,
			expected: true,
		},
		{
			name:     "point on right edge",
			x:        110,
			y:        50,
			expected: true,
		},
		{
			name:     "point on top edge",
			x:        50,
			y:        20,
			expected: true,
		},
		{
			name:     "point on bottom edge",
			x:        50,
			y:        100,
			expected: true,
		},
		{
			name:     "point outside left",
			x:        5,
			y:        50,
			expected: false,
		},
		{
			name:     "point outside right",
			x:        115,
			y:        50,
			expected: false,
		},
		{
			name:     "point outside top",
			x:        50,
			y:        15,
			expected: false,
		},
		{
			name:     "point outside bottom",
			x:        50,
			y:        105,
			expected: false,
		},
		{
			name:     "point far outside",
			x:        1000,
			y:        1000,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Contains(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Contains(%v, %v) = %v, want %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestRectangle_Center(t *testing.T) {
	tests := []struct {
		name     string
		rect     gamemath.Rectangle
		expected gamemath.Vector2
	}{
		{
			name:     "square at origin",
			rect:     gamemath.Rectangle{X: 0, Y: 0, Width: 100, Height: 100},
			expected: gamemath.Vector2{X: 50, Y: 50},
		},
		{
			name:     "rectangle offset",
			rect:     gamemath.Rectangle{X: 10, Y: 20, Width: 100, Height: 80},
			expected: gamemath.Vector2{X: 60, Y: 60},
		},
		{
			name:     "small rectangle",
			rect:     gamemath.Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			expected: gamemath.Vector2{X: 5, Y: 5},
		},
		{
			name:     "rectangle with negative position",
			rect:     gamemath.Rectangle{X: -50, Y: -50, Width: 100, Height: 100},
			expected: gamemath.Vector2{X: 0, Y: 0},
		},
		{
			name:     "wide rectangle",
			rect:     gamemath.Rectangle{X: 0, Y: 0, Width: 200, Height: 50},
			expected: gamemath.Vector2{X: 100, Y: 25},
		},
		{
			name:     "tall rectangle",
			rect:     gamemath.Rectangle{X: 0, Y: 0, Width: 50, Height: 200},
			expected: gamemath.Vector2{X: 25, Y: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.Center()
			if result.X != tt.expected.X || result.Y != tt.expected.Y {
				t.Errorf("Center() = %v, want %v", result, tt.expected)
			}
		})
	}
}
