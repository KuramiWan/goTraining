package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
)

var PlaySprite = mustLoadImage("Sprite/playerShip1_blue.png")

func mustLoadImage(name string) *ebiten.Image {
	p, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer p.Close()
	img, _, err := image.Decode(p)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
