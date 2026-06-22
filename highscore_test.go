package main

import (
	"testing"
)

func TestHighScoreManagerNew(t *testing.T) {
	h := NewHighScoreManager()
	if h == nil {
		t.Fatal("NewHighScoreManager returned nil")
	}
	if len(h.Entries) != 0 {
		t.Errorf("new manager should have 0 entries, got %d", len(h.Entries))
	}
}

func TestHighScoreIsHighScore(t *testing.T) {
	h := NewHighScoreManager()

	// Empty list: any positive score qualifies
	if !h.IsHighScore(100) {
		t.Error("100 should qualify for empty list")
	}
	if h.IsHighScore(0) {
		t.Error("0 should not qualify")
	}
	if h.IsHighScore(-1) {
		t.Error("negative should not qualify")
	}
}

func TestHighScoreAddAndSort(t *testing.T) {
	h := NewHighScoreManager()
	h.Add(HighScoreEntry{Name: "A", Score: 100})
	h.Add(HighScoreEntry{Name: "B", Score: 300})
	h.Add(HighScoreEntry{Name: "C", Score: 200})

	if len(h.Entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(h.Entries))
	}
	if h.Entries[0].Score != 300 {
		t.Errorf("first should be highest (300), got %d", h.Entries[0].Score)
	}
	if h.Entries[1].Score != 200 {
		t.Errorf("second should be 200, got %d", h.Entries[1].Score)
	}
	if h.Entries[2].Score != 100 {
		t.Errorf("third should be 100, got %d", h.Entries[2].Score)
	}
}

func TestHighScoreMaxEntries(t *testing.T) {
	h := NewHighScoreManager()
	for i := 0; i < 20; i++ {
		h.Add(HighScoreEntry{Name: "X", Score: i * 10})
	}
	// maxHighScoreEntries = 10
	if len(h.Entries) > maxHighScoreEntries {
		t.Errorf("exceeded max entries: %d > %d", len(h.Entries), maxHighScoreEntries)
	}
	// Lowest scores should be dropped
	if h.Entries[len(h.Entries)-1].Score < 100 {
		t.Error("lowest should be high scores only")
	}
}

func TestHighScoreTop(t *testing.T) {
	h := NewHighScoreManager()
	h.Add(HighScoreEntry{Name: "1", Score: 500})
	h.Add(HighScoreEntry{Name: "2", Score: 300})
	h.Add(HighScoreEntry{Name: "3", Score: 100})

	top := h.Top(2)
	if len(top) != 2 {
		t.Fatalf("Top(2) want 2, got %d", len(top))
	}
	if top[0].Score != 500 {
		t.Errorf("top[0] want 500, got %d", top[0].Score)
	}
}

func TestHighScoreRank(t *testing.T) {
	h := NewHighScoreManager()
	h.Add(HighScoreEntry{Name: "1", Score: 500})
	h.Add(HighScoreEntry{Name: "2", Score: 300})

	if r := h.Rank(500); r != 1 {
		t.Errorf("rank of 500 want 1, got %d", r)
	}
	if r := h.Rank(300); r != 2 {
		t.Errorf("rank of 300 want 2, got %d", r)
	}
	if r := h.Rank(999); r != 0 {
		t.Errorf("not-in-list score want rank 0, got %d", r)
	}
}

func TestHighScorePersistence(t *testing.T) {
	h := NewHighScoreManager()
	// Override path for test
	tmp := t.TempDir()
	h.filePath = tmp + "/highscores_test.json"

	h.Add(HighScoreEntry{Name: "Test", Score: 777})
	if err := h.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	h2 := NewHighScoreManager()
	h2.filePath = tmp + "/highscores_test.json"
	if err := h2.Load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(h2.Entries) != 1 {
		t.Fatalf("want 1 loaded entry, got %d", len(h2.Entries))
	}
	if h2.Entries[0].Score != 777 {
		t.Errorf("score mismatch: %d", h2.Entries[0].Score)
	}
}
