package main

import (
	"github.com/faiface/pixel"
)

type Player struct {
	health float64
	bounds pixel.Rect
	sprites []*pixel.Sprite
	offSet pixel.Vec
}

func NewPlayer() *Player {
	pic, err := loadPicture("assets/tilemap.png")
	if err != nil {
		panic(err)
	}

	p := Player{
		health: 100,
		bounds: pixel.R(-4, -4, 4, 4),
		sprites: []*pixel.Sprite{
			pixel.NewSprite(pic, pixel.R(0, 0, 16, 16)),
			pixel.NewSprite(pic, pixel.R(16, 0, 32, 16)),
			pixel.NewSprite(pic, pixel.R(32, 0, 48, 16)),
			pixel.NewSprite(pic, pixel.R(48, 0, 64, 16)),
		},
	}

	return &p
}

func (p *Player) CanMove(delta pixel.Vec) bool {
	return !lvlMan.CurrentLevel().Collides(p.bounds.Moved(p.offSet.Add(delta)))
}

func (p *Player) Update(dt float64, offset pixel.Vec) {
	p.offSet = offset
}

func (p *Player) Draw(target pixel.Target) {
	// TODO animate
	p.sprites[0].Draw(target, pixel.IM.Moved(p.offSet))
}
