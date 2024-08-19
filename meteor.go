package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
)

type Meteor struct {
	sprite        *ebiten.Image
	position      Vector
	rotationSpeed float64
	movement      Vector
}

var target = Vector{
	X: ScreenWidth / 2,
	Y: ScreenHeight / 2,
}

func (m *Meteor) Update() {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
}

func (m *Meteor) Draw(I *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	bounds := m.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	options.GeoM.Translate(-halfW, -halfH)
	options.GeoM.Rotate(m.rotationSpeed)
	options.GeoM.Translate(halfW, halfH)
	options.GeoM.Translate(m.position.X, m.position.Y)
	I.DrawImage(m.sprite, options)
}

func newMeteor() *Meteor {
	sprite := MeteorSprites[rand.Intn(len(MeteorSprites))]
	r := ScreenWidth / 2.0
	angle := rand.Float64() * 2 * math.Pi
	p := Vector{
		X: target.X + r*math.Cos(angle),
		Y: -target.Y + r*math.Sin(angle),
	}
	velocity := 1.0 + rand.Float64()*1.5
	direction := Vector{X: target.X - p.X, Y: target.Y - p.Y}
	normalized := direction.Normalize()
	move := Vector{X: normalized.X * velocity, Y: normalized.Y * velocity}
	rotation := -0.02 + rand.Float64()*0.04
	return &Meteor{
		sprite:        sprite,
		position:      p,
		movement:      move,
		rotationSpeed: rotation,
	}
}
