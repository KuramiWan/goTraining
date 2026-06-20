package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
)

type PowerUpType int

const (
	PowerUpSilver PowerUpType = iota
	PowerUpBronze
	PowerUpGold
)

type PowerUp struct {
	pType    PowerUpType
	position Vector
	sprite   *ebiten.Image
	movement Vector
}

func getPowerUpSprite(t PowerUpType) *ebiten.Image {
	switch t {
	case PowerUpSilver:
		return PowerUpSilverSprite
	case PowerUpBronze:
		return PowerUpBronzeSprite
	case PowerUpGold:
		return PowerUpGoldSprite
	}
	return PowerUpSilverSprite
}

func newPowerUp(t PowerUpType) *PowerUp {
	side := rand.Intn(4)
	var pos Vector
	switch side {
	case 0:
		pos = Vector{X: rand.Float64() * ScreenWidth, Y: -60}
	case 1:
		pos = Vector{X: rand.Float64() * ScreenWidth, Y: ScreenHeight + 60}
	case 2:
		pos = Vector{X: -60, Y: rand.Float64() * ScreenHeight}
	case 3:
		pos = Vector{X: ScreenWidth + 60, Y: rand.Float64() * ScreenHeight}
	}

	targetX := rand.Float64() * ScreenWidth
	targetY := rand.Float64() * ScreenHeight
	dx := targetX - pos.X
	dy := targetY - pos.Y
	mag := math.Sqrt(dx*dx + dy*dy)
	speed := 0.8 + rand.Float64()*1.2

	return &PowerUp{
		pType:    t,
		position: pos,
		sprite:   getPowerUpSprite(t),
		movement: Vector{X: dx / mag * speed, Y: dy / mag * speed},
	}
}

func (pu *PowerUp) Update() {
	pu.position.X += pu.movement.X
	pu.position.Y += pu.movement.Y
}

func (pu *PowerUp) Draw(s *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	bounds := pu.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	options.GeoM.Translate(-halfW, -halfH)
	options.GeoM.Translate(halfW, halfH)
	options.GeoM.Translate(pu.position.X, pu.position.Y)
	s.DrawImage(pu.sprite, options)
}

func (pu *PowerUp) Collider() *Rect {
	return newRect(pu.position, pu.sprite)
}

func (pu *PowerUp) IsOffScreen() bool {
	margin := 120.0
	return pu.position.X < -margin ||
		pu.position.X > ScreenWidth+margin ||
		pu.position.Y < -margin ||
		pu.position.Y > ScreenHeight+margin
}
