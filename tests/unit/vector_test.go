package unit

import (
	"math"
	"testing"

	gamemath "github.com/dshills/gogame/engine/math"
)

func TestVector2_Add(t *testing.T) {
	tests := []struct {
		name     string
		v1       gamemath.Vector2
		v2       gamemath.Vector2
		expected gamemath.Vector2
	}{
		{
			name:     "positive vectors",
			v1:       gamemath.Vector2{X: 1.0, Y: 2.0},
			v2:       gamemath.Vector2{X: 3.0, Y: 4.0},
			expected: gamemath.Vector2{X: 4.0, Y: 6.0},
		},
		{
			name:     "negative vectors",
			v1:       gamemath.Vector2{X: -1.0, Y: -2.0},
			v2:       gamemath.Vector2{X: -3.0, Y: -4.0},
			expected: gamemath.Vector2{X: -4.0, Y: -6.0},
		},
		{
			name:     "mixed signs",
			v1:       gamemath.Vector2{X: 5.0, Y: -3.0},
			v2:       gamemath.Vector2{X: -2.0, Y: 7.0},
			expected: gamemath.Vector2{X: 3.0, Y: 4.0},
		},
		{
			name:     "zero vector",
			v1:       gamemath.Vector2{X: 5.0, Y: 3.0},
			v2:       gamemath.Vector2{X: 0.0, Y: 0.0},
			expected: gamemath.Vector2{X: 5.0, Y: 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Add(tt.v2)
			if result.X != tt.expected.X || result.Y != tt.expected.Y {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVector2_Sub(t *testing.T) {
	tests := []struct {
		name     string
		v1       gamemath.Vector2
		v2       gamemath.Vector2
		expected gamemath.Vector2
	}{
		{
			name:     "positive vectors",
			v1:       gamemath.Vector2{X: 5.0, Y: 8.0},
			v2:       gamemath.Vector2{X: 2.0, Y: 3.0},
			expected: gamemath.Vector2{X: 3.0, Y: 5.0},
		},
		{
			name:     "negative result",
			v1:       gamemath.Vector2{X: 2.0, Y: 3.0},
			v2:       gamemath.Vector2{X: 5.0, Y: 8.0},
			expected: gamemath.Vector2{X: -3.0, Y: -5.0},
		},
		{
			name:     "zero vector",
			v1:       gamemath.Vector2{X: 5.0, Y: 3.0},
			v2:       gamemath.Vector2{X: 0.0, Y: 0.0},
			expected: gamemath.Vector2{X: 5.0, Y: 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Sub(tt.v2)
			if result.X != tt.expected.X || result.Y != tt.expected.Y {
				t.Errorf("Sub() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVector2_Scale(t *testing.T) {
	tests := []struct {
		name     string
		v        gamemath.Vector2
		factor   float64
		expected gamemath.Vector2
	}{
		{
			name:     "scale by 2",
			v:        gamemath.Vector2{X: 3.0, Y: 4.0},
			factor:   2.0,
			expected: gamemath.Vector2{X: 6.0, Y: 8.0},
		},
		{
			name:     "scale by 0.5",
			v:        gamemath.Vector2{X: 10.0, Y: 20.0},
			factor:   0.5,
			expected: gamemath.Vector2{X: 5.0, Y: 10.0},
		},
		{
			name:     "scale by negative",
			v:        gamemath.Vector2{X: 3.0, Y: 4.0},
			factor:   -1.0,
			expected: gamemath.Vector2{X: -3.0, Y: -4.0},
		},
		{
			name:     "scale by zero",
			v:        gamemath.Vector2{X: 3.0, Y: 4.0},
			factor:   0.0,
			expected: gamemath.Vector2{X: 0.0, Y: 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Scale(tt.factor)
			if result.X != tt.expected.X || result.Y != tt.expected.Y {
				t.Errorf("Scale() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVector2_Length(t *testing.T) {
	tests := []struct {
		name     string
		v        gamemath.Vector2
		expected float64
	}{
		{
			name:     "3-4-5 triangle",
			v:        gamemath.Vector2{X: 3.0, Y: 4.0},
			expected: 5.0,
		},
		{
			name:     "unit vector X",
			v:        gamemath.Vector2{X: 1.0, Y: 0.0},
			expected: 1.0,
		},
		{
			name:     "unit vector Y",
			v:        gamemath.Vector2{X: 0.0, Y: 1.0},
			expected: 1.0,
		},
		{
			name:     "zero vector",
			v:        gamemath.Vector2{X: 0.0, Y: 0.0},
			expected: 0.0,
		},
		{
			name:     "negative components",
			v:        gamemath.Vector2{X: -3.0, Y: -4.0},
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Length()
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Length() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVector2_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		v        gamemath.Vector2
		expected gamemath.Vector2
	}{
		{
			name:     "3-4-5 triangle",
			v:        gamemath.Vector2{X: 3.0, Y: 4.0},
			expected: gamemath.Vector2{X: 0.6, Y: 0.8},
		},
		{
			name:     "already normalized",
			v:        gamemath.Vector2{X: 1.0, Y: 0.0},
			expected: gamemath.Vector2{X: 1.0, Y: 0.0},
		},
		{
			name:     "zero vector",
			v:        gamemath.Vector2{X: 0.0, Y: 0.0},
			expected: gamemath.Vector2{X: 0.0, Y: 0.0},
		},
		{
			name:     "negative components",
			v:        gamemath.Vector2{X: -5.0, Y: 0.0},
			expected: gamemath.Vector2{X: -1.0, Y: 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Normalize()
			if math.Abs(result.X-tt.expected.X) > 1e-9 || math.Abs(result.Y-tt.expected.Y) > 1e-9 {
				t.Errorf("Normalize() = %v, want %v", result, tt.expected)
			}
			// Verify that normalized vector has length 1 (except zero vector)
			if tt.v.Length() > 0 {
				length := result.Length()
				if math.Abs(length-1.0) > 1e-9 {
					t.Errorf("Normalized vector length = %v, want 1.0", length)
				}
			}
		})
	}
}

func TestVector2_Distance(t *testing.T) {
	tests := []struct {
		name     string
		v1       gamemath.Vector2
		v2       gamemath.Vector2
		expected float64
	}{
		{
			name:     "3-4-5 triangle",
			v1:       gamemath.Vector2{X: 0.0, Y: 0.0},
			v2:       gamemath.Vector2{X: 3.0, Y: 4.0},
			expected: 5.0,
		},
		{
			name:     "horizontal distance",
			v1:       gamemath.Vector2{X: 0.0, Y: 5.0},
			v2:       gamemath.Vector2{X: 10.0, Y: 5.0},
			expected: 10.0,
		},
		{
			name:     "vertical distance",
			v1:       gamemath.Vector2{X: 5.0, Y: 0.0},
			v2:       gamemath.Vector2{X: 5.0, Y: 10.0},
			expected: 10.0,
		},
		{
			name:     "same point",
			v1:       gamemath.Vector2{X: 5.0, Y: 5.0},
			v2:       gamemath.Vector2{X: 5.0, Y: 5.0},
			expected: 0.0,
		},
		{
			name:     "negative coordinates",
			v1:       gamemath.Vector2{X: -3.0, Y: -4.0},
			v2:       gamemath.Vector2{X: 0.0, Y: 0.0},
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Distance(tt.v2)
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Distance() = %v, want %v", result, tt.expected)
			}
		})
	}
}
