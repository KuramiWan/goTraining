package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Timer struct {
	currentTicks int
	targetTicks  int
}

func (t *Timer) UpdateTicks() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

func (t *Timer) RestTicks() {
	t.currentTicks = 0
}

func (t *Timer) IsReadyAttack() bool {
	return t.currentTicks >= t.targetTicks
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks:  int(duration.Milliseconds()) * ebiten.TPS() / 1000,
	}
}
