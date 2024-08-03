package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
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
}

func newPlayer() *Player {
	sprite := PlaySprite
	HalfW := sprite.Bounds().Dx()
	HalfH := sprite.Bounds().Dy()
	return &Player{
		playPosition: Vector{X: float64(ScreenWidth-HalfW) / 2, Y: float64(ScreenHeight-HalfH) / 2},
		sprite:       sprite,
	}
}

func (p *Player) Update() {
	p.movement()
	p.rotate()
}

func (p *Player) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	options.GeoM.Translate(-halfW, -halfH)
	options.GeoM.Rotate(p.rotation)
	options.GeoM.Translate(halfW, halfH)
	options.GeoM.Translate(p.playPosition.X, p.playPosition.Y)
	screen.DrawImage(p.sprite, options)
}

// update
func (p *Player) movement() {
	speed := 5.0
	var vector Vector
	if vector.X != 0 || vector.Y != 0 {
		factor := speed / math.Sqrt(vector.X*vector.X+vector.Y*vector.Y)
		vector.X *= factor
		vector.Y *= factor
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		p.playPosition.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.playPosition.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.playPosition.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.playPosition.X += speed
	}
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
