package input

import "github.com/veandco/go-sdl2/sdl"

// KeyCode represents a keyboard key or mouse button.
type KeyCode int

// Keyboard keys (wrapping SDL scancodes for type safety).
const (
	// Letters
	KeyA KeyCode = KeyCode(sdl.SCANCODE_A)
	KeyB KeyCode = KeyCode(sdl.SCANCODE_B)
	KeyC KeyCode = KeyCode(sdl.SCANCODE_C)
	KeyD KeyCode = KeyCode(sdl.SCANCODE_D)
	KeyE KeyCode = KeyCode(sdl.SCANCODE_E)
	KeyF KeyCode = KeyCode(sdl.SCANCODE_F)
	KeyG KeyCode = KeyCode(sdl.SCANCODE_G)
	KeyH KeyCode = KeyCode(sdl.SCANCODE_H)
	KeyI KeyCode = KeyCode(sdl.SCANCODE_I)
	KeyJ KeyCode = KeyCode(sdl.SCANCODE_J)
	KeyK KeyCode = KeyCode(sdl.SCANCODE_K)
	KeyL KeyCode = KeyCode(sdl.SCANCODE_L)
	KeyM KeyCode = KeyCode(sdl.SCANCODE_M)
	KeyN KeyCode = KeyCode(sdl.SCANCODE_N)
	KeyO KeyCode = KeyCode(sdl.SCANCODE_O)
	KeyP KeyCode = KeyCode(sdl.SCANCODE_P)
	KeyQ KeyCode = KeyCode(sdl.SCANCODE_Q)
	KeyR KeyCode = KeyCode(sdl.SCANCODE_R)
	KeyS KeyCode = KeyCode(sdl.SCANCODE_S)
	KeyT KeyCode = KeyCode(sdl.SCANCODE_T)
	KeyU KeyCode = KeyCode(sdl.SCANCODE_U)
	KeyV KeyCode = KeyCode(sdl.SCANCODE_V)
	KeyW KeyCode = KeyCode(sdl.SCANCODE_W)
	KeyX KeyCode = KeyCode(sdl.SCANCODE_X)
	KeyY KeyCode = KeyCode(sdl.SCANCODE_Y)
	KeyZ KeyCode = KeyCode(sdl.SCANCODE_Z)

	// Numbers
	Key0 KeyCode = KeyCode(sdl.SCANCODE_0)
	Key1 KeyCode = KeyCode(sdl.SCANCODE_1)
	Key2 KeyCode = KeyCode(sdl.SCANCODE_2)
	Key3 KeyCode = KeyCode(sdl.SCANCODE_3)
	Key4 KeyCode = KeyCode(sdl.SCANCODE_4)
	Key5 KeyCode = KeyCode(sdl.SCANCODE_5)
	Key6 KeyCode = KeyCode(sdl.SCANCODE_6)
	Key7 KeyCode = KeyCode(sdl.SCANCODE_7)
	Key8 KeyCode = KeyCode(sdl.SCANCODE_8)
	Key9 KeyCode = KeyCode(sdl.SCANCODE_9)

	// Arrow keys
	KeyArrowUp    KeyCode = KeyCode(sdl.SCANCODE_UP)
	KeyArrowDown  KeyCode = KeyCode(sdl.SCANCODE_DOWN)
	KeyArrowLeft  KeyCode = KeyCode(sdl.SCANCODE_LEFT)
	KeyArrowRight KeyCode = KeyCode(sdl.SCANCODE_RIGHT)

	// Special keys
	KeySpace  KeyCode = KeyCode(sdl.SCANCODE_SPACE)
	KeyEnter  KeyCode = KeyCode(sdl.SCANCODE_RETURN)
	KeyEscape KeyCode = KeyCode(sdl.SCANCODE_ESCAPE)
	KeyTab    KeyCode = KeyCode(sdl.SCANCODE_TAB)
	KeyShift  KeyCode = KeyCode(sdl.SCANCODE_LSHIFT)
	KeyCtrl   KeyCode = KeyCode(sdl.SCANCODE_LCTRL)
	KeyAlt    KeyCode = KeyCode(sdl.SCANCODE_LALT)

	// Mouse buttons (using high values to avoid conflicts with keyboard scancodes).
	KeyMouseLeft   KeyCode = 1000
	KeyMouseRight  KeyCode = 1001
	KeyMouseMiddle KeyCode = 1002
)
