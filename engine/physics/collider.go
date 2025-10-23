// Package physics provides collision detection and physics simulation.
package physics

import (
	gamemath "github.com/dshills/gogame/engine/math"
)

// Collider provides AABB collision detection with layer masks.
type Collider struct {
	Bounds         gamemath.Rectangle // Local bounds (relative to entity)
	Offset         gamemath.Vector2   // Offset from entity position
	IsTrigger      bool               // If true, collisions don't block movement
	CollisionLayer int                // Which layer this collider is on (bit position)
	CollisionMask  int                // Which layers this collider can collide with (bitmask)
}

// NewCollider creates a collider with centered bounds.
//
// Parameters:
//
//	width, height: Collider dimensions
//
// Returns:
//
//	*Collider: New collider on layer 0, colliding with all layers
//
// Example:
//
//	collider := physics.NewCollider(32, 32)
//	collider.CollisionLayer = 1 // Player layer
//	collider.CollisionMask = 2 | 4 // Collides with enemies (2) and walls (4)
func NewCollider(width, height float64) *Collider {
	return &Collider{
		Bounds: gamemath.Rectangle{
			X:      -width / 2,  // Centered on entity
			Y:      -height / 2, // Centered on entity
			Width:  width,
			Height: height,
		},
		Offset:         gamemath.Vector2{X: 0, Y: 0},
		IsTrigger:      false,
		CollisionLayer: 0,
		CollisionMask:  0xFFFFFFFF, // Collide with all layers by default
	}
}

// GetWorldBounds transforms local bounds to world space.
//
// Parameters:
//
//	transform: Entity's transform (position, rotation, scale)
//
// Returns:
//
//	gamemath.Rectangle: World-space AABB bounds
//
// Note:
//
//	Currently ignores rotation. Supports position, scale, and offset.
//
// Example:
//
//	worldBounds := collider.GetWorldBounds(entity.Transform)
//	if worldBounds.Contains(point) { ... }
func (c *Collider) GetWorldBounds(transform gamemath.Transform) gamemath.Rectangle {
	// Apply scale to bounds
	scaledWidth := c.Bounds.Width * transform.Scale.X
	scaledHeight := c.Bounds.Height * transform.Scale.Y

	// Apply offset and position
	worldX := transform.Position.X + (c.Offset.X * transform.Scale.X) + (c.Bounds.X * transform.Scale.X)
	worldY := transform.Position.Y + (c.Offset.Y * transform.Scale.Y) + (c.Bounds.Y * transform.Scale.Y)

	return gamemath.Rectangle{
		X:      worldX,
		Y:      worldY,
		Width:  scaledWidth,
		Height: scaledHeight,
	}
}

// Intersects tests AABB overlap with layer mask filtering.
//
// Parameters:
//
//	other: Other collider to test
//	thisTransform: This entity's transform
//	otherTransform: Other entity's transform
//
// Returns:
//
//	bool: True if colliders overlap AND layers are compatible
//
// Example:
//
//	if collider.Intersects(otherCollider, entity.Transform, other.Transform) {
//	    // Handle collision
//	}
func (c *Collider) Intersects(other *Collider, thisTransform, otherTransform gamemath.Transform) bool {
	// Check layer masks - must be on compatible layers
	thisLayerBit := 1 << c.CollisionLayer
	otherLayerBit := 1 << other.CollisionLayer

	// Check if this collider's mask includes other's layer
	// AND other collider's mask includes this layer
	if (c.CollisionMask&otherLayerBit) == 0 || (other.CollisionMask&thisLayerBit) == 0 {
		return false // Layers incompatible
	}

	// Get world bounds
	thisBounds := c.GetWorldBounds(thisTransform)
	otherBounds := other.GetWorldBounds(otherTransform)

	// AABB intersection test
	return thisBounds.Intersects(otherBounds)
}
