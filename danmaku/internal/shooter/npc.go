package shooter

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/yohamta/godanmaku/danmaku/internal/shared"
	"github.com/yohamta/godanmaku/danmaku/internal/util"
)

type NPCController struct{}

func (c *NPCController) init(sh *Shooter) {
	c.updateDestination(sh)
}

func (c *NPCController) update(sh *Shooter) {
	sh.SetPosition(sh.x+sh.vx, sh.y+sh.vy)

	if c.isArrived(sh) {
		c.updateDestination(sh)
	}

	target := sh.target

	if rand.Float64() < 0.05 {
		sh.degree = util.RadToDeg(math.Atan2(target.GetY()-sh.y, target.GetX()-sh.x))
	}
}

func (c *NPCController) draw(sh *Shooter, screen *ebiten.Image) {
	sh.spr.SetPosition(sh.x-shared.OffsetX, sh.y-shared.OffsetY)
	sh.spr.SetIndex(util.DegreeToDirectionIndex(sh.degree))
	sh.spr.Draw(screen)
}

func (c *NPCController) isArrived(sh *Shooter) bool {
	return math.Abs(sh.y-sh.destination.y) < sh.GetHeight() &&
		math.Abs(sh.x-sh.destination.x) < sh.GetWidth()
}

func (c *NPCController) updateDestination(sh *Shooter) {
	f := sh.field
	x := (f.GetRight() - f.GetLeft()) * rand.Float64()
	y := (f.GetBottom() - f.GetTop()) * rand.Float64()
	sh.destination.x = x
	sh.destination.y = y
	rad := math.Atan2(y-sh.y, x-sh.x)
	sh.vx = math.Cos(rad) * sh.speed
	sh.vy = math.Sin(rad) * sh.speed
}
