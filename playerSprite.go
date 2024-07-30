package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"io/fs"
	"path"
)

var PlaySprite = mustLoadImage("Sprite/playerShip1_blue.png")

func mustLoadImage(name string) *ebiten.Image {
	p, err := data.Open(path.Join("assets", name))
	defer p.Close()
	img, _, err := image.Decode(p)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

// 获取dir路径
func getDir(dirName string, fs embed.FS) []fs.DirEntry {
	dir, err := fs.ReadDir(dirName)
	if err != nil {
		panic(err)
		return nil
	}
	return dir
}
