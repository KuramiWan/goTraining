package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Meteors struct {
	array            []*Meteor
	meteorSpawnTimer *Timer
}

func (meteors *Meteors) Update() {
	meteors.meteorSpawnTimer.UpdateTicks()
	if meteors.meteorSpawnTimer.IsReadyAttack() {
		meteors.meteorSpawnTimer.RestTicks()
		meteor := newMeteor()
		meteors.array = append(meteors.array, meteor)
	}
	for _, meteor := range meteors.array {
		meteor.Update()
	}
}

func (meteors *Meteors) Draw(s *ebiten.Image) {
	for _, meteor := range meteors.array {
		meteor.Draw(s)
	}
}
func newMeteors() *Meteors {
	return &Meteors{make([]*Meteor, 0), NewTimer(5 * time.Second)}
}
