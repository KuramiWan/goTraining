package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// WeaponType represents the player's weapon category.
type WeaponType int

const (
	WpnBlue  WeaponType = iota // default, balanced
	WpnGreen                    // faster, lower damage
	WpnRed                      // slower, higher damage
)

// WeaponInfo describes a weapon's properties.
type WeaponInfo struct {
	wType        WeaponType
	fireRate     time.Duration
	damage       int
	bulletSpeed  float64
	bulletSprites []*ebiten.Image
	spriteIdx    int
}

// CurrentWeapon returns the active weapon config for the player.
func (p *Player) CurrentWeapon() *WeaponInfo {
	if p.weapon == nil {
		return GetWeaponInfo(WpnBlue)
	}
	return p.weapon
}

// GetWeaponInfo returns weapon info for a given type.
func GetWeaponInfo(wType WeaponType) *WeaponInfo {
	switch wType {
	case WpnGreen:
		return &WeaponInfo{
			wType:        WpnGreen,
			fireRate:     600 * time.Millisecond,
			damage:       1,
			bulletSpeed:  1200,
			bulletSprites: LaserGreenSprites,
			spriteIdx:    0,
		}
	case WpnRed:
		return &WeaponInfo{
			wType:        WpnRed,
			fireRate:     2 * time.Second,
			damage:       3,
			bulletSpeed:  700,
			bulletSprites: LaserRedSprites,
			spriteIdx:    0,
		}
	default: // WpnBlue
		return &WeaponInfo{
			wType:        WpnBlue,
			fireRate:     1 * time.Second,
			damage:       1,
			bulletSpeed:  1000,
			bulletSprites: LaserBlueSprites,
			spriteIdx:    0,
		}
	}
}

// SwitchWeapon changes the player's weapon type.
func (p *Player) SwitchWeapon(wType WeaponType) {
	p.weapon = GetWeaponInfo(wType)
}

// UpgradeWeapon improves the current weapon's fire rate (up to 5 stacks).
func (p *Player) UpgradeWeapon() {
	if p.weaponUpgrades >= 5 {
		return
	}
	p.weaponUpgrades++
	mult := 1.0 - float64(p.weaponUpgrades)*0.1 // 0.9, 0.8, 0.7, 0.6, 0.5
	p.cdMultiplier *= mult
}
