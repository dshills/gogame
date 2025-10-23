// Package input provides keyboard and mouse input handling with action mapping.
package input

// Action represents a game action that can be bound to multiple keys.
type Action int

// Common game actions (users can define their own).
const (
	ActionNone Action = iota

	// Movement
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight

	// Player actions
	ActionJump
	ActionAttack
	ActionInteract
	ActionPause

	// UI actions
	ActionConfirm
	ActionCancel
	ActionMenu
)
