package main

import (
	"testing"
	"time"
)

func TestGameStatsNew(t *testing.T) {
	s := &GameStats{}
	if s.Kills != 0 || s.ShotsFired != 0 || s.ShotsHit != 0 {
		t.Error("new GameStats should be zero")
	}
}

func TestGameStatsRecordKill(t *testing.T) {
	s := &GameStats{}
	s.RecordKill(false) // meteor kill
	s.RecordKill(true)  // enemy kill
	if s.Kills != 2 {
		t.Errorf("kills want 2, got %d", s.Kills)
	}
	if s.MeteorsDestroyed != 1 {
		t.Errorf("meteors want 1, got %d", s.MeteorsDestroyed)
	}
	if s.EnemiesDestroyed != 1 {
		t.Errorf("enemies want 1, got %d", s.EnemiesDestroyed)
	}
}

func TestGameStatsCombo(t *testing.T) {
	s := &GameStats{}
	s.RecordKill(false)
	if s.CurrentCombo != 1 {
		t.Errorf("first kill combo want 1, got %d", s.CurrentCombo)
	}
	s.RecordKill(false) // quick kill within window
	if s.CurrentCombo != 2 {
		t.Errorf("rapid kill combo want 2, got %d", s.CurrentCombo)
	}
	if s.MaxCombo != 2 {
		t.Errorf("max combo want 2, got %d", s.MaxCombo)
	}
	// Simulate gap
	s.LastKillTime = time.Now().Add(-3 * time.Second)
	s.RecordKill(false)
	if s.CurrentCombo != 1 {
		t.Errorf("after gap, combo should reset to 1, got %d", s.CurrentCombo)
	}
}

func TestGameStatsRecordShotHitBoss(t *testing.T) {
	s := &GameStats{}
	s.RecordShot()
	s.RecordShot()
	s.RecordHit()
	if s.ShotsFired != 2 {
		t.Errorf("shots fired want 2, got %d", s.ShotsFired)
	}
	if s.ShotsHit != 1 {
		t.Errorf("shots hit want 1, got %d", s.ShotsHit)
	}
	s.RecordBossKill()
	if s.BossKills != 1 {
		t.Errorf("boss kills want 1, got %d", s.BossKills)
	}
}

func TestGameStatsAccuracy(t *testing.T) {
	s := &GameStats{}
	if s.Accuracy() != 0 {
		t.Error("accuracy with 0 shots should be 0")
	}
	s.ShotsFired = 10
	s.ShotsHit = 7
	if !floatApprox(s.Accuracy(), 0.7, 0.001) {
		t.Errorf("accuracy want 0.7, got %f", s.Accuracy())
	}
}

func TestGameStatsRecordDamage(t *testing.T) {
	s := &GameStats{}
	s.RecordDamage()
	s.RecordDamage()
	if s.DamageTaken != 2 {
		t.Errorf("damage taken want 2, got %d", s.DamageTaken)
	}
}

func TestGameStatsRecordPowerUp(t *testing.T) {
	s := &GameStats{}
	s.RecordPowerUp()
	s.RecordPowerUp()
	s.RecordPowerUp()
	if s.PowerUpsCollected != 3 {
		t.Errorf("powerups want 3, got %d", s.PowerUpsCollected)
	}
}

func TestGameStatsReset(t *testing.T) {
	s := &GameStats{
		Kills:       10,
		ShotsFired:  100,
		ShotsHit:    50,
		MaxCombo:    5,
	}
	s.Reset()
	if s.Kills != 0 || s.ShotsFired != 0 {
		t.Error("reset should zero all fields")
	}
}

func TestGameStatsPiercingHit(t *testing.T) {
	s := &GameStats{}
	s.RecordPiercingHit()
	s.RecordPiercingHit()
	if s.PiercingHits != 2 {
		t.Errorf("piercing hits want 2, got %d", s.PiercingHits)
	}
}
