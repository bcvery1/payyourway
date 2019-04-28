package main

import (
	"github.com/faiface/pixel"
)

var (
	flares = make(map[int]*flare)

	flareCount   int
	flareDieRate = 30.

	flareSprites []*pixel.Sprite
)

type flare struct {
	id   int
	pos  pixel.Vec
	life float64
}

func NewFlare(pos pixel.Vec) {
	f := flare{
		id:   flareCount,
		pos:  pos,
		life: 250,
	}

	flareCount++

	flares[f.id] = &f
}

func (f *flare) update(dt float64) {
	f.life -= dt * flareDieRate

	if f.life < 0 {
		delete(flares, f.id)
	}
}

func (f *flare) draw(target pixel.Target) {
	flareSprites[0].Draw(target, pixel.IM.Moved(f.pos))
}

func UpdateFlares(dt float64) {
	for _, f := range flares {
		f.update(dt)
	}
}

func DrawFlares(target pixel.Target) {
	for _, f := range flares {
		f.draw(target)
	}
}
