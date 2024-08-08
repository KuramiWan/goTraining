package main

import "github.com/hajimehoshi/ebiten/v2"

type Rect struct {
	X, Y, Width, Height float64
}

func (r Rect) MaxX() float64 {
	return r.X + r.Width
}
func (r Rect) MaxY() float64 {
	return r.Y + r.Height
}
func (r Rect) Intersect(Other *Rect) bool {
	return r.X <= Other.MaxX() && Other.X <= r.MaxX() && r.Y <= Other.MaxY() && Other.Y <= r.MaxY()
}

func newRect(pos Vector, sprite *ebiten.Image) *Rect {
	return &Rect{
		X:      pos.X,
		Y:      pos.Y,
		Width:  float64(sprite.Bounds().Dx()),
		Height: float64(sprite.Bounds().Dy()),
	}
}
