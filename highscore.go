package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxHighScoreEntries = 10

// HighScoreEntry represents a single high score record.
type HighScoreEntry struct {
	Name  string    `json:"name"`
	Score int       `json:"score"`
	Date  time.Time `json:"date"`
}

// HighScoreManager handles loading, saving, and querying high scores.
type HighScoreManager struct {
	Entries  []HighScoreEntry `json:"entries"`
	filePath string
}

// NewHighScoreManager creates a manager and loads existing scores from disk.
func NewHighScoreManager() *HighScoreManager {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	appDir := filepath.Join(dir, "space-shooter")
	os.MkdirAll(appDir, 0755)
	return &HighScoreManager{
		Entries:  make([]HighScoreEntry, 0, maxHighScoreEntries),
		filePath: filepath.Join(appDir, "highscores.json"),
	}
}

// Load reads high scores from the JSON file.
func (h *HighScoreManager) Load() error {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &h.Entries)
}

// Save writes high scores to the JSON file.
func (h *HighScoreManager) Save() error {
	data, err := json.MarshalIndent(h.Entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.filePath, data, 0644)
}

// IsHighScore checks if the given score qualifies for the high score list.
func (h *HighScoreManager) IsHighScore(score int) bool {
	if score <= 0 {
		return false
	}
	if len(h.Entries) < maxHighScoreEntries {
		return true
	}
	return score > h.Entries[len(h.Entries)-1].Score
}

// Add inserts a high score entry in sorted order (highest first) and trims to max.
func (h *HighScoreManager) Add(entry HighScoreEntry) {
	entry.Date = time.Now()
	h.Entries = append(h.Entries, entry)
	sort.Slice(h.Entries, func(i, j int) bool {
		return h.Entries[i].Score > h.Entries[j].Score
	})
	if len(h.Entries) > maxHighScoreEntries {
		h.Entries = h.Entries[:maxHighScoreEntries]
	}
}

// Top returns the top n entries.
func (h *HighScoreManager) Top(n int) []HighScoreEntry {
	if n > len(h.Entries) {
		n = len(h.Entries)
	}
	return h.Entries[:n]
}

// Rank returns the 1-based rank of a score, or 0 if not in the list.
func (h *HighScoreManager) Rank(score int) int {
	for i, e := range h.Entries {
		if e.Score == score {
			return i + 1
		}
	}
	return 0
}
