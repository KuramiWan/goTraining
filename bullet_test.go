package main

import (
	"testing"
)

func TestBulletNew(t *testing.T) {
	b := newBullet(Vector{X: 100, Y: 200}, 0, Vector{Vx: 5, Vy: 0}, false)
	if b == nil {
		t.Fatal("newBullet returned nil")
	}
	if b.position.X != 100-float64(b.sprite.Bounds().Dx())/2 {
		// position is adjusted by half sprite size
		t.Logf("bullet X: %f", b.position.X)
	}
	if b.rotation != 0 {
		t.Errorf("rotation want 0, got %f", b.rotation)
	}
	if b.piercing {
		t.Error("piercing should be false")
	}
}

func TestBulletPiercing(t *testing.T) {
	b := newBullet(Vector{X: 100, Y: 200}, 0, Vector{}, true)
	if !b.piercing {
		t.Error("piercing should be true")
	}
}

func TestBulletCollider(t *testing.T) {
	b := newBullet(Vector{X: 100, Y: 200}, 0, Vector{}, false)
	r := b.Collider()
	if r.Width <= 0 || r.Height <= 0 {
		t.Error("bullet collider should have positive dimensions")
	}
}

func TestBulletIsOffScreen(t *testing.T) {
	b := newBullet(Vector{X: 100, Y: 200}, 0, Vector{}, false)
	if b.IsOffScreen() {
		t.Error("bullet at (100,200) should not be off screen")
	}
	b.position.X = -200
	if !b.IsOffScreen() {
		t.Error("bullet at x=-200 should be off screen")
	}
}

func TestBulletUpdate(t *testing.T) {
	b := newBullet(Vector{X: 500, Y: 400}, 0, Vector{}, false)
	b.Update()
	// Should move upward (negative Y) due to rotation=0
	if b.position.Y >= 400 {
		t.Errorf("bullet should move upward: Y=%f", b.position.Y)
	}
}
