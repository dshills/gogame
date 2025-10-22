package unit

import (
	"testing"

	gamemath "github.com/dshills/gogame/engine/math"
)

func TestTransform_Translate(t *testing.T) {
	tests := []struct {
		name      string
		initial   gamemath.Transform
		dx        float64
		dy        float64
		expectedX float64
		expectedY float64
	}{
		{
			name:      "translate from origin",
			initial:   gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
			dx:        10,
			dy:        20,
			expectedX: 10,
			expectedY: 20,
		},
		{
			name:      "translate positive",
			initial:   gamemath.Transform{Position: gamemath.Vector2{X: 5, Y: 10}},
			dx:        3,
			dy:        7,
			expectedX: 8,
			expectedY: 17,
		},
		{
			name:      "translate negative",
			initial:   gamemath.Transform{Position: gamemath.Vector2{X: 10, Y: 20}},
			dx:        -5,
			dy:        -8,
			expectedX: 5,
			expectedY: 12,
		},
		{
			name:      "translate by zero",
			initial:   gamemath.Transform{Position: gamemath.Vector2{X: 10, Y: 20}},
			dx:        0,
			dy:        0,
			expectedX: 10,
			expectedY: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transform := tt.initial
			transform.Translate(tt.dx, tt.dy)
			if transform.Position.X != tt.expectedX || transform.Position.Y != tt.expectedY {
				t.Errorf("After Translate(%v, %v), Position = (%v, %v), want (%v, %v)",
					tt.dx, tt.dy, transform.Position.X, transform.Position.Y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

func TestTransform_Rotate(t *testing.T) {
	tests := []struct {
		name             string
		initial          gamemath.Transform
		degrees          float64
		expectedRotation float64
	}{
		{
			name:             "rotate from zero",
			initial:          gamemath.Transform{Rotation: 0},
			degrees:          90,
			expectedRotation: 90,
		},
		{
			name:             "rotate positive",
			initial:          gamemath.Transform{Rotation: 45},
			degrees:          45,
			expectedRotation: 90,
		},
		{
			name:             "rotate negative",
			initial:          gamemath.Transform{Rotation: 90},
			degrees:          -45,
			expectedRotation: 45,
		},
		{
			name:             "rotate full circle",
			initial:          gamemath.Transform{Rotation: 0},
			degrees:          360,
			expectedRotation: 360,
		},
		{
			name:             "rotate beyond 360",
			initial:          gamemath.Transform{Rotation: 270},
			degrees:          180,
			expectedRotation: 450,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transform := tt.initial
			transform.Rotate(tt.degrees)
			if transform.Rotation != tt.expectedRotation {
				t.Errorf("After Rotate(%v), Rotation = %v, want %v",
					tt.degrees, transform.Rotation, tt.expectedRotation)
			}
		})
	}
}

func TestTransform_InitialState(t *testing.T) {
	// Test default zero values
	transform := gamemath.Transform{}
	if transform.Position.X != 0 || transform.Position.Y != 0 {
		t.Errorf("Default Position = (%v, %v), want (0, 0)", transform.Position.X, transform.Position.Y)
	}
	if transform.Rotation != 0 {
		t.Errorf("Default Rotation = %v, want 0", transform.Rotation)
	}
	if transform.Scale.X != 0 || transform.Scale.Y != 0 {
		t.Errorf("Default Scale = (%v, %v), want (0, 0)", transform.Scale.X, transform.Scale.Y)
	}
}

func TestTransform_CombinedOperations(t *testing.T) {
	// Test multiple operations in sequence
	transform := gamemath.Transform{
		Position: gamemath.Vector2{X: 100, Y: 100},
		Rotation: 0,
		Scale:    gamemath.Vector2{X: 1, Y: 1},
	}

	transform.Translate(50, 25)
	transform.Rotate(90)
	transform.Translate(-20, -10)
	transform.Rotate(45)

	if transform.Position.X != 130 || transform.Position.Y != 115 {
		t.Errorf("After combined operations, Position = (%v, %v), want (130, 115)",
			transform.Position.X, transform.Position.Y)
	}
	if transform.Rotation != 135 {
		t.Errorf("After combined operations, Rotation = %v, want 135", transform.Rotation)
	}
}

func TestColor_PredefinedColors(t *testing.T) {
	tests := []struct {
		name  string
		color gamemath.Color
		r     uint8
		g     uint8
		b     uint8
		a     uint8
	}{
		{"White", gamemath.White, 255, 255, 255, 255},
		{"Black", gamemath.Black, 0, 0, 0, 255},
		{"Red", gamemath.Red, 255, 0, 0, 255},
		{"Green", gamemath.Green, 0, 255, 0, 255},
		{"Blue", gamemath.Blue, 0, 0, 255, 255},
		{"Transparent", gamemath.Transparent, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color.R != tt.r || tt.color.G != tt.g || tt.color.B != tt.b || tt.color.A != tt.a {
				t.Errorf("%s = RGBA(%v, %v, %v, %v), want RGBA(%v, %v, %v, %v)",
					tt.name, tt.color.R, tt.color.G, tt.color.B, tt.color.A, tt.r, tt.g, tt.b, tt.a)
			}
		})
	}
}

func TestColor_CustomColors(t *testing.T) {
	tests := []struct {
		name  string
		color gamemath.Color
		r     uint8
		g     uint8
		b     uint8
		a     uint8
	}{
		{
			name:  "custom opaque",
			color: gamemath.Color{R: 100, G: 150, B: 200, A: 255},
			r:     100,
			g:     150,
			b:     200,
			a:     255,
		},
		{
			name:  "custom semi-transparent",
			color: gamemath.Color{R: 255, G: 128, B: 64, A: 128},
			r:     255,
			g:     128,
			b:     64,
			a:     128,
		},
		{
			name:  "custom fully transparent",
			color: gamemath.Color{R: 255, G: 255, B: 255, A: 0},
			r:     255,
			g:     255,
			b:     255,
			a:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color.R != tt.r || tt.color.G != tt.g || tt.color.B != tt.b || tt.color.A != tt.a {
				t.Errorf("Color = RGBA(%v, %v, %v, %v), want RGBA(%v, %v, %v, %v)",
					tt.color.R, tt.color.G, tt.color.B, tt.color.A, tt.r, tt.g, tt.b, tt.a)
			}
		})
	}
}
