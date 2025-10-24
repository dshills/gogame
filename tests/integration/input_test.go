package integration

import (
	"runtime"
	"testing"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/input"
	gamemath "github.com/dshills/gogame/engine/math"
)

// TestInputInGameLoop tests input handling within game loop context.
func TestInputInGameLoop(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	engine, err := core.NewEngine("Input Test", 800, 600, false)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Shutdown()

	scene := core.NewScene()
	inputMgr := engine.Input()

	// Bind actions
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
	inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
	inputMgr.BindAction(input.ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
	inputMgr.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)

	// Create entity with input-driven behavior
	moved := false
	entity := &core.Entity{
		Active:    true,
		Transform: gamemath.Transform{Position: gamemath.Vector2{X: 400, Y: 300}},
		Behavior:  &inputTestBehavior{inputMgr: inputMgr, moved: &moved},
		Layer:     0,
	}
	scene.AddEntity(entity)

	// Simulate key press
	// inputMgr.SetKeyState(input.KeyW, true)
	scene.Update(0.016)

	// Verify entity responded to input
	if !moved {
		t.Error("Expected entity to respond to input")
	}
}

// inputTestBehavior is a test behavior that tracks input responses.
type inputTestBehavior struct {
	inputMgr *input.InputManager
	moved    *bool
}

func (itb *inputTestBehavior) Update(entity *core.Entity, dt float64) {
	if itb.inputMgr.ActionHeld(input.ActionMoveUp) {
		entity.Transform.Position.Y -= 100 * dt
		*itb.moved = true
	}
}

// TestSimultaneousInputs tests handling multiple keys at once.
func TestSimultaneousInputs(t *testing.T) {
	inputMgr := input.NewInputManager()
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW)
	inputMgr.BindAction(input.ActionMoveRight, input.KeyD)

	// Press both keys
	// inputMgr.SetKeyState(input.KeyW, true)
	// inputMgr.SetKeyState(input.KeyD, true)

	// Both should be detected
	if !inputMgr.ActionPressed(input.ActionMoveUp) {
		t.Error("Expected ActionMoveUp to be pressed")
	}
	if !inputMgr.ActionPressed(input.ActionMoveRight) {
		t.Error("Expected ActionMoveRight to be pressed")
	}
}
