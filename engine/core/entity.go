package core

import (
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// Behavior defines custom per-frame logic for an entity.
type Behavior interface {
	// Update is called every frame
	//
	// Parameters:
	//   entity: The entity this behavior is attached to
	//   dt: Delta time in seconds (typically 0.016 at 60 FPS)
	//
	// Example:
	//   func (pc *PlayerController) Update(entity *Entity, dt float64) {
	//       if input.ActionHeld(ActionMoveRight) {
	//           entity.Transform.Position.X += pc.Speed * dt
	//       }
	//   }
	Update(entity *Entity, dt float64)
}

// CollisionCallback is called when collision events occur.
// Parameters:
//   - self: The entity this callback is attached to
//   - other: The entity we collided with
type CollisionCallback func(self, other *Entity)

// Entity represents a game object with position, optional visuals, and behavior.
type Entity struct {
	ID        uint64             // Unique identifier (assigned by Scene)
	Active    bool               // Update/render only if true
	Transform gamemath.Transform // Position, rotation, scale (required)
	Sprite    *graphics.Sprite   // Optional visual representation
	Collider  *physics.Collider  // Optional collision detection
	Behavior  Behavior           // Optional custom update logic
	Layer     int                // Z-order (higher renders on top)

	// Collision callbacks (optional)
	OnCollisionEnter CollisionCallback // Called when collision starts
	OnCollisionStay  CollisionCallback // Called while collision continues
	OnCollisionExit  CollisionCallback // Called when collision ends
}

// Update updates the entity's transform and behavior
//
// Parameters:
//
//	dt: Delta time in seconds
//
// Behavior:
//   - Calls Behavior.Update() if non-nil
//   - Called automatically by Scene during update phase
//
// Example:
//
//	// Typically called by engine, not user code
//	entity.Update(0.016)  // 16ms frame
func (e *Entity) Update(dt float64) {
	if e.Behavior != nil {
		e.Behavior.Update(e, dt)
	}
}

// Render draws the entity's sprite
//
// Parameters:
//
//	renderer: Renderer
//	camera: Camera for view transform
//
// Behavior:
//   - Renders Sprite if non-nil
//   - Applies transform (position, rotation, scale)
//   - Called automatically by Scene during render phase
//
// Example:
//
//	// Typically called by engine, not user code
//	entity.Render(renderer, camera)
func (e *Entity) Render(renderer *graphics.Renderer, camera *graphics.Camera) error {
	if e.Sprite != nil {
		return renderer.DrawSprite(e.Sprite, e.Transform, camera)
	}
	return nil
}

// GetBounds returns world-space bounding box
//
// Returns:
//
//	math.Rectangle: World-space bounds (or zero rect if no collider)
//
// Example:
//
//	bounds := entity.GetBounds()
//	if bounds.Contains(clickX, clickY) {
//	    fmt.Println("Entity clicked!")
//	}
func (e *Entity) GetBounds() gamemath.Rectangle {
	if e.Collider != nil {
		return e.Collider.GetWorldBounds(e.Transform)
	}

	// No collider - return zero-size rectangle at entity position
	return gamemath.Rectangle{
		X:      e.Transform.Position.X,
		Y:      e.Transform.Position.Y,
		Width:  0,
		Height: 0,
	}
}

// GetID returns the entity's unique identifier.
func (e *Entity) GetID() uint64 {
	return e.ID
}

// GetTransform returns the entity's transform.
func (e *Entity) GetTransform() gamemath.Transform {
	return e.Transform
}

// GetCollider returns the entity's collider.
func (e *Entity) GetCollider() *physics.Collider {
	return e.Collider
}

// IsActive returns whether the entity is active.
func (e *Entity) IsActive() bool {
	return e.Active
}
