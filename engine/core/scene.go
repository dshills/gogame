package core

import (
	"github.com/dshills/gogame/engine/graphics"
	gamemath "github.com/dshills/gogame/engine/math"
)

// Scene represents a container for entities (game level or screen)
type Scene struct {
	entities         []*Entity
	nextEntityID     uint64
	camera           *graphics.Camera
	backgroundColor  gamemath.Color
	entitiesToRemove []uint64 // Deferred removal during Update
}

// NewScene creates an empty scene
//
// Returns:
//
//	*Scene: New scene with no entities
//
// Example:
//
//	scene := core.NewScene()
func NewScene() *Scene {
	return &Scene{
		entities:         make([]*Entity, 0),
		nextEntityID:     1,
		camera:           graphics.NewCamera(),
		backgroundColor:  gamemath.Black,
		entitiesToRemove: make([]uint64, 0),
	}
}

// AddEntity adds an entity to the scene
//
// Parameters:
//
//	entity: Entity to add
//
// Returns:
//
//	uint64: Assigned entity ID (unique within scene)
//
// Behavior:
//   - Entity begins updating/rendering immediately if Active
//   - ID assigned sequentially starting from 1
//
// Example:
//
//	player := &core.Entity{Active: true}
//	playerID := scene.AddEntity(player)
func (s *Scene) AddEntity(entity *Entity) uint64 {
	entity.ID = s.nextEntityID
	s.nextEntityID++
	s.entities = append(s.entities, entity)
	return entity.ID
}

// RemoveEntity removes an entity by ID
//
// Parameters:
//
//	id: Entity ID to remove
//
// Behavior:
//   - Entity removed immediately (doesn't update/render next frame)
//   - Safe to call during Update() (deferred removal)
//   - No-op if ID not found
//
// Example:
//
//	scene.RemoveEntity(enemyID)
func (s *Scene) RemoveEntity(id uint64) {
	// Queue for deferred removal (safe during Update)
	s.entitiesToRemove = append(s.entitiesToRemove, id)
}

// processDeferredRemovals removes queued entities after update phase
func (s *Scene) processDeferredRemovals() {
	if len(s.entitiesToRemove) == 0 {
		return
	}

	// Create a map for O(1) lookup
	toRemove := make(map[uint64]bool)
	for _, id := range s.entitiesToRemove {
		toRemove[id] = true
	}

	// Filter out entities to remove
	filtered := make([]*Entity, 0, len(s.entities))
	for _, entity := range s.entities {
		if !toRemove[entity.ID] {
			filtered = append(filtered, entity)
		}
	}

	s.entities = filtered
	s.entitiesToRemove = s.entitiesToRemove[:0] // Clear removal queue
}

// GetEntity retrieves an entity by ID
//
// Parameters:
//
//	id: Entity ID to query
//
// Returns:
//
//	*Entity: Entity with matching ID, or nil if not found
//
// Example:
//
//	entity := scene.GetEntity(playerID)
//	if entity != nil {
//	    entity.Transform.Position.X += 10
//	}
func (s *Scene) GetEntity(id uint64) *Entity {
	for _, entity := range s.entities {
		if entity.ID == id {
			return entity
		}
	}
	return nil
}

// GetEntitiesAt finds all entities at a world position
//
// Parameters:
//
//	x, y: World coordinates
//
// Returns:
//
//	[]*Entity: Entities whose bounds contain the point (may be empty)
//
// Behavior:
//   - Returns entities in arbitrary order
//   - Empty slice if no matches
//
// Example:
//
//	mouseWorldX, mouseWorldY := camera.ScreenToWorld(mouseX, mouseY)
//	entities := scene.GetEntitiesAt(mouseWorldX, mouseWorldY)
func (s *Scene) GetEntitiesAt(x, y float64) []*Entity {
	result := make([]*Entity, 0)
	for _, entity := range s.entities {
		if entity.Active {
			bounds := entity.GetBounds()
			if bounds.Contains(x, y) {
				result = append(result, entity)
			}
		}
	}
	return result
}

// Camera returns the scene's camera
//
// Returns:
//
//	*graphics.Camera: Scene camera for view transform
func (s *Scene) Camera() *graphics.Camera {
	return s.camera
}

// SetBackgroundColor sets the clear color
//
// Parameters:
//
//	color: RGBA color to clear screen with
//
// Example:
//
//	scene.SetBackgroundColor(math.Color{R: 135, G: 206, B: 235, A: 255})  // Sky blue
func (s *Scene) SetBackgroundColor(color gamemath.Color) {
	s.backgroundColor = color
}

// GetBackgroundColor returns the current background color
func (s *Scene) GetBackgroundColor() gamemath.Color {
	return s.backgroundColor
}

// Update updates all active entities
func (s *Scene) Update(dt float64) {
	for _, entity := range s.entities {
		if entity.Active {
			entity.Update(dt)
		}
	}
	// Process any entities queued for removal during Update
	s.processDeferredRemovals()
}

// Render renders all active entities
func (s *Scene) Render(renderer *graphics.Renderer) error {
	// Sort entities by layer for correct draw order (lower layers first)
	// For now, we'll render in the order they were added (simple implementation)
	// TODO: Add layer sorting for proper z-ordering

	for _, entity := range s.entities {
		if entity.Active {
			if err := entity.Render(renderer, s.camera); err != nil {
				return err
			}
		}
	}
	return nil
}
