package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// BossManager controls boss spawning and lifecycle.
type BossManager struct {
	currentBoss      *Boss
	bossSpawnTimer   *Timer
	bossActive       bool
	bossWarningTimer *Timer
	warningTick      int
}

// Boss represents a boss enemy.
type Boss struct {
	sprite        *ebiten.Image
	position      Vector
	health        int
	maxHealth     int
	phase         int // 0=entry, 1=phase1, 2=phase2, 3=enrage
	attackTimer   *Timer
	patternTimer  *Timer
	bullets       []*EnemyBullet
	rotation      float64
	entryComplete bool
	flashCounter  int
	scoreValue    int
	currentPattern int
	phaseAge      int
}

const (
	BossPhaseEntry = iota
	BossPhase1
	BossPhase2
	BossPhase3
)

// NewBossManager creates a boss manager.
func NewBossManager() *BossManager {
	return &BossManager{
		bossSpawnTimer:   NewTimer(90 * time.Second),
		bossActive:       false,
		bossWarningTimer: NewTimer(3 * time.Second),
	}
}

// SpawnBoss creates a new boss.
func (bm *BossManager) SpawnBoss(difficultyLevel int) {
	sprites := BossSprites
	if len(sprites) == 0 {
		// Fall back to UFO sprites
		sprites = []*ebiten.Image{UfoBlueSprite, UfoRedSprite, UfoGreenSprite, UfoYellowSprite}
	}
	var sprite *ebiten.Image
	if len(sprites) > 0 {
		sprite = sprites[rand.Intn(len(sprites))]
	}

	boss := &Boss{
		sprite:     sprite,
		position:   Vector{X: float64(ScreenWidth) / 2, Y: -150},
		health:     50 + difficultyLevel*10,
		maxHealth:  50 + difficultyLevel*10,
		phase:      BossPhaseEntry,
		attackTimer: NewTimer(1500 * time.Millisecond),
		patternTimer: NewTimer(5 * time.Second),
		bullets:    make([]*EnemyBullet, 0),
		scoreValue: 100 + difficultyLevel*10,
	}

	bm.currentBoss = boss
	bm.bossActive = true
	bm.warningTick = 0
}

// Update advances the boss state.
func (bm *BossManager) Update() {
	if bm.currentBoss == nil && bm.bossActive {
		bm.bossActive = false
		return
	}
	if !bm.bossActive {
		return
	}

	boss := bm.currentBoss
	boss.phaseAge++

	// Phase transitions
	switch boss.phase {
	case BossPhaseEntry:
		// Move down to center area
		targetY := float64(ScreenHeight) / 3
		boss.position.Y += 1.5
		if boss.position.Y >= targetY {
			boss.position.Y = targetY
			boss.entryComplete = true
			boss.phase = BossPhase1
		}
	case BossPhase1:
		if float64(boss.health) <= float64(boss.maxHealth)*0.5 {
			boss.phase = BossPhase2
		}
	case BossPhase2:
		if float64(boss.health) <= float64(boss.maxHealth)*0.25 {
			boss.phase = BossPhase3
		}
	}

	// Side-to-side movement in later phases
	if boss.phase >= BossPhase1 {
		boss.position.X += math.Sin(float64(boss.phaseAge)*0.02) * 2.0
		// Clamp to screen
		boss.position.X = math.Max(100, math.Min(ScreenWidth-100, boss.position.X))
	}

	// Attack patterns
	boss.attackTimer.UpdateTicks()
	boss.patternTimer.UpdateTicks()

	if boss.patternTimer.IsReadyAttack() && boss.phase >= BossPhase1 {
		boss.patternTimer.RestTicks()
		boss.currentPattern = rand.Intn(4)
	}

	if boss.attackTimer.IsReadyAttack() {
		boss.attackTimer.RestTicks()
		switch boss.currentPattern {
		case 0:
			boss.attackRadial()
		case 1:
			boss.attackSpiral()
		case 2:
			boss.attackAimed()
		case 3:
			boss.attackBurst()
		}
		// Enrage: faster attacks
		if boss.phase == BossPhase3 {
			boss.attackTimer.SetDuration(500 * time.Millisecond)
		}
	}

	// Update boss bullets
	for _, b := range boss.bullets {
		b.position.X += math.Sin(b.rotation) * b.speed
		b.position.Y += -math.Cos(b.rotation) * b.speed
		b.age++
	}
	surviving := boss.bullets[:0]
	for _, b := range boss.bullets {
		if b.age < 300 {
			surviving = append(surviving, b)
		}
	}
	boss.bullets = surviving

	// Flash counter for hit effect
	if boss.flashCounter > 0 {
		boss.flashCounter--
	}
}

// attackRadial fires bullets in all directions.
func (b *Boss) attackRadial() {
	count := 12
	for i := 0; i < count; i++ {
		angle := float64(i) * 2 * math.Pi / float64(count)
		b.spawnBullet(angle, 3.0)
	}
}

// attackSpiral fires bullets in a rotating spiral.
func (b *Boss) attackSpiral() {
	baseAngle := b.rotation
	count := 8
	for i := 0; i < count; i++ {
		angle := baseAngle + float64(i)*2*math.Pi/float64(count)
		b.spawnBullet(angle, 2.5)
	}
	b.rotation += math.Pi / 8
}

// attackAimed fires multiple bullets aimed at center.
func (b *Boss) attackAimed() {
	targetX := ScreenWidth / 2.0
	targetY := ScreenHeight / 2.0
	for i := -1; i <= 1; i++ {
		dx := targetX - b.position.X
		dy := targetY - b.position.Y
		angle := math.Atan2(dx, -dy) + float64(i)*0.2
		b.spawnBullet(angle, 4.0)
	}
}

// attackBurst fires a tight spread of fast bullets.
func (b *Boss) attackBurst() {
	count := 5
	baseAngle := b.rotation
	for i := 0; i < count; i++ {
		angle := baseAngle - 0.3 + float64(i)*0.15
		b.spawnBullet(angle, 5.0)
	}
	b.rotation += math.Pi / 6
}

func (b *Boss) spawnBullet(angle, speed float64) {
	sprites := LaserRedSprites
	var sprite *ebiten.Image
	if len(sprites) > 0 {
		sprite = sprites[rand.Intn(len(sprites))]
	}
	b.bullets = append(b.bullets, &EnemyBullet{
		position: Vector{X: b.position.X, Y: b.position.Y},
		sprite:   sprite,
		rotation: angle,
		speed:    speed,
	})
}

// TakeDamage applies damage to the boss. Returns true if the boss died.
func (b *Boss) TakeDamage(amount int) bool {
	b.health -= amount
	b.flashCounter = 4
	if b.health <= 0 {
		return true
	}
	return false
}

// Draw renders the boss and its bullets.
func (bm *BossManager) Draw(s *ebiten.Image) {
	if bm.currentBoss == nil {
		return
	}
	b := bm.currentBoss

	// Flash on hit
	if b.flashCounter > 0 && b.flashCounter%2 == 0 {
		return
	}

	// Draw boss
	if b.sprite != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := b.sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(b.rotation * 0.1)
		op.GeoM.Translate(halfW, halfH)
		op.GeoM.Translate(b.position.X, b.position.Y)
		s.DrawImage(b.sprite, op)
	}

	// Draw boss bullets
	for _, bul := range b.bullets {
		if bul.sprite == nil {
			continue
		}
		bop := &ebiten.DrawImageOptions{}
		bb := bul.sprite.Bounds()
		bhw := float64(bb.Dx()) / 2
		bhh := float64(bb.Dy()) / 2
		bop.GeoM.Translate(-bhw, -bhh)
		bop.GeoM.Rotate(bul.rotation)
		bop.GeoM.Translate(bhw, bhh)
		bop.GeoM.Translate(bul.position.X, bul.position.Y)
		s.DrawImage(bul.sprite, bop)
	}

	// Boss name label
	text.Draw(s, "BOSS", UIFont, int(b.position.X)-40, int(b.position.Y)-80, color.RGBA{255, 50, 50, 255})
}

// DrawHealthBar renders the boss health bar.
func (bm *BossManager) DrawHealthBar(s *ebiten.Image) {
	if bm.currentBoss == nil {
		return
	}
	b := bm.currentBoss
	barW := 600.0
	barH := 24.0
	x := float64(ScreenWidth)/2 - barW/2
	y := 20.0

	// Border
	border := ebiten.NewImage(int(barW)+4, int(barH)+4)
	border.Fill(color.RGBA{80, 80, 80, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-2, y-2)
	s.DrawImage(border, op)

	// Background
	bg := ebiten.NewImage(int(barW), int(barH))
	bg.Fill(color.RGBA{40, 40, 40, 255})
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(x, y)
	s.DrawImage(bg, op2)

	// Health
	pct := float64(b.health) / float64(b.maxHealth)
	if pct < 0 {
		pct = 0
	}
	fillW := barW * pct
	if fillW > 0 {
		var c color.RGBA
		switch {
		case pct > 0.5:
			c = color.RGBA{0, 200, 0, 255}
		case pct > 0.25:
			c = color.RGBA{220, 220, 0, 255}
		default:
			c = color.RGBA{220, 0, 0, 255}
		}
		fill := ebiten.NewImage(int(fillW), int(barH))
		fill.Fill(c)
		s.DrawImage(fill, op2)
	}

	// Text
	label := fmt.Sprintf("%d / %d", b.health, b.maxHealth)
	text.Draw(s, label, UIFont, int(x+barW/2-50), int(y+barH-5), color.White)
}

// IsActive returns whether a boss fight is in progress.
func (bm *BossManager) IsActive() bool {
	return bm.bossActive
}

// Clear removes the current boss.
func (bm *BossManager) Clear() {
	bm.currentBoss = nil
	bm.bossActive = false
}

// SetSpawnInterval adjusts the boss spawn interval.
func (bm *BossManager) SetSpawnInterval(d time.Duration) {
	bm.bossSpawnTimer.SetDuration(d)
}
