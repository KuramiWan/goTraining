package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type Meteor struct {
	sprite   *ebiten.Image
	position Vector
}

func (m *Meteor) Update() {

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
