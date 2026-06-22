package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"time"
)

//go:embed assets
var data embed.FS

type GameState int

const (
	StateStart GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

const (
	ScreenWidth       = 1600
	ScreenHeight      = 900
	InitialLives      = 3
	InvulnDuration    = 2 * time.Second
	BaseSpawnInterval = 5 * time.Second
	MinSpawnInterval  = 500 * time.Millisecond
)

type Game struct {
	player         *Player
	meteors        *Meteors
	Score          int
	Lives          int
	State          GameState
	DifficultyLevel int
	invulnTimer    *Timer
	firstSpawnDone bool
	flashCounter   int
}

func (g *Game) Update() error {
	switch g.State {
	case StateStart:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startGame()
		}
		return nil
	case StatePlaying:
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.State = StatePaused
			return nil
		}
		g.updatePlaying()
		return nil
	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.State = StatePlaying
		}
		return nil
	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset()
		}
		return nil
	}
	return nil
}

func (g *Game) startGame() {
	g.State = StatePlaying
	g.Score = 0
	g.Lives = InitialLives
	g.DifficultyLevel = 0
	g.firstSpawnDone = false
	g.invulnTimer = nil
	g.flashCounter = 0
	g.player = newPlayer()
	g.meteors = newMeteors()
}

func (g *Game) updatePlaying() {
	g.player.Update()
	g.meteors.Update(!g.firstSpawnDone)
	if !g.firstSpawnDone {
		g.firstSpawnDone = true
	}
	g.bulletCollisions()
	g.meteorCollisions()
	g.updateDifficulty()
	g.cleanupOffScreen()
	g.updateInvulnerability()
	g.flashCounter++
}

func (g *Game) updateInvulnerability() {
	if g.invulnTimer != nil {
		g.invulnTimer.UpdateTicks()
	}
}

func (g *Game) updateDifficulty() {
	newLevel := g.Score / 10
	if newLevel > g.DifficultyLevel {
		g.DifficultyLevel = newLevel
		interval := BaseSpawnInterval - time.Duration(g.DifficultyLevel)*500*time.Millisecond
		if interval < MinSpawnInterval {
			interval = MinSpawnInterval
		}
		g.meteors.meteorSpawnTimer.SetDuration(interval)
	}
}

func (g *Game) meteorCollisions() {
	if g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack() {
		return
	}
	pRect := g.player.Collider()
	hitIndex := -1
	for i, meteor := range g.meteors.value {
		if meteor.Collider().Intersect(pRect) {
			hitIndex = i
			break
		}
	}
	if hitIndex >= 0 {
		g.meteors.value = append(g.meteors.value[:hitIndex], g.meteors.value[hitIndex+1:]...)
		g.Lives--
		if g.Lives <= 0 {
			g.State = StateGameOver
		} else {
			g.invulnTimer = NewTimer(InvulnDuration)
		}
	}
}

func (g *Game) bulletCollisions() {
	if len(g.meteors.value) == 0 || len(g.player.bullets) == 0 {
		return
	}
	meteorHit := make([]bool, len(g.meteors.value))
	bulletHit := make([]bool, len(g.player.bullets))

	for i, meteor := range g.meteors.value {
		mRect := meteor.Collider()
		for j, bullet := range g.player.bullets {
			if !meteorHit[i] && !bulletHit[j] && mRect.Intersect(bullet.Collider()) {
				meteorHit[i] = true
				bulletHit[j] = true
				g.Score++
			}
		}
	}

	survivingMeteors := make([]*Meteor, 0, len(g.meteors.value))
	for i, m := range g.meteors.value {
		if !meteorHit[i] {
			survivingMeteors = append(survivingMeteors, m)
		}
	}
	g.meteors.value = survivingMeteors

	survivingBullets := make([]*Bullet, 0, len(g.player.bullets))
	for j, b := range g.player.bullets {
		if !bulletHit[j] {
			survivingBullets = append(survivingBullets, b)
		}
	}
	g.player.bullets = survivingBullets
}

func (g *Game) cleanupOffScreen() {
	pb := g.player.bullets[:0]
	for _, b := range g.player.bullets {
		if !b.IsOffScreen() {
			pb = append(pb, b)
		}
	}
	g.player.bullets = pb

	ms := g.meteors.value[:0]
	for _, m := range g.meteors.value {
		if !m.IsOffScreen() {
			ms = append(ms, m)
		}
	}
	g.meteors.value = ms
}

func (g *Game) reset() {
	g.player = newPlayer()
	g.meteors = newMeteors()
	g.Score = 0
	g.Lives = InitialLives
	g.DifficultyLevel = 0
	g.firstSpawnDone = false
	g.invulnTimer = nil
	g.flashCounter = 0
	g.State = StatePlaying
}

func (g *Game) Draw(s *ebiten.Image) {
	switch g.State {
	case StateStart:
		g.drawStartScreen(s)
	case StateGameOver:
		g.drawPlayField(s)
		g.drawGameOverOverlay(s)
	case StatePlaying, StatePaused:
		g.drawPlayField(s)
		if g.State == StatePaused {
			g.drawPauseOverlay(s)
		}
	}
}

func (g *Game) drawStartScreen(s *ebiten.Image) {
	title := "SPACE SHOOTER"
	text.Draw(s, title, ScoreFont, ScreenWidth/2-250, ScreenHeight/2-120, color.White)
	prompt := "Press ENTER or SPACE to Start"
	text.Draw(s, prompt, ScoreFont, ScreenWidth/2-350, ScreenHeight/2-20, color.White)
	controls := "WASD: Move | Mouse L/R: Rotate | SPACE: Shoot | P: Pause"
	text.Draw(s, controls, ScoreFont, ScreenWidth/2-520, ScreenHeight/2+80, color.RGBA{180, 180, 180, 255})
}

func (g *Game) drawPlayField(s *ebiten.Image) {
	isInvuln := g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack()
	drawPlayer := !isInvuln || (g.flashCounter/6)%2 == 0

	if drawPlayer {
		g.player.Draw(s)
	}
	g.meteors.Draw(s)

	text.Draw(s, fmt.Sprintf("%06d", g.Score), ScoreFont, ScreenWidth/2-100, 50, color.White)
	livesText := fmt.Sprintf("Lives: %d", g.Lives)
	text.Draw(s, livesText, ScoreFont, 50, 50, color.White)
	levelText := fmt.Sprintf("Level: %d", g.DifficultyLevel)
	text.Draw(s, levelText, ScoreFont, ScreenWidth-300, 50, color.White)
}

func (g *Game) drawGameOverOverlay(s *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	s.DrawImage(overlay, nil)

	gameOver := "GAME OVER"
	text.Draw(s, gameOver, ScoreFont, ScreenWidth/2-180, ScreenHeight/2-120, color.RGBA{255, 50, 50, 255})
	scoreText := fmt.Sprintf("Score: %d  Level: %d", g.Score, g.DifficultyLevel)
	text.Draw(s, scoreText, ScoreFont, ScreenWidth/2-250, ScreenHeight/2-40, color.White)
	prompt := "Press ENTER or SPACE to Restart"
	text.Draw(s, prompt, ScoreFont, ScreenWidth/2-350, ScreenHeight/2+40, color.White)
}

func (g *Game) drawPauseOverlay(s *ebiten.Image) {
	pauseText := "PAUSED"
	text.Draw(s, pauseText, ScoreFont, ScreenWidth/2-120, ScreenHeight/2-40, color.White)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Space Shooter")
	g := &Game{
		player:  newPlayer(),
		meteors: newMeteors(),
		State:   StateStart,
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
