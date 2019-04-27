package main

import (
	"fmt"
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
	itemBorder = color.RGBA{R: 0x7e, G: 0x7e, B: 0x7e, A: 0x00}
	itemSelectBorder = color.RGBA{R: 0x7e, G: 0x7e, B: 0x7e, A: 0xaa}
)

type Item struct {
	disabled bool
	highlighted bool
	cost float64
	gridPos pixel.Vec
	name string
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
	imd.Rectangle(0)

	imd.Color = itemBorder
	if i.highlighted {
		imd.Color = itemSelectBorder
	}
	imd.Push(box.Min, box.Max)
	imd.Rectangle(4)
}

type Shop struct {
	imd *imdraw.IMDraw
	items []*Item
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

func (s *Shop) Collides(pixel.Rect) bool {
	return false
}

func (s *Shop) Hurt(pixel.Rect) {}

func (s *Shop) Init(pixel.Rect) {
	s.imd = imdraw.New(nil)

	s.AddItem(10, "a")
	s.AddItem(20, "b")
	s.AddItem(300, "c")
	s.AddItem(10, "a")
	s.AddItem(20, "b")
	s.AddItem(10, "a")
	s.AddItem(20, "b")
	fmt.Println(s)
}

func (s *Shop) AddItem(cost float64, name string) {
	if len(s.items) >= buttonsAcross * buttonsDown {
		return
	}

	i := Item{
		name: name,
		cost: cost,
		gridPos: pixel.V(
			float64(len(s.items)%buttonsAcross),
			float64(len(s.items)/buttonsAcross),
			),
	}

	s.items = append(s.items, &i)
}

func (s *Shop) Draw(target pixel.Target) {
	s.imd.Clear()

	for _, i := range s.items {
		i.Draw(s.imd)
	}

	s.imd.Draw(target)
}

func (s *Shop) Start() {
	camPos = winBounds.Center()
}
