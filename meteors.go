package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Meteors struct {
	value            []*Meteor
	meteorSpawnTimer *Timer
}

func (ms *Meteors) Update(spawnNow bool) {
	if spawnNow {
		ms.value = append(ms.value, newMeteor())
	}
	ms.meteorSpawnTimer.UpdateTicks()
	if ms.meteorSpawnTimer.IsReadyAttack() {
		ms.meteorSpawnTimer.RestTicks()
		ms.value = append(ms.value, newMeteor())
	}
	for _, meteor := range ms.value {
		meteor.Update()
	}
}

func (ms *Meteors) Draw(s *ebiten.Image) {
	for _, meteor := range ms.value {
		meteor.Draw(s)
	}
}

func newMeteors() *Meteors {
	return &Meteors{
		value:            make([]*Meteor, 0),
		meteorSpawnTimer: NewTimer(5 * time.Second),
	}
}
