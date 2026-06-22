package main

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	// Check key values match the original game constants
	if cfg.Player.Acceleration != 1.2 {
		t.Errorf("acceleration want 1.2 got %f", cfg.Player.Acceleration)
	}
	if cfg.Player.Friction != 0.8 {
		t.Errorf("friction want 0.8 got %f", cfg.Player.Friction)
	}
	if cfg.Player.MaxLives != 5 {
		t.Errorf("max lives want 5 got %d", cfg.Player.MaxLives)
	}
}

func TestConfigJSONRoundTrip(t *testing.T) {
	cfg := DefaultConfig()
	tmp := t.TempDir()
	path := tmp + "/test_config.json"

	if err := cfg.SaveConfig(path); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.Player.Acceleration != cfg.Player.Acceleration {
		t.Errorf("acceleration not preserved: %f vs %f", loaded.Player.Acceleration, cfg.Player.Acceleration)
	}
	if loaded.Player.Friction != cfg.Player.Friction {
		t.Errorf("friction not preserved")
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/test_config.json")
	if err != nil {
		t.Fatalf("LoadConfig should not error on missing file: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig should return defaults on missing file")
	}
	// Should be defaults
	if cfg.Player.Acceleration != 1.2 {
		t.Errorf("should have default acceleration")
	}
}

func TestConfigSaveLoad(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Difficulty.ScorePerLevel = 20 // non-default value
	cfg.Debug.ShowFPS = true

	tmp := t.TempDir()
	path := tmp + "/config.json"
	if err := cfg.SaveConfig(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	// The config system overwrites with loaded values
	if loaded.Difficulty.ScorePerLevel != 20 {
		t.Errorf("ScorePerLevel: want 20 got %d", loaded.Difficulty.ScorePerLevel)
	}
	if !loaded.Debug.ShowFPS {
		t.Error("ShowFPS should be true")
	}
}

func TestVolumeClamping(t *testing.T) {
	if clampVolume(-0.5) != 0 {
		t.Error("negative should clamp to 0")
	}
	if clampVolume(1.5) != 1 {
		t.Error(">1 should clamp to 1")
	}
	if clampVolume(0.5) != 0.5 {
		t.Error("valid value should pass through")
	}
}

// Avoid unused import
var _ = os.O_RDONLY
