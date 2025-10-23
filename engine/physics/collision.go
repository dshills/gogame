package physics

import (
	gamemath "github.com/dshills/gogame/engine/math"
)

// Entity interface for collision detection (to avoid circular imports).
type Entity interface {
	GetID() uint64
	GetTransform() gamemath.Transform
	GetCollider() *Collider
	IsActive() bool
}

// CollisionPair represents two entities that are colliding.
type CollisionPair struct {
	EntityA Entity
	EntityB Entity
}

// DetectCollisions performs O(n²) broad-phase collision detection.
//
// Parameters:
//
//	entities: Slice of entities to check
//
// Returns:
//
//	[]CollisionPair: All colliding pairs
//
// Note:
//
//	This is a simple implementation suitable for <1000 entities.
//	For larger worlds, use spatial partitioning (quadtree/grid).
//
// Example:
//
//	collisions := physics.DetectCollisions(scene.GetEntities())
//	for _, pair := range collisions {
//	    // Handle collision between pair.EntityA and pair.EntityB
//	}
func DetectCollisions(entities []Entity) []CollisionPair {
	var collisions []CollisionPair

	// O(n²) broad phase - check all pairs
	for i := 0; i < len(entities); i++ {
		entityA := entities[i]

		// Skip inactive entities or entities without colliders
		if !entityA.IsActive() || entityA.GetCollider() == nil {
			continue
		}

		for j := i + 1; j < len(entities); j++ {
			entityB := entities[j]

			// Skip inactive entities or entities without colliders
			if !entityB.IsActive() || entityB.GetCollider() == nil {
				continue
			}

			// Test collision
			colliderA := entityA.GetCollider()
			colliderB := entityB.GetCollider()

			if colliderA.Intersects(colliderB, entityA.GetTransform(), entityB.GetTransform()) {
				collisions = append(collisions, CollisionPair{
					EntityA: entityA,
					EntityB: entityB,
				})
			}
		}
	}

	return collisions
}
