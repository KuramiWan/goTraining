package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// backgroundDarkColor is the fill color for the background.
var backgroundDarkColor = color.RGBA{8, 4, 24, 255}

// Background renders an animated space background with floating stars.
type Background struct {
	stars []Star
}

// Star is a tiny decorative marker on the background.
type Star struct {
	x, y       float64
	speed      float64
	brightness float32
	sprite     *ebiten.Image
}

// NewBackground creates a new space background with randomly placed stars.
func NewBackground() *Background {
	bg := &Background{
		stars: make([]Star, 0, 100),
	}

	ss := []*ebiten.Image{Star1Sprite, Star2Sprite, Star3Sprite}
	op := &ebiten.DrawImageOptions{}
	for sprite := range ss {
		op.GeoM.Scale()
	}

	for i := 0; i < 100; i++ {
		var sp *ebiten.Image
		if len(ss) > 0 {
			sp = ss[i%len(ss)]
		}
		bg.stars = append(bg.stars, Star{
			x:          float64(rand.Intn(ScreenWidth)),
			y:          float64(rand.Intn(ScreenHeight)),
			speed:      0.15 + float64(i%10)*0.12,
			brightness: 0.25 + float32(i%8)*0.09,
			sprite:     sp,
		})
	}
	return bg
}

// Update scrolls the stars downward.
func (bg *Background) Update() {
	for i := range bg.stars {
		bg.stars[i].y += bg.stars[i].speed
		if bg.stars[i].y > ScreenHeight {
			bg.stars[i].y = 0
			bg.stars[i].x = float64(rand.Intn(ScreenWidth))
		}
	}
}

// Draw fills the screen with a solid dark color then draws the stars on top.
func (bg *Background) Draw(s *ebiten.Image) {
	// Solid fill — simple, reliable, and looks like deep space.
	s.Fill(backgroundDarkColor)

	for _, st := range bg.stars {
		if st.sprite == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.SetA(st.brightness)
		op.GeoM.Translate(st.x, st.y)
		s.DrawImage(st.sprite, op)
	}
}
