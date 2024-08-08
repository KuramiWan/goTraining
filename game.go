package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

//go:embed assets
var data embed.FS

type Game struct {
	player  *Player
	meteors *Meteors
	score   int
}

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func (g *Game) Update() error {
	g.player.Update()
	g.meteors.Update()
	g.BulletCollisions()
	g.MeteorCollisions()
	return nil
}

func (g *Game) MeteorCollisions() {
	p := g.player.Collider()
	for i, meteor := range g.meteors.value {
		m := meteor.Collider()
		if m.Intersect(p) {
			g.meteors.value = append(g.meteors.value[:i], g.meteors.value[i+1:]...)
			g.reset()
		}
	}
}

func (g *Game) BulletCollisions() {
	for i, meteor := range g.meteors.value {
		for j, bullet := range g.player.bullets {
			m := meteor.Collider()
			b := bullet.Collider()
			if m.Intersect(b) {
				g.meteors.value = append(g.meteors.value[:i], g.meteors.value[i+1:]...)
				g.player.bullets = append(g.player.bullets[:j], g.player.bullets[j+1:]...)
				g.score++
			}
		}
	}
}

func (g *Game) reset() {
	g.player = newPlayer()
	g.meteors = newMeteors()
	g.player.bullets = newBullets()
	g.score = 0
}

func (g *Game) Draw(s *ebiten.Image) {
	g.player.Draw(s)
	g.meteors.Draw(s)
	text.Draw(s, fmt.Sprintf("%06d", g.score), ScoreFont, ScreenWidth/2-100, 50, color.White)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	g := &Game{
		player:  newPlayer(),
		meteors: newMeteors(),
		score:   0,
	}
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
