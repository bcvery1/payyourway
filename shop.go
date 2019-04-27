package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	buttonsAcross = 4
	buttonsDown = 3
)

var (
	itemWidth = (winBounds.Max.X-30)/buttonsAcross
	itemHeight = (winBounds.Max.Y-130)/buttonsDown
	itemSize = pixel.V(itemWidth, itemHeight)

	itemColour = color.RGBA{R: 0x5c, G: 0x72, B: 0x7e, A: 0x00}
	itemSelectColour = color.RGBA{R: 0x80, G: 0xb9, B: 0xda, A: 0x00}
)

type Item struct {
	disabled bool
	highlighted bool
	cost float64
	gridPos pixel.Vec
}

func (i *Item) winPos() pixel.Rect {
	min := pixel.V(itemWidth, itemHeight).ScaledXY(i.gridPos).Add(pixel.V(15, 65))
	return pixel.Rect{
		Min: min,
		Max: min.Add(itemSize),
	}
}

func (i *Item) Draw (imd *imdraw.IMDraw) {
	imd.Color = itemColour
	if i.highlighted {
		imd.Color = itemSelectColour
	}

	box := i.winPos()
	imd.Push(box.Min, box.Max)
}

type Shop struct {
	imd *imdraw.IMDraw
	items []Item
}

func (s *Shop) Update(dt float64, win *pixelgl.Window) {
	mousePos := win.MousePosition()

	for _, i := range s.items {
		if i.cost > player.health {
			i.disabled = true
		}

		i.highlighted = i.winPos().Contains(mousePos)
	}
}

func (s *Shop) Collides(pixel.Rect) bool {}

func (s *Shop) Hurt(pixel.Rect) {}

func (s *Shop) Init(pixel.Rect) {
	camPos = pixel.ZV

	s.imd = imdraw.New(nil)
}

func (s *Shop) Draw(target pixel.Target) {
	s.imd.Clear()

	for _, i := range s.items {
		i.Draw(s.imd)
	}

	s.imd.Draw(target)
}

func (s *Shop) Start() {
}
