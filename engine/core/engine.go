package core

import (
	"fmt"

	"github.com/dshills/gogame/engine/graphics"
	"github.com/veandco/go-sdl2/sdl"
)

// Engine is the root game engine managing window, rendering, and game loop
type Engine struct {
	window      *sdl.Window
	renderer    *graphics.Renderer
	scene       *Scene
	time        *Time
	running     bool
	width       int
	height      int
	assetMgr    *graphics.AssetManager
	initialized bool
}

// NewEngine creates a new game engine instance
//
// IMPORTANT: Must be called from the main OS thread. Call runtime.LockOSThread()
// in your main() function before calling NewEngine.
//
// Parameters:
//
//	title: Window title
//	width: Window width in pixels
//	height: Window height in pixels
//	fullscreen: Start in fullscreen mode (uses desktop resolution)
//
// Returns:
//
//	*Engine: Initialized engine
//	error: Non-nil if window/renderer creation fails
//
// Example:
//
//	import "runtime"
//
//	func main() {
//	    runtime.LockOSThread() // CRITICAL: SDL requires main thread
//	    engine, err := core.NewEngine("My Game", 800, 600, false)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer engine.Shutdown()
//	    // ...
//	}
func NewEngine(title string, width, height int, fullscreen bool) (*Engine, error) {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, fmt.Errorf("failed to initialize SDL: %w", err)
	}

	// Create window
	windowFlags := sdl.WINDOW_SHOWN
	if fullscreen {
		// Use desktop fullscreen for smoother mode switching
		windowFlags |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}

	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(width),
		int32(height),
		uint32(windowFlags),
	)
	if err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("failed to create SDL window: %w", err)
	}

	// Create hardware-accelerated renderer with vsync
	sdlRenderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		window.Destroy()
		sdl.Quit()
		return nil, fmt.Errorf("failed to create SDL renderer: %w", err)
	}

	// Wrap SDL renderer
	renderer := graphics.NewRenderer(sdlRenderer)

	// Create asset manager
	assetMgr := graphics.NewAssetManager(sdlRenderer)

	return &Engine{
		window:      window,
		renderer:    renderer,
		scene:       nil,
		time:        NewTime(),
		running:     false,
		width:       width,
		height:      height,
		assetMgr:    assetMgr,
		initialized: true,
	}, nil
}

// SetScene sets the active scene
//
// Parameters:
//
//	scene: Scene to activate
//
// Behavior:
//   - Previous scene (if any) is not destroyed (developer must manage)
//   - New scene begins updating/rendering immediately
//
// Example:
//
//	menuScene := core.NewScene()
//	engine.SetScene(menuScene)
func (e *Engine) SetScene(scene *Scene) {
	e.scene = scene
	// Update camera screen size
	if scene != nil && scene.camera != nil {
		scene.camera.SetScreenSize(e.width, e.height)
	}
}

// GetScene returns the currently active scene
//
// Returns:
//
//	*Scene: Active scene, or nil if none set
func (e *Engine) GetScene() *Scene {
	return e.scene
}

// Run starts the game loop (blocking)
//
// IMPORTANT: Must be called from the main OS thread. Call runtime.LockOSThread()
// in your main() function before creating the engine.
//
// Behavior:
//   - Runs until window closed or Stop() called
//   - Fixed 60 FPS update rate
//   - Variable rendering rate (vsync if enabled)
//   - Calls scene Update() and Render() each frame
//   - Returns error if rendering fails
//
// Example:
//
//	import "runtime"
//
//	func main() {
//	    runtime.LockOSThread() // Required for SDL
//	    engine, _ := core.NewEngine("Game", 800, 600, false)
//	    defer engine.Shutdown()
//	    if err := engine.Run(); err != nil {
//	        log.Fatal(err)
//	    }
//	}
func (e *Engine) Run() error {
	if !e.initialized {
		return nil
	}

	e.running = true
	defer func() { e.running = false }()

	const maxUpdateSteps = 8 // Prevent spiral of death

	for e.running {
		// Handle SDL events
		if !e.handleEvents() {
			break
		}

		// Prevent busy loop when no scene is active
		if e.scene == nil {
			sdl.Delay(1) // Sleep 1ms to avoid maxing CPU
			continue
		}

		// Update with fixed timestep (capped to prevent spiral of death)
		updateCount, dt := e.time.Tick()
		if updateCount > maxUpdateSteps {
			updateCount = maxUpdateSteps
		}

		for i := 0; i < updateCount; i++ {
			e.scene.Update(dt)
		}

		// Render
		// Clear screen with background color
		bgColor := e.scene.GetBackgroundColor()
		if err := e.renderer.Clear(bgColor); err != nil {
			return fmt.Errorf("failed to clear screen: %w", err)
		}

		// Render scene
		if err := e.scene.Render(e.renderer); err != nil {
			return fmt.Errorf("failed to render scene: %w", err)
		}

		// Present frame
		e.renderer.Present()
	}

	return nil
}

// handleEvents processes SDL events and returns false if should quit
func (e *Engine) handleEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch evt := event.(type) {
		case *sdl.QuitEvent:
			return false

		case *sdl.WindowEvent:
			if evt.Event == sdl.WINDOWEVENT_RESIZED {
				e.width = int(evt.Data1)
				e.height = int(evt.Data2)
				// Update camera dimensions
				if e.scene != nil && e.scene.camera != nil {
					e.scene.camera.SetScreenSize(e.width, e.height)
				}
			}
		}
	}
	return true
}

// Stop signals the game loop to exit
//
// Behavior:
//   - Game loop exits after current frame completes
//   - Resources remain allocated (call Shutdown to cleanup)
//
// Example:
//
//	engine.Stop()
func (e *Engine) Stop() {
	e.running = false
}

// Shutdown releases all engine resources
//
// Behavior:
//   - Destroys SDL window and renderer
//   - Unloads all textures
//   - Must be called before program exit
//   - Engine unusable after this call
//
// Example:
//
//	defer engine.Shutdown()
func (e *Engine) Shutdown() {
	if !e.initialized {
		return
	}

	// Destroy asset manager (unloads all textures)
	if e.assetMgr != nil {
		e.assetMgr.Destroy()
	}

	// Destroy renderer
	if e.renderer != nil {
		e.renderer.Destroy()
	}

	// Destroy window
	if e.window != nil {
		e.window.Destroy()
	}

	// Quit SDL
	sdl.Quit()

	e.initialized = false
}

// Assets returns the asset manager
//
// Returns:
//
//	*graphics.AssetManager: Asset loading subsystem
func (e *Engine) Assets() *graphics.AssetManager {
	return e.assetMgr
}

// Width returns the window width
func (e *Engine) Width() int {
	return e.width
}

// Height returns the window height
func (e *Engine) Height() int {
	return e.height
}
