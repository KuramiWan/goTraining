package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var data embed.FS

type Game struct {
	player  *Player
	meteors *Meteors
}

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func (g *Game) Update() error {
	g.player.Update()
	//g.bulletColdTimer.Update()

	g.meteors.Update()
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	//options := &ebiten.DrawImageOptions{}
	//width, height := PlaySprite.Bounds().Dx(), PlaySprite.Bounds().Dy()
	//halfW, halfH := float64(width/2), float64(height/2)
	//options.GeoM.Translate(-halfW, -halfH)
	//options.GeoM.Rotate(45.0 * math.Pi / 180)
	//m := colorm.ColorM{}
	//m.Translate(1, 1, 1, 1)
	//options.GeoM.Translate(g.playPosition.X, g.playPosition.Y)
	//options.GeoM.Scale(1, -1)
	//s.DrawImage(PlaySprite, options)
	g.player.Draw(s)
	g.meteors.Draw(s)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	g := &Game{
		player:  newPlayer(),
		meteors: newMeteors(),
	}
	//open, err := assets.Open("Sprite/playerShip1_blue.png")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(open)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
