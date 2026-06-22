package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyType represents different enemy ship types with varying stats.
type EnemyType int

const (
	EnemyBlack EnemyType = iota
	EnemyBlue
	EnemyGreen
	EnemyRed
)

// EnemyBehavior determines how the enemy moves.
type EnemyBehavior int

const (
	BehStraight EnemyBehavior = iota // fly straight across
	BehWeave                          // sinusoidal movement
	BehCircle                         // orbit and shoot
	BehKamikaze                       // charge at player
)

// Enemy is an AI-controlled enemy ship.
type Enemy struct {
	sprite      *ebiten.Image
	position    Vector
	movement    Vector
	health      int
	maxHealth   int
	shootTimer  *Timer
	bullets     []*EnemyBullet
	scoreValue  int
	enemyType   EnemyType
	behavior    EnemyBehavior
	phaseOffset float64
	age         float64
}

// EnemyBullet is a bullet fired by an enemy.
type EnemyBullet struct {
	position Vector
	sprite   *ebiten.Image
	rotation float64
	speed    float64
	age      int
}

// EnemyManager manages all enemies in the game.
type EnemyManager struct {
	value      []*Enemy
	spawnTimer *Timer
}

// NewEnemyManager creates an enemy manager.
func NewEnemyManager() *EnemyManager {
	return &EnemyManager{
		value:      make([]*Enemy, 0),
		spawnTimer: NewTimer(8 * time.Second),
	}
}

// spawnEnemy creates a new enemy based on difficulty level.
func (em *EnemyManager) spawnEnemy(difficultyLevel int) {
	var eType EnemyType
	var behavior EnemyBehavior

	// Determine type based on difficulty
	r := rand.Intn(100)
	switch {
	case difficultyLevel >= 10 && r < 20:
		eType = EnemyRed
	case difficultyLevel >= 6 && r < 40:
		eType = EnemyGreen
	case difficultyLevel >= 3 && r < 65:
		eType = EnemyBlue
	default:
		eType = EnemyBlack
	}

	// Pick behavior
	br := rand.Intn(100)
	switch {
	case br < 40:
		behavior = BehStraight
	case br < 70:
		behavior = BehWeave
	case br < 90:
		behavior = BehCircle
	default:
		behavior = BehKamikaze
	}

	// Select sprite
	sprites := [][]*ebiten.Image{EnemyBlackSprites, EnemyBlueSprites, EnemyGreenSprites, EnemyRedSprites}
	var sprite *ebiten.Image
	if len(sprites[eType]) > 0 {
		sprite = sprites[eType][rand.Intn(len(sprites[eType]))]
	}

	// Spawn from screen edge, moving toward a random point on screen
	side := rand.Intn(4)
	var pos Vector
	var targetX, targetY float64

	switch side {
	case 0: // top
		pos = Vector{X: rand.Float64() * ScreenWidth, Y: -80}
		targetX = pos.X + rand.Float64()*400 - 200
		targetY = ScreenHeight/2 + rand.Float64()*300
	case 1: // bottom
		pos = Vector{X: rand.Float64() * ScreenWidth, Y: ScreenHeight + 80}
		targetX = pos.X + rand.Float64()*400 - 200
		targetY = ScreenHeight/2 - rand.Float64()*300
	case 2: // left
		pos = Vector{X: -80, Y: rand.Float64() * ScreenHeight}
		targetX = ScreenWidth/2 + rand.Float64()*300
		targetY = pos.Y + rand.Float64()*300 - 150
	case 3: // right
		pos = Vector{X: ScreenWidth + 80, Y: rand.Float64() * ScreenHeight}
		targetX = ScreenWidth/2 - rand.Float64()*300
		targetY = pos.Y + rand.Float64()*300 - 150
	}

	dx := targetX - pos.X
	dy := targetY - pos.Y
	mag := math.Sqrt(dx*dx + dy*dy)
	speed := 1.5 + rand.Float64()*2.0
	speedByType := []float64{1.5, 2.5, 1.8, 3.0}

	healthByType := []int{2, 3, 4, 6}
	scoreByType := []int{5, 8, 12, 20}
	fireRateByType := []time.Duration{3 * time.Second, 2 * time.Second, 1500 * time.Millisecond, 1 * time.Second}

	if int(eType) < len(speedByType) {
		speed = speedByType[eType]
	}

	e := &Enemy{
		sprite:    sprite,
		position:  pos,
		movement:  Vector{X: dx / mag * speed, Y: dy / mag * speed},
		health:    healthByType[eType],
		maxHealth: healthByType[eType],
		shootTimer: NewTimer(fireRateByType[eType]),
		bullets:   make([]*EnemyBullet, 0),
		scoreValue: scoreByType[eType],
		enemyType: eType,
		behavior:  behavior,
		phaseOffset: rand.Float64() * 2 * math.Pi,
	}
	em.value = append(em.value, e)
}

// Update advances all enemies.
func (em *EnemyManager) Update() {
	em.spawnTimer.UpdateTicks()
	for _, e := range em.value {
		e.Update()
	}
	// Cleanup off-screen and dead enemies
	surviving := em.value[:0]
	for _, e := range em.value {
		if e.health > 0 && !e.IsOffScreen() {
			surviving = append(surviving, e)
		}
	}
	em.value = surviving
}

// TrySpawn spawns an enemy if the timer is ready.
func (em *EnemyManager) TrySpawn(difficultyLevel int) bool {
	if em.spawnTimer.IsReadyAttack() {
		em.spawnTimer.RestTicks()
		em.spawnEnemy(difficultyLevel)
		return true
	}
	return false
}

// Update advances enemy position and behavior.
func (e *Enemy) Update() {
	e.age++

	// Apply behavior-based movement adjustments
	switch e.behavior {
	case BehWeave:
		e.movement.Y += math.Sin((e.age+e.phaseOffset)*0.05) * 0.3
	case BehCircle:
		e.movement.X += math.Cos((e.age+e.phaseOffset)*0.03) * 0.2
		e.movement.Y += math.Sin((e.age+e.phaseOffset)*0.03) * 0.2
	case BehKamikaze:
		// Gradually accelerate toward center
		e.movement.X *= 1.01
		e.movement.Y *= 1.01
	}

	// Clamp speed
	speed := math.Sqrt(e.movement.X*e.movement.X + e.movement.Y*e.movement.Y)
	maxSpeed := 6.0
	if speed > maxSpeed {
		e.movement.X = e.movement.X / speed * maxSpeed
		e.movement.Y = e.movement.Y / speed * maxSpeed
	}

	e.position.X += e.movement.X
	e.position.Y += e.movement.Y

	// Update bullets
	e.shootTimer.UpdateTicks()
	for _, b := range e.bullets {
		b.position.X += math.Sin(b.rotation) * b.speed
		b.position.Y += -math.Cos(b.rotation) * b.speed
		b.age++
	}

	// Cleanup old bullets
	surviving := e.bullets[:0]
	for _, b := range e.bullets {
		if b.age < 180 { // 3 seconds at 60 TPS
			surviving = append(surviving, b)
		}
	}
	e.bullets = surviving
}

// Shoot fires a bullet toward the given target position.
func (e *Enemy) Shoot(targetX, targetY float64) {
	if !e.shootTimer.IsReadyAttack() {
		return
	}
	e.shootTimer.RestTicks()

	// Aim at target
	dx := targetX - e.position.X
	dy := targetY - e.position.Y
	angle := math.Atan2(dx, -dy)

	bulletSprites := LaserRedSprites
	if len(bulletSprites) == 0 {
		return
	}
	sprite := bulletSprites[rand.Intn(len(bulletSprites))]

	e.bullets = append(e.bullets, &EnemyBullet{
		position: Vector{X: e.position.X, Y: e.position.Y},
		sprite:   sprite,
		rotation: angle,
		speed:    3.0 + rand.Float64()*2.0,
	})
}

// Draw renders the enemy.
func (e *Enemy) Draw(s *ebiten.Image) {
	if e.sprite == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	bounds := e.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(e.position.X, e.position.Y)
	s.DrawImage(e.sprite, op)

	// Draw enemy bullets
	for _, b := range e.bullets {
		if b.sprite == nil {
			continue
		}
		bop := &ebiten.DrawImageOptions{}
		bb := b.sprite.Bounds()
		bhw := float64(bb.Dx()) / 2
		bhh := float64(bb.Dy()) / 2
		bop.GeoM.Translate(-bhw, -bhh)
		bop.GeoM.Rotate(b.rotation)
		bop.GeoM.Translate(bhw, bhh)
		bop.GeoM.Translate(b.position.X, b.position.Y)
		s.DrawImage(b.sprite, bop)
	}
}

// TakeDamage applies damage to the enemy. Returns true if the enemy died.
func (e *Enemy) TakeDamage(amount int) bool {
	e.health -= amount
	return e.health <= 0
}

// Collider returns the enemy's collision rectangle.
func (e *Enemy) Collider() *Rect {
	return newRect(e.position, e.sprite)
}

// BulletCollider returns a bullet's collision rectangle.
func (b *EnemyBullet) Collider() *Rect {
	return &Rect{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  float64(b.sprite.Bounds().Dx()),
		Height: float64(b.sprite.Bounds().Dy()),
	}
}

// IsOffScreen checks if the enemy is off screen.
func (e *Enemy) IsOffScreen() bool {
	margin := 200.0
	return e.position.X < -margin ||
		e.position.X > ScreenWidth+margin ||
		e.position.Y < -margin ||
		e.position.Y > ScreenHeight+margin
}

// Count returns the number of active enemies.
func (em *EnemyManager) Count() int {
	return len(em.value)
}

// SetSpawnInterval adjusts the spawn interval.
func (em *EnemyManager) SetSpawnInterval(d time.Duration) {
	em.spawnTimer.SetDuration(d)
}
