package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

type Bullets struct {
	coldTimer Timer
	position  Vector
	sprite    *ebiten.Image
}

func (b *Bullets) Update() {
	b.coldTimer.UpdateTicks()
	if b.coldTimer.IsReadyAttack() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		b.coldTimer.RestTicks()
		// Do Attack
	}
}

func (b *Bullets) Draw(s *ebiten.Image) {
	sprite := LaserSprite
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(b.position.X, b.position.Y)
	s.DrawImage(sprite, options)
}

func newBullets(p *Player) *Bullets {
	return &Bullets{
		coldTimer: *NewTimer(1 * time.Second),
		position: Vector{
			X: p.playPosition.X + math.Cos(p.rotation)*float64(p.sprite.Bounds().Dx()/2.0),
			Y: p.playPosition.Y + math.Sin(p.rotation)*float64(p.sprite.Bounds().Dy()/2.0),
		},
	}
}
