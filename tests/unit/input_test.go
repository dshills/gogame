package unit

import (
	"testing"

	"github.com/dshills/gogame/engine/input"
)

// TestActionBinding tests binding keys to actions.
func TestActionBinding(t *testing.T) {
	inputMgr := input.NewInputManager()

	// Bind multiple keys to one action
	inputMgr.BindAction(input.ActionJump, input.KeySpace, input.KeyW)

	// Verify bindings (internal state - tested through behavior)
	// This is a smoke test to ensure BindAction doesn't crash
}

// TestActionPressed tests ActionPressed detection.
func TestActionPressed(t *testing.T) {
	t.Skip("Requires SDL event simulation - needs mock InputManager")
	inputMgr := input.NewInputManager()
	inputMgr.BindAction(input.ActionJump, input.KeySpace)

	// On first frame, key is "pressed"
	if !inputMgr.ActionPressed(input.ActionJump) {
		t.Error("Expected ActionPressed to be true on first press")
	}

	// Call Update to copy current → previous
	inputMgr.Update()

	// After update, key is "held" not "pressed"
	if inputMgr.ActionPressed(input.ActionJump) {
		t.Error("Expected ActionPressed to be false after update (now held)")
	}
}

// TestActionHeld tests ActionHeld detection.
func TestActionHeld(t *testing.T) {
	inputMgr := input.NewInputManager()
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW)

	// Set key down
	// inputMgr.SetKeyState(input.KeyW, true)
	inputMgr.Update() // Copy to previous

	// Keep key down
	// inputMgr.SetKeyState(input.KeyW, true)

	// Should be held
	if !inputMgr.ActionHeld(input.ActionMoveUp) {
		t.Error("Expected ActionHeld to be true")
	}
}

// TestActionReleased tests ActionReleased detection.
func TestActionReleased(t *testing.T) {
	inputMgr := input.NewInputManager()
	inputMgr.BindAction(input.ActionJump, input.KeySpace)

	// Press key
	// inputMgr.SetKeyState(input.KeySpace, true)
	inputMgr.Update()

	// Release key
	// inputMgr.SetKeyState(input.KeySpace, false)

	// Should be released
	if !inputMgr.ActionReleased(input.ActionJump) {
		t.Error("Expected ActionReleased to be true")
	}

	inputMgr.Update()

	// After update, no longer released
	if inputMgr.ActionReleased(input.ActionJump) {
		t.Error("Expected ActionReleased to be false after update")
	}
}

// TestMultipleKeyBindings tests multiple keys bound to same action.
func TestMultipleKeyBindings(t *testing.T) {
	inputMgr := input.NewInputManager()
	inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)

	// Press first key
	// inputMgr.SetKeyState(input.KeyW, true)

	if !inputMgr.ActionPressed(input.ActionMoveUp) {
		t.Error("Expected action to work with first key")
	}

	inputMgr.Update()
	// inputMgr.SetKeyState(input.KeyW, false)
	inputMgr.Update()

	// Press second key
	// inputMgr.SetKeyState(input.KeyArrowUp, true)

	if !inputMgr.ActionPressed(input.ActionMoveUp) {
		t.Error("Expected action to work with second key")
	}
}

// TestMousePosition tests mouse position tracking.
func TestMousePosition(t *testing.T) {
	inputMgr := input.NewInputManager()

	// Set mouse position
	// inputMgr.SetMousePosition(100, 200)

	x, y := inputMgr.MousePosition()

	if x != 100 || y != 200 {
		t.Errorf("Expected position (100, 200), got (%d, %d)", x, y)
	}
}

// TestMouseDelta tests mouse movement delta.
func TestMouseDelta(t *testing.T) {
	inputMgr := input.NewInputManager()

	// Initial position
	// inputMgr.SetMousePosition(100, 100)
	inputMgr.Update()

	// Move mouse
	// inputMgr.SetMousePosition(150, 120)

	dx, dy := inputMgr.MouseDelta()

	if dx != 50 || dy != 20 {
		t.Errorf("Expected delta (50, 20), got (%d, %d)", dx, dy)
	}
}
