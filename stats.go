package main

import "time"

// GameStats tracks all game statistics. All fields exported for test access.
type GameStats struct {
	Kills            int           `json:"kills"`
	ShotsFired       int           `json:"shots_fired"`
	ShotsHit         int           `json:"shots_hit"`
	BossKills        int           `json:"boss_kills"`
	PowerUpsCollected int          `json:"powerups_collected"`
	DamageTaken      int           `json:"damage_taken"`
	PlayTime         time.Duration `json:"play_time"`
	MaxCombo         int           `json:"max_combo"`
	CurrentCombo     int           `json:"current_combo"`
	LastKillTime     time.Time     `json:"-"`
	PiercingHits     int           `json:"piercing_hits"`
	MeteorsDestroyed int           `json:"meteors_destroyed"`
	EnemiesDestroyed int           `json:"enemies_destroyed"`
}

// ComboWindow is the time within which a new kill extends the combo.
const ComboWindow = 2 * time.Second

// RecordKill records a kill and updates the combo counter.
func (s *GameStats) RecordKill(isEnemy bool) {
	s.Kills++
	if isEnemy {
		s.EnemiesDestroyed++
	} else {
		s.MeteorsDestroyed++
	}

	now := time.Now()
	if now.Sub(s.LastKillTime) < ComboWindow {
		s.CurrentCombo++
	} else {
		s.CurrentCombo = 1
	}
	s.LastKillTime = now
	if s.CurrentCombo > s.MaxCombo {
		s.MaxCombo = s.CurrentCombo
	}
}

// RecordShot increments the shots fired counter.
func (s *GameStats) RecordShot() {
	s.ShotsFired++
}

// RecordHit increments the shots hit counter.
func (s *GameStats) RecordHit() {
	s.ShotsHit++
}

// RecordPiercingHit records that a piercing bullet hit a target.
func (s *GameStats) RecordPiercingHit() {
	s.PiercingHits++
}

// RecordBossKill records a boss kill.
func (s *GameStats) RecordBossKill() {
	s.BossKills++
}

// RecordPowerUp records a power-up pickup.
func (s *GameStats) RecordPowerUp() {
	s.PowerUpsCollected++
}

// RecordDamage records damage taken.
func (s *GameStats) RecordDamage() {
	s.DamageTaken++
}

// Accuracy returns the hit ratio (0.0-1.0), or 0 if no shots fired.
func (s *GameStats) Accuracy() float64 {
	if s.ShotsFired == 0 {
		return 0
	}
	return float64(s.ShotsHit) / float64(s.ShotsFired)
}

// Reset zeroes all counters (called on new game).
func (s *GameStats) Reset() {
	*s = GameStats{}
}
