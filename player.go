package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

const (
	bulletOffset = 50
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Normalize() Vector {
	magnitude := math.Sqrt(v.X*v.X + v.Y*v.Y)
	if magnitude != 0 {
		v.X /= magnitude
		v.Y /= magnitude
	}
	return v
}

type Player struct {
	playPosition Vector
	sprite       *ebiten.Image
	rotation     float64
	bullets      []*Bullet
	coldTimer    Timer
	speed        float64
}

func newPlayer() *Player {
	sprite := PlaySprite
	HalfW := sprite.Bounds().Dx()
	HalfH := sprite.Bounds().Dy()
	p := &Player{
		playPosition: Vector{X: float64(ScreenWidth-HalfW) / 2, Y: float64(ScreenHeight-HalfH) / 2},
		sprite:       sprite,
		coldTimer:    *NewTimer(1 * time.Second),
		bullets:      newBullets(),
	}
	return p
}

func (p *Player) Update() {
	move := p.movement()
	p.rotate()
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	pos := Vector{p.playPosition.X + halfW + math.Sin(p.rotation)*bulletOffset, p.playPosition.Y + halfH - math.Cos(p.rotation)*bulletOffset}
	p.coldTimer.UpdateTicks()
	if p.coldTimer.IsReadyAttack() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.coldTimer.RestTicks()
		p.bullets = append(p.bullets, newBullet(pos, p.rotation, move))
	}
	for _, b := range p.bullets {
		b.Update()
	}
}
func (p *Player) Draw(s *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	options.GeoM.Translate(-halfW, -halfH)
	options.GeoM.Rotate(p.rotation)
	options.GeoM.Translate(halfW, halfH)
	options.GeoM.Translate(p.playPosition.X, p.playPosition.Y)
	s.DrawImage(p.sprite, options)
	for _, bullet := range p.bullets {
		bullet.Draw(s)
	}
	for _, b := range p.bullets {
		b.Draw(s)
	}
}

// update return moved Vector
func (p *Player) movement() Vector {
	speed := 5.0
	var move Vector
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		move.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		move.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		move.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		move.X += speed
	}
	if move.X != 0 || move.Y != 0 {
		factor := speed / math.Sqrt(move.X*move.X+move.Y*move.Y)
		move.X *= factor
		move.Y *= factor
	}
	p.playPosition.X += move.X
	p.playPosition.Y += move.Y
	return move
}

func (p *Player) rotate() {
	speed := math.Pi / float64(ebiten.TPS())
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		p.rotation -= speed
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
		p.rotation += speed
	}
}

func (p *Player) Collider() *Rect {
	return newRect(p.playPosition, p.sprite)
}
