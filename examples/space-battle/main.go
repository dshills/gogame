// Package main implements Space Battle - a complete demo game showcasing the gogame engine
package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/dshills/gogame/engine/core"
	"github.com/dshills/gogame/engine/graphics"
	"github.com/dshills/gogame/engine/input"
	gamemath "github.com/dshills/gogame/engine/math"
	"github.com/dshills/gogame/engine/physics"
)

// Game constants
const (
	ScreenWidth  = 800
	ScreenHeight = 600

	PlayerSpeed   = 300.0
	BulletSpeed   = 400.0
	EnemySpeed    = 100.0
	ShootCooldown = 0.25 // Seconds between shots

	EnemySpawnInterval = 1.5  // Seconds between enemy spawns
	StarSpawnInterval  = 0.1  // Seconds between star spawns
	MaxStars           = 50   // Maximum number of background stars
	StarSpeed          = 50.0 // Pixels per second

	CollisionLayerPlayer = 1
	CollisionLayerEnemy  = 2
	CollisionLayerBullet = 4
)

// Game states
type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
)

// Game manages the overall game state
type Game struct {
	engine              *core.Engine
	scene               *core.Scene
	inputMgr            *input.InputManager
	state               GameState
	score               int
	player              *core.Entity
	enemies             []*core.Entity
	bullets             []*core.Entity
	stars               []*core.Entity
	lastShot            float64
	enemySpawnTimer     float64
	starSpawnTimer      float64
	gameTime            float64
	playerTexture       *graphics.Texture
	enemyTexture        *graphics.Texture
	bulletTexture       *graphics.Texture
	starTexture         *graphics.Texture
	playerStartPosition gamemath.Vector2
}

// PlayerController handles player movement and shooting
type PlayerController struct {
	game *Game
}

func (pc *PlayerController) Update(entity *core.Entity, dt float64) {
	if pc.game.state != StatePlaying {
		return
	}

	inputMgr := pc.game.inputMgr

	// Movement
	moveSpeed := PlayerSpeed * dt
	if inputMgr.ActionHeld(input.ActionMoveLeft) {
		entity.Transform.Position.X -= moveSpeed
	}
	if inputMgr.ActionHeld(input.ActionMoveRight) {
		entity.Transform.Position.X += moveSpeed
	}
	if inputMgr.ActionHeld(input.ActionMoveUp) {
		entity.Transform.Position.Y -= moveSpeed
	}
	if inputMgr.ActionHeld(input.ActionMoveDown) {
		entity.Transform.Position.Y += moveSpeed
	}

	// Constrain to screen bounds
	if entity.Transform.Position.X < 50 {
		entity.Transform.Position.X = 50
	}
	if entity.Transform.Position.X > ScreenWidth-50 {
		entity.Transform.Position.X = ScreenWidth - 50
	}
	if entity.Transform.Position.Y < 50 {
		entity.Transform.Position.Y = 50
	}
	if entity.Transform.Position.Y > ScreenHeight-50 {
		entity.Transform.Position.Y = ScreenHeight - 50
	}

	// Shooting
	if inputMgr.KeyHeld(input.KeySpace) {
		pc.game.tryShoot()
	}
}

// EnemyBehavior moves enemies downward
type EnemyBehavior struct {
	game *Game
}

func (eb *EnemyBehavior) Update(entity *core.Entity, dt float64) {
	// Move down
	entity.Transform.Position.Y += EnemySpeed * dt

	// Remove if off screen
	if entity.Transform.Position.Y > ScreenHeight+50 {
		eb.game.removeEnemy(entity)
	}
}

// BulletBehavior moves bullets upward
type BulletBehavior struct {
	game *Game
}

func (bb *BulletBehavior) Update(entity *core.Entity, dt float64) {
	// Move up
	entity.Transform.Position.Y -= BulletSpeed * dt

	// Remove if off screen
	if entity.Transform.Position.Y < -50 {
		bb.game.removeBullet(entity)
	}
}

// StarBehavior moves stars downward for parallax effect
type StarBehavior struct {
	game *Game
}

func (sb *StarBehavior) Update(entity *core.Entity, dt float64) {
	// Move down slowly
	entity.Transform.Position.Y += StarSpeed * dt

	// Wrap around when off screen
	if entity.Transform.Position.Y > ScreenHeight+10 {
		entity.Transform.Position.Y = -10
		entity.Transform.Position.X = rand.Float64() * ScreenWidth
	}
}

// NewGame creates a new game instance
func NewGame(engine *core.Engine) *Game {
	return &Game{
		engine:              engine,
		inputMgr:            engine.Input(),
		state:               StatePlaying,
		score:               0,
		enemies:             make([]*core.Entity, 0),
		bullets:             make([]*core.Entity, 0),
		stars:               make([]*core.Entity, 0),
		lastShot:            0,
		enemySpawnTimer:     0,
		starSpawnTimer:      0,
		gameTime:            0,
		playerStartPosition: gamemath.Vector2{X: ScreenWidth / 2, Y: ScreenHeight - 100},
	}
}

// Initialize sets up the game
func (g *Game) Initialize() error {
	log.Println("╔═══════════════════════════════════════════════════════════╗")
	log.Println("║                    SPACE BATTLE                           ║")
	log.Println("╚═══════════════════════════════════════════════════════════╝")
	log.Println()
	log.Println("Controls:")
	log.Println("  WASD / Arrow Keys - Move")
	log.Println("  SPACE - Shoot")
	log.Println("  R - Restart (when game over)")
	log.Println("  ESC - Quit")
	log.Println()
	log.Println("Objective: Destroy enemies and survive!")
	log.Println()

	// Setup input bindings
	g.inputMgr.BindAction(input.ActionMoveUp, input.KeyW, input.KeyArrowUp)
	g.inputMgr.BindAction(input.ActionMoveDown, input.KeyS, input.KeyArrowDown)
	g.inputMgr.BindAction(input.ActionMoveLeft, input.KeyA, input.KeyArrowLeft)
	g.inputMgr.BindAction(input.ActionMoveRight, input.KeyD, input.KeyArrowRight)

	// Create scene
	g.scene = core.NewScene()
	g.scene.SetBackgroundColor(gamemath.Color{R: 10, G: 10, B: 30, A: 255}) // Dark space blue

	// Position camera at screen center
	camera := g.scene.Camera()
	camera.Position = gamemath.Vector2{X: ScreenWidth / 2, Y: ScreenHeight / 2}

	g.engine.SetScene(g.scene)

	// Load textures
	assets := g.engine.Assets()

	var err error
	g.playerTexture, err = assets.LoadTexture("examples/space-battle/assets/player.png")
	if err != nil {
		return fmt.Errorf("failed to load player texture: %v", err)
	}

	g.enemyTexture, err = assets.LoadTexture("examples/space-battle/assets/enemy.png")
	if err != nil {
		return fmt.Errorf("failed to load enemy texture: %v", err)
	}

	g.bulletTexture, err = assets.LoadTexture("examples/space-battle/assets/bullet.png")
	if err != nil {
		return fmt.Errorf("failed to load bullet texture: %v", err)
	}

	g.starTexture, err = assets.LoadTexture("examples/space-battle/assets/star.png")
	if err != nil {
		return fmt.Errorf("failed to load star texture: %v", err)
	}

	// Create game manager entity (invisible, just runs game logic)
	gameManager := &core.Entity{
		Active:   true,
		Behavior: &GameManagerBehavior{game: g},
		Layer:    0,
	}
	g.scene.AddEntity(gameManager)

	// Create player
	g.createPlayer()

	// Spawn initial stars
	for i := 0; i < MaxStars; i++ {
		g.spawnStar(rand.Float64() * ScreenHeight)
	}

	log.Println("Game initialized! Good luck!")
	log.Println()

	return nil
}

// createPlayer creates the player entity
func (g *Game) createPlayer() {
	sprite := graphics.NewSprite(g.playerTexture)
	sprite.SetColor(gamemath.Color{R: 100, G: 200, B: 255, A: 255})

	g.player = &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: g.playerStartPosition,
			Scale:    gamemath.Vector2{X: 2, Y: 2},
		},
		Sprite:   sprite,
		Collider: physics.NewCollider(32, 32),
		Behavior: &PlayerController{game: g},
		Layer:    2,
	}
	g.player.Collider.CollisionLayer = CollisionLayerPlayer
	g.player.Collider.CollisionMask = CollisionLayerEnemy // Collide with enemies only

	// Collision callbacks
	g.player.OnCollisionEnter = func(self, other *core.Entity) {
		// Check if collided with enemy
		if other.Collider != nil && other.Collider.CollisionLayer == CollisionLayerEnemy {
			g.onPlayerHit()
		}
	}

	g.scene.AddEntity(g.player)
}

// tryShoot attempts to shoot a bullet
func (g *Game) tryShoot() {
	if g.gameTime-g.lastShot < ShootCooldown {
		return
	}

	g.lastShot = g.gameTime

	// Create bullet at player position
	sprite := graphics.NewSprite(g.bulletTexture)
	sprite.SetColor(gamemath.Color{R: 255, G: 255, B: 150, A: 255})

	bullet := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{
				X: g.player.Transform.Position.X,
				Y: g.player.Transform.Position.Y - 30,
			},
			Scale: gamemath.Vector2{X: 1.5, Y: 1.5},
		},
		Sprite:   sprite,
		Collider: physics.NewCollider(8, 16),
		Behavior: &BulletBehavior{game: g},
		Layer:    2,
	}
	bullet.Collider.CollisionLayer = CollisionLayerBullet
	bullet.Collider.CollisionMask = CollisionLayerEnemy // Collide with enemies only

	// Collision callback
	bullet.OnCollisionEnter = func(self, other *core.Entity) {
		// Check if hit enemy
		if other.Collider != nil && other.Collider.CollisionLayer == CollisionLayerEnemy {
			g.onEnemyHit(other)
			g.removeBullet(self)
		}
	}

	g.bullets = append(g.bullets, bullet)
	g.scene.AddEntity(bullet)
}

// spawnEnemy creates a new enemy at a random position at the top
func (g *Game) spawnEnemy() {
	x := 50 + rand.Float64()*(ScreenWidth-100)

	sprite := graphics.NewSprite(g.enemyTexture)
	sprite.SetColor(gamemath.Color{R: 255, G: 100, B: 100, A: 255})

	enemy := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: x, Y: -30},
			Scale:    gamemath.Vector2{X: 2, Y: 2},
		},
		Sprite:   sprite,
		Collider: physics.NewCollider(32, 32),
		Behavior: &EnemyBehavior{game: g},
		Layer:    2,
	}
	enemy.Collider.CollisionLayer = CollisionLayerEnemy
	enemy.Collider.CollisionMask = CollisionLayerPlayer | CollisionLayerBullet // Collide with player and bullets

	g.enemies = append(g.enemies, enemy)
	g.scene.AddEntity(enemy)
}

// spawnStar creates a background star
func (g *Game) spawnStar(y float64) {
	if len(g.stars) >= MaxStars {
		return
	}

	x := rand.Float64() * ScreenWidth

	sprite := graphics.NewSprite(g.starTexture)
	sprite.Alpha = 0.3 + rand.Float64()*0.4 // Random alpha for depth effect

	star := &core.Entity{
		Active: true,
		Transform: gamemath.Transform{
			Position: gamemath.Vector2{X: x, Y: y},
			Scale:    gamemath.Vector2{X: 1, Y: 1},
		},
		Sprite:   sprite,
		Behavior: &StarBehavior{game: g},
		Layer:    0, // Background layer
	}

	g.stars = append(g.stars, star)
	g.scene.AddEntity(star)
}

// onEnemyHit is called when an enemy is hit by a bullet
func (g *Game) onEnemyHit(enemy *core.Entity) {
	g.score += 10

	// Visual feedback - flash white
	if enemy.Sprite != nil {
		enemy.Sprite.SetColor(gamemath.Color{R: 255, G: 255, B: 255, A: 255})
	}

	g.removeEnemy(enemy)
	log.Printf("Enemy destroyed! Score: %d", g.score)
}

// onPlayerHit is called when player is hit by an enemy
func (g *Game) onPlayerHit() {
	log.Println()
	log.Println("═══════════════════════════════════")
	log.Println("         GAME OVER!")
	log.Printf("       Final Score: %d", g.score)
	log.Println("═══════════════════════════════════")
	log.Println("Press R to restart or ESC to quit")
	log.Println()

	g.state = StateGameOver

	// Change player color to red
	if g.player.Sprite != nil {
		g.player.Sprite.SetColor(gamemath.Color{R: 255, G: 50, B: 50, A: 255})
	}
}

// removeEnemy removes an enemy from the game
func (g *Game) removeEnemy(enemy *core.Entity) {
	g.scene.RemoveEntity(enemy.ID)

	// Remove from enemies list
	for i, e := range g.enemies {
		if e.ID == enemy.ID {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			break
		}
	}
}

// removeBullet removes a bullet from the game
func (g *Game) removeBullet(bullet *core.Entity) {
	g.scene.RemoveEntity(bullet.ID)

	// Remove from bullets list
	for i, b := range g.bullets {
		if b.ID == bullet.ID {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			break
		}
	}
}

// GameManagerBehavior handles global game logic
type GameManagerBehavior struct {
	game *Game
}

func (gmb *GameManagerBehavior) Update(entity *core.Entity, dt float64) {
	g := gmb.game
	g.gameTime += dt

	// Check for restart
	if g.state == StateGameOver {
		if g.inputMgr.KeyPressed(input.KeyR) {
			g.restart()
		}
		return
	}

	// Check for quit
	if g.inputMgr.KeyPressed(input.KeyEscape) {
		g.engine.Stop()
		return
	}

	// Spawn enemies
	g.enemySpawnTimer += dt
	if g.enemySpawnTimer >= EnemySpawnInterval {
		g.enemySpawnTimer = 0
		g.spawnEnemy()
	}

	// Spawn stars
	g.starSpawnTimer += dt
	if g.starSpawnTimer >= StarSpawnInterval {
		g.starSpawnTimer = 0
		if len(g.stars) < MaxStars {
			g.spawnStar(-10)
		}
	}
}

// restart restarts the game
func (g *Game) restart() {
	log.Println()
	log.Println("Restarting game...")
	log.Println()

	// Clear all entities except stars
	for _, enemy := range g.enemies {
		g.scene.RemoveEntity(enemy.ID)
	}
	for _, bullet := range g.bullets {
		g.scene.RemoveEntity(bullet.ID)
	}

	g.enemies = make([]*core.Entity, 0)
	g.bullets = make([]*core.Entity, 0)

	// Reset player
	g.scene.RemoveEntity(g.player.ID)
	g.createPlayer()

	// Reset game state
	g.state = StatePlaying
	g.score = 0
	g.gameTime = 0
	g.lastShot = 0
	g.enemySpawnTimer = 0

	log.Println("Game restarted! Good luck!")
}

func main() {
	// CRITICAL: SDL requires running on the main OS thread
	runtime.LockOSThread()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create engine
	engine, err := core.NewEngine("Space Battle - gogame Demo", ScreenWidth, ScreenHeight, false)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	// Create game
	game := NewGame(engine)
	if err := game.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Run game loop (game manager behavior handles updates)
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println()
	log.Println("Thanks for playing Space Battle!")
	log.Printf("Final Score: %d", game.score)
}
