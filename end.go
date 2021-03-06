package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type End struct {
	text *text.Text
}

func (e *End) Init(pixel.Rect) {
	e.text = text.New(winBounds.Center().Sub(pixel.V(390, -150)), atlas)
}

func (e *End) Start() {
	camPos = winBounds.Center()
}

func (e *End) Update(float64, *pixelgl.Window) {
}

func (e *End) Draw(win *pixelgl.Window) {
	e.text.Clear()
	_, _ = fmt.Fprintf(e.text, "\t\tThe\n\t\tEnd\n\nYou got %.0f points", player.health)
	e.text.Draw(win, pixel.IM.Scaled(e.text.Orig, 6))
}

func (e *End) Collides(pixel.Rect) bool {
	return false
}

func (e *End) Hurt(pixel.Rect) {

}

func (e *End) ReachedShop() string {
	return ""
}
