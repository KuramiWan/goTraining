package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

//go:embed assets
var data embed.FS

type Game struct {
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
}

type Timer struct {
	currentTicks int
	targetTicks  int
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks:  int(duration.Milliseconds()) * ebiten.TPS() / 1000,
	}
}

func (g *Game) Update() error {
	g.player.Update()
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReadyAttack() {
		g.meteorSpawnTimer.RestTicks()
		meteor := newMeteor()
		g.meteors = append(g.meteors, meteor)
	}
	for _, meteor := range g.meteors {
		meteor.Update()
	}
	return nil
}

func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

func (t *Timer) RestTicks() {
	t.currentTicks = 0
}

func (t *Timer) IsReadyAttack() bool {
	return t.currentTicks >= t.targetTicks
}

func (g *Game) Draw(screen *ebiten.Image) {
	//options := &ebiten.DrawImageOptions{}
	//width, height := PlaySprite.Bounds().Dx(), PlaySprite.Bounds().Dy()
	//halfW, halfH := float64(width/2), float64(height/2)
	//options.GeoM.Translate(-halfW, -halfH)
	//options.GeoM.Rotate(45.0 * math.Pi / 180)
	//m := colorm.ColorM{}
	//m.Translate(1, 1, 1, 1)
	//options.GeoM.Translate(g.playPosition.X, g.playPosition.Y)
	//options.GeoM.Scale(1, -1)
	//screen.DrawImage(PlaySprite, options)
	g.player.Draw(screen)
	for _, meteor := range g.meteors {
		meteor.Draw(screen)
	}
}

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	var meteors = make([]*Meteor, 0)
	g := &Game{
		meteorSpawnTimer: NewTimer(5 * time.Second),
		player:           newPlayer(),
		meteors:          meteors,
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
