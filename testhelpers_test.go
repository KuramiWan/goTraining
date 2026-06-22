package main

import (
	"math"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Test helpers for setting up game state ---

// newTestGame creates a Game instance suitable for testing.
func newTestGame() *Game {
	g := &Game{
		player:            newPlayer(),
		meteors:           newMeteors(),
		State:             StatePlaying,
		Lives:             InitialLives,
		powerUpSpawnTimer: NewTimer(PowerUpSpawnPeriod),
		Config:            DefaultConfig(),
		Stats:             &GameStats{},
	}
	return g
}

// newTestGameWithConfig creates a Game with a custom config.
func newTestGameWithConfig(cfg *Config) *Game {
	g := newTestGame()
	g.Config = cfg
	return g
}

// setPlayerPosition places the player at specific coordinates.
func setPlayerPosition(g *Game, x, y float64) {
	g.player.playPosition.X = x
	g.player.playPosition.Y = y
}

// getPlayerPosition returns the player's current position.
func getPlayerPosition(g *Game) (float64, float64) {
	return g.player.playPosition.X, g.player.playPosition.Y
}

// spawnTestMeteor creates a meteor at a specific position with a given tier.
func spawnTestMeteor(g *Game, x, y float64, tier int) *Meteor {
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
	m := &Meteor{
		sprite:   sprite,
		position: Vector{X: x, Y: y},
		movement: Vector{X: 0, Y: 0},
		canSplit: tier < 3,
		tier:     tier,
	}
	g.meteors.value = append(g.meteors.value, m)
	return m
}

// clearMeteors removes all meteors from the game.
func clearMeteors(g *Game) {
	g.meteors.value = nil
}

// spawnTestBullet creates a bullet at a position with an angle.
func spawnTestBullet(g *Game, x, y, angle float64) *Bullet {
	b := newBullet(
		Vector{X: x, Y: y},
		angle,
		Vector{},
		false,
	)
	g.player.bullets = append(g.player.bullets, b)
	return b
}

// clearBullets removes all player bullets.
func clearBullets(g *Game) {
	g.player.bullets = nil
}

// setScore sets the game score to a specific value.
func setScore(g *Game, score int) {
	g.Score = score
}

// setLives sets the remaining lives.
func setLives(g *Game, lives int) {
	g.Lives = lives
}

// setInvulnerability makes the player invulnerable for a duration.
func setInvulnerability(g *Game, d time.Duration) {
	g.invulnTimer = NewTimer(d)
}

// clearInvulnerability removes invulnerability.
func clearInvulnerability(g *Game) {
	g.invulnTimer = nil
}

// advanceGameTicks simulates n ticks of game updates (for the playing state).
func advanceGameTicks(g *Game, n int) {
	for i := 0; i < n; i++ {
		g.updatePlaying()
	}
}

// fireBullet manually fires a bullet from the player position.
func fireBullet(g *Game) {
	pos := g.player.playPosition
	g.player.bullets = append(g.player.bullets, newBullet(
		Vector{X: pos.X, Y: pos.Y},
		0,
		Vector{},
		false,
	))
}

// newExpiredTimer creates a timer that is already ready.
func newExpiredTimer() *Timer {
	return &Timer{
		currentTicks: 9999,
		targetTicks:  10,
	}
}

// newTimerAt creates a timer with current ticks set.
func newTimerAt(ticks int) *Timer {
	t := &Timer{
		currentTicks: ticks,
		targetTicks:  ebiten.TPS(), // 1 second
	}
	return t
}

// floatApprox checks if two floats are approximately equal.
func floatApprox(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// Ensure imports are used - reference internal types.
var _ = (*testing.T).TempDir
var _ = time.Now
var _ = ebiten.TPS
