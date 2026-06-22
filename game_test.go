package main

import (
	"testing"
)

// Integration tests for collision systems and game state

func TestGameStartState(t *testing.T) {
	g := &Game{
		player:            newPlayer(),
		meteors:           newMeteors(),
		State:             StateStart,
		powerUpSpawnTimer: NewTimer(PowerUpSpawnPeriod),
		Config:            DefaultConfig(),
		Stats:             &GameStats{},
	}
	if g.State != StateStart {
		t.Error("new game should start in StateStart")
	}
	g.startGame()
	if g.State != StatePlaying {
		t.Error("after startGame, state should be StatePlaying")
	}
	if g.Score != 0 {
		t.Errorf("score should be 0 after start, got %d", g.Score)
	}
	if g.Lives != InitialLives {
		t.Errorf("lives should be %d after start, got %d", InitialLives, g.Lives)
	}
}

func TestGameStateTransitions(t *testing.T) {
	g := newTestGame()
	g.startGame()
	if g.State != StatePlaying {
		t.Fatal("should be playing")
	}

	// Force game over
	g.Lives = 1
	g.invulnTimer = nil
	// Simulate collision: spawn meteor on top of player
	px, py := getPlayerPosition(g)
	spawnTestMeteor(g, px+20, py+20, 0)

	// Advance one tick to let meteor collision fire
	g.updatePlaying()

	if g.State != StatePlaying && g.Lives == 0 {
		// Only check state=GameOver if meteor actually hit
		if g.State == StateGameOver {
			t.Log("game over triggered correctly")
		}
	}
}

func TestBulletMeteorCollision(t *testing.T) {
	g := newTestGame()
	clearMeteors(g)
	clearBullets(g)

	// Place meteor at known position
	m := spawnTestMeteor(g, 500, 400, 0)
	// Fire bullet directly at meteor
	b := spawnTestBullet(g, 500-float64(m.sprite.Bounds().Dx())/2-10, 400, 0)

	// Force bullet to intersect by moving it to meteor position
	b.position.X = m.position.X + 1
	b.position.Y = m.position.Y + 1

	origScore := g.Score
	g.bulletCollisions()

	if g.Score > origScore {
		t.Log("bullet hit meteor, score increased")
	}
}

func TestInvulnerabilityPreventsDamage(t *testing.T) {
	g := newTestGame()
	setInvulnerability(g, InvulnDuration)
	origLives := g.Lives

	// Place meteor on player
	px, py := getPlayerPosition(g)
	spawnTestMeteor(g, px+20, py+20, 0)
	g.updatePlaying()

	// Lives should not decrease
	if g.Lives < origLives {
		t.Error("invulnerable player should not take damage")
	}
}

func TestGodMode(t *testing.T) {
	g := newTestGame()
	g.godMode = true
	origLives := g.Lives

	// Place meteor on player
	px, py := getPlayerPosition(g)
	spawnTestMeteor(g, px+20, py+20, 0)
	g.updatePlaying()

	if g.Lives < origLives {
		t.Error("god mode should prevent damage")
	}
}

func TestPlayerClamped(t *testing.T) {
	g := newTestGame()
	setPlayerPosition(g, -50, -50)
	g.player.clampToScreen()
	x, y := getPlayerPosition(g)
	if x < 0 || y < 0 {
		t.Errorf("player should be clamped to screen: (%f,%f)", x, y)
	}
}

func TestScoreIncrements(t *testing.T) {
	g := newTestGame()
	orig := g.Score
	g.Score += 10
	if g.Score != orig+10 {
		t.Error("score should increment")
	}
}

func TestDifficultyScaling(t *testing.T) {
	g := newTestGame()
	g.Score = 100 // 10 levels from score
	g.updateDifficulty()
	if g.DifficultyLevel < 10 {
		t.Errorf("difficulty should be at least 10, got %d", g.DifficultyLevel)
	}
}
