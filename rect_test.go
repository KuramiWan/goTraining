package main

import (
	"testing"
)

func TestRectMaxXY(t *testing.T) {
	r := &Rect{X: 10, Y: 20, Width: 30, Height: 40}
	if r.MaxX() != 40 {
		t.Errorf("MaxX want 40, got %f", r.MaxX())
	}
	if r.MaxY() != 60 {
		t.Errorf("MaxY want 60, got %f", r.MaxY())
	}
}

func TestRectIntersectOverlap(t *testing.T) {
	a := &Rect{X: 10, Y: 10, Width: 20, Height: 20}
	b := &Rect{X: 20, Y: 20, Width: 20, Height: 20}
	if !a.Intersect(b) {
		t.Error("overlapping rects should intersect")
	}
	if !b.Intersect(a) {
		t.Error("intersection should be symmetric")
	}
}

func TestRectIntersectSeparated(t *testing.T) {
	a := &Rect{X: 10, Y: 10, Width: 20, Height: 20}
	b := &Rect{X: 100, Y: 100, Width: 20, Height: 20}
	if a.Intersect(b) {
		t.Error("separated rects should not intersect")
	}
}

func TestRectIntersectEdgeTouch(t *testing.T) {
	a := &Rect{X: 0, Y: 0, Width: 10, Height: 10}
	b := &Rect{X: 10, Y: 0, Width: 10, Height: 10}
	if !a.Intersect(b) {
		t.Error("edge-touching rects should intersect (MaxX=10 meets X=10)")
	}
}

func TestRectIntersectContained(t *testing.T) {
	a := &Rect{X: 0, Y: 0, Width: 100, Height: 100}
	b := &Rect{X: 25, Y: 25, Width: 50, Height: 50}
	if !a.Intersect(b) {
		t.Error("smaller inside larger should intersect")
	}
}

func TestNewRect(t *testing.T) {
	// Use a sprite that exists
	pos := Vector{X: 100, Y: 200}
	r := newRect(pos, LaserSprite)
	if r.X != 100 || r.Y != 200 {
		t.Errorf("position mismatch: got (%f,%f) want (100,200)", r.X, r.Y)
	}
	if r.Width <= 0 || r.Height <= 0 {
		t.Errorf("dimensions should be positive: %fx%f", r.Width, r.Height)
	}
}
