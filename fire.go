package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

const (
	minRate = 1.
	maxRate = 10.
)

var (
	fireSprites  []*pixel.Sprite
	smokeSprites []*pixel.Sprite

	fireCount  int
	smokeCount int

	fires     = make(map[int]*fire)
	fireIndex int

	fadeRate = 0.5
	growRate = 0.8
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
	inc(fireCount)

	f := fire{
		id:       fireIndex,
		pos:      pos,
		vel:      vel,
		rate:     (rand.Float64() * (maxRate - minRate)) + minRate,
		fireFade: 0.5,
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
}

func DrawFires(target pixel.Target) {
	for _, f := range fires {
		f.Draw(target)
	}
}
