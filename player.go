package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Player struct {
	health     float64
	maxHealth  float64
	bounds     pixel.Rect
	sprites    []*pixel.Sprite
	offSet     pixel.Vec
	imd        *imdraw.IMDraw
	hitFade    uint8
	drownFade  uint8
	inventory  []Item
	shield     float64
	maxShield  float64
	boatHealth float64
}

func NewPlayer() *Player {
	p := Player{
		health:    200,
		maxHealth: 200,
		bounds:    pixel.R(-4, -4, 4, 4),
		sprites: []*pixel.Sprite{
			pixel.NewSprite(tilemapPic, pixel.R(0, 0, 16, 16)),
			pixel.NewSprite(tilemapPic, pixel.R(16, 0, 32, 16)),
			pixel.NewSprite(tilemapPic, pixel.R(32, 0, 48, 16)),
			pixel.NewSprite(tilemapPic, pixel.R(48, 0, 64, 16)),
		},
		imd:     imdraw.New(nil),
		hitFade: 255,
		inventory: []Item{
			{name: "Flares"},
		},
	}

	return &p
}

func (p *Player) useItem(item Item) {
	switch item.name {
	case "Boat":
		Announce("Using boat")
		p.boatHealth += 50
	case "Flares":
		Announce("Deployed flares")
		NewFlare(p.CollisionBox().Center())
	}
}

func (p *Player) CanMove(delta pixel.Vec) bool {
	return !lvlMan.CurrentLevel().Collides(p.CollisionBox().Moved(delta))
}

func (p *Player) Update(dt float64, offset pixel.Vec) {
	p.offSet = offset

	if p.hitFade < 255 {
		hf := int(p.hitFade)
		hf += int(500 * dt)
		if hf > 255 {
			p.hitFade = 255
		} else {
			p.hitFade = uint8(hf)
		}
	}

	if p.drownFade < 255 {
		hf := int(p.drownFade)
		hf += int(500 * dt)
		if hf > 255 {
			p.drownFade = 255
		} else {
			p.drownFade = uint8(hf)
		}
	}
}

func (p *Player) Draw(win *pixelgl.Window) {
	// TODO animate
	p.sprites[0].Draw(win, pixel.IM.Moved(p.offSet))

	p.drawHUD(win)
	p.drawInv(win)

	if p.hitFade < 255 {
		win.SetColorMask(color.RGBA{R: p.hitFade, G: 0x44, B: 0x44, A: 0x00})
	} else if p.drownFade < 255 {
		win.SetColorMask(color.RGBA{R: 0x44, G: 0x44, B: p.drownFade, A: 0x00})
	} else {
		win.SetColorMask(nil)
	}
}

func (p *Player) drawInv(target pixel.Target) {
	for i, item := range p.inventory {
		j := i + 1
		if j > 9 {
			break
		}

		t := text.New(pixel.V(16*float64(j), 8), atlas)
		_, _ = fmt.Fprint(t, j)

		t.Draw(target, pixel.IM.Moved(p.offSet.Sub(winBounds.Center()).Add(t.Orig)))

		s, ok := invSprites[item.name]
		if !ok {
			continue
		}

		shift := pixel.V(32*float64(j), 40)
		s.Draw(target, pixel.IM.Moved(p.offSet.Sub(winBounds.Center()).Add(shift)))
	}
}

func (p *Player) drawHUD(target pixel.Target) {
	p.imd.Clear()

	startV := winBounds.Center().Add(p.offSet).Sub(pixel.V(25, 210))
	size := pixel.V(20, 200)

	// Health backing
	p.imd.Color = color.Black
	p.imd.Push(
		startV,
		startV.Add(size),
	)
	p.imd.Rectangle(0)

	// Health indicator
	hSize := size.ScaledXY(pixel.V(1, p.health/p.maxHealth))
	p.imd.Color = color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0x00}
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

	if p.shield > 0 {
		// Shield
		startV := winBounds.Center().Add(p.offSet).Sub(pixel.V(55, 210))

		// Shield backing
		p.imd.Color = color.Black
		p.imd.Push(
			startV,
			startV.Add(size),
		)
		p.imd.Rectangle(0)

		// Shield indicator
		hSize := size.ScaledXY(pixel.V(1, p.shield/p.maxShield))
		p.imd.Color = color.RGBA{R: 0x43, G: 0x6d, B: 0xda, A: 0x00}
		p.imd.Push(
			startV,
			startV.Add(hSize),
		)
		p.imd.Rectangle(0)

		// Shield container
		p.imd.Color = pixel.RGB(0x71, 0x25, 0x16)
		p.imd.Push(
			startV,
			startV.Add(size),
		)
		p.imd.Rectangle(2)
	}

	p.imd.Draw(target)
}

func (p *Player) Hurt(hp float64) {
	PlaySound(hurtSound)

	if p.hitFade < 255 {
		return
	}

	p.hitFade = 150

	excess := hp - p.shield
	p.shield -= hp
	if excess > 0 {
		p.shield = 0
		p.maxShield = 0
		hp = excess
	} else {
		return
	}

	p.health -= hp

	if p.health <= 0 {
		p.Die()
	}
}

func (p *Player) Drown(hp float64) {
	PlaySound(hurtSound)

	if p.drownFade < 255 {
		return
	}

	p.drownFade = 150

	if p.boatHealth > 0 {
		p.boatHealth -= hp
		if p.boatHealth < 0 {
			Announce("The boat broke")
			p.boatHealth = 0
		} else {
			return
		}
	}

	p.health -= hp

	if p.health <= 0 {
		p.Die()
	}
}

func (p *Player) Die() {
	// TODO add death screen with timeout
	p.health = p.maxHealth

	lvlMan.RestartLevel()
}

func (p *Player) CollisionBox() pixel.Rect {
	return p.bounds.Moved(p.offSet)
}

func (p *Player) UpdateInventory(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.Key1) {
		if len(p.inventory) < 1 {
			return
		}
		p.useItem(p.inventory[0])
		p.inventory = p.inventory[1:]

	} else if win.JustPressed(pixelgl.Key2) {
		if len(p.inventory) < 2 {
			return
		}
		p.useItem(p.inventory[1])
		p.inventory = append(p.inventory[:1], p.inventory[2:]...)

	} else if win.JustPressed(pixelgl.Key3) {
		if len(p.inventory) < 3 {
			return
		}
		p.useItem(p.inventory[2])
		p.inventory = append(p.inventory[:2], p.inventory[3:]...)

	} else if win.JustPressed(pixelgl.Key4) {
		if len(p.inventory) < 4 {
			return
		}
		p.useItem(p.inventory[3])
		p.inventory = append(p.inventory[:3], p.inventory[4:]...)

	} else if win.JustPressed(pixelgl.Key5) {
		if len(p.inventory) < 5 {
			return
		}
		p.useItem(p.inventory[4])
		p.inventory = append(p.inventory[:4], p.inventory[5:]...)

	} else if win.JustPressed(pixelgl.Key6) {
		if len(p.inventory) < 6 {
			return
		}
		p.useItem(p.inventory[5])
		p.inventory = append(p.inventory[:5], p.inventory[6:]...)

	} else if win.JustPressed(pixelgl.Key7) {
		if len(p.inventory) < 7 {
			return
		}
		p.useItem(p.inventory[6])
		p.inventory = append(p.inventory[:6], p.inventory[7:]...)

	} else if win.JustPressed(pixelgl.Key8) {
		if len(p.inventory) < 8 {
			return
		}
		p.useItem(p.inventory[7])
		p.inventory = append(p.inventory[:7], p.inventory[8:]...)

	} else if win.JustPressed(pixelgl.Key9) {
		if len(p.inventory) < 9 {
			return
		}
		p.useItem(p.inventory[8])
		p.inventory = p.inventory[:8]
	}
}
