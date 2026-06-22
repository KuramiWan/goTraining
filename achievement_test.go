package main

import (
	"testing"
	"time"
)

func TestAchievementTrackerNew(t *testing.T) {
	tr := NewAchievementTracker()
	if tr == nil {
		t.Fatal("NewAchievementTracker returned nil")
	}
	if len(tr.Achievements) == 0 {
		t.Fatal("achievements should be pre-populated")
	}
	// All should be locked
	for _, ach := range tr.Achievements {
		if ach.Unlocked {
			t.Errorf("achievement %s should start locked", ach.Name)
		}
	}
}

func TestAchievementUnlock(t *testing.T) {
	tr := NewAchievementTracker()

	// Already locked
	if !tr.Unlock(AchFirstKill) {
		t.Fatal("first unlock should return true")
	}

	// Verify it's now unlocked
	if !tr.IsUnlocked(AchFirstKill) {
		t.Error("AchFirstKill should be unlocked")
	}

	// Double unlock should return false
	if tr.Unlock(AchFirstKill) {
		t.Error("double unlock should return false")
	}
}

func TestAchievementRecordKill(t *testing.T) {
	tr := NewAchievementTracker()
	tr.RecordKill(1)
	if !tr.IsUnlocked(AchFirstKill) {
		t.Error("first kill should unlock AchFirstKill")
	}

	tr.RecordKill(50)
	if !tr.IsUnlocked(AchKill50Enemies) {
		t.Error("50 kills should unlock AchKill50Enemies")
	}

	tr.RecordKill(100)
	if !tr.IsUnlocked(AchKill100Enemies) {
		t.Error("100 kills should unlock AchKill100Enemies")
	}
}

func TestAchievementCheckStats(t *testing.T) {
	tr := NewAchievementTracker()
	stats := &GameStats{}

	tr.CheckStats(stats, 100, 0, 0)
	if !tr.IsUnlocked(AchScore100) {
		t.Error("score 100 should unlock AchScore100")
	}

	tr.CheckStats(stats, 1000, 0, 0)
	if !tr.IsUnlocked(AchScore1000) {
		t.Error("score 1000 should unlock AchScore1000")
	}
}

func TestAchievementSurvival(t *testing.T) {
	tr := NewAchievementTracker()
	stats := &GameStats{}

	tr.CheckStats(stats, 0, 0, 61*time.Second) // 1 second = 1e9 ns
	if !tr.IsUnlocked(AchSurvive1Min) {
		t.Error("surviving 61s should unlock AchSurvive1Min")
	}
}

func TestAchievementMaxDifficulty(t *testing.T) {
	tr := NewAchievementTracker()
	stats := &GameStats{}

	tr.CheckStats(stats, 0, 20, 0)
	if !tr.IsUnlocked(AchMaxDifficulty) {
		t.Error("difficulty 20 should unlock AchMaxDifficulty")
	}
}

func TestAchievementRecordPowerUp(t *testing.T) {
	tr := NewAchievementTracker()

	// Record all power-up types
	tr.RecordPowerUpCollected(PowerUpSilver)
	tr.RecordPowerUpCollected(PowerUpBronze)
	tr.RecordPowerUpCollected(PowerUpGold)
	tr.RecordPowerUpCollected(PowerUpShield)
	tr.RecordPowerUpCollected(PowerUpHealth)
	tr.RecordPowerUpCollected(PowerUpSpeed)

	if !tr.IsUnlocked(AchCollectAllPowerUps) {
		t.Error("collecting all types should unlock AchCollectAllPowerUps")
	}
}

func TestAchievementBronzeMilestone(t *testing.T) {
	tr := NewAchievementTracker()
	for i := 0; i < 5; i++ {
		tr.RecordPowerUpCollected(PowerUpBronze)
	}
	if !tr.IsUnlocked(AchBronzeMilestone) {
		t.Error("5 bronze pickups should unlock AchBronzeMilestone")
	}
}

func TestAchievementRecordDamage(t *testing.T) {
	tr := NewAchievementTracker()
	tr.RecordDamage()
	if tr.lastDamageTime.IsZero() {
		t.Error("RecordDamage should set lastDamageTime")
	}
}
