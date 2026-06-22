package main

import (
	"testing"
)

func TestPowerUpNew(t *testing.T) {
	pu := newPowerUp(PowerUpSilver)
	if pu == nil {
		t.Fatal("newPowerUp returned nil")
	}
	if pu.pType != PowerUpSilver {
		t.Errorf("want Silver, got %d", pu.pType)
	}
	if pu.sprite == nil {
		t.Error("power-up should have a sprite")
	}
	// Should be outside or on edge of screen
	// (spawned from edges)
}

func TestPowerUpNewShield(t *testing.T) {
	pu := newPowerUp(PowerUpShield)
	if pu.pType != PowerUpShield {
		t.Errorf("want Shield, got %d", pu.pType)
	}
}

func TestPowerUpNewHealth(t *testing.T) {
	pu := newPowerUp(PowerUpHealth)
	if pu.pType != PowerUpHealth {
		t.Errorf("want Health, got %d", pu.pType)
	}
}

func TestPowerUpNewSpeed(t *testing.T) {
	pu := newPowerUp(PowerUpSpeed)
	if pu.pType != PowerUpSpeed {
		t.Errorf("want Speed, got %d", pu.pType)
	}
}

func TestPowerUpAt(t *testing.T) {
	pu := newPowerUpAt(500, 400)
	if pu == nil {
		t.Fatal("newPowerUpAt returned nil")
	}
	if pu.position.X != 500 || pu.position.Y != 400 {
		t.Errorf("position mismatch: (%f,%f)", pu.position.X, pu.position.Y)
	}
}

func TestPowerUpUpdate(t *testing.T) {
	pu := newPowerUpAt(500, 400)
	origX := pu.position.X
	origY := pu.position.Y
	pu.Update()
	// Should have moved
	if pu.position.X == origX && pu.position.Y == origY {
		t.Error("power-up should move after update")
	}
}

func TestPowerUpCollider(t *testing.T) {
	pu := newPowerUp(PowerUpSilver)
	r := pu.Collider()
	if r.Width <= 0 || r.Height <= 0 {
		t.Error("power-up collider should have positive dimensions")
	}
}

func TestPowerUpIsOffScreen(t *testing.T) {
	pu := newPowerUpAt(500, 400)
	if pu.IsOffScreen() {
		t.Error("power-up at center should not be off screen")
	}
	pu.position.X = -200
	if !pu.IsOffScreen() {
		t.Error("power-up far left should be off screen")
	}
}

func TestGetPowerUpSprite(t *testing.T) {
	types := []PowerUpType{PowerUpSilver, PowerUpBronze, PowerUpGold, PowerUpShield, PowerUpHealth, PowerUpSpeed}
	for _, typ := range types {
		s := getPowerUpSprite(typ)
		if s == nil {
			t.Errorf("power-up type %d returned nil sprite", typ)
		}
	}
}
