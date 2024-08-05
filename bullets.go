package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

type Bullets struct {
	coldTimer Timer
	movement  Vector
}

func (b *Bullets) Update() {

}

func (b *Bullets) Draw(s *ebiten.Image) {

}

func newBullets(p *Player) *Bullets {
	return &Bullets{
		coldTimer: *NewTimer(1 * time.Second),
		movement: Vector{
			X: p.playPosition.X + math.Cos(p.rotation)*float64(p.sprite.Bounds().Dx()/2.0),
			Y: p.playPosition.Y + math.Sin(p.rotation)*float64(p.sprite.Bounds().Dy()/2.0),
		},
	}
}
