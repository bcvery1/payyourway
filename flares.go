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
	id               int
	pos              pixel.Vec
	life             float64
	aniRate          float64
	aniSince         float64
	currentSpriteInd int
}

func (f *flare) Sprites() []*pixel.Sprite {
	return flareSprites
}

func (f *flare) PreviousSpriteInd() int {
	return f.currentSpriteInd
}

func (f *flare) Tickers() (sinceLast float64, rate float64) {
	return f.aniSince, f.aniRate
}

func (f *flare) SetLast(last float64) {
	f.aniSince = last
}

func NewFlare(pos pixel.Vec) {
	f := flare{
		id:      flareCount,
		pos:     pos,
		life:    250,
		aniRate: 0.04,
	}

	flareCount++

	flares[f.id] = &f
}

func (f *flare) update(dt float64) {
	f.life -= dt * flareDieRate

	if f.life < 0 {
		delete(flares, f.id)
		NewSmoke(f.pos, pixel.ZV)
	}

	f.currentSpriteInd = Sprite(f, dt)
}

func (f *flare) draw(target pixel.Target) {
	flareSprites[f.currentSpriteInd].Draw(target, pixel.IM.Moved(f.pos))
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
