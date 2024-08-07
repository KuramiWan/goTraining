package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Meteors struct {
	value            []*Meteor
	meteorSpawnTimer *Timer
}

func (meteors *Meteors) Update() {
	meteors.meteorSpawnTimer.UpdateTicks()
	if meteors.meteorSpawnTimer.IsReadyAttack() {
		meteors.meteorSpawnTimer.RestTicks()
		meteor := newMeteor()
		meteors.value = append(meteors.value, meteor)
	}
	for _, meteor := range meteors.value {
		meteor.Update()
	}
}
func (meteors *Meteors) Draw(s *ebiten.Image) {
	for _, meteor := range meteors.value {
		meteor.Draw(s)
	}
}
func newMeteors() *Meteors {
	return &Meteors{make([]*Meteor, 0), NewTimer(5 * time.Second)}
}
