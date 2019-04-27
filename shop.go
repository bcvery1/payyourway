package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
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
	itemDisabledColour = color.RGBA{R: 0x55, G: 0x55, B: 0x55, A: 0x00}

	itemBorder = color.RGBA{R: 0x7e, G: 0x7e, B: 0x7e, A: 0x00}
	itemSelectBorder = color.RGBA{R: 0x7e, G: 0x7e, B: 0x7e, A: 0xaa}
)

type Item struct {
	disabled bool
	highlighted bool
	cost float64
	gridPos pixel.Vec
	name string
	text *text.Text
}

func (i *Item) winPos() pixel.Rect {
	min := pixel.V(itemWidth, itemHeight).ScaledXY(i.gridPos).Add(pixel.V(15, 65))
	return pixel.Rect{
		Min: min,
		Max: min.Add(itemSize),
	}
}

func (i *Item) Buy() {
	player.Hurt(i.cost)

	switch i.name {
	case "Light Shield":
		player.shield += 20
		player.maxShield += 20
	case "Shield":
		player.shield += 40
		player.maxShield += 40
	case "Heavy Shield":
		player.shield += 60
		player.maxShield+= 60
	case "Max HP Boost":
		player.maxHealth += 20
	case "Flares":
		player.inventory = append(player.inventory, *i)
	default:
		fmt.Println(i.name)
	}
}

func (i *Item) Draw (imd *imdraw.IMDraw, target pixel.Target) {
	imd.Color = itemColour
	if i.highlighted {
		imd.Color = itemSelectColour
	}
	if i.disabled {
		imd.Color = itemDisabledColour
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

	i.text.Draw(target, pixel.IM.Scaled(i.text.Orig, 2))
}

type Shop struct {
	imd *imdraw.IMDraw
	items []*Item
	nextLevel int
	returnText *text.Text
}

func (s *Shop) Update(dt float64, win *pixelgl.Window) {
	mousePos := win.MousePosition()

	for j, i := range s.items {
		if i.cost >= player.health {
			i.disabled = true
		}

		i.highlighted = i.winPos().Contains(mousePos) && !i.disabled

		if win.JustPressed(pixelgl.MouseButtonLeft) && i.highlighted {
			i.Buy()
			s.items = append(s.items[:j], s.items[j+1:]...)
		}
	}

	if win.JustPressed(pixelgl.MouseButtonLeft) && pixel.R(winBounds.W()-120, 0, winBounds.W(), 50).Contains(mousePos) {
		player.health = player.maxHealth
		lvlMan.StartLevel(s.nextLevel)
		s.ClearItems()
	}
}

func (s *Shop) Collides(pixel.Rect) bool {
	return false
}

func (s *Shop) Hurt(pixel.Rect) {}

func (s *Shop) Init(pixel.Rect) {
	s.imd = imdraw.New(nil)

	s.returnText = text.New(pixel.V(winBounds.W()-120, 5), atlas)
	s.returnText.Color = color.White
	_, _ = fmt.Fprint(s.returnText, "Return")
}

func (s *Shop) AddItem(cost float64, name, desc string) {
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

	i.text = text.New(i.winPos().Min.Add(pixel.V(20, itemHeight-30)), atlas)
	i.text.Color = color.RGBA{R: 0xbb, G: 0xc5, B: 0xda, A: 0x00}
	_, _ = fmt.Fprintf(i.text, "%s\n\nCost (hp): %.2f", i.name, i.cost)

	s.items = append(s.items, &i)
}

func (s *Shop) Draw(win *pixelgl.Window) {
	s.imd.Clear()

	for _, i := range s.items {
		i.Draw(s.imd, win)
	}

	s.returnText.Draw(win, pixel.IM.Scaled(s.returnText.Orig, 2))

	s.imd.Draw(win)
}

func (s *Shop) Start() {
	camPos = winBounds.Center()
}

func (s *Shop) ClearItems() {
	s.items = []*Item{}
}

func (s *Shop) ReachedShop() string {
	return ""
}

func (s *Shop) Setup(shopName string) {
	switch shopName {
	case "Main":
		if lvlMan.previousLevel == Level4Ind {
			lvlMan.StartLevel(EndInd)
		}
		s.nextLevel = lvlMan.previousLevel+1
	case "FirstLevelMid":
		s.nextLevel = Level1Ind

		s.AddItem(45, "Light Shield", "Provides 20 points of protection")
		s.AddItem(45, "Light Shield", "Provides 20 points of protection")
		s.AddItem(99, "Max HP Boost", "Permanently adds 20 extra HP")
		s.AddItem(150, "Flares", "One use deploy flares")
	case "SecondLevelMid":
		s.nextLevel = Level2Ind
	case "ThirdLevelMid":
		s.nextLevel = Level3Ind
	default:
		panic("Unrecognised shop name " + shopName)
	}
}
