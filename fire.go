package main

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
)

const (
	minRate = -10.
	maxRate = 10.
)

var (
	fireSprites  []*pixel.Sprite
	smokeSprites []*pixel.Sprite

	fireCount  int
	smokeCount int

	fires     = make(map[int]*fire)
	fireIndex int

	smokes     = make(map[int]*Smoke)
	smokeIndex int

	fadeRate = 0.9
	growRate = 0.6

	initFireFade = 1.
)

// inc relies on having the same number of fire sprites as smoke sprites
func inc(count int) {
	count += 1
	if count > len(fireSprites) {
		count = 0
	}
}

type fire struct {
	id       int
	pos      pixel.Vec
	vel      pixel.Vec
	ind      int
	angle    float64
	rate     float64
	fireFade float64
	scale    float64
}

func NewFire(pos, vel pixel.Vec) {
	fmt.Println(vel)
	inc(fireCount)

	f := fire{
		id:       fireIndex,
		pos:      pos,
		vel:      vel,
		rate:     (rand.Float64() * (maxRate - minRate)) + minRate,
		fireFade: initFireFade,
		scale:    1,
	}

	fires[f.id] = &f

	fireIndex++
}

func (f *fire) Update(dt float64) {
	f.pos = f.pos.Add(f.vel.Scaled(dt))

	f.angle = f.angle + f.rate*dt

	f.scale += dt * growRate

	f.fireFade -= dt * fadeRate

	if f.fireFade < initFireFade*0.99 {
		NewSmoke(f.pos, f.vel)
	}

	if f.fireFade < 0 {
		delete(fires, f.id)
	}
}

func (f *fire) Draw(target pixel.Target) {
	fireSprites[f.ind].DrawColorMask(target, pixel.IM.Moved(f.pos).Rotated(f.pos, f.angle).Scaled(f.pos, f.scale), pixel.Alpha(f.fireFade))
}

func UpdateFires(dt float64) {
	for _, f := range fires {
		f.Update(dt)
	}

	for _, s := range smokes {
		s.Update(dt)
	}
}

func DrawFires(target pixel.Target) {
	for _, f := range fires {
		f.Draw(target)
	}

	for _, s := range smokes {
		s.Draw(target)
	}
}

type Smoke struct {
	id        int
	pos       pixel.Vec
	vel       pixel.Vec
	ind       int
	angle     float64
	rate      float64
	smokeFade float64
	scale     float64
}

func NewSmoke(pos, vel pixel.Vec) {
	inc(smokeCount)

	s := Smoke{
		id:        smokeIndex,
		pos:       diviate(pos),
		vel:       diviate(vel.Scaled(10)),
		rate:      ((rand.Float64() * (maxRate - minRate)) + minRate) / 2,
		smokeFade: 1.2,
		scale:     2,
	}

	smokeIndex++

	smokes[s.id] = &s
}

func (f *Smoke) Update(dt float64) {
	f.pos = f.pos.Add(f.vel.Scaled(dt))

	f.angle = f.angle + f.rate*dt

	f.scale += dt * growRate * 2

	f.smokeFade -= dt * fadeRate * 25

	if f.smokeFade < 0 {
		delete(smokes, f.id)
	}
}

func (f *Smoke) Draw(target pixel.Target) {
	smokeSprites[f.ind].DrawColorMask(target, pixel.IM.Moved(f.pos).Rotated(f.pos, f.angle).Scaled(f.pos, f.scale), pixel.Alpha(f.smokeFade))
}

func diviate(v pixel.Vec) pixel.Vec {
	return pixel.V(
		v.X+float64n(-0.2, 0.2),
		v.Y+float64n(-0.2, 0.2),
	)
}

func float64n(min, max float64) float64 {
	return (rand.Float64() * (max - min)) + min
}
