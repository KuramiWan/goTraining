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
	canSplit      bool
	tier          int // 0=big, 1=med, 2=small, 3=tiny
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
	if m.sprite == nil {
		return
	}
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

func randomMeteorSpriteAndTier() (*ebiten.Image, int) {
	allTiers := [][]*ebiten.Image{
		MeteorBigSprites,
		MeteorMedSprites,
		MeteorSmallSprites,
	}
	total := len(MeteorBigSprites) + len(MeteorMedSprites) + len(MeteorSmallSprites)
	if total == 0 {
		return nil, 0
	}
	idx := rand.Intn(total)
	for tier, sprites := range allTiers {
		if idx < len(sprites) {
			return sprites[idx], tier
		}
		idx -= len(sprites)
	}
	return MeteorBigSprites[0], 0
}

func newMeteor() *Meteor {
	sprite, tier := randomMeteorSpriteAndTier()
	r := ScreenWidth / 2.0
	angle := rand.Float64() * 2 * math.Pi
	p := Vector{
		X: target.X + r*math.Cos(angle),
		Y: target.Y + r*math.Sin(angle),
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
		canSplit:      tier < 3,
		tier:          tier,
	}
}

func newSplitMeteor(pos Vector, parentTier int) *Meteor {
	childTier := parentTier + 1
	if childTier > 3 {
		return nil
	}
	tiers := [][]*ebiten.Image{
		MeteorBigSprites,
		MeteorMedSprites,
		MeteorSmallSprites,
		MeteorTinySprites,
	}
	sprites := tiers[childTier]
	if len(sprites) == 0 {
		return nil
	}
	sprite := sprites[rand.Intn(len(sprites))]

	angle := rand.Float64() * 2 * math.Pi
	speed := 1.5 + rand.Float64()*2.5
	move := Vector{
		X: math.Cos(angle) * speed,
		Y: math.Sin(angle) * speed,
	}
	rotation := -0.08 + rand.Float64()*0.16
	return &Meteor{
		sprite:        sprite,
		position:      Vector{X: pos.X, Y: pos.Y},
		movement:      move,
		rotationSpeed: rotation,
		canSplit:      childTier < 3,
		tier:          childTier,
	}
}

func (m *Meteor) Collider() *Rect {
	return newRect(m.position, m.sprite)
}

func (m *Meteor) IsOffScreen() bool {
	margin := 200.0
	return m.position.X < -margin ||
		m.position.X > ScreenWidth+margin ||
		m.position.Y < -margin ||
		m.position.Y > ScreenHeight+margin
}
