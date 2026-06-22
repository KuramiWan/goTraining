package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// HUD renders sprite-based UI elements.
type HUD struct {
	numeralDigits []*ebiten.Image
	lifeSprite    *ebiten.Image
}

// NewHUD creates a new HUD instance.
func NewHUD() *HUD {
	return &HUD{
		numeralDigits: NumeralSprites,
		lifeSprite:    LifeBlueSprite,
	}
}

// DrawScore renders the score using numeral sprites.
func (h *HUD) DrawScore(s *ebiten.Image, score int, x, y float64) {
	digits := fmt.Sprintf("%d", score)
	digitWidth := 40.0
	for i, ch := range digits {
		idx := int(ch - '0')
		if idx < 0 || idx >= len(h.numeralDigits) || h.numeralDigits[idx] == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x+float64(i)*digitWidth, y)
		s.DrawImage(h.numeralDigits[idx], op)
	}
}

// DrawLives renders life icons.
func (h *HUD) DrawLives(s *ebiten.Image, lives int, x, y float64) {
	if h.lifeSprite == nil {
		return
	}
	for i := 0; i < lives; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x+float64(i)*45, y)
		s.DrawImage(h.lifeSprite, op)
	}
}

// DrawBossHealth renders a boss health bar.
func (h *HUD) DrawBossHealth(s *ebiten.Image, health, maxHealth int) {
	if maxHealth <= 0 {
		return
	}
	barW := 600.0
	barH := 20.0
	x := float64(ScreenWidth)/2 - barW/2
	y := 30.0

	// Background
	bg := ebiten.NewImage(int(barW), int(barH))
	bg.Fill(color.RGBA{60, 60, 60, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	s.DrawImage(bg, op)

	// Health fill
	pct := float64(health) / float64(maxHealth)
	if pct < 0 {
		pct = 0
	}
	fillW := barW * pct
	if fillW > 0 {
		fill := ebiten.NewImage(int(fillW), int(barH))
		if pct > 0.5 {
			fill.Fill(color.RGBA{0, 200, 0, 255})
		} else if pct > 0.25 {
			fill.Fill(color.RGBA{200, 200, 0, 255})
		} else {
			fill.Fill(color.RGBA{200, 0, 0, 255})
		}
		s.DrawImage(fill, op)
	}
}

// DrawWeaponIndicator shows the current weapon type.
func (h *HUD) DrawWeaponIndicator(s *ebiten.Image, weaponIdx int) {
	x := float64(ScreenWidth - 80)
	y := float64(ScreenHeight - 80)
	text.Draw(s, "武器", UIFont, int(x-20), int(y), color.White)

	laserSprites := [][]*ebiten.Image{LaserBlueSprites, LaserGreenSprites, LaserRedSprites}
	if weaponIdx >= 0 && weaponIdx < len(laserSprites) && len(laserSprites[weaponIdx]) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y+30)
		s.DrawImage(laserSprites[weaponIdx][0], op)
	}
}

// DrawAchievementPopup shows a brief achievement notification.
func (h *HUD) DrawAchievementPopup(s *ebiten.Image, ach *Achievement) {
	if ach == nil {
		return
	}
	x := float64(ScreenWidth)/2 - 150
	y := 120.0

	bg := ebiten.NewImage(300, 50)
	bg.Fill(color.RGBA{0, 0, 0, 200})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	s.DrawImage(bg, op)

	text.Draw(s, "成就解锁: "+ach.Name, UIFont, int(x+10), int(y+32), color.RGBA{255, 215, 0, 255})
}

// DrawDebugOverlay shows debug information.
func (h *HUD) DrawDebugOverlay(s *ebiten.Image, g *Game) {
	y := ScreenHeight - 120
	lineH := 20
	lines := []string{
		fmt.Sprintf("FPS: %.1f  TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()),
		fmt.Sprintf("实体: M=%d B=%d E=%d",
			len(g.meteors.value), len(g.player.bullets),
			func() int {
				if g.enemies != nil {
					return len(g.enemies.value)
				}
				return 0
			}()),
		fmt.Sprintf("命中率: %.1f%%  穿透: %d", g.Stats.Accuracy()*100, g.Stats.PiercingHits),
		fmt.Sprintf("GM:%v SM:%v 分数:%d 生命:%d", g.godMode, g.slowMotion, g.Score, g.Lives),
	}
	for i, line := range lines {
		text.Draw(s, line, UIFont, 10, y+i*lineH, color.RGBA{0, 255, 0, 200})
	}
}

// DrawHitboxes visualizes collision rectangles for debugging.
func (h *HUD) DrawHitboxes(s *ebiten.Image, g *Game) {
	drawRect := func(r *Rect, c color.Color) {
		// Top
		top := ebiten.NewImage(int(r.Width), 1)
		top.Fill(c)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(r.X, r.Y)
		s.DrawImage(top, op)
		// Bottom
		bot := ebiten.NewImage(int(r.Width), 1)
		bot.Fill(c)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(r.X, r.MaxY())
		s.DrawImage(bot, op2)
		// Left
		left := ebiten.NewImage(1, int(r.Height))
		left.Fill(c)
		op3 := &ebiten.DrawImageOptions{}
		op3.GeoM.Translate(r.X, r.Y)
		s.DrawImage(left, op3)
		// Right
		right := ebiten.NewImage(1, int(r.Height))
		right.Fill(c)
		op4 := &ebiten.DrawImageOptions{}
		op4.GeoM.Translate(r.MaxX(), r.Y)
		s.DrawImage(right, op4)
	}

	// Player hitbox (green)
	pRect := g.player.Collider()
	drawRect(pRect, color.RGBA{0, 255, 0, 200})

	// Meteor hitboxes (red)
	for _, m := range g.meteors.value {
		drawRect(m.Collider(), color.RGBA{255, 0, 0, 200})
	}
}
