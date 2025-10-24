package benchmarks

import (
	"testing"

	"github.com/dshills/gogame/engine/core"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// BenchmarkCollisionDetection50Entities benchmarks collision detection with 50 entities.
// Target: <16ms per frame
func BenchmarkCollisionDetection50Entities(b *testing.B) {
	scene := core.NewScene()

	// Create 50 entities with colliders
	for i := 0; i < 50; i++ {
		collider := physics.NewCollider(32, 32)
		collider.CollisionLayer = 0
		collider.CollisionMask = 0xFF // Collide with all

		entity := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{
					X: float64((i % 10) * 50),
					Y: float64((i / 10) * 50),
				},
			},
			Collider: collider,
			Layer:    0,
		}
		scene.AddEntity(entity)
	}

	b.ResetTimer()

	// Benchmark collision detection
	for i := 0; i < b.N; i++ {
		scene.Update(0.016) // Includes collision detection
	}
}

// BenchmarkCollisionDetection100Entities tests with more entities.
func BenchmarkCollisionDetection100Entities(b *testing.B) {
	scene := core.NewScene()

	for i := 0; i < 100; i++ {
		collider := physics.NewCollider(32, 32)
		collider.CollisionLayer = 0
		collider.CollisionMask = 0xFF

		entity := &core.Entity{
			Active: true,
			Transform: gamemath.Transform{
				Position: gamemath.Vector2{
					X: float64((i % 10) * 40),
					Y: float64((i / 10) * 40),
				},
			},
			Collider: collider,
			Layer:    0,
		}
		scene.AddEntity(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		scene.Update(0.016)
	}
}
