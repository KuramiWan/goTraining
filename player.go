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
	return Player{
		playPosition: Vector{X: 100, Y: 100},
		sprite:       PlaySprite,
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
