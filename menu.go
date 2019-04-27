package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

const (
	titleText        = "Pay Your Way"
	instructionsText = "Welcome to 'Pay Your Way', in which you must get through all four\nlevels with as much health" +
		" at the end as you can.  You start with\n100 health, and can buy upgrades using that at shops. Leaving" +
		" a\nshop will take you back to the beginning of the level with full\nhealth.  Reach the center shop to" +
		" advance to the next level.\nReach the center shop from each level to complete the\ngame\n\nPress space to start!"
)

type Menu struct {
	title        *text.Text
	instructions *text.Text
}

func (m *Menu) Init(pixel.Rect) {
	m.title = text.New(pixel.V(350, winBounds.H()-120), atlas)
	m.title.Color = color.RGBA{R: 0x00, G: 0x2f, B: 0x5f, A: 0xff}
	_, _ = fmt.Fprint(m.title, titleText)

	m.instructions = text.New(pixel.V(50, winBounds.H()-180), atlas)
	m.instructions.Color = color.RGBA{R: 0x8e, G: 0xb4, B: 0xda, A: 0xff}
	_, _ = fmt.Fprint(m.instructions, instructionsText)

}

func (m *Menu) Start() {
	camPos = winBounds.Center()
}

func (m *Menu) Update(dt float64, win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeySpace) {
		lvlMan.StartLevel(Level1Ind)
	}
}

func (m *Menu) Draw(win *pixelgl.Window) {
	m.title.Draw(win, pixel.IM.Scaled(m.title.Orig, 5))
	m.instructions.Draw(win, pixel.IM.Scaled(m.instructions.Orig, 2))
}

func (m *Menu) Collides(pixel.Rect) bool {
	return false
}

func (m *Menu) Hurt(pixel.Rect) {}

func (m *Menu) ReachedShop() string {
	return ""
}
