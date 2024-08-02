package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"io/fs"
	"path"
)

var PlaySprite = mustLoadImage("Sprite/playerShip1_blue.png")

func mustLoadImage(name string) *ebiten.Image {
	p, err := data.Open(path.Join("assets", name))
	if err != nil {
		panic(err)
	}
	defer func(p fs.File) {
		err := p.Close()
		if err != nil {
			panic(err)
		}
	}(p)
	img, _, err := image.Decode(p)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
