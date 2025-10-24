package unit

import (
	"testing"

	"github.com/dshills/gogame/engine/core"
	gamemath "github.com/dshills/gogame/engine/math"
)

// mockBehavior is a test behavior that tracks update calls.
type mockBehavior struct {
	updateCalled bool
	lastDt       float64
	updateCount  int
}

func (mb *mockBehavior) Update(entity *core.Entity, dt float64) {
	mb.updateCalled = true
	mb.lastDt = dt
	mb.updateCount++
}

// TestEntityUpdate_WithBehavior tests that Entity.Update calls Behavior.Update.
func TestEntityUpdate_WithBehavior(t *testing.T) {
	behavior := &mockBehavior{}
	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
		Behavior:  behavior,
	}

	// Call update
	dt := 0.016667
	entity.Update(dt)

	// Verify behavior was called
	if !behavior.updateCalled {
		t.Error("Expected Behavior.Update to be called")
	}

	if behavior.lastDt != dt {
		t.Errorf("Expected dt=%f, got %f", dt, behavior.lastDt)
	}
}

// TestEntityUpdate_WithoutBehavior tests that Update works with nil behavior.
func TestEntityUpdate_WithoutBehavior(t *testing.T) {
	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
		Behavior:  nil,
	}

	// Should not crash
	entity.Update(0.016)
}

// TestEntityUpdate_MultipleFrames tests behavior across multiple updates.
func TestEntityUpdate_MultipleFrames(t *testing.T) {
	behavior := &mockBehavior{}
	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
		Behavior:  behavior,
	}

	// Update multiple times
	for i := 0; i < 10; i++ {
		entity.Update(0.016)
	}

	if behavior.updateCount != 10 {
		t.Errorf("Expected 10 updates, got %d", behavior.updateCount)
	}
}

// TestEntityUpdate_DeltaTimeVariation tests with different delta times.
func TestEntityUpdate_DeltaTimeVariation(t *testing.T) {
	tests := []struct {
		name string
		dt   float64
	}{
		{"60 FPS", 0.016667},
		{"30 FPS", 0.033333},
		{"120 FPS", 0.008333},
		{"slow frame", 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			behavior := &mockBehavior{}
			entity := &core.Entity{
				Active:    true,
				Transform: gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
				Behavior:  behavior,
			}

			entity.Update(tt.dt)

			if behavior.lastDt != tt.dt {
				t.Errorf("Expected dt=%f, got %f", tt.dt, behavior.lastDt)
			}
		})
	}
}

// velocityBehavior moves entity by velocity each frame.
type velocityBehavior struct {
	velocity gamemath.Vector2
}

func (vb *velocityBehavior) Update(entity *core.Entity, dt float64) {
	entity.Transform.Position.X += vb.velocity.X * dt
	entity.Transform.Position.Y += vb.velocity.Y * dt
}

// TestEntityUpdate_VelocityIntegration tests realistic velocity-based movement.
func TestEntityUpdate_VelocityIntegration(t *testing.T) {
	behavior := &velocityBehavior{
		velocity: gamemath.Vector2{X: 100, Y: 50}, // pixels per second
	}

	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
		Behavior:  behavior,
	}

	// Update for 1 second at 60 FPS
	dt := 1.0 / 60.0
	for i := 0; i < 60; i++ {
		entity.Update(dt)
	}

	// After 1 second, entity should have moved by velocity
	expectedX := 100.0
	expectedY := 50.0

	// Allow small floating point error
	if !almostEqual(entity.Transform.Position.X, expectedX, 0.01) {
		t.Errorf("Expected X=%f, got %f", expectedX, entity.Transform.Position.X)
	}

	if !almostEqual(entity.Transform.Position.Y, expectedY, 0.01) {
		t.Errorf("Expected Y=%f, got %f", expectedY, entity.Transform.Position.Y)
	}
}

func almostEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}
