package main

import "github.com/hajimehoshi/ebiten/v2"

// EffectType categorizes animation effects.
type EffectType int

const (
	EffExplosion EffectType = iota
	EffShield
	EffStar
	EffSpeed
)

// AnimationEffect plays a sequence of frames at a position.
type AnimationEffect struct {
	effectType EffectType
	position   Vector
	frames     []*ebiten.Image
	frameIdx   int
	frameTimer *Timer
	finished   bool
	rotation   float64
	scale      float64
}

// EffectManager manages all active animation effects.
type EffectManager struct {
	effects []*AnimationEffect
}

// NewEffectManager creates a new effect manager.
func NewEffectManager() *EffectManager {
	return &EffectManager{
		effects: make([]*AnimationEffect, 0),
	}
}

// AddExplosion adds an explosion effect at a position.
func (em *EffectManager) AddExplosion(x, y, scale float64) {
	// Stub: will be implemented in Phase 3
}

// AddEffect adds a generic effect.
func (em *EffectManager) AddEffect(effectType EffectType, x, y, scale float64) {
	// Stub: will be implemented in Phase 3
}

// Update advances all active effects.
func (em *EffectManager) Update() {
	surviving := em.effects[:0]
	for _, e := range em.effects {
		e.frameTimer.UpdateTicks()
		if e.frameTimer.IsReadyAttack() {
			e.frameTimer.RestTicks()
			e.frameIdx++
			if e.frameIdx >= len(e.frames) {
				e.finished = true
			}
		}
		if !e.finished {
			surviving = append(surviving, e)
		}
	}
	em.effects = surviving
}

// Draw renders all active effects.
func (em *EffectManager) Draw(s *ebiten.Image) {
	for _, e := range em.effects {
		if e.frameIdx < len(e.frames) {
			op := &ebiten.DrawImageOptions{}
			img := e.frames[e.frameIdx]
			bounds := img.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2
			op.GeoM.Translate(-halfW, -halfH)
			if e.rotation != 0 {
				op.GeoM.Rotate(e.rotation)
			}
			if e.scale != 1.0 {
				op.GeoM.Scale(e.scale, e.scale)
			}
			op.GeoM.Translate(halfW, halfH)
			op.GeoM.Translate(e.position.X, e.position.Y)
			s.DrawImage(img, op)
		}
	}
}

// Count returns the number of active effects.
func (em *EffectManager) Count() int {
	return len(em.effects)
}
