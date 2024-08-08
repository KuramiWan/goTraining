package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	bulletSpeedPerSecond = 300
)

type Bullet struct {
	position Vector
	sprite   *ebiten.Image
	rotation float64
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())
	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += -math.Cos(b.rotation) * speed

}

func (b *Bullet) Draw(s *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	bounds := b.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	options.GeoM.Translate(-halfW, -halfH)
	options.GeoM.Rotate(b.rotation)
	options.GeoM.Translate(halfW, halfH)
	options.GeoM.Translate(b.position.X, b.position.Y)
	s.DrawImage(b.sprite, options)
}

func newBullet(p Vector, r float64) *Bullet {
	sprite := LaserSprite
	halfW := float64(sprite.Bounds().Dx() / 2.0)
	halfH := float64(sprite.Bounds().Dy() / 2.0)
	p.Y -= halfH
	p.X -= halfW
	return &Bullet{
		position: Vector{p.X, p.Y},
		rotation: r,
		sprite:   sprite,
	}
}

func newBullets() []*Bullet {
	return make([]*Bullet, 0)
}

func (b *Bullet) Collider() *Rect {
	return newRect(b.position, b.sprite)
}
