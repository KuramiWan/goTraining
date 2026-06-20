package main

import (
	"embed"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	ScreenWidth        = 1600
	ScreenHeight       = 900
	InitialLives       = 3
	InvulnDuration     = 2 * time.Second
	BaseSpawnInterval  = 5 * time.Second
	MinSpawnInterval   = 500 * time.Millisecond
	BaseFireRate       = 1 * time.Second
	MinFireRate        = 200 * time.Millisecond
	PowerUpSpawnPeriod = 60 * time.Second
)

type Game struct {
	player              *Player
	meteors             *Meteors
	powerUps            []*PowerUp
	Score               int
	Lives               int
	State               GameState
	DifficultyLevel     int
	invulnTimer         *Timer
	firstSpawnDone      bool
	flashCounter        int
	confirmQuit         bool
	elapsedTime         time.Duration
	powerUpSpawnTimer   *Timer
	lastBronzeMilestone int
}

func (g *Game) Update() error {
	switch g.State {
	case StateStart:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startGame()
		}
		return nil
	case StatePlaying:
		if g.confirmQuit {
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				g.State = StateStart
				g.confirmQuit = false
				return nil
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				g.confirmQuit = false
				return nil
			}
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.confirmQuit = true
			return nil
		}
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
	g.confirmQuit = false
	g.elapsedTime = 0
	g.powerUps = nil
	g.powerUpSpawnTimer = NewTimer(PowerUpSpawnPeriod)
	g.lastBronzeMilestone = 10
	g.player = newPlayer()
	g.meteors = newMeteors()
}

func (g *Game) updatePlaying() {
	g.player.Update()
	g.meteors.Update(!g.firstSpawnDone)
	if !g.firstSpawnDone {
		g.firstSpawnDone = true
	}
	g.updatePowerUps()
	g.bulletCollisions()
	g.meteorCollisions()
	g.powerUpCollisions()
	g.bronzeMilestoneSpawn()
	g.updateDifficulty()
	g.cleanupOffScreen()
	g.updateInvulnerability()
	g.flashCounter++
	g.elapsedTime += time.Second / time.Duration(ebiten.TPS())
}

func (g *Game) updatePowerUps() {
	g.powerUpSpawnTimer.UpdateTicks()
	if g.powerUpSpawnTimer.IsReadyAttack() {
		g.powerUpSpawnTimer.RestTicks()
		g.spawnRandomPowerUp()
	}
	for _, pu := range g.powerUps {
		pu.Update()
	}
}

func (g *Game) spawnRandomPowerUp() {
	r := rand.Intn(10) // 0-9
	var t PowerUpType
	switch {
	case r < 5:
		t = PowerUpSilver // 50%
	case r < 9:
		t = PowerUpBronze // 40%
	default:
		t = PowerUpGold // 10%
	}
	g.powerUps = append(g.powerUps, newPowerUp(t))
}

func (g *Game) bronzeMilestoneSpawn() {
	if g.Score >= g.lastBronzeMilestone {
		g.powerUps = append(g.powerUps, newPowerUp(PowerUpBronze))
		g.lastBronzeMilestone *= 10
	}
}

func (g *Game) powerUpCollisions() {
	pRect := g.player.Collider()
	surviving := g.powerUps[:0]
	for _, pu := range g.powerUps {
		if pu.Collider().Intersect(pRect) {
			g.player.ApplyPowerUp(pu.pType)
		} else {
			surviving = append(surviving, pu)
		}
	}
	g.powerUps = surviving
}

func (g *Game) updateInvulnerability() {
	if g.invulnTimer != nil {
		g.invulnTimer.UpdateTicks()
	}
}

func (g *Game) updateDifficulty() {
	scoreLevel := g.Score / 10
	timeLevel := int(g.elapsedTime.Seconds()) / 30
	newLevel := scoreLevel + timeLevel
	if newLevel > g.DifficultyLevel {
		g.DifficultyLevel = newLevel
		interval := BaseSpawnInterval - time.Duration(g.DifficultyLevel)*500*time.Millisecond
		if interval < MinSpawnInterval {
			interval = MinSpawnInterval
		}
		g.meteors.meteorSpawnTimer.SetDuration(interval)

		fireRate := BaseFireRate - time.Duration(g.DifficultyLevel)*100*time.Millisecond
		if fireRate < MinFireRate {
			fireRate = MinFireRate
		}
		g.player.SetFireRate(fireRate)
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
	var newMeteors []*Meteor

	for i, meteor := range g.meteors.value {
		mRect := meteor.Collider()
		for j, bullet := range g.player.bullets {
			if !meteorHit[i] && !bulletHit[j] && mRect.Intersect(bullet.Collider()) {
				meteorHit[i] = true
				if !bullet.piercing {
					bulletHit[j] = true
				}
				g.Score++
				if meteor.canSplit {
					count := 2 + rand.Intn(3)
					for k := 0; k < count; k++ {
						split := newSplitMeteor(meteor.position, meteor.tier)
						if split != nil {
							newMeteors = append(newMeteors, split)
						}
					}
				}
			}
		}
	}

	survivingMeteors := make([]*Meteor, 0, len(g.meteors.value))
	for i, m := range g.meteors.value {
		if !meteorHit[i] {
			survivingMeteors = append(survivingMeteors, m)
		}
	}
	survivingMeteors = append(survivingMeteors, newMeteors...)
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

	ps := g.powerUps[:0]
	for _, pu := range g.powerUps {
		if !pu.IsOffScreen() {
			ps = append(ps, pu)
		}
	}
	g.powerUps = ps
}

func (g *Game) reset() {
	g.player = newPlayer()
	g.meteors = newMeteors()
	g.powerUps = nil
	g.Score = 0
	g.Lives = InitialLives
	g.DifficultyLevel = 0
	g.firstSpawnDone = false
	g.invulnTimer = nil
	g.flashCounter = 0
	g.confirmQuit = false
	g.elapsedTime = 0
	g.powerUpSpawnTimer = NewTimer(PowerUpSpawnPeriod)
	g.lastBronzeMilestone = 10
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
		g.drawPowerUps(s)
		if g.State == StatePaused {
			g.drawPauseOverlay(s)
		}
		if g.confirmQuit {
			g.drawConfirmOverlay(s)
		}
	}
}

func (g *Game) drawStartScreen(s *ebiten.Image) {
	text.Draw(s, "太空射击", UIFont, ScreenWidth/2-120, ScreenHeight/2-120, color.White)
	text.Draw(s, "按 回车 或 空格 开始游戏", UIFont, ScreenWidth/2-200, ScreenHeight/2-30, color.White)
	text.Draw(s, "WASD: 移动 | 鼠标左/右键: 旋转 | 空格: 射击 | P: 暂停 | ESC: 返回菜单", UIFont, ScreenWidth/2-460, ScreenHeight/2+60, color.RGBA{180, 180, 180, 255})
}

func (g *Game) drawPlayField(s *ebiten.Image) {
	isInvuln := g.invulnTimer != nil && !g.invulnTimer.IsReadyAttack()
	drawPlayer := !isInvuln || (g.flashCounter/6)%2 == 0

	if drawPlayer {
		g.player.Draw(s)
	}
	g.meteors.Draw(s)

	text.Draw(s, fmt.Sprintf("%06d", g.Score), ScoreFont, ScreenWidth/2-100, 50, color.White)
	text.Draw(s, fmt.Sprintf("生命: %d", g.Lives), UIFont, 50, 50, color.White)
	text.Draw(s, fmt.Sprintf("等级: %d", g.DifficultyLevel), UIFont, ScreenWidth-200, 50, color.White)
	text.Draw(s, fmt.Sprintf("时间: %.0fs", g.elapsedTime.Seconds()), UIFont, ScreenWidth/2-100, 80, color.RGBA{180, 180, 180, 255})

	buffY := 110
	if g.player.piercingActive {
		text.Draw(s, "穿透", UIFont, 50, buffY, color.RGBA{255, 215, 0, 255})
		buffY += 30
	}
	if g.player.extraSpread > 0 {
		text.Draw(s, fmt.Sprintf("散射 +%d", g.player.extraSpread), UIFont, 50, buffY, color.RGBA{205, 127, 50, 255})
		buffY += 30
	}
	if g.player.cdMultiplier < 1.0 {
		text.Draw(s, fmt.Sprintf("射速 x%.0f", 1.0/g.player.cdMultiplier), UIFont, 50, buffY, color.RGBA{192, 192, 192, 255})
	}
}

func (g *Game) drawPowerUps(s *ebiten.Image) {
	for _, pu := range g.powerUps {
		pu.Draw(s)
	}
}

func (g *Game) drawGameOverOverlay(s *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	s.DrawImage(overlay, nil)

	text.Draw(s, "游戏结束", UIFont, ScreenWidth/2-100, ScreenHeight/2-120, color.RGBA{255, 50, 50, 255})
	text.Draw(s, fmt.Sprintf("得分: %d  等级: %d", g.Score, g.DifficultyLevel), UIFont, ScreenWidth/2-180, ScreenHeight/2-40, color.White)
	text.Draw(s, "按 回车 或 空格 重新开始", UIFont, ScreenWidth/2-200, ScreenHeight/2+40, color.White)
}

func (g *Game) drawPauseOverlay(s *ebiten.Image) {
	text.Draw(s, "已暂停", UIFont, ScreenWidth/2-70, ScreenHeight/2-40, color.White)
}

func (g *Game) drawConfirmOverlay(s *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 140})
	s.DrawImage(overlay, nil)

	text.Draw(s, "确认返回菜单？", UIFont, ScreenWidth/2-130, ScreenHeight/2-50, color.White)
	text.Draw(s, "Enter 确认  |  ESC 取消", UIFont, ScreenWidth/2-200, ScreenHeight/2+10, color.RGBA{220, 220, 220, 255})
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("太空射击")
	g := &Game{
		player:            newPlayer(),
		meteors:           newMeteors(),
		State:             StateStart,
		powerUpSpawnTimer: NewTimer(PowerUpSpawnPeriod),
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
