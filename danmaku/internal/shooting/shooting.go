package shooting

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/yohamta/godanmaku/danmaku/internal/ui"

	"github.com/hajimehoshi/ebiten"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/actors"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/fields"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/inputs"
	"github.com/yohamta/godanmaku/danmaku/internal/shooting/weapons"
)

const (
	maxPlayerBulletsNum = 80
	maxEnemyNum         = 50
)

type gameState int

const (
	gameStateLoading gameState = iota
	gameStatePlaying
)

// PlayerWeapon represents interface of Player Weapon
type PlayerWeapon interface {
	Shot(x, y float64, degree int, playerShots []*actors.PlayerBullet)
}

var (
	screenWidth  = 0
	screenHeight = 0

	input *inputs.Input
	field *fields.Field

	uiBackground      *ui.Box
	uiBackgroundColor = color.RGBA{0x20, 0x20, 0x40, 0xff}

	player       *actors.Player
	playerWeapon PlayerWeapon

	playerShots [maxPlayerBulletsNum]*actors.PlayerBullet
	enemies     [maxEnemyNum]*actors.Enemy

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
	player = actors.NewPlayer()
	for i := 0; i < maxPlayerBulletsNum; i++ {
		playerShots[i] = actors.NewPlayerShot()
	}

	for i := 0; i < maxEnemyNum; i++ {
		enemies[i] = actors.NewEnemy()
	}

	playerWeapon = &weapons.PlayerWeapon1{}

	initEnemies()
}

// Update updates the scene
func (stg *Shooting) Update() {
	input.Update()

	// player
	player.Move(input.Horizontal, input.Vertical, input.Fire)
	if input.Fire {
		x, y := player.GetPosition()
		playerWeapon.Shot(x, y, player.GetNormalizedDegree(), playerShots[:])
	}

	// player shots
	for i := 0; i < len(playerShots); i++ {
		p := playerShots[i]
		if p.IsActive() == false {
			continue
		}
		p.Move()
	}

	// enemies
	for i := 0; i < len(enemies); i++ {
		e := enemies[i]
		if e.IsActive() == false {
			continue
		}
		e.Move(player)
	}
}

// Draw draws the scene
func (stg *Shooting) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x30, 0xff})

	field.Draw(screen)

	// enemies
	for i := 0; i < len(enemies); i++ {
		e := enemies[i]
		if e.IsActive() == false {
			continue
		}
		e.Draw(screen)
	}

	player.Draw(screen)

	// player shots
	for i := 0; i < len(playerShots); i++ {
		p := playerShots[i]
		if p.IsActive() == false {
			continue
		}
		p.Draw(screen)
	}

	uiBackground.Draw(screen)
	input.Draw(screen)
}

func initEnemies() {
	enemyCount := 1

	for i := 0; i < enemyCount; i++ {
		enemy := enemies[i]
		enemy.InitEnemy(actors.EnemyKindBall)
	}
}
