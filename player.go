package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Player struct {
	health float64
	maxHealth float64
	bounds pixel.Rect
	sprites []*pixel.Sprite
	offSet pixel.Vec
	imd *imdraw.IMDraw
	hitFade uint8
}

func NewPlayer() *Player {
	pic, err := loadPicture("assets/tilemap.png")
	if err != nil {
		panic(err)
	}

	p := Player{
		health: 100,
		maxHealth: 100,
		bounds: pixel.R(-4, -4, 4, 4),
		sprites: []*pixel.Sprite{
			pixel.NewSprite(pic, pixel.R(0, 0, 16, 16)),
			pixel.NewSprite(pic, pixel.R(16, 0, 32, 16)),
			pixel.NewSprite(pic, pixel.R(32, 0, 48, 16)),
			pixel.NewSprite(pic, pixel.R(48, 0, 64, 16)),
		},
		imd: imdraw.New(nil),
		hitFade: 255,
	}

	return &p
}

func (p *Player) CanMove(delta pixel.Vec) bool {
	return !lvlMan.CurrentLevel().Collides(p.bounds.Moved(p.offSet.Add(delta)))
}

func (p *Player) Update(dt float64, offset pixel.Vec) {
	p.offSet = offset
	if p.hitFade < 255 {
		hf := int(p.hitFade)
		hf += int(500*dt)
		if hf > 255 {
			p.hitFade = 255
		} else {
			p.hitFade = uint8(hf)
		}
	}
}

func (p *Player) Draw(win *pixelgl.Window) {
	// TODO animate
	p.sprites[0].Draw(win, pixel.IM.Moved(p.offSet))

	p.drawHUD(win)

	if p.hitFade < 255 {
		win.SetColorMask(color.RGBA{R: p.hitFade, B: 0x00, G: 0x00, A: 0x00})
	} else {
		win.SetColorMask(nil)
	}
}

func (p *Player) drawHUD(target pixel.Target) {
	p.imd.Clear()

	startV := winBounds.Center().Add(p.offSet).Sub(pixel.V(70, 240))
	size := pixel.V(30, 200)

	// Health backing
	p.imd.Color = color.Black
	p.imd.Push(
		startV,
		startV.Add(size),
		)
	p.imd.Rectangle(0)

	// Health indicator
	hSize := size.ScaledXY(pixel.V(1, p.health/p.maxHealth))
	p.imd.Color = color.RGBA{R: 0x00, G: 0xff, B:0x00, A:0x00}
	p.imd.Push(
		startV,
		startV.Add(hSize),
	)
	p.imd.Rectangle(0)

	// health container
	p.imd.Color = pixel.RGB(0x71, 0x25, 0x16)
	p.imd.Push(
		startV,
		startV.Add(size),
		)
	p.imd.Rectangle(2)

	p.imd.Draw(target)
}

func (p *Player) Hurt(hp float64) {
	if p.hitFade < 255 {
		return
	}
	p.health -= hp
	p.hitFade = 150

	if p.health <= 0 {
		p.Die()
	}
}

func (p *Player) Die() {
	panic("You died")
}
