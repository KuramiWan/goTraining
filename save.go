package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// SaveData holds the serializable game state.
type SaveData struct {
	Score           int           `json:"score"`
	Lives           int           `json:"lives"`
	DifficultyLevel int           `json:"difficulty_level"`
	ElapsedTime     time.Duration `json:"elapsed_time"`
	Timestamp       time.Time     `json:"timestamp"`
}

// SaveManager handles saving and loading game state.
type SaveManager struct {
	filePath string
}

// NewSaveManager creates a SaveManager for the app's save directory.
func NewSaveManager() *SaveManager {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	appDir := filepath.Join(dir, "space-shooter")
	os.MkdirAll(appDir, 0755)
	return &SaveManager{
		filePath: filepath.Join(appDir, "save.json"),
	}
}

// Save writes the current game state to disk.
func (s *SaveManager) Save(g *Game) error {
	data := SaveData{
		Score:           g.Score,
		Lives:           g.Lives,
		DifficultyLevel: g.DifficultyLevel,
		ElapsedTime:     g.elapsedTime,
		Timestamp:       time.Now(),
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, bytes, 0644)
}

// Load reads the saved game state from disk.
func (s *SaveManager) Load() (*SaveData, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}
	var sd SaveData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

// Exists returns whether a save file exists.
func (s *SaveManager) Exists() bool {
	_, err := os.Stat(s.filePath)
	return err == nil
}

// Delete removes the save file.
func (s *SaveManager) Delete() error {
	return os.Remove(s.filePath)
}

// CanSave returns whether the game is in a savable state.
func CanSave(g *Game) bool {
	return g.State == StatePlaying &&
		(g.invulnTimer == nil || g.invulnTimer.IsReadyAttack()) &&
		!g.bossActive
}
