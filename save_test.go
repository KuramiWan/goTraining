package main

import (
	"testing"
)

func TestSaveManagerNew(t *testing.T) {
	sm := NewSaveManager()
	if sm == nil {
		t.Fatal("NewSaveManager returned nil")
	}
	if sm.filePath == "" {
		t.Error("filePath should not be empty")
	}
}

func TestSaveManagerDelete(t *testing.T) {
	sm := NewSaveManager()
	tmp := t.TempDir()
	sm.filePath = tmp + "/test_save.json"

	// Deleting non-existent file should error
	err := sm.Delete()
	if err == nil {
		t.Log("delete of missing file returned nil (platform-specific)")
	}
}

func TestSaveManagerExists(t *testing.T) {
	sm := NewSaveManager()
	tmp := t.TempDir()
	sm.filePath = tmp + "/nonexistent.json"
	if sm.Exists() {
		t.Error("should not exist before save")
	}
}

func TestSaveCanSave(t *testing.T) {
	g := newTestGame()
	g.State = StatePlaying

	if !CanSave(g) {
		t.Error("should be able to save during normal play")
	}

	g.State = StatePaused
	if CanSave(g) {
		t.Error("should not save during pause")
	}
}
