package main

import (
	"testing"
	"time"
)

func TestVectorNormalize(t *testing.T) {
	v := Vector{X: 3, Y: 4}
	n := v.Normalize()
	if !floatApprox(n.X, 0.6, 0.001) {
		t.Errorf("normalized X want 0.6, got %f", n.X)
	}
	if !floatApprox(n.Y, 0.8, 0.001) {
		t.Errorf("normalized Y want 0.8, got %f", n.Y)
	}
}

func TestVectorNormalizeZero(t *testing.T) {
	v := Vector{X: 0, Y: 0}
	n := v.Normalize()
	if n.X != 0 || n.Y != 0 {
		t.Errorf("zero vector normalizing should return zero: (%f,%f)", n.X, n.Y)
	}
}

func TestPlayerNew(t *testing.T) {
	p := newPlayer()
	if p == nil {
		t.Fatal("newPlayer returned nil")
	}
	if p.sprite == nil {
		t.Error("player should have a sprite")
	}
	if p.weapon == nil {
		t.Error("player should have a weapon")
	}
	if p.speedMultiplier != 1.0 {
		t.Errorf("speedMult want 1.0, got %f", p.speedMultiplier)
	}
	if p.cdMultiplier != 1.0 {
		t.Errorf("cdMult want 1.0, got %f", p.cdMultiplier)
	}
}

func TestPlayerCollider(t *testing.T) {
	p := newPlayer()
	r := p.Collider()
	if r.Width <= 0 || r.Height <= 0 {
		t.Error("player collider should have positive dimensions")
	}
}

func TestPlayerSetFireRate(t *testing.T) {
	p := newPlayer()
	orig := p.coldTimer.targetTicks
	p.SetFireRate(500 * time.Millisecond)
	if p.coldTimer.targetTicks >= orig {
		t.Errorf("fire rate should decrease: %d >= %d", p.coldTimer.targetTicks, orig)
	}
}

func TestPlayerApplyPowerUpSilver(t *testing.T) {
	p := newPlayer()
	orig := p.cdMultiplier
	p.ApplyPowerUp(PowerUpSilver)
	if p.cdMultiplier >= orig {
		t.Errorf("silver should reduce cdMultiplier: %f >= %f", p.cdMultiplier, orig)
	}
}

func TestPlayerApplyPowerUpBronze(t *testing.T) {
	p := newPlayer()
	p.ApplyPowerUp(PowerUpBronze)
	if p.extraSpread != 1 {
		t.Errorf("bronze should add spread: want 1 got %d", p.extraSpread)
	}
}

func TestPlayerApplyPowerUpGold(t *testing.T) {
	p := newPlayer()
	p.ApplyPowerUp(PowerUpGold)
	if !p.piercingActive {
		t.Error("gold should activate piercing")
	}
}

func TestPlayerApplyPowerUpShield(t *testing.T) {
	p := newPlayer()
	p.ApplyPowerUp(PowerUpShield)
	if !p.shieldActive {
		t.Error("shield power-up should activate shield")
	}
	if p.shieldTimer == nil {
		t.Error("shield should have a timer")
	}
}

func TestPlayerApplyPowerUpHealth(t *testing.T) {
	p := newPlayer()
	p.ApplyPowerUp(PowerUpHealth)
	// Health is handled in game.go (g.Lives++), not in player
	// We just verify no crash
}

func TestPlayerApplyPowerUpSpeed(t *testing.T) {
	p := newPlayer()
	p.ApplyPowerUp(PowerUpSpeed)
	if !p.speedBoostActive {
		t.Error("speed power-up should activate speed boost")
	}
	if p.speedMultiplier != 1.5 {
		t.Errorf("speed mult want 1.5, got %f", p.speedMultiplier)
	}
}

func TestPlayerClampToScreen(t *testing.T) {
	p := newPlayer()
	// Move off-screen left
	p.playPosition.X = -100
	p.clampToScreen()
	if p.playPosition.X < 0 {
		t.Errorf("should clamp to 0, got %f", p.playPosition.X)
	}
	// Move off-screen bottom
	p.playPosition.Y = 2000
	p.clampToScreen()
	if p.playPosition.Y > ScreenHeight {
		t.Errorf("should clamp to screen height, got %f", p.playPosition.Y)
	}
}

func TestPlayerResetBuffs(t *testing.T) {
	p := newPlayer()
	p.ApplyPowerUp(PowerUpSilver)
	p.ApplyPowerUp(PowerUpBronze)
	p.ApplyPowerUp(PowerUpGold)
	p.ResetBuffs()
	if p.cdMultiplier != 1.0 || p.extraSpread != 0 || p.piercingActive {
		t.Error("ResetBuffs should clear all buffs")
	}
}
