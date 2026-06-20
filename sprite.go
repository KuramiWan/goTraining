package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	_ "image/png"
	"io/fs"
	"os"
	"path"
	"strings"
)

var PlaySprite = mustLoadImage("Sprite/playerShip1_blue.png")

var LaserSprite = mustLoadImage("PNG/Lasers/laserBlue01.png")

var PowerUpSilverSprite = mustLoadImage("PNG/Power-ups/bold_silver.png")
var PowerUpBronzeSprite = mustLoadImage("PNG/Power-ups/bolt_bronze.png")
var PowerUpGoldSprite = mustLoadImage("PNG/Power-ups/bolt_gold.png")

var ScoreFont = mustLoadFont("assets/font/Kenney Blocks.ttf")

var UIFont = mustLoadSystemFont(`C:\Windows\Fonts\simhei.ttf`)

var MeteorBigSprites = mustLoadFilteredImages("assets/PNG/Meteors", "big")

var MeteorMedSprites = mustLoadFilteredImages("assets/PNG/Meteors", "med")

var MeteorSmallSprites = mustLoadFilteredImages("assets/PNG/Meteors", "small")

var MeteorTinySprites = mustLoadFilteredImages("assets/PNG/Meteors", "tiny")

func mustLoadFilteredImages(dir, pattern string) []*ebiten.Image {
	entries, err := data.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var images []*ebiten.Image
	for _, entry := range entries {
		if !strings.Contains(entry.Name(), pattern) {
			continue
		}
		p, err := data.Open(path.Join(dir, entry.Name()))
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(p)
		p.Close()
		if err != nil {
			panic(err)
		}
		images = append(images, ebiten.NewImageFromImage(img))
	}
	return images
}

func mustLoadImage(n string) *ebiten.Image {
	p, err := data.Open(path.Join("assets", n))
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

func mustLoadImages(n string) []*ebiten.Image {
	dir := getDir(n)
	var images []*ebiten.Image
	for _, dirEntry := range dir {
		func() {
			p, err := data.Open(path.Join(n, dirEntry.Name()))
			if err != nil {
				panic(err)
			}
			defer p.Close()
			img, _, err := image.Decode(p)
			if err != nil {
				panic(err)
			}
			images = append(images, ebiten.NewImageFromImage(img))
		}()
	}
	return images
}

func getDir(d string) []fs.DirEntry {
	dir, err := data.ReadDir(d)
	if err != nil {
		panic(err)
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

func mustLoadSystemFont(path string) font.Face {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(file)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	return face
}
