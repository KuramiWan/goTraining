package main

import (
	"testing"
)

func TestWeaponGetInfo(t *testing.T) {
	w := GetWeaponInfo(WpnBlue)
	if w == nil {
		t.Fatal("Blue weapon should not be nil")
	}
	if w.wType != WpnBlue {
		t.Errorf("type want Blue, got %d", w.wType)
	}
	if w.damage != 1 {
		t.Errorf("blue damage want 1, got %d", w.damage)
	}

	g := GetWeaponInfo(WpnGreen)
	if g == nil {
		t.Fatal("Green weapon should not be nil")
	}
	if g.damage != 1 {
		t.Errorf("green damage want 1, got %d", g.damage)
	}

	r := GetWeaponInfo(WpnRed)
	if r == nil {
		t.Fatal("Red weapon should not be nil")
	}
	if r.damage != 3 {
		t.Errorf("red damage want 3, got %d", r.damage)
	}
}

func TestWeaponSwitch(t *testing.T) {
	p := newPlayer()
	orig := p.weapon.wType
	if orig != WpnBlue {
		t.Fatalf("default weapon should be Blue, got %d", orig)
	}

	p.SwitchWeapon(WpnRed)
	if p.weapon.wType != WpnRed {
		t.Errorf("after switch, want Red, got %d", p.weapon.wType)
	}

	p.SwitchWeapon(WpnGreen)
	if p.weapon.wType != WpnGreen {
		t.Errorf("after switch, want Green, got %d", p.weapon.wType)
	}
}

func TestWeaponUpgrade(t *testing.T) {
	p := newPlayer()
	origMult := p.cdMultiplier
	if origMult != 1.0 {
		t.Fatalf("default cdMult want 1.0, got %f", origMult)
	}

	p.UpgradeWeapon()
	if p.weaponUpgrades != 1 {
		t.Errorf("upgrades want 1, got %d", p.weaponUpgrades)
	}
	if p.cdMultiplier >= origMult {
		t.Error("cdMultiplier should decrease after upgrade")
	}

	// Upgrade to max (5)
	for i := 0; i < 4; i++ {
		p.UpgradeWeapon()
	}
	if p.weaponUpgrades != 5 {
		t.Errorf("upgrades want 5, got %d", p.weaponUpgrades)
	}
	// Should cap at 5
	p.UpgradeWeapon()
	if p.weaponUpgrades != 5 {
		t.Errorf("upgrades should cap at 5, got %d", p.weaponUpgrades)
	}
}

func TestWeaponPlayerDefault(t *testing.T) {
	p := newPlayer()
	w := p.CurrentWeapon()
	if w == nil {
		t.Fatal("CurrentWeapon returned nil for default player")
	}
	if w.wType != WpnBlue {
		t.Errorf("default weapon type want Blue, got %d", w.wType)
	}
}
