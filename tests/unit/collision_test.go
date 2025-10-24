package unit

import (
	"testing"

	"github.com/dshills/gogame/engine/core"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// TestAABBIntersection tests basic AABB collision detection.
func TestAABBIntersection(t *testing.T) {
	// Create two overlapping colliders
	collider1 := physics.NewCollider(50, 50)
	collider1.CollisionLayer = 0
	collider1.CollisionMask = 1 << 1 // Can collide with layer 1

	collider2 := physics.NewCollider(50, 50)
	collider2.CollisionLayer = 1
	collider2.CollisionMask = 1 << 0 // Can collide with layer 0

	// Create entities
	entity1 := &core.Entity{
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 100, Y: 100}},
		Collider:  collider1,
	}

	entity2 := &core.Entity{
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 120, Y: 100}},
		Collider:  collider2,
	}

	// Get world bounds
	bounds1 := collider1.GetWorldBounds(entity1.Transform)
	bounds2 := collider2.GetWorldBounds(entity2.Transform)

	// Test intersection
	if !bounds1.Intersects(bounds2) {
		t.Error("Expected overlapping colliders to intersect")
	}
}

// TestColliderLayerMask tests collision layer filtering.
func TestColliderLayerMask(t *testing.T) {
	// Layer 0 can only collide with layer 1
	collider1 := physics.NewCollider(50, 50)
	collider1.CollisionLayer = 0
	collider1.CollisionMask = 1 << 1

	// Layer 2 (should not collide with layer 0)
	collider2 := physics.NewCollider(50, 50)
	collider2.CollisionLayer = 2
	collider2.CollisionMask = 1 << 0

	// Even if overlapping, layer mask should prevent collision
	entity1 := &core.Entity{
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 100, Y: 100}},
		Collider:  collider1,
	}

	entity2 := &core.Entity{
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 100, Y: 100}},
		Collider:  collider2,
	}

	// Verify collision check respects masks
	shouldCollide := collider1.Intersects(collider2, entity1.Transform, entity2.Transform)
	if shouldCollide {
		t.Error("Expected layer mask to prevent collision")
	}
}

// TestColliderContains tests point containment.
func TestColliderContains(t *testing.T) {
	collider := physics.NewCollider(100, 100)
	entity := &core.Entity{
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 200, Y: 200}},
		Collider:  collider,
	}

	bounds := collider.GetWorldBounds(entity.Transform)

	// Point inside
	if !bounds.Contains(200, 200) {
		t.Error("Expected center point to be inside collider")
	}

	// Point outside
	if bounds.Contains(1000, 1000) {
		t.Error("Expected distant point to be outside collider")
	}
}
