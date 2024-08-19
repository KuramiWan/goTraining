package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Meteors struct {
	value            []*Meteor
	meteorSpawnTimer *Timer
}

const (
	//first meteors generate
	firstMeteors = 1 * time.Second
)

func (meteors *Meteors) Update(t time.Duration) {
	if t == firstMeteors {
		meteor := newMeteor()
		meteors.value = append(meteors.value, meteor)
	}
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

func (m *Meteor) Collider() *Rect {
	return newRect(m.position, m.sprite)
}
