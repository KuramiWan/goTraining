package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

const (
	bulletOffset  = 50
	spreadDegrees = 70
)

type Vector struct {
	X  float64
	Y  float64
	Vx float64
	Vy float64
	Ax float64
	Ay float64
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
	playPosition    Vector
	sprite          *ebiten.Image
	rotation        float64
	bullets         []*Bullet
	coldTimer       Timer
	speed           float64
	cdMultiplier    float64
	extraSpread     int
	piercingActive  bool
	// Phase 3: weapon system
	weapon          *WeaponInfo
	weaponUpgrades  int
	// Phase 3: shield
	shieldActive    bool
	shieldTimer     *Timer
	// Phase 3: speed boost
	speedBoostActive bool
	speedBoostTimer  *Timer
	speedMultiplier  float64
}

func newPlayer() *Player {
	sprite := PlaySprite
	HalfW := sprite.Bounds().Dx()
	HalfH := sprite.Bounds().Dy()
	p := &Player{
		playPosition:   Vector{X: float64(ScreenWidth-HalfW) / 2, Y: float64(ScreenHeight-HalfH) / 2},
		sprite:         sprite,
		coldTimer:      *NewTimer(1 * time.Second),
		bullets:        newBullets(),
		cdMultiplier:   1.0,
		extraSpread:    0,
		piercingActive: false,
		weapon:         GetWeaponInfo(WpnBlue),
		speedMultiplier: 1.0,
	}
	return p
}

func (p *Player) Update() {
	move := p.movement()
	p.clampToScreen()
	p.rotate()
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	pos := Vector{p.playPosition.X + halfW + math.Sin(p.rotation)*bulletOffset, p.playPosition.Y + halfH - math.Cos(p.rotation)*bulletOffset, 0, 0, 0, 0}
	p.coldTimer.UpdateTicks()
	if p.coldTimer.IsReadyAttack() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.coldTimer.RestTicks()
		p.shoot(pos, move)
	}
	for _, b := range p.bullets {
		b.Update()
	}
}

func (p *Player) shoot(pos, move Vector) {
	totalBullets := 1 + p.extraSpread
	halfSpreadRad := (spreadDegrees / 2.0) * math.Pi / 180.0

	for i := 0; i < totalBullets; i++ {
		var angle float64
		if totalBullets == 1 {
			angle = p.rotation
		} else {
			t := float64(i) / float64(totalBullets-1)
			angle = p.rotation - halfSpreadRad + 2*halfSpreadRad*t
		}
		p.bullets = append(p.bullets, newBullet(pos, angle, move, p.piercingActive))
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
}

func (p *Player) movement() Vector {
	a := 1.2
	move := &p.playPosition
	const friction = 0.8

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		move.Ax += -a
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		move.Ax += a
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		move.Ay += -a
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		move.Ay += a
	}

	move.Vx += move.Ax
	move.Vy += move.Ay
	move.X += move.Vx
	move.Y += move.Vy
	move.Vx *= friction
	move.Vy *= friction
	move.Ax = 0
	move.Ay = 0
	p.playPosition = *move
	return Vector{move.Vx, move.Vy, 0, 0, 0, 0}
}

func (p *Player) clampToScreen() {
	w := float64(p.sprite.Bounds().Dx())
	h := float64(p.sprite.Bounds().Dy())
	p.playPosition.X = math.Max(0, math.Min(float64(ScreenWidth)-w, p.playPosition.X))
	p.playPosition.Y = math.Max(0, math.Min(float64(ScreenHeight)-h, p.playPosition.Y))
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

func (p *Player) SetFireRate(d time.Duration) {
	effective := time.Duration(float64(d) * p.cdMultiplier)
	p.coldTimer.SetDuration(effective)
}

func (p *Player) ApplyPowerUp(t PowerUpType) {
	switch t {
	case PowerUpSilver:
		p.cdMultiplier *= 0.5
	case PowerUpBronze:
		p.extraSpread++
	case PowerUpGold:
		p.piercingActive = true
	case PowerUpShield:
		p.shieldActive = true
		p.shieldTimer = NewTimer(10 * time.Second)
	case PowerUpHealth:
		// Health is handled in game.go (g.Lives++)
	case PowerUpSpeed:
		p.speedBoostActive = true
		p.speedMultiplier = 1.5
		p.speedBoostTimer = NewTimer(15 * time.Second)
	}
}

func (p *Player) ResetBuffs() {
	p.cdMultiplier = 1.0
	p.extraSpread = 0
	p.piercingActive = false
}
