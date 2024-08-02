package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
)

type Meteor struct {
	sprite   *ebiten.Image
	position Vector
}

func (m *Meteor) Update() {
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}
	r := ScreenWidth / 2.0
	angle := rand.Float64() * 2 * math.Pi
	pos := Vector{
		X: target.X + r*math.Cos(angle),
		Y: target.Y + r*math.Sin(angle),
	}
	velocity := 0.25 + rand.Float64()*1.5
	direction := Vector{X: target.X - pos.X, Y: target.X - pos.Y}
	direction.Normalize()
}

func (m *Meteor) Draw(screen *ebiten.Image) {

}

func newMeteor() Meteor {
	sprite := MeteorSprites[rand.Intn(len(MeteorSprites))]
	return Meteor{
		sprite:   sprite,
		position: Vector{},
	}
}
