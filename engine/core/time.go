package core

import "time"

// Time manages the game loop timing with fixed timestep
//
// Based on "Fix Your Timestep" pattern:
// - Fixed update rate (60 FPS = 16.67ms per update)
// - Variable render rate (as fast as possible with vsync)
// - Accumulator prevents spiral of death
type Time struct {
	dt           float64   // Fixed delta time in seconds (1/60 = 0.016667)
	accumulator  float64   // Time accumulated since last update
	lastTime     time.Time // Last frame timestamp
	targetFPS    float64   // Target updates per second (60.0)
	maxFrameTime float64   // Maximum frame time to prevent spiral of death (0.25 seconds)
}

// NewTime creates a new time manager with 60 FPS target
func NewTime() *Time {
	targetFPS := 60.0
	return &Time{
		dt:           1.0 / targetFPS,
		accumulator:  0.0,
		lastTime:     time.Now(),
		targetFPS:    targetFPS,
		maxFrameTime: 0.25, // Cap at 4 FPS minimum to prevent spiral of death
	}
}

// Tick advances the timer and returns how many fixed updates should run
//
// Returns:
//
//	int: Number of fixed updates to execute this frame (0-N)
//	float64: Fixed delta time for each update (always 1/60)
//
// Example:
//
//	updateCount, dt := time.Tick()
//	for i := 0; i < updateCount; i++ {
//	    scene.Update(dt)
//	}
func (t *Time) Tick() (updateCount int, dt float64) {
	now := time.Now()
	frameTime := now.Sub(t.lastTime).Seconds()
	t.lastTime = now

	// Clamp frame time to prevent spiral of death
	// (if game freezes for 5 seconds, don't try to catch up with 300 updates)
	if frameTime > t.maxFrameTime {
		frameTime = t.maxFrameTime
	}

	t.accumulator += frameTime

	// Consume accumulator in fixed timesteps
	updateCount = 0
	for t.accumulator >= t.dt {
		t.accumulator -= t.dt
		updateCount++
	}

	return updateCount, t.dt
}

// DeltaTime returns the fixed delta time in seconds
func (t *Time) DeltaTime() float64 {
	return t.dt
}

// FPS returns the target FPS
func (t *Time) FPS() float64 {
	return t.targetFPS
}
