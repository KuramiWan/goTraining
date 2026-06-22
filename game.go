package main

import (
	"embed"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

//go:embed assets
var data embed.FS

type GameState int

const (
	StateStart GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

const (
	ScreenWidth        = 1600
	ScreenHeight       = 900
	InitialLives       = 3
	InvulnDuration     = 2 * time.Second
	BaseSpawnInterval  = 5 * time.Second
	MinSpawnInterval   = 500 * time.Millisecond
	BaseFireRate       = 1 * time.Second
	MinFireRate        = 200 * time.Millisecond
	PowerUpSpawnPeriod = 60 * time.Second
)

type Game struct {
	player              *Player
	meteors             *Meteors
	powerUps            []*PowerUp
	Score               int
	Lives               int
	State               GameState
	DifficultyLevel     int
	invulnTimer         *Timer
	firstSpawnDone      bool
	flashCounter        int
	confirmQuit         bool
	elapsedTime         time.Duration
	powerUpSpawnTimer   *Timer
	lastBronzeMilestone int
	// Phase 1: testability infrastructure
	Config       *Config
	HighScores   *HighScoreManager
	Achievements *AchievementTracker
	Stats        *GameStats
	SaveMgr      *SaveManager
	DebugMode    bool
	Audio        *AudioManager
	// Phase 1: debug toggles
	godMode          bool
	showHitboxes     bool
	slowMotion       bool
	// Phase 3: boss tracking (referenced in save.go)
	bossActive bool
	// Phase 3: effects
	effects    *EffectManager
	background *Background
	hud        *HUD
	// Phase 3: enemies
	enemies  *EnemyManager
	bossMgr  *BossManager
	// Phase 1: achievement popup
	lastAchievementID AchievementID
	achievementPopup  bool
	achievementTimer  *Timer
	// Phase 1: high score entry
	highScoreEntryName string
}

func (g *Game) Update() error {
	switch g.State {
	case StateStart:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startGame()
		}
		return nil
	case StatePlaying:
		if g.confirmQuit {
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				g.State = StateStart
				g.confirmQuit = false
				return nil
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				g.confirmQuit = false
				return nil
			}
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.confirmQuit = true
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.State = StatePaused
			return nil
		}
		g.updatePlaying()
		return nil
	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.State = StatePlaying
		}
		return nil
	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset()
		}
		return nil
	}
	return nil
}

func (g *Game) startGame() {
	g.State = StatePlaying
	g.Score = 0
	g.Lives = InitialLives
	g.DifficultyLevel = 0
	g.firstSpawnDone = false
	g.invulnTimer = nil
	g.flashCounter = 0
	g.confirmQuit = false
	g.elapsedTime = 0
	g.powerUps = nil
	g.powerUpSpawnTimer = NewTimer(PowerUpSpawnPeriod)
	g.lastBronzeMilestone = 10
	g.player = newPlayer()
	g.meteors = newMeteors()
	// Phase 1: init/load persistent systems
	if g.Stats != nil {
		g.Stats.Reset()
	} else {
		g.Stats = &GameStats{}
	}
	if g.Achievements != nil {
		g.Achievements.Load()
	}
	if g.HighScores != nil {
		g.HighScores.Load()
	}
	// Phase 1: debug toggles start off
	g.godMode = false
	g.showHitboxes = false
	g.slowMotion = false
	// Phase 3: init game systems
	g.effects = NewEffectManager()
	g.background = NewBackground()
	if g.hud == nil {
		g.hud = NewHUD()
	}
	g.enemies = NewEnemyManager()
	g.bossMgr = NewBossManager()
	g.bossActive = false
	// Phase 1: init achievement popup state
	g.achievementPopup = false
	g.achievementTimer = nil
}

func (g *Game) updatePlaying() {
	// Phase 1: debug keys (checked before regular input)
	g.checkDebugKeys()

	// Slow motion affects tick rate
	if g.slowMotion && g.flashCounter%2 != 0 {
		g.flashCounter++
		g.elapsedTime += time.Second / time.Duration(ebiten.TPS())
		return
	}

	g.player.Update()
	g.meteors.Update(!g.firstSpawnDone)
	if !g.firstSpawnDone {
		g.firstSpawnDone = true
	}
	// Phase 3: enemy spawning and AI
	g.updateEnemies()
	// Phase 3: boss spawning and AI
	g.updateBoss()
	// Phase 3: effects (explosions etc.)
	if g.effects != nil {
		g.effects.Update()
	}
	// Phase 3: background scroll
	if g.background != nil {
		g.background.Update()
	}
	g.updatePowerUps()
	g.bulletCollisions()
	g.bulletEnemyCollisions()
	g.enemyBulletCollisions()
	g.meteorCollisions()
	g.enemyCollisions()
	g.bossBulletCollisions()
	g.powerUpCollisions()
	g.bronzeMilestoneSpawn()
	g.updateDifficulty()
	g.cleanupOffScreen()
	g.updateInvulnerability()
	g.updateShield()
	g.updateSpeedBoost()
	// Phase 1: check achievements
	g.checkAchievements()
	g.flashCounter++
	g.elapsedTime += time.Second / time.Duration(ebiten.TPS())
}

func (g *Game) updatePowerUps() {
	g.powerUpSpawnTimer.UpdateTicks()
	if g.powerUpSpawnTimer.IsReadyAttack() {
		g.powerUpSpawnTimer.RestTicks()
		g.spawnRandomPowerUp()
	}
	for _, pu := range g.powerUps {
		pu.Update()
	}
}

func (g *Game) spawnRandomPowerUp() {
	r := rand.Intn(100) // 0-99
	var t PowerUpType
	switch {
	case r < 35:
		t = PowerUpSilver
	case r < 60:
		t = PowerUpBronze
	case r < 70:
		t = PowerUpGold
	case r < 85:
		t = PowerUpShield
	case r < 95:
		t = PowerUpSpeed
	default:
		t = PowerUpHealth
	}
	g.powerUps = append(g.powerUps, newPowerUp(t))
}

func (g *Game) bronzeMilestoneSpawn() {
	if g.Score >= g.lastBronzeMilestone {
		g.powerUps = append(g.powerUps, newPowerUp(PowerUpBronze))
		g.lastBronzeMilestone *= 10
	}
}

// Phase 1: Debug mode key handling
func (g *Game) checkDebugKeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.DebugMode = !g.DebugMode
	}
	if !g.DebugMode {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.godMode = !g.godMode
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		// Stats overlay toggled via showHitboxes reusing the overlay area
		g.showHitboxes = !g.showHitboxes
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		g.slowMotion = !g.slowMotion
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		// Save game
		if g.SaveMgr != nil {
			g.SaveMgr.Save(g)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		// Spawn random power-up for testing
		g.spawnRandomPowerUp()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF7) {
		// Toggle hitbox visualization
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF8) {
		// Force boss spawn (will be wired in Phase 3)
	}
	// Number keys to spawn test entities
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
		g.meteors.value = append(g.meteors.value, newMeteor())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
		// Spawn big meteor at player position
		m := spawnTestMeteorAt(g.player.playPosition.X, g.player.playPosition.Y, 0)
		g.meteors.value = append(g.meteors.value, m)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
		// Spawn split-capable meteor
		m := spawnTestMeteorAt(g.player.playPosition.X+100, g.player.playPosition.Y, 1)
		g.meteors.value = append(g.meteors.value, m)
	}
}

// spawnTestMeteorAt is a debug helper to create a specific tier meteor at a position.
func spawnTestMeteorAt(x, y float64, tier int) *Meteor {
	tiers := [][]*ebiten.Image{
		MeteorBigSprites,
		MeteorMedSprites,
		MeteorSmallSprites,
		MeteorTinySprites,
	}
	var sprite *ebiten.Image
	if tier >= 0 && tier < len(tiers) && len(tiers[tier]) > 0 {
		sprite = tiers[tier][0]
	}
	return &Meteor{
		sprite:   sprite,
		position: Vector{X: x, Y: y},
		movement: Vector{X: 1, Y: 0},
		canSplit: tier < 3,
		tier:     tier,
	}
}

// Phase 1: Achievement checking
func (g *Game) checkAchievements() {
	if g.Achievements == nil {
		return
	}
	g.Achievements.CheckStats(g.Stats, g.Score, g.DifficultyLevel, g.elapsedTime)

	// Check no-damage achievement
	if g.Achievements.TimeSinceLastDamage() >= 1*time.Minute && !g.Achievements.IsUnlocked(AchNoDamage1Min) {
		g.Achievements.Unlock(AchNoDamage1Min)
	}

	// Check for newly unlocked achievements to show popup
	// (Popup logic handled in Draw)
}

// Phase 1: Record player shot for stats
func (g *Game) recordShot() {
	if g.Stats != nil {
		g.Stats.RecordShot()
	}
	// Phase 2: play laser sound
	if g.Audio != nil {
		g.Audio.PlayLaser1()
	}
}

// Phase 1: Record hit for stats
func (g *Game) recordHit(isEnemy bool) {
	if g.Stats != nil {
		g.Stats.RecordHit()
		g.Stats.RecordKill(isEnemy)
	}
	if g.Achievements != nil {
		g.Achievements.RecordKill(g.Stats.Kills)
	}
}

// Phase 1: Record damage for stats
func (g *Game) recordDamage() {
	if g.Stats != nil {
		g.Stats.RecordDamage()
	}
	if g.Achievements != nil {
		g.Achievements.RecordDamage()
	}
}

// Phase 1: Record power-up collection for stats
func (g *Game) recordPowerUpCollected(pType PowerUpType) {
	if g.Stats != nil {
		g.Stats.RecordPowerUp()
	}
	if g.Achievements != nil {
		g.Achievements.RecordPowerUpCollected(pType)
	}
}

func (g *Game) powerUpCollisions() {
	pRect := g.player.Collider()
	surviving := g.powerUps[:0]
	for _, pu := range g.powerUps {
		if pu.Collider().Intersect(pRect) {
			g.player.ApplyPowerUp(pu.pType)
			g.recordPowerUpCollected(pu.pType)
			// Phase 3: special handling for health
			if pu.pType == PowerUpHealth {
				g.Lives++
				if g.Lives > g.Config.Player.MaxLives {
					g.Lives = g.Config.Player.MaxLives
				}
			}
			// Phase 2: audio
			if g.Audio != nil {
				if pu.pType == PowerUpShield {
					g.Audio.PlayShieldUp()
				} else {
					g.Audio.PlayPowerUpPickup()
				}
			}
		} else {
			surviving = append(surviving, pu)
		}
	}
	g.powerUps = surviving
}

func (g *Game) updateInvulnerability() {
	if g.invulnTimer != nil {
		g.invulnTimer.UpdateTicks()
	}
}

func (g *Game) updateDifficulty() {
	scoreLevel := g.Score / 10
	timeLevel := int(g.elapsedTime.Seconds()) / 30
	newLevel := scoreLevel + timeLevel
	if newLevel > g.DifficultyLevel {
		g.DifficultyLevel = newLevel
		interval := BaseSpawnInterval - time.Duration(g.DifficultyLevel)*500*time.Millisecond
		if interval < MinSpawnInterval {
			interval = MinSpawnInterval
		}
		g.meteors.meteorSpawnTimer.SetDuration(interval)

		fireRate := BaseFireRate - time.Duration(g.DifficultyLevel)*100*time.Millisecond
		if fireRate < MinFireRate {
			fireRate = MinFireRate
		}
		g.player.SetFireRate(fireRate)
	}
}

func (g *Game) meteorCollisions() {
	// Phase 1: god mode skips damage
	if g.godMode {
		return
	}
	if g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack() {
		return
	}
	pRect := g.player.Collider()
	hitIndex := -1
	for i, meteor := range g.meteors.value {
		if meteor.Collider().Intersect(pRect) {
			hitIndex = i
			break
		}
	}
	if hitIndex >= 0 {
		g.meteors.value = append(g.meteors.value[:hitIndex], g.meteors.value[hitIndex+1:]...)
		g.Lives--
		g.recordDamage()
		if g.Lives <= 0 {
			g.State = StateGameOver
		} else {
			g.invulnTimer = NewTimer(InvulnDuration)
		}
	}
}

func (g *Game) bulletCollisions() {
	if len(g.meteors.value) == 0 || len(g.player.bullets) == 0 {
		return
	}
	meteorHit := make([]bool, len(g.meteors.value))
	bulletHit := make([]bool, len(g.player.bullets))
	var newMeteors []*Meteor

	for i, meteor := range g.meteors.value {
		mRect := meteor.Collider()
		for j, bullet := range g.player.bullets {
			if !meteorHit[i] && !bulletHit[j] && mRect.Intersect(bullet.Collider()) {
				meteorHit[i] = true
				if !bullet.piercing {
					bulletHit[j] = true
				}
				g.Score++
				g.recordHit(false) // Phase 1: stats tracking (meteor = not enemy)
				if bullet.piercing {
					g.Stats.RecordPiercingHit()
				}
				if meteor.canSplit {
					count := 2 + rand.Intn(3)
					for k := 0; k < count; k++ {
						split := newSplitMeteor(meteor.position, meteor.tier)
						if split != nil {
							newMeteors = append(newMeteors, split)
						}
					}
				}
			}
		}
	}

	survivingMeteors := make([]*Meteor, 0, len(g.meteors.value))
	for i, m := range g.meteors.value {
		if !meteorHit[i] {
			survivingMeteors = append(survivingMeteors, m)
		}
	}
	survivingMeteors = append(survivingMeteors, newMeteors...)
	g.meteors.value = survivingMeteors

	survivingBullets := make([]*Bullet, 0, len(g.player.bullets))
	for j, b := range g.player.bullets {
		if !bulletHit[j] {
			survivingBullets = append(survivingBullets, b)
		}
	}
	g.player.bullets = survivingBullets
}

func (g *Game) cleanupOffScreen() {
	pb := g.player.bullets[:0]
	for _, b := range g.player.bullets {
		if !b.IsOffScreen() {
			pb = append(pb, b)
		}
	}
	g.player.bullets = pb

	ms := g.meteors.value[:0]
	for _, m := range g.meteors.value {
		if !m.IsOffScreen() {
			ms = append(ms, m)
		}
	}
	g.meteors.value = ms

	ps := g.powerUps[:0]
	for _, pu := range g.powerUps {
		if !pu.IsOffScreen() {
			ps = append(ps, pu)
		}
	}
	g.powerUps = ps
}

// Phase 3: Enemy spawning and AI
func (g *Game) updateEnemies() {
	if g.enemies == nil {
		return
	}
	// Spawn new enemies
	enemyInterval := 8*time.Second - time.Duration(g.DifficultyLevel)*500*time.Millisecond
	if enemyInterval < 1*time.Second {
		enemyInterval = 1 * time.Second
	}
	g.enemies.SetSpawnInterval(enemyInterval)
	g.enemies.TrySpawn(g.DifficultyLevel)

	// Update enemies and have them shoot at player
	g.enemies.Update()
	px, py := g.player.playPosition.X+float64(g.player.sprite.Bounds().Dx())/2,
		g.player.playPosition.Y+float64(g.player.sprite.Bounds().Dy())/2
	for _, e := range g.enemies.value {
		e.Shoot(px, py)
	}
}

// Phase 3: Boss spawning and AI
func (g *Game) updateBoss() {
	if g.bossMgr == nil {
		return
	}
	// Try spawning boss
	if !g.bossMgr.IsActive() {
		g.bossMgr.bossSpawnTimer.UpdateTicks()
		if g.bossMgr.bossSpawnTimer.IsReadyAttack() {
			g.bossMgr.bossSpawnTimer.RestTicks()
			g.bossMgr.SpawnBoss(g.DifficultyLevel)
			g.bossActive = true
		}
	} else {
		g.bossMgr.Update()
		if g.bossMgr.currentBoss != nil && g.bossMgr.currentBoss.health <= 0 {
			// Boss defeated!
			if g.Stats != nil {
				g.Stats.RecordBossKill()
			}
			if g.Achievements != nil {
				g.Achievements.Unlock(AchKillBoss)
			}
			if g.Audio != nil {
				g.Audio.PlayExplosion()
			}
			g.Score += g.bossMgr.currentBoss.scoreValue
			// Drop power-ups
			for i := 0; i < 3; i++ {
				p := newPowerUpAt(g.bossMgr.currentBoss.position.X+float64(i-1)*100,
					g.bossMgr.currentBoss.position.Y)
				g.powerUps = append(g.powerUps, p)
			}
			// Add explosion effects
			if g.effects != nil {
				g.effects.AddExplosion(g.bossMgr.currentBoss.position.X,
					g.bossMgr.currentBoss.position.Y, 2.0)
			}
			g.bossMgr.Clear()
			g.bossActive = false
		}
	}
}

// Phase 3: Enemy vs player collisions
func (g *Game) enemyCollisions() {
	if g.enemies == nil || g.godMode {
		return
	}
	if g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack() {
		return
	}
	pRect := g.player.Collider()
	hitIdx := -1
	for i, e := range g.enemies.value {
		if e.Collider().Intersect(pRect) {
			hitIdx = i
			break
		}
	}
	if hitIdx >= 0 {
		g.enemies.value = append(g.enemies.value[:hitIdx], g.enemies.value[hitIdx+1:]...)
		g.Lives--
		g.recordDamage()
		if g.effects != nil {
			e := g.enemies.value[hitIdx] // already removed but this is the one that hit
			_ = e
		}
		if g.Lives <= 0 {
			g.State = StateGameOver
		} else {
			g.invulnTimer = NewTimer(InvulnDuration)
		}
	}
}

// Phase 3: Player bullets vs enemies
func (g *Game) bulletEnemyCollisions() {
	if g.enemies == nil || len(g.enemies.value) == 0 || len(g.player.bullets) == 0 {
		return
	}
	enemyHit := make([]bool, len(g.enemies.value))
	bulletHit := make([]bool, len(g.player.bullets))

	for i, e := range g.enemies.value {
		eRect := e.Collider()
		for j, b := range g.player.bullets {
			if !enemyHit[i] && !bulletHit[j] && eRect.Intersect(b.Collider()) {
				if !b.piercing {
					bulletHit[j] = true
				}
				damage := 1
				if g.player.weapon != nil {
					damage = g.player.weapon.damage
				}
				if e.TakeDamage(damage) {
					enemyHit[i] = true
					g.Score += e.scoreValue
					g.recordHit(true)
					if g.effects != nil {
						g.effects.AddExplosion(e.position.X, e.position.Y, 1.0)
					}
					if g.Audio != nil {
						g.Audio.PlayExplosion()
					}
				}
				if b.piercing {
					g.Stats.RecordPiercingHit()
				}
			}
		}
	}

	// Remove hit enemies and bullets
	survivingEnemies := make([]*Enemy, 0, len(g.enemies.value))
	for i, e := range g.enemies.value {
		if !enemyHit[i] {
			survivingEnemies = append(survivingEnemies, e)
		}
	}
	g.enemies.value = survivingEnemies

	survivingBullets := make([]*Bullet, 0, len(g.player.bullets))
	for j, b := range g.player.bullets {
		if !bulletHit[j] {
			survivingBullets = append(survivingBullets, b)
		}
	}
	g.player.bullets = survivingBullets
}

// Phase 3: Player bullets vs boss
func (g *Game) bossBulletCollisions() {
	if g.bossMgr == nil || !g.bossMgr.IsActive() || g.bossMgr.currentBoss == nil {
		return
	}
	boss := g.bossMgr.currentBoss
	if boss.sprite == nil {
		return
	}
	bossRect := newRect(boss.position, boss.sprite)
	bulletHit := make([]bool, len(g.player.bullets))

	for j, b := range g.player.bullets {
		if bossRect.Intersect(b.Collider()) {
			if !b.piercing {
				bulletHit[j] = true
			}
			damage := 1
			if g.player.weapon != nil {
				damage = g.player.weapon.damage
			}
			boss.TakeDamage(damage)
			g.recordHit(false)
		}
	}

	survivingBullets := make([]*Bullet, 0, len(g.player.bullets))
	for j, b := range g.player.bullets {
		if !bulletHit[j] {
			survivingBullets = append(survivingBullets, b)
		}
	}
	g.player.bullets = survivingBullets
}

// Phase 3: Enemy bullets vs player
func (g *Game) enemyBulletCollisions() {
	if g.godMode {
		return
	}
	if g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack() {
		return
	}
	// Shield: absorb bullets but destroy enemy bullets
	if g.player.shieldActive {
		// With shield, enemy bullets are destroyed on contact
		// (handled by iterating and removing)
	}

	pRect := g.player.Collider()
	hit := false

	// Check enemy bullets
	if g.enemies != nil {
		for _, e := range g.enemies.value {
			surviving := e.bullets[:0]
			for _, b := range e.bullets {
				if !hit && b.Collider().Intersect(pRect) {
					hit = true
					continue // remove this bullet
				}
				surviving = append(surviving, b)
			}
			e.bullets = surviving
		}
	}

	// Check boss bullets
	if g.bossMgr != nil && g.bossMgr.IsActive() && g.bossMgr.currentBoss != nil {
		surviving := g.bossMgr.currentBoss.bullets[:0]
		for _, b := range g.bossMgr.currentBoss.bullets {
			if !hit && b.Collider().Intersect(pRect) {
				hit = true
				continue
			}
			surviving = append(surviving, b)
		}
		g.bossMgr.currentBoss.bullets = surviving
	}

	if hit {
		g.Lives--
		g.recordDamage()
		if g.Lives <= 0 {
			g.State = StateGameOver
			if g.Audio != nil {
				g.Audio.PlayLose()
			}
		} else {
			g.invulnTimer = NewTimer(InvulnDuration)
		}
	}
}

// Phase 3: Shield update
func (g *Game) updateShield() {
	if !g.player.shieldActive {
		return
	}
	g.player.shieldTimer.UpdateTicks()
	if g.player.shieldTimer.IsReadyAttack() {
		g.player.shieldActive = false
		g.player.shieldTimer = nil
		if g.Audio != nil {
			g.Audio.PlayShieldDown()
		}
	}
}

// Phase 3: Speed boost update
func (g *Game) updateSpeedBoost() {
	if !g.player.speedBoostActive {
		return
	}
	g.player.speedBoostTimer.UpdateTicks()
	if g.player.speedBoostTimer.IsReadyAttack() {
		g.player.speedBoostActive = false
		g.player.speedBoostTimer = nil
		g.player.speedMultiplier = 1.0
	}
}

func (g *Game) reset() {
	g.player = newPlayer()
	g.meteors = newMeteors()
	g.powerUps = nil
	g.Score = 0
	g.Lives = InitialLives
	g.DifficultyLevel = 0
	g.firstSpawnDone = false
	g.invulnTimer = nil
	g.flashCounter = 0
	g.confirmQuit = false
	g.elapsedTime = 0
	g.powerUpSpawnTimer = NewTimer(PowerUpSpawnPeriod)
	g.lastBronzeMilestone = 10
	g.State = StatePlaying
	// Phase 1: reset stats
	if g.Stats != nil {
		g.Stats.Reset()
	}
	g.godMode = false
	g.showHitboxes = false
	g.slowMotion = false
	g.bossActive = false
	// Phase 3: reset game systems
	g.enemies = NewEnemyManager()
	g.bossMgr = NewBossManager()
	g.effects = NewEffectManager()
	g.background = NewBackground()
}

func (g *Game) Draw(s *ebiten.Image) {
	switch g.State {
	case StateStart:
		g.drawStartScreen(s)
	case StateGameOver:
		g.drawPlayField(s)
		g.drawGameOverOverlay(s)
	case StatePlaying, StatePaused:
		g.drawPlayField(s)
		g.drawPowerUps(s)
		if g.State == StatePaused {
			g.drawPauseOverlay(s)
		}
		if g.confirmQuit {
			g.drawConfirmOverlay(s)
		}
	}
}

func (g *Game) drawStartScreen(s *ebiten.Image) {
	text.Draw(s, "太空射击", UIFont, ScreenWidth/2-120, ScreenHeight/2-120, color.White)
	text.Draw(s, "按 回车 或 空格 开始游戏", UIFont, ScreenWidth/2-200, ScreenHeight/2-30, color.White)
	text.Draw(s, "WASD: 移动 | 鼠标左/右键: 旋转 | 空格: 射击 | P: 暂停 | ESC: 返回菜单", UIFont, ScreenWidth/2-460, ScreenHeight/2+60, color.RGBA{180, 180, 180, 255})
}

func (g *Game) drawPlayField(s *ebiten.Image) {
	// Phase 3: Draw background
	if g.background != nil {
		g.background.Draw(s)
	}

	isInvuln := g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack()
	drawPlayer := !isInvuln || (g.flashCounter/6)%2 == 0

	if drawPlayer {
		g.player.Draw(s)
	}
	g.meteors.Draw(s)

	// Phase 3: Draw enemies
	if g.enemies != nil {
		for _, e := range g.enemies.value {
			e.Draw(s)
		}
	}
	// Phase 3: Draw boss
	if g.bossMgr != nil {
		g.bossMgr.Draw(s)
	}
	// Phase 3: Draw effects
	if g.effects != nil {
		g.effects.Draw(s)
	}

	text.Draw(s, fmt.Sprintf("%06d", g.Score), ScoreFont, ScreenWidth/2-100, 50, color.White)
	text.Draw(s, fmt.Sprintf("生命: %d", g.Lives), UIFont, 50, 50, color.White)
	text.Draw(s, fmt.Sprintf("等级: %d", g.DifficultyLevel), UIFont, ScreenWidth-200, 50, color.White)
	text.Draw(s, fmt.Sprintf("时间: %.0fs", g.elapsedTime.Seconds()), UIFont, ScreenWidth/2-100, 80, color.RGBA{180, 180, 180, 255})

	// Phase 3: HUD sprite displays
	if g.hud != nil {
		g.hud.DrawLives(s, g.Lives, 50, 50)
	}

	buffY := 110
	if g.player.piercingActive {
		text.Draw(s, "穿透", UIFont, 50, buffY, color.RGBA{255, 215, 0, 255})
		buffY += 30
	}
	if g.player.extraSpread > 0 {
		text.Draw(s, fmt.Sprintf("散射 +%d", g.player.extraSpread), UIFont, 50, buffY, color.RGBA{205, 127, 50, 255})
		buffY += 30
	}
	if g.player.cdMultiplier < 1.0 {
		text.Draw(s, fmt.Sprintf("射速 x%.0f", 1.0/g.player.cdMultiplier), UIFont, 50, buffY, color.RGBA{192, 192, 192, 255})
		buffY += 30
	}
	if g.player.shieldActive {
		text.Draw(s, "护盾", UIFont, 50, buffY, color.RGBA{0, 200, 255, 255})
		buffY += 30
	}
	if g.player.speedBoostActive {
		text.Draw(s, "加速", UIFont, 50, buffY, color.RGBA{0, 255, 100, 255})
	}

	// Phase 1: debug overlays
	if g.DebugMode && g.hud != nil {
		g.hud.DrawDebugOverlay(s, g)
	}
	if g.showHitboxes && g.hud != nil {
		g.hud.DrawHitboxes(s, g)
	}
	// Phase 1: achievement popup
	if g.achievementPopup && g.achievementTimer != nil {
		// Show popup briefly (3 seconds = 180 ticks)
		g.achievementTimer.UpdateTicks()
		if !g.achievementTimer.IsReadyAttack() {
			ach := g.Achievements.Achievements[g.lastAchievementID]
			if g.hud != nil {
				g.hud.DrawAchievementPopup(s, ach)
			}
		} else {
			g.achievementPopup = false
			g.achievementTimer = nil
		}
	}
}

func (g *Game) drawPowerUps(s *ebiten.Image) {
	for _, pu := range g.powerUps {
		pu.Draw(s)
	}
}

func (g *Game) drawGameOverOverlay(s *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	s.DrawImage(overlay, nil)

	text.Draw(s, "游戏结束", UIFont, ScreenWidth/2-100, ScreenHeight/2-120, color.RGBA{255, 50, 50, 255})
	text.Draw(s, fmt.Sprintf("得分: %d  等级: %d", g.Score, g.DifficultyLevel), UIFont, ScreenWidth/2-180, ScreenHeight/2-40, color.White)
	text.Draw(s, "按 回车 或 空格 重新开始", UIFont, ScreenWidth/2-200, ScreenHeight/2+40, color.White)
}

func (g *Game) drawPauseOverlay(s *ebiten.Image) {
	text.Draw(s, "已暂停", UIFont, ScreenWidth/2-70, ScreenHeight/2-40, color.White)
}

func (g *Game) drawConfirmOverlay(s *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 140})
	s.DrawImage(overlay, nil)

	text.Draw(s, "确认返回菜单？", UIFont, ScreenWidth/2-130, ScreenHeight/2-50, color.White)
	text.Draw(s, "Enter 确认  |  ESC 取消", UIFont, ScreenWidth/2-200, ScreenHeight/2+10, color.RGBA{220, 220, 220, 255})
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("太空射击")
	g := &Game{
		player:            newPlayer(),
		meteors:           newMeteors(),
		State:             StateStart,
		powerUpSpawnTimer: NewTimer(PowerUpSpawnPeriod),
		Config:            DefaultConfig(),
		HighScores:        NewHighScoreManager(),
		Achievements:      NewAchievementTracker(),
		Stats:             &GameStats{},
		SaveMgr:           NewSaveManager(),
		Audio:             NewAudioManager(0.7, 0.4),
	}
	// Load persisted data
	g.HighScores.Load()
	g.Achievements.Load()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
