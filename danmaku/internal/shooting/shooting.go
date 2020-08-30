package shooting

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/yohamta/godanmaku/danmaku/internal/shooting/effects"

	"github.com/yohamta/godanmaku/danmaku/internal/ui"

	"github.com/hajimehoshi/ebiten"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/actors"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/fields"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/inputs"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/weapons"
)

const (
	maxPlayerShot = 80
	maxEnemyShot  = 70
	maxEnemy      = 50
	maxHitEffects = 30
	maxExplosions = 30
)

type gameState int

const (
	gameStateLoading gameState = iota
	gameStatePlaying
)

// PlayerShooter represents interface of Player Weapon
type PlayerShooter interface {
	Shot(x, y float64, degree int, playerShots []*actors.PlayerShot)
}

var (
	screenWidth  = 0
	screenHeight = 0

	input *inputs.Input
	field *fields.Field

	uiBackground      *ui.Box
	uiBackgroundColor = color.RGBA{0x00, 0x00, 0x00, 0xff}

	player       *actors.Player
	playerWeapon PlayerShooter

	playerShots [maxPlayerShot]*actors.PlayerShot
	enemyShots  [maxEnemyShot]*actors.EnemyShot
	enemies     [maxEnemy]*actors.Enemy
	hitEffects  [maxHitEffects]*effects.Hit
	explosions  [maxExplosions]*effects.Explosion

	state gameState = gameStateLoading
)

// Shooting represents shooting scene
type Shooting struct {
}

// NewShootingOptions represents options for New func
type NewShootingOptions struct {
	ScreenWidth  int
	ScreenHeight int
}

// NewShooting returns new Shooting struct
func NewShooting(options NewShootingOptions) *Shooting {
	stg := &Shooting{}

	screenWidth = options.ScreenWidth
	screenHeight = options.ScreenHeight

	state = gameStateLoading
	initGame()
	state = gameStatePlaying

	return stg
}

func initGame() {
	rand.Seed(time.Now().Unix())
	input = inputs.NewInput(screenWidth, screenHeight)
	field = fields.NewField()
	uiBackground = ui.NewBox(0, field.GetBottom(),
		screenWidth, screenHeight-(field.GetBottom()-field.GetTop()),
		uiBackgroundColor)

	actors.SetBoundary(field)

	// player
	player = actors.NewPlayer()
	playerWeapon = &weapons.PlayerWeapon1{}

	// enemies
	for i := 0; i < len(enemies); i++ {
		enemies[i] = actors.NewEnemy()
	}

	// shots
	for i := 0; i < len(playerShots); i++ {
		playerShots[i] = actors.NewPlayerShot()
	}

	// enemyShots
	for i := 0; i < len(enemyShots); i++ {
		enemyShots[i] = actors.NewEnemyShot()
	}

	// effects
	for i := 0; i < len(hitEffects); i++ {
		hitEffects[i] = effects.NewHit()
	}
	for i := 0; i < len(explosions); i++ {
		explosions[i] = effects.NewExplosion()
	}

	// Setup stage
	initEnemies()
}

// Update updates the scene
func (stg *Shooting) Update() {
	input.Update()

	checkCollision()

	// player
	if player.IsDead() == false {
		player.Move(input.Horizontal, input.Vertical, input.Fire)
		if input.Fire {
			x, y := player.GetPosition()
			playerWeapon.Shot(x, y, player.GetNormalizedDegree(), playerShots[:])
		}
	}

	// player shots
	for i := 0; i < len(playerShots); i++ {
		p := playerShots[i]
		if p.IsActive() == false {
			continue
		}
		p.Move()
	}

	// enemy shots
	for i := 0; i < len(enemyShots); i++ {
		e := enemyShots[i]
		if e.IsActive() == false {
			continue
		}
		e.Move()
	}

	// enemies
	for i := 0; i < len(enemies); i++ {
		e := enemies[i]
		if e.IsActive() == false {
			continue
		}
		e.Move(player)
		if e.ShouldAttack() {
			weapons.EnemyAttack(e, player, enemyShots[:])
		}
	}

	// hitEffects
	for i := 0; i < len(hitEffects); i++ {
		h := hitEffects[i]
		if h.IsActive() == false {
			continue
		}
		h.Update()
	}

	// explosions
	for i := 0; i < len(explosions); i++ {
		e := explosions[i]
		if e.IsActive() == false {
			continue
		}
		e.Update()
	}
}

// Draw draws the scene
func (stg *Shooting) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x30, 0xff})

	field.Draw(screen)

	// player shots
	for i := 0; i < len(playerShots); i++ {
		p := playerShots[i]
		if p.IsActive() == false {
			continue
		}
		p.Draw(screen)
	}

	// enemy shots
	for i := 0; i < len(enemyShots); i++ {
		e := enemyShots[i]
		if e.IsActive() == false {
			continue
		}
		e.Draw(screen)
	}

	// enemies
	for i := 0; i < len(enemies); i++ {
		e := enemies[i]
		if e.IsActive() == false {
			continue
		}
		e.Draw(screen)
	}

	if player.IsDead() == false {
		player.Draw(screen)
	}

	// explosions
	for i := 0; i < len(explosions); i++ {
		e := explosions[i]
		if e.IsActive() == false {
			continue
		}
		e.Draw(screen)
	}

	// hitEffects
	for i := 0; i < len(hitEffects); i++ {
		h := hitEffects[i]
		if h.IsActive() == false {
			continue
		}
		h.Draw(screen)
	}

	uiBackground.Draw(screen)
	input.Draw(screen)
}

func initEnemies() {
	enemyCount := 30

	for i := 0; i < enemyCount; i++ {
		enemy := enemies[i]
		enemy.Init(actors.EnemyKindBall)
	}
}

func checkCollision() {
	// player shots
	for i := 0; i < len(playerShots); i++ {
		p := playerShots[i]
		if p.IsActive() == false {
			continue
		}
		for j := 0; j < len(enemies); j++ {
			e := enemies[j]
			if e.IsActive() == false {
				continue
			}
			if actors.IsCollideWith(e, p) == false {
				continue
			}
			e.AddDamage(1)
			p.SetInactive()
			createHitEffect(p.GetX(), p.GetY())
			if e.IsDead() {
				createExplosion(e.GetX(), e.GetY())
			}
		}
	}

	// enemy shots
	if player.IsDead() == false {
		for i := 0; i < len(enemyShots); i++ {
			e := enemyShots[i]
			if e.IsActive() == false {
				continue
			}
			if actors.IsCollideWith(e, player) == false {
				continue
			}
			player.AddDamage(1)
			e.SetInactive()
			createHitEffect(player.GetX(), player.GetY())
			if player.IsDead() {
				createExplosion(player.GetX(), player.GetY())
			}
		}
	}
}

func createHitEffect(x, y int) {
	for i := 0; i < len(hitEffects); i++ {
		h := hitEffects[i]
		if h.IsActive() {
			continue
		}
		h.StartEffect(x, y)
		break
	}
}

func createExplosion(x, y int) {
	for i := 0; i < len(explosions); i++ {
		e := explosions[i]
		if e.IsActive() {
			continue
		}
		e.StartEffect(x, y)
		break
	}
}
