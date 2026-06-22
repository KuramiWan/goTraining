package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// backgroundDarkColor is the fill color for the background.
var backgroundDarkColor = color.RGBA{8, 4, 24, 255}

// Background represents the game's background.
type Background struct{}

// NewBackground creates a new background.
func NewBackground() *Background {
	return &Background{}
}

// Update is a no-op for the solid-color background.
func (bg *Background) Update() {}

// Draw fills the screen with a solid dark color.
func (bg *Background) Draw(s *ebiten.Image) {
	s.Fill(backgroundDarkColor)
}
