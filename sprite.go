package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	_ "image/png"
	"io/fs"
	"path"
)

var PlaySprite = mustLoadImage("Sprite/playerShip1_blue.png")

var MeteorSprites = mustLoadImages("assets/PNG/Meteors")

var LaserSprite = mustLoadImage("PNG/Lasers/laserBlue01.png")

var ScoreFont = mustLoadFont("assets/font/Kenney Blocks.ttf")

func mustLoadImage(n string) *ebiten.Image {
	p, err := data.Open(path.Join("assets", n))
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

func mustLoadImages(n string) []*ebiten.Image {
	dir := getDir(n)
	var images []*ebiten.Image
	for _, dirEntry := range dir {
		p, err := data.Open(path.Join(n, dirEntry.Name()))
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
		images = append(images, ebiten.NewImageFromImage(img))
	}
	return images
}

// 获取dir路径
func getDir(d string) []fs.DirEntry {
	dir, err := data.ReadDir(d)
	if err != nil {
		panic(err)
		return nil
	}
	return dir
}

func mustLoadFont(name string) font.Face {
	file, err := data.ReadFile(name)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(file)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}
	return face
}
