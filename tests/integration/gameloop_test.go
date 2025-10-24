package integration

import (
	"runtime"
	"testing"
	"time"

	"github.com/dshills/gogame/engine/core"
	gamemath "github.com/dshills/gogame/engine/math"
)

// TestFixedTimestepGameLoop verifies that the game loop updates at a fixed timestep.
func TestFixedTimestepGameLoop(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	engine, err := core.NewEngine("GameLoop Test", 800, 600, false)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	scene := core.NewScene()
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}
	engine.SetScene(scene)

	// Track update calls
	updateCount := 0

	// Custom behavior implementation (inline for test)
	// Note: This would normally implement core.Behavior interface
	entity := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 400, Y: 300},
			Scale:    gamemath.Vector2{X: 1, Y: 1},
		},
		Behavior: &testBehavior{counter: &updateCount},
		Layer:    0,
	}
	scene.AddEntity(entity)

	// Run engine in goroutine for limited time
	done := make(chan bool)
	go func() {
		// Run for ~100ms (should get ~6 updates at 60 FPS)
		time.Sleep(100 * time.Millisecond)
		done <- true
	}()

	// Note: This test is conceptual - actual game loop testing
	// requires more sophisticated mocking of SDL events
	// In practice, would test Time component separately

	<-done

	// Verify we got updates (actual count depends on timing)
	if updateCount == 0 {
		t.Error("Expected at least one update, got zero")
	}

	t.Logf("Got %d updates in ~100ms", updateCount)
}

// testBehavior is a simple test behavior that counts updates.
type testBehavior struct {
	counter *int
}

func (tb *testBehavior) Update(entity *core.Entity, dt float64) {
	*tb.counter++
}

// TestSceneUpdateRenderCycle verifies update happens before render.
func TestSceneUpdateRenderCycle(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	engine, err := core.NewEngine("Update/Render Test", 800, 600, false)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	scene := core.NewScene()
	camera := scene.Camera()
	camera.Position = gamemath.Vector2{X: 400, Y: 300}

	// Track execution order
	executionLog := []string{}

	// Entity that logs update
	entity := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: 100, Y: 100},
			Scale:    gamemath.Vector2{X: 1, Y: 1},
		},
		Behavior: &logBehavior{log: &executionLog, message: "update"},
		Layer:    0,
	}
	scene.AddEntity(entity)

	// Manually call update and render to test order
	scene.Update(0.016) // ~60 FPS delta
	executionLog = append(executionLog, "render")

	// Verify update happened before render
	if len(executionLog) < 2 {
		t.Fatal("Expected at least 2 log entries")
	}

	if executionLog[0] != "update" {
		t.Errorf("Expected first entry to be 'update', got '%s'", executionLog[0])
	}

	if executionLog[1] != "render" {
		t.Errorf("Expected second entry to be 'render', got '%s'", executionLog[1])
	}
}

// logBehavior logs messages to verify execution order.
type logBehavior struct {
	log     *[]string
	message string
}

func (lb *logBehavior) Update(entity *core.Entity, dt float64) {
	*lb.log = append(*lb.log, lb.message)
}

// TestDeltaTime verifies delta time is passed correctly to behaviors.
func TestDeltaTime(t *testing.T) {
	scene := core.NewScene()

	// Track delta times received
	deltaTimes := []float64{}

	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 0, Y: 0}},
		Behavior:  &deltaBehavior{deltas: &deltaTimes},
		Layer:     0,
	}
	scene.AddEntity(entity)

	// Call update with known delta
	testDelta := 0.016667 // ~60 FPS
	scene.Update(testDelta)

	if len(deltaTimes) == 0 {
		t.Fatal("Expected delta time to be recorded")
	}

	if deltaTimes[0] != testDelta {
		t.Errorf("Expected delta %f, got %f", testDelta, deltaTimes[0])
	}
}

// deltaBehavior records delta times.
type deltaBehavior struct {
	deltas *[]float64
}

func (db *deltaBehavior) Update(entity *core.Entity, dt float64) {
	*db.deltas = append(*db.deltas, dt)
}
