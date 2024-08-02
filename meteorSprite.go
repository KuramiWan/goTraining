package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"io/fs"
	"path"
)

var MeteorSprites = mustLoadImages("assets/PNG/Meteors")

func mustLoadImages(name string) []*ebiten.Image {

	dir := getDir(name)
	var images []*ebiten.Image
	for _, dirEntry := range dir {
		p, err := data.Open(path.Join(name, dirEntry.Name()))
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
func getDir(dirName string) []fs.DirEntry {
	dir, err := data.ReadDir(dirName)
	if err != nil {
		panic(err)
		return nil
	}
	return dir
}
