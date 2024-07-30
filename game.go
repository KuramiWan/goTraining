package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

//go:embed assets
var data embed.FS

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	playPosition Vector
}

func (g *Game) Update() error {
	speed := 5.0
	var vector Vector
	if vector.X != 0 || vector.Y != 0 {
		factor := speed / math.Sqrt(vector.X*vector.X+vector.Y*vector.Y)
		vector.X *= factor
		vector.Y *= factor
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.playPosition.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.playPosition.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playPosition.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playPosition.X += speed
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	options := &ebiten.DrawImageOptions{}
	//width, height := PlaySprite.Bounds().Dx(), PlaySprite.Bounds().Dy()
	//halfW, halfH := float64(width/2), float64(height/2)
	//options.GeoM.Translate(-halfW, -halfH)
	//options.GeoM.Rotate(45.0 * math.Pi / 180)
	//m := colorm.ColorM{}
	//m.Translate(1, 1, 1, 1)
	options.GeoM.Translate(g.playPosition.X, g.playPosition.Y)
	//options.GeoM.Scale(1, -1)
	screen.DrawImage(PlaySprite, options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{
		playPosition: Vector{X: 100, Y: 100},
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
