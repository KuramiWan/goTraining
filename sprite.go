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

// Laser sprites by color (for weapon system)
var LaserBlueSprites = mustLoadFilteredImages("assets/PNG/Lasers", "Blue")
var LaserGreenSprites = mustLoadFilteredImages("assets/PNG/Lasers", "Green")
var LaserRedSprites = mustLoadFilteredImages("assets/PNG/Lasers", "Red")

// Power-up sprites
var PowerUpSilverSprite = mustLoadImage("PNG/Power-ups/bold_silver.png")
var PowerUpBronzeSprite = mustLoadImage("PNG/Power-ups/bolt_bronze.png")
var PowerUpGoldSprite = mustLoadImage("PNG/Power-ups/bolt_gold.png")
var PowerUpShieldSprite = mustLoadImage("PNG/Power-ups/shield_silver.png")
var PowerUpHealthSprite = mustLoadImage("PNG/Power-ups/pill_red.png")
var PowerUpSpeedSprite = mustLoadImage("PNG/Effects/speed.png")

// Fonts
var ScoreFont = mustLoadFont("assets/font/Kenney Blocks.ttf")
var UIFont = mustLoadSystemFont(`C:\Windows\Fonts\simhei.ttf`)

// Meteor sprites by tier
var MeteorBigSprites = mustLoadFilteredImages("assets/PNG/Meteors", "big")
var MeteorMedSprites = mustLoadFilteredImages("assets/PNG/Meteors", "med")
var MeteorSmallSprites = mustLoadFilteredImages("assets/PNG/Meteors", "small")
var MeteorTinySprites = mustLoadFilteredImages("assets/PNG/Meteors", "tiny")

// Enemy sprites
var EnemyBlackSprites = mustLoadFilteredImages("assets/PNG/Enemies", "Black")
var EnemyBlueSprites = mustLoadFilteredImages("assets/PNG/Enemies", enemyColorPattern("Blue"))
var EnemyGreenSprites = mustLoadFilteredImages("assets/PNG/Enemies", enemyColorPattern("Green"))
var EnemyRedSprites = mustLoadFilteredImages("assets/PNG/Enemies", enemyColorPattern("Red"))

// Boss sprites (UFOs)
var BossSprites = mustLoadFilteredImages("assets/PNG", "ufo")
var UfoBlueSprite = mustLoadImage("PNG/ufoBlue.png")
var UfoRedSprite = mustLoadImage("PNG/ufoRed.png")
var UfoGreenSprite = mustLoadImage("PNG/ufoGreen.png")
var UfoYellowSprite = mustLoadImage("PNG/ufoYellow.png")

// Effect sprites
var ExplosionSprites = mustLoadFilteredImages("assets/PNG/Effects", "fire")
var ShieldSprites = mustLoadFilteredImages("assets/PNG/Effects", "shield")
var StarSprites = mustLoadFilteredImages("assets/PNG/Effects", "star")
var Star1Sprite = mustLoadImage("PNG/Effects/star1.png")
var Star2Sprite = mustLoadImage("PNG/Effects/star2.png")
var Star3Sprite = mustLoadImage("PNG/Effects/star3.png")
var SpeedSprite = mustLoadImage("PNG/Effects/speed.png")

// Background sprites
var BgBlack = mustLoadImage("Backgrounds/black.png")
var BgBlue = mustLoadImage("Backgrounds/blue.png")
var BgDarkPurple = mustLoadImage("Backgrounds/darkPurple.png")
var BgPurple = mustLoadImage("Backgrounds/purple.png")

// UI sprites
var NumeralSprites = loadNumeralSprites()
var LifeBlueSprite = mustLoadImage("PNG/UI/playerLife1_blue.png")

// Player variant sprites
var PlayerShip1Blue = mustLoadImage("Sprite/playerShip1_blue.png")
var PlayerShip1Green = mustLoadImage("PNG/playerShip1_green.png")
var PlayerShip1Orange = mustLoadImage("PNG/playerShip1_orange.png")
var PlayerShip1Red = mustLoadImage("PNG/playerShip1_red.png")
var PlayerShip2Blue = mustLoadImage("PNG/playerShip2_blue.png")
var PlayerShip2Green = mustLoadImage("PNG/playerShip2_green.png")
var PlayerShip2Orange = mustLoadImage("PNG/playerShip2_orange.png")
var PlayerShip2Red = mustLoadImage("PNG/playerShip2_red.png")
var PlayerShip3Blue = mustLoadImage("PNG/playerShip3_blue.png")
var PlayerShip3Green = mustLoadImage("PNG/playerShip3_green.png")
var PlayerShip3Orange = mustLoadImage("PNG/playerShip3_orange.png")
var PlayerShip3Red = mustLoadImage("PNG/playerShip3_red.png")

// Damage sprites
var DamageSprites = mustLoadFilteredImages("assets/PNG/Damage", "damage")

// enemyColorPattern returns the exact case-sensitive pattern for enemy sprites.
func enemyColorPattern(color string) string {
	return color
}

func loadNumeralSprites() []*ebiten.Image {
	names := []string{
		"numeral0.png", "numeral1.png", "numeral2.png", "numeral3.png", "numeral4.png",
		"numeral5.png", "numeral6.png", "numeral7.png", "numeral8.png", "numeral9.png",
	}
	var sprites []*ebiten.Image
	for _, name := range names {
		sprite := mustLoadImageSilent(path.Join("PNG/UI", name))
		sprites = append(sprites, sprite)
	}
	return sprites
}

// mustLoadImageSilent loads an image but returns nil on error.
func mustLoadImageSilent(n string) *ebiten.Image {
	p, err := data.Open(path.Join("assets", n))
	if err != nil {
		return nil
	}
	defer p.Close()
	img, _, err := image.Decode(p)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

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
