package integration

import (
	"testing"

	"github.com/dshills/gogame/engine/core"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// TestCollisionCallbacks tests OnCollisionEnter and OnCollisionExit.
func TestCollisionCallbacks(t *testing.T) {
	scene := core.NewScene()

	// Track collision events
	enterCalled := false
	exitCalled := false

	// Entity 1
	collider1 := physics.NewCollider(50, 50)
	collider1.CollisionLayer = 0
	collider1.CollisionMask = 0xFF

	entity1 := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 100, Y: 100}},
		Collider:  collider1,
		OnCollisionEnter: func(self, other *core.Entity) {
			enterCalled = true
		},
		OnCollisionExit: func(self, other *core.Entity) {
			exitCalled = true
		},
		Layer: 0,
	}
	scene.AddEntity(entity1)

	// Entity 2 (starts far away)
	collider2 := physics.NewCollider(50, 50)
	collider2.CollisionLayer = 0
	collider2.CollisionMask = 0xFF

	entity2 := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 200, Y: 200}},
		Collider:  collider2,
		Layer:     0,
	}
	scene.AddEntity(entity2)

	// Update - no collision yet
	scene.Update(0.016)

	if enterCalled {
		t.Error("OnCollisionEnter should not be called yet")
	}

	// Move entity2 to overlap entity1
	entity2.Transform.Position = gamemath.Vector2{X: 110, Y: 100}
	scene.Update(0.016)

	// OnCollisionEnter should be called
	if !enterCalled {
		t.Error("Expected OnCollisionEnter to be called")
	}

	// Move entity2 away
	entity2.Transform.Position = gamemath.Vector2{X: 300, Y: 300}
	scene.Update(0.016)

	// OnCollisionExit should be called
	if !exitCalled {
		t.Error("Expected OnCollisionExit to be called")
	}
}

// TestCollisionStay tests OnCollisionStay callback.
func TestCollisionStay(t *testing.T) {
	scene := core.NewScene()

	stayCount := 0

	collider1 := physics.NewCollider(50, 50)
	collider1.CollisionLayer = 0
	collider1.CollisionMask = 0xFF

	entity1 := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 100, Y: 100}},
		Collider:  collider1,
		OnCollisionStay: func(self, other *core.Entity) {
			stayCount++
		},
		Layer: 0,
	}
	scene.AddEntity(entity1)

	collider2 := physics.NewCollider(50, 50)
	collider2.CollisionLayer = 0
	collider2.CollisionMask = 0xFF

	entity2 := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 110, Y: 100}},
		Collider:  collider2,
		Layer:     0,
	}
	scene.AddEntity(entity2)

	// Update multiple times while overlapping
	for i := 0; i < 5; i++ {
		scene.Update(0.016)
	}

	// OnCollisionStay should be called multiple times
	if stayCount < 5 {
		t.Errorf("Expected at least 5 OnCollisionStay calls, got %d", stayCount)
	}
}
