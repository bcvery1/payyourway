package main

import "github.com/faiface/pixel"

type Animator interface {
	Sprites() []*pixel.Sprite
	PreviousSpriteInd() int
	Tickers() (sinceLast float64, rate float64)
	SetLast(float64)
}

func Sprite(a Animator, dt float64) int {
	since, rate := a.Tickers()
	since += dt

	if since < rate {
		a.SetLast(since)
		return a.PreviousSpriteInd()
	}

	a.SetLast(0)
	ind := a.PreviousSpriteInd()
	if ind+1 == len(a.Sprites()) {
		return 0
	}

	return a.PreviousSpriteInd() + 1
}
