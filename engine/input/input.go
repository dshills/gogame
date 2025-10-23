package input

import "github.com/veandco/go-sdl2/sdl"

// InputManager manages keyboard and mouse input state with action mapping.
type InputManager struct {
	currentKeys  map[KeyCode]bool     // Current frame key state
	previousKeys map[KeyCode]bool     // Previous frame key state
	actionMap    map[Action][]KeyCode // Action to key bindings
	mouseX       int32                // Current mouse X position
	mouseY       int32                // Current mouse Y position
	prevMouseX   int32                // Previous mouse X position
	prevMouseY   int32                // Previous mouse Y position
}

// NewInputManager creates a new input manager.
//
// Returns:
//
//	*InputManager: New input manager with empty bindings
//
// Example:
//
//	input := input.NewInputManager()
//	input.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
func NewInputManager() *InputManager {
	return &InputManager{
		currentKeys:  make(map[KeyCode]bool),
		previousKeys: make(map[KeyCode]bool),
		actionMap:    make(map[Action][]KeyCode),
		mouseX:       0,
		mouseY:       0,
		prevMouseX:   0,
		prevMouseY:   0,
	}
}

// BindAction binds an action to one or more keys.
//
// Parameters:
//
//	action: Action to bind
//	keys: One or more keys that trigger this action
//
// Behavior:
//   - Replaces existing bindings for this action
//   - Multiple keys can trigger the same action
//
// Example:
//
//	input.BindAction(input.ActionJump, input.KeySpace)
//	input.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)
func (im *InputManager) BindAction(action Action, keys ...KeyCode) {
	im.actionMap[action] = keys
}

// ActionPressed returns true if action was just pressed this frame.
//
// Parameters:
//
//	action: Action to query
//
// Returns:
//
//	bool: True if any bound key went from up to down this frame
//
// Example:
//
//	if input.ActionPressed(input.ActionJump) {
//	    player.Jump()
//	}
func (im *InputManager) ActionPressed(action Action) bool {
	keys, exists := im.actionMap[action]
	if !exists {
		return false
	}

	for _, key := range keys {
		if im.currentKeys[key] && !im.previousKeys[key] {
			return true
		}
	}
	return false
}

// ActionReleased returns true if action was just released this frame.
//
// Parameters:
//
//	action: Action to query
//
// Returns:
//
//	bool: True if any bound key went from down to up this frame
//
// Example:
//
//	if input.ActionReleased(input.ActionAttack) {
//	    player.StopAttacking()
//	}
func (im *InputManager) ActionReleased(action Action) bool {
	keys, exists := im.actionMap[action]
	if !exists {
		return false
	}

	for _, key := range keys {
		if !im.currentKeys[key] && im.previousKeys[key] {
			return true
		}
	}
	return false
}

// ActionHeld returns true if action is currently being held.
//
// Parameters:
//
//	action: Action to query
//
// Returns:
//
//	bool: True if any bound key is currently down
//
// Example:
//
//	if input.ActionHeld(input.ActionMoveRight) {
//	    player.Transform.Position.X += speed * dt
//	}
func (im *InputManager) ActionHeld(action Action) bool {
	keys, exists := im.actionMap[action]
	if !exists {
		return false
	}

	for _, key := range keys {
		if im.currentKeys[key] {
			return true
		}
	}
	return false
}

// KeyPressed returns true if key was just pressed this frame.
//
// Parameters:
//
//	key: Key to query
//
// Returns:
//
//	bool: True if key went from up to down this frame
//
// Example:
//
//	if input.KeyPressed(input.KeyEscape) {
//	    game.Pause()
//	}
func (im *InputManager) KeyPressed(key KeyCode) bool {
	return im.currentKeys[key] && !im.previousKeys[key]
}

// KeyReleased returns true if key was just released this frame.
//
// Parameters:
//
//	key: Key to query
//
// Returns:
//
//	bool: True if key went from down to up this frame
func (im *InputManager) KeyReleased(key KeyCode) bool {
	return !im.currentKeys[key] && im.previousKeys[key]
}

// KeyHeld returns true if key is currently being held.
//
// Parameters:
//
//	key: Key to query
//
// Returns:
//
//	bool: True if key is currently down
func (im *InputManager) KeyHeld(key KeyCode) bool {
	return im.currentKeys[key]
}

// MousePosition returns the current mouse position.
//
// Returns:
//
//	x, y: Screen coordinates
//
// Example:
//
//	mouseX, mouseY := input.MousePosition()
//	worldX, worldY := camera.ScreenToWorld(mouseX, mouseY)
func (im *InputManager) MousePosition() (int32, int32) {
	return im.mouseX, im.mouseY
}

// MouseDelta returns mouse movement since last frame.
//
// Returns:
//
//	dx, dy: Movement in pixels
//
// Example:
//
//	dx, dy := input.MouseDelta()
//	camera.Position.X -= float64(dx) * sensitivity
func (im *InputManager) MouseDelta() (int32, int32) {
	return im.mouseX - im.prevMouseX, im.mouseY - im.prevMouseY
}

// Update swaps input buffers - call at end of frame.
//
// Behavior:
//   - Copies current state to previous state
//   - Should be called by Engine after update/render
//   - Do NOT call manually in game code
func (im *InputManager) Update() {
	// Copy current to previous
	im.previousKeys = make(map[KeyCode]bool, len(im.currentKeys))
	for key, state := range im.currentKeys {
		im.previousKeys[key] = state
	}

	// Update mouse delta tracking
	im.prevMouseX = im.mouseX
	im.prevMouseY = im.mouseY
}

// ProcessKeyEvent updates key state from SDL event.
func (im *InputManager) ProcessKeyEvent(event *sdl.KeyboardEvent) {
	scancode := KeyCode(event.Keysym.Scancode)
	im.currentKeys[scancode] = (event.State == sdl.PRESSED)
}

// ProcessMouseButtonEvent updates mouse button state from SDL event.
func (im *InputManager) ProcessMouseButtonEvent(event *sdl.MouseButtonEvent) {
	var key KeyCode
	switch event.Button {
	case sdl.BUTTON_LEFT:
		key = KeyMouseLeft
	case sdl.BUTTON_RIGHT:
		key = KeyMouseRight
	case sdl.BUTTON_MIDDLE:
		key = KeyMouseMiddle
	default:
		return
	}
	im.currentKeys[key] = (event.State == sdl.PRESSED)
}

// ProcessMouseMotionEvent updates mouse position from SDL event.
func (im *InputManager) ProcessMouseMotionEvent(event *sdl.MouseMotionEvent) {
	im.mouseX = event.X
	im.mouseY = event.Y
}
