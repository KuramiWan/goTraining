package main

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	d := 1 * time.Second
	timer := NewTimer(d)
	if timer == nil {
		t.Fatal("NewTimer returned nil")
	}
	if timer.currentTicks != 0 {
		t.Errorf("currentTicks want 0, got %d", timer.currentTicks)
	}
	if timer.targetTicks <= 0 {
		t.Errorf("targetTicks should be positive, got %d", timer.targetTicks)
	}
}

func TestTimerUpdateTicks(t *testing.T) {
	timer := &Timer{currentTicks: 0, targetTicks: 10}
	timer.UpdateTicks()
	if timer.currentTicks != 1 {
		t.Errorf("after update, want 1 got %d", timer.currentTicks)
	}
}

func TestTimerUpdateTicksCap(t *testing.T) {
	timer := &Timer{currentTicks: 10, targetTicks: 10}
	timer.UpdateTicks()
	if timer.currentTicks != 10 {
		t.Errorf("should not exceed target, got %d", timer.currentTicks)
	}
}

func TestTimerRestTicks(t *testing.T) {
	timer := &Timer{currentTicks: 8, targetTicks: 10}
	timer.RestTicks()
	if timer.currentTicks != 0 {
		t.Errorf("after reset, want 0 got %d", timer.currentTicks)
	}
}

func TestTimerIsReadyAttack(t *testing.T) {
	timer := &Timer{currentTicks: 5, targetTicks: 10}
	if timer.IsReadyAttack() {
		t.Error("should not be ready at 5/10")
	}
	timer.currentTicks = 10
	if !timer.IsReadyAttack() {
		t.Error("should be ready at 10/10")
	}
	timer.currentTicks = 15
	if !timer.IsReadyAttack() {
		t.Error("should be ready when over target")
	}
}

func TestTimerSetDuration(t *testing.T) {
	timer := &Timer{currentTicks: 5, targetTicks: 10}
	timer.SetDuration(2 * time.Second)
	if timer.targetTicks <= 0 {
		t.Errorf("targetTicks should be positive after SetDuration")
	}
	// Current ticks should not be reset by SetDuration alone
	if timer.currentTicks != 5 {
		t.Errorf("currentTicks should be preserved, got %d", timer.currentTicks)
	}
}

func TestTimerProgress(t *testing.T) {
	timer := &Timer{currentTicks: 5, targetTicks: 10}
	if p := timer.Progress(); !floatApprox(p, 0.5, 0.001) {
		t.Errorf("progress want 0.5, got %f", p)
	}

	timer.currentTicks = 10
	if p := timer.Progress(); !floatApprox(p, 1.0, 0.001) {
		t.Errorf("progress want 1.0, got %f", p)
	}

	empty := &Timer{currentTicks: 0, targetTicks: 0}
	if p := empty.Progress(); !floatApprox(p, 1.0, 0.001) {
		t.Errorf("progress for 0-target timer want 1.0, got %f", p)
	}
}
