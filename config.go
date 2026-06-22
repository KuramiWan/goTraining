package main

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds all tunable game parameters. Exported fields support
// JSON round-tripping and direct test access.
type Config struct {
	Player     PlayerConfig     `json:"player"`
	Meteor     MeteorConfig     `json:"meteor"`
	Enemy      EnemyConfig      `json:"enemy"`
	Boss       BossConfig       `json:"boss"`
	PowerUp    PowerUpConfig    `json:"powerup"`
	Difficulty DifficultyConfig `json:"difficulty"`
	Audio      AudioConfig      `json:"audio"`
	Debug      DebugConfig      `json:"debug"`
}

type PlayerConfig struct {
	Speed              float64       `json:"speed"`
	Acceleration       float64       `json:"acceleration"`
	Friction           float64       `json:"friction"`
	MaxLives           int           `json:"max_lives"`
	InvulnDuration     time.Duration `json:"invuln_duration"`
	BaseFireRate       time.Duration `json:"base_fire_rate"`
	MinFireRate        time.Duration `json:"min_fire_rate"`
	ShieldDuration     time.Duration `json:"shield_duration"`
	SpeedBoostDuration time.Duration `json:"speed_boost_duration"`
	SpeedBoostMult     float64       `json:"speed_boost_mult"`
}

type MeteorConfig struct {
	BaseSpawnInterval time.Duration `json:"base_spawn_interval"`
	MinSpawnInterval  time.Duration `json:"min_spawn_interval"`
	BaseSpeed         float64       `json:"base_speed"`
	SpeedVariance     float64       `json:"speed_variance"`
	MaxTier           int           `json:"max_tier"` // 0-3, 3 = tiny (no split)
}

type EnemyConfig struct {
	BaseSpawnInterval  time.Duration `json:"base_spawn_interval"`
	MinSpawnInterval   time.Duration `json:"min_spawn_interval"`
	MaxOnScreen        int           `json:"max_on_screen"`
	BaseHealth         int           `json:"base_health"`
	BulletSpeed        float64       `json:"bullet_speed"`
	ShootInterval      time.Duration `json:"shoot_interval"`
	ScoreValue         int           `json:"score_value"`
	UnlockDifficulty   map[string]int `json:"unlock_difficulty"` // enemy type -> min difficulty
}

type BossConfig struct {
	SpawnInterval       time.Duration `json:"spawn_interval"`
	WarningDuration     time.Duration `json:"warning_duration"`
	BaseHealth          int           `json:"base_health"`
	HealthScale         int           `json:"health_scale"` // per difficulty level
	ScoreValue          int           `json:"score_value"`
	PowerUpDropCount    int           `json:"powerup_drop_count"`
	Phase2HealthPct     float64       `json:"phase2_health_pct"` // 0.5 = 50%
	Phase3HealthPct     float64       `json:"phase3_health_pct"`
}

type PowerUpConfig struct {
	SpawnPeriod  time.Duration `json:"spawn_period"`
	SpeedMin     float64       `json:"speed_min"`
	SpeedMax     float64       `json:"speed_max"`
	BronzeMilestoneBase int    `json:"bronze_milestone_base"`
	DropChanceSilver    int    `json:"drop_chance_silver"` // out of 100
	DropChanceBronze    int    `json:"drop_chance_bronze"`
	DropChanceGold      int    `json:"drop_chance_gold"`
	DropChanceShield    int    `json:"drop_chance_shield"`
	DropChanceHealth    int    `json:"drop_chance_health"`
	DropChanceSpeed     int    `json:"drop_chance_speed"`
	HealthOnlyBelowLives int   `json:"health_only_below_lives"` // only spawn health if lives < this
}

type DifficultyConfig struct {
	ScorePerLevel      int           `json:"score_per_level"`
	TimePerLevel       time.Duration `json:"time_per_level"`
	SpawnRateStep      time.Duration `json:"spawn_rate_step"`
	FireRateStep       time.Duration `json:"fire_rate_step"`
	EnemyBulletSpeedUp float64       `json:"enemy_bullet_speed_up"` // per level
}

type AudioConfig struct {
	SFXVolume  float64 `json:"sfx_volume"`  // 0.0 - 1.0
	BGMVolume  float64 `json:"bgm_volume"`  // 0.0 - 1.0
	Muted      bool    `json:"muted"`
	Enabled    bool    `json:"enabled"`
}

type DebugConfig struct {
	ShowFPS     bool `json:"show_fps"`
	ShowHitboxes bool `json:"show_hitboxes"`
	VerboseLog  bool `json:"verbose_log"`
}

// DefaultConfig returns a Config with all the current hardcoded game values.
func DefaultConfig() *Config {
	return &Config{
		Player: PlayerConfig{
			Speed:              0,
			Acceleration:       1.2,
			Friction:           0.8,
			MaxLives:           5,
			InvulnDuration:     2 * time.Second,
			BaseFireRate:       1 * time.Second,
			MinFireRate:        200 * time.Millisecond,
			ShieldDuration:     10 * time.Second,
			SpeedBoostDuration: 15 * time.Second,
			SpeedBoostMult:     1.5,
		},
		Meteor: MeteorConfig{
			BaseSpawnInterval: 5 * time.Second,
			MinSpawnInterval:  500 * time.Millisecond,
			BaseSpeed:         1.0,
			SpeedVariance:     1.5,
			MaxTier:           3,
		},
		Enemy: EnemyConfig{
			BaseSpawnInterval: 8 * time.Second,
			MinSpawnInterval:  1 * time.Second,
			MaxOnScreen:       10,
			BaseHealth:        3,
			BulletSpeed:       400,
			ShootInterval:     2 * time.Second,
			ScoreValue:        5,
			UnlockDifficulty: map[string]int{
				"black": 0,
				"blue":  3,
				"green": 6,
				"red":   10,
			},
		},
		Boss: BossConfig{
			SpawnInterval:    90 * time.Second,
			WarningDuration:  3 * time.Second,
			BaseHealth:       50,
			HealthScale:      10,
			ScoreValue:       100,
			PowerUpDropCount: 3,
			Phase2HealthPct:  0.5,
			Phase3HealthPct:  0.25,
		},
		PowerUp: PowerUpConfig{
			SpawnPeriod:          60 * time.Second,
			SpeedMin:             0.8,
			SpeedMax:             2.0,
			BronzeMilestoneBase:  10,
			DropChanceSilver:     35,
			DropChanceBronze:     25,
			DropChanceGold:       10,
			DropChanceShield:     15,
			DropChanceHealth:     5,
			DropChanceSpeed:      10,
			HealthOnlyBelowLives: 3,
		},
		Difficulty: DifficultyConfig{
			ScorePerLevel:      10,
			TimePerLevel:       30 * time.Second,
			SpawnRateStep:      500 * time.Millisecond,
			FireRateStep:       100 * time.Millisecond,
			EnemyBulletSpeedUp: 20,
		},
		Audio: AudioConfig{
			SFXVolume: 0.7,
			BGMVolume: 0.4,
			Muted:     false,
			Enabled:   true,
		},
		Debug: DebugConfig{
			ShowFPS:      false,
			ShowHitboxes: false,
			VerboseLog:   false,
		},
	}
}

// LoadConfig reads a JSON config file. Missing keys fall back to defaults.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // use defaults
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// SaveConfig writes the current config to a JSON file.
func (c *Config) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
