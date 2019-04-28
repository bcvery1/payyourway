package main

import (
	"github.com/faiface/pixel"
)

const (
	maxSafeDist = 48
	midSafeDist = 32
	minSafeDist = 16

	explodeDist = 8

	gunFireDist = 200
)

var (
	guns = make([]*Gun, 0)

	rockets     = make(map[int]*rocket)
	rocketCount = 0

	rocketAcc = .5

	rocketSprite *pixel.Sprite
)

func SetupGuns() {
	for _, p := range tmxMap.GetObjectLayerByName("GunLocs").Objects {
		if p, err := p.GetPoint(); err == nil {
			NewGun(p)
		}
	}
}

type Gun struct {
	pos       pixel.Vec
	sinceLast float64
	speed     float64
}

func NewGun(pos pixel.Vec) {
	g := Gun{
		pos:   pos,
		speed: 5,
	}

	guns = append(guns, &g)
}

func (g *Gun) Update(dt float64) {
	if g.sinceLast < g.speed {
		g.sinceLast += dt
	}

	if g.pos.To(player.CollisionBox().Center()).Len() > gunFireDist {
		return
	}

	if g.sinceLast > g.speed {
		g.Fire()
	}
}

func (g *Gun) Fire() {
	PlaySound(rocketSound)

	g.sinceLast = 0

	r := rocket{
		id:   rocketCount,
		pos:  g.pos,
		life: 5,
	}

	r.dir = r.toPlayer().Unit().Scaled(rocketAcc)

	rockets[rocketCount] = &r
	rocketCount++
}

func UpdateGuns(dt float64) {
	for _, g := range guns {
		g.Update(dt)
	}

	UpdateRockets(dt)
}

func DrawGuns(target pixel.Target) {
	DrawRockets(target)
}

type rocket struct {
	id   int
	pos  pixel.Vec
	dir  pixel.Vec
	life float64
}

func (r *rocket) Draw(target pixel.Target) {
	rocketSprite.Draw(target, pixel.IM.Moved(r.pos).Rotated(r.pos, r.dir.Angle()))
}

func (r *rocket) Update(dt float64) {
	r.life -= dt
	if r.life < 0 {
		r.explode()
		return
	}

	r.pos = r.pos.Add(r.dir)

	if r.toPlayer().Len() < explodeDist {
		r.explode()
	}

	var newDir pixel.Vec
	if len(flares) > 0 {
		var closestInd int
		closestDist := float64(gunFireDist * 2)
		for i, f := range flares {
			if r.pos.To(f.pos).Len() < closestDist {
				closestInd = i
			}
		}

		newDir = r.pos.To(flares[closestInd].pos)
	} else {
		newDir = r.toPlayer()
	}

	r.dir = r.dir.Add(newDir.Unit().Scaled(rocketAcc))
}

func (r *rocket) toPlayer() pixel.Vec {
	return r.pos.To(player.CollisionBox().Center())
}

func (r *rocket) explode() {
	PlaySound(explosionSound)

	delete(rockets, r.id)

	if r.toPlayer().Len() < minSafeDist {
		player.Hurt(50)
	} else if r.toPlayer().Len() < midSafeDist {
		player.Hurt(30)
	} else if r.toPlayer().Len() < maxSafeDist {
		player.Hurt(15)
	}

	NewFire(r.pos, pixel.ZV)
}

func UpdateRockets(dt float64) {
	for _, r := range rockets {
		r.Update(dt)
	}
}

func DrawRockets(target pixel.Target) {
	for _, r := range rockets {
		r.Draw(target)
	}
}
