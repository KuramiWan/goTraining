package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Vector struct {
	X float64
	Y float64
}

type Player struct {
	playPosition Vector
	sprite       *ebiten.Image
}

func newPlayer() Player {
	sprite := PlaySprite
	HalfW := sprite.Bounds().Dx()
	HalfH := sprite.Bounds().Dy()
	return Player{
		playPosition: Vector{X: float64(ScreenWidth-HalfW) / 2, Y: float64(ScreenHeight-HalfH) / 2},
		sprite:       sprite,
	}
}

func (p *Player) Update() {
	p.movement()
}

func (p *Player) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(p.playPosition.X, p.playPosition.Y)
	screen.DrawImage(p.sprite, options)
}
func (p *Player) movement() {
	speed := 5.0
	var vector Vector
	if vector.X != 0 || vector.Y != 0 {
		factor := speed / math.Sqrt(vector.X*vector.X+vector.Y*vector.Y)
		vector.X *= factor
		vector.Y *= factor
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.playPosition.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.playPosition.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.playPosition.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.playPosition.X += speed
	}
}
