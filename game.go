package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var data embed.FS

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	//options.GeoM.Translate(150, 200)
	//options.GeoM.Rotate(25.0 * math.Pi / 180)
	//options.GeoM.Scale(1, -1)
	screen.DrawImage(PlaySprite, options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	//open, err := assets.Open("Sprite/playerShip1_blue.png")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(open)
	g := &Game{}
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
